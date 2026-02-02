# Story 5.3: Implement WASM-Powered In-Browser Analysis

Status: ready-for-dev

## Story

As a **user**,
I want **analysis to run entirely in my browser**,
So that **my code never leaves my machine**.

## Acceptance Criteria

### AC1: WASM Module Loading
**Given** uploaded workspace files are validated
**When** analysis is triggered
**Then**:
- WASM module loads with a progress indicator (percentage or spinner)
- Loading happens once and instance is cached for subsequent analyses
- If WASM fails to load, a clear error message is shown with retry option
- WASM file is loaded from `public/monoguard.wasm`

### AC2: Zero Network Analysis
**Given** the analysis is executing
**When** I check browser DevTools Network tab
**Then**:
- No network requests are made to external servers for analysis
- Only local WASM file loading occurs (from same origin)
- Analysis runs entirely in the WASM module
- Privacy indicator confirms "Analyzed locally"

### AC3: Analysis Progress Display
**Given** the analysis is running
**When** I view the analysis page
**Then**:
- Progress shows current phase: "Loading WASM..." → "Analyzing packages..." → "Calculating health score..."
- Package progress displayed: "Analyzing X/Y packages" (if available from WASM)
- Estimated time remaining or elapsed time shown
- Cancel button available to abort analysis

### AC4: Analysis Timing Requirements
**Given** a workspace with up to 100 packages
**When** analysis runs in the browser
**Then**:
- First feedback (WASM loaded indicator) appears in < 0.5 seconds
- Complete analysis finishes in < 3 seconds for 100 packages
- Results appear progressively (health score first, then details)
- No UI freezing during analysis (responsive interactions)

### AC5: Result Handling
**Given** analysis completes successfully
**When** results are available
**Then**:
- Results are stored in Zustand analysis store
- Automatic navigation to results view (or in-page transition)
- Results include: dependency graph, circular dependencies, version conflicts, health score
- Results match `AnalysisResult` type from `@monoguard/types`

### AC6: Error Handling
**Given** analysis encounters an error
**When** the error occurs
**Then**:
- Error message is user-friendly (not raw Go/WASM error)
- Error uses `AnalysisError` class with technical + user message separation
- Retry button is available
- Previous valid results (if any) are preserved
- Error details available in expandable section for debugging

### AC7: Privacy Confirmation
**Given** analysis has completed
**When** I view the results
**Then**:
- Privacy badge displays "Analyzed locally - no data uploaded"
- Badge is visible on the results page
- Clicking badge shows details: "All analysis ran in your browser via WebAssembly"

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create WASM Loader Service (AC: 1)
  - [ ] Create `apps/web/app/lib/wasmLoader.ts`
  - [ ] Implement singleton WASM loading with caching (load once, reuse)
  - [ ] Add progress callback for loading state
  - [ ] Handle loading failures with retry mechanism
  - [ ] Load from `/monoguard.wasm` (public directory)
  - [ ] Follow project-context.md pattern: cache instance, don't reload

- [ ] Task 2: Create Analysis Execution Service (AC: 2, 4, 5)
  - [ ] Create `apps/web/app/lib/analysisEngine.ts`
  - [ ] Wrap WASM `analyze()` call with TypeScript types
  - [ ] Handle `Result<AnalysisResult>` response (check error field)
  - [ ] Transform WASM output to typed `AnalysisResult`
  - [ ] Ensure all JSON uses camelCase (Go struct tags already configured)

- [ ] Task 3: Create Analysis Progress Component (AC: 3)
  - [ ] Create `apps/web/app/components/analysis/AnalysisProgress.tsx`
  - [ ] Display phase indicators (Loading WASM → Analyzing → Calculating)
  - [ ] Show elapsed time counter
  - [ ] Add cancel button with abort handling
  - [ ] Animate transitions between phases
  - [ ] Use Radix UI progress component or custom Tailwind progress bar

- [ ] Task 4: Create Privacy Badge Component (AC: 7)
  - [ ] Create `apps/web/app/components/common/PrivacyBadge.tsx`
  - [ ] Display lock icon + "Analyzed locally" text
  - [ ] Clickable to show details dialog/tooltip
  - [ ] Reusable across landing page and results page

- [ ] Task 5: Wire Up Analysis Flow in Analyze Route (AC: 3, 4, 5, 6)
  - [ ] Update `apps/web/app/routes/analyze.tsx`
  - [ ] Connect workspace data from upload (Zustand store) to WASM analysis
  - [ ] Handle loading → analyzing → results state machine
  - [ ] Navigate to results on completion
  - [ ] Show error state with retry on failure
  - [ ] Implement cancel functionality

- [ ] Task 6: Create/Update Analysis Store (AC: 5)
  - [ ] Create `apps/web/app/stores/analysis.ts` (Zustand with devtools)
  - [ ] Store: workspaceData, analysisResult, isAnalyzing, error, progress phase
  - [ ] Actions: startAnalysis, setResult, setError, clearResult, cancelAnalysis
  - [ ] Use selectors for component subscriptions (prevent re-render)

