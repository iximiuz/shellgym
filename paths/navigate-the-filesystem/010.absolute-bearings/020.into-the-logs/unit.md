---
title: Head for the logs
vars:
  SVC: { pick: [relay, beacon, harvest] }
init:
  - name: create_logs
    run: |
      rm -rf /var/log/gym
      mkdir -p "/var/log/gym/$SVC"
      touch "/var/log/gym/$SVC/access.log" "/var/log/gym/$SVC/error.log"
      chmod -R a+rx /var/log/gym
tasks:
  arrive:
    check: |
      wait_cwd "/var/log/gym/$SVC"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "Currently in $CWD. Logs on Linux live under /var/log - the ${SVC} service keeps its logs in /var/log/gym/${SVC}."
    solve: |
      cd /var/log/gym/$SVC
---

The `${SVC}` service is misbehaving and its logs are in
`/var/log/gym/${SVC}`. Investigations start at the scene - move your
shell there.

`/var/log` is where most Linux services keep their logs; you will `cd`
into it for the rest of your career.

::task{name="arrive"}
#active
Waiting for your shell in `/var/log/gym/${SVC}`...
#completed
On the scene. `access.log` and `error.log` are sitting right here -
reading them is a later path's business.
::

::tip
---
title: Let the shell type for you
---
Type `cd /var/lo` and press `Tab` - the shell completes the directory
name. Tab completion works on every path segment and never misspells.
::
