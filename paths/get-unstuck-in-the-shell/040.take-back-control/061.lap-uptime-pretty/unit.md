---
title: Victory lap
requires: [readline]
variant: lap=uptime-pretty
vars:
  N: { pick: ["25", "30", "40"] }
  GHOST: { pick: [flurble, znark, sort, tac] }
  NAP: { pick: ["600", "900"] }
tasks:
  discovered:
    check: |
      wait_exec '(^|/)(uptime --help|man uptime)$'
    hint: |
      echo "Ask uptime itself with uptime --help (or man uptime) and look for 'pretty'."
    solve: |
      uptime --help
  applied:
    needs: [discovered]
    timeout: 60
    check: |
      wait_exec "(^|/)uptime (-p|--pretty)\$"
    hint: |
      echo "It is an option with no value: uptime, then the pretty-format option."
    solve: |
      uptime -p
  looked_up:
    timeout: 60
    check: |
      LINE=$(wait_line --latest "^((type|command|which)( +-[a-zA-Z]+)* +)?${GHOST}\b") || exit 1
      case "$LINE" in
        type\ *) exit 0 ;;
        command\ *-[vV]\ *) exit 0 ;;
        command\ *) hint_exit "That line runs ${GHOST} instead of asking about it. The -v option turns command into a lookup: command -v ${GHOST}." ;;
        which\ *) hint_exit "The which tool only searches for program files and cannot see builtins. Ask with type or command -v instead." ;;
        *) hint_exit "That line ran ${GHOST} instead of asking about it. If the prompt has not come back, press Ctrl-C. Then look the name up with type or command -v." ;;
      esac
    hint: |
      echo "Ask about the name instead of running it: type ${GHOST}, or command -v ${GHOST}."
    solve: |
      type $GHOST
  unstuck:
    timeout: 90
    check: |
      wait_exec "(^|/)sleep ${NAP}s?\$"
      wait_proc --timeout 15 "^(/usr/bin/)?slee[p] ${NAP}s?\$" || true
      wait_proc_gone "^(/usr/bin/)?slee[p] ${NAP}s?\$"
    hint: |
      echo "Start sleep ${NAP}, then free yourself with Ctrl-C."
    solve: |
      #!type sleep $NAP
      #!keys enter
      #!wait 1
      #!keys C-c
      #!wait 1
---

This is the last unit, and it combines everything from this gym in
three quick jobs:

1. The `uptime` command can report how long the machine has been up in a short,
   human-friendly form, like `up 2 hours, 5 minutes`. Find that option
   yourself, in `--help` or in the manual page, and print the uptime
   that way.
2. Does a command named `${GHOST}` exist on this machine? Find out
   without running it.
3. Start `sleep ${NAP}`, then get your prompt back without waiting.

::task{name="discovered"}
#active
Waiting for you to consult the documentation of `uptime`...
#completed
You looked the option up instead of guessing.
::

::task{name="applied"}
#active
Waiting for the uptime in the human-friendly form...
#completed
You read the documentation and applied the option.
::

::task{name="looked_up"}
#active
Waiting for a lookup of `${GHOST}` that does not run it...
#completed
You asked about the name instead of running it.
::

::task{name="unstuck"}
#active
Waiting for a runaway `sleep ${NAP}` to be dealt with...
#completed
You took control back, and that is the whole gym.
::
