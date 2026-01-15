# Story 1.8: Setup Testing Framework and Code Quality

Status: ready-for-dev

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

- [ ] **Task 1: Configure Vitest for Web App** (AC: #1)
  - [ ] 1.1 Install Vitest dependencies:
    ```bash
    pnpm add -D vitest @vitest/coverage-v8 @vitest/ui @testing-library/react @testing-library/jest-dom jsdom --filter @monoguard/web
    ```
  - [ ] 1.2 Create `apps/web/vitest.config.ts`:
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
  - [ ] 1.3 Create `apps/web/src/test/setup.ts`:
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
  - [ ] 1.4 Create sample test `apps/web/src/test/sample.test.tsx`:
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

- [ ] **Task 2: Configure Vitest for Types Package** (AC: #1)
  - [ ] 2.1 Install Vitest dependencies:
    ```bash
    pnpm add -D vitest --filter @monoguard/types
    ```
  - [ ] 2.2 Create `packages/types/vitest.config.ts`:
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
  - [ ] 2.3 Update existing tests to use Vitest syntax if needed

- [ ] **Task 3: Configure Go Testing with Testify** (AC: #2)
  - [ ] 3.1 Add Testify to `packages/analysis-engine`:
    ```bash
    cd packages/analysis-engine
    go get github.com/stretchr/testify
    ```
  - [ ] 3.2 Create sample test `packages/analysis-engine/internal/result/result_test.go`:
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
  - [ ] 3.3 Add Testify to `apps/cli`:
    ```bash
    cd apps/cli
    go get github.com/stretchr/testify
    ```
  - [ ] 3.4 Update Makefile to include coverage:
    ```makefile
    # In both analysis-engine and cli Makefiles
    test:
    	go test -v -race ./...

    test-coverage:
    	go test -v -race -coverprofile=coverage.out ./...
    	go tool cover -html=coverage.out -o coverage.html
    ```

- [ ] **Task 4: Configure Biome** (AC: #4)
  - [ ] 4.1 Install Biome at root:
    ```bash
    pnpm add -D @biomejs/biome
    ```
  - [ ] 4.2 Create `biome.json` at repository root:
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
  - [ ] 4.3 Add Biome scripts to root `package.json`:
    ```json
    {
      "scripts": {
        "lint": "biome check .",
        "lint:fix": "biome check --write .",
        "format": "biome format --write ."
      }
    }
    ```
  - [ ] 4.4 Remove ESLint/Prettier if present (optional - can coexist)

- [ ] **Task 5: Configure Pre-commit Hooks** (AC: #5)
  - [ ] 5.1 Install Husky and lint-staged:
    ```bash
    pnpm add -D husky lint-staged
    pnpm exec husky init
    ```
  - [ ] 5.2 Create `.husky/pre-commit`:
    ```bash
    #!/usr/bin/env sh
    . "$(dirname -- "$0")/_/husky.sh"

    pnpm lint-staged
    ```
  - [ ] 5.3 Add lint-staged config to `package.json`:
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

- [ ] **Task 6: Update Nx Project Configurations** (AC: #6)
  - [ ] 6.1 Update `apps/web/project.json`:
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
  - [ ] 6.2 Update `packages/types/project.json` with test targets
  - [ ] 6.3 Verify Go projects have test targets in project.json

- [ ] **Task 7: Create Coverage Configuration** (AC: #3)
  - [ ] 7.1 Add coverage reporting to CI workflow:
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
  - [ ] 7.2 Create `.codecov.yml` (optional):
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

- [ ] **Task 8: Verification** (AC: #1, #2, #4, #6)
  - [ ] 8.1 Run `pnpm nx test web` - verify Vitest runs
  - [ ] 8.2 Run `pnpm nx test types` - verify types tests run
  - [ ] 8.3 Run `pnpm nx test analysis-engine` - verify Go tests run
  - [ ] 8.4 Run `pnpm nx test cli` - verify CLI Go tests run
  - [ ] 8.5 Run `pnpm biome check` - verify formatting check
  - [ ] 8.6 Run `pnpm biome check --write` - verify auto-fix
  - [ ] 8.7 Make a commit - verify pre-commit hooks run
  - [ ] 8.8 Check coverage reports are generated

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

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

