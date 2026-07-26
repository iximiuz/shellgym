---
title: The dash trap
init:
  - name: create_post
    run: |
      rm -rf /tmp/gym-post
      mkdir -p /tmp/gym-post
      chown "$GYM_USER" /tmp/gym-post
tasks:
  urgent_filed:
    check: |
      wait_file "/tmp/gym-post/-urgent.txt"
    hint: |
      echo "touch reads -urgent.txt as options, not a name. Two ways out: give it a path form it can't mistake (./-urgent.txt), or put -- before the name to declare 'no more options'."
    solve: |
      cd /tmp/gym-post
      touch ./-urgent.txt
---

A message must be filed in `/tmp/gym-post` as `-urgent.txt` - dash
first, as some tool upstream insists.

Try the naive `touch -urgent.txt` there and watch it fail: anything
starting with `-` looks like an *option* to the command. The file name
never gets a chance.

Two standard escapes:

- Path form: `./-urgent.txt` - a leading `./` makes it unmistakably a
  path.
- The option terminator: `--` tells the command "everything after this
  is a name, not an option".

::task{name="urgent_filed"}
#active
Waiting for `-urgent.txt` to exist in `/tmp/gym-post`...
#completed
Filed. Both escapes work on nearly every command, and you will want
them again the day a stray `-rf` shows up as a *file name*.
::
