# Package.json Parser Module Technical Specification

## Overview

The Package.json Parser Module is responsible for parsing and extracting dependency information from package.json files across different monorepo workspace formats (npm, yarn, pnpm). This module serves as the foundation for all dependency analysis operations in MonoGuard Phase 1.

## Technical Requirements

### REQ-PKG-001: Multi-Format Workspace Support

**Description:** The parser shall support all major monorepo workspace formats with their specific configuration patterns.

**Supported Formats:**
- **npm workspaces**: `workspaces` field in root package.json
- **yarn workspaces**: `workspaces` field with yarn-specific features
- **pnpm workspaces**: `pnpm-workspace.yaml` file configuration
- **Lerna**: `lerna.json` configuration (legacy support)
- **Nx**: `workspace.json` and `nx.json` configuration

**Technical Specifications:**
```go
type WorkspaceType string

const (
    WorkspaceTypeNPM   WorkspaceType = "npm"
    WorkspaceTypeYarn  WorkspaceType = "yarn"
    WorkspaceTypePnpm  WorkspaceType = "pnpm"
    WorkspaceTypeLerna WorkspaceType = "lerna"
    WorkspaceTypeNx    WorkspaceType = "nx"
)

type WorkspaceConfig struct {
    Type        WorkspaceType `json:"type"`
    RootPath    string        `json:"root_path"`
    Packages    []string      `json:"packages"`        // Glob patterns
    PackageDirs []string      `json:"package_dirs"`    // Resolved directories
    Metadata    interface{}   `json:"metadata"`        // Type-specific metadata
}
```

### REQ-PKG-002: Package.json Parsing Engine

**Description:** Core parsing engine with comprehensive field extraction and validation.

**Data Structures:**
```go
type PackageInfo struct {
    // Basic Package Information
    Name            string            `json:"name"`
    Version         string            `json:"version"`
    Description     string            `json:"description"`
    Private         bool              `json:"private"`
    
    // File System Information
    Path            string            `json:"path"`             // Absolute path to package.json
    RelativePath    string            `json:"relative_path"`    // Relative to workspace root
    
    // Dependencies
    Dependencies         map[string]string `json:"dependencies"`
    DevDependencies      map[string]string `json:"dev_dependencies"`
    PeerDependencies     map[string]string `json:"peer_dependencies"`
    OptionalDependencies map[string]string `json:"optional_dependencies"`
    BundledDependencies  []string          `json:"bundled_dependencies"`
    
    // Workspace-specific fields
    Workspaces      []string          `json:"workspaces"`       // For workspace roots
    
    // Scripts and Configuration
    Scripts         map[string]string `json:"scripts"`
    Main            string            `json:"main"`
    Module          string            `json:"module"`
    Types           string            `json:"types"`
    Exports         interface{}       `json:"exports"`
    
    // Build and Tooling Configuration
    SideEffects     interface{}       `json:"side_effects"`     // boolean or string array
    Engines         map[string]string `json:"engines"`
    
    // MonoGuard-specific metadata
    ParsedAt        time.Time         `json:"parsed_at"`
    Checksum        string            `json:"checksum"`         // MD5 hash of file content
    Errors          []ParseError      `json:"errors"`
    Warnings        []ParseWarning    `json:"warnings"`
}

type ParseError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}

type ParseWarning struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}
```

### REQ-PKG-003: Version Range Parsing

**Description:** Parse and normalize version ranges according to semantic versioning and package manager specific formats.

**Version Range Types:**
```go
type VersionRange struct {
    Raw         string      `json:"raw"`          // Original string from package.json
    Normalized  string      `json:"normalized"`   // Normalized format
    Type        RangeType   `json:"type"`         // Type of version range
    Lower       *Version    `json:"lower"`        // Lower bound (if applicable)
    Upper       *Version    `json:"upper"`        // Upper bound (if applicable)
    Inclusive   bool        `json:"inclusive"`    // Whether bounds are inclusive
}

type RangeType string

const (
    RangeTypeExact     RangeType = "exact"      // 1.0.0
    RangeTypeCaret     RangeType = "caret"      // ^1.0.0
    RangeTypeTilde     RangeType = "tilde"      // ~1.0.0
    RangeTypeGTE       RangeType = "gte"        // >=1.0.0
    RangeTypeLTE       RangeType = "lte"        // <=1.0.0
    RangeTypeRange     RangeType = "range"      // 1.0.0 - 2.0.0
    RangeTypeWildcard  RangeType = "wildcard"   // *
    RangeTypeLatest    RangeType = "latest"     // latest
    RangeTypeFile      RangeType = "file"       // file:../local-pkg
    RangeTypeGit       RangeType = "git"        // git+https://...
    RangeTypeWorkspace RangeType = "workspace"  // workspace:*
)

type Version struct {
    Major int    `json:"major"`
    Minor int    `json:"minor"`
    Patch int    `json:"patch"`
    Pre   string `json:"pre"`    // Pre-release identifier
    Build string `json:"build"`  // Build metadata
}
```

