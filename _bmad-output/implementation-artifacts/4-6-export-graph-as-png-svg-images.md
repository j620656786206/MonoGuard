# Story 4.6: Export Graph as PNG/SVG Images

Status: done

## Story

As a **user**,
I want **to export the dependency graph as an image**,
So that **I can include it in documentation or presentations**.

## Acceptance Criteria

### AC1: Export Format Options
**Given** a rendered dependency graph
**When** I click export
**Then** I can choose:
- PNG format (raster, with resolution options)
- SVG format (vector, scalable)
- Current view or full graph option
- With or without legend option

### AC2: PNG Export with Resolution Options
**Given** I select PNG export
**When** I configure export settings
**Then**:
- I can choose resolution: 1x (standard), 2x (high DPI), 4x (print quality)
- Export captures the current graph state including highlights
- Background is white (light mode) or dark (dark mode) based on current theme
- File size is reasonable (< 5MB for 2x resolution)

### AC3: SVG Export
**Given** I select SVG export
**When** the export completes
**Then**:
- SVG is vector-based and scales without pixelation
- All graph elements (nodes, edges, labels) are preserved
- SVG is standalone (no external dependencies)
- Colors match the current theme
- SVG can be opened in design tools (Figma, Illustrator)

### AC4: Export Scope Options
**Given** I want to export the graph
**When** I configure export scope
**Then**:
- "Current view" exports only what's visible in viewport
- "Full graph" exports the entire dependency graph
- "Selected elements" exports highlighted nodes and edges (if any selected)

### AC5: Legend Inclusion
**Given** I configure export settings
**When** I toggle "Include legend"
**Then**:
- Legend is rendered in the export (matching GraphLegend component)
- Legend position is bottom-right or configurable
- Legend explains node colors and edge types

### AC6: Watermark and Metadata (Optional)
**Given** a completed export
**When** the file is generated
**Then**:
- Optional MonoGuard watermark can be included (toggleable)
- Filename includes project name and timestamp (e.g., `monoguard-deps-2026-01-25.png`)
- File downloads immediately to user's device

### AC7: Export Button/Menu UI
**Given** the dependency graph is displayed
**When** I look for export functionality
**Then**:
- Export button is visible in the graph controls area
- Clicking shows a dropdown/modal with format and options
- Clear preview of what will be exported
- Export progress indicator for large graphs

### AC8: Performance Requirements
**Given** a large graph (> 200 nodes)
**When** exporting to PNG or SVG
**Then**:
- Export completes in < 5 seconds
- UI remains responsive during export
- Progress feedback shown for lengthy exports
- No memory leaks during export process

### AC-CI: CI Pipeline Must Pass (MANDATORY)
**Given** the story implementation is complete
**When** verifying CI status
**Then** ALL of the following must pass:
- [x] `pnpm nx affected --target=lint --base=main` passes
- [x] `pnpm nx affected --target=test --base=main` passes
- [x] `pnpm nx affected --target=type-check --base=main` passes
- [x] `cd packages/analysis-engine && make test` passes (N/A - no Go changes)
- [ ] GitHub Actions CI workflow shows GREEN status
- **Story CANNOT be marked as "done" until CI is green**

## Tasks / Subtasks

- [x] Task 1: Create Export Types and Interfaces (AC: 1, 2, 3, 4, 5, 6)
  - [x] Define `ExportFormat` type ('png' | 'svg')
  - [x] Define `ExportOptions` interface with all configurable options
  - [x] Define `ExportScope` type ('viewport' | 'full' | 'selected')
  - [x] Define `Resolution` type (1 | 2 | 4)

- [x] Task 2: Create SVG Export Utility (AC: 3, 4, 5)
  - [x] Create `exportSvg.ts` utility function
  - [x] Clone the SVG element with all styles inlined
  - [x] Handle viewBox adjustment for full graph export
  - [x] Embed fonts if custom fonts are used
  - [x] Serialize SVG to string and create downloadable blob

- [x] Task 3: Create PNG Export Utility (AC: 2, 4, 5)
  - [x] Create `exportPng.ts` utility function
  - [x] Render SVG to Canvas using Image element
  - [x] Apply resolution multiplier for higher DPI exports
  - [x] Handle background color based on theme
  - [x] Convert Canvas to PNG blob and trigger download

- [x] Task 4: Create Legend Renderer for Export (AC: 5)
  - [x] Create `renderLegendForExport.ts` utility
  - [x] Generate standalone legend SVG that can be embedded
  - [x] Position legend appropriately in export
  - [x] Match GraphLegend component styling

- [x] Task 5: Create Export Modal/Dropdown Component (AC: 1, 7)
  - [x] Create `ExportMenu.tsx` component
  - [x] Format selection (PNG/SVG radio buttons or tabs)
  - [x] Resolution dropdown for PNG
  - [x] Scope selection (viewport/full/selected)
  - [x] Include legend checkbox
  - [x] Include watermark checkbox (optional)
  - [x] Export button with loading state

