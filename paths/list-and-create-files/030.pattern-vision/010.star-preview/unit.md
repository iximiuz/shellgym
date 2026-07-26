---
title: The star
vars:
  APP: { pick: [webapp, worker, gateway] }
init:
  - name: create_logs
    run: |
      rm -rf /tmp/gym-logs
      mkdir -p /tmp/gym-logs
      touch "/tmp/gym-logs/$APP-mon.log" "/tmp/gym-logs/$APP-tue.log" "/tmp/gym-logs/$APP-wed.log"
      touch /tmp/gym-logs/notes.txt /tmp/gym-logs/readme.md
      chmod -R a+rx /tmp/gym-logs
tasks:
  in_place:
    check: |
      wait_cwd /tmp/gym-logs
    hint: |
      echo "Globs expand against the current directory - cd to /tmp/gym-logs first."
    solve: |
      cd /tmp/gym-logs
  expand:
    needs: [in_place]
    check: |
      wait_exec "(^|/)(ls|echo).* (\./)?$APP-mon\.log (\./)?$APP-tue\.log (\./)?$APP-wed\.log$"
    hint: |
      echo "List (or echo) exactly the log files using one pattern: *.log - the star stands for any sequence of characters, so it matches every name ending in .log and nothing else."
    solve: |
      ls *.log
---

`/tmp/gym-logs` mixes three `.log` files with unrelated notes. Go
there:

::task{name="in_place"}
#active
Waiting for your shell in `/tmp/gym-logs`...
#completed
Five files, two kinds.
::

Now list **only the log files** - without typing their names. The
pattern `*.log` means "any name ending in `.log`": the shell expands it
into the matching names and hands those to the command.

::task{name="expand"}
#active
Waiting for a listing of exactly the `.log` files...
#completed
The command never saw a `*` - by the time it ran, the shell had already
turned the pattern into the three names. That is the key mental model
for every glob you will ever write.
::

::tip
---
title: Preview before you leap
---
`echo *.log` prints what a pattern expands to without doing anything
else - the safe way to check a glob before handing it to a destructive
command.
::
