---
title: The hidden passage
vars:
  VAULT: { pick: [vault, cache, stash] }
init:
  - name: create_vault
    run: |
      mkdir -p "$GYM_USER_HOME/.$VAULT"
      chown "$GYM_USER" "$GYM_USER_HOME/.$VAULT"
tasks:
  inside:
    check: |
      wait_cwd "$GYM_USER_HOME/.$VAULT"
    hint: |
      echo "Hidden only means the name starts with a dot. cd ~/.${VAULT} works exactly like any other path - type the dot as part of the name."
    solve: |
      cd ~/.$VAULT
---

A directory whose name starts with a dot - like `.${VAULT}` - is
*hidden*: plain listings skip it. That is a display convention, nothing
more. `cd` does not care.

There is a hidden directory `.${VAULT}` in your home. Step inside:

::task{name="inside"}
#active
Waiting for your shell in `~/.${VAULT}`...
#completed
Inside the hidden passage. Half of your home directory is dot-dirs -
`.ssh`, `.config`, `.cache` - and now none of them can hide from you.
::

::tip
---
title: Hidden from Tab too?
---
Completion still finds dot-names - but only after you type the leading
dot: `cd ~/.` then `Tab` lists every hidden entry.
::
