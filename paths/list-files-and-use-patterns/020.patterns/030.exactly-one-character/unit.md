---
title: Exactly one character
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
      LINE=$(wait_line --latest '^(ls|echo)( +-\S+)* +(\S*/)?sales-\S*[*?]\S* *$') || exit 1
      CWD=$(shell_cwd 2>/dev/null || true)
      ELSEWHERE="Your shell is in ${CWD:-another directory}, so the pattern was expanded there. Move to ~/projects/data first, or put the directory's path in front of the pattern."
      case "$LINE" in
        *\"*|*\'*) hint_exit "Quotes switch pattern matching off. Type the pattern bare, without quotes." ;;
        *\?*) ;;
        *) hint_exit "A star matches any number of characters, so that pattern also picks up sales-2025.csv. Use the wildcard that stands for exactly one character." ;;
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
          sales-q[1-4].csv) N=$((N+1)) ;;
          *) hint_exit "The last ls listed more than the four quarterly files. Run it with the quarterly pattern." ;;
        esac
      done
      [ "$N" -gt 0 ] || hint_exit "The last ls had no pattern. Run it with the quarterly pattern."
    hint: |
      echo "The question mark stands for one character. Put it where the quarter's digit goes."
    solve: |
      cd ~/projects/data
      ls sales-q?.csv
---

The four quarterly files are `sales-q1.csv` through `sales-q4.csv`. A
`sales-*` pattern would also match `sales-2025.csv`, which is not a
quarterly file. The `?` wildcard is stricter: it stands for exactly one
character.

List only the four quarterly files in `~/projects/data` with a single
pattern.

::task{name="matched"}
#active
Waiting for a pattern with a one-character wildcard...
#completed
Four names matched, and `sales-2025.csv` stayed out.
::
