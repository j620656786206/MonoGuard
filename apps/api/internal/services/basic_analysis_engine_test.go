package services

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBasicAnalysisEngine tests the basic analysis engine functionality
func TestBasicAnalysisEngine(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel) // Reduce noise in tests

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "mono-guard-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test package.json files
	createTestPackageFiles(t, tempDir)

	// Initialize the analysis engine
	engine := NewBasicAnalysisEngine(logger)

	// Run the analysis
	results, err := engine.AnalyzeRepository(context.Background(), tempDir, "test-project")
	require.NoError(t, err)
	require.NotNil(t, results)

	// Verify the results
	assert.GreaterOrEqual(t, len(results.DuplicateDependencies), 1, "Should find duplicate dependencies")
	assert.GreaterOrEqual(t, len(results.VersionConflicts), 1, "Should find version conflicts")
	assert.NotNil(t, results.BundleImpact, "Should generate bundle impact report")
	assert.NotNil(t, results.Summary, "Should generate summary")

	// Verify specific duplicate dependency
	foundLodash := false
	for _, duplicate := range results.DuplicateDependencies {
		if duplicate.PackageName == "lodash" {
			foundLodash = true
			assert.Equal(t, 2, len(duplicate.Versions), "Should find 2 versions of lodash")
			assert.Contains(t, duplicate.Versions, "^4.17.21")
			assert.Contains(t, duplicate.Versions, "^4.16.0")
			break
		}
	}
	assert.True(t, foundLodash, "Should find lodash as a duplicate dependency")

	// Verify specific version conflict
	foundReact := false
	for _, conflict := range results.VersionConflicts {
		if conflict.PackageName == "react" {
			foundReact = true
			assert.Equal(t, 2, len(conflict.ConflictingVersions), "Should find 2 conflicting versions of react")
			break
		}
	}
	assert.True(t, foundReact, "Should find react as a version conflict")

	// Verify summary
	assert.Equal(t, 3, results.Summary.TotalPackages, "Should count 3 packages")
	assert.GreaterOrEqual(t, results.Summary.DuplicateCount, 1, "Should have duplicate count")
	assert.GreaterOrEqual(t, results.Summary.ConflictCount, 1, "Should have conflict count")
	assert.LessOrEqual(t, results.Summary.HealthScore, 100.0, "Health score should not exceed 100")
	assert.GreaterOrEqual(t, results.Summary.HealthScore, 0.0, "Health score should not be negative")
}

// TestMonoGuardPackageParser tests the package parser functionality
func TestMonoGuardPackageParser(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	tempDir, err := os.MkdirTemp("", "mono-guard-parser-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test package.json
	packageJSON := map[string]interface{}{
		"name":    "test-package",
		"version": "1.0.0",
		"dependencies": map[string]string{
			"lodash": "^4.17.21",
			"react":  "^18.0.0",
		},
		"devDependencies": map[string]string{
			"jest": "^29.0.0",
		},
	}

	packagePath := filepath.Join(tempDir, "package.json")
	writeJSONFile(t, packagePath, packageJSON)

	parser := NewMonoGuardPackageParser(logger)
	packages, err := parser.ParseRepository(context.Background(), tempDir)
	require.NoError(t, err)
	require.Len(t, packages, 1)

	pkg := packages[0]
	assert.Equal(t, "test-package", pkg.Name)
	assert.Equal(t, "1.0.0", pkg.Version)
	assert.Equal(t, 2, len(pkg.Dependencies))
	assert.Equal(t, 1, len(pkg.DevDependencies))
	assert.Equal(t, "^4.17.21", pkg.Dependencies["lodash"])
	assert.Equal(t, "^18.0.0", pkg.Dependencies["react"])
	assert.Equal(t, "^29.0.0", pkg.DevDependencies["jest"])
}

// TestMonoGuardDuplicateDetector tests the duplicate detector functionality
func TestMonoGuardDuplicateDetector(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	detector := NewMonoGuardDuplicateDetector(logger)

	packages := []*MonoGuardPackageInfo{
		{
			Name: "package-a",
			Dependencies: map[string]string{
				"lodash": "^4.17.21",
				"react":  "^18.0.0",
			},
		},
		{
			Name: "package-b",
			Dependencies: map[string]string{
				"lodash": "^4.16.0",
				"react":  "^18.0.0",
			},
		},
	}

	duplicates, err := detector.FindDuplicates(packages)
	require.NoError(t, err)
	require.Len(t, duplicates, 1) // Only lodash should be a duplicate

	duplicate := duplicates[0]
	assert.Equal(t, "lodash", duplicate.PackageName)
	assert.Equal(t, 2, len(duplicate.Versions))
	assert.Contains(t, duplicate.Versions, "^4.17.21")
	assert.Contains(t, duplicate.Versions, "^4.16.0")
	assert.NotEmpty(t, duplicate.Recommendation)
	assert.NotEmpty(t, duplicate.MigrationSteps)
}

// TestMonoGuardConflictAnalyzer tests the version conflict analyzer functionality
func TestMonoGuardConflictAnalyzer(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	analyzer := NewMonoGuardConflictAnalyzer(logger)

	packages := []*MonoGuardPackageInfo{
		{
			Name: "package-a",
			Dependencies: map[string]string{
				"react": "^17.0.0", // Major version 17
			},
		},
		{
			Name: "package-b",
			Dependencies: map[string]string{
				"react": "^18.0.0", // Major version 18 - conflict!
			},
		},
	}

	conflicts, err := analyzer.FindConflicts(packages)
	require.NoError(t, err)
	require.Len(t, conflicts, 1)

	conflict := conflicts[0]
	assert.Equal(t, "react", conflict.PackageName)
	assert.Equal(t, 2, len(conflict.ConflictingVersions))
	assert.NotEmpty(t, conflict.Resolution)
	assert.NotEmpty(t, conflict.Impact)
}

// createTestPackageFiles creates test package.json files for testing
func createTestPackageFiles(t *testing.T, tempDir string) {
	// Root package.json
	rootPackage := map[string]interface{}{
		"name":    "mono-repo-root",
		"version": "1.0.0",
		"workspaces": []string{
			"packages/*",
		},
		"devDependencies": map[string]string{
			"jest": "^29.0.0",
		},
	}
	writeJSONFile(t, filepath.Join(tempDir, "package.json"), rootPackage)

	// Package A
	packageADir := filepath.Join(tempDir, "packages", "package-a")
	err := os.MkdirAll(packageADir, 0755)
	require.NoError(t, err)

	packageA := map[string]interface{}{
		"name":    "package-a",
		"version": "1.0.0",
		"dependencies": map[string]string{
			"lodash": "^4.17.21",
			"react":  "^17.0.0", // Different major version
		},
	}
	writeJSONFile(t, filepath.Join(packageADir, "package.json"), packageA)

	// Package B
	packageBDir := filepath.Join(tempDir, "packages", "package-b")
	err = os.MkdirAll(packageBDir, 0755)
	require.NoError(t, err)

	packageB := map[string]interface{}{
		"name":    "package-b",
		"version": "1.0.0",
		"dependencies": map[string]string{
			"lodash": "^4.16.0",  // Different minor version (duplicate)
			"react":  "^18.0.0",  // Different major version (conflict)
		},
	}
	writeJSONFile(t, filepath.Join(packageBDir, "package.json"), packageB)
}

// writeJSONFile writes a JSON object to a file
func writeJSONFile(t *testing.T, path string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	require.NoError(t, err)

	err = os.WriteFile(path, jsonData, 0644)
	require.NoError(t, err)
}