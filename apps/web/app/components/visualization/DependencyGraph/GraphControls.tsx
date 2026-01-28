/**
 * GraphControls - Depth-based expand/collapse controls for dependency graph
 *
 * Provides UI controls for collapsing/expanding nodes at specific depth levels.
 *
 * @see Story 4.3: Implement Node Expand/Collapse Functionality (AC3)
 */

'use client'

/**
 * Props for GraphControls component
 */
export interface GraphControlsProps {
  /** Current depth level (1-5 or 'all') */
  currentDepth: number | 'all'
  /** Maximum depth in the graph */
  maxDepth: number
  /** Callback when depth changes */
  onDepthChange: (depth: number | 'all') => void
  /** Callback to expand all nodes */
  onExpandAll: () => void
  /** Callback to collapse all nodes */
  onCollapseAll: () => void
}

/**
 * GraphControls component
 *
 * Renders depth-based controls for the dependency graph visualization.
 */
export function GraphControls({
  currentDepth,
  maxDepth,
  onDepthChange,
  onExpandAll,
  onCollapseAll,
}: GraphControlsProps) {
  // Generate depth options: 'all' plus depth levels 1 through min(maxDepth, 5)
  const depthOptions: Array<'all' | number> = [
    'all',
    ...Array.from({ length: Math.min(maxDepth, 5) }, (_, i) => i + 1),
  ]

  return (
    <fieldset
      className="absolute top-4 right-4 bg-white/90 dark:bg-gray-800/90 rounded-lg shadow-lg p-3 text-sm space-y-2 z-10 border-0 m-0"
      aria-label="Graph depth controls"
    >
      <legend className="font-semibold text-gray-700 dark:text-gray-200 p-0 float-left w-full mb-2">
        Depth Control
      </legend>

      <div className="flex gap-1 flex-wrap">
        {depthOptions.map((depth) => (
          <button
            key={depth}
            type="button"
            onClick={() => onDepthChange(depth)}
            className={`px-2 py-1 rounded text-xs transition-colors ${
              currentDepth === depth
                ? 'bg-indigo-600 text-white'
                : 'bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200'
            }`}
            aria-pressed={currentDepth === depth}
            aria-label={depth === 'all' ? 'Show all depths' : `Show depth level ${depth}`}
          >
            {depth === 'all' ? 'All' : `L${depth}`}
          </button>
        ))}
      </div>

      <div className="flex gap-2 pt-1 border-t border-gray-200 dark:border-gray-700">
        <button
          type="button"
          onClick={onExpandAll}
          className="px-2 py-1 rounded text-xs bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-200 hover:bg-green-200 dark:hover:bg-green-800 transition-colors"
          aria-label="Expand all nodes"
        >
          Expand All
        </button>
        <button
          type="button"
          onClick={onCollapseAll}
          className="px-2 py-1 rounded text-xs bg-orange-100 dark:bg-orange-900 text-orange-700 dark:text-orange-200 hover:bg-orange-200 dark:hover:bg-orange-800 transition-colors"
          aria-label="Collapse all nodes"
        >
          Collapse All
        </button>
      </div>
    </fieldset>
  )
}