### REQ-PKG-004: Workspace Discovery Algorithm

**Description:** Implement efficient workspace discovery with caching and incremental updates.

**Algorithm Specification:**
```go
type WorkspaceDiscovery interface {
    // Discover workspaces from a root directory
    DiscoverWorkspaces(rootPath string) (*WorkspaceConfig, error)
    
    // Validate workspace configuration
    ValidateWorkspace(config *WorkspaceConfig) []ValidationError
    
    // Get all package.json files in workspace
    FindPackageFiles(config *WorkspaceConfig) ([]string, error)
    
    // Check if workspace configuration has changed
    HasChanged(config *WorkspaceConfig, lastChecksum string) (bool, error)
}

type ValidationError struct {
    Path    string `json:"path"`
    Message string `json:"message"`
    Level   string `json:"level"` // error, warning, info
}
```

**Discovery Priority Order:**
1. `pnpm-workspace.yaml` (highest priority)
2. `package.json` with `workspaces` field
3. `lerna.json` configuration
4. `nx.json` and `workspace.json`
5. Fallback to single package detection

## API Interfaces

### REQ-PKG-005: Parser Interface

**Primary Interface:**
```go
type PackageParser interface {
    // Parse a single package.json file
    ParsePackageFile(filePath string) (*PackageInfo, error)
    
    // Parse all packages in a workspace
    ParseWorkspace(rootPath string) (*WorkspaceAnalysis, error)
    
    // Parse with options for customization
    ParseWithOptions(rootPath string, options ParseOptions) (*WorkspaceAnalysis, error)
    
    // Validate parsed package information
    ValidatePackage(pkg *PackageInfo) []ValidationError
    
    // Get parser statistics and metrics
    GetStats() ParseStats
}

type ParseOptions struct {
    // Workspace configuration
    WorkspaceType       *WorkspaceType `json:"workspace_type"`    // Auto-detect if nil
    IncludeDevDeps      bool           `json:"include_dev_deps"`
    IncludePeerDeps     bool           `json:"include_peer_deps"`
    IncludeOptionalDeps bool           `json:"include_optional_deps"`
    
    // Performance options
    MaxConcurrency      int            `json:"max_concurrency"`   // Default: runtime.NumCPU()
    EnableCaching       bool           `json:"enable_caching"`
    CacheDir           string         `json:"cache_dir"`
    
    // Filtering options
    IncludePatterns    []string       `json:"include_patterns"`  // Glob patterns
    ExcludePatterns    []string       `json:"exclude_patterns"`  // Glob patterns
    MaxDepth           int            `json:"max_depth"`         // Max directory depth
    
    // Error handling
    ContinueOnError    bool           `json:"continue_on_error"`
    MaxErrors          int            `json:"max_errors"`        // Stop after N errors
}

type WorkspaceAnalysis struct {
    Config    WorkspaceConfig  `json:"config"`
    Packages  []*PackageInfo   `json:"packages"`
    Stats     ParseStats       `json:"stats"`
    Errors    []ParseError     `json:"errors"`
    Warnings  []ParseWarning   `json:"warnings"`
}

type ParseStats struct {
    TotalPackages     int           `json:"total_packages"`
    ParsedPackages    int           `json:"parsed_packages"`
    FailedPackages    int           `json:"failed_packages"`
    TotalDependencies int           `json:"total_dependencies"`
    UniquePackages    int           `json:"unique_packages"`
    ParseDuration     time.Duration `json:"parse_duration"`
    CacheHits         int           `json:"cache_hits"`
    CacheMisses       int           `json:"cache_misses"`
}
```

### REQ-PKG-006: Version Range Parser Interface

```go
type VersionRangeParser interface {
    // Parse a version range string
    ParseRange(rangeStr string) (*VersionRange, error)
    
    // Parse multiple version ranges (for conflicts)
    ParseRanges(ranges []string) ([]*VersionRange, error)
    
    // Check if a version satisfies a range
    Satisfies(version *Version, vrange *VersionRange) bool
    
    // Find intersection of version ranges
    Intersect(ranges []*VersionRange) (*VersionRange, error)
    
    // Compare version ranges for conflict detection
    AreCompatible(range1, range2 *VersionRange) bool
}
```

