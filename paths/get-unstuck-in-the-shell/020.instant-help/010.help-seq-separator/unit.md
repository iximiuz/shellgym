---
title: Numbers on one line
variant: help1=seq-separator
vars:
  N: { pick: ["5", "7", "9"] }
tasks:
  consulted:
    check: |
      wait_exec '(^|/)seq --help$'
    hint: |
      echo "Run: seq --help"
    solve: |
      seq --help
  used:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)seq (-s|--separator=?) ?, ${N}\$"
    hint: |
      echo "The separator option takes the separator as its value: seq, the option with a comma right after it, then ${N}."
    solve: |
      seq -s, $N
---

By default, `seq` prints one number per line. It can also print all numbers on a
single line, separated by a character of your choice, and the option
for that is in its help:

```
seq --help
```

Find the option that sets the separator. Then make `seq` count from 1
to **${N}** with commas between the numbers, like `1,2,3`.

::task{name="consulted"}
#active
Waiting for you to ask `seq` for its own help...
#completed
The help is up. Find the option you need in it.
::

::task{name="used"}
#active
Waiting for a comma-separated count to ${N}...
#completed
That is the count, with commas between the numbers.
::

::tip
The `--help` output can scroll past quickly. **Shift-PgUp** scrolls the
terminal back in most terminals. A better tool for long text comes in
the next module.
::
