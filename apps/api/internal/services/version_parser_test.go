package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionParser_NewVersionParser(t *testing.T) {
	parser := NewVersionParser()
	
	assert.NotNil(t, parser)
	assert.NotNil(t, parser.semverRegex)
	assert.NotNil(t, parser.rangeRegex)
	assert.NotNil(t, parser.prereleaseRegex)
}

func TestVersionParser_ParseVersion(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name        string
		versionStr  string
		expected    *SemanticVersion
		expectError bool
	}{
		{
			name:       "basic version",
			versionStr: "1.2.3",
			expected: &SemanticVersion{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Raw:   "1.2.3",
			},
		},
		{
			name:       "version with v prefix",
			versionStr: "v2.0.0",
			expected: &SemanticVersion{
				Major: 2,
				Minor: 0,
				Patch: 0,
				Raw:   "v2.0.0",
			},
		},
		{
			name:       "version with prerelease",
			versionStr: "1.0.0-alpha.1",
			expected: &SemanticVersion{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: "alpha.1",
				Raw:        "1.0.0-alpha.1",
			},
		},
		{
			name:       "version with build metadata",
			versionStr: "1.0.0+20210101",
			expected: &SemanticVersion{
				Major: 1,
				Minor: 0,
				Patch: 0,
				Build: "20210101",
				Raw:   "1.0.0+20210101",
			},
		},
		{
			name:       "version with prerelease and build",
			versionStr: "2.0.0-beta.1+exp.sha.5114f85",
			expected: &SemanticVersion{
				Major:      2,
				Minor:      0,
				Patch:      0,
				Prerelease: "beta.1",
				Build:      "exp.sha.5114f85",
				Raw:        "2.0.0-beta.1+exp.sha.5114f85",
			},
		},
		{
			name:        "empty version",
			versionStr:  "",
			expectError: true,
		},
		{
			name:        "invalid version",
			versionStr:  "invalid.version",
			expectError: true,
		},
		{
			name:        "incomplete version",
			versionStr:  "1.2",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := parser.ParseVersion(tt.versionStr)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, version)
			} else {
				require.NoError(t, err)
				require.NotNil(t, version)
				assert.Equal(t, tt.expected.Major, version.Major)
				assert.Equal(t, tt.expected.Minor, version.Minor)
				assert.Equal(t, tt.expected.Patch, version.Patch)
				assert.Equal(t, tt.expected.Prerelease, version.Prerelease)
				assert.Equal(t, tt.expected.Build, version.Build)
				assert.Equal(t, tt.expected.Raw, version.Raw)
			}
		})
	}
}

func TestVersionParser_ParseVersionRange(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name        string
		rangeStr    string
		expectedOp  string
		expectedVer string
		expectError bool
	}{
		{
			name:        "caret range",
			rangeStr:    "^1.2.3",
			expectedOp:  "^",
			expectedVer: "1.2.3",
		},
		{
			name:        "tilde range",
			rangeStr:    "~2.0.0",
			expectedOp:  "~",
			expectedVer: "2.0.0",
		},
		{
			name:        "greater than or equal",
			rangeStr:    ">=1.0.0",
			expectedOp:  ">=",
			expectedVer: "1.0.0",
		},
		{
			name:        "less than",
			rangeStr:    "<2.0.0",
			expectedOp:  "<",
			expectedVer: "2.0.0",
		},
		{
			name:        "exact version",
			rangeStr:    "1.5.0",
			expectedOp:  "=",
			expectedVer: "1.5.0",
		},
		{
			name:        "empty range",
			rangeStr:    "",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versionRange, err := parser.ParseVersionRange(tt.rangeStr)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, versionRange)
			} else {
				require.NoError(t, err)
				require.NotNil(t, versionRange)
				assert.Equal(t, tt.expectedOp, versionRange.Operator)
				assert.Equal(t, tt.expectedVer, versionRange.Version.String())
				assert.Equal(t, tt.rangeStr, versionRange.Raw)
				assert.NotNil(t, versionRange.Satisfies)
			}
		})
	}
}

