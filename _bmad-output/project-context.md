---
project_name: 'mono-guard'
user_name: 'Alexyu'
date: '2026-01-12'
sections_completed:
  [
    'technology_stack',
    'language_rules',
    'framework_rules',
    'testing_rules',
    'code_quality_rules',
    'development_workflow_rules',
    'critical_rules',
  ]
existing_patterns_found: 15
status: 'complete'
rule_count: 140
optimized_for_llm: true
---

# Project Context for AI Agents

_This file contains critical rules and patterns that AI agents must follow when implementing code in this project. Focus on unobvious details that agents might otherwise miss._

---

## Technology Stack & Versions

**Core Stack (Target Architecture):**

- TanStack Start 0.34+ (SSG) + React 19.0.0
- Zustand 4.4+ (state management, < 5KB)
- Tailwind CSS 3.3+ (JIT mode) + D3.js 7.x (visualization)
- Go 1.21+ (WASM analysis engine + native CLI)
- Dexie.js 5.x (IndexedDB wrapper for local storage)
- TypeScript 5.9.2

**Development Tools:**

- Nx 21.4.1 (monorepo orchestration)
- pnpm 10.14.0 (package manager)
- Node.js >= 18.0.0
- Vitest (frontend unit tests) + Go testing (backend) + Playwright (E2E)
- ESLint 9.0 + Prettier 3.2.0 + golangci-lint (Go)

**Critical Version Constraints:**

- ⚠️ React 19.0.0 requires TypeScript 5.9+ (older versions lack new type support)
- ⚠️ TanStack Start 0.34+ required for stable SSG support
- ⚠️ Go 1.21+ required for complete WASM functionality
- ⚠️ Zustand 4.4+ required for React 19 compatibility

## Critical Implementation Rules

### Language-Specific Rules

**TypeScript:**

- **Naming:** camelCase (variables/functions), PascalCase (types/interfaces/components)
- **Files:** PascalCase.tsx (React components), camelCase.ts (utilities)
- **Imports:** Use Nx workspace paths (`@monoguard/*`), avoid deep relative paths (>2 levels)
- **Errors:** Use `AnalysisError` class with separated technical + user messages
- **WASM calls:** Always wrap returns with `Result<T>` type
- **Example:**

  ```typescript
  // ✅ CORRECT: Layered error handling
  catch (error) {
    if (error instanceof AnalysisError) {
      toast.error(error.userMessage); // User-friendly
      Sentry.captureException(error, { extra: { technical: error.technicalMessage } });
    }
  }

  // ❌ WRONG: Exposing technical errors to users
  catch (error) {
    toast.error(error.message);
  }
  ```

**Go:**

- **Naming:** PascalCase (exported), camelCase (unexported), snake_case (files)
- **Files:** snake_case.go, \*\_test.go for tests
- **JSON struct tags:** Always camelCase (`json:"healthScore"` NOT `json:"health_score"`)
- **Dates:** ISO 8601 strings, NEVER Unix timestamps
- **WASM functions:** Must return `Result` type with {Data, Error} structure
- **Error codes:** UPPER_SNAKE_CASE (e.g., PARSE_ERROR, CIRCULAR_DETECTED)
- **Example:**

  ```go
  // ✅ CORRECT: Unified Result type + camelCase JSON
  type AnalysisResult struct {
      HealthScore int `json:"healthScore"`
      CreatedAt   string `json:"createdAt"` // ISO 8601
  }

  func AnalyzeWorkspace(jsonData string) string {
      if err != nil {
          return toJSON(NewError("ANALYSIS_FAILED", err.Error()))
      }
      return toJSON(NewSuccess(result))
  }

  // ❌ WRONG: No Result wrapper, snake_case JSON
  type AnalysisResult struct {
      HealthScore int `json:"health_score"` // Should be camelCase
  }
  ```

**Cross-Language Consistency:**

- **JSON format:** camelCase everywhere (TypeScript ↔ Go ↔ JSON)
- **Date format:** ISO 8601 strings (e.g., "2026-01-12T10:30:00Z")
- **Error handling:** Layered approach (technical logs, user-friendly UI messages)

### Framework-Specific Rules

**React 19:**

