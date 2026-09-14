# One Tree

Every file on a Linux machine lives in a single tree of directories. The
tree starts at one directory called the **root**, written as `/`.
Every other directory is inside it, or inside a directory
inside it, and so on down. There are no drive letters. A disk that gets
attached to the machine simply shows up as a directory somewhere in
the same tree.

Your shell is always standing in one of those directories, its
**working directory**. The `pwd` command prints which one, and the `cd`
command moves the shell to another. This module walks you through the
important top-level directories and gets you moving between them by
their full names.
