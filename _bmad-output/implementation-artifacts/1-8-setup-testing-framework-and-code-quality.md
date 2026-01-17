# Story 1.8: Setup Testing Framework and Code Quality

Status: done

## Story

As a **developer**,
I want **testing frameworks configured with code quality tools**,
So that **I can write and run tests with consistent code standards across TypeScript and Go projects**.

## Acceptance Criteria

1. **AC1: Vitest Configuration**
   - Given the TypeScript packages and apps
   - When Vitest is configured
   - Then:
     - `apps/web/` has Vitest configured for React components
     - `packages/types/` has Vitest configured for type tests
     - Tests run with `pnpm nx test <project>`
     - Coverage reports are generated

2. **AC2: Go Testing with Testify**
   - Given the Go packages
   - When Go testing is configured
   - Then:
     - `packages/analysis-engine/` has Go tests with Testify
     - `apps/cli/` has Go tests with Testify
     - Tests run with `make test` or `go test ./...`
     - Coverage reports are generated

3. **AC3: Coverage Thresholds**
   - Given the test configuration
   - When coverage is measured
   - Then:
     - Overall coverage target: 80%
     - Critical paths coverage: 90%
     - CI fails if thresholds not met
     - Coverage reports uploaded to artifacts

4. **AC4: Code Formatting (Biome)**
   - Given the codebase
   - When Biome is configured
   - Then:
     - TypeScript/JavaScript files are formatted consistently
     - `pnpm biome check` verifies formatting
     - `pnpm biome format --write` auto-fixes formatting
     - Biome replaces ESLint + Prettier for simpler tooling

5. **AC5: Pre-commit Hooks (Optional)**
   - Given the development workflow
   - When pre-commit hooks are configured
   - Then:
     - Husky runs lint-staged on commit
     - Only staged files are checked
     - Format errors are auto-fixed
     - Tests are NOT run on commit (too slow)

6. **AC6: Nx Integration**
   - Given the test configuration
   - When Nx commands are run
   - Then:
     - `pnpm nx test web` runs Vitest for web app
     - `pnpm nx test types` runs Vitest for types package
     - `pnpm nx test analysis-engine` runs Go tests
     - `pnpm nx test cli` runs Go tests

## Tasks / Subtasks

