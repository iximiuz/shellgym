# Pattern vision

Type `*.log` and the *shell* - before the command even runs - replaces
it with every matching name in the directory. That substitution is
called globbing, and it is why `ls *.log` works with no special support
from `ls`.

Three patterns cover most days: `*` matches any run of characters, `?`
matches exactly one, and `[...]` matches one character from a set.
Learn to predict the expansion and you will never fear putting a glob
in a command again.
