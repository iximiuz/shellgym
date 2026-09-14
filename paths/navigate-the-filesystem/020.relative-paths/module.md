# Relative Paths

A path that starts with `/` is an **absolute** path. It spells out the
whole way from the root, and it means the same thing no matter where
your shell stands.

Any other path is **relative**. The shell reads it starting from the
working directory, so `reports` means "the `reports` directory inside
the directory I am in right now". Two special names make relative
paths go anywhere: `..` is the parent of the current directory, and
`.` is the current directory itself. Names can be chained with slashes,
so `../data` means "up one, then into `data`".

This module uses a small work tree in your home directory. Every unit
creates it if it is missing, so nothing you do can break it:

```
~/work
├── archive
│   └── 2025
├── data
│   ├── exports
│   └── imports
├── notes
└── reports
    ├── drafts
    └── final
```
