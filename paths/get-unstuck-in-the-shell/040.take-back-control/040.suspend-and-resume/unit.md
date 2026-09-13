---
title: Suspend and resume
vars:
  PAGE: { pick: [seq, date, uptime] }
tasks:
  reading:
    check: |
      wait_exec "(^|/)man ${PAGE}\$"
    hint: |
      echo "Open the manual: man ${PAGE}"
    solve: |
      #!type man $PAGE
      #!keys enter
      #!wait 2
  suspended:
    needs: [reading]
    timeout: 60
    check: |
      wait_proc_state "^(/usr/bin/)?ma[n] ${PAGE}\$" T
    hint: |
      echo "While the page is open, press Ctrl-Z, the letter Z. The shell reports Stopped and returns your prompt."
    solve: |
      #!keys C-z
      #!wait 1
  tried:
    needs: [suspended]
    timeout: 60
    check: |
      wait_exec "(^|/)${PAGE} (.* )?-[A-Za-z-]"
    hint: |
      echo "You have a prompt while the manual waits. Run ${PAGE} with any option from the page, for example the first one listed under OPTIONS."
    solve: |
      case $PAGE in seq) seq -w 5 ;; date) date -u ;; *) uptime -p ;; esac
  resumed_and_closed:
    needs: [tried]
    timeout: 60
    check: |
      wait_proc --timeout 15 "^(/usr/bin/)?ma[n] ${PAGE}\$" || true
      wait_proc_gone "^(/usr/bin/)?ma[n] ${PAGE}\$"
    hint: |
      echo "Bring the manual back with fg. You land exactly where you left it. Then quit it properly with q."
    solve: |
      #!type fg
      #!keys enter
      #!wait 1
      #!keys q
      #!wait 1
---

Pressing **Ctrl-C** terminates the command, and if you need it again you
will have to start it over. Sometimes you want it back exactly where it was before interruption. For instance, you
are reading a manual page and want to try an option you just read
about, without losing your place.

**Ctrl-Z** suspends the foreground command and keeps it around, still alive, while you get the prompt.
The command `fg` brings it back.

Practice the full cycle on a manual page:

1. Open `man ${PAGE}` and find its options.
2. Press **Ctrl-Z**. The shell prints `Stopped` and the prompt returns.
3. Run `${PAGE}` with any option from the page.
4. Bring the manual back with `fg`. You land on the exact spot you
   left. Then quit it properly with `q`.

::task{name="reading"}
#active
Waiting for `man ${PAGE}` to open...
#completed
You are in the middle of the page, and now you want to try something
from it.
::

::task{name="suspended"}
#active
Waiting for the pager to be suspended with Ctrl-Z...
#completed
The pager is suspended, and your prompt is back.
::

::task{name="tried"}
#active
Waiting for `${PAGE}` to run with an option at the reclaimed prompt...
#completed
The option worked, and the manual page waited the whole time.
::

::task{name="resumed_and_closed"}
#active
Waiting for `fg` and then a proper `q`...
#completed
The page came back where you left it, and now it is closed for good.
::

::hint{title="Ctrl-Z did nothing?"}
Ctrl-Z only affects the foreground command, the one that currently
holds your prompt. If you can see the prompt, there is nothing to
suspend.
::
