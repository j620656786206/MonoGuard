# Story 4.5: Implement Hover Details and Tooltips

Status: ready-for-review

## Story

As a **user**,
I want **to see package details when hovering over nodes**,
So that **I can quickly understand each package without clicking**.

## Acceptance Criteria

### AC1: Tooltip Content Display
**Given** a dependency graph
**When** I hover over a node
**Then** I see a tooltip with:
- Package name and path
- Dependency count (in/out - incoming and outgoing)
- Health contribution score
- Circular dependency involvement (if any)

### AC2: Tooltip Timing
**Given** a dependency graph
**When** I hover over a node
**Then**:
- Tooltip appears within 200ms
- Tooltip has smooth fade-in animation
- Tooltip disappears when mouse leaves the node

### AC3: Tooltip Positioning
**Given** a tooltip is displayed
**When** viewing the tooltip
**Then**:
- Tooltip follows mouse cursor OR anchors to node (developer choice)
- Tooltip stays within viewport bounds (doesn't clip off edges)
- Tooltip position adjusts when near viewport edges

### AC4: Edge Highlighting on Hover
**Given** a dependency graph
**When** I hover over a node
**Then**:
- Connected edges (both incoming and outgoing) are highlighted
- Non-connected edges are dimmed
- Connected nodes are subtly highlighted
- Clear visual distinction between highlighted and dimmed elements

### AC5: Performance Requirements
**Given** rapid mouse movements over nodes
**When** quickly moving between nodes
**Then**:
- No visual lag or stuttering
- Tooltips update smoothly
- Edge highlighting transitions smoothly
- No memory leaks from rapid hover events

### AC6: Edge Tooltip (Optional Enhancement)
**Given** a dependency graph
**When** I hover over an edge/link
**Then**:
- Optional tooltip shows dependency type (production/dev/peer)
- Shows source and target package names

### AC7: Accessibility
**Given** a tooltip is displayed
**When** using screen readers or keyboard
**Then**:
- Tooltip content is accessible via ARIA attributes
- Focus-based tooltip trigger for keyboard users (optional)

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

- [x] Task 1: Create Tooltip Component (AC: 1, 2, 3, 7)
  - [x] Create `NodeTooltip.tsx` component in DependencyGraph folder
  - [x] Implement tooltip content with all required fields
  - [x] Add fade-in animation (200ms or less)
  - [x] Implement viewport boundary detection and repositioning
  - [x] Add ARIA attributes for accessibility

- [x] Task 2: Create useNodeHover Hook (AC: 2, 4, 5)
  - [x] Create `useNodeHover.ts` custom hook
  - [x] Implement mouse enter/leave event handling
  - [x] Track hovered node ID
  - [x] Debounce rapid hover events if needed
  - [x] Manage hover state with proper cleanup

- [x] Task 3: Implement Edge Highlighting (AC: 4)
  - [x] Add `isHighlighted` state to edges based on hovered node
  - [x] Compute connected edges (incoming and outgoing) for hovered node
  - [x] Add CSS transitions for highlight/dim effect
  - [x] Dim non-connected edges (reduce opacity)
  - [x] Slightly highlight connected nodes

- [x] Task 4: Calculate Tooltip Data (AC: 1)
  - [x] Calculate incoming dependency count (edges where node is target)
  - [x] Calculate outgoing dependency count (edges where node is source)
  - [x] Extract health contribution score from node data
  - [x] Determine circular dependency involvement from cycle data

- [x] Task 5: Integrate with Main DependencyGraph Component (AC: all)
  - [x] Import and use NodeTooltip component
  - [x] Wire up useNodeHover hook
  - [x] Pass hover state to edge/node rendering
  - [x] Handle tooltip positioning relative to SVG container

- [ ] Task 6: Optional Edge Tooltip (AC: 6) - SKIPPED (Optional Enhancement)
  - [ ] Create EdgeTooltip component (simpler than node tooltip)
  - [ ] Show dependency type and connected packages
  - [ ] Implement hover detection on edges

- [x] Task 7: Write Unit Tests (AC: all)
  - [x] Test tooltip renders with correct content
  - [x] Test tooltip positioning logic
  - [x] Test edge highlighting computation
  - [x] Test hover state management
  - [x] Test performance with rapid hover events

- [x] Task 8: Verify CI passes (AC-CI)
  - [x] Run `pnpm nx affected --target=lint --base=main`
  - [x] Run `pnpm nx affected --target=test --base=main`
  - [x] Run `pnpm nx affected --target=type-check --base=main`
  - [x] Build passes successfully

## Dev Notes

### Architecture Patterns & Constraints

**Dependency on Stories 4.1-4.4:** This story extends the DependencyGraph component with hover interactions and tooltips.

**File Location:** `apps/web/app/components/visualization/DependencyGraph/`

**Updated Component Structure:**
```
apps/web/app/components/visualization/DependencyGraph/
├── index.tsx                      # Main component (extended with hover)
├── types.ts                       # Extended with hover/tooltip types
├── useForceSimulation.ts          # Force simulation hook
├── useCycleHighlight.ts           # From Story 4.2
├── useNodeExpandCollapse.ts       # From Story 4.3
├── useZoomPan.ts                  # From Story 4.4
├── useNodeHover.ts                # NEW: Hover state management
├── GraphLegend.tsx                # From Story 4.2
├── GraphControls.tsx              # From Story 4.3
├── ZoomControls.tsx               # From Story 4.4
├── GraphMinimap.tsx               # From Story 4.4
├── NodeTooltip.tsx                # NEW: Node tooltip component
├── EdgeTooltip.tsx                # NEW: Optional edge tooltip
├── styles.ts                      # Updated with tooltip styles
├── utils/
│   ├── computeVisibleNodes.ts     # From Story 4.3
│   ├── calculateDepth.ts          # From Story 4.3
│   ├── calculateBounds.ts         # From Story 4.4
│   └── computeConnectedElements.ts # NEW: Calculate connected edges/nodes
└── __tests__/
    ├── DependencyGraph.test.tsx
    ├── useCycleHighlight.test.ts
    ├── useNodeExpandCollapse.test.ts
    ├── useZoomPan.test.ts
    ├── useNodeHover.test.ts        # NEW
    ├── NodeTooltip.test.tsx        # NEW
    └── computeConnectedElements.test.ts # NEW
```

### Key Implementation Details

**Tooltip Type Definitions:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/types.ts
// Add to existing types

export interface TooltipData {
  packageName: string;
  packagePath: string;
  incomingCount: number;      // Dependencies pointing TO this node
  outgoingCount: number;      // Dependencies this node points TO
  healthContribution: number; // Impact on overall health score
  inCycle: boolean;           // Whether node is in a circular dependency
  cycleInfo?: {               // If in cycle, which cycle(s)
    cycleCount: number;
    packages: string[];
  };
}

export interface TooltipPosition {
  x: number;
  y: number;
  placement: 'top' | 'bottom' | 'left' | 'right';
}

export interface HoverState {
  nodeId: string | null;
  position: { x: number; y: number } | null;
}
```

**useNodeHover Hook:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/useNodeHover.ts
import { useState, useCallback, useMemo } from 'react';
import type { D3Node, D3Link, HoverState } from './types';

interface UseNodeHoverProps {
  nodes: D3Node[];
  links: D3Link[];
}

interface UseNodeHoverResult {
  hoverState: HoverState;
  connectedNodeIds: Set<string>;
  connectedLinkIndices: Set<number>;
  handleNodeMouseEnter: (nodeId: string, event: MouseEvent) => void;
  handleNodeMouseLeave: () => void;
  handleNodeMouseMove: (event: MouseEvent) => void;
}

export function useNodeHover({ nodes, links }: UseNodeHoverProps): UseNodeHoverResult {
  const [hoverState, setHoverState] = useState<HoverState>({
    nodeId: null,
    position: null,
  });

  // Compute connected elements when hover changes
  const { connectedNodeIds, connectedLinkIndices } = useMemo(() => {
    if (!hoverState.nodeId) {
      return { connectedNodeIds: new Set<string>(), connectedLinkIndices: new Set<number>() };
    }

    const nodeIds = new Set<string>([hoverState.nodeId]);
    const linkIndices = new Set<number>();

    links.forEach((link, index) => {
      const sourceId = typeof link.source === 'string' ? link.source : link.source.id;
      const targetId = typeof link.target === 'string' ? link.target : link.target.id;

      if (sourceId === hoverState.nodeId || targetId === hoverState.nodeId) {
        linkIndices.add(index);
        nodeIds.add(sourceId);
        nodeIds.add(targetId);
      }
    });

    return { connectedNodeIds: nodeIds, connectedLinkIndices: linkIndices };
  }, [hoverState.nodeId, links]);

  const handleNodeMouseEnter = useCallback((nodeId: string, event: MouseEvent) => {
    setHoverState({
      nodeId,
      position: { x: event.clientX, y: event.clientY },
    });
  }, []);

  const handleNodeMouseMove = useCallback((event: MouseEvent) => {
    setHoverState(prev => ({
      ...prev,
      position: { x: event.clientX, y: event.clientY },
    }));
  }, []);

  const handleNodeMouseLeave = useCallback(() => {
    setHoverState({ nodeId: null, position: null });
  }, []);

  return {
    hoverState,
    connectedNodeIds,
    connectedLinkIndices,
    handleNodeMouseEnter,
    handleNodeMouseLeave,
    handleNodeMouseMove,
  };
}
```

**NodeTooltip Component:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/NodeTooltip.tsx
import React, { useMemo, useRef, useEffect, useState } from 'react';
import type { TooltipData, TooltipPosition } from './types';

interface NodeTooltipProps {
  data: TooltipData | null;
  position: { x: number; y: number } | null;
  containerRef: React.RefObject<HTMLDivElement>;
}

const TOOLTIP_OFFSET = 12;
const ANIMATION_DURATION = 150; // ms, under 200ms requirement

export function NodeTooltip({ data, position, containerRef }: NodeTooltipProps) {
  const tooltipRef = useRef<HTMLDivElement>(null);
  const [calculatedPosition, setCalculatedPosition] = useState<TooltipPosition | null>(null);
  const [isVisible, setIsVisible] = useState(false);

  // Calculate position to keep tooltip in viewport
  useEffect(() => {
    if (!data || !position || !tooltipRef.current || !containerRef.current) {
      setIsVisible(false);
      return;
    }

    const tooltipRect = tooltipRef.current.getBoundingClientRect();
    const containerRect = containerRef.current.getBoundingClientRect();

    let x = position.x - containerRect.left + TOOLTIP_OFFSET;
    let y = position.y - containerRect.top + TOOLTIP_OFFSET;
    let placement: 'top' | 'bottom' | 'left' | 'right' = 'right';

    // Adjust if tooltip would clip right edge
    if (x + tooltipRect.width > containerRect.width) {
      x = position.x - containerRect.left - tooltipRect.width - TOOLTIP_OFFSET;
      placement = 'left';
    }

    // Adjust if tooltip would clip bottom edge
    if (y + tooltipRect.height > containerRect.height) {
      y = position.y - containerRect.top - tooltipRect.height - TOOLTIP_OFFSET;
      placement = placement === 'left' ? 'left' : 'top';
    }

    // Adjust if tooltip would clip left edge
    if (x < 0) {
      x = TOOLTIP_OFFSET;
    }

    // Adjust if tooltip would clip top edge
    if (y < 0) {
      y = TOOLTIP_OFFSET;
    }

    setCalculatedPosition({ x, y, placement });
    setIsVisible(true);
  }, [data, position, containerRef]);

  if (!data) return null;

  const shortPath = data.packagePath.split('/').slice(-2).join('/');

  return (
    <div
      ref={tooltipRef}
      role="tooltip"
      aria-live="polite"
      className={`
        absolute z-50 pointer-events-none
        bg-white dark:bg-gray-800 rounded-lg shadow-xl
        border border-gray-200 dark:border-gray-700
        p-3 min-w-[200px] max-w-[300px]
        transition-opacity duration-150 ease-out
        ${isVisible ? 'opacity-100' : 'opacity-0'}
      `}
      style={{
        left: calculatedPosition?.x ?? 0,
        top: calculatedPosition?.y ?? 0,
        transitionDuration: `${ANIMATION_DURATION}ms`,
      }}
    >
      {/* Package Name */}
      <div className="font-semibold text-gray-900 dark:text-white truncate">
        {data.packageName}
      </div>

      {/* Package Path */}
      <div className="text-xs text-gray-500 dark:text-gray-400 mb-2 truncate">
        {shortPath}
      </div>

      {/* Dependency Counts */}
      <div className="flex gap-4 text-sm mb-2">
        <div>
          <span className="text-gray-500 dark:text-gray-400">In:</span>{' '}
          <span className="font-medium text-green-600 dark:text-green-400">
            {data.incomingCount}
          </span>
        </div>
        <div>
          <span className="text-gray-500 dark:text-gray-400">Out:</span>{' '}
          <span className="font-medium text-blue-600 dark:text-blue-400">
            {data.outgoingCount}
          </span>
        </div>
      </div>

      {/* Health Contribution */}
      <div className="text-sm mb-2">
        <span className="text-gray-500 dark:text-gray-400">Health Impact:</span>{' '}
        <span className={`font-medium ${
          data.healthContribution >= 0
            ? 'text-green-600 dark:text-green-400'
            : 'text-red-600 dark:text-red-400'
        }`}>
          {data.healthContribution >= 0 ? '+' : ''}{data.healthContribution}
        </span>
      </div>

      {/* Circular Dependency Warning */}
      {data.inCycle && (
        <div className="flex items-center gap-1 text-sm text-red-600 dark:text-red-400
                        bg-red-50 dark:bg-red-900/20 rounded px-2 py-1">
          <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd"
                  d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                  clipRule="evenodd" />
          </svg>
          <span>
            In {data.cycleInfo?.cycleCount ?? 1} circular
            {(data.cycleInfo?.cycleCount ?? 1) > 1 ? ' dependencies' : ' dependency'}
          </span>
        </div>
      )}
    </div>
  );
}
```

**Compute Connected Elements Utility:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/utils/computeConnectedElements.ts

import type { D3Link } from '../types';

interface ConnectedElements {
  nodeIds: Set<string>;
  linkIndices: Set<number>;
}

/**
 * Computes all nodes and links connected to a given node.
 * Includes both incoming (dependencies of this node) and outgoing (dependents).
 */
export function computeConnectedElements(
  nodeId: string,
  links: D3Link[]
): ConnectedElements {
  const nodeIds = new Set<string>([nodeId]);
  const linkIndices = new Set<number>();

  links.forEach((link, index) => {
    const sourceId = typeof link.source === 'string' ? link.source : link.source.id;
    const targetId = typeof link.target === 'string' ? link.target : link.target.id;

    if (sourceId === nodeId || targetId === nodeId) {
      linkIndices.add(index);
      nodeIds.add(sourceId);
      nodeIds.add(targetId);
    }
  });

  return { nodeIds, linkIndices };
}

/**
 * Computes incoming and outgoing dependency counts for a node.
 */
export function computeDependencyCounts(
  nodeId: string,
  links: D3Link[]
): { incoming: number; outgoing: number } {
  let incoming = 0;
  let outgoing = 0;

  links.forEach(link => {
    const sourceId = typeof link.source === 'string' ? link.source : link.source.id;
    const targetId = typeof link.target === 'string' ? link.target : link.target.id;

    if (targetId === nodeId) {
      incoming++;
    }
    if (sourceId === nodeId) {
      outgoing++;
    }
  });

  return { incoming, outgoing };
}
```

**Compute Tooltip Data:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/utils/computeTooltipData.ts

import type { D3Node, D3Link, TooltipData, CircularDependency } from '../types';
import { computeDependencyCounts } from './computeConnectedElements';

interface ComputeTooltipDataParams {
  node: D3Node;
  links: D3Link[];
  circularDependencies: CircularDependency[];
  totalNodes: number;
}

/**
 * Computes all data needed for the node tooltip.
 */
export function computeTooltipData({
  node,
  links,
  circularDependencies,
  totalNodes,
}: ComputeTooltipDataParams): TooltipData {
  // Compute dependency counts
  const { incoming, outgoing } = computeDependencyCounts(node.id, links);

  // Check if node is in any cycles
  const cyclesContainingNode = circularDependencies.filter(cycle =>
    cycle.path.includes(node.id)
  );

  const inCycle = cyclesContainingNode.length > 0;

  // Calculate health contribution
  // Negative if in cycles, more negative with more connections in cycles
  // Positive if well-structured with reasonable dependency count
  let healthContribution = 0;

  if (inCycle) {
    // Each cycle involvement reduces health
    healthContribution = -5 * cyclesContainingNode.length;
    // Additional penalty for being hub in cycle
    if (incoming > 3 || outgoing > 3) {
      healthContribution -= 2;
    }
  } else {
    // Good node contributes positively
    // More connections = less contribution (coupling concern)
    const totalConnections = incoming + outgoing;
    if (totalConnections <= 3) {
      healthContribution = 2;
    } else if (totalConnections <= 6) {
      healthContribution = 1;
    } else {
      healthContribution = 0; // High coupling
    }
  }

  // Get packages in the same cycle(s)
  const cyclePackages = new Set<string>();
  cyclesContainingNode.forEach(cycle => {
    cycle.path.forEach(pkg => {
      if (pkg !== node.id) {
        cyclePackages.add(pkg);
      }
    });
  });

  return {
    packageName: node.name,
    packagePath: node.path,
    incomingCount: incoming,
    outgoingCount: outgoing,
    healthContribution,
    inCycle,
    cycleInfo: inCycle
      ? {
          cycleCount: cyclesContainingNode.length,
          packages: Array.from(cyclePackages),
        }
      : undefined,
  };
}
```

### Integration with Main Component

```typescript
// In apps/web/app/components/visualization/DependencyGraph/index.tsx
// Add to existing component from Stories 4.1-4.4

import { useNodeHover } from './useNodeHover';
import { NodeTooltip } from './NodeTooltip';
import { computeTooltipData } from './utils/computeTooltipData';

export const DependencyGraphViz = React.memo(function DependencyGraphViz({
  data,
  circularDependencies,
}: DependencyGraphProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const svgRef = useRef<SVGSVGElement>(null);
  // ... existing refs from previous stories

  const {
    hoverState,
    connectedNodeIds,
    connectedLinkIndices,
    handleNodeMouseEnter,
    handleNodeMouseLeave,
    handleNodeMouseMove,
  } = useNodeHover({
    nodes: visibleNodes,
    links: visibleLinks,
  });

  // Compute tooltip data when hovering
  const tooltipData = useMemo(() => {
    if (!hoverState.nodeId) return null;

    const node = visibleNodes.find(n => n.id === hoverState.nodeId);
    if (!node) return null;

    return computeTooltipData({
      node,
      links: visibleLinks,
      circularDependencies,
      totalNodes: visibleNodes.length,
    });
  }, [hoverState.nodeId, visibleNodes, visibleLinks, circularDependencies]);

  // In D3 setup useEffect, add hover handlers to nodes
  useEffect(() => {
    // ... existing D3 setup ...

    // Add hover event handlers to nodes
    node
      .on('mouseenter', function(event, d) {
        handleNodeMouseEnter(d.id, event);
      })
      .on('mousemove', function(event) {
        handleNodeMouseMove(event);
      })
      .on('mouseleave', function() {
        handleNodeMouseLeave();
      });

    // Apply highlighting styles based on hover state
    // This is reactive via React state change

    return () => {
      simulation.stop();
      svg.on('.zoom', null);
      svg.selectAll('*').remove();
    };
  }, [data, handleNodeMouseEnter, handleNodeMouseLeave, handleNodeMouseMove]);

  // Effect to update visual highlighting when hover changes
  useEffect(() => {
    if (!svgRef.current) return;

    const svg = d3.select(svgRef.current);

    if (hoverState.nodeId) {
      // Dim non-connected elements
      svg.selectAll('line')
        .transition()
        .duration(150)
        .attr('stroke-opacity', (d, i) =>
          connectedLinkIndices.has(i) ? 0.8 : 0.15
        )
        .attr('stroke-width', (d, i) =>
          connectedLinkIndices.has(i) ? 2 : 1
        );

      svg.selectAll('circle')
        .transition()
        .duration(150)
        .attr('opacity', (d: D3Node) =>
          connectedNodeIds.has(d.id) ? 1 : 0.3
        );

      svg.selectAll('text')
        .transition()
        .duration(150)
        .attr('opacity', (d: D3Node) =>
          connectedNodeIds.has(d.id) ? 1 : 0.3
        );
    } else {
      // Reset all elements
      svg.selectAll('line')
        .transition()
        .duration(150)
        .attr('stroke-opacity', 0.6)
        .attr('stroke-width', 1);

      svg.selectAll('circle')
        .transition()
        .duration(150)
        .attr('opacity', 1);

      svg.selectAll('text')
        .transition()
        .duration(150)
        .attr('opacity', 1);
    }
  }, [hoverState.nodeId, connectedNodeIds, connectedLinkIndices]);

  return (
    <div ref={containerRef} className="relative w-full h-full">
      <svg ref={svgRef} className="w-full h-full min-h-[500px]">
        {/* SVG content rendered by D3 */}
      </svg>

      {/* Tooltip - rendered in React for easier positioning */}
      <NodeTooltip
        data={tooltipData}
        position={hoverState.position}
        containerRef={containerRef}
      />

      {/* Other controls from previous stories */}
      <GraphMinimap ... />
      <ZoomControls ... />
      <GraphLegend />
      <GraphControls ... />
    </div>
  );
});
```

### Testing Requirements

**Test File:** `apps/web/src/__tests__/useNodeHover.test.ts`

```typescript
import { renderHook, act } from '@testing-library/react';
import { useNodeHover } from '@/components/visualization/DependencyGraph/useNodeHover';

describe('useNodeHover', () => {
  const mockNodes = [
    { id: 'A', name: 'Package A', path: '/a', dependencyCount: 2 },
    { id: 'B', name: 'Package B', path: '/b', dependencyCount: 1 },
    { id: 'C', name: 'Package C', path: '/c', dependencyCount: 1 },
  ];

  const mockLinks = [
    { source: 'A', target: 'B', type: 'production' as const },
    { source: 'A', target: 'C', type: 'production' as const },
  ];

  it('should initialize with null hover state', () => {
    const { result } = renderHook(() =>
      useNodeHover({ nodes: mockNodes as any, links: mockLinks as any })
    );

    expect(result.current.hoverState.nodeId).toBeNull();
    expect(result.current.hoverState.position).toBeNull();
  });

  it('should update hover state on mouse enter', () => {
    const { result } = renderHook(() =>
      useNodeHover({ nodes: mockNodes as any, links: mockLinks as any })
    );

    act(() => {
      result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent);
    });

    expect(result.current.hoverState.nodeId).toBe('A');
    expect(result.current.hoverState.position).toEqual({ x: 100, y: 200 });
  });

  it('should clear hover state on mouse leave', () => {
    const { result } = renderHook(() =>
      useNodeHover({ nodes: mockNodes as any, links: mockLinks as any })
    );

    act(() => {
      result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent);
    });

    act(() => {
      result.current.handleNodeMouseLeave();
    });

    expect(result.current.hoverState.nodeId).toBeNull();
  });

  it('should compute connected elements correctly', () => {
    const { result } = renderHook(() =>
      useNodeHover({ nodes: mockNodes as any, links: mockLinks as any })
    );

    act(() => {
      result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent);
    });

    expect(result.current.connectedNodeIds.has('A')).toBe(true);
    expect(result.current.connectedNodeIds.has('B')).toBe(true);
    expect(result.current.connectedNodeIds.has('C')).toBe(true);
    expect(result.current.connectedLinkIndices.has(0)).toBe(true);
    expect(result.current.connectedLinkIndices.has(1)).toBe(true);
  });

  it('should update position on mouse move', () => {
    const { result } = renderHook(() =>
      useNodeHover({ nodes: mockNodes as any, links: mockLinks as any })
    );

    act(() => {
      result.current.handleNodeMouseEnter('A', { clientX: 100, clientY: 200 } as MouseEvent);
    });

    act(() => {
      result.current.handleNodeMouseMove({ clientX: 150, clientY: 250 } as MouseEvent);
    });

    expect(result.current.hoverState.position).toEqual({ x: 150, y: 250 });
  });
});
```

**Test File:** `apps/web/src/__tests__/NodeTooltip.test.tsx`

```typescript
import { render, screen } from '@testing-library/react';
import { NodeTooltip } from '@/components/visualization/DependencyGraph/NodeTooltip';
import React from 'react';

