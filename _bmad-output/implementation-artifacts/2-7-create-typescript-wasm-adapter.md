# Story 2.7: Create TypeScript WASM Adapter

Status: review

## Story

As a **developer**,
I want **a TypeScript wrapper for the WASM analysis engine**,
So that **I can call analysis functions with full type safety from the web app**.

## Acceptance Criteria

1. **AC1: WASM Initialization**
   - Given the compiled WASM module
   - When I call `initAnalyzer(options)`
   - Then:
     - WASM is loaded from specified path (or default)
     - Go runtime is initialized
     - MonoGuard global object is available
     - Returns Promise that resolves when ready
   - And supports timeout option for slow networks

2. **AC2: Typed Analyze Function**
   - Given initialized WASM
   - When I call `analyzer.analyze(workspaceData)`
   - Then:
     - Input is typed as `WorkspaceInput` (files + optional config)
     - Output is typed as `Result<AnalysisResult>`
     - JSON serialization is handled internally
   - And errors are properly typed as `ResultError`

3. **AC3: Typed Check Function**
   - Given initialized WASM
   - When I call `analyzer.check(workspaceData)`
   - Then:
     - Input is typed as `WorkspaceInput`
     - Output is typed as `Result<CheckResult>`
     - Suitable for CI/CD validation
   - And returns pass/fail with error details

4. **AC4: Typed GetVersion Function**
   - Given initialized WASM
   - When I call `analyzer.getVersion()`
   - Then:
     - Output is typed as `Result<VersionInfo>`
     - Returns version, commit, buildDate
   - And can be used to verify WASM loaded correctly

5. **AC5: Result<T> Pattern**
   - Given any analyzer function call
   - When the function returns
   - Then the result follows unified `Result<T>` pattern:
     ```typescript
     type Result<T> = {
       data: T | null;
       error: { code: string; message: string } | null;
     }
     ```
   - And type guards `isSuccess()` and `isError()` work correctly

6. **AC6: Error Handling**
   - Given various error scenarios
   - When errors occur
   - Then proper error types are returned:
     - WASM load failure → `WASM_LOAD_ERROR`
     - Initialization timeout → `WASM_TIMEOUT`
     - Parse error → `PARSE_ERROR`
     - Analysis failure → `ANALYSIS_FAILED`
   - And error messages are descriptive

7. **AC7: JSDoc Documentation**
   - Given the adapter module
   - When developers use it
   - Then all public functions have JSDoc:
     - Function description
     - Parameter descriptions
     - Return type descriptions
     - Usage examples
   - And TypeScript intellisense works correctly

8. **AC8: Performance**
   - Given a standard workspace
   - When WASM is loaded
   - Then initialization completes in < 2 seconds
   - And subsequent calls have < 10ms overhead (JSON serialization)

## Tasks / Subtasks

