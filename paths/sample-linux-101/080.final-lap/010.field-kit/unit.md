---
title: Assemble a field kit
vars:
  KIT: { pick: [kit-alpha, kit-bravo, kit-charlie] }
  WORD: { pick: [compass, lantern, rope, flint] }
init:
  - name: seed_manifest
    run: |
      rm -f /tmp/manifest.txt
      for i in $(seq 1 30); do echo "item $i: ration" >> /tmp/manifest.txt; done
      for i in $(seq 1 7); do echo "special $i: $WORD" >> /tmp/manifest.txt; done
      chmod a+r /tmp/manifest.txt
tasks:
  kit_dir:
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_file "$HOME_DIR/$KIT"
    solve: |
      mkdir -p ~/$KIT
  manifest_copied:
    needs: [kit_dir]
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_file "$HOME_DIR/$KIT/manifest.txt"
    hint: |
      echo "The original manifest lives at /tmp/manifest.txt. Copy it, keeping the name."
    solve: |
      cp /tmp/manifest.txt ~/$KIT/
  specials_counted:
    needs: [manifest_copied]
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_file_contains "$HOME_DIR/$KIT/specials.txt" "^7\s*$"
    hint: |
      echo "Count the manifest lines mentioning '${WORD}' and put that number into specials.txt. A pipe from a filter into a line counter did this for you once before."
    solve: |
      grep -c $WORD ~/$KIT/manifest.txt > ~/$KIT/specials.txt
---

Assemble a field kit from memory:

1. Create a directory `${KIT}` in your home directory.
2. Copy `/tmp/manifest.txt` into it.
3. Count how many manifest lines mention `${WORD}` and store the number in
   `~/${KIT}/specials.txt`.

::task{name="kit_dir"}
#active
Waiting for `~/${KIT}`...
#completed
Kit directory ready.
::

::task{name="manifest_copied"}
#active
Waiting for the manifest copy...
#completed
Manifest secured.
::

::task{name="specials_counted"}
#active
Waiting for the special-item count in `specials.txt`...
#completed
Full kit, correct count. That was `mkdir`, `cp`, `grep`, `wc`, and
redirection - all from memory.
::
