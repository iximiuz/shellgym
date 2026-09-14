---
title: Find the file
vars:
  NAME: { pick: [budget-2026.txt, checklist.txt, invoice-1042.txt, access-log.txt] }
  HIDE: { pick: [reports/drafts, reports/final, archive/2025, data/imports, data/exports, notes] }
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
  - name: hide_the_file
    run: |
      W="$GYM_USER_HOME/work"
      rm -f "$W"/*/"$NAME" "$W"/*/*/"$NAME"
      echo "You found it." > "$W/$HIDE/$NAME"
      chown "$GYM_USER:$GYM_USER" "$W/$HIDE/$NAME"
tasks:
  at_work:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/work") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/work. From anywhere, cd ~/work gets you there."
    solve: |
      cd ~/work
  found:
    needs: [at_work]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/$HIDE"
    hint: |
      echo "Look into each directory with ls until ${NAME} shows up, then move into that directory."
    solve: |
      cd ~/work/$HIDE
---

A file named `${NAME}` is somewhere in the `~/work` tree, and nobody
remembers where. Every directory in the tree is a candidate.

Start in `~/work`. Look into the directories until you spot the file,
then move your shell into the directory that holds it.

::task{name="at_work"}
#active
Waiting for your shell to be in `~/work`...
#completed
You are in `~/work`. Now find the file.
::

::task{name="found"}
#active
Waiting for your shell to arrive in the directory that holds `${NAME}`...
#completed
The file was hiding in `${HIDE}`.
::

::tip
The `ls` command accepts a path, so `ls reports/final` shows the
contents of that directory without moving there. It saves a lot of
back and forth while you search.
::
