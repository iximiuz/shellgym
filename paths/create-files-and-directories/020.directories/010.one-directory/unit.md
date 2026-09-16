---
title: One directory
vars:
  DIR: { pick: [staging, incoming, outgoing, scratch] }
init:
  - name: build_work_tree
    run: |
      W="$GYM_USER_HOME/projects"
      mkdir -p "$W/reports/drafts" "$W/reports/final" "$W/archive/2025" "$W/data/imports" "$W/data/exports" "$W/notes" "$W/backups"
      [ -f "$W/reports/final/summary.txt" ] || echo "Final report, approved." > "$W/reports/final/summary.txt"
      [ -f "$W/reports/drafts/outline.txt" ] || echo "Report outline, work in progress." > "$W/reports/drafts/outline.txt"
      [ -f "$W/archive/2025/report-2025.txt" ] || echo "Last year's report." > "$W/archive/2025/report-2025.txt"
      [ -f "$W/data/imports/inventory.csv" ] || echo "id,name,quantity" > "$W/data/imports/inventory.csv"
      [ -f "$W/data/exports/inventory-export.csv" ] || echo "id,name,quantity" > "$W/data/exports/inventory-export.csv"
      [ -f "$W/notes/todo.txt" ] || echo "- review the final report" > "$W/notes/todo.txt"
      set -- 1200 24000 2100000 350000
      for q in q1 q2 q3 q4; do [ -f "$W/data/sales-$q.csv" ] || { yes 'id,name,quantity,location' | head -c "$1" > "$W/data/sales-$q.csv" || true; }; shift; done
      [ -f "$W/data/sales-2025.csv" ] || { yes 'id,name,quantity,location' | head -c 96000 > "$W/data/sales-2025.csv" || true; }
      [ -f "$W/data/summary.txt" ] || echo "Sales summary, all quarters." > "$W/data/summary.txt"
      [ -f "$W/data/summary.bak" ] || echo "Sales summary, all quarters." > "$W/data/summary.bak"
      [ -f "$W/data/.last-import" ] || date > "$W/data/.last-import"
      for b in web db mail files; do [ -f "$W/backups/backup-$b.bak" ] || echo "Backup of the $b settings." > "$W/backups/backup-$b.bak"; done
      chown -R "$GYM_USER:$GYM_USER" "$W"
  - name: clear_target
    run: |
      rm -rf "$GYM_USER_HOME/projects/$DIR"
tasks:
  created:
    check: |
      wait_dir "$GYM_USER_HOME/projects/$DIR"
    hint: |
      if [ -f "$GYM_USER_HOME/projects/$DIR" ]; then
        echo "${DIR} exists, but it is a regular file. The task needs a directory, which is what mkdir creates."
      else
        echo "Give mkdir the path of the directory to create: ~/projects/${DIR}."
      fi
    solve: |
      mkdir ~/projects/$DIR
---

Create a directory named `${DIR}` inside `~/projects`. The `mkdir` command
takes the path of the new directory. Absolute or relative, it is your
choice.

::task
#active
Waiting for the directory `~/projects/${DIR}` to appear...
#completed
The directory is there, empty and ready.
::

::tip
A long listing of `~/projects` shows the new entry with `d` as its first
character. That is how you tell it from a file called `${DIR}`.
::
