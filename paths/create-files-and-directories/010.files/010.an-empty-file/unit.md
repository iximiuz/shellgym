---
title: An empty file
vars:
  NAME: { pick: [status.txt, checklist.txt, changelog.txt, questions.txt] }
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
      rm -f "$GYM_USER_HOME/projects/notes/$NAME"
tasks:
  created:
    check: |
      wait_file "$GYM_USER_HOME/projects/notes/$NAME" >/dev/null || exit 1
      ARGV=$(wait_exec --latest '(^|/)touch +\S+$') || exit 1
      case "$ARGV" in
        *" /"*) exit 0 ;;
        *) hint_exit "The file exists, but touch got a relative path. Run it again with the full path ~/projects/notes/${NAME}." ;;
      esac
    hint: |
      echo "Run touch with the full path of the new file: ~/projects/notes/${NAME}, starting with the tilde."
    solve: |
      touch ~/projects/notes/$NAME
---

Create an empty file by running `touch` with the absolute path
`~/projects/notes/${NAME}`. The `touch` command takes the path of the
file and creates it, empty, if nothing exists there yet. A full path
tells it exactly where to put the file, wherever your shell stands.

::task
#active
Waiting for `~/projects/notes/${NAME}` to appear...
#completed
The file exists and is empty.
::

::tip
Listing contents of `~/projects/notes` with `ls -l` shows the new file with a size of `0`.
Empty files are useful as markers, and as a starting point for
something you will write later.
::
