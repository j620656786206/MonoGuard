package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageJSONParser_NewPackageJSONParser(t *testing.T) {
	logger := logrus.New()
	config := DefaultParserConfig()
	
	parser := NewPackageJSONParser(logger, config)
	
	assert.NotNil(t, parser)
	assert.NotNil(t, parser.workspaceParser)
	assert.NotNil(t, parser.versionParser)
	assert.NotNil(t, parser.logger)
	assert.NotNil(t, parser.parseCache)
	assert.Equal(t, config, parser.config)
}

func TestPackageJSONParser_DefaultConfig(t *testing.T) {
	config := DefaultParserConfig()
	
	assert.True(t, config.EnableCaching)
	assert.Equal(t, time.Hour*1, config.CacheExpiry)
	assert.Equal(t, 10, config.MaxConcurrency)
	assert.True(t, config.EnableUsageAnalysis)
	assert.True(t, config.SkipNodeModules)
	assert.True(t, config.IncludeDevDeps)
	assert.False(t, config.DeepAnalysis)
	assert.True(t, config.MemoryOptimized)
	assert.Contains(t, config.ExcludePatterns, "node_modules")
	assert.Contains(t, config.ExcludePatterns, ".git")
	assert.Contains(t, config.ExcludePatterns, "dist")
}

func TestPackageJSONParser_BuildDependencyGraph(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	// Create test packages
	packages := []*WorkspacePackage{
		{
			Name:    "package-a",
			Version: "1.0.0",
			Path:    "/test/packages/a/package.json",
			Dependencies: map[string]string{
				"package-b": "^1.0.0",
				"lodash":    "^4.17.0",
			},
			DevDependencies: map[string]string{
				"typescript": "^4.0.0",
			},
		},
		{
			Name:    "package-b",
			Version: "1.0.0",
			Path:    "/test/packages/b/package.json",
			Dependencies: map[string]string{
				"react": "^18.0.0",
			},
			PeerDependencies: map[string]string{
				"react-dom": "^18.0.0",
			},
		},
		{
			Name:    "package-c",
			Version: "2.0.0",
			Path:    "/test/packages/c/package.json",
			Dependencies: map[string]string{
				"package-a": "^1.0.0",
			},
		},
	}
	
	graph, err := parser.buildDependencyGraph(packages)
	require.NoError(t, err)
	require.NotNil(t, graph)
	
	// Check nodes
	assert.Equal(t, 3, len(graph.Nodes))
	assert.Contains(t, graph.Nodes, "package-a")
	assert.Contains(t, graph.Nodes, "package-b")
	assert.Contains(t, graph.Nodes, "package-c")
	
	// Check package-a node
	nodeA := graph.Nodes["package-a"]
	assert.Equal(t, "package-a", nodeA.PackageName)
	assert.Equal(t, "1.0.0", nodeA.Version)
	assert.True(t, nodeA.IsWorkspace)
	assert.Contains(t, nodeA.Dependencies, "package-b")
	assert.Contains(t, nodeA.Dependencies, "lodash")
	assert.Contains(t, nodeA.DevDependencies, "typescript")
	assert.Contains(t, nodeA.Dependents, "package-c")
	
	// Check edges
	assert.True(t, len(graph.Edges) >= 5) // At least the dependencies we defined
	
	// Find specific edges
	var foundAtoB, foundCtoA bool
	for _, edge := range graph.Edges {
		if edge.From == "package-a" && edge.To == "package-b" {
			foundAtoB = true
			assert.Equal(t, "^1.0.0", edge.VersionRange)
			assert.Equal(t, "production", edge.Type)
		}
		if edge.From == "package-c" && edge.To == "package-a" {
			foundCtoA = true
			assert.Equal(t, "^1.0.0", edge.VersionRange)
			assert.Equal(t, "production", edge.Type)
		}
	}
	assert.True(t, foundAtoB)
	assert.True(t, foundCtoA)
}

