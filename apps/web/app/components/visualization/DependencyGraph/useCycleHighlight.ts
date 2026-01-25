/**
 * useCycleHighlight - Hook for managing cycle highlighting state
 *
 * Processes CircularDependencyInfo to identify which nodes and edges
 * are part of cycles, enabling visual highlighting in the graph.
 *
 * @see Story 4.2: Highlight Circular Dependencies in Graph
 */

import type { CircularDependencyInfo } from '@monoguard/types'
import { useMemo } from 'react'

/**
 * Result of cycle highlight processing
 */
export interface CycleHighlightResult {
  /** Set of node IDs that are part of any cycle */
  cycleNodeIds: Set<string>
  /** Map of edge keys ("from->to") to cycle indices */
  cycleEdges: Map<string, number[]>
  /** Retrieve cycle info by index */
  getCycleById: (id: number) => CircularDependencyInfo | undefined
  /** Get all cycle indices for a given node */
  getNodeCycleIds: (nodeId: string) => number[]
  /** Get all cycle indices for a given edge */
  getEdgeCycleIds: (from: string, to: string) => number[]
}

/**
 * Creates a normalized edge key for consistent lookup
 */
function createEdgeKey(from: string, to: string): string {
  return `${from}->${to}`
}

/**
 * Hook to process circular dependencies and provide lookup functions
 *
 * @param circularDependencies - Array of circular dependency info from analysis
 * @returns Processed cycle data for highlighting
 */
export function useCycleHighlight(
  circularDependencies: CircularDependencyInfo[] | undefined
): CycleHighlightResult {
  return useMemo(() => {
    const cycleNodeIds = new Set<string>()
    const cycleEdges = new Map<string, number[]>()
    const nodeToCycles = new Map<string, number[]>()

    if (!circularDependencies || circularDependencies.length === 0) {
      return {
        cycleNodeIds,
        cycleEdges,
        getCycleById: () => undefined,
        getNodeCycleIds: () => [],
        getEdgeCycleIds: () => [],
      }
    }

    circularDependencies.forEach((cycle, cycleIndex) => {
      // Process each node in the cycle
      // Note: cycle.cycle ends with the first package (e.g., [A, B, C, A])
      cycle.cycle.forEach((nodeName) => {
        cycleNodeIds.add(nodeName)

        // Track which cycles each node belongs to
        const existing = nodeToCycles.get(nodeName) || []
        if (!existing.includes(cycleIndex)) {
          nodeToCycles.set(nodeName, [...existing, cycleIndex])
        }
      })

      // Process edges in the cycle (consecutive pairs)
      // For cycle [A, B, C, A], edges are: A->B, B->C, C->A
      for (let i = 0; i < cycle.cycle.length - 1; i++) {
        const from = cycle.cycle[i]
        const to = cycle.cycle[i + 1]
        const edgeKey = createEdgeKey(from, to)

        const existingCycles = cycleEdges.get(edgeKey) || []
        if (!existingCycles.includes(cycleIndex)) {
          cycleEdges.set(edgeKey, [...existingCycles, cycleIndex])
        }
      }
    })

    return {
      cycleNodeIds,
      cycleEdges,
      getCycleById: (id: number) => circularDependencies[id],
      getNodeCycleIds: (nodeId: string) => nodeToCycles.get(nodeId) || [],
      getEdgeCycleIds: (from: string, to: string) => cycleEdges.get(createEdgeKey(from, to)) || [],
    }
  }, [circularDependencies])
}

/**
 * Hook for managing selected cycle state for click-to-highlight feature
 *
 * @see AC5: Click-to-Highlight Cycle
 * @see AC6: Dim Non-Cycle Elements on Selection
 */
export interface CycleSelectionState {
  /** Currently selected cycle index (null if none selected) */
  selectedCycleIndex: number | null
  /** Select a specific cycle by index */
  selectCycle: (index: number) => void
  /** Clear the current selection */
  clearSelection: () => void
  /** Check if a node is in the selected cycle */
  isNodeInSelectedCycle: (nodeId: string) => boolean
  /** Check if an edge is in the selected cycle */
  isEdgeInSelectedCycle: (from: string, to: string) => boolean
}
