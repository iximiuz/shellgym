---
title: "The time, twice"
tasks:
  local_time:
    check: |
      wait_exec '(^|/)date$'
    hint: |
      echo "Just the bare command first: date, no arguments."
    solve: |
      date
  utc_time:
    needs: [local_time]
    check: |
      wait_exec '(^|/)date -u$'
    hint: |
      echo "Same command, then a space, then the short option: a dash and the letter u."
    solve: |
      date -u
---

`date` prints the current date and time. Run it plain first.

Then run it again with the **option** `-u`: a dash and one letter. It
switches the output to UTC, the shared reference time all machines on
the planet can agree on. Options are just arguments too - the dash is
what tells the command "this one changes my behavior".

Compare the two answers.

::task{name="local_time"}
#active
Waiting for the plain local time...
#completed
That's this machine's local time.
::

::task{name="utc_time"}
#active
Now the same clock in UTC, using the short option...
#completed
Same moment, different presentation - one single-letter option made
the difference. Single-dash-single-letter options are called **short
options**.
::
