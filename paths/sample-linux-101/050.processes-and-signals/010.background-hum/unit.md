---
title: Run a command in the background
tasks:
  humming:
    check: |
      wait_proc '(^|/)sleep 543'
    hint: |
      echo "Append a single & to the end of the command line to run it in the background. The shell prints the job number and PID, then gives the prompt back."
    solve: |
      sleep 543 &
---

A command normally holds your terminal until it finishes. Appending `&`
sends it to the background instead, and the prompt returns immediately.

Start this long sleeper in the background:

```
sleep 543 &
```

::task{name="humming"}
#active
Waiting for a background `sleep 543` process...
#completed
It runs, you type. The `jobs` command lists your background jobs, and `fg`
brings one back to the foreground.
::
