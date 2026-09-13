---
title: The date in ISO form
variant: help2=date-iso
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
      wait_exec "(^|/)date (-I(date)?|--iso-8601(=date)?)\$"
    hint: |
      echo "Look for the line that mentions ISO 8601. The short option is a single capital letter."
    solve: |
      date -I
---

Plain `date` prints the current moment in a long, human-friendly form. Many
tools and file names want the short international form instead, with
the year first, like `2026-09-11`. The `date --help` output lists an
option that prints exactly that form.

Look the option up, then make `date` print today's date in the ISO
8601 format.

::task{name="consulted"}
#active
Waiting for you to open `date --help`...
#completed
The option is in the list. Find the line that describes what you need.
::

::task{name="applied"}
#active
Waiting for today's date in the year-month-day form...
#completed
That is today's date in the ISO 8601 form.
::

::tip{title="Why this form?"}
Dates written year first sort correctly as plain text, and they leave
no doubt about which part is the month.
::
