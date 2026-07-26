---
title: A quiet note
vars:
  SECRET: { pick: [todo, wishlist, plans] }
init:
  - name: create_hush
    run: |
      rm -rf /tmp/gym-hush
      mkdir -p /tmp/gym-hush
      touch /tmp/gym-hush/readme.txt
      chown "$GYM_USER" /tmp/gym-hush
tasks:
  hidden_note:
    check: |
      wait_file "/tmp/gym-hush/.$SECRET"
    hint: |
      echo "A file is hidden when its name starts with a dot - create .${SECRET} (dot included) in /tmp/gym-hush with touch."
    solve: |
      touch /tmp/gym-hush/.$SECRET
  double_check:
    needs: [hidden_note]
    check: |
      wait_exec '(^|/)ls.* (-[a-zA-Z]*[aA][a-zA-Z]*|--all|--almost-all)( |$)'
    hint: |
      echo "Plain ls won't show your new file - which listing option reveals hidden entries? You used it in the attic."
    solve: |
      ls -a /tmp/gym-hush
---

Leave yourself a private note in `/tmp/gym-hush`: an empty file named
`.${SECRET}` - with the leading dot, so casual listings pass it by.

Hiding a file is nothing but naming: start the name with `.` and it is
hidden. Same `touch` as always.

::task{name="hidden_note"}
#active
Waiting for `.${SECRET}` to appear in `/tmp/gym-hush`...
#completed
Note left. `readme.txt` shows in every listing; your `.${SECRET}` only
shows to those who ask properly.
::

Now confirm it is really there - list the directory in the way that
shows hidden entries:

::task{name="double_check"}
#active
Waiting for a listing that would reveal `.${SECRET}`...
#completed
Created blind, verified with `-a`. Make-then-check is a habit worth
keeping: commands that succeed silently still deserve a glance.
::
