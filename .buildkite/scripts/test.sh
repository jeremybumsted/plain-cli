#!/bin/bash

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/setup-go.sh"

echo "--- Running :go: Tests"

mise x -- gotestsum --format testname --junitfile unit-tests.xml --junitfile-testcase-classname relative -- -coverprofile=cover.out ./...

echo "--- Uploading artifacts"
buildkite-agent artifact upload "cover.out"
buildkite-agent artifact upload "unit-tests.xml"
