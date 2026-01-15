# MonoGuard E2E Tests

End-to-end testing suite for MonoGuard using Playwright.

## Quick Start

```bash
# Install dependencies
pnpm install

# Install Playwright browsers
pnpm exec playwright install

# Run tests
pnpm nx run web-e2e:e2e

# Run tests in UI mode
pnpm exec playwright test --ui

# Run tests in headed mode
pnpm exec playwright test --headed
```

## Architecture Overview

This test suite follows the **TEA (Test Engineer Architecture)** patterns from the BMAD knowledge base.

### Directory Structure

```
apps/web-e2e/
├── playwright.config.ts      # Playwright configuration
├── src/
│   ├── example.spec.ts      # Landing page tests
│   ├── analysis.spec.ts     # Analysis functionality tests
│   ├── upload-flow.spec.ts  # Upload flow tests (P0)
│   ├── dashboard.spec.ts    # Dashboard tests (P1)
│   └── support/
│       ├── fixtures/        # Test fixtures (mergeTests pattern)
│       │   ├── index.ts     # Merged fixture export
│       │   ├── analysis-fixture.ts
│       │   └── factories/   # Data factories
│       │       ├── workspace-factory.ts
│       │       └── package-json-factory.ts
│       ├── helpers/         # Pure function helpers
│       │   ├── navigation.ts
│       │   ├── assertions.ts
│       │   └── index.ts
│       └── page-objects/    # (Optional) Page objects
```

### Test Files Overview

| File | Priority | Coverage |
|------|----------|----------|
| `upload-flow.spec.ts` | P0-P1 | File upload, analysis trigger, error handling |
| `dashboard.spec.ts` | P1-P2 | Health score, projects, issues, quick actions |
| `analysis.spec.ts` | P1-P2 | Analysis results, workspace factory tests |
| `example.spec.ts` | P1-P2 | Landing page, navigation, responsiveness |

### Key Patterns

#### 1. Fixture Architecture

We use the **pure function → fixture → mergeTests** composition pattern:

```typescript
// Import the merged fixtures
import { test, expect } from './support/fixtures';

test('my test', async ({ page, workspaceFactory }) => {
  // Use fixtures directly in tests
  const workspace = workspaceFactory.create({ projectCount: 5 });
});
```

#### 2. Data Factories

Factory functions with overrides for flexible test data:

```typescript
import { createWorkspaceJson, createCircularWorkspace } from './support/fixtures/factories/workspace-factory';

// Default workspace
const workspace = createWorkspaceJson();

// Custom projects
const customWorkspace = createWorkspaceJson({
  projects: {
    'my-app': { projectType: 'application' },
    'my-lib': { tags: ['shared'] },
  },
});

// Specialized factories
const circularDeps = createCircularWorkspace();
```

#### 3. Helper Functions

Pure, framework-agnostic functions:

```typescript
import { goToUpload, assertHealthScore } from './support/helpers';

test('analysis flow', async ({ page }) => {
  await goToUpload(page);
  await assertHealthScore(page, 70, 100);
});
```

## Timeout Standards

| Timeout Type | Value | Description |
|-------------|-------|-------------|
| Action | 15s | Clicks, fills, etc. |
| Navigation | 30s | page.goto, reload |
| Expect | 10s | Assertions |
| Test | 60s | Overall test timeout |

## Running Tests

### Local Development

```bash
# Run all tests
pnpm nx run web-e2e:e2e

# Run specific test file
pnpm exec playwright test src/analysis.spec.ts

# Run tests with UI
pnpm exec playwright test --ui

# Debug mode
pnpm exec playwright test --debug

# Headed mode (see browser)
pnpm exec playwright test --headed
```

### Specific Browsers

```bash
# Chromium only
pnpm exec playwright test --project=chromium

# Firefox only
pnpm exec playwright test --project=firefox

# WebKit (Safari) only
pnpm exec playwright test --project=webkit
```

### CI Environment

Tests automatically detect CI via the `CI` environment variable:

- Retries: 2 (vs 0 locally)
- Workers: 1 (serial execution for stability)
- Artifacts: Captured on failure only

## Artifacts

On test failure, the following are captured:

- **Screenshots**: `test-results/*.png`
- **Videos**: `test-results/*.webm`
- **Traces**: `test-results/*.zip` (viewable with `npx playwright show-trace`)
- **HTML Report**: `playwright-report/index.html`

### Viewing Reports

```bash
# Open HTML report
pnpm exec playwright show-report apps/web-e2e/playwright-report

# View a trace file
pnpm exec playwright show-trace apps/web-e2e/test-results/trace.zip
```

## Best Practices

### Selector Strategy

Always use `data-testid` attributes:

```typescript
// ✅ Good
await page.click('[data-testid="submit-button"]');
await page.fill('[data-testid="email-input"]', 'test@example.com');

// ❌ Avoid
await page.click('.btn-primary');
await page.click('#submit');
await page.click('button:has-text("Submit")');
```

### Test Isolation

Each test should:
- Create its own test data via factories
- Not depend on other tests
- Clean up after itself (handled by fixture teardown)

### No Hard Waits

```typescript
// ✅ Good - wait for element
await expect(page.locator('[data-testid="result"]')).toBeVisible();

// ❌ Bad - arbitrary wait
await page.waitForTimeout(3000);
```

## Knowledge Base References

This test architecture is based on:

- `fixture-architecture.md` - Composable fixtures with mergeTests
- `data-factories.md` - Factory pattern with overrides
- `playwright-config.md` - Timeout standards and artifact configuration
- `test-quality.md` - Deterministic test design principles

## Troubleshooting

### Tests timeout frequently

1. Check if the dev server is running: `pnpm nx run web:dev`
2. Increase timeouts in `playwright.config.ts` for slow CI
3. Check network conditions

### Browser installation issues

```bash
# Reinstall browsers
pnpm exec playwright install --with-deps
```

### Test results not uploading in CI

Ensure your CI workflow uploads the `test-results/` and `playwright-report/` directories on failure.

## Related Documentation

- [Playwright Documentation](https://playwright.dev/docs/intro)
- [Nx Playwright Plugin](https://nx.dev/packages/playwright)
- [MonoGuard Architecture](/_bmad-output/planning-artifacts/)