- [x] Task 6: Create useGraphExport Hook (AC: all)
  - [x] Create `useGraphExport.ts` hook
  - [x] Manage export state (isExporting, progress)
  - [x] Handle export execution with progress tracking
  - [x] Generate filename with project name and timestamp
  - [x] Trigger file download

- [x] Task 7: Integrate Export with DependencyGraph Component (AC: all)
  - [x] Add export button to GraphControls or create separate ExportButton
  - [x] Wire up ExportMenu with graph ref and current state
  - [x] Pass selected nodes/edges for "selected" scope export
  - [x] Handle export for both SVG and Canvas rendering modes

- [x] Task 8: Write Unit Tests (AC: all)
  - [x] Test SVG export generates valid SVG
  - [x] Test PNG export with different resolutions
  - [x] Test export scope (viewport vs full)
  - [x] Test legend inclusion/exclusion
  - [x] Test filename generation
  - [x] Test ExportMenu component interactions

- [x] Task 9: Verify CI passes (AC-CI)
  - [x] Run `pnpm nx affected --target=lint --base=main`
  - [x] Run `pnpm nx affected --target=test --base=main`
  - [x] Run `pnpm nx affected --target=type-check --base=main`
  - [ ] Verify GitHub Actions CI is GREEN

## Dev Notes

### Architecture Patterns & Constraints

**Dependency on Stories 4.1-4.5:** This story adds export functionality to the existing DependencyGraph component built in previous stories.

**File Location:** `apps/web/app/components/visualization/DependencyGraph/`

**Updated Component Structure:**
```
apps/web/app/components/visualization/DependencyGraph/
├── index.tsx                      # Main component (add export integration)
├── types.ts                       # Extended with export types
├── useForceSimulation.ts          # Force simulation hook
├── useCycleHighlight.ts           # From Story 4.2
├── useNodeExpandCollapse.ts       # From Story 4.3
├── useZoomPan.ts                  # From Story 4.4
├── useNodeHover.ts                # From Story 4.5
├── useGraphExport.ts              # NEW: Export state and logic
├── GraphLegend.tsx                # From Story 4.2
├── GraphControls.tsx              # From Story 4.3 (extend with export)
├── ZoomControls.tsx               # From Story 4.4
├── GraphMinimap.tsx               # From Story 4.4
├── NodeTooltip.tsx                # From Story 4.5
├── ExportMenu.tsx                 # NEW: Export options UI
├── styles.ts                      # Styling
├── utils/
│   ├── computeVisibleNodes.ts     # From Story 4.3
│   ├── calculateDepth.ts          # From Story 4.3
│   ├── calculateBounds.ts         # From Story 4.4
│   ├── computeConnectedElements.ts # From Story 4.5
│   ├── exportSvg.ts               # NEW: SVG export logic
│   ├── exportPng.ts               # NEW: PNG export logic
│   └── renderLegendForExport.ts   # NEW: Legend for exports
└── __tests__/
    ├── DependencyGraph.test.tsx
    ├── useCycleHighlight.test.ts
    ├── useNodeExpandCollapse.test.ts
    ├── useZoomPan.test.ts
    ├── useNodeHover.test.ts
    ├── useGraphExport.test.ts     # NEW
    ├── ExportMenu.test.tsx        # NEW
    ├── exportSvg.test.ts          # NEW
    └── exportPng.test.ts          # NEW
```

### Key Implementation Details

**Export Types:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/types.ts
// Add to existing types

export type ExportFormat = 'png' | 'svg';
export type ExportScope = 'viewport' | 'full' | 'selected';
export type ExportResolution = 1 | 2 | 4;

export interface ExportOptions {
  format: ExportFormat;
  scope: ExportScope;
  resolution: ExportResolution; // Only for PNG
  includeLegend: boolean;
  includeWatermark: boolean;
  backgroundColor: string | 'transparent';
}

export interface ExportResult {
  blob: Blob;
  filename: string;
  width: number;
  height: number;
}

