---
title: Odd numbers only
vars:
  ODD: { pick: ["9", "11", "13", "15"] }
tasks:
  counted:
    check: |
      wait_exec "(^|/)seq 1 2 ${ODD}\$"
    hint: |
      echo "Three arguments in this order: the start (1), the step (2), the stop (${ODD})."
    solve: |
      seq 1 2 $ODD
---

One more trick from `seq`: with **three** arguments they mean *start*,
*step*, and *stop*. So a step of 2 skips every other number.

Print the odd numbers from **1** to **${ODD}**.

::task
#active
Waiting for 1, 3, 5, ... up to ${ODD}...
#completed
Same command, and the meaning of each argument depended on how many
you gave. Commands read their manuals; you soon will too.
::
