# ATDD Checklist - Epic 1, Story 1-4: Setup Go CLI Project with Cobra

**Date:** 2026-01-16
**Author:** Murat (TEA Agent) via Alexyu
**Primary Test Level:** Integration (Go Cobra) + E2E (Binary execution)

---

## Story Summary

**As a** developer,
**I want** a Go CLI project using Cobra for command management,
**So that** I can build command-line tools for dependency analysis with native performance and cross-platform distribution.

---

## Acceptance Criteria

1. **AC1**: Go module initialized with `go.mod`, Go 1.21+, Cobra and Viper dependencies
2. **AC2**: Project directory structure with cmd/, pkg/config/, pkg/output/, main.go, Makefile
3. **AC3**: Cobra command structure with monoguard, analyze/check/fix/init subcommands, global flags
4. **AC4**: Viper configuration integration (config file, home directory, environment variables, CLI flag override)
5. **AC5**: Build configuration with Makefile, cross-compilation, version via ldflags, < 15MB binary
6. **AC6**: Placeholder commands work with text and JSON output formats
7. **AC7**: Nx integration with `pnpm nx build cli` and `pnpm nx test cli`

---

## Failing Tests Created (RED Phase)

### Unit Tests (2 packages)

**File:** `apps/cli/pkg/config/config_test.go` (110 lines)

| Test                              | Status | Verifies                                   |
| --------------------------------- | ------ | ------------------------------------------ |
| `TestLoadConfig`                  | 🔴 RED | AC4 - Config loading from multiple sources |
| `TestConfigStructure`             | 🔴 RED | AC4 - Config struct fields                 |
| `TestEnvironmentVariableOverride` | 🔴 RED | AC4 - MONOGUARD\_ env prefix               |
| `TestCLIFlagOverride`             | 🔴 RED | AC4 - CLI flags override config            |

**File:** `apps/cli/pkg/output/output_test.go` (115 lines)

| Test                      | Status | Verifies                  |
| ------------------------- | ------ | ------------------------- |
| `TestFormatterText`       | 🔴 RED | AC6 - Text output format  |
| `TestFormatterJSON`       | 🔴 RED | AC6 - JSON output format  |
| `TestFormatterPrettyJSON` | 🔴 RED | AC6 - Pretty-printed JSON |
| `TestNewFormatter`        | 🔴 RED | Formatter creation        |
| `TestFormatterPrint`      | 🔴 RED | Print to stdout           |

### Integration Tests (5 files)

**File:** `apps/cli/cmd/root_test.go` (90 lines)

| Test                         | Status | Verifies                            |
| ---------------------------- | ------ | ----------------------------------- |
| `TestRootCommandHelp`        | 🔴 RED | AC3 - Help output structure         |
| `TestRootCommandVersion`     | 🔴 RED | AC3 - Version information           |
| `TestGlobalFlagsRegistered`  | 🔴 RED | AC3 - --config, --verbose, --format |
| `TestRootCommandDescription` | 🔴 RED | AC3 - Short/Long descriptions       |
| `TestSubcommandsRegistered`  | 🔴 RED | AC3 - analyze, check, fix, init     |
| `TestViperInitialization`    | 🔴 RED | AC4 - Viper setup                   |

**File:** `apps/cli/cmd/analyze_test.go` (85 lines)

| Test                           | Status | Verifies                 |
| ------------------------------ | ------ | ------------------------ |
| `TestAnalyzeCommandRegistered` | 🔴 RED | AC3 - Command registered |
| `TestAnalyzeCommandUsage`      | 🔴 RED | AC6 - Usage string       |
| `TestAnalyzeCommandTextOutput` | 🔴 RED | AC6 - Placeholder text   |
| `TestAnalyzeCommandJSONOutput` | 🔴 RED | AC6 - Valid JSON         |
| `TestAnalyzeCommandWithPath`   | 🔴 RED | AC6 - Path argument      |

**File:** `apps/cli/cmd/check_test.go` (95 lines)

| Test                         | Status | Verifies                     |
| ---------------------------- | ------ | ---------------------------- |
| `TestCheckCommandRegistered` | 🔴 RED | AC3 - Command registered     |
| `TestCheckCommandExitCode`   | 🔴 RED | AC6 - Exit code 0            |
| `TestCheckCommandTextOutput` | 🔴 RED | AC6 - Placeholder text       |
| `TestCheckCommandJSONOutput` | 🔴 RED | AC6 - Valid JSON with passed |
| `TestCheckCommandFlags`      | 🔴 RED | AC6 - --fail-on, --threshold |
| `TestCheckCommandCIMode`     | 🔴 RED | AC6 - CI/CD integration      |

**File:** `apps/cli/cmd/fix_test.go` (75 lines)

