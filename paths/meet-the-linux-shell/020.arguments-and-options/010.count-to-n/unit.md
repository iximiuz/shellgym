---
title: Count up
vars:
  LIMIT: { pick: ["7", "9", "12", "14"] }
tasks:
  counted:
    check: |
      wait_exec "(^|/)seq ${LIMIT}\$"
    hint: |
      echo "The command name, one space, then the number ${LIMIT}. Nothing else."
    solve: |
      seq $LIMIT
---

Meet `seq`. Given a single number as its argument, it counts out loud
from 1 up to that number, one per line.

Make it count to **${LIMIT}**.

::task
#active
Waiting for a count from 1 to ${LIMIT}...
#completed
The number was an argument: the shell handed it to `seq`, and `seq`
decided what it meant. That division of labor never changes.
::
