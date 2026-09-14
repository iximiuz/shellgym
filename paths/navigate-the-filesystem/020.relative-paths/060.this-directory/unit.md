---
title: This directory
requires: [readline]
init:
  - name: build_work_tree
    run: |
      W="$GYM_USER_HOME/work"
      mkdir -p "$W/reports/drafts" "$W/reports/final" "$W/archive/2025" "$W/data/imports" "$W/data/exports" "$W/notes"
      [ -f "$W/reports/final/summary.txt" ] || echo "Final report, approved." > "$W/reports/final/summary.txt"
      [ -f "$W/reports/drafts/outline.txt" ] || echo "Report outline, work in progress." > "$W/reports/drafts/outline.txt"
      [ -f "$W/archive/2025/report-2025.txt" ] || echo "Last year's report." > "$W/archive/2025/report-2025.txt"
      [ -f "$W/data/imports/inventory.csv" ] || echo "id,name,quantity" > "$W/data/imports/inventory.csv"
      [ -f "$W/data/exports/inventory-export.csv" ] || echo "id,name,quantity" > "$W/data/exports/inventory-export.csv"
      [ -f "$W/notes/todo.txt" ] || echo "- review the final report" > "$W/notes/todo.txt"
      chown -R "$GYM_USER:$GYM_USER" "$W"
tasks:
  at_work:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/work") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/work. From anywhere, cd ~/work gets you there."
    solve: |
      cd ~/work
  arrived:
    needs: [at_work]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/notes" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        ./*) exit 0 ;;
        *) hint_exit "You are in notes, but the path did not start with the one-dot name. Go back to ~/work and write the path as a dot, a slash, and the name." ;;
      esac
    hint: |
      echo "Start the path with a single dot and a slash, then the directory name."
    solve: |
      cd ./notes
---

Next to `..` there is a second special name, `.` (a single dot),
which means the current directory itself. The path `./notes` reads "in this
directory, the thing called `notes`". It reaches the same place as
plain `notes`.

Start in `~/work` and move into `notes` with a path that begins with
the one-dot name.

::task{name="at_work"}
#active
Waiting for your shell to be in `~/work`...
#completed
You are in `~/work`. Now enter `notes` through the one-dot name.
::

::task{name="arrived"}
#active
Waiting for a `cd` into `./notes`...
#completed
The dot stood for the directory you were in.
::

::note{title="Why bother with the dot?"}
For `cd` the dot changes nothing. It matters when a name could be
mistaken for something else: `./date` means the file called `date`
right here, while a bare `date` runs the program from `/usr/bin`. You
will see it in front of scripts all the time.
::
