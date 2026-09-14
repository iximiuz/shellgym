---
title: Long names
init:
  - name: build_work_tree
    run: |
      W="$GYM_USER_HOME/work"
      mkdir -p "$W/reports/drafts" "$W/reports/final" "$W/archive/2025" "$W/data/imports" "$W/data/exports" "$W/notes"
      mkdir -p "$W/Project Plans" "$W/2026 (draft)" "$W/quarterly-financial-statements-2026"
      [ -f "$W/reports/final/summary.txt" ] || echo "Final report, approved." > "$W/reports/final/summary.txt"
      [ -f "$W/reports/drafts/outline.txt" ] || echo "Report outline, work in progress." > "$W/reports/drafts/outline.txt"
      [ -f "$W/archive/2025/report-2025.txt" ] || echo "Last year's report." > "$W/archive/2025/report-2025.txt"
      [ -f "$W/data/imports/inventory.csv" ] || echo "id,name,quantity" > "$W/data/imports/inventory.csv"
      [ -f "$W/data/exports/inventory-export.csv" ] || echo "id,name,quantity" > "$W/data/exports/inventory-export.csv"
      [ -f "$W/notes/todo.txt" ] || echo "- review the final report" > "$W/notes/todo.txt"
      [ -f "$W/Project Plans/roadmap.txt" ] || echo "Roadmap for the next quarter." > "$W/Project Plans/roadmap.txt"
      [ -f "$W/2026 (draft)/budget.txt" ] || echo "Draft budget, numbers not final." > "$W/2026 (draft)/budget.txt"
      [ -f "$W/quarterly-financial-statements-2026/q1.txt" ] || echo "Q1 statement." > "$W/quarterly-financial-statements-2026/q1.txt"
      chown -R "$GYM_USER:$GYM_USER" "$W"
tasks:
  at_work:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/work") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/work. From anywhere, cd ~/work gets you there."
    solve: |
      cd ~/work
  arrived:
    needs: [at_work]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/quarterly-financial-statements-2026"
    hint: |
      echo "Type cd q and press Tab. The shell fills in the rest of the name."
    solve: |
      cd quarterly-financial-statements-2026
---

The `~/work` directory also holds one with a long name:
`quarterly-financial-statements-2026`. Typing it out invites a typo
somewhere in the middle.

Start in `~/work`, type `cd q`, and press **Tab**. The name is the only
one here that starts with `q`, so the shell completes it in full. Press
Enter to move in.

::task{name="at_work"}
#active
Waiting for your shell to be in `~/work`...
#completed
You are in `~/work`. Now let Tab type the long name.
::

::task{name="arrived"}
#active
Waiting for your shell to arrive in `quarterly-financial-statements-2026`...
#completed
That is one long name you never had to type in full.
::

::tip{title="When Tab stops early"}
If several names share the letters you typed, Tab completes only the
common part and stops. Press **Tab** twice to see the candidates, type
one more letter to tell them apart, and press Tab again.
::
