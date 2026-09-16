---
title: A hidden workspace
vars:
  NAME: { pick: [scratchpad, worklog, playbook, toolbox] }
init:
  - name: clear_target
    run: |
      rm -rf "$GYM_USER_HOME/.$NAME"
tasks:
  directory:
    check: |
      wait_dir "$GYM_USER_HOME/.$NAME"
    hint: |
      echo "Create ~/.${NAME} with mkdir. The dot is part of the name."
    solve: |
      mkdir ~/.$NAME
  file:
    needs: [directory]
    timeout: 60
    check: |
      wait_file "$GYM_USER_HOME/.$NAME/settings"
    hint: |
      echo "Create an empty file named settings inside ~/.${NAME}."
    solve: |
      touch ~/.$NAME/settings
  listed:
    needs: [file]
    timeout: 60
    check: |
      ARGV=$(wait_exec --latest "(^|/)ls( +\S+)* +\S*\.${NAME}/? *\$") || exit 1
      S=' -[a-zA-Z0-9]*l[a-zA-Z0-9]* '
      if [[ " $ARGV " =~ $S ]]; then exit 0; fi
      hint_exit "That listing shows names only. List the contents of the directory in the long format to see the size of settings."
    hint: |
      echo "List the contents of ~/.${NAME} in the long format. The path is the argument, so your shell can stay where it is."
    solve: |
      ls -l ~/.$NAME
---

Prepare a hidden workspace in your home directory: create a hidden
directory named `.${NAME}`, create an empty file named `settings`
inside it, and list its contents in the long format to confirm the
file is there with a size of `0`.

::task{name="directory"}
#active
Waiting for `~/.${NAME}` to appear...
#completed
The hidden directory exists. Now put the file in it.
::

::task{name="file"}
#active
Waiting for `settings` to appear inside `~/.${NAME}`...
#completed
The file is in place. Now confirm it with a long listing.
::

::task{name="listed"}
#active
Waiting for the contents of `~/.${NAME}` in the long format...
#completed
There it is, with a size of `0` and a fresh date.
::
