# Story 4.4: Add Zoom, Pan, and Navigation Controls

Status: ready-for-dev

## Story

As a **user**,
I want **to zoom and pan the dependency graph**,
So that **I can navigate large graphs effectively**.

## Acceptance Criteria

### AC1: Scroll Zoom Behavior
**Given** a dependency graph
**When** I scroll (mouse wheel) on the graph
**Then**:
- Graph zooms in/out centered on cursor position
- Zoom is smooth with animation
- Zoom respects min/max limits (10% to 400%)

### AC2: Click and Drag Pan
**Given** a dependency graph
**When** I click and drag on the graph background
**Then**:
- Graph pans in the direction of drag
- Pan is smooth and responsive
- Nodes remain interactive during pan

### AC3: Zoom Control Buttons
**Given** a dependency graph
**When** I use zoom control buttons
**Then** I can:
- Click "+" to zoom in by a fixed increment (e.g., 20%)
- Click "-" to zoom out by a fixed increment
- Buttons are always visible in a fixed position
- Buttons show disabled state at zoom limits

### AC4: Fit to Screen
**Given** a dependency graph of any size
**When** I click "Fit to screen" button
**Then**:
- Entire graph becomes visible in the viewport
- Graph is centered in the container
- Zoom level adjusts automatically
- Smooth transition animation

### AC5: Minimap Navigation
**Given** a dependency graph with > 50 nodes
**When** viewing the graph
**Then** I see:
- A minimap in the corner showing the entire graph
- Current viewport highlighted on minimap
- Can click/drag on minimap to navigate
- Minimap is hidden for small graphs (< 50 nodes)

### AC6: Zoom Level Display
**Given** any zoom operation
**When** zoom level changes
**Then**:
- Current zoom percentage is displayed (e.g., "100%")
- Display updates in real-time during zoom
- Display is positioned near zoom controls

### AC7: Zoom Range Limits
**Given** zoom operations
**When** attempting to zoom beyond limits
**Then**:
- Zoom stops at 10% minimum (prevent too small)
- Zoom stops at 400% maximum (prevent too large)
- Visual/haptic feedback at limits (optional)

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

- [ ] Task 1: Implement D3 zoom behavior (AC: 1, 2, 7)
  - [ ] Configure `d3.zoom()` with scale extent [0.1, 4]
  - [ ] Attach zoom behavior to SVG element
  - [ ] Handle wheel events for zoom
  - [ ] Handle drag events for pan
  - [ ] Add zoom transform to graph container group

- [ ] Task 2: Create zoom controls component (AC: 3, 6)
  - [ ] Create `ZoomControls.tsx` component
  - [ ] Implement zoom in button with increment logic
  - [ ] Implement zoom out button with decrement logic
  - [ ] Add current zoom level display
  - [ ] Style controls with proper positioning

- [ ] Task 3: Implement fit-to-screen functionality (AC: 4)
  - [ ] Calculate bounding box of all nodes
  - [ ] Compute required scale and translation
  - [ ] Implement smooth transition to fit view
  - [ ] Add "Fit" button to zoom controls

- [ ] Task 4: Create minimap component (AC: 5)
  - [ ] Create `GraphMinimap.tsx` component
  - [ ] Render scaled-down version of graph
  - [ ] Show viewport indicator rectangle
  - [ ] Implement click-to-navigate on minimap
  - [ ] Implement drag-to-navigate on minimap
  - [ ] Add conditional rendering (> 50 nodes only)

- [ ] Task 5: Integrate zoom state with React (AC: 1, 2, 3, 6)
  - [ ] Create `useZoomPan.ts` custom hook
  - [ ] Sync D3 zoom state with React state
  - [ ] Expose zoom controls API (zoomIn, zoomOut, fitToScreen, resetZoom)
  - [ ] Handle zoom level as percentage for display

