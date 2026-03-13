# filename: scripts/test_coverage.sh
#!/usr/bin/env bash
set -euo pipefail

# Ensure we're at repo root (adjust if needed)
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

COVERAGE_DIR="coverage"
PROFILE="$COVERAGE_DIR/coverage.out"
HTML="$COVERAGE_DIR/coverage.html"

mkdir -p "$COVERAGE_DIR"

# Run tests with coverage across all modules/packages
# -covermode=atomic is safer for parallel tests
go test ./... -covermode=atomic -coverprofile="$PROFILE"

# Generate human-readable HTML report
go tool cover -html="$PROFILE" -o "$HTML"

echo "Coverage profile: $PROFILE"
echo "HTML report: $HTML"
