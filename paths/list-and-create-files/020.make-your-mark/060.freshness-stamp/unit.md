---
title: Freshen the stamp
vars:
  REPORT: { pick: [quarterly, weekly, annual] }
init:
  - name: create_inbox
    run: |
      rm -rf /tmp/gym-inbox
      mkdir -p /tmp/gym-inbox /run/gym-refs
      echo "figures pending" > "/tmp/gym-inbox/$REPORT-report.txt"
      touch -d "2 days ago" "/tmp/gym-inbox/$REPORT-report.txt"
      chown "$GYM_USER" "/tmp/gym-inbox/$REPORT-report.txt"
      touch /run/gym-refs/freshness
tasks:
  stamped:
    check: |
      wait_file_newer "/tmp/gym-inbox/$REPORT-report.txt" /run/gym-refs/freshness
    hint: |
      TS=$(stat -c %y "/tmp/gym-inbox/$REPORT-report.txt" 2>/dev/null | cut -d. -f1)
      echo "The report's timestamp is still $TS. touch on an EXISTING file does not change its content - it updates the modification time to now."
    solve: |
      touch /tmp/gym-inbox/$REPORT-report.txt
---

The `${REPORT}-report.txt` in `/tmp/gym-inbox` was last modified two
days ago - see for yourself with `ls -l`. An auditing script upstream
refuses files older than a day, and the content is fine as it is.

Here is `touch`'s second talent: on a file that **already exists**, it
simply updates the modification time to *now*, content untouched.

::task{name="stamped"}
#active
Waiting for `${REPORT}-report.txt` to carry a fresh timestamp...
#completed
Stamped today. That is actually `touch`'s original job - the
create-if-missing behavior is the side effect that became famous.
Compare `ls -l` before and after if you kept the old output around.
::
