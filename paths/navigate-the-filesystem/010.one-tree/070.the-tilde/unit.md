---
title: The tilde
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
  away:
    check: |
      P=$(wait_cwd /etc) || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Move to /etc first."
    solve: |
      cd /etc
  arrived:
    needs: [away]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        "~"*) exit 0 ;;
        *) hint_exit "You got there, but without the tilde. Go back to /etc and enter work in one command that starts with ~/ instead of the full path of your home." ;;
      esac
    hint: |
      echo "Write the path as the tilde, a slash, and the directory name: ~/work. One command does it."
    solve: |
      cd ~/work
---

The tilde, `~`, is the shell's short name for your home directory.
It works inside paths too: `~` stands for the full path of your home,
so `~/work` means the `work` directory inside it.

A `work` directory has been prepared in your home for this gym. First
move to `/etc`. It holds the system's configuration files. Then jump
straight into `~/work` with a single command, using the tilde instead
of typing out the path of your home.

::task{name="away"}
#active
Waiting for your shell to arrive in `/etc`...
#completed
You are in `/etc`. Now jump into `~/work` in one move.
::

::task{name="arrived"}
#active
Waiting for a `cd` through the tilde to land you in `~/work`...
#completed
The tilde expanded to your home and the rest of the path followed.
::

::tip
The tilde is expanded by the shell before `cd` ever sees it. It only
works at the start of a path, and only when it is not inside quotes.
::
