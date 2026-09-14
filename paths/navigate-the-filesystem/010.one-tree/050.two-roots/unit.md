---
title: Two things called root
requires: [readline]
tasks:
  denied:
    check: |
      wait_line '^cd +/root/? *$'
    hint: |
      echo "Try to enter /root with cd. It will not work, and the error message is the point."
    solve: |
      cd /root
  arrived:
    needs: [denied]
    timeout: 60
    check: |
      wait_cwd /
    hint: |
      echo "The root directory is a lone slash. Move there with cd."
    solve: |
      cd /
---

The word "root" means two different things on a Linux machine, and
both show up in paths.

The **root directory** is `/`, the top of the tree. The **root user**
is the administrator account, and its home directory is `/root`. That
directory is private to the administrator.

See the difference for yourself. First, try to move into `/root` and
read what the shell says. Then move to the root directory `/`, which
is open to everyone.

::task{name="denied"}
#active
Waiting for an attempt to enter `/root`...
#completed
The shell refused because the administrator's files are off limits to you. Now move to `/`.
::

::task{name="arrived"}
#active
Waiting for your shell to arrive at `/`...
#completed
The root directory is open to everyone.
::

::tip{title="Read the error"}
A failed `cd` leaves your shell exactly where it was. The message
names the directory it could not enter and the reason, and the exit
status in `$?` is non-zero.
::
