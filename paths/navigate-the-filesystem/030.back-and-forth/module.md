# Back and Forth

Real work rarely happens in one directory. You go somewhere to look at
a file, come back, go somewhere else, and return again. Typing the
full path every time gets old fast.

The shell keeps three shortcuts for exactly this. A bare `cd` returns
home. The variable `$HOME` holds the path of your home directory, so
it can be used inside any path. And `cd -` takes you back to the
directory you were in just before the last move. This module
practices all three.

The units use the same `~/projects` tree as the previous module. It is
recreated on every unit if needed.
