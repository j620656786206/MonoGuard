# Story 5.5: Create Fix Suggestions Side Panel

Status: ready-for-dev

## Story

As a **user**,
I want **a fix suggestions panel alongside the dependency graph**,
So that **I can see problems and solutions together**.

## Acceptance Criteria

### AC1: Side Panel Layout
**Given** analysis results with issues detected
**When** I view the results page
**Then** I see:
- Main area: Dependency graph visualization (left/center)
- Side panel: Fix suggestions list (right)
- Panel has a clear header "Fix Suggestions" with issue count
- Panel is visually distinct from the graph area

### AC2: Panel Interactivity
**Given** the fix suggestions side panel
**When** I interact with it
**Then**:
- Panel is collapsible (toggle button to hide/show)
- Panel is resizable (drag handle on the divider)
- Panel width is persisted across sessions (Zustand persist)
- Minimum width: 280px, maximum width: 50% of viewport

### AC3: Suggestion List Content
**Given** analysis results with circular dependencies and version conflicts
**When** the side panel renders
**Then** each suggestion shows:
- Issue title (e.g., "Circular: @app/a → @app/b → @app/a")
- Severity badge (Critical / High / Medium / Low)
- Quick description (1-line summary)
- Best fix strategy recommendation
- Expand arrow for detailed view
- Items sorted by severity (highest first) then by impact

### AC4: Suggestion-Graph Interaction
**Given** the side panel with suggestions and the graph
**When** I click on a suggestion
**Then**:
- Related nodes in the graph are highlighted
- Related edges glow/pulse to show the dependency path
- Graph viewport auto-pans to center on the affected nodes
- Other nodes are dimmed
- Clicking away or another suggestion clears the highlight

### AC5: Expanded Suggestion Detail
**Given** a fix suggestion in the panel
**When** I expand it
**Then** I see:
- Full cycle path visualization (text or mini diagram)
- Up to 3 fix strategy options with suitability scores
- Complexity score (1-10) with human label
- Impact assessment (packages affected, risk level)
- "View full diagnostic report" link (to DiagnosticReportModal from Epic 4)
- Each strategy has a brief description

### AC6: Responsive Behavior
**Given** the results page with side panel
**When** viewed on tablet
**Then**:
- Panel collapses to bottom sheet or overlay mode
- Toggle button clearly accessible
- Graph takes full width when panel is collapsed

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Side Panel Container (AC: 1, 2)
  - [ ] Create `apps/web/app/components/dashboard/SidePanel.tsx`
  - [ ] Implement collapsible panel with toggle button
  - [ ] Add resize handle (drag to resize)
  - [ ] Persist panel width in settings store (Zustand persist)
  - [ ] Set min/max width constraints

- [ ] Task 2: Create Suggestion List Component (AC: 3)
  - [ ] Create `apps/web/app/components/dashboard/FixSuggestionList.tsx`
  - [ ] Render list of suggestions from analysis results
  - [ ] Each item: title, severity badge, description, strategy name
  - [ ] Implement sorting by severity then impact
  - [ ] Add expand/collapse toggle per item

- [ ] Task 3: Create Suggestion Detail Component (AC: 5)
  - [ ] Create `apps/web/app/components/dashboard/FixSuggestionDetail.tsx`
  - [ ] Show full cycle path (text representation)
  - [ ] List fix strategies with suitability scores
  - [ ] Show complexity and impact scores
  - [ ] Link to DiagnosticReportModal (from `apps/web/app/components/diagnostics/`)

- [ ] Task 4: Implement Graph-Panel Interaction (AC: 4)
  - [ ] Use Zustand store for selected suggestion state
  - [ ] When suggestion clicked → update selectedNodeId and highlight cycle in graph
  - [ ] Auto-pan graph viewport to affected nodes
  - [ ] Dim non-related nodes
  - [ ] Clear highlight on deselect or click-away

