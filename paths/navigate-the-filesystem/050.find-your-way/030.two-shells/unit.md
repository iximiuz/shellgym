---
title: Two shells, two directories
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
  parked:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/work/notes") || exit 1
      set_var OUTER_PID "$P"
    hint: |
      echo "Start in ~/work/notes. From anywhere, cd ~/work/notes gets you there."
    solve: |
      cd ~/work/notes
  nested:
    needs: [parked]
    timeout: 60
    check: |
      wait_line '^bash *$'
    hint: |
      echo "Start a second shell inside the first by running bash with no arguments."
    solve: |
      bash
  moved:
    needs: [nested]
    timeout: 60
    check: |
      P=$(wait_cwd /tmp) || exit 1
      [ "$P" != "$OUTER_PID" ] || hint_exit "The shell that moved to /tmp is the outer one. Start the inner shell with bash first, then move inside it."
    hint: |
      echo "Inside the inner shell, move to /tmp."
    solve: |
      cd /tmp
  returned:
    needs: [moved]
    timeout: 60
    check: |
      wait_line '^exit *$' >/dev/null || exit 1
      wait_cwd "$OUTER_PID" "$GYM_USER_HOME/work/notes"
    hint: |
      echo "Leave the inner shell with exit. The outer shell is still in notes."
    solve: |
      exit
---

The working directory belongs to a shell, and every shell has its own.
A shell inside a shell demonstrates this clearly.

Start in `~/work/notes`. Then run `bash` with no arguments: this starts
a second shell inside the first, and the prompt you see now belongs to
the inner one. Move the inner shell to `/tmp`. Finally, run `exit` to
close the inner shell and look at the prompt: the outer shell never
moved.

::task{name="parked"}
#active
Waiting for your shell to be in `~/work/notes`...
#completed
The outer shell is parked in `notes`. Now start the inner one.
::

::task{name="nested"}
#active
Waiting for you to start an inner `bash`...
#completed
You are talking to the inner shell now. Move it to `/tmp`.
::

::task{name="moved"}
#active
Waiting for the inner shell to move to `/tmp`...
#completed
The inner shell is in `/tmp`, and the outer one is still in `notes`.
::

::task{name="returned"}
#active
Waiting for `exit` to bring the outer shell's prompt back...
#completed
The outer shell is exactly where you left it.
::

::note{title="One working directory per process"}
Every terminal window, every script, and every program you start
gets a working directory of its own. Moving in one of them never
moves another. That is why a `cd` inside a script does not affect
the shell that ran the script.
::
