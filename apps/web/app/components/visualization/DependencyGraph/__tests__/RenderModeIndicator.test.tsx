/**
 * Tests for RenderModeIndicator component
 *
 * @see Story 4.9: Implement Hybrid SVG/Canvas Rendering
 * @see AC2: Mode Indicator Display
 */

import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { RenderModeIndicator } from '../RenderModeIndicator'

describe('RenderModeIndicator', () => {
  it('should display node count', () => {
    render(<RenderModeIndicator mode="svg" nodeCount={42} isForced={false} />)

    expect(screen.getByText('42 nodes')).toBeTruthy()
  })

  it('should display SVG mode indicator', () => {
    render(<RenderModeIndicator mode="svg" nodeCount={100} isForced={false} />)

    expect(screen.getByText('SVG mode')).toBeTruthy()
  })

  it('should display Canvas mode indicator', () => {
    render(<RenderModeIndicator mode="canvas" nodeCount={600} isForced={false} />)

    expect(screen.getByText('CANVAS mode')).toBeTruthy()
  })

  it('should show "Forced" badge when mode is forced by user', () => {
    render(<RenderModeIndicator mode="svg" nodeCount={600} isForced={true} />)

    expect(screen.getByText('Forced')).toBeTruthy()
  })

  it('should not show "Forced" badge in auto mode', () => {
    render(<RenderModeIndicator mode="svg" nodeCount={100} isForced={false} />)

    expect(screen.queryByText('Forced')).toBeNull()
  })

  it('should have proper ARIA label', () => {
    render(<RenderModeIndicator mode="canvas" nodeCount={500} isForced={true} />)

    const indicator = screen.getByLabelText('500 nodes, CANVAS rendering mode, forced')
    expect(indicator).toBeTruthy()
  })

  it('should have ARIA label without forced suffix in auto mode', () => {
    render(<RenderModeIndicator mode="svg" nodeCount={100} isForced={false} />)

    const indicator = screen.getByLabelText('100 nodes, SVG rendering mode')
    expect(indicator).toBeTruthy()
  })

  it('should display performance warning when warningMessage is provided (AC3)', () => {
    render(
      <RenderModeIndicator
        mode="svg"
        nodeCount={600}
        isForced={true}
        warningMessage="SVG mode may be slow with 600 nodes"
      />
    )

    expect(screen.getByRole('alert')).toBeTruthy()
    expect(screen.getByText('SVG mode may be slow with 600 nodes')).toBeTruthy()
  })

  it('should not display warning when warningMessage is null', () => {
    render(
      <RenderModeIndicator mode="svg" nodeCount={100} isForced={false} warningMessage={null} />
    )

    expect(screen.queryByRole('alert')).toBeNull()
  })
})
