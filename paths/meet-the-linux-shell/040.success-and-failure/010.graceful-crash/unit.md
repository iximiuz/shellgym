---
title: Break something (safely)
vars:
  BOGUS: { pick: [pumpkin, sideways, banana, upstairs] }
tasks:
  broke_it:
    check: |
      wait_exec "(^|/)date --${BOGUS}\$"
    hint: |
      echo "Give date the long option --${BOGUS} - yes, really. It's guaranteed not to exist."
    solve: |
      date --$BOGUS
---

Sooner or later every command you type will fail. Let's get the first
failure over with, on purpose: run `date` with the made-up option
`--${BOGUS}`.

Read the complaint carefully - error messages follow a pattern worth
recognizing: *who* is complaining (`date:`), *what* it didn't like, and
usually a pointer toward help.

By the way: since anything starting with a dash looks like an option,
commands accept a lone `--` argument meaning "options end here -
whatever follows is a plain value, dashes and all". You won't need it
today, but you'll be glad it exists the day a value of yours starts
with a `-`.

::task
#active
Waiting for you to feed `date` the nonsense option `--${BOGUS}`...
#completed
The command refused, explained itself, and the prompt came right back.
An error message is a conversation, not a catastrophe. (Mistyping a
command's *name* is equally harmless - the shell itself answers
`command not found`.)
::

::tip
Changed your mind halfway through typing a line? **Ctrl-C** abandons
the line and gives you a fresh prompt - nothing gets run.
::
