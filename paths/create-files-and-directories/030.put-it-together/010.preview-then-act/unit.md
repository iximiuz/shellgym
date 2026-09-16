---
title: Preview, then act
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
  - name: age_the_files
    run: |
      D="$GYM_USER_HOME/projects/data"
      mkdir -p /run/shellgym-refs
      for f in sales-q1.csv sales-q2.csv sales-q3.csv sales-q4.csv sales-2025.csv summary.txt summary.bak; do touch -d '3 days ago' "$D/$f"; done
      touch /run/shellgym-refs/preview-then-act
tasks:
  refreshed:
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/data"
      M=/run/shellgym-refs/preview-then-act
      CSVS="sales-q1.csv sales-q2.csv sales-q3.csv sales-q4.csv sales-2025.csv"
      OTHERS="summary.txt summary.bak"
      while :; do
        for f in $OTHERS; do
          if wait_file_newer --now "$D/$f" "$M" >/dev/null 2>&1; then
            for g in $CSVS $OTHERS; do touch -d '3 days ago' "$D/$g"; done
            hint_exit "$f was refreshed too, so the pattern matched more than the csv files. Every date is back to three days ago. Preview the pattern with echo, then touch again."
          fi
        done
        OK=1
        for f in $CSVS; do wait_file_newer --now "$D/$f" "$M" >/dev/null 2>&1 || OK=0; done
        [ "$OK" = 1 ] && exit 0
        sleep 1
      done
    hint: |
      echo "Preview the pattern with echo until it lists exactly the five csv files. Then give the same pattern to touch."
    solve: |
      echo ~/projects/data/*.csv
      touch ~/projects/data/*.csv
---

A pattern works with `touch` exactly as it does with `ls`: the shell
expands it into the matching names before the command runs, so
`touch` receives a list of files and refreshes every one of them.

Every file in `~/projects/data` is three days old. Refresh the modification
time of the five `.csv` files, and only those, with one `touch`.

Preview the pattern with `echo` first and check the list. If it names
exactly the five files, hand the same pattern to `touch`.

::task{name="refreshed"}
#active
Waiting for the five `.csv` files, and only those, to get a fresh timestamp...
#completed
Five files refreshed, and the other two kept their old date.
::

::tip
The habit is: type the pattern after `echo`, read the list, then press
the **Up** arrow, replace `echo` with the real command, and press Enter.
With `Ctrl-A` you jump to the start of the line to make that edit.
::
