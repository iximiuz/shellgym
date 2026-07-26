# Relative moves

An absolute path spells out the whole address. A **relative path** starts
from where you already are: `cd docs` enters the `docs` directory inside
the current one, `cd ..` climbs to the parent, and pieces combine freely -
`cd ../sibling/deeper` is a perfectly good move.

Two names are always defined: `.` is the current directory and `..` is
its parent. Around a project tree, relative moves are usually shorter
than absolute ones - this module makes them automatic.
