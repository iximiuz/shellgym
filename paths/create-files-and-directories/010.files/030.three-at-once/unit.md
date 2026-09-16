---
title: Three at once
vars:
  A: { pick: [agenda.txt, minutes.txt] }
  B: { pick: [actions.txt, decisions.txt] }
  C: { pick: [followups.txt, questions.txt] }
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
  - name: clear_targets
    run: |
      rm -f "$GYM_USER_HOME/projects/notes/$A" "$GYM_USER_HOME/projects/notes/$B" "$GYM_USER_HOME/projects/notes/$C"
tasks:
  created:
    timeout: 60
    check: |
      N="$GYM_USER_HOME/projects/notes"
      wait_file "$N/$A" >/dev/null || exit 1
      wait_file "$N/$B" >/dev/null || exit 1
      wait_file "$N/$C" >/dev/null || exit 1
      ARGV=$(wait_exec --latest '(^|/)touch( +\S+)+$') || exit 1
      case " $ARGV " in
        *" $A "*|*"/$A "*) ;;
        *) hint_exit "All three files exist, but the last touch did not name all of them. Run one touch with the three names as its arguments." ;;
      esac
      case " $ARGV " in
        *" $B "*|*"/$B "*) ;;
        *) hint_exit "All three files exist, but the last touch did not name all of them. Run one touch with the three names as its arguments." ;;
      esac
      case " $ARGV " in
        *" $C "*|*"/$C "*) exit 0 ;;
        *) hint_exit "All three files exist, but the last touch did not name all of them. Run one touch with the three names as its arguments." ;;
      esac
    hint: |
      echo "One touch, three arguments: ${A}, ${B}, and ${C}, separated by spaces. Run it in ~/projects/notes, or give each name with the directory's path in front."
    solve: |
      cd ~/projects/notes
      touch $A $B $C
---

The `touch` command takes as many paths as you give it and creates a
file for each. That is one command instead of three.

Create `${A}`, `${B}`, and `${C}` in `~/projects/notes` with a single
command.

::task{name="created"}
#active
Waiting for one `touch` to create all three files...
#completed
Three files from one command.
::
