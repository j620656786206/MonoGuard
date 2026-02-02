# Story 5.4: Build Analysis Results Dashboard

Status: ready-for-dev

## Story

As a **user**,
I want **a clear dashboard showing my analysis results**,
So that **I can quickly understand my monorepo's health**.

## Acceptance Criteria

### AC1: Progressive Disclosure (L1/L2/L3)
**Given** completed analysis results
**When** I view the results dashboard
**Then** information is presented in three progressive levels:
- **L1 (Overview)**: Large health score number with color coding, problem count badges (circular deps, version conflicts), "Analyzed locally" privacy badge
- **L2 (Problem List)**: Expandable problem categories with counts, severity indicators, package names involved
- **L3 (Detailed View)**: Full dependency graph visualization (from Epic 4), detailed fix recommendations, diagnostic reports

### AC2: Health Score Display
**Given** analysis results with a health score
**When** the dashboard renders
**Then**:
- Health score displayed prominently (large font, centered or top-left)
- Color gradient matches thresholds: Excellent (85-100 green), Good (70-84), Fair (50-69 yellow), Poor (30-49 orange), Critical (0-29 red)
- Score counting animation from 0 to actual value (on first load)
- Score breakdown is available (click to expand factors: circular deps weight, version conflicts weight, etc.)

### AC3: Problem Summary Cards
**Given** analysis results
**When** the dashboard renders
**Then** I see summary cards for:
- Circular Dependencies: count + severity badge
- Version Conflicts: count + severity badge
- Total Packages Analyzed: count
- Analysis Duration: time taken
- Each card is clickable to jump to detailed section (L2/L3)

### AC4: Dashboard Layout
**Given** the results dashboard
**When** viewed on desktop
**Then**:
- Two-column layout: main content area (graph + details) + side panel (fix suggestions - Story 5.5)
- Top section: Health score + summary cards
- Middle section: Dependency graph visualization (DependencyGraphViz from Epic 4)
- Bottom section: Detailed lists (circular deps, version conflicts)
- Responsive: stacks to single column on tablet/mobile

### AC5: Navigation Between Levels
**Given** the progressive disclosure dashboard
**When** I click on summary items
**Then**:
- Clicking health score → expands score breakdown (L1→L2)
- Clicking problem count badge → scrolls to detailed list (L1→L2)
- Clicking specific problem → opens detailed view with graph highlighting (L2→L3)
- Breadcrumb or back button to return to higher level
- Transitions are smooth with skeleton loading states

### AC6: Integration with Existing Visualization
**Given** the results dashboard
**When** the dependency graph section renders
**Then**:
- Uses `DependencyGraphViz` component from Epic 4
- All Epic 4 features work: zoom, pan, hover, expand/collapse, hybrid rendering
- Circular dependencies are highlighted in the graph
- Clicking a node in the graph updates the side panel details

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Create Dashboard Layout Component (AC: 4)
  - [ ] Create `apps/web/app/components/dashboard/DashboardLayout.tsx`
  - [ ] Implement two-column layout (main + sidebar)
  - [ ] Make sidebar collapsible and responsive
  - [ ] Add top section for health score and summary cards

- [ ] Task 2: Create Health Score Hero Component (AC: 2)
  - [ ] Create `apps/web/app/components/dashboard/HealthScoreHero.tsx`
  - [ ] Implement large score display with color gradient
  - [ ] Add counting animation (0 → actual score on mount)
  - [ ] Add expandable score breakdown panel
  - [ ] Use health score thresholds from project-context.md

- [ ] Task 3: Create Summary Cards Component (AC: 3)
  - [ ] Create `apps/web/app/components/dashboard/SummaryCards.tsx`
  - [ ] Card: Circular Dependencies (count + severity)
  - [ ] Card: Version Conflicts (count + severity)
  - [ ] Card: Total Packages
  - [ ] Card: Analysis Duration
  - [ ] Each card clickable with hover states

- [ ] Task 4: Create Problem List Section (AC: 1, 5)
  - [ ] Create `apps/web/app/components/dashboard/ProblemList.tsx`
  - [ ] List circular dependencies with expandable details
  - [ ] List version conflicts with package info
  - [ ] Add severity indicators and sorting (by severity, by impact)
  - [ ] Clicking items triggers graph highlighting via store

- [ ] Task 5: Refactor Results Route (AC: 1, 4, 6)
  - [ ] Refactor `apps/web/app/routes/results.tsx` to use new dashboard components
  - [ ] Connect to analysis store for result data
  - [ ] Integrate `DependencyGraphViz` from Epic 4 components
  - [ ] Add progressive disclosure navigation (L1 → L2 → L3)
  - [ ] Add skeleton loading states during transitions

