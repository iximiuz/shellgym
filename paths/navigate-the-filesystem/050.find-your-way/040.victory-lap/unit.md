---
title: Victory lap
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
  etc:
    check: |
      P=$(wait_cwd /etc) || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Move to /etc with its full path."
    solve: |
      cd /etc
  drafts:
    needs: [etc]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects/reports/drafts" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        "~"*) exit 0 ;;
        *) hint_exit "You reached drafts, but not in one command through the tilde. Go back to /etc and write the path from ~ down." ;;
      esac
    hint: |
      echo "One command: cd and the path starting with the tilde, ~/projects/reports/drafts."
    solve: |
      cd ~/projects/reports/drafts
  final:
    needs: [drafts]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects/reports/final" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        ../*) exit 0 ;;
        *) hint_exit "You reached final, but not through the parent. Go back to drafts and use a path that starts with the two-dot name." ;;
      esac
    hint: |
      echo "Up one with the two-dot name, then down into final, in one relative path."
    solve: |
      cd ../final
  back:
    needs: [final]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects/reports/drafts" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        -) exit 0 ;;
        *) hint_exit "You are back in drafts, but you typed the path. Return to final and come back with the dash." ;;
      esac
    hint: |
      echo "The dash remembers the previous directory: cd -"
    solve: |
      cd -
  home:
    needs: [back]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      case "$LINE" in
        cd) exit 0 ;;
        *) hint_exit "You are home, but with an argument. Go back to drafts and finish with the two-letter form." ;;
      esac
    hint: |
      echo "Run cd with nothing after it."
    solve: |
      cd
---

This is the last unit, and it is a lap through every kind of move in
this gym:

1. Move to `/etc` by its full path.
2. Jump to `~/projects/reports/drafts` in one command, through the tilde.
3. Cross over to the sibling `final` with a relative path through the parent.
4. Bounce back to `drafts` with the dash.
5. Go home the shortest way.

::task{name="etc"}
#active
Waiting for your shell to arrive in `/etc`...
#completed
That is the absolute path done.
::

::task{name="drafts"}
#active
Waiting for a single `cd` through `~` into `drafts`...
#completed
That is the tilde done.
::

::task{name="final"}
#active
Waiting for a `cd` through `..` into `final`...
#completed
That is the relative path done.
::

::task{name="back"}
#active
Waiting for `cd -` to bring you back to `drafts`...
#completed
That is the dash done.
::

::task{name="home"}
#active
Waiting for a bare `cd` to bring you home...
#completed
You are home, and that is the whole gym.
::
