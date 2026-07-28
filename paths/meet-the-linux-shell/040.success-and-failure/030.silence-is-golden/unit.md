---
title: Silence is golden
tasks:
  napped:
    check: |
      wait_exec '(^|/)sleep 2s?$'
    hint: |
      echo "Run: sleep 2"
    solve: |
      sleep 2
  reported:
    needs: [napped]
    check: |
      wait_exec '(^|/)whoami$'
    hint: |
      echo "sleep printed nothing but its exit status is 0 - success. So the correct report is whoami."
    solve: |
      echo $?
      whoami
---

Remember `sleep`? It prints nothing at all. So... did it *work*? Don't
guess - ask for the verdict.

1. Run `sleep 2`.
2. Reveal its exit status right away.
3. Report: if the status is `0`, run `whoami`; anything else, run
   `hostname`.

::task{name="napped"}
#active
Waiting for the two-second nap...
#completed
Done napping. Now: success or failure? The screen won't tell you - the
exit status will.
::

::task{name="reported"}
#active
Waiting for your report: `whoami` for `0`, `hostname` otherwise...
#completed
Status `0`: total success, in total silence. **No output does not mean
failure** - the exit status is the truth, and plenty of well-behaved
commands say nothing when everything goes fine.
::
