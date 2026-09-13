#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

echo "=== Running tests with coverage ==="
go test -coverprofile=coverage.out ./...
go test -race -coverprofile=coverage.out ./...

echo "=== Generating coverage report ==="
go tool cover -func=coverage.out | tee coverage.txt

# Extract total coverage percentage
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | tr -d '%')
echo ""
echo "Total coverage: ${coverage}%"

# Generate badge
color="red"
if (( $(echo "$coverage >= 80" | bc -l) )); then
    color="brightgreen"
elif (( $(echo "$coverage >= 60" | bc -l) )); then
    color="green"
elif (( $(echo "$coverage >= 40" | bc -l) )); then
    color="yellow"
elif (( $(echo "$coverage >= 20" | bc -l) )); then
    color="orange"
fi

# Download badge from shields.io
curl -s "https://img.shields.io/badge/coverage-${coverage}%25-${color}" > coverage-badge.svg

echo ""
echo "Badge saved to coverage-badge.svg"
echo "Coverage report saved to coverage.txt"
