#!/usr/bin/env bash
# This loops forever, restarting the server whenever any kind of change under
# src/ is detected

while true; do
  find src/ | entr -d ./makerun.sh;
done
