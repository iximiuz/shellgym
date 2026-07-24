---
title: Copy a file
vars:
  TOKEN: { shell: "head -c4 /dev/urandom | od -An -tx1 | tr -d ' \\n'" }
init:
  - name: seed_original
    run: |
      mkdir -p /tmp/vault
      printf 'artifact %s\n' "$TOKEN" > /tmp/vault/original.txt
      chmod a+r /tmp/vault/original.txt
tasks:
  copied:
    check: |
      wait_file_contains "$GYM_USER_HOME/original.txt" "artifact $TOKEN"
    hint: |
      if [ -f "$GYM_USER_HOME/original.txt" ]; then
        echo "A file named original.txt is in your home directory, but its content differs from /tmp/vault/original.txt. Copy the file again."
      else
        echo "cp takes two arguments: the source path and the destination. The destination can be a directory."
      fi
    solve: |
      cp /tmp/vault/original.txt ~/
---

There is a file at `/tmp/vault/original.txt`. Copy it into your home
directory, keeping the same file name. The copy command is `cp source
destination`.

::task{name="copied"}
#active
Waiting for a faithful copy of `original.txt` in your home directory...
#completed
Copied byte for byte. `cp -r` does the same for whole directories.
::
