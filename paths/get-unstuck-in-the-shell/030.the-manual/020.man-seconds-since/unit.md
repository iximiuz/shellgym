---
title: Search inside a page
variant: man1=seconds-since
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
      wait_exec '(^|/)date \+%s$'
    hint: |
      echo "Inside the page, type /seconds since and press Enter. The match names a two-character format sequence. Quit with q and use that sequence after a + sign."
    solve: |
      #!type /seconds since
      #!keys enter
      #!wait 1
      #!keys q
      #!wait 1
      date +%s
---

In [Meet the Linux Shell](https://labs.iximiuz.com/shell-gyms/meet-the-linux-shell) you used formatting recipes with `date`, such
as `+%A` for the weekday and `%B` for the month. There is also a
sequence that prints the number of seconds since 1970-01-01, the
timestamp format used all over computing. The manual page for `date` lists every sequence.

The pager can search. Press `/`, type the text to look for, and press
Enter.

Open `man date`, search for `seconds since`, and find the sequence. Then
quit the pager and run `date` with that recipe.

::task{name="opened"}
#active
Waiting for `man date` to open...
#completed
The page is open. Now search it, because the list of sequences is long.
::

::task{name="applied"}
#active
Waiting for `date` to print the seconds-since-1970 number...
#completed
That number is the Unix timestamp.
::

::tip
Inside the pager, `/some-text` searches forward from the current position. **n** jumps to the
next match and **N** to the previous one. The search wraps around the
page, so it does not matter where you start.
::
