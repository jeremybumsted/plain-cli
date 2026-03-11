#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y

mise install

echo "--- Downloading :go: modules"
mise x -- go mod download

echo "--- Running :go: Tests"

mise x -- gotestsum --format testname --junitfile unit-tests.xml --junitfile-testcase-classname relative -- -coverprofile=cover.out ./...

echo "--- Uploading artifacts"
buildkite-agent artifact upload "cover.out"
buildkite-agent artifact upload "unit-tests.xml"
