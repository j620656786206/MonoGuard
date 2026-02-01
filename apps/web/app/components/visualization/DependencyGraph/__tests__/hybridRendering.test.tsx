/**
 * Integration tests for Hybrid SVG/Canvas Rendering (Story 4.9)
 *
 * Validates that the DependencyGraphViz component correctly switches between
 * SVG and Canvas rendering modes, and that all overlay features work in both modes.
 *
 * @see AC8: Feature Parity Verification
 */

import type { DependencyGraph } from '@monoguard/types'
import { render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useSettingsStore } from '../../../../stores/settings'
import { DependencyGraphViz } from '../index'

// Mock ResizeObserver
beforeEach(() => {
  global.ResizeObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
  }))

  // Mock canvas context for Canvas mode
  HTMLCanvasElement.prototype.getContext = vi.fn(() => ({
    save: vi.fn(),
    restore: vi.fn(),
    clearRect: vi.fn(),
    setTransform: vi.fn(),
    translate: vi.fn(),
    scale: vi.fn(),
    beginPath: vi.fn(),
    moveTo: vi.fn(),
    lineTo: vi.fn(),
    arc: vi.fn(),
    closePath: vi.fn(),
    fill: vi.fn(),
    stroke: vi.fn(),
    fillText: vi.fn(),
    rotate: vi.fn(),
    fillStyle: '',
    strokeStyle: '',
    lineWidth: 1,
    globalAlpha: 1,
    font: '',
    textAlign: '',
  })) as unknown as typeof HTMLCanvasElement.prototype.getContext

  vi.stubGlobal('devicePixelRatio', 1)

  // Reset store to auto
  useSettingsStore.setState({ visualizationMode: 'auto' })
})

afterEach(() => {
  vi.restoreAllMocks()
  vi.unstubAllGlobals()
  useSettingsStore.setState({ visualizationMode: 'auto' })
})

/**
 * Create mock graph data with specified node count
 */
const createMockData = (nodeCount: number): DependencyGraph => {
  const nodes: Record<string, DependencyGraph['nodes'][string]> = {}
  const edges: DependencyGraph['edges'] = []

  for (let i = 0; i < nodeCount; i++) {
    const name = `@app/package-${i}`
    nodes[name] = {
      name,
      version: '1.0.0',
      path: `packages/package-${i}`,
      dependencies: i > 0 ? [`@app/package-${i - 1}`] : [],
      devDependencies: [],
      peerDependencies: [],
    }

    if (i > 0) {
      edges.push({
        from: name,
        to: `@app/package-${i - 1}`,
        type: 'production',
        versionRange: '^1.0.0',
      })
    }
  }

  return {
    nodes,
    edges,
    rootPath: '/workspace',
    workspaceType: 'npm',
  }
}

describe('Hybrid Rendering Integration (Story 4.9)', () => {
  describe('AC1: Automatic Mode Selection', () => {
    it('should render SVG mode for small graphs (< 500 nodes)', () => {
      const data = createMockData(10)
      const { container } = render(<DependencyGraphViz data={data} />)

      expect(container.querySelector('svg')).toBeTruthy()
      expect(container.querySelector('canvas')).toBeNull()
    })

    it('should render Canvas mode for large graphs (>= 500 nodes)', () => {
      const data = createMockData(500)
      const { container } = render(<DependencyGraphViz data={data} />)

      expect(container.querySelector('canvas')).toBeTruthy()
      // The main graph SVG (h-full w-full) should not exist; icon SVGs in overlays are expected
      expect(container.querySelector('svg.h-full')).toBeNull()
    })
  })

  describe('AC2: Mode Indicator', () => {
    it('should show SVG mode indicator for small graphs', () => {
      const data = createMockData(10)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByText('SVG mode')).toBeTruthy()
      expect(screen.getByText('10 nodes')).toBeTruthy()
    })

    it('should show Canvas mode indicator for large graphs', () => {
      const data = createMockData(500)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByText('CANVAS mode')).toBeTruthy()
    })

    it('should show "Forced" badge when user forces a mode', () => {
      useSettingsStore.setState({ visualizationMode: 'force-svg' })
      const data = createMockData(10)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByText('Forced')).toBeTruthy()
    })
  })

  describe('AC3: User Override', () => {
    it('should force SVG mode when user selects force-svg', () => {
      useSettingsStore.setState({ visualizationMode: 'force-svg' })
      const data = createMockData(600) // Would be Canvas in auto mode
      const { container } = render(<DependencyGraphViz data={data} />)

      expect(container.querySelector('svg')).toBeTruthy()
      expect(container.querySelector('canvas')).toBeNull()
    })

    it('should force Canvas mode when user selects force-canvas', () => {
      useSettingsStore.setState({ visualizationMode: 'force-canvas' })
      const data = createMockData(10) // Would be SVG in auto mode
      const { container } = render(<DependencyGraphViz data={data} />)

      expect(container.querySelector('canvas')).toBeTruthy()
      // The main graph SVG (h-full w-full) should not exist; icon SVGs in overlays are expected
      expect(container.querySelector('svg.h-full')).toBeNull()
    })
  })

  describe('AC8: Feature Parity - Shared Overlay Components', () => {
    it('should render GraphControls in SVG mode', () => {
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      // GraphControls renders "All" and level buttons
      expect(screen.getByText('All')).toBeTruthy()
    })

    it('should render GraphControls in Canvas mode', () => {
      useSettingsStore.setState({ visualizationMode: 'force-canvas' })
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByText('All')).toBeTruthy()
    })

    it('should render ZoomControls in SVG mode', () => {
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByLabelText('Zoom in')).toBeTruthy()
      expect(screen.getByLabelText('Zoom out')).toBeTruthy()
    })

    it('should render ZoomControls in Canvas mode', () => {
      useSettingsStore.setState({ visualizationMode: 'force-canvas' })
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByLabelText('Zoom in')).toBeTruthy()
      expect(screen.getByLabelText('Zoom out')).toBeTruthy()
    })

    it('should render GraphLegend in SVG mode', () => {
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByText('Normal Package')).toBeTruthy()
    })

    it('should render GraphLegend in Canvas mode', () => {
      useSettingsStore.setState({ visualizationMode: 'force-canvas' })
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByText('Normal Package')).toBeTruthy()
    })

    it('should render Export button in SVG mode', () => {
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByLabelText('Export graph')).toBeTruthy()
    })

    it('should render Export button in Canvas mode', () => {
      useSettingsStore.setState({ visualizationMode: 'force-canvas' })
      const data = createMockData(5)
      render(<DependencyGraphViz data={data} />)

      expect(screen.getByLabelText('Export graph')).toBeTruthy()
    })
  })

  describe('Empty data handling', () => {
    it('should not render mode indicator for empty data', () => {
      const data: DependencyGraph = {
        nodes: {},
        edges: [],
        rootPath: '/workspace',
        workspaceType: 'npm',
      }
      render(<DependencyGraphViz data={data} />)

      expect(screen.queryByText('SVG mode')).toBeNull()
      expect(screen.queryByText('CANVAS mode')).toBeNull()
    })
  })
})
