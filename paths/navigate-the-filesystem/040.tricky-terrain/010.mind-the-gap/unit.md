---
title: Mind the gap
vars:
  ALBUM: { pick: [release notes, meeting minutes, field reports] }
init:
  - name: create_shelf
    run: |
      rm -rf /tmp/gym-shelf
      mkdir -p "/tmp/gym-shelf/$ALBUM"
      chmod -R a+rx /tmp/gym-shelf
tasks:
  arrive:
    check: |
      wait_cwd "/tmp/gym-shelf/$ALBUM"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. The space in '${ALBUM}' splits an unquoted path into two arguments. Wrap the path in quotes, or escape the space with a backslash."
    solve: |
      cd "/tmp/gym-shelf/$ALBUM"
---

Someone created a directory called `${ALBUM}` - with a space in the
middle - under `/tmp/gym-shelf`. It happens all the time.

The catch: the shell splits command lines on whitespace, so an unquoted
space ends one argument and starts another. `cd` would receive two
half-names and find neither.

Get your shell inside `/tmp/gym-shelf/${ALBUM}`:

::task{name="arrive"}
#active
Waiting for your shell in `/tmp/gym-shelf/${ALBUM}`...
#completed
In. Quotes hold the pieces together into one argument - the shell never
even sees a "special" character inside them.
::

::tip
---
title: Three ways past a space
---
Quote the whole path (`cd "/tmp/a b"`), escape just the space
(`cd /tmp/a\ b`), or type a few letters and press `Tab` - completion
inserts the escaping for you.
::
