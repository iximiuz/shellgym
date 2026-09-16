---
title: Victory lap
vars:
  YEAR: { pick: ["2023", "2024", "2026"] }
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
  at_tmp:
    check: |
      P=$(wait_cwd "/tmp") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in /tmp. Move there with its full path."
    solve: |
      cd /tmp
  directories:
    needs: [at_tmp]
    timeout: 60
    check: |
      wait_dir "$GYM_USER_HOME/projects/archive/$YEAR/reports" >/dev/null || exit 1
      wait_cwd --now "$SHELL_PID" /tmp >/dev/null || hint_exit "The directories exist, but your shell left /tmp. Go back to /tmp and do the rest from there, with full paths."
    hint: |
      echo "From /tmp, one mkdir with the parents option and the full path ~/projects/archive/${YEAR}/reports."
    solve: |
      mkdir -p ~/projects/archive/$YEAR/reports
  file:
    needs: [directories]
    timeout: 60
    check: |
      wait_file "$GYM_USER_HOME/projects/archive/$YEAR/reports/Quarterly summary.txt" >/dev/null || exit 1
      wait_cwd --now "$SHELL_PID" /tmp >/dev/null || hint_exit "The file exists, but your shell left /tmp. Go back to /tmp and do the rest from there, with full paths."
    hint: |
      R="$GYM_USER_HOME/projects/archive/$YEAR/reports"
      if [ -e "$R/Quarterly" ] || [ -e "$R/summary.txt" ]; then
        echo "The space split the name in two. Quote the file name so it stays one argument, and keep the tilde outside the quotes."
      else
        echo "From /tmp, touch the file by its full path: ~/projects/archive/${YEAR}/reports/ followed by the quoted name."
      fi
    solve: |
      touch ~/projects/archive/$YEAR/reports/"Quarterly summary.txt"
  listed:
    needs: [file]
    timeout: 60
    check: |
      ARGV=$(wait_exec --cwd /tmp --latest "(^|/)ls( +\S+)* +\S*archive/${YEAR}/reports/? *\$") || exit 1
      S=' -[a-zA-Z0-9]*l[a-zA-Z0-9]* '
      L=' --format=(long|verbose) '
      if [[ " $ARGV " =~ $S ]] || [[ " $ARGV " =~ $L ]]; then exit 0; fi
      hint_exit "That listing shows names only. Use the long format."
    hint: |
      echo "List the contents of ~/projects/archive/${YEAR}/reports in the long format, from /tmp."
    solve: |
      ls -l ~/projects/archive/$YEAR/reports
---

This is the last unit, and it is a lap through the whole gym. Do
everything from `/tmp`, without moving your shell:

1. Create `~/projects/archive/${YEAR}/reports`. The `${YEAR}` directory
   does not exist yet.
2. Inside it, create an empty file named `Quarterly summary.txt`, with
   the space.
3. List the contents of that directory in the long format to see the
   new file.

::task{name="at_tmp"}
#active
Waiting for your shell to be in `/tmp`...
#completed
You are in `/tmp`. Stay here for the rest of the lap.
::

::task{name="directories"}
#active
Waiting for `~/projects/archive/${YEAR}/reports` to appear...
#completed
That is the nested directory done.
::

::task{name="file"}
#active
Waiting for `Quarterly summary.txt` to appear in `reports`...
#completed
That is the file with a space done.
::

::task{name="listed"}
#active
Waiting for the contents of the new directory in the long format...
#completed
An empty file with a fresh date, and that is the whole gym.
::