- [ ] Task 7: Write Unit Tests (AC: all)
  - [ ] Test WASM loader singleton behavior
  - [ ] Test WASM loader error handling and retry
  - [ ] Test analysis execution with mock WASM (Result<T> success and error)
  - [ ] Test progress component state transitions
  - [ ] Test privacy badge rendering and interaction
  - [ ] Test analysis store state management
  - [ ] Test error handling with AnalysisError class
  - [ ] Mock WASM calls with `vi.mock` returning `Result<T>` structure
  - [ ] Target: >80% coverage

- [ ] Task 8: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- WASM loader: `apps/web/app/lib/wasmLoader.ts` (new)
- Analysis engine service: `apps/web/app/lib/analysisEngine.ts` (new)
- Analysis progress component: `apps/web/app/components/analysis/AnalysisProgress.tsx` (new)
- Analysis store: `apps/web/app/stores/analysis.ts` (new)
- Analyze route: `apps/web/app/routes/analyze.tsx` (existing, refactor)

**WASM Loading Pattern (from project-context.md):**
```typescript
// CORRECT: Load once, cache instance
let wasmInstance: MonoGuardAnalyzer | null = null;

async function getAnalyzer(): Promise<MonoGuardAnalyzer> {
  if (!wasmInstance) {
    wasmInstance = await loadWASM(); // Initialize once
  }
  return wasmInstance; // Reuse instance
}
```

**Result<T> Type Handling:**
```typescript
// WASM returns Result<T> structure
const result = await analyzer.analyze(workspaceData);
if (result.error) {
  // Handle error: { code: string, message: string }
  throw new AnalysisError(result.error.code, result.error.message, userFriendlyMessage);
}
// Use result.data (AnalysisResult)
```

**WASM File Location:**
- Build output: `packages/analysis-engine/` → `monoguard.wasm`
- Serve from: `apps/web/public/monoguard.wasm`
- Note: For development, WASM may need to be copied/symlinked from build output

**Error Handling Pattern (from project-context.md):**
```typescript
catch (error) {
  if (error instanceof AnalysisError) {
    toast.error(error.userMessage); // User-friendly
    // Technical details in expandable section
  }
}
```

### Zustand Store Pattern (from project-context.md)

```typescript
// Analysis store with devtools
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

interface AnalysisState {
  workspaceData: WorkspaceData | null;
  result: AnalysisResult | null;
  isAnalyzing: boolean;
  error: AnalysisError | null;
  progress: 'idle' | 'loading-wasm' | 'analyzing' | 'calculating' | 'complete' | 'error';

  // Actions (verb naming)
  startAnalysis: (data: WorkspaceData) => void;
  setResult: (result: AnalysisResult) => void;
  setError: (error: AnalysisError) => void;
  clearResult: () => void;
  cancelAnalysis: () => void;
}

// Components use selectors:
const { result, isAnalyzing } = useAnalysisStore(state => ({
  result: state.result,
  isAnalyzing: state.isAnalyzing,
}));
```

### Previous Story Intelligence

**From Story 5.2:**
- Workspace data arrives from file upload pipeline
- Data is stored in Zustand store for analysis consumption
- FileProcessing outputs `WorkspaceData` type

**From Epic 2 (WASM Adapter):**
- Story 2.7 created TypeScript WASM adapter
- `analyzer.init()` → `analyzer.analyze(workspaceData)` → `Result<AnalysisResult>`
- JSON serialization uses camelCase throughout
- WASM file compiled by Go build (`GOOS=js GOARCH=wasm`)

### Critical Don't-Miss Rules

1. **NEVER send code to external services** - All analysis via WASM locally
2. **Wrap WASM returns with Result<T>** - Always check error field
3. **Use AnalysisError class** - Separate technical + user messages
4. **Cache WASM instance** - Don't reload on every analysis
5. **Zustand selectors** - Prevent over-rendering
6. **D3.js cleanup** not directly applicable here but analysis results feed into Epic 4 visualization

### Testing Requirements

- **Mock WASM calls** with `vi.mock` returning `Result<T>` structure
- **Don't load real WASM in tests** - Mock the loader
- Test both success and error paths
- Test progress state machine transitions
- Test store state management with actions

### References

- [FR30: WASM browser execution] `_bmad-output/planning-artifacts/epics.md`
- [NFR9: Zero code upload] `_bmad-output/planning-artifacts/epics.md`
- [NFR1: Performance targets] `_bmad-output/planning-artifacts/epics.md`
- [Project Context: WASM Patterns] `_bmad-output/project-context.md#wasm-initialization`
- [Project Context: Error Handling] `_bmad-output/project-context.md#language-specific-rules`
- [Story 2.7: WASM Adapter] `_bmad-output/implementation-artifacts/2-7-create-typescript-wasm-adapter.md`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