- [ ] Task 6: Handle edge cases (AC: 7)
  - [ ] Prevent zoom beyond limits
  - [ ] Handle empty graph (no nodes)
  - [ ] Handle single node graph
  - [ ] Ensure touch device compatibility (pinch-to-zoom)

- [ ] Task 7: Write unit tests (AC: all)
  - [ ] Test zoom controls render correctly
  - [ ] Test zoom in/out functions
  - [ ] Test fit-to-screen calculation
  - [ ] Test minimap visibility logic
  - [ ] Test zoom limits enforcement

- [ ] Task 8: Verify CI passes (AC-CI)
  - [ ] Run `pnpm nx affected --target=lint --base=main`
  - [ ] Run `pnpm nx affected --target=test --base=main`
  - [ ] Run `pnpm nx affected --target=type-check --base=main`
  - [ ] Verify GitHub Actions CI is GREEN

## Dev Notes

### Architecture Patterns & Constraints

**Dependency on Stories 4.1, 4.2, 4.3:** This story extends the DependencyGraph component with zoom/pan capabilities.

**File Location:** `apps/web/app/components/visualization/DependencyGraph/`

**Updated Component Structure:**
```
apps/web/app/components/visualization/DependencyGraph/
├── index.tsx                      # Main component (extended with zoom)
├── types.ts                       # Extended with zoom types
├── useForceSimulation.ts          # Force simulation hook
├── useCycleHighlight.ts           # From Story 4.2
├── useNodeExpandCollapse.ts       # From Story 4.3
├── useZoomPan.ts                  # NEW: Zoom and pan state management
├── GraphLegend.tsx                # From Story 4.2
├── GraphControls.tsx              # From Story 4.3 (extended)
├── ZoomControls.tsx               # NEW: Zoom buttons and display
├── GraphMinimap.tsx               # NEW: Minimap navigation
├── styles.ts                      # Updated with zoom control styles
├── utils/
│   ├── computeVisibleNodes.ts     # From Story 4.3
│   ├── calculateDepth.ts          # From Story 4.3
│   └── calculateBounds.ts         # NEW: Bounding box calculation
└── __tests__/
    ├── DependencyGraph.test.tsx
    ├── useCycleHighlight.test.ts
    ├── useNodeExpandCollapse.test.ts
    ├── useZoomPan.test.ts         # NEW
    ├── ZoomControls.test.tsx      # NEW
    └── GraphMinimap.test.tsx      # NEW
```

### Key Implementation Details

**D3 Zoom Configuration:**
```typescript
// In main component - D3 zoom setup
import * as d3 from 'd3';

const ZOOM_CONFIG = {
  scaleExtent: [0.1, 4] as [number, number],  // 10% to 400%
  zoomIncrement: 0.2,                          // 20% per button click
  transitionDuration: 300,                     // ms
};

// Create zoom behavior
const zoom = d3.zoom<SVGSVGElement, unknown>()
  .scaleExtent(ZOOM_CONFIG.scaleExtent)
  .on('zoom', (event) => {
    // Apply transform to container group
    container.attr('transform', event.transform);
    // Update React state for UI display
    onZoomChange(event.transform.k);
  });

// Attach to SVG
svg.call(zoom);

// Prevent double-click zoom (conflicts with expand/collapse)
svg.on('dblclick.zoom', null);
```

**Zoom Pan Hook:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/useZoomPan.ts
import { useState, useCallback, useRef, useEffect } from 'react';
import * as d3 from 'd3';

interface ZoomPanState {
  scale: number;
  translateX: number;
  translateY: number;
}

interface UseZoomPanResult {
  zoomState: ZoomPanState;
  zoomPercent: number;
  zoomIn: () => void;
  zoomOut: () => void;
  resetZoom: () => void;
  fitToScreen: () => void;
  setZoomBehavior: (zoom: d3.ZoomBehavior<SVGSVGElement, unknown>) => void;
  canZoomIn: boolean;
  canZoomOut: boolean;
}

