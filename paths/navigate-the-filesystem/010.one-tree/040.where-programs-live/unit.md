---
title: Where the programs live
vars:
  CMD: { pick: [date, seq, sleep, whoami] }
tasks:
  arrived:
    check: |
      wait_cwd /usr/bin
    hint: |
      echo "Ask type ${CMD} for the full path of the program. Everything before the last slash is the directory that holds it."
    solve: |
      type $CMD
      cd /usr/bin
---

In [Get Unstuck in the Shell](https://labs.iximiuz.com/shell-gyms/get-unstuck-in-the-shell) you learned that `type ${CMD}` prints the
full path of the program file that runs when you type `${CMD}`.

A path to a file works like a path to a directory: the last name is
the file, and everything before it is the directory the file sits in.
Look up `${CMD}` with `type`, read the directory from its path, and
move your shell into that directory.

::task
#active
Waiting for your shell to arrive in the directory that holds `${CMD}`...
#completed
This is where most of the programs on the machine live.
::

::tip
Run `ls` here. Every name in that long list is a command you can
type. Some of them you already know.
::
