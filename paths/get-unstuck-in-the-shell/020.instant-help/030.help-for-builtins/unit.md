---
title: Help for builtins
requires: [readline]
tasks:
  documented:
    timeout: 45
    check: |
      wait_line '^help( +-[a-zA-Z]+)* +echo\b'
    hint: |
      echo "The --help option is answered by a program, and echo is a builtin. Ask the shell for its documentation: help echo."
    solve: |
      echo --help
      help echo
---

Ask `echo` for its help the usual way:

```
echo --help
```

It printed `--help` back. The `--help` option is answered by the
program behind a command, and `echo` is a builtin. There is no program
behind it, so the shell ran its own `echo`, which treats `--help` as
one more word to print. Check with `type echo` if you like.

The shell has its own documentation command for builtins: the `help`
builtin. It only knows builtins. Ask it about a program such as `date`
and it fails with *no help topics match* and a non-zero exit status.

Get the documentation of `echo` from the shell.

::task
#active
Waiting for you to ask the shell for the documentation of `echo`...
#completed
The shell documented its own `echo`.
::

::tip
**Tab** completion works for command names too. Type `hel` and press
Tab. It saves typing and catches misspellings before Enter does.
::
