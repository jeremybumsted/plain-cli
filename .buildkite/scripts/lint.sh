#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y
mise install

echo "--- Downloading :go: modules"
mise x -- go mod download

echo "--- Running :golangci-lint:"
mise x golangci-lint -- golangci-lint run --verbose --timeout 2m
