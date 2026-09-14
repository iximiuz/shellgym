---
title: Next door
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
      P=$(wait_cwd "$GYM_USER_HOME/work/reports/final") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/work/reports/final. From anywhere, cd ~/work/reports/final gets you there."
    solve: |
      cd ~/work/reports/final
  arrived:
    needs: [at_start]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/reports/drafts" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        ../*) exit 0 ;;
        *) hint_exit "You made it next door, but not in one relative move through the parent. Go back to final and try a path that starts with the two-dot name." ;;
      esac
    hint: |
      echo "Go up to the parent and down into the sibling in one path: the two-dot name, a slash, and the sibling's name."
    solve: |
      cd ../drafts
---

The `reports` directory holds two directories, `drafts` and `final`.
They are siblings: both have `reports` as their parent.

Start in `~/work/reports/final`. Then move to `drafts` with a
**single** relative path that goes up to the parent and down into the
sibling.

::task{name="at_start"}
#active
Waiting for your shell to be in `~/work/reports/final`...
#completed
You are in `final`. Now cross over to `drafts` in one move.
::

::task{name="arrived"}
#active
Waiting for a `cd` through `..` into `drafts`...
#completed
One path took you up to the parent and down into the sibling.
::
