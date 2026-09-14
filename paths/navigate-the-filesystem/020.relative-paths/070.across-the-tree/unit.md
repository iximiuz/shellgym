---
title: Across the tree
requires: [readline]
vars:
  FROM: { pick: [data/imports, data/exports, reports/final, reports/drafts] }
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
      P=$(wait_cwd "$GYM_USER_HOME/work/$FROM") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/work/${FROM}. From anywhere, cd ~/work/${FROM} gets you there."
    solve: |
      cd ~/work/$FROM
  arrived:
    needs: [at_start]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/archive/2025" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        ../../*) exit 0 ;;
        *) hint_exit "You reached archive/2025, but not in one relative move. Go back and write the whole route in one path: up twice, then down twice." ;;
      esac
    hint: |
      echo "Climb two levels with two two-dot names joined by a slash, then continue down into archive and 2025, all in one path."
    solve: |
      cd ../../archive/2025
---

The two-dot name can be chained. The path `../..` means "the parent of
the parent", and you can keep going down from there.

Start in `~/work/${FROM}`. Then move to `~/work/archive/2025` with a
single relative path: two levels up, then two levels down.

::task{name="at_start"}
#active
Waiting for your shell to be in `~/work/${FROM}`...
#completed
You are in `${FROM}`. Now cross to `archive/2025` in one move.
::

::task{name="arrived"}
#active
Waiting for a single `cd` through `../..` into `archive/2025`...
#completed
That was the whole route across the tree in one path.
::

::tip
Tab completion understands `..` too. After typing `../../a`, one
**Tab** fills in `archive/`.
::
