---
title: Tab does the escaping
requires: [readline]
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
  at_home:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start from your home directory. A bare cd takes you there."
    solve: |
      cd
  arrived:
    needs: [at_home]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/Project Plans" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        *\\\ *) exit 0 ;;
        *\"*|*\'*) hint_exit "Quotes work too, and this rep is about the backslash. Go back home, type cd work/Proj, and press Tab to see the shell write the escaped name for you." ;;
        *) hint_exit "You got in, but the name was not protected. Go back home and try again with a backslash in front of the space, or let Tab type it." ;;
      esac
    hint: |
      echo "Put a backslash right before the space, or type cd work/Proj and press Tab to let the shell do it."
    solve: |
      cd work/Project\ Plans
---

Quotes protect a whole word. A **backslash** protects a single
character: `Project\ Plans` tells the shell that this one space is
part of the name.

You rarely need to type that by hand. Start in your home directory, type
`cd work/Proj`, and press **Tab**. The shell completes the name and
inserts the backslash for you. Press Enter to move in.

::task{name="at_home"}
#active
Waiting for your shell to be in your home directory...
#completed
You are home. Now let Tab finish the awkward name.
::

::task{name="arrived"}
#active
Waiting for a backslash-escaped `cd` into `Project Plans`...
#completed
The backslash protected the space, and Tab typed it for you.
::

::tip
Tab completion escapes every awkward character it meets, spaces
included. Whenever a name looks odd, type the first few letters and
let Tab spell the rest.
::
