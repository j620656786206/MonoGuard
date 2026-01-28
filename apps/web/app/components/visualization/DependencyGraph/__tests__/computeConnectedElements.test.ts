/**
 * Tests for computeConnectedElements utilities
 *
 * @see Story 4.5: Implement Hover Details and Tooltips (AC1, AC4)
 *
 * Following Given-When-Then format with priority tags.
 */

import type { CircularDependencyInfo } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import type { D3Link, D3Node } from '../types'
import {
  computeConnectedElements,
  computeDependencyCounts,
  computeTooltipData,
} from '../utils/computeConnectedElements'

describe('computeConnectedElements', () => {
  const createMockLinks = (): D3Link[] => [
    { source: 'A', target: 'B', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'A', target: 'C', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'B', target: 'D', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'C', target: 'A', type: 'production', inCycle: true, cycleIds: [0] }, // Cycle
  ]

  describe('Node Connection Finding (AC4)', () => {
    it('[P1] should find all connected nodes for a given node', () => {
      // GIVEN: Links where A connects to B and C
      // WHEN: Computing connected elements for A
      const result = computeConnectedElements('A', createMockLinks())

      // THEN: Should include A, B, C (A's direct connections)
      expect(result.nodeIds.has('A')).toBe(true)
      expect(result.nodeIds.has('B')).toBe(true)
      expect(result.nodeIds.has('C')).toBe(true)
      // D is not directly connected to A
      expect(result.nodeIds.has('D')).toBe(false)
    })

    it('[P1] should find all connected link indices', () => {
      // GIVEN: Links [A->B, A->C, B->D, C->A]
      // WHEN: Computing connected elements for A
      const result = computeConnectedElements('A', createMockLinks())

      // THEN: Should include links 0 (A->B), 1 (A->C), and 3 (C->A)
      expect(result.linkIndices.has(0)).toBe(true)
      expect(result.linkIndices.has(1)).toBe(true)
      expect(result.linkIndices.has(3)).toBe(true)
      // B->D is not directly connected to A
      expect(result.linkIndices.has(2)).toBe(false)
    })

    it('[P1] should handle nodes with no connections', () => {
      // GIVEN: Links that don't include 'isolated' node
      // WHEN: Computing connected elements for isolated node
      const result = computeConnectedElements('isolated', createMockLinks())

      // THEN: Should only contain the node itself
      expect(result.nodeIds.size).toBe(1)
      expect(result.nodeIds.has('isolated')).toBe(true)
      expect(result.linkIndices.size).toBe(0)
    })

    it('[P2] should handle links with D3Node objects', () => {
      // GIVEN: Links with node objects instead of string IDs
      const nodeA: D3Node = {
        id: 'A',
        name: 'A',
        path: '/a',
        dependencyCount: 0,
        inCycle: false,
        cycleIds: [],
      }
      const nodeB: D3Node = {
        id: 'B',
        name: 'B',
        path: '/b',
        dependencyCount: 0,
        inCycle: false,
        cycleIds: [],
      }
      const links: D3Link[] = [
        { source: nodeA, target: nodeB, type: 'production', inCycle: false, cycleIds: [] },
      ]

      // WHEN: Computing connected elements for A
      const result = computeConnectedElements('A', links)

      // THEN: Should correctly identify connections
      expect(result.nodeIds.has('A')).toBe(true)
      expect(result.nodeIds.has('B')).toBe(true)
      expect(result.linkIndices.has(0)).toBe(true)
    })
  })
})

describe('computeDependencyCounts', () => {
  const createMockLinks = (): D3Link[] => [
    { source: 'A', target: 'B', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'A', target: 'C', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'D', target: 'A', type: 'production', inCycle: false, cycleIds: [] },
  ]

  describe('Dependency Counting (AC1)', () => {
    it('[P1] should count incoming dependencies correctly', () => {
      // GIVEN: Links where D->A (A has 1 incoming)
      // WHEN: Computing dependency counts for A
      const result = computeDependencyCounts('A', createMockLinks())

      // THEN: Should have 1 incoming (D->A)
      expect(result.incoming).toBe(1)
    })

    it('[P1] should count outgoing dependencies correctly', () => {
      // GIVEN: Links where A->B and A->C (A has 2 outgoing)
      // WHEN: Computing dependency counts for A
      const result = computeDependencyCounts('A', createMockLinks())

      // THEN: Should have 2 outgoing (A->B, A->C)
      expect(result.outgoing).toBe(2)
    })

    it('[P1] should return zero counts for isolated nodes', () => {
      // GIVEN: Links that don't include 'isolated' node
      // WHEN: Computing dependency counts for isolated node
      const result = computeDependencyCounts('isolated', createMockLinks())

      // THEN: Should have zero for both
      expect(result.incoming).toBe(0)
      expect(result.outgoing).toBe(0)
    })

    it('[P2] should handle node as both source and target', () => {
      // GIVEN: Links where B has incoming from A
      // WHEN: Computing dependency counts for B
      const result = computeDependencyCounts('B', createMockLinks())

      // THEN: Should have 1 incoming (A->B), 0 outgoing
      expect(result.incoming).toBe(1)
      expect(result.outgoing).toBe(0)
    })
  })
})

