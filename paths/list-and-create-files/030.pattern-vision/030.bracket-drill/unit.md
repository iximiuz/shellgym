---
title: Pick from a set
init:
  - name: create_snaps
    run: |
      rm -rf /tmp/gym-snaps
      mkdir -p /tmp/gym-snaps
      touch /tmp/gym-snaps/snap-1.dat /tmp/gym-snaps/snap-2.dat /tmp/gym-snaps/snap-7.dat
      touch /tmp/gym-snaps/snap-a.dat /tmp/gym-snaps/snap-b.dat
      chmod -R a+rx /tmp/gym-snaps
tasks:
  in_place:
    check: |
      wait_cwd /tmp/gym-snaps
    hint: |
      echo "cd to /tmp/gym-snaps first."
    solve: |
      cd /tmp/gym-snaps
  digits_only:
    needs: [in_place]
    check: |
      wait_exec '(^|/)(ls|echo).* (\./)?snap-1\.dat (\./)?snap-2\.dat (\./)?snap-7\.dat$'
    hint: |
      echo "? matches ANY single character - letters included. Square brackets restrict the match to a set: [0-9] means 'one digit'."
    solve: |
      ls snap-[0-9].dat
---

The snapshots in `/tmp/gym-snaps` come in two flavors: numbered
(`snap-1`, `snap-2`, `snap-7`) and lettered (`snap-a`, `snap-b`). Get
there:

::task{name="in_place"}
#active
Waiting for your shell in `/tmp/gym-snaps`...
#completed
Numbers and letters, all `.dat`.
::

List **only the numbered** snapshots. `snap-?.dat` matches all five -
`?` doesn't care whether the character is a digit. Square brackets do:
`[0-9]` matches one character *from that range*.

::task{name="digits_only"}
#active
Waiting for a listing of just the numbered snapshots...
#completed
Character classes finish the glob trio: `*` for any run, `?` for any
one, `[...]` for one from a set. `[a-z]`, `[0-9]`, `[abc]` - the
bracket takes ranges and lists alike.
::
