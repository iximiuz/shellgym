---
title: Keep it in one piece
tasks:
  quoted:
    timeout: 45
    check: |
      wait_exec --argc 2 '(^|/)date \+%A %d %B$'
    hint: |
      echo "Wrap the whole recipe - plus sign included - in quotes: a single argument with spaces inside."
    solve: |
      date '+%A %d %B'
---

`date` accepts a formatting *recipe* as an argument. A recipe starts
with `+` and mixes placeholders: `%A` is the weekday name, `%d` the day
of the month, `%B` the month name. This recipe prints something like
"Monday 28 July":

```
+%A %d %B
```

But there's a catch. The recipe must reach `date` as **one single
argument** - and it contains spaces, which the shell uses to *split*
arguments. Try it bare and `date` will complain about extra operands.

The fix: wrap the recipe in quotes. Quotes tell the shell "keep this
together, spaces and all" - and they are removed before the command
sees the argument.

::task
#active
Waiting for `date` to receive the three-part recipe as a single
argument...
#completed
The quotes never reached `date` - the shell consumed them, and `date`
received one argument with its spaces preserved. Quoting is how you
hand over any value that contains spaces.
::

::hint{title="date says: extra operand"}
That's the unquoted attempt: the shell split your recipe into three
separate arguments. Quote the whole thing - single or double quotes
both work here.
::
