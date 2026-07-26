---
title: Shelf space
vars:
  BOX: { pick: [invoices, receipts, manifests] }
init:
  - name: create_depot
    run: |
      rm -rf /tmp/gym-depot
      mkdir -p /tmp/gym-depot
      chown "$GYM_USER" /tmp/gym-depot
tasks:
  shelf_up:
    check: |
      wait_dir "/tmp/gym-depot/$BOX"
    hint: |
      if [ -f "/tmp/gym-depot/$BOX" ]; then
        echo "There is a regular FILE named ${BOX} at /tmp/gym-depot - the task needs a directory. Remove the file (rm /tmp/gym-depot/${BOX}) and use the directory-making command instead."
      else
        echo "mkdir <path> creates a directory. The depot is /tmp/gym-depot; the new shelf is ${BOX}."
      fi
    solve: |
      mkdir /tmp/gym-depot/$BOX
---

The depot at `/tmp/gym-depot` needs a new shelf for `${BOX}` - a
**directory**, ready to hold files.

`touch` makes files; directories have their own command: `mkdir`.

::task{name="shelf_up"}
#active
Waiting for the directory `/tmp/gym-depot/${BOX}` to exist...
#completed
Shelf mounted. Files and directories are different things even when
their names look alike - the check would not have accepted a plain file
here.
::
