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
    await expect(page.locator('h1')).toBeVisible()
  })

  test('should have navigation to upload page', async ({ page }) => {
    await goToHome(page)

    // Look for a link or button to the upload functionality
    const uploadLink = page.locator('a[href*="upload"], [data-testid="start-analysis"]')
    await expect(uploadLink).toBeVisible()
  })

  test('should be responsive on mobile viewport', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 })
    await goToHome(page)

    // Page should still render correctly
    await expect(page.locator('body')).toBeVisible()
  })
})
