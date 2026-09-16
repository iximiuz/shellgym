---
title: Only the directories
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
      D="$GYM_USER_HOME/projects"
      LINE=$(wait_line --latest '^(ls|echo)( +\S+)* +(\S*/)?\*/ *$') || exit 1
      CWD=$(shell_cwd 2>/dev/null || true)
      ELSEWHERE="Your shell is in ${CWD:-another directory}, so the pattern was expanded there. Move to ~/projects first, or put the directory's path in front of the pattern."
      case "$LINE" in
        echo*)
          case "$LINE" in *"projects/*/"*) exit 0 ;; esac
          [ "$CWD" = "$D" ] || hint_exit "$ELSEWHERE"
          exit 0 ;;
      esac
      DOPT=' -[a-zA-Z0-9]*d'
      if [[ "$LINE" =~ $DOPT ]] || [[ "$LINE" == *"--directory"* ]]; then
        # A builtin echo leaves no exec event, so only an ls line can be placed:
        # its expanded arguments tell where the pattern was applied and
        # whether it matched what it should.
        ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
        set -f; set -- $ARGV; shift
        N=0
        for a in "$@"; do
          case "$a" in
            -*) continue ;;
            */) ;;
            *) hint_exit "The last ls listed something other than directories. Run ls -d with the */ pattern." ;;
          esac
          B="${a%/}"
          case "$B" in
            */*) case "$B" in /*) P=$(realpath -m "$B") ;; *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$B") ;; esac
                 [ "$(dirname "$P")" = "$D" ] || hint_exit "The last ls listed $a, which is outside ~/projects. Point the pattern at ~/projects." ;;
            *) [ "$CWD" = "$D" ] || hint_exit "$ELSEWHERE" ;;
          esac
          N=$((N+1))
        done
        [ "$N" -gt 0 ] || hint_exit "The last ls had no pattern. Run ls -d with */."
        exit 0
      fi
      hint_exit "The pattern matched the directories, but ls then listed the contents of each one. Add the option that makes ls show a directory's name instead of its contents. Look for 'directory' in ls --help."
    hint: |
      echo "The pattern is a star followed by a slash. Preview it with echo, or give it to ls together with the option that lists directories themselves."
    solve: |
      ls -d ~/projects/*/
---

A `/` at the end of a pattern makes it match directories only: `*/`
means "every directory in the working directory", and files are left
out.

List only the directories of `~/projects`, using that pattern.
There is one snag: given a directory, `ls` lists what is inside it.
Its `--help` has an option that makes it show the directory itself
instead.

::task{name="matched"}
#active
Waiting for a pattern that matches only directories...
#completed
Only the directories came out, each with a trailing slash.
::
