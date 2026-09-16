---
title: Sizes you can read
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
tasks:
  listed:
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
      LONG=' -[a-zA-Z0-9]*[ls][a-zA-Z0-9]* '
      LONGL=' --(format=(long|verbose)|size) '
      HUM=' -[a-zA-Z0-9]*h[a-zA-Z0-9]* '
      HUML=' --(human-readable|si) '
      if [[ " $ARGV " =~ $LONG ]] || [[ " $ARGV " =~ $LONGL ]]; then
        if [[ " $ARGV " =~ $HUM ]] || [[ " $ARGV " =~ $HUML ]]; then exit 0; fi
        hint_exit "Those sizes are in bytes. Add the option that prints them in K, M, and G. Look for 'human' in ls --help."
      fi
      hint_exit "The short format shows no sizes at all. Use the long format, plus the option that prints sizes in K, M, and G."
    hint: |
      CWD=$(shell_cwd 2>/dev/null)
      if [ -n "$CWD" ] && [ "$CWD" != "$GYM_USER_HOME/projects/data" ]; then
        echo "Your shell is in $CWD. Move to ~/projects/data first, or give ls the path ~/projects/data as its argument."
        exit 0
      fi
      echo "Combine the long format with the human-readable option. Look for 'human' in ls --help."
    solve: |
      cd ~/projects/data
      ls -lh
---

The long format prints sizes in bytes, and a number like `2100000` is hard to read.
The `ls` command has an option to round the sizes to `K`, `M`, and `G`.

List the contents of `~/projects/data` with readable sizes. Move there
first and list the working directory, or stay where you are and give
`ls` the path. Then find the largest file by glancing at the output.

::task{name="listed"}
#active
Waiting for a long listing with human-readable sizes...
#completed
The largest file is the one with an `M` next to it.
::

::tip
Cannot recall the option? `ls --help` is your friend.
::
