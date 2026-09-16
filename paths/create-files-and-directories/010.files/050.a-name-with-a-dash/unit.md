---
title: A name that starts with a dash
requires: [readline]
vars:
  NAME: { pick: [-draft.txt, -old.txt, -archive.txt, -notes.txt] }
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
      rm -f "$GYM_USER_HOME/projects/notes/$NAME"
tasks:
  tried:
    timeout: 60
    check: |
      R="^touch +${NAME//./\\.} *\$"
      wait_line "$R"
    hint: |
      echo "Run touch ${NAME} exactly like that, with the bare name. It will fail, and the error message is the point."
    solve: |
      cd ~/projects/notes
      touch $NAME
  created:
    needs: [tried]
    timeout: 60
    check: |
      wait_file "$GYM_USER_HOME/projects/notes/$NAME"
    hint: |
      echo "End the options with -- before the name, or write the name as a path: ./${NAME} from inside notes, or its full path from anywhere."
    solve: |
      touch -- $NAME
---

You know from [the listing files gym](https://labs.iximiuz.com/shell-gyms/list-files-and-use-patterns) that a name which starts with a dash
looks like an option to the command that gets it. The same rule
applies when creating a file.

Create an empty file named `${NAME}` in `~/projects/notes`. Try the
bare name first and read the error: `touch` complains about an option
in the same way `ls` did. Then protect the name. A bare `--` argument
ends the options, so everything after it is a name. Or write the name
as a path, `./${NAME}` from inside `notes` or its full path from
anywhere, so it no longer starts with a dash.

::task{name="tried"}
#active
Waiting for the doomed `touch ${NAME}`...
#completed
That error is `touch` reading the name as options. Now protect the name.
::

::task{name="created"}
#active
Waiting for `${NAME}` to appear in `notes`...
#completed
The file exists, dash and all.
::

::tip
Quotes do not help here. They protect a name from the shell, and the
shell is not the problem: `touch` itself decides that a dash starts an
option. The `--` marker works with almost every Linux command.
::
