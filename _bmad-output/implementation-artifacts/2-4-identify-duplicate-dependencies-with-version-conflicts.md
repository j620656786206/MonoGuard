# Story 2.4: Identify Duplicate Dependencies with Version Conflicts

Status: ready-for-dev

## Story

As a **user**,
I want **to see which dependencies have version conflicts across packages**,
So that **I can resolve version mismatches that may cause issues**.

## Acceptance Criteria

1. **AC1: Duplicate Dependency Detection**
   - Given parsed workspace data from Story 2.1
   - When I analyze for duplicate dependencies
   - Then the analysis identifies all external dependencies used by multiple packages
   - And groups them by dependency name

2. **AC2: Version Conflict Identification**
   - Given dependencies used by multiple packages
   - When versions differ across packages
   - Then a version conflict is reported with:
     - `packageName` - the dependency with conflict (e.g., "lodash")
     - `conflictingVersions` - array of version/packages pairs
     - Each version shows which workspace packages use it
   - Example: lodash@4.17.21 used by pkg-a, lodash@4.17.19 used by pkg-b

3. **AC3: Severity Classification**
   - Given version conflicts
   - When severity is calculated
   - Then classification follows semver differences:
     - `critical` - Major version differences (e.g., v3 vs v4) - breaking changes likely
     - `warning` - Minor version differences (e.g., v4.17 vs v4.18) - new features, possible issues
     - `info` - Patch version differences (e.g., v4.17.21 vs v4.17.19) - bug fixes only
   - And severity is based on the largest version gap in the conflict

4. **AC4: Affected Packages Report**
   - Given a version conflict
   - When results are returned
   - Then each conflict includes:
     - List of all workspace packages affected
     - Which version each package uses
     - Whether the conflict is in dependencies, devDependencies, or peerDependencies

5. **AC5: VersionConflict Type Compliance**
   - Given the analysis output
   - When serialized to JSON
   - Then the output matches the expected type structure:
     ```typescript
     interface VersionConflict {
       packageName: string;
       conflictingVersions: ConflictingVersion[];
       severity: 'critical' | 'warning' | 'info';
       resolution: string;
       impact: string;
     }
     ```
   - And all JSON uses camelCase field names

6. **AC6: Performance Requirements**
   - Given a workspace with 100 packages
   - When conflict detection completes
   - Then it finishes in < 1 second

## Tasks / Subtasks