describe('NodeTooltip', () => {
  const mockContainerRef = { current: document.createElement('div') };

  beforeEach(() => {
    mockContainerRef.current.getBoundingClientRect = vi.fn(() => ({
      left: 0,
      top: 0,
      width: 800,
      height: 600,
      right: 800,
      bottom: 600,
      x: 0,
      y: 0,
      toJSON: () => {},
    }));
  });

  const mockData = {
    packageName: '@app/core',
    packagePath: 'packages/core',
    incomingCount: 3,
    outgoingCount: 5,
    healthContribution: 2,
    inCycle: false,
  };

  it('should render nothing when data is null', () => {
    const { container } = render(
      <NodeTooltip data={null} position={null} containerRef={mockContainerRef as any} />
    );
    expect(container.firstChild).toBeNull();
  });

  it('should display package name', () => {
    render(
      <NodeTooltip
        data={mockData}
        position={{ x: 100, y: 100 }}
        containerRef={mockContainerRef as any}
      />
    );
    expect(screen.getByText('@app/core')).toBeInTheDocument();
  });

  it('should display dependency counts', () => {
    render(
      <NodeTooltip
        data={mockData}
        position={{ x: 100, y: 100 }}
        containerRef={mockContainerRef as any}
      />
    );
    expect(screen.getByText('3')).toBeInTheDocument(); // incoming
    expect(screen.getByText('5')).toBeInTheDocument(); // outgoing
  });

  it('should display health contribution', () => {
    render(
      <NodeTooltip
        data={mockData}
        position={{ x: 100, y: 100 }}
        containerRef={mockContainerRef as any}
      />
    );
    expect(screen.getByText('+2')).toBeInTheDocument();
  });

  it('should show circular dependency warning when in cycle', () => {
    const cycleData = {
      ...mockData,
      inCycle: true,
      cycleInfo: { cycleCount: 2, packages: ['@app/utils', '@app/shared'] },
    };

    render(
      <NodeTooltip
        data={cycleData}
        position={{ x: 100, y: 100 }}
        containerRef={mockContainerRef as any}
      />
    );
    expect(screen.getByText(/circular dependencies/i)).toBeInTheDocument();
  });

  it('should have proper accessibility attributes', () => {
    render(
      <NodeTooltip
        data={mockData}
        position={{ x: 100, y: 100 }}
        containerRef={mockContainerRef as any}
      />
    );
    expect(screen.getByRole('tooltip')).toBeInTheDocument();
  });
});
```

**Test File:** `apps/web/src/__tests__/computeConnectedElements.test.ts`

```typescript
import {
  computeConnectedElements,
  computeDependencyCounts,
} from '@/components/visualization/DependencyGraph/utils/computeConnectedElements';