- [ ] Task 5: Add Panel Width Persistence (AC: 2)
  - [ ] Add `sidePanelWidth` and `sidePanelCollapsed` to settings store
  - [ ] Persist via Zustand persist middleware
  - [ ] Read on mount, update on resize

- [ ] Task 6: Implement Responsive Behavior (AC: 6)
  - [ ] At tablet breakpoint (< 1024px): panel becomes overlay/bottom sheet
  - [ ] Toggle button visible in both modes
  - [ ] Graph takes full width when panel collapsed

- [ ] Task 7: Integrate with Dashboard Layout (AC: 1)
  - [ ] Wire SidePanel into `DashboardLayout.tsx` from Story 5.4
  - [ ] Pass analysis results to suggestion list
  - [ ] Connect node selection between graph and panel

- [ ] Task 8: Write Unit Tests (AC: all)
  - [ ] Test panel collapse/expand toggle
  - [ ] Test suggestion list rendering and sorting
  - [ ] Test suggestion detail expansion
  - [ ] Test panel width persistence
  - [ ] Test graph interaction (mock DependencyGraphViz)
  - [ ] Test responsive behavior
  - [ ] Target: >80% coverage

- [ ] Task 9: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Side panel: `apps/web/app/components/dashboard/SidePanel.tsx` (new)
- Suggestion list: `apps/web/app/components/dashboard/FixSuggestionList.tsx` (new)
- Suggestion detail: `apps/web/app/components/dashboard/FixSuggestionDetail.tsx` (new)
- Settings store: `apps/web/app/stores/settings.ts` (extend with panel width)
- Diagnostic modal: `apps/web/app/components/diagnostics/DiagnosticReportModal.tsx` (existing)

**Data Source:**
- Fix suggestions come from Epic 3 analysis engine output
- `AnalysisResult.cycles` → circular dependencies with fix strategies
- `AnalysisResult.conflicts` → version conflicts
- Each `CircularDependency` includes `fixStrategies`, `complexityScore`, `impactAssessment`

**Graph Interaction Pattern:**
```typescript
// In dashboard, connect suggestion selection to graph highlighting
const handleSuggestionClick = (suggestion: FixSuggestion) => {
  // Update store with selected cycle's node IDs
  useAnalysisStore.setState({
    selectedCycleNodes: suggestion.affectedPackages,
    highlightedEdges: suggestion.cyclePath,
  });
};

// DependencyGraphViz reads from store to apply highlighting
```

**Resize Implementation:**
- Use CSS resize or custom drag handler
- Store width in Zustand settings store with persist
- `useCallback` for resize handler to avoid re-renders

### Previous Story Intelligence

**From Story 5.4:**
- Dashboard layout provides the two-column structure
- Side panel is the right column of the dashboard
- Health score and summary cards in the top section

**From Epic 3 (Fix Suggestions Engine):**
- Stories 3.1-3.8 built the fix suggestion system
- Three strategies: Extract Shared Module, Dependency Injection, Module Boundary Refactoring
- Each has suitability score, complexity score, impact assessment
- Before/after explanations available

**From Epic 4 (Visualization):**
- DependencyGraphViz supports: selectedNodeId, onNodeSelect, circularDependencies props
- Cycle highlighting already built into graph (red nodes/edges)
- Graph can dim non-selected elements

### Testing Requirements

- Mock analysis store with fixture data including fix suggestions
- Test panel resize drag interaction
- Test suggestion sorting logic
- Test graph interaction via store state changes
- Use `@testing-library/user-event` for click and resize interactions

### References

- [FR31: Fix suggestions panel alongside graph] `_bmad-output/planning-artifacts/epics.md`
- [UX Spec: Side panel + main view] `_bmad-output/planning-artifacts/ux-design-specification.md`
- [Epic 3: Fix Suggestions] `_bmad-output/planning-artifacts/epics.md#epic-3-stories`
- [Story 4.2: Cycle Highlighting] `_bmad-output/implementation-artifacts/4-2-*.md`
- [Project Context: Zustand Persist] `_bmad-output/project-context.md#zustand-state-management`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