## Algorithm Specifications

### REQ-PKG-007: Workspace Discovery Algorithm

**Algorithm:** Multi-stage workspace detection with priority-based resolution

```
ALGORITHM: DiscoverWorkspace(rootPath)
INPUT: rootPath (string) - Path to potential workspace root
OUTPUT: WorkspaceConfig - Discovered workspace configuration

1. INITIALIZE candidate configurations array
2. SET priority_order = [pnpm, yarn/npm, lerna, nx]

3. FOR each workspace_type in priority_order:
   a. TRY detect configuration for workspace_type
   b. IF configuration found:
      - VALIDATE configuration syntax
      - RESOLVE package glob patterns
      - ADD to candidates with priority score
   
4. IF candidates is empty:
   a. FALLBACK to single package detection
   b. RETURN single package workspace config

5. SELECT highest priority valid configuration
6. RESOLVE all package directories using glob patterns
7. VALIDATE resolved directories contain package.json files
8. RETURN final WorkspaceConfig
```

**Complexity:** O(n) where n is the number of directories in workspace

### REQ-PKG-008: Package.json Parsing Algorithm

**Algorithm:** Concurrent parsing with error resilience

```
ALGORITHM: ParseWorkspace(config)
INPUT: config (WorkspaceConfig) - Validated workspace configuration
OUTPUT: WorkspaceAnalysis - Complete analysis results

1. INITIALIZE worker pool with MaxConcurrency workers
2. CREATE job queue with package.json file paths
3. INITIALIZE result collectors (packages, errors, warnings)

4. FOR each worker in parallel:
   a. WHILE jobs available:
      - DEQUEUE file path
      - CALL ParsePackageFile(path)
      - IF parsing successful:
        * ADD package to results
        * UPDATE statistics
      - ELSE:
        * ADD error to error collection
        * IF ContinueOnError is false: TERMINATE
      
5. WAIT for all workers to complete
6. AGGREGATE results and compute final statistics
7. VALIDATE workspace consistency (no duplicate names, etc.)
8. RETURN WorkspaceAnalysis
```

**Error Recovery:**
- Individual package parsing failures don't stop workspace analysis
- Malformed JSON handled gracefully with detailed error messages
- Missing required fields result in warnings, not errors
- File permission issues logged with specific remediation steps

### REQ-PKG-009: Version Range Normalization Algorithm

**Algorithm:** Semantic version range parsing with npm/yarn/pnpm compatibility

```
ALGORITHM: ParseVersionRange(rangeString)
INPUT: rangeString (string) - Raw version range from package.json
OUTPUT: VersionRange - Normalized version range object

1. TRIM whitespace and normalize input
2. IDENTIFY range type using regex patterns:
   - Match exact version: /^\d+\.\d+\.\d+$/
   - Match caret range: /^\^[\d\.]+/
   - Match tilde range: /^~[\d\.]+/
   - Match comparison: /^[><=]+[\d\.]+/
   - Match range: /^[\d\.]+ - [\d\.]+/
   - Match special: workspace:, file:, git+, latest, *

3. SWITCH on detected range type:
   a. CASE exact: PARSE single version
   b. CASE caret: COMPUTE compatible range bounds
   c. CASE tilde: COMPUTE patch-level range bounds
   d. CASE comparison: PARSE operator and version
   e. CASE range: PARSE lower and upper bounds
   f. CASE special: HANDLE type-specific parsing

4. VALIDATE parsed components:
   - Check semantic version format
   - Verify range bounds are logical
   - Ensure compatibility with package manager

5. RETURN normalized VersionRange object
```

**Caret Range Logic (^1.2.3):**
- Allow changes that do not modify the major version
- Range: [1.2.3, 2.0.0)

**Tilde Range Logic (~1.2.3):**
- Allow patch-level changes if minor version is specified
- Range: [1.2.3, 1.3.0)

## Data Structures

### REQ-PKG-010: Core Data Models

