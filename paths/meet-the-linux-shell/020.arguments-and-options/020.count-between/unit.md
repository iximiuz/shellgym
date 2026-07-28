---
title: Count between two numbers
vars:
  LO: { pick: ["2", "3", "4"] }
  HI: { pick: ["11", "12", "13"] }
tasks:
  counted:
    check: |
      wait_exec "(^|/)seq ${LO} ${HI}\$"
    hint: |
      echo "Two arguments this time: first ${LO}, then ${HI}, separated by whitespace."
    solve: |
      seq $LO $HI
---

`seq` also accepts **two** arguments: where to start and where to stop.

Count from **${LO}** to **${HI}**. Remember: it is whitespace that
separates one argument from the next - one space is enough, extra
spaces do no harm.

::task
#active
Waiting for a count from ${LO} to ${HI}...
#completed
Two arguments, split on whitespace. That splitting is done by the
shell before `seq` ever sees them - keep that in mind, it will matter
soon.
::

::tip
Don't retype: **Up** arrow, edit the numbers, Enter.
::
