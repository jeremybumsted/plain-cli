#!/bin/bash

set -euo pipefail

echo "--- Checking dependencies up to date"
mise trust -y
mise install

echo "--- :golang: Building binary with goreleaser"

# Set GOARCH if not already set (default to amd64)
export GOARCH="${GOARCH:-amd64}"

# Validate GOOS is set
if [ -z "${GOOS:-}" ]; then
  echo "ERROR: GOOS environment variable must be set"
  exit 1
fi

echo "Building for GOOS=$GOOS GOARCH=$GOARCH"

# Determine build mode based on tag presence
if [ -n "${BUILDKITE_TAG:-}" ]; then
  echo "Building for tag release: $BUILDKITE_TAG"
  BUILD_MODE="release"
else
  echo "CI build - snapshot mode"
  BUILD_MODE="snapshot"
fi

# Run goreleaser with single-target flag
if [ "$BUILD_MODE" = "release" ]; then
  echo "Running goreleaser release (single-target)"
  mise x -- goreleaser release --clean --single-target
else
  echo "Running goreleaser build (snapshot, single-target)"
  mise x -- goreleaser build --snapshot --clean --single-target
fi

echo "--- :white_check_mark: Build completed successfully"
