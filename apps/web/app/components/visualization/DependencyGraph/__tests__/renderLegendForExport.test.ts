/**
 * Tests for renderLegendForExport utility
 *
 * @see Story 4.6: Export Graph as PNG/SVG Images - AC5
 */
import { describe, expect, it } from 'vitest'

import { LEGEND_COLORS } from '../styles'
import { renderLegendSvg } from '../utils/renderLegendForExport'

describe('renderLegendSvg', () => {
  it('should return a valid SVG string', () => {
    const svg = renderLegendSvg(false)

    expect(svg).toContain('<svg')
    expect(svg).toContain('</svg>')
    expect(svg).toContain('xmlns="http://www.w3.org/2000/svg"')
  })

  it('should include "Legend" title text', () => {
    const svg = renderLegendSvg(false)

    expect(svg).toContain('>Legend</text>')
  })

  it('should include node color entries matching LEGEND_COLORS', () => {
    const svg = renderLegendSvg(false)

    expect(svg).toContain(`fill="${LEGEND_COLORS.normalNode}"`)
    expect(svg).toContain(`fill="${LEGEND_COLORS.cycleNode}"`)
  })

  it('should include dependency edge entry', () => {
    const svg = renderLegendSvg(false)

    expect(svg).toContain(`stroke="${LEGEND_COLORS.normalEdge}"`)
    expect(svg).toContain('>Dependency</text>')
  })

  it('should use light mode colors when isDarkMode is false', () => {
    const svg = renderLegendSvg(false)

    // Light mode background
    expect(svg).toContain('fill="#f3f4f6"')
    // Light mode text
    expect(svg).toContain('fill="#1f2937"')
  })

  it('should use dark mode colors when isDarkMode is true', () => {
    const svg = renderLegendSvg(true)

    // Dark mode background
    expect(svg).toContain('fill="#374151"')
    // Dark mode text
    expect(svg).toContain('fill="#f3f4f6"')
  })

  it('should include Package and Circular Dep labels', () => {
    const svg = renderLegendSvg(false)

    expect(svg).toContain('>Package</text>')
    expect(svg).toContain('>Circular Dep</text>')
  })

  it('should have consistent dimensions', () => {
    const svg = renderLegendSvg(false)

    expect(svg).toContain('width="150"')
    expect(svg).toContain('height="100"')
  })
})
