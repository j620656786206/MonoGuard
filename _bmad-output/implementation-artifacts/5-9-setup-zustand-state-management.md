# Story 5.9: Setup Zustand State Management

Status: ready-for-dev

## Story

As a **developer**,
I want **Zustand configured for application state management**,
So that **state is predictable and debuggable**.

## Acceptance Criteria

### AC1: Analysis Store
**Given** the web application needs to manage analysis state
**When** Zustand analysis store is configured
**Then**:
- Analysis store manages: workspaceData, analysisResult, isAnalyzing, error, progress
- Actions use verb naming: `startAnalysis`, `setResult`, `setError`, `clearResult`, `cancelAnalysis`
- State transitions are predictable (idle → loading → analyzing → complete/error)
- DevTools integration shows state changes in browser extension

### AC2: Settings Store (Enhanced)
**Given** the existing settings store
**When** enhanced for full app needs
**Then** it manages:
- Theme preference: 'light' | 'dark' | 'system'
- Visualization mode: 'auto' | 'force-svg' | 'force-canvas'
- Side panel width and collapsed state
- Recent commands for command palette
- All settings persist across sessions via `persist` middleware

### AC3: UI Store
**Given** the web application has various UI states
**When** Zustand UI store is configured
**Then** it manages:
- Selected node ID in the dependency graph
- Highlighted cycle path
- Active panel/tab in dashboard
- Command palette open/closed state
- Modal states (export dialog, diagnostic report)
- UI state does NOT persist (resets on reload)

### AC4: Store Performance
**Given** components use Zustand stores
**When** state updates occur
**Then**:
- Components use selectors to subscribe only to needed properties
- State updates are < 16ms (maintain 60fps)
- No unnecessary re-renders (verified with React DevTools profiler)
- Store subscriptions use shallow equality checks where appropriate

### AC5: DevTools Integration
**Given** the development environment
**When** Zustand stores are active
**Then**:
- All stores visible in Redux DevTools browser extension
- Each store has a descriptive name (e.g., "analysis", "settings", "ui")
- State changes are logged with action names
- Time-travel debugging works for state inspection

### AC6: Store Testing Support
**Given** Zustand stores are implemented
**When** writing tests
**Then**:
- Stores can be easily mocked in component tests
- Store state can be set directly in tests via `setState()`
- Store reset function available for test cleanup
- No shared state between test cases

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Analysis Store (AC: 1, 5)
  - [ ] Create `apps/web/app/stores/analysis.ts`
  - [ ] Define `AnalysisState` interface with typed state and actions
  - [ ] Implement with `create()` + `devtools()` middleware
  - [ ] State: workspaceData, result, isAnalyzing, error, progress phase
  - [ ] Actions: startAnalysis, setResult, setError, clearResult, cancelAnalysis, setWorkspaceData
  - [ ] DO NOT use persist for analysis state (too large for localStorage)

- [ ] Task 2: Enhance Settings Store (AC: 2, 5)
  - [ ] Update `apps/web/app/stores/settings.ts`
  - [ ] Verify/add: theme, visualizationMode, sidePanelWidth, sidePanelCollapsed, recentCommands
  - [ ] All settings use `persist` middleware (already configured)
  - [ ] Add `resetSettings` action for defaults restoration
  - [ ] Store key remains `monoguard-settings`

- [ ] Task 3: Create UI Store (AC: 3, 5)
  - [ ] Create `apps/web/app/stores/ui.ts`
  - [ ] Define `UIState` interface
  - [ ] State: selectedNodeId, highlightedCyclePath, activeTab, commandPaletteOpen, exportDialogOpen, diagnosticReportOpen
  - [ ] Actions: selectNode, highlightCycle, clearHighlight, setActiveTab, toggleCommandPalette, openExportDialog, closeDiagnosticReport
  - [ ] Use `devtools()` middleware only (NO persist - UI resets on reload)

- [ ] Task 4: Create Store Barrel Export (AC: 1, 2, 3)
  - [ ] Create `apps/web/app/stores/index.ts`
  - [ ] Export all stores from single entry point
  - [ ] Export store types for component usage

- [ ] Task 5: Add Selector Patterns (AC: 4)
  - [ ] Create typed selector helpers or document patterns
  - [ ] Example: `useAnalysisStore(state => state.result)` - single value
  - [ ] Example: `useAnalysisStore(state => ({ result: state.result, isAnalyzing: state.isAnalyzing }))` - multiple values
  - [ ] Verify no components subscribe to entire store (code review)

- [ ] Task 6: Migrate Existing State (AC: 1, 3)
  - [ ] Move any local state in existing components to appropriate stores
  - [ ] Migrate `HealthScoreContext` (`apps/web/app/contexts/HealthScoreContext.tsx`) to analysis store
  - [ ] Evaluate if React Query state can complement Zustand (for future API calls)
  - [ ] Remove unused context providers after migration

- [ ] Task 7: Write Unit Tests (AC: 1, 2, 3, 6)
  - [ ] Test analysis store state transitions (idle → loading → analyzing → complete)
  - [ ] Test analysis store error handling (setError, clearResult)
  - [ ] Test settings store persistence (mock localStorage)
  - [ ] Test UI store state changes (selectNode, highlightCycle)
  - [ ] Test store reset functions
  - [ ] Test DevTools integration (store names visible)
  - [ ] Verify no state leaks between test cases
  - [ ] Target: >80% coverage

