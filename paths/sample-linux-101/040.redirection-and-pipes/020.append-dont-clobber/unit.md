---
title: Append without overwriting
needs: [write-it-down]
vars:
  MOTTO: { from: write-it-down.MOTTO }
tasks:
  appended:
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_file_contains "$HOME_DIR/motto.txt" "^$MOTTO\n.+"
    hint: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      LINES=$(wc -l < "$HOME_DIR/motto.txt" 2>/dev/null || echo 0)
      FIRST=$(head -1 "$HOME_DIR/motto.txt" 2>/dev/null || true)
      if [ "$FIRST" != "$MOTTO" ]; then
        echo "The first line of motto.txt is no longer '${MOTTO}' - it looks like the file got overwritten. Recreate the first line, then append with >> (two arrows)."
      elif [ "$LINES" -lt 2 ]; then
        echo "motto.txt still has a single line. A double arrow >> appends to the end instead of replacing."
      fi
    solve: |
      echo "second line" >> ~/motto.txt
---

Your `motto.txt` still holds the line `${MOTTO}`. Add a second line to it,
any text you like, without losing the first one. A single `>` would wipe
the file; the double `>>` appends.

::task{name="appended"}
#active
Waiting for a second line in `~/motto.txt` (first line intact)...
#completed
Appended. This check would have caught an accidental overwrite, a classic
slip when working fast.
::
