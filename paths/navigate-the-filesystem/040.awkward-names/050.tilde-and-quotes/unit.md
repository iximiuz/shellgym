---
title: The tilde and the quotes
requires: [readline]
init:
  - name: build_work_tree
    run: |
      W="$GYM_USER_HOME/projects"
      mkdir -p "$W/reports/drafts" "$W/reports/final" "$W/archive/2025" "$W/data/imports" "$W/data/exports" "$W/notes"
      mkdir -p "$W/Vendor Contracts" "$W/2026 (draft)" "$W/quarterly-financial-statements-2026"
      [ -f "$W/reports/final/summary.txt" ] || echo "Final report, approved." > "$W/reports/final/summary.txt"
      [ -f "$W/reports/drafts/outline.txt" ] || echo "Report outline, work in progress." > "$W/reports/drafts/outline.txt"
      [ -f "$W/archive/2025/report-2025.txt" ] || echo "Last year's report." > "$W/archive/2025/report-2025.txt"
      [ -f "$W/data/imports/inventory.csv" ] || echo "id,name,quantity" > "$W/data/imports/inventory.csv"
      [ -f "$W/data/exports/inventory-export.csv" ] || echo "id,name,quantity" > "$W/data/exports/inventory-export.csv"
      [ -f "$W/notes/todo.txt" ] || echo "- review the final report" > "$W/notes/todo.txt"
      [ -f "$W/Vendor Contracts/renewals.txt" ] || echo "Contracts up for renewal this year." > "$W/Vendor Contracts/renewals.txt"
      [ -f "$W/2026 (draft)/budget.txt" ] || echo "Draft budget, numbers not final." > "$W/2026 (draft)/budget.txt"
      [ -f "$W/quarterly-financial-statements-2026/q1.txt" ] || echo "Q1 statement." > "$W/quarterly-financial-statements-2026/q1.txt"
      chown -R "$GYM_USER:$GYM_USER" "$W"
tasks:
  away:
    check: |
      P=$(wait_cwd /tmp) || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Move to /tmp with its full path."
    solve: |
      cd /tmp
  arrived:
    needs: [away]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects/Vendor Contracts" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        "~"*) exit 0 ;;
        *) hint_exit "You got there, but not in one command that starts with the tilde. Go back to /tmp and write the path from ~ down, protecting only the part with the space." ;;
      esac
    hint: |
      echo "Start the path with the bare tilde and quote only the last part: ~/projects/ followed by the quoted name. Or type ~/projects/Vend and press Tab."
    solve: |
      cd ~/projects/"Vendor Contracts"
---

Combine two things you know. From `/tmp`, get into `~/projects/Vendor
Contracts` with a single command that starts with the tilde.

There is a trap. Quotes turn off the tilde: `"~/projects/Vendor Contracts"`
is taken literally, tilde and all, and no directory by that name
exists. Keep the tilde outside the quotes and protect only the part
that needs it.

::task{name="away"}
#active
Waiting for your shell to arrive in `/tmp`...
#completed
You are in `/tmp`. Now jump into the spaced directory in one move.
::

::task{name="arrived"}
#active
Waiting for a single `cd` from `~` into `Vendor Contracts`...
#completed
The tilde expanded, the quotes held the space, and `cd` got one clean path.
::

::hint{title="No such file or directory: ~/projects/Vendor Contracts"}
The message shows the tilde unexpanded. It was inside the quotes.
Move the opening quote to the right, just before the awkward name,
or drop the quotes and let Tab escape the space.
::
