# MonoGuard CI/CD Pipeline

This document describes the CI/CD pipeline configuration and usage for MonoGuard.

## Overview

The CI pipeline runs on **GitHub Actions** and includes:

- **Lint & Type Check**: Code quality validation
- **Unit Tests**: Fast feedback on logic errors
- **E2E Tests**: Browser-based testing with Playwright (4 parallel shards)
- **Burn-In Testing**: Flaky test detection (5 iterations)

## Pipeline Triggers

| Trigger | When | Stages |
|---------|------|--------|
| Pull Request | On PR to main/develop | All stages |
| Push | Direct push to main/develop | All stages |
| Schedule | Every Monday 6am UTC | Full burn-in (10 iterations) |

## Pipeline Stages

```
┌─────────────────────────────────────────────────────────────────┐
│                        INSTALL                                   │
│                    (Cache dependencies)                          │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
┌─────────────────────────┐     ┌─────────────────────────┐
│         LINT            │     │      UNIT TESTS         │
│    (ESLint + Types)     │     │    (Vitest + Jest)      │
└─────────────────────────┘     └─────────────────────────┘
              │                               │
              └───────────────┬───────────────┘
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     E2E TESTS (4 Shards)                         │
│         Shard 1  │  Shard 2  │  Shard 3  │  Shard 4             │
│        (parallel execution, fail-fast: false)                    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     BURN-IN (PR only)                            │
│               5 iterations × full test suite                     │
│           (detects flaky tests before merge)                     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     MERGE REPORTS                                │
│              (Combine shard results into HTML)                   │
└─────────────────────────────────────────────────────────────────┘
```

## Performance Targets

| Stage | Target | Notes |
|-------|--------|-------|
| Install | <3 min | Cached on subsequent runs |
| Lint | <2 min | Parallel across projects |
| Unit Tests | <5 min | Parallel execution |
| E2E (per shard) | <10 min | 4 shards = ~75% faster |
| Burn-In | <30 min | Only on PRs |
| **Total Pipeline** | <45 min | With parallelization |

## Running Locally

### Full CI Mirror

Run the same stages as CI locally:

```bash
./scripts/ci-local.sh
```

This runs:
1. Lint & Type Check
2. Unit Tests
3. E2E Tests (chromium)
4. Mini Burn-In (3 iterations)

### Burn-In Testing

Test stability before pushing:

```bash
# Default: 5 iterations
./scripts/burn-in.sh

# Custom iterations
./scripts/burn-in.sh 10

# Specific browser
./scripts/burn-in.sh 5 firefox
```

### Selective Testing

Run only tests affected by your changes:

```bash
# Compare to main
./scripts/test-changed.sh

# Compare to different branch
./scripts/test-changed.sh develop
```

## Artifacts

### On Failure

The following artifacts are uploaded when tests fail:

| Artifact | Retention | Contents |
|----------|-----------|----------|
| `e2e-results-shard-N` | 30 days | Test results, screenshots, videos, traces |
| `burn-in-failures` | 7 days | Failure artifacts from burn-in |
| `merged-playwright-report` | 30 days | Combined HTML report |

### Viewing Reports

Download artifacts from GitHub Actions and view locally:

```bash
# View HTML report
pnpm exec playwright show-report ./playwright-report

# View trace file
pnpm exec playwright show-trace ./test-results/trace.zip
```

## Caching Strategy

| Cache | Key | Contents |
|-------|-----|----------|
| pnpm store | `pnpm-lock.yaml` hash | Package cache |
| node_modules | `pnpm-lock.yaml` hash | Installed packages |
| Playwright browsers | `pnpm-lock.yaml` hash | Browser binaries |

Cache hit saves ~3-5 minutes per run.

## Debugging CI Failures

### 1. Check the Logs

View the job logs in GitHub Actions for error messages.

### 2. Download Artifacts

Download `e2e-results-shard-N` artifacts for:
- Screenshots at failure point
- Video recordings
- Playwright traces (best for debugging)

### 3. Reproduce Locally

```bash
# Run the exact same test
pnpm nx run web-e2e:e2e --project=chromium

# Debug mode
pnpm exec playwright test --debug

# With UI
pnpm exec playwright test --ui
```

### 4. Check for Flakiness

If a test passes locally but fails in CI:

```bash
# Run burn-in to detect flakiness
./scripts/burn-in.sh 10

# Check for timing issues
pnpm exec playwright test --repeat-each=5
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CI` | `true` (in CI) | Indicates CI environment |
| `BASE_URL` | `http://localhost:4200` | App URL for tests |
| `TEST_ENV` | `local` | Test environment name |

## Adding New Tests

When adding new E2E tests:

1. **Create test file** in `apps/web-e2e/src/`
2. **Use fixtures** from `support/fixtures`
3. **Run locally** with `pnpm nx run web-e2e:e2e`
4. **Run burn-in** with `./scripts/burn-in.sh 5`
5. **Push** when stable

## Maintenance

### Weekly Tasks

- Review burn-in results from Monday cron job
- Check artifact storage usage
- Update browser versions if needed

### Monthly Tasks

- Review test execution times
- Adjust shard count if needed
- Clean up flaky tests

## Related Documentation

- [E2E Test README](../apps/web-e2e/README.md)
- [Playwright Documentation](https://playwright.dev/docs/intro)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
