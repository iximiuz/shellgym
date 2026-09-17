---
name: shellgym-authoring
description: Author ShellGym learning paths - path/module/unit format, verification tasks (check/hint/solve), built-in wait_* checks, vars, components, guidelines, and the test workflow. Use when creating or editing ShellGym content.
---

# Authoring ShellGym Learning Paths

ShellGym is a daemon that trains Linux command-line skills through
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
(and `wait_env`, `wait_line`) only counts commands run after the unit's
latest activation, while state-based checks (`wait_file`, `wait_cwd`,
...) pass on whatever is true when they look. That horizon is shared by
ALL tasks of the unit: a task gated with `needs:` still judges commands
buffered while it was locked. "Run date, then hostname" as two gated
`wait_exec` tasks completes both the instant `date` passes for a student
who ran `hostname` first; "cd to X, then ls -l" is completed by an
`ls -l` typed anywhere before the cd. Task-level `needs:` orders the
task boxes, never the commands - a gated check must demand something the
earlier steps cannot have produced: chain by sequence number (step one
`wait_exec ... || exit 1` then `set_var STEP1_SEQ "$(event_seq)"`, step
two `wait_exec --after "$STEP1_SEQ" ...`; exec and line events share one
clock, so a `wait_line --after` a mark taken after an exec works too), verify an
effect, or add a match condition the stray command fails (`--cwd` for
location reps, an argument the first step created); or the steps
become separate units (a unit's activation is a fresh horizon). The UI auto-activates just the next unit in path
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
variant: scene=forest                 # optional: shown only when "forest" is drawn for "scene"
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
  Do NOT use vars in `title:` - the path map lists all units including
  never-activated ones, whose vars have no values yet, so the map shows
  the raw `${NAME}` text. Keep titles literal; use vars in the body.
- **Filters** apply at load time; filtered-out units do not exist for
  that host. `labels` matches `ID`/`ID_LIKE` from `/etc/os-release`
  (`ubuntu`, `debian`, `rocky`, ...); `requires` matches detected
  capabilities (currently: `systemd`, `python3`, `readline`). A unit with
  an unmet requirement is shown but never activated.
- **Variants** swap whole units where vars only vary details:
  `variant: <key>=<value>` tags a unit (key/value: lowercase letters,
  digits, dashes). When the path is first served, one value per key is
  drawn uniformly from the values used across the path; units tagged
  with the other values are hidden for that attempt. The draw is
  persisted with the progress (survives restarts, fresh per attempt),
  units without a variant are always shown, and keys draw
  independently. A key can carry a storyline across modules: every
  `scene=forest` unit appears in the same attempts. Rules: a unit may
  `needs:`/`from:` only units that are always shown with it (no
  variant, or the same key=value) - load-time error otherwise; a key
  with one value is always drawn (validate warns).
- **edge tasks** ("the student did X") run until they first exit 0,
  then stay completed forever; use blocking `wait_*` checks.
  **level tasks** ("X is currently true") are re-polled about once a
  second (use `--now` on the `wait_*` calls) and may flip back. The
  unit completes when all edge tasks are completed AND all level tasks
  are simultaneously satisfied; completion is terminal. Edge tasks may
  not depend on level tasks.
- **attempts**: a check run that exits non-zero on its own (`hint_exit`,
  or a plain `exit 1` on the wrong branch) is a rejected attempt; a run
  killed by the task timeout is not. After a rejected attempt the task's
  exec/line horizon moves to now - the restarted check judges only newer
  commands, so a wrong answer is judged once and its hint stays until
  the next try - and `GYM_CHECK_ATTEMPT` (1 on the first run) goes up by
  one in the task's `check:`/`hint:` scripts; a rejection without
  `hint_exit` runs `hint:` at once. Escalate with it:
  `[ "$GYM_CHECK_ATTEMPT" -gt 2 ] && hint_exit "<the answer>"` before
  the usual nudge. An expired `wait_* --timeout N || exit 1` is the
  check's own exit and counts - for an idle nudge use a plain blocking
  wait and `hint:`.
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
  Every solve line is typed exactly as written, except that references
  to the unit's vars (`$NAME`, `${NAME}`) are replaced by their values
  first - a student types the literal value, and `wait_line` checks see
  the line as typed. Other `$` syntax is left to bash (the vars are also
  exported into the solve shell). Lines starting with `#!` are
  **directives** for interactions a plain typed line cannot express -
  signal keys, pager keystrokes, and commands that hold the foreground
  (a normal solve line waits for the shell's next prompt, so it would
  block forever behind `sleep 300` or an open pager):
  - `#!type TEXT` - write TEXT to the pty verbatim, no Enter (unit vars
    are substituted as in a normal line);
  - `#!keys K K...` - send named keys: `enter`, `tab`, `space`, `esc`,
    `C-c`/`C-z`/any `C-<letter>`, or a single literal character
    (`q`, `/`);
  - `#!wait SECONDS` - pause to let the previous keystrokes take
    effect (start a process, redraw a pager) before the next line.
  Example - start a foreground sleep and interrupt it:
  `#!type sleep 300` / `#!keys enter` / `#!wait 1` / `#!keys C-c`

## Built-in checks (on PATH inside check/hint scripts)

- `wait_cwd [shell-pid] <path|regex>` - a student shell has the given
  working directory (read live from `/proc/<pid>/cwd`): without a PID
  any open shell may match, with a PID only that specific shell
  counts. The path is an exact absolute path; if it contains regex
  metacharacters and compiles, it is matched against the whole path
  as a regex, auto-anchored `^(...)$`. On success prints the matched
  shell's PID - publish it as a task var when a later task or a
  dependent unit must watch the SAME shell:
  `TRAVELER=$(wait_cwd "/tmp/gym/$D") || exit 1` then
  `set_var TRAVELER "$TRAVELER"` in one unit, and
  `wait_cwd "$TRAVELER" "$GYM_USER_HOME"` in the unit that `needs:` it
- `wait_exec [--argc N] [--latest] [--cwd <path>] [--after N] <regex>` - the student ran a command
  matching regex (matched against full argv joined with spaces; only
  tty-attached processes of the observed user, executed after the
  unit's activation, count; matched commands are buffered, so a command
  run just before the check restarted still passes; commands typed at
  human speed are captured reliably, but processes spawned in tight
  machine-speed loops can be missed - do not depend on catching those).
  On success prints the matched command's argv, so a check can branch
  on WHICH command satisfied the pattern. `--argc N` also requires
  exactly N argv elements - the way to tell a quoted space-containing
  argument from the same text as separate arguments (identical when
  joined). `--latest` prefers the newest buffered match over the
  oldest - use it for right/wrong-branch checks
  (`REPORT=$(wait_exec --latest '(^|/)(right|wrong)$')` then `case` +
  `hint_exit` on the wrong branch), so that when several answers are
  buffered the newest one is judged. `--cwd <path>` also requires the
  command to have run FROM that directory (exact path, or a regex
  matched against the whole path when it contains metacharacters - the
  `wait_cwd` rules; an unreadable cwd never matches) - for reps where
  the location is the point: a listing from inside a directory, a
  relative path, a file created "right here". `--after N` accepts only
  events observed after the mark N taken with `event_seq` - this chains
  ordered steps across tasks (see the activation note above). IMPORTANT: shells exec only
  EXTERNAL commands - builtins (`echo`, `printf`, `true`, `false`,
  `pwd`, `type`, `cd`, ...) produce no exec event and are invisible to
  `wait_exec`; anchor such reps on an external command (`whoami`,
  `date`, `seq`, `/bin/echo`, ...) or on an effect - or observe the
  typed line itself with `wait_line`. Exec events also carry no trace of
  the shell line: `a; b` and `a && b` are identical to `wait_exec`
- `wait_line [--latest] [--after N] <regex>` - the student typed a command line
  matching regex (the line as bash's readline returned it, surrounding
  whitespace trimmed, otherwise verbatim: operators, quotes, pipes, and
  builtins included). Same scoping and buffering as `wait_exec`; prints
  the matched line so a check can branch on it (`--latest` for
  right/wrong branching, as above; `--after` as for `wait_exec`, on
  the same clock). There is no `--cwd` for lines: when
  the location matters, pair the line check with a `wait_exec --cwd` on
  the command the line runs (external commands only - a builtin line
  such as `echo *` cannot be placed). Keep regexes permissive about
  whitespace (`sleep 2&&hostname` is the same command) and pair with
  `wait_exec` when the command must also have run. OPTIONAL CAPABILITY:
  every unit using it MUST declare `requires: [readline]` - hosts without
  the readline uprobe (no tracefs, non-bash login shell) mark such units
  unsupported instead of running them; without the declaration the check
  fails at once with exit code 2 there. Only bash is observed
- `wait_env [--cwd <path>] [--after N] <NAME> [regex]` - a command was observed with the env var
  set; this is how exports are verified (ask the student to run any
  command after exporting). `--cwd` and `--after` as for `wait_exec`
- `event_seq` - prints the current event sequence number (one clock for
  exec and line events): the mark for `--after`. Take it right after the
  wait that saw step one's event and publish it with `set_var`
- `wait_file <path|glob>` / `wait_file_gone <path|glob>`
- `wait_dir <path|glob>` - like `wait_file`, but only a directory
  satisfies it (use for `mkdir` tasks so a plain file at the path does
  not pass)
- `wait_file_contains <path> <regex>` (multiline mode: `^...$` = a line)
- `wait_file_mode <path> <octal>` - the file's permission bits equal the
  octal mode; 4-digit form compares setuid/setgid/sticky too (`4755`),
  3-digit form only matches when all special bits are clear
- `wait_file_newer <path> <reference-path>` - the file's mtime is
  strictly newer than the reference file's (verifying `touch` on an
  existing file, "newer than" scenes); init typically plants a
  reference marker next to an artificially aged target
  (`touch -d '2 days ago' target`)
