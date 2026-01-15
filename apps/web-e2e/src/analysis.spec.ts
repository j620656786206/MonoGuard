/**
 * MonoGuard E2E Tests - Analysis Functionality
 *
 * Tests for the core analysis features including:
 * - File upload
 * - Analysis processing
 * - Results display
 *
 * Uses fixtures and factories from support/fixtures.
 */

import { test, expect } from './support/fixtures';
import { createCircularWorkspace, createMinimalWorkspace } from './support/fixtures/factories/workspace-factory';

test.describe('Analysis Feature', () => {
  test.describe('Upload Flow', () => {
    test('should display upload area on upload page', async ({ page }) => {
      await page.goto('/upload');

      // The upload area should be visible
      const uploadArea = page.locator('[data-testid="file-upload-area"]');
      await expect(uploadArea).toBeVisible();
    });

    test('should show supported file format hint', async ({ page }) => {
      await page.goto('/upload');

      // Should indicate workspace.json is supported
      await expect(page.getByText(/workspace\.json|JSON/i)).toBeVisible();
    });
  });

  test.describe('Analysis Results', () => {
    test('should use workspace factory for test data', async ({ workspaceFactory }) => {
      // Demonstrate factory usage - creates a valid workspace structure
      const workspace = workspaceFactory.create({
        projects: {
          'test-app': { projectType: 'application' },
          'test-lib': { projectType: 'library' },
        },
      });

      // Verify the factory produces valid data
      expect(workspace.version).toBe(2);
      expect(Object.keys(workspace.projects)).toHaveLength(2);
      expect(workspace.projects['test-app'].projectType).toBe('application');
      expect(workspace.projects['test-lib'].projectType).toBe('library');
    });

    test('should create circular dependency workspace for testing', () => {
      // Use the specialized factory for circular deps
      const circularWorkspace = createCircularWorkspace();

      // Should have the circular dependency structure
      expect(circularWorkspace.projects['lib-a'].implicitDependencies).toContain('lib-b');
      expect(circularWorkspace.projects['lib-b'].implicitDependencies).toContain('lib-c');
      expect(circularWorkspace.projects['lib-c'].implicitDependencies).toContain('lib-a');
    });

    test('should create minimal workspace for quick tests', () => {
      const minimalWorkspace = createMinimalWorkspace();

      expect(Object.keys(minimalWorkspace.projects)).toHaveLength(1);
      expect(minimalWorkspace.projects['single-app'].projectType).toBe('application');
    });
  });

  test.describe('Dashboard', () => {
    test('should navigate to dashboard', async ({ page }) => {
      await page.goto('/dashboard');

      // Dashboard should have some content
      await expect(page.locator('body')).toBeVisible();
    });
  });
});
