---
title: Home by variable
vars:
  BENCH: { pick: [drafts, blueprints, sketches] }
init:
  - name: create_workshop
    run: |
      mkdir -p "$GYM_USER_HOME/workshop/$BENCH"
      chown -R "$GYM_USER" "$GYM_USER_HOME/workshop"
tasks:
  far_away:
    check: |
      wait_cwd /srv
    hint: |
      echo "Start far from home: cd to /srv."
    solve: |
      cd /srv
  home_by_var:
    needs: [far_away]
    check: |
      wait_cwd "$GYM_USER_HOME/workshop"
    hint: |
      echo "The HOME variable holds your home directory's absolute path. Use it in a path: \$HOME/workshop - the shell substitutes the value before cd runs."
    solve: |
      cd $HOME/workshop
---

The `~` shortcut has a longhand twin: the shell keeps your home
directory's path in the `HOME` variable, and `$HOME` anywhere on a
command line is replaced by its value. Scripts tend to prefer `$HOME`;
fingers tend to prefer `~`. Both name the same place.

First get some distance - head to `/srv`:

::task{name="far_away"}
#active
Waiting for your shell in `/srv`...
#completed
Suitably far.
::

Now come back to your workshop, this time spelling home as `$HOME`:

::task{name="home_by_var"}
#active
Waiting for your shell in `~/workshop`...
#completed
Same destination, three spellings: the absolute path, `~/workshop`, and
`$HOME/workshop`. You now own all three.
::
