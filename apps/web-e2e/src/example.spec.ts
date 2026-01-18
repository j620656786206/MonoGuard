/**
 * MonoGuard E2E Tests - Landing Page
 *
 * Tests for the main landing page functionality.
 * Uses the merged fixture system from support/fixtures.
 */

import { expect, test } from './support/fixtures'
import { goToHome } from './support/helpers'

test.describe('Landing Page', () => {
  test('should display welcome message', async ({ page }) => {
    await goToHome(page)

    // Verify the page title or main heading
    await expect(page.locator('h1')).toContainText('MonoGuard')
  })

  test('should have navigation to analyze page', async ({ page }) => {
    await goToHome(page)

    // Look for a link to the analyze functionality
    const analyzeLink = page.locator('a[href*="analyze"]')
    await expect(analyzeLink).toBeVisible()
  })

  test('should have Start Analysis button', async ({ page }) => {
    await goToHome(page)

    // Check for Start Analysis button
    await expect(page.getByText(/Start Analysis/i)).toBeVisible()
  })

  test('should have View Results link', async ({ page }) => {
    await goToHome(page)

    // Check for View Results link
    await expect(page.getByText(/View Results/i)).toBeVisible()
  })

  test('should navigate to analyze page when clicking Start Analysis', async ({ page }) => {
    await goToHome(page)

    // Click Start Analysis
    await page.getByText(/Start Analysis/i).click()

    // Should navigate to analyze page
    await expect(page).toHaveURL(/analyze/)
  })

  test('should be responsive on mobile viewport', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 })
    await goToHome(page)

    // Page should still render correctly
    await expect(page.locator('h1')).toContainText('MonoGuard')
  })
})
