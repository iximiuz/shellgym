---
title: Look around
init:
  - name: seed_home
    run: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      mkdir -p "$HOME_DIR/notes" "$HOME_DIR/projects"
      touch "$HOME_DIR/readme.txt"
      chown "$GYM_USER" "$HOME_DIR/notes" "$HOME_DIR/projects" "$HOME_DIR/readme.txt"
tasks:
  ran_ls:
    check: |
      wait_exec '(^|/)ls($| )'
    hint: |
      echo "Type the two-letter listing command and press Enter. Flags are optional."
    solve: |
      ls
---

The `ls` command lists what is inside a directory. Run it (with or without
extra flags) to see what your home directory contains:

::task{name="ran_ls"}
#active
Waiting for you to list a directory...
#completed
You should have spotted `notes`, `projects`, and `readme.txt`. Add `-l` next
time for a detailed view.
::

::hint{title="Want more detail?"}
`ls -l` shows permissions, owners, sizes, and dates. `ls -a` also reveals
hidden files, the ones whose names start with a dot.
::
