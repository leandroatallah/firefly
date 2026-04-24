#!/usr/bin/env bash
set -euo pipefail

file=$(jq -r '.tool_input.file_path // empty')
[[ -z "$file" ]] && exit 0
[[ "$file" != *.go ]] && exit 0

gofmt -w "$file"
cd "$(dirname "$file")"
golangci-lint run
