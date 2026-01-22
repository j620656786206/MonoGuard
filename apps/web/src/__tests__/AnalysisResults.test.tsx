import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import { AnalysisResults } from '@/components/analysis/AnalysisResults'
import {
  createAnalysisSummary,
  createArchitectureValidationResults,
  createBundleImpactReport,
  createCircularDependency,
  createCompletedAnalysis,
  createDuplicateDetectionResults,
  createHealthScore,
  createInProgressAnalysis,
  createVersionConflict,
} from './factories/analysis.factory'

describe('AnalysisResults', () => {
  describe('In Progress State', () => {
    it('[P1] should display progress bar when analysis is in progress', () => {
      // GIVEN: Analysis in progress
      const analysis = createInProgressAnalysis({ progress: 45 })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Progress indicator should be visible
      expect(screen.getByText('Analysis in Progress')).toBeInTheDocument()
      expect(screen.getByText('45% complete')).toBeInTheDocument()
    })

    it('[P1] should display current step message', () => {
      // GIVEN: Analysis with current step
      const analysis = createInProgressAnalysis({
        currentStep: 'Analyzing circular dependencies...',
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Current step should be visible
      expect(screen.getByText('Analyzing circular dependencies...')).toBeInTheDocument()
    })

    it('[P1] should display default message when no current step', () => {
      // GIVEN: Analysis without current step
      const analysis = createInProgressAnalysis({ currentStep: undefined })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Default message should be visible
      expect(screen.getByText('Starting analysis...')).toBeInTheDocument()
    })

    it('[P1] should show progress bar with correct width', () => {
      // GIVEN: Analysis at 75% progress
      const analysis = createInProgressAnalysis({ progress: 75 })

      const { container } = render(<AnalysisResults analysis={analysis} />)

      // THEN: Progress bar should have correct width
      const progressBar = container.querySelector('[style*="width: 75%"]')
      expect(progressBar).toBeInTheDocument()
    })
  })

  describe('Completed State - Header', () => {
    it('[P1] should display Analysis Results header when complete', () => {
      // GIVEN: Completed analysis
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Header should be visible
      expect(screen.getByText('Analysis Results')).toBeInTheDocument()
    })

    it('[P1] should display completion time', () => {
      // GIVEN: Completed analysis with specific time
      const completedAt = '2024-01-15T10:30:00Z'
      const analysis = createCompletedAnalysis({ completedAt })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Completion info should be visible
      expect(screen.getByText(/Completed/)).toBeInTheDocument()
    })

    it('[P1] should display New Analysis button when callback provided', () => {
      // GIVEN: Completed analysis with callback
      const onNewAnalysis = vi.fn()
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} onNewAnalysis={onNewAnalysis} />)

      // THEN: Button should be visible
      expect(screen.getByText('New Analysis')).toBeInTheDocument()
    })

    it('[P1] should call onNewAnalysis when button clicked', () => {
      // GIVEN: Completed analysis with callback
      const onNewAnalysis = vi.fn()
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} onNewAnalysis={onNewAnalysis} />)

      // WHEN: Button is clicked
      fireEvent.click(screen.getByText('New Analysis'))

      // THEN: Callback should be called
      expect(onNewAnalysis).toHaveBeenCalled()
    })

    it('[P1] should not show New Analysis button when no callback', () => {
      // GIVEN: Completed analysis without callback
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Button should not be visible
      expect(screen.queryByText('New Analysis')).not.toBeInTheDocument()
    })
  })

  describe('Tab Navigation', () => {
    it('[P1] should display all tab buttons', () => {
      // GIVEN: Completed analysis
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} />)

      // THEN: All tabs should be visible (check within nav element)
      const nav = screen.getByRole('navigation', { name: 'Tabs' })
      expect(nav).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Overview/ })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Dependencies/ })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Architecture/ })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Duplicates/ })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Bundle Impact/ })).toBeInTheDocument()
      expect(screen.getByRole('button', { name: /Health Score/ })).toBeInTheDocument()
    })

    it('[P1] should display count badges on tabs with issues', () => {
      // GIVEN: Analysis with issues
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary(),
          dependencyAnalysis: {
            duplicateDependencies: [],
            versionConflicts: [createVersionConflict()],
            unusedDependencies: [],
            circularDependencies: [createCircularDependency(), createCircularDependency()],
            bundleImpact: createBundleImpactReport(),
            summary: createAnalysisSummary(),
          },
          architectureValidation: createArchitectureValidationResults(),
          duplicateDetection: createDuplicateDetectionResults({ totalDuplicates: 5 }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Count badges should show correct numbers
      // Dependencies tab should show circular deps count (2)
      const depsTab = screen.getByText('Dependencies').closest('button')
      expect(depsTab?.textContent).toContain('2')

      // Architecture tab should show violations count (1)
      const archTab = screen.getByText('Architecture').closest('button')
      expect(archTab?.textContent).toContain('1')

      // Duplicates tab should show duplicates count (5)
      const dupsTab = screen.getByText('Duplicates').closest('button')
      expect(dupsTab?.textContent).toContain('5')
    })

    it('[P1] should show Overview tab content by default', () => {
      // GIVEN: Completed analysis
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Overview stats should be visible (use getAllByText for repeated labels)
      // Health Score appears both as a tab and stat label
      const healthScoreLabels = screen.getAllByText('Health Score')
      expect(healthScoreLabels.length).toBeGreaterThanOrEqual(1)
      expect(screen.getByText('Circular Dependencies')).toBeInTheDocument()
      expect(screen.getByText('Version Conflicts')).toBeInTheDocument()
    })

    it('[P1] should switch to Dependencies tab on click', () => {
      // GIVEN: Completed analysis
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Dependencies tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Dependencies/ }))

      // THEN: Dependencies content should be visible (CircularDependencyViz component)
      expect(screen.getByText('Dependency Analysis')).toBeInTheDocument()
    })

    it('[P1] should switch to Architecture tab on click', () => {
      // GIVEN: Completed analysis
      const analysis = createCompletedAnalysis()

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Architecture tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Architecture/ }))

      // THEN: Architecture panel should be visible
      // The ArchitectureValidationPanel shows violations
      const archElements = screen.getAllByText(/Architecture/)
      expect(archElements.length).toBeGreaterThanOrEqual(1)
    })
  })

  describe('Overview Panel', () => {
    it('[P1] should display health score with color coding', () => {
      // GIVEN: Analysis with good health score (use unique value 89 to avoid collisions)
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ healthScore: 89 }),
          healthScore: createHealthScore({ overall: 89 }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Health score should be displayed with suffix as combined text "89/100"
      expect(screen.getByText('89/100')).toBeInTheDocument()
    })

    it('[P1] should display circular dependencies count', () => {
      // GIVEN: Analysis with circular deps (use unique number to avoid ambiguity)
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ circularCount: 7 }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Count should be displayed
      expect(screen.getByText('7')).toBeInTheDocument()
    })

    it('[P1] should display version conflicts count', () => {
      // GIVEN: Analysis with version conflicts (use unique number)
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ conflictCount: 5 }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Count should be displayed
      expect(screen.getByText('5')).toBeInTheDocument()
    })

    it('[P1] should display duplicate packages count', () => {
      // GIVEN: Analysis with duplicates (use unique number)
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ duplicateCount: 8 }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Count should be displayed
      expect(screen.getByText('8')).toBeInTheDocument()
    })

    it('[P1] should display architecture violations count', () => {
      // GIVEN: Analysis with 4 violations (unique count to avoid collisions)
      const analysis = createCompletedAnalysis({
        results: {
          architectureValidation: createArchitectureValidationResults({
            violations: [
              {
                ruleName: 'test1',
                severity: 'high',
                description: 'test1',
                violatingFile: 'test1.ts',
                violatingImport: 'test1',
                expectedLayer: 'domain',
                actualLayer: 'ui',
                suggestion: 'fix it',
              },
              {
                ruleName: 'test2',
                severity: 'medium',
                description: 'test2',
                violatingFile: 'test2.ts',
                violatingImport: 'test2',
                expectedLayer: 'domain',
                actualLayer: 'ui',
                suggestion: 'fix it',
              },
              {
                ruleName: 'test3',
                severity: 'low',
                description: 'test3',
                violatingFile: 'test3.ts',
                violatingImport: 'test3',
                expectedLayer: 'domain',
                actualLayer: 'ui',
                suggestion: 'fix it',
              },
              {
                ruleName: 'test4',
                severity: 'low',
                description: 'test4',
                violatingFile: 'test4.ts',
                violatingImport: 'test4',
                expectedLayer: 'domain',
                actualLayer: 'ui',
                suggestion: 'fix it',
              },
            ],
          }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Violations count should be displayed (4)
      // May appear in multiple places (stat card + tab badge), use getAllByText
      const fourElements = screen.getAllByText('4')
      expect(fourElements.length).toBeGreaterThanOrEqual(1)
    })

    it('[P1] should display potential savings', () => {
      // GIVEN: Analysis with bundle impact
      const analysis = createCompletedAnalysis({
        results: {
          bundleImpact: createBundleImpactReport({ potentialSavings: '500 KB' }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Potential savings should be displayed
      expect(screen.getByText('500 KB')).toBeInTheDocument()
    })

    it('[P1] should display warnings when present', () => {
      // GIVEN: Analysis with warnings
      const analysis = createCompletedAnalysis({
        warnings: ['Some packages could not be analyzed', 'Lock file is outdated'],
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Warnings should be displayed
      expect(screen.getByText('Warnings')).toBeInTheDocument()
      expect(screen.getByText('Some packages could not be analyzed')).toBeInTheDocument()
      expect(screen.getByText('Lock file is outdated')).toBeInTheDocument()
    })
  })

  describe('Empty States', () => {
    it('[P1] should display empty state for missing dependency analysis', () => {
      // GIVEN: Analysis without dependency results
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary(),
          dependencyAnalysis: undefined,
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Dependencies tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Dependencies/ }))

      // THEN: Empty state should be shown
      expect(screen.getByText('No dependency analysis results available')).toBeInTheDocument()
    })

    it('[P1] should display empty state for missing architecture validation', () => {
      // GIVEN: Analysis without architecture results
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary(),
          architectureValidation: undefined,
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Architecture tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Architecture/ }))

      // THEN: Empty state should be shown
      expect(screen.getByText('No architecture validation results available')).toBeInTheDocument()
    })

    it('[P1] should display empty state for missing duplicate detection', () => {
      // GIVEN: Analysis without duplicate results
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary(),
          duplicateDetection: undefined,
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Duplicates tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Duplicates/ }))

      // THEN: Empty state should be shown
      expect(screen.getByText('No duplicate detection results available')).toBeInTheDocument()
    })

    it('[P1] should display empty state for missing bundle impact', () => {
      // GIVEN: Analysis without bundle results
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary(),
          bundleImpact: undefined,
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Bundle Impact tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Bundle Impact/ }))

      // THEN: Empty state should be shown
      expect(screen.getByText('No bundle impact analysis available')).toBeInTheDocument()
    })
  })

  describe('Health Score Tab', () => {
    it('[P1] should display HealthScoreDisplay component', () => {
      // GIVEN: Analysis with health score
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary(),
          healthScore: createHealthScore({ overall: 75 }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Health Score tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Health Score/ }))

      // THEN: HealthScoreDisplay should be rendered
      // It should show the overall score
      const scoreElements = screen.getAllByText('75')
      expect(scoreElements.length).toBeGreaterThanOrEqual(1)
    })

    it('[P1] should handle numeric health score', () => {
      // GIVEN: Analysis with numeric health score
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ healthScore: 80 }),
          healthScore: 80 as unknown as ReturnType<typeof createHealthScore>, // Numeric value instead of object
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // WHEN: Health Score tab is clicked
      fireEvent.click(screen.getByRole('button', { name: /Health Score/ }))

      // THEN: Should handle gracefully and show score
      const scoreElements = screen.getAllByText('80')
      expect(scoreElements.length).toBeGreaterThanOrEqual(1)
    })
  })

  describe('Color Coding', () => {
    it('[P2] should apply green color for good health score (>= 80)', () => {
      // GIVEN: Analysis with good health score - override both summary and healthScore
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ healthScore: 92 }),
          healthScore: createHealthScore({ overall: 92 }),
        },
      })

      const { container } = render(<AnalysisResults analysis={analysis} />)

      // THEN: Green styling should be applied (bg-green-50)
      const greenElements = container.querySelectorAll('[class*="bg-green"]')
      expect(greenElements.length).toBeGreaterThan(0)
    })

    it('[P2] should apply yellow color for medium health score (60-79)', () => {
      // GIVEN: Analysis with medium health score
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ healthScore: 65 }),
        },
      })

      const { container } = render(<AnalysisResults analysis={analysis} />)

      // THEN: Yellow styling should be applied
      const yellowElements = container.querySelectorAll('[class*="bg-yellow"]')
      expect(yellowElements.length).toBeGreaterThan(0)
    })

    it('[P2] should apply red color for poor health score (< 60)', () => {
      // GIVEN: Analysis with poor health score
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({ healthScore: 45 }),
        },
      })

      const { container } = render(<AnalysisResults analysis={analysis} />)

      // THEN: Red styling should be applied
      const redElements = container.querySelectorAll('[class*="bg-red"]')
      expect(redElements.length).toBeGreaterThan(0)
    })
  })

  describe('Edge Cases', () => {
    it('[P2] should handle analysis with no results object', () => {
      // GIVEN: Analysis without results
      const analysis = createInProgressAnalysis({ progress: 0 })

      // WHEN: Rendered
      render(<AnalysisResults analysis={analysis} />)

      // THEN: Should show progress view without errors
      expect(screen.getByText('Analysis in Progress')).toBeInTheDocument()
    })

    it('[P2] should handle zero values correctly', () => {
      // GIVEN: Analysis with all zeros - override healthScore object to match summary
      const analysis = createCompletedAnalysis({
        results: {
          summary: createAnalysisSummary({
            healthScore: 100,
            circularCount: 0,
            conflictCount: 0,
            duplicateCount: 0,
          }),
          healthScore: createHealthScore({ overall: 100 }),
          architectureValidation: createArchitectureValidationResults({
            violations: [],
          }),
          bundleImpact: createBundleImpactReport({ potentialSavings: '0 KB' }),
        },
      })

      render(<AnalysisResults analysis={analysis} />)

      // THEN: Zero values should be displayed correctly
      // Health score shows as "100/100"
      expect(screen.getByText('100/100')).toBeInTheDocument()
      expect(screen.getByText('0 KB')).toBeInTheDocument()
    })
  })
})