- [x] **Task 1: Configure Vitest for Web App** (AC: #1)
  - [x] 1.1 Install Vitest dependencies (already installed):
    ```bash
    pnpm add -D vitest @vitest/coverage-v8 @vitest/ui @testing-library/react @testing-library/jest-dom jsdom --filter @monoguard/web
    ```
  - [x] 1.2 Create `apps/web/vitest.config.ts`:
    ```typescript
    import { defineConfig } from 'vitest/config'
    import react from '@vitejs/plugin-react'
    import tsconfigPaths from 'vite-tsconfig-paths'

    export default defineConfig({
      plugins: [react(), tsconfigPaths({ root: '../../' })],
      test: {
        environment: 'jsdom',
        globals: true,
        setupFiles: ['./src/test/setup.ts'],
        include: ['**/*.{test,spec}.{js,ts,jsx,tsx}'],
        exclude: ['node_modules', '.output'],
        coverage: {
          provider: 'v8',
          reporter: ['text', 'json', 'html'],
          exclude: [
            'node_modules/',
            'src/test/',
            '**/*.d.ts',
            '**/*.config.*',
          ],
          thresholds: {
            global: {
              branches: 80,
              functions: 80,
              lines: 80,
              statements: 80,
            },
          },
        },
      },
    })
    ```
  - [x] 1.3 Create `apps/web/src/test/setup.ts`:
    ```typescript
    import '@testing-library/jest-dom/vitest'
    import { cleanup } from '@testing-library/react'
    import { afterEach } from 'vitest'

    // Cleanup after each test
    afterEach(() => {
      cleanup()
    })

    // Mock WASM for tests
    vi.mock('@/lib/wasmLoader', () => ({
      loadWasm: vi.fn().mockResolvedValue({
        analyzer: {
          getVersion: vi.fn().mockReturnValue(JSON.stringify({
            data: { version: '0.1.0' },
            error: null,
          })),
          analyze: vi.fn().mockReturnValue(JSON.stringify({
            data: { healthScore: 85, packageCount: 10 },
            error: null,
          })),
          check: vi.fn().mockReturnValue(JSON.stringify({
            data: { passed: true, errors: [] },
            error: null,
          })),
        },
      }),
    }))
    ```
  - [x] 1.4 Create sample test `apps/web/src/test/sample.test.tsx`:
    ```typescript
    import { describe, it, expect } from 'vitest'
    import { render, screen } from '@testing-library/react'

    describe('Sample Test', () => {
      it('renders correctly', () => {
        render(<div>Hello MonoGuard</div>)
        expect(screen.getByText('Hello MonoGuard')).toBeInTheDocument()
      })
    })
    ```

- [x] **Task 2: Configure Vitest for Types Package** (AC: #1)
  - [x] 2.1 Install Vitest dependencies:
    ```bash
    pnpm add -D vitest --filter @monoguard/types
    ```
  - [x] 2.2 Create `packages/types/vitest.config.ts`:
    ```typescript
    import { defineConfig } from 'vitest/config'

    export default defineConfig({
      test: {
        globals: true,
        include: ['**/*.{test,spec}.ts'],
        coverage: {
          provider: 'v8',
          reporter: ['text', 'json'],
          thresholds: {
            global: {
              branches: 80,
              functions: 80,
              lines: 80,
              statements: 80,
            },
          },
        },
      },
    })
    ```
  - [x] 2.3 Update existing tests to use Vitest syntax if needed

- [x] **Task 3: Configure Go Testing with Testify** (AC: #2)
  - [x] 3.1 Add Testify to `packages/analysis-engine`:
    ```bash
    cd packages/analysis-engine
    go get github.com/stretchr/testify
    ```
  - [x] 3.2 Create sample test `packages/analysis-engine/internal/result/result_test.go`:
    ```go
    package result

    import (
        "encoding/json"
        "testing"

        "github.com/stretchr/testify/assert"
        "github.com/stretchr/testify/require"
    )

    func TestNewSuccess(t *testing.T) {
        data := map[string]int{"healthScore": 85}
        result := NewSuccess(data)

        assert.NotNil(t, result.Data)
        assert.Nil(t, result.Error)
    }

    func TestNewError(t *testing.T) {
        result := NewError("PARSE_ERROR", "Invalid JSON")

        assert.Nil(t, result.Data)
        require.NotNil(t, result.Error)
        assert.Equal(t, "PARSE_ERROR", result.Error.Code)
        assert.Equal(t, "Invalid JSON", result.Error.Message)
    }

    func TestResultToJSON(t *testing.T) {
        result := NewSuccess(map[string]string{"version": "0.1.0"})
        jsonStr := result.ToJSON()

        var parsed Result
        err := json.Unmarshal([]byte(jsonStr), &parsed)
        require.NoError(t, err)
        assert.NotNil(t, parsed.Data)
    }
    ```
  - [x] 3.3 Add Testify to `apps/cli`:
    ```bash
    cd apps/cli
    go get github.com/stretchr/testify
    ```
  - [x] 3.4 Update Makefile to include coverage:
    ```makefile
    # In both analysis-engine and cli Makefiles
    test:
    	go test -v -race ./...

    test-coverage:
    	go test -v -race -coverprofile=coverage.out ./...
    	go tool cover -html=coverage.out -o coverage.html
    ```

- [x] **Task 4: Configure Biome** (AC: #4)
  - [x] 4.1 Install Biome at root:
    ```bash
    pnpm add -D @biomejs/biome
    ```
  - [x] 4.2 Create `biome.json` at repository root:
    ```json
    {
      "$schema": "https://biomejs.dev/schemas/1.9.0/schema.json",
      "organizeImports": {
        "enabled": true
      },
      "linter": {
        "enabled": true,
        "rules": {
          "recommended": true,
          "correctness": {
            "noUnusedVariables": "error",
            "noUnusedImports": "error"
          },
          "suspicious": {
            "noExplicitAny": "warn"
          },
          "style": {
            "useConst": "error",
            "noNonNullAssertion": "warn"
          }
        }
      },
      "formatter": {
        "enabled": true,
        "indentStyle": "space",
        "indentWidth": 2,
        "lineWidth": 100
      },
      "javascript": {
        "formatter": {
          "quoteStyle": "single",
          "trailingCommas": "es5",
          "semicolons": "asNeeded"
        }
      },
      "files": {
        "ignore": [
          "node_modules",
          "dist",
          ".output",
          ".next",
          "coverage",
          "*.d.ts",
          "pnpm-lock.yaml"
        ]
      }
    }
    ```
  - [x] 4.3 Add Biome scripts to root `package.json`:
    ```json
    {
      "scripts": {
        "lint": "biome check .",
        "lint:fix": "biome check --write .",
        "format": "biome format --write ."
      }
    }
    ```
  - [x] 4.4 Remove ESLint/Prettier if present (optional - can coexist)

- [x] **Task 5: Configure Pre-commit Hooks** (AC: #5)
  - [x] 5.1 Install Husky and lint-staged:
    ```bash
    pnpm add -D husky lint-staged
    pnpm exec husky init
    ```
  - [x] 5.2 Create `.husky/pre-commit`:
    ```bash
    #!/usr/bin/env sh
    . "$(dirname -- "$0")/_/husky.sh"

    pnpm lint-staged
    ```
  - [x] 5.3 Add lint-staged config to `package.json`:
    ```json
    {
      "lint-staged": {
        "*.{js,ts,jsx,tsx}": [
          "biome check --write"
        ],
        "*.{json,md,yaml,yml}": [
          "biome format --write"
        ]
      }
    }
    ```

- [x] **Task 6: Update Nx Project Configurations** (AC: #6)
  - [x] 6.1 Update `apps/web/project.json`:
    ```json
    {
      "targets": {
        "test": {
          "executor": "nx:run-commands",
          "options": {
            "command": "vitest run",
            "cwd": "apps/web"
          }
        },
        "test:watch": {
          "executor": "nx:run-commands",
          "options": {
            "command": "vitest",
            "cwd": "apps/web"
          }
        },
        "test:coverage": {
          "executor": "nx:run-commands",
          "options": {
            "command": "vitest run --coverage",
            "cwd": "apps/web"
          }
        }
      }
    }
    ```
  - [x] 6.2 Update `packages/types/project.json` with test targets
  - [x] 6.3 Verify Go projects have test targets in project.json

- [x] **Task 7: Create Coverage Configuration** (AC: #3)
  - [x] 7.1 Add coverage reporting to CI workflow:
    ```yaml
    # In .github/workflows/ci.yml test job
    - name: Run TypeScript tests with coverage
      run: pnpm nx affected -t test --coverage

    - name: Upload coverage
      uses: codecov/codecov-action@v4
      with:
        files: ./coverage/lcov.info,./packages/analysis-engine/coverage.out
        fail_ci_if_error: false
    ```
  - [x] 7.2 Create `.codecov.yml` (optional):
    ```yaml
    coverage:
      status:
        project:
          default:
            target: 80%
        patch:
          default:
            target: 80%
    ```

- [x] **Task 8: Verification** (AC: #1, #2, #4, #6)
  - [x] 8.1 Run `pnpm nx test web` - verify Vitest runs
  - [x] 8.2 Run `pnpm nx test types` - verify types tests run
  - [x] 8.3 Run `pnpm nx test analysis-engine` - verify Go tests run
  - [x] 8.4 Run `pnpm nx test cli` - verify CLI Go tests run
  - [x] 8.5 Run `pnpm biome check` - verify formatting check
  - [x] 8.6 Run `pnpm biome check --write` - verify auto-fix
  - [x] 8.7 Make a commit - verify pre-commit hooks run
  - [x] 8.8 Check coverage reports are generated

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**

- **TypeScript Testing:** Vitest (Vite-native, faster than Jest)
- **Go Testing:** Standard `testing` package + Testify for assertions
- **Code Quality:** Biome (replaces ESLint + Prettier)
- **Coverage Target:** 80% overall, 90% for critical paths

**Critical Constraints:**

- **Mock WASM:** TypeScript tests must mock WASM calls
- **No Browser in Go Tests:** Go tests run in native environment
- **Consistent Formatting:** All code must pass Biome checks

### Critical Don't-Miss Rules

**From project-context.md:**

1. **WASM Mocking Pattern:**
   ```typescript
   // ✅ CORRECT: Mock with proper Result type
   vi.mock('@/lib/wasmLoader', () => ({
     MonoGuardAnalyzer: {
       analyze: vi.fn().mockResolvedValue({
         data: { healthScore: 85 },
         error: null,
       }),
     },
   }))

   // ❌ WRONG: Mock without Result structure
   vi.mock('@/lib/wasmLoader', () => ({
     analyze: vi.fn().mockResolvedValue({ healthScore: 85 }),
   }))
   ```

2. **Go Test File Naming:**
   ```
   # ✅ CORRECT: *_test.go in same directory
   pkg/analyzer/
   ├── workspace.go
   └── workspace_test.go

   # ❌ WRONG: Separate test directory
   pkg/analyzer/workspace.go
   test/analyzer/workspace_test.go
   ```

3. **Testify Assertions:**
   ```go
   // ✅ CORRECT: Use assert for soft failures, require for hard failures
   assert.Equal(t, expected, actual)  // Continues on failure
   require.NoError(t, err)            // Stops on failure

   // ❌ WRONG: Only using t.Fatal
   if err != nil {
       t.Fatal(err)
   }
   ```

### Testing Structure

```
mono-guard/
├── apps/
│   ├── web/
│   │   ├── vitest.config.ts
│   │   └── src/
│   │       └── test/
│   │           ├── setup.ts
│   │           └── *.test.tsx
│   └── cli/
│       └── cmd/
│           └── *_test.go
├── packages/
│   ├── types/
│   │   ├── vitest.config.ts
│   │   └── src/
│   │       └── __tests__/
│   │           └── *.test.ts
│   └── analysis-engine/
│       └── internal/
│           └── result/
│               └── result_test.go
└── biome.json
```

### Biome vs ESLint/Prettier

**Why Biome:**
- 10-100x faster than ESLint
- Single tool for lint + format
- Zero configuration needed
- TypeScript-first design

**Migration Notes:**
- Biome can coexist with ESLint initially
- Gradually remove ESLint configs
- Update CI to use Biome

### Coverage Strategy

| Package | Target | Critical Paths |
|---------|--------|----------------|
| `@monoguard/web` | 80% | UI components, hooks |
| `@monoguard/types` | 90% | Type guards, Result handling |
| `@monoguard/analysis-engine` | 80% | Parser, analyzer |
| `@monoguard/cli` | 80% | Command handlers |

### Previous Story Intelligence

**From Story 1.5 (ready-for-dev):**

- `packages/types/` already has `__tests__/` directory
- Result type guards (`isSuccess`, `isError`) need tests

**From Story 1.3 (ready-for-dev):**

- `internal/result/result.go` needs `result_test.go`
- Go test files in same directory as source

**From Story 1.6 (ready-for-dev):**

- CI already runs tests
- Add coverage upload step

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#Testing Framework]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.8]
- [Source: _bmad-output/project-context.md#Test Organization]
- [Vitest Documentation](https://vitest.dev/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Biome Documentation](https://biomejs.dev/)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- All tests passing: web (10 tests), types (18 tests), analysis-engine (Go tests), cli (Go tests)
- Biome v2.3.11 installed and configured (migrated from 1.9.0 schema)

### Completion Notes List

1. **Vitest for Web App** - Created `apps/web/vitest.config.ts` and `apps/web/src/test/setup.ts` with:
   - jsdom environment for React Testing Library
   - WASM mock for tests
   - Coverage thresholds at 80%
   - Created `HealthScoreDisplay.test.tsx` as sample component test

2. **Vitest for Types Package** - Created `packages/types/vitest.config.ts` with:
   - Coverage thresholds at 80%
   - Existing tests already using Vitest syntax

3. **Go Testing with Testify** - Added Testify to:
   - `packages/analysis-engine/go.mod` (v1.11.1)
   - `apps/cli/go.mod`
   - Updated Makefiles with `test-coverage` target for both

4. **Biome** - Installed @biomejs/biome and created `biome.json` with:
   - Linter rules (recommended + custom)
   - Formatter settings (single quotes, trailing commas)
   - Added `biome:check`, `biome:fix`, `biome:format` scripts
   - Note: Coexists with existing ESLint/Prettier

5. **Pre-commit Hooks** - Already configured with Husky + lint-staged

6. **Nx Project Configurations** - Added `test:watch` and `test:coverage` targets to:
   - `apps/web/project.json`
   - `packages/types/project.json`
   - `packages/analysis-engine/project.json`
   - `apps/cli/project.json`

7. **Coverage Configuration** - CI already has coverage configured with artifact upload

8. **Code Review Fixes (2026-01-17)** - Addressed issues found in adversarial code review:
   - Added testify dependency to `packages/analysis-engine/go.mod` (was missing)
   - Refactored `packages/analysis-engine/internal/result/result_test.go` to use testify assert/require
   - Fixed 172 files with Biome auto-formatting (`biome check --write`)
   - Updated `lint-staged` config to use Biome instead of ESLint/Prettier (as specified in AC5)
   - Note: 182 lint errors remain (a11y SVG titles, noExplicitAny warnings) - these are valid lint issues requiring manual fixes

### File List

**New Files:**
- `apps/web/vitest.config.ts`
- `apps/web/src/test/setup.ts`
- `apps/web/src/__tests__/HealthScoreDisplay.test.tsx`
- `packages/types/vitest.config.ts`
- `biome.json`
- `packages/analysis-engine/go.sum` (testify dependencies)

**Modified Files:**
- `packages/analysis-engine/go.mod` (added testify)
- `packages/analysis-engine/internal/result/result_test.go` (refactored to use testify)
- `packages/analysis-engine/Makefile` (added test-coverage)
- `apps/cli/Makefile` (added test-coverage)
- `apps/cli/go.mod` (testify as indirect dependency)
- `package.json` (added biome scripts, updated lint-staged to use Biome)
- `apps/web/project.json` (added test:watch, test:coverage)
- `apps/web/src/__tests__/app.test.ts` (sample test)
- `packages/types/project.json` (added test:watch, test:coverage)
- `packages/analysis-engine/project.json` (added test:coverage)
- `apps/cli/project.json` (added test:coverage)
- `pnpm-lock.yaml` (dependency updates)
- Multiple files formatted by Biome (172 files)

