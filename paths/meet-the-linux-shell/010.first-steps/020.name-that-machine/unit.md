---
title: Name that machine
tasks:
  asked:
    check: |
      wait_exec '(^|/)hostname$'
    hint: |
      echo "The command is a single word: hostname. Type it and press Enter."
    solve: |
      hostname
---

Every Linux machine has a name of its own, so that people (and other
machines) can tell them apart. The command that prints it is `hostname`.

Ask this machine for its name.

::task
#active
Waiting for you to ask the machine for its name...
#completed
One machine, one name. Many prompts show both at once as
`user@machine` - now you know where each half comes from.
::

::tip
Press the **Up** arrow to bring back your previous command. Command
history saves a lot of typing - you'll lean on it constantly.
::
