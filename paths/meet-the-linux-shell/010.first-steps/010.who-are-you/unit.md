---
title: Ask who you are
tasks:
  asked:
    check: |
      wait_exec '(^|/)whoami$'
    hint: |
      echo "Type exactly: whoami - one word, no spaces - and press Enter."
    solve: |
      whoami
---

Click into the terminal so it receives your keystrokes. You are logged in
under some user name - the machine knows which one, and there is a
command that asks it: `whoami` (read it as "who am I").

Type it after the prompt, all in one word, and press **Enter**.

::task
#active
Waiting for you to ask the machine who you are...
#completed
That name is your user account. Notice the order of events: prompt,
your command, the answer, and a fresh prompt ready for the next one.
::

::tip
Made a typo? **Backspace** erases, and the **Left** and **Right** arrow
keys move the cursor within the line before you press Enter.
::
