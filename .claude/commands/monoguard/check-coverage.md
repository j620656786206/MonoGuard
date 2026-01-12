---
description: Check test coverage and verify it meets >80% target
argument-hint: [package-path]
---

Check test coverage across MonoGuard packages and verify compliance with >80% target.

**Usage:**

- `/monoguard:check-coverage` - Check all packages
- `/monoguard:check-coverage apps/web` - Check web app coverage
- `/monoguard:check-coverage packages/analysis-engine` - Check Go package coverage

**Coverage Targets (from project-context.md):**

- 🎯 Unit Tests: **>80% coverage**
- 🎯 Integration Tests: Core WASM bridge paths must be tested
- 🎯 E2E Tests: Minimum 3-5 critical user flows

**What This Checks:**

**1. TypeScript/JavaScript Coverage (Vitest):**

- Line coverage: >80%
- Branch coverage: >75%
- Function coverage: >80%
- Uncovered files report

**2. Go Coverage:**

- Unit test coverage per package
- WASM bridge functions coverage
- Critical path coverage (>90%)

**3. E2E Test Coverage:**

- Critical user flows documented
- Minimum 3-5 E2E tests present
- Key user journeys covered

**4. Critical Path Coverage:**

- ✅ WASM bridge error handling (must be >80%)
- ✅ Zustand store state transitions
- ✅ D3.js rendering (no errors)
- ✅ IndexedDB operations

**Analysis Process:**

1. **Run coverage tools:**

   ```bash
   pnpm nx run-many --target=test --all --coverage
   ```

2. **Parse coverage reports:**
   - Read coverage/coverage-summary.json
   - Extract line/branch/function coverage
   - Identify uncovered critical paths

3. **Check against targets:**
   - Unit tests: >80% coverage required
   - Integration tests: Core WASM bridge paths
   - E2E tests: 3-5 critical user flows

4. **Report gaps:**
   - Files below 80% coverage
   - Untested WASM bridge functions
   - Missing integration tests
   - Missing E2E test scenarios

**Report Format:**

```
🧪 MonoGuard Test Coverage Analysis

📊 Overall Coverage:
- Lines: 78% (Target: >80%)
- Branches: 72%
- Functions: 85%
- Statements: 78%

❌ Packages Below Target (<80%):

packages/analysis-engine/
  Coverage: 65% (Target: >80%)
  Missing tests:
  - pkg/analyzer/workspace.go:45-67 (cycle detection)
  - pkg/analyzer/health.go:112-145 (score calculation)

apps/web/app/lib/wasmBridge.ts
  Coverage: 72% (Target: >80%)
  Missing tests:
  - Error handling paths (lines 45-52)
  - Edge cases for empty input (lines 78-85)

✅ Packages meeting coverage target (>80%):
  - packages/types: 95%
  - apps/web/app/stores: 92%

📊 Overall Coverage:
- Statements: 76% (Target: >80%)
- Branches: 68% (Target: >80%)
- Functions: 82% (Target: >80%)
- Lines: 74% (Target: >80%)

🎯 Action Items:
1. Add tests for WASM bridge error cases
2. Increase coverage in analysis-engine/pkg/analyzer
3. Add E2E tests for complete user flows
```

**Commands to Check Coverage:**

```bash
# Run all tests with coverage
pnpm nx run-many --target=test --all --coverage

# Check specific package coverage
pnpm nx test types --coverage

# Generate coverage report
pnpm nx run-many --target=test --all --coverage

# View coverage report (HTML)
open coverage/index.html
```

**Coverage Targets:**

✅ **Unit Tests:** >80% coverage

- Line coverage
- Branch coverage
- Function coverage
- Statement coverage

✅ **Critical Paths (Must be 100%):**

- WASM bridge (all Result<T> error cases)
- Zustand stores (all state transitions)
- IndexedDB operations (save/load/error handling)
- AnalysisError class (layered error handling)

✅ **Integration Tests:**

- WASM bridge → TypeScript integration
- Store + Component integration
- D3.js rendering with real data
- IndexedDB persistence flows

✅ **E2E Tests:**

- Upload → Analyze → Visualize flow
- Error handling flows
- Data persistence verification

**Report Format:**

```
📊 Test Coverage Report

Overall Coverage: 78% (Target: >80%)

❌ Below Target:
packages/analysis-engine/pkg/analyzer/
  Lines:      120/150 (80.0%)
  Statements: 85/100 (85.0%)
  Functions:  12/15 (80.0%)
  Branches:   45/60 (75.0%) ⚠️ Below target

apps/web/app/lib/wasmBridge.ts
  Coverage: 65% ❌ (Target: >80%)
  Missing: Error handling paths (lines 45-52)

🎯 Coverage Summary:
- Overall: 78% (Target: >80%)
- TypeScript: 82% ✅
- Go: 75% ⚠️
- Tests: 45/60 files

❌ Files Below 80% Coverage:
1. packages/analysis-engine/pkg/analyzer/workspace.go: 65%
2. apps/web/app/lib/wasmBridge.ts: 72%
3. apps/web/app/stores/analysis.ts: 75%

💡 Recommendations:
- Add tests for error handling paths
- Test edge cases (empty workspace, circular deps)
- Add integration tests for WASM bridge
```

Let me check the test coverage: **$ARGUMENTS**

I'll analyze test coverage across the project.
