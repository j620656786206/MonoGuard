import type { AnalysisResult, ComprehensiveAnalysisResult, HealthScore } from '@monoguard/types'

/**
 * Report format options for export
 * AC1: Format Selection
 */
export type ReportFormat = 'json' | 'html' | 'markdown'

/**
 * Available report sections that can be toggled
 * AC5-8: Report Content sections
 */
export interface ReportSections {
  healthScore: boolean
  circularDependencies: boolean
  versionConflicts: boolean
  fixRecommendations: boolean
  packageList: boolean
  graphSummary: boolean
}

/**
 * Options for configuring report generation
 * AC1, AC10: Format selection and content customization
 */
export interface ReportOptions {
  format: ReportFormat
  sections: ReportSections
  includeMetadata: boolean
  includeTimestamp: boolean
  projectName: string
}

/**
 * Metadata included in every report
 * AC9: Report Metadata
 */
export interface ReportMetadata {
  /** ISO 8601 generation timestamp */
  generatedAt: string
  /** MonoGuard version string */
  monoguardVersion: string
  /** Name of the project/workspace analyzed */
  projectName: string
  /** Analysis duration in milliseconds */
  analysisDuration: number
  /** Total packages analyzed */
  packageCount: number
  /** Number of graph nodes */
  nodeCount: number
  /** Number of graph edges */
  edgeCount: number
}

/**
 * Health score data structured for report rendering
 * AC5: Health Score Summary
 */
export interface HealthScoreReport {
  overall: number
  breakdown: {
    category: string
    score: number
    weight: number
  }[]
  rating: 'excellent' | 'good' | 'fair' | 'poor' | 'critical'
  ratingThresholds: Record<string, number>
}

/**
 * Circular dependency data structured for report rendering
 * AC6: Circular Dependencies
 */
export interface CircularDependencyReport {
  totalCount: number
  bySeverity: {
    critical: number
    high: number
    medium: number
    low: number
  }
  cycles: {
    id: string
    packages: string[]
    severity: string
    type: 'direct' | 'indirect'
  }[]
}

/**
 * Version conflict data structured for report rendering
 * AC7: Version Conflicts
 */
export interface VersionConflictReport {
  totalCount: number
  byRiskLevel: {
    critical: number
    high: number
    medium: number
    low: number
  }
  conflicts: {
    packageName: string
    versions: string[]
    riskLevel: string
    recommendedVersion: string
  }[]
}

/**
 * Fix recommendation data structured for report rendering
 * AC8: Fix Recommendations
 */
export interface FixRecommendationReport {
  totalCount: number
  quickWins: number
  recommendations: {
    id: string
    title: string
    description: string
    effort: 'low' | 'medium' | 'high'
    impact: 'low' | 'medium' | 'high'
    priority: number
    affectedPackages: string[]
  }[]
}

/**
 * Complete report data aggregating all sections
 */
export interface ReportData {
  metadata: ReportMetadata
  healthScore: HealthScoreReport
  circularDependencies: CircularDependencyReport
  versionConflicts: VersionConflictReport
  fixRecommendations: FixRecommendationReport
}

/**
 * Result of report generation
 */
export interface ReportResult {
  blob: Blob
  filename: string
  format: ReportFormat
  sizeBytes: number
}

/**
 * Default report sections configuration - all major sections enabled
 */
export const DEFAULT_REPORT_SECTIONS: ReportSections = {
  healthScore: true,
  circularDependencies: true,
  versionConflicts: true,
  fixRecommendations: true,
  packageList: false,
  graphSummary: false,
}

/**
 * MonoGuard version used in report metadata
 */
export const MONOGUARD_VERSION = '0.1.0'

/**
 * Rating thresholds for health score classification
 */
export const RATING_THRESHOLDS: Record<string, number> = {
  excellent: 85,
  good: 70,
  fair: 50,
  poor: 30,
  critical: 0,
}

/**
 * Compute health score rating from a numeric score
 */