- [ ] Task 6: Wire Graph Interaction to Dashboard (AC: 6)
  - [ ] Connect node selection in graph to sidebar/detail updates
  - [ ] Clicking a circular dep in the list → highlight in graph
  - [ ] Clicking a node in graph → show details in sidebar
  - [ ] Use Zustand store for cross-component communication

- [ ] Task 7: Write Unit Tests (AC: all)
  - [ ] Test health score rendering and color coding
  - [ ] Test counting animation
  - [ ] Test summary cards with various data states
  - [ ] Test problem list rendering and sorting
  - [ ] Test progressive disclosure navigation
  - [ ] Test responsive layout behavior
  - [ ] Test integration with DependencyGraphViz (mock graph component)
  - [ ] Target: >80% coverage

- [ ] Task 8: Verify CI passes (AC-CI)
  - [ ] All lint, test, type-check targets pass

## Dev Notes

### Architecture Patterns & Constraints

**File Locations:**
- Results route: `apps/web/app/routes/results.tsx` (existing, major refactor)
- Dashboard components: `apps/web/app/components/dashboard/` (new directory)
- Existing analysis components: `apps/web/app/components/analysis/` (reuse)
- Graph visualization: `apps/web/app/components/visualization/DependencyGraph/` (from Epic 4)
- Analysis store: `apps/web/app/stores/analysis.ts` (from Story 5.3)

**Existing Components to Integrate:**
- `AnalysisResults.tsx` (21.8 KB) - Has tabbed interface with overview, deps, architecture, etc. Evaluate reuse vs. rebuild
- `HealthScoreDisplay.tsx` - Existing health score display, enhance with animation
- `CircularDependencyViz.tsx` - Existing circular dep visualization
- `VersionConflictTable.tsx` - Existing version conflict display
- `DependencyGraphViz` - Full Epic 4 graph (index.tsx in visualization/)

**Layout Pattern:**
```
┌──────────────────────────────────────────────┐
│  Nav Bar                                      │
├──────────────────────────────────────────────┤
│  Health Score (L1)  │  Summary Cards          │
├─────────────────────┼────────────────────────┤
│                     │                        │
│  Dependency Graph   │  Side Panel            │
│  (Epic 4 Viz)       │  (Fix Suggestions 5.5) │
│                     │                        │
├─────────────────────┴────────────────────────┤
│  Problem Details (L2/L3)                      │
└──────────────────────────────────────────────┘
```

**Health Score Color Thresholds:**
- Excellent: 85-100 → `text-green-600` / `bg-green-50`
- Good: 70-84 → `text-emerald-600` / `bg-emerald-50`
- Fair: 50-69 → `text-yellow-600` / `bg-yellow-50`
- Poor: 30-49 → `text-orange-600` / `bg-orange-50`
- Critical: 0-29 → `text-red-600` / `bg-red-50`

### UX Design Requirements

- Progressive disclosure: L1 → L2 → L3 (from UX spec)
- Desktop-first, tablet-friendly
- Skeleton loading during transitions (not spinners)
- Animation: score counting animation on first load
- Click-to-highlight: clicking problems → highlights in graph

### Previous Story Intelligence

**From Story 5.3:**
- Analysis results stored in Zustand analysis store
- `AnalysisResult` type includes: graph, cycles, conflicts, healthScore
- Privacy badge component created (reuse here)

**From Epic 4 (Visualization):**
- `DependencyGraphViz` accepts: data, circularDependencies, selectedNodeId, onNodeSelect, onNodeHover
- All zoom, pan, expand/collapse, hover features built in
- Hybrid SVG/Canvas rendering handles large graphs
- Component uses React.memo for performance

### Testing Requirements

- Mock analysis store with test data
- Mock DependencyGraphViz component for dashboard tests
- Test color coding matches health score thresholds
- Test responsive layout with container queries or viewport mocks
- Test click interactions between components

### References

- [UX Spec: Progressive Disclosure] `_bmad-output/planning-artifacts/ux-design-specification.md`
- [UX Spec: Health Score Visual] `_bmad-output/planning-artifacts/ux-design-specification.md`
- [FR31: Fix suggestions panel alongside graph] `_bmad-output/planning-artifacts/epics.md`
- [Epic 4 Visualization] `_bmad-output/implementation-artifacts/4-1-*.md` through `4-9-*.md`
- [Project Context: Zustand Selectors] `_bmad-output/project-context.md#zustand-state-management`

## Dev Agent Record

### Agent Model Used

### Debug Log References

### Completion Notes List

### File List
