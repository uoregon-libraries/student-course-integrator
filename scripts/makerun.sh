#!/usr/bin/env bash
#
# This builds and starts up the server

make -e
if [[ $? == 0 ]]; then
  docker-compose stop web
  docker-compose rm -f web
  docker-compose up -d web
fi