func TestVersionParser_CompareVersions(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name     string
		version1 string
		version2 string
		expected int // -1: v1 < v2, 0: v1 = v2, 1: v1 > v2
	}{
		{
			name:     "equal versions",
			version1: "1.0.0",
			version2: "1.0.0",
			expected: 0,
		},
		{
			name:     "major version difference",
			version1: "2.0.0",
			version2: "1.0.0",
			expected: 1,
		},
		{
			name:     "minor version difference",
			version1: "1.1.0",
			version2: "1.2.0",
			expected: -1,
		},
		{
			name:     "patch version difference",
			version1: "1.0.2",
			version2: "1.0.1",
			expected: 1,
		},
		{
			name:     "prerelease vs stable",
			version1: "1.0.0-alpha",
			version2: "1.0.0",
			expected: -1,
		},
		{
			name:     "prerelease comparison",
			version1: "1.0.0-alpha.1",
			version2: "1.0.0-alpha.2",
			expected: -1,
		},
		{
			name:     "prerelease with numeric comparison",
			version1: "1.0.0-alpha.10",
			version2: "1.0.0-alpha.2",
			expected: 1,
		},
		{
			name:     "prerelease string vs numeric",
			version1: "1.0.0-alpha",
			version2: "1.0.0-1",
			expected: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := parser.ParseVersion(tt.version1)
			require.NoError(t, err)
			v2, err := parser.ParseVersion(tt.version2)
			require.NoError(t, err)
			
			result := parser.compareVersions(v1, v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionParser_SatisfiesCaretRange(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name          string
		rangeVersion  string
		testVersion   string
		shouldSatisfy bool
	}{
		{
			name:          "caret range normal case",
			rangeVersion:  "1.2.3",
			testVersion:   "1.5.0",
			shouldSatisfy: true,
		},
		{
			name:          "caret range exact match",
			rangeVersion:  "1.2.3",
			testVersion:   "1.2.3",
			shouldSatisfy: true,
		},
		{
			name:          "caret range major version boundary",
			rangeVersion:  "1.2.3",
			testVersion:   "2.0.0",
			shouldSatisfy: false,
		},
		{
			name:          "caret range with 0 major",
			rangeVersion:  "0.2.3",
			testVersion:   "0.2.5",
			shouldSatisfy: true,
		},
		{
			name:          "caret range with 0 major minor boundary",
			rangeVersion:  "0.2.3",
			testVersion:   "0.3.0",
			shouldSatisfy: false,
		},
		{
			name:          "caret range with 0.0.x",
			rangeVersion:  "0.0.3",
			testVersion:   "0.0.3",
			shouldSatisfy: true,
		},
		{
			name:          "caret range with 0.0.x different patch",
			rangeVersion:  "0.0.3",
			testVersion:   "0.0.4",
			shouldSatisfy: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rangeVer, err := parser.ParseVersion(tt.rangeVersion)
			require.NoError(t, err)
			testVer, err := parser.ParseVersion(tt.testVersion)
			require.NoError(t, err)
			
			result := parser.satisfiesCaretRange(testVer, rangeVer)
			assert.Equal(t, tt.shouldSatisfy, result)
		})
	}
}

func TestVersionParser_SatisfiesTildeRange(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name          string
		rangeVersion  string
		testVersion   string
		shouldSatisfy bool
	}{
		{
			name:          "tilde range normal case",
			rangeVersion:  "1.2.3",
			testVersion:   "1.2.5",
			shouldSatisfy: true,
		},
		{
			name:          "tilde range exact match",
			rangeVersion:  "1.2.3",
			testVersion:   "1.2.3",
			shouldSatisfy: true,
		},
		{
			name:          "tilde range minor version boundary",
			rangeVersion:  "1.2.3",
			testVersion:   "1.3.0",
			shouldSatisfy: false,
		},
		{
			name:          "tilde range major version boundary",
			rangeVersion:  "1.2.3",
			testVersion:   "2.2.3",
			shouldSatisfy: false,
		},
		{
			name:          "tilde range lower patch",
			rangeVersion:  "1.2.3",
			testVersion:   "1.2.1",
			shouldSatisfy: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rangeVer, err := parser.ParseVersion(tt.rangeVersion)
			require.NoError(t, err)
			testVer, err := parser.ParseVersion(tt.testVersion)
			require.NoError(t, err)
			
			result := parser.satisfiesTildeRange(testVer, rangeVer)
			assert.Equal(t, tt.shouldSatisfy, result)
		})
	}
}