describe('computeConnectedElements', () => {
  const mockLinks = [
    { source: 'A', target: 'B', type: 'production' as const },
    { source: 'A', target: 'C', type: 'production' as const },
    { source: 'B', target: 'D', type: 'production' as const },
    { source: 'C', target: 'A', type: 'production' as const }, // Cycle
  ];

  it('should find all connected nodes', () => {
    const result = computeConnectedElements('A', mockLinks as any);

    expect(result.nodeIds.has('A')).toBe(true);
    expect(result.nodeIds.has('B')).toBe(true);
    expect(result.nodeIds.has('C')).toBe(true);
    expect(result.nodeIds.has('D')).toBe(false); // Not directly connected
  });

  it('should find all connected link indices', () => {
    const result = computeConnectedElements('A', mockLinks as any);

    expect(result.linkIndices.has(0)).toBe(true); // A -> B
    expect(result.linkIndices.has(1)).toBe(true); // A -> C
    expect(result.linkIndices.has(2)).toBe(false); // B -> D (not connected to A directly)
    expect(result.linkIndices.has(3)).toBe(true); // C -> A
  });

  it('should handle nodes with no connections', () => {
    const result = computeConnectedElements('isolated', mockLinks as any);

    expect(result.nodeIds.size).toBe(1);
    expect(result.nodeIds.has('isolated')).toBe(true);
    expect(result.linkIndices.size).toBe(0);
  });
});

