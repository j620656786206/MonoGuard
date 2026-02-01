import type React from 'react'
import { useCallback, useEffect } from 'react'
import type { DiagnosticReport } from '../../lib/diagnostics/types'

export interface DiagnosticReportModalProps {
  isOpen: boolean
  onClose: () => void
  onExportHtml: () => void
  report: DiagnosticReport | null
  isGenerating: boolean
}

/**
 * DiagnosticReportModal - Modal for displaying a detailed diagnostic report
 * AC: Modal with sections, HTML export, keyboard navigation
 */
export const DiagnosticReportModal: React.FC<DiagnosticReportModalProps> = ({
  isOpen,
  onClose,
  onExportHtml,
  report,
  isGenerating,
}) => {
  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    },
    [onClose]
  )

  useEffect(() => {
    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown)
      return () => document.removeEventListener('keydown', handleKeyDown)
    }
  }, [isOpen, handleKeyDown])

  if (!isOpen) return null

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      onClick={onClose}
      onKeyDown={(e) => {
        if (e.key === 'Escape') onClose()
      }}
      role="dialog"
      aria-modal="true"
      aria-label="Diagnostic Report"
      data-testid="diagnostic-modal"
    >
      <div
        className="relative mx-4 flex max-h-[90vh] w-full max-w-4xl flex-col rounded-lg bg-white shadow-xl dark:bg-gray-800"
        onClick={(e) => e.stopPropagation()}
        onKeyDown={(e) => e.stopPropagation()}
        role="document"
      >
        {/* Header */}
        <div className="flex items-center justify-between border-b border-gray-200 px-6 py-4 dark:border-gray-700">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100">
            Diagnostic Report
            {report && (
              <span className="ml-2 text-sm font-normal text-gray-500 dark:text-gray-400">
                {report.cycleId}
              </span>
            )}
          </h2>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={() => window.print()}
              disabled={!report || isGenerating}
              className="rounded-md border border-gray-300 px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-100 disabled:cursor-not-allowed disabled:opacity-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
              data-testid="print-button"
            >
              Print
            </button>
            <button
              type="button"
              onClick={onExportHtml}
              disabled={!report || isGenerating}
              className="rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-50"
              data-testid="export-html-button"
            >
              Export HTML
            </button>
            <button
              type="button"
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              aria-label="Close diagnostic report"
              data-testid="close-modal-button"
            >
              <svg
                className="h-5 w-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto px-6 py-4">
          {isGenerating && (
            <div
              className="flex items-center justify-center py-12"
              data-testid="generating-indicator"
            >
              <div className="h-8 w-8 animate-spin rounded-full border-4 border-blue-600 border-t-transparent" />
              <span className="ml-3 text-gray-600 dark:text-gray-400">Generating report...</span>
            </div>
          )}

          {!isGenerating && !report && (
            <div className="py-12 text-center text-gray-500 dark:text-gray-400">
              No report generated yet.
            </div>
          )}

          {!isGenerating && report && (
            <div className="space-y-6">
              {/* Executive Summary */}
              <section data-testid="section-executive-summary">
                <h3 className="mb-3 text-base font-semibold text-gray-900 dark:text-gray-100">
                  Executive Summary
                </h3>
                <div className="rounded-md bg-gray-50 p-4 dark:bg-gray-700/50">
                  <p className="mb-3 text-sm text-gray-700 dark:text-gray-300">
                    {report.executiveSummary.description}
                  </p>
                  <div className="flex flex-wrap gap-3">
                    <SeverityBadge severity={report.executiveSummary.severity} />
                    <span className="inline-block rounded-md bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-800 dark:bg-blue-900/30 dark:text-blue-300">
                      {report.executiveSummary.cycleLength} packages
                    </span>
                    <span className="inline-block rounded-md bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-700 dark:bg-gray-600 dark:text-gray-300">
                      Effort: {report.executiveSummary.estimatedEffort}
                    </span>
                  </div>
                  <p className="mt-3 text-sm text-gray-600 dark:text-gray-400">
                    <strong>Recommendation:</strong> {report.executiveSummary.recommendation}
                  </p>
                </div>
              </section>

              {/* Cycle Path */}
              <section data-testid="section-cycle-path">
                <h3 className="mb-3 text-base font-semibold text-gray-900 dark:text-gray-100">
                  Cycle Path
                </h3>
                <InternalSvg html={report?.cyclePath.svgDiagram} />
                <div className="rounded-md bg-gray-50 p-3 dark:bg-gray-700/50">
                  <p className="text-xs font-medium text-gray-500 dark:text-gray-400">
                    Breaking Point:
                  </p>
                  <p className="text-sm text-gray-700 dark:text-gray-300">
                    <code className="rounded bg-red-100 px-1 dark:bg-red-900/30">
                      {report.cyclePath.breakingPoint.fromPackage}
                    </code>
                    {' -> '}
                    <code className="rounded bg-red-100 px-1 dark:bg-red-900/30">
                      {report.cyclePath.breakingPoint.toPackage}
                    </code>
                  </p>
                  <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    {report.cyclePath.breakingPoint.reason}
                  </p>
                </div>
              </section>

              {/* Root Cause */}
              <section data-testid="section-root-cause">
                <h3 className="mb-3 text-base font-semibold text-gray-900 dark:text-gray-100">
                  Root Cause Analysis
                </h3>
                <div className="rounded-md bg-gray-50 p-4 dark:bg-gray-700/50">
                  <div className="mb-2 flex items-center gap-2">
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                      Confidence: {report.rootCause.confidenceScore}%
                    </span>
                    <span className="text-sm text-gray-500 dark:text-gray-400">
                      | Origin: <code>{report.rootCause.originatingPackage}</code>
                    </span>
                  </div>
                  <p className="text-sm text-gray-700 dark:text-gray-300">
                    {report.rootCause.explanation}
                  </p>
                  {report.rootCause.codeReferences.length > 0 && (
                    <div className="mt-3">
                      <p className="mb-1 text-xs font-medium text-gray-500 dark:text-gray-400">
                        Code References:
                      </p>
                      {report.rootCause.codeReferences.map((ref) => (
                        <div
                          key={`${ref.file}:${ref.line}`}
                          className="text-xs text-gray-600 dark:text-gray-400"
                        >
                          <code>
                            {ref.file}:{ref.line}
                          </code>{' '}
                          - {ref.importStatement}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </section>

              {/* Fix Strategies */}
              {report.fixStrategies.length > 0 && (
                <section data-testid="section-fix-strategies">
                  <h3 className="mb-3 text-base font-semibold text-gray-900 dark:text-gray-100">
                    Fix Strategies ({report.fixStrategies.length})
                  </h3>
                  <div className="space-y-3">
                    {report.fixStrategies.map((fs) => (
                      <div
                        key={`strategy-${fs.strategy}-${fs.title}`}
                        className="rounded-md border border-gray-200 p-4 dark:border-gray-600"
                      >
                        <div className="mb-2 flex items-center justify-between">
                          <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                            {fs.title}
                          </h4>
                          <div className="flex items-center gap-2">
                            <span className="text-xs text-gray-500 dark:text-gray-400">
                              {fs.suitabilityScore}/10
                            </span>
                            <span className="inline-block rounded bg-blue-100 px-1.5 py-0.5 text-xs font-medium text-blue-800 dark:bg-blue-900/30 dark:text-blue-300">
                              {fs.estimatedEffort}
                            </span>
                          </div>
                        </div>
                        <p className="mb-2 text-sm text-gray-600 dark:text-gray-400">
                          {fs.description}
                        </p>
                        {fs.steps.length > 0 && (
                          <details className="mt-2">
                            <summary className="cursor-pointer text-xs font-medium text-blue-600 dark:text-blue-400">
                              View {fs.steps.length} steps
                            </summary>
                            <ol className="mt-2 space-y-1 pl-4 text-xs text-gray-600 dark:text-gray-400">
                              {fs.steps.map((step) => (
                                <li key={`step-${step.number}`}>
                                  <strong>{step.title}:</strong> {step.description}
                                </li>
                              ))}
                            </ol>
                          </details>
                        )}
                      </div>
                    ))}
                  </div>
                </section>
              )}

              {/* Impact Assessment */}
              <section data-testid="section-impact-assessment">
                <h3 className="mb-3 text-base font-semibold text-gray-900 dark:text-gray-100">
                  Impact Assessment
                </h3>
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
                  <MetricCard
                    label="Direct"
                    value={report.impactAssessment.directParticipantsCount}
                  />
                  <MetricCard
                    label="Indirect"
                    value={report.impactAssessment.indirectDependentsCount}
                  />
                  <MetricCard label="Total" value={report.impactAssessment.totalAffectedCount} />
                  <MetricCard
                    label="% of Monorepo"
                    value={`${report.impactAssessment.percentageOfMonorepo}%`}
                  />
                </div>
                <div className="mt-3 rounded-md bg-gray-50 p-3 dark:bg-gray-700/50">
                  <SeverityBadge severity={report.impactAssessment.riskLevel} />
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
                    {report.impactAssessment.riskExplanation}
                  </p>
                </div>
              </section>

              {/* Related Cycles */}
              {report.relatedCycles.length > 0 && (
                <section data-testid="section-related-cycles">
                  <h3 className="mb-3 text-base font-semibold text-gray-900 dark:text-gray-100">
                    Related Cycles ({report.relatedCycles.length})
                  </h3>
                  <div className="space-y-2">
                    {report.relatedCycles.map((rc) => (
                      <div
                        key={rc.cycleId}
                        className="rounded-md border border-gray-200 p-3 dark:border-gray-600"
                      >
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                            {rc.cycleId}
                          </span>
                          <span className="text-xs text-gray-500 dark:text-gray-400">
                            {rc.overlapPercentage}% overlap
                          </span>
                        </div>
                        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                          Shared: {rc.sharedPackages.join(', ')}
                        </p>
                        {rc.recommendFixTogether && (
                          <p className="mt-1 text-xs font-medium text-blue-600 dark:text-blue-400">
                            Recommended to fix together
                          </p>
                        )}
                      </div>
                    ))}
                  </div>
                </section>
              )}

              {/* Metadata */}
              <section
                className="border-t border-gray-200 pt-3 dark:border-gray-700"
                data-testid="section-metadata"
              >
                <p className="text-xs text-gray-400 dark:text-gray-500">
                  Report ID: {report.id} | Generated: {report.metadata.generatedAt} | Duration:{' '}
                  {report.metadata.generationDurationMs}ms | MonoGuard v{report.monoguardVersion}
                </p>
              </section>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function SeverityBadge({ severity }: { severity: string }) {
  const colorMap: Record<string, string> = {
    critical: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
    high: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-300',
    medium: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300',
    low: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
  }

  return (
    <span
      className={`inline-block rounded-full px-2 py-0.5 text-xs font-semibold uppercase ${colorMap[severity] || colorMap.low}`}
    >
      {severity}
    </span>
  )
}

/** Renders internally-generated SVG content (not user input) */
function InternalSvg({ html }: { html?: string }) {
  if (!html) return <div className="mb-3 flex justify-center" />
  // biome-ignore lint/security/noDangerouslySetInnerHtml: SVG is generated internally by cycleSvg.ts, not from user input
  return <div className="mb-3 flex justify-center" dangerouslySetInnerHTML={{ __html: html }} />
}

function MetricCard({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="rounded-md bg-gray-50 p-3 text-center dark:bg-gray-700/50">
      <div className="text-lg font-bold text-gray-900 dark:text-gray-100">{value}</div>
      <div className="text-xs text-gray-500 dark:text-gray-400">{label}</div>
    </div>
  )
}
