---
title: Start your own listener
vars:
  MYPORT: { pick: ["8451", "8452", "8453"] }
tasks:
  listening:
    check: |
      wait_port "$MYPORT"
    hint: |
      echo "python3 -m http.server ${MYPORT} starts a web server on that port. Remember & from the processes module if you want your prompt back."
    solve: |
      cd /tmp
      python3 -m http.server $MYPORT &
---

Your turn to serve. Start a web server listening on port `${MYPORT}`.
Python's built-in one is the classic choice for quick file sharing:

```
python3 -m http.server ${MYPORT} &
```

(The trailing `&` keeps it in the background so you can keep typing.)

::task
#active
Waiting for a listener on port `${MYPORT}`...
#completed
Serving. Anything in your working directory is now browsable over HTTP -
handy, and worth remembering when that directory holds private files.
::
