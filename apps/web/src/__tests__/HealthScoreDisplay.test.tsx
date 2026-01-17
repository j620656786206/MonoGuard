import type { HealthScore } from '@monoguard/types'
import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { HealthScoreDisplay } from '@/components/analysis/HealthScoreDisplay'

const mockHealthScore: HealthScore = {
  overall: 85,
  dependencies: 90,
  architecture: 80,
  maintainability: 85,
  security: 88,
  performance: 82,
  lastUpdated: '2026-01-15T10:00:00Z',
  trend: 'improving',
  factors: [
    {
      name: 'Dependency Health',
      score: 90,
      weight: 0.25,
      description: 'Dependencies are well-maintained',
      recommendations: ['Consider updating lodash to latest version'],
    },
    {
      name: 'Code Structure',
      score: 80,
      weight: 0.25,
      description: 'Good code organization',
      recommendations: [],
    },
  ],
}

describe('HealthScoreDisplay', () => {
  it('renders the overall health score', () => {
    render(<HealthScoreDisplay healthScore={mockHealthScore} />)

    // Check that the overall score is displayed (may appear multiple times in different contexts)
    const scoreElements = screen.getAllByText('85')
    expect(scoreElements.length).toBeGreaterThan(0)
    expect(screen.getByText('/ 100')).toBeInTheDocument()
  })

  it('displays the overall health score title', () => {
    render(<HealthScoreDisplay healthScore={mockHealthScore} />)

    expect(screen.getByText('Overall Health Score')).toBeInTheDocument()
  })

  it('shows the trend indicator', () => {
    render(<HealthScoreDisplay healthScore={mockHealthScore} />)

    expect(screen.getByText('improving')).toBeInTheDocument()
  })

  it('renders category score cards', () => {
    render(<HealthScoreDisplay healthScore={mockHealthScore} />)

    expect(screen.getByText('Dependencies')).toBeInTheDocument()
    expect(screen.getByText('Architecture')).toBeInTheDocument()
    expect(screen.getByText('Maintainability')).toBeInTheDocument()
    expect(screen.getByText('Security')).toBeInTheDocument()
    expect(screen.getByText('Performance')).toBeInTheDocument()
  })

  it('renders health factors section', () => {
    render(<HealthScoreDisplay healthScore={mockHealthScore} />)

    expect(screen.getByText('Health Factors')).toBeInTheDocument()
    // Factor names may appear multiple times (in factors list and recommendations)
    expect(screen.getAllByText('Dependency Health').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Code Structure').length).toBeGreaterThan(0)
  })

  it('renders recommendations section', () => {
    render(<HealthScoreDisplay healthScore={mockHealthScore} />)

    expect(screen.getByText('Key Recommendations')).toBeInTheDocument()
  })

  it('applies correct color class for high score', () => {
    const highScoreHealth: HealthScore = {
      ...mockHealthScore,
      overall: 95,
    }
    render(<HealthScoreDisplay healthScore={highScoreHealth} />)

    // Get the main score element (the large one with text-6xl)
    const scoreElements = screen.getAllByText('95')
    const mainScore = scoreElements.find((el) => el.classList.contains('text-6xl'))
    expect(mainScore).toHaveClass('text-green-600')
  })

  it('applies correct color class for low score', () => {
    const lowScoreHealth: HealthScore = {
      ...mockHealthScore,
      overall: 45,
    }
    render(<HealthScoreDisplay healthScore={lowScoreHealth} />)

    // Get the main score element (the large one with text-6xl)
    const scoreElements = screen.getAllByText('45')
    const mainScore = scoreElements.find((el) => el.classList.contains('text-6xl'))
    expect(mainScore).toHaveClass('text-red-500')
  })
})
