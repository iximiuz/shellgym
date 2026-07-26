---
title: The grand circuit
vars:
  DOCK1: { pick: [pier-north, pier-east] }
  DOCK2: { pick: [pier-south, pier-west] }
init:
  - name: create_harbor
    run: |
      rm -rf /opt/gym-harbor
      mkdir -p "/opt/gym-harbor/$DOCK1" "/opt/gym-harbor/$DOCK2"
      chmod -R a+rx /opt/gym-harbor
tasks:
  step_1:
    check: |
      wait_cwd "/opt/gym-harbor/$DOCK1"
    hint: |
      echo "Leg 1: absolute jump to /opt/gym-harbor/${DOCK1}."
    solve: |
      cd /opt/gym-harbor/$DOCK1
  step_2:
    needs: [step_1]
    check: |
      wait_cwd "/opt/gym-harbor/$DOCK2"
    hint: |
      echo "Leg 2: ${DOCK2} is a sibling of ${DOCK1} - one relative move through '..'."
    solve: |
      cd ../$DOCK2
  step_3:
    needs: [step_2]
    check: |
      wait_cwd "$GYM_USER_HOME"
    hint: |
      echo "Leg 3: home. The shortest command you know."
    solve: |
      cd
  step_4:
    needs: [step_3]
    check: |
      wait_cwd "/opt/gym-harbor/$DOCK2"
    hint: |
      echo "Leg 4: back to ${DOCK2} without typing its path - the back button."
    solve: |
      cd -
---

Final inspection round at the harbor: four legs, and each leg has a
*shortest* move you have already mastered. Use the best tool for each.

**Leg 1** - report to `/opt/gym-harbor/${DOCK1}`:

::task{name="step_1"}
#active
Waiting for your shell at `${DOCK1}`...
#completed
Absolute path - the opening move.
::

**Leg 2** - cross to the neighboring `${DOCK2}`:

::task{name="step_2"}
#active
Waiting for your shell at `${DOCK2}`...
#completed
A sibling hop through `..` - no full address needed.
::

**Leg 3** - return to base (your home directory):

::task{name="step_3"}
#active
Waiting for your shell at home...
#completed
The bare `cd`.
::

**Leg 4** - one last check at `${DOCK2}`, without retyping its path:

::task{name="step_4"}
#active
Waiting for your shell back at `${DOCK2}`...
#completed
Circuit complete: absolute jump, relative hop, bare `cd`, `cd -`. That
is the entire navigation toolkit, and it is now in your fingers - not
your notes.
::
