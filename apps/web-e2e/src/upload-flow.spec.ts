/**
 * MonoGuard E2E Tests - Upload Flow
 *
 * Tests for the analyze page and file upload functionality.
 * Priority: P0 (Critical user path)
 *
 * Test Coverage:
 * - Analyze page UI display
 * - Upload instructions
 * - Navigation flow
 */

import { expect, test } from './support/fixtures'

test.describe('Analyze Flow', () => {
  test.describe('Analyze Page UI', () => {
    test('[P0] should display analyze page with upload area', async ({ page }) => {
      // GIVEN: User navigates to analyze page
      await page.goto('/analyze')

      // THEN: Page should be visible with proper heading
      await expect(page.locator('h1')).toContainText('Analyze')
    })

    test('[P1] should display upload instructions', async ({ page }) => {
      // GIVEN: User is on analyze page
      await page.goto('/analyze')

      // THEN: Instructions about supported formats should be visible
      await expect(page.getByText(/workspace\.json|drop/i)).toBeVisible()
    })

    test('[P1] should show supported workspace types', async ({ page }) => {
      // GIVEN: User is on analyze page
      await page.goto('/analyze')

      // THEN: Supported workspace types should be mentioned
      await expect(page.getByText(/Nx|Lerna|Turborepo/i)).toBeVisible()
    })

    test('[P1] should have Select File button', async ({ page }) => {
      // GIVEN: User is on analyze page
      await page.goto('/analyze')

      // THEN: Select File button should be visible
      await expect(page.getByText(/Select File/i)).toBeVisible()
    })
  })

  test.describe('Navigation', () => {
    test('[P0] should navigate from home to analyze page', async ({ page }) => {
      // GIVEN: User is on home page
      await page.goto('/')

      // WHEN: User clicks Start Analysis
      await page.getByText(/Start Analysis/i).click()

      // THEN: Should navigate to analyze page
      await expect(page).toHaveURL(/analyze/)
      await expect(page.locator('h1')).toContainText('Analyze')
    })

    test('[P1] should navigate from results to analyze page', async ({ page }) => {
      // GIVEN: User is on results page
      await page.goto('/results')

      // WHEN: User clicks Start New Analysis
      await page.getByText(/Start New Analysis/i).click()

      // THEN: Should navigate to analyze page
      await expect(page).toHaveURL(/analyze/)
    })
  })

  test.describe('Results Page', () => {
    test('[P0] should display results page', async ({ page }) => {
      // GIVEN: User navigates to results page
      await page.goto('/results')

      // THEN: Results page should be visible
      await expect(page.locator('h1')).toContainText('Analysis Results')
    })

    test('[P1] should display health score section', async ({ page }) => {
      // GIVEN: User is on results page
      await page.goto('/results')

      // THEN: Health score section should be visible
      await expect(page.getByText(/Health Score/i)).toBeVisible()
    })

    test('[P1] should display circular dependencies section', async ({ page }) => {
      // GIVEN: User is on results page
      await page.goto('/results')

      // THEN: Circular dependencies section should be visible
      await expect(page.getByText(/Circular Dependencies/i)).toBeVisible()
    })

    test('[P1] should display total packages section', async ({ page }) => {
      // GIVEN: User is on results page
      await page.goto('/results')

      // THEN: Total packages section should be visible
      await expect(page.getByText(/Total Packages/i)).toBeVisible()
    })

    test('[P1] should display dependency graph placeholder', async ({ page }) => {
      // GIVEN: User is on results page
      await page.goto('/results')

      // THEN: Dependency graph section should be visible
      await expect(page.getByText(/Dependency Graph/i)).toBeVisible()
    })
  })

  test.describe('Responsive Design', () => {
    test('[P2] should be responsive on tablet viewport', async ({ page }) => {
      // GIVEN: Tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 })

      // WHEN: User views analyze page
      await page.goto('/analyze')

      // THEN: Page should render correctly
      await expect(page.locator('h1')).toContainText('Analyze')
    })

    test('[P2] should be responsive on mobile viewport', async ({ page }) => {
      // GIVEN: Mobile viewport
      await page.setViewportSize({ width: 375, height: 667 })

      // WHEN: User views analyze page
      await page.goto('/analyze')

      // THEN: Page should render correctly
      await expect(page.locator('h1')).toContainText('Analyze')
    })
  })
})
