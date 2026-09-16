---
title: Newest first
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
      B="$GYM_USER_HOME/projects/backups"
      IFS=, read -r B1 B2 B3 B4 <<< "$ORDER"
      touch -d '1 hour ago' "$B/backup-$B1.bak"
      touch -d '2 days ago' "$B/backup-$B2.bak"
      touch -d '5 days ago' "$B/backup-$B3.bak"
      touch -d '9 days ago' "$B/backup-$B4.bak"
tasks:
  newest:
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/backups"
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      set -f; set -- $ARGV; shift
      CWD=$(shell_cwd 2>/dev/null || true)
      N=0
      for a in "$@"; do
        case "$a" in -*) ;; *) N=$((N+1)); T="$a" ;; esac
      done
      if [ "$N" -eq 0 ]; then
        [ "$CWD" = "$D" ] || hint_exit "Your shell is in ${CWD:-another directory}, so that listed the wrong directory. Move to ~/projects/backups first, or give ls the path ~/projects/backups as its argument."
      elif [ "$N" -eq 1 ]; then
        case "$T" in /*) P=$(realpath -m "$T") ;; *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$T") ;; esac
        [ "$P" = "$D" ] || hint_exit "That listed $T instead of ~/projects/backups."
      else
        hint_exit "That listed several paths at once. List only ~/projects/backups."
      fi
      T=' -[a-zA-Z0-9]*t[a-zA-Z0-9]* '
      TL=' --sort=time '
      if [[ " $ARGV " =~ $T ]] || [[ " $ARGV " =~ $TL ]]; then exit 0; fi
      if [ "${GYM_CHECK_ATTEMPT:-1}" -gt 2 ]; then
        hint_exit "That listing is still sorted by name. The option that sorts by modification time is -t, so the command is ls -lt."
      fi
      hint_exit "That listing is sorted by name. Add the option that sorts by modification time. Look for 'time' in ls --help."
    hint: |
      CWD=$(shell_cwd 2>/dev/null)
      if [ -n "$CWD" ] && [ "$CWD" != "$GYM_USER_HOME/projects/backups" ]; then
        echo "Your shell is in $CWD. Move to ~/projects/backups first, or give ls the path ~/projects/backups as its argument."
        exit 0
      fi
      if [ "${GYM_CHECK_ATTEMPT:-1}" -gt 2 ]; then
        echo "The option that sorts by modification time is -t. Run ls -lt."
      else
        echo "Look for 'sort by' and 'time' in ls --help. Combine that option with the long format so you can see the dates."
      fi
    solve: |
      ls -lt ~/projects/backups
  oldest:
    needs: [newest]
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/backups"
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      set -f; set -- $ARGV; shift
      CWD=$(shell_cwd 2>/dev/null || true)
      N=0
      for a in "$@"; do
        case "$a" in -*) ;; *) N=$((N+1)); T="$a" ;; esac
      done
      if [ "$N" -eq 0 ]; then
        [ "$CWD" = "$D" ] || hint_exit "Your shell is in ${CWD:-another directory}, so that listed the wrong directory. Move to ~/projects/backups first, or give ls the path ~/projects/backups as its argument."
      elif [ "$N" -eq 1 ]; then
        case "$T" in /*) P=$(realpath -m "$T") ;; *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$T") ;; esac
        [ "$P" = "$D" ] || hint_exit "That listed $T instead of ~/projects/backups."
      else
        hint_exit "That listed several paths at once. List only ~/projects/backups."
      fi
      T=' -[a-zA-Z0-9]*t[a-zA-Z0-9]* '
      TL=' --sort=time '
      R=' -[a-zA-Z0-9]*r[a-zA-Z0-9]* '
      RL=' --reverse '
      if [[ " $ARGV " =~ $T ]] || [[ " $ARGV " =~ $TL ]]; then
        if [[ " $ARGV " =~ $R ]] || [[ " $ARGV " =~ $RL ]]; then exit 0; fi
        hint_exit "That is newest first again. Add the option that reverses the order."
      fi
      hint_exit "Keep the sort by time and add the option that reverses the order."
    hint: |
      CWD=$(shell_cwd 2>/dev/null)
      if [ -n "$CWD" ] && [ "$CWD" != "$GYM_USER_HOME/projects/backups" ]; then
        echo "Your shell is in $CWD. Move to ~/projects/backups first, or give ls the path ~/projects/backups as its argument."
        exit 0
      fi
      echo "Keep the time sort and add the reverse option. Look for 'reverse' in ls --help."
    solve: |
      ls -ltr ~/projects/backups
---

Four backup files sit in `~/projects/backups`, and their names say nothing
about their age. The long format shows each file's date, but reading
dates and comparing them in your head can be slow.
Sorting does the comparison for you.

List the contents of `~/projects/backups` sorted by modification time,
newest first. Then reverse the order so the oldest backup comes first.

::task{name="newest"}
#active
Waiting for a listing sorted by time...
#completed
The top line is the newest backup. Now flip the order.
::

::task{name="oldest"}
#active
Waiting for the same listing in reverse...
#completed
Now the oldest backup is on top.
::

::tip
Sorting newest first with the long format is another combination worth
remembering. Reversed, it puts the most recent change at the bottom,
right above your prompt, which is handy in a directory with hundreds of
files.
::
