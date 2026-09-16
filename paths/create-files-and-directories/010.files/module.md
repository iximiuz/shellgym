# Files

The `touch` command takes a path and creates an empty file there. If a
file already exists at that path, `touch` leaves its contents alone and
only updates its modification time to now. It accepts several paths in
one command, absolute or relative, and it is subject to the quoting
rules you already know: a name with a space in it needs quotes to stay
one argument.

The units use the same small `projects` tree in your home directory as
[List Files and Use Patterns](https://labs.iximiuz.com/shell-gyms/list-files-and-use-patterns).
Every unit creates it if it is missing, so nothing you do can break it:

```
~/projects
├── archive
│   └── 2025
├── backups
├── data
│   ├── exports
│   └── imports
├── notes
└── reports
    ├── drafts
    └── final
```
