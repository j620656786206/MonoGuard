package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock external resolver for testing
type mockExternalResolver struct {
	packages map[string]*PackageMetadata
	versions map[string][]*SemanticVersion
}

func newMockExternalResolver() *mockExternalResolver {
	return &mockExternalResolver{
		packages: make(map[string]*PackageMetadata),
		versions: make(map[string][]*SemanticVersion),
	}
}

func (m *mockExternalResolver) ResolvePackageVersion(name string, vRange *VersionRange) (*SemanticVersion, error) {
	versions := m.versions[name]
	if versions == nil {
		return nil, fmt.Errorf("package %s not found", name)
	}
	
	// Return the latest version for simplicity
	if len(versions) > 0 {
		return versions[len(versions)-1], nil
	}
	
	return nil, fmt.Errorf("no versions found for %s", name)
}

func (m *mockExternalResolver) GetPackageMetadata(name, version string) (*PackageMetadata, error) {
	key := name + "@" + version
	if metadata, exists := m.packages[key]; exists {
		return metadata, nil
	}
	return nil, fmt.Errorf("metadata not found for %s@%s", name, version)
}

func (m *mockExternalResolver) PackageExists(name, version string) (bool, error) {
	key := name + "@" + version
	_, exists := m.packages[key]
	return exists, nil
}

func (m *mockExternalResolver) GetAvailableVersions(name string) ([]*SemanticVersion, error) {
	versions := m.versions[name]
	if versions == nil {
		return nil, fmt.Errorf("package %s not found", name)
	}
	return versions, nil
}

func (m *mockExternalResolver) GetPackageDependencies(name, version string) (map[string]*VersionRange, error) {
	metadata, err := m.GetPackageMetadata(name, version)
	if err != nil {
		return nil, err
	}
	
	deps := make(map[string]*VersionRange)
	for depName, versionRange := range metadata.Dependencies {
		vRange, _ := parseVersionRange(versionRange)
		deps[depName] = vRange
	}
	
	return deps, nil
}

func (m *mockExternalResolver) addPackage(name, version string, deps map[string]string) {
	key := name + "@" + version
	m.packages[key] = &PackageMetadata{
		Name:         name,
		Version:      version,
		Dependencies: deps,
		PublishedAt:  time.Now(),
	}
	
	// Add version to versions list
	versionObj := &Version{
		Raw:   version,
		Major: extractMajorVersion(version),
		Minor: extractMinorVersion(version),
		Patch: extractPatchVersion(version),
	}
	
	m.versions[name] = append(m.versions[name], versionObj)
}

func parseVersionRange(versionRange string) (*VersionRange, error) {
	// Simplified version range parsing for tests
	return &VersionRange{
		Raw:      versionRange,
		Operator: "^",
	}, nil
}

func extractMajorVersion(version string) int {
	// Simplified version extraction for tests
	return 1
}

func extractMinorVersion(version string) int {
	return 0
}

func extractPatchVersion(version string) int {
	return 0
}

func createTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests
	return logger
}

func createTestPackages() []*PackageInfo {
	return []*PackageInfo{
		{
			Path: "/test/package-a/package.json",
			PackageJSON: PackageJSON{
				Name:    "package-a",
				Version: "1.0.0",
				Dependencies: map[string]string{
					"lodash":  "^4.17.21",
					"express": "^4.18.0",
				},
				DevDependencies: map[string]string{
					"jest": "^29.0.0",
				},
			},
		},
		{
			Path: "/test/package-b/package.json",
			PackageJSON: PackageJSON{
				Name:    "package-b",
				Version: "1.0.0",
				Dependencies: map[string]string{
					"lodash":  "^4.17.20", // Different version
					"moment": "^2.29.0",
				},
			},
		},
		{
			Path: "/test/package-c/package.json",
			PackageJSON: PackageJSON{
				Name:    "package-c",
				Version: "1.0.0",
				Dependencies: map[string]string{
					"package-a": "workspace:*",
					"lodash":    "^3.10.1", // Major version conflict
				},
			},
		},
	}
}