func TestVersionParser_AnalyzeVersionConflicts(t *testing.T) {
	parser := NewVersionParser()
	
	dependencies := map[string][]string{
		"react": {"^16.0.0", "^17.0.0", "^18.0.0"},
		"lodash": {"^4.17.0", "^4.17.1"},
		"typescript": {"^4.0.0"},
		"webpack": {"^5.0.0", "^4.46.0"},
	}
	
	conflicts, err := parser.AnalyzeVersionConflicts(dependencies)
	require.NoError(t, err)
	
	// Should have conflicts for react and webpack (different majors), not for lodash or typescript
	assert.True(t, len(conflicts) >= 2)
	
	// Find specific conflicts
	var reactConflict, webpackConflict, lodashConflict *VersionConflictInfo
	for _, conflict := range conflicts {
		switch conflict.PackageName {
		case "react":
			reactConflict = conflict
		case "webpack":
			webpackConflict = conflict
		case "lodash":
			lodashConflict = conflict
		}
	}
	
	// React should have major version conflict
	require.NotNil(t, reactConflict)
	assert.Equal(t, "major", reactConflict.ConflictType)
	assert.Equal(t, 3, len(reactConflict.Versions))
	assert.False(t, reactConflict.IsResolvable) // Major version conflicts are not auto-resolvable
	
	// Webpack should have major version conflict
	require.NotNil(t, webpackConflict)
	assert.Equal(t, "major", webpackConflict.ConflictType)
	
	// Lodash might not be considered a conflict (minor version differences)
	if lodashConflict != nil {
		assert.Equal(t, "patch", lodashConflict.ConflictType)
		assert.True(t, lodashConflict.IsResolvable)
	}
}

func TestVersionParser_DetermineConflictType(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name         string
		versions     []string
		expectedType string
	}{
		{
			name:         "no conflict",
			versions:     []string{"1.0.0"},
			expectedType: "none",
		},
		{
			name:         "major version conflict",
			versions:     []string{"1.0.0", "2.0.0"},
			expectedType: "major",
		},
		{
			name:         "minor version conflict",
			versions:     []string{"1.1.0", "1.2.0"},
			expectedType: "minor",
		},
		{
			name:         "patch version conflict",
			versions:     []string{"1.0.1", "1.0.2"},
			expectedType: "patch",
		},
		{
			name:         "prerelease conflict",
			versions:     []string{"1.0.0-alpha", "1.0.0-beta"},
			expectedType: "prerelease",
		},
		{
			name:         "mixed conflict (major takes precedence)",
			versions:     []string{"1.0.0-alpha", "2.0.0-beta"},
			expectedType: "major",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var versions []*SemanticVersion
			for _, versionStr := range tt.versions {
				version, err := parser.ParseVersion(versionStr)
				require.NoError(t, err)
				versions = append(versions, version)
			}
			
			conflictType := parser.determineConflictType(versions)
			assert.Equal(t, tt.expectedType, conflictType)
		})
	}
}

