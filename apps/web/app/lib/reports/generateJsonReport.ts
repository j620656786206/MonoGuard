import type { ReportData, ReportOptions, ReportResult } from './types'
import { MONOGUARD_VERSION } from './types'

/**
 * Generate a JSON report from analysis data
 * AC2: JSON Export - machine-readable, complete data, pretty-printed
 */
export function generateJsonReport(data: ReportData, options: ReportOptions): ReportResult {
  const report: Record<string, unknown> = {}

  if (options.includeMetadata) {
    report.metadata = {
      ...data.metadata,
      monoguardVersion: MONOGUARD_VERSION,
      generatedAt: new Date().toISOString(),
    }
  }

  if (options.sections.healthScore) {
    report.healthScore = data.healthScore
  }

  if (options.sections.circularDependencies) {
    report.circularDependencies = data.circularDependencies
  }

  if (options.sections.versionConflicts) {
    report.versionConflicts = data.versionConflicts
  }

  if (options.sections.fixRecommendations) {
    report.fixRecommendations = data.fixRecommendations
  }

  const jsonString = JSON.stringify(report, null, 2)
  const blob = new Blob([jsonString], { type: 'application/json' })
  const timestamp = new Date().toISOString().split('T')[0]
  const filename = `${options.projectName}-analysis-report-${timestamp}.json`

  return {
    blob,
    filename,
    format: 'json',
    sizeBytes: blob.size,
  }
}
