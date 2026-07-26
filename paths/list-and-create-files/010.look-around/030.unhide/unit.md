---
title: What the attic hides
vars:
  TRINKET: { pick: [locket, compass, spyglass] }
init:
  - name: create_attic
    run: |
      rm -rf /tmp/gym-attic
      mkdir -p /tmp/gym-attic
      touch /tmp/gym-attic/lamp.txt /tmp/gym-attic/chair.txt
      touch "/tmp/gym-attic/.$TRINKET"
      chmod -R a+rx /tmp/gym-attic
tasks:
  up_there:
    check: |
      wait_cwd /tmp/gym-attic
    hint: |
      echo "Climb up first: cd /tmp/gym-attic."
    solve: |
      cd /tmp/gym-attic
  reveal:
    needs: [up_there]
    check: |
      wait_exec '(^|/)ls.* (-[a-zA-Z]*[aA][a-zA-Z]*|--all|--almost-all)( |$)'
    hint: |
      echo "Plain ls skips names that start with a dot. One short option makes it show all entries - a is for 'all'."
    solve: |
      ls -a
---

The attic at `/tmp/gym-attic` looks nearly empty - a lamp, a chair.
But atticks keep secrets. Go up and look:

::task{name="up_there"}
#active
Waiting for your shell in `/tmp/gym-attic`...
#completed
Dusty. A plain `ls` shows the lamp and the chair... and that's it?
::

Names that start with a dot are hidden from plain listings. `ls` has an
option that reveals **all** entries, hidden ones included. Find what
the attic is hiding:

::task{name="reveal"}
#active
Waiting for a listing that shows hidden entries...
#completed
There it is: `.${TRINKET}`. Note the `.` and `..` entries in the
listing too - the directory itself and its parent, present in every
directory on the system.
::
