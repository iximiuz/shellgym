# Built-in Checks Reference

Task scripts (`check:`, `hint:`, and `init:` blocks) run with a set of
built-in commands on their PATH. Each is a tiny shim that execs
`shellgym check <name>`, which either inspects the system directly or
talks to the daemon over its unix socket (see
[detection.md](detection.md) for the mechanisms underneath).

## Common behavior

- **Exit codes**: `0` = condition met, `1` = not met (deadline passed or
  one-shot check false), `2` = usage or communication error.
- **Blocking**: every `wait_*` check blocks until its condition is met.
  Two flags adjust that:
  - `--timeout <sec>` - give up (exit 1) after this many seconds;
    fractional values are accepted.
  - `--now` - a single instant evaluation, no waiting. This is the form
    to use in `level` tasks, which are re-polled by the engine anyway.
- **Polling checks** (`wait_cwd`, `wait_file*`, `wait_proc*`,
  `wait_port*`) re-evaluate their condition every 200 ms in the shim
  process. **Event checks** (`wait_exec`, `wait_env`) block inside the
  daemon and wake up on the next matching exec event.
- **Regex flavor** is Go's RE2 (no backreferences). Patterns with shell
  metacharacters must be quoted.
- **Environment**: the engine exports the unit's vars (including task
  vars published with `set_var`, and the vars of every unit listed in
  the unit's `needs:`) plus `GYM_UNIT`, `GYM_TASK`, `GYM_USER` (the
  observed login user), `GYM_USER_HOME`, `GYM_SINCE_SEQ` (exec-event
  horizon, see `wait_exec`), and `GYM_SOCK` (the daemon socket path)
  into every script. `hint:` scripts additionally get
  `GYM_TASK_EXIT`, `GYM_TASK_STDOUT`, `GYM_TASK_STDERR` from the last
  failed check run (streams clipped to 1 KiB).

Checks compose freely with shell:

```yaml
check: |
  wait_file --timeout 15 "$GYM_USER_HOME/junk.tmp" || \
    hint_exit "junk.tmp never appeared - have you tried creating it?"
  wait_file_gone "$GYM_USER_HOME/junk.tmp"
```

## Shell state

### `shell_cwd [shell-pid]`

Prints the current working directory of the interactive shell of the observed user -
if no shell PID provided, the **most recently started** shell is used.
Non-blocking. Exits 2 if no (matching) shell is found. Typical use is inside `hint:` scripts:

```yaml
hint: |
  echo "Your shell is still in $(shell_cwd)."
```

### `shells`

Lists all observed interactive shells, one per line, tab-separated:
`PID exe tty cwd`, most recently started first. A debugging aid; rarely
needed in real checks.

### `wait_cwd [shell-pid] <path-or-regex>`

Waits until an interactive shell of the observed user has the given
working directory. Which shell counts is picked by the optional first
argument:

- `wait_cwd <path>` - **any** open shell may match (the student may
  legitimately have several terminals open);
- `wait_cwd <shell-pid> <path>` - only the shell with that **specific**
  PID counts.

The path is an exact absolute path; if it contains regex metacharacters
(`?*[](|^$+`) and compiles, it is treated as a regex matched against the
whole path (auto-anchored as `^(...)$`):

```yaml
check: |
  wait_cwd "/tmp/gym/$DIRNAME"          # exact
  wait_cwd "/tmp/gym/(alpha|bravo)"     # regex
  wait_cwd "$TRAVELER" "$GYM_USER_HOME" # only this one shell
```

The cwd is read live from `/proc/<pid>/cwd`, so this reflects where the
shell is *now*, not where it has been.

On success `wait_cwd` prints the PID of the shell that matched. Publish
it as a [task var](#task-vars) when a later task (or a dependent unit)
must track the **same** shell instead of accepting a match in any
terminal:

```yaml
# earlier unit
check: |
  TRAVELER=$(wait_cwd "/tmp/gym/$DIRNAME") || exit 1
  set_var TRAVELER "$TRAVELER"
# dependent unit (needs: [earlier-unit])
check: |
  wait_cwd "$TRAVELER" "$GYM_USER_HOME"
```

## Command execution

### `wait_exec <regex>`

Waits until the student runs a command whose **full argv** (joined with
single spaces) matches the regex. Matching is scoped to student
activity:

- only commands executed **after the current unit was activated** count
  (earlier history never satisfies a fresh unit);
- only processes of the observed user with a controlling terminal count;
  the daemon's own scripts can never match (see
  [detection.md](detection.md));
- matched commands are buffered: if the student already ran the command
  a moment before the check (re)started, the check still passes
  immediately - checks do not miss commands between their own restarts.

```yaml
check: |
  wait_exec '(^|/)ls (-[a-zA-Z]+ )*-l'   # the student ran ls with -l
```

Match on argv text, not on outcomes: `wait_exec` proves the command was
*run*, not that it succeeded. When the effect matters, verify the effect
(`wait_file`, `wait_port`, ...) and use `wait_exec` for commands that
leave no trace (`ls`, `cat`, `ps`, `curl`).

Very short-lived processes are harvested from `/proc` right after the
exec event; in the rare case the process vanishes before its argv could
be read, the event is dropped. Commands typed at human speed are
reliably captured.

### `wait_env <NAME> [regex]`

