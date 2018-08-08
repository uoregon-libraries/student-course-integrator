#!/usr/bin/env bash
set -eu

dest=${1:-production.uoregon.edu:/usr/local/sci}
src=${2:-.}

make clean
make
./scripts/rsync_production_files.sh "$src" "$dest"
