import type { CircularDependencyInfo } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import { generateCyclePath } from '../sections/cyclePath'

const baseCycle: CircularDependencyInfo = {
  cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
  type: 'indirect',
  severity: 'warning',
  depth: 3,
  impact: 'Test impact',
  complexity: 5,
  priorityScore: 5,
  importTraces: [
    {
      fromPackage: 'pkg-a',
      toPackage: 'pkg-b',
      filePath: 'src/index.ts',
      lineNumber: 1,
      statement: "import { foo } from 'pkg-b'",
      importType: 'esm-named',
    },
    {
      fromPackage: 'pkg-b',
      toPackage: 'pkg-c',
      filePath: 'src/main.ts',
      lineNumber: 5,
      statement: "import { bar } from 'pkg-c'",
      importType: 'esm-named',
    },
    {
      fromPackage: 'pkg-c',
      toPackage: 'pkg-a',
      filePath: 'src/util.ts',
      lineNumber: 10,
      statement: "import { baz } from 'pkg-a'",
      importType: 'esm-named',
    },
  ],
}

describe('generateCyclePath', () => {
  it('should generate nodes for all unique packages', () => {
    const path = generateCyclePath(baseCycle)
    expect(path.nodes).toHaveLength(3)
    expect(path.nodes.map((n) => n.id)).toEqual(['pkg-a', 'pkg-b', 'pkg-c'])
  })

  it('should generate edges between consecutive packages', () => {
    const path = generateCyclePath(baseCycle)
    expect(path.edges).toHaveLength(3)
    expect(path.edges[0]).toMatchObject({ from: 'pkg-a', to: 'pkg-b' })
    expect(path.edges[1]).toMatchObject({ from: 'pkg-b', to: 'pkg-c' })
    expect(path.edges[2]).toMatchObject({ from: 'pkg-c', to: 'pkg-a' })
  })

  it('should identify breaking point', () => {
    const path = generateCyclePath(baseCycle)
    expect(path.breakingPoint).toBeDefined()
    expect(path.breakingPoint.fromPackage).toBeTruthy()
    expect(path.breakingPoint.toPackage).toBeTruthy()
    expect(path.breakingPoint.reason).toBeTruthy()
  })

  it('should use criticalEdge as breaking point when available', () => {
    const withCriticalEdge: CircularDependencyInfo = {
      ...baseCycle,
      rootCause: {
        originatingPackage: 'pkg-a',
        problematicDependency: { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
        confidence: 90,
        explanation: 'Test',
        chain: [],
        criticalEdge: { from: 'pkg-b', to: 'pkg-c', type: 'production', critical: true },
      },
    }
    const path = generateCyclePath(withCriticalEdge)
    expect(path.breakingPoint.fromPackage).toBe('pkg-b')
    expect(path.breakingPoint.toPackage).toBe('pkg-c')
  })

  it('should generate valid SVG diagram', () => {
    const path = generateCyclePath(baseCycle)
    expect(path.svgDiagram).toContain('<svg')
    expect(path.svgDiagram).toContain('</svg>')
    expect(path.svgDiagram).toContain('pkg-a')
  })

  it('should generate ASCII diagram', () => {
    const path = generateCyclePath(baseCycle)
    expect(path.asciiDiagram).toContain('Cycle Path')
    expect(path.asciiDiagram).toContain('pkg-a')
    expect(path.asciiDiagram).toContain('BREAK HERE')
  })

  it('should support dark mode', () => {
    const lightPath = generateCyclePath(baseCycle, false)
    const darkPath = generateCyclePath(baseCycle, true)
    expect(lightPath.svgDiagram).toContain('#ffffff')
    expect(darkPath.svgDiagram).toContain('#1f2937')
  })

  it('should include import traces in edges', () => {
    const path = generateCyclePath(baseCycle)
    expect(path.edges[0].importStatement).toBe("import { foo } from 'pkg-b'")
    expect(path.edges[0].filePath).toBe('src/index.ts')
    expect(path.edges[0].lineNumber).toBe(1)
  })

  it('should position nodes in a circle', () => {
    const path = generateCyclePath(baseCycle)
    for (const node of path.nodes) {
      expect(node.position.x).toBeGreaterThan(0)
      expect(node.position.y).toBeGreaterThan(0)
    }
  })

  it('should mark exactly one edge as breaking point', () => {
    const path = generateCyclePath(baseCycle)
    const breakingEdges = path.edges.filter((e) => e.isBreakingPoint)
    expect(breakingEdges).toHaveLength(1)
  })

  it('should handle cycle without import traces', () => {
    const cycleNoTraces: CircularDependencyInfo = {
      ...baseCycle,
      importTraces: undefined,
    }
    const path = generateCyclePath(cycleNoTraces)
    expect(path.nodes).toHaveLength(3)
    expect(path.edges).toHaveLength(3)
    expect(path.edges[0].importStatement).toBeUndefined()
  })
})
