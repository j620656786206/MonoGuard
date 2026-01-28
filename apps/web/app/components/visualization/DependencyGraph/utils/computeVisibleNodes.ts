/**
 * computeVisibleNodes - Calculates which nodes and links should be visible
 * based on collapsed node state
 *
 * Handles transitive dependency hiding: when a node is collapsed, its descendants
 * are hidden ONLY if they have no alternate visible path to a root node.
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality
 */

import type { D3Link, D3Node } from '../types'

/**
 * Result of computing visible nodes
 */
export interface ComputeVisibleResult {
  /** Nodes that should be rendered */
  visibleNodes: D3Node[]
  /** Links that should be rendered */
  visibleLinks: D3Link[]
  /** Map of collapsed node ID to count of hidden children */
  hiddenChildCounts: Map<string, number>
}

/**
 * Helper to extract node ID from link source/target (handles both string and node object)
 */
function getLinkNodeId(nodeRef: string | D3Node): string {
  return typeof nodeRef === 'string' ? nodeRef : nodeRef.id
}

/**
 * Get all descendants of a node via DFS
 */
function getDescendants(
  nodeId: string,
  outgoingEdges: Map<string, string[]>,
  visited: Set<string>
): Set<string> {
  const descendants = new Set<string>()
  const children = outgoingEdges.get(nodeId) || []

  for (const child of children) {
    if (visited.has(child)) continue
    visited.add(child)
    descendants.add(child)

    const childDescendants = getDescendants(child, outgoingEdges, visited)
    for (const d of childDescendants) {
      descendants.add(d)
    }
  }

  return descendants
}

/**
 * Check if a node has any visible path to a root (not through collapsed nodes)
 */
function hasVisiblePath(
  nodeId: string,
  collapsedNodeIds: Set<string>,
  incomingEdges: Map<string, string[]>,
  visited: Set<string>
): boolean {
  if (visited.has(nodeId)) return false
  visited.add(nodeId)

  const parents = incomingEdges.get(nodeId) || []

  // If node has no parents, it's a root - it's visible
  if (parents.length === 0) return true

  // Check if any parent provides a visible path
  for (const parent of parents) {
    // If parent is not collapsed, check recursively
    if (!collapsedNodeIds.has(parent)) {
      if (hasVisiblePath(parent, collapsedNodeIds, incomingEdges, visited)) {
        return true
      }
    }
  }

  return false
}

/**
 * Compute which nodes and links should be visible based on collapsed state
 *
 * @param allNodes - All nodes in the graph
 * @param allLinks - All links in the graph
 * @param collapsedNodeIds - Set of node IDs that are collapsed
 * @returns Visible nodes, links, and hidden child counts per collapsed node
 */
export function computeVisibleNodes(
  allNodes: D3Node[],
  allLinks: D3Link[],
  collapsedNodeIds: Set<string>
): ComputeVisibleResult {
  // Handle empty graph
  if (allNodes.length === 0) {
    return {
      visibleNodes: [],
      visibleLinks: [],
      hiddenChildCounts: new Map(),
    }
  }

  // Handle no collapsed nodes
  if (collapsedNodeIds.size === 0) {
    return {
      visibleNodes: allNodes,
      visibleLinks: allLinks,
      hiddenChildCounts: new Map(),
    }
  }

  // Build adjacency lists for efficient traversal
  const outgoingEdges = new Map<string, string[]>()
  const incomingEdges = new Map<string, string[]>()

  for (const link of allLinks) {
    const sourceId = getLinkNodeId(link.source)
    const targetId = getLinkNodeId(link.target)

    const sourceList = outgoingEdges.get(sourceId) ?? []
    sourceList.push(targetId)
    outgoingEdges.set(sourceId, sourceList)

    const targetList = incomingEdges.get(targetId) ?? []
    targetList.push(sourceId)
    incomingEdges.set(targetId, targetList)
  }

  // Track which nodes should be hidden
  const hiddenNodeIds = new Set<string>()
  const hiddenChildCounts = new Map<string, number>()

  // For each collapsed node, hide its descendants that have no other path
  collapsedNodeIds.forEach((collapsedId) => {
    // Start with the collapsed node in visited to prevent cycles from including it as descendant
    const visited = new Set<string>([collapsedId])
    const descendants = getDescendants(collapsedId, outgoingEdges, visited)
    let count = 0

    descendants.forEach((descendantId) => {
      // Check if descendant has any visible path (not through collapsed nodes)
      const hasAlternatePath = hasVisiblePath(
        descendantId,
        collapsedNodeIds,
        incomingEdges,
        new Set()
      )

      if (!hasAlternatePath) {
        hiddenNodeIds.add(descendantId)
        count++
      }
    })

    hiddenChildCounts.set(collapsedId, count)
  })

  // Filter nodes and links
  const visibleNodes = allNodes.filter((node) => !hiddenNodeIds.has(node.id))
  const visibleNodeIds = new Set(visibleNodes.map((n) => n.id))

  const visibleLinks = allLinks.filter((link) => {
    const sourceId = getLinkNodeId(link.source)
    const targetId = getLinkNodeId(link.target)
    return visibleNodeIds.has(sourceId) && visibleNodeIds.has(targetId)
  })

  return { visibleNodes, visibleLinks, hiddenChildCounts }
}