- Use hooks (avoid class components), `React.memo()` for D3 components to prevent re-renders
- Custom hooks: Must prefix with `use` (e.g., `useWasmLoader`, `useAnalysis`)
- Component structure: One file per component (PascalCase.tsx), tests in `__tests__/`
- Props: TypeScript interfaces named `ComponentNameProps`
- ⚠️ React 19 `use` hook can be used in conditionals (new feature)

**Zustand (State Management):**

- All global state via Zustand stores (`apps/web/app/stores/`)
- Use `devtools` middleware (dev) + `persist` middleware (if data needs persistence)
- Actions: Verb naming (e.g., `startAnalysis`, `clearResult`, not `handleClick`)
- Components: Use selectors to avoid over-rendering

  ```typescript
  // ✅ CORRECT: Selector usage
  const { result, isAnalyzing } = useAnalysisStore((state) => ({
    result: state.result,
    isAnalyzing: state.isAnalyzing,
  }));

  // ❌ WRONG: Subscribe to entire store (causes unnecessary re-renders)
  const store = useAnalysisStore();
  ```

**TanStack Start (SSG):**

- File-based routing in `apps/web/app/routes/` (index.tsx, analysis.$id.tsx, \_\_root.tsx)
- SSG only (no SSR) - all data loaded client-side or at build time
- ❌ Do NOT use `getStaticProps` or `getServerSideProps` (those are Next.js, not TanStack Start)
- All analysis runs in browser via WASM (zero backend architecture)

**D3.js Integration:**

- Must use `useEffect` + `useRef` for D3 initialization
- **Critical:** Always cleanup - remove event listeners in useEffect return function
- Performance rule: SVG rendering (<500 nodes), Canvas rendering (>500 nodes)

  ```typescript
  // ✅ CORRECT: D3 with proper cleanup
  const DependencyGraph = React.memo(({ data }: Props) => {
    const svgRef = useRef<SVGSVGElement>(null);

    useEffect(() => {
      if (!svgRef.current || !data) return;

      const svg = d3.select(svgRef.current);
      // ... D3 rendering logic

      return () => {
        svg.selectAll('*').remove();
        svg.on('zoom', null); // Remove event listeners
      };
    }, [data]);

    return <svg ref={svgRef} />;
  });

  // ❌ WRONG: No cleanup (memory leak)
  useEffect(() => {
    const svg = d3.select(svgRef.current);
    svg.call(d3.zoom().on('zoom', handleZoom));
    // Missing: return cleanup function
  }, [data]);
  ```

**Dexie.js (IndexedDB):**

- Centralized database management in `apps/web/app/lib/persistence.ts`
- Schema versioning: `.version(1).stores({...})`
- Table names: Use plural (e.g., `analyses`, `settings`)
- **Storage rules:**
  - Large analysis results (>100KB): IndexedDB via Dexie.js
  - Small settings (<5KB): Zustand persist (uses localStorage automatically)
  - ❌ NEVER use `localStorage.setItem()` directly for analysis results (performance + size limits)

  ```typescript
  // ✅ CORRECT: Dexie initialization
  class MonoGuardDB extends Dexie {
    analyses!: Table<AnalysisRecord>;
    settings!: Table<SettingRecord>;

    constructor() {
      super('monoguard');
      this.version(1).stores({
        analyses: '++id, timestamp, workspaceName, [workspaceName+timestamp]',
        settings: 'key',
      });
    }
  }

  // ❌ WRONG: Using localStorage for large data
  localStorage.setItem('analysis', JSON.stringify(largeAnalysisResult)); // Violates architecture
  ```

### Testing Rules

**Test Organization:**

- TypeScript: `__tests__/` directory next to source files
  ```
  packages/types/src/analysis/
  ├── index.ts
  └── __tests__/
      └── index.test.ts
  ```
- Go: `*_test.go` in same directory as source
  ```
  packages/analysis-engine/pkg/analyzer/
  ├── workspace.go
  └── workspace_test.go
  ```
- ❌ Do NOT mix patterns (e.g., TypeScript using `*.test.ts` in source dir violates architecture)

**Testing Frameworks:**

- **Vitest** (frontend unit tests) - Config: `vitest.config.ts` per package
- **Go testing** (backend unit tests) - Native `testing` package + optional `testify/assert`
- **Playwright** (E2E tests) - Location: `apps/web-e2e/src/*.spec.ts`

**Test Coverage Targets:**