- [ ] **Task 1: Define VersionConflict Types in Go** (AC: #5)
  - [ ] 1.1 Create `pkg/types/version_conflict.go`:
    ```go
    package types

    // VersionConflict represents a dependency with multiple versions across packages.
    // Matches @monoguard/types VersionConflict.
    type VersionConflict struct {
        PackageName         string               `json:"packageName"`
        ConflictingVersions []*ConflictingVersion `json:"conflictingVersions"`
        Severity            ConflictSeverity     `json:"severity"`
        Resolution          string               `json:"resolution"`
        Impact              string               `json:"impact"`
    }

    // ConflictingVersion represents one version and which packages use it.
    type ConflictingVersion struct {
        Version    string   `json:"version"`
        Packages   []string `json:"packages"`   // Workspace packages using this version
        IsBreaking bool     `json:"isBreaking"` // True if major version differs from others
        DepType    string   `json:"depType"`    // "production", "development", "peer"
    }

    // ConflictSeverity indicates how serious the version mismatch is
    type ConflictSeverity string

    const (
        ConflictSeverityCritical ConflictSeverity = "critical" // Major version difference
        ConflictSeverityWarning  ConflictSeverity = "warning"  // Minor version difference
        ConflictSeverityInfo     ConflictSeverity = "info"     // Patch version difference
    )
    ```
  - [ ] 1.2 Add JSON serialization tests in `pkg/types/version_conflict_test.go`
  - [ ] 1.3 Update old VersionConflict type in types.go if needed

- [ ] **Task 2: Implement Semver Parsing** (AC: #3)
  - [ ] 2.1 Create `pkg/analyzer/semver.go`:
    ```go
    package analyzer

    // SemVer represents a parsed semantic version
    type SemVer struct {
        Major      int
        Minor      int
        Patch      int
        Prerelease string
        Raw        string
    }

    // ParseSemVer parses a version string like "4.17.21" or "^4.17.0"
    func ParseSemVer(version string) (*SemVer, error)

    // StripRange removes semver range prefixes (^, ~, >=, etc.)
    func StripRange(version string) string

    // CompareVersions returns the type of difference between two versions
    func CompareVersions(v1, v2 *SemVer) VersionDifference

    // VersionDifference represents the type of difference
    type VersionDifference int

    const (
        VersionDifferenceNone  VersionDifference = iota // Same version
        VersionDifferencePatch                          // Only patch differs
        VersionDifferenceMinor                          // Minor or patch differs
        VersionDifferenceMajor                          // Major differs
    )
    ```
  - [ ] 2.2 Handle common version formats:
    - Exact: `4.17.21`
    - Caret: `^4.17.0`
    - Tilde: `~4.17.0`
    - Range: `>=4.0.0 <5.0.0`
    - Wildcard: `4.x`, `4.17.x`
  - [ ] 2.3 Create tests in `pkg/analyzer/semver_test.go`

- [ ] **Task 3: Implement Conflict Detector** (AC: #1, #2, #4)
  - [ ] 3.1 Create `pkg/analyzer/conflict_detector.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // ConflictDetector finds version conflicts in workspace dependencies
    type ConflictDetector struct {
        workspace *types.WorkspaceData
    }

    // NewConflictDetector creates a new detector
    func NewConflictDetector(workspace *types.WorkspaceData) *ConflictDetector

    // DetectConflicts finds all version conflicts across packages
    func (cd *ConflictDetector) DetectConflicts() []*types.VersionConflict

    // collectDependencies gathers all dependencies across workspace
    // Returns map[depName]map[version][]packageNames
    func (cd *ConflictDetector) collectDependencies() map[string]map[string][]string

    // buildConflict creates a VersionConflict from collected data
    func (cd *ConflictDetector) buildConflict(
        depName string,
        versionMap map[string][]string,
    ) *types.VersionConflict

    // determineSeverity calculates severity based on version differences
    func determineSeverity(versions []string) types.ConflictSeverity

    // generateResolution suggests how to resolve the conflict
    func generateResolution(conflict *types.VersionConflict) string

    // generateImpact describes the impact of the conflict
    func generateImpact(conflict *types.VersionConflict) string
    ```
  - [ ] 3.2 Implement dependency collection from all packages:
    ```go
    func (cd *ConflictDetector) collectDependencies() map[string]map[string][]string {
        // depName -> version -> []packageNames
        deps := make(map[string]map[string][]string)

        for pkgName, pkg := range cd.workspace.Packages {
            // Collect production dependencies
            for depName, version := range pkg.Dependencies {
                cd.addDependency(deps, depName, version, pkgName)
            }
            // Collect dev dependencies
            for depName, version := range pkg.DevDependencies {
                cd.addDependency(deps, depName, version, pkgName)
            }
            // Collect peer dependencies
            for depName, version := range pkg.PeerDependencies {
                cd.addDependency(deps, depName, version, pkgName)
            }
        }

        return deps
    }
    ```
  - [ ] 3.3 Filter to only dependencies with 2+ different versions
  - [ ] 3.4 Create tests in `pkg/analyzer/conflict_detector_test.go`

- [ ] **Task 4: Implement Resolution Suggestions** (AC: #5)
  - [ ] 4.1 Implement `generateResolution`:
    ```go
    func generateResolution(conflict *types.VersionConflict) string {
        // Find the highest version
        highestVersion := findHighestVersion(conflict.ConflictingVersions)

        switch conflict.Severity {
        case types.ConflictSeverityCritical:
            return fmt.Sprintf(
                "Major version conflict detected. Review breaking changes before upgrading all packages to %s",
                highestVersion,
            )
        case types.ConflictSeverityWarning:
            return fmt.Sprintf(
                "Consider upgrading all packages to %s for consistency",
                highestVersion,
            )
        case types.ConflictSeverityInfo:
            return fmt.Sprintf(
                "Patch version difference. Safe to upgrade all packages to %s",
                highestVersion,
            )
        }
        return ""
    }
    ```
  - [ ] 4.2 Implement `generateImpact` with bundle size and runtime considerations
  - [ ] 4.3 Add tests for resolution and impact generation

- [ ] **Task 5: Wire to Analyzer** (AC: all)
  - [ ] 5.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) Analyze(workspace *types.WorkspaceData) (*types.AnalysisResult, error) {
        // Build dependency graph (Story 2.2)
        graph, err := a.graphBuilder.Build(workspace)
        if err != nil {
            return nil, err
        }

        // Detect circular dependencies (Story 2.3)
        cycleDetector := NewCycleDetector(graph)
        cycles := cycleDetector.DetectCycles()

        // Detect version conflicts (Story 2.4)
        conflictDetector := NewConflictDetector(workspace)
        conflicts := conflictDetector.DetectConflicts()

        return &types.AnalysisResult{
            HealthScore:          100, // Will be calculated in Story 2.5
            Packages:             len(graph.Nodes),
            CircularDependencies: cycles,
            VersionConflicts:     conflicts,
            Graph:                graph,
            CreatedAt:            time.Now().UTC().Format(time.RFC3339),
        }, nil
    }
    ```
  - [ ] 5.2 Update AnalysisResult type to include VersionConflicts field
  - [ ] 5.3 Update handler and WASM tests

- [ ] **Task 6: Performance Testing** (AC: #6)
  - [ ] 6.1 Create `pkg/analyzer/conflict_detector_benchmark_test.go`:
    ```go
    func BenchmarkDetectConflicts100Packages(b *testing.B) {
        workspace := generateWorkspaceWithConflicts(100, 20)
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            detector := NewConflictDetector(workspace)
            detector.DetectConflicts()
        }
    }

    func generateWorkspaceWithConflicts(packageCount, conflictCount int) *types.WorkspaceData {
        // Generate realistic workspace with version conflicts
    }
    ```
  - [ ] 6.2 Verify 100 packages < 1 second

- [ ] **Task 7: Integration Verification** (AC: all)
  - [ ] 7.1 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 7.2 Update smoke test to verify conflict detection
  - [ ] 7.3 Test with known conflict scenarios
  - [ ] 7.4 Verify all tests pass: `make test`

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Conflict detector in `pkg/analyzer/`
- **Input:** WorkspaceData (from Story 2.1, uses Package.Dependencies which are `map[string]string`)
- **Output:** List of VersionConflict matching TypeScript types

**Data Flow:**
```
WorkspaceData.Packages[name].Dependencies (map[string]string)
    ↓
collectDependencies() → map[depName]map[version][]packageNames
    ↓
Filter: only deps with 2+ different versions
    ↓
buildConflict() → VersionConflict[]
```

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **External Only:** Only analyze external dependencies (internal workspace deps don't have version conflicts in the same way)
- **Semver Parsing:** Handle common version formats including ranges

### Critical Don't-Miss Rules

**From project-context.md:**

1. **Version Format Handling:**
   ```go
   // ✅ CORRECT: Strip range prefixes for comparison
   "^4.17.0"  → "4.17.0"
   "~4.17.0"  → "4.17.0"
   ">=4.0.0"  → "4.0.0"

   // ❌ WRONG: Comparing raw version strings with ranges
   "^4.17.0" != "4.17.0" // Would incorrectly report as conflict
   ```

2. **Severity Rules:**
   ```go
   // Major difference: 3.x vs 4.x → critical
   // Minor difference: 4.17.x vs 4.18.x → warning
   // Patch difference: 4.17.19 vs 4.17.21 → info

   func determineSeverity(versions []string) ConflictSeverity {
       parsed := parseAllVersions(versions)
       maxDiff := findMaxDifference(parsed)

       switch maxDiff {
       case VersionDifferenceMajor:
           return ConflictSeverityCritical
       case VersionDifferenceMinor:
           return ConflictSeverityWarning
       default:
           return ConflictSeverityInfo
       }
   }
   ```

3. **Don't Count Internal Packages:**
   ```go
   // ✅ CORRECT: Skip internal workspace packages
   for depName, version := range pkg.Dependencies {
       if isInternalPackage(depName) {
           continue // Internal packages don't cause npm version conflicts
       }
       addDependency(deps, depName, version, pkgName)
   }
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Add conflict detection
│   │   ├── analyzer_test.go                   # UPDATE
│   │   ├── graph_builder.go                   # From Story 2.2
│   │   ├── cycle_detector.go                  # From Story 2.3
│   │   ├── conflict_detector.go               # NEW: Version conflict detection
│   │   ├── conflict_detector_test.go          # NEW: Conflict tests
│   │   ├── conflict_detector_benchmark_test.go # NEW: Performance tests
│   │   ├── semver.go                          # NEW: Semver parsing
│   │   └── semver_test.go                     # NEW: Semver tests
│   └── types/
│       ├── types.go                           # UPDATE: Add VersionConflicts to AnalysisResult
│       ├── version_conflict.go                # NEW: VersionConflict types
│       ├── version_conflict_test.go           # NEW: Type tests
│       ├── circular.go                        # From Story 2.3
│       └── graph.go                           # From Story 2.2
└── ...
```

### Input/Output Format

**Input (WorkspaceData - external deps from Package):**
```json
{
  "packages": {
    "@mono/app": {
      "name": "@mono/app",
      "dependencies": { "lodash": "^4.17.21", "react": "^18.2.0" },
      "devDependencies": { "typescript": "^5.0.0" }
    },
    "@mono/utils": {
      "name": "@mono/utils",
      "dependencies": { "lodash": "^4.17.19" },
      "devDependencies": { "typescript": "^4.9.0" }
    }
  }
}
```

**Output (VersionConflict[]):**
```json
[
  {
    "packageName": "lodash",
    "conflictingVersions": [
      { "version": "^4.17.21", "packages": ["@mono/app"], "isBreaking": false, "depType": "production" },
      { "version": "^4.17.19", "packages": ["@mono/utils"], "isBreaking": false, "depType": "production" }
    ],
    "severity": "info",
    "resolution": "Patch version difference. Safe to upgrade all packages to ^4.17.21",
    "impact": "Minor bundle size increase. No breaking changes expected."
  },
  {
    "packageName": "typescript",
    "conflictingVersions": [
      { "version": "^5.0.0", "packages": ["@mono/app"], "isBreaking": true, "depType": "development" },
      { "version": "^4.9.0", "packages": ["@mono/utils"], "isBreaking": false, "depType": "development" }
    ],
    "severity": "critical",
    "resolution": "Major version conflict detected. Review breaking changes before upgrading all packages to ^5.0.0",
    "impact": "TypeScript 5.x has breaking changes. Check compatibility before upgrading."
  }
]
```

### Test Scenarios

| Scenario | Versions | Expected Severity |
|----------|----------|-------------------|
| No conflict | lodash@4.17.21 everywhere | No conflict reported |
| Patch diff | 4.17.19 vs 4.17.21 | info |
| Minor diff | 4.17.x vs 4.18.x | warning |
| Major diff | 3.x vs 4.x | critical |
| Mixed diff | 3.x, 4.17.x, 4.18.x | critical (worst case) |
| Range formats | ^4.17.0 vs ~4.17.0 | info (same base) |
| Dev vs prod | same dep different sections | Reported with depType |

### Previous Story Intelligence

**From Story 2.1 (ready-for-dev):**
- WorkspaceData.Packages has `Dependencies map[string]string`
- These are ALL dependencies including external ones
- Version is the raw string from package.json (may include ^, ~, etc.)

**From Story 2.2 (ready-for-dev):**
- PackageNode separates internal deps ([]string) from external (map[string]string)
- `ExternalDeps` and `ExternalDevDeps` contain external dependencies with versions

**Key Data Source:**
- Use `WorkspaceData.Packages[name].Dependencies` for analysis
- This includes both internal and external, filter to external only
- OR use `graph.Nodes[name].ExternalDeps` which is pre-filtered

### Semver Parsing Notes

**Common Version Formats:**
| Format | Example | Meaning |
|--------|---------|---------|
| Exact | `4.17.21` | Exactly this version |
| Caret | `^4.17.0` | Compatible with 4.x.x |
| Tilde | `~4.17.0` | Compatible with 4.17.x |
| Range | `>=4.0.0 <5.0.0` | Between versions |
| Wildcard | `4.x` | Any 4.x.x |
| Latest | `latest` | Latest version |
| Tag | `next`, `beta` | npm dist-tag |

**Parsing Strategy:**
1. Strip range prefix (^, ~, >=, etc.)
2. Parse major.minor.patch
3. Handle prerelease (-alpha, -beta)
4. For comparison, use normalized version

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 2.4]
- [Source: packages/types/src/domain.ts#VersionConflict]
- [Source: _bmad-output/implementation-artifacts/2-1-implement-workspace-configuration-parser.md]
- [Semantic Versioning Spec](https://semver.org/)
- [npm Semver Package](https://github.com/npm/node-semver)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List
