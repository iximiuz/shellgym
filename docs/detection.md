# Detection Mechanisms

This page describes how Shell Gym observes the student without touching
their shell. Everything the built-in checks (see [checks.md](checks.md))
report is derived from four mechanisms: student shell discovery, exec
watching, direct system-state polling, and the check API socket that
glues them to task scripts.

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
several terminals open).

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

### The activation horizon (`GYM_SINCE_SEQ`)

When a unit activates, the engine snapshots the current event sequence
number and exports it to every task script as `GYM_SINCE_SEQ`. Exec
checks only consider events **newer** than this horizon, so a command
the student ran before ever seeing the unit cannot satisfy it. Within
one unit attempt the horizon is fixed - a check that is restarted (for
example after a `hint_exit`) still sees everything the student did since
activation, so no command is lost between check restarts.

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

## Direct system-state polling (procfs and friends)

Used by: `wait_file`, `wait_file_gone`, `wait_file_contains`,
`wait_proc`, `wait_proc_gone`, `wait_port`, `wait_port_free` built-ins.

These checks do not go through the daemon at all - the check process
inspects the system directly, re-evaluating every 200 ms until the
condition holds (or `--timeout`/`--now` says stop):

- **Files** - `filepath.Glob` against the given path or glob pattern;
  existence of *anything* at the path counts (files, directories,
  sockets). `wait_file_contains` reads the file and applies an RE2
  regex in multiline mode.
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
| `/hint` | `hint_exit` | push a hint to the UI |

The socket path and the activation horizon reach the shims through the
`GYM_SOCK` and `GYM_SINCE_SEQ` environment variables the engine sets for
every script.

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
- per-attempt timeouts: 30 s for edge checks, 10 s for level polls, 60 s
  for init scripts, 10 s for hint scripts (a task can override its check
  timeout via `timeout:` in the frontmatter); on expiry the entire
  process group receives SIGKILL;
- stdout/stderr are captured through real pipes, so a script may leave
  long-lived background children without wedging the runner, and every
  attempt is recorded (exit code, streams, duration) for the debug
  drawer and `/api/debug`.
