---
title: Mark the largest file
vars:
  LARGEST: { pick: [q1, q2, q3, q4] }
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
  - name: size_the_files
    run: |
      D="$GYM_USER_HOME/projects/data"
      set -- 1200 24000 350000
      for q in q1 q2 q3 q4; do
        if [ "$q" = "$LARGEST" ]; then S=2100000; else S=$1; shift; fi
        yes 'id,name,quantity,location' | head -c "$S" > "$D/sales-$q.csv" || true
      done
      chown "$GYM_USER:$GYM_USER" "$D"/sales-q?.csv
      rm -f "$D"/*.reviewed
tasks:
  marked:
    timeout: 60
    check: |
      wait_file "$GYM_USER_HOME/projects/data/sales-$LARGEST.csv.reviewed"
    hint: |
      D="$GYM_USER_HOME/projects/data"
      for f in "$D"/*.reviewed; do
        [ -e "$f" ] || continue
        case "$f" in
          *"sales-$LARGEST.csv.reviewed") ;;
          *) echo "$(basename "$f") exists, but $(basename "${f%.reviewed}") is not the largest file. Sort the listing by size and read the first line."; exit 0 ;;
        esac
      done
      echo "Sort ~/projects/data by size, largest first. Then create an empty file named after the top entry, with .reviewed appended."
    solve: |
      ls -lS ~/projects/data
      touch ~/projects/data/sales-$LARGEST.csv.reviewed
---

Find the largest file in `~/projects/data` and mark it as reviewed: create
an empty file next to it with the same name plus `.reviewed`. For a
file called `report.txt`, the marker would be `report.txt.reviewed`.

Sort the listing by size to find the file. Readable sizes make the
answer obvious.

::task
#active
Waiting for a `.reviewed` marker next to the largest file...
#completed
The largest file is marked.
::
