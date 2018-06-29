#!/usr/bin/env bash
#
# This script grabs all necessary dependent packages and installs them, then
# installs this project's packages.  This helps with tools (like vim-go,
# deoplete, etc.) that aren't vgo-aware yet.
cat go.mod | grep "^	" | sed "s|^\t\(.*\) v.*|go get \1|" | bash >/dev/null 2>&1
cat go.mod | grep "^	" | sed "s|^\t\(.*\) v.*|go get \1/...|" | bash >/dev/null 2>&1
go install ./src/...
