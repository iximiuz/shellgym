# Design

Shell Gym is a background daemon that turns any Linux box into an
interactive command-line trainer. The student works in a completely
ordinary shell; a split screen (browser tab, or a future TUI) shows small,
fast-changing assignments - **reps**, in the traditional gym sense - and
reacts live as the student completes them.

This page covers architecture and internals. For the observation machinery
in depth see [detection.md](detection.md); for the check commands see
[checks.md](checks.md); for the content format see
[authoring-guide.md](authoring-guide.md).

## Goals

- Train fundamental Linux operations: navigation, files, permissions,
  redirection, pipes, processes, signals, networking basics.
- Fast scene changes; most reps take seconds to a minute.
- Zero interference with the student's shell: no prompt hooks, no
  wrappers, no PROMPT_COMMAND. All observation happens from the outside
  (procfs, the kernel proc connector, plain checks).
- Content-agnostic: the tool ships a format, not a curriculum. The
  `sample-linux-101` path is the reference/test path.
- Persistent progress: stop any time, resume days later.
- Repetition and review are intentionally external: the unit of
  repetition is the whole learning path - re-running it re-randomizes the
  parametric reps. There is no in-daemon skill tracking, review flagging,
  or scheduling.

## Subsystems

Three independent subsystems, glued by a small event bus:

1. **Content engine** (`internal/content`) - loads a learning path from
   disk, defines the unit format (markdown + YAML frontmatter), resolves
   parameters (vars), renders markdown with custom components, strips
   solve scripts in live mode. Knows nothing about task execution.
2. **Validation engine** (`internal/engine`) - executes init and check
   scripts, tracks per-task status with full run records (stdout, stderr,
   exit code, duration), provides built-in checks to task scripts over a
   unix socket, receives `hint_exit` messages. Knows nothing about
   markdown or rendering.
3. **UI** (`ui/`) - the `ui.UI` interface is the pluggable frontend
   contract; `ui/webui` is the built-in web UI (vanilla JS, `<template>`
   elements, no build step, embedded via `go:embed`); `ui/tui` is reserved
   for a future terminal UI.

`internal/state` persists progress; `internal/bus` is a tiny pub/sub the
engine publishes events through and any number of UIs consume.

## CLI

