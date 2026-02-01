import type { CircularDependencyInfo } from '@monoguard/types'
import type { RootCauseDetails } from '../types'

/**
 * Format root cause analysis data for the diagnostic report
 * AC3: Root Cause Analysis Details
 */
export function renderRootCauseAnalysis(cycle: CircularDependencyInfo): RootCauseDetails {
  if (!cycle.rootCause) {
    return {
      explanation: 'Root cause analysis not available for this cycle.',
      confidenceScore: 0,
      originatingPackage: cycle.cycle[0],
      originatingReason: 'Unable to determine root cause with available information.',
      alternativeCandidates: [],
      codeReferences: buildCodeReferences(cycle),
    }
  }

  const { rootCause } = cycle
  const alternativeCandidates: RootCauseDetails['alternativeCandidates'] = []

  // Show alternative candidates when confidence < 80%
  if (rootCause.confidence < 80) {
    // Chain edges that are not the originating package are potential alternatives
    for (const edge of rootCause.chain) {
      if (edge.from !== rootCause.originatingPackage && edge.critical) {
        alternativeCandidates.push({
          package: edge.from,
          reason: `Also contributes to the cycle via dependency on ${edge.to}.`,
          confidence: Math.max(0, rootCause.confidence - 20),
        })
      }
    }
  }

  return {
    explanation: rootCause.explanation,
    confidenceScore: rootCause.confidence,
    originatingPackage: rootCause.originatingPackage,
    originatingReason: `This package introduces the problematic dependency from \`${rootCause.problematicDependency.from}\` to \`${rootCause.problematicDependency.to}\`.`,
    alternativeCandidates,
    codeReferences: buildCodeReferences(cycle),
  }
}

function buildCodeReferences(cycle: CircularDependencyInfo): RootCauseDetails['codeReferences'] {
  if (!cycle.importTraces) return []

  return cycle.importTraces.map((trace) => ({
    file: trace.filePath,
    line: trace.lineNumber,
    importStatement: trace.statement,
  }))
}
