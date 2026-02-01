/**
 * Tests for useRenderMode hook
 *
 * @see Story 4.9: Implement Hybrid SVG/Canvas Rendering
 * @see AC1: Automatic Render Mode Selection
 * @see AC3: User Override in Settings
 */

import { renderHook } from '@testing-library/react'
import { afterEach, describe, expect, it } from 'vitest'
import { useSettingsStore } from '../../../../stores/settings'
import { NODE_THRESHOLD } from '../types'
import { useRenderMode } from '../useRenderMode'

describe('useRenderMode', () => {
  afterEach(() => {
    // Reset to default after each test
    useSettingsStore.setState({ visualizationMode: 'auto' })
  })

  describe('Auto mode (AC1)', () => {
    it('should return SVG mode for small graphs (< NODE_THRESHOLD)', () => {
      const { result } = renderHook(() => useRenderMode(100))

      expect(result.current.mode).toBe('svg')
      expect(result.current.isAutoMode).toBe(true)
      expect(result.current.isForced).toBe(false)
      expect(result.current.shouldShowWarning).toBe(false)
    })

    it('should return Canvas mode for large graphs (>= NODE_THRESHOLD)', () => {
      const { result } = renderHook(() => useRenderMode(NODE_THRESHOLD))

      expect(result.current.mode).toBe('canvas')
      expect(result.current.isAutoMode).toBe(true)
      expect(result.current.isForced).toBe(false)
    })

    it('should use NODE_THRESHOLD of 500 as the switching point', () => {
      expect(NODE_THRESHOLD).toBe(500)

      const { result: below } = renderHook(() => useRenderMode(499))
      expect(below.current.mode).toBe('svg')

      const { result: atThreshold } = renderHook(() => useRenderMode(500))
      expect(atThreshold.current.mode).toBe('canvas')

      const { result: above } = renderHook(() => useRenderMode(1000))
      expect(above.current.mode).toBe('canvas')
    })

    it('should handle zero nodes', () => {
      const { result } = renderHook(() => useRenderMode(0))
      expect(result.current.mode).toBe('svg')
    })
  })

  describe('Force SVG mode (AC3)', () => {
    it('should return SVG mode regardless of node count', () => {
      useSettingsStore.setState({ visualizationMode: 'force-svg' })
      const { result } = renderHook(() => useRenderMode(1000))

      expect(result.current.mode).toBe('svg')
      expect(result.current.isAutoMode).toBe(false)
      expect(result.current.isForced).toBe(true)
    })

    it('should show performance warning for large graphs in forced SVG', () => {
      useSettingsStore.setState({ visualizationMode: 'force-svg' })
      const { result } = renderHook(() => useRenderMode(NODE_THRESHOLD))

      expect(result.current.shouldShowWarning).toBe(true)
      expect(result.current.warningMessage).toContain('SVG mode may be slow')
      expect(result.current.warningMessage).toContain(`${NODE_THRESHOLD}`)
    })

    it('should not show warning for small graphs in forced SVG', () => {
      useSettingsStore.setState({ visualizationMode: 'force-svg' })
      const { result } = renderHook(() => useRenderMode(100))

      expect(result.current.shouldShowWarning).toBe(false)
      expect(result.current.warningMessage).toBeNull()
    })
  })

  describe('Force Canvas mode (AC3)', () => {
    it('should return Canvas mode regardless of node count', () => {
      useSettingsStore.setState({ visualizationMode: 'force-canvas' })
      const { result } = renderHook(() => useRenderMode(10))

      expect(result.current.mode).toBe('canvas')
      expect(result.current.isAutoMode).toBe(false)
      expect(result.current.isForced).toBe(true)
    })

    it('should not show warning for forced Canvas mode', () => {
      useSettingsStore.setState({ visualizationMode: 'force-canvas' })
      const { result } = renderHook(() => useRenderMode(10))

      expect(result.current.shouldShowWarning).toBe(false)
      expect(result.current.warningMessage).toBeNull()
    })
  })
})
