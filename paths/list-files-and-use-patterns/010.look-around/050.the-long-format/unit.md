---
title: The long format
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
  listed:
    timeout: 60
    check: |
      D="$GYM_USER_HOME/projects/data"
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      set -f; set -- $ARGV; shift
      CWD=$(shell_cwd 2>/dev/null || true)
      N=0
      for a in "$@"; do
        case "$a" in -*) ;; *) N=$((N+1)); T="$a" ;; esac
      done
      if [ "$N" -eq 0 ]; then
        [ "$CWD" = "$D" ] || hint_exit "Your shell is in ${CWD:-another directory}, so that listed the wrong directory. Move to ~/projects/data first, or give ls the path ~/projects/data as its argument."
      elif [ "$N" -eq 1 ]; then
        case "$T" in /*) P=$(realpath -m "$T") ;; *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$T") ;; esac
        [ "$P" = "$D" ] || hint_exit "That listed $T instead of ~/projects/data."
      else
        hint_exit "That listed several paths at once. List only ~/projects/data."
      fi
      S=' -[a-zA-Z0-9]*l[a-zA-Z0-9]* '
      L=' --format=(long|verbose) '
      if [[ " $ARGV " =~ $S ]] || [[ " $ARGV " =~ $L ]]; then exit 0; fi
      hint_exit "That is the short format, names only. Add the long-format option, a single lowercase letter."
    hint: |
      CWD=$(shell_cwd 2>/dev/null)
      if [ -n "$CWD" ] && [ "$CWD" != "$GYM_USER_HOME/projects/data" ]; then
        echo "Your shell is in $CWD. Move to ~/projects/data first, or give ls the path ~/projects/data as its argument."
        exit 0
      fi
      echo "The long-format option is the lowercase letter l."
    solve: |
      ls -l ~/projects/data
---

By default, `ls` prints only the names of the entries. The **long format** `ls -l`
prints one line per entry with the details of that entry. A typical line
looks like this:

```
-rw-r--r-- 1 laborant laborant 1200 Sep 15 10:42 sales-q1.csv
```

From left to right, the columns are: the entry type (`-` for a regular file, `d` for a directory)
together with its permissions (e.g., `rw` for **r**ead-**w**rite), a count of names for the entry,
the **owner**, the **group**, the **size in bytes**, the date and time of the last change, and the name.
All of it will be covered in more details by later gyms.

List the contents of `~/projects/data` in the long format. Move there
first and list the working directory, or stay where you are and give
`ls` the path. Either way works.

::task{name="listed"}
#active
Waiting for a long-format listing of `~/projects/data`...
#completed
One line per entry, with the type, owner, size, and date of each.
::

::tip
The long format and the hidden names combine: `-l` and `-a` can be
written together as `-la`, exactly as you learned for short options in
Meet the Linux Shell.
::