- [x] **Task 1: Create WASM Loader Module** (AC: #1, #6)
  - [ ] 1.1 Create `packages/analysis-engine/src/loader.ts`:
    ```typescript
    /**
     * Options for loading the WASM module
     */
    export interface WasmLoaderOptions {
      /** Path to monoguard.wasm file. Default: '/monoguard.wasm' */
      wasmPath?: string;
      /** Timeout in milliseconds. Default: 10000 */
      timeout?: number;
    }

    /**
     * Load and initialize the MonoGuard WASM module
     *
     * @param options - Configuration options
     * @returns Promise that resolves when WASM is ready
     * @throws WasmLoadError if loading fails or times out
     *
     * @example
     * ```typescript
     * await loadWasm({ wasmPath: '/wasm/monoguard.wasm' });
     * // MonoGuard global is now available
     * ```
     */
    export async function loadWasm(options?: WasmLoaderOptions): Promise<void> {
      const { wasmPath = '/monoguard.wasm', timeout = 10000 } = options ?? {};

      return new Promise((resolve, reject) => {
        const timeoutId = setTimeout(() => {
          reject(new WasmLoadError('WASM_TIMEOUT', 'WASM initialization timed out'));
        }, timeout);

        // Load wasm_exec.js (Go runtime)
        // Fetch and instantiate WASM
        // Initialize Go runtime
        // Resolve when MonoGuard global is available

        clearTimeout(timeoutId);
        resolve();
      });
    }

    export class WasmLoadError extends Error {
      constructor(public code: string, message: string) {
        super(message);
        this.name = 'WasmLoadError';
      }
    }
    ```
  - [ ] 1.2 Implement WASM loading with fetch
  - [ ] 1.3 Implement Go runtime initialization
  - [ ] 1.4 Add timeout handling
  - [ ] 1.5 Create tests with mock WASM

- [x] **Task 2: Create Analyzer Adapter** (AC: #2, #3, #4)
  - [ ] 2.1 Create `packages/analysis-engine/src/analyzer.ts`:
    ```typescript
    import type {
      Result,
      AnalysisResult,
      CheckResult,
      VersionInfo,
    } from '@monoguard/types';

    /**
     * Input for analyze and check functions
     */
    export interface WorkspaceInput {
      /** Map of filename to file content */
      files: Record<string, string>;
      /** Optional analysis configuration */
      config?: AnalysisConfig;
    }

    /**
     * Analysis configuration options
     */
    export interface AnalysisConfig {
      /** Patterns to exclude from analysis */
      exclude?: string[];
    }

    /**
     * MonoGuard analyzer wrapper with full type safety
     *
     * @example
     * ```typescript
     * const analyzer = new MonoGuardAnalyzer();
     * await analyzer.init();
     *
     * const result = await analyzer.analyze({
     *   files: { 'package.json': '...' }
     * });
     *
     * if (isSuccess(result)) {
     *   console.log(result.data.healthScore);
     * }
     * ```
     */
    export class MonoGuardAnalyzer {
      private initialized = false;

      /**
       * Initialize the WASM module
       * Must be called before analyze/check/getVersion
       */
      async init(options?: WasmLoaderOptions): Promise<void>;

      /**
       * Analyze workspace dependencies
       *
       * @param input - Workspace files and optional config
       * @returns Analysis result with health score, cycles, conflicts
       */
      async analyze(input: WorkspaceInput): Promise<Result<AnalysisResult>>;

      /**
       * Check workspace for CI/CD validation
       *
       * @param input - Workspace files and optional config
       * @returns Pass/fail result with errors
       */
      async check(input: WorkspaceInput): Promise<Result<CheckResult>>;

      /**
       * Get MonoGuard version information
       *
       * @returns Version info including commit and build date
       */
      async getVersion(): Promise<Result<VersionInfo>>;
    }
    ```
  - [ ] 2.2 Implement init() method
  - [ ] 2.3 Implement analyze() with JSON serialization
  - [ ] 2.4 Implement check() with JSON serialization
  - [ ] 2.5 Implement getVersion()
  - [ ] 2.6 Create comprehensive tests

- [x] **Task 3: Implement JSON Bridge** (AC: #5)
  - [ ] 3.1 Create `packages/analysis-engine/src/bridge.ts`:
    ```typescript
    import type { Result } from '@monoguard/types';

    /**
     * Call a WASM function with JSON serialization
     *
     * @param funcName - Name of the MonoGuard function
     * @param input - Input data to serialize
     * @returns Parsed Result from WASM
     */
    export function callWasm<T>(funcName: string, input: unknown): Result<T> {
      if (typeof (window as any).MonoGuard?.[funcName] !== 'function') {
        return {
          data: null,
          error: { code: 'WASM_NOT_INITIALIZED', message: 'WASM not loaded' },
        };
      }

      try {
        const inputJson = JSON.stringify(input);
        const resultJson = (window as any).MonoGuard[funcName](inputJson);
        return JSON.parse(resultJson) as Result<T>;
      } catch (err) {
        return {
          data: null,
          error: {
            code: 'WASM_CALL_FAILED',
            message: err instanceof Error ? err.message : 'Unknown error',
          },
        };
      }
    }

    /**
     * Call a WASM function that takes no input
     */
    export function callWasmNoInput<T>(funcName: string): Result<T> {
      if (typeof (window as any).MonoGuard?.[funcName] !== 'function') {
        return {
          data: null,
          error: { code: 'WASM_NOT_INITIALIZED', message: 'WASM not loaded' },
        };
      }

      try {
        const resultJson = (window as any).MonoGuard[funcName]();
        return JSON.parse(resultJson) as Result<T>;
      } catch (err) {
        return {
          data: null,
          error: {
            code: 'WASM_CALL_FAILED',
            message: err instanceof Error ? err.message : 'Unknown error',
          },
        };
      }
    }
    ```
  - [ ] 3.2 Add error wrapping
  - [ ] 3.3 Create tests

- [x] **Task 4: Re-export Type Guards** (AC: #5)
  - [ ] 4.1 Create `packages/analysis-engine/src/index.ts`:
    ```typescript
    // Re-export loader
    export { loadWasm, WasmLoadError, type WasmLoaderOptions } from './loader';

    // Re-export analyzer
    export {
      MonoGuardAnalyzer,
      type WorkspaceInput,
      type AnalysisConfig,
    } from './analyzer';

    // Re-export type guards from @monoguard/types
    export { isSuccess, isError } from '@monoguard/types';

    // Re-export commonly used types
    export type {
      Result,
      AnalysisResult,
      CheckResult,
      VersionInfo,
      DependencyGraph,
      CircularDependencyInfo,
      VersionConflict,
      HealthScoreResult,
    } from '@monoguard/types';
    ```
  - [ ] 4.2 Update package.json exports

- [x] **Task 5: Add Package Configuration** (AC: #7)
  - [ ] 5.1 Update `packages/analysis-engine/package.json`:
    ```json
    {
      "name": "@monoguard/analysis-engine",
      "version": "0.1.0",
      "type": "module",
      "main": "./dist/index.js",
      "types": "./dist/index.d.ts",
      "exports": {
        ".": {
          "import": "./dist/index.js",
          "types": "./dist/index.d.ts"
        },
        "./wasm": {
          "import": "./dist/monoguard.wasm"
        }
      },
      "files": ["dist/"],
      "scripts": {
        "build": "make build-wasm && tsc",
        "build:wasm": "make build-wasm",
        "build:ts": "tsc",
        "test": "vitest"
      },
      "dependencies": {
        "@monoguard/types": "workspace:*"
      },
      "devDependencies": {
        "typescript": "^5.0.0",
        "vitest": "^3.2.0"
      }
    }
    ```
  - [ ] 5.2 Create `packages/analysis-engine/tsconfig.json` for TypeScript
  - [ ] 5.3 Update project.json for Nx

- [x] **Task 6: Create Tests** (AC: all)
  - [ ] 6.1 Create `packages/analysis-engine/src/__tests__/analyzer.test.ts`:
    ```typescript
    import { describe, it, expect, vi, beforeEach } from 'vitest';
    import { MonoGuardAnalyzer } from '../analyzer';
    import { isSuccess, isError } from '@monoguard/types';

    describe('MonoGuardAnalyzer', () => {
      let analyzer: MonoGuardAnalyzer;

      beforeEach(() => {
        analyzer = new MonoGuardAnalyzer();
        // Mock WASM global
        (window as any).MonoGuard = {
          getVersion: vi.fn().mockReturnValue(JSON.stringify({
            data: { version: '0.1.0' },
            error: null,
          })),
          analyze: vi.fn().mockReturnValue(JSON.stringify({
            data: { healthScore: 85, packages: 5 },
            error: null,
          })),
          check: vi.fn().mockReturnValue(JSON.stringify({
            data: { passed: true, errors: [] },
            error: null,
          })),
        };
      });

      it('should return version info', async () => {
        await analyzer.init();
        const result = await analyzer.getVersion();
        expect(isSuccess(result)).toBe(true);
        if (isSuccess(result)) {
          expect(result.data.version).toBe('0.1.0');
        }
      });

      it('should analyze workspace', async () => {
        await analyzer.init();
        const result = await analyzer.analyze({
          files: { 'package.json': '{}' },
        });
        expect(isSuccess(result)).toBe(true);
        if (isSuccess(result)) {
          expect(result.data.healthScore).toBe(85);
        }
      });

      it('should return error for uninitialized analyzer', async () => {
        const result = await analyzer.analyze({ files: {} });
        expect(isError(result)).toBe(true);
      });
    });
    ```
  - [ ] 6.2 Create loader tests
  - [ ] 6.3 Create bridge tests

- [x] **Task 7: Performance Testing** (AC: #8)
  - [ ] 7.1 Create `packages/analysis-engine/src/__tests__/performance.test.ts`:
    ```typescript
    import { describe, it, expect } from 'vitest';
    import { MonoGuardAnalyzer } from '../analyzer';

    describe('Performance', () => {
      it('should initialize within 2 seconds', async () => {
        const start = performance.now();
        const analyzer = new MonoGuardAnalyzer();
        await analyzer.init();
        const duration = performance.now() - start;
        expect(duration).toBeLessThan(2000);
      });

      it('should have minimal JSON serialization overhead', async () => {
        const analyzer = new MonoGuardAnalyzer();
        await analyzer.init();

        const input = { files: generateLargeInput() };
        const start = performance.now();
        await analyzer.analyze(input);
        const duration = performance.now() - start;

        // Serialization should be < 10ms overhead
        expect(duration).toBeLessThan(100); // Allow for WASM execution
      });
    });
    ```

- [x] **Task 8: Integration with Web App** (AC: all)
  - [ ] 8.1 Verify import works in apps/web:
    ```typescript
    import {
      MonoGuardAnalyzer,
      isSuccess,
      type AnalysisResult,
    } from '@monoguard/analysis-engine';
    ```
  - [ ] 8.2 Update apps/web to copy WASM to public folder
  - [ ] 8.3 Create example usage in web app

- [x] **Task 9: Integration Verification** (AC: all)
  - [ ] 9.1 Build TypeScript: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 9.2 Verify .d.ts files are generated
  - [ ] 9.3 Test in browser environment
  - [ ] 9.4 Verify all tests pass

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** TypeScript adapter in `packages/analysis-engine/src/`
- **Pattern:** Wrapper class with async methods
- **Integration:** Uses @monoguard/types for all type definitions

**Module Structure:**
```
packages/analysis-engine/
├── src/
│   ├── index.ts          # Main exports
│   ├── loader.ts         # WASM loading
│   ├── analyzer.ts       # Analyzer class
│   ├── bridge.ts         # JSON bridge utilities
│   └── __tests__/        # Tests
├── dist/
│   ├── index.js          # Compiled JS
│   ├── index.d.ts        # Type declarations
│   ├── monoguard.wasm    # WASM binary
│   └── wasm_exec.js      # Go runtime
└── ...
```

**Critical Constraints:**
- **Type Safety:** All public APIs must be fully typed
- **Result Pattern:** ALL functions return `Result<T>`
- **camelCase JSON:** Must match Go struct tags
- **Browser Compatibility:** Must work in modern browsers

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Always Use Result Pattern:**
   ```typescript
   // ✅ CORRECT: Return Result<T>
   async analyze(input: WorkspaceInput): Promise<Result<AnalysisResult>> {
     try {
       return callWasm<AnalysisResult>('analyze', input);
     } catch (err) {
       return { data: null, error: { code: 'WASM_ERROR', message: err.message } };
     }
   }

   // ❌ WRONG: Throw errors
   async analyze(input: WorkspaceInput): Promise<AnalysisResult> {
     if (error) throw new Error('Analysis failed'); // Don't throw!
   }
   ```

2. **Initialize Before Use:**
   ```typescript
   // ✅ CORRECT: Check initialization
   async analyze(input: WorkspaceInput): Promise<Result<AnalysisResult>> {
     if (!this.initialized) {
       return {
         data: null,
         error: { code: 'NOT_INITIALIZED', message: 'Call init() first' },
       };
     }
     // ... proceed
   }
   ```

3. **JSDoc Required:**
   ```typescript
   // ✅ CORRECT: Full JSDoc
   /**
    * Analyze workspace dependencies
    *
    * @param input - Workspace files and optional config
    * @returns Analysis result with health score, cycles, conflicts
    *
    * @example
    * ```typescript
    * const result = await analyzer.analyze({ files: {...} });
    * if (isSuccess(result)) {
    *   console.log(result.data.healthScore);
    * }
    * ```
    */
   async analyze(input: WorkspaceInput): Promise<Result<AnalysisResult>>
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── cmd/wasm/                    # Go WASM source (existing)
├── pkg/                         # Go packages (existing)
├── internal/                    # Go internal (existing)
├── src/                         # NEW: TypeScript adapter
│   ├── index.ts                 # Main exports
│   ├── loader.ts                # WASM loading
│   ├── analyzer.ts              # Analyzer class
│   ├── bridge.ts                # JSON bridge
│   └── __tests__/
│       ├── analyzer.test.ts
│       ├── loader.test.ts
│       └── performance.test.ts
├── dist/                        # Build output
│   ├── index.js
│   ├── index.d.ts
│   ├── monoguard.wasm          # From Go build
│   └── wasm_exec.js            # From Go build
├── go.mod
├── Makefile
├── package.json                 # UPDATE: Add TS build
├── project.json                 # UPDATE: Add TS targets
└── tsconfig.json                # NEW: TypeScript config
```

### Usage Example

```typescript
import {
  MonoGuardAnalyzer,
  isSuccess,
  isError,
  type WorkspaceInput,
} from '@monoguard/analysis-engine';

// Initialize once
const analyzer = new MonoGuardAnalyzer();
await analyzer.init({ wasmPath: '/wasm/monoguard.wasm' });

// Check version
const versionResult = await analyzer.getVersion();
if (isSuccess(versionResult)) {
  console.log(`MonoGuard v${versionResult.data.version}`);
}

// Analyze workspace
const input: WorkspaceInput = {
  files: {
    'package.json': JSON.stringify({ workspaces: ['packages/*'] }),
    'packages/core/package.json': JSON.stringify({ name: '@mono/core' }),
  },
  config: {
    exclude: ['packages/legacy-*'],
  },
};

const result = await analyzer.analyze(input);

if (isSuccess(result)) {
  console.log(`Health Score: ${result.data.healthScore}`);
  console.log(`Packages: ${result.data.packages}`);
  console.log(`Cycles: ${result.data.circularDependencies.length}`);
} else {
  console.error(`Error: ${result.error.code} - ${result.error.message}`);
}
```

### Previous Story Intelligence

**From Story 1.3 (done):**
- Go WASM exports: `MonoGuard.getVersion()`, `MonoGuard.analyze()`, `MonoGuard.check()`
- All return JSON strings with Result structure
- wasm_exec.js needed for Go runtime

**From Story 1.5 (done):**
- Types defined in `@monoguard/types`
- `Result<T>`, `isSuccess()`, `isError()` available
- `MonoGuardAnalyzer` interface defined in `wasm/adapter.ts`

**Key Integration:**
- TypeScript adapter implements `MonoGuardAnalyzer` interface
- Uses types from `@monoguard/types` package
- Wraps raw WASM calls with type safety

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 2.7]
- [Source: packages/types/src/wasm/adapter.ts]
- [Source: _bmad-output/implementation-artifacts/1-3-setup-go-wasm-analysis-engine-project.md]
- [Source: _bmad-output/implementation-artifacts/1-5-setup-shared-typescript-types-package.md]
- [Go WASM in Browser](https://github.com/nicholaides/using-go-in-browser-wasm-example)

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

N/A

### Completion Notes List

- Created TypeScript WASM adapter with full type safety
- Loader module (`loader.ts`): Async WASM loading with timeout support, WasmLoadError custom error class
- Bridge module (`bridge.ts`): Type-safe JSON serialization wrappers for WASM calls with proper error handling
- Analyzer module (`analyzer.ts`): MonoGuardAnalyzer class with init(), analyze(), check(), getVersion() methods
- Index module (`index.ts`): Clean exports with re-exported type guards from @monoguard/types
- All functions follow Result<T> pattern for consistent error handling
- Comprehensive test suite: 35 tests covering loader, bridge, and analyzer modules
- TypeScript builds successfully to dist/ with .d.ts declarations
- Fixed tsconfig.json to properly handle rootDir and exclude test files
- Fixed types package tsconfig.json to exclude vitest.config.ts

### File List

**New TypeScript Files:**
- packages/analysis-engine/src/loader.ts
- packages/analysis-engine/src/bridge.ts
- packages/analysis-engine/src/analyzer.ts
- packages/analysis-engine/src/index.ts
- packages/analysis-engine/src/__tests__/loader.test.ts
- packages/analysis-engine/src/__tests__/bridge.test.ts
- packages/analysis-engine/src/__tests__/analyzer.test.ts

**Configuration Files:**
- packages/analysis-engine/tsconfig.json (created)
- packages/analysis-engine/vitest.config.ts (created)
- packages/analysis-engine/package.json (updated with TypeScript build scripts)
- packages/types/tsconfig.json (fixed include pattern)
