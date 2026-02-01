/**
 * Tests for useViewportState hook
 *
 * @see Story 4.9: Implement Hybrid SVG/Canvas Rendering
 * @see AC5: Viewport State Preservation
 */

import { act, renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { DEFAULT_VIEWPORT } from '../types'
import { useViewportState } from '../useViewportState'

describe('useViewportState', () => {
  it('should initialize with default viewport state', () => {
    const { result } = renderHook(() => useViewportState())

    expect(result.current.viewport).toEqual(DEFAULT_VIEWPORT)
    expect(result.current.viewport.zoom).toBe(1)
    expect(result.current.viewport.panX).toBe(0)
    expect(result.current.viewport.panY).toBe(0)
  })

  it('should accept custom initial state', () => {
    const initial = { zoom: 2, panX: 100, panY: -50 }
    const { result } = renderHook(() => useViewportState(initial))

    expect(result.current.viewport).toEqual(initial)
  })

  it('should update viewport via setViewport', () => {
    const { result } = renderHook(() => useViewportState())

    act(() => {
      result.current.setViewport({ zoom: 1.5, panX: 50, panY: 75 })
    })

    expect(result.current.viewport).toEqual({ zoom: 1.5, panX: 50, panY: 75 })
  })

  it('should reset viewport to default values', () => {
    const { result } = renderHook(() => useViewportState({ zoom: 3, panX: 200, panY: -100 }))

    act(() => {
      result.current.resetViewport()
    })

    expect(result.current.viewport).toEqual(DEFAULT_VIEWPORT)
  })

  it('should clamp zoom between 0.1 and 4', () => {
    const { result } = renderHook(() => useViewportState())

    act(() => {
      result.current.setZoom(10)
    })
    expect(result.current.viewport.zoom).toBe(4)

    act(() => {
      result.current.setZoom(0.01)
    })
    expect(result.current.viewport.zoom).toBe(0.1)
  })

  it('should update pan independently', () => {
    const { result } = renderHook(() => useViewportState())

    act(() => {
      result.current.setPan(150, -75)
    })

    expect(result.current.viewport.panX).toBe(150)
    expect(result.current.viewport.panY).toBe(-75)
    expect(result.current.viewport.zoom).toBe(1) // zoom unchanged
  })

  it('should support functional updates via setViewport', () => {
    const { result } = renderHook(() => useViewportState({ zoom: 1, panX: 0, panY: 0 }))

    act(() => {
      result.current.setViewport((prev) => ({
        ...prev,
        zoom: prev.zoom * 1.5,
      }))
    })

    expect(result.current.viewport.zoom).toBe(1.5)
  })
})
