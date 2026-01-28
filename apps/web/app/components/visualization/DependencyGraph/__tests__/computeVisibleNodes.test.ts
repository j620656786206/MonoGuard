/**
 * Tests for computeVisibleNodes utility
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality
 */

import { describe, expect, it } from 'vitest'
import type { D3Link, D3Node } from '../types'
import { computeVisibleNodes } from '../utils/computeVisibleNodes'

describe('computeVisibleNodes', () => {
  const createNode = (id: string, overrides?: Partial<D3Node>): D3Node => ({
    id,
    name: id,
    path: `/${id}`,
    dependencyCount: 0,
    inCycle: false,
    cycleIds: [],
    ...overrides,
  })

  const createLink = (source: string, target: string): D3Link => ({
    source,
    target,
    type: 'production',
    inCycle: false,
    cycleIds: [],
  })

  const mockNodes: D3Node[] = [
    createNode('root', { dependencyCount: 2 }),
    createNode('child1', { dependencyCount: 1 }),
    createNode('child2', { dependencyCount: 0 }),
    createNode('grandchild', { dependencyCount: 0 }),
  ]

  const mockLinks: D3Link[] = [
    createLink('root', 'child1'),
    createLink('root', 'child2'),
    createLink('child1', 'grandchild'),
  ]

  it('should return all nodes when none are collapsed', () => {
    const result = computeVisibleNodes(mockNodes, mockLinks, new Set())

    expect(result.visibleNodes.length).toBe(4)
    expect(result.visibleLinks.length).toBe(3)
    expect(result.hiddenChildCounts.size).toBe(0)
  })

  it('should hide descendants when parent is collapsed', () => {
    const result = computeVisibleNodes(mockNodes, mockLinks, new Set(['child1']))

    expect(result.visibleNodes.map((n) => n.id)).toContain('root')
    expect(result.visibleNodes.map((n) => n.id)).toContain('child1')
    expect(result.visibleNodes.map((n) => n.id)).toContain('child2')
    expect(result.visibleNodes.map((n) => n.id)).not.toContain('grandchild')
    expect(result.hiddenChildCounts.get('child1')).toBe(1)
  })

  it('should keep nodes visible if they have alternate paths', () => {
    // Add alternate path to grandchild
    const linksWithAlternatePath: D3Link[] = [...mockLinks, createLink('root', 'grandchild')]

    const result = computeVisibleNodes(mockNodes, linksWithAlternatePath, new Set(['child1']))

    // grandchild should still be visible because root has direct edge to it
    expect(result.visibleNodes.map((n) => n.id)).toContain('grandchild')
    expect(result.hiddenChildCounts.get('child1')).toBe(0)
  })

  it('should handle empty collapsed set', () => {
    const result = computeVisibleNodes(mockNodes, mockLinks, new Set())

    expect(result.visibleNodes).toEqual(mockNodes)
    expect(result.visibleLinks).toEqual(mockLinks)
    expect(result.hiddenChildCounts.size).toBe(0)
  })

  it('should handle multiple collapsed nodes', () => {
    // Create a more complex graph
    const nodes: D3Node[] = [
      createNode('root'),
      createNode('a'),
      createNode('b'),
      createNode('a1'),
      createNode('a2'),
      createNode('b1'),
    ]

    const links: D3Link[] = [
      createLink('root', 'a'),
      createLink('root', 'b'),
      createLink('a', 'a1'),
      createLink('a', 'a2'),
      createLink('b', 'b1'),
    ]

    const result = computeVisibleNodes(nodes, links, new Set(['a', 'b']))

    // Only root, a, and b should be visible
    expect(result.visibleNodes.map((n) => n.id).sort()).toEqual(['a', 'b', 'root'])
    expect(result.hiddenChildCounts.get('a')).toBe(2)
    expect(result.hiddenChildCounts.get('b')).toBe(1)
  })

  it('should handle nested collapsed nodes', () => {
    const nodes: D3Node[] = [
      createNode('root'),
      createNode('child'),
      createNode('grandchild'),
      createNode('greatgrandchild'),
    ]

    const links: D3Link[] = [
      createLink('root', 'child'),
      createLink('child', 'grandchild'),
      createLink('grandchild', 'greatgrandchild'),
    ]

    // Collapse both child and grandchild
    const result = computeVisibleNodes(nodes, links, new Set(['child', 'grandchild']))

    // Only root and child should be visible
    expect(result.visibleNodes.map((n) => n.id).sort()).toEqual(['child', 'root'])
    expect(result.hiddenChildCounts.get('child')).toBe(2) // grandchild + greatgrandchild
  })

  it('should handle circular dependencies', () => {
    const nodes: D3Node[] = [
      createNode('a', { inCycle: true }),
      createNode('b', { inCycle: true }),
      createNode('c', { inCycle: true }),
      createNode('d'),
    ]

    const links: D3Link[] = [
      createLink('a', 'b'),
      createLink('b', 'c'),
      createLink('c', 'a'), // Creates cycle
      createLink('a', 'd'),
    ]

    const result = computeVisibleNodes(nodes, links, new Set(['a']))

    // b, c are in cycle with a but should be hidden when a is collapsed
    // d is child of a and should be hidden
    expect(result.visibleNodes.map((n) => n.id)).toContain('a')
    expect(result.visibleNodes.map((n) => n.id)).not.toContain('d')
  })

  it('should handle empty graph', () => {
    const result = computeVisibleNodes([], [], new Set())

    expect(result.visibleNodes).toEqual([])
    expect(result.visibleLinks).toEqual([])
    expect(result.hiddenChildCounts.size).toBe(0)
  })

  it('should handle graph with no edges', () => {
    const isolatedNodes: D3Node[] = [createNode('a'), createNode('b'), createNode('c')]

    const result = computeVisibleNodes(isolatedNodes, [], new Set(['a']))

    // All nodes should remain visible since they're all roots
    expect(result.visibleNodes.length).toBe(3)
    expect(result.hiddenChildCounts.get('a')).toBe(0)
  })

  it('should handle D3 resolved links (with node objects)', () => {
    const nodeA = createNode('a')
    const nodeB = createNode('b')
    const nodeC = createNode('c')

    const nodes = [nodeA, nodeB, nodeC]

    // Links with resolved node objects (as D3 does after simulation)
    const resolvedLinks: D3Link[] = [
      { source: nodeA, target: nodeB, type: 'production', inCycle: false, cycleIds: [] },
      { source: nodeB, target: nodeC, type: 'production', inCycle: false, cycleIds: [] },
    ]

    const result = computeVisibleNodes(nodes, resolvedLinks, new Set(['b']))

    expect(result.visibleNodes.map((n) => n.id)).toContain('a')
    expect(result.visibleNodes.map((n) => n.id)).toContain('b')
    expect(result.visibleNodes.map((n) => n.id)).not.toContain('c')
    expect(result.hiddenChildCounts.get('b')).toBe(1)
  })
})
