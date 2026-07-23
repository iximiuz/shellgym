---
title: "Warm-up lap: navigate and collect"
vars:
  BASECAMP: { pick: [basecamp-1, basecamp-2, basecamp-3] }
tasks:
  camp_ready:
    check: |
      wait_file "/tmp/trek/$BASECAMP/inventory.txt"
    hint: |
      echo "Create the /tmp/trek/${BASECAMP} directory first (mkdir -p makes the whole chain), then produce inventory.txt inside it with a redirect."
    solve: |
      mkdir -p /tmp/trek/$BASECAMP
      cd /tmp/trek/$BASECAMP
      ls ~ > inventory.txt
  standing_in_camp:
    mode: level
    check: |
      wait_cwd --now "/tmp/trek/$BASECAMP"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "Finish the task while standing in the camp: your shell is in $CWD, expected /tmp/trek/${BASECAMP}."
---

A quick recap lap, no new commands. Combine what you practiced so far:

1. Create the directory `/tmp/trek/${BASECAMP}` (one command can do it).
2. Move your shell into it.
3. While standing there, create `inventory.txt` containing anything, for
   example the listing of your home directory.

::task{name="camp_ready"}
#active
Waiting for `/tmp/trek/${BASECAMP}/inventory.txt`...
#completed
Inventory in place.
::

::task{name="standing_in_camp"}
#active
And your shell must be standing in `/tmp/trek/${BASECAMP}` when the file is
there.
#completed
Rehearsed: `mkdir -p`, `cd`, and redirection in one breath.
::
