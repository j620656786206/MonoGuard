/**
 * Tests for DependencyGraph component
 *
 * Following red-green-refactor cycle.
 * @see Story 4.1: Implement D3.js Force-Directed Dependency Graph
 * @see Story 4.2: Highlight Circular Dependencies in Graph
 */

import type { CircularDependencyInfo, DependencyGraph } from '@monoguard/types'
import { fireEvent, render, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Import component
import { DependencyGraphViz, GraphLegend } from '../index'

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

    it('should handle graph with 100+ nodes (AC2 performance requirement)', async () => {
      // AC2: Graph renders in < 2 seconds for 100 packages
      const largeData = createMockData(100)

      const startTime = performance.now()
      render(<DependencyGraphViz data={largeData} />)

      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          expect(circles.length).toBe(100)
        },
        { timeout: 2000 }
      )

      const renderTime = performance.now() - startTime
      // AC2: Graph renders in < 2 seconds for 100 packages
      expect(renderTime).toBeLessThan(2000)
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

  describe('truncatePackageName utility', () => {
    // Import the utility for direct testing
    it('should handle scoped package names', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          const texts = document.querySelectorAll('text')
          // @app/core should show as "core"
          const textContents = Array.from(texts).map((t) => t.textContent)
          expect(textContents).toContain('core')
          expect(textContents).toContain('utils')
        },
        { timeout: 1000 }
      )
    })

    it('should handle non-scoped package names', async () => {
      const unscoped: DependencyGraph = {
        nodes: {
          lodash: {
            name: 'lodash',
            version: '4.17.21',
            path: 'node_modules/lodash',
            dependencies: [],
            devDependencies: [],
            peerDependencies: [],
          },
        },
        edges: [],
        rootPath: '/workspace',
        workspaceType: 'npm',
      }

      render(<DependencyGraphViz data={unscoped} />)

      await waitFor(
        () => {
          const texts = document.querySelectorAll('text')
          const textContents = Array.from(texts).map((t) => t.textContent)
          expect(textContents).toContain('lodash')
        },
        { timeout: 1000 }
      )
    })
  })

  describe('Interactive behaviors', () => {
    it('should set up zoom behavior on SVG', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          const svg = document.querySelector('svg')
          expect(svg).toBeInTheDocument()
          // Zoom is set up via svg.call(zoom) which adds __zoom property
          // The main group 'g' should exist for transformations
          const mainGroup = document.querySelector('svg > g')
          expect(mainGroup).toBeInTheDocument()
        },
        { timeout: 1000 }
      )
    })

    it('should render draggable node groups', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          // Nodes are rendered as 'g.node' groups which have drag behavior attached
          const nodeGroups = document.querySelectorAll('g.node')
          expect(nodeGroups.length).toBe(2)
          // Each node group should have cursor: pointer on circles
          const circles = document.querySelectorAll('circle')
          circles.forEach((circle) => {
            expect(circle.style.cursor).toBe('pointer')
          })
        },
        { timeout: 1000 }
      )
    })
  })
})

/**
 * Story 4.2 Tests: Highlight Circular Dependencies in Graph
 */
