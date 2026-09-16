# Patterns

Typing every file name by hand does not scale. The shell can name a
whole group of files with a **pattern**: `*` stands for any sequence of
characters (including none), `?` stands for exactly one character, and
`[...]` stands for one character out of a set. A pattern such as
`*.csv` is turned into the list of matching names before the command
runs. The command never sees the pattern, only the names.

Because the shell does the matching, patterns work with every command:
`ls *.csv` and `echo *.csv` receive the same list, and so will every
command you learn later.

The units use the familiar `~/projects/data` tree from the previous module to let you practice with name patterns:

```
~/projects/data
├── exports
├── imports
├── sales-2025.csv
├── sales-q1.csv
├── sales-q2.csv
...
```
