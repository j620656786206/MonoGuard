import type { CircularDependencyInfo, DependencyGraph } from '@monoguard/types'
import { generateCyclePath } from './sections/cyclePath'
import { generateExecutiveSummary } from './sections/executiveSummary'
import { renderFixStrategies } from './sections/fixStrategies'
import { generateImpactAssessment } from './sections/impactAssessment'
import { findRelatedCycles } from './sections/relatedCycles'
import { renderRootCauseAnalysis } from './sections/rootCauseAnalysis'
import { getDiagnosticHtmlTemplate } from './templates/diagnosticHtmlTemplate'
import type { DiagnosticReport } from './types'
import { MONOGUARD_VERSION } from './types'

/**
 * Options for generating a diagnostic report
 */
export interface GenerateDiagnosticReportOptions {
  cycle: CircularDependencyInfo
  graph: DependencyGraph
  allCycles: CircularDependencyInfo[]
  totalPackages: number
  projectName: string
  isDarkMode?: boolean
}

/**
 * Generate a complete diagnostic report for a circular dependency
 * AC1-8: All diagnostic report sections
 */
export function generateDiagnosticReport(
  options: GenerateDiagnosticReportOptions
): DiagnosticReport {
  const startTime = performance.now()
  const { cycle, graph, allCycles, totalPackages, projectName, isDarkMode = false } = options

  const cycleId = buildCycleId(cycle)
  const reportId = `diag-${cycleId}-${Date.now()}`

  const executiveSummary = generateExecutiveSummary(cycle)
  const cyclePath = generateCyclePath(cycle, isDarkMode)
  const rootCause = renderRootCauseAnalysis(cycle)
  const fixStrategies = renderFixStrategies(cycle)
  const impactAssessment = generateImpactAssessment(cycle, graph, totalPackages)
  const relatedCycles = findRelatedCycles(cycle, allCycles)

  const endTime = performance.now()

  return {
    id: reportId,
    cycleId,
    generatedAt: new Date().toISOString(),
    monoguardVersion: MONOGUARD_VERSION,
    projectName,

    executiveSummary,
    cyclePath,
    rootCause,
    fixStrategies,
    impactAssessment,
    relatedCycles,

    metadata: {
      generatedAt: new Date().toISOString(),
      generationDurationMs: Math.round(endTime - startTime),
      monoguardVersion: MONOGUARD_VERSION,
      projectName,
      analysisConfigHash: 'default',
    },
  }
}

/**
 * Export a diagnostic report as self-contained HTML
 * AC7: PDF-Ready HTML Export
 */
export function exportDiagnosticReportAsHtml(report: DiagnosticReport): {
  blob: Blob
  filename: string
} {
  const html = getDiagnosticHtmlTemplate(report)
  const blob = new Blob([html], { type: 'text/html' })
  const timestamp = new Date().toISOString().split('T')[0]
  const filename = `${report.projectName}-diagnostic-${report.cycleId}-${timestamp}.html`

  return { blob, filename }
}

function buildCycleId(cycle: CircularDependencyInfo): string {
  const packages = cycle.cycle
  const uniquePackages =
    packages.length > 1 && packages[packages.length - 1] === packages[0]
      ? packages.slice(0, -1)
      : packages

  return uniquePackages.map((p) => p.split('/').pop() || p).join('-')
}
