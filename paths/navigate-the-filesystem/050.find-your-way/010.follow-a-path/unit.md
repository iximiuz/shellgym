---
title: Follow a path to a file
vars:
  FILE: { pick: [data/imports/inventory.csv, reports/final/summary.txt, archive/2025/report-2025.txt, notes/todo.txt, data/exports/inventory-export.csv] }
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
  arrived:
    check: |
      P=$(wait_cwd "$(dirname "$GYM_USER_HOME/projects/$FILE")") || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Everything before the last slash is the directory. Move there and leave the file name off the path."
    solve: |
      cd ~/projects/$(dirname $FILE)
  listed:
    needs: [arrived]
    timeout: 60
    check: |
      wait_exec '(^|/)ls( .*)?$'
    hint: |
      echo "Run ls to see the files in the directory you are in. The one from the path should be among them."
    solve: |
      ls
---

A teammate asks you to check on a file and gives you its path:

```
~/projects/${FILE}
```

The last part of a path is the file. Everything before it is the
directory the file sits in. Move your shell into that directory, then
list its contents to confirm the file is really there.

::task{name="arrived"}
#active
Waiting for your shell to arrive in the directory that holds the file...
#completed
This is the directory from the path. Now confirm the file is in it.
::

::task{name="listed"}
#active
Waiting for you to list the directory...
#completed
The file is there, right where the path said.
::

::tip
Paste the path, then use the **Left** arrow or **Ctrl-W**, which
erases the word before the cursor, to remove the file name before
pressing Enter.
::