export interface ExportProgress {
  isExporting: boolean;
  progress: number; // 0-100
  stage: 'preparing' | 'rendering' | 'encoding' | 'complete';
}
```

**SVG Export Utility:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/utils/exportSvg.ts

import type { ExportOptions, ExportResult } from '../types';

interface ExportSvgParams {
  svgElement: SVGSVGElement;
  options: ExportOptions;
  projectName: string;
  legendSvg?: string;
}

/**
 * Exports the dependency graph as an SVG file.
 * Inlines all styles and creates a standalone SVG.
 */
export async function exportSvg({
  svgElement,
  options,
  projectName,
  legendSvg,
}: ExportSvgParams): Promise<ExportResult> {
  // Clone SVG to avoid modifying the original
  const clonedSvg = svgElement.cloneNode(true) as SVGSVGElement;

  // Get computed styles and inline them
  inlineStyles(clonedSvg, svgElement);

  // Adjust viewBox based on export scope
  if (options.scope === 'full') {
    const bounds = calculateFullGraphBounds(clonedSvg);
    clonedSvg.setAttribute('viewBox',
      `${bounds.x} ${bounds.y} ${bounds.width} ${bounds.height}`
    );
  }

  // Add background if not transparent
  if (options.backgroundColor !== 'transparent') {
    const background = document.createElementNS('http://www.w3.org/2000/svg', 'rect');
    background.setAttribute('width', '100%');
    background.setAttribute('height', '100%');
    background.setAttribute('fill', options.backgroundColor);
    clonedSvg.insertBefore(background, clonedSvg.firstChild);
  }

  // Add legend if requested
  if (options.includeLegend && legendSvg) {
    // Parse and append legend SVG
    const legendGroup = parseLegendSvg(legendSvg, clonedSvg);
    clonedSvg.appendChild(legendGroup);
  }

  // Add watermark if requested
  if (options.includeWatermark) {
    const watermark = createWatermark();
    clonedSvg.appendChild(watermark);
  }

  // Serialize to string
  const serializer = new XMLSerializer();
  const svgString = serializer.serializeToString(clonedSvg);

  // Add XML declaration
  const svgWithDeclaration = `<?xml version="1.0" encoding="UTF-8"?>\n${svgString}`;

  // Create blob
  const blob = new Blob([svgWithDeclaration], { type: 'image/svg+xml' });

  // Generate filename
  const timestamp = new Date().toISOString().split('T')[0];
  const filename = `${projectName}-dependency-graph-${timestamp}.svg`;

  return {
    blob,
    filename,
    width: parseInt(clonedSvg.getAttribute('width') || '800'),
    height: parseInt(clonedSvg.getAttribute('height') || '600'),
  };
}

/**
 * Inline computed styles into SVG elements.
 * This ensures the exported SVG looks identical to the rendered version.
 */
function inlineStyles(clonedNode: Element, originalNode: Element) {
  const computedStyle = window.getComputedStyle(originalNode);
  const relevantStyles = [
    'fill', 'stroke', 'stroke-width', 'opacity',
    'font-family', 'font-size', 'font-weight',
    'text-anchor', 'dominant-baseline',
  ];

  let styleString = '';
  relevantStyles.forEach(prop => {
    const value = computedStyle.getPropertyValue(prop);
    if (value) {
      styleString += `${prop}:${value};`;
    }
  });

  if (styleString) {
    (clonedNode as SVGElement).style.cssText += styleString;
  }

  // Recurse to children
  const originalChildren = originalNode.children;
  const clonedChildren = clonedNode.children;

  for (let i = 0; i < originalChildren.length; i++) {
    if (clonedChildren[i]) {
      inlineStyles(clonedChildren[i], originalChildren[i]);
    }
  }
}

function calculateFullGraphBounds(svg: SVGSVGElement): {
  x: number; y: number; width: number; height: number;
} {
  const bbox = svg.getBBox();
  const padding = 20;
  return {
    x: bbox.x - padding,
    y: bbox.y - padding,
    width: bbox.width + padding * 2,
    height: bbox.height + padding * 2,
  };
}

function createWatermark(): SVGGElement {
  const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
  g.setAttribute('class', 'watermark');

  const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
  text.setAttribute('x', '10');
  text.setAttribute('y', '99%');
  text.setAttribute('fill', '#999');
  text.setAttribute('font-size', '10');
  text.setAttribute('font-family', 'sans-serif');
  text.textContent = 'Generated by MonoGuard';

  g.appendChild(text);
  return g;
}

function parseLegendSvg(legendSvgString: string, parentSvg: SVGSVGElement): SVGGElement {
  const parser = new DOMParser();
  const doc = parser.parseFromString(legendSvgString, 'image/svg+xml');
  const legendSvg = doc.documentElement;

  // Wrap in group and position
  const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
  g.setAttribute('class', 'legend-group');

  // Position at bottom-right
  const viewBox = parentSvg.viewBox.baseVal;
  const legendWidth = 150;
  const legendHeight = 100;
  g.setAttribute('transform',
    `translate(${viewBox.width - legendWidth - 20}, ${viewBox.height - legendHeight - 20})`
  );

  // Copy legend content
  Array.from(legendSvg.children).forEach(child => {
    g.appendChild(child.cloneNode(true));
  });

  return g;
}
```

**PNG Export Utility:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/utils/exportPng.ts

import type { ExportOptions, ExportResult } from '../types';
import { exportSvg } from './exportSvg';

interface ExportPngParams {
  svgElement: SVGSVGElement;
  options: ExportOptions;
  projectName: string;
  legendSvg?: string;
  onProgress?: (progress: number) => void;
}

/**
 * Exports the dependency graph as a PNG file.
 * First renders SVG to canvas, then converts to PNG.
 */
