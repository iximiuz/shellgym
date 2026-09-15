---
title: One step down
requires: [readline]
init:
  - name: build_work_tree
    run: |
      W="$GYM_USER_HOME/projects"
      mkdir -p "$W/reports/drafts" "$W/reports/final" "$W/archive/2025" "$W/data/imports" "$W/data/exports" "$W/notes"
      [ -f "$W/reports/final/summary.txt" ] || echo "Final report, approved." > "$W/reports/final/summary.txt"
      [ -f "$W/reports/drafts/outline.txt" ] || echo "Report outline, work in progress." > "$W/reports/drafts/outline.txt"
      [ -f "$W/archive/2025/report-2025.txt" ] || echo "Last year's report." > "$W/archive/2025/report-2025.txt"
      [ -f "$W/data/imports/inventory.csv" ] || echo "id,name,quantity" > "$W/data/imports/inventory.csv"
      [ -f "$W/data/exports/inventory-export.csv" ] || echo "id,name,quantity" > "$W/data/exports/inventory-export.csv"
      [ -f "$W/notes/todo.txt" ] || echo "- review the final report" > "$W/notes/todo.txt"
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
  arrived:
    needs: [at_home]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/projects" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        /*|"~"*|"\$"*) hint_exit "You got there with a full path. From home, the name alone is enough: the shell looks for it inside the current directory. Go back home and try just the name." ;;
        *) exit 0 ;;
      esac
    hint: |
      echo "Give cd only the name of the directory, with no slash in front of it."
    solve: |
      cd projects
---

Start in your home directory. The `projects` directory is right inside
it, so you do not need its full path. A path that does not start with
`/` is **relative**: the shell looks for it inside the directory you
are standing in.

Move into `projects` by its name alone.

::task{name="at_home"}
#active
Waiting for your shell to be in your home directory...
#completed
You are home. Now step into `projects`.
::

::task{name="arrived"}
#active
Waiting for a relative `cd` into `projects`...
#completed
The shell found `projects` inside the current directory.
::

::tip
Type `cd p` and press **Tab**. The shell completes the name for you,
because `projects` is the only name here that starts with `p`.
::
