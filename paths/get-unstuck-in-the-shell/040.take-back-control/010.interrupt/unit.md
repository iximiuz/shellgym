---
title: Interrupt a command
vars:
  NAP: { pick: ["300", "400", "500"] }
tasks:
  launched:
    check: |
      wait_exec "(^|/)sleep ${NAP}s?\$"
    hint: |
      echo "Start the long sleep first: sleep ${NAP}"
    solve: |
      #!type sleep $NAP
      #!keys enter
      #!wait 1
  interrupted:
    needs: [launched]
    timeout: 60
    check: |
      wait_proc --timeout 15 "^(/usr/bin/)?slee[p] ${NAP}s?\$" || true
      wait_proc_gone "^(/usr/bin/)?slee[p] ${NAP}s?\$"
    hint: |
      echo "Do not wait the ${NAP} seconds out. Hold Ctrl and press C to interrupt the sleep."
    solve: |
      #!keys C-c
      #!wait 1
---

You know `sleep` from Meet the Linux Shell, and you know that the
prompt only returns when the command finishes. There is another way
out.

Start a sleep that you have no intention of waiting for:

```
sleep ${NAP}
```

That is ${NAP} seconds of nothing. Press **Ctrl-C** instead of
waiting. It interrupts the foreground command, which asks it to stop
right now, and gives you the prompt back.

::task{name="launched"}
#active
Waiting for the ${NAP}-second sleep to start...
#completed
The prompt is gone for the next several minutes, unless you act.
::

::task{name="interrupted"}
#active
Waiting for you to interrupt it with **Ctrl-C**...
#completed
The prompt is back.
::

::tip{title="What Ctrl-C leaves alone"}
Only the running _foreground_ command stops. Your terminal and your shell are fine,
and the `^C` on the screen marks where you pressed the key.
::
