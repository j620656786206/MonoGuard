import type { EffortLevel, FixStrategyType } from '@monoguard/types'

/**
 * MonoGuard version used in diagnostic report metadata
 */
export const MONOGUARD_VERSION = '0.1.0'

/**
 * DiagnosticReport - Complete diagnostic report for a single circular dependency
 * AC: 1-8 - Aggregates all diagnostic sections
 */
export interface DiagnosticReport {
  /** Unique report identifier */
  id: string
  /** Cycle identifier being diagnosed */
  cycleId: string
  /** ISO 8601 generation timestamp */
  generatedAt: string
  /** MonoGuard version */
  monoguardVersion: string
  /** Project/workspace name */
  projectName: string

  executiveSummary: ExecutiveSummary
  cyclePath: CyclePathVisualization
  rootCause: RootCauseDetails
  fixStrategies: FixStrategyGuide[]
  impactAssessment: ImpactAssessmentDetails
  relatedCycles: RelatedCycleInfo[]

  metadata: DiagnosticMetadata
}

/**
 * ExecutiveSummary - Quick overview of the circular dependency
 * AC1: Executive Summary Generation
 */
export interface ExecutiveSummary {
  /** 1-2 sentence description of the cycle */
  description: string
  /** Severity classification */
  severity: 'critical' | 'high' | 'medium' | 'low'
  /** Quick recommendation */
  recommendation: string
  /** Estimated fix effort */
  estimatedEffort: EffortLevel
  /** Number of packages affected */
  affectedPackagesCount: number
  /** Number of packages in the cycle */
  cycleLength: number
}

/**
 * CyclePathVisualization - Visual representation of the cycle path
 * AC2: Complete Cycle Path Visualization
 */
export interface CyclePathVisualization {
  /** Nodes in the cycle diagram */
  nodes: CycleNode[]
  /** Edges connecting nodes */
  edges: CycleEdge[]
  /** Recommended edge to remove */
  breakingPoint: {
    fromPackage: string
    toPackage: string
    reason: string
  }
  /** SVG diagram string */
  svgDiagram: string
  /** ASCII diagram for text/markdown */
  asciiDiagram: string
}

/**
 * CycleNode - A package node in the cycle visualization
 */
export interface CycleNode {
  /** Package identifier */
  id: string
  /** Display name (short) */
  name: string
  /** Full package path */
  path: string
  /** Whether this node is part of the cycle */
  isInCycle: boolean
  /** Position for visualization */
  position: { x: number; y: number }
}

/**
 * CycleEdge - A dependency edge in the cycle visualization
 */
export interface CycleEdge {
  /** Source package */
  from: string
  /** Target package */
  to: string
  /** Whether this is the recommended breaking point */
  isBreakingPoint: boolean
  /** Import statement text */
  importStatement?: string
  /** Source file path */
  filePath?: string
  /** Line number of import */
  lineNumber?: number
}

/**
 * RootCauseDetails - Formatted root cause analysis
 * AC3: Root Cause Analysis Details
 */
export interface RootCauseDetails {
  /** Human-readable root cause explanation */
  explanation: string
  /** Confidence score (0-100) */
  confidenceScore: number
  /** Package identified as the root cause source */
  originatingPackage: string
  /** Why this package is identified as the source */
  originatingReason: string
  /** Alternative candidates (shown when confidence < 80%) */
  alternativeCandidates: {
    package: string
    reason: string
    confidence: number
  }[]
  /** Code references for the imports */
  codeReferences: {
    file: string
    line: number
    importStatement: string
  }[]
}

/**
 * FixStrategyGuide - Formatted fix strategy with step-by-step guide
 * AC4: All Fix Strategies with Full Guides
 */
export interface FixStrategyGuide {
  /** Strategy type identifier */
  strategy: FixStrategyType
  /** Human-readable title */
  title: string
  /** Strategy description */
  description: string
  /** Suitability score (1-10) */
  suitabilityScore: number
  /** Estimated effort level */
  estimatedEffort: EffortLevel
  /** Estimated time */
  estimatedTime: string
  /** Advantages */
  pros: string[]
  /** Disadvantages */
  cons: string[]
  /** Ordered implementation steps */
  steps: FixStrategyStep[]
  /** Before/after code snippets */
  codeSnippets: {
    before: string
    after: string
  }
}

/**
 * FixStrategyStep - A single step in a fix strategy guide
 */
export interface FixStrategyStep {
  /** Step number (1-based) */
  number: number
  /** Step title */
  title: string
  /** Step description */
  description: string
  /** Code snippet for this step */
  codeSnippet?: string
  /** File path to modify */
  filePath?: string
  /** Whether this step is optional */
  isOptional: boolean
}

/**
 * ImpactAssessmentDetails - Formatted impact assessment
 * AC5: Impact Assessment
 */
export interface ImpactAssessmentDetails {
  /** Packages directly in the cycle */
  directParticipants: string[]
  /** Count of direct participants */
  directParticipantsCount: number
  /** Package names of indirect dependents */
  indirectDependents: string[]
  /** Count of indirect dependents */
  indirectDependentsCount: number
  /** Total affected packages (direct + indirect) */
  totalAffectedCount: number
  /** Percentage of monorepo affected */
  percentageOfMonorepo: number
  /** Risk classification */
  riskLevel: 'critical' | 'high' | 'medium' | 'low'
  /** Risk explanation text */
  riskExplanation: string
  /** Tree structure for ripple effect visualization */
  rippleEffectTree: RippleNode
}

/**
 * RippleNode - Node in the ripple effect tree
 */
export interface RippleNode {
  /** Package name */
  package: string
  /** Depth from cycle (0 = cycle itself) */
  depth: number
  /** Packages that depend on this one */
  dependents: RippleNode[]
}

/**
 * RelatedCycleInfo - Information about a related cycle
 * AC6: Related Cycles Detection
 */
export interface RelatedCycleInfo {
  /** Identifier of the related cycle */
  cycleId: string
  /** Packages shared between cycles */
  sharedPackages: string[]
  /** Percentage of overlap */
  overlapPercentage: number
  /** Whether to fix together */
  recommendFixTogether: boolean
  /** Reason for recommendation */
  reason?: string
}

/**
 * DiagnosticMetadata - Report metadata
 * AC8: Report Metadata
 */
export interface DiagnosticMetadata {
  /** ISO 8601 generation timestamp */
  generatedAt: string
  /** Report generation time in ms */
  generationDurationMs: number
  /** MonoGuard version */
  monoguardVersion: string
  /** Project/workspace name */
  projectName: string
  /** Hash of analysis configuration */
  analysisConfigHash: string
}
