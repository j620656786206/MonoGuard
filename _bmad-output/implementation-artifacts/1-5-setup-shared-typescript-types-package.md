# Story 1.5: Setup Shared TypeScript Types Package

Status: review

## Story

As a **developer**,
I want **a shared TypeScript types package with WASM-compatible type definitions**,
So that **I can share type definitions between web app, WASM adapter, and ensure cross-language consistency with Go**.

## Acceptance Criteria

1. **AC1: Result<T> Type Definition**
   - Given the types package
   - When I define the Result type
   - Then it matches the Go structure exactly:
     ```typescript
     interface Result<T> {
       data: T | null;
       error: { code: string; message: string } | null;
     }
     ```
   - And it can be used to type WASM function returns

2. **AC2: Core Analysis Types**
   - Given the types package
   - When I verify the type definitions
   - Then I have properly typed:
     - `DependencyGraph` - graph structure with nodes and edges
     - `Package` - basic package information
     - `CircularDependency` - cycle detection result (enhance existing)
     - `HealthScore` - health calculation result (enhance existing)
     - `AnalysisResult` - comprehensive analysis output

3. **AC3: WASM Adapter Interface**
   - Given the types package
   - When I define the MonoGuardAnalyzer interface
   - Then it defines:
     ```typescript
     interface MonoGuardAnalyzer {
       getVersion(): Promise<Result<{ version: string }>>;
       analyze(input: string): Promise<Result<AnalysisResult>>;
       check(input: string): Promise<Result<CheckResult>>;
     }
     ```
   - And it matches the Go WASM exports from Story 1.3

4. **AC4: Build Configuration**
   - Given the types package
   - When I run `pnpm nx build types`
   - Then:
     - TypeScript compiles without errors
     - Type declarations (.d.ts) are generated
     - Output goes to `packages/types/dist/`
     - ESM module format is used

5. **AC5: Import Verification**
   - Given the built types package
   - When I import types in apps/web
   - Then:
     - All types are accessible via `@monoguard/types`
     - TypeScript intellisense works correctly
     - No circular import issues

6. **AC6: Cross-Language Consistency Documentation**
   - Given the type definitions
   - When I review the package
   - Then:
     - JSDoc comments explain Go equivalent for each type
     - Error codes enum matches Go constants
     - Date fields documented as ISO 8601 strings

## Tasks / Subtasks