| Test                          | Status | Verifies                     |
| ----------------------------- | ------ | ---------------------------- |
| `TestFixCommandRegistered`    | 🔴 RED | AC3 - Command registered     |
| `TestFixCommandTextOutput`    | 🔴 RED | AC6 - Placeholder text       |
| `TestFixCommandJSONOutput`    | 🔴 RED | AC6 - Valid JSON with dryRun |
| `TestFixCommandDryRunFlag`    | 🔴 RED | AC6 - --dry-run flag         |
| `TestFixCommandWithoutDryRun` | 🔴 RED | AC6 - dryRun=false default   |

**File:** `apps/cli/cmd/init_cmd_test.go` (55 lines)

| Test                         | Status | Verifies                 |
| ---------------------------- | ------ | ------------------------ |
| `TestInitCommandRegistered`  | 🔴 RED | AC3 - Command registered |
| `TestInitCommandTextOutput`  | 🔴 RED | AC6 - Placeholder text   |
| `TestInitCommandJSONOutput`  | 🔴 RED | AC6 - Valid JSON         |
| `TestInitCommandDescription` | 🔴 RED | AC3 - Description        |

### E2E Tests (1 file)

**File:** `apps/cli/tests/cli_e2e_test.go` (180 lines)

| Test                 | Status | Verifies                   |
| -------------------- | ------ | -------------------------- |
| `TestBinaryExists`   | 🔴 RED | AC5 - Binary built, < 15MB |
| `TestHelpOutput`     | 🔴 RED | AC3 - --help output        |
| `TestVersionOutput`  | 🔴 RED | AC3 - --version output     |
| `TestAnalyzeCommand` | 🔴 RED | AC6 - analyze text/JSON    |
| `TestCheckCommand`   | 🔴 RED | AC6 - check exit code 0    |
| `TestFixCommand`     | 🔴 RED | AC6 - fix --dry-run        |
| `TestInitCommand`    | 🔴 RED | AC6 - init output          |
| `TestFormatFlag`     | 🔴 RED | AC6 - --format json\|text  |

---

## Data Factories Created

Not applicable for this story (CLI project - no test data factories needed).

---

## Fixtures Created

Not applicable for this story (Go CLI uses table-driven tests, not fixtures).

---

## Mock Requirements

None - placeholder commands return static responses.

---

## Required data-testid Attributes

Not applicable for this story (CLI project - no UI).

---

## Implementation Checklist

### Test: Config Package Unit Tests

**File:** `apps/cli/pkg/config/config_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/pkg/config/config.go` with Config struct
- [ ] Implement `Load()` function using Viper
- [ ] Support `.monoguard.yaml` in current directory
- [ ] Support `~/.monoguard/config.yaml` fallback
- [ ] Support `MONOGUARD_` environment variable prefix
- [ ] Run test: `cd apps/cli && go test ./pkg/config/...`
- [ ] ✅ Tests pass (green phase)

### Test: Output Package Unit Tests

**File:** `apps/cli/pkg/output/output_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/pkg/output/output.go` with Formatter struct
- [ ] Implement `NewFormatter(format string)` constructor
- [ ] Implement `Print(data interface{})` for stdout
- [ ] Implement `PrintTo(w io.Writer, data interface{})` for testing
- [ ] Support "text" format (human-readable)
- [ ] Support "json" format (pretty-printed)
- [ ] Run test: `cd apps/cli && go test ./pkg/output/...`
- [ ] ✅ Tests pass (green phase)

### Test: Root Command Integration Tests

**File:** `apps/cli/cmd/root_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/cmd/root.go` with rootCmd
- [ ] Set Use to "monoguard"
- [ ] Set Short and Long descriptions
- [ ] Set Version to "0.1.0" (via ldflags)
- [ ] Register --config, --verbose, --format global flags
- [ ] Initialize Viper in init() via cobra.OnInitialize
- [ ] Run test: `cd apps/cli && go test ./cmd/... -run TestRoot`
- [ ] ✅ Tests pass (green phase)

### Test: Analyze Command Integration Tests

**File:** `apps/cli/cmd/analyze_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/cmd/analyze.go` with analyzeCmd
- [ ] Set Use to "analyze [path]"
- [ ] Accept optional path argument (default ".")
- [ ] Output placeholder text format
- [ ] Output placeholder JSON format with status, path, message
- [ ] Register command in init()
- [ ] Run test: `cd apps/cli && go test ./cmd/... -run TestAnalyze`
- [ ] ✅ Tests pass (green phase)

### Test: Check Command Integration Tests

**File:** `apps/cli/cmd/check_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/cmd/check.go` with checkCmd
- [ ] Set Use to "check [path]"
- [ ] Return exit code 0 for placeholder success
- [ ] Output placeholder text with passed indicator
- [ ] Output placeholder JSON with passed=true
- [ ] Register --fail-on and --threshold flags
- [ ] Register command in init()
- [ ] Run test: `cd apps/cli && go test ./cmd/... -run TestCheck`
- [ ] ✅ Tests pass (green phase)

