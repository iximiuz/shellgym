---
title: Terminate a process
requires: [systemd]
needs: [hunt-the-daemon]
vars:
  DAEMON: { from: hunt-the-daemon.DAEMON }
tasks:
  daemon_down:
    check: |
      # Baseline first: the decoy must be alive before we watch for its death.
      wait_proc --timeout 15 "$DAEMON-gym-decoy" || {
        echo "decoy process is not running (was it killed already, or did init fail?)" >&2
        exit 2
      }
      wait_proc_gone "$DAEMON-gym-decoy"
    hint: |
      echo "kill sends a signal to a PID. The default signal (TERM) politely asks the process to exit. You already saved the PID in ~/daemon.pid."
    solve: |
      kill $(cat ~/daemon.pid)
---

You located `${DAEMON}-gym-decoy` and saved its PID. Now terminate it
with the `kill` command. Despite the name, `kill` just sends a signal; the
default `TERM` signal asks the process to shut down.

::task{name="daemon_down"}
#active
Waiting for `${DAEMON}-gym-decoy` to terminate...
#completed
Process gone. If a process ignores `TERM`, the last resort is `kill -9`
(`KILL`), which the process cannot ignore or handle.
::

::hint{title="Reading the PID back"}
`cat ~/daemon.pid` shows the number again. Command substitution can inline
it: `kill $(cat ~/daemon.pid)`.
::
