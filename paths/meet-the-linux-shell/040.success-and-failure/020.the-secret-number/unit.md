---
title: The secret number
vars:
  BOGUS: { pick: [pineapple, backwards, moonshine] }
tasks:
  broke_it:
    check: |
      wait_exec "(^|/)date --${BOGUS}\$" || exit 1
      set_var BROKE_SEQ "$(event_seq)"
    hint: |
      echo "First run the failing command: date --${BOGUS}"
    solve: |
      date --$BOGUS
  reported:
    needs: [broke_it]
    check: |
      # only a report given after the failing date counts
      REPORT=$(wait_exec --after "$BROKE_SEQ" --latest '(^|/)(hostname|whoami)$') || exit 1
      [[ "$REPORT" == *hostname* ]] && exit 0
      hint_exit "You reported with whoami, but whoami is the answer for a status of 0, and the failing date left a different number. Reveal it again with echo \$? right after the failing command, then report accordingly."
    hint: |
      echo "Reveal the number with echo \$? right after the failing command. Since the command failed, the number is not zero, so report accordingly."
    solve: |
      echo $?
      hostname
---

Besides its output, every command leaves behind a hidden verdict: a
number called the **exit status**. The value `0` means "all went well". Anything
else signals a problem. The shell stores the latest verdict in `$?`,
and `echo` can reveal it:

```
echo $?
```

Each command overwrites the verdict, so read it *right after* the
command you care about.

The assignment has three steps:

1. Run `date` with the bad option `--${BOGUS}` again.
2. Immediately reveal its exit status.
3. Report your finding. If the number is `0`, run `whoami`. If it is
   anything else, run `hostname`.

::task{name="broke_it"}
#active
Waiting for the failing `date --${BOGUS}`...
#completed
It failed, as ordered. Now check the verdict it left behind.
::

::task{name="reported"}
#active
Waiting for your report: `whoami` if the status was `0`, `hostname`
otherwise...
#completed
Correct, the status was `1`, a failure verdict.
::
