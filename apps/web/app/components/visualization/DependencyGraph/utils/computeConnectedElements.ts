/**
 * computeConnectedElements - Utilities for computing connected graph elements
 *
 * Provides functions for calculating which nodes and edges are connected
 * to a given node, as well as computing tooltip data for display.
 *
 * @see Story 4.5: Implement Hover Details and Tooltips (AC1, AC4)
 */

import type { CircularDependencyInfo } from '@monoguard/types'
import type { D3Link, D3Node, TooltipData } from '../types'

/**
 * Result of computing connected elements
 */
export interface ConnectedElements {
  /** Set of node IDs connected to the target node */
  nodeIds: Set<string>
  /** Set of link indices connected to the target node */
  linkIndices: Set<number>
}

/**
 * Helper to extract node ID from link source/target
 * D3 replaces string IDs with node objects during simulation
 */
function getNodeId(nodeOrId: string | D3Node): string {
  return typeof nodeOrId === 'string' ? nodeOrId : nodeOrId.id
}

/**
 * Computes all nodes and links connected to a given node.
 * Includes both incoming (dependencies of this node) and outgoing (dependents).
 *
 * @param nodeId - The ID of the node to find connections for
 * @param links - All links in the graph
 * @returns ConnectedElements containing connected node IDs and link indices
 */
export function computeConnectedElements(nodeId: string, links: D3Link[]): ConnectedElements {
  const nodeIds = new Set<string>([nodeId])
  const linkIndices = new Set<number>()

  links.forEach((link, index) => {
    const sourceId = getNodeId(link.source)
    const targetId = getNodeId(link.target)

    if (sourceId === nodeId || targetId === nodeId) {
      linkIndices.add(index)
      nodeIds.add(sourceId)
      nodeIds.add(targetId)
    }
  })

  return { nodeIds, linkIndices }
}

/**
 * Result of computing dependency counts
 */
export interface DependencyCounts {
  /** Number of incoming dependencies (edges pointing to this node) */
  incoming: number
  /** Number of outgoing dependencies (edges from this node) */
  outgoing: number
}

/**
 * Computes incoming and outgoing dependency counts for a node.
 *
 * @param nodeId - The ID of the node to count dependencies for
 * @param links - All links in the graph
 * @returns Object with incoming and outgoing counts
 */
export function computeDependencyCounts(nodeId: string, links: D3Link[]): DependencyCounts {
  let incoming = 0
  let outgoing = 0

  links.forEach((link) => {
    const sourceId = getNodeId(link.source)
    const targetId = getNodeId(link.target)

    if (targetId === nodeId) {
      incoming++
    }
    if (sourceId === nodeId) {
      outgoing++
    }
  })

  return { incoming, outgoing }
}

/**
 * Parameters for computing tooltip data
 */
export interface ComputeTooltipDataParams {
  /** The node to compute tooltip data for */
  node: D3Node
  /** All links in the graph */
  links: D3Link[]
  /** Circular dependency information from analysis */
  circularDependencies: CircularDependencyInfo[]
}

/**
 * Computes all data needed for the node tooltip.
 *
 * Calculates dependency counts, health contribution, and cycle information
 * for display in the tooltip.
 *
 * @param params - The node, links, and circular dependencies
 * @returns TooltipData ready for display
 */
export function computeTooltipData({
  node,
  links,
  circularDependencies,
}: ComputeTooltipDataParams): TooltipData {
  // Compute dependency counts
  const { incoming, outgoing } = computeDependencyCounts(node.id, links)

  // Check if node is in any cycles
  const cyclesContainingNode = circularDependencies.filter((cycleInfo) =>
    cycleInfo.cycle.includes(node.id)
  )

  const inCycle = cyclesContainingNode.length > 0

  // Calculate health contribution
  // Negative if in cycles, more negative with more connections in cycles
  // Positive if well-structured with reasonable dependency count
  let healthContribution = 0

  if (inCycle) {
    // Each cycle involvement reduces health
    healthContribution = -5 * cyclesContainingNode.length
    // Additional penalty for being hub in cycle
    if (incoming > 3 || outgoing > 3) {
      healthContribution -= 2
    }
  } else {
    // Good node contributes positively
    // More connections = less contribution (coupling concern)
    const totalConnections = incoming + outgoing
    if (totalConnections <= 3) {
      healthContribution = 2
    } else if (totalConnections <= 6) {
      healthContribution = 1
    } else {
      healthContribution = 0 // High coupling
    }
  }

  // Get packages in the same cycle(s), excluding the current node
  const cyclePackages = new Set<string>()
  cyclesContainingNode.forEach((cycleInfo) => {
    cycleInfo.cycle.forEach((pkg: string) => {
      if (pkg !== node.id) {
        cyclePackages.add(pkg)
      }
    })
  })

  return {
    packageName: node.name,
    packagePath: node.path,
    incomingCount: incoming,
    outgoingCount: outgoing,
    healthContribution,
    inCycle,
    cycleInfo: inCycle
      ? {
          cycleCount: cyclesContainingNode.length,
          packages: Array.from(cyclePackages),
        }
      : undefined,
  }
}
