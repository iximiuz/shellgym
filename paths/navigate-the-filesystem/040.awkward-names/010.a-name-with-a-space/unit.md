---
title: A name with a space
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
  at_work:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/projects") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/projects. From anywhere, cd ~/projects gets you there."
    solve: |
      cd ~/projects
  tried:
    needs: [at_work]
    timeout: 60
    check: |
      wait_line '^cd +Vendor +Contracts *$'
    hint: |
      echo "Run cd Vendor Contracts as two bare words, without quotes. It is going to fail, and the error is the lesson."
    solve: |
      cd Vendor Contracts
  arrived:
    needs: [tried]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects/Vendor Contracts" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        *\"*|*\'*|*\\*) exit 0 ;;
        *) hint_exit "You got in, but the name was not protected. Go back to ~/projects and wrap the name in quotes so it reaches cd as one argument." ;;
      esac
    hint: |
      echo "Wrap the whole name in quotes, single or double, so the space stays inside one argument."
    solve: |
      cd "Vendor Contracts"
---

The `~/projects` directory now contains a directory named `Vendor Contracts`,
with a space in the middle. To the shell, a space separates arguments,
so `cd Vendor Contracts` hands `cd` two arguments, and `cd` only takes one.

Start in `~/projects` and try exactly that. Read the complaint. Then get
into the directory the way you learned in Meet the Linux Shell: wrap
the name in quotes so it stays one argument.

::task{name="at_work"}
#active
Waiting for your shell to be in `~/projects`...
#completed
You are in `~/projects`. Now try the bare name with the space.
::

::task{name="tried"}
#active
Waiting for the doomed `cd Vendor Contracts`...
#completed
The space split the name in two, so `cd` got too many arguments. Now quote it.
::

::task{name="arrived"}
#active
Waiting for a quoted name to get you into `Vendor Contracts`...
#completed
The quotes held the name together, space and all.
::