interface UseZoomPanProps {
  svgRef: React.RefObject<SVGSVGElement>;
  containerRef: React.RefObject<SVGGElement>;
  minScale?: number;
  maxScale?: number;
  zoomIncrement?: number;
}

export function useZoomPan({
  svgRef,
  containerRef,
  minScale = 0.1,
  maxScale = 4,
  zoomIncrement = 0.2,
}: UseZoomPanProps): UseZoomPanResult {
  const [zoomState, setZoomState] = useState<ZoomPanState>({
    scale: 1,
    translateX: 0,
    translateY: 0,
  });

  const zoomBehaviorRef = useRef<d3.ZoomBehavior<SVGSVGElement, unknown> | null>(null);

  const zoomPercent = Math.round(zoomState.scale * 100);
  const canZoomIn = zoomState.scale < maxScale;
  const canZoomOut = zoomState.scale > minScale;

  const setZoomBehavior = useCallback((zoom: d3.ZoomBehavior<SVGSVGElement, unknown>) => {
    zoomBehaviorRef.current = zoom;
  }, []);

  const zoomIn = useCallback(() => {
    if (!svgRef.current || !zoomBehaviorRef.current) return;

    const svg = d3.select(svgRef.current);
    const newScale = Math.min(zoomState.scale + zoomIncrement, maxScale);

    svg.transition()
      .duration(200)
      .call(zoomBehaviorRef.current.scaleTo, newScale);
  }, [svgRef, zoomState.scale, zoomIncrement, maxScale]);

  const zoomOut = useCallback(() => {
    if (!svgRef.current || !zoomBehaviorRef.current) return;

    const svg = d3.select(svgRef.current);
    const newScale = Math.max(zoomState.scale - zoomIncrement, minScale);

    svg.transition()
      .duration(200)
      .call(zoomBehaviorRef.current.scaleTo, newScale);
  }, [svgRef, zoomState.scale, zoomIncrement, minScale]);

  const resetZoom = useCallback(() => {
    if (!svgRef.current || !zoomBehaviorRef.current) return;

    const svg = d3.select(svgRef.current);

    svg.transition()
      .duration(300)
      .call(zoomBehaviorRef.current.transform, d3.zoomIdentity);
  }, [svgRef]);

  const fitToScreen = useCallback(() => {
    if (!svgRef.current || !containerRef.current || !zoomBehaviorRef.current) return;

    const svg = d3.select(svgRef.current);
    const svgNode = svgRef.current;
    const containerNode = containerRef.current;

    // Get SVG dimensions
    const { width: svgWidth, height: svgHeight } = svgNode.getBoundingClientRect();

    // Get container bounds (all nodes)
    const bounds = containerNode.getBBox();

    if (bounds.width === 0 || bounds.height === 0) return;

    // Calculate scale to fit with padding
    const padding = 40;
    const scaleX = (svgWidth - padding * 2) / bounds.width;
    const scaleY = (svgHeight - padding * 2) / bounds.height;
    const scale = Math.min(scaleX, scaleY, maxScale);

    // Calculate translation to center
    const translateX = (svgWidth - bounds.width * scale) / 2 - bounds.x * scale;
    const translateY = (svgHeight - bounds.height * scale) / 2 - bounds.y * scale;

    const transform = d3.zoomIdentity
      .translate(translateX, translateY)
      .scale(scale);

    svg.transition()
      .duration(500)
      .call(zoomBehaviorRef.current.transform, transform);
  }, [svgRef, containerRef, maxScale]);

  // Update state when zoom changes
  const handleZoomChange = useCallback((transform: d3.ZoomTransform) => {
    setZoomState({
      scale: transform.k,
      translateX: transform.x,
      translateY: transform.y,
    });
  }, []);

  return {
    zoomState,
    zoomPercent,
    zoomIn,
    zoomOut,
    resetZoom,
    fitToScreen,
    setZoomBehavior,
    canZoomIn,
    canZoomOut,
  };
}
```

**Zoom Controls Component:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/ZoomControls.tsx
import React from 'react';

interface ZoomControlsProps {
  zoomPercent: number;
  onZoomIn: () => void;
  onZoomOut: () => void;
  onFitToScreen: () => void;
  onResetZoom: () => void;
  canZoomIn: boolean;
  canZoomOut: boolean;
}

export function ZoomControls({
  zoomPercent,
  onZoomIn,
  onZoomOut,
  onFitToScreen,
  onResetZoom,
  canZoomIn,
  canZoomOut,
}: ZoomControlsProps) {
  return (
    <div className="absolute bottom-4 right-4 bg-white/90 dark:bg-gray-800/90
                    rounded-lg shadow-lg p-2 flex items-center gap-1">
      {/* Zoom Out Button */}
      <button
        onClick={onZoomOut}
        disabled={!canZoomOut}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
                   disabled:opacity-40 disabled:cursor-not-allowed"
        aria-label="Zoom out"
        title="Zoom out"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
        </svg>
      </button>

      {/* Zoom Level Display */}
      <div className="w-14 text-center text-sm font-medium tabular-nums">
        {zoomPercent}%
      </div>

      {/* Zoom In Button */}
      <button
        onClick={onZoomIn}
        disabled={!canZoomIn}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors
                   disabled:opacity-40 disabled:cursor-not-allowed"
        aria-label="Zoom in"
        title="Zoom in"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
        </svg>
      </button>

      {/* Divider */}
      <div className="w-px h-6 bg-gray-200 dark:bg-gray-700 mx-1" />

      {/* Fit to Screen Button */}
      <button
        onClick={onFitToScreen}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        aria-label="Fit to screen"
        title="Fit to screen"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5" />
        </svg>
      </button>

      {/* Reset Zoom Button */}
      <button
        onClick={onResetZoom}
        className="w-8 h-8 flex items-center justify-center rounded
                   hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        aria-label="Reset zoom to 100%"
        title="Reset zoom to 100%"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>
  );
}
```

