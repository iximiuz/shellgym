---
title: Stop the listener
needs: [serve-yourself]
vars:
  MYPORT: { from: serve-yourself.MYPORT }
tasks:
  port_free:
    check: |
      # Baseline: the port must be busy before we watch for it to free up.
      wait_port --timeout 15 "$MYPORT" || hint_exit "Nothing is listening on port ${MYPORT} - the server from the previous rep is already gone, so there is nothing to stop. Revisit the previous rep to start it again."
      wait_port_free "$MYPORT"
    hint: |
      echo "Find the server's PID (pgrep -f http.server, or the jobs command since you started it from this shell) and send it a signal with kill."
    solve: |
      pkill -f "http.server $MYPORT"
---

The server you started on port `${MYPORT}` is still running. Stop it, so
the port becomes free again.

You know the tools already: find the process, send it a signal. Since you
started it from your own shell, `jobs` will list it too, and `kill %1`
addresses a job by number instead of PID.

::task{name="port_free"}
#active
Waiting for port `${MYPORT}` to become free...
#completed
Port released. Start, inspect, stop - you have now walked a service
through its whole lifecycle.
::
