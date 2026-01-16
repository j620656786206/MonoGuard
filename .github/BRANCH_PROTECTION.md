# Branch Protection Rules

This document describes the recommended branch protection settings for the MonoGuard repository.

## Protected Branches

### `main` Branch

The `main` branch is the production branch and requires the following protection rules:

#### Required Status Checks

The following CI jobs must pass before merging:

| Status Check | Description                         |
| ------------ | ----------------------------------- |
| `Lint`       | ESLint and TypeScript type checking |
| `Test`       | Unit and integration tests          |
| `Build`      | Build verification for all packages |
| `CI Summary` | Aggregated pipeline status          |

#### Branch Protection Settings

1. **Require a pull request before merging**
   - Require approvals: 1 (optional for solo development)
   - Dismiss stale pull request approvals when new commits are pushed: Yes
   - Require review from Code Owners: No (optional)

2. **Require status checks to pass before merging**
   - Require branches to be up to date before merging: Yes
   - Status checks that are required:
     - `Lint`
     - `Test`
     - `Build`
     - `CI Summary`

3. **Require conversation resolution before merging**: Yes

4. **Require signed commits**: Optional (recommended for team projects)

5. **Require linear history**: Optional (enables cleaner git history)

6. **Do not allow bypassing the above settings**: Recommended

### `develop` Branch

The `develop` branch is the integration branch with similar but slightly relaxed rules:

1. **Require status checks to pass before merging**: Yes
2. **Required status checks**: Same as `main`
3. **Require approvals**: Optional (0 for fast iteration)

## Setting Up Branch Protection

### Via GitHub Web Interface

1. Go to **Settings** > **Branches**
2. Click **Add branch protection rule**
3. Enter `main` as the branch name pattern
4. Configure settings as described above
5. Click **Create**
6. Repeat for `develop` branch

### Via GitHub CLI

```bash
# For main branch
gh api repos/{owner}/{repo}/branches/main/protection -X PUT -f required_status_checks='{"strict":true,"contexts":["Lint","Test","Build","CI Summary"]}' -f enforce_admins=true -f required_pull_request_reviews='{"required_approving_review_count":1}'

# For develop branch
gh api repos/{owner}/{repo}/branches/develop/protection -X PUT -f required_status_checks='{"strict":true,"contexts":["Lint","Test","Build","CI Summary"]}'
```

## E2E Tests

The E2E test workflow runs separately and includes:

- Parallel sharding (4 shards)
- Burn-in flaky detection
- Comprehensive Playwright tests

E2E tests are run on PRs but may be configured as optional for faster iteration.

## Notes

- The CI pipeline is designed to complete in < 5 minutes for typical PRs
- Caching is enabled for pnpm, Nx, and Go modules
- Use `nx affected` commands to only test/build changed projects