- Unit Tests: >80% coverage
- Integration Tests: Core WASM Bridge paths must be tested
- E2E Tests: Minimum 3-5 critical user flows
  - Example flows: Upload → Analyze → Visualize → Export

**Critical Test Points (Must Test):**

- WASM Bridge error handling (`Result<T>` error cases, not just success)
- Zustand store state transitions (especially error states)
- D3.js rendering (verify no errors + correct node count)
- IndexedDB save/load operations (data persistence)

**Mock Patterns:**

- **Mock WASM calls:** Return `Result<T>` structure

  ```typescript
  // ✅ CORRECT: Mock with proper Result type
  vi.mock('@/lib/wasmLoader', () => ({
    MonoGuardAnalyzer: {
      analyze: vi.fn().mockResolvedValue({
        data: { healthScore: 85 },
        error: null,
      }),
      analyzeWithError: vi.fn().mockResolvedValue({
        data: null,
        error: { code: 'PARSE_ERROR', message: 'Invalid JSON' },
      }),
    },
  }));

  // ❌ WRONG: Real WASM loading (slow + unstable tests)
  const realResult = await analyzer.analyze(data);
  ```

- **Mock Zustand stores:** Provide fake state + actions

  ```typescript
  const mockStore = {
    result: mockData,
    isAnalyzing: false,
    startAnalysis: vi.fn(),
  };
  vi.mock('@/stores/analysis', () => ({
    useAnalysisStore: () => mockStore,
  }));
  ```

- **Mock IndexedDB:** Use `fake-indexeddb` for in-memory testing
  ```typescript
  import 'fake-indexeddb/auto'; // Auto-mocks IndexedDB
  ```

**Test Boundaries:**

- **Unit Tests:** Single function/component, mock all external dependencies (< 1s execution)
- **Integration Tests:** Multiple modules working together, minimal mocking (5-10s acceptable)
- **E2E Tests:** Full user flow in real browser environment (30s-1min acceptable)

**Go Testing Patterns:**

- Use table-driven tests for multiple scenarios
- Test function naming: `TestFunctionName`

  ```go
  // ✅ CORRECT: Table-driven test
  func TestAnalyzeWorkspace(t *testing.T) {
      tests := []struct {
          name    string
          input   string
          wantErr bool
      }{
          {"valid workspace", validJSON, false},
          {"invalid JSON", "{invalid", true},
          {"empty input", "", true},
      }

      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              result, err := AnalyzeWorkspace(tt.input)
              if (err != nil) != tt.wantErr {
                  t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
              }
          })
      }
  }
  ```

### Development Workflow Rules

**Git Branch Naming:**

- **Main branches:** `main` (production), `develop` (integration)
- **Feature branches:** `feature/ISSUE-123-short-description` (e.g., `feature/GH-45-wasm-error-handling`)
- **Bug fixes:** `fix/ISSUE-123-short-description`
- **Refactors:** `refactor/ISSUE-123-short-description`
- **Hot fixes:** `hotfix/critical-description`
- ❌ Do NOT use: `my-branch`, `test`, `temp`, `wip`

**Commit Message Format:**

```
<type>(<scope>): <subject>

<body (optional)>

<footer (optional)>
```

**Types:**

- `feat`: New feature (e.g., `feat(wasm): add circular dependency detection`)
- `fix`: Bug fix (e.g., `fix(analysis): handle empty workspace.json`)
- `refactor`: Code change that neither fixes bug nor adds feature
- `perf`: Performance improvement
- `test`: Adding/updating tests
- `docs`: Documentation only
- `chore`: Build process, dependencies, tooling

**Scopes:** `wasm`, `ui`, `visualization`, `storage`, `analysis`, `types`

**Examples:**

```
✅ CORRECT:
feat(wasm): add health score calculation
fix(ui): resolve D3 memory leak on unmount
refactor(storage): migrate from localStorage to IndexedDB

❌ WRONG:
Update stuff
Fixed bug
WIP
```

**Pull Request Requirements:**

