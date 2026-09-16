---
title: A project tree
vars:
  NAME: { pick: [website, inventory, newsletter, migration] }
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
      rm -rf "$GYM_USER_HOME/projects/$NAME"
tasks:
  created:
    timeout: 60
    check: |
      P="$GYM_USER_HOME/projects/$NAME"
      wait_dir "$P/data" >/dev/null || exit 1
      wait_dir "$P/reports" >/dev/null || exit 1
      wait_dir "$P/notes"
    hint: |
      P="$GYM_USER_HOME/projects/$NAME"
      if [ ! -d "$P" ]; then
        echo "${NAME} does not exist yet, so mkdir needs the option that creates missing parents, or create ${NAME} first."
      else
        for d in data reports notes; do [ -d "$P/$d" ] || echo "$P/$d is still missing."; done
      fi
    solve: |
      mkdir -p ~/projects/$NAME/data ~/projects/$NAME/reports ~/projects/$NAME/notes
  listed:
    needs: [created]
    timeout: 60
    check: |
      wait_exec "(^|/)ls( +-\S+)* +\S*projects/${NAME}/?\$"
    hint: |
      echo "Give ls the path of the new project directory, ~/projects/${NAME}, to see what is in it."
    solve: |
      ls ~/projects/$NAME
---

Set up a place for a new project. Create the directory
`~/projects/${NAME}` with three directories inside it: `data`,
`reports`, and `notes`. The `${NAME}` directory does not exist yet.

Then list the contents of `${NAME}` to check the result.

::task{name="created"}
#active
Waiting for `data`, `reports`, and `notes` to appear in `~/projects/${NAME}`...
#completed
All three directories exist. Now check the tree.
::

::task{name="listed"}
#active
Waiting for the contents of the project directory to be listed...
#completed
The tree is ready for the project.
::

::tip
The parents option lets you create the whole tree with one command:
three full paths as arguments, and every missing level on the way is
created too.
::
