#!/usr/bin/env bash
set -eu

dest=${1:-production.uoregon.edu:/usr/local/sci}
src=${2:-.}

make clean
make

rsync -av --delete $src/bin/ $dest/bin/
rsync -av --delete $src/db/migrations $dest/db/
rsync -av --delete $src/scripts/ $dest/scripts/
rsync -av --delete $src/static/ $dest/static/
rsync -av --delete $src/templates/ $dest/templates/
