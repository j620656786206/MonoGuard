#!/bin/bash
# ============================================
# MonoGuard Burn-In Test Runner
# ============================================
#
# Runs E2E tests multiple times to detect flaky tests.
# Usage: ./scripts/burn-in.sh [iterations] [project]
#
# Examples:
#   ./scripts/burn-in.sh          # 5 iterations, chromium
#   ./scripts/burn-in.sh 10       # 10 iterations, chromium
#   ./scripts/burn-in.sh 5 firefox # 5 iterations, firefox
#
# Based on TEA Knowledge Base: ci-burn-in.md

set -e

# Configuration
ITERATIONS=${1:-5}
PROJECT=${2:-chromium}
RESULTS_DIR="burn-in-results"

echo ""
echo "🔥 MonoGuard Burn-In Test Runner"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Iterations: $ITERATIONS"
echo "Project:    $PROJECT"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Create results directory
mkdir -p "$RESULTS_DIR"

# Track failures
FAILURES=()
PASSED=0

# Run burn-in loop
for i in $(seq 1 $ITERATIONS); do
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "🔄 Iteration $i/$ITERATIONS"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  # Run tests
  if pnpm nx run web-e2e:e2e --project="$PROJECT" 2>&1 | tee "$RESULTS_DIR/iteration-$i.log"; then
    echo "✅ Iteration $i passed"
    PASSED=$((PASSED + 1))
  else
    echo "❌ Iteration $i FAILED"
    FAILURES+=($i)

    # Save failure artifacts
    mkdir -p "$RESULTS_DIR/failure-$i"
    cp -r apps/web-e2e/test-results/* "$RESULTS_DIR/failure-$i/" 2>/dev/null || true
    cp -r apps/web-e2e/playwright-report/* "$RESULTS_DIR/failure-$i/" 2>/dev/null || true

    echo ""
    echo "🛑 BURN-IN FAILED on iteration $i"
    echo "Failure artifacts saved to: $RESULTS_DIR/failure-$i/"
    echo ""

    # Fail fast on first failure
    exit 1
  fi
done

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎉 BURN-IN COMPLETE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Passed: $PASSED/$ITERATIONS"
echo ""

if [ ${#FAILURES[@]} -eq 0 ]; then
  echo "✅ All tests are stable and ready to merge!"

  # Cleanup logs on success
  rm -rf "$RESULTS_DIR"
  exit 0
else
  echo "❌ Flaky tests detected on iterations: ${FAILURES[*]}"
  echo "Review failure artifacts in: $RESULTS_DIR/"
  exit 1
fi
