#!/bin/bash
# ============================================
# MonoGuard Selective Test Runner
# ============================================
#
# Runs only tests affected by changed files.
# Useful for faster feedback during development.
#
# Usage: ./scripts/test-changed.sh [base-branch]
#
# Examples:
#   ./scripts/test-changed.sh         # Compare to main
#   ./scripts/test-changed.sh develop # Compare to develop
#
# Based on TEA Knowledge Base: selective-testing.md

set -e

# Configuration
BASE_BRANCH=${1:-main}

echo ""
echo "🎯 MonoGuard Selective Test Runner"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Base branch: $BASE_BRANCH"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Detect changed files
echo "📋 Detecting changed files..."
CHANGED_FILES=$(git diff --name-only origin/$BASE_BRANCH...HEAD 2>/dev/null || git diff --name-only HEAD~1)

if [ -z "$CHANGED_FILES" ]; then
  echo "✅ No files changed. Skipping tests."
  exit 0
fi

echo "Changed files:"
echo "$CHANGED_FILES" | sed 's/^/  - /'
echo ""

# Detect changed test files
CHANGED_SPECS=$(echo "$CHANGED_FILES" | grep -E '\.(spec|test)\.(ts|js)$' || echo "")

# Critical file changes = run all tests
if echo "$CHANGED_FILES" | grep -qE '(package\.json|pnpm-lock\.yaml|playwright\.config|nx\.json)'; then
  echo "⚠️  Critical configuration files changed."
  echo "Running FULL test suite..."
  echo ""
  pnpm nx run web-e2e:e2e
  exit $?
fi

# Test file changes = run those specific tests
if [ -n "$CHANGED_SPECS" ]; then
  echo "🧪 Test files changed:"
  echo "$CHANGED_SPECS" | sed 's/^/  - /'
  echo ""
  echo "Running changed test files..."

  # Convert to space-separated list for playwright
  SPEC_LIST=$(echo "$CHANGED_SPECS" | tr '\n' ' ')
  pnpm exec playwright test $SPEC_LIST
  exit $?
fi

# Source file changes = run affected tests using Nx
if echo "$CHANGED_FILES" | grep -qE '\.(ts|tsx|js|jsx)$'; then
  echo "📦 Source files changed."
  echo "Running affected tests via Nx..."
  echo ""
  pnpm nx affected --target=e2e --base=origin/$BASE_BRANCH
  exit $?
fi

# Documentation/config only = skip tests
if echo "$CHANGED_FILES" | grep -qE '\.(md|json|yml|yaml)$'; then
  echo "📝 Only documentation/config files changed."
  echo "Skipping E2E tests (run manually if needed)."
  exit 0
fi

# Default: run all tests
echo "⚙️  Running full test suite..."
pnpm nx run web-e2e:e2e
