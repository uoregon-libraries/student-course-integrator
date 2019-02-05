#!/usr/bin/env bash
set -eu

command -v goose >/dev/null 2>&1 || {
  echo >&2 "goose needs to be installed in order to run migrations:"
  echo >&2
  echo >&2 "    go get github.com/pressly/goose/cmd/goose"
  echo >&2
  echo >&2 'You must also have $GOPATH/bin in your $PATH or else copy the goose binary'
  echo >&2 'from $GOPATH/bin to a location in your $PATH'
  exit 1
}

DB="${SCI_DB:-}"
if [[ "$DB" != "" ]]; then
  echo "Using db config from environment"
elif [[ -f /etc/sci.conf ]]; then
  echo "Reading db config from /etc/sci.conf"
  source /etc/sci.conf
elif [[ -f ./sci.conf ]]; then
  echo "Reading db config from ./sci.conf"
  source ./sci.conf
else
  echo >&2 "You must have /etc/sci.conf or ./sci.conf for migrate.sh to work"
  exit 1
fi

echo 'Executing goose'
goose -dir ./db/migrations mysql "$DB" up
echo "Database migration completed"
