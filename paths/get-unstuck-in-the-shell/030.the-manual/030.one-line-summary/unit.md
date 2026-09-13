---
title: Summaries and sections
vars:
  CMD: { pick: [sleep, hostname, tty] }
tasks:
  summarized:
    check: |
      wait_exec "(^|/)whatis ${CMD}\$"
    hint: |
      echo "Run: whatis ${CMD}"
    solve: |
      whatis $CMD
  sectioned:
    needs: [summarized]
    timeout: 60
    check: |
      wait_exec "(^|/)man 1 ${CMD}\$"
    hint: |
      echo "Put the section number between the command and the name: man, then 1, then ${CMD}."
    solve: |
      #!type man 1 $CMD
      #!keys enter
      #!wait 2
      #!keys q
      #!wait 1
---

Sometimes you only need a reminder of what a command is. The `whatis`
command prints just the one-line summary from the top of the manual page.

Ask it about `${CMD}`. You get several answers, each with a number in
parentheses. Those numbers are manual sections. The manual is split
into numbered parts: `1` holds user commands, `5` file formats, `8`
administration tools, and so on. The same name can have a different
page in each section.

Plain `man ${CMD}` opens the lowest-numbered section, which is the
command. To open a specific section, put its number before the name.
One day the default page will be the wrong one, and the section number
is how you pick the right page. Open the section 1 page for `${CMD}`
that way, then quit it.

::task{name="summarized"}
#active
Waiting for the one-line summary of `${CMD}`...
#completed
That is one line per section. Now open the section 1 page by number.
::

::task{name="sectioned"}
#active
Waiting for you to open the section 1 page for `${CMD}` (and `q` works
as always)...
#completed
That was the section 1 page, opened by its number.
::
