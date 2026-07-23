---
title: Descend a deep path
vars:
  LEVEL1: { pick: [north, south, east, west] }
  LEVEL2: { pick: [upper, lower, middle] }
  LEVEL3: { pick: [red, green, blue, gold] }
init:
  - name: create_tree
    run: |
      mkdir -p "/tmp/depths/$LEVEL1/$LEVEL2/$LEVEL3"
      chmod -R a+rx /tmp/depths
tasks:
  at_bottom:
    check: |
      wait_cwd "/tmp/depths/$LEVEL1/$LEVEL2/$LEVEL3"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. The target is /tmp/depths/${LEVEL1}/${LEVEL2}/${LEVEL3}. You can cd there in one jump by giving the whole path."
    solve: |
      cd /tmp/depths/$LEVEL1/$LEVEL2/$LEVEL3
---

Directories nest. A path like `/tmp/depths/${LEVEL1}/${LEVEL2}/${LEVEL3}`
describes the whole chain from the root. You can walk it one `cd` at a
time, or jump straight there with a single `cd` and the full path.

Get your shell to the bottom:

::task{name="at_bottom"}
#active
Waiting for your shell in `/tmp/depths/${LEVEL1}/${LEVEL2}/${LEVEL3}`...
#completed
Single jump or step by step, you made it. Tab completion makes long paths
painless: type a few letters and press Tab.
::

::hint
---
title: Typing long paths
---
Press Tab after a few letters of each directory name and the shell completes
it for you. Double-Tab lists the options.
::
