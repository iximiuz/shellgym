---
title: Back where you were
requires: [readline]
tasks:
  at_home:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start from your home directory. A bare cd takes you there."
    solve: |
      cd
  away:
    needs: [at_home]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" /tmp
    hint: |
      echo "Move to /tmp with its full path."
    solve: |
      cd /tmp
  back:
    needs: [away]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        -) exit 0 ;;
        *) hint_exit "You are home, but not through the dash. Go to /tmp again and come back with cd and a single dash as the argument." ;;
      esac
    hint: |
      echo "The argument is a single dash, so the command is cd -."
    solve: |
      cd -
---

The shell remembers the directory you were in before the last move,
and `cd -` (a single dash as the argument) takes you back there. It
also prints where you landed, so you can check the jump.

Start at home, move to `/tmp`, and then come back with the dash.

::task{name="at_home"}
#active
Waiting for your shell to be in your home directory...
#completed
You are home. Now go to `/tmp`.
::

::task{name="away"}
#active
Waiting for your shell to arrive in `/tmp`...
#completed
You are in `/tmp`. Now come back with the dash.
::

::task{name="back"}
#active
Waiting for `cd -` to return you home...
#completed
The dash brought you back to where you were.
::
