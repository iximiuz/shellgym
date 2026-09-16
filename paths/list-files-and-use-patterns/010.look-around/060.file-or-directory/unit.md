---
title: File or directory?
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
  marked:
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/data"
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      set -f; set -- $ARGV; shift
      CWD=$(shell_cwd 2>/dev/null || true)
      N=0
      for a in "$@"; do
        case "$a" in -*) ;; *) N=$((N+1)); T="$a" ;; esac
      done
      if [ "$N" -eq 0 ]; then
        [ "$CWD" = "$D" ] || hint_exit "Your shell is in ${CWD:-another directory}, so that listed the wrong directory. Move to ~/projects/data first, or give ls the path ~/projects/data as its argument."
      elif [ "$N" -eq 1 ]; then
        case "$T" in /*) P=$(realpath -m "$T") ;; *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$T") ;; esac
        [ "$P" = "$D" ] || hint_exit "That listed $T instead of ~/projects/data."
      else
        hint_exit "That listed several paths at once. List only ~/projects/data."
      fi
      S=' -[a-zA-Z0-9]*[Fp][a-zA-Z0-9]* '
      L=' --(classify|file-type|indicator-style=[a-z-]+) '
      LONG=' -[a-zA-Z0-9]*l[a-zA-Z0-9]* '
      if [[ " $ARGV " =~ $S ]] || [[ " $ARGV " =~ $L ]]; then exit 0; fi
      if [[ " $ARGV " =~ $LONG ]]; then
        hint_exit "The long format shows the type too, in its first character. There is a shorter way: an option that appends a slash to every directory name. Look for 'classify' in ls --help."
      fi
      hint_exit "In that listing a file and a directory look the same. Look for 'classify' in ls --help and add that option."
    hint: |
      CWD=$(shell_cwd 2>/dev/null)
      if [ -n "$CWD" ] && [ "$CWD" != "$GYM_USER_HOME/projects/data" ]; then
        echo "Your shell is in $CWD. Move to ~/projects/data first, or give ls the path ~/projects/data as its argument."
        exit 0
      fi
      echo "Look for 'classify' in ls --help. The option appends / to directory names."
    solve: |
      cd ~/projects/data
      ls -F
---

In a plain listing a file and a directory look alike. The long format
tells them apart, but it is a lot of output when all you want is the
type. There is a lighter option: it keeps the short format and appends
a `/` to every directory name.

Find the option in `ls --help` (search for "classify") and list the
contents of `~/projects/data` with it.

::task{name="marked"}
#active
Waiting for a listing that marks directories with a slash...
#completed
Two names end in a slash. Those are the directories.
::
