# Story 4.2: Highlight Circular Dependencies in Graph

Status: ready-for-dev

## Story

As a **user**,
I want **circular dependencies to be visually highlighted in the graph**,
So that **I can immediately identify problematic relationships**.

## Acceptance Criteria

### AC1: Visual Highlighting of Cycle Nodes
**Given** a dependency graph with circular dependencies
**When** the graph renders
**Then** nodes in cycles have:
- Red border/glow effect distinguishing them from normal nodes
- Visual indicator (icon or badge) showing they're part of a cycle
- Consistent styling that stands out clearly

### AC2: Visual Highlighting of Cycle Edges
**Given** a dependency graph with circular dependencies
**When** the graph renders
**Then** edges forming cycles are:
- Colored red (distinct from normal gray edges)
- Thicker than normal edges for visibility
- Clearly directional (arrows visible)

### AC3: Animated Cycle Paths
**Given** highlighted circular dependencies
**When** viewing the graph
**Then** cycle paths have animation:
- Pulsing or flowing animation effect
- Animation draws attention without being distracting
- Animation can be toggled off in settings (optional)

### AC4: Color Legend
**Given** a graph with circular dependencies
**When** viewing the visualization
**Then** a legend is displayed showing:
- Normal node color and meaning
- Cycle node color (red) and meaning
- Normal edge color and meaning
- Cycle edge color (red) and meaning

### AC5: Click-to-Highlight Cycle
**Given** multiple circular dependencies in the graph
**When** I click on a node or edge that is part of a cycle
**Then**:
- Only that specific cycle's path is highlighted
- Other cycles remain in their default highlighted state (dimmed red)
- The selected cycle is emphasized (brighter/thicker)

### AC6: Dim Non-Cycle Elements on Selection
**Given** a cycle is selected (clicked)
**When** viewing the graph
**Then**:
- Non-cycle nodes and edges are dimmed (reduced opacity)
- The selected cycle path stands out clearly
- Clicking elsewhere or pressing Escape deselects

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [ ] `pnpm nx affected --target=lint --base=main` passes
- [ ] `pnpm nx affected --target=test --base=main` passes
- [ ] `pnpm nx affected --target=type-check --base=main` passes
- [ ] `cd packages/analysis-engine && make test` passes (if Go changes)
- [ ] GitHub Actions CI workflow shows GREEN status
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [ ] Task 1: Extend DependencyGraph component props (AC: 1, 2)
  - [ ] Add `circularDependencies?: CircularDependencyInfo[]` prop
  - [ ] Create cycle detection helper to identify nodes/edges in cycles
  - [ ] Update D3Node and D3Link types to include `inCycle: boolean` flag

- [ ] Task 2: Implement cycle node highlighting (AC: 1)
  - [ ] Apply red border/glow styling to nodes in cycles
  - [ ] Add CSS classes for cycle vs normal nodes
  - [ ] Ensure styling works with both light and dark themes

- [ ] Task 3: Implement cycle edge highlighting (AC: 2)
  - [ ] Apply red color to edges forming cycles
  - [ ] Increase stroke-width for cycle edges
  - [ ] Add separate arrow marker for cycle edges (red)

- [ ] Task 4: Implement cycle path animation (AC: 3)
  - [ ] Add CSS animation for pulsing effect on cycle edges
  - [ ] Use SVG stroke-dasharray + animation for flowing effect
  - [ ] Ensure animation is performant (60fps)

- [ ] Task 5: Create legend component (AC: 4)
  - [ ] Create `GraphLegend` component
  - [ ] Show color coding for nodes and edges
  - [ ] Position legend in corner (configurable)

- [ ] Task 6: Implement click-to-highlight interaction (AC: 5, 6)
  - [ ] Add click handler to cycle nodes and edges
  - [ ] Track selected cycle in component state
  - [ ] Apply highlight class to selected cycle
  - [ ] Apply dimmed class to non-selected elements
  - [ ] Handle Escape key to deselect

