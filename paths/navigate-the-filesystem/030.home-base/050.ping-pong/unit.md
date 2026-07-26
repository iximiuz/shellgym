---
title: Config, data, config, data
vars:
  APP: { pick: [ledger, courier, sentry] }
init:
  - name: create_app_dirs
    run: |
      rm -rf "/etc/gym-$APP" "/var/lib/gym-$APP"
      mkdir -p "/etc/gym-$APP" "/var/lib/gym-$APP"
      chmod -R a+rx "/etc/gym-$APP" "/var/lib/gym-$APP"
tasks:
  visit_config:
    check: |
      wait_cwd "/etc/gym-$APP"
    hint: |
      echo "First stop: the app's config directory /etc/gym-${APP}."
    solve: |
      cd /etc/gym-$APP
  visit_data:
    needs: [visit_config]
    check: |
      wait_cwd "/var/lib/gym-$APP"
    hint: |
      echo "Now its data directory: /var/lib/gym-${APP}."
    solve: |
      cd /var/lib/gym-$APP
  back_to_config:
    needs: [visit_data]
    check: |
      wait_cwd "/etc/gym-$APP"
    hint: |
      echo "Back to the config directory - this hop and every following one is just 'cd -'."
    solve: |
      cd -
  data_again:
    needs: [back_to_config]
    check: |
      wait_cwd "/var/lib/gym-$APP"
    hint: |
      echo "And over to data once more. 'cd -' again - each use swaps you to the other side."
    solve: |
      cd -
---

The `${APP}` service keeps its configuration in `/etc/gym-${APP}` and
its data in `/var/lib/gym-${APP}` - a classic Linux split. Debugging
means bouncing between the two, and `cd -` turns the bounce into a
reflex.

Visit the config directory, then the data directory:

::task{name="visit_config"}
#active
Waiting for your shell in `/etc/gym-${APP}`...
#completed
Config side.
::

::task{name="visit_data"}
#active
Waiting for your shell in `/var/lib/gym-${APP}`...
#completed
Data side. Both endpoints established - now let `cd -` do the walking.
::

Bounce back to config, then to data again - two more hops, no paths
typed:

::task{name="back_to_config"}
#active
Waiting for your shell back in `/etc/gym-${APP}`...
#completed
One `cd -`.
::

::task{name="data_again"}
#active
Waiting for your shell back in `/var/lib/gym-${APP}`...
#completed
Two `cd -`. Each swap flips you to the other side - config, data,
config, data, for as long as the debugging takes.
::