**Minimap Component:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/GraphMinimap.tsx
import React, { useRef, useEffect, useMemo } from 'react';
import * as d3 from 'd3';
import type { D3Node, D3Link } from './types';

interface GraphMinimapProps {
  nodes: D3Node[];
  links: D3Link[];
  viewportBounds: { x: number; y: number; width: number; height: number };
  graphBounds: { x: number; y: number; width: number; height: number };
  onNavigate: (x: number, y: number) => void;
  width?: number;
  height?: number;
}

const MINIMAP_SIZE = { width: 150, height: 100 };
const MIN_NODES_FOR_MINIMAP = 50;

export function GraphMinimap({
  nodes,
  links,
  viewportBounds,
  graphBounds,
  onNavigate,
  width = MINIMAP_SIZE.width,
  height = MINIMAP_SIZE.height,
}: GraphMinimapProps) {
  const svgRef = useRef<SVGSVGElement>(null);

  // Don't render for small graphs
  if (nodes.length < MIN_NODES_FOR_MINIMAP) {
    return null;
  }

  // Calculate scale to fit graph in minimap
  const scale = useMemo(() => {
    if (graphBounds.width === 0 || graphBounds.height === 0) return 1;
    const scaleX = width / graphBounds.width;
    const scaleY = height / graphBounds.height;
    return Math.min(scaleX, scaleY) * 0.9; // 90% to add padding
  }, [graphBounds, width, height]);

  // Calculate offset to center graph in minimap
  const offset = useMemo(() => ({
    x: (width - graphBounds.width * scale) / 2 - graphBounds.x * scale,
    y: (height - graphBounds.height * scale) / 2 - graphBounds.y * scale,
  }), [width, height, graphBounds, scale]);

  // Calculate viewport rectangle in minimap coordinates
  const viewportRect = useMemo(() => ({
    x: viewportBounds.x * scale + offset.x,
    y: viewportBounds.y * scale + offset.y,
    width: viewportBounds.width * scale,
    height: viewportBounds.height * scale,
  }), [viewportBounds, scale, offset]);

  const handleClick = (event: React.MouseEvent<SVGSVGElement>) => {
    if (!svgRef.current) return;

    const rect = svgRef.current.getBoundingClientRect();
    const clickX = event.clientX - rect.left;
    const clickY = event.clientY - rect.top;

    // Convert click to graph coordinates
    const graphX = (clickX - offset.x) / scale;
    const graphY = (clickY - offset.y) / scale;

    onNavigate(graphX, graphY);
  };

  return (
    <div className="absolute top-4 left-4 bg-white/90 dark:bg-gray-800/90
                    rounded-lg shadow-lg p-1 border border-gray-200 dark:border-gray-700">
      <svg
        ref={svgRef}
        width={width}
        height={height}
        className="cursor-pointer"
        onClick={handleClick}
      >
        {/* Background */}
        <rect
          width={width}
          height={height}
          fill="transparent"
        />

        {/* Graph content group */}
        <g transform={`translate(${offset.x}, ${offset.y}) scale(${scale})`}>
          {/* Links */}
          {links.map((link, i) => {
            const source = typeof link.source === 'string'
              ? nodes.find(n => n.id === link.source)
              : link.source;
            const target = typeof link.target === 'string'
              ? nodes.find(n => n.id === link.target)
              : link.target;

            if (!source || !target) return null;

            return (
              <line
                key={i}
                x1={source.x || 0}
                y1={source.y || 0}
                x2={target.x || 0}
                y2={target.y || 0}
                stroke="#9ca3af"
                strokeWidth={0.5 / scale}
                strokeOpacity={0.5}
              />
            );
          })}

          {/* Nodes */}
          {nodes.map(node => (
            <circle
              key={node.id}
              cx={node.x || 0}
              cy={node.y || 0}
              r={3 / scale}
              fill={node.inCycle ? '#ef4444' : '#4f46e5'}
            />
          ))}
        </g>

        {/* Viewport indicator */}
        <rect
          x={viewportRect.x}
          y={viewportRect.y}
          width={Math.max(viewportRect.width, 10)}
          height={Math.max(viewportRect.height, 10)}
          fill="rgba(99, 102, 241, 0.2)"
          stroke="#6366f1"
          strokeWidth={1.5}
          rx={2}
        />
      </svg>
    </div>
  );
}
```

**Calculate Bounds Utility:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/utils/calculateBounds.ts

export interface Bounds {
  x: number;
  y: number;
  width: number;
  height: number;
}

export function calculateNodeBounds(
  nodes: Array<{ x?: number; y?: number }>
): Bounds {
  if (nodes.length === 0) {
    return { x: 0, y: 0, width: 0, height: 0 };
  }

  let minX = Infinity;
  let minY = Infinity;
  let maxX = -Infinity;
  let maxY = -Infinity;

  nodes.forEach(node => {
    const x = node.x ?? 0;
    const y = node.y ?? 0;
    minX = Math.min(minX, x);
    minY = Math.min(minY, y);
    maxX = Math.max(maxX, x);
    maxY = Math.max(maxY, y);
  });

  // Add padding for node radius
  const padding = 20;
  return {
    x: minX - padding,
    y: minY - padding,
    width: maxX - minX + padding * 2,
    height: maxY - minY + padding * 2,
  };
}

export function calculateViewportBounds(
  transform: { k: number; x: number; y: number },
  svgWidth: number,
  svgHeight: number
): Bounds {
  // Inverse transform to get viewport in graph coordinates
  return {
    x: -transform.x / transform.k,
    y: -transform.y / transform.k,
    width: svgWidth / transform.k,
    height: svgHeight / transform.k,
  };
}
```

