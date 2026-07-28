---
title: The echo twins
tasks:
  twin_ran:
    check: |
      wait_exec '^(/usr)?/bin/echo .+'
    hint: |
      echo "Type the full location /bin/echo, a space, and then any words you like."
    solve: |
      /bin/echo hello from the standalone twin
---

Meet `echo` - it simply prints back whatever arguments you give it:

```
echo have we met before
```

Try it. Then here's the twist: `echo` is a **builtin**. When you type
it, the shell doesn't launch any program - it prints the words itself,
instantly. And yet a standalone program file with the same name *also*
exists, at `/bin/echo`. Typing a command's full location runs that exact
file, bypassing the builtin.

Print a greeting using the standalone twin at `/bin/echo`.

::task
#active
Waiting for a message printed by the standalone `/bin/echo`...
#completed
Both twins look identical from the outside - but this time a real
separate program ran. The shell prefers its builtin only when you use
the bare name.
::
