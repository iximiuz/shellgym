---
title: When the name escapes you
variant: apropos=free
tasks:
  searched:
    check: |
      wait_exec '(?i)(^|/)apropos .*(free|memory)'
    hint: |
      echo "Search the descriptions with apropos followed by the phrase. Quote the phrase, since it contains a space."
    solve: |
      apropos 'free and used memory'
  ran:
    needs: [searched]
    timeout: 60
    check: |
      wait_exec '(^|/)free( .*)?$'
    hint: |
      echo "In each result line, the command name is on the left of the dash. Spot the one about memory, and run it."
    solve: |
      free
---

The hardest command to look up is the one whose name you do not know.
You remember what it does. There is a command that shows how much memory is free and how much is in use. You do not
remember what it is called.

The `apropos` command searches the one-line descriptions of every manual
page and prints the ones that match. Search for the phrase `free and used memory`. The words must reach `apropos` as one argument, so remember your quoting.

Then run the command you discovered.

::task{name="searched"}
#active
Waiting for you to search the manual descriptions...
#completed
One of the matches is clearly the command you are after.
::

::task{name="ran"}
#active
Waiting for you to run the command you just discovered...
#completed
That is the memory table, and reading it is a topic for later.
::

::hint{title="Too many or zero results?"}
The `apropos` command matches words in the descriptions, so pick words that a
description would use. Zero results on a fresh machine can also mean
that the manual index is missing, which is not the case in this gym.
::
