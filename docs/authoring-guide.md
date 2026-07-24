# Authoring Guide

This page is the human-facing reference for writing Shell Gym learning
paths. The same material, condensed for AI agents, is embedded in the
binary - `shellgym skills authoring` prints it (drop the output into
`.claude/skills/shellgym-authoring/SKILL.md` to teach an agent the
format).

## Mental model

A **learning path** is a directory tree of small exercises that one
Shell Gym daemon serves as a whole. The student opens the gym page next
to an ordinary terminal and works through the path exercise by
exercise. Shell Gym never instruments the student's shell - the daemon
observes the system from the outside (interactive shells and their
working directories, exec events, files, processes, ports) and
completes tasks the moment the expected state change is seen.

The exercise itself is called a **unit**, and a unit is one rep - a
markdown page the student reads, plus scripts the daemon runs:

- `init:` (optional) scripts build the scene (the files, processes, and ports the
  student will act on);
- `check:` (one or more required) scripts recognize the student's accomplishment;
- `hint:` scripts and `hint_exit` calls coach on failure;
- a hidden `solve:` script proves the rep is solvable and powers
  automated acceptance testing.

Most reps should take a student well under a minute. A path teaches not
through any single rep but through many small, varied reps that keep
the student typing real commands until the motions become automatic.

## Glossary

- **Path** - the top-level entity: `path.yaml` plus modules (and their units).
  One daemon serves one path.
- **Module** - a themed group of units with an optional intro scene
  (`module.md`).
- **Unit** (a **rep**) - one exercise: a `unit.md` file with YAML
  frontmatter (scripts and metadata) and a markdown body (the page).
- **Task** - one verifiable condition inside a unit. A unit completes
  when all of its tasks are met.
- **Check** - a task's shell script that exits 0 when the condition is
  met; usually built around a blocking `wait_*` built-in.
- **Scene** - the state of the system for the student to act on.
- **Vars** - per-unit parameters (fixed, random, or inherited) that
  make a rep look different on every attempt.
- **Activation** - the moment a unit becomes the student's current
  exercise; init scripts run and vars resolve once per activation.
- **Observed user** - the login user declared as `shellUser` in
  `path.yaml`; only this user's shells and commands count as student
  activity (exported to scripts as `$GYM_USER`, with the home directory
  as `$GYM_USER_HOME`).
- **Solve script** - the hidden reference solution, typed into a real
  pty shell by `shellgym solve` during acceptance testing.

## Path layout

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

The complete field reference:

```yaml
---
title: Change into a directory        # required
labels: [ubuntu, debian]              # optional distro filter
requires: [systemd]                   # optional host capability filter
needs: [earlier-unit]                 # optional same-module state deps
vars:
  DIRNAME: { pick: [alpha, bravo] }   # random choice, sticky per attempt
  TOKEN:   { shell: "head -c4 /dev/urandom | od -An -tx1 | tr -d ' \\n'" }
  PORT:    { value: "8080" }          # fixed value
  SAME:    { from: earlier-unit.DIRNAME }  # inherit from a preceding unit
init:                                 # ordered root scripts, run once on activation
  - name: create_tree
    run: |
      mkdir -p /tmp/gym/$DIRNAME
tasks:
  chdir:                              # task name = map key
    mode: edge                        # edge (default) | level
    needs: []                         # other tasks in this unit
    timeout: 45                       # per-attempt seconds (default: 30 edge / 10 level)
    check: |                          # exit 0 = condition met
      wait_cwd "/tmp/gym/$DIRNAME"
    hint: |                          # optional dynamic-hint script
      echo "Your shell is still in $(shell_cwd)."
    solve: |                          # hidden reference solution
      cd /tmp/gym/$DIRNAME
---
```

### Filtering: `labels` and `requires`

Both filters apply at load time; filtered-out units simply do not exist
for that host.

- `labels` matches the distro: the `ID` and `ID_LIKE` values from
  `/etc/os-release` (`ubuntu`, `debian`, `rocky`, `rhel`, ...). Empty
  means "any distro". Use it for package-manager reps and other
  distro-specific material, ideally providing a sibling unit per family
  (see `sample-linux-101/070.package-tools` for the pattern).
- `requires` matches detected host capabilities. Currently detected:
  `systemd` (a reachable system systemd instance).

### Unit dependencies: `needs`

`needs` declares that this unit builds on the *state* left behind by
earlier units - the listed units must be **preceding units in the same
module**. Keep chains short (under 5); long chains make it painful to
reset or jump around. Use `vars: { X: { from: other-unit.X } }` to share
randomized values along the chain, so the follow-up unit talks about the
same directory or token the student saw before.

