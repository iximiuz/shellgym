---
title: Two directories at once
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
tasks:
  listed:
    timeout: 60
    check: |
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)+$') || exit 1
      T=0; H=0
      case " $ARGV " in *" /tmp "*|*" /tmp/ "*) T=1 ;; esac
      case " $ARGV " in *" $GYM_USER_HOME "*|*" $GYM_USER_HOME/ "*) H=1 ;; esac
      [ "$T$H" = 11 ] && exit 0
      hint_exit "That listing did not cover both directories. Give ls the two paths in a single command: /tmp and your home."
    hint: |
      echo "One ls, two arguments: /tmp and your home directory (the tilde works)."
    solve: |
      ls /tmp ~
---

The `ls` command accepts any number of paths. With more than one, it
prints each directory in turn, under a header line with the
directory's name.

List the contents of `/tmp` and of your home directory with a single
command.

::task
#active
Waiting for one `ls` that covers both `/tmp` and your home...
#completed
Both directories came out of one command, each under its own header.
::