### Integration with Main Component

```typescript
// In apps/web/app/components/visualization/DependencyGraph/index.tsx
// Add to existing component from Stories 4.1-4.3

import { useZoomPan } from './useZoomPan';
import { ZoomControls } from './ZoomControls';
import { GraphMinimap } from './GraphMinimap';
import { calculateNodeBounds, calculateViewportBounds } from './utils/calculateBounds';

export const DependencyGraphViz = React.memo(function DependencyGraphViz({
  data,
  circularDependencies,
}: DependencyGraphProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<SVGGElement>(null);

  const {
    zoomState,
    zoomPercent,
    zoomIn,
    zoomOut,
    resetZoom,
    fitToScreen,
    setZoomBehavior,
    canZoomIn,
    canZoomOut,
  } = useZoomPan({
    svgRef,
    containerRef,
  });

  // Track graph bounds for minimap
  const [graphBounds, setGraphBounds] = useState<Bounds>({ x: 0, y: 0, width: 0, height: 0 });

  useEffect(() => {
    if (!svgRef.current) return;

    const svg = d3.select(svgRef.current);
    const container = svg.append('g').attr('class', 'graph-container');
    containerRef.current = container.node();

    // Setup zoom behavior
    const zoom = d3.zoom<SVGSVGElement, unknown>()
      .scaleExtent([0.1, 4])
      .on('zoom', (event) => {
        container.attr('transform', event.transform);
        // Update zoom state in hook
        handleZoomChange(event.transform);
      });

    svg.call(zoom);

    // Disable double-click zoom (conflicts with expand/collapse)
    svg.on('dblclick.zoom', null);

    // Store zoom behavior for external control
    setZoomBehavior(zoom);

    // ... rest of D3 setup ...

    // After simulation stabilizes, update bounds
    simulation.on('end', () => {
      const bounds = calculateNodeBounds(nodes);
      setGraphBounds(bounds);
    });

    return () => {
      simulation.stop();
      svg.on('.zoom', null);
      svg.selectAll('*').remove();
    };
  }, [data]);

  // Calculate viewport bounds for minimap
  const viewportBounds = useMemo(() => {
    if (!svgRef.current) return { x: 0, y: 0, width: 0, height: 0 };
    const { width, height } = svgRef.current.getBoundingClientRect();
    return calculateViewportBounds(
      { k: zoomState.scale, x: zoomState.translateX, y: zoomState.translateY },
      width,
      height
    );
  }, [zoomState]);

  // Navigate from minimap
  const handleMinimapNavigate = useCallback((x: number, y: number) => {
    if (!svgRef.current || !zoomBehaviorRef.current) return;

    const svg = d3.select(svgRef.current);
    const { width, height } = svgRef.current.getBoundingClientRect();

    // Center viewport on clicked position
    const transform = d3.zoomIdentity
      .translate(width / 2 - x * zoomState.scale, height / 2 - y * zoomState.scale)
      .scale(zoomState.scale);

    svg.transition()
      .duration(300)
      .call(zoomBehaviorRef.current.transform, transform);
  }, [zoomState.scale]);

  return (
    <div className="relative w-full h-full">
      <svg ref={svgRef} className="w-full h-full min-h-[500px]">
        <g ref={containerRef} />
      </svg>

      {/* Minimap - only for large graphs */}
      <GraphMinimap
        nodes={visibleNodes}
        links={visibleLinks}
        viewportBounds={viewportBounds}
        graphBounds={graphBounds}
        onNavigate={handleMinimapNavigate}
      />

      {/* Zoom Controls */}
      <ZoomControls
        zoomPercent={zoomPercent}
        onZoomIn={zoomIn}
        onZoomOut={zoomOut}
        onFitToScreen={fitToScreen}
        onResetZoom={resetZoom}
        canZoomIn={canZoomIn}
        canZoomOut={canZoomOut}
      />

      {/* Legend and other controls from previous stories */}
      <GraphLegend />
      <GraphControls ... />
    </div>
  );
});
```

