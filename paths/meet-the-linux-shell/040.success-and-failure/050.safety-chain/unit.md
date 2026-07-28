---
title: Only if it worked
vars:
  PAUSE: { pick: ["2", "3"] }
tasks:
  chained:
    timeout: 45
    check: |
      wait_exec "(^|/)sleep ${PAUSE}s?\$"
      wait_exec '(^|/)hostname$'
    hint: |
      echo "One line: the ${PAUSE}-second sleep, then &&, then hostname."
    solve: |
      sleep $PAUSE && hostname
---

Here is where the exit status starts working for you. The operator `&&`
("and-and") chains two commands with a condition: the second one runs
**only if the first succeeded** (exit status `0`).

In one line: pause for **${PAUSE} seconds**, and then - only if the
pause finished properly - print the machine's name.

::task
#active
Waiting for a ${PAUSE}-second pause followed by the machine name, one
line, chained on success...
#completed
The name appeared only after the pause delivered its `0`. Had the
first command failed, the shell would have skipped the second
entirely - that's `&&`: "and, if that worked, ...".
::
