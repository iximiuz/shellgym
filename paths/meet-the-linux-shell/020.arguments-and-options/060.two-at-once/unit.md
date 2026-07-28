---
title: Two options, one dash
tasks:
  combined:
    check: |
      wait_exec '(^|/)date -(uR|Ru)$'
    hint: |
      echo "One dash, then both letters back to back - order doesn't matter."
    solve: |
      date -uR
---

`date` knows another short option: `-R` prints the time in the classic
e-mail header style (`Mon, 28 Jul 2025 ...`).

Here's a convenience worth knowing: **compatible short options can share
one dash**. Instead of writing `-u -R` as two arguments, you can squeeze
both letters into a single dash group.

Print the e-mail-style time in UTC, using one dash for both options.

::task
#active
Waiting for `-u` and `-R` combined into a single dash group...
#completed
Two behaviors switched on with one compact argument. Long options never
combine like this - one word per `--` option.
::
