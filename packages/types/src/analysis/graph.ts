/**
 * DependencyGraph - Core graph data structure
 *
 * Matches Go: pkg/types/graph.go
 * All field names use camelCase for JSON serialization compatibility.
 */
export interface DependencyGraph {
  /** Map of package name to Package details */
  nodes: Record<string, PackageNode>;
  /** List of dependency edges */
  edges: DependencyEdge[];
  /** Workspace root path */
  rootPath: string;
  /** Workspace type detected */
  workspaceType: WorkspaceType;
}

/**
 * PackageNode - Information about a single package in the graph
 *
 * Matches Go: pkg/types/package.go
 */
export interface PackageNode {
  /** Package name (e.g., "@monoguard/types") */
  name: string;
  /** Package version */
  version: string;
  /** Relative path from workspace root */
  path: string;
  /** Direct dependencies */
  dependencies: string[];
  /** Dev dependencies */
  devDependencies: string[];
  /** Peer dependencies */
  peerDependencies: string[];
}

/**
 * DependencyEdge - Represents a dependency relationship between packages
 *
 * Matches Go: pkg/types/edge.go
 */
export interface DependencyEdge {
  /** Source package name */
  from: string;
  /** Target package name */
  to: string;
  /** Dependency type */
  type: DependencyType;
  /** Version range specified */
  versionRange: string;
}

/**
 * DependencyType - Classification of dependency relationship
 *
 * Matches Go: pkg/types/dependency_type.go
 */
export type DependencyType = 'production' | 'development' | 'peer' | 'optional';

/**
 * WorkspaceType - Type of monorepo workspace detected
 *
 * Matches Go: pkg/types/workspace_type.go
 */
export type WorkspaceType = 'npm' | 'yarn' | 'pnpm' | 'nx' | 'unknown';
