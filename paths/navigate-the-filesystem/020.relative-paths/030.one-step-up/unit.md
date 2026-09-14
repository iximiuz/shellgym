---
title: One step up
requires: [readline]
vars:
  START: { pick: [reports/final, reports/drafts, data/imports, archive/2025] }
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
  at_start:
    check: |
      P=$(wait_cwd "$GYM_USER_HOME/work/$START") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Start in ~/work/${START}. From anywhere, cd ~/work/${START} gets you there."
    solve: |
      cd ~/work/$START
  arrived:
    needs: [at_start]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$(dirname "$GYM_USER_HOME/work/$START")" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        ..|../) exit 0 ;;
        *) hint_exit "You are in the parent, but you spelled its path out. Go back down and climb with the two-dot name instead." ;;
      esac
    hint: |
      echo "The parent of any directory is called .. (two dots). Give that name to cd."
    solve: |
      cd ..
---

Every directory has a parent, the directory that contains it. The
relative name of the parent is always `..` (two dots). It works from
anywhere, without knowing what the parent is called.

Start in `~/work/${START}`. Then move one level up, into the parent,
using the two-dot name.

::task{name="at_start"}
#active
Waiting for your shell to be in `~/work/${START}`...
#completed
You are in `${START}`. Now go one level up.
::

::task{name="arrived"}
#active
Waiting for a `cd` to `..`...
#completed
You are in the parent, and `pwd` confirms which directory that is.
::

::note{title="The one exception"}
The root directory `/` is the top of the tree, so there is nothing
above it. Its `..` points back to `/` itself, and `cd ..` from there
leaves you where you are.
::
