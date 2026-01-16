/**
 * Navigation Helpers
 *
 * Pure functions for common navigation patterns.
 * Framework-agnostic design following TEA fixture-architecture.md patterns.
 */

import type { Page } from '@playwright/test';

/**
 * Navigate to a page and wait for it to be ready
 */
export async function navigateAndWait(page: Page, url: string): Promise<void> {
  await page.goto(url);
  await page.waitForLoadState('domcontentloaded');
}

/**
 * Navigate to the upload page
 */
export async function goToUpload(page: Page): Promise<void> {
  await navigateAndWait(page, '/upload');
}

/**
 * Navigate to the dashboard
 */
export async function goToDashboard(page: Page): Promise<void> {
  await navigateAndWait(page, '/dashboard');
}

/**
 * Navigate to the home page
 */
export async function goToHome(page: Page): Promise<void> {
  await navigateAndWait(page, '/');
}

/**
 * Wait for navigation to complete after a click
 */
export async function clickAndWaitForNavigation(
  page: Page,
  selector: string
): Promise<void> {
  await Promise.all([
    page.waitForNavigation({ waitUntil: 'domcontentloaded' }),
    page.click(selector),
  ]);
}
