---
title: Look without moving
vars:
  TARGET: { pick: [reports, data, notes, backups, archive] }
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
  at_home:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start from your home directory. A bare cd takes you there."
    solve: |
      cd
  listed:
    needs: [at_home]
    timeout: 60
    check: |
      wait_exec --cwd "$GYM_USER_HOME" "(^|/)ls( +-\S+)* +\S*projects/${TARGET}/?\$"
    hint: |
      HERE=$(shell_cwd "$SHELL_PID")
      if [ "$HERE" != "$GYM_USER_HOME" ]; then
        echo "Your shell has moved to $HERE. Go back home and give ls the path of the directory as its argument instead of moving there."
      else
        echo "Give ls the path projects/${TARGET} as its argument, and stay where you are."
      fi
    solve: |
      ls projects/$TARGET
---

You do not have to enter a directory to see what is in it. Give `ls`
a path, absolute or relative, and it shows the contents of that directory while your
shell stays where it is.

Start in your home directory and list the contents of
`~/projects/${TARGET}` without moving.

::task{name="at_home"}
#active
Waiting for your shell to be in your home directory...
#completed
You are home. Now list the contents of `projects/${TARGET}` from here.
::

::task{name="listed"}
#active
Waiting for the contents of `projects/${TARGET}` to be listed from your home directory...
#completed
The listing came to you, and your shell never moved.
::

::tip
Tab completion works on the argument of `ls` just as it does on `cd`.
Type `ls pr`, press **Tab**, and keep going one part at a time.
::
