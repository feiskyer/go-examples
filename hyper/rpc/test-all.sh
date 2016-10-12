#!/bin/bash
set -e
packages=$(find . -name '*_test.go' -print0 | xargs -0n1 dirname | sed 's|^\./||' | sort -u)
for p in $packages; do
  go test ./...$p -check.vv -v
done
