package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceParser_NewWorkspaceParser(t *testing.T) {
	parser := NewWorkspaceParser()
	
	assert.NotNil(t, parser)
	assert.NotNil(t, parser.cache)
	assert.NotNil(t, parser.packageCache)
	assert.NotNil(t, parser.stringInternMap)
}

func TestWorkspaceParser_StringInterning(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Test string interning works correctly
	str1 := parser.intern("test-string")
	str2 := parser.intern("test-string")
	str3 := parser.intern("different-string")
	
	// Same strings should have same pointer
	assert.Equal(t, str1, str2)
	assert.NotEqual(t, str1, str3)
	
	// Verify internal map
	assert.Equal(t, 2, len(parser.stringInternMap))
}

func TestWorkspaceParser_InternSlice(t *testing.T) {
	parser := NewWorkspaceParser()
	
	original := []string{"package1", "package2", "package1"}
	interned := parser.internSlice(original)
	
	assert.Equal(t, len(original), len(interned))
	assert.Equal(t, original, interned)
	
	// Verify interning was applied
	assert.Equal(t, 2, len(parser.stringInternMap))
}

func TestWorkspaceParser_InternMap(t *testing.T) {
	parser := NewWorkspaceParser()
	
	original := map[string]string{
		"dep1": "^1.0.0",
		"dep2": "~2.0.0",
		"dep1": "^1.0.0", // Duplicate key to test interning
	}
	
	interned := parser.internMap(original)
	
	assert.Equal(t, len(original), len(interned))
	assert.Equal(t, original, interned)
	
	// Should have interned both keys and values
	assert.True(t, len(parser.stringInternMap) >= 3) // At least dep1, dep2, ^1.0.0, ~2.0.0
}

