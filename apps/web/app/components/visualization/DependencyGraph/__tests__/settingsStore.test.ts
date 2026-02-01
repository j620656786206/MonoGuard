/**
 * Tests for settings store
 *
 * @see Story 4.9: Implement Hybrid SVG/Canvas Rendering
 * @see AC3: User Override in Settings
 */

import { afterEach, describe, expect, it } from 'vitest'
import { useSettingsStore } from '../../../../stores/settings'

// Clear store state before each test
const originalState = useSettingsStore.getState()

describe('useSettingsStore', () => {
  afterEach(() => {
    useSettingsStore.setState(originalState)
  })

  it('should have default visualizationMode of "auto"', () => {
    const state = useSettingsStore.getState()
    expect(state.visualizationMode).toBe('auto')
  })

  it('should update visualizationMode to "force-svg"', () => {
    useSettingsStore.getState().setVisualizationMode('force-svg')
    expect(useSettingsStore.getState().visualizationMode).toBe('force-svg')
  })

  it('should update visualizationMode to "force-canvas"', () => {
    useSettingsStore.getState().setVisualizationMode('force-canvas')
    expect(useSettingsStore.getState().visualizationMode).toBe('force-canvas')
  })

  it('should update visualizationMode back to "auto"', () => {
    useSettingsStore.getState().setVisualizationMode('force-canvas')
    useSettingsStore.getState().setVisualizationMode('auto')
    expect(useSettingsStore.getState().visualizationMode).toBe('auto')
  })
})
