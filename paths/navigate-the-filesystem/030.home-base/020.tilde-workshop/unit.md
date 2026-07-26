---
title: The tilde shortcut
vars:
  BENCH: { pick: [drafts, blueprints, sketches] }
init:
  - name: create_workshop
    run: |
      mkdir -p "$GYM_USER_HOME/workshop/$BENCH"
      chown -R "$GYM_USER" "$GYM_USER_HOME/workshop"
tasks:
  at_bench:
    check: |
      wait_cwd "$GYM_USER_HOME/workshop/$BENCH"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. The ~ character expands to your home directory in any path: ~/workshop/${BENCH} names the bench without spelling out /home/..."
    solve: |
      cd ~/workshop/$BENCH
---

A workshop has been set up inside your home directory, with a
`${BENCH}` bench in it.

You could type the full absolute path - but the shell expands `~` to
your home directory wherever it starts a path. `~/workshop` means "the
workshop directory in my home", whoever you are and wherever home is.

Get to the bench using the tilde:

::task{name="at_bench"}
#active
Waiting for your shell at `~/workshop/${BENCH}`...
#completed
`~` saves you from ever typing `/home/yourname` again - and the same
command works for any user on any machine.
::

::tip
---
title: Where is the tilde?
---
On most keyboards `~` shares a key with the backtick, left of the `1`
key. Worth finding by touch - you will type it constantly.
::