### Testing Requirements

**Test File:** `apps/web/src/__tests__/useZoomPan.test.ts`

```typescript
import { renderHook, act } from '@testing-library/react';
import { useZoomPan } from '@/components/visualization/DependencyGraph/useZoomPan';

describe('useZoomPan', () => {
  const mockSvgRef = { current: document.createElementNS('http://www.w3.org/2000/svg', 'svg') };
  const mockContainerRef = { current: document.createElementNS('http://www.w3.org/2000/svg', 'g') };

  it('should initialize with default zoom state', () => {
    const { result } = renderHook(() =>
      useZoomPan({
        svgRef: mockSvgRef as any,
        containerRef: mockContainerRef as any,
      })
    );

    expect(result.current.zoomState.scale).toBe(1);
    expect(result.current.zoomPercent).toBe(100);
    expect(result.current.canZoomIn).toBe(true);
    expect(result.current.canZoomOut).toBe(true);
  });

  it('should report correct zoom limits', () => {
    const { result } = renderHook(() =>
      useZoomPan({
        svgRef: mockSvgRef as any,
        containerRef: mockContainerRef as any,
        minScale: 0.5,
        maxScale: 2,
      })
    );

    // At scale 1, should be able to zoom both directions
    expect(result.current.canZoomIn).toBe(true);
    expect(result.current.canZoomOut).toBe(true);
  });

  it('should calculate zoom percentage correctly', () => {
    const { result } = renderHook(() =>
      useZoomPan({
        svgRef: mockSvgRef as any,
        containerRef: mockContainerRef as any,
      })
    );

    expect(result.current.zoomPercent).toBe(100);
  });
});
```

