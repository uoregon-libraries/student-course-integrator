#!/usr/bin/env bash
set -eu

command -v goose >/dev/null 2>&1 || {
  echo >&2 "goose needs to be installed in order to run migrations:"
  echo >&2
  echo >&2 "    go get -u bitbucket.org/liamstask/goose/..."
  echo >&2
  echo >&2 'You must also have $GOPATH/bin in your $PATH or else copy the goose binary'
  echo >&2 'from $GOPATH/bin to a location in your $PATH'
  exit 1
}

if [[ -f /etc/sci.conf ]]; then
  source /etc/sci.conf
elif [[ -f ./sci.conf ]]; then
  source ./sci.conf
else
  echo >&2 "You must have /etc/sci.conf or ./sci.conf for migrate.sh to work"
  exit 1
fi

echo "development:" >db/dbconf.yml
echo "  driver: mysql" >>db/dbconf.yml
echo "  open: $DB" >>db/dbconf.yml

goose up
rm db/dbconf.yml
