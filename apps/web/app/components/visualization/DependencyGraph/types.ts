/**
 * D3-specific types for the DependencyGraph visualization
 *
 * These types extend D3's simulation types for use with our dependency graph data.
 */

import type { DependencyType } from '@monoguard/types'
import type { SimulationLinkDatum, SimulationNodeDatum } from 'd3'

/**
 * D3Node - Node representation for D3 force simulation
 *
 * Extends D3's SimulationNodeDatum to include package-specific data.
 * D3 will add x, y, vx, vy, fx, fy during simulation.
 */
export interface D3Node extends SimulationNodeDatum {
  /** Package identifier (usually the package name) */
  id: string
  /** Display name for the node */
  name: string
  /** Relative path from workspace root */
  path: string
  /** Total number of dependencies (used for node sizing) */
  dependencyCount: number
  /** True if node is part of any circular dependency cycle (Story 4.2) */
  inCycle: boolean
  /** Indices of cycles this node belongs to (Story 4.2) */
  cycleIds: number[]
}

/**
 * D3Link - Link representation for D3 force simulation
 *
 * Extends D3's SimulationLinkDatum to include dependency-specific data.
 */
export interface D3Link extends SimulationLinkDatum<D3Node> {
  /** Source node id or node object (D3 replaces string with node ref) */
  source: string | D3Node
  /** Target node id or node object (D3 replaces string with node ref) */
  target: string | D3Node
  /** Type of dependency relationship */
  type: DependencyType
  /** True if edge is part of any circular dependency cycle (Story 4.2) */
  inCycle: boolean
  /** Indices of cycles this edge belongs to (Story 4.2) */
  cycleIds: number[]
}

/**
 * D3GraphData - Transformed graph data ready for D3 consumption
 */
export interface D3GraphData {
  nodes: D3Node[]
  links: D3Link[]
}

/**
 * DependencyGraphProps - Props for the DependencyGraph component
 */
export interface DependencyGraphProps {
  /** Dependency graph data from analysis results */
  data: import('@monoguard/types').DependencyGraph
  /** Circular dependency information for highlighting (Story 4.2) */
  circularDependencies?: import('@monoguard/types').CircularDependencyInfo[]
  /** Optional class name for styling */
  className?: string
  /** Width of the graph container (default: 100%) */
  width?: number | string
  /** Height of the graph container (default: 500px) */
  height?: number | string
}

/**
 * ForceSimulationConfig - Configuration for D3 force simulation
 */
export interface ForceSimulationConfig {
  /** Distance between linked nodes */
  linkDistance: number
  /** Charge strength (negative = repel, positive = attract) */
  chargeStrength: number
  /** Collision radius for nodes */
  collisionRadius: number
  /** Alpha decay rate for simulation stabilization */
  alphaDecay: number
}

/**
 * Default configuration for force simulation
 */
export const DEFAULT_SIMULATION_CONFIG: ForceSimulationConfig = {
  linkDistance: 100,
  chargeStrength: -200,
  collisionRadius: 30,
  alphaDecay: 0.02,
}

/**
 * TooltipData - Data structure for node tooltip content (Story 4.5)
 */
export interface TooltipData {
  /** Package name for display */
  packageName: string
  /** Package path relative to workspace root */
  packagePath: string
  /** Number of dependencies pointing TO this node */
  incomingCount: number
  /** Number of dependencies this node points TO */
  outgoingCount: number
  /** Impact on overall health score (positive = good, negative = bad) */
  healthContribution: number
  /** Whether this node is part of a circular dependency */
  inCycle: boolean
  /** Cycle information if node is in a cycle */
  cycleInfo?: {
    /** Number of cycles this node is involved in */
    cycleCount: number
    /** Other packages in the same cycle(s) */
    packages: string[]
  }
}

/**
 * TooltipPosition - Calculated position for tooltip with placement hint
 */
export interface TooltipPosition {
  /** X coordinate relative to container */
  x: number
  /** Y coordinate relative to container */
  y: number
  /** Placement direction (used for styling) */
  placement: 'top' | 'bottom' | 'left' | 'right'
}

/**
 * HoverState - Tracks which node is being hovered (Story 4.5)
 */
export interface HoverState {
  /** ID of the currently hovered node, null if none */
  nodeId: string | null
  /** Mouse position for tooltip placement */
  position: { x: number; y: number } | null
}
