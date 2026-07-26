---
title: Follow the trail
vars:
  L1: { pick: [gate, arch, portal] }
  L2: { pick: [tunnel, stair, ramp] }
  L3: { pick: [chamber, grotto, cellar] }
init:
  - name: create_trail
    run: |
      rm -rf /tmp/gym-trail
      mkdir -p "/tmp/gym-trail/$L1/$L2/$L3"
      chmod -R a+rx /tmp/gym-trail
tasks:
  at_the_end:
    check: |
      wait_cwd "/tmp/gym-trail/$L1/$L2/$L3"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. Type 'cd /tmp/gym-trail/' and press Tab - with exactly one way forward, completion fills in the next name. Repeat until the trail ends, three levels down."
    solve: |
      cd /tmp/gym-trail/$L1/$L2/$L3
---

A trail starts at `/tmp/gym-trail`. This time you are not told the
directory names - but the trail has exactly **one way forward at every
step**, and it is three levels deep.

You do not need to see a map: type `cd /tmp/gym-trail/` and press
`Tab`. With a single candidate, completion fills in the name instantly.
Keep pressing until there is nowhere deeper to go.

::task{name="at_the_end"}
#active
Waiting for your shell at the end of the trail, three levels below
`/tmp/gym-trail`...
#completed
Trail's end - and you likely never typed a single directory name.
Navigating an *unfamiliar* tree by Tab is exactly how experienced
operators move through systems they have never seen.
::

::tip
---
title: When Tab shows nothing
---
One `Tab` completes when there is a single match; two quick `Tab`s list
all candidates when there are several. If nothing appears, the name so
far has no matches - check for a typo.
::