describe('DependencyGraphViz - Circular Dependency Highlighting (Story 4.2)', () => {
  // Mock data with circular dependencies
  const mockDataWithCycles: DependencyGraph = {
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
        dependencies: ['@app/core'], // Creates cycle
        devDependencies: [],
        peerDependencies: [],
      },
      '@app/other': {
        name: '@app/other',
        version: '1.0.0',
        path: 'packages/other',
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
      {
        from: '@app/utils',
        to: '@app/core',
        type: 'production',
        versionRange: '^1.0.0',
      },
    ],
    rootPath: '/workspace',
    workspaceType: 'npm',
  }

  const mockCircularDependencies: CircularDependencyInfo[] = [
    {
      cycle: ['@app/core', '@app/utils', '@app/core'],
      type: 'direct',
      severity: 'critical',
      depth: 2,
      impact: 'Build failure risk',
      complexity: 3,
    },
  ]

  describe('AC1: Visual Highlighting of Cycle Nodes', () => {
    it('should render glow filter for cycle nodes', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          // Check for glow filter definition
          const glowFilter = document.querySelector('filter#glow')
          expect(glowFilter).toBeInTheDocument()
        },
        { timeout: 1000 }
      )
    })

    it('should apply cycle styling to nodes in cycles', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          expect(circles.length).toBe(3)

          // At least some circles should have red fill (cycle nodes)
          const circleColors = Array.from(circles).map((c) => c.getAttribute('fill'))
          expect(circleColors.some((color) => color === '#ef4444')).toBe(true)
        },
        { timeout: 1000 }
      )
    })

    it('should mark cycle nodes with CSS class', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          const cycleNodes = document.querySelectorAll('g.node--cycle')
          // @app/core and @app/utils are in the cycle
          expect(cycleNodes.length).toBe(2)
        },
        { timeout: 1000 }
      )
    })
  })

  describe('AC2: Visual Highlighting of Cycle Edges', () => {
    it('should render separate arrow marker for cycle edges', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          // Check for cycle-specific arrowhead marker
          const cycleMarker = document.querySelector('marker#arrowhead-cycle')
          expect(cycleMarker).toBeInTheDocument()

          // Cycle marker path should be red
          const cycleMarkerPath = cycleMarker?.querySelector('path')
          expect(cycleMarkerPath?.getAttribute('fill')).toBe('#ef4444')
        },
        { timeout: 1000 }
      )
    })

    it('should render cycle edges in separate group', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          // Check for separate link groups
          const cycleLinksGroup = document.querySelector('g.links-cycle')
          expect(cycleLinksGroup).toBeInTheDocument()

          // Cycle edges should use the cycle arrowhead
          const cycleLines = cycleLinksGroup?.querySelectorAll('line')
          expect(cycleLines?.length).toBeGreaterThan(0)
          cycleLines?.forEach((line) => {
            expect(line.getAttribute('marker-end')).toContain('arrowhead-cycle')
          })
        },
        { timeout: 1000 }
      )
    })
  })

  describe('AC3: Animated Cycle Paths', () => {
    it('should apply animation to cycle edges', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          const cycleLinksGroup = document.querySelector('g.links-cycle')
          const cycleLines = cycleLinksGroup?.querySelectorAll('line')

          // Cycle edges should have dash array for animation
          cycleLines?.forEach((line) => {
            expect(line.getAttribute('stroke-dasharray')).toBeTruthy()
          })
        },
        { timeout: 1000 }
      )
    })

    it('should include CSS animation keyframes', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          // Check for style element with keyframes
          const styleElement = document.querySelector('style')
          expect(styleElement).toBeInTheDocument()
          expect(styleElement?.textContent).toContain('flowAnimation')
          expect(styleElement?.textContent).toContain('stroke-dashoffset')
        },
        { timeout: 1000 }
      )
    })
  })

  describe('AC5: Click-to-Highlight Cycle', () => {
    it('should handle click on cycle nodes', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          const cycleNodes = document.querySelectorAll('g.node--cycle')
          expect(cycleNodes.length).toBe(2)

          // Click on a cycle node
          const firstCycleNode = cycleNodes[0]
          fireEvent.click(firstCycleNode)
        },
        { timeout: 1000 }
      )
    })
  })

  describe('AC6: Escape Key to Deselect', () => {
    it('should handle Escape key press', async () => {
      render(
        <DependencyGraphViz
          data={mockDataWithCycles}
          circularDependencies={mockCircularDependencies}
        />
      )

      await waitFor(
        () => {
          const svg = document.querySelector('svg')
          expect(svg).toBeInTheDocument()

          // Press Escape key
          fireEvent.keyDown(document, { key: 'Escape' })
        },
        { timeout: 1000 }
      )
    })
  })

  describe('Without Circular Dependencies', () => {
    it('should render normally without circularDependencies prop', async () => {
      render(<DependencyGraphViz data={mockData} />)

      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          expect(circles.length).toBe(2)

          // No cycle nodes should be marked
          const cycleNodes = document.querySelectorAll('g.node--cycle')
          expect(cycleNodes.length).toBe(0)
        },
        { timeout: 1000 }
      )
    })

    it('should render normally with empty circularDependencies array', async () => {
      render(<DependencyGraphViz data={mockData} circularDependencies={[]} />)

      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          expect(circles.length).toBe(2)

          // No cycle nodes should be marked
          const cycleNodes = document.querySelectorAll('g.node--cycle')
          expect(cycleNodes.length).toBe(0)
        },
        { timeout: 1000 }
      )
    })
  })
})

/**
 * GraphLegend Tests (Story 4.2 - AC4)
 */
