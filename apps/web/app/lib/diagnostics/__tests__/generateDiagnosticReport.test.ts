import type { CircularDependencyInfo, DependencyGraph } from '@monoguard/types'
import { describe, expect, it } from 'vitest'
import { exportDiagnosticReportAsHtml, generateDiagnosticReport } from '../generateDiagnosticReport'
import { getDiagnosticHtmlTemplate } from '../templates/diagnosticHtmlTemplate'
import { MONOGUARD_VERSION } from '../types'

const mockGraph: DependencyGraph = {
  nodes: {
    'pkg-a': {
      name: 'pkg-a',
      version: '1.0.0',
      path: '/a',
      dependencies: ['pkg-b'],
      devDependencies: [],
      peerDependencies: [],
    },
    'pkg-b': {
      name: 'pkg-b',
      version: '1.0.0',
      path: '/b',
      dependencies: ['pkg-c'],
      devDependencies: [],
      peerDependencies: [],
    },
    'pkg-c': {
      name: 'pkg-c',
      version: '1.0.0',
      path: '/c',
      dependencies: ['pkg-a'],
      devDependencies: [],
      peerDependencies: [],
    },
  },
  edges: [
    { from: 'pkg-a', to: 'pkg-b', type: 'production', versionRange: '^1.0.0' },
    { from: 'pkg-b', to: 'pkg-c', type: 'production', versionRange: '^1.0.0' },
    { from: 'pkg-c', to: 'pkg-a', type: 'production', versionRange: '^1.0.0' },
  ],
  rootPath: '/workspace',
  workspaceType: 'pnpm',
}

const mockCycle: CircularDependencyInfo = {
  cycle: ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a'],
  type: 'indirect',
  severity: 'warning',
  depth: 3,
  impact: 'Moderate impact',
  complexity: 5,
  priorityScore: 6,
  rootCause: {
    originatingPackage: 'pkg-a',
    problematicDependency: { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
    confidence: 85,
    explanation: 'pkg-a creates dependency on pkg-b',
    chain: [
      { from: 'pkg-a', to: 'pkg-b', type: 'production', critical: true },
      { from: 'pkg-b', to: 'pkg-c', type: 'production', critical: false },
      { from: 'pkg-c', to: 'pkg-a', type: 'production', critical: false },
    ],
  },
  fixStrategies: [
    {
      type: 'extract-module',
      name: 'Extract Shared Module',
      description: 'Move shared code to a new package',
      suitability: 8,
      effort: 'medium',
      pros: ['Clean separation', 'Testable'],
      cons: ['New package to maintain'],
      recommended: true,
      targetPackages: ['pkg-a', 'pkg-b'],
    },
  ],
}

describe('generateDiagnosticReport', () => {
  it('should generate a complete report with all sections', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.id).toMatch(/^diag-/)
    expect(report.cycleId).toBeTruthy()
    expect(report.generatedAt).toMatch(/^\d{4}-\d{2}-\d{2}T/)
    expect(report.monoguardVersion).toBe(MONOGUARD_VERSION)
    expect(report.projectName).toBe('test-project')
  })

  it('should include executive summary', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.executiveSummary.description).toBeTruthy()
    expect(report.executiveSummary.severity).toBeTruthy()
    expect(report.executiveSummary.recommendation).toBeTruthy()
    expect(report.executiveSummary.cycleLength).toBe(3)
  })

  it('should include cycle path visualization', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.cyclePath.nodes).toHaveLength(3)
    expect(report.cyclePath.edges).toHaveLength(3)
    expect(report.cyclePath.svgDiagram).toContain('<svg')
    expect(report.cyclePath.asciiDiagram).toBeTruthy()
    expect(report.cyclePath.breakingPoint).toBeDefined()
  })

  it('should include root cause details', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.rootCause.explanation).toBeTruthy()
    expect(report.rootCause.confidenceScore).toBe(85)
    expect(report.rootCause.originatingPackage).toBe('pkg-a')
  })

  it('should include fix strategies', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.fixStrategies).toHaveLength(1)
    expect(report.fixStrategies[0].strategy).toBe('extract-module')
    expect(report.fixStrategies[0].title).toBe('Extract Shared Module')
  })

  it('should include impact assessment', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.impactAssessment.directParticipants).toContain('pkg-a')
    expect(report.impactAssessment.totalAffectedCount).toBeGreaterThan(0)
  })

  it('should include metadata with generation duration', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.metadata.generationDurationMs).toBeGreaterThanOrEqual(0)
    expect(report.metadata.monoguardVersion).toBe(MONOGUARD_VERSION)
    expect(report.metadata.projectName).toBe('test-project')
  })

  it('should support dark mode', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
      isDarkMode: true,
    })

    expect(report.cyclePath.svgDiagram).toContain('#1f2937')
  })

  it('should handle non-repeating cycle array for cycle ID', () => {
    const nonRepeatingCycle: CircularDependencyInfo = {
      ...mockCycle,
      cycle: ['pkg-a', 'pkg-b', 'pkg-c'],
    }
    const report = generateDiagnosticReport({
      cycle: nonRepeatingCycle,
      graph: mockGraph,
      allCycles: [nonRepeatingCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    expect(report.cycleId).toBe('pkg-a-pkg-b-pkg-c')
  })

  it('should handle scoped package names in cycle ID', () => {
    const scopedCycle: CircularDependencyInfo = {
      ...mockCycle,
      cycle: ['@scope/pkg-a', '@scope/pkg-b', '@scope/pkg-a'],
    }
    const report = generateDiagnosticReport({
      cycle: scopedCycle,
      graph: mockGraph,
      allCycles: [scopedCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    // buildCycleId should extract the part after the last /
    expect(report.cycleId).toBe('pkg-a-pkg-b')
  })
})

describe('exportDiagnosticReportAsHtml', () => {
  it('should generate HTML blob and filename', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    const result = exportDiagnosticReportAsHtml(report)
    expect(result.blob).toBeInstanceOf(Blob)
    expect(result.blob.type).toBe('text/html')
    expect(result.filename).toMatch(/^test-project-diagnostic-.*\.html$/)
  })

  it('should produce self-contained HTML', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('<!DOCTYPE html>')
    expect(html).toContain('<style>')
    expect(html).toContain('Executive Summary')
    expect(html).toContain('Cycle Path')
    expect(html).toContain('Root Cause')
    expect(html).toContain('Fix Strategies')
    expect(html).toContain('Impact Assessment')
    expect(html).toContain('Related Cycles')
    expect(html).toContain('Table of Contents')
    // Should NOT reference external stylesheets
    expect(html).not.toContain('link rel="stylesheet"')
  })

  it('should include print-friendly CSS', () => {
    const report = generateDiagnosticReport({
      cycle: mockCycle,
      graph: mockGraph,
      allCycles: [mockCycle],
      totalPackages: 3,
      projectName: 'test-project',
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).toContain('@media print')
    expect(html).toContain('page-break')
  })

  it('should escape HTML in user-provided content', () => {
    const xssCycle: CircularDependencyInfo = {
      ...mockCycle,
      cycle: ['<script>alert(1)</script>', 'pkg-b', '<script>alert(1)</script>'],
    }

    const report = generateDiagnosticReport({
      cycle: xssCycle,
      graph: mockGraph,
      allCycles: [xssCycle],
      totalPackages: 3,
      projectName: '<script>xss</script>',
    })

    const html = getDiagnosticHtmlTemplate(report)

    expect(html).not.toContain('<script>alert')
    expect(html).not.toContain('<script>xss')
    expect(html).toContain('&lt;script&gt;')
  })
})
