---
title: The date for an email header
variant: help2=date-email
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
      wait_exec "(^|/)date (-R|--rfc-email)\$"
    hint: |
      echo "Look for the line that mentions the format used in email, RFC 5322. The short option is a single capital letter."
    solve: |
      date -R
---

The `date` command can also print the current moment in the exact format
that email headers use, like `Fri, 11 Sep 2026 17:50:43 +0000`. The
`date --help` output lists an option for it.

Look the option up, then make `date` print the current date and time
in that email format.

::task{name="consulted"}
#active
Waiting for you to open `date --help`...
#completed
The option is in the list. Find the line that describes what you need.
::

::task{name="applied"}
#active
Waiting for the date in the email header format...
#completed
That is the date in the email header format.
::
