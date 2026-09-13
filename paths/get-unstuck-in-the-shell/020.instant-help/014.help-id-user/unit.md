---
title: Your numeric user ID
variant: help1=id-user
tasks:
  consulted:
    check: |
      wait_exec '(^|/)id --help$'
    hint: |
      echo "Run: id --help"
    solve: |
      id --help
  used:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)id (-u|--user)\$"
    hint: |
      echo "Look for the line about the effective user ID. The option is a single lowercase letter."
    solve: |
      id -u
---

The `id` command prints your identity on this machine: your user name and the
groups you belong to, each with a number next to it. The number is how
the system knows you internally, and names are for people. Scripts
often need only the number, and `id` has an option for that. Ask:

```
id --help
```

Find the option that prints only your numeric user ID. Then make `id`
print it.

::task{name="consulted"}
#active
Waiting for you to ask `id` for its own help...
#completed
The help is up. Find the option you need in it.
::

::task{name="used"}
#active
Waiting for `id` to print your numeric user ID...
#completed
That number is your user ID.
::

::tip
The `--help` output can scroll past quickly. **Shift-PgUp** scrolls the
terminal back in most terminals. A better tool for long text comes in
the next module.
::
