#!/usr/bin/env bash
set -euo pipefail

COVERAGE_FILE="${1:-coverage.out}"

if [[ ! -f "$COVERAGE_FILE" ]]; then
    echo "Coverage file not found: $COVERAGE_FILE"
    echo "Run: go test ./... -coverprofile=coverage.out"
    exit 1
fi

check_threshold() {
    local pattern="$1"
    local threshold="$2"
    local label="$3"

    local pct
    pct=$(go tool cover -func="$COVERAGE_FILE" \
        | { grep "$pattern" || true; } \
        | awk '{print $NF}' \
        | sed 's/%//' \
        | awk '{sum += $1; count++} END {if (count > 0) printf "%.1f", sum/count; else print "0.0"}')

    if [[ "$pct" == "0.0" ]]; then
        echo "WARN: No coverage data found for $label"
        return 0
    fi

    local pass
    pass=$(awk "BEGIN {print ($pct >= $threshold) ? 1 : 0}")

    if [[ "$pass" -eq 1 ]]; then
        echo "PASS: $label coverage ${pct}% >= ${threshold}%"
    else
        echo "FAIL: $label coverage ${pct}% < ${threshold}%"
        return 1
    fi
}

echo "Checking coverage thresholds..."
echo

failures=0

# Overall coverage
total_pct=$(go tool cover -func="$COVERAGE_FILE" | grep "^total:" | awk '{print $NF}' | sed 's/%//')
total_pass=$(awk "BEGIN {print ($total_pct >= 70.0) ? 1 : 0}")
if [[ "$total_pass" -eq 1 ]]; then
    echo "PASS: Overall coverage ${total_pct}% >= 70%"
else
    echo "FAIL: Overall coverage ${total_pct}% < 70%"
    failures=$((failures + 1))
fi

# Package-specific thresholds
check_threshold "internal/data/" 80 "internal/data" || failures=$((failures + 1))
check_threshold "internal/tui/preview.go" 80 "internal/tui/preview.go" || failures=$((failures + 1))

echo
if [[ "$failures" -gt 0 ]]; then
    echo "Coverage check failed with $failures threshold(s) not met"
    exit 1
fi

echo "All coverage thresholds met"
