---
title: Silence is golden
tasks:
  napped:
    check: |
      wait_exec '(^|/)sleep 2s?$' || exit 1
      set_var NAP_SEQ "$(event_seq)"
    hint: |
      echo "Run: sleep 2"
    solve: |
      sleep 2
  reported:
    needs: [napped]
    check: |
      # only a report given after the nap counts
      wait_exec --after "$NAP_SEQ" '(^|/)whoami$'
    hint: |
      echo "sleep printed nothing, but its exit status is 0, which means success. So the correct report is whoami."
    solve: |
      echo $?
      whoami
---

Remember `sleep`? It prints nothing at all, so there is no visible sign
of whether it *worked*. Ask for the verdict instead of guessing.

1. Run `sleep 2`.
2. Reveal its exit status right away.
3. Report your finding. If the status is `0`, run `whoami`. If it is
   anything else, run `hostname`.

::task{name="napped"}
#active
Waiting for the two-second nap...
#completed
The nap is over. Now ask for the verdict.
::

::task{name="reported"}
#active
Waiting for your report: `whoami` for `0`, `hostname` otherwise...
#completed
Correct, the status was `0`, success in total silence.
::

::tip{title="No output does not mean failure"}
Plenty of well-behaved commands say nothing when everything goes
fine. The exit status is the only universal way to tell.
::
