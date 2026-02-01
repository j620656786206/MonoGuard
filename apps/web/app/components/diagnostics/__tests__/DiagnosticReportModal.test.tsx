import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import type { DiagnosticReport } from '../../../lib/diagnostics/types'
import { DiagnosticReportModal } from '../DiagnosticReportModal'

const mockReport: DiagnosticReport = {
  id: 'diag-test-123',
  cycleId: 'pkg-a-pkg-b',
  generatedAt: '2026-01-25T10:00:00Z',
  monoguardVersion: '0.1.0',
  projectName: 'test-project',
  executiveSummary: {
    description: 'Direct circular dependency between pkg-a and pkg-b.',
    severity: 'medium',
    recommendation: 'Use dependency injection to break the cycle.',
    estimatedEffort: 'medium',
    affectedPackagesCount: 3,
    cycleLength: 2,
  },
  cyclePath: {
    nodes: [
      { id: 'pkg-a', name: 'pkg-a', path: '/a', isInCycle: true, position: { x: 100, y: 200 } },
      { id: 'pkg-b', name: 'pkg-b', path: '/b', isInCycle: true, position: { x: 300, y: 200 } },
    ],
    edges: [
      { from: 'pkg-a', to: 'pkg-b', isBreakingPoint: false },
      { from: 'pkg-b', to: 'pkg-a', isBreakingPoint: true },
    ],
    breakingPoint: { fromPackage: 'pkg-b', toPackage: 'pkg-a', reason: 'Least downstream impact.' },
    svgDiagram: '<svg data-testid="cycle-svg"><circle cx="100" cy="100" r="10"/></svg>',
    asciiDiagram: 'pkg-a -> pkg-b -> pkg-a',
  },
  rootCause: {
    explanation: 'pkg-a imports directly from pkg-b.',
    confidenceScore: 90,
    originatingPackage: 'pkg-a',
    originatingReason: 'Primary importer.',
    alternativeCandidates: [],
    codeReferences: [
      { file: 'src/index.ts', line: 1, importStatement: "import { foo } from 'pkg-b'" },
    ],
  },
  fixStrategies: [
    {
      strategy: 'extract-module',
      title: 'Extract Shared Module',
      description: 'Move shared code to new package.',
      suitabilityScore: 8,
      estimatedEffort: 'medium',
      estimatedTime: '1-2 hours',
      pros: ['Clean separation'],
      cons: ['New package'],
      steps: [
        {
          number: 1,
          title: 'Create package',
          description: 'Create shared package',
          isOptional: false,
        },
      ],
      codeSnippets: { before: '', after: '' },
    },
  ],
  impactAssessment: {
    directParticipants: ['pkg-a', 'pkg-b'],
    directParticipantsCount: 2,
    indirectDependents: ['pkg-c'],
    indirectDependentsCount: 1,
    totalAffectedCount: 3,
    percentageOfMonorepo: 50,
    riskLevel: 'medium',
    riskExplanation: 'Medium risk: 3 packages affected.',
    rippleEffectTree: { package: 'Cycle', depth: 0, dependents: [] },
  },
  relatedCycles: [
    {
      cycleId: 'cycle-2',
      sharedPackages: ['pkg-b'],
      overlapPercentage: 50,
      recommendFixTogether: true,
      reason: 'Share pkg-b.',
    },
  ],
  metadata: {
    generatedAt: '2026-01-25T10:00:00Z',
    generationDurationMs: 15,
    monoguardVersion: '0.1.0',
    projectName: 'test-project',
    analysisConfigHash: 'default',
  },
}

describe('DiagnosticReportModal', () => {
  it('should not render when closed', () => {
    render(
      <DiagnosticReportModal
        isOpen={false}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.queryByTestId('diagnostic-modal')).not.toBeInTheDocument()
  })

  it('should render when open with report', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('diagnostic-modal')).toBeInTheDocument()
    expect(screen.getByText('Diagnostic Report')).toBeInTheDocument()
    expect(screen.getByText('pkg-a-pkg-b')).toBeInTheDocument()
  })

  it('should display executive summary section', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('section-executive-summary')).toBeInTheDocument()
    expect(screen.getByText(/Direct circular dependency/)).toBeInTheDocument()
    expect(screen.getAllByText('medium').length).toBeGreaterThan(0)
  })

  it('should display cycle path section', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('section-cycle-path')).toBeInTheDocument()
  })

  it('should display root cause section', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('section-root-cause')).toBeInTheDocument()
    expect(screen.getByText('Confidence: 90%')).toBeInTheDocument()
  })

  it('should display fix strategies section', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('section-fix-strategies')).toBeInTheDocument()
    expect(screen.getByText('Extract Shared Module')).toBeInTheDocument()
  })

  it('should display impact assessment section', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('section-impact-assessment')).toBeInTheDocument()
    expect(screen.getByText('50%')).toBeInTheDocument()
  })

  it('should display related cycles section', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('section-related-cycles')).toBeInTheDocument()
    expect(screen.getByText('cycle-2')).toBeInTheDocument()
  })

  it('should display metadata section', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    expect(screen.getByTestId('section-metadata')).toBeInTheDocument()
  })

  it('should call onClose when close button clicked', () => {
    const onClose = vi.fn()
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={onClose}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    fireEvent.click(screen.getByTestId('close-modal-button'))
    expect(onClose).toHaveBeenCalled()
  })

  it('should call onExportHtml when export button clicked', () => {
    const onExportHtml = vi.fn()
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={onExportHtml}
        report={mockReport}
        isGenerating={false}
      />
    )
    fireEvent.click(screen.getByTestId('export-html-button'))
    expect(onExportHtml).toHaveBeenCalled()
  })

  it('should show generating indicator when isGenerating', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={null}
        isGenerating={true}
      />
    )
    expect(screen.getByTestId('generating-indicator')).toBeInTheDocument()
    expect(screen.getByText('Generating report...')).toBeInTheDocument()
  })

  it('should disable export button when generating', () => {
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={vi.fn()}
        onExportHtml={vi.fn()}
        report={null}
        isGenerating={true}
      />
    )
    const exportBtn = screen.getByTestId('export-html-button')
    expect(exportBtn).toBeDisabled()
  })

  it('should close on backdrop click', () => {
    const onClose = vi.fn()
    render(
      <DiagnosticReportModal
        isOpen={true}
        onClose={onClose}
        onExportHtml={vi.fn()}
        report={mockReport}
        isGenerating={false}
      />
    )
    fireEvent.click(screen.getByTestId('diagnostic-modal'))
    expect(onClose).toHaveBeenCalled()
  })
})
