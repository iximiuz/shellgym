---
title: Redirect output into a file
vars:
  MOTTO: { pick: [festina lente, carpe diem, sapere aude] }
tasks:
  wrote_file:
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_file_contains "$HOME_DIR/motto.txt" "^$MOTTO$"
    hint: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      if [ -f "$HOME_DIR/motto.txt" ]; then
        echo "motto.txt exists but its content is not exactly '${MOTTO}'. Overwrite it with a fresh redirect (a single > replaces the whole file)."
      else
        echo "echo prints its arguments; the > arrow captures that output into a file instead of the screen."
      fi
    solve: |
      echo "$MOTTO" > ~/motto.txt
---

The `echo` command prints text. Combined with `>` it becomes the quickest
way to create a file with known content.

Create a file `motto.txt` in your home directory containing exactly one
line: `${MOTTO}`.

::task{name="wrote_file"}
#active
Waiting for `~/motto.txt` with the motto...
#completed
Captured. `>` always replaces the file's previous content, which segues
into the next unit.
::