**Test File:** `apps/web/src/__tests__/ZoomControls.test.tsx`

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { ZoomControls } from '@/components/visualization/DependencyGraph/ZoomControls';

describe('ZoomControls', () => {
  const mockProps = {
    zoomPercent: 100,
    onZoomIn: vi.fn(),
    onZoomOut: vi.fn(),
    onFitToScreen: vi.fn(),
    onResetZoom: vi.fn(),
    canZoomIn: true,
    canZoomOut: true,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should display current zoom percentage', () => {
    render(<ZoomControls {...mockProps} zoomPercent={150} />);
    expect(screen.getByText('150%')).toBeInTheDocument();
  });

  it('should call onZoomIn when plus button is clicked', () => {
    render(<ZoomControls {...mockProps} />);
    fireEvent.click(screen.getByLabelText('Zoom in'));
    expect(mockProps.onZoomIn).toHaveBeenCalledTimes(1);
  });

  it('should call onZoomOut when minus button is clicked', () => {
    render(<ZoomControls {...mockProps} />);
    fireEvent.click(screen.getByLabelText('Zoom out'));
    expect(mockProps.onZoomOut).toHaveBeenCalledTimes(1);
  });

  it('should call onFitToScreen when fit button is clicked', () => {
    render(<ZoomControls {...mockProps} />);
    fireEvent.click(screen.getByLabelText('Fit to screen'));
    expect(mockProps.onFitToScreen).toHaveBeenCalledTimes(1);
  });

  it('should disable zoom in button when canZoomIn is false', () => {
    render(<ZoomControls {...mockProps} canZoomIn={false} />);
    expect(screen.getByLabelText('Zoom in')).toBeDisabled();
  });

  it('should disable zoom out button when canZoomOut is false', () => {
    render(<ZoomControls {...mockProps} canZoomOut={false} />);
    expect(screen.getByLabelText('Zoom out')).toBeDisabled();
  });
});
```

**Test File:** `apps/web/src/__tests__/GraphMinimap.test.tsx`

```typescript
import { render, screen } from '@testing-library/react';
import { GraphMinimap } from '@/components/visualization/DependencyGraph/GraphMinimap';

