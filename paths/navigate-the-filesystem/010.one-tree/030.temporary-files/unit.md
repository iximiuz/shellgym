---
title: The scratch space
tasks:
  arrived:
    check: |
      wait_cwd /tmp
    hint: |
      echo "The directory is /tmp, directly under the root."
    solve: |
      cd /tmp
---

Every Linux machine has a `/tmp` directory for files that nobody
needs to keep. Programs drop their scratch files there, and so can you.
Anything in it may be removed when the machine restarts, which is the
whole point.

Move your shell there, then take a look with `ls` if you like.

::task
#active
Waiting for your shell to arrive in `/tmp`...
#completed
You are in the scratch space of this machine.
::

::tip
If your prompt shows the working directory, it changed when you moved.
A glance at it tells you roughly where you are, and `pwd` is the way
to make sure.
::
