---
title: Over to the next wing
vars:
  WING1: { pick: [brewery, foundry, roastery] }
  WING2: { pick: [storehouse, workshop, cellar] }
init:
  - name: create_plant
    run: |
      rm -rf /tmp/gym-plant
      mkdir -p "/tmp/gym-plant/$WING1" "/tmp/gym-plant/$WING2"
      chmod -R a+rx /tmp/gym-plant
tasks:
  in_first:
    check: |
      wait_cwd "/tmp/gym-plant/$WING1"
    hint: |
      echo "Start in the ${WING1}: /tmp/gym-plant/${WING1}."
    solve: |
      cd /tmp/gym-plant/$WING1
  across:
    needs: [in_first]
    check: |
      wait_cwd "/tmp/gym-plant/$WING2"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. The ${WING2} sits next to the ${WING1} - up to the parent and down the other side, in one path that starts with '..'."
    solve: |
      cd ../$WING2
---

The plant at `/tmp/gym-plant` has two wings side by side: the
`${WING1}` and the `${WING2}`. Start your shift in the `${WING1}`:

::task{name="in_first"}
#active
Waiting for your shell in the `${WING1}`...
#completed
Clocked in.
::

Now walk across to the `${WING2}`. You could spell out the absolute
path - but the neighboring wing is just *up one, down the other side*,
and `..` composes with a name into a single move.

::task{name="across"}
#active
Waiting for your shell in the `${WING2}`...
#completed
`../${WING2}` - sideways moves between siblings are one of the most
common walks in any project tree.
::
