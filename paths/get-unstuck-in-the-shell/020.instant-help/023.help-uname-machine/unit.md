---
title: The machine hardware name
variant: help2=uname-machine
tasks:
  consulted:
    check: |
      wait_exec '(^|/)uname --help$'
    hint: |
      echo "Start with: uname --help"
    solve: |
      uname --help
  applied:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)uname (-m|--machine)\$"
    hint: |
      echo "Look for the line about the machine hardware name. The option is a single lowercase letter."
    solve: |
      uname -m
---

The `uname` command prints the name of the operating system. It can also print the
machine hardware name, such as `x86_64` or `aarch64`, which tells you
which builds of a program can run here. Download pages call it the
architecture. The `uname --help` output lists an option for it.

Look the option up, then make `uname` print the machine hardware name.

::task{name="consulted"}
#active
Waiting for you to open `uname --help`...
#completed
The option is in the list. Find the line that describes what you need.
::

::task{name="applied"}
#active
Waiting for `uname` to print the machine hardware name...
#completed
That is the machine hardware name.
::