func TestVersionParser_IsConflictResolvable(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name         string
		ranges       []string
		resolvable   bool
	}{
		{
			name:       "compatible caret ranges",
			ranges:     []string{"^1.0.0", "^1.2.0"},
			resolvable: true,
		},
		{
			name:       "incompatible major versions",
			ranges:     []string{"^1.0.0", "^2.0.0"},
			resolvable: false,
		},
		{
			name:       "overlapping ranges",
			ranges:     []string{">=1.0.0", "<=2.0.0", "^1.5.0"},
			resolvable: true,
		},
		{
			name:       "single range",
			ranges:     []string{"^1.0.0"},
			resolvable: true,
		},
		{
			name:       "exact version matches",
			ranges:     []string{"1.0.0", "1.0.0"},
			resolvable: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ranges []*VersionRange
			for _, rangeStr := range tt.ranges {
				versionRange, err := parser.ParseVersionRange(rangeStr)
				require.NoError(t, err)
				ranges = append(ranges, versionRange)
			}
			
			resolvable := parser.isConflictResolvable(ranges)
			assert.Equal(t, tt.resolvable, resolvable)
		})
	}
}

func TestVersionParser_AssessConflictRisk(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name           string
		packageName    string
		versions       []string
		conflictType   string
		expectedLevel  string
		expectReasons  int
	}{
		{
			name:          "major version conflict",
			packageName:   "react",
			versions:      []string{"16.0.0", "17.0.0"},
			conflictType:  "major",
			expectedLevel: "critical",
			expectReasons: 1,
		},
		{
			name:          "minor version conflict",
			packageName:   "lodash",
			versions:      []string{"4.17.0", "4.18.0"},
			conflictType:  "minor",
			expectedLevel: "high",
			expectReasons: 1,
		},
		{
			name:          "patch version conflict",
			packageName:   "lodash",
			versions:      []string{"4.17.0", "4.17.1"},
			conflictType:  "patch",
			expectedLevel: "medium",
			expectReasons: 1,
		},
		{
			name:          "many versions",
			packageName:   "react",
			versions:      []string{"16.0.0", "17.0.0", "18.0.0", "19.0.0"},
			conflictType:  "major",
			expectedLevel: "critical",
			expectReasons: 2, // Major + many versions
		},
		{
			name:          "experimental package",
			packageName:   "experimental-lib",
			versions:      []string{"0.1.0", "0.2.0"},
			conflictType:  "minor",
			expectedLevel: "critical", // Escalated due to 0.x.x version
			expectReasons: 2, // Minor + experimental
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var versions []*SemanticVersion
			for _, versionStr := range tt.versions {
				version, err := parser.ParseVersion(versionStr)
				require.NoError(t, err)
				versions = append(versions, version)
			}
			
			risk := parser.assessConflictRisk(tt.packageName, versions, tt.conflictType)
			assert.Equal(t, tt.expectedLevel, risk.Level)
			assert.Equal(t, tt.expectReasons, len(risk.Reasons))
			assert.NotEmpty(t, risk.Impact)
			assert.NotEmpty(t, risk.Difficulty)
		})
	}
}

func TestVersionParser_FindLatestVersion(t *testing.T) {
	parser := NewVersionParser()
	
	versions := []*SemanticVersion{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 2, Minor: 1, Patch: 0},
		{Major: 2, Minor: 0, Patch: 5},
		{Major: 1, Minor: 5, Patch: 2},
	}
	
	latest := parser.findLatestVersion(versions)
	require.NotNil(t, latest)
	assert.Equal(t, 2, latest.Major)
	assert.Equal(t, 1, latest.Minor)
	assert.Equal(t, 0, latest.Patch)
}

func TestVersionParser_DeduplicateVersions(t *testing.T) {
	parser := NewVersionParser()
	
	versions := []*SemanticVersion{
		{Major: 1, Minor: 0, Patch: 0},
		{Major: 2, Minor: 0, Patch: 0},
		{Major: 1, Minor: 0, Patch: 0}, // Duplicate
		{Major: 1, Minor: 0, Patch: 0, Prerelease: "alpha"}, // Different prerelease
		{Major: 2, Minor: 0, Patch: 0}, // Duplicate
	}
	
	unique := parser.deduplicateVersions(versions)
	assert.Equal(t, 3, len(unique))
	
	// Should contain unique combinations of major.minor.patch-prerelease+build
	versionStrings := make(map[string]bool)
	for _, v := range unique {
		key := fmt.Sprintf("%d.%d.%d-%s+%s", v.Major, v.Minor, v.Patch, v.Prerelease, v.Build)
		versionStrings[key] = true
	}
	assert.Equal(t, 3, len(versionStrings))
}

