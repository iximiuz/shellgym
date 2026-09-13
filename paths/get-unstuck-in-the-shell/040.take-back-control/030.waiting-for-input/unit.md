---
title: A command waiting for input
tasks:
  hung:
    check: |
      wait_exec --argc 1 '(^|/)cat$'
    hint: |
      echo "Type cat with nothing after it and press Enter. It will look like nothing is happening, and that is the exercise."
    solve: |
      #!type cat
      #!keys enter
      #!wait 1
  recovered:
    needs: [hung]
    timeout: 60
    check: |
      wait_proc --timeout 15 '^(/usr/bin/|/bin/)?ca[t]$' || true
      wait_proc_gone '^(/usr/bin/|/bin/)?ca[t]$'
    hint: |
      echo "The terminal is fine. cat is waiting for input that never comes, and Ctrl-C gets you out."
    solve: |
      #!keys C-c
      #!wait 1
---

Sooner or later you run a command and nothing happens. There is no
output, no error, and no prompt. The command is rarely frozen. In most
cases it is waiting for input, and it will keep waiting until it gets
some.

Meet the situation on your own terms. Run `cat` with nothing after it:

```
cat
```

The `cat` command is a file-reading tool that you will meet properly later.
Given no file, it reads from your keyboard. Type a word and press Enter, and
`cat` repeats it back. It keeps doing that forever, because it has no
way to know that you are done.

You already know the way out. Take it. The next time a terminal seems
to freeze, ask yourself whether the command is waiting for you. Often
it is, and if you do not want to give it anything, Ctrl-C is one way
out.

::task{name="hung"}
#active
Waiting for a bare `cat` with no arguments...
#completed
Now `cat` is waiting for your input.
::

::task{name="recovered"}
#active
Waiting for you to break out of it...
#completed
The prompt is recovered.
::
