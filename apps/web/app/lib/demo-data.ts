/**
 * Demo mock data for the MonoGuard demo assembly.
 * Provides pre-built analysis results and dependency graph data
 * so the demo runs entirely client-side with no backend or WASM.
 */

import type {
  CircularDependencyInfo,
  ComprehensiveAnalysisResult,
  DependencyGraph,
} from '@monoguard/types'
import { createCompletedAnalysis } from '../../src/__tests__/factories/analysis.factory'

// ---------------------------------------------------------------------------
// 1. Demo analysis result (from test factory)
// ---------------------------------------------------------------------------
export const demoAnalysis: ComprehensiveAnalysisResult = createCompletedAnalysis()

// ---------------------------------------------------------------------------
// 2. Demo dependency graph — 12 packages, 2 circular dependency cycles
// ---------------------------------------------------------------------------
export const demoDependencyGraph: DependencyGraph = {
  rootPath: '/workspace',
  workspaceType: 'pnpm',
  nodes: {
    '@acme/app': {
      name: '@acme/app',
      version: '1.0.0',
      path: 'apps/app',
      dependencies: ['@acme/ui', '@acme/auth', '@acme/api'],
      devDependencies: ['@acme/test-utils'],
      peerDependencies: [],
    },
    '@acme/ui': {
      name: '@acme/ui',
      version: '2.1.0',
      path: 'packages/ui',
      dependencies: ['@acme/theme', '@acme/icons'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/auth': {
      name: '@acme/auth',
      version: '1.3.0',
      path: 'packages/auth',
      dependencies: ['@acme/core', '@acme/logger'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/core': {
      name: '@acme/core',
      version: '3.0.0',
      path: 'packages/core',
      dependencies: ['@acme/types', '@acme/config'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/types': {
      name: '@acme/types',
      version: '1.0.0',
      path: 'packages/types',
      dependencies: ['@acme/auth'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/api': {
      name: '@acme/api',
      version: '2.0.0',
      path: 'packages/api',
      dependencies: ['@acme/core', '@acme/logger'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/logger': {
      name: '@acme/logger',
      version: '1.1.0',
      path: 'packages/logger',
      dependencies: ['@acme/config'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/config': {
      name: '@acme/config',
      version: '1.0.0',
      path: 'packages/config',
      dependencies: ['@acme/logger'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/theme': {
      name: '@acme/theme',
      version: '1.2.0',
      path: 'packages/theme',
      dependencies: ['@acme/types'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/icons': {
      name: '@acme/icons',
      version: '1.0.0',
      path: 'packages/icons',
      dependencies: [],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/test-utils': {
      name: '@acme/test-utils',
      version: '1.0.0',
      path: 'packages/test-utils',
      dependencies: ['@acme/core'],
      devDependencies: [],
      peerDependencies: [],
    },
    '@acme/docs': {
      name: '@acme/docs',
      version: '1.0.0',
      path: 'apps/docs',
      dependencies: ['@acme/ui', '@acme/theme'],
      devDependencies: [],
      peerDependencies: [],
    },
  },
  edges: [
    // @acme/app edges
    { from: '@acme/app', to: '@acme/ui', type: 'production', versionRange: '^2.1.0' },
    { from: '@acme/app', to: '@acme/auth', type: 'production', versionRange: '^1.3.0' },
    { from: '@acme/app', to: '@acme/api', type: 'production', versionRange: '^2.0.0' },
    { from: '@acme/app', to: '@acme/test-utils', type: 'development', versionRange: '^1.0.0' },
    // @acme/ui edges
    { from: '@acme/ui', to: '@acme/theme', type: 'production', versionRange: '^1.2.0' },
    { from: '@acme/ui', to: '@acme/icons', type: 'production', versionRange: '^1.0.0' },
    // @acme/auth edges
    { from: '@acme/auth', to: '@acme/core', type: 'production', versionRange: '^3.0.0' },
    { from: '@acme/auth', to: '@acme/logger', type: 'production', versionRange: '^1.1.0' },
    // @acme/core edges
    { from: '@acme/core', to: '@acme/types', type: 'production', versionRange: '^1.0.0' },
    { from: '@acme/core', to: '@acme/config', type: 'production', versionRange: '^1.0.0' },
    // @acme/types → @acme/auth (creates indirect cycle: auth → core → types → auth)
    { from: '@acme/types', to: '@acme/auth', type: 'production', versionRange: '^1.3.0' },
    // @acme/api edges
    { from: '@acme/api', to: '@acme/core', type: 'production', versionRange: '^3.0.0' },
    { from: '@acme/api', to: '@acme/logger', type: 'production', versionRange: '^1.1.0' },
    // @acme/logger ↔ @acme/config (direct cycle)
    { from: '@acme/logger', to: '@acme/config', type: 'production', versionRange: '^1.0.0' },
    { from: '@acme/config', to: '@acme/logger', type: 'production', versionRange: '^1.1.0' },
    // @acme/theme edges
    { from: '@acme/theme', to: '@acme/types', type: 'production', versionRange: '^1.0.0' },
    // @acme/test-utils edges
    { from: '@acme/test-utils', to: '@acme/core', type: 'production', versionRange: '^3.0.0' },
    // @acme/docs edges
    { from: '@acme/docs', to: '@acme/ui', type: 'production', versionRange: '^2.1.0' },
    { from: '@acme/docs', to: '@acme/theme', type: 'production', versionRange: '^1.2.0' },
  ],
}

// ---------------------------------------------------------------------------
// 3. Demo circular dependencies (matching graph cycles)
// ---------------------------------------------------------------------------
export const demoCircularDependencies: CircularDependencyInfo[] = [
  {
    cycle: ['@acme/auth', '@acme/core', '@acme/types', '@acme/auth'],
    type: 'indirect',
    severity: 'critical',
    depth: 3,
    impact:
      'Creates tight coupling between auth, core, and types packages. Prevents independent versioning and causes cascading rebuilds.',
    complexity: 7,
    priorityScore: 90,
  },
  {
    cycle: ['@acme/logger', '@acme/config', '@acme/logger'],
    type: 'direct',
    severity: 'warning',
    depth: 2,
    impact:
      'Logger and config mutually depend on each other. May cause initialization ordering issues at runtime.',
    complexity: 4,
    priorityScore: 65,
  },
]
