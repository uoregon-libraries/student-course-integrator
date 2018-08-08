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

rsync -av --delete $src/bin/ $dest/bin/
rsync -av --delete $src/db/migrations $dest/db/
rsync -av --delete $src/scripts/ $dest/scripts/
rsync -av --delete $src/static/ $dest/static/
rsync -av --delete $src/templates/ $dest/templates/
