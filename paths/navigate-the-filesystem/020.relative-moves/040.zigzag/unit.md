---
title: Zigzag through the maze
vars:
  NOOK: { pick: [nook, alcove, cranny] }
  DEN: { pick: [den, lair, burrow] }
init:
  - name: create_maze
    run: |
      rm -rf /tmp/gym-maze
      mkdir -p "/tmp/gym-maze/entry/left/$NOOK" "/tmp/gym-maze/entry/right/$DEN"
      chmod -R a+rx /tmp/gym-maze
tasks:
  at_nook:
    check: |
      wait_cwd "/tmp/gym-maze/entry/left/$NOOK"
    hint: |
      echo "Start deep in the left branch: /tmp/gym-maze/entry/left/${NOOK}."
    solve: |
      cd /tmp/gym-maze/entry/left/$NOOK
  over_to_den:
    needs: [at_nook]
    check: |
      wait_cwd "/tmp/gym-maze/entry/right/$DEN"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD; the ${DEN} is at /tmp/gym-maze/entry/right/${DEN}. Count the climb: two levels up to entry, then down through right. '..' pairs chain with names in one path."
    solve: |
      cd ../../right/$DEN
---

The maze under `/tmp/gym-maze` splits at `entry` into `left` and
`right` branches. You start tucked away in the left branch:

::task{name="at_nook"}
#active
Waiting for your shell in `/tmp/gym-maze/entry/left/${NOOK}`...
#completed
Cozy. But you are needed on the other side.
::

Cross over to the `${DEN}`, which is at `entry/right/${DEN}` - in a
**single** relative move. Climb as many levels as needed with `..`
segments, then descend the other branch, all in one path.

::task{name="over_to_den"}
#active
Waiting for your shell in the `${DEN}` on the right branch...
#completed
Up, up, across, down - one path. When you can zigzag like this without
sketching the tree on paper, relative paths are yours.
::

::hint
---
title: How many dots?
---
Each `..` climbs exactly one level. From `entry/left/${NOOK}`, one pair
reaches `left`, two reach `entry` - and from there the right branch is a
plain descent.
::
