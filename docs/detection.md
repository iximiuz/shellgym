# Detection Mechanisms

This page describes how Shell Gym observes the student without touching
their shell. Everything the built-in checks (see [checks.md](checks.md))
report is derived from five mechanisms: student shell discovery, exec
watching, command line watching, direct system-state polling, and the
check API socket that glues them to task scripts.

The guiding constraint is **zero instrumentation**: no prompt hooks, no
shell wrappers, no `PROMPT_COMMAND`, no pty interception. The student's
shell is a stock process; all observation happens from the outside.

## Student shell discovery (procfs scan)

Used by: `shell_cwd`, `shells`, `wait_cwd` built-ins.

The daemon scans `/proc` for processes that qualify as the student's
interactive shells. A process qualifies when all of the following hold:

1. it is owned by the observed user (`--user` / `shellUser` from
   `path.yaml`), by uid;
2. its `comm` is a known shell: `bash`, `zsh`, `sh`, `fish`, `dash`,
   `ash` (a leading `-`, as in login shells, is ignored);
3. it has a **controlling terminal** - field 7 (`tty_nr`) of
   `/proc/<pid>/stat` is non-zero. This is what separates interactive
   shells from shell-executed scripts, cron jobs, and the daemon's own
   task runners.

For each match the daemon reads the working directory via
`/proc/<pid>/cwd` (a live symlink - always the shell's *current* cwd)
and the tty name via `/proc/<pid>/fd/0`. Shells are ordered by process
start time; "the student's shell" for single-shell checks like
`shell_cwd` means the most recently started one, while `wait_cwd`
accepts a match in **any** of them (the student may legitimately have
several terminals open). Both accept an optional shell PID argument to
pin the check to one specific shell - `wait_cwd` prints the matched
shell's PID on success, which is how a unit hands "the shell that did
it" over to a dependent unit (via a `set_var` task var).

The scan runs on demand per check evaluation - there is no shell
tracking state to go stale.

## Exec watching (kernel proc connector)

Used by: `wait_exec`, `wait_env` built-ins.

The daemon subscribes to the **kernel proc connector** - a netlink
channel (`NETLINK_CONNECTOR`, `CN_IDX_PROC`) over which the kernel
multicasts process lifecycle events. Shell Gym listens for exec events
only. This requires `CONFIG_PROC_EVENTS` (standard everywhere) and
`CAP_NET_ADMIN`, which is the main reason the daemon runs as root. If
the connector cannot be opened, the daemon refuses to start rather than
run with silently broken exec checks.

### From event to record

An exec event carries just the pid. The connector fires *after* the new
program image is installed, so the daemon immediately harvests the rest
from procfs:

- `/proc/<pid>/cmdline` - the argv (NUL-separated);
- `/proc/<pid>/status` - uid and ppid;
- `/proc/<pid>/stat` - the controlling tty (`tty_nr`);
- `/proc/<pid>/environ` - the environment, captured **eagerly** for
  tty-attached processes and bounded at 32 KiB. Eager capture matters:
  fast commands are gone long before a `wait_env` check could read
  their environment lazily.

Events are recorded into a bounded in-memory ring buffer (4096 entries)
with monotonically increasing sequence numbers. Checks wait on the
buffer: a `wait_exec` first scans buffered events, then blocks until a
new matching one arrives.

### What counts as student activity

A recorded event matches only if:

- the uid is the observed user's - or **unknown**. A very short-lived
  process can exit between the exec event and the procfs read; such
  reads yield "unknown" (`-1`), and only a *confirmed* foreign uid
  disqualifies. Fast interactive commands (`ls`, `pwd`) routinely fall
  into the unknown bucket and must still count.
- the process has (or may have) a controlling tty. A confirmed
  tty-less process is rejected.

### Why checks never match themselves

Check scripts frequently contain the very pattern they search for
(`wait_exec 'ls -l'` - the string `ls -l` is right there in the
script's argv). Two design choices prevent self-matching:

1. every task script runs in its own session (`setsid`), which means
   **no controlling terminal** - so the daemon's own processes are
   confirmed tty-less and always rejected by the tty filter;
2. task scripts run as root, while the student is a non-root user - the
   uid filter rejects them independently.

### The activation horizon (`GYM_SINCE_EXEC_SEQ`)

When a unit activates, the engine snapshots the current event sequence
number and exports it to every task script as `GYM_SINCE_EXEC_SEQ`. Exec
checks only consider events **newer** than this horizon, so a command
the student ran before ever seeing the unit cannot satisfy it. (Command
line watching keeps its own ring and its own horizon,
`GYM_SINCE_LINE_SEQ`, exported the same way.) The horizon is per task.
It moves to the present when a check run rejects an answer - exits
non-zero on its own, see [Attempts](authoring-guide.md#attempts) - so
the restarted check judges only newer commands. A run killed by the task
timeout keeps its horizon, so no command is lost across such restarts.

### Limits

- Processes that exec and exit within the procfs-harvest window (well
  under a millisecond) can be missed entirely. Human-typed commands are
  reliably captured; content should not depend on catching commands
  spawned in tight machine-speed loops.
- The ring buffer holds 4096 events. On a quiet training box that is
  hours of activity; a busy background workload could evict old events
  faster.
- Matching is textual, against argv. `wait_exec` proves a command was
  run, not that it succeeded - verify effects where effects exist.

## Command line watching (readline uprobe)

Used by: `wait_line` built-in. Optional - gates the `readline` capability.

Exec events cannot tell `sleep 2; hostname` from `sleep 2 && hostname`:
both fork the same two children with the same argv. The operator exists
only inside the shell, in the line it parsed. Bash reads that line
through `readline()`, and the daemon observes the call from the kernel
side with a **uretprobe** (the technique behind bpftrace's
`bashreadline`): a return probe on the symbol whose fetch argument
copies the returned string into the trace buffer, registered through
tracefs (`uprobe_events`) with no BPF program involved. The student's
shell is still a stock process - the probe is a kernel breakpoint the
process never notices.

### Setup

At start the daemon:

1. resolves the observed user's login shell (`/etc/passwd`) and requires
   it to be bash;
2. finds the `readline` symbol - in the bash binary itself when readline
   is linked statically (Debian, Ubuntu), otherwise in the `libreadline`
   shared object bash loads (Fedora, Rocky) - and converts its virtual
   address to a file offset via the ELF program headers;
3. writes `r:shellgym/readline <binary>:0x<offset> line=+0($retval):string`
   to `uprobe_events`, creates the private tracefs instance
   `instances/shellgym` so other tracing users see none of it, enables
   the event there, and streams `trace_pipe` from the instance;
4. on every record, harvests the shell's uid and controlling tty from
   `/proc/<pid>` (shells are long-lived, so this practically never
   misses) and publishes the line into a ring buffer with the same
   sequence-number semantics as exec events.

A daemon that dies leaves the probe and the instance behind (they are
kernel state); the next start removes them before registering afresh.
`Close` disables and deletes the event and the instance.

### Capability gating

When any step fails - not root, no tracefs or `CONFIG_UPROBE_EVENTS`, a
non-bash login shell, no readline symbol - the daemon logs the reason and
runs without the watcher. Units that declare `requires: [readline]` are
then marked unsupported (browsable, never activated, checks off), like
any other unmet `requires:`. The `/line/*` endpoints answer `501` so a
`wait_line` in a unit that forgot the declaration fails at once instead
of hanging until its timeout.

### What is (and is not) observed

- Every line an interactive bash of the observed user reads: the text
  as typed, before history expansion, alias expansion, or parsing.
  Continuation lines are separate events. `wait_line` trims surrounding
  whitespace and matches the rest verbatim.
- Non-interactive bash (`bash -c`, scripts, the daemon's own task
  runners) never calls readline, so there is nothing to filter out.
  `shellgym solve` types into a real interactive bash and is observed
  like a student - which is why the solve driver syncs through its own
  `PROMPT_COMMAND` rather than appending anything to the typed line.
- Other shells (zsh, fish, dash) produce no events; the login shell must
  be bash for the capability to be detected at all.
- The activation horizon is a separate sequence number
  (`GYM_SINCE_LINE_SEQ`, alongside `GYM_SINCE_EXEC_SEQ` for exec events).

## Direct system-state polling (procfs and friends)

Used by: `wait_file`, `wait_file_gone`, `wait_file_contains`,
`wait_dir`, `wait_file_mode`, `wait_file_newer`,
`wait_proc`, `wait_proc_gone`, `wait_port`, `wait_port_free` built-ins.

These checks do not go through the daemon at all - the check process
inspects the system directly, re-evaluating every 200 ms until the
condition holds (or `--timeout`/`--now` says stop):

- **Files** - `filepath.Glob` against the given path or glob pattern;
  existence of *anything* at the path counts (files, directories,
  sockets), except `wait_dir`, which additionally `stat`s the matches
  and requires a directory. `wait_file_contains` reads the file and
  applies an RE2 regex in multiline mode. `wait_file_mode` compares
  `stat` mode bits (`mode & 07777`) against the given octal;
  `wait_file_newer` compares the two files' mtimes.
- **Processes** - a scan of `/proc/<pid>/cmdline` for every process on
  the box, NULs replaced with spaces, matched against the regex. Only
  the check's own process is excluded - notably *not* the bash running
  the check script, whose cmdline contains the script text (hence the
  bracket-trick advice in [checks.md](checks.md)).
- **Ports** - `/proc/net/tcp` and `/proc/net/tcp6` parsed for sockets
  in `LISTEN` state (`0A`) on the given port, any local address. TCP
  only.

Since checks run as root, no permission games are needed to inspect
other users' processes and files.

## The check API socket

The built-in checks are PATH shims: the engine writes tiny
`#!/bin/sh` wrappers (one per check name) into `<run>/bin/` and prepends
that directory to every task script's PATH. Each shim execs
`shellgym check <name> ...` - the same binary as the daemon, in a
one-shot client role.

Checks that need daemon-side state talk to it over a unix socket
(`<run>/gym.sock`, mode 0600, root-only):

| Endpoint | Used by | Purpose |
|---|---|---|
| `/shells` | `shell_cwd`, `shells`, `wait_cwd` | current student-shell list |
| `/exec/wait` | `wait_exec`, `wait_env` | block until a matching exec event |
| `/exec/seq`, `/exec/snapshot` | debugging | event-stream introspection |
| `/line/wait` | `wait_line` | block until a matching command line is read (`501` without the `readline` capability) |
| `/line/seq`, `/line/snapshot` | debugging | line-stream introspection |
| `/hint` | `hint_exit` | push a hint to the UI |
| `/vars` | `set_var` | publish a task var on the current unit |
| `/units/watch` | external observers | SSE stream of unit statuses: a full snapshot on connect, then one event per change |

The socket path and the activation horizons reach the shims through the
`GYM_SOCK`, `GYM_SINCE_EXEC_SEQ` (exec events), and `GYM_SINCE_LINE_SEQ`
(command lines) environment variables the engine sets for every script.

`/units/watch` is a Server-Sent Events stream, not a check. It opens with a
`snapshot` event (`{"path": "...", "units": {"<unit id>": "pending|active|
completed|unsupported|hidden", ...}}` - `hidden` units are left out of this
attempt by the variant draw and need not complete), then emits a `unit` event
(`{"id": "...", "status": "..."}`) whenever a unit's status changes, plus a
`: keepalive` comment every 15 seconds so a client can tell a quiet daemon
from a dead one. A client resyncs from the snapshot on every (re)connect, so a
dropped connection never loses an event.

`hint_exit` gets one extra piece of machinery: the engine prepends a
shell function to every script that calls the real binary and then
`exit 42`s the script itself. That is how a plain
`wait_x ... || hint_exit "look here"` both reports the hint and stops
the attempt with a distinct, recognizable exit code.

## Script execution environment

The runner that executes `init:`, `check:`, and `hint:` scripts is part
of the detection story - it guarantees the isolation properties above:

- scripts run as root under `bash -o pipefail`;
- each script gets its own session via `setsid`: no controlling tty (the
  self-match guard) and a clean process group that can be killed as a
  whole tree on timeout;
- per-run timeouts: 30 s for edge checks, 10 s for level polls, 60 s
  for init scripts, 10 s for hint scripts (a task can override its check
  timeout via `timeout:` in the frontmatter); on expiry the entire
  process group receives SIGKILL;
- stdout/stderr are captured through real pipes, so a script may leave
  long-lived background children without wedging the runner, and every
  run is recorded (exit code, streams, duration) for the debug
  drawer and `/api/debug`.
