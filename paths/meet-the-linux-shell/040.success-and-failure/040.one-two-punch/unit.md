---
title: Three commands, one line
tasks:
  rollcall:
    timeout: 45
    check: |
      wait_exec '(^|/)whoami$'
      wait_exec '(^|/)hostname$'
      wait_exec '(^|/)tty$'
    hint: |
      echo "One line, three commands, two semicolons between them. Who, which machine, which terminal."
    solve: |
      whoami; hostname; tty
---

You don't have to run commands one prompt at a time. A semicolon `;`
separates several commands on a single line; the shell runs them left
to right, each one after the previous finishes - no matter whether it
succeeded or failed.

Do the full identity roll call - user name, machine name, terminal
name - in **one line**.

::task
#active
Waiting for the one-line roll call: who you are, what machine you are on, on
which terminal...
#completed
Three answers from one Enter press. `;` is the "and then" of the
shell - blind sequencing, no questions asked. Next up: chaining that
*does* ask questions.
::
