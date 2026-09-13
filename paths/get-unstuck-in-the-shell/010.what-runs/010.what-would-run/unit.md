---
title: What would run?
vars:
  CMD: { pick: [date, whoami, tty] }
tasks:
  proved:
    timeout: 45
    check: |
      wait_exec "^(/usr)?/bin/${CMD}\$"
    hint: |
      echo "Ask first: type ${CMD}. The answer ends with a location that starts with a slash. Run that location exactly as printed."
    solve: |
      type $CMD
      /usr/bin/$CMD
---

Ask the shell what the name `${CMD}` means to it:

```
type ${CMD}
```

The answer names a program file on disk. When you type `${CMD}`, the
shell looks that file up and runs it.

Now run `${CMD}` by its full location, exactly as `type` printed it.
The full location names the file directly and skips the lookup.

::task
#active
Waiting for you to run `${CMD}` by its full on-disk location...
#completed
The file ran directly, with no lookup involved.
::

::hint{title="type says the command is hashed"}
If you ran `${CMD}` recently, `type` may say that it is *hashed*,
because the shell remembered the lookup. The location in parentheses
is still the file you want.
::
