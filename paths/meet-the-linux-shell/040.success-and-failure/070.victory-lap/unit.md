---
title: Victory lap
vars:
  BOGUS: { pick: [teapot, polkadot, upside-down] }
tasks:
  punctual:
    timeout: 45
    check: |
      wait_exec '(^|/)sleep 2s?$'
      wait_exec '(^|/)date (-u|--utc)$'
    hint: |
      echo "Line one: a 2-second sleep, && , then date with its UTC option (short or long)."
    solve: |
      sleep 2 && date -u
  resilient:
    timeout: 45
    check: |
      wait_exec "(^|/)date --${BOGUS}\$"
      wait_exec '(^|/)tty$'
    hint: |
      echo "Line two: date with the bad option --${BOGUS}, || , then the terminal-name command."
    solve: |
      date --$BOGUS || tty
---

Last unit - everything at once, two one-liners:

1. Pause for **2 seconds**, and only if the pause succeeds, print the
   current time **in UTC** (either spelling of the option is fine).
2. Attempt `date` with the bad option `--${BOGUS}`, and when it fails,
   fall back to printing this terminal's name.

::task{name="punctual"}
#active
Waiting for line one: a successful pause chained to the UTC time...
#completed
Patience rewarded, on the condition of success.
::

::task{name="resilient"}
#active
Waiting for line two: a doomed `date` rescued by a fallback...
#completed
And failure handled gracefully. That's the whole first path: commands,
arguments, options, quoting, output, errors, exit statuses, and chains.
The next path teaches you how to get *unstuck*: finding help, reading
manuals, and taming commands that won't finish. See you there.
::
