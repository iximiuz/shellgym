---
title: Set an environment variable
vars:
  COLOR: { pick: [teal, amber, violet, crimson] }
tasks:
  var_seen:
    check: |
      wait_env GYM_COLOR "^$COLOR$"
    hint: |
      echo "Set the variable with export GYM_COLOR=${COLOR} (no spaces around =), then run any command, e.g. env, so the new variable can be observed."
    solve: |
      export GYM_COLOR=$COLOR
      env
---

Environment variables are named values that every command you launch
inherits from the shell. The `export` keyword sets one for the current
shell session.

Set a variable named `GYM_COLOR` to `${COLOR}`, then run any command
(for example `env`, which prints the environment it received).

::task{name="var_seen"}
#active
Waiting to observe `GYM_COLOR=${COLOR}` in a command you run...
#completed
Inherited. This is exactly how settings like `PATH`, `HOME`, and `LANG`
reach every program you start. Note the check could not see the variable
until you launched a command - exported variables travel with new
processes, they are not broadcast.
::
