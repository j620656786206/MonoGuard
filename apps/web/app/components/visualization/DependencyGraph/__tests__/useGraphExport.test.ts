/**
 * Tests for useGraphExport hook
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images
 */
import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { useGraphExport } from '../useGraphExport'

// Mock the export utilities
vi.mock('../utils/exportSvg', () => ({
  exportSvg: vi.fn().mockResolvedValue({
    blob: new Blob(['<svg></svg>'], { type: 'image/svg+xml' }),
    filename: 'test.svg',
    width: 800,
    height: 600,
  }),
}))

vi.mock('../utils/exportPng', () => ({
  exportPng: vi.fn().mockResolvedValue({
    blob: new Blob([''], { type: 'image/png' }),
    filename: 'test.png',
    width: 1600,
    height: 1200,
  }),
}))

vi.mock('../utils/renderLegendForExport', () => ({
  renderLegendSvg: vi.fn().mockReturnValue('<svg>legend</svg>'),
}))

describe('useGraphExport', () => {
  const mockSvgElement = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
  const mockSvgRef = { current: mockSvgElement } as React.RefObject<SVGSVGElement | null>

  beforeEach(() => {
    vi.clearAllMocks()

    // Mock URL and DOM methods for download
    globalThis.URL.createObjectURL = vi.fn().mockReturnValue('blob:test')
    globalThis.URL.revokeObjectURL = vi.fn()

    vi.spyOn(document, 'createElement').mockImplementation((tag: string) => {
      if (tag === 'a') {
        return {
          href: '',
          download: '',
          click: vi.fn(),
          style: {},
        } as unknown as HTMLAnchorElement
      }
      return document.createElementNS('http://www.w3.org/1999/xhtml', tag)
    })

    vi.spyOn(document.body, 'appendChild').mockImplementation((node: Node) => node)
    vi.spyOn(document.body, 'removeChild').mockImplementation((node: Node) => node)
  })

  it('should initialize with not exporting state', () => {
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef,
        projectName: 'test',
        isDarkMode: false,
      })
    )

    expect(result.current.exportProgress.isExporting).toBe(false)
    expect(result.current.exportProgress.progress).toBe(0)
    expect(result.current.exportProgress.stage).toBe('preparing')
  })

  it('should complete SVG export successfully', async () => {
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef,
        projectName: 'test',
        isDarkMode: false,
      })
    )

    await act(async () => {
      await result.current.startExport({
        format: 'svg',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      })
    })

    // After completion, should no longer be exporting
    expect(result.current.exportProgress.isExporting).toBe(false)
    expect(result.current.exportProgress.stage).toBe('complete')
  })

  it('should complete PNG export successfully', async () => {
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef,
        projectName: 'test',
        isDarkMode: false,
      })
    )

    await act(async () => {
      await result.current.startExport({
        format: 'png',
        scope: 'viewport',
        resolution: 2,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      })
    })

    expect(result.current.exportProgress.isExporting).toBe(false)
  })

  it('should throw error when SVG ref is null', async () => {
    const nullRef = { current: null } as React.RefObject<SVGSVGElement | null>
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: nullRef,
        projectName: 'test',
        isDarkMode: false,
      })
    )

    await expect(
      act(async () => {
        await result.current.startExport({
          format: 'svg',
          scope: 'viewport',
          resolution: 1,
          includeLegend: false,
          includeWatermark: false,
          backgroundColor: '#ffffff',
        })
      })
    ).rejects.toThrow('SVG element not found')
  })

  it('should handle cancellation', () => {
    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef,
        projectName: 'test',
        isDarkMode: false,
      })
    )

    act(() => {
      result.current.cancelExport()
    })

    expect(result.current.exportProgress.isExporting).toBe(false)
    expect(result.current.exportProgress.progress).toBe(0)
  })

  it('should pass legend when includeLegend is true', async () => {
    const { renderLegendSvg } = await import('../utils/renderLegendForExport')

    const { result } = renderHook(() =>
      useGraphExport({
        svgRef: mockSvgRef,
        projectName: 'test',
        isDarkMode: true,
      })
    )

    await act(async () => {
      await result.current.startExport({
        format: 'svg',
        scope: 'viewport',
        resolution: 1,
        includeLegend: true,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      })
    })

    expect(renderLegendSvg).toHaveBeenCalledWith(true)
  })
})
