---
title: Find out where you are
tasks:
  saved_location:
    check: |
      wait_file_contains /tmp/here.txt "^${GYM_USER_HOME}$"
    hint: |
      if [ -f /tmp/here.txt ]; then
        echo "The file /tmp/here.txt exists but does not contain your home directory path. Check its content with: cat /tmp/here.txt"
      else
        echo "No /tmp/here.txt yet. Remember: the '>' arrow sends a command's output into a file."
      fi
    solve: |
      pwd
      pwd > /tmp/here.txt
---

Every shell always sits in some directory, called the working directory.
The `pwd` command (print working directory) shows it. Try it in the
terminal:

```
pwd
```

You should see your home directory. Now save that location into a file so
the check can verify it. The `>` operator redirects a command's output into
a file:

```
pwd > /tmp/here.txt
```

::task
#active
Waiting for `/tmp/here.txt` to contain your working directory...
#completed
That is your home directory, and your first completed task.
::
