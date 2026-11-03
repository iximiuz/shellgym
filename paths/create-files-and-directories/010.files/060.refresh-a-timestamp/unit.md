---
title: Refresh a timestamp
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
  - name: age_the_report
    run: |
      mkdir -p /run/shellgym-refs
      touch -d '3 days ago' "$GYM_USER_HOME/projects/reports/final/summary.txt"
      touch /run/shellgym-refs/refresh-a-timestamp
tasks:
  looked:
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/reports/final"
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      set -f; set -- $ARGV; shift
      CWD=$(shell_cwd 2>/dev/null || true)
      N=0; SEEN=0
      for a in "$@"; do
        case "$a" in
          -*) continue ;;
          /*) P=$(realpath -m "$a") ;;
          *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$a") ;;
        esac
        N=$((N+1)); T="$a"
        if [ "$P" = "$D" ] || [ "$P" = "$D/summary.txt" ]; then SEEN=1; fi
      done
      if [ "$N" -eq 0 ]; then
        [ "$CWD" = "$D" ] || hint_exit "Your shell is in ${CWD:-another directory}, so that listing does not show summary.txt. Move to ~/projects/reports/final first, or give ls the path ~/projects/reports/final as its argument."
      elif [ "$SEEN" -eq 0 ]; then
        hint_exit "That listed $T, so it does not show summary.txt. List ~/projects/reports/final instead."
      fi
      S=' -[a-zA-Z0-9]*l[a-zA-Z0-9]* '
      L=' --format=(long|verbose) '
      if [[ " $ARGV " =~ $S ]] || [[ " $ARGV " =~ $L ]]; then exit 0; fi
      hint_exit "That listing does not show dates. Use the long format."
    hint: |
      CWD=$(shell_cwd 2>/dev/null)
      if [ -n "$CWD" ] && [ "$CWD" != "$GYM_USER_HOME/projects/reports/final" ]; then
        echo "Your shell is in $CWD. Give ls the path ~/projects/reports/final as its argument, and use the long format to see the date of summary.txt."
        exit 0
      fi
      echo "List the contents of ~/projects/reports/final in the long format and read the date of summary.txt."
    solve: |
      ls -l ~/projects/reports/final
  touched:
    needs: [looked]
    timeout: 60
    check: |
      wait_file_newer "$GYM_USER_HOME/projects/reports/final/summary.txt" /run/shellgym-refs/refresh-a-timestamp
    hint: |
      echo "Run touch on the existing file: ~/projects/reports/final/summary.txt. It updates the time and keeps the contents."
    solve: |
      touch ~/projects/reports/final/summary.txt
---

The file `~/projects/reports/final/summary.txt` was last changed three
days ago. Check that first: list the contents of its directory in the
long format and
read the date.

Then run `touch` on the file. It already exists, so nothing is created.
Instead, its modification time becomes now, and its contents stay
exactly as they were. List the contents again to see the new date.

::task{name="looked"}
#active
Waiting for a long listing that shows the old date...
#completed
That date is three days old. Now refresh it.
::

::task{name="touched"}
#active
Waiting for the modification time of `summary.txt` to change...
#completed
The date is now, and the contents are unchanged.
::

::tip{title="Why refresh a timestamp?"}
Some tools decide what to do by comparing dates: a backup program copies
files newer than its last run, a build tool rebuilds what changed. A
`touch` makes a file look freshly changed without editing it.
::
