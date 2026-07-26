---
title: Room for a name
vars:
  MEMO: { pick: [meeting notes, status report, team roster] }
init:
  - name: create_office
    run: |
      rm -rf /tmp/gym-office
      mkdir -p /tmp/gym-office
      chown "$GYM_USER" /tmp/gym-office
tasks:
  memo_filed:
    check: |
      wait_file "/tmp/gym-office/$MEMO.txt"
    hint: |
      echo "Unquoted, the space splits '${MEMO}.txt' into separate arguments and touch creates several wrong files. Quote the whole name so it stays one argument - then ls the office to spot any stray files to clean later."
    solve: |
      touch "/tmp/gym-office/$MEMO.txt"
---

The office at `/tmp/gym-office` needs a file called `${MEMO}.txt` -
space in the name, as delivered by whoever wrote the ticket.

You met this in the navigation path from the *reading* side; now you
are on the *writing* side. Unquoted, `touch` would happily create two
or three wrongly-named files instead of one right one.

::task{name="memo_filed"}
#active
Waiting for `${MEMO}.txt` (one file, space included) in
`/tmp/gym-office`...
#completed
One file, space and all. If a stray unquoted attempt left fragments
behind, `ls` will show them - a cleanup rep for the next path.
::

::tip
---
title: Quotes or backslash
---
`touch "a b.txt"` and `touch a\ b.txt` create the same file. Quotes
read better; the backslash is what Tab completion inserts.
::
