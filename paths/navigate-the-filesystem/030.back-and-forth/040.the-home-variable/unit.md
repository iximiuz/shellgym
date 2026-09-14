---
title: The home variable
requires: [readline]
tasks:
  printed:
    check: |
      wait_line '^echo +"?\$\{?HOME\}?"? *$'
    hint: |
      echo "Print the variable with echo and a dollar sign in front of its name: echo \$HOME"
    solve: |
      echo $HOME
  away:
    needs: [printed]
    timeout: 60
    check: |
      P=$(wait_cwd /usr) || exit 1
      set_var SHELL_PID "$P"
    hint: |
      echo "Move to /usr with its full path."
    solve: |
      cd /usr
  home:
    needs: [away]
    timeout: 60
    check: |
      wait_cwd "$SHELL_PID" "$GYM_USER_HOME" >/dev/null || exit 1
      LINE=$(wait_line --latest '^cd( +.*)?$') || exit 1
      ARG=$(printf '%s' "$LINE" | sed -E 's/^cd +//')
      case "$ARG" in
        "\$HOME"|"\$HOME/"|"\${HOME}"|"\${HOME}/"|"\"\$HOME\""|"\"\${HOME}\""|"\"\$HOME\"/"|"\"\${HOME}\"/") exit 0 ;;
        *) hint_exit "You are home, but not through the variable. Go back to /usr and give cd the variable, dollar sign and all, as its argument." ;;
      esac
    hint: |
      echo "Give cd the variable as its argument, written with the dollar sign in front of the name."
    solve: |
      cd $HOME
---

The shell keeps the path of your home directory in a variable named
`HOME`. Writing `$HOME` anywhere on a command line inserts that path,
just like `~` does. The variable form is the one you will see in
scripts and in documentation, because it also works in places where
the tilde does not.

First print the variable with `echo` to see its value. Then move to
`/usr` and return home by giving `cd` the variable as its argument.

::task{name="printed"}
#active
Waiting for you to print the value of `HOME`...
#completed
That path is where a bare `cd` takes you.
::

::task{name="away"}
#active
Waiting for your shell to arrive in `/usr`...
#completed
You are in `/usr`. Now go home through the variable.
::

::task{name="home"}
#active
Waiting for a `cd` through `$HOME` to bring you home...
#completed
The shell replaced the variable with its value before `cd` ran.
::

::tip
The dollar sign belongs to the shell. Inside single quotes it is just
a character, so `'$HOME'` means a dollar sign followed by four
letters, and `cd` complains that no such directory exists.
::
