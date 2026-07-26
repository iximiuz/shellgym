---
title: No place like ~
tasks:
  drift_out:
    check: |
      wait_cwd /var/tmp
    hint: |
      echo "First wander off somewhere: /var/tmp will do."
    solve: |
      cd /var/tmp
  come_home:
    needs: [drift_out]
    check: |
      wait_cwd "$GYM_USER_HOME"
    hint: |
      echo "cd with no argument at all - just the two letters - always brings you home."
    solve: |
      cd
---

Wherever your work takes you, home is one command away. First, wander
off to `/var/tmp`:

::task{name="drift_out"}
#active
Waiting for your shell in `/var/tmp`...
#completed
Far from home.
::

Now return home - and here is the trick: `cd` **without any argument**
goes to your home directory. No path to remember, no typing.

::task{name="come_home"}
#active
Waiting for your shell back in your home directory...
#completed
Home. A bare `cd` is the fastest escape hatch in the shell - from
anywhere, ever.
::
