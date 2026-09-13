---
title: Your first man page
vars:
  PAGE: { pick: [seq, sleep, tty] }
tasks:
  opened:
    check: |
      wait_exec "(^|/)man ${PAGE}\$"
    hint: |
      echo "Run: man ${PAGE}"
    solve: |
      #!type man $PAGE
      #!keys enter
      #!wait 2
  closed:
    needs: [opened]
    timeout: 60
    check: |
      wait_proc --timeout 20 "^(/usr/bin/)?ma[n] ${PAGE}\$" || true
      wait_proc_gone "^(/usr/bin/)?ma[n] ${PAGE}\$"
    hint: |
      echo "Still inside the page? A single q quits the pager and brings your prompt back."
    solve: |
      #!keys q
      #!wait 1
---

Open the manual page for `${PAGE}`:

```
man ${PAGE}
```

Your prompt is gone. You are inside the pager now, and it answers to
keystrokes. Have a look around. Every manual page follows the same
layout: NAME, SYNOPSIS, DESCRIPTION, and so on.

Then quit the pager and return to your shell.

::task{name="opened"}
#active
Waiting for you to open the manual page for `${PAGE}`...
#completed
You are in the pager now. The keys in the tip below are your controls.
::

::task{name="closed"}
#active
Waiting for you to quit back to the prompt (one key does it)...
#completed
The prompt is back.
::

::tip{title="Pager keys"}
**Space** and **b** move a screen down and up. The **arrow keys** move
one line at a time. **g** and **G** jump to the beginning and the end.
**q** quits.
::

::tip{title="Worth learning today"}
You will open, read, and quit pages like this a lot, so
the way out is worth learning today.
::
