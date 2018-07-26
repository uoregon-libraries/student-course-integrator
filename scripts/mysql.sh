#!/usr/bin/env bash
set -eu

source sci.conf
dbuser=${DB%%:*}

temp=${DB#*:}
dbpass=${temp%%@tcp*}

temp=${temp##*@tcp(}
dbhost=${temp%%:3306*}

dbname=${temp##*/}

mysql -u$dbuser -p$dbpass -h$dbhost -D$dbname
