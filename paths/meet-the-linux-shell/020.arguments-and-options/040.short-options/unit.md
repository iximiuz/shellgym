---
title: "The time, twice"
tasks:
  local_time:
    check: |
      wait_exec '(^|/)date$' || exit 1
      set_var LOCAL_SEQ "$(event_seq)"
    hint: |
      echo "Run the bare command first: date with no arguments."
    solve: |
      date
  utc_time:
    needs: [local_time]
    check: |
      # only a date -u run after the plain date counts
      wait_exec --after "$LOCAL_SEQ" '(^|/)date -u$'
    hint: |
      echo "Run the same command, then a space, then the short option: a dash and the letter u."
    solve: |
      date -u
---

The `date` command prints the current date and time. Run it plain first.

Then run it again with the **option** `-u`: a dash and one letter,
which is called a **short option**. It switches the output to UTC, the shared reference time all machines on
the planet can agree on. Options are just arguments too. The dash is
what tells the command "this one changes my behavior".

Compare the two answers.

::task{name="local_time"}
#active
Waiting for the plain local time...
#completed
That is this machine's local time.
::

::task{name="utc_time"}
#active
Now the same clock in UTC, using the short option...
#completed
It is the same moment in a different presentation.
::
