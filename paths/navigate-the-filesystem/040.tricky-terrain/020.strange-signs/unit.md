---
title: Strange signs
vars:
  ODD: { pick: ["data[2024]", "cache+static", "logs (old)"] }
init:
  - name: create_store
    run: |
      rm -rf /tmp/gym-store
      mkdir -p "/tmp/gym-store/$ODD"
      chmod -R a+rx /tmp/gym-store
tasks:
  arrive:
    check: |
      wait_cwd "/tmp/gym-store/$ODD"
    hint: |
      CWD=$(shell_cwd 2>/dev/null || echo "?")
      echo "You are in $CWD. Characters like [, ], (, ) and + can mean something to the shell before cd ever runs. Single quotes around the whole name make every character literal."
    solve: |
      cd "/tmp/gym-store/$ODD"
---

Under `/tmp/gym-store` sits a directory named `${ODD}`. Brackets,
parentheses, plus signs - characters the shell may try to interpret
before `cd` gets a say.

The armor-plated answer is single quotes: **everything** between `'`
and `'` is taken literally, no exceptions.

Get inside:

::task{name="arrive"}
#active
Waiting for your shell in `/tmp/gym-store/${ODD}`...
#completed
Made it. When a name looks even slightly suspicious, single-quote it
and move on - cheaper than remembering which characters are special.
::

::tip
---
title: Tab quotes for you
---
Type `cd /tmp/gym-store/` and press `Tab`: completion fills in the odd
name with all the escaping already in place.
::
