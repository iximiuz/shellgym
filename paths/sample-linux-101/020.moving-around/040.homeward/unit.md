---
title: Return home from anywhere
tasks:
  left_home:
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      while :; do
        CWD=$(shell_cwd 2>/dev/null || true)
        case "$CWD" in
          ""|"$HOME_DIR"|"$HOME_DIR"/*) sleep 0.5 ;;
          *) exit 0 ;;
        esac
      done
    hint: |
      echo "First wander off: cd into /tmp, /var, or anywhere outside your home directory."
    solve: |
      cd /var/log
  back_home:
    needs: [left_home]
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_cwd "$HOME_DIR"
    hint: |
      echo "Now come home. Plain cd with no arguments takes you straight there."
    solve: |
      cd
---

`cd` with no arguments takes you home from anywhere, no matter how deep you
wandered. Prove it in two moves.

First, go somewhere outside your home directory, for example `/var/log` or
`/tmp`:

::task{name="left_home"}
#active
Waiting for your shell to leave the home directory...
#completed
Far from home now.
::

Then return home with a single, argument-less command:

::task{name="back_home"}
#active
Waiting for your shell back in your home directory...
#completed
Home again. This move works from any depth of any tree.
::
