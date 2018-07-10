#!/usr/bin/env bash
#
# This builds and starts up the server

if [[ -f /tmp/sci-pid ]]; then
  echo "*** Terminating process"
  kill $(cat /tmp/sci-pid)
  rm /tmp/sci-pid
fi

export INSTALL=1
make -e
if [[ $? == 0 ]]; then
  ./bin/sci server &
  echo $! > /tmp/sci-pid
fi
