---
title: List without leaving
vars:
  DEPOT: { pick: [parcels, freight, cargo] }
init:
  - name: create_depot
    run: |
      rm -rf /var/gym-depot
      mkdir -p "/var/gym-depot/$DEPOT"
      touch "/var/gym-depot/$DEPOT/manifest.txt" "/var/gym-depot/$DEPOT/waybill.txt"
      chmod -R a+rx /var/gym-depot
tasks:
  remote_look:
    check: |
      wait_exec '(^|/)ls .*gym-depot'
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "Stay in $CWD - ls takes a path argument: point it at /var/gym-depot/${DEPOT} and it lists that directory instead of the current one."
    solve: |
      ls /var/gym-depot/$DEPOT
---

What is in `/var/gym-depot/${DEPOT}`? You could `cd` there and look -
but you do not have to. `ls` accepts a path and lists **that**
directory, while your shell stays put.

Check the depot's contents from wherever you are:

::task{name="remote_look"}
#active
Waiting for you to list `/var/gym-depot/${DEPOT}` by its path...
#completed
Contents inspected, shell unmoved. Listing by path keeps your bearings
when you only need a peek - no round trip required.
::
