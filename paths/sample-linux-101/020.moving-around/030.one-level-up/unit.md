---
title: Go one level up
needs: [deep-dive]
vars:
  LEVEL1: { from: deep-dive.LEVEL1 }
  LEVEL2: { from: deep-dive.LEVEL2 }
tasks:
  went_up:
    check: |
      # wait_cwd prints the PID of the shell that matched; publish it as a
      # task var so the homeward unit can watch this same shell come home.
      TRAVELER=$(wait_cwd "/tmp/depths/$LEVEL1/$LEVEL2") || exit 1
      set_var TRAVELER "$TRAVELER"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. Two dots (..) always mean the parent of the current directory."
    solve: |
      cd ..
---

Every directory has a parent, referred to as `..` (two dots). You are still
somewhere in the `/tmp/depths` tree from the previous unit. Move so that
your working directory becomes `/tmp/depths/${LEVEL1}/${LEVEL2}`.

::task{name="went_up"}
#active
Waiting for your shell in `/tmp/depths/${LEVEL1}/${LEVEL2}`...
#completed
`..` works from anywhere, and it chains: `../..` climbs two levels at once.
::
