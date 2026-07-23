---
title: Find a process by name
requires: [systemd]
vars:
  DAEMON: { pick: [murmur-daemon, drone-daemon, hum-daemon] }
init:
  - name: start_daemon
    run: |
      systemctl stop gym-decoy.service 2>/dev/null || true
      pkill -f 'gym-deco[y] 86400' 2>/dev/null || true
      mkdir -p /run/gym-bin
      # The decoy is a tiny named sh wrapper (argv0 tricks like `exec -a`
      # or a renamed copy of sleep break on multi-call coreutils distros).
      # systemd-run gives it its own unit outside the daemon's cgroup, and
      # --uid runs it as the student so a plain `kill` works.
      cat > "/run/gym-bin/$DAEMON-gym-decoy" <<'WRAP'
      #!/bin/sh
      sleep "$1" &
      CHILD=$!
      trap 'kill "$CHILD" 2>/dev/null; exit 0' TERM INT HUP
      wait "$CHILD"
      WRAP
      chmod +x "/run/gym-bin/$DAEMON-gym-decoy"
      systemd-run --collect --quiet --unit=gym-decoy --uid="$GYM_USER" \
        "/run/gym-bin/$DAEMON-gym-decoy" 86400
      for i in $(seq 1 30); do
        pgrep -f 'gym-deco[y] 86400' >/dev/null && exit 0
        sleep 0.2
      done
      echo "decoy did not start" >&2
      systemctl status gym-decoy --no-pager 2>&1 | tail -5 >&2
      exit 1
tasks:
  pid_found:
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      PID=$(pgrep -f "$DAEMON-gym-decoy" | head -1)
      if [ -z "$PID" ]; then
        echo "decoy process is not running" >&2
        exit 2
      fi
      wait_file_contains "$HOME_DIR/daemon.pid" "^$PID\s*$"
    hint: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      if [ -f "$HOME_DIR/daemon.pid" ]; then
        echo "daemon.pid exists but does not hold the right PID. Find the process named ${DAEMON}-gym-decoy - pgrep -f matches against the full command line."
      else
        echo "List processes with ps aux (or ask pgrep directly), find ${DAEMON}-gym-decoy, and save its PID into ~/daemon.pid."
      fi
    solve: |
      pgrep -f $DAEMON-gym-decoy > ~/daemon.pid
      cat ~/daemon.pid
---

Somewhere on this machine a process named `${DAEMON}-gym-decoy` is
running. Find its PID and save it into `~/daemon.pid`.

`ps aux` shows every process (pipe it through `grep` to filter). `pgrep -f
name` skips the reading and prints matching PIDs directly.

::task{name="pid_found"}
#active
Waiting for the daemon's PID in `~/daemon.pid`...
#completed
Target acquired. Keep the PID around - the next unit finishes the job.
::
