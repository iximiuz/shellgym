---
title: Refresh the package index (apt)
labels: [ubuntu, debian]
tasks:
  index_refreshed:
    check: |
      wait_exec '(^|/)(sudo )?apt(-get)? update'
    hint: |
      echo "The command is apt update, and package management needs root: prefix it with sudo."
    solve: |
      sudo apt-get update
---

Before installing anything, `apt` needs a fresh index of what is
available. Refreshing it is the customary first move on any Debian or
Ubuntu box. Package management requires root, so `sudo` comes first:

::task{name="index_refreshed"}
#active
Waiting for you to refresh the package index...
#completed
Index refreshed. From here, `apt install <name>` installs and `apt search
<term>` explores.
::
