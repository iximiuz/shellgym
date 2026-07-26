# Home base

Every user has a home directory, and the shell gives it privileged
treatment: `cd` with no argument goes there, `~` abbreviates it in any
path, and the `$HOME` variable holds it for scripts and commands alike.

One more shortcut rounds out the set: `cd -` returns to wherever you
were *before* the last move - the "back" button of the terminal.