func TestPackageJSONParser_AnalyzeVersionConflicts(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	packages := []*WorkspacePackage{
		{
			Name: "package-a",
			Dependencies: map[string]string{
				"react":  "^16.0.0",
				"lodash": "^4.17.0",
			},
		},
		{
			Name: "package-b",
			Dependencies: map[string]string{
				"react":  "^18.0.0",
				"lodash": "^4.17.1",
			},
		},
		{
			Name: "package-c",
			Dependencies: map[string]string{
				"react": "^17.0.0",
			},
			DevDependencies: map[string]string{
				"typescript": "^4.0.0",
			},
		},
	}
	
	conflicts, err := parser.analyzeVersionConflicts(packages)
	require.NoError(t, err)
	
	// Should have conflicts for react (major version differences)
	assert.True(t, len(conflicts) >= 1)
	
	// Find react conflict
	var reactConflict *VersionConflictInfo
	for _, conflict := range conflicts {
		if conflict.PackageName == "react" {
			reactConflict = conflict
			break
		}
	}
	
	require.NotNil(t, reactConflict)
	assert.Equal(t, "major", reactConflict.ConflictType)
	assert.Equal(t, 3, len(reactConflict.Versions))
	assert.NotNil(t, reactConflict.RiskAssessment)
	assert.Equal(t, "critical", reactConflict.RiskAssessment.Level)
}

func TestPackageJSONParser_FindDuplicateDependencies(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	packages := []*WorkspacePackage{
		{
			Name: "package-a",
			Dependencies: map[string]string{
				"lodash": "^4.17.0",
				"react":  "^18.0.0",
			},
		},
		{
			Name: "package-b",
			Dependencies: map[string]string{
				"lodash": "^4.17.1",
				"react":  "^18.0.0",
			},
		},
		{
			Name: "package-c",
			Dependencies: map[string]string{
				"lodash": "^4.18.0",
			},
			DevDependencies: map[string]string{
				"typescript": "^4.0.0",
			},
		},
	}
	
	duplicates, err := parser.findDuplicateDependencies(packages)
	require.NoError(t, err)
	
	// Should have duplicates for lodash (multiple versions)
	assert.Contains(t, duplicates, "lodash")
	
	lodashDuplicate := duplicates["lodash"]
	assert.Equal(t, "lodash", lodashDuplicate.PackageName)
	assert.Equal(t, 3, len(lodashDuplicate.Versions))
	assert.Equal(t, 3, len(lodashDuplicate.AffectedPackages))
	assert.NotNil(t, lodashDuplicate.EstimatedWaste)
	assert.NotNil(t, lodashDuplicate.ConsolidationPlan)
	
	// React should not be a duplicate (same version)
	assert.NotContains(t, duplicates, "react")
}

func TestPackageJSONParser_EstimateResourceWaste(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	waste := parser.estimateResourceWaste("lodash", 3)
	
	assert.NotEmpty(t, waste.DiskSpaceWaste)
	assert.NotEmpty(t, waste.BundleSizeWaste)
	assert.NotEmpty(t, waste.InstallTimeWaste)
	assert.NotEmpty(t, waste.MemoryWaste)
	assert.True(t, waste.WastePercentage > 0)
	assert.True(t, waste.WastePercentage < 100)
}

func TestPackageJSONParser_CreateConsolidationPlan(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	versions := []*SemanticVersion{
		{Major: 4, Minor: 17, Patch: 0, Raw: "4.17.0"},
		{Major: 4, Minor: 17, Patch: 1, Raw: "4.17.1"},
		{Major: 4, Minor: 18, Patch: 0, Raw: "4.18.0"},
	}
	
	ranges := []*VersionRange{
		{Raw: "^4.17.0", Version: versions[0]},
		{Raw: "^4.17.1", Version: versions[1]},
		{Raw: "^4.18.0", Version: versions[2]},
	}
	
	plan := parser.createConsolidationPlan("lodash", versions, ranges)
	
	assert.NotNil(t, plan)
	assert.NotNil(t, plan.TargetVersion)
	assert.True(t, len(plan.MigrationSteps) > 0)
	assert.NotNil(t, plan.RiskAssessment)
	assert.NotEmpty(t, plan.EstimatedEffort)
	assert.NotNil(t, plan.ExpectedSavings)
	
	// Target version should be the latest
	assert.Equal(t, 4, plan.TargetVersion.Major)
	assert.Equal(t, 18, plan.TargetVersion.Minor)
	assert.Equal(t, 0, plan.TargetVersion.Patch)
	
	// Should have migration steps for versions that need to change
	foundSteps := 0
	for _, step := range plan.MigrationSteps {
		if step.NewVersion == "4.18.0" && step.OldVersion != "4.18.0" {
			foundSteps++
		}
	}
	assert.Equal(t, 2, foundSteps) // Two versions need to be updated
}

