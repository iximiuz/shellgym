---
title: One of a set
requires: [readline]
vars:
  LO: { pick: ["1", "2"] }
  HI: { pick: ["3", "4"] }
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
      LINE=$(wait_line --latest '^(ls|echo)( +-\S+)* +(\S*/)?sales-q\S*[]*?[]\S* *$') || exit 1
      CWD=$(shell_cwd 2>/dev/null || true)
      ELSEWHERE="Your shell is in ${CWD:-another directory}, so the pattern was expanded there. Move to ~/projects/data first, or put the directory's path in front of the pattern."
      case "$LINE" in
        *\"*|*\'*) hint_exit "Quotes switch pattern matching off. Type the pattern bare, without quotes." ;;
        *"[${LO}${HI}]"*|*"[${HI}${LO}]"*) ;;
        *\[*\]*) hint_exit "The brackets are right, but the set inside them is not ${LO} and ${HI}. List exactly those two digits between the brackets." ;;
        *) hint_exit "A star or a question mark matches too much here. Put the wanted digits between square brackets, right where the digit goes." ;;
      esac
      case "$LINE" in
        echo*)
          case "$LINE" in *"projects/data/"*) exit 0 ;; esac
          [ "$CWD" = "$D" ] || hint_exit "$ELSEWHERE"
          exit 0 ;;
      esac
      # A builtin echo leaves no exec event, so only an ls line can be placed:
      # its expanded arguments tell where the pattern was applied and
      # whether it matched what it should.
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
          sales-q${LO}.csv|sales-q${HI}.csv) N=$((N+1)) ;;
          *) hint_exit "The last ls listed more than sales-q${LO}.csv and sales-q${HI}.csv. Run it with the bracket pattern." ;;
        esac
      done
      [ "$N" -gt 0 ] || hint_exit "The last ls had no pattern. Run it with the bracket pattern."
    hint: |
      echo "Square brackets hold a set of characters, and the pattern matches one of them: sales-q, then the two digits inside brackets, then .csv."
    solve: |
      ls ~/projects/data/sales-q[${LO}${HI}].csv
---

Square brackets match one character out of a set: `[13]` matches a
`1` or a `3`, and nothing else. The set goes exactly where that one
character sits in the name.

List `sales-q${LO}.csv` and `sales-q${HI}.csv` in `~/projects/data` with
a single pattern that matches no other file.

::task{name="matched"}
#active
Waiting for a pattern with a character set...
#completed
Exactly two names matched.
::

::note
A dash inside the brackets makes a range: `[1-4]` matches any digit
from 1 to 4, and `[a-z]` any lowercase letter.
::