- [ ] Task 7: Write unit tests (AC: all)
  - [ ] Test cycle detection helper function
  - [ ] Test node highlighting renders correctly
  - [ ] Test edge highlighting renders correctly
  - [ ] Test click interaction works
  - [ ] Test legend displays correctly

- [ ] Task 8: Verify CI passes (AC-CI)
  - [ ] Run `pnpm nx affected --target=lint --base=main`
  - [ ] Run `pnpm nx affected --target=test --base=main`
  - [ ] Run `pnpm nx affected --target=type-check --base=main`
  - [ ] Verify GitHub Actions CI is GREEN

## Dev Notes

### Architecture Patterns & Constraints

**Dependency on Story 4.1:** This story extends the DependencyGraph component created in Story 4.1.

**File Location:** `apps/web/app/components/visualization/DependencyGraph/`

**Updated Component Structure:**
```
apps/web/app/components/visualization/DependencyGraph/
├── index.tsx                    # Main component (from 4.1, extended)
├── types.ts                     # Extended with cycle flags
├── useForceSimulation.ts        # Force simulation hook
├── useCycleHighlight.ts         # NEW: Cycle highlighting logic
├── GraphLegend.tsx              # NEW: Legend component
├── styles.ts                    # NEW: D3 styling constants
└── __tests__/
    ├── DependencyGraph.test.tsx
    └── useCycleHighlight.test.ts  # NEW
```

### Data Types

**CircularDependencyInfo (from @monoguard/types):**
```typescript
interface CircularDependencyInfo {
  cycle: string[]           // Package names in cycle order
  type: 'direct' | 'indirect'
  severity: 'critical' | 'warning' | 'info'
  depth: number
  impact: string
  complexity: number
  priorityScore: number
  // ... additional fields from Epic 3
}
```

**Extended D3 Types:**
```typescript
interface D3Node extends d3.SimulationNodeDatum {
  id: string
  name: string
  path: string
  dependencyCount: number
  inCycle: boolean         // NEW: true if node is part of any cycle
  cycleIds: number[]       // NEW: which cycles this node belongs to
}

interface D3Link extends d3.SimulationLinkDatum<D3Node> {
  source: string | D3Node
  target: string | D3Node
  type: DependencyType
  inCycle: boolean         // NEW: true if edge is part of any cycle
  cycleIds: number[]       // NEW: which cycles this edge belongs to
}
```

### Implementation Patterns

**Cycle Detection Helper:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/useCycleHighlight.ts
import type { CircularDependencyInfo } from '@monoguard/types';

interface CycleHighlightResult {
  cycleNodeIds: Set<string>;
  cycleEdges: Map<string, number[]>; // "from->to" => cycleIds
  getCycleById: (id: number) => CircularDependencyInfo | undefined;
}

