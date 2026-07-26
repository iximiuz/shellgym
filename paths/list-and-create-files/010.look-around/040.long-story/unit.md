---
title: The long story
init:
  - name: create_archive
    run: |
      rm -rf /tmp/gym-archive
      mkdir -p /tmp/gym-archive
      head -c 48 /dev/zero > /tmp/gym-archive/stub.dat
      head -c 4096 /dev/zero > /tmp/gym-archive/bundle.dat
      head -c 12288 /dev/zero > /tmp/gym-archive/payload.dat
      touch -d "3 days ago" /tmp/gym-archive/stub.dat
      touch -d "yesterday" /tmp/gym-archive/payload.dat
      touch -d "2 hours ago" /tmp/gym-archive/bundle.dat
      touch /tmp/gym-archive/manifest.txt
      chmod -R a+rx /tmp/gym-archive
tasks:
  go_look:
    check: |
      wait_cwd /tmp/gym-archive
    hint: |
      echo "Head to /tmp/gym-archive first."
    solve: |
      cd /tmp/gym-archive
  detail_view:
    needs: [go_look]
    check: |
      wait_exec '(^|/)ls.* (-[a-zA-Z]*l[a-zA-Z]*|--format=long)( |$)'
    hint: |
      echo "The long format is one short option away: l is for 'long'. Each line shows permissions, owner, size, and modification time."
    solve: |
      ls -l
  by_time:
    needs: [detail_view]
    check: |
      wait_exec '(^|/)ls.* (-[a-zA-Z]*t[a-zA-Z]*|--sort=time)( |$)'
    hint: |
      echo "Add t to sort by modification time, newest first - -lt is the everyday spelling."
    solve: |
      ls -lt
---

`/tmp/gym-archive` holds a few data files. A bare `ls` gives you names;
sometimes you need the whole story - who owns what, how big it is,
when it changed.

Go to the archive:

::task{name="go_look"}
#active
Waiting for your shell in `/tmp/gym-archive`...
#completed
Four files, no details - yet.
::

Get the **long** listing - the one with a line of metadata per file:

::task{name="detail_view"}
#active
Waiting for a long-format listing...
#completed
Left to right: permissions, link count, owner, group, size in bytes,
modification time, name. `payload.dat` is the biggest. Add `-h` next
time and sizes turn human-readable (`12K` instead of `12288`).
::

Which file changed most recently? Don't squint at dates - ask for the
listing **sorted by modification time**, newest first:

::task{name="by_time"}
#active
Waiting for a time-sorted listing...
#completed
`manifest.txt` on top - it was touched moments ago, when the scene was
built. "What changed here recently?" is a question `ls -lt` answers
daily in real operations.
::

::tip
---
title: Combining short options
---
Short options stack behind one dash: `-l` plus `-t` is `-lt`, plus
`-h` is `-lth`. Order rarely matters.
::