**Memory-Efficient Structures:**
```go
// Optimize for memory usage in large monorepos
type CompactPackageInfo struct {
    // Interned strings to reduce memory usage
    NameID          uint32            `json:"name_id"`          // String interning
    Version         string            `json:"version"`
    Path            string            `json:"path"`
    
    // Compressed dependency maps
    Dependencies    CompressedDepMap  `json:"dependencies"`
    DevDependencies CompressedDepMap  `json:"dev_dependencies"`
    
    // Flags for boolean fields (pack into single uint32)
    Flags           uint32            `json:"flags"`
    
    // Checksum for change detection
    Checksum        [16]byte          `json:"checksum"`         // MD5 hash
}

type CompressedDepMap struct {
    // Use string interning and compact representation
    Names    []uint32 `json:"names"`     // Interned package name IDs
    Versions []uint32 `json:"versions"`  // Interned version string IDs
}

type StringInterningTable struct {
    strings   []string
    lookup    map[string]uint32
    nextID    uint32
    mutex     sync.RWMutex
}
```

**Workspace Cache Structure:**
```go
type WorkspaceCache struct {
    // Cache configuration
    CacheDir     string        `json:"cache_dir"`
    TTL          time.Duration `json:"ttl"`
    MaxSize      int64         `json:"max_size"`     // Max cache size in bytes
    
    // Cache entries
    entries      map[string]*CacheEntry
    access       *lru.Cache                          // LRU eviction
    mutex        sync.RWMutex
    
    // Statistics
    hits         int64
    misses       int64
    evictions    int64
}

type CacheEntry struct {
    Key         string                `json:"key"`          // Hash of workspace config
    Data        *WorkspaceAnalysis    `json:"data"`
    CreatedAt   time.Time            `json:"created_at"`
    AccessedAt  time.Time            `json:"accessed_at"`
    FileHashes  map[string]string    `json:"file_hashes"`  // For invalidation
    Size        int64                `json:"size"`         // Memory usage
}
```

## Error Handling

### REQ-PKG-011: Error Classification

**Error Types and Handling:**
```go
type ParseErrorType string

const (
    ErrorTypeFileNotFound     ParseErrorType = "file_not_found"
    ErrorTypePermissionDenied ParseErrorType = "permission_denied"
    ErrorTypeInvalidJSON      ParseErrorType = "invalid_json"
    ErrorTypeMissingField     ParseErrorType = "missing_field"
    ErrorTypeInvalidVersion   ParseErrorType = "invalid_version"
    ErrorTypeCircularDep      ParseErrorType = "circular_dependency"
    ErrorTypeWorkspaceConfig  ParseErrorType = "workspace_config"
    ErrorTypeMemoryLimit      ParseErrorType = "memory_limit"
    ErrorTypeTimeout          ParseErrorType = "timeout"
)

type DetailedError struct {
    Type        ParseErrorType `json:"type"`
    Message     string         `json:"message"`
    Path        string         `json:"path"`
    Field       string         `json:"field"`
    Value       interface{}    `json:"value"`
    Suggestions []string       `json:"suggestions"`
    Code        int            `json:"code"`
    Recoverable bool           `json:"recoverable"`
}
```

**Error Recovery Strategies:**
1. **File Access Errors**: Graceful skipping with detailed logging
2. **JSON Parse Errors**: Attempt partial parsing and report specific issues
3. **Version Format Errors**: Use fallback parsing with warnings
4. **Memory Limit Errors**: Implement streaming parsing for large files
5. **Timeout Errors**: Allow partial results with completion status

### REQ-PKG-012: Graceful Degradation

**Degradation Levels:**
```go
type DegradationLevel string

const (
    DegradationNone    DegradationLevel = "none"      // Full functionality
    DegradationPartial DegradationLevel = "partial"   // Some features disabled
    DegradationMinimal DegradationLevel = "minimal"   // Basic parsing only
    DegradationFailed  DegradationLevel = "failed"    // Unable to process
)

type DegradationStrategy struct {
    MemoryThreshold     int64             `json:"memory_threshold"`      // Bytes
    TimeoutThreshold    time.Duration     `json:"timeout_threshold"`
    ErrorThreshold      int               `json:"error_threshold"`
    FallbackBehavior    map[string]string `json:"fallback_behavior"`
}
```

## Testing Requirements

### REQ-PKG-013: Unit Test Coverage

**Test Categories:**
1. **Parser Core Logic**: 95% code coverage required
2. **Version Range Parsing**: Comprehensive semver test cases
3. **Workspace Detection**: All supported workspace formats
4. **Error Handling**: All error paths and recovery scenarios
5. **Performance**: Memory usage and parsing speed benchmarks