func TestPackageJSONParser_AnalyzePackageManagerInfo(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	// Create temporary directory with different lock files
	tmpDir, err := os.MkdirTemp("", "package-manager-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Test with pnpm-lock.yaml
	pnpmLockPath := filepath.Join(tmpDir, "pnpm-lock.yaml")
	require.NoError(t, os.WriteFile(pnpmLockPath, []byte("lockfileVersion: 5.4"), 0644))
	
	info, err := parser.analyzePackageManagerInfo(tmpDir)
	require.NoError(t, err)
	require.NotNil(t, info)
	
	assert.Equal(t, "pnpm", info.Type)
	assert.True(t, info.LockFilePresent)
	assert.Equal(t, pnpmLockPath, info.LockFilePath)
	
	// Remove pnpm lock and test with yarn.lock
	require.NoError(t, os.Remove(pnpmLockPath))
	yarnLockPath := filepath.Join(tmpDir, "yarn.lock")
	require.NoError(t, os.WriteFile(yarnLockPath, []byte("# yarn.lock"), 0644))
	
	info, err = parser.analyzePackageManagerInfo(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "yarn", info.Type)
	assert.Equal(t, yarnLockPath, info.LockFilePath)
	
	// Remove yarn lock and test with package-lock.json
	require.NoError(t, os.Remove(yarnLockPath))
	npmLockPath := filepath.Join(tmpDir, "package-lock.json")
	require.NoError(t, os.WriteFile(npmLockPath, []byte("{}"), 0644))
	
	info, err = parser.analyzePackageManagerInfo(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "npm", info.Type)
	assert.Equal(t, npmLockPath, info.LockFilePath)
	
	// Remove all locks and test default
	require.NoError(t, os.Remove(npmLockPath))
	
	info, err = parser.analyzePackageManagerInfo(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "npm", info.Type)
	assert.False(t, info.LockFilePresent)
	assert.Empty(t, info.LockFilePath)
}

func TestPackageJSONParser_GenerateRecoveryActions(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	errors := []error{
		fmt.Errorf("workspace discovery failed: permission denied"),
		fmt.Errorf("package discovery error: malformed JSON"),
		fmt.Errorf("version conflict analysis failed: invalid range"),
		fmt.Errorf("unknown error occurred"),
	}
	
	actions := parser.generateRecoveryActions(errors)
	
	assert.Equal(t, 4, len(actions))
	assert.Contains(t, actions, "Manual workspace configuration validation recommended")
	assert.Contains(t, actions, "Check package.json file permissions and syntax")
	assert.Contains(t, actions, "Manual version range validation needed")
	assert.Contains(t, actions, "Review error logs for specific issues")
}

func TestPackageJSONParser_CacheManagement(t *testing.T) {
	logger := logrus.New()
	config := &ParserConfig{
		EnableCaching: true,
		CacheExpiry:   time.Millisecond * 100, // Short expiry for testing
	}
	parser := NewPackageJSONParser(logger, config)
	
	// Test cache miss
	result := parser.getCachedResult("/test/path")
	assert.Nil(t, result)
	
	// Test cache set and hit
	testResult := &ParsedRepository{
		RootPath: "/test/path",
		CachedAt: time.Now(),
	}
	parser.cacheResult("/test/path", testResult)
	
	cached := parser.getCachedResult("/test/path")
	assert.NotNil(t, cached)
	assert.Equal(t, "/test/path", cached.RootPath)
	
	// Test cache expiry
	time.Sleep(time.Millisecond * 150) // Wait for cache to expire
	expired := parser.getCachedResult("/test/path")
	assert.Nil(t, expired)
	
	// Test cache clear
	parser.cacheResult("/test/path", testResult)
	parser.ClearCache()
	cleared := parser.getCachedResult("/test/path")
	assert.Nil(t, cleared)
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
		{1536*1024*1024, "1.50 GB"},
	}
	
	for _, tt := range tests {
		t.Run(fmt.Sprintf("bytes_%d", tt.bytes), func(t *testing.T) {
			result := formatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Integration tests

func TestPackageJSONParser_IntegrationParseRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests
	
	config := &ParserConfig{
		EnableCaching:        true,
		CacheExpiry:         time.Hour,
		MaxConcurrency:      2,
		EnableUsageAnalysis: false, // Disable to speed up tests
		SkipNodeModules:     true,
		IncludeDevDeps:      true,
		DeepAnalysis:        false,
		MemoryOptimized:     true,
		ExcludePatterns:     []string{"node_modules", ".git", "dist"},
	}
	
	parser := NewPackageJSONParser(logger, config)
	
	// Create a test monorepo structure
	tmpDir, err := os.MkdirTemp("", "integration-parse-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Create workspace configuration
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	packageJSONContent := `{
		"name": "test-monorepo",
		"private": true,
		"workspaces": ["packages/*", "apps/*"]
	}`
	require.NoError(t, os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644))
	
	// Create packages
	packagesDir := filepath.Join(tmpDir, "packages")
	require.NoError(t, os.MkdirAll(packagesDir, 0755))
	
	// Package A
	pkgADir := filepath.Join(packagesDir, "package-a")
	require.NoError(t, os.MkdirAll(pkgADir, 0755))
	pkgAJSON := `{
		"name": "package-a",
		"version": "1.0.0",
		"dependencies": {
			"lodash": "^4.17.0",
			"package-b": "^1.0.0"
		},
		"devDependencies": {
			"typescript": "^4.0.0"
		}
	}`
	require.NoError(t, os.WriteFile(filepath.Join(pkgADir, "package.json"), []byte(pkgAJSON), 0644))
	
	// Package B
	pkgBDir := filepath.Join(packagesDir, "package-b")
	require.NoError(t, os.MkdirAll(pkgBDir, 0755))
	pkgBJSON := `{
		"name": "package-b",
		"version": "1.0.0",
		"dependencies": {
			"lodash": "^4.18.0",
			"react": "^18.0.0"
		},
		"peerDependencies": {
			"react-dom": "^18.0.0"
		}
	}`
	require.NoError(t, os.WriteFile(filepath.Join(pkgBDir, "package.json"), []byte(pkgBJSON), 0644))
	
	// Package C (with conflicts)
	pkgCDir := filepath.Join(packagesDir, "package-c")
	require.NoError(t, os.MkdirAll(pkgCDir, 0755))
	pkgCJSON := `{
		"name": "package-c",
		"version": "2.0.0",
		"dependencies": {
			"react": "^16.0.0",
			"lodash": "^4.17.21"
		}
	}`
	require.NoError(t, os.WriteFile(filepath.Join(pkgCDir, "package.json"), []byte(pkgCJSON), 0644))
	
	// Create apps directory
	appsDir := filepath.Join(tmpDir, "apps")
	require.NoError(t, os.MkdirAll(appsDir, 0755))
	
	// App 1
	app1Dir := filepath.Join(appsDir, "app1")
	require.NoError(t, os.MkdirAll(app1Dir, 0755))
	app1JSON := `{
		"name": "app1",
		"version": "0.1.0",
		"dependencies": {
			"package-a": "^1.0.0",
			"react": "^17.0.0"
		}
	}`
	require.NoError(t, os.WriteFile(filepath.Join(app1Dir, "package.json"), []byte(app1JSON), 0644))
	
	// Parse the repository
	ctx := context.Background()
	result, err := parser.ParseRepository(ctx, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, result)
	
	// Verify basic structure
	assert.Equal(t, tmpDir, result.RootPath)
	assert.Equal(t, 1, len(result.WorkspaceConfigs))
	assert.Equal(t, 4, len(result.Packages)) // 3 packages + 1 app
	assert.NotNil(t, result.DependencyGraph)
	assert.NotNil(t, result.PackageManagerInfo)
	assert.NotNil(t, result.ParseMetadata)
	
	// Verify workspace configuration
	workspace := result.WorkspaceConfigs[0]
	assert.Equal(t, WorkspaceTypeNpm, workspace.Type)
	assert.Contains(t, workspace.Packages, "packages/*")
	assert.Contains(t, workspace.Packages, "apps/*")
	
	// Verify packages
	packageNames := make(map[string]bool)
	for _, pkg := range result.Packages {
		packageNames[pkg.Name] = true
	}
	assert.True(t, packageNames["package-a"])
	assert.True(t, packageNames["package-b"])
	assert.True(t, packageNames["package-c"])
	assert.True(t, packageNames["app1"])
	
	// Verify dependency graph
	assert.True(t, len(result.DependencyGraph.Nodes) >= 4)
	assert.True(t, len(result.DependencyGraph.Edges) > 0)
	
	// Verify conflicts (should have react version conflicts)
	assert.True(t, len(result.VersionConflicts) > 0)
	
	var reactConflict *VersionConflictInfo
	for _, conflict := range result.VersionConflicts {
		if conflict.PackageName == "react" {
			reactConflict = conflict
			break
		}
	}
	require.NotNil(t, reactConflict, "Should have react version conflict")
	assert.Equal(t, "major", reactConflict.ConflictType)
	
	// Verify duplicates (should have lodash duplicates)
	assert.Contains(t, result.DuplicateDependencies, "lodash")
	lodashDup := result.DuplicateDependencies["lodash"]
	assert.Equal(t, "lodash", lodashDup.PackageName)
	assert.True(t, len(lodashDup.Versions) >= 2)
	
	// Verify metadata
	assert.True(t, result.ParseMetadata.Duration > 0)
	assert.Equal(t, 4, result.ParseMetadata.PackagesProcessed)
	assert.True(t, result.ParseMetadata.FilesScanned >= 4)
	
	// Test caching - second call should be faster
	start := time.Now()
	cachedResult, err := parser.ParseRepository(ctx, tmpDir)
	duration := time.Since(start)
	
	require.NoError(t, err)
	require.NotNil(t, cachedResult)
	assert.True(t, duration < time.Millisecond*100, "Cached result should be much faster")
	assert.Equal(t, result.RootPath, cachedResult.RootPath)
	assert.Equal(t, len(result.Packages), len(cachedResult.Packages))
}

// Benchmark tests

func BenchmarkPackageJSONParser_ParseRepository(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	
	config := &ParserConfig{
		EnableCaching:        false, // Disable caching for benchmarks
		MaxConcurrency:      4,
		EnableUsageAnalysis: false,
		SkipNodeModules:     true,
		IncludeDevDeps:      true,
		DeepAnalysis:        false,
		MemoryOptimized:     true,
	}
	
	parser := NewPackageJSONParser(logger, config)
	
	// Create a simple test structure (reuse from integration test)
	tmpDir, err := os.MkdirTemp("", "benchmark-parse-test-*")
	require.NoError(b, err)
	defer os.RemoveAll(tmpDir)
	
	// Create minimal workspace
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	packageJSONContent := `{"name": "test", "workspaces": ["packages/*"]}`
	require.NoError(b, os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644))
	
	// Create test packages
	for i := 0; i < 5; i++ {
		pkgDir := filepath.Join(tmpDir, "packages", fmt.Sprintf("package-%d", i))
		require.NoError(b, os.MkdirAll(pkgDir, 0755))
		pkgJSON := fmt.Sprintf(`{
			"name": "package-%d",
			"version": "1.0.0",
			"dependencies": {
				"lodash": "^4.17.0",
				"react": "^18.0.0"
			}
		}`, i)
		require.NoError(b, os.WriteFile(filepath.Join(pkgDir, "package.json"), []byte(pkgJSON), 0644))
	}
	
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.ClearCache() // Ensure no caching between runs
		_, err := parser.ParseRepository(ctx, tmpDir)
		require.NoError(b, err)
	}
}

func BenchmarkPackageJSONParser_BuildDependencyGraph(b *testing.B) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	// Create test packages
	packages := make([]*WorkspacePackage, 10)
	for i := 0; i < 10; i++ {
		packages[i] = &WorkspacePackage{
			Name:    fmt.Sprintf("package-%d", i),
			Version: "1.0.0",
			Path:    fmt.Sprintf("/test/packages/%d/package.json", i),
			Dependencies: map[string]string{
				"lodash": "^4.17.0",
				"react":  "^18.0.0",
			},
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.buildDependencyGraph(packages)
		require.NoError(b, err)
	}
}

func BenchmarkPackageJSONParser_FindDuplicates(b *testing.B) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	// Create packages with duplicates
	packages := make([]*WorkspacePackage, 20)
	for i := 0; i < 20; i++ {
		packages[i] = &WorkspacePackage{
			Name: fmt.Sprintf("package-%d", i),
			Dependencies: map[string]string{
				"lodash": fmt.Sprintf("^4.%d.0", 17+(i%3)), // Create some duplicates
				"react":  fmt.Sprintf("^%d.0.0", 16+(i%3)), // Create major version conflicts
			},
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.findDuplicateDependencies(packages)
		require.NoError(b, err)
	}
}

// Error handling tests

func TestPackageJSONParser_ErrorHandling(t *testing.T) {
	logger := logrus.New()
	parser := NewPackageJSONParser(logger, nil)
	
	ctx := context.Background()
	
	// Test with non-existent directory
	result, err := parser.ParseRepository(ctx, "/non/existent/path")
	assert.Error(t, err)
	assert.Nil(t, result)
	
	// Test with empty directory
	tmpDir, err := os.MkdirTemp("", "empty-dir-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	result, err = parser.ParseRepository(ctx, tmpDir)
	assert.Error(t, err)
	assert.Nil(t, result)
	
	// Test with malformed package.json
	malformedDir, err := os.MkdirTemp("", "malformed-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(malformedDir)
	
	malformedJSON := filepath.Join(malformedDir, "package.json")
	require.NoError(t, os.WriteFile(malformedJSON, []byte(`{invalid json`), 0644))
	
	result, err = parser.ParseRepository(ctx, malformedDir)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPackageJSONParser_GracefulDegradation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise
	
	parser := NewPackageJSONParser(logger, nil)
	
	// Create directory with mixed valid and invalid packages
	tmpDir, err := os.MkdirTemp("", "graceful-degradation-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Create workspace config
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	packageJSONContent := `{"name": "test", "workspaces": ["packages/*"]}`
	require.NoError(t, os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644))
	
	// Create packages directory
	packagesDir := filepath.Join(tmpDir, "packages")
	require.NoError(t, os.MkdirAll(packagesDir, 0755))
	
	// Valid package
	validPkgDir := filepath.Join(packagesDir, "valid")
	require.NoError(t, os.MkdirAll(validPkgDir, 0755))
	validJSON := `{"name": "valid", "version": "1.0.0"}`
	require.NoError(t, os.WriteFile(filepath.Join(validPkgDir, "package.json"), []byte(validJSON), 0644))
	
	// Invalid package (malformed JSON)
	invalidPkgDir := filepath.Join(packagesDir, "invalid")
	require.NoError(t, os.MkdirAll(invalidPkgDir, 0755))
	invalidJSON := `{invalid json`
	require.NoError(t, os.WriteFile(filepath.Join(invalidPkgDir, "package.json"), []byte(invalidJSON), 0644))
	
	ctx := context.Background()
	
	// Should handle graceful degradation - valid package should still be processed
	result, err := parser.ParseRepository(ctx, tmpDir)
	
	// Should still return results even with some errors
	require.NoError(t, err)
	require.NotNil(t, result)
	
	// Should have processed the valid package
	assert.Equal(t, 1, len(result.Packages))
	assert.Equal(t, "valid", result.Packages[0].Name)
	
	// Should have recorded errors in metadata
	assert.True(t, result.ParseMetadata.ErrorsEncountered >= 0)
	assert.True(t, len(result.ParseMetadata.RecoveryActions) >= 0)
}