### Vars

Vars make reps parametric so repetition stays honest:

- `value` - a fixed string;
- `pick` - a uniformly random choice from a list;
- `shell` - the stdout of a shell command (trimmed), e.g. a random
  token: `head -c4 /dev/urandom | od -An -tx1 | tr -d ' \n'` (prefer
  this over openssl, which is not installed everywhere);
- `from` - inherit a var resolved by a preceding unit.

Vars resolve once, when the unit is first rendered or activated, and
persist for the whole attempt (across daemon restarts too). They are
exported into every `init:`/`check:`/`hint:` script as environment
variables and interpolate into the markdown body, titles, and component
text as `${DIRNAME}`. Scripts additionally see the vars of every unit
listed in the unit's `needs:`.

**Task vars.** A check can also publish a var at runtime with the
`set_var <NAME> <value>` built-in - for values that only exist once the
student acts, like the PID of the shell that completed a step (see
`wait_cwd` in [checks.md](checks.md)). Task vars join the unit's vars:
later runs of the unit's own scripts and the scripts of dependent
(`needs:`) units see them in the environment. They are the way to pass
values between tasks and units - never stash such values in files; use
a file only when the data is BLOB-like (content rather than a value).
Because their value does not exist at render time, task vars cannot be
interpolated into markdown or referenced with `from:`.

### Init scripts

Init scripts run in order, as root, once per activation, with a 60 s
timeout each. If one fails, the unit's tasks do not start (they would
misfire on a half-built scene), the failure surfaces in the author UI,
and the next activation retries from scratch - so keep init scripts
idempotent.

Conventions that save debugging time:

- create files as root, then `chown` them to `$GYM_USER`; scenes in the
  student's home go under `$GYM_USER_HOME` (the observed user's home
  directory, resolved by the engine);
- student-owned *processes* need `systemd-run --uid=$GYM_USER` and a
  named sh wrapper script; argv0 tricks like `exec -a` break on
  multi-call coreutils distros;
- beware `pgrep -f` (and `wait_proc`) matching the init script's own
  text - anchor patterns on argv text only the target has (bracket
  trick plus an argument: `'my-nam[e] 86400'`).

### Tasks

Task names are the map keys; the UI shows tasks in dependency-first
(topological) order, which is also the order `solve` types them.

- **`mode: edge`** (default) - "the student did X". The check runs
  (typically blocking on a `wait_*` built-in) until it first exits 0;
  then the task is completed forever. Failed attempts are recorded and
  the check restarts after a short delay.
- **`mode: level`** - "X is currently true". The check is re-polled
  about once a second (use `--now` on the `wait_*` calls) and may flip
  between satisfied and unsatisfied.
- `needs` - other tasks in the same unit that must be completed
  (edge) or currently satisfied (level) before this one starts. Edge
  tasks may not depend on level tasks (load-time error). The UI shows
  such a task with a lock icon and an "unlocks after" note until its
  dependencies pass.
- `timeout` - per-attempt seconds, overriding the defaults (30 for
  edge attempts, 10 for level polls).

A unit completes when all edge tasks are completed AND all level tasks
are simultaneously satisfied. Completion is terminal - checks stop and
the statuses freeze.

**Baseline negative checks.** A task verifying that something
*disappears* (file removed, process killed, port freed) must first
confirm it existed, or it auto-solves before the scene is even built:

```yaml
check: |
  wait_file --timeout 15 "$GYM_USER_HOME/junk.tmp" || exit 1
  wait_file_gone "$GYM_USER_HOME/junk.tmp"
```

### Hints

Two mechanisms, freely mixable; both update the task box live:

- a **`hint:` script** runs after failing check attempts (rate-limited
  to one refresh per ~10 s); its stdout replaces the task's hint area.
  It sees the unit env plus `GYM_TASK_EXIT`, `GYM_TASK_STDOUT`,
  `GYM_TASK_STDERR` from the failed run - enough to diagnose *why* the
  check failed and say something specific.
- **`hint_exit [task] <message>`** inside a `check:` pushes a message
  immediately and terminates the check with exit code 42:

  ```yaml
  check: |
    wait_port --timeout 60 "$PORT" || \
      hint_exit "Nothing is listening on $PORT yet. Is the server running?"
  ```

Hints point at what is wrong or where to look. Never paste the solution
command.

### Solve scripts

