/**
 * Custom Assertion Helpers
 *
 * Reusable assertion patterns for MonoGuard E2E tests.
 * These supplement Playwright's built-in expect() assertions.
 */

import type { Locator, Page } from '@playwright/test'
import { expect } from '@playwright/test'

/**
 * Assert that the health score is displayed with a specific value range
 */
export async function assertHealthScore(
  page: Page,
  minScore: number,
  maxScore: number = 100
): Promise<void> {
  const healthScoreElement = page.locator('[data-testid="health-score-value"]')
  await expect(healthScoreElement).toBeVisible()

  const scoreText = await healthScoreElement.textContent()
  const score = parseInt(scoreText || '0', 10)

  expect(score).toBeGreaterThanOrEqual(minScore)
  expect(score).toBeLessThanOrEqual(maxScore)
}

/**
 * Assert that a specific number of issues are displayed
 */
export async function assertIssueCount(page: Page, expectedCount: number): Promise<void> {
  const issueCountElement = page.locator('[data-testid="issue-count"]')
  await expect(issueCountElement).toBeVisible()
  await expect(issueCountElement).toHaveText(expectedCount.toString())
}

/**
 * Assert that a toast notification appears with specific text
 */
export async function assertToastMessage(page: Page, messageContains: string): Promise<void> {
  const toast = page.locator('[data-testid="toast-message"]')
  await expect(toast).toBeVisible()
  await expect(toast).toContainText(messageContains)
}

/**
 * Assert that an element is not present in the DOM
 */
export async function assertNotPresent(page: Page, selector: string): Promise<void> {
  const element = page.locator(selector)
  await expect(element).toHaveCount(0)
}

/**
 * Assert that loading state is complete
 */
export async function assertLoadingComplete(page: Page): Promise<void> {
  // Wait for any loading spinners to disappear
  const loadingIndicator = page.locator('[data-testid="loading-indicator"]')
  await expect(loadingIndicator).toBeHidden({ timeout: 30000 })
}

/**
 * Assert that the page has no console errors
 */
export async function assertNoConsoleErrors(page: Page, errors: string[]): Promise<void> {
  const criticalErrors = errors.filter(
    (error) => !error.includes('[React DevTools]') && !error.includes('favicon.ico')
  )

  expect(criticalErrors).toHaveLength(0)
}
