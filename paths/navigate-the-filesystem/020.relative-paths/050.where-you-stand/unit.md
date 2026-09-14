---
title: It depends on where you stand
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
  at_home:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start from your home directory. A bare cd takes you there."
    solve: |
      cd
  tried:
    needs: [at_home]
    timeout: 60
    check: |
      wait_line '^cd +reports/final/? *$'
    hint: |
      echo "Run cd reports/final from your home directory, even though it is going to fail."
    solve: |
      cd reports/final
  arrived:
    needs: [tried]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME/work/reports/final" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        /*|"~"*|"\$"*) hint_exit "You got there with a full path. Go back home and use a relative one: the whole way from home, three names joined by slashes." ;;
        *) exit 0 ;;
      esac
    hint: |
      echo "From home, the way to final passes through work first. Put all three names in the path."
    solve: |
      cd work/reports/final
---

A relative path is only as good as the place you use it from. The same
`reports/final` works from `~/work` and fails from anywhere else.

Start in your home directory and run `cd reports/final` anyway. Read
the error. Your shell has not moved.

Then get to `~/work/reports/final` with a relative path that is correct
from where you actually stand.

::task{name="at_home"}
#active
Waiting for your shell to be in your home directory...
#completed
You are home. Now try the path that only works from `~/work`.
::

::task{name="tried"}
#active
Waiting for the doomed `cd reports/final`...
#completed
There is no `reports` in your home, so the shell refused. Now take the route that works from here.
::

::task{name="arrived"}
#active
Waiting for a relative `cd` that reaches `final` from home...
#completed
The full relative route from home worked.
::

::tip
When a relative path fails, run `pwd` first. Often the path is fine
and the shell is simply standing in the wrong place.
::
