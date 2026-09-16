---
title: Hidden stays hidden
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
tasks:
  everything:
    timeout: 60
    check: |
      LINE=$(wait_line --latest '^echo +(\S*/)?\* *$') || exit 1
      case "$LINE" in *"projects/data/"*) exit 0 ;; esac
      CWD=$(shell_cwd 2>/dev/null || true)
      [ "$CWD" = "$GYM_USER_HOME/projects/data" ] || hint_exit "Your shell is in ${CWD:-another directory}, so the star was expanded there. Move to ~/projects/data first, or put the directory's path in front of the star."
    hint: |
      echo "Preview the bare star with echo: echo, a space, and a star, in ~/projects/data."
    solve: |
      cd ~/projects/data
      echo *
  hidden:
    needs: [everything]
    timeout: 60
    check: |
      LINE=$(wait_line --latest '^(ls|echo)( +-\S+)* +(\S*/)?\.\S*\*\S* *$') || exit 1
      case "$LINE" in *"projects/data/"*) exit 0 ;; esac
      CWD=$(shell_cwd 2>/dev/null || true)
      [ "$CWD" = "$GYM_USER_HOME/projects/data" ] || hint_exit "Your shell is in ${CWD:-another directory}, so the pattern was expanded there. Move to ~/projects/data first, or put the directory's path in front of the pattern."
    hint: |
      echo "Start the pattern with a dot, then the star. Preview it with echo."
    solve: |
      echo .*
---

The `*` matches any sequence of characters, with one exception: it never
matches a leading dot. Hidden names stay out of patterns unless the
pattern spells the dot out.

First preview a bare `*` in `~/projects/data` with `echo` and note that
`.last-import` is missing. Then preview a pattern that starts with a dot
and see the hidden name appear.

::task{name="everything"}
#active
Waiting for `echo *`...
#completed
No `.last-import` in that list. Now spell out the dot.
::

::task{name="hidden"}
#active
Waiting for a pattern that starts with a dot...
#completed
The dot pattern reached the hidden name.
::

::note{title="Two extra matches"}
A pattern that starts with a dot also matches `.` and `..`, the
current and parent directories. That is harmless for `echo` and
dangerous for commands that change things. You will meet a safer form
in a later gym.
::
