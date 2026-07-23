---
title: The night watch
requires: [systemd]
vars:
  WATCH: { pick: [night-watch, dawn-watch, dusk-watch] }
init:
  - name: start_watcher
    run: |
      systemctl stop gym-watcher.service 2>/dev/null || true
      pkill -f 'gym-watche[r] 86400' 2>/dev/null || true
      mkdir -p /run/gym-bin
      # The decoy is a tiny named sh wrapper (argv0 tricks like `exec -a`
      # or a renamed copy of sleep break on multi-call coreutils distros).
      # systemd-run gives it its own unit outside the daemon's cgroup, and
      # --uid runs it as the student so a plain `kill` works.
      cat > "/run/gym-bin/$WATCH-gym-watcher" <<'WRAP'
      #!/bin/sh
      sleep "$1" &
      CHILD=$!
      trap 'kill "$CHILD" 2>/dev/null; exit 0' TERM INT HUP
      wait "$CHILD"
      WRAP
      chmod +x "/run/gym-bin/$WATCH-gym-watcher"
      systemd-run --collect --quiet --unit=gym-watcher --uid="$GYM_USER" \
        "/run/gym-bin/$WATCH-gym-watcher" 86400
      for i in $(seq 1 30); do
        pgrep -f 'gym-watche[r] 86400' >/dev/null && exit 0
        sleep 0.2
      done
      echo "watcher did not start" >&2
      exit 1
tasks:
  own_sentry:
    check: |
      wait_proc '(^|/)sleep 7207'
    hint: |
      echo "Start 'sleep 7207' so it keeps running while your prompt stays free. One extra character at the end of the line does it."
    solve: |
      sleep 7207 &
  watcher_relieved:
    check: |
      wait_proc --timeout 15 "$WATCH-gym-watcher" || {
        echo "the watcher process is not running" >&2
        exit 2
      }
      wait_proc_gone "$WATCH-gym-watcher"
    hint: |
      echo "Find the ${WATCH}-gym-watcher process by name and terminate it. No PID file this time - a name-based lookup gets you there."
    solve: |
      kill $(pgrep -f $WATCH-gym-watcher)
---

Two duties tonight, both from memory:

First, post your own sentry: a background `sleep 7207`.

::task{name="own_sentry"}
#active
Waiting for a background `sleep 7207`...
#completed
Sentry posted without losing your prompt.
::

Second, relieve the old guard: a process named `${WATCH}-gym-watcher` is
still on duty somewhere. Find it and terminate it.

::task{name="watcher_relieved"}
#active
Waiting for `${WATCH}-gym-watcher` to be terminated...
#completed
Watch rotated. Finding a process by name and signaling it is a move you
will make for the rest of your Linux life - it is now yours.
::