func TestWorkspaceParser_ExpandPackagePattern(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Create temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "workspace-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Create test directories
	packagesDir := filepath.Join(tmpDir, "packages")
	require.NoError(t, os.MkdirAll(packagesDir, 0755))
	
	pkg1Dir := filepath.Join(packagesDir, "package1")
	pkg2Dir := filepath.Join(packagesDir, "package2")
	require.NoError(t, os.MkdirAll(pkg1Dir, 0755))
	require.NoError(t, os.MkdirAll(pkg2Dir, 0755))
	
	// Create package.json files
	require.NoError(t, os.WriteFile(filepath.Join(pkg1Dir, "package.json"), []byte(`{"name":"package1"}`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(pkg2Dir, "package.json"), []byte(`{"name":"package2"}`), 0644))
	
	// Test pattern expansion
	pattern := "packages/*"
	matches, err := parser.expandPackagePattern(tmpDir, pattern)
	require.NoError(t, err)
	
	assert.Equal(t, 2, len(matches))
	assert.Contains(t, matches, pkg1Dir)
	assert.Contains(t, matches, pkg2Dir)
}

func TestWorkspaceParser_ParseNpmWorkspace(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Create temporary package.json with workspaces
	tmpDir, err := os.MkdirTemp("", "npm-workspace-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	packageJSONContent := `{
		"name": "test-monorepo",
		"workspaces": [
			"packages/*",
			"apps/*"
		]
	}`
	require.NoError(t, os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644))
	
	workspace, err := parser.parseNpmWorkspace(packageJSONPath, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, workspace)
	
	assert.Equal(t, WorkspaceTypeNpm, workspace.Type)
	assert.Equal(t, tmpDir, workspace.RootPath)
	assert.Equal(t, packageJSONPath, workspace.ConfigPath)
	assert.Equal(t, []string{"packages/*", "apps/*"}, workspace.Packages)
	assert.Equal(t, 4, workspace.Priority)
}

func TestWorkspaceParser_ParsePnpmWorkspace(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Create temporary pnpm-workspace.yaml
	tmpDir, err := os.MkdirTemp("", "pnpm-workspace-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	workspaceYAMLPath := filepath.Join(tmpDir, "pnpm-workspace.yaml")
	workspaceYAMLContent := `packages:
  - 'packages/*'
  - 'apps/*'
  - '!**/test/**'`
	require.NoError(t, os.WriteFile(workspaceYAMLPath, []byte(workspaceYAMLContent), 0644))
	
	workspace, err := parser.parsePnpmWorkspace(workspaceYAMLPath, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, workspace)
	
	assert.Equal(t, WorkspaceTypePnpm, workspace.Type)
	assert.Equal(t, tmpDir, workspace.RootPath)
	assert.Equal(t, workspaceYAMLPath, workspace.ConfigPath)
	assert.Equal(t, []string{"packages/*", "apps/*", "!**/test/**"}, workspace.Packages)
	assert.Equal(t, "pnpm", workspace.PackageManager)
	assert.Equal(t, 1, workspace.Priority)
}

func TestWorkspaceParser_ParseLernaWorkspace(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Create temporary lerna.json
	tmpDir, err := os.MkdirTemp("", "lerna-workspace-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	lernaJSONPath := filepath.Join(tmpDir, "lerna.json")
	lernaJSONContent := `{
		"version": "1.0.0",
		"packages": [
			"packages/*",
			"tools/*"
		],
		"npmClient": "npm",
		"useWorkspaces": true
	}`
	require.NoError(t, os.WriteFile(lernaJSONPath, []byte(lernaJSONContent), 0644))
	
	workspace, err := parser.parseLernaWorkspace(lernaJSONPath, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, workspace)
	
	assert.Equal(t, WorkspaceTypeLerna, workspace.Type)
	assert.Equal(t, tmpDir, workspace.RootPath)
	assert.Equal(t, lernaJSONPath, workspace.ConfigPath)
	assert.Equal(t, []string{"packages/*", "tools/*"}, workspace.Packages)
	assert.Equal(t, "npm", workspace.PackageManager)
	assert.Equal(t, 2, workspace.Priority)
	
	// Check metadata
	metadata := workspace.Metadata
	assert.Equal(t, "1.0.0", metadata["version"])
	assert.Equal(t, true, metadata["useWorkspaces"])
}

func TestWorkspaceParser_ParseNxWorkspace(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Create temporary nx.json
	tmpDir, err := os.MkdirTemp("", "nx-workspace-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	nxJSONPath := filepath.Join(tmpDir, "nx.json")
	nxJSONContent := `{
		"version": 2,
		"projects": {
			"app1": "apps/app1",
			"lib1": "libs/lib1"
		},
		"workspaceLayout": {
			"appsDir": "apps",
			"libsDir": "libs"
		}
	}`
	require.NoError(t, os.WriteFile(nxJSONPath, []byte(nxJSONContent), 0644))
	
	workspace, err := parser.parseNxWorkspace(nxJSONPath, tmpDir)
	require.NoError(t, err)
	require.NotNil(t, workspace)
	
	assert.Equal(t, WorkspaceTypeNx, workspace.Type)
	assert.Equal(t, tmpDir, workspace.RootPath)
	assert.Equal(t, nxJSONPath, workspace.ConfigPath)
	assert.Equal(t, "npm", workspace.PackageManager)
	assert.Equal(t, 3, workspace.Priority)
	
	// Check that packages contain project names and workspace layout dirs
	assert.Contains(t, workspace.Packages, "app1")
	assert.Contains(t, workspace.Packages, "lib1")
	assert.Contains(t, workspace.Packages, "libs/*")
	assert.Contains(t, workspace.Packages, "apps/*")
}

func TestWorkspaceParser_ResolveWorkspaceConflicts(t *testing.T) {
	parser := NewWorkspaceParser()
	
	workspaces := []*WorkspaceConfiguration{
		{
			Type:     WorkspaceTypePnpm,
			RootPath: "/test/path",
			Priority: 1,
		},
		{
			Type:     WorkspaceTypeLerna,
			RootPath: "/test/path",
			Priority: 2,
		},
		{
			Type:     WorkspaceTypeNpm,
			RootPath: "/test/path",
			Priority: 4,
		},
		{
			Type:     WorkspaceTypeNpm,
			RootPath: "/different/path",
			Priority: 4,
		},
	}
	
	resolved := parser.resolveWorkspaceConflicts(workspaces)
	
	assert.Equal(t, 2, len(resolved))
	
	// Should keep highest priority (lowest number) for each root path
	var testPathWorkspace, differentPathWorkspace *WorkspaceConfiguration
	for _, ws := range resolved {
		if ws.RootPath == "/test/path" {
			testPathWorkspace = ws
		} else if ws.RootPath == "/different/path" {
			differentPathWorkspace = ws
		}
	}
	
	require.NotNil(t, testPathWorkspace)
	require.NotNil(t, differentPathWorkspace)
	
	assert.Equal(t, WorkspaceTypePnpm, testPathWorkspace.Type) // Highest priority
	assert.Equal(t, WorkspaceTypeNpm, differentPathWorkspace.Type)
}

func TestWorkspaceParser_HasYarnLock(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Create temporary directory with yarn.lock
	tmpDir, err := os.MkdirTemp("", "yarn-lock-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	yarnLockPath := filepath.Join(tmpDir, "yarn.lock")
	require.NoError(t, os.WriteFile(yarnLockPath, []byte("# yarn.lock"), 0644))
	
	assert.True(t, parser.hasYarnLock(tmpDir))
	
	// Test directory without yarn.lock
	tmpDir2, err := os.MkdirTemp("", "no-yarn-lock-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir2)
	
	assert.False(t, parser.hasYarnLock(tmpDir2))
}

func TestWorkspaceParser_ParseWorkspacePackage(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Create temporary package.json
	tmpDir, err := os.MkdirTemp("", "workspace-package-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	packageJSONContent := `{
		"name": "test-package",
		"version": "1.0.0",
		"private": true,
		"dependencies": {
			"react": "^18.0.0",
			"lodash": "~4.17.0"
		},
		"devDependencies": {
			"typescript": ">=4.0.0"
		},
		"peerDependencies": {
			"react-dom": "^18.0.0"
		}
	}`
	require.NoError(t, os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644))
	
	workspace := &WorkspaceConfiguration{
		Type: WorkspaceTypeNpm,
	}
	
	pkg, err := parser.parseWorkspacePackage(packageJSONPath, workspace)
	require.NoError(t, err)
	require.NotNil(t, pkg)
	
	assert.Equal(t, "test-package", pkg.Name)
	assert.Equal(t, "1.0.0", pkg.Version)
	assert.Equal(t, packageJSONPath, pkg.Path)
	assert.True(t, pkg.Private)
	assert.Equal(t, workspace, pkg.Workspace)
	
	// Check dependencies
	assert.Equal(t, "^18.0.0", pkg.Dependencies["react"])
	assert.Equal(t, "~4.17.0", pkg.Dependencies["lodash"])
	assert.Equal(t, ">=4.0.0", pkg.DevDependencies["typescript"])
	assert.Equal(t, "^18.0.0", pkg.PeerDependencies["react-dom"])
}

func TestWorkspaceParser_DeduplicatePackages(t *testing.T) {
	parser := NewWorkspaceParser()
	
	packages := []*WorkspacePackage{
		{Path: "/path/to/package1"},
		{Path: "/path/to/package2"},
		{Path: "/path/to/package1"}, // Duplicate
		{Path: "/path/to/package3"},
		{Path: "/path/to/package2"}, // Duplicate
	}
	
	unique := parser.deduplicatePackages(packages)
	
	assert.Equal(t, 3, len(unique))
	
	// Check that we have unique paths
	paths := make(map[string]bool)
	for _, pkg := range unique {
		paths[pkg.Path] = true
	}
	
	assert.True(t, paths["/path/to/package1"])
	assert.True(t, paths["/path/to/package2"])
	assert.True(t, paths["/path/to/package3"])
}

func TestWorkspaceParser_CacheManagement(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Test initial cache state
	stats := parser.GetCacheStats()
	assert.Equal(t, 0, stats["workspaces"])
	assert.Equal(t, 0, stats["packages"])
	assert.Equal(t, 0, stats["internedStrings"])
	
	// Add some test data to caches
	parser.cache["/test/path"] = &WorkspaceConfiguration{RootPath: "/test/path"}
	parser.packageCache["/test/package"] = &WorkspacePackage{Path: "/test/package"}
	parser.intern("test-string")
	
	// Check updated stats
	stats = parser.GetCacheStats()
	assert.Equal(t, 1, stats["workspaces"])
	assert.Equal(t, 1, stats["packages"])
	assert.Equal(t, 1, stats["internedStrings"])
	
	// Test cache clearing
	parser.ClearCache()
	stats = parser.GetCacheStats()
	assert.Equal(t, 0, stats["workspaces"])
	assert.Equal(t, 0, stats["packages"])
	assert.Equal(t, 0, stats["internedStrings"])
}

// Integration tests

func TestWorkspaceParser_IntegrationDiscoverWorkspaces(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	parser := NewWorkspaceParser()
	
	// Create a complex workspace structure
	tmpDir, err := os.MkdirTemp("", "integration-workspace-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Create package.json with npm workspaces
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	packageJSONContent := `{
		"name": "test-monorepo",
		"workspaces": ["packages/*"]
	}`
	require.NoError(t, os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644))
	
	// Create pnpm-workspace.yaml (should take priority)
	pnpmWorkspacePath := filepath.Join(tmpDir, "pnpm-workspace.yaml")
	pnpmWorkspaceContent := `packages:
  - 'packages/*'
  - 'apps/*'`
	require.NoError(t, os.WriteFile(pnpmWorkspacePath, []byte(pnpmWorkspaceContent), 0644))
	
	workspaces, err := parser.DiscoverWorkspaces(tmpDir)
	require.NoError(t, err)
	require.Equal(t, 1, len(workspaces))
	
	// Should resolve to pnpm workspace due to higher priority
	workspace := workspaces[0]
	assert.Equal(t, WorkspaceTypePnpm, workspace.Type)
	assert.Equal(t, []string{"packages/*", "apps/*"}, workspace.Packages)
}

func TestWorkspaceParser_IntegrationDiscoverPackages(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	parser := NewWorkspaceParser()
	
	// Create workspace structure with actual packages
	tmpDir, err := os.MkdirTemp("", "integration-packages-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	// Create workspace config
	workspaces := []*WorkspaceConfiguration{
		{
			Type:     WorkspaceTypeNpm,
			RootPath: tmpDir,
			Packages: []string{"packages/*"},
		},
	}
	
	// Create packages directory structure
	packagesDir := filepath.Join(tmpDir, "packages")
	require.NoError(t, os.MkdirAll(packagesDir, 0755))
	
	// Create package1
	pkg1Dir := filepath.Join(packagesDir, "package1")
	require.NoError(t, os.MkdirAll(pkg1Dir, 0755))
	pkg1JSON := `{
		"name": "package1",
		"version": "1.0.0",
		"dependencies": {
			"react": "^18.0.0"
		}
	}`
	require.NoError(t, os.WriteFile(filepath.Join(pkg1Dir, "package.json"), []byte(pkg1JSON), 0644))
	
	// Create package2
	pkg2Dir := filepath.Join(packagesDir, "package2")
	require.NoError(t, os.MkdirAll(pkg2Dir, 0755))
	pkg2JSON := `{
		"name": "package2",
		"version": "2.0.0",
		"dependencies": {
			"lodash": "^4.0.0"
		},
		"devDependencies": {
			"typescript": "^4.0.0"
		}
	}`
	require.NoError(t, os.WriteFile(filepath.Join(pkg2Dir, "package.json"), []byte(pkg2JSON), 0644))
	
	packages, err := parser.DiscoverPackages(workspaces)
	require.NoError(t, err)
	require.Equal(t, 2, len(packages))
	
	// Find packages by name
	var pkg1, pkg2 *WorkspacePackage
	for _, pkg := range packages {
		if pkg.Name == "package1" {
			pkg1 = pkg
		} else if pkg.Name == "package2" {
			pkg2 = pkg
		}
	}
	
	require.NotNil(t, pkg1)
	require.NotNil(t, pkg2)
	
	assert.Equal(t, "1.0.0", pkg1.Version)
	assert.Equal(t, "^18.0.0", pkg1.Dependencies["react"])
	
	assert.Equal(t, "2.0.0", pkg2.Version)
	assert.Equal(t, "^4.0.0", pkg2.Dependencies["lodash"])
	assert.Equal(t, "^4.0.0", pkg2.DevDependencies["typescript"])
}

// Benchmark tests

func BenchmarkWorkspaceParser_StringInterning(b *testing.B) {
	parser := NewWorkspaceParser()
	testStrings := []string{
		"react", "lodash", "typescript", "webpack", "babel",
		"eslint", "jest", "prettier", "husky", "lint-staged",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str := testStrings[i%len(testStrings)]
		parser.intern(str)
	}
}

func BenchmarkWorkspaceParser_InternSlice(b *testing.B) {
	parser := NewWorkspaceParser()
	testSlice := []string{
		"packages/*", "apps/*", "libs/*", "tools/*",
		"packages/*/src", "!**/test/**", "!**/node_modules/**",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.internSlice(testSlice)
	}
}

func BenchmarkWorkspaceParser_InternMap(b *testing.B) {
	parser := NewWorkspaceParser()
	testMap := map[string]string{
		"react":      "^18.0.0",
		"lodash":     "^4.17.0",
		"typescript": "^4.0.0",
		"webpack":    "^5.0.0",
		"babel":      "^7.0.0",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.internMap(testMap)
	}
}

// Test error cases

func TestWorkspaceParser_ErrorHandling(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Test with non-existent directory
	workspaces, err := parser.DiscoverWorkspaces("/non/existent/path")
	assert.Error(t, err)
	assert.Empty(t, workspaces)
	
	// Test with malformed JSON
	tmpDir, err := os.MkdirTemp("", "malformed-json-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	
	malformedJSON := filepath.Join(tmpDir, "package.json")
	require.NoError(t, os.WriteFile(malformedJSON, []byte(`{invalid json`), 0644))
	
	workspace, err := parser.parseNpmWorkspace(malformedJSON, tmpDir)
	assert.Error(t, err)
	assert.Nil(t, workspace)
	
	// Test with malformed YAML
	malformedYAML := filepath.Join(tmpDir, "pnpm-workspace.yaml")
	require.NoError(t, os.WriteFile(malformedYAML, []byte(`packages:\n  - invalid: yaml: syntax`), 0644))
	
	workspace, err = parser.parsePnpmWorkspace(malformedYAML, tmpDir)
	assert.Error(t, err)
	assert.Nil(t, workspace)
}

// Test concurrent access

func TestWorkspaceParser_ConcurrentAccess(t *testing.T) {
	parser := NewWorkspaceParser()
	
	// Test concurrent string interning
	const numGoroutines = 10
	const stringsPerGoroutine = 100
	
	done := make(chan bool)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < stringsPerGoroutine; j++ {
				str := fmt.Sprintf("test-string-%d-%d", id, j)
				parser.intern(str)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// Should have interned all unique strings
	expectedCount := numGoroutines * stringsPerGoroutine
	assert.Equal(t, expectedCount, len(parser.stringInternMap))
}

func TestWorkspaceParser_ConcurrentCacheAccess(t *testing.T) {
	parser := NewWorkspaceParser()
	
	const numGoroutines = 10
	const operationsPerGoroutine = 100
	
	done := make(chan bool)
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				key := fmt.Sprintf("cache-key-%d-%d", id, j)
				value := &WorkspaceConfiguration{RootPath: key}
				
				// Write to cache
				parser.cache[key] = value
				
				// Read from cache
				if cached, exists := parser.cache[key]; exists {
					assert.Equal(t, key, cached.RootPath)
				}
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// Verify cache integrity
	stats := parser.GetCacheStats()
	expectedCount := numGoroutines * operationsPerGoroutine
	assert.Equal(t, expectedCount, stats["workspaces"])
}