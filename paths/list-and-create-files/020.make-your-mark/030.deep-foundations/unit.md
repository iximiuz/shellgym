---
title: Deep foundations
vars:
  SITE: { pick: [seattle, lisbon, osaka] }
  TIER: { pick: [prod, staging] }
init:
  - name: create_warehouse
    run: |
      rm -rf /tmp/gym-warehouse
      mkdir -p /tmp/gym-warehouse
      chown "$GYM_USER" /tmp/gym-warehouse
tasks:
  racks_up:
    check: |
      wait_dir "/tmp/gym-warehouse/$SITE/$TIER/racks"
    hint: |
      echo "Plain mkdir refuses to create a directory whose parent is missing. One option tells it to create every missing parent along the way - check mkdir --help for 'parents'."
    solve: |
      mkdir -p /tmp/gym-warehouse/$SITE/$TIER/racks
---

New site coming online: the warehouse needs the whole chain
`/tmp/gym-warehouse/${SITE}/${TIER}/racks` - three nested directories,
none of which exist yet.

Try the obvious `mkdir` and it will complain: the parents are missing.
There is an option that makes `mkdir` create **every missing level** in
one go.

::task{name="racks_up"}
#active
Waiting for `/tmp/gym-warehouse/${SITE}/${TIER}/racks` to exist...
#completed
Whole chain in one command. `mkdir -p` is also quiet when the
directory already exists - which makes it a favorite in scripts that
must be safe to run twice.
::
