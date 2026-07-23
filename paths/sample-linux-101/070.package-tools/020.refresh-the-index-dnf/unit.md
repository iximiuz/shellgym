---
title: Refresh the package metadata (dnf)
labels: [rocky, fedora, rhel, centos, almalinux]
tasks:
  metadata_refreshed:
    check: |
      wait_exec '(^|/)(sudo )?dnf.*(makecache|check-update|repolist)'
    hint: |
      echo "dnf makecache downloads fresh repository metadata. Package management needs root, so use sudo."
    solve: |
      sudo dnf makecache
---

On Fedora-family systems (including Rocky Linux), `dnf` manages packages.
Refresh its repository metadata, the counterpart of `apt update`:

::task{name="metadata_refreshed"}
#active
Waiting for you to refresh the dnf metadata (`makecache`, `check-update`,
and `repolist` all count)...
#completed
Metadata refreshed. From here, `dnf install <name>` installs and `dnf
search <term>` explores.
::