export function useCycleHighlight(
  circularDependencies: CircularDependencyInfo[] | undefined
): CycleHighlightResult {
  return useMemo(() => {
    const cycleNodeIds = new Set<string>();
    const cycleEdges = new Map<string, number[]>();

    if (!circularDependencies) {
      return { cycleNodeIds, cycleEdges, getCycleById: () => undefined };
    }

    circularDependencies.forEach((cycle, index) => {
      // Add all nodes in cycle
      cycle.cycle.forEach(node => cycleNodeIds.add(node));

      // Add all edges in cycle (consecutive pairs)
      for (let i = 0; i < cycle.cycle.length - 1; i++) {
        const edgeKey = `${cycle.cycle[i]}->${cycle.cycle[i + 1]}`;
        const existing = cycleEdges.get(edgeKey) || [];
        cycleEdges.set(edgeKey, [...existing, index]);
      }
    });

    return {
      cycleNodeIds,
      cycleEdges,
      getCycleById: (id) => circularDependencies[id],
    };
  }, [circularDependencies]);
}
```

**Styling Constants:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/styles.ts
export const GRAPH_COLORS = {
  node: {
    normal: {
      fill: '#4f46e5',      // Indigo
      stroke: '#ffffff',
    },
    cycle: {
      fill: '#ef4444',      // Red-500
      stroke: '#fecaca',    // Red-200 (glow effect)
    },
    selected: {
      fill: '#dc2626',      // Red-600
      stroke: '#ffffff',
    },
    dimmed: {
      fill: '#9ca3af',      // Gray-400
      stroke: '#d1d5db',    // Gray-300
    },
  },
  edge: {
    normal: {
      stroke: '#9ca3af',    // Gray-400
      width: 1,
    },
    cycle: {
      stroke: '#ef4444',    // Red-500
      width: 2.5,
    },
    selected: {
      stroke: '#dc2626',    // Red-600
      width: 3,
    },
    dimmed: {
      stroke: '#d1d5db',    // Gray-300
      width: 0.5,
    },
  },
};

export const ANIMATION = {
  pulseDuration: '1.5s',
  flowDuration: '2s',
};
```

**SVG Animation for Cycle Edges:**
```typescript
// In D3 rendering code
const cycleEdge = svg.selectAll('.cycle-edge')
  .data(cycleLinks)
  .join('line')
  .attr('class', 'cycle-edge')
  .attr('stroke', GRAPH_COLORS.edge.cycle.stroke)
  .attr('stroke-width', GRAPH_COLORS.edge.cycle.width)
  .attr('stroke-dasharray', '10,5')
  .style('animation', `flowAnimation ${ANIMATION.flowDuration} linear infinite`);

// CSS animation (add to component or global CSS)
// @keyframes flowAnimation {
//   0% { stroke-dashoffset: 15; }
//   100% { stroke-dashoffset: 0; }
// }
```

**Legend Component:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/GraphLegend.tsx
import React from 'react';
import { GRAPH_COLORS } from './styles';

export function GraphLegend() {
  return (
    <div className="absolute bottom-4 left-4 bg-white/90 dark:bg-gray-800/90
                    rounded-lg shadow-lg p-3 text-xs">
      <div className="font-semibold mb-2">Legend</div>
      <div className="space-y-1">
        <div className="flex items-center gap-2">
          <div
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: GRAPH_COLORS.node.normal.fill }}
          />
          <span>Normal Package</span>
        </div>
        <div className="flex items-center gap-2">
          <div
            className="w-3 h-3 rounded-full"
            style={{
              backgroundColor: GRAPH_COLORS.node.cycle.fill,
              boxShadow: `0 0 4px ${GRAPH_COLORS.node.cycle.stroke}`
            }}
          />
          <span>In Circular Dependency</span>
        </div>
        <div className="flex items-center gap-2">
          <div
            className="w-6 h-0.5"
            style={{ backgroundColor: GRAPH_COLORS.edge.normal.stroke }}
          />
          <span>Normal Dependency</span>
        </div>
        <div className="flex items-center gap-2">
          <div
            className="w-6 h-0.5"
            style={{ backgroundColor: GRAPH_COLORS.edge.cycle.stroke }}
          />
          <span>Circular Dependency</span>
        </div>
      </div>
    </div>
  );
}
```

### Testing Requirements

**Test File:** `apps/web/src/__tests__/useCycleHighlight.test.ts`

```typescript
import { renderHook } from '@testing-library/react';
import { useCycleHighlight } from '@/components/visualization/DependencyGraph/useCycleHighlight';
import type { CircularDependencyInfo } from '@monoguard/types';

