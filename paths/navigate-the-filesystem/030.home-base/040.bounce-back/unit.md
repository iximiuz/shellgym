---
title: The back button
vars:
  STAGE: { pick: [rehearsal, recording, mixdown] }
init:
  - name: create_studio
    run: |
      rm -rf /tmp/gym-studio
      mkdir -p "/tmp/gym-studio/$STAGE"
      chmod -R a+rx /tmp/gym-studio
tasks:
  on_stage:
    check: |
      wait_cwd "/tmp/gym-studio/$STAGE"
    hint: |
      echo "Start in the studio: /tmp/gym-studio/${STAGE}."
    solve: |
      cd /tmp/gym-studio/$STAGE
  quick_home:
    needs: [on_stage]
    check: |
      wait_cwd "$GYM_USER_HOME"
    hint: |
      echo "Pop home for a moment - remember, a bare cd is enough."
    solve: |
      cd
  and_back:
    needs: [quick_home]
    check: |
      wait_cwd "/tmp/gym-studio/$STAGE"
    hint: |
      echo "Return to the studio without retyping its path: 'cd -' jumps to the directory you were in before the last move."
    solve: |
      cd -
---

Deep in a session at the `${STAGE}` studio, you will constantly need to
pop home and come right back. The shell has a back button for exactly
this.

First, get to work:

::task{name="on_stage"}
#active
Waiting for your shell in `/tmp/gym-studio/${STAGE}`...
#completed
Session started.
::

Pop home for a moment:

::task{name="quick_home"}
#active
Waiting for your shell in your home directory...
#completed
Quick errand done.
::

Now jump straight back to the studio - **without typing its path**.
`cd -` returns to the previous working directory, whatever it was:

::task{name="and_back"}
#active
Waiting for your shell back in the studio...
#completed
It even printed where it took you. `cd -` remembers exactly one step
back - but that one step is the one you need surprisingly often.
::
