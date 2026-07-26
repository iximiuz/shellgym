---
title: Scaffold a project
vars:
  PROJ: { pick: [atlas, nimbus, quasar] }
init:
  - name: create_workdir
    run: |
      rm -rf "$GYM_USER_HOME/gymwork"
      mkdir -p "$GYM_USER_HOME/gymwork"
      chown "$GYM_USER" "$GYM_USER_HOME/gymwork"
tasks:
  skeleton:
    check: |
      wait_dir "$GYM_USER_HOME/gymwork/$PROJ/src"
      wait_dir "$GYM_USER_HOME/gymwork/$PROJ/docs"
    hint: |
      echo "Two directories under ~/gymwork/${PROJ}: src and docs. mkdir -p handles missing parents, and mkdir takes several paths in one command."
    solve: |
      mkdir -p ~/gymwork/$PROJ/src ~/gymwork/$PROJ/docs
  seeded:
    needs: [skeleton]
    check: |
      wait_file "$GYM_USER_HOME/gymwork/$PROJ/src/main.txt"
      wait_file "$GYM_USER_HOME/gymwork/$PROJ/docs/notes.txt"
    hint: |
      echo "Now the starter files: src/main.txt and docs/notes.txt inside ~/gymwork/${PROJ}. touch, two paths, one command."
    solve: |
      touch ~/gymwork/$PROJ/src/main.txt ~/gymwork/$PROJ/docs/notes.txt
---

New project `${PROJ}` kicks off today. Set up its standard skeleton
under `~/gymwork`:

```
${PROJ}/
  src/
  docs/
```

::task{name="skeleton"}
#active
Waiting for `src` and `docs` under `~/gymwork/${PROJ}`...
#completed
Skeleton raised - ideally with one `mkdir -p` and two arguments.
::

Now drop in the starter files: `main.txt` in `src`, `notes.txt` in
`docs`:

::task{name="seeded"}
#active
Waiting for `src/main.txt` and `docs/notes.txt`...
#completed
Scaffolded and seeded. You just performed the opening ritual of ten
thousand real projects - `ls -R ${PROJ}` shows the whole tree at once
if you want to admire it.
::
