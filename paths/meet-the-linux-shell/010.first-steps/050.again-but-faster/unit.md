---
title: Again - but faster
tasks:
  reran_user:
    check: |
      wait_exec '(^|/)whoami$'
    solve: |
      whoami
  reran_host:
    check: |
      wait_exec '(^|/)hostname$'
    solve: |
      hostname
  reran_tty:
    check: |
      wait_exec '(^|/)tty$'
    solve: |
      tty
---

A quick second lap. Run all three commands you've met - the one that
prints your user name, the one that prints the machine's name, and the
one that prints this terminal's name - but this time, **don't retype
them**.

Press the **Up** arrow until an earlier command reappears, adjust if
needed, and press Enter.

::task{name="reran_user"}
#active
Waiting for the user-name command again...
#completed
One down.
::

::task{name="reran_host"}
#active
Waiting for the machine-name command again...
#completed
Two down.
::

::task{name="reran_tty"}
#active
Waiting for the terminal-name command again...
#completed
All three - hopefully with far fewer keystrokes this time.
::

::tip
**Up**/**Down** walk through history. **Ctrl-A** jumps to the beginning
of the line, **Ctrl-E** to the end. And the **Tab** key completes
half-typed command names: try `hostn` + Tab.
::
