# Awkward Names

Some directory names need extra care on the command line. A name with
a space in it looks like two arguments to the shell. Parentheses and a
few other characters have a meaning of their own. And some names are
simply too long to type without a mistake.

You met the fix for spaces in [Meet the Linux Shell](https://labs.iximiuz.com/shell-gyms/meet-the-linux-shell): quotes keep a value
together. This module applies it to paths, adds the backslash as a
second way to protect a single character, and shows how **Tab**
completion types the awkward part for you.

The units add a few oddly named directories to the `~/projects` tree:

```
~/projects
├── 2026 (draft)
├── Vendor Contracts
├── quarterly-financial-statements-2026
└── ... the directories from the previous modules
```
