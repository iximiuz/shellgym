---
title: Take stock of the shed
vars:
  SHED: { pick: [toolshed, boathouse, greenhouse] }
init:
  - name: create_shed
    run: |
      rm -rf /tmp/gym-yard
      mkdir -p "/tmp/gym-yard/$SHED/crates"
      touch "/tmp/gym-yard/$SHED/hammer.txt" "/tmp/gym-yard/$SHED/rope.txt" "/tmp/gym-yard/$SHED/ladder.txt"
      chmod -R a+rx /tmp/gym-yard
tasks:
  get_there:
    check: |
      wait_cwd "/tmp/gym-yard/$SHED"
    hint: |
      echo "First cd to /tmp/gym-yard/${SHED} - navigation is a solved problem for you now."
    solve: |
      cd /tmp/gym-yard/$SHED
  look:
    needs: [get_there]
    check: |
      wait_exec '(^|/)ls( |$)'
    hint: |
      echo "Two letters: ls. It lists the contents of the current directory."
    solve: |
      ls
---

You have inherited a ${SHED} at `/tmp/gym-yard/${SHED}` and nobody
told you what is inside. Go there:

::task{name="get_there"}
#active
Waiting for your shell in `/tmp/gym-yard/${SHED}`...
#completed
At the door.
::

Now look around. `ls` lists what the current directory contains:

::task{name="look"}
#active
Waiting for you to list the shed's contents...
#completed
Tools and a crate. Notice `crates` is a directory and the rest are
files - many terminals color directories differently, and `ls -F` marks
them with a trailing `/`.
::

::tip
---
title: The arrival reflex
---
`cd` somewhere, then `ls` - the pair is so common it should fire as one
motion. You will do it thousands of times.
::
