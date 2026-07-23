---
title: Create an empty file
needs: [make-a-home]
vars:
  PROJECT: { from: make-a-home.PROJECT }
tasks:
  file_created:
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_file "$HOME_DIR/$PROJECT/*.log"
    hint: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      if [ -d "$HOME_DIR/$PROJECT" ]; then
        echo "The ${PROJECT} directory is there, but no .log file inside it yet. touch creates empty files."
      else
        echo "The ~/${PROJECT} directory disappeared. Recreate it first with mkdir."
      fi
    solve: |
      touch ~/$PROJECT/build.log
---

Inside your new `${PROJECT}` directory, create an empty file. Name it
anything you like, as long as the name ends with `.log`. The `touch`
command creates an empty file (or updates the timestamp of an existing
one).

::task{name="file_created"}
#active
Waiting for a `.log` file inside `~/${PROJECT}`...
#completed
File planted. Note that the check accepted any name ending in `.log` - the
checks can match patterns, not just exact names.
::