**Test Data Sets:**
```go
// Comprehensive test fixtures
var TestWorkspaces = []struct {
    Name        string
    Path        string
    Type        WorkspaceType
    Expected    int  // Expected package count
    ShouldError bool
}{
    {"npm-basic", "fixtures/npm-basic", WorkspaceTypeNPM, 5, false},
    {"yarn-complex", "fixtures/yarn-complex", WorkspaceTypeYarn, 25, false},
    {"pnpm-nested", "fixtures/pnpm-nested", WorkspaceTypePnpm, 12, false},
    {"mixed-invalid", "fixtures/mixed-invalid", WorkspaceTypeNPM, 0, true},
    {"large-monorepo", "fixtures/large-monorepo", WorkspaceTypeNx, 100, false},
}
```

### REQ-PKG-014: Integration Test Requirements

**Integration Test Scenarios:**
1. **Real Monorepo Testing**: Test against actual open-source monorepos
2. **Performance Benchmarks**: Test with varying workspace sizes (10, 50, 100, 500+ packages)
3. **Concurrent Access**: Test parser under concurrent load
4. **Cache Behavior**: Test cache hit/miss scenarios and invalidation
5. **Error Resilience**: Test behavior with partially corrupted workspaces

**Benchmark Requirements:**
- Parse 100 packages in <2 seconds
- Memory usage <500MB for 100 packages
- Cache hit ratio >80% for repeated parsing
- Error recovery success rate >95%

### REQ-PKG-015: Performance Test Specifications

**Performance Benchmarks:**
```go
func BenchmarkPackageParser(b *testing.B) {
    testCases := []struct {
        name     string
        packages int
    }{
        {"small", 10},
        {"medium", 50},
        {"large", 100},
        {"xlarge", 500},
    }
    
    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            // Benchmark parsing performance
        })
    }
}
```

**Memory Profiling Requirements:**
- Profile memory allocation patterns
- Identify memory leaks in long-running scenarios  
- Test garbage collection behavior under load
- Monitor memory fragmentation with large workspaces

## Dependencies and Integration Points

### REQ-PKG-016: External Dependencies

**Required Libraries:**
```go
// Go standard library
import (
    "encoding/json"
    "path/filepath" 
    "sync"
    "time"
)

// Third-party libraries
import (
    "github.com/gobwas/glob"         // Glob pattern matching
    "github.com/hashicorp/golang-lru" // LRU cache implementation
    "gopkg.in/yaml.v3"               // YAML parsing for pnpm-workspace
    "github.com/Masterminds/semver/v3" // Semantic version parsing
)
```

**File System Dependencies:**
- Read access to workspace directories
- Temporary directory for caching (configurable)
- Watch file system for changes (optional)

### REQ-PKG-017: Integration Points

**Integration with Other Phase 1 Components:**

1. **Dependency Tree Resolver:**
   ```go
   // Package parser provides input to tree resolver
   type ParserIntegration interface {
       GetWorkspacePackages() []*PackageInfo
       GetPackageDependencies(pkg *PackageInfo) map[string]*VersionRange
       ResolvePackagePath(name string) (string, error)
   }
   ```

2. **Duplicate Dependency Detector:**
   ```go
   // Parser provides structured dependency data
   type DependencyData interface {
       GetAllDependencies() map[string][]*PackageRef
       GetPackageVersions(packageName string) []string
       GetDependencyGraph() *DependencyGraph
   }
   ```

3. **Unused Dependency Detector:**
   ```go
   // Parser provides declared dependencies for analysis
   type UnusedDepAnalysis interface {
       GetDeclaredDependencies(pkg *PackageInfo) []string
       GetPackageSourceFiles(pkg *PackageInfo) ([]string, error)
       IsDevDependency(pkg, dep string) bool
   }
   ```

## Success Metrics

### REQ-PKG-018: Performance Metrics

**Target Metrics:**
- **Parse Speed**: <2 seconds for 100 packages
- **Memory Usage**: <4GB for 500 packages  
- **Accuracy**: >99% for well-formed package.json files
- **Error Recovery**: >95% success rate for partially corrupted workspaces
- **Cache Efficiency**: >80% hit rate for repeated parsing

**Monitoring and Alerting:**
```go
type ParserMetrics struct {
    ParseDuration    prometheus.Histogram
    MemoryUsage      prometheus.Gauge
    ErrorRate        prometheus.Counter
    CacheHitRate     prometheus.Gauge
    PackagesProcessed prometheus.Counter
}
```

### REQ-PKG-019: Quality Metrics

**Code Quality Targets:**
- **Test Coverage**: >95% for core parsing logic
- **Cyclomatic Complexity**: <10 for individual functions
- **Documentation**: 100% public API documentation
- **Static Analysis**: Zero critical issues from linters

This specification provides comprehensive technical details for implementing the Package.json Parser Module as the foundation for MonoGuard's Phase 1 Core Engine.