export function getHealthScoreRating(
  score: number
): 'excellent' | 'good' | 'fair' | 'poor' | 'critical' {
  if (score >= RATING_THRESHOLDS.excellent) return 'excellent'
  if (score >= RATING_THRESHOLDS.good) return 'good'
  if (score >= RATING_THRESHOLDS.fair) return 'fair'
  if (score >= RATING_THRESHOLDS.poor) return 'poor'
  return 'critical'
}

/**
 * Build ReportData from an AnalysisResult
 * Transforms domain types into report-ready structures
 */
export function buildReportData(analysisResult: AnalysisResult, projectName: string): ReportData {
  const metadata: ReportMetadata = {
    generatedAt: new Date().toISOString(),
    monoguardVersion: MONOGUARD_VERSION,
    projectName,
    analysisDuration: analysisResult.metadata?.durationMs ?? 0,
    packageCount: analysisResult.packages,
    nodeCount: analysisResult.graph ? Object.keys(analysisResult.graph.nodes).length : 0,
    edgeCount: analysisResult.graph?.edges.length ?? 0,
  }

  const healthScore = buildHealthScoreReport(analysisResult)
  const circularDependencies = buildCircularDependencyReport(analysisResult)
  const versionConflicts = buildVersionConflictReport(analysisResult)
  const fixRecommendations = buildFixRecommendationReport(analysisResult)

  return {
    metadata,
    healthScore,
    circularDependencies,
    versionConflicts,
    fixRecommendations,
  }
}

function buildHealthScoreReport(analysisResult: AnalysisResult): HealthScoreReport {
  const details = analysisResult.healthScoreDetails
  if (details) {
    return {
      overall: details.overall,
      breakdown: details.factors.map((f) => ({
        category: f.name,
        score: f.score,
        weight: Math.round(f.weight * 100),
      })),
      rating: details.rating,
      ratingThresholds: RATING_THRESHOLDS,
    }
  }

  // Fallback when detailed breakdown not available
  const score = analysisResult.healthScore
  return {
    overall: score,
    breakdown: [{ category: 'Overall', score, weight: 100 }],
    rating: getHealthScoreRating(score),
    ratingThresholds: RATING_THRESHOLDS,
  }
}

function buildCircularDependencyReport(analysisResult: AnalysisResult): CircularDependencyReport {
  const deps = analysisResult.circularDependencies ?? []
  const bySeverity = { critical: 0, high: 0, medium: 0, low: 0 }

  for (const dep of deps) {
    const sev = dep.severity
    if (sev === 'critical') bySeverity.critical++
    else if (sev === 'warning') bySeverity.high++
    else if (sev === 'info') bySeverity.medium++
    else bySeverity.low++
  }

  return {
    totalCount: deps.length,
    bySeverity,
    cycles: deps.map((dep, index) => ({
      id: `cycle-${index + 1}`,
      packages: dep.cycle,
      severity: dep.severity,
      type: dep.type,
    })),
  }
}

function buildVersionConflictReport(analysisResult: AnalysisResult): VersionConflictReport {
  const conflicts = analysisResult.versionConflicts ?? []
  const byRiskLevel = { critical: 0, high: 0, medium: 0, low: 0 }

  for (const conflict of conflicts) {
    const sev = conflict.severity
    if (sev === 'critical') byRiskLevel.critical++
    else if (sev === 'warning') byRiskLevel.high++
    else if (sev === 'info') byRiskLevel.medium++
    else byRiskLevel.low++
  }

  return {
    totalCount: conflicts.length,
    byRiskLevel,
    conflicts: conflicts.map((c) => ({
      packageName: c.packageName,
      versions: c.conflictingVersions.map((v) => v.version),
      riskLevel: c.severity,
      recommendedVersion: c.resolution,
    })),
  }
}