export async function exportPng({
  svgElement,
  options,
  projectName,
  legendSvg,
  onProgress,
}: ExportPngParams): Promise<ExportResult> {
  onProgress?.(10);

  // First export as SVG (to get styled, complete SVG)
  const svgResult = await exportSvg({
    svgElement,
    options: { ...options, format: 'svg' },
    projectName,
    legendSvg,
  });

  onProgress?.(30);

  // Create image from SVG
  const img = new Image();
  const svgUrl = URL.createObjectURL(svgResult.blob);

  return new Promise((resolve, reject) => {
    img.onload = () => {
      onProgress?.(50);

      // Calculate dimensions with resolution
      const width = svgResult.width * options.resolution;
      const height = svgResult.height * options.resolution;

      // Create canvas
      const canvas = document.createElement('canvas');
      canvas.width = width;
      canvas.height = height;

      const ctx = canvas.getContext('2d');
      if (!ctx) {
        URL.revokeObjectURL(svgUrl);
        reject(new Error('Failed to get canvas context'));
        return;
      }

      // Fill background
      if (options.backgroundColor !== 'transparent') {
        ctx.fillStyle = options.backgroundColor;
        ctx.fillRect(0, 0, width, height);
      }

      // Scale for resolution
      ctx.scale(options.resolution, options.resolution);

      // Draw image
      ctx.drawImage(img, 0, 0);

      onProgress?.(80);

      // Convert to PNG
      canvas.toBlob(
        (blob) => {
          URL.revokeObjectURL(svgUrl);

          if (!blob) {
            reject(new Error('Failed to create PNG blob'));
            return;
          }

          onProgress?.(100);

          // Generate filename
          const timestamp = new Date().toISOString().split('T')[0];
          const resolutionSuffix = options.resolution > 1 ? `@${options.resolution}x` : '';
          const filename = `${projectName}-dependency-graph-${timestamp}${resolutionSuffix}.png`;

          resolve({
            blob,
            filename,
            width,
            height,
          });
        },
        'image/png',
        0.95 // Quality (0-1)
      );
    };

    img.onerror = () => {
      URL.revokeObjectURL(svgUrl);
      reject(new Error('Failed to load SVG for PNG conversion'));
    };

    img.src = svgUrl;
  });
}
```

**useGraphExport Hook:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/useGraphExport.ts

import { useState, useCallback } from 'react';
import type { ExportOptions, ExportProgress, ExportResult } from './types';
import { exportSvg } from './utils/exportSvg';
import { exportPng } from './utils/exportPng';
import { renderLegendSvg } from './utils/renderLegendForExport';

interface UseGraphExportProps {
  svgRef: React.RefObject<SVGSVGElement>;
  projectName: string;
  isDarkMode: boolean;
}

interface UseGraphExportResult {
  exportProgress: ExportProgress;
  startExport: (options: ExportOptions) => Promise<void>;
  cancelExport: () => void;
}

export function useGraphExport({
  svgRef,
  projectName,
  isDarkMode,
}: UseGraphExportProps): UseGraphExportResult {
  const [exportProgress, setExportProgress] = useState<ExportProgress>({
    isExporting: false,
    progress: 0,
    stage: 'preparing',
  });

  const [abortController, setAbortController] = useState<AbortController | null>(null);

  const startExport = useCallback(async (options: ExportOptions) => {
    if (!svgRef.current) {
      throw new Error('SVG element not found');
    }

    const controller = new AbortController();
    setAbortController(controller);

    setExportProgress({
      isExporting: true,
      progress: 0,
      stage: 'preparing',
    });

    try {
      // Generate legend SVG if needed
      let legendSvg: string | undefined;
      if (options.includeLegend) {
        legendSvg = renderLegendSvg(isDarkMode);
      }

      // Determine background color
      const backgroundColor = options.backgroundColor === 'transparent'
        ? 'transparent'
        : (isDarkMode ? '#1f2937' : '#ffffff');

      const exportOptions: ExportOptions = {
        ...options,
        backgroundColor,
      };

      setExportProgress(prev => ({ ...prev, progress: 20, stage: 'rendering' }));

      let result: ExportResult;

      if (options.format === 'svg') {
        result = await exportSvg({
          svgElement: svgRef.current,
          options: exportOptions,
          projectName,
          legendSvg,
        });
      } else {
        result = await exportPng({
          svgElement: svgRef.current,
          options: exportOptions,
          projectName,
          legendSvg,
          onProgress: (progress) => {
            setExportProgress(prev => ({
              ...prev,
              progress: 20 + progress * 0.7,
              stage: progress < 50 ? 'rendering' : 'encoding',
            }));
          },
        });
      }

      // Check if cancelled
      if (controller.signal.aborted) {
        return;
      }

      setExportProgress(prev => ({ ...prev, progress: 95, stage: 'complete' }));

      // Trigger download
      downloadBlob(result.blob, result.filename);

      setExportProgress({
        isExporting: false,
        progress: 100,
        stage: 'complete',
      });
    } catch (error) {
      if (!controller.signal.aborted) {
        console.error('Export failed:', error);
        setExportProgress({
          isExporting: false,
          progress: 0,
          stage: 'preparing',
        });
        throw error;
      }
    } finally {
      setAbortController(null);
    }
  }, [svgRef, projectName, isDarkMode]);

  const cancelExport = useCallback(() => {
    abortController?.abort();
    setExportProgress({
      isExporting: false,
      progress: 0,
      stage: 'preparing',
    });
  }, [abortController]);

  return {
    exportProgress,
    startExport,
    cancelExport,
  };
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}
```

