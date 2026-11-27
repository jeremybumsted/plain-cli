#!/bin/bash

echo "--- Checking dependencies up to date"
mise trust -y
mise install

echo "--- Downloading :go: modules"
if [[ "${BUILDKITE:-}" == true ]]; then
  export GOMODCACHE="$PWD/.buildkite/cache/go-mod"
  buildkite-agent cache restore --name go_modules
fi
mise x -- go mod download
if [[ "${BUILDKITE:-}" == true ]]; then
  buildkite-agent cache save --name go_modules
fi
