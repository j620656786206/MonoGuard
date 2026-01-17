/**
 * MonoGuard E2E Test Fixtures
 *
 * This module provides the merged fixture system following the TEA knowledge base patterns:
 * - Pure function → fixture → mergeTests composition
 * - Auto-cleanup for created resources
 * - Single responsibility per fixture
 *
 * Usage:
 *   import { test, expect } from '../support/fixtures';
 *
 *   test('my test', async ({ page, analysisFactory }) => {
 *     // Use fixtures in tests
 *   });
 */

import { test as base, expect, mergeTests } from '@playwright/test'
import { test as analysisFixture } from './analysis-fixture'

// Merge all fixtures into a single test object
export const test = mergeTests(base, analysisFixture)

export { expect }