- **PR Title:** Same format as commit messages (e.g., `feat(wasm): add circular dependency detection`)
- **PR Description Template:**

  ```markdown
  ## Summary

  Brief description of changes (1-3 sentences)

  ## Related Issues

  Closes #123

  ## Changes Made

  - Changed X to Y
  - Added new feature Z
  - Refactored component W

  ## Testing

  - [ ] Unit tests pass (pnpm test)
  - [ ] E2E tests pass (pnpm e2e)
  - [ ] Manual testing completed
  - [ ] No ESLint/Prettier warnings

  ## Screenshots (if UI changes)

  [Add screenshots]

  ## Breaking Changes

  None / [Describe breaking changes]
  ```

**PR Checklist:**

- [ ] Branch is up-to-date with `main`
- [ ] All tests pass locally
- [ ] No console.log/console.error left in code
- [ ] Code follows naming conventions
- [ ] Added/updated tests for new features
- [ ] Documentation updated (if needed)

**CI/CD Workflow (Nx + GitHub Actions):**

- **Build Pipeline:**
  1. Lint: `pnpm nx run-many --target=lint --all`
  2. Test: `pnpm nx run-many --target=test --all`
  3. Build: `pnpm nx run-many --target=build --all`
  4. E2E: `pnpm nx run web-e2e:e2e` (on `main` only)

**🚨 MANDATORY CI Verification Before Story Completion:**

Dev agents MUST verify CI passes BEFORE marking any story as "done":

```bash
# REQUIRED: Run this before marking story as done
pnpm nx affected --target=lint,test,type-check --base=main

# If any E2E tests might be affected:
pnpm nx run web-e2e:e2e

# Go tests (if analysis-engine changes):
cd packages/analysis-engine && make test
```

**CI Verification Checklist (add to every story):**
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- [ ] Go tests pass: `cd packages/analysis-engine && make test`
- [ ] E2E tests pass (if UI affected): `pnpm nx run web-e2e:e2e`

**⚠️ CRITICAL: A story is NOT done if CI fails. Local Go tests (`make test`) alone are INSUFFICIENT.**
- **Deployment Triggers:**
  - `main` branch: Auto-deploy to production (Render)
  - Pull requests: Preview deployments (manual trigger on Render)
- **Cache Strategy:**
  - Nx computation cache enabled (`.nx/cache`)
  - pnpm store cache: `.pnpm-store/`
  - ⚠️ Do NOT commit cache directories

**Monorepo Workflow with Nx:**

```bash
# Lint specific app
pnpm nx lint web

# Test specific package
pnpm nx test types

# Build everything
pnpm nx run-many --target=build --all

# Run affected (changed code only)
pnpm nx affected:test
pnpm nx affected:build
```

**Dependency Management:**

- Add to workspace root: `pnpm add -D <package> -w`
- Add to specific app: `pnpm add <package> --filter @monoguard/web`
- Update all: `pnpm update -r`
- ❌ Never use `npm install` (project uses pnpm)

**Deployment Considerations:**

- **Production Build:**
  - WASM file must be in `public/monoguard.wasm` (statically served)
  - Environment variables: None required (zero backend)
  - Build output: Static site (SSG)
- **Performance Checks:**
  - Lighthouse score > 90 (Performance, Accessibility, Best Practices, SEO)
  - First Contentful Paint < 1.8s
  - Time to Interactive < 3.8s
  - WASM bundle size < 2MB (compressed)

### Critical Don't-Miss Rules

**Anti-Patterns (Forbidden Code):**

**❌ NEVER use localStorage for analysis results**

```typescript
// ❌ WRONG: Will fail with QuotaExceededError on large data
localStorage.setItem('analysis', JSON.stringify(analysisResult));

// ✅ CORRECT: Use IndexedDB via Dexie.js
await db.analyses.add({
  timestamp: new Date().toISOString(),
  result: analysisResult,
});
```

**❌ NEVER use snake_case in JSON**

```go
// ❌ WRONG: Breaks frontend TypeScript
type AnalysisResult struct {
    HealthScore int `json:"health_score"` // WRONG
}

// ✅ CORRECT: Use camelCase for all JSON
type AnalysisResult struct {
    HealthScore int `json:"healthScore"` // CORRECT
}
```

**❌ NEVER forget D3.js cleanup**

```typescript
// ❌ WRONG: Memory leak - event listeners never removed
useEffect(() => {
  const svg = d3.select(svgRef.current);
  svg.call(d3.zoom().on('zoom', handleZoom));
}, [data]);

// ✅ CORRECT: Always cleanup in return function
useEffect(() => {
  const svg = d3.select(svgRef.current);
  const zoom = d3.zoom().on('zoom', handleZoom);
  svg.call(zoom);

  return () => {
    svg.on('.zoom', null); // Remove zoom listener
    svg.selectAll('*').remove(); // Clean DOM
  };
}, [data]);
```

