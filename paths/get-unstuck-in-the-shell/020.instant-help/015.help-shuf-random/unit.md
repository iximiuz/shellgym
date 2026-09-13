---
title: A random number
variant: help1=shuf-random
vars:
  N: { pick: ["10", "20", "100"] }
tasks:
  consulted:
    check: |
      wait_exec '(^|/)shuf --help$'
    hint: |
      echo "Run: shuf --help"
    solve: |
      shuf --help
  used:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)shuf ((-i|--input-range=?) ?1-${N} (-n|--head-count=?) ?1|(-n|--head-count=?) ?1 (-i|--input-range=?) ?1-${N})\$"
    hint: |
      echo "Two options: the input range, given as 1-${N}, and the head count, given as 1. Their order does not matter."
    solve: |
      shuf -i 1-$N -n 1
---

The `shuf` command shuffles the lines it is given and prints them in random
order. It can also generate the numbers itself and limit how many lines
it prints, and both options are in its help:

```
shuf --help
```

Find the option that takes a range of numbers as input and the option
that limits the output. Then make `shuf` print **one** random number
between 1 and **${N}**.

::task{name="consulted"}
#active
Waiting for you to ask `shuf` for its own help...
#completed
The help is up. Find the option you need in it.
::

::task{name="used"}
#active
Waiting for one random number between 1 and ${N}...
#completed
That is one random number, and it changes on every run.
::

::tip
The `--help` output can scroll past quickly. **Shift-PgUp** scrolls the
terminal back in most terminals. A better tool for long text comes in
the next module.
::
