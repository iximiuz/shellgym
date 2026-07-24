---
title: Return home from anywhere
needs: [one-level-up]
tasks:
  back_home:
    check: |
      # TRAVELER (a task var set by the one-level-up unit) is the PID of
      # the shell that walked the /tmp/depths tree; only THAT shell landing
      # exactly in $HOME counts - an idle second terminal already sitting
      # at home does not.
      if [ -z "$TRAVELER" ] || ! [ -d "/proc/$TRAVELER" ]; then
        sleep 5  # pace the restart loop - this state only clears on reset
        hint_exit "Lost track of the shell that walked the depths (was its terminal closed?). Reset the previous unit and climb again."
      fi
      wait_cwd "$TRAVELER" "$GYM_USER_HOME"
    hint: |
      CWD=$(shell_cwd "$TRAVELER" 2>/dev/null || echo "?")
      echo "Your shell is still in $CWD. Plain cd with no arguments takes you straight home."
    solve: |
      cd
---

Your shell is still parked somewhere in the `/tmp/depths` tree. You could
climb back with `..` steps, or spell out the full path - but there is a
shortcut: `cd` with no arguments takes you home from anywhere, no matter
how deep you wandered.

Come home with a single, argument-less command:

::task{name="back_home"}
#active
Waiting for your shell back in your home directory...
#completed
Home again. This move works from any depth of any tree.
::
