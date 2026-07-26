---
title: Jump to a project directory
vars:
  APP: { pick: [signalhub, chronos, meshmon] }
init:
  - name: create_project
    run: |
      rm -rf /opt/gymtrack
      mkdir -p "/opt/gymtrack/$APP"
      chmod -R a+rx /opt/gymtrack
tasks:
  arrive:
    check: |
      wait_cwd "/opt/gymtrack/$APP"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "Your shell is in $CWD. An absolute path starts with / - give cd the full address /opt/gymtrack/${APP} and it works from anywhere."
    solve: |
      cd /opt/gymtrack/$APP
---

A project called `${APP}` has just been checked out under
`/opt/gymtrack/${APP}`. Get your shell there.

If you are ever unsure where you currently are, `pwd` tells you - it
prints the shell's working directory.

::task{name="arrive"}
#active
Waiting for your shell in `/opt/gymtrack/${APP}`...
#completed
You're in. One `cd` with a full path takes you anywhere on the system,
no matter where you started.
::

::tip
---
title: Reading the ground
---
Run `pwd` whenever you feel lost. Many prompts also show the current
directory, but `pwd` is the source of truth.
::
