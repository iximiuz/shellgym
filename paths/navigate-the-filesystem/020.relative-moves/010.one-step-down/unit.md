---
title: One step down
vars:
  ROOM: { pick: [assets, content, styles] }
init:
  - name: create_site
    run: |
      rm -rf /tmp/gym-site
      mkdir -p "/tmp/gym-site/$ROOM" /tmp/gym-site/shared
      chmod -R a+rx /tmp/gym-site
tasks:
  at_site:
    check: |
      wait_cwd /tmp/gym-site
    hint: |
      echo "Start at the top: cd to /tmp/gym-site with an absolute path."
    solve: |
      cd /tmp/gym-site
  down_one:
    needs: [at_site]
    check: |
      wait_cwd "/tmp/gym-site/$ROOM"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. From /tmp/gym-site, the ${ROOM} directory is one step below - name it without any leading slash and cd resolves it relative to where you stand."
    solve: |
      cd $ROOM
---

A small website project lives in `/tmp/gym-site`. Enter the project
first:

::task{name="at_site"}
#active
Waiting for your shell in `/tmp/gym-site`...
#completed
At the project root.
::

Now descend into its `${ROOM}` directory - but this time skip the
`/tmp/gym-site` part. A path that does not start with `/` is resolved
from your current directory.

::task{name="down_one"}
#active
Waiting for your shell in the `${ROOM}` subdirectory...
#completed
That is a relative move: three words shorter than the absolute one, and
it works the same in any project you ever enter.
::