- `wait_proc <regex>` / `wait_proc_gone <regex>` (full-cmdline match,
  all processes on the box)
- `wait_proc_state <regex> <state-letters>` - a process matching regex
  is in one of the given `/proc/<pid>/stat` states (`T` = stopped by a
  job-control signal - the Ctrl-Z check; `SRD` = alive/running - the
  "resumed" check). Prints the matched state letter on success. Same
  full-cmdline matching (and self-match caveat) as `wait_proc`
- `wait_port <port>` / `wait_port_free <port>` (listening TCP, v4+v6)
- `shell_cwd [shell-pid]` - prints the cwd of the most recently started
  student shell, or of the shell with the given PID (for hint scripts;
  note: `wait_cwd` accepts a match in any shell, so the two may
  disagree when several terminals are open)
- `shells` - lists all observed shells, one per line (`PID exe tty
  cwd`, most recent first); a debugging aid
- `hint_exit [task] <message>` - pushes a hint to the UI immediately
  and TERMINATES the check script with exit code 42; the task defaults
  to the current one (`$GYM_TASK`); the run counts as a rejected
  attempt (see **attempts** above)
- `set_var <NAME> <value>` - publishes a task var on the current unit:
  it joins the unit's vars, exported into later runs of the unit's own
  scripts and into scripts of units that `needs:` this unit (persists
  across daemon restarts; cleared by unit reset). THE way to pass
  small values between tasks and units - never stash them in files;
  use a file only for BLOB-like data (content, not a value). Task vars
  are script-env only: not interpolatable in markdown, not usable with
  `from:`. `GYM_` prefix reserved. Note it exits 0 on success, so do
  not let it be the last command after an unverified condition

