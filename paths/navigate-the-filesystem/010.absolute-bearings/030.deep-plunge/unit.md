---
title: One jump, three levels
vars:
  L1: { pick: [north, south, east] }
  L2: { pick: [rack-a, rack-b, rack-c] }
  L3: { pick: [box-red, box-blue, box-gold] }
init:
  - name: create_depot
    run: |
      rm -rf /srv/depot
      mkdir -p "/srv/depot/$L1/$L2/$L3"
      chmod -R a+rx /srv/depot
tasks:
  at_bottom:
    check: |
      wait_cwd "/srv/depot/$L1/$L2/$L3"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. The whole chain /srv/depot/${L1}/${L2}/${L3} fits in a single cd - no need to enter each level separately."
    solve: |
      cd /srv/depot/$L1/$L2/$L3
---

Directories nest. The path `/srv/depot/${L1}/${L2}/${L3}` describes a
chain three levels deep - and one `cd` can walk the whole chain at once.

Get to the bottom:

::task{name="at_bottom"}
#active
Waiting for your shell in `/srv/depot/${L1}/${L2}/${L3}`...
#completed
Straight to the bottom in one move. Stepping level by level works too,
but the single jump is faster - and Tab completion makes it painless.
::

::tip
---
title: Tab, tab, tab
---
Press `Tab` after each `/` and the shell fills in the next directory
name. A long path becomes a few keystrokes.
::
