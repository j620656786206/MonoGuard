# Story 1.4: Setup Go CLI Project with Cobra

Status: done

## Story

As a **developer**,
I want **a Go CLI project using Cobra for command management**,
So that **I can build command-line tools for dependency analysis with native performance and cross-platform distribution**.

## Acceptance Criteria

1. **AC1: Go Module Initialization**
   - Given the Nx monorepo from Story 1.1
   - When I initialize the Go CLI project in apps/cli
   - Then I have:
     - `go.mod` with module path `github.com/j620656786206/MonoGuard/apps/cli`
     - Go version 1.21+ specified in go.mod
     - Cobra and Viper as dependencies

2. **AC2: Project Directory Structure**
   - Given the Go module is initialized
   - When I verify the project structure
   - Then I have:
     ```
     apps/cli/
     ├── cmd/
     │   ├── root.go       # Root command with Viper config
     │   ├── analyze.go    # Placeholder analyze command
     │   ├── check.go      # Placeholder check command
     │   ├── fix.go        # Placeholder fix command
     │   └── init.go       # Placeholder init command
     ├── pkg/
     │   ├── config/       # Viper configuration management
     │   └── output/       # Output formatting (JSON, text)
     ├── main.go           # Entry point
     ├── go.mod
     ├── go.sum
     └── Makefile
     ```

3. **AC3: Cobra Command Structure**
   - Given the CLI project structure
   - When I run `./monoguard --help`
   - Then I see:
     - Command name: `monoguard`
     - Available commands: `analyze`, `check`, `fix`, `init`
     - Global flags: `--config`, `--verbose`, `--format`
     - Version information via `--version`

4. **AC4: Viper Configuration Integration**
   - Given the Cobra commands
   - When I run commands
   - Then Viper is configured to:
     - Read from `.monoguard.yaml` in current directory
     - Read from `~/.monoguard/config.yaml` (global config)
     - Support environment variables with `MONOGUARD_` prefix
     - Allow CLI flags to override config file values

5. **AC5: Build Configuration**
   - Given the Go project
   - When I run `make build`
   - Then:
     - Build produces executable binary `dist/monoguard`
     - Binary can be cross-compiled for macOS, Linux, Windows
     - Build includes version information via ldflags
     - Binary size is reasonable (< 15MB)

6. **AC6: Placeholder Commands Work**
   - Given the built CLI binary
   - When I run each command
   - Then:
     - `./monoguard analyze` outputs placeholder message
     - `./monoguard check` exits with code 0 (placeholder success)
     - `./monoguard fix --dry-run` outputs placeholder message
     - `./monoguard init` outputs placeholder message
     - All commands accept `--format json|text` flag

7. **AC7: Nx Integration**
   - Given the Go CLI project
   - When I run Nx commands
   - Then:
     - `pnpm nx build cli` runs the Makefile
     - `pnpm nx test cli` runs Go tests
     - Project appears correctly in `pnpm nx graph`

## Tasks / Subtasks

