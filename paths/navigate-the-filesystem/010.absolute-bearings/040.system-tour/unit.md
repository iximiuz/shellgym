---
title: A short system tour
vars:
  STOP2: { pick: [/usr/local, /usr/share, /var/tmp] }
tasks:
  visit_etc:
    check: |
      wait_cwd /etc
    hint: |
      echo "First stop: /etc - the classic home of system configuration. One absolute cd away."
    solve: |
      cd /etc
  visit_next:
    needs: [visit_etc]
    check: |
      wait_cwd "$STOP2"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. Second stop: ${STOP2}. Same move, different address."
    solve: |
      cd $STOP2
---

Real systems have a standard layout, and knowing a few landmarks pays
off daily. Take a two-stop tour - these are real system directories,
you are only visiting.

First stop: `/etc`, where system configuration lives.

::task{name="visit_etc"}
#active
Waiting for your shell in `/etc`...
#completed
This is where the system keeps its settings - you will come back here
often.
::

Second stop: `${STOP2}`.

::task{name="visit_next"}
#active
Waiting for your shell in `${STOP2}`...
#completed
Tour complete. Absolute paths took you across the system in two moves -
no map needed, just addresses.
::

::tip
---
title: Retracing your steps
---
Press the Up arrow to bring back your previous command, edit it with
the Left/Right arrows, and press Enter. Most navigation is a recycled
command with a small change.
::
