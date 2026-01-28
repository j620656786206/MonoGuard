/**
 * Tests for useZoomPan hook
 *
 * @see Story 4.4: Add Zoom, Pan, and Navigation Controls
 * @vitest-environment jsdom
 */
import { act, renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { UseZoomPanProps } from '../useZoomPan'
import { useZoomPan } from '../useZoomPan'

// Mock D3
vi.mock('d3', () => ({
  select: vi.fn(() => ({
    transition: vi.fn().mockReturnThis(),
    duration: vi.fn().mockReturnThis(),
    call: vi.fn().mockReturnThis(),
  })),
  zoom: vi.fn(() => ({
    scaleExtent: vi.fn().mockReturnThis(),
    on: vi.fn().mockReturnThis(),
    scaleTo: vi.fn(),
    transform: vi.fn(),
  })),
  zoomIdentity: {
    translate: vi.fn().mockReturnThis(),
    scale: vi.fn().mockReturnThis(),
  },
}))

describe('useZoomPan', () => {
  let mockSvgElement: SVGSVGElement
  let mockContainerElement: SVGGElement

  beforeEach(() => {
    mockSvgElement = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
    mockContainerElement = document.createElementNS('http://www.w3.org/2000/svg', 'g')

    // Mock getBoundingClientRect
    mockSvgElement.getBoundingClientRect = vi.fn(() => ({
      width: 800,
      height: 600,
      top: 0,
      left: 0,
      right: 800,
      bottom: 600,
      x: 0,
      y: 0,
      toJSON: () => ({}),
    }))

    mockContainerElement.getBBox = vi.fn(() => ({
      x: 0,
      y: 0,
      width: 400,
      height: 300,
      toJSON: () => ({}),
    }))
  })

  const createProps = (overrides: Partial<UseZoomPanProps> = {}): UseZoomPanProps => ({
    svgRef: { current: mockSvgElement },
    containerRef: { current: mockContainerElement },
    ...overrides,
  })

  describe('initialization', () => {
    it('should initialize with default zoom state at scale 1', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(result.current.zoomState.scale).toBe(1)
      expect(result.current.zoomState.translateX).toBe(0)
      expect(result.current.zoomState.translateY).toBe(0)
    })

    it('should display zoom percentage as 100% initially', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(result.current.zoomPercent).toBe(100)
    })

    it('should allow zoom in at initial scale', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(result.current.canZoomIn).toBe(true)
    })

    it('should allow zoom out at initial scale', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(result.current.canZoomOut).toBe(true)
    })
  })

  describe('zoom limits (AC7)', () => {
    it('should use default min scale of 0.1 (10%)', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(result.current.minScale).toBe(0.1)
    })

    it('should use default max scale of 4 (400%)', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(result.current.maxScale).toBe(4)
    })

    it('should respect custom min scale', () => {
      const { result } = renderHook(() => useZoomPan(createProps({ minScale: 0.2 })))

      expect(result.current.minScale).toBe(0.2)
    })

    it('should respect custom max scale', () => {
      const { result } = renderHook(() => useZoomPan(createProps({ maxScale: 3 })))

      expect(result.current.maxScale).toBe(3)
    })
  })

  describe('zoom increment', () => {
    it('should use default zoom increment of 0.2 (20%)', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(result.current.zoomIncrement).toBe(0.2)
    })

    it('should respect custom zoom increment', () => {
      const { result } = renderHook(() => useZoomPan(createProps({ zoomIncrement: 0.3 })))

      expect(result.current.zoomIncrement).toBe(0.3)
    })
  })

  describe('zoom state updates', () => {
    it('should update zoom state when handleZoomChange is called', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      act(() => {
        result.current.handleZoomChange({ k: 2, x: 100, y: 50 })
      })

      expect(result.current.zoomState.scale).toBe(2)
      expect(result.current.zoomState.translateX).toBe(100)
      expect(result.current.zoomState.translateY).toBe(50)
    })

    it('should update zoom percentage when scale changes', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      act(() => {
        result.current.handleZoomChange({ k: 1.5, x: 0, y: 0 })
      })

      expect(result.current.zoomPercent).toBe(150)
    })

    it('should round zoom percentage to nearest integer', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      act(() => {
        result.current.handleZoomChange({ k: 1.555, x: 0, y: 0 })
      })

      expect(result.current.zoomPercent).toBe(156)
    })
  })

  describe('canZoomIn/canZoomOut', () => {
    it('should report canZoomIn as false when at max scale', () => {
      const { result } = renderHook(() => useZoomPan(createProps({ maxScale: 4 })))

      act(() => {
        result.current.handleZoomChange({ k: 4, x: 0, y: 0 })
      })

      expect(result.current.canZoomIn).toBe(false)
    })

    it('should report canZoomOut as false when at min scale', () => {
      const { result } = renderHook(() => useZoomPan(createProps({ minScale: 0.1 })))

      act(() => {
        result.current.handleZoomChange({ k: 0.1, x: 0, y: 0 })
      })

      expect(result.current.canZoomOut).toBe(false)
    })

    it('should report both true when within bounds', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      act(() => {
        result.current.handleZoomChange({ k: 2, x: 0, y: 0 })
      })

      expect(result.current.canZoomIn).toBe(true)
      expect(result.current.canZoomOut).toBe(true)
    })
  })

  describe('exposed API functions', () => {
    it('should expose zoomIn function', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(typeof result.current.zoomIn).toBe('function')
    })

    it('should expose zoomOut function', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(typeof result.current.zoomOut).toBe('function')
    })

    it('should expose resetZoom function', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(typeof result.current.resetZoom).toBe('function')
    })

    it('should expose fitToScreen function', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(typeof result.current.fitToScreen).toBe('function')
    })

    it('should expose setZoomBehavior function', () => {
      const { result } = renderHook(() => useZoomPan(createProps()))

      expect(typeof result.current.setZoomBehavior).toBe('function')
    })
  })

  describe('null ref handling', () => {
    it('should not throw when svgRef is null', () => {
      const { result } = renderHook(() =>
        useZoomPan({
          svgRef: { current: null },
          containerRef: { current: mockContainerElement },
        })
      )

      expect(() => {
        act(() => {
          result.current.zoomIn()
        })
      }).not.toThrow()
    })

    it('should not throw when containerRef is null', () => {
      const { result } = renderHook(() =>
        useZoomPan({
          svgRef: { current: mockSvgElement },
          containerRef: { current: null },
        })
      )

      expect(() => {
        act(() => {
          result.current.fitToScreen()
        })
      }).not.toThrow()
    })
  })
})
