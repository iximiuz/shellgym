---
title: When the name escapes you
variant: apropos=shuf
tasks:
  searched:
    check: |
      wait_exec '(?i)(^|/)apropos .*(random|permutation)'
    hint: |
      echo "Search the descriptions with apropos followed by the phrase. Quote the phrase, since it contains a space."
    solve: |
      apropos 'random permutation'
  ran:
    needs: [searched]
    timeout: 60
    check: |
      wait_exec '(^|/)shuf .*(-i|--input-range)'
    hint: |
      echo "In each result line, the command name is on the left of the dash. Spot the one about permutations, and run it with -i 1-10."
    solve: |
      shuf -i 1-10
---

The hardest command to look up is the one whose name you do not know.
You remember what it does. There is a command that shuffles its input into a random order, which the manual calls a random permutation. You do not
remember what it is called.

The `apropos` command searches the one-line descriptions of every manual
page and prints the ones that match. Search for the phrase `random permutation`. The words must reach `apropos` as one argument, so remember your quoting.

Then run the command you discovered with the option `-i 1-10`, which makes it shuffle the numbers from 1 to 10.

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
Those are the numbers from 1 to 10 in a random order.
::

::hint{title="Too many or zero results?"}
The `apropos` command matches words in the descriptions, so pick words that a
description would use. Zero results on a fresh machine can also mean
that the manual index is missing, which is not the case in this gym.
::
