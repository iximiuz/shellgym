---
title: Numbers of equal width
variant: help1=seq-width
vars:
  N: { pick: ["12", "15", "20"] }
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
      wait_exec "(^|/)seq (-w|--equal-width) ${N}\$"
    hint: |
      echo "The option takes no value: seq, the equal-width option, then ${N}."
    solve: |
      seq -w $N
---

By default, `seq` prints one number per line, and each number is only as wide as
it needs to be: `1`, `2`, and later `10`, `11`. It can also pad the
short ones with leading zeros so that every number has the same width,
and the option for that is in its help:

```
seq --help
```

Find the option that equalizes the width. Then make `seq` count from 1
to **${N}** with the numbers padded, like `01`, `02`, and so on.

::task{name="consulted"}
#active
Waiting for you to ask `seq` for its own help...
#completed
The help is up. Find the option you need in it.
::

::task{name="used"}
#active
Waiting for a zero-padded count to ${N}...
#completed
The numbers line up now.
::

::tip
The `--help` output can scroll past quickly. **Shift-PgUp** scrolls the
terminal back in most terminals. A better tool for long text comes in
the next module.
::
