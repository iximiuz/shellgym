---
title: Adjust file permissions
init:
  - name: seed_files
    run: |
      echo "the launch code is 0000" > "$GYM_USER_HOME/secret.txt"
      chmod 644 "$GYM_USER_HOME/secret.txt"
      printf '#!/bin/sh\necho "greetings from greet.sh" > "$HOME/greet.done"\n' > "$GYM_USER_HOME/greet.sh"
      chmod 644 "$GYM_USER_HOME/greet.sh"
      rm -f "$GYM_USER_HOME/greet.done"
      chown "$GYM_USER" "$GYM_USER_HOME/secret.txt" "$GYM_USER_HOME/greet.sh"
tasks:
  lock_secret:
    # The check's own poll loop runs up to ~50s before falling through to
    # hint_exit; the attempt timeout must exceed that or the hint never fires.
    timeout: 60
    check: |
      PERMS=""
      for i in $(seq 1 50); do
        PERMS=$(stat -c %a "$GYM_USER_HOME/secret.txt" 2>/dev/null || true)
        [ "$PERMS" = "600" ] && exit 0
        sleep 1
      done
      if [ -z "$PERMS" ]; then
        hint_exit "secret.txt disappeared from your home directory - recreate it (any content) and set its permissions."
      else
        hint_exit "secret.txt permissions are currently $PERMS - the target is 600 (owner read+write, nothing for anyone else)."
      fi
    solve: |
      chmod 600 ~/secret.txt
  run_greeter:
    # wait_file blocks 45s before hint_exit; see lock_secret.
    timeout: 60
    check: |
      wait_file --timeout 45 "$GYM_USER_HOME/greet.done" || {
        hint_exit run_greeter "No greet.done yet. Check ls -l ~/greet.sh - does it have an x among the permissions? Only then can ./greet.sh run."
      }
    solve: |
      chmod +x ~/greet.sh
      ~/greet.sh
---

Two files just landed in your home directory, and both have the wrong
permissions.

First, `secret.txt` is readable by everyone on the machine (mode `644`).
Restrict it with `chmod` so that only you can read and write it - mode
`600`. You can inspect any file's current mode with `ls -l` or `stat`.

::task{name="lock_secret"}
#active
Waiting for `~/secret.txt` to become private (mode `600`)...
#completed
Locked down. The three digits are owner/group/others; `600` leaves the
last two empty.
::

Second, `greet.sh` is a script, but it is not runnable yet: executing a
file requires the execute bit. Add it, then run the script:

::task{name="run_greeter"}
#active
Waiting for `greet.sh` to be made executable and run...
#completed
It ran and left `greet.done` behind. `chmod +x` is the move you will make
for every script you ever write.
::

::hint
---
title: chmod in ten seconds
---
`chmod 600 file` sets an exact mode; `chmod +x file` adds one bit to the
current mode. Both forms are everyday tools.
::
