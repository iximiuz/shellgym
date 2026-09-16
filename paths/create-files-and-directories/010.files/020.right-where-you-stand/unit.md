---
title: Right where you stand
vars:
  NAME: { pick: [review.txt, comments.txt, figures.txt, sources.txt] }
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
  - name: clear_target
    run: |
      find "$GYM_USER_HOME/projects" -name "$NAME" -type f -delete
tasks:
  at_drafts:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/projects/reports/drafts") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/projects/reports/drafts. From anywhere, cd ~/projects/reports/drafts gets you there."
    solve: |
      cd ~/projects/reports/drafts
  created:
    needs: [at_drafts]
    timeout: 60
    check: |
      wait_file "$GYM_USER_HOME/projects/reports/drafts/$NAME" >/dev/null || exit 1
      ARGV=$(wait_exec --latest '(^|/)touch( +\S+)+$') || exit 1
      set -f; set -- $ARGV; shift
      for a in "$@"; do
        case "$a" in
          -*) ;;
          */*) hint_exit "The file exists, but touch got a path. The point here is the bare name: stand in ~/projects/reports/drafts and give touch just ${NAME}." ;;
        esac
      done
      case " $ARGV " in *" $NAME "*) exit 0 ;; esac
      hint_exit "The last touch did not name ${NAME}. Give touch the bare name from inside drafts."
    hint: |
      F=$(find "$GYM_USER_HOME" -name "$NAME" -not -path "*/reports/drafts/*" 2>/dev/null | head -n 1)
      if [ -n "$F" ]; then
        echo "There is a ${NAME} at $F, but the assignment is the drafts directory. Check where your shell stands with pwd."
      else
        echo "Move into ~/projects/reports/drafts first, then give touch just the file name."
      fi
    solve: |
      touch $NAME
---

A relative path works for `touch`, too: a bare file
name means "in the directory I am standing in".

Start in `~/projects/reports/drafts` and create an empty file named
`${NAME}` there, using the name alone.

::task{name="at_drafts"}
#active
Waiting for your shell to be in `~/projects/reports/drafts`...
#completed
You are in `drafts`. Now create the file by its name alone.
::

::task{name="created"}
#active
Waiting for `${NAME}` to appear in `drafts`...
#completed
The file landed in the directory your shell was standing in.
::

::tip
When a file ends up in the wrong place, the cause is almost always the
working directory. Run `pwd` before creating anything if you are not
sure where you are.
::