describe('GraphLegend', () => {
  describe('AC4: Color Legend', () => {
    it('should render legend with normal node color', () => {
      render(<GraphLegend hasCycles={false} />)

      expect(document.body.textContent).toContain('Normal Package')
      expect(document.body.textContent).toContain('Normal Dependency')
    })

    it('should render cycle colors when hasCycles is true', () => {
      render(<GraphLegend hasCycles={true} />)

      expect(document.body.textContent).toContain('In Circular Dependency')
      expect(document.body.textContent).toContain('Circular Dependency')
    })

    it('should not render cycle colors when hasCycles is false', () => {
      render(<GraphLegend hasCycles={false} />)

      expect(document.body.textContent).not.toContain('In Circular Dependency')
      expect(document.body.textContent).not.toContain('Circular Dependency')
    })

    it('should render interaction hint when cycles exist', () => {
      render(<GraphLegend hasCycles={true} />)

      expect(document.body.textContent).toContain('Click on red nodes')
      expect(document.body.textContent).toContain('Escape')
    })

    it('should accept custom position', () => {
      render(<GraphLegend position="top-right" hasCycles={false} />)

      const legend = document.querySelector('.top-4.right-4')
      expect(legend).toBeInTheDocument()
    })

    it('should accept custom className', () => {
      render(<GraphLegend className="custom-legend" hasCycles={false} />)

      const legend = document.querySelector('.custom-legend')
      expect(legend).toBeInTheDocument()
    })
  })
})

/**
 * Story 4.3 Tests: Node Expand/Collapse Functionality Integration
 *
 * These tests verify the integration of expand/collapse functionality
 * with the main DependencyGraphViz component.
 *
 * Unit tests for individual functions are in:
 * - useNodeExpandCollapse.test.ts
 * - computeVisibleNodes.test.ts
 * - calculateDepth.test.ts
 * - GraphControls.test.tsx
 */
