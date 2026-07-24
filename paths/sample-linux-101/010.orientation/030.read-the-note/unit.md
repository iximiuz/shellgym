---
title: Read a file
vars:
  SECRET: { shell: "head -c3 /dev/urandom | od -An -tx1 | tr -d ' \\n'" }
init:
  - name: drop_note
    run: |
      printf 'A note from Shell Gym.\nYour code word is: %s\n' "$SECRET" > "$GYM_USER_HOME/note.txt"
      chown "$GYM_USER" "$GYM_USER_HOME/note.txt"
tasks:
  read_note:
    check: |
      wait_exec '(^|/)(cat|less|more|head|tail|nano|vim?) .*note\.txt'
    hint: |
      echo "The file is note.txt in your home directory. The most common way to print a file is the 'cat' command."
    solve: |
      cat ~/note.txt
---

A file named `note.txt` just appeared in your home directory. It contains a
code word. Print the file's content in the terminal to read it.

::task{name="read_note"}
#active
Waiting for you to open `note.txt`...
#completed
Code word retrieved: `${SECRET}`. `cat` is the go-to command for quickly
printing a file.
::
