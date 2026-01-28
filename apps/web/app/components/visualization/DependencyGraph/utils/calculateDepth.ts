/**
 * calculateNodeDepths - Calculates the depth of each node from root nodes
 *
 * Uses BFS to find minimum depth from any root node (nodes with no incoming edges).
 * For disconnected nodes, assigns maxDepth + 1.
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality
 */

/**
 * Edge representation for depth calculation
 */
export interface DepthEdge {
  source: string
  target: string
}

/**
 * Calculate the depth of each node from root nodes
 *
 * @param nodeIds - All node IDs in the graph
 * @param edges - All edges in the graph
 * @returns Map of node ID to depth (0 = root)
 */
export function calculateNodeDepths(nodeIds: string[], edges: DepthEdge[]): Map<string, number> {
  const depths = new Map<string, number>()

  // Handle empty graph
  if (nodeIds.length === 0) {
    return depths
  }

  // Build adjacency lists
  const outgoingEdges = new Map<string, string[]>()
  const incomingEdges = new Map<string, string[]>()

  for (const edge of edges) {
    const sourceList = outgoingEdges.get(edge.source) ?? []
    sourceList.push(edge.target)
    outgoingEdges.set(edge.source, sourceList)

    const targetList = incomingEdges.get(edge.target) ?? []
    targetList.push(edge.source)
    incomingEdges.set(edge.target, targetList)
  }

  // Find root nodes (no incoming edges)
  const roots = nodeIds.filter(
    (id) => !incomingEdges.has(id) || (incomingEdges.get(id)?.length ?? 0) === 0
  )

  // If no roots (all nodes have incoming edges - cycles), pick first node
  const startNodes = roots.length > 0 ? roots : [nodeIds[0]]

  // BFS to calculate depths
  const queue: Array<{ id: string; depth: number }> = startNodes.map((id) => ({ id, depth: 0 }))
  const visited = new Set<string>()

  while (queue.length > 0) {
    const item = queue.shift()
    if (!item) continue
    const { id, depth } = item

    // Skip if already visited with a shorter or equal path
    if (visited.has(id)) {
      // Update depth if current path is shorter
      const existingDepth = depths.get(id)
      if (existingDepth !== undefined && depth < existingDepth) {
        depths.set(id, depth)
      }
      continue
    }

    visited.add(id)
    depths.set(id, depth)

    // Add children to queue
    const children = outgoingEdges.get(id) || []
    for (const child of children) {
      if (!visited.has(child)) {
        queue.push({ id: child, depth: depth + 1 })
      }
    }
  }

  // Handle disconnected nodes (assign max depth + 1)
  const maxDepth = Math.max(0, ...depths.values())
  for (const nodeId of nodeIds) {
    if (!depths.has(nodeId)) {
      depths.set(nodeId, maxDepth + 1)
    }
  }

  return depths
}
