---
title: Builtin or program?
requires: [readline]
vars:
  BUILTIN: { pick: [cd, pwd, echo, export] }
  PROGRAM: { pick: [seq, sleep, cat, tac] }
tasks:
  builtin_asked:
    timeout: 45
    check: |
      wait_line "^type( +-[a-zA-Z]+)* +${BUILTIN} *$"
    hint: |
      echo "Ask the shell: type ${BUILTIN}. The answer is 'is a shell builtin'."
    solve: |
      type $BUILTIN
  program_asked:
    timeout: 45
    check: |
      wait_line "^type( +-[a-zA-Z]+)* +${PROGRAM} *$"
    hint: |
      echo "Ask the shell: type ${PROGRAM}. The answer is the location of a program file."
    solve: |
      type $PROGRAM
---

Some names run no program file at all. The shell handles them itself,
and `type` reports *is a shell builtin* for them. This is also why
`which`, the tool you learned about in [Meet the Linux Shell](https://labs.iximiuz.com/shell-gyms/meet-the-linux-shell),
cannot find `cd`: it only searches for program files, and `cd` is a builtin.

The `type` command sees both kinds, which makes it the right first question to ask
about any name.

Run `type` twice, once for `${BUILTIN}` and once for `${PROGRAM}`, and
compare the two answers. One of them is *is a shell builtin*. The other
one is the location of a program file.

::task{name="builtin_asked"}
#active
Waiting for you to ask `type` about `${BUILTIN}`...
#completed
That one is a builtin.
::

::task{name="program_asked"}
#active
Waiting for you to ask `type` about `${PROGRAM}`...
#completed
That one is a program file.
::
