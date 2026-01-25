/**
 * GraphLegend - Color legend for dependency graph visualization
 *
 * Displays color coding explanation for nodes and edges,
 * helping users understand the visual representation.
 *
 * @see Story 4.2: Highlight Circular Dependencies in Graph - AC4
 */

import React from 'react'
import { LEGEND_COLORS } from './styles'

/**
 * Props for GraphLegend component
 */
export interface GraphLegendProps {
  /** Position of the legend (default: bottom-left) */
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right'
  /** Whether circular dependencies exist in the graph */
  hasCycles?: boolean
  /** Optional additional class name */
  className?: string
}

/**
 * Legend item component for a single color/meaning pair
 */
interface LegendItemProps {
  color: string
  label: string
  type: 'node' | 'edge'
  glow?: boolean
}

function LegendItem({ color, label, type, glow }: LegendItemProps) {
  if (type === 'node') {
    return (
      <div className="flex items-center gap-2">
        <div
          className="h-3 w-3 rounded-full"
          style={{
            backgroundColor: color,
            boxShadow: glow ? `0 0 4px ${color}` : undefined,
          }}
        />
        <span>{label}</span>
      </div>
    )
  }

  return (
    <div className="flex items-center gap-2">
      <div
        className="h-0.5 w-6"
        style={{
          backgroundColor: color,
          height: type === 'edge' && glow ? '3px' : '2px',
        }}
      />
      <span>{label}</span>
    </div>
  )
}

/**
 * GraphLegend component
 *
 * Displays a legend showing the meaning of different colors used in the graph.
 * Only shows cycle-related legend items when cycles exist.
 */
export const GraphLegend = React.memo(function GraphLegend({
  position = 'bottom-left',
  hasCycles = false,
  className = '',
}: GraphLegendProps) {
  // Position classes
  const positionClasses = {
    'top-left': 'top-4 left-4',
    'top-right': 'top-4 right-4',
    'bottom-left': 'bottom-4 left-4',
    'bottom-right': 'bottom-4 right-4',
  }

  return (
    <div
      className={`absolute ${positionClasses[position]} rounded-lg bg-white/90 p-3 text-xs shadow-lg dark:bg-gray-800/90 ${className}`}
    >
      <div className="mb-2 font-semibold text-gray-900 dark:text-gray-100">Legend</div>
      <div className="space-y-1.5 text-gray-700 dark:text-gray-300">
        {/* Node legends */}
        <LegendItem color={LEGEND_COLORS.normalNode} label="Normal Package" type="node" />
        {hasCycles && (
          <LegendItem
            color={LEGEND_COLORS.cycleNode}
            label="In Circular Dependency"
            type="node"
            glow
          />
        )}

        {/* Edge legends */}
        <LegendItem color={LEGEND_COLORS.normalEdge} label="Normal Dependency" type="edge" />
        {hasCycles && (
          <LegendItem
            color={LEGEND_COLORS.cycleEdge}
            label="Circular Dependency"
            type="edge"
            glow
          />
        )}
      </div>

      {/* Interaction hint when cycles exist */}
      {hasCycles && (
        <div className="mt-2 border-t border-gray-200 pt-2 text-[10px] text-gray-500 dark:border-gray-600 dark:text-gray-400">
          Click on red nodes to highlight specific cycles.
          <br />
          Press Escape to clear selection.
        </div>
      )}
    </div>
  )
})

export default GraphLegend