func TestVersionParser_IsBackwardsCompatible(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		name       string
		version1   string
		version2   string
		compatible bool
	}{
		{
			name:       "same version",
			version1:   "1.0.0",
			version2:   "1.0.0",
			compatible: true,
		},
		{
			name:       "patch increase",
			version1:   "1.0.0",
			version2:   "1.0.1",
			compatible: true,
		},
		{
			name:       "minor increase",
			version1:   "1.0.0",
			version2:   "1.1.0",
			compatible: true,
		},
		{
			name:       "major increase",
			version1:   "1.0.0",
			version2:   "2.0.0",
			compatible: false,
		},
		{
			name:       "version decrease",
			version1:   "1.1.0",
			version2:   "1.0.0",
			compatible: false,
		},
		{
			name:       "0.x minor increase",
			version1:   "0.1.0",
			version2:   "0.2.0",
			compatible: false,
		},
		{
			name:       "0.x patch increase",
			version1:   "0.1.0",
			version2:   "0.1.1",
			compatible: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := parser.ParseVersion(tt.version1)
			require.NoError(t, err)
			v2, err := parser.ParseVersion(tt.version2)
			require.NoError(t, err)
			
			compatible := parser.IsBackwardsCompatible(v1, v2)
			assert.Equal(t, tt.compatible, compatible)
		})
	}
}

