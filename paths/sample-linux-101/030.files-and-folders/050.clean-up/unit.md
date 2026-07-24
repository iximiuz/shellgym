---
title: Delete a file
init:
  - name: scatter_junk
    run: |
      echo scrap > "$GYM_USER_HOME/junk.tmp"
      chown "$GYM_USER" "$GYM_USER_HOME/junk.tmp"
tasks:
  junk_gone:
    check: |
      # Baseline: the file must exist before we watch for its removal,
      # otherwise this check would pass without the student doing anything.
      wait_file --timeout 15 "$GYM_USER_HOME/junk.tmp" ||  hint_exit "junk.tmp is not in your home directory (was it removed before the rep started?). Reset the rep to re-create it."
      wait_file_gone "$GYM_USER_HOME/junk.tmp"
    hint: |
      echo "rm followed by the path removes a file. There is no trash bin on the command line - gone is gone."
    solve: |
      rm ~/junk.tmp
---

A useless file `junk.tmp` has appeared in your home directory. Remove it
with `rm`.

Be deliberate with this command: the shell has no trash bin, and `rm` does
not ask for confirmation.

::task{name="junk_gone"}
#active
Waiting for `junk.tmp` to be removed...
#completed
Clean. For directories, `rmdir` removes empty ones and `rm -r` removes
recursively - the latter deserves a second look before you press Enter.
::