describe('computeTooltipData', () => {
  const createMockNode = (id: string, inCycle = false, cycleIds: number[] = []): D3Node => ({
    id,
    name: `Package ${id}`,
    path: `packages/${id.toLowerCase()}`,
    dependencyCount: 2,
    inCycle,
    cycleIds,
  })

  const createMockLinks = (): D3Link[] => [
    { source: 'A', target: 'B', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'A', target: 'C', type: 'production', inCycle: false, cycleIds: [] },
    { source: 'D', target: 'A', type: 'production', inCycle: false, cycleIds: [] },
  ]

  const createMockCircularDeps = (): CircularDependencyInfo[] => [
    {
      cycle: ['A', 'B', 'C', 'A'],
      type: 'direct',
      severity: 'critical',
      depth: 3,
      impact: 'High impact on build times',
      complexity: 5,
      priorityScore: 75,
    },
  ]

  describe('Tooltip Data Generation (AC1)', () => {
    it('[P1] should include package name and path', () => {
      // GIVEN: A node with name and path
      const node = createMockNode('A')

      // WHEN: Computing tooltip data
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: [],
      })

      // THEN: Should include name and path
      expect(result.packageName).toBe('Package A')
      expect(result.packagePath).toBe('packages/a')
    })

    it('[P1] should calculate incoming and outgoing counts', () => {
      // GIVEN: Node A with 1 incoming (D->A) and 2 outgoing (A->B, A->C)
      const node = createMockNode('A')

      // WHEN: Computing tooltip data
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: [],
      })

      // THEN: Should have correct counts
      expect(result.incomingCount).toBe(1)
      expect(result.outgoingCount).toBe(2)
    })

    it('[P1] should set inCycle to false when not in cycle', () => {
      // GIVEN: Node not in any cycle
      const node = createMockNode('A')

      // WHEN: Computing tooltip data with no circular dependencies
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: [],
      })

      // THEN: inCycle should be false
      expect(result.inCycle).toBe(false)
      expect(result.cycleInfo).toBeUndefined()
    })

    it('[P1] should set inCycle to true when in cycle', () => {
      // GIVEN: Node that is in a cycle
      const node = createMockNode('A', true, [0])
      const cycles = createMockCircularDeps()

      // WHEN: Computing tooltip data
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: cycles,
      })

      // THEN: inCycle should be true with cycle info
      expect(result.inCycle).toBe(true)
      expect(result.cycleInfo).toBeDefined()
      expect(result.cycleInfo?.cycleCount).toBe(1)
    })

    it('[P1] should calculate positive health contribution for healthy nodes', () => {
      // GIVEN: Node not in cycle with few connections
      const node = createMockNode('D') // D has 0 incoming, 1 outgoing

      // WHEN: Computing tooltip data
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: [],
      })

      // THEN: Health contribution should be positive
      expect(result.healthContribution).toBeGreaterThanOrEqual(0)
    })

    it('[P1] should calculate negative health contribution for nodes in cycles', () => {
      // GIVEN: Node in a cycle
      const node = createMockNode('A', true, [0])
      const cycles = createMockCircularDeps()

      // WHEN: Computing tooltip data
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: cycles,
      })

      // THEN: Health contribution should be negative
      expect(result.healthContribution).toBeLessThan(0)
    })

    it('[P2] should list other packages in the cycle', () => {
      // GIVEN: Node A in cycle with B and C
      const node = createMockNode('A', true, [0])
      const cycles = createMockCircularDeps()

      // WHEN: Computing tooltip data
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: cycles,
      })

      // THEN: Cycle info should list other packages (not including A itself)
      expect(result.cycleInfo?.packages).toContain('B')
      expect(result.cycleInfo?.packages).toContain('C')
      expect(result.cycleInfo?.packages).not.toContain('A')
    })

    it('[P2] should handle node in multiple cycles', () => {
      // GIVEN: Node in 2 cycles
      const node = createMockNode('A', true, [0, 1])
      const cycles: CircularDependencyInfo[] = [
        {
          cycle: ['A', 'B', 'A'],
          type: 'direct',
          severity: 'warning',
          depth: 2,
          impact: 'Medium',
          complexity: 3,
          priorityScore: 50,
        },
        {
          cycle: ['A', 'C', 'A'],
          type: 'direct',
          severity: 'warning',
          depth: 2,
          impact: 'Medium',
          complexity: 3,
          priorityScore: 50,
        },
      ]

      // WHEN: Computing tooltip data
      const result = computeTooltipData({
        node,
        links: createMockLinks(),
        circularDependencies: cycles,
      })

      // THEN: Should report 2 cycles
      expect(result.cycleInfo?.cycleCount).toBe(2)
    })
  })
})
