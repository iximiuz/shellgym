---
title: The status of an interrupted command
requires: [readline]
vars:
  NAP: { pick: ["350", "450", "550"] }
tasks:
  interrupted:
    timeout: 60
    check: |
      wait_exec "(^|/)sleep ${NAP}s?\$"
      wait_proc --timeout 15 "^(/usr/bin/)?slee[p] ${NAP}s?\$" || true
      wait_proc_gone "^(/usr/bin/)?slee[p] ${NAP}s?\$"
    hint: |
      echo "Start sleep ${NAP}, then press Ctrl-C."
    solve: |
      #!type sleep $NAP
      #!keys enter
      #!wait 1
      #!keys C-c
      #!wait 1
  revealed:
    needs: [interrupted]
    timeout: 45
    check: |
      wait_line '^echo +("\$\?"|\$\?|\$\{\?\}) *$'
    hint: |
      echo "Right after the ^C, print the status of the last command: echo \$?"
    solve: |
      echo $?
---

An interrupted command still leaves an exit status behind, and for
most commands it is a non-zero one.

1. Start `sleep ${NAP}` and interrupt it with **Ctrl-C**.
2. Reveal the status it left right away: `echo $?`

::task{name="interrupted"}
#active
Waiting for a `sleep ${NAP}` cut short by Ctrl-C...
#completed
The sleep was interrupted. Now check what it left in `$?`.
::

::task{name="revealed"}
#active
Waiting for you to print `$?`...
#completed
That is the status the interrupted `sleep` left behind.
::

::tip{title="Not every command dies on Ctrl-C"}
**Ctrl-C** only sends a signal. A program can catch it, clean up, and
exit with a status of its own choosing, even `0`. The `130` you see
here is what you get when the command does nothing special and the
signal ends it.
::
