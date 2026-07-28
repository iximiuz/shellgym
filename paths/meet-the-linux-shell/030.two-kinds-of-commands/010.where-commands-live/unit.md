---
title: Where does date live?
tasks:
  located:
    check: |
      wait_exec '(^|/)which (-a )?date$'
    hint: |
      echo "Run which with the command name you're curious about: date."
    solve: |
      which date
---

If most commands are programs stored on disk, they must live
*somewhere*. The command `which` answers exactly that question: give it
a command name as an argument and it prints the location of the program
file that would run.

Ask it where `date` lives.

::task
#active
Waiting for you to look up the home of `date`...
#completed
`/usr/bin/date` - a real file on disk. When you type `date`, the shell
finds that file and runs it for you. Most commands work this way.
::
