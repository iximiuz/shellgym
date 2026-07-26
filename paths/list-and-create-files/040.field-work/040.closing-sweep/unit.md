---
title: Closing sweep
vars:
  STATION: { pick: [helios, borealis, meridian] }
init:
  - name: create_field
    run: |
      rm -rf /tmp/gym-field
      mkdir -p /tmp/gym-field
      chown "$GYM_USER" /tmp/gym-field
tasks:
  base_camp:
    check: |
      wait_dir "/tmp/gym-field/$STATION/samples"
    hint: |
      echo "Station first: the directory chain /tmp/gym-field/${STATION}/samples in one command - you know the option."
    solve: |
      mkdir -p /tmp/gym-field/$STATION/samples
  field_notes:
    needs: [base_camp]
    check: |
      wait_file "/tmp/gym-field/$STATION/.fieldlog"
      wait_file "/tmp/gym-field/$STATION/day one.txt"
    hint: |
      echo "Two files in /tmp/gym-field/${STATION}: a hidden .fieldlog (dot first) and 'day one.txt' (space - quote it). touch handles both, even in one command."
    solve: |
      cd /tmp/gym-field/$STATION
      touch .fieldlog "day one.txt"
  final_survey:
    needs: [field_notes]
    check: |
      wait_exec '(^|/)ls.* (-[a-zA-Z]*[aA][a-zA-Z]*|--all|--almost-all)( |$)'
    hint: |
      echo "Survey the station with a listing that shows EVERYTHING - hidden .fieldlog included. Long format plus all entries is the full picture."
    solve: |
      ls -la
---

Last assignment: set up research station `${STATION}` from scratch and
survey it. Everything in this rep is a move you already own.

**1.** Raise the base: the directory `/tmp/gym-field/${STATION}` with a
`samples` directory inside - one command.

::task{name="base_camp"}
#active
Waiting for `/tmp/gym-field/${STATION}/samples`...
#completed
Base camp standing.
::

**2.** Start the records in `/tmp/gym-field/${STATION}`: a hidden
`.fieldlog`, and a `day one.txt` - space in the name.

::task{name="field_notes"}
#active
Waiting for `.fieldlog` and `day one.txt` at the station...
#completed
Records open - a dot-name and a quoted name, back to back.
::

**3.** Full survey: produce a listing of the station that shows every
entry, hidden ones included.

::task{name="final_survey"}
#active
Waiting for the full listing...
#completed
Station built, records started, survey filed. `mkdir`, `touch`, `ls`
and the naming rules - that is the whole path, executed from memory.
Next stop: copying, moving, and removing what you now know how to
create.
::
