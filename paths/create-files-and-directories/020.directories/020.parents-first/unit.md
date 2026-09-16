---
title: Parents first
vars:
  YEAR: { pick: ["2023", "2024", "2026"] }
  Q: { pick: [q1, q2, q3, q4] }
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
      rm -rf "$GYM_USER_HOME/projects/archive/$YEAR"
tasks:
  created:
    check: |
      wait_dir "$GYM_USER_HOME/projects/archive/$YEAR/$Q"
    hint: |
      echo "mkdir refuses to create ${Q} while ${YEAR} is missing. Look for 'parents' in mkdir --help: one option creates every missing directory on the way."
    solve: |
      mkdir -p ~/projects/archive/$YEAR/$Q
---

Create the directory `~/projects/archive/${YEAR}/${Q}`. There is a catch:
`${YEAR}` does not exist yet, and by default `mkdir` refuses to create
a directory inside a parent that is missing. Try it and read the error.

Then look in `mkdir --help` for the option that creates the missing
parents along the way, and create both levels with one command.

::task
#active
Waiting for `~/projects/archive/${YEAR}/${Q}` to appear...
#completed
Both levels were created in one go.
::

::tip
That option is also safe to use when the directory already exists. It
turns "create this path" into "make sure this path exists", which is
exactly what most scripts want.
::