func TestDependencyTreeResolver_BuildTree(t *testing.T) {
	tests := []struct {
		name           string
		packages       []*PackageInfo
		options        BuildOptions
		setupMocks     func(*mockExternalResolver)
		expectedNodes  int
		expectedErrors bool
	}{
		{
			name:     "simple_tree",
			packages: createTestPackages(),
			options: BuildOptions{
				MaxDepth:        5,
				IncludeDevDeps:  false,
				PreferWorkspace: true,
				UseNpmRegistry:  true,
			},
			setupMocks: func(resolver *mockExternalResolver) {
				resolver.addPackage("lodash", "4.17.21", map[string]string{})
				resolver.addPackage("express", "4.18.0", map[string]string{
					"accepts": "^1.3.8",
				})
				resolver.addPackage("moment", "2.29.0", map[string]string{})
				resolver.addPackage("accepts", "1.3.8", map[string]string{})
			},
			expectedNodes:  3, // Workspace packages
			expectedErrors: false,
		},
		{
			name:     "with_dev_dependencies",
			packages: createTestPackages(),
			options: BuildOptions{
				MaxDepth:        3,
				IncludeDevDeps:  true,
				PreferWorkspace: true,
				UseNpmRegistry:  true,
			},
			setupMocks: func(resolver *mockExternalResolver) {
				resolver.addPackage("lodash", "4.17.21", map[string]string{})
				resolver.addPackage("jest", "29.0.0", map[string]string{})
			},
			expectedNodes:  3,
			expectedErrors: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := createTestLogger()
			resolver := NewDependencyTreeResolver(logger)
			
			// Replace external resolver with mock
			mockResolver := newMockExternalResolver()
			if tt.setupMocks != nil {
				tt.setupMocks(mockResolver)
			}
			
			// Create a new resolver with mock
			dtr := &DependencyTreeResolver{
				logger:           logger,
				externalResolver: mockResolver,
				cache:           NewResolutionCache(logger),
			}

			tree, err := dtr.BuildTree(context.Background(), tt.packages, tt.options)

			if tt.expectedErrors {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, tree)
				assert.GreaterOrEqual(t, len(tree.AllNodes), tt.expectedNodes)
				assert.NotNil(t, tree.Metadata)
				assert.Greater(t, tree.Metadata.TotalNodes, 0)
			}
		})
	}
}

func TestDependencyTreeResolver_ConflictDetection(t *testing.T) {
	logger := createTestLogger()
	resolver := NewDependencyTreeResolver(logger)
	mockResolver := newMockExternalResolver()
	
	// Setup mock data with conflicts
	mockResolver.addPackage("lodash", "4.17.21", map[string]string{})
	mockResolver.addPackage("lodash", "4.17.20", map[string]string{})
	mockResolver.addPackage("lodash", "3.10.1", map[string]string{})
	
	resolver.(*DependencyTreeResolver).externalResolver = mockResolver

	packages := createTestPackages()
	options := BuildOptions{
		MaxDepth:        3,
		UseNpmRegistry:  true,
		PreferWorkspace: true,
	}

	tree, err := resolver.BuildTree(context.Background(), packages, options)
	require.NoError(t, err)
	require.NotNil(t, tree)

	// Should detect version conflicts
	assert.Greater(t, len(tree.Conflicts), 0, "Expected to find version conflicts")
	
	// Check for lodash conflicts specifically
	var lodashConflict *EnhancedVersionConflict
	for _, conflict := range tree.Conflicts {
		if conflict.PackageName == "lodash" {
			lodashConflict = conflict
			break
		}
	}
	
	require.NotNil(t, lodashConflict, "Expected to find lodash version conflict")
	assert.Equal(t, ConflictTypeMajorVersion, lodashConflict.ConflictType)
	assert.Equal(t, SeverityCritical, lodashConflict.Severity)
	assert.Greater(t, len(lodashConflict.ResolutionOptions), 0)
}

func TestVersionRangeParsing(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		major    int
		minor    int
		patch    int
	}{
		{"^4.17.21", "^", 4, 17, 21},
		{"~2.29.0", "~", 2, 29, 0},
		{">=1.0.0", ">=", 1, 0, 0},
		{"1.5.3", "=", 1, 5, 3},
		{"<2.0.0", "<", 2, 0, 0},
	}

	logger := createTestLogger()
	resolver := NewDependencyTreeResolver(logger)
	dtr := resolver.(*DependencyTreeResolver)

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			vRange, err := dtr.parseVersionRange(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.operator, vRange.Operator)
			assert.Equal(t, tt.major, vRange.Major)
			assert.Equal(t, tt.minor, vRange.Minor)
			assert.Equal(t, tt.patch, vRange.Patch)
		})
	}
}

func TestVersionComparison(t *testing.T) {
	logger := createTestLogger()
	resolver := NewDependencyTreeResolver(logger)
	dtr := resolver.(*DependencyTreeResolver)

	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0", "1.0.0-alpha", 1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			result := dtr.compareVersions(tt.v1, tt.v2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolutionCache(t *testing.T) {
	logger := createTestLogger()
	cache := NewResolutionCache(logger)

	// Test basic get/set
	key := "test-key"
	node := &TreeNode{
		Name: "test-package",
		ResolvedVersion: &Version{
			Raw:   "1.0.0",
			Major: 1,
			Minor: 0,
			Patch: 0,
		},
	}

	// Should miss initially
	_, found := cache.Get(key)
	assert.False(t, found)

	// Set and retrieve
	cache.Set(key, node, "input-hash")
	retrieved, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, node.Name, retrieved.Name)

	// Test stats
	stats := cache.GetStats()
	assert.Equal(t, int64(1), stats.Misses)
	assert.Equal(t, int64(1), stats.L1Hits) // L1 hit when retrieved
}

