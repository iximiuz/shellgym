---
title: Wait for it
vars:
  NAP: { pick: ["4", "5", "6"] }
tasks:
  napped:
    check: |
      wait_exec "(^|/)sleep ${NAP}s?\$"
    hint: |
      echo "Run sleep with the number ${NAP} as its argument, separated by a space."
    solve: |
      sleep $NAP
---

Not every command answers instantly. Some work for a while - and the
prompt only comes back when they are **done**.

The command `sleep` does nothing on purpose: you give it a number of
seconds, and it simply takes that long. Make this shell sleep for
**${NAP} seconds**. The number is an *argument* - an extra word after
the command name, separated by a space.

::task
#active
Waiting for a ${NAP}-second nap...
#completed
Did you notice? No output at all, just the prompt returning after
${NAP} seconds. Silence does not mean failure - many commands say
nothing when all went well.
::

::hint{title="Nothing seems to happen?"}
That is the point: while `sleep` runs, the prompt is gone and the shell
is busy. Count to ${NAP} - the prompt will be back.
::
