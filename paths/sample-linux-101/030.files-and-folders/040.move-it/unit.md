---
title: Move and rename a file
vars:
  CRATE: { pick: [crate-a, crate-b, crate-c] }
init:
  - name: seed_box
    run: |
      mkdir -p "$GYM_USER_HOME/inbox" "$GYM_USER_HOME/outbox"
      echo cargo > "$GYM_USER_HOME/inbox/$CRATE.tmp"
      chown -R "$GYM_USER" "$GYM_USER_HOME/inbox" "$GYM_USER_HOME/outbox"
tasks:
  moved:
    check: |
      wait_file "$GYM_USER_HOME/outbox/$CRATE.txt"
    hint: |
      if [ -f "$GYM_USER_HOME/inbox/$CRATE.tmp" ]; then
        echo "The file is still in inbox/. mv can move and rename in a single step: give it the old path and the new path."
      else
        echo "inbox/${CRATE}.tmp is gone but outbox/${CRATE}.txt has not appeared. Check where the file ended up with ls."
      fi
    solve: |
      mv ~/inbox/$CRATE.tmp ~/outbox/$CRATE.txt
  old_gone:
    check: |
      wait_file_gone "$GYM_USER_HOME/inbox/$CRATE.tmp"
---

Your home directory now has `inbox` and `outbox` directories, and a file
`inbox/${CRATE}.tmp`. Move the file into `outbox`, renaming it to
`${CRATE}.txt` along the way. Unlike `cp`, the `mv` command leaves no file
behind.

::task{name="moved"}
#active
Waiting for `outbox/${CRATE}.txt`...
#completed
Landed in `outbox` under the new name.
::

::task{name="old_gone"}
#active
Waiting for `inbox/${CRATE}.tmp` to disappear...
#completed
And the original is gone - that is what makes it a move, not a copy.
::