`solve:` is the hidden reference solution, used only by the
`shellgym solve` acceptance driver, which types it line by line into a
real pty shell. Write each line to be independently typable: no
heredocs, no multi-line constructs, no `\` continuations. Blank lines
and `#` comments are skipped. In student-facing deployments
(`serve --live`) solve blocks are stripped from the on-disk files.

Every task should have a solve script - a unit without one cannot be
acceptance-tested and fails the `solve` run.

## Markdown body

Standard markdown (GFM) plus block components, written either with
inline attributes or an MDC-style YAML block:

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
  patched into it over WebSocket. If the unit has exactly **one** task,
  the `name` attribute may be omitted.
- `::hint` renders folded by default - for static nudges the student
  opts into.
- `::image` and plain `![...](file.png)` serve unit-local files.
- `${VAR}` interpolates everywhere: body, titles, component text.
- Nested components use identical `::` fences.

Module intros (`module.md`) are plain markdown scenes; the first `#`
heading becomes the module title.

## Authoring guidelines

### What makes a good path

- **Stay focused.** A good path trains one coherent skill area (shell
  navigation, file management, process control) rather than touring
  everything at once. Depth of practice beats breadth of coverage.
- **Plan for lots of repetition.** A skill is not formed by doing
  something once - a path that asks for a single `cd` teaches nothing
  durable. Budget many reps per command, spread across the path.
- **Make the repetition naturally diverse.** Do not simply ask for five
  different `cd`s in a row - that reads as a chore, not training.
  Instead, build small realistic scenarios in which the command keeps
  coming up on its own: exploring a project tree, chasing a log file,
  cleaning up after a misbehaving script. Each scenario exercises the
  same muscle from a slightly different angle (relative vs. absolute
  paths, going back up, jumping home).
- **Sometimes revisit earlier commands in later units.** Follow-up
  modules should occasionally weave commands from preceding units into
  their scenarios - e.g., a `find` rep that ends with removing what was
  found. Keep it natural: the revisit should serve the scenario, not
  read like a scheduled review session.

### Rules for individual units

- **One skill per rep.** Most reps should complete in under a minute;
  if a unit needs three tasks and a page of prose, split it.
- **Verify effects, not keystrokes,** where an effect exists; reserve
  `wait_exec` for commands that leave no trace. Never require one exact
  command form when several are correct.
- **Never put the exact solution in the problem statement.** Hints may
  point, not paste.
- **Prefer distro-neutral commands**; label distro-specific units.
- **Pick scene locations by the unit's lifetime.** Scenes should live
  where a reset is cheap and collisions with the real system are
  impossible:
  - **Standalone units** (no later unit `needs` their state) can build
    scenes in `/tmp` (conventionally `/tmp/gym`) or similar volatile
    locations like `/run`.
  - **Durable units** - ones whose files later units build on, or whose
    state should survive a host reboot in a resumable path - should
    prefer the student's home directory: `/tmp` and `/run` may be wiped
    on reboot, taking half-done progress with them.
  - **System locations** (`/etc`, `/usr/local/bin`, ...) are fine only
    when the task itself is about them (e.g., a rep on editing a config
    under `/etc`); never park auxiliary scene files there.
- Write in simple, clear International English; the tone is a friendly
  coach, not a quiz master.

## Testing a path

```sh
shellgym validate --path my-path        # lint + render every unit
sudo systemd-run --unit=shellgym --collect \
  shellgym serve --path "$PWD/my-path" --addr :63636 --user laborant
shellgym solve --path my-path           # acceptance: types every solve script
shellgym solve --path my-path --unit first-module/some-unit
```

`validate` catches format errors: missing titles, empty or misdeclared
tasks, bad task graphs (unknown `needs`, edge depending on level,
cycles), invalid unit deps and `from:` references, and markdown that
does not render.

`solve` is the real test: it spawns an interactive bash on a pty
(indistinguishable from a student to the daemon), activates each unit
through the API, types the solve lines, and waits for completion. It
reports PASS/FAIL per unit.

When a unit fails, look at the recorded check runs - every attempt's
exit code, stdout, stderr, and duration is kept:

- in the UI's **debug drawer** (press `d`; hidden in live mode);
- via `GET /api/debug/<unit-id>`;
- on disk under `<state>/<path-id>/runs/`.

Iterate with `POST /api/reset/<unit-id>` to forget a unit's progress
and re-run its init from scratch.

Shell Gym content for iximiuz Labs runs on playgrounds - see the
[development notes in the README](../README.md#development-workflow) for
the sync-build-test loop and the e2e playground setup.
