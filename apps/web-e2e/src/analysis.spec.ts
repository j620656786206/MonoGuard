/**
 * MonoGuard E2E Tests - Analysis Functionality
 *
 * Tests for the core analysis features including:
 * - Analyze page display
 * - Workspace upload
 * - Results display
 *
 * Uses fixtures and factories from support/fixtures.
 */

import { expect, test } from './support/fixtures'
import {
  createCircularWorkspace,
  createMinimalWorkspace,
} from './support/fixtures/factories/workspace-factory'

test.describe('Analysis Feature', () => {
  test.describe('Analyze Page', () => {
    test('should display analyze page heading', async ({ page }) => {
      await page.goto('/analyze')

      // The analyze page should have proper heading
      await expect(page.locator('h1')).toContainText('Analyze')
    })

    test('should show upload instructions', async ({ page }) => {
      await page.goto('/analyze')

      // Should indicate workspace.json is supported
      await expect(page.getByText(/workspace\.json/)).toBeVisible()
    })

    test('should have Select File button', async ({ page }) => {
      await page.goto('/analyze')

      // Should have a button to select file
      await expect(page.getByText(/Select File/i)).toBeVisible()
    })
  })

  test.describe('Results Page', () => {
    test('should display results page heading', async ({ page }) => {
      await page.goto('/results')

      // Results page should have proper heading
      await expect(page.locator('h1')).toContainText('Analysis Results')
    })

    test('should show Health Score card', async ({ page }) => {
      await page.goto('/results')

      // Health Score section should be visible
      await expect(page.getByText(/Health Score/i)).toBeVisible()
    })

    test('should show Circular Dependencies card', async ({ page }) => {
      await page.goto('/results')

      // Circular Dependencies section should be visible
      await expect(page.getByText(/Circular Dependencies/i)).toBeVisible()
    })

    test('should show Total Packages card', async ({ page }) => {
      await page.goto('/results')

      // Total Packages section should be visible
      await expect(page.getByText(/Total Packages/i)).toBeVisible()
    })

    test('should have Start New Analysis link', async ({ page }) => {
      await page.goto('/results')

      // Should have a link to start new analysis
      await expect(page.getByText(/Start New Analysis/i)).toBeVisible()
    })

    test('should navigate to analyze page when clicking Start New Analysis', async ({ page }) => {
      await page.goto('/results')

      // Click Start New Analysis
      await page.getByText(/Start New Analysis/i).click()

      // Should navigate to analyze page
      await expect(page).toHaveURL(/analyze/)
    })
  })

  test.describe('Workspace Factory Tests', () => {
    test('should use workspace factory for test data', async ({ workspaceFactory }) => {
      // Demonstrate factory usage - creates a valid workspace structure
      const workspace = workspaceFactory.create({
        projects: {
          'test-app': { projectType: 'application' },
          'test-lib': { projectType: 'library' },
        },
      })

      // Verify the factory produces valid data
      expect(workspace.version).toBe(2)
      expect(Object.keys(workspace.projects)).toHaveLength(2)
      expect(workspace.projects['test-app'].projectType).toBe('application')
      expect(workspace.projects['test-lib'].projectType).toBe('library')
    })

    test('should create circular dependency workspace for testing', () => {
      // Use the specialized factory for circular deps
      const circularWorkspace = createCircularWorkspace()

      // Should have the circular dependency structure
      expect(circularWorkspace.projects['lib-a'].implicitDependencies).toContain('lib-b')
      expect(circularWorkspace.projects['lib-b'].implicitDependencies).toContain('lib-c')
      expect(circularWorkspace.projects['lib-c'].implicitDependencies).toContain('lib-a')
    })

    test('should create minimal workspace for quick tests', () => {
      const minimalWorkspace = createMinimalWorkspace()

      expect(Object.keys(minimalWorkspace.projects)).toHaveLength(1)
      expect(minimalWorkspace.projects['single-app'].projectType).toBe('application')
    })
  })
})
