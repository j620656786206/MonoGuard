import { describe, it, expect } from 'vitest';
import type {
  AnalysisResult,
  DependencyGraph,
  CircularDependencyInfo,
  CheckResult,
  PackageNode,
  DependencyEdge,
} from '../analysis';

describe('Analysis types', () => {
  describe('AnalysisResult', () => {
    it('can be instantiated with valid data', () => {
      const graph: DependencyGraph = {
        nodes: {
          '@monoguard/types': {
            name: '@monoguard/types',
            version: '0.1.0',
            path: 'packages/types',
            dependencies: [],
            devDependencies: ['typescript'],
            peerDependencies: [],
          },
        },
        edges: [],
        rootPath: '/workspace',
        workspaceType: 'pnpm',
      };

      const result: AnalysisResult = {
        healthScore: 85,
        packageCount: 10,
        circularDependencies: [],
        graph,
        metadata: {
          version: '0.1.0',
          durationMs: 1500,
          filesProcessed: 50,
          workspaceType: 'pnpm',
        },
        createdAt: '2026-01-15T10:30:00Z',
      };

      expect(result.healthScore).toBe(85);
      expect(result.graph.workspaceType).toBe('pnpm');
      expect(result.metadata.durationMs).toBe(1500);
    });
  });

  describe('DependencyGraph', () => {
    it('supports Record<string, PackageNode> for nodes', () => {
      const graph: DependencyGraph = {
        nodes: {
          'package-a': {
            name: 'package-a',
            version: '1.0.0',
            path: 'packages/a',
            dependencies: ['package-b'],
            devDependencies: [],
            peerDependencies: [],
          },
          'package-b': {
            name: 'package-b',
            version: '1.0.0',
            path: 'packages/b',
            dependencies: [],
            devDependencies: [],
            peerDependencies: [],
          },
        },
        edges: [
          {
            from: 'package-a',
            to: 'package-b',
            type: 'production',
            versionRange: '^1.0.0',
          },
        ],
        rootPath: '/workspace',
        workspaceType: 'npm',
      };

      expect(Object.keys(graph.nodes)).toHaveLength(2);
      expect(graph.edges).toHaveLength(1);
    });
  });

  describe('CircularDependencyInfo', () => {
    it('can represent direct circular dependency', () => {
      const circular: CircularDependencyInfo = {
        cycle: ['package-a', 'package-b', 'package-a'],
        type: 'direct',
        severity: 'critical',
        impact: 'Build failure due to circular dependency',
        complexity: 5,
      };

      expect(circular.type).toBe('direct');
      expect(circular.severity).toBe('critical');
    });

    it('can include fix strategy', () => {
      const circular: CircularDependencyInfo = {
        cycle: ['package-a', 'package-b', 'package-c', 'package-a'],
        type: 'indirect',
        severity: 'warning',
        impact: 'Potential build issues',
        complexity: 7,
        fixStrategy: {
          type: 'extract_module',
          description: 'Extract shared code into a new package',
          steps: [
            'Create packages/shared',
            'Move common code to shared',
            'Update imports in package-a and package-c',
          ],
          affectedFiles: ['packages/a/src/index.ts', 'packages/c/src/index.ts'],
        },
      };

      expect(circular.fixStrategy?.type).toBe('extract_module');
      expect(circular.fixStrategy?.steps).toHaveLength(3);
    });
  });

  describe('CheckResult', () => {
    it('can represent passing check', () => {
      const checkResult: CheckResult = {
        passed: true,
        errors: [],
        warnings: [],
        healthScore: 95,
      };

      expect(checkResult.passed).toBe(true);
      expect(checkResult.errors).toHaveLength(0);
    });

    it('can represent failing check with errors', () => {
      const checkResult: CheckResult = {
        passed: false,
        errors: [
          {
            code: 'CIRCULAR_DETECTED',
            message: 'Circular dependency found: A -> B -> A',
            file: 'packages/a/package.json',
          },
        ],
        warnings: [
          {
            code: 'LOW_HEALTH_SCORE',
            message: 'Health score below threshold',
          },
        ],
        healthScore: 45,
      };

      expect(checkResult.passed).toBe(false);
      expect(checkResult.errors).toHaveLength(1);
      expect(checkResult.warnings).toHaveLength(1);
    });
  });

  describe('PackageNode', () => {
    it('correctly types all dependency categories', () => {
      const node: PackageNode = {
        name: '@monoguard/web',
        version: '1.0.0',
        path: 'apps/web',
        dependencies: ['react', 'react-dom'],
        devDependencies: ['typescript', 'vitest'],
        peerDependencies: ['react'],
      };

      expect(node.dependencies).toContain('react');
      expect(node.devDependencies).toContain('typescript');
      expect(node.peerDependencies).toContain('react');
    });
  });

  describe('DependencyEdge', () => {
    it('supports all dependency types', () => {
      const edges: DependencyEdge[] = [
        { from: 'a', to: 'b', type: 'production', versionRange: '^1.0.0' },
        { from: 'a', to: 'c', type: 'development', versionRange: '*' },
        { from: 'a', to: 'd', type: 'peer', versionRange: '>=16.0.0' },
        { from: 'a', to: 'e', type: 'optional', versionRange: '~2.0.0' },
      ];

      expect(edges[0].type).toBe('production');
      expect(edges[1].type).toBe('development');
      expect(edges[2].type).toBe('peer');
      expect(edges[3].type).toBe('optional');
    });
  });
});

describe('WorkspaceType', () => {
  it('supports all workspace types', () => {
    const graph1: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'npm',
    };
    const graph2: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'yarn',
    };
    const graph3: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'pnpm',
    };
    const graph4: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'nx',
    };
    const graph5: DependencyGraph = {
      nodes: {},
      edges: [],
      rootPath: '/',
      workspaceType: 'unknown',
    };

    expect(graph1.workspaceType).toBe('npm');
    expect(graph2.workspaceType).toBe('yarn');
    expect(graph3.workspaceType).toBe('pnpm');
    expect(graph4.workspaceType).toBe('nx');
    expect(graph5.workspaceType).toBe('unknown');
  });
});
