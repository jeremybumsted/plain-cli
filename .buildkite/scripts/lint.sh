#!/bin/bash

set -uo pipefail

mise x golangci-lint -- golangci-lint run --verbose --timeout 2m
