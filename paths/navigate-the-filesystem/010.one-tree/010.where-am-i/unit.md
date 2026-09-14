---
title: Where am I?
requires: [readline]
tasks:
  asked:
    check: |
      wait_line '^((/usr)?/bin/)?pwd *$'
    hint: |
      echo "The command is pwd. Type it and press Enter."
    solve: |
      pwd
---

Your shell is standing in some directory right now. Ask which one with
`pwd`, which stands for "print working directory".

The answer is a **path**: the full name of a directory, written as
every directory on the way down from the root, separated by slashes.
It starts with `/`, the root itself, and ends with the directory you
are in. Right now that is your **home directory**, `/home/` followed
by your user name. It is the place set aside for your own files.

::task
#active
Waiting for you to ask where you are...
#completed
That path is your home directory.
::

::tip
Many prompts show the working directory, or at least the last part of
its path. But the prompt looks different on every system, so `pwd` is
the one sure way to see where you are.
::
