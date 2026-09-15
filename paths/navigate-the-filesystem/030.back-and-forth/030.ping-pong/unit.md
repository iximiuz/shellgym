---
title: Ping-pong
requires: [readline]
init:
  - name: build_work_tree
    run: |
      W="$GYM_USER_HOME/projects"
      mkdir -p "$W/reports/drafts" "$W/reports/final" "$W/archive/2025" "$W/data/imports" "$W/data/exports" "$W/notes"
      [ -f "$W/reports/final/summary.txt" ] || echo "Final report, approved." > "$W/reports/final/summary.txt"
      [ -f "$W/reports/drafts/outline.txt" ] || echo "Report outline, work in progress." > "$W/reports/drafts/outline.txt"
      [ -f "$W/archive/2025/report-2025.txt" ] || echo "Last year's report." > "$W/archive/2025/report-2025.txt"
      [ -f "$W/data/imports/inventory.csv" ] || echo "id,name,quantity" > "$W/data/imports/inventory.csv"
      [ -f "$W/data/exports/inventory-export.csv" ] || echo "id,name,quantity" > "$W/data/exports/inventory-export.csv"
      [ -f "$W/notes/todo.txt" ] || echo "- review the final report" > "$W/notes/todo.txt"
      chown -R "$GYM_USER:$GYM_USER" "$W"
tasks:
  first:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/projects/archive/2025") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/projects/archive/2025. From anywhere, cd ~/projects/archive/2025 gets you there."
    solve: |
      cd ~/projects/archive/2025
  second:
    needs: [first]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" /var/tmp
    hint: |
      echo "Move to /var/tmp with its full path."
    solve: |
      cd /var/tmp
  third:
    needs: [second]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects/archive/2025" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        -) exit 0 ;;
        *) hint_exit "You are back in 2025, but not through the dash. Return to /var/tmp and bounce back with cd -." ;;
      esac
    hint: |
      echo "Bounce back with cd and a single dash."
    solve: |
      cd -
  fourth:
    needs: [third]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" /var/tmp >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        -) exit 0 ;;
        *) hint_exit "You are in /var/tmp, but you typed the path. The dash now points there: go back to 2025 and bounce with cd - again." ;;
      esac
    hint: |
      echo "The dash works in both directions. Use cd - once more."
    solve: |
      cd -
---

Every `cd -` also updates the memory, so the dash works in both
directions: two directories can be alternated with the same two
characters, again and again.

Start in `~/projects/archive/2025` and move to `/var/tmp` (a second
scratch space, whose files survive a restart). Then bounce back to
`2025`, and forward to `/var/tmp` again, using only the dash.

::task{name="first"}
#active
Waiting for your shell to be in `~/projects/archive/2025`...
#completed
You are in `2025`. Now go to `/var/tmp`.
::

::task{name="second"}
#active
Waiting for your shell to arrive in `/var/tmp`...
#completed
You are in `/var/tmp`. Now start bouncing.
::

::task{name="third"}
#active
Waiting for `cd -` to bring you back to `2025`...
#completed
You are back in `2025`. Now bounce once more, in the other direction.
::

::task{name="fourth"}
#active
Waiting for a second `cd -` to land you in `/var/tmp`...
#completed
The dash can bounce between two directories any number of times.
::

::tip
Press the **Up** arrow to get the previous `cd -` back instead of
retyping it. Each bounce then takes two keys.
::
