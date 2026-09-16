---
title: What is here
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
  at_work:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/projects") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/projects. From anywhere, cd ~/projects gets you there."
    solve: |
      cd ~/projects
  listed:
    needs: [at_work]
    timeout: 60
    check: |
      ARGV=$(wait_exec --cwd "$GYM_USER_HOME/projects" --latest '(^|/)ls( +\S+)*$') || exit 1
      ARGS=${ARGV#* }
      [ "$ARGS" = "$ARGV" ] && ARGS=""
      for ARG in $ARGS; do
        case "$ARG" in
          -*) ;;
          *) hint_exit "That ls got a path as its argument. This unit requires listing the contents of the working directory itself: move to ~/projects and run ls with no arguments." ;;
        esac
      done
      HERE=$(shell_cwd "$SHELL_PID")
      if [ "$HERE" != "$GYM_USER_HOME/projects" ]; then
        hint_exit "Your shell is in $HERE, so that ls showed the contents of $HERE. Move to ~/projects and run ls again."
      fi
    hint: |
      echo "Run ls with no arguments. It shows the contents of the directory your shell is standing in."
    solve: |
      ls
---

First, a quick recap. With no arguments, `ls` shows the contents of the working directory: the names of the
files and directories in it, sorted alphabetically and arranged in
columns.

Start in `~/projects` and list its contents.

::task{name="at_work"}
#active
Waiting for your shell to be in `~/projects`...
#completed
You are in the `~/projects` directory. Now list its contents.
::

::task{name="listed"}
#active
Waiting for you to list the contents of the working directory...
#completed
Those names are the contents of `~/projects`.
::

::tip
On most systems `ls` prints directories in a different color than
files. The colors are a convenience of the terminal. The names are the
real answer.
::
