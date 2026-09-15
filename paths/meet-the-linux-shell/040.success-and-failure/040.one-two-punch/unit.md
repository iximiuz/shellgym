---
title: Three commands on one line
requires: [readline]
tasks:
  rollcall:
    timeout: 45
    check: |
      LINE=$(wait_line --latest '^[a-z]+ *(;|&&|\|\|) *[a-z]+ *(;|&&|\|\|) *[a-z]+$') || exit 1
      case "$LINE" in
        *'&&'*|*'||'*) hint_exit "Three commands ran from one line, but they were joined with && or ||. Join them with the semicolon instead. It runs the next command no matter what happened before." ;;
      esac
      for cmd in whoami hostname tty; do
        case "$LINE" in
          *"$cmd"*) ;;
          *) hint_exit "The line has three commands, but $cmd is missing from it." ;;
        esac
      done
      wait_exec '(^|/)whoami$'
      wait_exec '(^|/)hostname$'
      wait_exec '(^|/)tty$'
    hint: |
      echo "You need one line with three commands (whoami, hostname, tty) and two semicolons between them."
    solve: |
      whoami; hostname; tty
---

You do not have to run commands one prompt at a time. A semicolon `;`
separates several commands on a single line. The shell runs them left
to right, each one after the previous finishes, no matter whether it
succeeded or failed.

Do the full identity roll call (user name, machine name, terminal
name) in **one line**.

::task
#active
Waiting for the one-line roll call: who you are, what machine you are on, on
which terminal...
#completed
One Enter press produced three answers.
::
