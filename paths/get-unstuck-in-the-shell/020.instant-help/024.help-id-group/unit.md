---
title: The name of your group
variant: help2=id-group
tasks:
  consulted:
    check: |
      wait_exec '(^|/)id --help$'
    hint: |
      echo "Start with: id --help"
    solve: |
      id --help
  applied:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)id (-gn|-ng|-g -n|-n -g|--group --name|--name --group)\$"
    hint: |
      echo "You need two options together: the one for the effective group and the one that prints a name instead of a number. Two short options can be combined into one word."
    solve: |
      id -gn
---

The `id` command prints your identity on this machine: your user name, your
primary group, and the other groups you belong to, each with a number
next to it. The `id --help` output lists options that print only one of these, and
an option that prints a name instead of a number.

Look them up, then make `id` print only the **name** of your primary
group. You need two of the options, and two short options can share
one dash.

::task{name="consulted"}
#active
Waiting for you to open `id --help`...
#completed
The option is in the list. Find the line that describes what you need.
::

::task{name="applied"}
#active
Waiting for `id` to print the name of your primary group...
#completed
That is the name of your primary group, and nothing else.
::