- [ ] Task 8: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Analysis store: `apps/web/app/stores/analysis.ts` (new)
- Settings store: `apps/web/app/stores/settings.ts` (existing, enhance)
- UI store: `apps/web/app/stores/ui.ts` (new)
- Barrel export: `apps/web/app/stores/index.ts` (new)
- Context to migrate: `apps/web/app/contexts/HealthScoreContext.tsx`

**Store Architecture:**
```
stores/
├── index.ts          # Barrel export
├── analysis.ts       # Analysis state (devtools only, no persist)
├── settings.ts       # User preferences (devtools + persist)
└── ui.ts             # Transient UI state (devtools only, no persist)
```

**Storage Rules (from project-context.md):**
- Large analysis results (>100KB): DO NOT persist to localStorage
- Small settings (<5KB): Zustand persist (uses localStorage automatically)
- Future: Large results should go to IndexedDB via Dexie.js (Epic 7)

**Analysis Store Pattern:**
```typescript
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import type { WorkspaceData, AnalysisResult } from '@monoguard/types';

type AnalysisProgress = 'idle' | 'loading-wasm' | 'analyzing' | 'calculating' | 'complete' | 'error';

interface AnalysisState {
  // State
  workspaceData: WorkspaceData | null;
  result: AnalysisResult | null;
  isAnalyzing: boolean;
  error: { code: string; userMessage: string; technicalMessage: string } | null;
  progress: AnalysisProgress;

  // Actions (verb naming - mandatory)
  setWorkspaceData: (data: WorkspaceData) => void;
  startAnalysis: () => void;
  setResult: (result: AnalysisResult) => void;
  setError: (error: AnalysisState['error']) => void;
  clearResult: () => void;
  cancelAnalysis: () => void;
  setProgress: (phase: AnalysisProgress) => void;
}

export const useAnalysisStore = create<AnalysisState>()(
  devtools(
    (set) => ({
      workspaceData: null,
      result: null,
      isAnalyzing: false,
      error: null,
      progress: 'idle',

      setWorkspaceData: (data) => set({ workspaceData: data }),
      startAnalysis: () => set({ isAnalyzing: true, error: null, progress: 'loading-wasm' }),
      setResult: (result) => set({ result, isAnalyzing: false, progress: 'complete' }),
      setError: (error) => set({ error, isAnalyzing: false, progress: 'error' }),
      clearResult: () => set({ result: null, error: null, progress: 'idle' }),
      cancelAnalysis: () => set({ isAnalyzing: false, progress: 'idle' }),
      setProgress: (phase) => set({ progress: phase }),
    }),
    { name: 'analysis' }
  )
);
```

**Settings Store Enhancement:**
```typescript
interface SettingsState {
  // Existing
  theme: 'light' | 'dark' | 'system';
  visualizationMode: RenderModePreference;

  // New additions
  sidePanelWidth: number;
  sidePanelCollapsed: boolean;
  recentCommands: string[];

  // Actions
  setTheme: (theme: SettingsState['theme']) => void;
  setVisualizationMode: (mode: RenderModePreference) => void;
  setSidePanelWidth: (width: number) => void;
  toggleSidePanel: () => void;
  addRecentCommand: (commandId: string) => void;
  resetSettings: () => void;
}
```

**UI Store Pattern:**
```typescript
interface UIState {
  selectedNodeId: string | null;
  highlightedCyclePath: string[] | null;
  activeTab: string;
  commandPaletteOpen: boolean;
  exportDialogOpen: boolean;
  diagnosticReportOpen: boolean;

  selectNode: (nodeId: string | null) => void;
  highlightCycle: (path: string[]) => void;
  clearHighlight: () => void;
  setActiveTab: (tab: string) => void;
  toggleCommandPalette: () => void;
  openExportDialog: () => void;
  closeExportDialog: () => void;
}
```

**Selector Pattern (MANDATORY - from project-context.md):**
```typescript
// CORRECT: Use selectors
const { result, isAnalyzing } = useAnalysisStore(state => ({
  result: state.result,
  isAnalyzing: state.isAnalyzing,
}));

// WRONG: Subscribe to entire store
const store = useAnalysisStore(); // Causes unnecessary re-renders
```

### Previous Story Intelligence

**From Story 4.9:**
- Settings store already has Zustand + devtools + persist pattern
- Store key: `monoguard-settings`
- Pattern for adding new fields to existing store is established

**Migration Notes:**
- `HealthScoreContext` can be replaced by analysis store's `result.healthScore`
- Components using `useHealthScore()` hook should migrate to `useAnalysisStore(state => state.result?.healthScore)`

### Testing Patterns

**Store Testing (from Story 4.9):**
```typescript
import { useSettingsStore } from '@/stores/settings';

beforeEach(() => {
  // Reset store state before each test
  useSettingsStore.setState({ theme: 'system', visualizationMode: 'auto' });
});

it('should update theme', () => {
  useSettingsStore.getState().setTheme('dark');
  expect(useSettingsStore.getState().theme).toBe('dark');
});
```

### References

- [Project Context: Zustand Rules] `_bmad-output/project-context.md#zustand-state-management`
- [Project Context: Storage Rules] `_bmad-output/project-context.md#dexiejs-indexeddb`
- [Story 4.9: Settings Store] `_bmad-output/implementation-artifacts/4-9-*.md`
- [Architecture: Zustand < 5KB] `_bmad-output/planning-artifacts/architecture.md`
- [NFR10: Browser data in IndexedDB] `_bmad-output/planning-artifacts/epics.md`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
