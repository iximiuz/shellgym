---
title: Discover a listening port
vars:
  PORT: { pick: ["7301", "7302", "7303", "7304"] }
  GREETING: { shell: "head -c4 /dev/urandom | od -An -tx1 | tr -d ' \\n'" }
init:
  - name: start_server
    run: |
      pkill -f 'gym-http-serve[r] -m' 2>/dev/null || true
      mkdir -p /tmp/gym-www
      echo "greeting: $GREETING" > /tmp/gym-www/index.html
      export SRV_PORT="$PORT"
      setsid bash -c 'cd /tmp/gym-www && exec -a gym-http-server python3 -m http.server "$SRV_PORT" --bind 127.0.0.1' >/dev/null 2>&1 &
      wait_port --timeout 15 "$PORT"
tasks:
  inspected:
    check: |
      wait_exec '(^|/)(ss|netstat)($| )'
    hint: |
      echo "ss -tlnp lists TCP sockets in listening state with the owning process. Look for the python3 listener on a 73xx port."
    solve: |
      ss -tlnp
  fetched:
    needs: [inspected]
    check: |
      HOME_DIR=$(getent passwd "$GYM_USER" | cut -d: -f6)
      wait_file_contains "$HOME_DIR/greeting.txt" "greeting: $GREETING"
    hint: |
      echo "Fetch http://127.0.0.1:<the-port-you-found>/ with curl and redirect the response into ~/greeting.txt."
    solve: |
      curl -s http://127.0.0.1:$PORT/ > ~/greeting.txt
---

A small web server was just started on this machine, listening on some
port between 7300 and 7310 on `127.0.0.1`. Find out which one.

The `ss` command (socket statistics) shows sockets; `ss -tlnp` narrows the
view to TCP listeners with their port numbers and process names.

::task{name="inspected"}
#active
Waiting for you to inspect the listening sockets...
#completed
Found in the list: a Python web server listening on port `${PORT}`.
::

Now talk to it. `curl URL` performs an HTTP request and prints the
response. Save the server's response into `~/greeting.txt`:

::task{name="fetched"}
#active
Waiting for the server's greeting in `~/greeting.txt`...
#completed
Port discovered, service queried. That is the daily bread of service
debugging.
::