- [x] **Task 1: Clean Up TypeScript CLI** (AC: #1)
  - [x] 1.1 Backup existing TypeScript CLI structure for reference
  - [x] 1.2 Remove TypeScript source files:
    ```bash
    rm -rf apps/cli/src
    rm apps/cli/tsconfig*.json
    rm apps/cli/eslint.config.cjs
    rm -rf apps/cli/node_modules
    ```
  - [x] 1.3 Keep package.json for npm distribution (will be updated)

- [x] **Task 2: Initialize Go Module** (AC: #1)
  - [x] 2.1 Initialize Go module:
    ```bash
    cd apps/cli
    go mod init github.com/j620656786206/MonoGuard/apps/cli
    ```
  - [x] 2.2 Add Cobra and Viper dependencies:
    ```bash
    go get github.com/spf13/cobra@latest
    go get github.com/spf13/viper@latest
    ```
  - [x] 2.3 Run `go mod tidy` to clean up dependencies

- [x] **Task 3: Create Project Structure** (AC: #2)
  - [x] 3.1 Create directory structure:
    ```bash
    mkdir -p cmd pkg/{config,output}
    ```
  - [x] 3.2 Create .gitignore for dist/ directory
  - [x] 3.3 Add placeholder README.md (skipped - not required)

- [x] **Task 4: Implement Root Command** (AC: #3, #4)
  - [x] 4.1 Create `main.go`:

    ```go
    package main

    import "github.com/j620656786206/MonoGuard/apps/cli/cmd"

    func main() {
        cmd.Execute()
    }
    ```

  - [x] 4.2 Create `cmd/root.go`:
        ```go
        package cmd

        import (
            "fmt"
            "os"

            "github.com/spf13/cobra"
            "github.com/spf13/viper"
        )

        var (
            cfgFile string
            verbose bool
            format  string
            version = "0.1.0" // Set via ldflags during build
        )

        var rootCmd = &cobra.Command{
            Use:     "monoguard",
            Short:   "MonoGuard - Monorepo dependency analysis and validation",
            Long:    `MonoGuard is a comprehensive tool for analyzing monorepo

    dependencies, detecting circular dependencies, and providing
    actionable fix suggestions.`,
    Version: version,
    }

        func Execute() {
            if err := rootCmd.Execute(); err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
            }
        }

        func init() {
            cobra.OnInitialize(initConfig)

            rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
                "config file (default is .monoguard.yaml)")
            rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
                "verbose output")
            rootCmd.PersistentFlags().StringVar(&format, "format", "text",
                "output format (text|json)")

            viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
            viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
        }

        func initConfig() {
            if cfgFile != "" {
                viper.SetConfigFile(cfgFile)
            } else {
                // Search for config in current directory
                viper.AddConfigPath(".")
                // Search in home directory
                home, _ := os.UserHomeDir()
                viper.AddConfigPath(home + "/.monoguard")
                viper.SetConfigType("yaml")
                viper.SetConfigName(".monoguard")
            }

            // Read environment variables with MONOGUARD_ prefix
            viper.SetEnvPrefix("MONOGUARD")
            viper.AutomaticEnv()

            // Read config file (ignore error if not found)
            viper.ReadInConfig()
        }
        ```

- [x] **Task 5: Implement Placeholder Commands** (AC: #6)
  - [x] 5.1 Create `cmd/analyze.go`:
        ```go
        package cmd

        import (
            "fmt"

            "github.com/spf13/cobra"
            "github.com/spf13/viper"
        )

        var analyzeCmd = &cobra.Command{
            Use:   "analyze [path]",
            Short: "Analyze monorepo dependencies",
            Long:  `Analyze the dependency structure of a monorepo and generate

    a comprehensive report including circular dependencies, health score,
    and fix suggestions.`,
    Args: cobra.MaximumNArgs(1),
    Run: func(cmd \*cobra.Command, args []string) {
    path := "."
    if len(args) > 0 {
    path = args[0]
    }

                format := viper.GetString("format")
                if format == "json" {
                    fmt.Printf(`{"status":"placeholder","path":"%s","message":"Analysis will be implemented in Epic 2"}`, path)
                    fmt.Println()
                } else {
                    fmt.Printf("🔍 MonoGuard Analysis (Placeholder)\n")
                    fmt.Printf("   Path: %s\n", path)
                    fmt.Printf("   Status: Will be implemented in Epic 2\n")
                }
            },
        }

        func init() {
            rootCmd.AddCommand(analyzeCmd)
        }
        ```

  - [x] 5.2 Create `cmd/check.go`:
        ```go
        package cmd

        import (
            "fmt"
            "os"

            "github.com/spf13/cobra"
            "github.com/spf13/viper"
        )

        var (
            failOn    string
            threshold int
        )

        var checkCmd = &cobra.Command{
            Use:   "check [path]",
            Short: "Validate dependencies for CI/CD",
            Long:  `Run validation checks on the monorepo dependencies.

    Returns exit code 0 on success, 1 on failure.
    Designed for CI/CD integration.`,
    Args: cobra.MaximumNArgs(1),
    Run: func(cmd \*cobra.Command, args []string) {
    path := "."
    if len(args) > 0 {
    path = args[0]
    }

                format := viper.GetString("format")
                if format == "json" {
                    fmt.Printf(`{"status":"placeholder","path":"%s","passed":true,"message":"Check will be implemented in Epic 2"}`, path)
                    fmt.Println()
                } else {
                    fmt.Printf("✅ MonoGuard Check (Placeholder)\n")
                    fmt.Printf("   Path: %s\n", path)
                    fmt.Printf("   Status: Passed (placeholder)\n")
                }
                os.Exit(0) // Placeholder success
            },
        }

        func init() {
            rootCmd.AddCommand(checkCmd)
            checkCmd.Flags().StringVar(&failOn, "fail-on", "all",
                "fail on: circular|boundary|all")
            checkCmd.Flags().IntVar(&threshold, "threshold", 0,
                "fail if health score below threshold (0-100)")
        }
        ```

  - [x] 5.3 Create `cmd/fix.go`:
        ```go
        package cmd

        import (
            "fmt"

            "github.com/spf13/cobra"
            "github.com/spf13/viper"
        )

        var dryRun bool

        var fixCmd = &cobra.Command{
            Use:   "fix [path]",
            Short: "Generate fix suggestions for issues",
            Long:  `Analyze the monorepo and generate fix suggestions

    for circular dependencies and other issues.`,
    Args: cobra.MaximumNArgs(1),
    Run: func(cmd \*cobra.Command, args []string) {
    path := "."
    if len(args) > 0 {
    path = args[0]
    }

                format := viper.GetString("format")
                if format == "json" {
                    fmt.Printf(`{"status":"placeholder","path":"%s","dryRun":%t,"message":"Fix will be implemented in Epic 3"}`, path, dryRun)
                    fmt.Println()
                } else {
                    fmt.Printf("🔧 MonoGuard Fix (Placeholder)\n")
                    fmt.Printf("   Path: %s\n", path)
                    fmt.Printf("   Dry Run: %t\n", dryRun)
                    fmt.Printf("   Status: Will be implemented in Epic 3\n")
                }
            },
        }

        func init() {
            rootCmd.AddCommand(fixCmd)
            fixCmd.Flags().BoolVar(&dryRun, "dry-run", false,
                "preview fixes without applying")
        }
        ```

  - [x] 5.4 Create `cmd/init_cmd.go` (avoiding conflict with init function):
        ```go
        package cmd

        import (
            "fmt"

            "github.com/spf13/cobra"
            "github.com/spf13/viper"
        )

        var initCmd = &cobra.Command{
            Use:   "init",
            Short: "Initialize MonoGuard configuration",
            Long:  `Create a .monoguard.yaml configuration file in the

    current directory with default settings.`,
        Run: func(cmd *cobra.Command, args []string) {
            format := viper.GetString("format")
            if format == "json" {
                fmt.Println(`{"status":"placeholder","message":"Init will be implemented in Epic 8"}`)
    } else {
    fmt.Printf("🚀 MonoGuard Init (Placeholder)\n")
    fmt.Printf(" Status: Will be implemented in Epic 8\n")
    }
    },
    }

        func init() {
            rootCmd.AddCommand(initCmd)
        }
        ```

- [x] **Task 6: Create pkg Utilities** (AC: #4, #6)
  - [x] 6.1 Create `pkg/config/config.go`:

    ```go
    // Package config provides configuration management using Viper
    package config

    import "github.com/spf13/viper"

    // Config represents the MonoGuard configuration structure
    type Config struct {
        Workspaces []string `mapstructure:"workspaces"`
        Rules      Rules    `mapstructure:"rules"`
        Thresholds Thresholds `mapstructure:"thresholds"`
    }

    type Rules struct {
        CircularDependencies string `mapstructure:"circularDependencies"`
        BoundaryViolations   string `mapstructure:"boundaryViolations"`
    }

    type Thresholds struct {
        HealthScore int `mapstructure:"healthScore"`
    }

    // Load reads configuration from Viper
    func Load() (*Config, error) {
        var cfg Config
        if err := viper.Unmarshal(&cfg); err != nil {
            return nil, err
        }
        return &cfg, nil
    }
    ```

  - [x] 6.2 Create `pkg/output/output.go`:

    ```go
    // Package output provides formatted output utilities
    package output

    import (
        "encoding/json"
        "fmt"
    )

    // Formatter handles output formatting
    type Formatter struct {
        Format string // "text" or "json"
    }

    // NewFormatter creates a new output formatter
    func NewFormatter(format string) *Formatter {
        return &Formatter{Format: format}
    }

    // Print outputs data in the configured format
    func (f *Formatter) Print(data interface{}) error {
        if f.Format == "json" {
            b, err := json.MarshalIndent(data, "", "  ")
            if err != nil {
                return err
            }
            fmt.Println(string(b))
        } else {
            fmt.Printf("%+v\n", data)
        }
        return nil
    }
    ```

- [x] **Task 7: Create Makefile** (AC: #5)
  - [x] 7.1 Create `Makefile`:

    ```makefile
    .PHONY: build build-all clean test install

    BINARY_NAME := monoguard
    DIST_DIR := dist
    VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
    COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

    LDFLAGS := -ldflags "-X github.com/j620656786206/MonoGuard/apps/cli/cmd.version=$(VERSION) \
        -X github.com/j620656786206/MonoGuard/apps/cli/cmd.commit=$(COMMIT) \
        -X github.com/j620656786206/MonoGuard/apps/cli/cmd.buildDate=$(BUILD_DATE)"

    # Build for current platform
    build: $(DIST_DIR)
    	go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) .
    	@echo "Built $(DIST_DIR)/$(BINARY_NAME)"
    	@ls -lh $(DIST_DIR)/$(BINARY_NAME)

    # Build for all platforms
    build-all: $(DIST_DIR)
    	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
    	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
    	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
    	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
    	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
    	@echo "Built binaries for all platforms"

    $(DIST_DIR):
    	mkdir -p $(DIST_DIR)

    # Run tests
    test:
    	go test -v ./...

    # Clean build artifacts
    clean:
    	rm -rf $(DIST_DIR)

    # Install to GOPATH/bin
    install:
    	go install $(LDFLAGS) .

    # Development: build and run help
    dev: build
    	$(DIST_DIR)/$(BINARY_NAME) --help
    ```

- [x] **Task 8: Update Package Configuration** (AC: #7)
  - [x] 8.1 Update `apps/cli/package.json` for npm distribution:
    ```json
    {
      "name": "@monoguard/cli",
      "version": "0.1.0",
      "description": "MonoGuard CLI - Monorepo dependency analysis tool",
      "private": true,
      "bin": {
        "monoguard": "./dist/monoguard"
      },
      "scripts": {
        "build": "make build",
        "build:all": "make build-all",
        "test": "make test",
        "clean": "make clean"
      },
      "keywords": ["monorepo", "dependency", "analysis", "cli"],
      "license": "MIT",
      "os": ["darwin", "linux", "win32"],
      "cpu": ["x64", "arm64"]
    }
    ```
  - [x] 8.2 Create or update `apps/cli/project.json`:
    ```json
    {
      "name": "@monoguard/cli",
      "projectType": "application",
      "sourceRoot": "apps/cli",
      "targets": {
        "build": {
          "executor": "nx:run-commands",
          "options": {
            "command": "make build",
            "cwd": "apps/cli"
          },
          "outputs": ["{projectRoot}/dist"]
        },
        "build-all": {
          "executor": "nx:run-commands",
          "options": {
            "command": "make build-all",
            "cwd": "apps/cli"
          },
          "outputs": ["{projectRoot}/dist"]
        },
        "test": {
          "executor": "nx:run-commands",
          "options": {
            "command": "make test",
            "cwd": "apps/cli"
          }
        },
        "clean": {
          "executor": "nx:run-commands",
          "options": {
            "command": "make clean",
            "cwd": "apps/cli"
          }
        }
      }
    }
    ```

- [x] **Task 9: Verification** (AC: #3, #5, #6, #7)
  - [x] 9.1 Run `make build` - verify binary builds successfully
  - [x] 9.2 Run `./dist/monoguard --help` - verify help output
  - [x] 9.3 Run `./dist/monoguard --version` - verify version info
  - [x] 9.4 Run `./dist/monoguard analyze` - verify placeholder output
  - [x] 9.5 Run `./dist/monoguard analyze --format json` - verify JSON output
  - [x] 9.6 Run `./dist/monoguard check` - verify exit code 0
  - [x] 9.7 Run `make test` - verify Go tests pass (all 40+ tests)
  - [x] 9.8 Check binary size: `du -h dist/monoguard` (7.3MB < 15MB limit)
  - [x] 9.9 Run `pnpm nx build cli` - verify Nx integration
  - [x] 9.10 Run `pnpm nx graph` - verify project appears correctly

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**

- **Framework:** Cobra for command management (industry standard - used by Docker, Kubernetes)
- **Configuration:** Viper for config file and environment variable management
- **Build:** Standard Go toolchain with Makefile
- **Distribution:** npm package wrapping native binaries

**Critical Constraints:**

- **Go Naming Conventions:** PascalCase exported, camelCase unexported, snake_case files
- **JSON Output:** Must use camelCase for all JSON fields
- **Exit Codes:** 0 for success, 1 for failure (CI/CD standard)
- **Config Priority:** CLI flags > Environment > Config file

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Go CLI Patterns:**

   ```go
   // ✅ CORRECT: Cobra command structure
   var analyzeCmd = &cobra.Command{
       Use:   "analyze [path]",
       Short: "Short description",
       Long:  `Long description`,
       Run: func(cmd *cobra.Command, args []string) {
           // Implementation
       },
   }

   // Register in init()
   func init() {
       rootCmd.AddCommand(analyzeCmd)
   }
   ```

2. **Viper Configuration:**

   ```go
   // Config file locations (in priority order):
   // 1. --config flag
   // 2. .monoguard.yaml in current directory
   // 3. ~/.monoguard/config.yaml

   // Environment variables with MONOGUARD_ prefix:
   // MONOGUARD_VERBOSE=true -> viper.GetBool("verbose")
   ```

3. **Output Format Pattern:**
   ```go
   // ✅ CORRECT: Support both text and JSON
   format := viper.GetString("format")
   if format == "json" {
       // Output JSON for machine parsing
   } else {
       // Output human-readable text
   }
   ```

### Project Structure Notes

**Target Directory Structure:**

```
apps/cli/
├── cmd/
│   ├── root.go           # Root command, Viper init
│   ├── analyze.go        # monoguard analyze
│   ├── check.go          # monoguard check
│   ├── fix.go            # monoguard fix
│   └── init_cmd.go       # monoguard init
├── pkg/
│   ├── config/           # Configuration management
│   │   └── config.go
│   └── output/           # Output formatting
│       └── output.go
├── dist/                 # Build output (gitignored)
│   └── monoguard
├── main.go               # Entry point
├── go.mod
├── go.sum
├── Makefile
├── package.json          # npm package config
├── project.json          # Nx project config
└── README.md
```

**Sharing Code with WASM (Future):**

- Analysis logic will be in `packages/analysis-engine/pkg/`
- CLI will import analysis packages directly (NOT via WASM)
- Same code, different entry points (native vs WASM)

### CLI Command Reference

**Commands to Implement (Placeholder):**

| Command                    | Description        | Exit Code      |
| -------------------------- | ------------------ | -------------- |
| `monoguard analyze [path]` | Full analysis      | N/A            |
| `monoguard check [path]`   | CI/CD validation   | 0=pass, 1=fail |
| `monoguard fix [path]`     | Fix suggestions    | N/A            |
| `monoguard init`           | Create config file | N/A            |

**Global Flags:**

| Flag            | Description               | Default           |
| --------------- | ------------------------- | ----------------- |
| `--config`      | Config file path          | `.monoguard.yaml` |
| `--verbose, -v` | Verbose output            | `false`           |
| `--format`      | Output format (text/json) | `text`            |
| `--version`     | Show version              | N/A               |

### Previous Story Intelligence

**From Story 1.1 (done):**

- `apps/cli/` exists with TypeScript implementation
- Package named `@monoguard/cli` in package.json
- TypeScript CLI uses Commander.js (being replaced)

**From Story 1.3 (ready-for-dev):**

- Go WASM project structure in `packages/analysis-engine/`
- Result type pattern established
- Same Go code can be shared with CLI

### Migration Notes (TypeScript → Go)

| TypeScript  | Go                                   |
| ----------- | ------------------------------------ |
| `commander` | `cobra`                              |
| `chalk`     | ANSI escape codes or `color` package |
| `inquirer`  | `survey` or `promptui` (future)      |
| `fs-extra`  | `os`, `io/fs`                        |
| `zod`       | Struct validation (future)           |

### References

- [Source: _bmad-output/planning-artifacts/architecture.md#CLI Tool (Go with Cobra)]
- [Source: _bmad-output/planning-artifacts/architecture.md#CLI Tool (Viper)]
- [Source: _bmad-output/project-context.md#Go Naming Conventions]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.4]
- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

None

### Completion Notes List

1. Successfully migrated TypeScript CLI to Go CLI with Cobra/Viper
2. All 40+ ATDD tests pass (unit tests + e2e tests)
3. Binary size is 7.3MB (well under 15MB limit)
4. Modified test infrastructure to use `ResetForTesting()` for test isolation
5. Commands write to `cmd.OutOrStdout()` for proper test capture
6. Removed `os.Exit(0)` from check command (Cobra handles exit codes)
7. Added reflection-based struct formatting for output package
8. [Code Review Fix] Added typed output structs with JSON marshaling to prevent JSON injection in analyze/check/fix commands
9. [Code Review Fix] Added JSON struct tags to pkg/config types for project-context.md compliance
10. [Code Review Fix] Removed outdated "Status: RED" comments from all test files
11. [Code Review Fix] Added assertions to TestViperInitialization for actual Viper env var verification
12. [Code Review Fix] Added assertions to TestEnvironmentVariableOverride to verify MONOGUARD\_ prefix works
13. [Code Review Fix] Added assertions to TestCLIFlagOverride to verify config file reading and flag override behavior

### File List

**Created/Modified:**

- `apps/cli/main.go` - Entry point
- `apps/cli/go.mod` - Go module (Cobra v1.10.2, Viper v1.21.0, Go 1.23.0)
- `apps/cli/go.sum` - Go dependencies lockfile
- `apps/cli/cmd/root.go` - Root command with Viper config and `ResetForTesting()`
- `apps/cli/cmd/analyze.go` - Analyze placeholder command (with JSON output struct)
- `apps/cli/cmd/check.go` - Check placeholder command (CI/CD friendly, with JSON output struct)
- `apps/cli/cmd/fix.go` - Fix placeholder command with --dry-run (with JSON output struct)
- `apps/cli/cmd/init_cmd.go` - Init placeholder command
- `apps/cli/cmd/root_test.go` - Root command tests
- `apps/cli/cmd/analyze_test.go` - Analyze command tests
- `apps/cli/cmd/check_test.go` - Check command tests
- `apps/cli/cmd/fix_test.go` - Fix command tests
- `apps/cli/cmd/init_cmd_test.go` - Init command tests
- `apps/cli/pkg/config/config.go` - Configuration types (with JSON tags)
- `apps/cli/pkg/config/config_test.go` - Configuration tests
- `apps/cli/pkg/output/output.go` - Output formatting with PrintTo method
- `apps/cli/pkg/output/output_test.go` - Output formatting tests
- `apps/cli/tests/cli_e2e_test.go` - End-to-end CLI tests
- `apps/cli/Makefile` - Build, test, clean targets
- `apps/cli/package.json` - npm package config (updated)
- `apps/cli/project.json` - Nx project config (updated)
- `apps/cli/.gitignore` - dist/ ignored

**Deleted:**

- `apps/cli/src/` - TypeScript source directory
- `apps/cli/tsconfig*.json` - TypeScript configs
- `apps/cli/eslint.config.cjs` - ESLint config
