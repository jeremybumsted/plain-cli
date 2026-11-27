#!/bin/bash

set -euo pipefail

source "$(dirname "${BASH_SOURCE[0]}")/setup-go.sh"

echo "--- Running :golangci-lint:"
mise x golangci-lint -- golangci-lint run --verbose --timeout 2m
