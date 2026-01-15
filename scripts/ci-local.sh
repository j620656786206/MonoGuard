#!/bin/bash
# ============================================
# MonoGuard Local CI Mirror
# ============================================
#
# Mirrors the CI pipeline locally for debugging.
# Runs the same stages as GitHub Actions.
#
# Usage: ./scripts/ci-local.sh
#
# Based on TEA Knowledge Base: ci-burn-in.md

set -e

echo ""
echo "🔍 MonoGuard Local CI Pipeline"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "This script mirrors the CI pipeline locally."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Track timing
START_TIME=$(date +%s)

# Stage 1: Lint
echo "📝 Stage 1: Lint & Type Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
pnpm nx run-many --target=lint --all --parallel=3 || {
  echo "❌ Lint failed"
  exit 1
}
echo "✅ Lint passed"
echo ""

# Stage 2: Unit Tests
echo "🧪 Stage 2: Unit Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
pnpm nx run-many --target=test --all --parallel=3 || {
  echo "❌ Unit tests failed"
  exit 1
}
echo "✅ Unit tests passed"
echo ""

# Stage 3: E2E Tests
echo "🎭 Stage 3: E2E Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
pnpm nx run web-e2e:e2e --project=chromium || {
  echo "❌ E2E tests failed"
  exit 1
}
echo "✅ E2E tests passed"
echo ""

# Stage 4: Mini Burn-In (3 iterations instead of 5)
echo "🔥 Stage 4: Mini Burn-In (3 iterations)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
for i in {1..3}; do
  echo "Burn-in iteration $i/3..."
  pnpm nx run web-e2e:e2e --project=chromium || {
    echo "❌ Burn-in failed on iteration $i"
    exit 1
  }
  echo "✅ Iteration $i passed"
done
echo ""

# Summary
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ LOCAL CI PIPELINE PASSED"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Duration: ${DURATION}s"
echo ""
echo "Your changes are ready for PR!"
