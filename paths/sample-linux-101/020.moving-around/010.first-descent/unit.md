---
title: Change into a directory
vars:
  DIRNAME: { pick: [expedition, workshop, archive, greenhouse] }
init:
  - name: create_dir
    run: |
      mkdir -p "$GYM_USER_HOME/$DIRNAME"
      chown "$GYM_USER" "$GYM_USER_HOME/$DIRNAME"
tasks:
  arrived:
    check: |
      wait_cwd "$GYM_USER_HOME/$DIRNAME"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "Your shell is currently in $CWD. Change into the ${DIRNAME} directory with cd."
    solve: |
      cd ~/$DIRNAME
---

A directory named `${DIRNAME}` was just created in your home directory. Use
`cd` to make it your working directory.

::task{name="arrived"}
#active
Waiting for your shell to arrive in `~/${DIRNAME}`...
#completed
You moved. The prompt usually reflects the new location too.
::
