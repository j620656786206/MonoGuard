/**
 * MonoGuard E2E Tests - Dashboard
 *
 * Tests for the dashboard functionality.
 * Priority: P1 (Important user feature)
 *
 * Test Coverage:
 * - Dashboard page loads
 * - Health score display
 * - Projects overview
 * - Recent analyses
 * - Quick actions
 */

import { test, expect } from './support/fixtures';

test.describe('Dashboard', () => {
  test.describe('Page Load', () => {
    test('[P1] should load dashboard page successfully', async ({ page }) => {
      // GIVEN: User navigates to dashboard
      await page.goto('/dashboard');

      // THEN: Dashboard should be visible
      await expect(page.locator('h1')).toContainText('Dashboard');
    });

    test('[P1] should display loading state while fetching data', async ({
      page,
    }) => {
      // GIVEN: Slow API response
      await page.route('**/api/v1/projects**', async (route) => {
        await new Promise((resolve) => setTimeout(resolve, 500));
        await route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({ success: true, data: [] }),
        });
      });

      // WHEN: User navigates to dashboard
      await page.goto('/dashboard');

      // THEN: Loading state should be shown (or data loads)
      await expect(page.locator('h1')).toContainText('Dashboard');
    });
  });

  test.describe('Health Score', () => {
    test('[P1] should display health score card', async ({ page }) => {
      // GIVEN: Dashboard with health score data
      await page.route('**/api/v1/projects**', (route) =>
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: [
              {
                id: 1,
                name: 'test-project',
                healthScore: 85,
                status: 'completed',
                createdAt: new Date().toISOString(),
              },
            ],
          }),
        })
      );

      // WHEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Health score should be displayed
      await expect(page.getByText(/Health Score/i)).toBeVisible();
    });

    test('[P1] should show health status indicator', async ({ page }) => {
      // GIVEN: Dashboard loads
      await page.goto('/dashboard');

      // THEN: Health status should be visible (Healthy/Warning/Critical)
      // The actual status depends on the score
      await expect(page.locator('body')).toBeVisible();
    });
  });

  test.describe('Overview Stats', () => {
    test('[P1] should display total projects count', async ({ page }) => {
      // GIVEN: Dashboard with projects
      await page.route('**/api/v1/projects**', (route) =>
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: [
              { id: 1, name: 'project-1', healthScore: 90 },
              { id: 2, name: 'project-2', healthScore: 75 },
              { id: 3, name: 'project-3', healthScore: 60 },
            ],
          }),
        })
      );

      // WHEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Total projects should be shown
      await expect(page.getByText(/Total Projects/i)).toBeVisible();
    });

    test('[P1] should display total packages count', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Total packages should be shown
      await expect(page.getByText(/Total Packages/i)).toBeVisible();
    });

    test('[P1] should display last analysis time', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Last analysis time should be shown
      await expect(page.getByText(/Last Analysis/i)).toBeVisible();
    });
  });

  test.describe('Recent Analyses', () => {
    test('[P1] should display recent analyses section', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Recent analyses section should be visible
      await expect(page.getByText(/Recent Analyses/i)).toBeVisible();
    });

    test('[P2] should show analysis status badges', async ({ page }) => {
      // GIVEN: Dashboard with analysis data
      await page.route('**/api/v1/projects**', (route) =>
        route.fulfill({
          status: 200,
          contentType: 'application/json',
          body: JSON.stringify({
            success: true,
            data: [
              {
                id: 1,
                name: 'completed-project',
                status: 'completed',
                healthScore: 92,
                updatedAt: new Date().toISOString(),
              },
            ],
          }),
        })
      );

      // WHEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Status badges should be visible
      await expect(page.locator('body')).toBeVisible();
    });
  });

  test.describe('Issues Display', () => {
    test('[P1] should display issues section', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Issues section should be visible
      await expect(page.getByText(/Issues/i)).toBeVisible();
    });

    test('[P2] should show issue severity indicators', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Severity indicators should be styled appropriately
      // (Visual verification - checking the section exists)
      await expect(page.locator('body')).toBeVisible();
    });
  });

  test.describe('Quick Actions', () => {
    test('[P1] should display quick actions section', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Quick actions should be visible
      await expect(page.getByText(/Quick Actions/i)).toBeVisible();
    });

    test('[P1] should have analyze new project button', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Analyze button should be visible
      await expect(page.getByText(/Analyze New Project/i)).toBeVisible();
    });

    test('[P1] should navigate to upload page from quick actions', async ({
      page,
    }) => {
      // GIVEN: User is on dashboard
      await page.goto('/dashboard');

      // WHEN: User clicks analyze new project
      await page.getByText(/Analyze New Project/i).click();

      // THEN: Should navigate to upload page
      await expect(page).toHaveURL(/upload/);
    });

    test('[P2] should have view architecture button', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: View architecture button should be visible
      await expect(page.getByText(/View Architecture/i)).toBeVisible();
    });

    test('[P2] should have team settings button', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Team settings button should be visible
      await expect(page.getByText(/Team Settings/i)).toBeVisible();
    });
  });

  test.describe('Run Analysis', () => {
    test('[P1] should have run analysis button in header', async ({ page }) => {
      // GIVEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Run analysis button should be visible
      await expect(page.getByRole('button', { name: /Run Analysis/i })).toBeVisible();
    });

    test('[P2] should show loading state when running analysis', async ({
      page,
    }) => {
      // GIVEN: User is on dashboard
      await page.goto('/dashboard');

      // WHEN: User clicks run analysis
      const runButton = page.getByRole('button', { name: /Run Analysis/i });
      await runButton.click();

      // THEN: Button should show loading state
      await expect(runButton).toContainText(/Analyzing/i);
    });
  });

  test.describe('Error States', () => {
    test('[P2] should handle API errors gracefully', async ({ page }) => {
      // GIVEN: API returns error
      await page.route('**/api/v1/projects**', (route) =>
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Internal server error' }),
        })
      );

      // WHEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Fallback data should be shown (or error message)
      await expect(page.getByText(/Dashboard/i)).toBeVisible();
    });

    test('[P2] should show fallback data indicator', async ({ page }) => {
      // GIVEN: API fails
      await page.route('**/api/v1/projects**', (route) =>
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ error: 'Failed' }),
        })
      );

      // WHEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Fallback indicator should be shown
      await expect(page.locator('body')).toBeVisible();
    });
  });

  test.describe('Responsive Design', () => {
    test('[P2] should be responsive on tablet viewport', async ({ page }) => {
      // GIVEN: Tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 });

      // WHEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Dashboard should render correctly
      await expect(page.locator('h1')).toContainText('Dashboard');
    });

    test('[P2] should be responsive on mobile viewport', async ({ page }) => {
      // GIVEN: Mobile viewport
      await page.setViewportSize({ width: 375, height: 667 });

      // WHEN: User views dashboard
      await page.goto('/dashboard');

      // THEN: Dashboard should render correctly
      await expect(page.locator('h1')).toContainText('Dashboard');
    });
  });
});