- [x] **Task 1: Create Result<T> Type** (AC: #1)
  - [x] 1.1 Create `packages/types/src/result.ts`:

    ````typescript
    /**
     * Result<T> - Unified response type for WASM functions
     *
     * This type MUST match the Go Result struct exactly:
     * ```go
     * type Result struct {
     *     Data  interface{} `json:"data"`
     *     Error *Error      `json:"error"`
     * }
     * ```
     *
     * @example
     * // Success case
     * { data: { healthScore: 85 }, error: null }
     *
     * // Error case
     * { data: null, error: { code: "PARSE_ERROR", message: "Invalid JSON" } }
     */
    export interface Result<T> {
      data: T | null;
      error: ResultError | null;
    }

    /**
     * ResultError - Error structure matching Go Error type
     *
     * Error codes use UPPER_SNAKE_CASE convention.
     */
    export interface ResultError {
      /** Error code in UPPER_SNAKE_CASE (e.g., PARSE_ERROR, CIRCULAR_DETECTED) */
      code: string;
      /** Human-readable error message */
      message: string;
    }

    /**
     * Standard error codes used by MonoGuard
     * Must match Go constants in internal/result/result.go
     */
    export const ErrorCodes = {
      PARSE_ERROR: 'PARSE_ERROR',
      INVALID_INPUT: 'INVALID_INPUT',
      CIRCULAR_DETECTED: 'CIRCULAR_DETECTED',
      ANALYSIS_FAILED: 'ANALYSIS_FAILED',
      WASM_ERROR: 'WASM_ERROR',
      TIMEOUT: 'TIMEOUT',
    } as const;

    export type ErrorCode = (typeof ErrorCodes)[keyof typeof ErrorCodes];

    // Type guards for Result handling
    export function isSuccess<T>(
      result: Result<T>
    ): result is Result<T> & { data: T; error: null } {
      return result.error === null && result.data !== null;
    }

    export function isError<T>(
      result: Result<T>
    ): result is Result<T> & { data: null; error: ResultError } {
      return result.error !== null;
    }
    ````

- [x] **Task 2: Create Core Analysis Types** (AC: #2)
  - [x] 2.1 Create `packages/types/src/analysis/graph.ts`:

    ```typescript
    /**
     * DependencyGraph - Core graph data structure
     *
     * Matches Go: pkg/types/graph.go
     */
    export interface DependencyGraph {
      /** Map of package name to Package details */
      nodes: Map<string, PackageNode> | Record<string, PackageNode>;
      /** List of dependency edges */
      edges: DependencyEdge[];
      /** Workspace root path */
      rootPath: string;
      /** Workspace type detected */
      workspaceType: WorkspaceType;
    }

    export interface PackageNode {
      /** Package name (e.g., "@monoguard/types") */
      name: string;
      /** Package version */
      version: string;
      /** Relative path from workspace root */
      path: string;
      /** Direct dependencies */
      dependencies: string[];
      /** Dev dependencies */
      devDependencies: string[];
      /** Peer dependencies */
      peerDependencies: string[];
    }

    export interface DependencyEdge {
      /** Source package name */
      from: string;
      /** Target package name */
      to: string;
      /** Dependency type */
      type: DependencyType;
      /** Version range specified */
      versionRange: string;
    }

    export type DependencyType =
      | 'production'
      | 'development'
      | 'peer'
      | 'optional';
    export type WorkspaceType = 'npm' | 'yarn' | 'pnpm' | 'unknown';
    ```

  - [x] 2.2 Create `packages/types/src/analysis/results.ts`:

    ```typescript
    import type { DependencyGraph, PackageNode } from './graph';

    /**
     * AnalysisResult - Complete analysis output
     *
     * Matches Go: pkg/types/analysis_result.go
     * All date fields use ISO 8601 format (e.g., "2026-01-15T10:30:00Z")
     */
    export interface AnalysisResult {
      /** Architecture health score (0-100) */
      healthScore: number;
      /** Total packages analyzed */
      packageCount: number;
      /** Detected circular dependencies */
      circularDependencies: CircularDependencyInfo[];
      /** Dependency graph data */
      graph: DependencyGraph;
      /** Analysis metadata */
      metadata: AnalysisMetadata;
      /** ISO 8601 timestamp */
      createdAt: string;
    }

    /**
     * CircularDependencyInfo - Enhanced circular dependency with fix suggestions
     *
     * Matches Go: pkg/types/circular.go
     */
    export interface CircularDependencyInfo {
      /** Packages involved in the cycle (in order) */
      cycle: string[];
      /** Type of circular dependency */
      type: 'direct' | 'indirect';
      /** Severity level */
      severity: 'critical' | 'warning' | 'info';
      /** Impact description */
      impact: string;
      /** Suggested fix strategy */
      fixStrategy?: FixStrategy;
      /** Refactoring complexity (1-10) */
      complexity: number;
    }

    export interface FixStrategy {
      /** Strategy type */
      type: 'extract_module' | 'dependency_injection' | 'boundary_refactor';
      /** Human-readable description */
      description: string;
      /** Step-by-step instructions */
      steps: string[];
      /** Files that need modification */
      affectedFiles: string[];
    }

    /**
     * CheckResult - Validation-only output for CI/CD
     *
     * Matches Go: pkg/types/check_result.go
     */
    export interface CheckResult {
      /** Overall pass/fail status */
      passed: boolean;
      /** List of errors found */
      errors: ValidationError[];
      /** List of warnings */
      warnings: ValidationWarning[];
      /** Health score (0-100) */
      healthScore: number;
    }

    export interface ValidationError {
      /** Error code */
      code: string;
      /** Error message */
      message: string;
      /** Related file path (optional) */
      file?: string;
      /** Line number (optional) */
      line?: number;
    }

    export interface ValidationWarning {
      /** Warning code */
      code: string;
      /** Warning message */
      message: string;
      /** Related file path (optional) */
      file?: string;
    }

    export interface AnalysisMetadata {
      /** MonoGuard version */
      version: string;
      /** Analysis duration in milliseconds */
      durationMs: number;
      /** Number of files processed */
      filesProcessed: number;
      /** Workspace type detected */
      workspaceType: WorkspaceType;
    }

    // Re-export for convenience
    export type { WorkspaceType } from './graph';
    ```

- [x] **Task 3: Create WASM Adapter Interface** (AC: #3)
  - [x] 3.1 Create `packages/types/src/wasm/adapter.ts`:

    ````typescript
    import type { Result } from '../result';
    import type { AnalysisResult, CheckResult } from '../analysis/results';

    /**
     * MonoGuardAnalyzer - WASM adapter interface
     *
     * This interface defines the contract between TypeScript and Go WASM.
     * All methods return Promise<Result<T>> to handle async WASM calls.
     *
     * Implementation is in packages/analysis-engine (Go WASM).
     * TypeScript adapter is implemented in Story 2.7.
     *
     * @example
     * ```typescript
     * const analyzer: MonoGuardAnalyzer = await loadWasm();
     * const result = await analyzer.analyze(JSON.stringify(workspaceConfig));
     * if (isSuccess(result)) {
     *   console.log(result.data.healthScore);
     * }
     * ```
     */
    export interface MonoGuardAnalyzer {
      /**
       * Get MonoGuard version
       * @returns Version information
       */
      getVersion(): Promise<Result<VersionInfo>>;

      /**
       * Analyze workspace dependencies
       * @param input JSON string of workspace configuration
       * @returns Complete analysis result
       */
      analyze(input: string): Promise<Result<AnalysisResult>>;

      /**
       * Check workspace for CI/CD validation
       * @param input JSON string of workspace configuration
       * @returns Pass/fail result with errors
       */
      check(input: string): Promise<Result<CheckResult>>;
    }

    export interface VersionInfo {
      version: string;
      commit?: string;
      buildDate?: string;
    }

    /**
     * WasmLoaderOptions - Options for loading WASM module
     */
    export interface WasmLoaderOptions {
      /** Path to monoguard.wasm file */
      wasmPath?: string;
      /** Timeout for WASM initialization (ms) */
      timeout?: number;
    }

    /**
     * WasmLoadResult - Result of loading WASM module
     */
    export interface WasmLoadResult {
      analyzer: MonoGuardAnalyzer;
      loadTimeMs: number;
    }
    ````

- [x] **Task 4: Update Index Exports** (AC: #4, #5)
  - [x] 4.1 Create `packages/types/src/analysis/index.ts`:
    ```typescript
    export * from './graph';
    export * from './results';
    ```
  - [x] 4.2 Create `packages/types/src/wasm/index.ts`:
    ```typescript
    export * from './adapter';
    ```
  - [x] 4.3 Update `packages/types/src/index.ts` to include new exports:

    ```typescript
    // Existing exports
    export * from './api';
    export * from './domain';
    export * from './auth';
    export * from './common';

    // New WASM-compatible types
    export * from './result';
    export * from './analysis';
    export * from './wasm';

    // Re-export commonly used types (add new ones)
    export type { Result, ResultError, ErrorCode } from './result';

    export { ErrorCodes, isSuccess, isError } from './result';

    export type {
      DependencyGraph,
      PackageNode,
      DependencyEdge,
      DependencyType,
      WorkspaceType,
    } from './analysis/graph';

    export type {
      AnalysisResult,
      CircularDependencyInfo,
      FixStrategy,
      CheckResult,
      ValidationError,
      ValidationWarning,
      AnalysisMetadata,
    } from './analysis/results';

    export type {
      MonoGuardAnalyzer,
      VersionInfo,
      WasmLoaderOptions,
      WasmLoadResult,
    } from './wasm/adapter';
    ```

- [x] **Task 5: Update Build Configuration** (AC: #4)
  - [x] 5.1 Verify `packages/types/tsconfig.json` is correct:
    ```json
    {
      "extends": "../../tsconfig.base.json",
      "compilerOptions": {
        "outDir": "./dist",
        "declaration": true,
        "declarationMap": true,
        "module": "ESNext",
        "moduleResolution": "bundler",
        "target": "ES2022",
        "strict": true,
        "esModuleInterop": true,
        "skipLibCheck": true,
        "forceConsistentCasingInFileNames": true,
        "rootDir": "./src"
      },
      "include": ["src/**/*"],
      "exclude": ["node_modules", "dist", "**/*.test.ts"]
    }
    ```
  - [x] 5.2 Create or update `packages/types/project.json`:
    ```json
    {
      "name": "@monoguard/types",
      "projectType": "library",
      "sourceRoot": "packages/types/src",
      "targets": {
        "build": {
          "executor": "nx:run-commands",
          "options": {
            "command": "tsc -p tsconfig.json",
            "cwd": "packages/types"
          },
          "outputs": ["{projectRoot}/dist"]
        },
        "lint": {
          "executor": "nx:run-commands",
          "options": {
            "command": "eslint src --ext ts",
            "cwd": "packages/types"
          }
        },
        "test": {
          "executor": "nx:run-commands",
          "options": {
            "command": "jest",
            "cwd": "packages/types"
          }
        }
      }
    }
    ```

- [x] **Task 6: Create Type Tests** (AC: #5)
  - [x] 6.1 Create `packages/types/src/__tests__/result.test.ts`:

    ```typescript
    import { describe, it, expect } from '@jest/globals';
    import { isSuccess, isError, type Result } from '../result';

    describe('Result type guards', () => {
      it('isSuccess returns true for successful result', () => {
        const result: Result<number> = { data: 42, error: null };
        expect(isSuccess(result)).toBe(true);
        expect(isError(result)).toBe(false);
      });

      it('isError returns true for error result', () => {
        const result: Result<number> = {
          data: null,
          error: { code: 'TEST_ERROR', message: 'Test error' },
        };
        expect(isError(result)).toBe(true);
        expect(isSuccess(result)).toBe(false);
      });
    });
    ```

  - [x] 6.2 Create `packages/types/src/__tests__/analysis.test.ts`:

    ```typescript
    import { describe, it, expect } from '@jest/globals';
    import type { AnalysisResult, DependencyGraph } from '../analysis';

    describe('Analysis types', () => {
      it('AnalysisResult type is correctly defined', () => {
        const result: AnalysisResult = {
          healthScore: 85,
          packageCount: 10,
          circularDependencies: [],
          graph: {
            nodes: {},
            edges: [],
            rootPath: '/workspace',
            workspaceType: 'pnpm',
          },
          metadata: {
            version: '0.1.0',
            durationMs: 1500,
            filesProcessed: 50,
            workspaceType: 'pnpm',
          },
          createdAt: '2026-01-15T10:30:00Z',
        };

        expect(result.healthScore).toBe(85);
        expect(result.graph.workspaceType).toBe('pnpm');
      });
    });
    ```

- [x] **Task 7: Verification** (AC: #4, #5, #6)
  - [x] 7.1 Run `pnpm nx build types` - verify build succeeds
  - [x] 7.2 Check dist/ output includes .d.ts files
  - [x] 7.3 Run `pnpm nx test types` - verify tests pass (18/18 passing)
  - [x] 7.4 Verify imports work in apps/web:
    ```typescript
    import {
      Result,
      isSuccess,
      AnalysisResult,
      MonoGuardAnalyzer,
    } from '@monoguard/types';
    ```
  - [x] 7.5 Run `pnpm nx graph` - verify project dependencies

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**

- **Purpose:** Cross-language type consistency between TypeScript and Go
- **Pattern:** All types must match Go struct definitions exactly
- **Constraint:** JSON field names must use camelCase (not snake_case)

**Critical Constraints:**

- **Result<T> is mandatory:** All WASM returns use this pattern
- **ISO 8601 dates:** All date strings, never Unix timestamps
- **Error codes:** UPPER_SNAKE_CASE (e.g., PARSE_ERROR)
- **Type guards:** Provide isSuccess/isError for safe type narrowing

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Cross-Language Type Matching:**

   ```typescript
   // ✅ TypeScript must match Go
   interface AnalysisResult {
     healthScore: number; // Go: int `json:"healthScore"`
     createdAt: string; // Go: string `json:"createdAt"` (ISO 8601)
   }

   // ❌ WRONG: snake_case breaks Go compatibility
   interface AnalysisResult {
     health_score: number; // Go expects camelCase JSON
   }
   ```

2. **Result Pattern Usage:**

   ```typescript
   // ✅ CORRECT: Use type guards
   const result = await analyzer.analyze(input);
   if (isSuccess(result)) {
     // TypeScript knows result.data is non-null here
     console.log(result.data.healthScore);
   } else {
     // TypeScript knows result.error is non-null here
     console.error(result.error.message);
   }

   // ❌ WRONG: Direct access without type guard
   console.log(result.data.healthScore); // May be null!
   ```

3. **Date Format:**

   ```typescript
   // ✅ CORRECT: ISO 8601 string
   createdAt: '2026-01-15T10:30:00Z';

   // ❌ WRONG: Unix timestamp
   createdAt: 1736939400;
   ```

### Project Structure Notes

**Target Directory Structure:**

```
packages/types/
├── src/
│   ├── analysis/
│   │   ├── graph.ts           # DependencyGraph, PackageNode, etc.
│   │   ├── results.ts         # AnalysisResult, CheckResult, etc.
│   │   └── index.ts           # Re-exports
│   ├── wasm/
│   │   ├── adapter.ts         # MonoGuardAnalyzer interface
│   │   └── index.ts           # Re-exports
│   ├── __tests__/
│   │   ├── result.test.ts     # Result type tests
│   │   └── analysis.test.ts   # Analysis type tests
│   ├── api.ts                 # Existing API types
│   ├── auth.ts                # Existing auth types
│   ├── common.ts              # Existing common types
│   ├── domain.ts              # Existing domain types
│   ├── result.ts              # NEW: Result<T> type
│   └── index.ts               # Main exports
├── dist/                      # Build output
├── package.json
├── project.json
└── tsconfig.json
```

### Type Mapping Reference

| Go Type         | TypeScript Type     |
| --------------- | ------------------- |
| `int`           | `number`            |
| `int64`         | `number`            |
| `float64`       | `number`            |
| `string`        | `string`            |
| `bool`          | `boolean`           |
| `[]string`      | `string[]`          |
| `map[string]T`  | `Record<string, T>` |
| `*T` (nullable) | `T \| null`         |
| `interface{}`   | `unknown`           |

### Existing Types to Keep

The existing types in `domain.ts`, `api.ts`, etc. should be kept. They are used for:

- Backend API communication (if needed in future)
- Legacy compatibility

New WASM-focused types are additive, not replacing existing types.

### Previous Story Intelligence

**From Story 1.1 (done):**

- `packages/types/` exists with existing type definitions
- Package named `@monoguard/types`
- Path mapping configured in tsconfig.base.json

**From Story 1.3 (ready-for-dev):**

- Go Result type defined in `internal/result/result.go`
- TypeScript Result<T> must match exactly
- Error codes must match Go constants

### References

- [Source: _bmad-output/project-context.md#Cross-Language Type Mapping]
- [Source: _bmad-output/project-context.md#Result Pattern]
- [Source: _bmad-output/planning-artifacts/epics.md#Story 1.5]
- [Source: _bmad-output/planning-artifacts/architecture.md#Integration Points]
- [TypeScript Handbook - Generics](https://www.typescriptlang.org/docs/handbook/2/generics.html)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

None required - implementation proceeded without issues.

### Completion Notes List

- **Task 1 Complete:** Created `result.ts` with `Result<T>`, `ResultError`, `ErrorCodes`, and type guards (`isSuccess`, `isError`)
- **Task 2 Complete:** Created `analysis/graph.ts` with `DependencyGraph`, `PackageNode`, `DependencyEdge`, `DependencyType`, `WorkspaceType`; Created `analysis/results.ts` with `AnalysisResult`, `CircularDependencyInfo`, `FixStrategy`, `CheckResult`, `ValidationError`, `ValidationWarning`, `AnalysisMetadata`
- **Task 3 Complete:** Created `wasm/adapter.ts` with `MonoGuardAnalyzer`, `VersionInfo`, `WasmLoaderOptions`, `WasmLoadResult`
- **Task 4 Complete:** Updated `index.ts` to export all new types; created `analysis/index.ts` and `wasm/index.ts` barrel files
- **Task 5 Complete:** Created `project.json` for Nx integration with build, lint, test targets
- **Task 6 Complete:** Created comprehensive tests in `__tests__/result.test.ts` (9 tests) and `__tests__/analysis.test.ts` (9 tests) - all 18 tests passing
- **Task 7 Complete:** Verified build succeeds, .d.ts files generated, tests pass, lint passes, Nx project recognized

### File List

**New Files:**

- packages/types/src/result.ts
- packages/types/src/analysis/graph.ts
- packages/types/src/analysis/results.ts
- packages/types/src/analysis/index.ts
- packages/types/src/wasm/adapter.ts
- packages/types/src/wasm/index.ts
- packages/types/src/**tests**/result.test.ts
- packages/types/src/**tests**/analysis.test.ts
- packages/types/project.json
- packages/types/jest.config.cjs

**Modified Files:**

- packages/types/src/index.ts (added new exports)
- packages/types/package.json (added ts-jest dependency)
- \_bmad-output/implementation-artifacts/sprint-status.yaml (status: in-progress → review)

## Change Log

- 2026-01-16: Implemented Story 1.5 - Created WASM-compatible TypeScript types package with Result<T>, analysis types, and WASM adapter interface. All 18 tests passing.
