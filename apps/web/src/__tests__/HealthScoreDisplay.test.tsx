import type { HealthScore } from '@monoguard/types'
import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { HealthScoreDisplay } from '@/components/analysis/HealthScoreDisplay'

// Factory function for creating test data
const createMockHealthScore = (overrides: Partial<HealthScore> = {}): HealthScore => ({
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
  ...overrides,
})

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

  describe('Trend Indicators', () => {
    it('[P1] should display stable trend indicator', () => {
      // GIVEN: Health score with stable trend
      const stableHealth = createMockHealthScore({ trend: 'stable' })

      render(<HealthScoreDisplay healthScore={stableHealth} />)

      // THEN: Stable trend should be displayed
      expect(screen.getByText('stable')).toBeInTheDocument()
    })

    it('[P1] should display declining trend indicator', () => {
      // GIVEN: Health score with declining trend
      const decliningHealth = createMockHealthScore({ trend: 'declining' })

      render(<HealthScoreDisplay healthScore={decliningHealth} />)

      // THEN: Declining trend should be displayed
      expect(screen.getByText('declining')).toBeInTheDocument()
    })
  })

  describe('Score Color Boundaries', () => {
    it('[P2] should apply green-500 for score 80-89', () => {
      // GIVEN: Score of exactly 80
      const health = createMockHealthScore({ overall: 80 })

      render(<HealthScoreDisplay healthScore={health} />)

      // THEN: Should have green-500 color
      const scoreElements = screen.getAllByText('80')
      const mainScore = scoreElements.find((el) => el.classList.contains('text-6xl'))
      expect(mainScore).toHaveClass('text-green-500')
    })

    it('[P2] should apply yellow-500 for score 70-79', () => {
      // GIVEN: Score in yellow range
      const health = createMockHealthScore({ overall: 75 })

      render(<HealthScoreDisplay healthScore={health} />)

      // THEN: Should have yellow-500 color
      const scoreElements = screen.getAllByText('75')
      const mainScore = scoreElements.find((el) => el.classList.contains('text-6xl'))
      expect(mainScore).toHaveClass('text-yellow-500')
    })

    it('[P2] should apply orange-500 for score 60-69', () => {
      // GIVEN: Score in orange range
      const health = createMockHealthScore({ overall: 65 })

      render(<HealthScoreDisplay healthScore={health} />)

      // THEN: Should have orange-500 color
      const scoreElements = screen.getAllByText('65')
      const mainScore = scoreElements.find((el) => el.classList.contains('text-6xl'))
      expect(mainScore).toHaveClass('text-orange-500')
    })

    it('[P2] should apply green-600 for score >= 90', () => {
      // GIVEN: Score at exactly 90 (boundary)
      const health = createMockHealthScore({ overall: 90 })

      render(<HealthScoreDisplay healthScore={health} />)

      // THEN: Should have green-600 color
      const scoreElements = screen.getAllByText('90')
      const mainScore = scoreElements.find((el) => el.classList.contains('text-6xl'))
      expect(mainScore).toHaveClass('text-green-600')
    })
  })

  describe('Health Factors Interaction', () => {
    it('[P1] should expand factor to show recommendations', () => {
      // GIVEN: Health score with factor that has recommendations
      render(<HealthScoreDisplay healthScore={mockHealthScore} />)

      // WHEN: Click "Show Tips" button
      const showTipsButton = screen.getByText('Show Tips')
      fireEvent.click(showTipsButton)

      // THEN: Recommendations section should be visible
      expect(screen.getByText('Recommendations:')).toBeInTheDocument()
      // The recommendation text may appear multiple times (in factor and Key Recommendations)
      expect(
        screen.getAllByText('Consider updating lodash to latest version').length
      ).toBeGreaterThan(0)
    })

    it('[P1] should collapse factor when clicking Hide Tips', () => {
      // GIVEN: Health score with expanded factor
      render(<HealthScoreDisplay healthScore={mockHealthScore} />)

      // WHEN: Click "Show Tips" then "Hide Tips"
      fireEvent.click(screen.getByText('Show Tips'))
      fireEvent.click(screen.getByText('Hide Tips'))

      // THEN: Recommendations should be hidden
      expect(screen.queryByText('Recommendations:')).not.toBeInTheDocument()
    })

    it('[P2] should not show Tips button for factors without recommendations', () => {
      // GIVEN: Health score with factor that has no recommendations
      const health = createMockHealthScore({
        factors: [
          {
            name: 'Clean Code',
            score: 100,
            weight: 0.25,
            description: 'Perfect code quality',
            recommendations: [],
          },
        ],
      })

      render(<HealthScoreDisplay healthScore={health} />)

      // THEN: No Tips button should be visible
      expect(screen.queryByText('Show Tips')).not.toBeInTheDocument()
    })
  })

  describe('Empty States', () => {
    it('[P2] should handle empty factors array', () => {
      // GIVEN: Health score with no factors
      const health = createMockHealthScore({ factors: [] })

      render(<HealthScoreDisplay healthScore={health} />)

      // THEN: Health Factors section should still render
      expect(screen.getByText('Health Factors')).toBeInTheDocument()
    })
  })

  describe('Date Display', () => {
    it('[P2] should display formatted last updated date', () => {
      // GIVEN: Health score with specific date
      const health = createMockHealthScore({ lastUpdated: '2026-01-15T10:00:00Z' })

      render(<HealthScoreDisplay healthScore={health} />)

      // THEN: Date should be displayed (format depends on locale)
      expect(screen.getByText(/Last updated/)).toBeInTheDocument()
    })
  })
})