**❌ NEVER use SSR features in TanStack Start**

```typescript
// ❌ WRONG: TanStack Start is SSG-only, not SSR
export function getServerSideProps() { ... } // This is Next.js, not TanStack Start

// ✅ CORRECT: Client-side data loading only
export default function AnalysisPage() {
  const [data, setData] = useState(null);

  useEffect(() => {
    // Load from IndexedDB or trigger analysis
    loadAnalysisData().then(setData);
  }, []);
}
```

**❌ NEVER return raw Go errors to WASM**

```go
// ❌ WRONG: Frontend gets unparseable response
func AnalyzeWorkspace(input string) string {
    data, err := parseJSON(input)
    if err != nil {
        return err.Error() // WRONG: Frontend expects Result type
    }
}

// ✅ CORRECT: Always wrap in Result type
func AnalyzeWorkspace(input string) string {
    data, err := parseJSON(input)
    if err != nil {
        return toJSON(NewError("PARSE_ERROR", err.Error()))
    }
    return toJSON(NewSuccess(data))
}
```

**Edge Cases & Gotchas:**

**⚠️ Empty Workspace Handling:**

```typescript
// Edge case: User uploads empty workspace.json
// ✅ CORRECT: Validate before passing to WASM
if (
  !workspaceData.projects ||
  Object.keys(workspaceData.projects).length === 0
) {
  throw new AnalysisError(
    'EMPTY_WORKSPACE',
    'Workspace contains no projects',
    'Please select a workspace.json with at least one project'
  );
}
```

**⚠️ Circular Dependencies Detection:**

```go
// Edge case: Self-referencing packages (A -> A)
// ✅ CORRECT: Check for self-loops before cycle detection
func detectCycles(graph DependencyGraph) []Cycle {
    // First check for self-loops
    for node, deps := range graph {
        for _, dep := range deps {
            if dep == node {
                return []Cycle{{Nodes: []string{node, node}}}
            }
        }
    }
    // Then run DFS for multi-node cycles
    return detectMultiNodeCycles(graph)
}
```

**⚠️ WASM Memory Limits:**

```go
// Edge case: Very large monorepos (>10,000 files)
// ⚠️ WASM has memory constraints (default 16MB stack)
// ✅ CORRECT: Process in chunks, stream results
func AnalyzeLargeWorkspace(input string) string {
    if estimatedSize(input) > WASM_SAFE_SIZE {
        return toJSON(NewError("WORKSPACE_TOO_LARGE",
            "Please analyze smaller portions of the workspace"))
    }
    // Process normally for smaller workspaces
}
```

**⚠️ Zustand Re-render Performance:**

```typescript
// Edge case: Subscribing to entire store causes re-renders on every state change
// ❌ WRONG: Component re-renders when ANY store property changes
const store = useAnalysisStore();

// ✅ CORRECT: Use selectors to subscribe only to needed properties
const { result, isAnalyzing } = useAnalysisStore((state) => ({
  result: state.result,
  isAnalyzing: state.isAnalyzing,
}));
```

**Security Rules:**

**🔒 Input Validation:**

```typescript
// Security: Always validate file uploads
const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB

function validateWorkspaceFile(file: File) {
  if (file.size > MAX_FILE_SIZE) {
    throw new Error('File too large (max 10MB)');
  }
  if (!file.name.endsWith('.json')) {
    throw new Error('Only JSON files allowed');
  }
  // Validate JSON structure before parsing
  const content = await file.text();
  try {
    const data = JSON.parse(content);
    if (!data.version || !data.projects) {
      throw new Error('Invalid workspace.json structure');
    }
  } catch {
    throw new Error('Invalid JSON format');
  }
}
```

**🔒 XSS Prevention:**

```typescript
// Security: D3.js can introduce XSS if not careful
// ❌ WRONG: Directly inserting user data as HTML
svg.append('text').html(userInput); // XSS vulnerability

// ✅ CORRECT: Use .text() or sanitize HTML
svg.append('text').text(userInput); // Safe - auto-escaped
```

**🔒 Privacy Protection:**

