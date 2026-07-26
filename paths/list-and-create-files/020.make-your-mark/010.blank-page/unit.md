---
title: A blank page
vars:
  PAGE: { pick: [journal, logbook, ledger] }
init:
  - name: create_desk
    run: |
      mkdir -p "$GYM_USER_HOME/desk"
      rm -f "$GYM_USER_HOME/desk/$PAGE.txt"
      chown "$GYM_USER" "$GYM_USER_HOME/desk"
tasks:
  page_exists:
    check: |
      wait_file "$GYM_USER_HOME/desk/$PAGE.txt"
    hint: |
      echo "touch <path> creates an empty file at that path. The desk is ~/desk; the page is ${PAGE}.txt."
    solve: |
      touch ~/desk/$PAGE.txt
---

There is a fresh desk at `~/desk`. Start a `${PAGE}` on it: create an
empty file named `${PAGE}.txt` there.

The tool is `touch` - given a path that does not exist yet, it creates
an empty file.

::task{name="page_exists"}
#active
Waiting for `~/desk/${PAGE}.txt` to appear...
#completed
A file made of nothing but a name - zero bytes, confirmed by `ls -l`.
Empty files matter more than you'd think: markers, locks, placeholders.
::
