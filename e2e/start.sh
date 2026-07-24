#!/bin/sh
# Shell Gym playground bootstrap: run as root from the unpacked dist dir.
# The e2e playground's init task (see playground.yaml) downloads
# shellgym-dist.tar.gz into /opt, unpacks it, and runs this script; the
# manifest's health check and the Gym UI tab both expect port 63636.
set -e
DIR=$(cd "$(dirname "$0")" && pwd)
systemd-run --unit=shellgym --collect "$DIR/bin/shellgym" serve \
  --path "$DIR/paths/sample-linux-101" --addr :63636 --user laborant
