---
title: A space in the name
vars:
  FIRST: { pick: [Meeting, Weekly, Planning, Review] }
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
  - name: clear_targets
    run: |
      N="$GYM_USER_HOME/projects/notes"
      rm -f "$N/$FIRST" "$N/notes.txt" "$N/$FIRST notes.txt"
      mkdir -p /run/shellgym-refs
      touch /run/shellgym-refs/a-space-in-the-name
tasks:
  created:
    timeout: 60
    check: |
      wait_file "$GYM_USER_HOME/projects/notes/$FIRST notes.txt"
    hint: |
      N="$GYM_USER_HOME/projects/notes"
      STRAY=$(find "$GYM_USER_HOME" -maxdepth 4 -name notes.txt -newer /run/shellgym-refs/a-space-in-the-name 2>/dev/null | head -n 1)
      if [ -e "$N/$FIRST" ] || [ -n "$STRAY" ]; then
        echo "You now have a file named ${FIRST} and one named notes.txt. The space split the name in two. Quote the whole name so it reaches touch as one argument."
      else
        echo "Wrap the whole name in quotes, single or double, so the space stays inside one argument."
      fi
    solve: |
      touch ~/projects/notes/"$FIRST notes.txt"
---

A file name can contain a space, but to the shell a space separates
arguments. Written bare, `${FIRST} notes.txt` is two arguments, and
`touch` would happily create two files.

Create a single empty file named `${FIRST} notes.txt` in
`~/projects/notes`. You know the fix from the earlier gyms: quotes keep
the name together.

::task{name="created"}
#active
Waiting for `${FIRST} notes.txt` to appear as one file...
#completed
One file, with a space in its name.
::

::tip
Spaces in names are legal and common, but they cost extra quoting.
Naming your files using `-` or `_` as word separators can be a handy alternative.
::
