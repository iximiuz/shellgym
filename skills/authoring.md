---
name: shellgym-authoring
description: Author Shell Gym learning paths - path/module/unit format, verification tasks (check/hint/solve), built-in wait_* checks, vars, components, guidelines, and the test workflow. Use when creating or editing Shell Gym content.
---

# Authoring Shell Gym Learning Paths

Shell Gym is a daemon that trains Linux command-line skills through
repetition. The student works in an ordinary terminal while small
assignments ("reps") appear in a split-screen web UI and complete
automatically as the system state changes. Nothing is injected into the
student's shell: the daemon observes from the outside (procfs for
shells, cwds, files, processes, and ports; the kernel proc connector
for exec events), so the practiced skills transfer one-to-one to any
real terminal.

You author **learning paths**. A path is the course: a directory tree
of **modules** (themed groups), each holding **units**. A unit is one
rep - a markdown page plus scripts: `init:` builds the scene, `check:`
scripts recognize the accomplishment, `hint:` scripts and `hint_exit`
coach on failure, and a hidden `solve:` script proves the rep is
solvable. A unit completes when all of its **tasks** are met. One
daemon serves one path.

A unit is armed by **activation**: init runs and checks start watching
only then. A student can have several units in progress, but the daemon
supervises only the most recently activated one - viewing a started
unit re-activates it and moves the watch there. Note that `wait_exec`
only counts commands run after the unit's latest activation, while
state-based checks (`wait_file`, `wait_cwd`, ...) pass on whatever is
true when they look. The UI auto-activates just the next unit in path
order; a student who jumps ahead must start the unit explicitly, and a
unit whose `needs:` dependencies are not all completed is **locked** -
it cannot be activated at all until they are.
Consequences for authoring: students may solve units out of order, so
never assume an earlier unit was solved unless you declare it in
`needs:`; and browsing a unit runs nothing, so scenes may safely do
disruptive setup in init.

## Layout

```
my-path/
  path.yaml                  # id, title, description, shellUser
  010.first-module/          # numeric prefix = order, stripped from ids
    module.md                # optional module intro (first # heading = title)
    010.some-unit/
      unit.md                # the unit: frontmatter + markdown body
      diagram.png            # unit-local assets, referenced relatively
    020.other-unit/
      unit.md
  020.second-module/
    ...
```

Ids are prefix-less folder names: `first-module/some-unit`. The numeric
prefixes only encode order and never appear in ids or URLs.

`path.yaml`:

```yaml
id: my-path
title: My Learning Path
description: >
  One paragraph shown as the path's summary.
shellUser: laborant     # the login user whose shells are observed
```

## unit.md frontmatter

```yaml
---
title: Change into a directory        # required
labels: [ubuntu, debian]              # optional distro filter (ID/ID_LIKE)
requires: [systemd]                   # optional host capability filter
needs: [earlier-unit]                 # optional same-module state deps
vars:
  DIRNAME: { pick: [alpha, bravo] }   # random choice, sticky per attempt
  TOKEN:   { shell: "head -c4 /dev/urandom | od -An -tx1 | tr -d ' \\n'" }
  PORT:    { value: "8080" }          # fixed value
  OTHER:    { from: earlier-unit.OTHER }  # inherit from a preceding unit
init:                                 # ordered root scripts, run once on activation
  - name: create_tree
    run: |
      mkdir -p /tmp/gym/$DIRNAME
tasks:
  chdir:                              # task name = map key
    mode: edge                        # edge (default) | level
    needs: []                         # other tasks in this unit (UI shows the task locked until they pass)
    timeout: 45                       # per-attempt seconds (default: 30 edge / 10 level)
    check: |                          # exit 0 = condition met
      wait_cwd "/tmp/gym/$DIRNAME"
    hint: |                           # optional dynamic-hint script
      echo "Your shell is still in $(shell_cwd)."
    solve: |                          # hidden reference solution (see below)
      cd /tmp/gym/$DIRNAME
---
```

Rules and behaviors:

- **Vars** make reps parametric so repetition stays honest. They
  resolve once per activation, persist for the whole attempt (across
  daemon restarts), are exported into every script as environment
  variables, and interpolate into the markdown as `${DIRNAME}`.
