/**
 * MonoGuard E2E Tests - Results Dashboard
 *
 * Tests for the results/dashboard functionality.
 * Priority: P1 (Important user feature)
 *
 * Test Coverage:
 * - Results page loads
 * - Health score display
 * - Circular dependencies display
 * - Total packages display
 * - Navigation to analyze
 */

import { expect, test } from './support/fixtures'

test.describe('Results Dashboard', () => {
  test.describe('Page Load', () => {
    test('[P1] should load results page successfully', async ({ page }) => {
      // GIVEN: User navigates to results
      await page.goto('/results')

      // THEN: Results page should be visible
      await expect(page.locator('h1')).toContainText('Analysis Results')
    })

    test('[P1] should display main content area', async ({ page }) => {
      // WHEN: User navigates to results
      await page.goto('/results')

      // THEN: Main content should be visible
      await expect(page.locator('main')).toBeVisible()
    })
  })

  test.describe('Health Score', () => {
    test('[P1] should display health score card', async ({ page }) => {
      // WHEN: User views results
      await page.goto('/results')

      // THEN: Health score should be displayed
      await expect(page.getByText(/Health Score/i)).toBeVisible()
    })

    test('[P1] should show placeholder when no data', async ({ page }) => {
      // WHEN: User views results without analysis
      await page.goto('/results')

      // THEN: Placeholder message should be shown
      await expect(page.getByText(/No analysis data/i)).toBeVisible()
    })
  })

  test.describe('Stats Cards', () => {
    test('[P1] should display circular dependencies card', async ({ page }) => {
      // WHEN: User views results
      await page.goto('/results')

      // THEN: Circular dependencies should be shown
      await expect(page.getByText(/Circular Dependencies/i)).toBeVisible()
    })

    test('[P1] should display total packages card', async ({ page }) => {
      // WHEN: User views results
      await page.goto('/results')

      // THEN: Total packages should be shown
      await expect(page.getByText(/Total Packages/i)).toBeVisible()
    })
  })

  test.describe('Dependency Graph', () => {
    test('[P1] should display dependency graph section', async ({ page }) => {
      // WHEN: User views results
      await page.goto('/results')

      // THEN: Dependency graph section should be visible
      await expect(page.getByText(/Dependency Graph/i)).toBeVisible()
    })

    test('[P1] should show placeholder for empty graph', async ({ page }) => {
      // WHEN: User views results without analysis
      await page.goto('/results')

      // THEN: Placeholder should be shown
      await expect(page.getByText(/Run analysis to visualize/i)).toBeVisible()
    })
  })

  test.describe('Quick Actions', () => {
    test('[P1] should have Start New Analysis button', async ({ page }) => {
      // WHEN: User views results
      await page.goto('/results')

      // THEN: Start New Analysis button should be visible
      await expect(page.getByText(/Start New Analysis/i)).toBeVisible()
    })

    test('[P1] should navigate to analyze page', async ({ page }) => {
      // GIVEN: User is on results page
      await page.goto('/results')

      // WHEN: User clicks Start New Analysis
      await page.getByText(/Start New Analysis/i).click()

      // THEN: Should navigate to analyze page
      await expect(page).toHaveURL(/analyze/)
    })
  })

  test.describe('Responsive Design', () => {
    test('[P2] should be responsive on tablet viewport', async ({ page }) => {
      // GIVEN: Tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 })

      // WHEN: User views results
      await page.goto('/results')

      // THEN: Results should render correctly
      await expect(page.locator('h1')).toContainText('Analysis Results')
    })

    test('[P2] should be responsive on mobile viewport', async ({ page }) => {
      // GIVEN: Mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })

      // WHEN: User views results
      await page.goto('/results')

      // THEN: Results should render correctly
      await expect(page.locator('h1')).toContainText('Analysis Results')
    })

    test('[P2] should stack cards on mobile', async ({ page }) => {
      // GIVEN: Mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })

      // WHEN: User views results
      await page.goto('/results')

      // THEN: Cards should be visible (layout adapts)
      await expect(page.getByText(/Health Score/i)).toBeVisible()
      await expect(page.getByText(/Circular Dependencies/i)).toBeVisible()
    })
  })
})
