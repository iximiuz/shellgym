---
title: Refresh the oldest backup
vars:
  ORDER: { pick: ["web,db,mail,files", "db,files,web,mail", "mail,web,files,db", "files,mail,db,web"] }
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
  - name: age_the_backups
    run: |
      mkdir -p /run/shellgym-refs
      B="$GYM_USER_HOME/projects/backups"
      IFS=, read -r B1 B2 B3 B4 <<< "$ORDER"
      touch -d '1 hour ago' "$B/backup-$B1.bak"
      touch -d '2 days ago' "$B/backup-$B2.bak"
      touch -d '5 days ago' "$B/backup-$B3.bak"
      touch -d '9 days ago' "$B/backup-$B4.bak"
      touch /run/shellgym-refs/refresh-the-oldest
tasks:
  sorted:
    timeout: 60
    check: |
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      T=' -[a-zA-Z0-9]*t[a-zA-Z0-9]* '
      TL=' --sort=time '
      if [[ " $ARGV " =~ $T ]] || [[ " $ARGV " =~ $TL ]]; then exit 0; fi
      hint_exit "That listing is not sorted by time. Sort ~/projects/backups by modification time to see which backup is the oldest."
    hint: |
      echo "List ~/projects/backups sorted by modification time. Newest first means the oldest is on the last line."
    solve: |
      ls -lt ~/projects/backups
  refreshed:
    needs: [sorted]
    timeout: 60
    check: |
      B="$GYM_USER_HOME/projects/backups"
      M=/run/shellgym-refs/refresh-the-oldest
      OLDEST=${ORDER##*,}
      while :; do
        for b in web db mail files; do
          if [ "$b" != "$OLDEST" ] && wait_file_newer --now "$B/backup-$b.bak" "$M" >/dev/null 2>&1; then
            B="$GYM_USER_HOME/projects/backups"
            IFS=, read -r B1 B2 B3 B4 <<< "$ORDER"
            touch -d '1 hour ago' "$B/backup-$B1.bak"
            touch -d '2 days ago' "$B/backup-$B2.bak"
            touch -d '5 days ago' "$B/backup-$B3.bak"
            touch -d '9 days ago' "$B/backup-$B4.bak"
            hint_exit "backup-$b.bak was refreshed, but it was not the oldest one. The dates are back to what they were. Sort the listing by time again and read the last line."
          fi
        done
        wait_file_newer --now "$B/backup-$OLDEST.bak" "$M" >/dev/null 2>&1 && exit 0
        sleep 1
      done
    hint: |
      echo "Find the oldest backup in the time-sorted listing, then run touch on that one file."
    solve: |
      touch ~/projects/backups/"$(ls -t ~/projects/backups | tail -n 1)"
---

Four backups sit in `~/projects/backups`. Find the oldest one with a
listing sorted by modification time, then refresh its timestamp with
`touch` so that it becomes the newest.

Touch only that one file. A sorted listing afterwards should show it on
top.

::task{name="sorted"}
#active
Waiting for a listing of the backups sorted by time...
#completed
The oldest backup is the last line of that listing. Now refresh it.
::

::task{name="refreshed"}
#active
Waiting for the oldest backup to get a fresh timestamp...
#completed
The oldest backup is now the newest.
::
