import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { ReportData, ReportOptions } from '../../lib/reports/types'
import { MONOGUARD_VERSION } from '../../lib/reports/types'
import { useReportExport } from '../useReportExport'

// Mock URL.createObjectURL and URL.revokeObjectURL
globalThis.URL.createObjectURL = vi.fn().mockReturnValue('blob:test-url')
globalThis.URL.revokeObjectURL = vi.fn()

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
    ratingThresholds: {
      excellent: 85,
      good: 70,
      fair: 50,
      poor: 30,
      critical: 0,
    },
  },
  circularDependencies: {
    totalCount: 0,
    bySeverity: { critical: 0, high: 0, medium: 0, low: 0 },
    cycles: [],
  },
  versionConflicts: {
    totalCount: 0,
    byRiskLevel: { critical: 0, high: 0, medium: 0, low: 0 },
    conflicts: [],
  },
  fixRecommendations: {
    totalCount: 0,
    quickWins: 0,
    recommendations: [],
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

describe('useReportExport', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('should initialize with default state', () => {
    const { result } = renderHook(() => useReportExport())

    expect(result.current.exportProgress.isExporting).toBe(false)
    expect(result.current.exportProgress.progress).toBe(0)
    expect(result.current.exportProgress.stage).toBe('preparing')
  })

  it('should export JSON report and trigger download', async () => {
    const { result } = renderHook(() => useReportExport())

    await act(async () => {
      await result.current.startExport(mockData, defaultOptions)
    })

    expect(result.current.exportProgress.isExporting).toBe(false)
    expect(result.current.exportProgress.progress).toBe(100)
    expect(result.current.exportProgress.stage).toBe('complete')

    expect(URL.createObjectURL).toHaveBeenCalled()
    expect(URL.revokeObjectURL).toHaveBeenCalledWith('blob:test-url')
  })

  it('should export HTML report', async () => {
    const { result } = renderHook(() => useReportExport())

    await act(async () => {
      await result.current.startExport(mockData, {
        ...defaultOptions,
        format: 'html',
      })
    })

    expect(result.current.exportProgress.stage).toBe('complete')
    expect(URL.createObjectURL).toHaveBeenCalled()
  })

  it('should export Markdown report', async () => {
    const { result } = renderHook(() => useReportExport())

    await act(async () => {
      await result.current.startExport(mockData, {
        ...defaultOptions,
        format: 'markdown',
      })
    })

    expect(result.current.exportProgress.stage).toBe('complete')
    expect(URL.createObjectURL).toHaveBeenCalled()
  })

  it('should reset on cancel', async () => {
    const { result } = renderHook(() => useReportExport())

    act(() => {
      result.current.cancelExport()
    })

    expect(result.current.exportProgress.isExporting).toBe(false)
    expect(result.current.exportProgress.progress).toBe(0)
    expect(result.current.exportProgress.stage).toBe('preparing')
  })

  it('should reset progress on export error', async () => {
    const { result } = renderHook(() => useReportExport())

    await expect(
      act(async () => {
        await result.current.startExport(mockData, {
          ...defaultOptions,
          format: 'invalid' as ReportOptions['format'],
        })
      })
    ).rejects.toThrow('Unknown format')

    expect(result.current.exportProgress.isExporting).toBe(false)
    expect(result.current.exportProgress.progress).toBe(0)
  })
})