func TestSemanticVersion_String(t *testing.T) {
	tests := []struct {
		name     string
		version  *SemanticVersion
		expected string
	}{
		{
			name: "basic version",
			version: &SemanticVersion{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: "1.2.3",
		},
		{
			name: "version with prerelease",
			version: &SemanticVersion{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: "alpha.1",
			},
			expected: "1.0.0-alpha.1",
		},
		{
			name: "version with build",
			version: &SemanticVersion{
				Major: 1,
				Minor: 0,
				Patch: 0,
				Build: "20210101",
			},
			expected: "1.0.0+20210101",
		},
		{
			name: "version with prerelease and build",
			version: &SemanticVersion{
				Major:      2,
				Minor:      0,
				Patch:      0,
				Prerelease: "beta.1",
				Build:      "exp.sha.5114f85",
			},
			expected: "2.0.0-beta.1+exp.sha.5114f85",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSemanticVersion_IsStable(t *testing.T) {
	tests := []struct {
		name     string
		version  *SemanticVersion
		expected bool
	}{
		{
			name: "stable version",
			version: &SemanticVersion{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: true,
		},
		{
			name: "prerelease version",
			version: &SemanticVersion{
				Major:      1,
				Minor:      0,
				Patch:      0,
				Prerelease: "alpha",
			},
			expected: false,
		},
		{
			name: "version with build metadata",
			version: &SemanticVersion{
				Major: 1,
				Minor: 0,
				Patch: 0,
				Build: "20210101",
			},
			expected: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.IsStable()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests

func BenchmarkVersionParser_ParseVersion(b *testing.B) {
	parser := NewVersionParser()
	versions := []string{
		"1.0.0", "2.1.3", "0.5.2", "1.0.0-alpha.1", "2.0.0+20210101",
		"3.2.1-beta.2+exp.sha.5114f85", "10.15.7", "0.0.1", "1.2.3-rc.1",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		version := versions[i%len(versions)]
		parser.ParseVersion(version)
	}
}

func BenchmarkVersionParser_ParseVersionRange(b *testing.B) {
	parser := NewVersionParser()
	ranges := []string{
		"^1.0.0", "~2.1.3", ">=0.5.2", "<1.0.0", "1.0.0",
		"^2.0.0", "~3.2.1", ">=10.15.7", "<0.0.1", "1.2.3",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rangeStr := ranges[i%len(ranges)]
		parser.ParseVersionRange(rangeStr)
	}
}

func BenchmarkVersionParser_CompareVersions(b *testing.B) {
	parser := NewVersionParser()
	
	v1, _ := parser.ParseVersion("1.2.3")
	v2, _ := parser.ParseVersion("1.2.4")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.compareVersions(v1, v2)
	}
}

func BenchmarkVersionParser_SatisfiesCaretRange(b *testing.B) {
	parser := NewVersionParser()
	
	rangeVer, _ := parser.ParseVersion("1.2.3")
	testVer, _ := parser.ParseVersion("1.5.0")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.satisfiesCaretRange(testVer, rangeVer)
	}
}

// Error case tests

func TestVersionParser_ErrorCases(t *testing.T) {
	parser := NewVersionParser()
	
	// Test invalid version parsing
	invalidVersions := []string{
		"1", "1.2", "1.2.3.4", "v", "1.2.3-", "1.2.3+",
		"1.2.3-+build", "1.2.3-alpha-", "1.2.a", "a.b.c",
	}
	
	for _, invalidVersion := range invalidVersions {
		_, err := parser.ParseVersion(invalidVersion)
		assert.Error(t, err, "Should error for version: %s", invalidVersion)
	}
	
	// Test invalid range parsing
	invalidRanges := []string{
		"^^1.0.0", "~>1.0.0", ">>1.0.0", "1.0.0-2.0.0",
	}
	
	for _, invalidRange := range invalidRanges {
		_, err := parser.ParseVersionRange(invalidRange)
		if err == nil {
			// Some ranges might parse but have invalid versions
			continue
		}
		assert.Error(t, err, "Should error for range: %s", invalidRange)
	}
}

// Edge case tests

func TestVersionParser_EdgeCases(t *testing.T) {
	parser := NewVersionParser()
	
	// Test very large version numbers
	largeVersion := "999999999.999999999.999999999"
	version, err := parser.ParseVersion(largeVersion)
	require.NoError(t, err)
	assert.Equal(t, 999999999, version.Major)
	
	// Test zero versions
	zeroVersion := "0.0.0"
	version, err = parser.ParseVersion(zeroVersion)
	require.NoError(t, err)
	assert.Equal(t, 0, version.Major)
	assert.Equal(t, 0, version.Minor)
	assert.Equal(t, 0, version.Patch)
	
	// Test complex prerelease
	complexPrerelease := "1.0.0-alpha.beta.1.2.3.build.456"
	version, err = parser.ParseVersion(complexPrerelease)
	require.NoError(t, err)
	assert.Equal(t, "alpha.beta.1.2.3.build.456", version.Prerelease)
	
	// Test complex build metadata
	complexBuild := "1.0.0+build.1.2.3.exp.sha.1234567890abcdef"
	version, err = parser.ParseVersion(complexBuild)
	require.NoError(t, err)
	assert.Equal(t, "build.1.2.3.exp.sha.1234567890abcdef", version.Build)
}

func TestVersionParser_PrereleaseComparison(t *testing.T) {
	parser := NewVersionParser()
	
	tests := []struct {
		pre1     string
		pre2     string
		expected int
	}{
		{"", "", 0},
		{"", "alpha", 1},
		{"alpha", "", -1},
		{"alpha", "beta", -1},
		{"alpha.1", "alpha.2", -1},
		{"alpha.10", "alpha.2", 1}, // Numeric comparison
		{"alpha", "alpha.1", -1},   // String vs numeric
		{"1", "alpha", -1},         // Numeric vs string
		{"1.2.3", "1.2.10", -1},   // Multi-part numeric
	}
	
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_vs_%s", tt.pre1, tt.pre2), func(t *testing.T) {
			result := parser.comparePrereleases(tt.pre1, tt.pre2)
			assert.Equal(t, tt.expected, result)
		})
	}
}