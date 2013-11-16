#!/usr/bin/env bash
set -eu

# Runs all tests and examples. Used for continuous integration.

echo "----> Running tests"
go test -v .
pushd examples > /dev/null
for example in $(find . -type f -name '*.go'); do
  echo "----> Running example: ${example}"
  bin="example-$(basename $(dirname "${example}"))"
  go build -o "${bin}" "${example}"
  "./${bin}" || true
  rm -rf "${bin}"
done
popd > /dev/null
