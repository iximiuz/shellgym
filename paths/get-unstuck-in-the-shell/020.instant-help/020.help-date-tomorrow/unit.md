---
title: Tomorrow's date today
variant: help2=date-tomorrow
tasks:
  consulted:
    check: |
      wait_exec '(^|/)date --help$'
    hint: |
      echo "Start with: date --help"
    solve: |
      date --help
  applied:
    needs: [consulted]
    timeout: 60
    check: |
      wait_exec "(^|/)date (-d|--date=?) ?tomorrow\$"
    hint: |
      echo "Near the top of the help there is an option to display a time other than now. Its value is the word tomorrow."
    solve: |
      date -d tomorrow
---

The `date` command prints the current moment, and it can print other moments too.
Near the top of `date --help` there is an option that takes a
description of a time and displays that time instead of now.

Look the option up, then make `date` print tomorrow's date. The
description to give it is the single word `tomorrow`.

::task{name="consulted"}
#active
Waiting for you to open `date --help`...
#completed
The option is in the list. Find the line that describes what you need.
::

::task{name="applied"}
#active
Waiting for tomorrow's date...
#completed
That is tomorrow's date.
::

::tip{title="More than tomorrow"}
The `date` command understands a small language of time descriptions, like
`next friday` and `2 hours ago`. You do not need to learn it now.
Knowing that it exists, and where you found the option, is enough.
::
