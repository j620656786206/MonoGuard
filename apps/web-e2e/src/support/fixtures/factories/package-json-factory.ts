/**
 * Package.json Factory
 *
 * Creates realistic package.json test data for MonoGuard upload tests.
 * Uses the factory pattern with overrides for flexible test data generation.
 */

export type PackageJson = {
  name: string
  version: string
  dependencies?: Record<string, string>
  devDependencies?: Record<string, string>
  peerDependencies?: Record<string, string>
  scripts?: Record<string, string>
  private?: boolean
}

export type PackageJsonOverrides = Partial<PackageJson> & {
  dependencyCount?: number
  devDependencyCount?: number
  includeVersionConflicts?: boolean
  includeDuplicates?: boolean
}

/**
 * Common dependencies used in monorepos
 */
const COMMON_DEPS: Record<string, string> = {
  react: '^18.2.0',
  'react-dom': '^18.2.0',
  next: '^14.0.0',
  typescript: '^5.3.0',
  lodash: '^4.17.21',
  axios: '^1.6.0',
  zod: '^3.22.0',
  '@tanstack/react-query': '^5.0.0',
  zustand: '^4.4.0',
  tailwindcss: '^3.4.0',
}

const COMMON_DEV_DEPS: Record<string, string> = {
  '@types/react': '^18.2.0',
  '@types/node': '^20.0.0',
  eslint: '^8.56.0',
  prettier: '^3.2.0',
  vitest: '^1.2.0',
  '@playwright/test': '^1.41.0',
  tsup: '^8.0.0',
}

/**
 * Creates a default package.json with sensible defaults
 */
export function createPackageJson(overrides: PackageJsonOverrides = {}): PackageJson {
  const {
    dependencyCount,
    devDependencyCount,
    includeVersionConflicts,
    includeDuplicates,
    ...rest
  } = overrides

  // Generate dependencies based on count or use defaults
  let dependencies: Record<string, string> = {}
  let devDependencies: Record<string, string> = {}

  if (dependencyCount !== undefined) {
    const depNames = Object.keys(COMMON_DEPS).slice(0, dependencyCount)
    dependencies = depNames.reduce(
      (acc, name) => {
        acc[name] = COMMON_DEPS[name]
        return acc
      },
      {} as Record<string, string>
    )
  } else {
    dependencies = {
      react: '^18.2.0',
      'react-dom': '^18.2.0',
      next: '^14.0.0',
    }
  }

  if (devDependencyCount !== undefined) {
    const devDepNames = Object.keys(COMMON_DEV_DEPS).slice(0, devDependencyCount)
    devDependencies = devDepNames.reduce(
      (acc, name) => {
        acc[name] = COMMON_DEV_DEPS[name]
        return acc
      },
      {} as Record<string, string>
    )
  } else {
    devDependencies = {
      '@types/react': '^18.2.0',
      typescript: '^5.3.0',
    }
  }

  // Add version conflicts if requested
  if (includeVersionConflicts) {
    dependencies['lodash'] = '^4.17.21'
    devDependencies['lodash'] = '^4.17.20' // Intentional version conflict
  }

  // Add duplicates if requested
  if (includeDuplicates) {
    dependencies['lodash'] = '^4.17.21'
    dependencies['lodash-es'] = '^4.17.21' // Duplicate functionality
  }

  return {
    name: 'test-package',
    version: '1.0.0',
    private: true,
    dependencies,
    devDependencies,
    scripts: {
      dev: 'next dev',
      build: 'next build',
      start: 'next start',
      test: 'vitest',
    },
    ...rest,
  }
}

/**
 * Creates a monorepo root package.json
 */
export function createMonorepoRootPackageJson(overrides: Partial<PackageJson> = {}): PackageJson {
  return {
    name: 'test-monorepo',
    version: '1.0.0',
    private: true,
    scripts: {
      build: 'nx run-many --target=build --all',
      test: 'nx run-many --target=test --all',
      lint: 'nx run-many --target=lint --all',
    },
    devDependencies: {
      nx: '^17.0.0',
      '@nx/js': '^17.0.0',
      '@nx/react': '^17.0.0',
      typescript: '^5.3.0',
    },
    ...overrides,
  }
}

/**
 * Creates a package.json with circular dependency indicators
 */
export function createPackageWithCircularDeps(name: string, dependsOn: string[]): PackageJson {
  const dependencies: Record<string, string> = {}
  dependsOn.forEach((dep) => {
    dependencies[dep] = 'workspace:*'
  })

  return {
    name,
    version: '1.0.0',
    private: true,
    dependencies,
    devDependencies: {
      typescript: '^5.3.0',
    },
  }
}

/**
 * Creates a package.json with security vulnerabilities (for testing detection)
 */
export function createVulnerablePackageJson(): PackageJson {
  return {
    name: 'vulnerable-package',
    version: '1.0.0',
    dependencies: {
      lodash: '4.17.4', // Known vulnerable version
      minimist: '0.0.8', // Known vulnerable version
      'node-fetch': '2.6.0', // Known vulnerable version
    },
    devDependencies: {},
  }
}

/**
 * Creates a large package.json for performance testing
 */
export function createLargePackageJson(dependencyCount: number = 100): PackageJson {
  const dependencies: Record<string, string> = {}
  const devDependencies: Record<string, string> = {}

  for (let i = 0; i < dependencyCount; i++) {
    dependencies[`dep-${i}`] = `^${Math.floor(Math.random() * 10)}.0.0`
  }

  for (let i = 0; i < Math.floor(dependencyCount / 2); i++) {
    devDependencies[`dev-dep-${i}`] = `^${Math.floor(Math.random() * 10)}.0.0`
  }

  return {
    name: 'large-package',
    version: '1.0.0',
    private: true,
    dependencies,
    devDependencies,
  }
}

/**
 * Converts a package.json object to a File-like object for upload tests
 */
export function packageJsonToFile(
  packageJson: PackageJson,
  filename: string = 'package.json'
): { name: string; mimeType: string; buffer: Buffer } {
  const content = JSON.stringify(packageJson, null, 2)
  return {
    name: filename,
    mimeType: 'application/json',
    buffer: Buffer.from(content),
  }
}
