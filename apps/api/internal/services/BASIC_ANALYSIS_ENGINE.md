# Basic Analysis Engine

The Basic Analysis Engine is a focused, single-purpose analysis system for mono-guard that provides core dependency analysis functionality.

## Components

### 1. Package.json Parser (`MonoGuardPackageParser`)

**Purpose**: Discovers and parses package.json files in a monorepo.

**Features**:
- Recursively discovers all package.json files
- Skips common directories (node_modules, .git, dist, build, etc.)
- Extracts dependencies, devDependencies, and peerDependencies
- Handles malformed JSON gracefully

**Usage**:
```go
parser := NewMonoGuardPackageParser(logger)
packages, err := parser.ParseRepository(ctx, "/path/to/repo")
```

### 2. Duplicate Dependency Detector (`MonoGuardDuplicateDetector`)

**Purpose**: Identifies duplicate dependencies across packages.

**Features**:
- Detects dependencies with multiple versions
- Calculates risk levels based on number of versions
- Estimates resource waste (disk space, bundle size)
- Generates consolidation recommendations
- Provides step-by-step migration plans

**Risk Levels**:
- **Low**: 1 version (no duplicates)
- **Medium**: 2 versions
- **High**: 3 versions
- **Critical**: 4+ versions

**Usage**:
```go
detector := NewMonoGuardDuplicateDetector(logger)
duplicates, err := detector.FindDuplicates(packages)
```

### 3. Version Conflict Analyzer (`MonoGuardConflictAnalyzer`)

**Purpose**: Analyzes version conflicts between dependencies.

**Features**:
- Identifies incompatible major versions
- Assesses breaking changes impact
- Calculates conflict risk levels
- Generates resolution strategies
- Provides impact assessments

**Conflict Detection**:
- Focuses on production dependencies
- Identifies major version conflicts (e.g., React 17 vs 18)
- Ignores minor/patch version differences for conflicts

**Usage**:
```go
analyzer := NewMonoGuardConflictAnalyzer(logger)
conflicts, err := analyzer.FindConflicts(packages)
```

### 4. Analysis Report Generator (`MonoGuardReportGenerator`)

**Purpose**: Combines analysis results into comprehensive reports.

**Features**:
- Generates bundle impact analysis
- Calculates health scores
- Creates dependency usage breakdowns
- Provides summary statistics

**Health Score Calculation**:
- Starts at 100 points
- Deducts 5 points per duplicate dependency
- Deducts 10 points per version conflict
- Minimum score: 0

**Usage**:
```go
generator := NewMonoGuardReportGenerator(logger)
results, err := generator.GenerateReport(packages, duplicates, conflicts)
```

## Complete Analysis Engine

### BasicAnalysisEngine

**Purpose**: Orchestrates all analysis components to perform complete repository analysis.

**Features**:
- Integrated workflow from parsing to reporting
- Structured error handling
- Performance logging
- Comprehensive results

**Usage**:
```go
engine := NewBasicAnalysisEngine(logger)
results, err := engine.AnalyzeRepository(ctx, "/path/to/repo", "project-id")
```

**Analysis Results** (`models.DependencyAnalysisResults`):
- `DuplicateDependencies`: List of duplicate dependency issues
- `VersionConflicts`: List of version conflict issues
- `BundleImpact`: Bundle size and waste analysis
- `Summary`: Overall analysis summary with health score

## Example Usage

See `basic_analysis_engine_example.go` for complete usage examples:

### Quick Analysis
```go
summary, err := QuickAnalysisExample("/path/to/repo")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Health Score: %.1f\n", summary.HealthScore)
fmt.Printf("Issues Found: %d\n", summary.IssueCount)
```

### Full Analysis
```go
engine := NewBasicAnalysisEngine(logger)
results, err := engine.AnalyzeRepository(ctx, repoPath, projectID)
if err != nil {
    log.Fatal(err)
}

// Access detailed results
for _, duplicate := range results.DuplicateDependencies {
    fmt.Printf("Duplicate: %s with versions %v\n", 
        duplicate.PackageName, duplicate.Versions)
}
```

## Testing

The engine includes comprehensive tests:

```bash
go test ./internal/services -v -run "TestMonoGuard|TestBasicAnalysisEngine"
```

Tests cover:
- Package parsing functionality
- Duplicate detection accuracy
- Version conflict identification
- End-to-end analysis workflow

## Design Principles

1. **Surgical Precision**: Each component has a single, focused responsibility
2. **Minimal Dependencies**: Uses only essential external libraries
3. **Error Resilience**: Graceful handling of malformed files and edge cases
4. **Performance**: Efficient file walking and parsing
5. **Extensibility**: Clean interfaces for future enhancements

## Integration

The Basic Analysis Engine integrates seamlessly with the existing mono-guard infrastructure:

- Uses standard `models.DependencyAnalysisResults` for output
- Compatible with existing API handlers
- Follows established logging patterns
- Supports context-based cancellation

## Limitations

This is a **basic** implementation focused on core functionality:

- **No unused dependency detection** (requires static analysis)
- **No circular dependency detection** (requires dependency graph analysis)
- **Simplified version comparison** (uses lexical sorting instead of semver)
- **Basic bundle size estimation** (uses fixed estimates)

For advanced features, consider using the full `DependencyAnalyzer` or `PackageJSONParser` services.