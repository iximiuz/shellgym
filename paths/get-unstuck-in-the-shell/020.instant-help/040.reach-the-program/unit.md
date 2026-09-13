---
title: "Bonus: reach the program"
tasks:
  program_help:
    timeout: 45
    check: |
      wait_exec '(^|/)env echo --help$'
    hint: |
      echo "Put env in front of the whole command: env echo --help. The env program searches the usual program locations for echo and runs the program version."
    solve: |
      type -a echo
      env echo --help
---

The name `echo` has two versions on this machine. Run `type -a echo`
to see both. The `-a` option lists every version of a name in the
order the shell would try them: first the builtin, then the program
file.

::note{title="The file may be listed twice"}
On many Linux distributions `/bin` is only a link to `/usr/bin`, so
the same program file can show up under both locations. Either way,
for a bare name the shell always takes the first entry.
::

The program version of `echo` answers `--help` like any other program. One way
to reach it is the full location from `type -a`, as in the first
module. Another way is the `env` program: it searches the known
program locations for the name you give it and runs the program file
it finds. Being a program itself, `env` does not know about builtins,
so `env echo` can only mean the program version.

Get the help of the program version of `echo` through `env`.

::task
#active
Waiting for `env` to run the program version of `echo` with `--help`...
#completed
That help page came from the program version of `echo`.
::
