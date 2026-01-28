/**
 * Styling constants for DependencyGraph visualization
 *
 * Defines colors for nodes and edges in different states:
 * - Normal (default)
 * - Cycle (part of circular dependency)
 * - Selected (currently highlighted)
 * - Dimmed (when another element is selected)
 *
 * @see Story 4.2: Highlight Circular Dependencies in Graph
 */

/**
 * Node color configurations
 */
export const NODE_COLORS = {
  normal: {
    fill: '#6366f1', // indigo-500
    stroke: '#ffffff',
  },
  cycle: {
    fill: '#ef4444', // red-500
    stroke: '#fecaca', // red-200 (glow effect)
  },
  selected: {
    fill: '#dc2626', // red-600
    stroke: '#ffffff',
  },
  dimmed: {
    fill: '#9ca3af', // gray-400
    stroke: '#d1d5db', // gray-300
  },
} as const

/**
 * Edge color and width configurations
 */
export const EDGE_COLORS = {
  normal: {
    stroke: '#94a3b8', // slate-400
    width: 1.5,
    opacity: 0.6,
  },
  cycle: {
    stroke: '#ef4444', // red-500
    width: 2.5,
    opacity: 1,
  },
  selected: {
    stroke: '#dc2626', // red-600
    width: 3,
    opacity: 1,
  },
  dimmed: {
    stroke: '#d1d5db', // gray-300
    width: 0.5,
    opacity: 0.3,
  },
} as const

/**
 * Animation configuration for cycle edges
 */
export const ANIMATION = {
  /** Duration of pulsing animation for cycle edges */
  pulseDuration: '1.5s',
  /** Duration of flowing animation along cycle paths */
  flowDuration: '2s',
  /** Dash array pattern for animated edges */
  dashArray: '10,5',
  /** Dash offset for animation start */
  dashOffset: 15,
} as const

/**
 * Glow filter configuration for cycle nodes
 */
export const GLOW_FILTER = {
  /** Blur radius for glow effect */
  stdDeviation: 3,
  /** Glow color */
  color: '#ef4444', // red-500
  /** Glow opacity */
  opacity: 0.6,
} as const

/**
 * Legend display colors (matching the graph colors)
 */
export const LEGEND_COLORS = {
  normalNode: NODE_COLORS.normal.fill,
  cycleNode: NODE_COLORS.cycle.fill,
  normalEdge: EDGE_COLORS.normal.stroke,
  cycleEdge: EDGE_COLORS.cycle.stroke,
} as const

/**
 * CSS class names for styling
 */
export const CSS_CLASSES = {
  node: 'dependency-node',
  nodeCircle: 'dependency-node-circle',
  nodeLabel: 'dependency-node-label',
  nodeInCycle: 'dependency-node--cycle',
  nodeSelected: 'dependency-node--selected',
  nodeDimmed: 'dependency-node--dimmed',
  nodeCollapsed: 'dependency-node--collapsed',
  edge: 'dependency-edge',
  edgeInCycle: 'dependency-edge--cycle',
  edgeSelected: 'dependency-edge--selected',
  edgeDimmed: 'dependency-edge--dimmed',
  edgeAnimated: 'dependency-edge--animated',
} as const

/**
 * Collapsed node styling (Story 4.3)
 */
export const COLLAPSED_STYLES = {
  node: {
    fill: '#6366f1', // indigo-500
    stroke: '#a5b4fc', // indigo-300
    strokeDasharray: '4,2', // Dashed border
    strokeWidth: 2,
  },
  badge: {
    fill: '#f97316', // orange-500
    textFill: '#ffffff',
    fontSize: '8px',
    fontWeight: 'bold',
    radius: 8,
    offsetX: 12,
    offsetY: -12,
  },
} as const

/**
 * Animation configuration for expand/collapse (Story 4.3)
 * Must complete in < 300ms per AC5
 */
export const EXPAND_COLLAPSE_ANIMATION = {
  duration: 250, // < 300ms as per AC5
  easing: 'easeCubicOut',
  alphaDecay: 0.02, // Slower decay for smoother settling
  alphaTarget: 0,
  velocityDecay: 0.4,
} as const