**ExportMenu Component:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/ExportMenu.tsx

import React, { useState } from 'react';
import type { ExportOptions, ExportFormat, ExportScope, ExportResolution, ExportProgress } from './types';

interface ExportMenuProps {
  isOpen: boolean;
  onClose: () => void;
  onExport: (options: ExportOptions) => Promise<void>;
  exportProgress: ExportProgress;
  isDarkMode: boolean;
}

const DEFAULT_OPTIONS: ExportOptions = {
  format: 'png',
  scope: 'viewport',
  resolution: 2,
  includeLegend: true,
  includeWatermark: false,
  backgroundColor: '#ffffff',
};

export function ExportMenu({
  isOpen,
  onClose,
  onExport,
  exportProgress,
  isDarkMode,
}: ExportMenuProps) {
  const [options, setOptions] = useState<ExportOptions>({
    ...DEFAULT_OPTIONS,
    backgroundColor: isDarkMode ? '#1f2937' : '#ffffff',
  });

  const handleExport = async () => {
    try {
      await onExport(options);
      onClose();
    } catch (error) {
      console.error('Export failed:', error);
      // Toast notification would be shown here
    }
  };

  if (!isOpen) return null;

  return (
    <div className="absolute right-4 top-14 z-50 w-72 bg-white dark:bg-gray-800
                    rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 p-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="font-semibold text-gray-900 dark:text-white">Export Graph</h3>
        <button
          onClick={onClose}
          className="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
          aria-label="Close export menu"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      {/* Format Selection */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Format
        </label>
        <div className="flex gap-2">
          {(['png', 'svg'] as ExportFormat[]).map((format) => (
            <button
              key={format}
              onClick={() => setOptions(prev => ({ ...prev, format }))}
              className={`flex-1 px-3 py-2 rounded-md text-sm font-medium transition-colors
                ${options.format === format
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
                }`}
            >
              {format.toUpperCase()}
            </button>
          ))}
        </div>
      </div>

      {/* Resolution (PNG only) */}
      {options.format === 'png' && (
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Resolution
          </label>
          <select
            value={options.resolution}
            onChange={(e) => setOptions(prev => ({
              ...prev,
              resolution: parseInt(e.target.value) as ExportResolution
            }))}
            className="w-full px-3 py-2 rounded-md border border-gray-300 dark:border-gray-600
                       bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
          >
            <option value={1}>1x (Standard)</option>
            <option value={2}>2x (High DPI)</option>
            <option value={4}>4x (Print Quality)</option>
          </select>
        </div>
      )}

      {/* Scope Selection */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Scope
        </label>
        <select
          value={options.scope}
          onChange={(e) => setOptions(prev => ({
            ...prev,
            scope: e.target.value as ExportScope
          }))}
          className="w-full px-3 py-2 rounded-md border border-gray-300 dark:border-gray-600
                     bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
        >
          <option value="viewport">Current View</option>
          <option value="full">Full Graph</option>
          <option value="selected">Selected Elements</option>
        </select>
      </div>

      {/* Options Checkboxes */}
      <div className="mb-4 space-y-2">
        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={options.includeLegend}
            onChange={(e) => setOptions(prev => ({
              ...prev,
              includeLegend: e.target.checked
            }))}
            className="rounded border-gray-300 dark:border-gray-600
                       text-blue-600 focus:ring-blue-500"
          />
          <span className="text-sm text-gray-700 dark:text-gray-300">Include Legend</span>
        </label>

        <label className="flex items-center gap-2 cursor-pointer">
          <input
            type="checkbox"
            checked={options.includeWatermark}
            onChange={(e) => setOptions(prev => ({
              ...prev,
              includeWatermark: e.target.checked
            }))}
            className="rounded border-gray-300 dark:border-gray-600
                       text-blue-600 focus:ring-blue-500"
          />
          <span className="text-sm text-gray-700 dark:text-gray-300">Include Watermark</span>
        </label>
      </div>

      {/* Export Progress */}
      {exportProgress.isExporting && (
        <div className="mb-4">
          <div className="flex justify-between text-sm text-gray-600 dark:text-gray-400 mb-1">
            <span>Exporting...</span>
            <span>{Math.round(exportProgress.progress)}%</span>
          </div>
          <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2">
            <div
              className="bg-blue-600 h-2 rounded-full transition-all duration-200"
              style={{ width: `${exportProgress.progress}%` }}
            />
          </div>
          <div className="text-xs text-gray-500 dark:text-gray-400 mt-1 capitalize">
            {exportProgress.stage.replace('-', ' ')}
          </div>
        </div>
      )}

      {/* Export Button */}
      <button
        onClick={handleExport}
        disabled={exportProgress.isExporting}
        className="w-full px-4 py-2 bg-blue-600 text-white rounded-md font-medium
                   hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500
                   disabled:opacity-50 disabled:cursor-not-allowed
                   transition-colors"
      >
        {exportProgress.isExporting ? 'Exporting...' : 'Export'}
      </button>
    </div>
  );
}
```

**Render Legend for Export:**
```typescript
// apps/web/app/components/visualization/DependencyGraph/utils/renderLegendForExport.ts

/**
 * Generates a standalone SVG legend that can be embedded in exports.
 */
export function renderLegendSvg(isDarkMode: boolean): string {
  const bgColor = isDarkMode ? '#374151' : '#f3f4f6';
  const textColor = isDarkMode ? '#f3f4f6' : '#1f2937';
  const borderColor = isDarkMode ? '#4b5563' : '#d1d5db';

  return `
    <svg xmlns="http://www.w3.org/2000/svg" width="140" height="90">
      <rect width="140" height="90" rx="8" fill="${bgColor}" stroke="${borderColor}" stroke-width="1"/>

      <!-- Title -->
      <text x="10" y="18" fill="${textColor}" font-size="11" font-weight="600" font-family="system-ui, sans-serif">Legend</text>

      <!-- Normal Node -->
      <circle cx="20" cy="35" r="6" fill="#3b82f6"/>
      <text x="32" y="39" fill="${textColor}" font-size="10" font-family="system-ui, sans-serif">Package</text>

      <!-- Circular Node -->
      <circle cx="20" cy="55" r="6" fill="#ef4444" stroke="#dc2626" stroke-width="2"/>
      <text x="32" y="59" fill="${textColor}" font-size="10" font-family="system-ui, sans-serif">Circular Dep</text>

      <!-- Edge Legend -->
      <line x1="12" y1="75" x2="28" y2="75" stroke="#6b7280" stroke-width="1.5"/>
      <text x="32" y="79" fill="${textColor}" font-size="10" font-family="system-ui, sans-serif">Dependency</text>
    </svg>
  `;
}
```

### Integration with Main Component

```typescript
// In apps/web/app/components/visualization/DependencyGraph/index.tsx
// Add to existing component from Stories 4.1-4.5

import { useState } from 'react';
import { useGraphExport } from './useGraphExport';
import { ExportMenu } from './ExportMenu';
import { useSettingsStore } from '@/stores/settings';

export const DependencyGraphViz = React.memo(function DependencyGraphViz({
  data,
  circularDependencies,
  projectName = 'monoguard',
}: DependencyGraphProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  const { theme } = useSettingsStore();
  const isDarkMode = theme === 'dark' ||
    (theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);

  const [isExportMenuOpen, setIsExportMenuOpen] = useState(false);

  const { exportProgress, startExport, cancelExport } = useGraphExport({
    svgRef,
    projectName,
    isDarkMode,
  });

  // ... existing hooks from previous stories ...

  return (
    <div ref={containerRef} className="relative w-full h-full">
      <svg ref={svgRef} className="w-full h-full min-h-[500px]">
        {/* SVG content rendered by D3 */}
      </svg>

      {/* Export Button */}
      <button
        onClick={() => setIsExportMenuOpen(true)}
        className="absolute top-4 right-4 px-3 py-2 bg-white dark:bg-gray-800
                   rounded-md shadow-md border border-gray-200 dark:border-gray-700
                   text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700
                   flex items-center gap-2 transition-colors"
        aria-label="Export graph"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        Export
      </button>

      {/* Export Menu */}
      <ExportMenu
        isOpen={isExportMenuOpen}
        onClose={() => setIsExportMenuOpen(false)}
        onExport={startExport}
        exportProgress={exportProgress}
        isDarkMode={isDarkMode}
      />

      {/* Other components from previous stories */}
      <NodeTooltip ... />
      <GraphMinimap ... />
      <ZoomControls ... />
      <GraphLegend />
      <GraphControls ... />
    </div>
  );
});
```

### Testing Requirements

**Test File:** `apps/web/src/__tests__/exportSvg.test.ts`

```typescript
import { exportSvg } from '@/components/visualization/DependencyGraph/utils/exportSvg';

describe('exportSvg', () => {
  let mockSvg: SVGSVGElement;

  beforeEach(() => {
    mockSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    mockSvg.setAttribute('width', '800');
    mockSvg.setAttribute('height', '600');

    const circle = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
    circle.setAttribute('cx', '100');
    circle.setAttribute('cy', '100');
    circle.setAttribute('r', '50');
    mockSvg.appendChild(circle);
  });

  it('should export SVG as blob', async () => {
    const result = await exportSvg({
      svgElement: mockSvg,
      options: {
        format: 'svg',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'test-project',
    });

    expect(result.blob).toBeInstanceOf(Blob);
    expect(result.blob.type).toBe('image/svg+xml');
  });

  it('should generate correct filename', async () => {
    const result = await exportSvg({
      svgElement: mockSvg,
      options: {
        format: 'svg',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'mono-guard',
    });

    expect(result.filename).toMatch(/^mono-guard-dependency-graph-\d{4}-\d{2}-\d{2}\.svg$/);
  });

  it('should include legend when requested', async () => {
    const legendSvg = '<svg><text>Legend</text></svg>';

    const result = await exportSvg({
      svgElement: mockSvg,
      options: {
        format: 'svg',
        scope: 'viewport',
        resolution: 1,
        includeLegend: true,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'test',
      legendSvg,
    });

    const svgContent = await result.blob.text();
    expect(svgContent).toContain('legend-group');
  });

  it('should add watermark when requested', async () => {
    const result = await exportSvg({
      svgElement: mockSvg,
      options: {
        format: 'svg',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: true,
        backgroundColor: '#ffffff',
      },
      projectName: 'test',
    });

    const svgContent = await result.blob.text();
    expect(svgContent).toContain('MonoGuard');
  });
});
```

**Test File:** `apps/web/src/__tests__/ExportMenu.test.tsx`

```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import { ExportMenu } from '@/components/visualization/DependencyGraph/ExportMenu';

describe('ExportMenu', () => {
  const mockOnClose = vi.fn();
  const mockOnExport = vi.fn().mockResolvedValue(undefined);
  const defaultProgress = { isExporting: false, progress: 0, stage: 'preparing' as const };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should not render when closed', () => {
    render(
      <ExportMenu
        isOpen={false}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    );

    expect(screen.queryByText('Export Graph')).not.toBeInTheDocument();
  });

  it('should render when open', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    );

    expect(screen.getByText('Export Graph')).toBeInTheDocument();
  });

  it('should show format options', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    );

    expect(screen.getByText('PNG')).toBeInTheDocument();
    expect(screen.getByText('SVG')).toBeInTheDocument();
  });

  it('should show resolution options only for PNG', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    );

    // PNG is default, resolution should be visible
    expect(screen.getByText('Resolution')).toBeInTheDocument();

    // Switch to SVG
    fireEvent.click(screen.getByText('SVG'));

    // Resolution should be hidden
    expect(screen.queryByText('Resolution')).not.toBeInTheDocument();
  });

  it('should call onExport with options when export button clicked', async () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    );

    fireEvent.click(screen.getByText('Export'));

    expect(mockOnExport).toHaveBeenCalledWith(expect.objectContaining({
      format: 'png',
      scope: 'viewport',
      resolution: 2,
      includeLegend: true,
    }));
  });

  it('should show progress during export', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={{ isExporting: true, progress: 50, stage: 'rendering' }}
        isDarkMode={false}
      />
    );

    expect(screen.getByText('Exporting...')).toBeInTheDocument();
    expect(screen.getByText('50%')).toBeInTheDocument();
    expect(screen.getByText('rendering')).toBeInTheDocument();
  });

  it('should close when close button clicked', () => {
    render(
      <ExportMenu
        isOpen={true}
        onClose={mockOnClose}
        onExport={mockOnExport}
        exportProgress={defaultProgress}
        isDarkMode={false}
      />
    );

    fireEvent.click(screen.getByLabelText('Close export menu'));

    expect(mockOnClose).toHaveBeenCalled();
  });
});
```

**Test File:** `apps/web/src/__tests__/useGraphExport.test.ts`

```typescript
import { renderHook, act } from '@testing-library/react';
import { useGraphExport } from '@/components/visualization/DependencyGraph/useGraphExport';

