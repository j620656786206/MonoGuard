/**
 * Settings store - Global application settings with persistence (Story 4.9)
 *
 * Uses Zustand with devtools and persist middleware per project-context.md.
 * Stores visualization preferences that persist across sessions.
 */

import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'
import type { RenderModePreference } from '../components/visualization/DependencyGraph/types'

export interface SettingsState {
  /** Visualization render mode preference: auto, force-svg, force-canvas */
  visualizationMode: RenderModePreference

  /** Actions */
  setVisualizationMode: (mode: RenderModePreference) => void
}

export const useSettingsStore = create<SettingsState>()(
  devtools(
    persist(
      (set) => ({
        visualizationMode: 'auto',

        setVisualizationMode: (mode) => set({ visualizationMode: mode }),
      }),
      {
        name: 'monoguard-settings',
      }
    ),
    {
      name: 'SettingsStore',
    }
  )
)
