/**
 * Tests for calculateNodeDepths utility
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality
 */

import { describe, expect, it } from 'vitest'
import { calculateNodeDepths } from '../utils/calculateDepth'

describe('calculateNodeDepths', () => {
  it('should calculate depths for a simple linear graph', () => {
    const nodeIds = ['a', 'b', 'c']
    const edges = [
      { source: 'a', target: 'b' },
      { source: 'b', target: 'c' },
    ]

    const depths = calculateNodeDepths(nodeIds, edges)

    expect(depths.get('a')).toBe(0)
    expect(depths.get('b')).toBe(1)
    expect(depths.get('c')).toBe(2)
  })

  it('should handle multiple root nodes', () => {
    const nodeIds = ['root1', 'root2', 'child1', 'child2']
    const edges = [
      { source: 'root1', target: 'child1' },
      { source: 'root2', target: 'child2' },
    ]

    const depths = calculateNodeDepths(nodeIds, edges)

    expect(depths.get('root1')).toBe(0)
    expect(depths.get('root2')).toBe(0)
    expect(depths.get('child1')).toBe(1)
    expect(depths.get('child2')).toBe(1)
  })

  it('should use minimum depth when node has multiple paths', () => {
    const nodeIds = ['root', 'a', 'b', 'shared']
    const edges = [
      { source: 'root', target: 'a' },
      { source: 'root', target: 'b' },
      { source: 'a', target: 'shared' }, // depth 2
      { source: 'b', target: 'shared' }, // depth 2
      { source: 'root', target: 'shared' }, // depth 1 - should win
    ]

    const depths = calculateNodeDepths(nodeIds, edges)

    expect(depths.get('root')).toBe(0)
    expect(depths.get('shared')).toBe(1) // Minimum depth
  })

  it('should handle disconnected nodes', () => {
    const nodeIds = ['root', 'child', 'orphan']
    const edges = [{ source: 'root', target: 'child' }]

    const depths = calculateNodeDepths(nodeIds, edges)

    expect(depths.get('root')).toBe(0)
    expect(depths.get('child')).toBe(1)
    // Orphan has no incoming edges, so it's also a root (depth 0)
    // In a dependency graph, a node with no incoming edges is a top-level package
    expect(depths.get('orphan')).toBe(0)
  })

  it('should handle empty graph', () => {
    const depths = calculateNodeDepths([], [])
    expect(depths.size).toBe(0)
  })

  it('should handle graph with no edges', () => {
    const nodeIds = ['a', 'b', 'c']
    const edges: Array<{ source: string; target: string }> = []

    const depths = calculateNodeDepths(nodeIds, edges)

    // All nodes are roots
    expect(depths.get('a')).toBe(0)
    expect(depths.get('b')).toBe(0)
    expect(depths.get('c')).toBe(0)
  })

  it('should handle circular dependencies', () => {
    const nodeIds = ['a', 'b', 'c']
    const edges = [
      { source: 'a', target: 'b' },
      { source: 'b', target: 'c' },
      { source: 'c', target: 'a' }, // Creates cycle
    ]

    const depths = calculateNodeDepths(nodeIds, edges)

    // In a cycle, node that appears first in BFS from nodes without external incoming edges
    // Here 'a' has incoming from 'c' but if we consider all nodes, none are true roots
    // Algorithm should still assign depths based on BFS order
    expect(depths.size).toBe(3)
    // All nodes should have some depth assigned
    expect(depths.has('a')).toBe(true)
    expect(depths.has('b')).toBe(true)
    expect(depths.has('c')).toBe(true)
  })

  it('should calculate correct max depth', () => {
    const nodeIds = ['root', 'l1a', 'l1b', 'l2a', 'l2b', 'l3']
    const edges = [
      { source: 'root', target: 'l1a' },
      { source: 'root', target: 'l1b' },
      { source: 'l1a', target: 'l2a' },
      { source: 'l1b', target: 'l2b' },
      { source: 'l2a', target: 'l3' },
    ]

    const depths = calculateNodeDepths(nodeIds, edges)

    expect(depths.get('root')).toBe(0)
    expect(depths.get('l1a')).toBe(1)
    expect(depths.get('l1b')).toBe(1)
    expect(depths.get('l2a')).toBe(2)
    expect(depths.get('l2b')).toBe(2)
    expect(depths.get('l3')).toBe(3)
  })
})
