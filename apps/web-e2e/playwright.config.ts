import { defineConfig, devices } from '@playwright/test';
import { workspaceRoot } from '@nx/devkit';
import path from 'path';

/**
 * MonoGuard E2E Test Configuration
 */

const baseURL = process.env['BASE_URL'] || 'http://localhost:3000';
const isCI = !!process.env['CI'];

export default defineConfig({
  testDir: './src',
  timeout: 60 * 1000,
  expect: {
    timeout: 10 * 1000,
  },
  fullyParallel: true,
  forbidOnly: isCI,
  retries: isCI ? 2 : 0,
  workers: isCI ? 1 : undefined,
  outputDir: path.join(__dirname, 'test-results'),
  reporter: [
    ['html', { outputFolder: path.join(__dirname, 'playwright-report'), open: 'never' }],
    ['list'],
  ],
  use: {
    baseURL,
    actionTimeout: 15 * 1000,
    navigationTimeout: 30 * 1000,
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  webServer: {
    command: 'pnpm exec nx run web:dev',
    url: 'http://localhost:3000',
    reuseExistingServer: !isCI,
    timeout: 120 * 1000,
    cwd: workspaceRoot,
  },
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        baseURL,
      },
    },
  ],
});