func TestExternalResolver(t *testing.T) {
	logger := createTestLogger()
	
	// Test mock resolver
	mockResolver := newMockExternalResolver()
	mockResolver.addPackage("lodash", "4.17.21", map[string]string{
		"core-js": "^2.6.12",
	})
	
	// Test version resolution
	vRange := &VersionRange{
		Raw:      "^4.17.0",
		Operator: "^",
		Major:    4,
		Minor:    17,
		Patch:    0,
	}
	
	version, err := mockResolver.ResolvePackageVersion("lodash", vRange)
	require.NoError(t, err)
	assert.Equal(t, "4.17.21", version.Raw)

	// Test metadata retrieval
	metadata, err := mockResolver.GetPackageMetadata("lodash", "4.17.21")
	require.NoError(t, err)
	assert.Equal(t, "lodash", metadata.Name)
	assert.Equal(t, "4.17.21", metadata.Version)
	assert.Contains(t, metadata.Dependencies, "core-js")

	// Test package existence
	exists, err := mockResolver.PackageExists("lodash", "4.17.21")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = mockResolver.PackageExists("nonexistent", "1.0.0")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestConflictSeverityCalculation(t *testing.T) {
	logger := createTestLogger()
	resolver := NewDependencyTreeResolver(logger)
	dtr := resolver.(*DependencyTreeResolver)

	tests := []struct {
		versions []string
		expected ConflictSeverity
	}{
		{[]string{"1.0.0", "2.0.0"}, SeverityCritical},  // Major version conflict
		{[]string{"1.0.0", "1.3.0"}, SeverityMedium},    // Minor version conflict
		{[]string{"1.0.0", "1.0.5"}, SeverityInfo},      // Patch version conflict
		{[]string{"1.0.0", "1.0.0"}, SeverityInfo},      // No conflict
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.versions), func(t *testing.T) {
			severity := dtr.calculateConflictSeverity(tt.versions)
			assert.Equal(t, tt.expected, severity)
		})
	}
}

func TestBuildOptions(t *testing.T) {
	logger := createTestLogger()
	resolver := NewDependencyTreeResolver(logger)
	packages := createTestPackages()
	
	mockResolver := newMockExternalResolver()
	mockResolver.addPackage("lodash", "4.17.21", map[string]string{})
	mockResolver.addPackage("jest", "29.0.0", map[string]string{})
	resolver.(*DependencyTreeResolver).externalResolver = mockResolver

	// Test max depth limiting
	options := BuildOptions{
		MaxDepth:       1, // Very shallow
		UseNpmRegistry: true,
	}

	tree, err := resolver.BuildTree(context.Background(), packages, options)
	require.NoError(t, err)
	assert.LessOrEqual(t, tree.Metadata.MaxDepth, 1)

	// Test dev dependency inclusion
	options = BuildOptions{
		MaxDepth:       5,
		IncludeDevDeps: true,
		UseNpmRegistry: true,
	}

	tree, err = resolver.BuildTree(context.Background(), packages, options)
	require.NoError(t, err)
	
	// Should have more nodes when including dev dependencies
	assert.Greater(t, tree.Metadata.TotalNodes, 3)
}

func BenchmarkTreeResolution(b *testing.B) {
	logger := createTestLogger()
	resolver := NewDependencyTreeResolver(logger)
	packages := createTestPackages()
	
	mockResolver := newMockExternalResolver()
	mockResolver.addPackage("lodash", "4.17.21", map[string]string{})
	mockResolver.addPackage("express", "4.18.0", map[string]string{})
	mockResolver.addPackage("moment", "2.29.0", map[string]string{})
	resolver.(*DependencyTreeResolver).externalResolver = mockResolver
	
	options := BuildOptions{
		MaxDepth:       5,
		UseNpmRegistry: true,
		EnableCaching:  true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := resolver.BuildTree(context.Background(), packages, options)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func TestCircularDependencyDetection(t *testing.T) {
	logger := createTestLogger()
	resolver := NewDependencyTreeResolver(logger)
	
	// Create packages with circular dependencies
	circularPackages := []*PackageInfo{
		{
			Path: "/test/pkg-a/package.json",
			PackageJSON: PackageJSON{
				Name:    "pkg-a",
				Version: "1.0.0",
				Dependencies: map[string]string{
					"pkg-b": "workspace:*",
				},
			},
		},
		{
			Path: "/test/pkg-b/package.json",
			PackageJSON: PackageJSON{
				Name:    "pkg-b",
				Version: "1.0.0",
				Dependencies: map[string]string{
					"pkg-a": "workspace:*",
				},
			},
		},
	}
	
	mockResolver := newMockExternalResolver()
	resolver.(*DependencyTreeResolver).externalResolver = mockResolver
	
	options := BuildOptions{
		MaxDepth:        5,
		PreferWorkspace: true,
		UseNpmRegistry:  false, // Don't fetch external packages for this test
	}

	tree, err := resolver.BuildTree(context.Background(), circularPackages, options)
	require.NoError(t, err)
	require.NotNil(t, tree)
	
	// The tree should be built successfully even with circular dependencies
	// The resolver should detect and handle the circular dependency gracefully
	assert.Greater(t, len(tree.AllNodes), 0)
}