---
title: One character exactly
init:
  - name: create_scans
    run: |
      rm -rf /tmp/gym-scans
      mkdir -p /tmp/gym-scans
      touch /tmp/gym-scans/scan1.txt /tmp/gym-scans/scan2.txt /tmp/gym-scans/scan3.txt
      touch /tmp/gym-scans/scan10.txt /tmp/gym-scans/scan11.txt
      chmod -R a+rx /tmp/gym-scans
tasks:
  in_place:
    check: |
      wait_cwd /tmp/gym-scans
    hint: |
      echo "cd to /tmp/gym-scans first."
    solve: |
      cd /tmp/gym-scans
  single_digit:
    needs: [in_place]
    check: |
      wait_exec '(^|/)(ls|echo).* (\./)?scan1\.txt (\./)?scan2\.txt (\./)?scan3\.txt$'
    hint: |
      echo "A * after 'scan' would also swallow scan10 and scan11. The ? wildcard matches exactly ONE character - scan?.txt fits the single-digit scans only."
    solve: |
      ls scan?.txt
---

`/tmp/gym-scans` holds scans numbered 1 through 3 - plus `scan10` and
`scan11` from another batch. Head over:

::task{name="in_place"}
#active
Waiting for your shell in `/tmp/gym-scans`...
#completed
Two batches, one directory. Classic.
::

List **only the single-digit** scans. `scan*.txt` is too greedy - it
matches the teens too. You need the wildcard that stands for *exactly
one* character.

::task{name="single_digit"}
#active
Waiting for a listing of just `scan1` through `scan3`...
#completed
`?` is the precise sibling of `*`: one character, no more, no less.
`scan10` needs two characters after `scan`, so it stays out.
::
