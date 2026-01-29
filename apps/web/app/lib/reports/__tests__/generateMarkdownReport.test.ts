import { describe, expect, it } from 'vitest'
import { generateMarkdownReport } from '../generateMarkdownReport'
import type { ReportData, ReportOptions } from '../types'
import { MONOGUARD_VERSION } from '../types'

/**
 * Read blob content as text (compatible with jsdom which may lack Blob.text())
 */
function readBlobAsText(blob: Blob): Promise<string> {
  return new Promise((resolve) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.readAsText(blob)
  })
}

const mockData: ReportData = {
  metadata: {
    generatedAt: '2026-01-25T10:00:00Z',
    monoguardVersion: MONOGUARD_VERSION,
    projectName: 'test-project',
    analysisDuration: 1234,
    packageCount: 50,
    nodeCount: 50,
    edgeCount: 120,
  },
  healthScore: {
    overall: 75,
    breakdown: [{ category: 'Dependencies', score: 80, weight: 40 }],
    rating: 'good',
    ratingThresholds: { excellent: 85, good: 70, fair: 50, poor: 30, critical: 0 },
  },
  circularDependencies: {
    totalCount: 1,
    bySeverity: { critical: 0, high: 1, medium: 0, low: 0 },
    cycles: [{ id: 'c-1', packages: ['a', 'b'], severity: 'warning', type: 'direct' }],
  },
  versionConflicts: {
    totalCount: 1,
    byRiskLevel: { critical: 0, high: 1, medium: 0, low: 0 },
    conflicts: [
      {
        packageName: 'lodash',
        versions: ['4.17.21', '4.17.15'],
        riskLevel: 'warning',
        recommendedVersion: '4.17.21',
      },
    ],
  },
  fixRecommendations: {
    totalCount: 1,
    quickWins: 1,
    recommendations: [
      {
        id: 'fix-1',
        title: 'Extract Module',
        description: 'Extract shared module',
        effort: 'low',
        impact: 'high',
        priority: 9,
        affectedPackages: ['pkg-a'],
      },
    ],
  },
}

const defaultOptions: ReportOptions = {
  format: 'markdown',
  sections: {
    healthScore: true,
    circularDependencies: true,
    versionConflicts: true,
    fixRecommendations: true,
    packageList: false,
    graphSummary: false,
  },
  includeMetadata: true,
  includeTimestamp: true,
  projectName: 'test-project',
}

describe('generateMarkdownReport', () => {
  it('should generate markdown blob', () => {
    const result = generateMarkdownReport(mockData, defaultOptions)

    expect(result.blob).toBeInstanceOf(Blob)
    expect(result.blob.type).toBe('text/markdown')
  })

  it('should start with proper heading', async () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toMatch(/^# test-project - Dependency Analysis Report/)
  })

  it('should include MonoGuard version', async () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain(`MonoGuard v${MONOGUARD_VERSION}`)
  })

  it('should include ISO 8601 timestamp', async () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toMatch(/\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
  })

  it('should generate GFM-compatible tables', async () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    // GFM tables have header separator rows with dashes
    expect(text).toMatch(/\|[-]+\|/)
    expect(text).toContain('| Category | Score | Weight |')
  })

  it('should include metadata table when enabled', async () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('## Report Metadata')
    expect(text).toContain('| Packages Analyzed | 50 |')
    expect(text).toContain('| Analysis Duration | 1234ms |')
  })

  it('should exclude metadata when disabled', async () => {
    const options: ReportOptions = {
      ...defaultOptions,
      includeMetadata: false,
    }

    const result = generateMarkdownReport(mockData, options)
    const text = await readBlobAsText(result.blob)

    expect(text).not.toContain('## Report Metadata')
  })

  it('should generate correct filename', () => {
    const result = generateMarkdownReport(mockData, defaultOptions)

    expect(result.filename).toMatch(/^test-project-analysis-report-\d{4}-\d{2}-\d{2}\.md$/)
  })

  it('should have correct format field', () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    expect(result.format).toBe('markdown')
  })

  it('should include all selected sections', async () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('## Health Score')
    expect(text).toContain('## Circular Dependencies')
    expect(text).toContain('## Version Conflicts')
    expect(text).toContain('## Fix Recommendations')
  })

  it('should exclude unselected sections', async () => {
    const options: ReportOptions = {
      ...defaultOptions,
      sections: {
        ...defaultOptions.sections,
        healthScore: false,
        fixRecommendations: false,
      },
    }

    const result = generateMarkdownReport(mockData, options)
    const text = await readBlobAsText(result.blob)

    expect(text).not.toContain('## Health Score')
    expect(text).toContain('## Circular Dependencies')
    expect(text).not.toContain('## Fix Recommendations')
  })

  it('should use horizontal rule separator', async () => {
    const result = generateMarkdownReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('---')
  })
})
