import type { CircularDependencyInfo, DependencyGraph } from '@monoguard/types'
import { act, renderHook } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import * as diagnosticModule from '../../lib/diagnostics/generateDiagnosticReport'
import { useDiagnosticReport } from '../useDiagnosticReport'

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
      dependencies: ['pkg-a'],
      devDependencies: [],
      peerDependencies: [],
    },
  },
  edges: [
    { from: 'pkg-a', to: 'pkg-b', type: 'production', versionRange: '^1.0.0' },
    { from: 'pkg-b', to: 'pkg-a', type: 'production', versionRange: '^1.0.0' },
  ],
  rootPath: '/workspace',
  workspaceType: 'pnpm',
}

const mockCycle: CircularDependencyInfo = {
  cycle: ['pkg-a', 'pkg-b', 'pkg-a'],
  type: 'direct',
  severity: 'warning',
  depth: 2,
  impact: 'Moderate',
  complexity: 3,
  priorityScore: 5,
}

const defaultOptions = {
  graph: mockGraph,
  allCycles: [mockCycle],
  totalPackages: 2,
  projectName: 'test-project',
}

describe('useDiagnosticReport', () => {
  it('should initialize with default state', () => {
    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))
    expect(result.current.state.report).toBeNull()
    expect(result.current.state.isGenerating).toBe(false)
    expect(result.current.state.isModalOpen).toBe(false)
    expect(result.current.state.error).toBeNull()
  })

  it('should generate report and open modal', () => {
    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })

    expect(result.current.state.report).not.toBeNull()
    expect(result.current.state.isGenerating).toBe(false)
    expect(result.current.state.isModalOpen).toBe(true)
    expect(result.current.state.error).toBeNull()
    expect(result.current.state.report?.cycleId).toBeTruthy()
  })

  it('should close modal', () => {
    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })
    expect(result.current.state.isModalOpen).toBe(true)

    act(() => {
      result.current.closeModal()
    })
    expect(result.current.state.isModalOpen).toBe(false)
    // Report should still be available
    expect(result.current.state.report).not.toBeNull()
  })

  it('should export as HTML', () => {
    // Mock URL.createObjectURL and related DOM methods
    const createObjectURLSpy = vi.fn().mockReturnValue('blob:test')
    const revokeObjectURLSpy = vi.fn()
    const clickSpy = vi.fn()
    const appendChildSpy = vi.spyOn(document.body, 'appendChild').mockImplementation((node) => node)
    const removeChildSpy = vi.spyOn(document.body, 'removeChild').mockImplementation((node) => node)

    global.URL.createObjectURL = createObjectURLSpy
    global.URL.revokeObjectURL = revokeObjectURLSpy

    const mockLink = document.createElement('a')
    mockLink.click = clickSpy
    vi.spyOn(document, 'createElement').mockReturnValue(mockLink)

    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })

    act(() => {
      result.current.exportAsHtml()
    })

    expect(createObjectURLSpy).toHaveBeenCalled()
    expect(clickSpy).toHaveBeenCalled()
    expect(revokeObjectURLSpy).toHaveBeenCalled()

    // Cleanup
    appendChildSpy.mockRestore()
    removeChildSpy.mockRestore()
    vi.restoreAllMocks()
  })

  it('should not export when no report', () => {
    const createObjectURLSpy = vi.fn()
    global.URL.createObjectURL = createObjectURLSpy

    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.exportAsHtml()
    })

    expect(createObjectURLSpy).not.toHaveBeenCalled()
    vi.restoreAllMocks()
  })

  it('should include all report sections', () => {
    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })

    const report = result.current.state.report
    expect(report).not.toBeNull()
    expect(report?.executiveSummary).toBeDefined()
    expect(report?.cyclePath).toBeDefined()
    expect(report?.rootCause).toBeDefined()
    expect(report?.fixStrategies).toBeDefined()
    expect(report?.impactAssessment).toBeDefined()
    expect(report?.relatedCycles).toBeDefined()
    expect(report?.metadata).toBeDefined()
  })

  it('should handle error during report generation', () => {
    const spy = vi.spyOn(diagnosticModule, 'generateDiagnosticReport').mockImplementation(() => {
      throw new Error('Generation failed')
    })

    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })

    expect(result.current.state.error).toBe('Generation failed')
    expect(result.current.state.isGenerating).toBe(false)
    expect(result.current.state.isModalOpen).toBe(false)

    spy.mockRestore()
  })

  it('should handle non-Error thrown during report generation', () => {
    const spy = vi.spyOn(diagnosticModule, 'generateDiagnosticReport').mockImplementation(() => {
      throw 'string error'
    })

    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })

    expect(result.current.state.error).toBe('Failed to generate report')
    expect(result.current.state.isGenerating).toBe(false)

    spy.mockRestore()
  })

  it('should handle error during HTML export', () => {
    const exportSpy = vi
      .spyOn(diagnosticModule, 'exportDiagnosticReportAsHtml')
      .mockImplementation(() => {
        throw new Error('Export failed')
      })

    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })
    expect(result.current.state.report).not.toBeNull()

    act(() => {
      result.current.exportAsHtml()
    })

    expect(result.current.state.error).toBe('Export failed')

    exportSpy.mockRestore()
  })

  it('should handle non-Error thrown during HTML export', () => {
    const exportSpy = vi
      .spyOn(diagnosticModule, 'exportDiagnosticReportAsHtml')
      .mockImplementation(() => {
        throw 42
      })

    const { result } = renderHook(() => useDiagnosticReport(defaultOptions))

    act(() => {
      result.current.generateReport(mockCycle)
    })

    act(() => {
      result.current.exportAsHtml()
    })

    expect(result.current.state.error).toBe('Failed to export report')

    exportSpy.mockRestore()
  })
})