- **Filters** apply at load time; filtered-out units do not exist for
  that host. `labels` matches `ID`/`ID_LIKE` from `/etc/os-release`
  (`ubuntu`, `debian`, `rocky`, ...); `requires` matches detected
  capabilities (currently: `systemd`).
- **edge tasks** ("the student did X") run until they first exit 0,
  then stay completed forever; use blocking `wait_*` checks.
  **level tasks** ("X is currently true") are re-polled about once a
  second (use `--now` on the `wait_*` calls) and may flip back. The
  unit completes when all edge tasks are completed AND all level tasks
  are simultaneously satisfied; completion is terminal. Edge tasks may
  not depend on level tasks.
- **needs (unit-level)** declares that this unit builds on the *state*
  left behind by earlier units; the listed units must be preceding
  units in the same module. Keep chains short (< 5). Use `vars.from`
  to share randomized values along the chain.
- **Init scripts** run in order, as root, once per activation, 60 s
  timeout each. If one fails the unit's tasks do not start and the
  next activation retries from scratch - keep them idempotent.
- **solve** is the hidden reference solution: plain shell lines typed
  one by one into a real pty shell by `shellgym solve` (acceptance
  testing). Write each line to be independently typable: no heredocs,
  no multi-line constructs, no `\` continuations. Blank lines and `#`
  comments are skipped. Every task needs one, or the unit cannot be
  acceptance-tested. In student-facing deployments (`serve --live`)
  solve blocks are stripped from the on-disk files.

## Built-in checks (on PATH inside check/hint scripts)

- `wait_cwd <path|regex>` - ANY of the student's open shells has the
  given working directory (read live from `/proc/<pid>/cwd`). The
  argument is an exact absolute path; if it contains regex
  metacharacters and compiles, it is matched against the whole path
  as a regex, auto-anchored `^(...)$`
- `wait_exec <regex>` - the student ran a command matching regex
  (matched against full argv; only tty-attached processes of the
  observed user, executed after the unit's activation, count; matched
  commands are buffered, so a command run just before the check
  restarted still passes; commands typed at human speed are captured
  reliably, but processes spawned in tight machine-speed loops can be
  missed - do not depend on catching those)
- `wait_env <NAME> [regex]` - a command was observed with the env var
  set; this is how exports are verified (ask the student to run any
  command after exporting)
- `wait_file <path|glob>` / `wait_file_gone <path|glob>`
- `wait_file_contains <path> <regex>` (multiline mode: `^...$` = a line)
- `wait_proc <regex>` / `wait_proc_gone <regex>` (full-cmdline match,
  all processes on the box)
- `wait_port <port>` / `wait_port_free <port>` (listening TCP, v4+v6)
- `shell_cwd` - prints the cwd of the most recently started student
  shell (for hint scripts; note: `wait_cwd` accepts a match in any
  shell, so the two may disagree when several terminals are open)
- `shells` - lists all observed shells, one per line (`PID exe tty
  cwd`, most recent first); a debugging aid
- `hint_exit [task] <message>` - pushes a hint to the UI immediately
  and TERMINATES the check script with exit code 42; the task defaults
  to the current one (`$GYM_TASK`)

All `wait_*` block until met; add `--timeout <sec>` for a bound or
`--now` for a single instant evaluation (the form for level tasks).
Regexes are Go RE2 (no backreferences). Checks compose freely with
shell:

```yaml
check: |
  wait_file --timeout 15 "$HOME_DIR/junk.tmp" || \
    hint_exit "junk.tmp never appeared - was the init OK?"
  wait_file_gone "$HOME_DIR/junk.tmp"
```

Choosing: verify **effects** (`wait_file`, `wait_port`, ...) over
commands; reserve `wait_exec` for commands that leave no trace (`ls`,
`cat`, `curl`). `wait_exec` proves the command was run, not that it
succeeded. Never require one exact command form when several are
correct - keep `wait_exec` regexes permissive.

All scripts (`init:`, `check:`, `hint:`) run as root under
`bash -o pipefail`, each in its own session (no controlling tty - the
guard that keeps checks from matching themselves).
Environment available to scripts: unit vars, `GYM_UNIT`, `GYM_TASK`,
`GYM_USER` (observed login user). `hint:` scripts additionally get
`GYM_TASK_EXIT`, `GYM_TASK_STDOUT`, `GYM_TASK_STDERR` from the last
failed check run - enough to diagnose why it failed and say something
specific. Hint script stdout replaces the task's hint area live
(rate-limited to one refresh per ~10 s).

