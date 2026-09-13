---
title: Does it even exist?
requires: [readline]
vars:
  NAME: { pick: [snarkle, blorf, sort, tac] }
tasks:
  looked_up:
    timeout: 45
    check: |
      LINE=$(wait_line --latest "^((type|command|which)( +-[a-zA-Z]+)* +)?${NAME}\b") || exit 1
      case "$LINE" in
        type\ *) exit 0 ;;
        command\ *-[vV]\ *) exit 0 ;;
        command\ *) hint_exit "That line runs ${NAME} instead of asking about it. The -v option turns command into a lookup: command -v ${NAME}." ;;
        which\ *) hint_exit "The which tool only searches for program files and cannot see builtins. Ask with type or command -v instead. They see both kinds." ;;
        *) hint_exit "That line ran ${NAME} instead of asking about it. If the prompt has not come back, press Ctrl-C. Then look the name up with type or command -v." ;;
      esac
    hint: |
      echo "Ask about the name instead of running it: type ${NAME}, or command -v ${NAME}."
    solve: |
      type $NAME
      command -v $NAME
---

A teammate mentions a command named `${NAME}`. Before you use it, you
want to know whether it is installed on this machine.

You could type the name and see whether the shell reports *command not
found*. But that is a risky way to check. If the command exists, it runs,
and you do not know yet what it does. It may change or delete something
before you can react, or it may wait for input and never give the
prompt back.

Ask about the name instead of running it. The `type` builtin answers
without running anything. The `command -v` builtin does the same but in a less verbose way -
for a name that does not exist it prints nothing and exits with a non-zero status.

Find out whether a command named `${NAME}` exists here, without running
it.

::task
#active
Waiting for a lookup of `${NAME}` that does not run it...
#completed
The lookup answered without running `${NAME}`.
::

::tip{title="Why two ways to ask the same question"}
For a program file, `command -v` prints the location of the file. For
a builtin, it prints the name itself. And when the name is unknown,
it prints nothing and exits with a non-zero status. The short output and
the exit status make `command -v` the usual choice in shell scripts.
Future Shell Gyms will cover that use.
::