describe('computeDependencyCounts', () => {
  const mockLinks = [
    { source: 'A', target: 'B', type: 'production' as const },
    { source: 'A', target: 'C', type: 'production' as const },
    { source: 'D', target: 'A', type: 'production' as const },
  ];

  it('should count incoming dependencies correctly', () => {
    const result = computeDependencyCounts('A', mockLinks as any);
    expect(result.incoming).toBe(1); // D -> A
  });

  it('should count outgoing dependencies correctly', () => {
    const result = computeDependencyCounts('A', mockLinks as any);
    expect(result.outgoing).toBe(2); // A -> B, A -> C
  });

  it('should return zero counts for isolated nodes', () => {
    const result = computeDependencyCounts('isolated', mockLinks as any);
    expect(result.incoming).toBe(0);
    expect(result.outgoing).toBe(0);
  });
});
```

### Critical Don't-Miss Rules (from project-context.md)

1. **NEVER forget D3.js cleanup** - Remove event listeners in cleanup
2. **Use React.memo** - Already in place from Story 4.1
3. **Tooltip appears within 200ms** - Animation duration set to 150ms
4. **Performance** - Memoize tooltip data computation, debounce if needed
5. **Accessibility** - Include ARIA attributes on tooltip element
6. **Dark mode support** - Use Tailwind dark: variants for all colors

### Previous Story Intelligence (Stories 4.1, 4.2, 4.3, 4.4)

**Key Patterns to Follow:**
- Component structure with separate hooks for each concern
- D3 initialization in useEffect with cleanup
- Comprehensive test coverage for each AC
- Consistent styling with existing components

**Integration Points:**
- Tooltips should work with zoom/pan from Story 4.4
- Circular dependency nodes should show cycle info in tooltip
- Highlight colors should match cycle highlighting from Story 4.2
- Tooltips should respect expanded/collapsed state from Story 4.3

**Existing Files to Extend:**
- `index.tsx` - Add hover handlers and tooltip rendering
- `types.ts` - Add TooltipData and HoverState types
- `styles.ts` - Add tooltip styles if needed

### UX Design Requirements (from ux-design-specification.md)

- **Timing:** Tooltip within 200ms (implemented as 150ms for snappier feel)
- **Content:** Package name, path, in/out counts, health, cycle status
- **Visual:** Follows mouse or anchors to node
- **Interaction:** Hovering highlights connected paths
- **Dark mode:** Proper contrast in both light and dark themes

### Performance Considerations

1. **Memoization:** Tooltip data computed with useMemo
2. **Debouncing:** Consider debouncing rapid hover events if performance issues
3. **D3 Transitions:** Use short durations (150ms) for smooth but quick
4. **Set operations:** Use Set for O(1) connected element lookup

### References

- [Story 4.1: DependencyGraph Implementation] `4-1-implement-d3js-force-directed-dependency-graph.md`
- [Story 4.2: Cycle Highlighting] `4-2-highlight-circular-dependencies-in-graph.md`
- [Story 4.3: Node Expand/Collapse] `4-3-implement-node-expand-collapse-functionality.md`
- [Story 4.4: Zoom/Pan Controls] `4-4-add-zoom-pan-and-navigation-controls.md`
- [Epic 4 Story 4.5 Requirements] `_bmad-output/planning-artifacts/epics.md` - Lines 1054-1075
- [UX Design: Interactive Graph] `_bmad-output/planning-artifacts/ux-design-specification.md` - Line 148
- [Project Context: D3.js Rules] `_bmad-output/project-context.md` - D3.js Integration section
- [D3 Selection Events] https://d3js.org/d3-selection/events

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Fixed lint errors in NodeTooltip.test.tsx (DOMRect type issue - replaced with MockRect interface)
- Fixed type errors in computeConnectedElements.ts (CircularDependencyInfo uses `cycle` not `path`)
- Fixed test mocks to match actual CircularDependencyInfo type structure

### Completion Notes List

- All 7 core tasks completed (Task 6 skipped as optional)
- 51 new tests added across 3 test files (NodeTooltip: 20, useNodeHover: 15, computeConnectedElements: 16)
- Total test count: 373 tests passing
- CI verification: lint (warnings only, no errors), test, type-check, build all pass
- Tooltip appears within 150ms (under 200ms requirement)
- Edge highlighting uses 150ms transitions for smooth animation
- ARIA attributes included for accessibility (role="tooltip", aria-live="polite")

### File List

**New Files Created:**
- `apps/web/app/components/visualization/DependencyGraph/NodeTooltip.tsx`
- `apps/web/app/components/visualization/DependencyGraph/useNodeHover.ts`
- `apps/web/app/components/visualization/DependencyGraph/utils/computeConnectedElements.ts`
- `apps/web/app/components/visualization/DependencyGraph/__tests__/NodeTooltip.test.tsx`
- `apps/web/app/components/visualization/DependencyGraph/__tests__/useNodeHover.test.ts`
- `apps/web/app/components/visualization/DependencyGraph/__tests__/computeConnectedElements.test.ts`

**Modified Files:**
- `apps/web/app/components/visualization/DependencyGraph/index.tsx` - Added hover integration
- `apps/web/app/components/visualization/DependencyGraph/types.ts` - Added TooltipData, TooltipPosition, HoverState types

