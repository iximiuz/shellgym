---
title: Create a directory
vars:
  PROJECT: { pick: [rocket, lighthouse, orchard, observatory] }
tasks:
  dir_created:
    check: |
      wait_file "$GYM_USER_HOME/$PROJECT"
    hint: |
      echo "Directories are made with mkdir followed by the name. Run it from your home directory, or give the full path."
    solve: |
      mkdir -p ~/$PROJECT
---

Time to build something. Create a directory named `${PROJECT}` in your home
directory. The command for making directories is `mkdir`.

::task
#active
Waiting for `~/${PROJECT}` to exist...
#completed
Created. `mkdir -p a/b/c` can even create a whole chain of nested
directories in one go.
::
