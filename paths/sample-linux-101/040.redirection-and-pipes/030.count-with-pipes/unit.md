---
title: Count matches with a pipe
vars:
  ANIMAL: { pick: [heron, badger, lynx, otter] }
init:
  - name: seed_journal
    run: |
      rm -f /tmp/journal.txt
      COUNT=$((RANDOM % 20 + 5))
      echo "$COUNT" > /tmp/journal_expected
      OTHERS=(fox crow deer hare)
      for i in $(seq 1 "$COUNT"); do
        echo "sighting: $ANIMAL near the river ($i)" >> /tmp/journal.txt
        echo "sighting: ${OTHERS[$((RANDOM % 4))]} in the field" >> /tmp/journal.txt
      done
      chmod a+r /tmp/journal.txt
tasks:
  counted:
    check: |
      EXPECTED=$(cat /tmp/journal_expected)
      wait_file_contains "$GYM_USER_HOME/count.txt" "^${EXPECTED}\s*$"
    hint: |
      if [ -f "$GYM_USER_HOME/count.txt" ]; then
        echo "count.txt exists but the number inside is not right. Count lines that mention ${ANIMAL} in /tmp/journal.txt - grep can filter, wc -l can count, and a pipe connects them. (grep -c does both in one step.)"
      else
        echo "Filter the journal with grep, count the matching lines, and redirect the resulting number into ~/count.txt."
      fi
    solve: |
      grep $ANIMAL /tmp/journal.txt | wc -l > ~/count.txt
---

A field journal at `/tmp/journal.txt` logs animal sightings, one per line.
Count how many lines mention `${ANIMAL}` and save that number into
`~/count.txt`.

Doing this by eye would be miserable. Instead, chain tools with a pipe:
`grep` selects the matching lines, `wc -l` counts lines fed to it, and `|`
connects the two. Finish with a `>` redirect to capture the number.

::task{name="counted"}
#active
Waiting for the correct count in `~/count.txt`...
#completed
Correct. `filter | count > file` is a pattern you will use weekly for logs,
processes, and inventories of all kinds.
::

::hint
---
title: How pipes work
---
`command-a | command-b` runs both at once: whatever `command-a` prints
becomes the input of `command-b`. Any number of commands can be chained
this way.
::