### Test: Fix Command Integration Tests

**File:** `apps/cli/cmd/fix_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/cmd/fix.go` with fixCmd
- [ ] Set Use to "fix [path]"
- [ ] Register --dry-run flag (default false)
- [ ] Output placeholder text with dry run indicator
- [ ] Output placeholder JSON with dryRun field
- [ ] Register command in init()
- [ ] Run test: `cd apps/cli && go test ./cmd/... -run TestFix`
- [ ] ✅ Tests pass (green phase)

### Test: Init Command Integration Tests

**File:** `apps/cli/cmd/init_cmd_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/cmd/init_cmd.go` with initCmd
- [ ] Set Use to "init"
- [ ] Set Short and Long descriptions mentioning configuration
- [ ] Output placeholder text
- [ ] Output placeholder JSON with status, message
- [ ] Register command in init()
- [ ] Run test: `cd apps/cli && go test ./cmd/... -run TestInit`
- [ ] ✅ Tests pass (green phase)

### Test: E2E Binary Tests

**File:** `apps/cli/tests/cli_e2e_test.go`

**Tasks to make tests pass:**

- [ ] Create `apps/cli/main.go` entry point
- [ ] Create `apps/cli/Makefile` with build target
- [ ] Run `make build` to create dist/monoguard binary
- [ ] Verify binary < 15MB
- [ ] Verify --help output
- [ ] Verify --version output
- [ ] Verify all commands work with --format text and --format json
- [ ] Run test: `cd apps/cli && go test ./tests/...`
- [ ] ✅ Tests pass (green phase)

### Test: Nx Integration

**No automated test - manual verification**

**Tasks:**

- [ ] Update `apps/cli/project.json` for Nx
- [ ] Update `apps/cli/package.json` for npm scripts
- [ ] Run `pnpm nx build cli` - verify success
- [ ] Run `pnpm nx test cli` - verify Go tests run
- [ ] Run `pnpm nx graph` - verify project appears
- [ ] ✅ Nx integration verified

---

## Running Tests

```bash
# Run all Go tests for CLI
cd apps/cli && go test -v ./...

# Run specific package tests
cd apps/cli && go test -v ./pkg/config/...
cd apps/cli && go test -v ./pkg/output/...
cd apps/cli && go test -v ./cmd/...

# Run E2E tests (requires binary built first)
cd apps/cli && make build && go test -v ./tests/...

# Run tests with coverage
cd apps/cli && go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run via Nx (after implementation)
pnpm nx test cli
```

---

## Red-Green-Refactor Workflow

### RED Phase (Complete) ✅

**TEA Agent Responsibilities:**

- ✅ All tests written and failing
- ✅ Test files created:
  - `pkg/config/config_test.go`
  - `pkg/output/output_test.go`
  - `cmd/root_test.go`
  - `cmd/analyze_test.go`
  - `cmd/check_test.go`
  - `cmd/fix_test.go`
  - `cmd/init_cmd_test.go`
  - `tests/cli_e2e_test.go`

---

### GREEN Phase (DEV Team - Next Steps)

**DEV Agent Responsibilities:**

1. **Pick one failing test** from implementation checklist
2. **Read the test** to understand expected behavior
3. **Implement minimal code** to make that specific test pass
4. **Run the test** to verify it now passes (green)
5. **Check off the task** in implementation checklist
6. **Move to next test** and repeat

**Recommended Order:**

1. Start with `pkg/config/config.go` and `pkg/output/output.go` (unit tests)
2. Then `cmd/root.go` (integration tests depend on it)
3. Then `cmd/analyze.go`, `cmd/check.go`, `cmd/fix.go`, `cmd/init_cmd.go`
4. Finally `main.go` and `Makefile` for E2E tests

---

### REFACTOR Phase (DEV Team - After All Tests Pass)

1. **Verify all tests pass** (green phase complete)
2. **Review code for quality** (Go best practices)
3. **Extract duplications** (shared helpers)
4. **Ensure tests still pass** after each refactor

---

## Next Steps

1. **Share this checklist** with the dev workflow
2. **Run failing tests** to confirm RED phase: `cd apps/cli && go test ./...`
3. **Begin implementation** using implementation checklist
4. **Work one test at a time** (red → green for each)
5. **When all tests pass**, update story status to 'done'

---

## Test Statistics

| Category          | Count  |
| ----------------- | ------ |
| Unit Tests        | 9      |
| Integration Tests | 22     |
| E2E Tests         | 8      |
| **Total Tests**   | **39** |
| Test Files        | 8      |
| Total Lines       | ~805   |

---

## Knowledge Base References Applied

- **test-levels-framework.md** - Unit vs Integration vs E2E selection
- **test-quality.md** - Table-driven tests, deterministic assertions
- **Go testing patterns** - `t.Run()`, table-driven tests from `result_test.go`

---

**Generated by BMad TEA Agent (Murat)** - 2026-01-16
