/**
 * GraphMinimap - Miniature overview of the dependency graph
 *
 * Shows a scaled-down version of the graph with viewport indicator.
 * Allows click/drag navigation to quickly move around large graphs.
 * Only displays for graphs with >= 50 nodes (AC5).
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 */
'use client'

import type React from 'react'
import { useCallback, useMemo, useRef } from 'react'

import type { D3Link, D3Node } from './types'
import type { Bounds } from './utils/calculateBounds'

/**
 * Default minimap configuration
 */
const MINIMAP_CONFIG = {
  /** Default width in pixels */
  width: 150,
  /** Default height in pixels */
  height: 100,
  /** Minimum nodes required to show minimap */
  minNodes: 50,
  /** Scale factor for internal padding (0.9 = 90% of space used) */
  scaleFactor: 0.9,
}

/**
 * Props for GraphMinimap component
 */
export interface GraphMinimapProps {
  /** All nodes in the graph */
  nodes: D3Node[]
  /** All links in the graph */
  links: D3Link[]
  /** Current viewport bounds in graph coordinates */
  viewportBounds: Bounds
  /** Total graph bounds */
  graphBounds: Bounds
  /** Callback when user clicks/drags to navigate */
  onNavigate: (x: number, y: number) => void
  /** Minimap width in pixels (default: 150) */
  width?: number
  /** Minimap height in pixels (default: 100) */
  height?: number
}

/**
 * GraphMinimap component
 *
 * Renders a miniature overview of the dependency graph with:
 * - Scaled-down representation of all nodes and links
 * - Viewport indicator showing current view
 * - Click-to-navigate functionality
 *
 * Only renders for graphs with >= 50 nodes (per AC5).
 */
export function GraphMinimap({
  nodes,
  links,
  viewportBounds,
  graphBounds,
  onNavigate,
  width = MINIMAP_CONFIG.width,
  height = MINIMAP_CONFIG.height,
}: GraphMinimapProps) {
  const svgRef = useRef<SVGSVGElement>(null)

  // Calculate scale to fit graph in minimap
  const scale = useMemo(() => {
    if (graphBounds.width === 0 || graphBounds.height === 0) return 1
    const scaleX = width / graphBounds.width
    const scaleY = height / graphBounds.height
    return Math.min(scaleX, scaleY) * MINIMAP_CONFIG.scaleFactor
  }, [graphBounds.width, graphBounds.height, width, height])

  // Calculate offset to center graph in minimap
  const offset = useMemo(
    () => ({
      x: (width - graphBounds.width * scale) / 2 - graphBounds.x * scale,
      y: (height - graphBounds.height * scale) / 2 - graphBounds.y * scale,
    }),
    [width, height, graphBounds, scale]
  )

  // Calculate viewport rectangle in minimap coordinates
  const viewportRect = useMemo(
    () => ({
      x: viewportBounds.x * scale + offset.x,
      y: viewportBounds.y * scale + offset.y,
      width: viewportBounds.width * scale,
      height: viewportBounds.height * scale,
    }),
    [viewportBounds, scale, offset]
  )

  // Handle click/keyboard to navigate
  const handleNavigate = useCallback(
    (clientX: number, clientY: number) => {
      if (!svgRef.current) return

      const rect = svgRef.current.getBoundingClientRect()
      const clickX = clientX - rect.left
      const clickY = clientY - rect.top

      // Convert click to graph coordinates
      const graphX = (clickX - offset.x) / scale
      const graphY = (clickY - offset.y) / scale

      onNavigate(graphX, graphY)
    },
    [offset.x, offset.y, scale, onNavigate]
  )

  const handleClick = useCallback(
    (event: React.MouseEvent<SVGSVGElement>) => {
      handleNavigate(event.clientX, event.clientY)
    },
    [handleNavigate]
  )

  const handleKeyDown = useCallback(
    (event: React.KeyboardEvent<SVGSVGElement>) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault()
        // Navigate to center when activated via keyboard
        if (svgRef.current) {
          const rect = svgRef.current.getBoundingClientRect()
          handleNavigate(rect.left + width / 2, rect.top + height / 2)
        }
      }
    },
    [handleNavigate, width, height]
  )

  // Create a lookup map for nodes by id for link rendering
  const nodeMap = useMemo(() => {
    const map = new Map<string, D3Node>()
    for (const node of nodes) {
      map.set(node.id, node)
    }
    return map
  }, [nodes])

  // AC5: Only show minimap for graphs with >= 50 nodes
  if (nodes.length < MINIMAP_CONFIG.minNodes) {
    return null
  }

  return (
    <div
      className="absolute top-4 left-4 bg-white/90 dark:bg-gray-800/90
                  rounded-lg shadow-lg p-1 border border-gray-200 dark:border-gray-700 z-10"
    >
      <svg
        ref={svgRef}
        width={width}
        height={height}
        className="cursor-pointer"
        onClick={handleClick}
        onKeyDown={handleKeyDown}
        role="img"
        aria-labelledby="minimap-title"
      >
        <title id="minimap-title">
          Graph minimap navigation. Click to navigate to a position in the graph.
        </title>
        {/* Background */}
        <rect width={width} height={height} fill="transparent" />

        {/* Graph content group */}
        <g transform={`translate(${offset.x}, ${offset.y}) scale(${scale})`}>
          {/* Links */}
          {links.map((link) => {
            // Get source and target nodes
            const sourceId = typeof link.source === 'string' ? link.source : link.source.id
            const targetId = typeof link.target === 'string' ? link.target : link.target.id
            const source =
              typeof link.source === 'string' ? nodeMap.get(sourceId) : (link.source as D3Node)
            const target =
              typeof link.target === 'string' ? nodeMap.get(targetId) : (link.target as D3Node)

            if (!source || !target) return null

            return (
              <line
                key={`link-${sourceId}-${targetId}`}
                x1={source.x ?? 0}
                y1={source.y ?? 0}
                x2={target.x ?? 0}
                y2={target.y ?? 0}
                stroke="#9ca3af"
                strokeWidth={0.5 / scale}
                strokeOpacity={0.5}
              />
            )
          })}

          {/* Nodes */}
          {nodes.map((node) => (
            <circle
              key={node.id}
              cx={node.x ?? 0}
              cy={node.y ?? 0}
              r={3 / scale}
              fill={node.inCycle ? '#ef4444' : '#4f46e5'}
            />
          ))}
        </g>

        {/* Viewport indicator */}
        <rect
          x={viewportRect.x}
          y={viewportRect.y}
          width={Math.max(viewportRect.width, 10)}
          height={Math.max(viewportRect.height, 10)}
          fill="rgba(99, 102, 241, 0.2)"
          stroke="#6366f1"
          strokeWidth={1.5}
          rx={2}
        />
      </svg>
    </div>
  )
}
