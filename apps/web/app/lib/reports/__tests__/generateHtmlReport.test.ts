import { describe, expect, it } from 'vitest'
import { generateHtmlReport } from '../generateHtmlReport'
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
  format: 'html',
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

describe('generateHtmlReport', () => {
  it('should generate HTML blob', () => {
    const result = generateHtmlReport(mockData, defaultOptions)

    expect(result.blob).toBeInstanceOf(Blob)
    expect(result.blob.type).toBe('text/html')
  })

  it('should be self-contained (no external links)', async () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).not.toMatch(/<link[^>]+href="http/)
    expect(text).not.toMatch(/<script[^>]+src="http/)
    expect(text).toContain('<style>')
  })

  it('should include dark mode media query', async () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('prefers-color-scheme: dark')
  })

  it('should include print-friendly styles', async () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('@media print')
  })

  it('should escape HTML in project name to prevent XSS', async () => {
    const options: ReportOptions = {
      ...defaultOptions,
      projectName: '<script>alert("xss")</script>',
    }

    const result = generateHtmlReport(mockData, options)
    const text = await readBlobAsText(result.blob)

    expect(text).not.toContain('<script>alert')
    expect(text).toContain('&lt;script&gt;')
  })

  it('should include collapsible section script', async () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('classList.toggle')
    expect(text).toContain('collapsed')
  })

  it('should include MonoGuard version in footer', async () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain(`MonoGuard v${MONOGUARD_VERSION}`)
  })

  it('should generate correct filename', () => {
    const result = generateHtmlReport(mockData, defaultOptions)

    expect(result.filename).toMatch(/^test-project-analysis-report-\d{4}-\d{2}-\d{2}\.html$/)
  })

  it('should have correct format field', () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    expect(result.format).toBe('html')
  })

  it('should include valid HTML structure', async () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('<!DOCTYPE html>')
    expect(text).toContain('<html lang="en">')
    expect(text).toContain('<head>')
    expect(text).toContain('</head>')
    expect(text).toContain('<body>')
    expect(text).toContain('</body>')
    expect(text).toContain('</html>')
  })

  it('should include selected sections', async () => {
    const result = generateHtmlReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('Health Score')
    expect(text).toContain('Circular Dependencies')
    expect(text).toContain('Version Conflicts')
    expect(text).toContain('Fix Recommendations')
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

    const result = generateHtmlReport(mockData, options)
    const text = await readBlobAsText(result.blob)

    expect(text).not.toContain('class="health-score')
    expect(text).toContain('Circular Dependencies')
  })
})