All `wait_*` block until met; add `--timeout <sec>` for a bound or
`--now` for a single instant evaluation (the form for level tasks).
Regexes are Go RE2 (no backreferences). Checks compose freely with
shell:

```yaml
check: |
  wait_file --timeout 15 "$GYM_USER_HOME/junk.tmp" || \
    hint_exit "junk.tmp never appeared - was the init OK?"
  wait_file_gone "$GYM_USER_HOME/junk.tmp"
```

Choosing: verify **effects** (`wait_file`, `wait_port`, ...) over
commands; reserve `wait_exec` for commands that leave no trace (`ls`,
`cat`, `curl`) and `wait_line` for reps where the shape of the line is
the skill (`&&` vs `;`, a pipe, quoting, a builtin). `wait_exec` proves
the command was run, not that it succeeded. Never require one exact
command form when several are correct - keep `wait_exec` and
`wait_line` regexes permissive.

All scripts (`init:`, `check:`, `hint:`) run as root under
`bash -o pipefail`, each in its own session (no controlling tty - the
guard that keeps checks from matching themselves).
Environment available to scripts: unit vars (including task vars from
`set_var` and the vars of every unit listed in `needs:`), `GYM_UNIT`,
`GYM_TASK`, `GYM_USER` (observed login user), `GYM_USER_HOME` (that
user's home directory - never shell out to `getent` for it), and, in a
task's `check:`/`hint:`, `GYM_CHECK_ATTEMPT` (see **attempts** above).
`hint:` scripts additionally get
`GYM_TASK_EXIT`, `GYM_TASK_STDOUT`, `GYM_TASK_STDERR` from the last
failed check run - enough to diagnose why it failed and say something
specific. Hint script stdout replaces the task's hint area live
(rate-limited to one refresh per ~10 s).

## Reliable checks

A check's verdict must follow from what the student did in THIS rep,
not from what happened to be around. Design every check against both
failure modes, and treat each as a bug, not a corner case:

- a **false positive** completes the task on something that is not the
  rep - a command typed earlier or elsewhere, state left over from a
  previous session or unit, a wrapper process whose argv mentions the
  command, the check's own text, a lucky buffered event;
- a **false negative** never completes (or rejects with a hint) although
  the rep was done - an over-exact regex, a judged-once wrong answer
  that keeps winning, a fast command that was missed, a builtin that
  never exec'ed, a check killed by its own `--timeout`.

The student sees a green box that lies, or a red box that will not
turn green. Both destroy trust in the gym. Mechanisms and gotchas:

- **The event horizon is per unit, not per task.** Event checks judge
  everything buffered since the unit was activated, so a task behind
  `needs:` sees commands typed while it was locked (see the activation
  note above). Order two steps with a mark: step one
  `wait_exec ... || exit 1` then `set_var STEP1_SEQ "$(event_seq)"`,
  step two `wait_exec --after "$STEP1_SEQ" ...` (`wait_line` and
  `wait_env` too; one clock for all events). Place a command with
  `--cwd <dir>` when where it ran is the point. Never rely on the gate
  alone for order.
- **Judge the newest answer** in right/wrong branches: `--latest` +
  `case` + `hint_exit`. Without it, the oldest buffered wrong answer is
  judged again on every restart until the horizon moves - and a correct
  retry after it is never seen.
- **A rejection moves the horizon.** `hint_exit` and a plain `exit 1`
  both do, so only reject on evidence about the student's LATEST
  answer, never on an ambiguous or stale one. `wait_* --timeout N ||
  exit 1` is a rejection too: it silently hides everything the student
  did before it fired. For an idle nudge use a plain blocking wait and
  the `hint:` block.
- **State checks pass on pre-existing state.** `wait_file`, `wait_dir`,
  `wait_cwd`, `wait_port` look at the world, not at the student: a
  target that already exists (from a previous session, a sibling unit,
  a system default) completes the task at activation. Init must clear
  the target (`rm -f`, `pkill`, a fresh random name via `vars:`) so
  only the rep can satisfy the check; and a task verifying that
  something disappears must first confirm it existed (baseline negative
  check).
- **Effects over keystrokes.** `wait_exec` proves a command ran, not
  that it worked (`mkdir` on an existing path, `touch` on a read-only
  dir). Verify the effect (`wait_file`, `wait_dir`, `wait_file_newer`,
  `wait_port`, ...) and add `wait_exec` only for the how, or for
  commands that leave no trace.
- **Builtins leave no exec event.** `echo`, `cd`, `pwd`, `type`,
  `export`, `printf`, `true`, `false`, `[`, `read` - `wait_exec` and
  `wait_env` never see them. Observe the typed line with `wait_line`
  (`requires: [readline]`), verify the effect, or anchor the rep on an
  external command (`whoami`, `date`, `/bin/echo`). And a `wait_line`
  cannot be placed with `--cwd`: pair it with a `wait_exec --cwd` on
  the external command the line runs.
- **Text is not shape.** Joined argv cannot tell `date '+%A %d'` from
  `date +%A %d` - use `--argc N`. Exec events cannot tell `a; b` from
  `a && b` - use `wait_line`. Argv cannot tell an expanded glob from
  names typed by hand - accept both, or judge the line.
- **Regexes: permissive in form, strict in anchors.** Accept every
  correct spelling: short and long options in any order, `--opt=v` and
  `--opt v`, `(^|/)` before the command (argv0 may be a path), an
  optional trailing `/` on directories, `2s?` for `sleep 2`. Anchor
  both ends (`(^|/)date -u$`) so `date -u --bogus` or `cat date` do not
  match. Escape `$` (`\$`) and `.` (`\.`) inside double-quoted regexes
  built from vars, and remember RE2 has no backreferences.
- **`--cwd` is the exec'ed process's cwd read a moment after the exec.**
  It is the shell's cwd for ordinary commands; a program that chdirs at
  once (`tar -C`, `make -C`, `git -C`) may already show its target, and
  a process gone before the read never matches. A regex path follows
  `wait_cwd` rules (whole path, `^(...)$`).
- **Process patterns match wrappers and the check itself.**
  `pgrep -f`, `wait_proc*` see every argv that CONTAINS the text - a
  `bash -c "sleep 600; ..."` wrapper, `script -qec`, the check script.
  Use the bracket trick plus an argument only the target has and anchor
  the whole cmdline (`"^(/usr/bin/)?slee[p] ${N}s?\$"`); otherwise
  `wait_proc_gone` never fires and `wait_proc_state` reads the wrapper.
- **Several shells may be open.** `wait_cwd <path>` matches ANY shell
  of the student; when later steps must continue in the same one, keep
  its PID (`P=$(wait_cwd ...)`, `set_var SHELL_PID "$P"`) and pass it
  as the first argument from then on. `wait_exec` cannot be pinned to a
  shell - use `--cwd` and `--after` to narrow it.
- **Timing.** Commands typed at human speed are captured reliably;
  processes spawned in machine-speed loops can be missed - never depend
  on catching those. `wait_cwd` and the file checks poll every 200 ms;
  a very short-lived process can be gone before `wait_proc` looks, so
  wait for it with a bounded `--timeout N || true` and rely on the
  `_gone` or the effect afterwards. An `event_seq` mark is taken a
  moment after the matched event: a command run in that same instant
  falls before it (never at typing speed).
- **Level tasks re-evaluate and may flip back.** Use `--now` on every
  `wait_*` in a level task; edge tasks (the default) settle once.
- **Prove it with the solver and with alternatives.** `shellgym solve`
  types fast and exactly as written: a racy check or an over-exact
  regex fails there first. Then try the plausible alternative answers
  by hand (long options, a different order, `echo` instead of `ls`,
  a path instead of a name) and make sure each one either passes or
  gets a specific hint - never silence.

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
- `::tip` renders an always-visible callout (optional `title`, default
  "Tip") - for terminal-interaction technique that checks cannot verify
  (Tab completion, `Ctrl-R`, arrow-key history, pager keys). Hints answer
  "I'm stuck"; tips teach the ergonomic way to do what the task asks.
- `::note` and `::warn` share the tip's always-visible shape (optional
  `title`, defaults "Note" / "Warning") in neutral and amber colors -
  notes for context and asides, warnings for destructive or surprising
  behavior. Keep them short and next to the task they concern.
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
  `pgrep -f`, `wait_proc`, `wait_proc_gone`, and `wait_proc_state`
  matching the script's own text - anchor patterns on argv text only
  the target has (bracket trick + an argument: `'my-nam[e] 86400'`).
  When the target's full cmdline is known, ALSO anchor `^...$`
  (`"^(/usr/bin/)?slee[p] ${N}s?\$"`): substring patterns match any
  wrapper (`bash -c "sleep 600 ..."`, `script -qec ...`) whose argv
  merely *mentions* the command - fatal for `wait_proc_gone` (never
  gone while the wrapper lives) and `wait_proc_state` (the wrapper's
  state, not the target's).
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
completion, reporting PASS/FAIL per unit. Units hidden by the variant
draw are solved too (`[key=value, hidden]`), so one run covers every
variant. A unit whose `needs:` are not solved cannot be activated - even
by `solve --unit`; solve its chain first (or run without `--unit`, which
walks the path in order).

Debugging a failing unit - every check attempt's exit code, stdout,
stderr, and duration is recorded:

- the UI's **debug drawer** (press `d`; hidden in live mode);
- `GET /api/debug/<unit-id>`;
- on disk under `<state>/<path-id>/runs/`.

`POST /api/reset/<unit-id>` forgets a unit's progress and re-runs its
init from scratch - the fast iteration loop while authoring. `POST
/api/variants/<key>/<value>` switches the variant draw to preview the
other units (the debug drawer offers the same switch).
