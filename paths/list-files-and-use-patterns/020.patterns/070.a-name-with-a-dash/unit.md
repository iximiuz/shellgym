---
title: A name that starts with a dash
requires: [readline]
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
  - name: plant_the_file
    run: |
      N="$GYM_USER_HOME/projects/notes"
      [ -f "$N/-report.txt" ] || echo "An oddly named report." > "$N/-report.txt"
      chown "$GYM_USER:$GYM_USER" "$N/-report.txt"
tasks:
  tried:
    timeout: 60
    check: |
      wait_line '^ls +-report\.txt *$'
    hint: |
      echo "Run ls -report.txt exactly like that, without protection. It will fail, and the error message is the point."
    solve: |
      cd ~/projects/notes
      ls -report.txt
  listed:
    needs: [tried]
    timeout: 60
    check: |
      ARGV=$(wait_exec --latest '(^|/)ls( +-\S+)* +(-- +-report\.txt|\./-report\.txt|\S*/notes/-report\.txt)/? *$') || exit 1
      case "$ARGV" in *"/notes/-report.txt"*) exit 0 ;; esac
      CWD=$(shell_cwd 2>/dev/null || true)
      [ "$CWD" = "$GYM_USER_HOME/projects/notes" ] || hint_exit "Your shell is in ${CWD:-another directory}, so ls looked for -report.txt there. Move to ~/projects/notes first, or give the file's full path instead of the bare name."
    hint: |
      echo "Two ways work: put -- before the name to end the options, or write the name as a path, ./-report.txt or its full path, so it no longer starts with a dash."
    solve: |
      ls -- -report.txt
---

A file named `-report.txt` sits in `~/projects/notes`. Its name starts with
a dash, and every command reads a leading dash as the start of an
option.

Run `ls -report.txt` and read the error: `ls` took the name for a bundle
of options. Then list the file for real. There are two ways. A bare `--`
argument tells a command that the options are over and everything after
it is a name. Or write the name as a path, `./-report.txt` from inside
`notes` or its full path from anywhere, so it no longer starts with a
dash.

::task{name="tried"}
#active
Waiting for the doomed `ls -report.txt`...
#completed
That error is `ls` complaining about an option it does not know. Now protect the name.
::

::task{name="listed"}
#active
Waiting for `-report.txt` to be listed as a file name...
#completed
The dash reached `ls` as part of a name this time.
::

::note
Quotes do not help here. They protect a name from the shell, and the
shell is not the problem: `ls` itself decides that a dash starts an
option. The `--` marker works with almost every Linux command.
::