## Markdown body

Standard markdown plus block components, with inline attributes or an
MDC-style YAML block:

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
Directories are entered with a two-letter command.
::

::image{src="diagram.png" alt="The tree"}
::
```

- `::task` renders the live task box: `#active` (or the unnamed body)
  shows while pending, `#completed` after success; dynamic hints are
  patched into it over WebSocket. If the unit has exactly ONE task, the
  `name` attribute may be omitted.
- `::hint` renders folded-by-default - for static nudges the student
  opts into.
- `::image` and plain `![...](file.png)` serve unit-local files.
- `${VAR}` interpolates in the body, titles, and component text.
- Nested components use identical `::` fences.

## Authoring guidelines

What makes a good path:

- **Stay focused**: train one coherent skill area per path; depth of
  practice beats breadth of coverage.
- **Plan for lots of repetition**: a single `cd` rep forms no skill;
  budget many reps per command across the path.
- **Make repetition naturally diverse**: not five bare `cd`s in a row,
  but small realistic scenarios where the command keeps coming up on
  its own (exploring a project tree, chasing a log file), each hitting
  the same muscle from a different angle.
- **Sometimes revisit earlier commands in later units**, woven
  naturally into the scenario (a `find` rep that ends with removing
  what was found) - not as a scheduled review session.

Rules for individual units:

- One skill per rep; most reps should complete in under a minute.
- Never put the exact solution in the problem statement; hints may
  point, not paste.
- **Baseline negative checks**: a task verifying something *disappears*
  must first confirm it existed (`wait_file --timeout 15 X || exit 1`
  then `wait_file_gone X`), or it auto-solves before the scene exists.
- **Pick scene locations by lifetime**: standalone units may build
  scenes in `/tmp` (conventionally `/tmp/gym`) or `/run`-style volatile
  locations; units whose state later units build on (or that should
  survive a reboot) should prefer the student's home, since `/tmp` may
  be wiped on reboot; system locations (`/etc`, ...) only when the task
  itself is about them.
- Init scripts run as root: create files as root and `chown` them to
  `$GYM_USER`. Student-owned *processes* need
  `systemd-run --uid=$GYM_USER` and a named sh wrapper script (argv0
  tricks like `exec -a` break on multi-call coreutils distros). Beware
  `pgrep -f`, `wait_proc`, and `wait_proc_gone` matching the script's
  own text - anchor patterns on argv text only the target has (bracket
  trick + an argument: `'my-nam[e] 86400'`).
- Prefer distro-neutral commands; label distro-specific units
  (`labels: [ubuntu, debian]`), ideally with a sibling unit per family.
- Random tokens: `head -cN /dev/urandom | od -An -tx1 | tr -d ' \n'`
  (openssl is not installed everywhere).
- Simple, clear International English; a friendly coach, not a quiz
  master.

## Test workflow

```sh
shellgym validate --path my-path          # lint + render every unit
sudo systemd-run --unit=shellgym --collect \
  shellgym serve --path "$PWD/my-path" --addr :63636 --user laborant
shellgym solve --path my-path             # acceptance: types every solve script
shellgym solve --path my-path --unit first-module/some-unit
```

`validate` catches format errors: missing titles, empty or misdeclared
tasks, bad task graphs, invalid unit deps and `from:` references, and
markdown that does not render. `solve` is the real test: it spawns an
interactive bash on a pty (indistinguishable from a student), activates
each unit through the API, types the solve lines, and waits for
completion, reporting PASS/FAIL per unit. A unit whose `needs:` are not
solved cannot be activated - even by `solve --unit`; solve its chain
first (or run without `--unit`, which walks the path in order).

Debugging a failing unit - every check attempt's exit code, stdout,
stderr, and duration is recorded:

- the UI's **debug drawer** (press `d`; hidden in live mode);
- `GET /api/debug/<unit-id>`;
- on disk under `<state>/<path-id>/runs/`.

`POST /api/reset/<unit-id>` forgets a unit's progress and re-runs its
init from scratch - the fast iteration loop while authoring.
