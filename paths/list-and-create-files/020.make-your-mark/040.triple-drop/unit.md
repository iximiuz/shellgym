---
title: Three at once
vars:
  CREW: { pick: [alpha, bravo, delta] }
init:
  - name: create_rosters
    run: |
      rm -rf /tmp/gym-rosters
      mkdir -p /tmp/gym-rosters
      chown "$GYM_USER" /tmp/gym-rosters
tasks:
  all_three:
    check: |
      wait_file "/tmp/gym-rosters/$CREW-mon.txt"
      wait_file "/tmp/gym-rosters/$CREW-wed.txt"
      wait_file "/tmp/gym-rosters/$CREW-fri.txt"
    hint: |
      MISSING=""
      for d in mon wed fri; do
        [ -e "/tmp/gym-rosters/$CREW-$d.txt" ] || MISSING="$MISSING $CREW-$d.txt"
      done
      echo "Still missing:$MISSING. touch accepts several paths in one command - list them all, separated by spaces."
    solve: |
      touch /tmp/gym-rosters/$CREW-mon.txt /tmp/gym-rosters/$CREW-wed.txt /tmp/gym-rosters/$CREW-fri.txt
---

Crew `${CREW}` trains three times a week, and each session needs a
roster file in `/tmp/gym-rosters`:

- `${CREW}-mon.txt`
- `${CREW}-wed.txt`
- `${CREW}-fri.txt`

`touch` happily takes **several** names in a single command - every
argument becomes a file. (Three separate commands work too, but why
type three times?)

::task{name="all_three"}
#active
Waiting for all three roster files...
#completed
Three files, ideally one command. Most file commands - `touch`,
`mkdir`, and the copy/move family you'll meet next path - accept
multiple arguments; batching is the shell's home turf.
::

::tip
---
title: Recycle the line
---
If you do go one file at a time: Up arrow brings the last command back,
and `Ctrl-W` deletes the word before the cursor - edit `mon` into
`wed` and go.
::