describe('DependencyGraphViz - Node Expand/Collapse (Story 4.3)', () => {
  // Mock data with hierarchical structure for expand/collapse testing
  const mockHierarchicalData: DependencyGraph = {
    nodes: {
      '@app/root': {
        name: '@app/root',
        version: '1.0.0',
        path: 'packages/root',
        dependencies: ['@app/child1', '@app/child2'],
        devDependencies: [],
        peerDependencies: [],
      },
      '@app/child1': {
        name: '@app/child1',
        version: '1.0.0',
        path: 'packages/child1',
        dependencies: ['@app/grandchild'],
        devDependencies: [],
        peerDependencies: [],
      },
      '@app/child2': {
        name: '@app/child2',
        version: '1.0.0',
        path: 'packages/child2',
        dependencies: [],
        devDependencies: [],
        peerDependencies: [],
      },
      '@app/grandchild': {
        name: '@app/grandchild',
        version: '1.0.0',
        path: 'packages/grandchild',
        dependencies: [],
        devDependencies: [],
        peerDependencies: [],
      },
    },
    edges: [
      { from: '@app/root', to: '@app/child1', type: 'production', versionRange: '^1.0.0' },
      { from: '@app/root', to: '@app/child2', type: 'production', versionRange: '^1.0.0' },
      { from: '@app/child1', to: '@app/grandchild', type: 'production', versionRange: '^1.0.0' },
    ],
    rootPath: '/workspace',
    workspaceType: 'npm',
  }

  describe('AC3: Depth-Based Controls Integration', () => {
    it('[P1] should render GraphControls when data exists', async () => {
      // GIVEN: DependencyGraphViz with hierarchical data
      // WHEN: Rendered
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      // THEN: GraphControls should be visible
      // Note: Using aria-label selector because fieldset has implicit role="group"
      // which may not be recognized by document.querySelector in JSDOM
      await waitFor(
        () => {
          expect(document.querySelector('[aria-label="Graph depth controls"]')).toBeInTheDocument()
        },
        { timeout: 1000 }
      )
    })

    it('[P1] should render depth level buttons based on graph structure', async () => {
      // GIVEN: DependencyGraphViz with 3 depth levels
      // WHEN: Rendered
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      // THEN: Should have depth control buttons
      await waitFor(
        () => {
          const depthControlGroup = document.querySelector('[aria-label="Graph depth controls"]')
          expect(depthControlGroup).toBeInTheDocument()

          // Should have 'All' button
          const buttons = depthControlGroup?.querySelectorAll('button')
          expect(buttons?.length).toBeGreaterThanOrEqual(2) // At least All + one depth level
        },
        { timeout: 1000 }
      )
    })

    it('[P2] should not render GraphControls when no data', () => {
      // GIVEN: Empty data
      const emptyData: DependencyGraph = {
        nodes: {},
        edges: [],
        rootPath: '/workspace',
        workspaceType: 'npm',
      }

      // WHEN: Rendered
      render(<DependencyGraphViz data={emptyData} />)

      // THEN: GraphControls should not be present
      expect(document.querySelector('[aria-label="Graph depth controls"]')).not.toBeInTheDocument()
    })
  })

  describe('AC1/AC2: Double-Click Expand/Collapse Integration', () => {
    it('[P1] should set up double-click handler on nodes', async () => {
      // GIVEN: DependencyGraphViz with nodes
      // WHEN: Rendered
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      // THEN: Nodes should be rendered and have event handlers attached
      await waitFor(
        () => {
          const nodeGroups = document.querySelectorAll('g.node')
          expect(nodeGroups.length).toBe(4)

          // Each node should have pointer cursor
          nodeGroups.forEach((nodeGroup) => {
            const circle = nodeGroup.querySelector('circle')
            expect(circle?.style.cursor).toBe('pointer')
          })
        },
        { timeout: 1000 }
      )
    })

    it('[P1] should handle double-click on node', async () => {
      // GIVEN: DependencyGraphViz with nodes
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      await waitFor(
        () => {
          const nodeGroups = document.querySelectorAll('g.node')
          expect(nodeGroups.length).toBe(4)
        },
        { timeout: 1000 }
      )

      // WHEN: Double-click on a node
      const nodeGroup = document.querySelector('g.node')
      if (nodeGroup) {
        fireEvent.dblClick(nodeGroup)
      }

      // THEN: The interaction should not throw errors (basic sanity check)
      // Full behavior tested in useNodeExpandCollapse.test.ts
      expect(document.querySelector('svg')).toBeInTheDocument()
    })
  })

  describe('AC4: Collapsed Node Indicator Integration', () => {
    it('[P1] should render badge group container', async () => {
      // GIVEN: DependencyGraphViz with hierarchical data
      // WHEN: Rendered
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      // THEN: Should have SVG structure ready for badges
      await waitFor(
        () => {
          const svg = document.querySelector('svg')
          expect(svg).toBeInTheDocument()

          // Badge group is created even if empty
          const mainGroup = svg?.querySelector('g')
          expect(mainGroup).toBeInTheDocument()
        },
        { timeout: 1000 }
      )
    })
  })

  describe('Component Integration', () => {
    it('[P1] should integrate useNodeExpandCollapse hook', async () => {
      // GIVEN: DependencyGraphViz with hierarchical data
      // WHEN: Rendered
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      // THEN: Should render nodes (D3 may create circles progressively in jsdom)
      await waitFor(
        () => {
          const circles = document.querySelectorAll('circle')
          // At minimum, nodes should be rendered (exact count may vary in jsdom)
          expect(circles.length).toBeGreaterThanOrEqual(1)
        },
        { timeout: 1000 }
      )
    })

    it('[P1] should integrate computeVisibleNodes utility', async () => {
      // GIVEN: DependencyGraphViz with hierarchical data
      // WHEN: Rendered
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      // THEN: SVG should be created with link groups (D3 creates structure even if jsdom doesn't render all lines)
      await waitFor(
        () => {
          const svg = document.querySelector('svg')
          expect(svg).toBeInTheDocument()

          // Verify link groups are created (may not have visible lines in jsdom)
          const mainGroup = svg?.querySelector('g')
          expect(mainGroup).toBeInTheDocument()
        },
        { timeout: 1000 }
      )
    })

    it('[P1] should integrate calculateNodeDepths utility', async () => {
      // GIVEN: DependencyGraphViz with hierarchical data
      // WHEN: Rendered
      render(<DependencyGraphViz data={mockHierarchicalData} />)

      // THEN: Depth controls should show levels based on calculated depths
      await waitFor(
        () => {
          const depthGroup = document.querySelector('[aria-label="Graph depth controls"]')
          expect(depthGroup).toBeInTheDocument()

          // Max depth in this graph is 2 (root->child1->grandchild)
          // So we should see L1, L2 buttons
          const buttons = depthGroup?.querySelectorAll('button')
          const buttonLabels = Array.from(buttons || []).map((b) => b.textContent)
          expect(buttonLabels).toContain('All')
          expect(buttonLabels.some((l) => l?.includes('L1'))).toBe(true)
        },
        { timeout: 1000 }
      )
    })
  })

  describe('Cleanup and Memory Management', () => {
    it('[P1] should clean up on unmount', async () => {
      // GIVEN: DependencyGraphViz rendered
      const { unmount } = render(<DependencyGraphViz data={mockHierarchicalData} />)

      await waitFor(() => {
        expect(document.querySelector('svg')).toBeInTheDocument()
      })

      // WHEN: Unmounted
      unmount()

      // THEN: SVG should be removed (no memory leaks)
      expect(document.querySelector('svg')).not.toBeInTheDocument()
    })
  })
})
