---
title: Look and return
requires: [readline]
init:
  - name: build_work_tree
    run: |
      W="$GYM_USER_HOME/work"
      mkdir -p "$W/reports/drafts" "$W/reports/final" "$W/archive/2025" "$W/data/imports" "$W/data/exports" "$W/notes"
      [ -f "$W/reports/final/summary.txt" ] || echo "Final report, approved." > "$W/reports/final/summary.txt"
      [ -f "$W/reports/drafts/outline.txt" ] || echo "Report outline, work in progress." > "$W/reports/drafts/outline.txt"
      [ -f "$W/archive/2025/report-2025.txt" ] || echo "Last year's report." > "$W/archive/2025/report-2025.txt"
      [ -f "$W/data/imports/inventory.csv" ] || echo "id,name,quantity" > "$W/data/imports/inventory.csv"
      [ -f "$W/data/exports/inventory-export.csv" ] || echo "id,name,quantity" > "$W/data/exports/inventory-export.csv"
      [ -f "$W/notes/todo.txt" ] || echo "- review the final report" > "$W/notes/todo.txt"
      chown -R "$GYM_USER:$GYM_USER" "$W"
tasks:
  at_start:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/work/reports/drafts") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/work/reports/drafts. From anywhere, cd ~/work/reports/drafts gets you there."
    solve: |
      cd ~/work/reports/drafts
  looked:
    needs: [at_start]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" /etc >/dev/null || exit 1
      wait_exec '(^|/)ls( .*)?$'
    hint: |
      echo "Move to /etc, then run ls to see what is in it."
    solve: |
      cd /etc
      ls
  back:
    needs: [looked]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/reports/drafts" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        -) exit 0 ;;
        *) hint_exit "You are back in drafts, but you typed the path. The dash remembers it for you: go to /etc again and return with cd -." ;;
      esac
    hint: |
      echo "One command brings you back to the directory you were in before /etc: cd with a single dash."
    solve: |
      cd -
---

Here is a typical detour. You are working somewhere deep in the tree,
you need a quick look at another directory, and you want to be back
where you were right after.

Start in `~/work/reports/drafts`. Move to `/etc`, list what is there,
and return to `drafts` without typing its path.

::task{name="at_start"}
#active
Waiting for your shell to be in `~/work/reports/drafts`...
#completed
You are in `drafts`. Now take the detour to `/etc`.
::

::task{name="looked"}
#active
Waiting for you to arrive in `/etc` and list it...
#completed
Those are the system's configuration files. Now return to `drafts`.
::

::task{name="back"}
#active
Waiting for `cd -` to bring you back to `drafts`...
#completed
One dash, and the detour is over.
::

::tip
The dash only remembers one step back. After two moves, `cd -` returns
to the middle one, and the start is forgotten.
::
