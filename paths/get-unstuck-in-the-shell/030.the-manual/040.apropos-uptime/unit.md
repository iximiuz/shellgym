---
title: When the name escapes you
variant: apropos=uptime
tasks:
  searched:
    check: |
      ARGV=$(wait_exec --latest '(^|/)apropos .+') || exit 1
      WORDS=${ARGV#* }
      if [[ "$WORDS" == *" "* ]] && [ "$(wait_exec --now --latest --argc 2 '(^|/)apropos .+')" != "$ARGV" ]; then
        hint_exit "The words reached apropos as separate arguments, so it searched for each word on its own. Quote the phrase to search for it as a whole."
      fi
      FOUND=$(apropos -- "$WORDS" 2>/dev/null)
      grep -q '^uptime ' <<< "$FOUND" || hint_exit "The command you are after is not among the results of that search. Search for the phrase how long."
      set_var SEARCHED_SEQ "$(event_seq)"
    hint: |
      echo "Search the descriptions with apropos followed by the phrase. Quote the phrase, since it contains a space."
    solve: |
      apropos 'how long'
  ran:
    needs: [searched]
    timeout: 60
    check: |
      wait_exec --after "$SEARCHED_SEQ" '(^|/)uptime$'
    hint: |
      echo "In each result line, the command name is on the left of the dash. Spot the one about the system running, and run it."
    solve: |
      uptime
---

The hardest command to look up is the one whose name you do not know.
You remember what it does. There is a command that tells how long the machine has been running. You do not
remember what it is called.

The `apropos` command searches the one-line descriptions of every manual
page and prints the ones that match. Search for the phrase `how long`. The words must reach `apropos` as one argument, so remember your quoting.

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
That is how long the machine has been up.
::

::hint{title="Too many or zero results?"}
The `apropos` command matches words in the descriptions, so pick words that a
description would use. Zero results on a fresh machine can also mean
that the manual index is missing, which is not the case in this gym.
::