The daemon CLI is a typical [spf13/cobra](https://github.com/spf13/cobra) app:

- `shellgym serve --path <dir> [--addr :63636] [--state <dir>]
  [--run <dir>] [--user <login>] [--live]` - the daemon
- `shellgym solve --path <dir> [--api URL] [--unit <id>]` - acceptance
  driver: spawns a real interactive bash on a pty (indistinguishable from
  a student terminal), walks the path via the daemon's API, types each
  task's `solve:` lines, reports PASS/FAIL per unit
- `shellgym validate --path <dir>` - lint + render without running
- `shellgym skills [name]` - dump embedded authoring skills (markdown)
- `shellgym check <name> [args...]` - hidden; the target of the PATH shims
  task scripts call (`wait_cwd`, `hint_exit`, ...)

*`--live` is the student-facing mode: `solve:` blocks are stripped from the
on-disk unit files at startup and the debug API/UI is disabled.

## Content format

A learning path is a folder:

```
paths/sample-linux-101/
  path.yaml                 # id, title, description, shellUser
  010.orientation/          # module (numeric prefix = order, stripped from id)
    module.md               # optional module intro scene
    010.where-am-i/         # unit
      unit.md               # frontmatter + markdown body
      img.png               # unit-local static assets, referenced relatively
  020.moving-around/
    ...
```

Units and modules are identified by their prefix-less folder names
(`orientation/where-am-i`); prefixes only encode order.

### unit.md frontmatter

```yaml
title: Change into a directory
labels: [ubuntu, debian]   # distro filter (os-release ID/ID_LIKE); empty = any
requires: [systemd]        # host capability filter; unmet -> unit dropped
needs: [make-a-home]       # units (same module) whose *state* this unit builds on
vars:
  DIRNAME: { pick: [alpha, bravo, charlie] }   # random choice, sticky per attempt
  TOKEN:   { shell: "head -c3 /dev/urandom | od -An -tx1 | tr -d ' \n'" }
  PORT:    { value: "8080" }                   # fixed value
  PROJECT: { from: make-a-home.PROJECT }       # inherited from a preceding unit
init:                       # ordered root scripts, run once when the unit activates
  - name: create_tree
    run: |
      mkdir -p /tmp/gym/$DIRNAME
tasks:
  chdir:
    mode: edge              # edge (default) | level
    needs: []               # other tasks in this unit
    check: |
      wait_cwd "/tmp/gym/$DIRNAME"
    hint: |                 # optional dynamic-hint script, stdout -> UI
      echo "Your shell is still in $(shell_cwd)."
    solve: |                # hidden reference solution (typed by `shellgym solve`)
      cd /tmp/gym/$DIRNAME
```

- **Vars** resolve once per activation (persisted, stable across daemon
  restarts), are exported into every task script, and interpolate into the
  markdown as `${DIRNAME}`. `from:` references let dependent units share
  randomized state.
- **Filtering**: `labels` matches the running distro; `requires` matches
  detected host capabilities (currently `systemd`). Both filters apply at
  load time.
- **Init scripts** run in order, as root, once per activation; they create
  files as root and `chown` to the student. Failures surface to the UI
  (author-facing) and block the unit's tasks; a later activation retries.
- **Task modes**:
  - `edge` (default): the check runs (typically blocking on a `wait_*`
    built-in but any shell commands can be used) until it exits 0 once;
    then the task is completed forever. Failed attempts are recorded and
    restarted after a short delay.
  - `level`: the check is polled (~1s); exit 0 = satisfied, non-zero =
    unsatisfied, may flip both ways until the unit completes.
  - Edge tasks may not depend on level tasks (load-time error). Tasks are
    ordered dependency-first (topologically) - the order the UI shows and
    `solve` executes.
- A unit completes when all edge tasks are completed and all level tasks
  are simultaneously satisfied. Completion is terminal.
- **Dynamic hints**, two mechanisms that can be mixed freely:
  - a `hint:` script runs after failing attempts; its
    stdout replaces the task's hint area live. It sees the unit env plus
    `GYM_TASK_EXIT/STDOUT/STDERR` of the failed check run.
  - `hint_exit [task] <message>` inside a `check:` pushes a message to the
    UI immediately (over the daemon's unix socket) and TERMINATES the
    check script with the distinct code 42 (the runner injects a shell
    function wrapping the binary). The task name defaults to `$GYM_TASK`.

### Markdown body

Standard markdown (goldmark, GFM) plus block components in two equivalent
forms - inline attrs or an MDC-style YAML block:

```
::task{name="chdir"}
#active
Waiting for your shell to arrive in `/tmp/gym/${DIRNAME}`...
#completed
There you are.
::

::hint
---
title: Forgot the command?
---
Directories are entered with a two-letter command from the 70s.
::

::image{src="tree.png" alt="The tree"}
::
```

- `::task` renders live task status (spinner/checkmark, `#active` and
  `#completed` sections, hint area patched over WebSocket). If the unit
  has exactly one task, `name` may be omitted.
- `::hint` renders folded by default; `::image` and plain `![...]()` serve
  unit-local files.
- `${VAR}` interpolation is applied server-side; task status/hints are
  patched live by the UI.

## Validation engine details

- **Student shell discovery**: processes owned by the configured user
  whose exe is a known shell and which have a controlling TTY. Checks
  operate on the most recently started one (or all).
- **Exec watching**: the daemon subscribes to the kernel proc connector
  (netlink, needs `CAP_NET_ADMIN`) and records every exec (pid, ppid, uid,
  tty, argv) into a bounded ring buffer. Only tty-attached processes of
  the observed user count as student activity; the daemon's own task
  scripts run under `setsid` (no controlling tty) so they can never match
  themselves.
- **Built-in checks** are PATH shims that exec `shellgym check <name>`,
  which talks to the daemon over `<run>/gym.sock` where needed. All
  `wait_*` checks block until met and accept `--timeout <sec>` or `--now`
  (instant, for level tasks). Full reference: [checks.md](checks.md).
- Task scripts run as root via `bash -o pipefail`, in their own session,
  with per-attempt timeouts and whole-tree kill; stdout/stderr/exit are
  captured without blocking on lingering background children.

The mechanisms (procfs scanning, the proc connector, direct state
polling, the check API socket) are described in depth in
[detection.md](detection.md).

## State

A directory per path under `--state` (default `/var/lib/shellgym`):

```
<state>/<path-id>/
  progress.json            # unit/task statuses, vars, timestamps (small,
                           # atomic tmp+rename writes)
  runs/<module>__<unit>/<task>.jsonl   # last N run records per task
```

Run records (potentially large stdout/stderr, truncated to 4 KiB per
stream) append to per-task JSONL files compacted to the last 20 entries -
they never bloat the progress document. The debug API reads them back.

## Web UI

- Vanilla ES modules + CSS, no framework, no build step; all markup lives
  in `<template>` elements in index.html, JS only clones and fills them.
- Scenes = module intros and units, loaded lazily one at a time and slid
  horizontally to suggest movement along a path. On unit completion: a
  checkmark animation, ~1.4s pause, auto-advance (toggleable).
- Debug drawer (author-facing, hidden in `--live`): per-task last runs
  with exit codes, streams, timings - fed by the run JSONL files.
- Live updates over `WS /api/events` (task status, hints, unit
  completion, init failures, run records).

## HTTP API (consumed by the web UI and `shellgym solve`)

- `GET  /api/path` - path tree + per-unit status + progress counters
- `GET  /api/unit/{id}` - rendered unit (HTML, tasks, vars)
- `POST /api/activate/{id}` - make unit current (resolve vars, run init).
  Returns 409 for a *locked* unit - one whose `needs:` dependencies are
  not all completed; there is no way around the gate. Browsing never
  activates: the web UI auto-activates only the next unit in path order
  and otherwise offers an explicit start button.
- `POST /api/reset/{id}` - forget unit progress, re-activate
- `POST /api/module-seen/{id}` - mark a module intro viewed
- `GET  /api/module/{id}` - rendered module intro
- `GET  /api/debug/{id}` - task run history (404 in live mode)
- `GET  /api/status` - distro, exec source, live flag, observed shells
- `GET  /api/events` - WebSocket event stream
- `GET  /unit-assets/{id}/...` - unit-local static files

## Deployment

The daemon runs on a plain Linux host - e.g., a VM or an [iximiuz Labs Linux playground](https://labs.iximiuz.com/playgrounds?category=linux&filter=official) -
as root, typically under systemd (`systemd-run --unit=shellgym --collect
...`). The student uses any terminal on the box; the UI is on port 63636.
