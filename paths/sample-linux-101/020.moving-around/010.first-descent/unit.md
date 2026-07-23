---
title: Change into a directory
vars:
  DIRNAME: { pick: [expedition, workshop, archive, greenhouse] }
init:
  - name: create_dir
    run: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      mkdir -p "$HOME_DIR/$DIRNAME"
      chown "$GYM_USER" "$HOME_DIR/$DIRNAME"
tasks:
  arrived:
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_cwd "$HOME_DIR/$DIRNAME"
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
