---
title: The long way to say it
tasks:
  spelled_out:
    check: |
      wait_exec '(^|/)date --utc$'
    hint: |
      echo "Two dashes, then the word: --utc. One argument, no spaces inside it."
    solve: |
      date --utc
---

Most options have two spellings. You used the short one, `-u`. The same
option can be spelled out as a **long option**: two dashes and a word,
`--utc`.

Print the UTC time again, this time using the long spelling.

::task
#active
Waiting for the UTC time via the long option...
#completed
Identical output. Short options are quicker to type; long options are
easier to read later. In saved scripts, seasoned users prefer the long
form for exactly that reason.
::
