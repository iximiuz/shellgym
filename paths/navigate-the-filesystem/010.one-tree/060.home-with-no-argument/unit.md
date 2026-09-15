---
title: Home with no argument
requires: [readline]
tasks:
  away:
    check: |
      P=$(wait_cwd /var) || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Move to /var first."
    solve: |
      cd /var
  home:
    needs: [away]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      case "$LINE" in
        cd) exit 0 ;;
        *) hint_exit "You are home, but you gave cd an argument. Move away and come back with cd on its own, with nothing after it." ;;
      esac
    hint: |
      echo "Run cd with nothing after it. No argument means home."
    solve: |
      cd
---

You will return to your home directory more often than to any other
place, so the shell makes it the easiest move. A bare `cd`, with
**no argument at all**, takes you home from anywhere.

First move to `/var`. It holds data that changes while the machine
runs, such as logs. Then go home the short way.

::task{name="away"}
#active
Waiting for your shell to arrive in `/var`...
#completed
You are in `/var`. Now go home without typing its path.
::

::task{name="home"}
#active
Waiting for a bare `cd` to bring you home...
#completed
You are home, and it took two letters.
::
