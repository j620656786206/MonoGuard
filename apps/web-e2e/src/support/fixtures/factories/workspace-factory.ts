/**
 * Workspace Factory
 *
 * Creates realistic workspace.json test data for MonoGuard analysis tests.
 * Uses the factory pattern with overrides for flexible test data generation.
 *
 * Pattern: Pure function with Partial<T> overrides (TEA knowledge base: data-factories.md)
 */

export type NxProject = {
  root: string
  sourceRoot: string
  projectType: 'application' | 'library'
  targets?: Record<string, { executor: string; options?: Record<string, unknown> }>
  tags?: string[]
  implicitDependencies?: string[]
}

export type WorkspaceJson = {
  version: number
  projects: Record<string, NxProject>
  defaultProject?: string
}

export type WorkspaceJsonOverrides = {
  version?: number
  projects?: Record<string, Partial<NxProject>>
  defaultProject?: string
  projectCount?: number
  includeCircularDeps?: boolean
}

/**
 * Creates a default NxProject with sensible defaults
 */
export function createProject(name: string, overrides: Partial<NxProject> = {}): NxProject {
  const projectType = overrides.projectType ?? 'library'
  const root = overrides.root ?? `packages/${name}`

  return {
    root,
    sourceRoot: `${root}/src`,
    projectType,
    targets: {
      build: {
        executor: '@nx/js:tsc',
        options: {
          outputPath: `dist/${root}`,
          main: `${root}/src/index.ts`,
          tsConfig: `${root}/tsconfig.lib.json`,
        },
      },
      test: {
        executor: '@nx/jest:jest',
        options: {
          jestConfig: `${root}/jest.config.ts`,
        },
      },
      lint: {
        executor: '@nx/eslint:lint',
        options: {
          lintFilePatterns: [`${root}/**/*.ts`],
        },
      },
    },
    tags: [],
    ...overrides,
  }
}

/**
 * Creates a workspace.json with configurable projects
 *
 * @example
 * // Default workspace with 3 projects
 * const workspace = createWorkspaceJson();
 *
 * @example
 * // Custom workspace with specific projects
 * const workspace = createWorkspaceJson({
 *   projects: {
 *     'my-app': { projectType: 'application' },
 *     'my-lib': { tags: ['shared'] },
 *   },
 * });
 *
 * @example
 * // Workspace with circular dependencies for testing detection
 * const workspace = createWorkspaceJson({ includeCircularDeps: true });
 */
export function createWorkspaceJson(overrides: WorkspaceJsonOverrides = {}): WorkspaceJson {
  const {
    version = 2,
    projects = {},
    defaultProject,
    includeCircularDeps = false,
    projectCount,
  } = overrides

  // If projectCount is specified, generate that many projects
  let finalProjects: Record<string, NxProject> = {}

  if (projectCount !== undefined) {
    for (let i = 1; i <= projectCount; i++) {
      const name = `project-${i}`
      finalProjects[name] = createProject(name, {
        projectType: i === 1 ? 'application' : 'library',
      })
    }
  } else if (Object.keys(projects).length === 0) {
    // Default projects if none specified
    finalProjects = {
      web: createProject('web', { projectType: 'application', root: 'apps/web' }),
      types: createProject('types', { root: 'packages/types' }),
      'ui-components': createProject('ui-components', {
        root: 'packages/ui-components',
        tags: ['scope:shared', 'type:ui'],
      }),
    }
  } else {
    // Use provided projects with createProject for defaults
    for (const [name, projectOverrides] of Object.entries(projects)) {
      finalProjects[name] = createProject(name, projectOverrides)
    }
  }

  // Add circular dependencies if requested (for testing detection)
  if (includeCircularDeps && Object.keys(finalProjects).length >= 2) {
    const projectNames = Object.keys(finalProjects)
    const first = projectNames[0]
    const second = projectNames[1]

    finalProjects[first] = {
      ...finalProjects[first],
      implicitDependencies: [second],
    }
    finalProjects[second] = {
      ...finalProjects[second],
      implicitDependencies: [first],
    }
  }

  return {
    version,
    projects: finalProjects,
    ...(defaultProject && { defaultProject }),
  }
}

/**
 * Creates a minimal workspace.json for quick tests
 */
export function createMinimalWorkspace(): WorkspaceJson {
  return {
    version: 2,
    projects: {
      'single-app': createProject('single-app', {
        projectType: 'application',
        root: 'apps/single-app',
      }),
    },
  }
}

/**
 * Creates a large workspace for performance testing
 */
export function createLargeWorkspace(projectCount: number = 50): WorkspaceJson {
  return createWorkspaceJson({ projectCount })
}

/**
 * Creates a workspace with known circular dependencies
 */
export function createCircularWorkspace(): WorkspaceJson {
  return createWorkspaceJson({
    projects: {
      'lib-a': { implicitDependencies: ['lib-b'] },
      'lib-b': { implicitDependencies: ['lib-c'] },
      'lib-c': { implicitDependencies: ['lib-a'] },
    },
    includeCircularDeps: false, // We're manually setting the cycle above
  })
}
