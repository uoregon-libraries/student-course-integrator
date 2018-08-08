#!/usr/bin/env bash
set -eu

src=${1:-}
dest=${2:-}

if [[ ! -d $src ]]; then
  echo "Invalid source for rsync: \"$src\""
  exit 1
fi

if [[ $dest == "" ]]; then
  echo "rsync destination must be set"
  exit 1
fi

rsync -a --delete $src/bin/ $dest/bin/
rsync -a --delete $src/db/migrations $dest/db/
rsync -a --delete $src/scripts/ $dest/scripts/
rsync -a --delete $src/static/ $dest/static/
rsync -a --delete $src/templates/ $dest/templates/
