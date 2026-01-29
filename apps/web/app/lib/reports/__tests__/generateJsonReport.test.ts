import { describe, expect, it } from 'vitest'
import { generateJsonReport } from '../generateJsonReport'
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
    breakdown: [
      { category: 'Dependencies', score: 80, weight: 40 },
      { category: 'Architecture', score: 70, weight: 30 },
    ],
    rating: 'good',
    ratingThresholds: { excellent: 85, good: 70, fair: 50, poor: 30, critical: 0 },
  },
  circularDependencies: {
    totalCount: 2,
    bySeverity: { critical: 0, high: 1, medium: 1, low: 0 },
    cycles: [{ id: 'c-1', packages: ['a', 'b'], severity: 'warning', type: 'direct' }],
  },
  versionConflicts: {
    totalCount: 3,
    byRiskLevel: { critical: 0, high: 1, medium: 2, low: 0 },
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
    totalCount: 5,
    quickWins: 2,
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
  format: 'json',
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

describe('generateJsonReport', () => {
  it('should generate valid JSON blob', () => {
    const result = generateJsonReport(mockData, defaultOptions)

    expect(result.blob).toBeInstanceOf(Blob)
    expect(result.blob.type).toBe('application/json')
  })

  it('should include selected sections only', async () => {
    const options: ReportOptions = {
      ...defaultOptions,
      sections: {
        ...defaultOptions.sections,
        versionConflicts: false,
      },
    }

    const result = generateJsonReport(mockData, options)
    const text = await readBlobAsText(result.blob)
    const json = JSON.parse(text)

    expect(json.healthScore).toBeDefined()
    expect(json.circularDependencies).toBeDefined()
    expect(json.versionConflicts).toBeUndefined()
    expect(json.fixRecommendations).toBeDefined()
  })

  it('should generate correct filename', () => {
    const result = generateJsonReport(mockData, defaultOptions)

    expect(result.filename).toMatch(/^test-project-analysis-report-\d{4}-\d{2}-\d{2}\.json$/)
  })

  it('should include metadata when enabled', async () => {
    const result = generateJsonReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)
    const json = JSON.parse(text)

    expect(json.metadata).toBeDefined()
    expect(json.metadata.projectName).toBe('test-project')
    expect(json.metadata.monoguardVersion).toBe(MONOGUARD_VERSION)
  })

  it('should exclude metadata when disabled', async () => {
    const options: ReportOptions = {
      ...defaultOptions,
      includeMetadata: false,
    }

    const result = generateJsonReport(mockData, options)
    const text = await readBlobAsText(result.blob)
    const json = JSON.parse(text)

    expect(json.metadata).toBeUndefined()
  })

  it('should format JSON with indentation (pretty-printed)', async () => {
    const result = generateJsonReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)

    expect(text).toContain('\n')
    expect(text).toContain('  ')
  })

  it('should have correct format field', () => {
    const result = generateJsonReport(mockData, defaultOptions)
    expect(result.format).toBe('json')
  })

  it('should report correct size', () => {
    const result = generateJsonReport(mockData, defaultOptions)
    expect(result.sizeBytes).toBe(result.blob.size)
    expect(result.sizeBytes).toBeGreaterThan(0)
  })

  it('should include all sections when all enabled', async () => {
    const result = generateJsonReport(mockData, defaultOptions)
    const text = await readBlobAsText(result.blob)
    const json = JSON.parse(text)

    expect(json.metadata).toBeDefined()
    expect(json.healthScore).toBeDefined()
    expect(json.circularDependencies).toBeDefined()
    expect(json.versionConflicts).toBeDefined()
    expect(json.fixRecommendations).toBeDefined()
  })

  it('should generate empty report when no sections selected', async () => {
    const options: ReportOptions = {
      ...defaultOptions,
      includeMetadata: false,
      sections: {
        healthScore: false,
        circularDependencies: false,
        versionConflicts: false,
        fixRecommendations: false,
        packageList: false,
        graphSummary: false,
      },
    }

    const result = generateJsonReport(mockData, options)
    const text = await readBlobAsText(result.blob)
    const json = JSON.parse(text)

    expect(Object.keys(json)).toHaveLength(0)
  })
})
