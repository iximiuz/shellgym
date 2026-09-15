---
title: Climb home
requires: [readline]
vars:
  FROM: { pick: [archive/2025, data/exports, reports/drafts, reports/final] }
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
  at_start:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/projects/$FROM") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/projects/${FROM}. From anywhere, cd ~/projects/${FROM} gets you there."
    solve: |
      cd ~/projects/$FROM
  arrived:
    needs: [at_start]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        ../../..|../../../) exit 0 ;;
        *) hint_exit "You are home, but not by climbing. Go back down to ~/projects/${FROM} and count the levels: each two-dot name takes you up one." ;;
      esac
    hint: |
      echo "Count how many directories lie between you and home, and chain that many two-dot names with slashes."
    solve: |
      cd ../../..
---

You know two ways home already: a bare `cd` and `cd ~`. Here is the
long way, for practice: climb there with `..` only.

Start in `~/projects/${FROM}`. Count how many levels separate you from
your home directory, and climb up all of them in one relative path.

::task{name="at_start"}
#active
Waiting for your shell to be in `~/projects/${FROM}`...
#completed
You are in `${FROM}`. Now climb home level by level, in one path.
::

::task{name="arrived"}
#active
Waiting for a `cd` made only of two-dot names to reach home...
#completed
Three levels up, and you are home.
::
