/**
 * Analysis Fixture
 *
 * Provides test utilities for MonoGuard's core analysis functionality.
 * Follows the pure function → fixture pattern with auto-cleanup.
 */

import { test as base, Page } from '@playwright/test';
import { createWorkspaceJson, WorkspaceJsonOverrides } from '../fixtures/factories/workspace-factory';
import {
  createPackageJson,
  packageJsonToFile,
  PackageJsonOverrides,
} from '../fixtures/factories/package-json-factory';

type AnalysisFixtures = {
  /**
   * Factory for creating workspace.json test data
   */
  workspaceFactory: {
    create: (overrides?: WorkspaceJsonOverrides) => ReturnType<typeof createWorkspaceJson>;
  };

  /**
   * Factory for creating package.json test data
   */
  packageJsonFactory: {
    create: (overrides?: PackageJsonOverrides) => ReturnType<typeof createPackageJson>;
    toFile: (
      packageJson: ReturnType<typeof createPackageJson>,
      filename?: string
    ) => ReturnType<typeof packageJsonToFile>;
  };

  /**
   * Helper to upload a file to the analysis page
   */
  uploadWorkspace: (workspaceData: object) => Promise<void>;

  /**
   * Helper to wait for analysis to complete
   */
  waitForAnalysis: () => Promise<void>;
};

export const test = base.extend<AnalysisFixtures>({
  workspaceFactory: async ({}, use) => {
    await use({
      create: (overrides) => createWorkspaceJson(overrides),
    });
  },

  packageJsonFactory: async ({}, use) => {
    await use({
      create: (overrides) => createPackageJson(overrides),
      toFile: (packageJson, filename) => packageJsonToFile(packageJson, filename),
    });
  },

  uploadWorkspace: async ({ page }, use) => {
    const uploadWorkspace = async (workspaceData: object) => {
      // Navigate to upload page if not already there
      const currentUrl = page.url();
      if (!currentUrl.includes('/upload')) {
        await page.goto('/upload');
      }

      // Create a virtual file from the workspace data
      const jsonContent = JSON.stringify(workspaceData, null, 2);
      const buffer = Buffer.from(jsonContent);

      // Use Playwright's file chooser to upload
      const fileChooserPromise = page.waitForEvent('filechooser');

      // Click the upload area/button to trigger file chooser
      await page.click('[data-testid="file-upload-area"]');

      const fileChooser = await fileChooserPromise;
      await fileChooser.setFiles({
        name: 'workspace.json',
        mimeType: 'application/json',
        buffer,
      });
    };

    await use(uploadWorkspace);
  },

  waitForAnalysis: async ({ page }, use) => {
    const waitForAnalysis = async () => {
      // Wait for the analysis to complete by checking for results
      await page.waitForSelector('[data-testid="analysis-results"]', {
        state: 'visible',
        timeout: 30000,
      });
    };

    await use(waitForAnalysis);
  },
});