function buildFixRecommendationReport(analysisResult: AnalysisResult): FixRecommendationReport {
  const deps = analysisResult.circularDependencies ?? []
  const recommendations: FixRecommendationReport['recommendations'] = []
  let quickWins = 0

  for (const dep of deps) {
    if (!dep.fixStrategies) continue
    for (const strategy of dep.fixStrategies) {
      const isQuickWin = strategy.effort === 'low' && strategy.suitability >= 7
      if (isQuickWin) quickWins++

      recommendations.push({
        id: `fix-${dep.cycle.join('-')}-${strategy.type}`,
        title: strategy.name,
        description: strategy.description,
        effort: strategy.effort,
        impact: strategy.suitability >= 7 ? 'high' : strategy.suitability >= 4 ? 'medium' : 'low',
        priority: strategy.suitability,
        affectedPackages: strategy.targetPackages,
      })
    }
  }

  // Sort by priority descending
  recommendations.sort((a, b) => b.priority - a.priority)

  return {
    totalCount: recommendations.length,
    quickWins,
    recommendations,
  }
}

/**
 * Build ReportData from ComprehensiveAnalysisResult (used by UI components)
 * Bridges the gap between the UI's data model and the report system
 */
export function buildReportDataFromComprehensive(
  analysis: ComprehensiveAnalysisResult,
  projectName: string
): ReportData {
  const results = analysis.results
  const healthScoreValue =
    typeof results?.healthScore === 'number'
      ? results.healthScore
      : ((results?.healthScore as HealthScore | undefined)?.overall ??
        results?.summary?.healthScore ??
        0)

  const healthScoreObj = results?.healthScore as HealthScore | undefined

  const metadata: ReportMetadata = {
    generatedAt: new Date().toISOString(),
    monoguardVersion: MONOGUARD_VERSION,
    projectName,
    analysisDuration: 0,
    packageCount: results?.summary?.totalPackages ?? 0,
    nodeCount: 0,
    edgeCount: 0,
  }

  const healthScore: HealthScoreReport = healthScoreObj?.factors
    ? {
        overall: healthScoreObj.overall,
        breakdown: healthScoreObj.factors.map((f) => ({
          category: f.name,
          score: f.score,
          weight: Math.round(f.weight * 100),
        })),
        rating: getHealthScoreRating(healthScoreObj.overall),
        ratingThresholds: RATING_THRESHOLDS,
      }
    : {
        overall: healthScoreValue,
        breakdown: [{ category: 'Overall', score: healthScoreValue, weight: 100 }],
        rating: getHealthScoreRating(healthScoreValue),
        ratingThresholds: RATING_THRESHOLDS,
      }

  const circularDeps = results?.circularDependencies ?? []
  const circBySeverity = { critical: 0, high: 0, medium: 0, low: 0 }
  for (const cd of circularDeps) {
    const sev = cd.severity as string
    if (sev === 'critical') circBySeverity.critical++
    else if (sev === 'high') circBySeverity.high++
    else if (sev === 'medium') circBySeverity.medium++
    else circBySeverity.low++
  }
  const circularDependencies: CircularDependencyReport = {
    totalCount: circularDeps.length,
    bySeverity: circBySeverity,
    cycles: circularDeps.map((cd, i) => ({
      id: `cycle-${i + 1}`,
      packages: cd.cycle,
      severity: cd.severity as string,
      type: cd.type,
    })),
  }

  const vcList = results?.versionConflicts ?? []
  const vcByRisk = { critical: 0, high: 0, medium: 0, low: 0 }
  for (const vc of vcList) {
    if (vc.riskLevel === 'critical') vcByRisk.critical++
    else if (vc.riskLevel === 'high') vcByRisk.high++
    else if (vc.riskLevel === 'medium') vcByRisk.medium++
    else vcByRisk.low++
  }
  const versionConflicts: VersionConflictReport = {
    totalCount: vcList.length,
    byRiskLevel: vcByRisk,
    conflicts: vcList.map((vc) => ({
      packageName: vc.packageName,
      versions: vc.conflictingVersions.map((v) => v.version),
      riskLevel: vc.riskLevel as string,
      recommendedVersion: vc.resolution,
    })),
  }

  const fixRecommendations: FixRecommendationReport = {
    totalCount: 0,
    quickWins: 0,
    recommendations: [],
  }

  return {
    metadata,
    healthScore,
    circularDependencies,
    versionConflicts,
    fixRecommendations,
  }
}
