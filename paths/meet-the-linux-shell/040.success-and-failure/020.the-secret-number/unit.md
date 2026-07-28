---
title: The secret number
vars:
  BOGUS: { pick: [pineapple, backwards, moonshine] }
tasks:
  broke_it:
    check: |
      wait_exec "(^|/)date --${BOGUS}\$"
    hint: |
      echo "First run the failing command: date --${BOGUS}"
    solve: |
      date --$BOGUS
  reported:
    needs: [broke_it]
    check: |
      wait_exec '(^|/)hostname$'
    hint: |
      echo "Reveal the number with: echo \$? - since the command failed, it is NOT zero. Act accordingly."
    solve: |
      echo $?
      hostname
---

Besides its output, every command leaves behind a hidden verdict: a
number called the **exit status**. `0` means "all went well"; anything
else signals a problem. The shell stores the latest verdict in `$?`,
and `echo` can reveal it:

```
echo $?
```

Each command overwrites the verdict - so read it *right after* the
command you care about.

Now, the assignment:

1. Run `date` with the bad option `--${BOGUS}` again.
2. Immediately reveal its exit status.
3. Report your finding: if the number is `0`, run `whoami`;
   if it is anything else, run `hostname`.

::task{name="broke_it"}
#active
Waiting for the failing `date --${BOGUS}`...
#completed
Failed, as ordered. Now check the verdict it left behind.
::

::task{name="reported"}
#active
Waiting for your report: `whoami` if the status was `0`, `hostname`
otherwise...
#completed
It was `1` - a failure verdict, so `hostname` was the right report.
You'll almost never *print* `$?` in daily work, but everything in the
rest of this module quietly runs on it.
::
