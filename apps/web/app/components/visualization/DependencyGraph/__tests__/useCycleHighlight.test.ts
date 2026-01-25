/**
 * Tests for useCycleHighlight hook
 *
 * @see Story 4.2: Highlight Circular Dependencies in Graph
 */

import type { CircularDependencyInfo } from '@monoguard/types'
import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { useCycleHighlight } from '../useCycleHighlight'

/**
 * Mock circular dependencies for testing
 */
const mockCycles: CircularDependencyInfo[] = [
  {
    cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
    type: 'indirect',
    severity: 'warning',
    depth: 3,
    impact: 'Build may fail',
    complexity: 5,
    priorityScore: 75,
  },
  {
    cycle: ['pkg-x', 'pkg-y', 'pkg-x'],
    type: 'direct',
    severity: 'critical',
    depth: 2,
    impact: 'Immediate failure',
    complexity: 3,
    priorityScore: 90,
  },
]

describe('useCycleHighlight', () => {
  describe('Cycle Node Detection', () => {
    it('should identify all nodes in cycles', () => {
      const { result } = renderHook(() => useCycleHighlight(mockCycles))

      // Cycle 1: pkg-a, pkg-b, pkg-c (pkg-a appears twice but is one node)
      expect(result.current.cycleNodeIds.has('pkg-a')).toBe(true)
      expect(result.current.cycleNodeIds.has('pkg-b')).toBe(true)
      expect(result.current.cycleNodeIds.has('pkg-c')).toBe(true)

      // Cycle 2: pkg-x, pkg-y
      expect(result.current.cycleNodeIds.has('pkg-x')).toBe(true)
      expect(result.current.cycleNodeIds.has('pkg-y')).toBe(true)

      // Non-cycle node
      expect(result.current.cycleNodeIds.has('pkg-z')).toBe(false)
    })

    it('should return correct cycle IDs for each node', () => {
      const { result } = renderHook(() => useCycleHighlight(mockCycles))

      // pkg-a is only in cycle 0
      expect(result.current.getNodeCycleIds('pkg-a')).toEqual([0])

      // pkg-x is only in cycle 1
      expect(result.current.getNodeCycleIds('pkg-x')).toEqual([1])

      // Non-cycle node returns empty array
      expect(result.current.getNodeCycleIds('pkg-z')).toEqual([])
    })
  })

  describe('Cycle Edge Detection', () => {
    it('should identify all edges in cycles', () => {
      const { result } = renderHook(() => useCycleHighlight(mockCycles))

      // Cycle 1 edges: pkg-a->pkg-b, pkg-b->pkg-c, pkg-c->pkg-a
      expect(result.current.cycleEdges.has('pkg-a->pkg-b')).toBe(true)
      expect(result.current.cycleEdges.has('pkg-b->pkg-c')).toBe(true)
      expect(result.current.cycleEdges.has('pkg-c->pkg-a')).toBe(true)

      // Cycle 2 edges: pkg-x->pkg-y, pkg-y->pkg-x
      expect(result.current.cycleEdges.has('pkg-x->pkg-y')).toBe(true)
      expect(result.current.cycleEdges.has('pkg-y->pkg-x')).toBe(true)

      // Non-cycle edge
      expect(result.current.cycleEdges.has('pkg-z->pkg-w')).toBe(false)
    })

    it('should return correct cycle IDs for each edge', () => {
      const { result } = renderHook(() => useCycleHighlight(mockCycles))

      // Edge in cycle 0
      expect(result.current.getEdgeCycleIds('pkg-a', 'pkg-b')).toEqual([0])

      // Edge in cycle 1
      expect(result.current.getEdgeCycleIds('pkg-x', 'pkg-y')).toEqual([1])

      // Non-cycle edge returns empty array
      expect(result.current.getEdgeCycleIds('pkg-z', 'pkg-w')).toEqual([])
    })
  })

  describe('Cycle Retrieval', () => {
    it('should return correct cycle by ID', () => {
      const { result } = renderHook(() => useCycleHighlight(mockCycles))

      expect(result.current.getCycleById(0)).toEqual(mockCycles[0])
      expect(result.current.getCycleById(1)).toEqual(mockCycles[1])
    })

    it('should return undefined for invalid cycle ID', () => {
      const { result } = renderHook(() => useCycleHighlight(mockCycles))

      expect(result.current.getCycleById(2)).toBeUndefined()
      expect(result.current.getCycleById(-1)).toBeUndefined()
    })
  })

  describe('Empty/Undefined Input', () => {
    it('should handle undefined circular dependencies', () => {
      const { result } = renderHook(() => useCycleHighlight(undefined))

      expect(result.current.cycleNodeIds.size).toBe(0)
      expect(result.current.cycleEdges.size).toBe(0)
      expect(result.current.getCycleById(0)).toBeUndefined()
      expect(result.current.getNodeCycleIds('any')).toEqual([])
      expect(result.current.getEdgeCycleIds('any', 'other')).toEqual([])
    })

    it('should handle empty circular dependencies array', () => {
      const { result } = renderHook(() => useCycleHighlight([]))

      expect(result.current.cycleNodeIds.size).toBe(0)
      expect(result.current.cycleEdges.size).toBe(0)
    })
  })

  describe('Overlapping Cycles', () => {
    it('should handle nodes that belong to multiple cycles', () => {
      const overlappingCycles: CircularDependencyInfo[] = [
        {
          cycle: ['shared', 'pkg-a', 'shared'],
          type: 'direct',
          severity: 'warning',
          depth: 2,
          impact: 'Minor',
          complexity: 2,
          priorityScore: 50,
        },
        {
          cycle: ['shared', 'pkg-b', 'shared'],
          type: 'direct',
          severity: 'warning',
          depth: 2,
          impact: 'Minor',
          complexity: 2,
          priorityScore: 50,
        },
      ]

      const { result } = renderHook(() => useCycleHighlight(overlappingCycles))

      // 'shared' appears in both cycles
      expect(result.current.cycleNodeIds.has('shared')).toBe(true)
      expect(result.current.getNodeCycleIds('shared')).toContain(0)
      expect(result.current.getNodeCycleIds('shared')).toContain(1)
    })
  })

  describe('Memoization', () => {
    it('should return same reference for same input', () => {
      const { result, rerender } = renderHook(({ deps }) => useCycleHighlight(deps), {
        initialProps: { deps: mockCycles },
      })

      const firstResult = result.current

      // Re-render with same reference
      rerender({ deps: mockCycles })

      expect(result.current).toBe(firstResult)
    })

    it('should update when input changes', () => {
      const { result, rerender } = renderHook(({ deps }) => useCycleHighlight(deps), {
        initialProps: { deps: mockCycles },
      })

      const firstNodeCount = result.current.cycleNodeIds.size

      // Re-render with different data
      const newCycles: CircularDependencyInfo[] = [
        {
          cycle: ['new-a', 'new-b', 'new-a'],
          type: 'direct',
          severity: 'info',
          depth: 2,
          impact: 'Minimal',
          complexity: 1,
          priorityScore: 30,
        },
      ]
      rerender({ deps: newCycles })

      expect(result.current.cycleNodeIds.has('new-a')).toBe(true)
      expect(result.current.cycleNodeIds.has('pkg-a')).toBe(false)
      expect(result.current.cycleNodeIds.size).not.toBe(firstNodeCount)
    })
  })
})
