# Look Around

The `ls` command lists what a directory contains. In
[Navigate the Filesystem](https://labs.iximiuz.com/shell-gyms/navigate-the-filesystem)
you used it to list the files in the current working directory.
This module covers the rest of the command: list a directory without moving into it, list several
directories at once, show the hidden names, read the size, owner, and date
of every entry, and sort the listing so the entry you need comes first.

The units use a small `projects` tree in your home directory. Every unit
creates it if it is missing, so nothing you do can break it:

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