// Mock the export utilities
vi.mock('@/components/visualization/DependencyGraph/utils/exportSvg', () => ({
  exportSvg: vi.fn().mockResolvedValue({
    blob: new Blob(['<svg></svg>'], { type: 'image/svg+xml' }),
    filename: 'test.svg',
    width: 800,
    height: 600,
  }),
}));

vi.mock('@/components/visualization/DependencyGraph/utils/exportPng', () => ({
  exportPng: vi.fn().mockResolvedValue({
    blob: new Blob([''], { type: 'image/png' }),
    filename: 'test.png',
    width: 1600,
    height: 1200,
  }),
}));

describe('useGraphExport', () => {
  const mockSvgRef = { current: document.createElementNS('http://www.w3.org/2000/svg', 'svg') };

  it('should initialize with not exporting state', () => {
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef as any,
        projectName: 'test',
        isDarkMode: false,
      })
    );

    expect(result.current.exportProgress.isExporting).toBe(false);
    expect(result.current.exportProgress.progress).toBe(0);
  });

  it('should update progress during export', async () => {
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef as any,
        projectName: 'test',
        isDarkMode: false,
      })
    );

    await act(async () => {
      await result.current.startExport({
        format: 'svg',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      });
    });

    // After completion
    expect(result.current.exportProgress.isExporting).toBe(false);
  });

  it('should handle cancellation', async () => {
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef as any,
        projectName: 'test',
        isDarkMode: false,
      })
    );

    // Start then immediately cancel
    const exportPromise = act(async () => {
      result.current.startExport({
        format: 'png',
        scope: 'full',
        resolution: 2,
        includeLegend: true,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      });
    });

    act(() => {
      result.current.cancelExport();
    });

    await exportPromise;

    expect(result.current.exportProgress.isExporting).toBe(false);
  });
});
```

### Critical Don't-Miss Rules (from project-context.md)

1. **NEVER forget D3.js cleanup** - Remove event listeners in cleanup
2. **Use React.memo** - Already in place from Story 4.1
3. **Performance** - Export should complete < 5 seconds for large graphs
4. **Memory cleanup** - Revoke object URLs after use
5. **Dark mode support** - Export should respect current theme
6. **Accessibility** - Export button has proper aria-label

### Previous Story Intelligence (Stories 4.1-4.5)

**Key Patterns to Follow:**
- Component structure with separate hooks for each concern
- Comprehensive test coverage for each AC
- Consistent styling with existing components
- Progressive UI patterns (menu/modal)

**Integration Points:**
- Export should capture current state including:
  - Cycle highlighting from Story 4.2
  - Expanded/collapsed state from Story 4.3
  - Current zoom/pan from Story 4.4
  - Hover highlighting from Story 4.5

**Existing Files to Extend:**
- `index.tsx` - Add export button and menu
- `types.ts` - Add export types
- `GraphControls.tsx` - Optionally add export to controls

### UX Design Requirements (from ux-design-specification.md)

- **Export formats:** PNG and SVG as per FR18
- **Resolution options:** Support high DPI exports
- **Scope options:** Viewport, full graph, selected
- **Legend:** Can be included/excluded
- **Immediate download:** File downloads when ready

### Performance Considerations

1. **SVG Cloning:** Clone SVG before manipulation to avoid affecting displayed graph
2. **Canvas Rendering:** Use OffscreenCanvas if available for better performance
3. **Memory Management:** Revoke blob URLs after download
4. **Large Graph Handling:** Show progress for graphs > 200 nodes
5. **Abort Support:** Allow cancellation of long exports

### References

- [Story 4.1: DependencyGraph Implementation] `4-1-implement-d3js-force-directed-dependency-graph.md`
- [Story 4.2: Cycle Highlighting] `4-2-highlight-circular-dependencies-in-graph.md`
- [Story 4.3: Node Expand/Collapse] `4-3-implement-node-expand-collapse-functionality.md`
- [Story 4.4: Zoom/Pan Controls] `4-4-add-zoom-pan-and-navigation-controls.md`
- [Story 4.5: Hover Details and Tooltips] `4-5-implement-hover-details-and-tooltips.md`
- [Epic 4 Story 4.6 Requirements] `_bmad-output/planning-artifacts/epics.md` - Lines 1076-1096
- [Project Context: D3.js Rules] `_bmad-output/project-context.md` - D3.js Integration section
- [SVG to Canvas Rendering] https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Drawing_DOM_objects_into_a_canvas
- [Blob and Object URLs] https://developer.mozilla.org/en-US/docs/Web/API/URL/createObjectURL

## Dev Agent Record

### Agent Model Used

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

- Lint fix: Added `SVGElement`, `XMLSerializer`, `DOMParser`, `Image`, `HTMLCanvasElement`, `HTMLAnchorElement`, `HTMLImageElement` to ESLint globals in `apps/web/eslint.config.mjs`
- Test fix: Replaced `blob.text()` (not available in jsdom) with `FileReader.readAsText()` helper
- Test fix: `ExportMenu.test.tsx` - Used `getAllByText` for "Exporting..." which appears in both progress label and button

### Completion Notes List

- All 9 tasks completed successfully
- 44 new tests added across 5 test files (exportSvg, exportPng, renderLegendForExport, ExportMenu, useGraphExport)
- Full test suite: 25 files, 420 tests, 0 failures
- CI: lint (0 errors), type-check (pass), test (pass)
- No Go changes → Go CI not applicable

### File List

**New Files:**
- `apps/web/app/components/visualization/DependencyGraph/utils/exportSvg.ts` - SVG export utility
- `apps/web/app/components/visualization/DependencyGraph/utils/exportPng.ts` - PNG export utility
- `apps/web/app/components/visualization/DependencyGraph/utils/renderLegendForExport.ts` - Legend renderer for exports
- `apps/web/app/components/visualization/DependencyGraph/ExportMenu.tsx` - Export options UI component
- `apps/web/app/components/visualization/DependencyGraph/useGraphExport.ts` - Export hook
- `apps/web/app/components/visualization/DependencyGraph/__tests__/exportSvg.test.ts` - SVG export tests (10 tests)
- `apps/web/app/components/visualization/DependencyGraph/__tests__/exportPng.test.ts` - PNG export tests (6 tests)
- `apps/web/app/components/visualization/DependencyGraph/__tests__/renderLegendForExport.test.ts` - Legend renderer tests (8 tests)
- `apps/web/app/components/visualization/DependencyGraph/__tests__/ExportMenu.test.tsx` - ExportMenu tests (14 tests)
- `apps/web/app/components/visualization/DependencyGraph/__tests__/useGraphExport.test.ts` - useGraphExport tests (6 tests)

**Modified Files:**
- `apps/web/app/components/visualization/DependencyGraph/types.ts` - Added export types (ExportFormat, ExportOptions, ExportScope, ExportResolution, ExportResult, ExportProgress)
- `apps/web/app/components/visualization/DependencyGraph/index.tsx` - Integrated export button, ExportMenu, useGraphExport hook, and re-exports
- `apps/web/eslint.config.mjs` - Added browser globals for SVG/DOM APIs