Waits for an executed command whose **environment** contains variable
`NAME` (with a value matching `regex`, if given). This is the way to
verify exports: `/proc/<pid>/environ` of a running shell shows only its
*initial* environment, but every command the shell spawns inherits the
current one - so ask the student to run any command after exporting:

```yaml
check: |
  wait_env GREETING '^hello$'
```

The same student-activity scoping as `wait_exec` applies. Environments
are captured at exec time (bounded at 32 KiB), so even fast commands
are inspected reliably.

## Files

### `wait_file <path-or-glob>`

Waits until a file (or anything - directory, socket, ...) matching the
path or glob exists. Globs use `filepath.Glob` syntax:

```yaml
check: |
  wait_file "/tmp/gym/reports/*.txt"
```

### `wait_file_gone <path-or-glob>`

Waits until **nothing** matches the path or glob. Pair it with a
baseline so the task cannot auto-solve before the setup ran:

```yaml
check: |
  wait_file --timeout 15 "$GYM_USER_HOME/junk.tmp" || exit 1
  wait_file_gone "$GYM_USER_HOME/junk.tmp"
```

### `wait_file_contains <path> <regex>`

Waits until the file exists and its content matches the regex. The
pattern is compiled in multiline mode (`(?m)`), so `^` and `$` anchor to
lines - `'^done$'` means "a line that is exactly `done`":

```yaml
check: |
  wait_file_contains "$GYM_USER_HOME/notes.txt" "^$TOKEN$"
```

## Processes

### `wait_proc <regex>`

Waits until some process's **full command line** (argv joined with
spaces) matches the regex. All processes on the box are scanned, not
just the student's.

### `wait_proc_gone <regex>`

Waits until **no** process command line matches the regex.

**Pitfall - self-matching.** The check script itself runs under a bash
process whose command line contains the whole script text, including
your pattern. Only the check shim's own process is excluded from the
scan. Anchor patterns with the bracket trick so the literal pattern text
cannot match itself:

```yaml
check: |
  wait_proc 'sleep 8640[0]'        # matches "sleep 86400", not this script
  wait_proc_gone 'hum-servic[e]'
```

The same trick applies to `pgrep -f` in init scripts.

## Network

### `wait_port <port>`

Waits until some local TCP socket is **listening** on the port (IPv4 or
IPv6, any address), per `/proc/net/tcp` and `/proc/net/tcp6`.

### `wait_port_free <port>`

Waits until nothing is listening on the port. Like `wait_file_gone`,
give it a baseline when the unit starts with the port occupied:

```yaml
check: |
  wait_port --timeout 15 "$PORT" || hint_exit "the server never started"
  wait_port_free "$PORT"
```

Only listening TCP sockets are considered - UDP and established
connections are invisible to these checks.

## Task vars

### `set_var <NAME> <value>`

Persists a **task var** on the current unit. The var joins the unit's
vars: it is exported into every subsequent run of the unit's own
scripts, and into the scripts of any unit that declares this unit in
its `needs:`. Vars survive daemon restarts and are cleared by a unit
reset, exactly like frontmatter vars.

This is the way to pass small values between tasks and on to dependent
units - which shell did it, a generated token, a chosen port. Do **not**
stash such values in files; reach for a file only when the data is
BLOB-like (a log to diff, a directory tree, anything that is content
rather than a value).

```yaml
check: |
  TRAVELER=$(wait_cwd "/tmp/gym/$DIRNAME") || exit 1
  set_var TRAVELER "$TRAVELER"
```

Task vars are script-environment only: their value does not exist when
the unit's markdown is rendered, so they cannot be interpolated into
the body (`${...}`) or referenced with `vars: { from: ... }`. Names
follow the usual `[A-Za-z_][A-Za-z0-9_]*` shape; the `GYM_` prefix is
reserved. Exits non-zero on error, but note that a *successful*
`set_var` exits 0 - end the check with it only when everything before
it already proved the condition.

## Hints

### `hint_exit [task] <message>`

Pushes a hint message to the UI **immediately** and terminates the whole
check script with the distinct exit code 42. The task defaults to the
current one (`$GYM_TASK`); pass a task name explicitly only when a
script needs to hint a different task.

The termination matters: the engine injects a shell-function wrapper
around the binary, so `hint_exit` behaves like `exit` for the calling
script. The canonical shape is a bounded wait with a fallback message:

```yaml
check: |
  wait_cwd --timeout 60 "/tmp/gym/$DIRNAME" || \
    hint_exit "Still waiting... check where your shell is with pwd."
```

The failed attempt is recorded, the engine restarts the check after a
short delay, and the hint stays visible in the task box until replaced.
`hint_exit` only works inside task scripts (it needs `GYM_UNIT`,
`GYM_TASK`, and the daemon socket).

## Choosing the right check

| To verify... | Use |
|---|---|
| the shell moved somewhere | `wait_cwd` |
| a command was run (no lasting effect) | `wait_exec` |
| a variable was exported | `wait_env` |
| a file/directory was created | `wait_file` |
| a file was removed | baseline + `wait_file_gone` |
| file content | `wait_file_contains` |
| a process is running | `wait_proc` |
| a process was stopped | baseline + `wait_proc_gone` |
| a server is up | `wait_port` |
| a server was stopped | baseline + `wait_port_free` |

Prefer verifying **effects** over commands: effects are what the student
actually needs to achieve, and they leave the student free to find their
own way there.
