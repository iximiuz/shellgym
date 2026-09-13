---
title: Search inside a page
variant: man1=day-of-week
tasks:
  opened:
    check: |
      wait_exec '(^|/)man date$'
    hint: |
      echo "Open the manual first: man date"
    solve: |
      #!type man date
      #!keys enter
      #!wait 2
  applied:
    needs: [opened]
    timeout: 90
    check: |
      wait_exec '(^|/)date \+%[uw]$'
    hint: |
      echo "Inside the page, type /day of week and press Enter. The match names a two-character format sequence. Quit with q and use that sequence after a + sign."
    solve: |
      #!type /day of week
      #!keys enter
      #!wait 1
      #!keys q
      #!wait 1
      date +%u
---

In [Meet the Linux Shell](https://labs.iximiuz.com/shell-gyms/meet-the-linux-shell) you used formatting recipes with `date`, such
as `+%A` for the weekday and `%B` for the month. There is also a
sequence that prints the day of the week as a single digit, with
Monday as 1. The manual page for `date` lists every sequence.

The pager can search. Press `/`, type the text to look for, and press
Enter.

Open `man date`, search for `day of week`, and find the sequence. Then
quit the pager and run `date` with that recipe.

::task{name="opened"}
#active
Waiting for `man date` to open...
#completed
The page is open. Now search it, because the list of sequences is long.
::

::task{name="applied"}
#active
Waiting for `date` to print the day of the week as a digit...
#completed
That is the day of the week as a digit.
::

::tip
Inside the pager, `/some-text` searches forward from the current position. **n** jumps to the
next match and **N** to the previous one. The search wraps around the
page, so it does not matter where you start.
::
