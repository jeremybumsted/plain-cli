#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y
mise install

echo "--- Downloading :go: modules"
mise x -- go mod download

echo "--- :rocket: Creating GitHub Release"

# Validate BUILDKITE_TAG is set
if [ -z "${BUILDKITE_TAG:-}" ]; then
  echo "ERROR: BUILDKITE_TAG environment variable must be set for releases"
  exit 1
fi

# Validate GITHUB_TOKEN is set
if [ -z "${GITHUB_TOKEN:-}" ]; then
  echo "ERROR: GITHUB_TOKEN environment variable must be set for releases"
  exit 1
fi

echo "Creating release for tag: $BUILDKITE_TAG"

# Run goreleaser in release mode (builds all platforms and creates GitHub release)
echo "Running goreleaser release --clean"
mise x -- goreleaser release --clean

echo "--- :white_check_mark: Release completed successfully"
echo "Release available at: https://github.com/jeremybumsted/plain-cli/releases/tag/$BUILDKITE_TAG"
