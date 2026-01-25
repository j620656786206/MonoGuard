/**
 * Tests for DependencyGraph component
 *
 * Following red-green-refactor cycle.
 */

import type { DependencyGraph } from '@monoguard/types'
import { render, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Import component (will fail until implemented - RED phase)
import { DependencyGraphViz } from '../index'

// Mock ResizeObserver for tests
beforeEach(() => {
  global.ResizeObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
  }))
})

afterEach(() => {
  vi.restoreAllMocks()
})

/**
 * Mock data for testing
 */
const createMockData = (nodeCount: number = 2): DependencyGraph => {
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

const mockData: DependencyGraph = {
  nodes: {
    '@app/core': {
      name: '@app/core',
      version: '1.0.0',
      path: 'packages/core',
      dependencies: ['@app/utils'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@app/utils': {
      name: '@app/utils',
      version: '1.0.0',
      path: 'packages/utils',
      dependencies: [],
      devDependencies: [],
      peerDependencies: [],
    },
  },
  edges: [
    {
      from: '@app/core',
      to: '@app/utils',
      type: 'production',
      versionRange: '^1.0.0',
    },
  ],
  rootPath: '/workspace',
  workspaceType: 'npm',
}

describe('DependencyGraphViz', () => {
  describe('AC4: Data Integration', () => {
    it('should render SVG element', () => {
      render(<DependencyGraphViz data={mockData} />)
      expect(document.querySelector('svg')).toBeInTheDocument()
    })

    it('should render correct number of nodes', async () => {
      render(<DependencyGraphViz data={mockData} />)

      // Wait for D3 to render nodes
      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          expect(circles.length).toBe(2)
        },
        { timeout: 1000 }
      )
    })

    it('should render correct number of links', async () => {
      render(<DependencyGraphViz data={mockData} />)

      // Wait for D3 to render links
      await waitFor(
        () => {
          const lines = document.querySelectorAll('line')
          expect(lines.length).toBe(1)
        },
        { timeout: 1000 }
      )
    })

    it('should render node labels', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          const texts = document.querySelectorAll('text')
          expect(texts.length).toBeGreaterThanOrEqual(2)
        },
        { timeout: 1000 }
      )
    })
  })

  describe('AC5: Node Visual Representation', () => {
    it('should display truncated package names', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          const texts = document.querySelectorAll('text')
          // Should show the last part of the package name (e.g., "core" instead of "@app/core")
          const textContents = Array.from(texts).map((t) => t.textContent)
          expect(textContents.some((t) => t?.includes('core'))).toBe(true)
        },
        { timeout: 1000 }
      )
    })

    it('should render nodes with consistent styling', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          circles.forEach((circle) => {
            expect(circle.getAttribute('r')).toBeTruthy()
            expect(circle.getAttribute('fill')).toBeTruthy()
          })
        },
        { timeout: 1000 }
      )
    })
  })

  describe('AC6: Edge Visual Representation', () => {
    it('should render edges with arrow markers', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          // Check for arrowhead marker definition
          const marker = document.querySelector('marker#arrowhead')
          expect(marker).toBeInTheDocument()

          // Check that lines use the marker
          const lines = document.querySelectorAll('line')
          lines.forEach((line) => {
            expect(line.getAttribute('marker-end')).toContain('arrowhead')
          })
        },
        { timeout: 1000 }
      )
    })

    it('should render edges with visible styling', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          const lines = document.querySelectorAll('line')
          lines.forEach((line) => {
            expect(line.getAttribute('stroke')).toBeTruthy()
            expect(line.getAttribute('stroke-width')).toBeTruthy()
          })
        },
        { timeout: 1000 }
      )
    })
  })

  describe('AC1: Force-Directed Layout', () => {
    it('should render with force-directed positions', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          // Nodes are positioned via transform on parent g element
          const nodeGroups = document.querySelectorAll('g.node')
          nodeGroups.forEach((nodeGroup) => {
            const transform = nodeGroup.getAttribute('transform')
            expect(transform).toBeTruthy()
            // Should have translate with non-zero values
            expect(transform).toMatch(/translate\(\d+\.?\d*,\d+\.?\d*\)/)
          })
        },
        { timeout: 2000 }
      )
    })
  })

  describe('AC3: Responsive Container', () => {
    it('should have responsive container classes', () => {
      render(<DependencyGraphViz data={mockData} />)
      const svg = document.querySelector('svg')
      expect(svg).toHaveClass('w-full')
    })

    it('should accept custom className', () => {
      render(<DependencyGraphViz data={mockData} className="custom-class" />)
      const container = document.querySelector('.custom-class')
      expect(container).toBeInTheDocument()
    })
  })

  describe('Data transformation', () => {
    it('should handle empty graph data', () => {
      const emptyData: DependencyGraph = {
        nodes: {},
        edges: [],
        rootPath: '/workspace',
        workspaceType: 'npm',
      }

      render(<DependencyGraphViz data={emptyData} />)
      const svg = document.querySelector('svg')
      expect(svg).toBeInTheDocument()
    })

    it('should handle graph with many nodes', async () => {
      const largeData = createMockData(50)

      render(<DependencyGraphViz data={largeData} />)

      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          expect(circles.length).toBe(50)
        },
        { timeout: 2000 }
      )
    })
  })

  describe('Cleanup and memory management', () => {
    it('should clean up on unmount', async () => {
      const { unmount } = render(<DependencyGraphViz data={mockData} />)

      await waitFor(() => {
        expect(document.querySelector('svg')).toBeInTheDocument()
      })

      // Unmount should clean up D3 elements
      unmount()

      // SVG should be removed from DOM
      expect(document.querySelector('svg')).not.toBeInTheDocument()
    })
  })
})
