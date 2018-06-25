#!/usr/bin/env bash
# This builds and starts up the server or dies loudly

if [[ -f /tmp/sci-pid ]]; then
  echo "*** Terminating process"
  kill $(cat /tmp/sci-pid)
  rm /tmp/sci-pid
fi

make
if [[ $? == 0 ]]; then
  ./bin/sci &
  echo $! > /tmp/sci-pid
else
  echo ""
  echo
  echo "ERROR"
  echo
fi
