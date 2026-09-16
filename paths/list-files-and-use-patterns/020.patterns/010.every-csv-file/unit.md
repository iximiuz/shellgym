---
title: Every csv file
requires: [readline]
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
  matched:
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/data"
      LINE=$(wait_line --latest '^ls( +-\S+)* +\S*\*\S*csv *$') || exit 1
      CWD=$(shell_cwd 2>/dev/null || true)
      ELSEWHERE="Your shell is in ${CWD:-another directory}, so the pattern was expanded there. Move to ~/projects/data first, or put the directory's path in front of the pattern."
      case "$LINE" in
        *\"*|*\'*) hint_exit "Quotes switch pattern matching off, so the command looked for a file literally named that. Type the pattern bare, without quotes." ;;
      esac
      # The expanded arguments of the ls tell where the pattern was applied
      # and whether it matched what it should.
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      set -f; set -- $ARGV; shift
      N=0
      for a in "$@"; do
        case "$a" in
          -*) continue ;;
          */*) case "$a" in /*) P=$(realpath -m "$a") ;; *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$a") ;; esac
               [ "$(dirname "$P")" = "$D" ] || hint_exit "The last ls listed $a, which is outside ~/projects/data. Point the pattern at ~/projects/data." ;;
          *) [ "$CWD" = "$D" ] || hint_exit "$ELSEWHERE" ;;
        esac
        case "${a##*/}" in
          *.csv) N=$((N+1)) ;;
          *) hint_exit "The last ls listed more than the .csv files. Run it with the *.csv pattern." ;;
        esac
      done
      [ "$N" -gt 0 ] || hint_exit "The last ls had no pattern. Run it with *.csv."
    hint: |
      echo "The pattern is a star followed by .csv, and it goes where a file name would go: after ls."
    solve: |
      cd ~/projects/data
      ls *.csv
---

The `*` matches any sequence of characters, so `*.csv` means "any name that
ends in `.csv`". Written as an argument, the shell replaces it with
every matching name in the working directory before the command runs.

List only the `.csv` files in `~/projects/data` by giving `ls` the
pattern. A pattern is matched in the working directory, so move there
first, or put the directory's path in front of the pattern.

::task{name="matched"}
#active
Waiting for an `ls` with a pattern that matches the `.csv` files...
#completed
The shell turned the pattern into five names before the command ran.
::

::tip
You can see what a pattern expands to before handing it to a command.
Run `echo *.csv` in the same directory: `echo` prints its arguments, so
it shows exactly the names the shell would pass to `ls`.
::
