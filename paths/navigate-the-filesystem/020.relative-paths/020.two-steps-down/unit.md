---
title: Two steps down
requires: [readline]
vars:
  TARGET: { pick: [reports/final, reports/drafts, data/imports, data/exports, archive/2025] }
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
  at_work:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/projects") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/projects. From anywhere, cd ~/projects gets you there."
    solve: |
      cd ~/projects
  arrived:
    needs: [at_work]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects/$TARGET" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        /*|"~"*|"\$"*) hint_exit "You got there with a full path. From ~/projects a relative path is enough. Go back and try the two names joined by a slash." ;;
        *) exit 0 ;;
      esac
    hint: |
      echo "Join the two names with a slash, exactly as written in the assignment, and give that to cd."
    solve: |
      cd $TARGET
---

A relative path can have several parts, just like an absolute one.
Written as `${TARGET}`, it means "into the first directory, then into
the second".

Start in `~/projects` and move into `${TARGET}` with a single command.

::task{name="at_work"}
#active
Waiting for your shell to be in `~/projects`...
#completed
You are in `~/projects`. Now go two levels down in one move.
::

::task{name="arrived"}
#active
Waiting for a relative `cd` into `${TARGET}`...
#completed
You went two directories down in one command.
::

::tip
Tab completion works one part at a time. Type the first letters of
the first name, press **Tab**, type a few letters of the second name,
and press **Tab** again.
::
