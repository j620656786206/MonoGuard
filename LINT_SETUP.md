# Husky + Lint-Staged Setup for MonoGuard

This document describes the pre-commit hooks and linting setup for the MonoGuard monorepo.

## Overview

The project uses Husky for git hooks and lint-staged to run linting and formatting only on staged files. This ensures code quality and consistency across the entire monorepo.

## Pre-commit Hooks Configuration

### What runs on commit:

- **TypeScript/JavaScript files**: Prettier formatting for all .ts, .tsx, .js, .jsx files
- **Go files (Backend)**: gofmt formatting for files in apps/api
- **Configuration files**: Prettier for JSON, Markdown, YAML files

### File patterns covered:

- `*.{js,jsx,ts,tsx}` - All TypeScript/JavaScript files (Prettier formatting)
- `apps/api/**/*.go` - Backend Go files (gofmt formatting)
- `*.{json,md,yml,yaml}` - Configuration and documentation files (Prettier formatting)

### Linting Strategy:

The pre-commit hook focuses on **formatting only** to ensure consistency. For linting (ESLint, etc.), use:
- Manual linting: `pnpm lint` or `pnpm lint:fix`
- CI/CD pipeline linting for comprehensive checks
- IDE integration for real-time linting feedback

## Setup Requirements

### For JavaScript/TypeScript linting:

All required tools are already installed via npm/pnpm:

- ESLint
- Prettier
- TypeScript

### For Go formatting:

Go formatting tools are included with the standard Go installation:
- `gofmt` - Standard Go formatter (used in pre-commit hooks)
- `go vet` - Go code analysis tool

### For Go linting (optional):

For additional Go linting beyond formatting, install golangci-lint:

```bash
# On macOS:
brew install golangci-lint

# On Linux:
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# On Windows:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
```

## Available Scripts

Run these scripts manually for formatting and linting:

```bash
# Format all files with Prettier
pnpm format

# Check formatting without changes
pnpm format:check

# Format Go files
pnpm go:fmt

# Vet Go code (requires Go installation)
pnpm go:vet

# Lint Go code with golangci-lint (requires golangci-lint installation)
pnpm go:lint

# Run lint-staged manually (same as pre-commit hook)
pnpm pre-commit

# Standard lint commands (uses Nx - may need ESLint config fixes)
pnpm lint          # Lint all projects
pnpm lint:fix      # Lint and fix all projects
```

## How it Works

1. **Pre-commit Hook**: When you commit, Husky triggers the pre-commit hook
2. **Lint-Staged**: Only processes files that are staged for commit
3. **Tool Execution**: Runs formatting tools based on file type:
   - TypeScript/JavaScript: Prettier formatting
   - Go: gofmt formatting
   - JSON/Markdown/YAML: Prettier formatting
4. **Commit Prevention**: If any formatting tool fails, the commit is blocked

## Configuration Files

- `.husky/pre-commit` - Git hook that runs lint-staged
- `package.json` - Contains lint-staged configuration and scripts
- Individual ESLint configs in each app directory
- Prettier config at project root level

## Troubleshooting

### Common Issues:

1. **Go tools not found**:

   ```bash
   # Ensure Go is installed and tools are in PATH
   go version
   which golangci-lint
   ```

2. **ESLint errors blocking commit**:

   ```bash
   # Run lint with fix manually
   pnpm lint:fix
   ```

3. **Prettier formatting conflicts**:

   ```bash
   # Format all files manually
   pnpm format
   ```

4. **Skip hooks temporarily** (not recommended):
   ```bash
   git commit --no-verify
   ```

### Manual lint-staged run:

```bash
pnpm lint-staged
```

## Benefits

- **Consistent Code Style**: All code follows the same formatting standards
- **Early Error Detection**: Linting catches issues before they reach CI/CD
- **Faster CI**: Only clean, formatted code gets committed
- **Monorepo Support**: Different tools for different parts of the codebase
- **Selective Processing**: Only staged files are processed, making it fast

## Integration with IDEs

Consider configuring your IDE to:

- Run ESLint and Prettier on save
- Show linting errors inline
- Format on save with Prettier

This reduces pre-commit hook failures and improves development experience.
