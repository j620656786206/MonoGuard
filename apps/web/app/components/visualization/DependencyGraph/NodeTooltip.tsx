/**
 * NodeTooltip - Tooltip component for displaying node details on hover
 *
 * Displays package information when hovering over nodes in the dependency graph.
 * Includes package name, path, dependency counts, health contribution, and cycle info.
 *
 * @see Story 4.5: Implement Hover Details and Tooltips (AC1, AC2, AC3, AC7)
 */
'use client'

import type React from 'react'
import { useEffect, useRef, useState } from 'react'
import type { TooltipData, TooltipPosition } from './types'

/**
 * Props for NodeTooltip component
 */
export interface NodeTooltipProps {
  /** Tooltip data to display, null if no tooltip should be shown */
  data: TooltipData | null
  /** Mouse position for tooltip placement */
  position: { x: number; y: number } | null
  /** Reference to the container element for bounds calculation */
  containerRef: React.RefObject<HTMLDivElement | null>
}

/** Offset from cursor to tooltip edge */
const TOOLTIP_OFFSET = 12

/** Animation duration in milliseconds (must be < 200ms per AC2) */
const ANIMATION_DURATION = 150

/**
 * NodeTooltip component
 *
 * Renders a tooltip with package details when hovering over graph nodes.
 * Automatically repositions to stay within viewport bounds.
 */
export function NodeTooltip({
  data,
  position,
  containerRef,
}: NodeTooltipProps): React.ReactElement | null {
  const tooltipRef = useRef<HTMLDivElement>(null)
  const [calculatedPosition, setCalculatedPosition] = useState<TooltipPosition | null>(null)
  const [isVisible, setIsVisible] = useState(false)

  // Calculate position to keep tooltip in viewport (AC3)
  useEffect(() => {
    if (!data || !position || !tooltipRef.current || !containerRef.current) {
      setIsVisible(false)
      return
    }

    const tooltipRect = tooltipRef.current.getBoundingClientRect()
    const containerRect = containerRef.current.getBoundingClientRect()

    let x = position.x - containerRect.left + TOOLTIP_OFFSET
    let y = position.y - containerRect.top + TOOLTIP_OFFSET
    let placement: TooltipPosition['placement'] = 'right'

    // Adjust if tooltip would clip right edge
    if (x + tooltipRect.width > containerRect.width) {
      x = position.x - containerRect.left - tooltipRect.width - TOOLTIP_OFFSET
      placement = 'left'
    }

    // Adjust if tooltip would clip bottom edge
    if (y + tooltipRect.height > containerRect.height) {
      y = position.y - containerRect.top - tooltipRect.height - TOOLTIP_OFFSET
      placement = placement === 'left' ? 'left' : 'top'
    }

    // Adjust if tooltip would clip left edge
    if (x < 0) {
      x = TOOLTIP_OFFSET
    }

    // Adjust if tooltip would clip top edge
    if (y < 0) {
      y = TOOLTIP_OFFSET
    }

    setCalculatedPosition({ x, y, placement })
    setIsVisible(true)
  }, [data, position, containerRef])

  // Don't render if no data or position
  if (!data || !position) {
    return null
  }

  // Shorten path to last 2 segments for display
  const pathSegments = data.packagePath.split('/')
  const shortPath = pathSegments.length > 2 ? pathSegments.slice(-2).join('/') : data.packagePath

  // Format health contribution with sign
  const formatHealthContribution = (value: number): string => {
    if (value > 0) return `+${value}`
    return String(value)
  }

  // Determine cycle count for display
  const cycleCount = data.cycleInfo?.cycleCount ?? 1

  return (
    <div
      ref={tooltipRef}
      role="tooltip"
      aria-live="polite"
      className={`
        absolute z-50 pointer-events-none
        bg-white dark:bg-gray-800 rounded-lg shadow-xl
        border border-gray-200 dark:border-gray-700
        p-3 min-w-[200px] max-w-[300px]
        transition-opacity ease-out
        ${isVisible ? 'opacity-100' : 'opacity-0'}
      `}
      style={{
        left: calculatedPosition?.x ?? 0,
        top: calculatedPosition?.y ?? 0,
        transitionDuration: `${ANIMATION_DURATION}ms`,
      }}
    >
      {/* Package Name (AC1) */}
      <div className="font-semibold text-gray-900 dark:text-white truncate">{data.packageName}</div>

      {/* Package Path (AC1) */}
      <div className="text-xs text-gray-500 dark:text-gray-400 mb-2 truncate">{shortPath}</div>

      {/* Dependency Counts (AC1) */}
      <div className="flex gap-4 text-sm mb-2">
        <div>
          <span className="text-gray-500 dark:text-gray-400">In:</span>{' '}
          <span className="font-medium text-green-600 dark:text-green-400">
            {data.incomingCount}
          </span>
        </div>
        <div>
          <span className="text-gray-500 dark:text-gray-400">Out:</span>{' '}
          <span className="font-medium text-blue-600 dark:text-blue-400">{data.outgoingCount}</span>
        </div>
      </div>

      {/* Health Contribution (AC1) */}
      <div className="text-sm mb-2">
        <span className="text-gray-500 dark:text-gray-400">Health Impact:</span>{' '}
        <span
          className={`font-medium ${
            data.healthContribution >= 0
              ? 'text-green-600 dark:text-green-400'
              : 'text-red-600 dark:text-red-400'
          }`}
        >
          {formatHealthContribution(data.healthContribution)}
        </span>
      </div>

      {/* Circular Dependency Warning (AC1) */}
      {data.inCycle && (
        <div
          className="flex items-center gap-1 text-sm text-red-600 dark:text-red-400
                      bg-red-50 dark:bg-red-900/20 rounded px-2 py-1"
        >
          <svg
            className="w-4 h-4 flex-shrink-0"
            fill="currentColor"
            viewBox="0 0 20 20"
            aria-hidden="true"
          >
            <path
              fillRule="evenodd"
              d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
              clipRule="evenodd"
            />
          </svg>
          <span>
            In {cycleCount} circular {cycleCount > 1 ? 'dependencies' : 'dependency'}
          </span>
        </div>
      )}
    </div>
  )
}
