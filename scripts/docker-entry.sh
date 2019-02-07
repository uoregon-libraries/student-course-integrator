#!/usr/bin/env bash
set -eu

echo "Waiting for database connectivity - \"Can't connect ...\" messages are normal for a few seconds"
wait_for_database

if [[ ! -f /setup ]]; then
  echo "Performing first-time setup"
  /app/db/migrate.sh
  mysql -usci -psci -Dsci -hdb -e "source /app/db/seed.sql"
fi

echo 'Executing "'$@'"'
exec $@