```typescript
// Security: Never send code to external services
// ❌ WRONG: Violates privacy promise (NFR9-NFR10)
await fetch('https://external-api.com/analyze', {
  body: JSON.stringify(codeData), // FORBIDDEN
});

// ✅ CORRECT: All analysis happens locally via WASM
const result = await MonoGuardAnalyzer.analyze(codeData); // Client-side only
```

**Performance Gotchas:**

**⚡ D3.js Rendering Threshold:**

```typescript
// Performance: SVG becomes slow at >500 nodes
// ✅ CORRECT: Switch to Canvas for large graphs
const DependencyGraph = ({ data }: Props) => {
  const nodeCount = data.nodes.length;

  if (nodeCount > 500) {
    return <CanvasGraph data={data} />; // Use Canvas API
  }
  return <D3SVGGraph data={data} />; // Use D3 + SVG
};
```

**⚡ IndexedDB Batch Operations:**

```typescript
// Performance: Don't insert one-by-one in loops
// ❌ WRONG: Slow - creates transaction per insert
for (const item of items) {
  await db.analyses.add(item); // Slow
}

// ✅ CORRECT: Use bulkAdd for batch operations
await db.analyses.bulkAdd(items); // Fast - single transaction
```

**⚡ WASM Initialization:**

```typescript
// Performance: Initialize WASM once, reuse instance
// ❌ WRONG: Loading WASM on every analysis (slow)
async function analyze(data: any) {
  const wasm = await loadWASM(); // Slow initialization
  return wasm.analyze(data);
}

// ✅ CORRECT: Load once, cache instance
let wasmInstance: MonoGuardAnalyzer | null = null;

async function analyze(data: any) {
  if (!wasmInstance) {
    wasmInstance = await loadWASM(); // Initialize once
  }
  return wasmInstance.analyze(data); // Reuse instance
}
```

**⚡ React Rendering Optimization:**

```typescript
// Performance: Prevent unnecessary re-renders of D3 components
// ❌ WRONG: D3 component re-renders on every parent update
function DependencyGraph({ data }: Props) {
  // D3 rendering logic
}

// ✅ CORRECT: Use React.memo to prevent re-renders unless data changes
const DependencyGraph = React.memo(
  ({ data }: Props) => {
    // D3 rendering logic
  },
  (prevProps, nextProps) => {
    return prevProps.data === nextProps.data; // Custom comparison
  }
);
```

**Critical Reminders:**

**🎯 Result<T> Type is MANDATORY for WASM:**

- ALL Go WASM functions MUST return `Result<T>` JSON structure
- Structure: `{ data: T | null, error: { code: string, message: string } | null }`
- Frontend expects this format - raw returns will break TypeScript

**🎯 Date Format Consistency:**

- ALWAYS use ISO 8601: `2026-01-12T10:30:00Z`
- NEVER use Unix timestamps (e.g., `1673520600`)
- NEVER use locale-specific formats (e.g., `1/12/2026`)

**🎯 Zero Backend Architecture:**

- NO server-side code execution
- NO external API calls with user code
- ALL analysis happens client-side via WASM
- This is a core privacy promise (NFR9-NFR10)

**🎯 Nx Monorepo Commands:**

- ALWAYS use `pnpm` (NEVER `npm` or `yarn`)
- ALWAYS use `pnpm nx` for builds/tests/lints
- Use `affected` commands when possible for performance

**🎯 Test Coverage is NOT Optional:**

- WASM bridge MUST have >80% coverage
- Core analysis logic MUST have integration tests
- E2E tests MUST cover upload → analyze → visualize flow

**🎯 CI Must Pass Before Story Completion:**

- NEVER mark a story as "done" if CI is failing
- Running `make test` in analysis-engine is INSUFFICIENT - must run full `pnpm nx affected`
- E2E tests MUST be verified if any UI or routing changes were made
- This is a BLOCKING requirement - no exceptions

---

## Usage Guidelines

**For AI Agents:**

- Read this file before implementing any code
- Follow ALL rules exactly as documented
- When in doubt, prefer the more restrictive option
- Update this file if new patterns emerge

**For Humans:**

- Keep this file lean and focused on agent needs
- Update when technology stack changes
- Review quarterly for outdated rules
- Remove rules that become obvious over time

Last Updated: 2026-01-18
