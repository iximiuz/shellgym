---
title: The hidden names
tasks:
  shown:
    timeout: 60
    check: |
      D="$GYM_USER_HOME"
      ARGV=$(wait_exec --latest '(^|/)ls( +\S+)*$') || exit 1
      set -f; set -- $ARGV; shift
      CWD=$(shell_cwd 2>/dev/null || true)
      N=0
      for a in "$@"; do
        case "$a" in -*) ;; *) N=$((N+1)); T="$a" ;; esac
      done
      if [ "$N" -eq 0 ]; then
        [ "$CWD" = "$D" ] || hint_exit "Your shell is in ${CWD:-another directory}, so that listed the wrong directory. Move to your home directory first, or give ls the tilde as its argument."
      elif [ "$N" -eq 1 ]; then
        case "$T" in /*) P=$(realpath -m "$T") ;; *) P=$(realpath -m "${CWD:-$GYM_USER_HOME}/$T") ;; esac
        [ "$P" = "$D" ] || hint_exit "That listed $T instead of your home directory."
      else
        hint_exit "That listed several paths at once. List only your home directory."
      fi
      S=' -[a-zA-Z0-9]*[aA][a-zA-Z0-9]* '
      L=' --(all|almost-all) '
      if [[ " $ARGV " =~ $S ]] || [[ " $ARGV " =~ $L ]]; then exit 0; fi
      hint_exit "That listing skipped the hidden names. Add the option that shows all entries. You can find it in ls --help."
    hint: |
      CWD=$(shell_cwd 2>/dev/null)
      if [ -n "$CWD" ] && [ "$CWD" != "$GYM_USER_HOME" ]; then
        echo "Your shell is in $CWD. Move to your home directory first, or give ls the tilde as its argument."
        exit 0
      fi
      echo "Look for 'all' in ls --help. The short form is a single letter."
    solve: |
      ls -a ~
---

A name that starts with a dot is **hidden**: `ls` leaves it out unless
you ask. Programs use hidden names for settings they keep in your home
directory, so a plain `ls` there looks almost empty while several files
are present.

List the contents of your home directory with the hidden names
included. You can move there first, or give `ls` the tilde as its
argument from wherever you are. Check `ls --help` for the right option.

::task{name="shown"}
#active
Waiting for a listing that includes the hidden names...
#completed
Those dotted names were there all along.
::

::tip{title="Two names that are always there"}
The listing starts with `.` and `..`, the current directory and its
parent. Every directory has them. The `-A` option shows the hidden
names but leaves those two out.
::