describe('useCycleHighlight', () => {
  const mockCycles: CircularDependencyInfo[] = [
    {
      cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
      type: 'indirect',
      severity: 'warning',
      depth: 3,
      impact: 'Build may fail',
      complexity: 5,
      priorityScore: 75,
    },
    {
      cycle: ['pkg-x', 'pkg-y', 'pkg-x'],
      type: 'direct',
      severity: 'critical',
      depth: 2,
      impact: 'Immediate failure',
      complexity: 3,
      priorityScore: 90,
    },
  ];

  it('should identify all nodes in cycles', () => {
    const { result } = renderHook(() => useCycleHighlight(mockCycles));

    expect(result.current.cycleNodeIds.has('pkg-a')).toBe(true);
    expect(result.current.cycleNodeIds.has('pkg-b')).toBe(true);
    expect(result.current.cycleNodeIds.has('pkg-c')).toBe(true);
    expect(result.current.cycleNodeIds.has('pkg-x')).toBe(true);
    expect(result.current.cycleNodeIds.has('pkg-y')).toBe(true);
    expect(result.current.cycleNodeIds.has('pkg-z')).toBe(false);
  });

  it('should identify all edges in cycles', () => {
    const { result } = renderHook(() => useCycleHighlight(mockCycles));

    expect(result.current.cycleEdges.has('pkg-a->pkg-b')).toBe(true);
    expect(result.current.cycleEdges.has('pkg-b->pkg-c')).toBe(true);
    expect(result.current.cycleEdges.has('pkg-c->pkg-a')).toBe(true);
    expect(result.current.cycleEdges.has('pkg-x->pkg-y')).toBe(true);
    expect(result.current.cycleEdges.has('pkg-y->pkg-x')).toBe(true);
  });

  it('should handle undefined circular dependencies', () => {
    const { result } = renderHook(() => useCycleHighlight(undefined));

    expect(result.current.cycleNodeIds.size).toBe(0);
    expect(result.current.cycleEdges.size).toBe(0);
  });

  it('should return correct cycle by id', () => {
    const { result } = renderHook(() => useCycleHighlight(mockCycles));

    expect(result.current.getCycleById(0)).toEqual(mockCycles[0]);
    expect(result.current.getCycleById(1)).toEqual(mockCycles[1]);
    expect(result.current.getCycleById(2)).toBeUndefined();
  });
});
```

### Critical Don't-Miss Rules (from project-context.md)

1. **NEVER forget D3.js cleanup** - Ensure all event listeners are removed
2. **Use React.memo** - Already in place from Story 4.1
3. **Performance** - Animation should not drop frames (60fps)
4. **Accessibility** - Color is not the only indicator (use patterns/thickness too)
5. **Dark mode** - Ensure colors work in both light and dark themes

### Previous Story Intelligence (Story 4.1)

**Key Patterns to Follow:**
- Component structure with separate hooks and types
- D3 initialization in useEffect with cleanup
- Data transformation helper functions
- Comprehensive test coverage

**Files Created in Story 4.1:**
- `apps/web/app/components/visualization/DependencyGraph/index.tsx`
- `apps/web/app/components/visualization/DependencyGraph/types.ts`
- `apps/web/app/components/visualization/DependencyGraph/useForceSimulation.ts`

**Integration Point:** This story extends the component from 4.1 by adding:
- New prop: `circularDependencies?: CircularDependencyInfo[]`
- New hook: `useCycleHighlight`
- New component: `GraphLegend`
- Updated styling logic based on cycle membership

### References

- [Story 4.1: DependencyGraph Implementation] `4-1-implement-d3js-force-directed-dependency-graph.md`
- [Types: CircularDependencyInfo] `packages/types/src/analysis/results.ts:37`
- [Architecture: D3.js Visualization] `_bmad-output/planning-artifacts/architecture.md` - Decision 6
- [Epic 4 Story 4.2 Requirements] `_bmad-output/planning-artifacts/epics.md` - Lines 998-1008
- [Project Context: D3.js Rules] `_bmad-output/project-context.md` - D3.js Integration section
- [Epic 3 Retrospective: CI Requirements] `_bmad-output/implementation-artifacts/epic-3-retrospective.md`

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

