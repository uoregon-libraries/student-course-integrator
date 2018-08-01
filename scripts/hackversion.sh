#!/usr/bin/env bash
set -eu

output="./src/version/commit.go"

echo "package version" > $output
echo "" >> $output

commit=$(git rev-parse --short=8 --verify HEAD 2>&1)
echo 'const commit = "'$commit'"' >> $output