describe('GraphMinimap', () => {
  const mockNodes = Array.from({ length: 60 }, (_, i) => ({
    id: `node-${i}`,
    name: `Node ${i}`,
    path: `/path/${i}`,
    dependencyCount: 1,
    x: Math.random() * 500,
    y: Math.random() * 500,
    inCycle: i < 5, // First 5 nodes are in cycle
  }));

  const mockLinks = mockNodes.slice(0, -1).map((node, i) => ({
    source: node,
    target: mockNodes[i + 1],
    type: 'production' as const,
  }));

  const defaultProps = {
    nodes: mockNodes,
    links: mockLinks,
    viewportBounds: { x: 0, y: 0, width: 800, height: 600 },
    graphBounds: { x: 0, y: 0, width: 500, height: 500 },
    onNavigate: vi.fn(),
  };

  it('should render minimap for graphs with > 50 nodes', () => {
    const { container } = render(<GraphMinimap {...defaultProps} />);
    expect(container.querySelector('svg')).toBeInTheDocument();
  });

  it('should not render minimap for graphs with < 50 nodes', () => {
    const smallNodes = mockNodes.slice(0, 40);
    const { container } = render(
      <GraphMinimap {...defaultProps} nodes={smallNodes} />
    );
    expect(container.querySelector('svg')).not.toBeInTheDocument();
  });

  it('should render viewport indicator rectangle', () => {
    const { container } = render(<GraphMinimap {...defaultProps} />);
    const rects = container.querySelectorAll('rect');
    // Should have background rect and viewport rect
    expect(rects.length).toBeGreaterThanOrEqual(2);
  });
});
```

### Critical Don't-Miss Rules (from project-context.md)

1. **NEVER forget D3.js cleanup** - Remove zoom event listeners in cleanup
2. **Use React.memo** - Already in place from Story 4.1
3. **Disable double-click zoom** - Conflicts with expand/collapse from Story 4.3
4. **Smooth transitions** - Use D3 transitions with appropriate duration
5. **Zoom limits** - Enforce 10% to 400% range strictly
6. **Touch support** - D3 zoom handles pinch-to-zoom automatically

### Previous Story Intelligence (Stories 4.1, 4.2, 4.3)

**Key Patterns to Follow:**
- Component structure with separate hooks for each concern
- D3 initialization in useEffect with cleanup
- Comprehensive test coverage
- Consistent styling with existing components

**Integration Points:**
- Disable D3 double-click zoom (conflicts with expand/collapse)
- Minimap should respect cycle highlighting colors
- Zoom controls should be positioned to not overlap with GraphControls (Story 4.3)

**Existing Files to Extend:**
- `index.tsx` - Add zoom behavior setup
- `types.ts` - Add zoom-related types
- `styles.ts` - Add zoom control styles

### Accessibility Considerations

1. **Keyboard Support**: Consider Cmd/Ctrl + Plus/Minus for zoom
2. **Screen Reader**: Buttons have aria-labels
3. **Focus Indicators**: Ensure zoom buttons have visible focus states
4. **Motion Sensitivity**: Respect `prefers-reduced-motion` for transitions

### References

- [Story 4.1: DependencyGraph Implementation] `4-1-implement-d3js-force-directed-dependency-graph.md`
- [Story 4.2: Cycle Highlighting] `4-2-highlight-circular-dependencies-in-graph.md`
- [Story 4.3: Node Expand/Collapse] `4-3-implement-node-expand-collapse-functionality.md`
- [Epic 4 Story 4.4 Requirements] `_bmad-output/planning-artifacts/epics.md` - Lines 1032-1052
- [Architecture: D3.js Visualization] `_bmad-output/planning-artifacts/architecture.md` - Decision 6
- [Project Context: D3.js Rules] `_bmad-output/project-context.md` - D3.js Integration section
- [D3 Zoom Documentation] https://d3js.org/d3-zoom

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

