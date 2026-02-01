import type { CycleEdge, CycleNode } from '../types'

/**
 * Generate an SVG diagram for a cycle visualization
 * AC2: SVG diagram rendering
 */
export function generateCycleSvg(
  nodes: CycleNode[],
  edges: CycleEdge[],
  isDarkMode: boolean
): string {
  const width = 400
  const height = 400

  const colors = isDarkMode
    ? {
        background: '#1f2937',
        node: '#3b82f6',
        nodeStroke: '#60a5fa',
        text: '#f9fafb',
        edge: '#6b7280',
        breakingEdge: '#ef4444',
        breakingEdgeGlow: '#fca5a5',
      }
    : {
        background: '#ffffff',
        node: '#3b82f6',
        nodeStroke: '#2563eb',
        text: '#1f2937',
        edge: '#9ca3af',
        breakingEdge: '#ef4444',
        breakingEdgeGlow: '#fecaca',
      }

  const nodeElements = nodes
    .map(
      (node) => `
    <g transform="translate(${node.position.x}, ${node.position.y})">
      <circle r="30" fill="${colors.node}" stroke="${colors.nodeStroke}" stroke-width="2"/>
      <text y="5" text-anchor="middle" fill="white" font-size="10" font-weight="500">
        ${escapeXml(node.name.substring(0, 10))}
      </text>
    </g>`
    )
    .join('\n')

  const edgeElements = edges
    .map((edge) => {
      const fromNode = nodes.find((n) => n.id === edge.from)
      const toNode = nodes.find((n) => n.id === edge.to)

      if (!fromNode || !toNode) return ''

      const color = edge.isBreakingPoint ? colors.breakingEdge : colors.edge
      const strokeWidth = edge.isBreakingPoint ? 3 : 2
      const dashArray = edge.isBreakingPoint ? '5,5' : 'none'

      const dx = toNode.position.x - fromNode.position.x
      const dy = toNode.position.y - fromNode.position.y
      const dist = Math.sqrt(dx * dx + dy * dy)

      if (dist === 0) return ''

      const endRatio = (dist - 35) / dist
      const endX = fromNode.position.x + dx * endRatio
      const endY = fromNode.position.y + dy * endRatio

      const startRatio = 35 / dist
      const startX = fromNode.position.x + dx * startRatio
      const startY = fromNode.position.y + dy * startRatio

      const glow = edge.isBreakingPoint
        ? `<line x1="${startX}" y1="${startY}" x2="${endX}" y2="${endY}"
                stroke="${colors.breakingEdgeGlow}" stroke-width="8" opacity="0.5"/>`
        : ''

      return `
      <g class="edge${edge.isBreakingPoint ? ' breaking-point' : ''}">
        ${glow}
        <line x1="${startX}" y1="${startY}" x2="${endX}" y2="${endY}"
              stroke="${color}" stroke-width="${strokeWidth}"
              stroke-dasharray="${dashArray}"
              marker-end="url(#arrowhead${edge.isBreakingPoint ? '-red' : ''})"/>
      </g>`
    })
    .join('\n')

  const legend = `
    <g transform="translate(10, ${height - 50})">
      <line x1="0" y1="0" x2="30" y2="0" stroke="${colors.breakingEdge}" stroke-width="3" stroke-dasharray="5,5"/>
      <text x="40" y="4" fill="${colors.text}" font-size="11">Recommended breaking point</text>
    </g>`

  return `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
  <defs>
    <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
      <polygon points="0 0, 10 3.5, 0 7" fill="${colors.edge}"/>
    </marker>
    <marker id="arrowhead-red" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
      <polygon points="0 0, 10 3.5, 0 7" fill="${colors.breakingEdge}"/>
    </marker>
  </defs>

  <rect width="${width}" height="${height}" fill="${colors.background}"/>

  ${edgeElements}
  ${nodeElements}
  ${legend}
</svg>`
}

function escapeXml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}
