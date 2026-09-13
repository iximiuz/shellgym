---
title: The machine's address
variant: help1=hostname-address
tasks:
  consulted:
    check: |
      wait_exec '(^|/)hostname --help$'
    hint: |
      echo "Run: hostname --help"
    solve: |
      hostname --help
  used:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)hostname (-I|-i|--all-ip-addresses|--ip-address)\$"
    hint: |
      echo "Look for the line about all addresses for the host. The option is a single capital letter."
    solve: |
      hostname -I
---

The `hostname` command prints the name of this machine. The same program can also
print the machine's network addresses, the ones other machines use to
reach this one, and the option for that is in its help:

```
hostname --help
```

Find the option that prints all addresses of the host. Then make
`hostname` print them.

::task{name="consulted"}
#active
Waiting for you to ask `hostname` for its own help...
#completed
The help is up. Find the option you need in it.
::

::task{name="used"}
#active
Waiting for `hostname` to print the machine's addresses...
#completed
Those are the machine's addresses.
::

::tip
The `--help` output can scroll past quickly. **Shift-PgUp** scrolls the
terminal back in most terminals. A better tool for long text comes in
the next module.
::
