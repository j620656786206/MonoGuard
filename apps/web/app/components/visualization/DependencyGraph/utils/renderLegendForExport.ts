/**
 * Legend Renderer for Export
 *
 * Generates a standalone SVG legend string that can be embedded
 * in exported SVG or PNG images. Matches the GraphLegend component styling.
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images - AC5
 */

import { LEGEND_COLORS } from '../styles'

/**
 * Generates a standalone SVG legend string for embedding in exports.
 * Matches the visual style of the GraphLegend React component.
 *
 * @param isDarkMode - Whether to use dark mode colors
 * @returns SVG markup string
 */
export function renderLegendSvg(isDarkMode: boolean): string {
  const bgColor = isDarkMode ? '#374151' : '#f3f4f6'
  const textColor = isDarkMode ? '#f3f4f6' : '#1f2937'
  const borderColor = isDarkMode ? '#4b5563' : '#d1d5db'

  return `<svg xmlns="http://www.w3.org/2000/svg" width="150" height="100">
  <rect width="150" height="100" rx="8" fill="${bgColor}" stroke="${borderColor}" stroke-width="1"/>
  <text x="10" y="18" fill="${textColor}" font-size="11" font-weight="600" font-family="system-ui, sans-serif">Legend</text>
  <circle cx="20" cy="38" r="6" fill="${LEGEND_COLORS.normalNode}"/>
  <text x="32" y="42" fill="${textColor}" font-size="10" font-family="system-ui, sans-serif">Package</text>
  <circle cx="20" cy="58" r="6" fill="${LEGEND_COLORS.cycleNode}" stroke="#dc2626" stroke-width="2"/>
  <text x="32" y="62" fill="${textColor}" font-size="10" font-family="system-ui, sans-serif">Circular Dep</text>
  <line x1="12" y1="78" x2="28" y2="78" stroke="${LEGEND_COLORS.normalEdge}" stroke-width="1.5"/>
  <text x="32" y="82" fill="${textColor}" font-size="10" font-family="system-ui, sans-serif">Dependency</text>
</svg>`
}
