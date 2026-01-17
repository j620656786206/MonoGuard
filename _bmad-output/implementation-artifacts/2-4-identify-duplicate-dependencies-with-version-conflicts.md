# Story 2.4: Identify Duplicate Dependencies with Version Conflicts

Status: review

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

- [x] **Task 1: Define VersionConflict Types in Go** (AC: #5)
  - [x] 1.1 Create `pkg/types/version_conflict.go`
  - [x] 1.2 Add JSON serialization tests in `pkg/types/version_conflict_test.go`
  - [x] 1.3 Update old VersionConflict type in types.go (marked as deprecated)

- [x] **Task 2: Implement Semver Parsing** (AC: #3)
  - [x] 2.1 Create `pkg/analyzer/semver.go`
  - [x] 2.2 Handle common version formats (Exact, Caret, Tilde, Range, Wildcard)
  - [x] 2.3 Create tests in `pkg/analyzer/semver_test.go`

- [x] **Task 3: Implement Conflict Detector** (AC: #1, #2, #4)
  - [x] 3.1 Create `pkg/analyzer/conflict_detector.go`
  - [x] 3.2 Implement dependency collection from all packages (uses DependencyGraph.ExternalDeps)
  - [x] 3.3 Filter to only dependencies with 2+ different versions
  - [x] 3.4 Create tests in `pkg/analyzer/conflict_detector_test.go`

- [x] **Task 4: Implement Resolution Suggestions** (AC: #5)
  - [x] 4.1 Implement `generateResolution` with severity-based messaging
  - [x] 4.2 Implement `generateImpact` with conflict count and dependency name
  - [x] 4.3 Add tests for resolution and impact generation

- [x] **Task 5: Wire to Analyzer** (AC: all)
  - [x] 5.1 Update `pkg/analyzer/analyzer.go` to call ConflictDetector
  - [x] 5.2 Update AnalysisResult type to include VersionConflicts field
  - [x] 5.3 Update analyzer and handler tests for version conflicts

- [x] **Task 6: Performance Testing** (AC: #6)
  - [x] 6.1 Create `pkg/analyzer/conflict_detector_benchmark_test.go`
  - [x] 6.2 Verify 100 packages < 1 second (achieved ~0.7ms, far exceeding requirement)

- [x] **Task 7: Integration Verification** (AC: all)
  - [x] 7.1 Build WASM: `pnpm nx build @monoguard/analysis-engine` ✓
  - [x] 7.2 Add handler tests for conflict detection via WASM interface
  - [x] 7.3 Test with known conflict scenarios (patch, minor, major differences)
  - [x] 7.4 Verify all tests pass: `make test` ✓ (coverage >80% all packages)

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

Claude Opus 4.5 (claude-opus-4-5-20251101)

### Debug Log References

N/A - Implementation completed without significant issues.

### Completion Notes List

1. **VersionConflictInfo Types (Task 1)**: Created new types in `version_conflict.go` with full JSON serialization support. Used `VersionConflictInfo` to avoid conflict with deprecated `VersionConflict` type. All JSON fields use camelCase per project conventions.

2. **Semver Parsing (Task 2)**: Implemented comprehensive semver parsing supporting exact versions, caret (^), tilde (~), comparison ranges (>=, <), and wildcards (x, X). Added `FindMaxDifference` and `FindHighestVersion` helper functions.

3. **Conflict Detection (Task 3)**: Used `DependencyGraph.ExternalDeps` (pre-filtered by GraphBuilder) rather than raw WorkspaceData. This leverages existing work from Story 2.2 and avoids duplicate internal/external classification logic.

4. **Resolution/Impact Messages (Task 4)**: Implemented context-aware messages based on severity level. Messages include the recommended version and guidance appropriate for the conflict severity.

5. **Analyzer Integration (Task 5)**: Wired ConflictDetector into main Analyzer.Analyze() flow. VersionConflicts now returned alongside CircularDependencies in AnalysisResult.

6. **Performance (Task 6)**: Achieved ~0.7ms for 100 packages with 20 conflict-prone dependencies. This is approximately 1400x faster than the 1-second requirement. Scales well: 500 packages completes in ~6ms.

7. **Test Coverage**: All packages maintain >80% coverage:
   - handlers: 93.1%
   - analyzer: 91.7%
   - types: 91.4%
   - parser: 84.6%
   - result: 88.5%

### File List

**New Files:**
- `packages/analysis-engine/pkg/types/version_conflict.go` - VersionConflictInfo, ConflictingVersion, ConflictSeverity types
- `packages/analysis-engine/pkg/types/version_conflict_test.go` - JSON serialization tests for version conflict types
- `packages/analysis-engine/pkg/analyzer/semver.go` - SemVer parsing and comparison functions
- `packages/analysis-engine/pkg/analyzer/semver_test.go` - Tests for semver parsing
- `packages/analysis-engine/pkg/analyzer/conflict_detector.go` - ConflictDetector implementation
- `packages/analysis-engine/pkg/analyzer/conflict_detector_test.go` - Conflict detection tests
- `packages/analysis-engine/pkg/analyzer/conflict_detector_benchmark_test.go` - Performance benchmarks

**Modified Files:**
- `packages/analysis-engine/pkg/types/types.go` - Added VersionConflicts field to AnalysisResult, deprecated old VersionConflict type
- `packages/analysis-engine/pkg/analyzer/analyzer.go` - Integrated ConflictDetector into Analyze() method
- `packages/analysis-engine/pkg/analyzer/analyzer_test.go` - Added version conflict detection tests
- `packages/analysis-engine/internal/handlers/handlers_test.go` - Added WASM interface tests for version conflicts

## Change Log

| Date | Changes |
|------|---------|
| 2026-01-17 | Initial implementation of version conflict detection (Story 2.4) |
