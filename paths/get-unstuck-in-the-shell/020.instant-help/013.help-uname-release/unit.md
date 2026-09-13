---
title: The kernel release
variant: help1=uname-release
tasks:
  consulted:
    check: |
      wait_exec '(^|/)uname --help$'
    hint: |
      echo "Run: uname --help"
    solve: |
      uname --help
  used:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)uname (-r|--kernel-release)\$"
    hint: |
      echo "Look for the line about the kernel release. The option is a single lowercase letter."
    solve: |
      uname -r
---

The `uname` command prints the name of the operating system, which is `Linux`
here. It can also print details about the running system, such as the
version of the kernel, and the options for that are in its help:

```
uname --help
```

Find the option that prints only the kernel release. Then make `uname`
print it. You will need that number whenever you report a problem.

::task{name="consulted"}
#active
Waiting for you to ask `uname` for its own help...
#completed
The help is up. Find the option you need in it.
::

::task{name="used"}
#active
Waiting for `uname` to print the kernel release...
#completed
That is the kernel release.
::

::tip
The `--help` output can scroll past quickly. **Shift-PgUp** scrolls the
terminal back in most terminals. A better tool for long text comes in
the next module.
::
