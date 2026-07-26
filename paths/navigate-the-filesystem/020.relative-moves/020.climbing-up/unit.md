---
title: Climbing out of the mine
vars:
  MINE: { pick: [copper, cobalt, quartz] }
init:
  - name: create_mine
    run: |
      rm -rf /tmp/gym-mine
      mkdir -p "/tmp/gym-mine/$MINE/gallery/shaft"
      chmod -R a+rx /tmp/gym-mine
tasks:
  at_bottom:
    check: |
      wait_cwd "/tmp/gym-mine/$MINE/gallery/shaft"
    hint: |
      echo "Descend first: the shaft is at /tmp/gym-mine/${MINE}/gallery/shaft."
    solve: |
      cd /tmp/gym-mine/$MINE/gallery/shaft
  up_one:
    needs: [at_bottom]
    check: |
      wait_cwd "/tmp/gym-mine/$MINE/gallery"
    hint: |
      echo "The parent of any directory is spelled '..' - two dots. cd there."
    solve: |
      cd ..
  up_two:
    needs: [up_one]
    check: |
      wait_cwd /tmp/gym-mine
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD and the surface is /tmp/gym-mine - two levels up. '..' segments chain: parent-of-parent is ../.."
    solve: |
      cd ../..
---

Time to practice the way *up*. First take the elevator down: the deepest
point of the ${MINE} mine is `/tmp/gym-mine/${MINE}/gallery/shaft`.

::task{name="at_bottom"}
#active
Waiting for your shell at the bottom of the shaft...
#completed
Deep underground. Now, about getting out...
::

`..` always names the parent of the current directory. Climb one level,
back up to the gallery:

::task{name="up_one"}
#active
Waiting for your shell one level up, in `gallery`...
#completed
One level up, and you never typed the word "gallery".
::

`..` segments chain like any other path. Climb the remaining **two**
levels in a single move:

::task{name="up_two"}
#active
Waiting for your shell at `/tmp/gym-mine`...
#completed
Back at the surface. `..`, `../..`, `../../..` - each extra pair of dots
climbs one more level.
::
