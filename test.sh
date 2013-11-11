#!/usr/bin/env bash
set -eu

# Runs all tests and examples. Used for continuous integration.

echo "----> Running tests"
go test -v .
pushd examples > /dev/null
for example in "$(find . -type f -name '*.go')"; do
  echo "----> Running example: ${example}"
  go build "${example}"
  example_bin="${example//.go/}"
  "${example_bin}" || true
  rm "${example_bin}"
done
popd > /dev/null
