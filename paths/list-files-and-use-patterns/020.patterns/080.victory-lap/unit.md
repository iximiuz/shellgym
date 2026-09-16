---
title: Victory lap
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
  at_tmp:
    check: |
      P=$(wait_cwd "/tmp") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in /tmp. Move there with its full path."
    solve: |
      cd /tmp
  listed:
    needs: [at_tmp]
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/data"
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      wait_cwd --now "$SHELL_PID" /tmp >/dev/null || hint_exit "Your shell left /tmp. Go back to /tmp and list the files from there, with the directory's path in front of the pattern."
      set -f; set -- $ARGV; shift
      N=0
      for a in "$@"; do
        case "$a" in
          -*) continue ;;
          *\**|*\?*|*\[*) hint_exit "The pattern reached ls unexpanded, so it was quoted or it matched nothing. Type it bare, and check the path in front of it." ;;
        esac
        case "$a" in /*) P=$(realpath -m "$a") ;; *) P=$(realpath -m "/tmp/$a") ;; esac
        [ "$P" = "$D" ] && hint_exit "That listed the whole ~/projects/data directory. Add the quarterly pattern after the path, so that only the four files reach ls."
        [ "$(dirname "$P")" = "$D" ] || hint_exit "The last ls listed $a, which is not in ~/projects/data. Put the path of ~/projects/data in front of the pattern."
        case "${a##*/}" in
          sales-q[1-4].csv) N=$((N+1)) ;;
          *) hint_exit "The last ls listed more than the four quarterly files. Use the pattern that matches sales-q1.csv through sales-q4.csv and nothing else." ;;
        esac
      done
      [ "$N" -gt 0 ] || hint_exit "The last ls had no pattern. Give it the quarterly pattern with the path of ~/projects/data in front."
      LONG=' -[a-zA-Z0-9]*l[a-zA-Z0-9]* '
      LONGL=' --format=(long|verbose) '
      HUM=' -[a-zA-Z0-9]*h[a-zA-Z0-9]* '
      HUML=' --(human-readable|si) '
      SIZE=' -[a-zA-Z0-9]*S[a-zA-Z0-9]* '
      SIZEL=' --sort=size '
      if ! [[ " $ARGV " =~ $LONG ]] && ! [[ " $ARGV " =~ $LONGL ]]; then
        hint_exit "The pattern is right, but the short format shows no sizes at all. Add the long format with readable sizes, sorted by size."
      fi
      if ! [[ " $ARGV " =~ $HUM ]] && ! [[ " $ARGV " =~ $HUML ]]; then
        hint_exit "Those sizes are in bytes. Add the option that prints them in K, M, and G."
      fi
      if ! [[ " $ARGV " =~ $SIZE ]] && ! [[ " $ARGV " =~ $SIZEL ]]; then
        hint_exit "That listing is not sorted by size. Add the option that puts the largest file first."
      fi
    hint: |
      echo "One ls from /tmp: the long format, readable sizes, and the sort by size, followed by the quarterly pattern with ~/projects/data/ in front of it."
    solve: |
      ls -lhS ~/projects/data/sales-q?.csv
---

This is the last unit, and it is a lap through the whole gym. Do it
from `/tmp`, without moving your shell.

List the four quarterly files of `~/projects/data`, and only those, in
the long format with readable sizes, sorted by size with the largest
first. That is a single `ls` with three options and one pattern, and
the pattern needs the directory's path in front of it because your
shell is somewhere else.

::task{name="at_tmp"}
#active
Waiting for your shell to be in `/tmp`...
#completed
You are in `/tmp`. Stay here for the lap.
::

::task{name="listed"}
#active
Waiting for a sorted long listing of the quarterly files, from `/tmp`...
#completed
The top line is the largest quarterly file, and that is the whole gym.
::
