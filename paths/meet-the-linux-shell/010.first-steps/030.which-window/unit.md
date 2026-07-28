---
title: Which window is this?
tasks:
  asked:
    check: |
      wait_exec '(^|/)tty$'
    hint: |
      echo "Three letters: tty. It prints something like /dev/pts/0."
    solve: |
      tty
---

You can open several terminal windows to the same machine at once, and
each one is a separate conversation. The command `tty` prints the name of
the terminal *this* conversation is running on.

Run it here.

::task
#active
Waiting for you to ask which terminal this is...
#completed
Something like `/dev/pts/0` - that's this window's own address. A second
terminal window would report a different one.
::

::tip
If the screen gets cluttered, **Ctrl-L** wipes it clean. Your command
history stays intact - only the pixels are cleared.
::
