---
title: The top of the tree
tasks:
  arrived:
    check: |
      wait_cwd /
    hint: |
      echo "Run cd with / as its argument."
    solve: |
      cd /
  looked:
    needs: [arrived]
    timeout: 60
    check: |
      wait_exec --cwd / '(^|/)ls( .*)?$'
    hint: |
      echo "Run ls with no arguments to list the names in the current directory."
    solve: |
      ls
---

The `cd` command ("change directory") moves your shell. It takes one
argument: the path of the directory to go to.

Move to the root directory. Its path is `/`, which is the top of the
filesystem tree.

Then look around with `ls`, which lists the names of everything in the
current directory. You will meet `ls` properly in the next gym. For
now it is your window into whatever directory you are standing in.

::task{name="arrived"}
#active
Waiting for your shell to arrive at `/`...
#completed
You are at the top of the tree. Have a look around.
::

::task{name="looked"}
#active
Waiting for you to list what is in the root directory...
#completed
Those names are the top-level directories of this machine.
::
