# What Runs When You Type

Every command starts with a name, and the shell has to decide what that
name means before anything runs. Most names point at a program file on
disk. Some name a builtin, and the shell does the work itself. A few
names exist in both forms at once: `echo` and `pwd` are builtins, and
`/bin/echo` and `/bin/pwd` are separate programs.

Most of the time you will not care which of the two you are running.
Sometimes, though, the behavior differs: the two versions of one name
may accept different options, and getting help for a builtin takes a
different action than getting help for a program.

The builtin `type` tells you which one a name runs. In this module you
will use it to look names up, run the version you mean, and check
whether a name exists at all without running it.
