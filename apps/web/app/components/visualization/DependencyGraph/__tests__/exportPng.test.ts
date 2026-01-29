/**
 * Tests for PNG export utility
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images - AC2, AC4, AC5
 */
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { exportPng } from '../utils/exportPng'

// Mock the exportSvg dependency
vi.mock('../utils/exportSvg', () => ({
  exportSvg: vi.fn().mockResolvedValue({
    blob: new Blob(['<svg></svg>'], { type: 'image/svg+xml' }),
    filename: 'test.svg',
    width: 800,
    height: 600,
  }),
}))

describe('exportPng', () => {
  let mockSvg: SVGSVGElement

  beforeEach(() => {
    vi.clearAllMocks()

    mockSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
    mockSvg.setAttribute('width', '800')
    mockSvg.setAttribute('height', '600')

    // Mock URL.createObjectURL and URL.revokeObjectURL
    globalThis.URL.createObjectURL = vi.fn().mockReturnValue('blob:test-url')
    globalThis.URL.revokeObjectURL = vi.fn()

    // Mock Image with onload behavior
    const mockImage = {
      onload: null as (() => void) | null,
      onerror: null as (() => void) | null,
      src: '',
      set width(_v: number) {
        /* noop */
      },
      set height(_v: number) {
        /* noop */
      },
    }

    vi.spyOn(globalThis, 'Image').mockImplementation(() => {
      // Trigger onload asynchronously when src is set
      const img = mockImage as unknown as HTMLImageElement
      Object.defineProperty(img, 'src', {
        set() {
          setTimeout(() => {
            if (mockImage.onload) mockImage.onload()
          }, 0)
        },
        get() {
          return ''
        },
      })
      return img
    })

    // Mock canvas context
    const mockContext = {
      fillStyle: '',
      fillRect: vi.fn(),
      scale: vi.fn(),
      drawImage: vi.fn(),
    }

    // Mock canvas with toBlob
    const mockCanvas = {
      width: 0,
      height: 0,
      getContext: vi.fn().mockReturnValue(mockContext),
      toBlob: vi.fn().mockImplementation((callback: (blob: Blob | null) => void) => {
        callback(new Blob(['png-data'], { type: 'image/png' }))
      }),
    }

    vi.spyOn(document, 'createElement').mockImplementation((tag: string) => {
      if (tag === 'canvas') {
        return mockCanvas as unknown as HTMLCanvasElement
      }
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
  })

  it('should export PNG with correct filename including resolution suffix', async () => {
    const result = await exportPng({
      svgElement: mockSvg,
      options: {
        format: 'png',
        scope: 'viewport',
        resolution: 2,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'mono-guard',
    })

    expect(result.filename).toMatch(/^mono-guard-dependency-graph-\d{4}-\d{2}-\d{2}@2x\.png$/)
  })

  it('should not include resolution suffix for 1x', async () => {
    const result = await exportPng({
      svgElement: mockSvg,
      options: {
        format: 'png',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'test',
    })

    expect(result.filename).not.toContain('@')
    expect(result.filename).toMatch(/\.png$/)
  })

  it('should scale dimensions by resolution', async () => {
    const result = await exportPng({
      svgElement: mockSvg,
      options: {
        format: 'png',
        scope: 'viewport',
        resolution: 4,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'test',
    })

    expect(result.width).toBe(800 * 4)
    expect(result.height).toBe(600 * 4)
  })

  it('should return a PNG blob', async () => {
    const result = await exportPng({
      svgElement: mockSvg,
      options: {
        format: 'png',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'test',
    })

    expect(result.blob).toBeInstanceOf(Blob)
    expect(result.blob.type).toBe('image/png')
  })

  it('should call onProgress callback', async () => {
    const onProgress = vi.fn()

    await exportPng({
      svgElement: mockSvg,
      options: {
        format: 'png',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'test',
      onProgress,
    })

    expect(onProgress).toHaveBeenCalledWith(10)
    expect(onProgress).toHaveBeenCalledWith(30)
  })

  it('should revoke object URL after export', async () => {
    await exportPng({
      svgElement: mockSvg,
      options: {
        format: 'png',
        scope: 'viewport',
        resolution: 1,
        includeLegend: false,
        includeWatermark: false,
        backgroundColor: '#ffffff',
      },
      projectName: 'test',
    })

    expect(URL.revokeObjectURL).toHaveBeenCalledWith('blob:test-url')
  })
})
