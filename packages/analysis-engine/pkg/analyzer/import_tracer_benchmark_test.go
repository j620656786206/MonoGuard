// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains benchmark tests for import tracing performance (Story 3.2).
package analyzer

import (
	"fmt"
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Import Tracer Benchmarks (Story 3.2 AC7)
// Goal: Add <5% overhead to existing analysis time
// ========================================

// BenchmarkImportTracer_SmallWorkspace benchmarks tracing for small workspaces (5 packages).
func BenchmarkImportTracer_SmallWorkspace(b *testing.B) {
	workspace, files := generateTestWorkspace(5)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/pkg-0", "@mono/pkg-1", "@mono/pkg-2", "@mono/pkg-0"},
	}

	tracer := NewImportTracer(workspace, files)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracer.Trace(cycle)
	}
}

// BenchmarkImportTracer_MediumWorkspace benchmarks tracing for medium workspaces (20 packages).
func BenchmarkImportTracer_MediumWorkspace(b *testing.B) {
	workspace, files := generateTestWorkspace(20)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/pkg-0", "@mono/pkg-1", "@mono/pkg-2", "@mono/pkg-3", "@mono/pkg-0"},
	}

	tracer := NewImportTracer(workspace, files)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracer.Trace(cycle)
	}
}

// BenchmarkImportTracer_LargeWorkspace benchmarks tracing for large workspaces (50 packages).
func BenchmarkImportTracer_LargeWorkspace(b *testing.B) {
	workspace, files := generateTestWorkspace(50)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/pkg-0", "@mono/pkg-1", "@mono/pkg-2", "@mono/pkg-3", "@mono/pkg-4", "@mono/pkg-0"},
	}

	tracer := NewImportTracer(workspace, files)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracer.Trace(cycle)
	}
}

// BenchmarkImportTracer_ManyFilesPerPackage benchmarks with many files per package.
func BenchmarkImportTracer_ManyFilesPerPackage(b *testing.B) {
	workspace, files := generateTestWorkspaceWithManyFiles(10, 50) // 10 packages, 50 files each
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/pkg-0", "@mono/pkg-1", "@mono/pkg-2", "@mono/pkg-0"},
	}

	tracer := NewImportTracer(workspace, files)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracer.Trace(cycle)
	}
}

// BenchmarkAnalyzeWithSources_Overhead measures import tracing overhead.
func BenchmarkAnalyzeWithSources_Overhead(b *testing.B) {
	workspace := generateCycleWorkspace()

	// Measure without sources
	b.Run("WithoutSources", func(b *testing.B) {
		analyzer := NewAnalyzer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = analyzer.Analyze(workspace)
		}
	})

	// Measure with sources
	b.Run("WithSources", func(b *testing.B) {
		analyzer := NewAnalyzer()
		files := generateSourceFiles(workspace)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = analyzer.AnalyzeWithSources(workspace, files)
		}
	})
}

// generateTestWorkspace creates a test workspace with specified number of packages.
func generateTestWorkspace(numPackages int) (*types.WorkspaceData, map[string][]byte) {
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      make(map[string]*types.PackageInfo),
	}

	files := make(map[string][]byte)

	for i := 0; i < numPackages; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		pkgPath := fmt.Sprintf("packages/pkg-%d", i)

		workspace.Packages[pkgName] = &types.PackageInfo{
			Name:    pkgName,
			Version: "1.0.0",
			Path:    pkgPath,
		}

		// Generate source files with imports to next package
		nextPkg := fmt.Sprintf("@mono/pkg-%d", (i+1)%numPackages)
		content := fmt.Sprintf(`import { something } from '%s';
import type { Type } from '%s';
const other = require('%s');
`, nextPkg, nextPkg, nextPkg)

		filePath := fmt.Sprintf("%s/src/index.ts", pkgPath)
		files[filePath] = []byte(content)
	}

	return workspace, files
}

// generateTestWorkspaceWithManyFiles creates a workspace with many files per package.
func generateTestWorkspaceWithManyFiles(numPackages, filesPerPackage int) (*types.WorkspaceData, map[string][]byte) {
	workspace := &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages:      make(map[string]*types.PackageInfo),
	}

	files := make(map[string][]byte)

	for i := 0; i < numPackages; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		pkgPath := fmt.Sprintf("packages/pkg-%d", i)

		workspace.Packages[pkgName] = &types.PackageInfo{
			Name:    pkgName,
			Version: "1.0.0",
			Path:    pkgPath,
		}

		nextPkg := fmt.Sprintf("@mono/pkg-%d", (i+1)%numPackages)

		for j := 0; j < filesPerPackage; j++ {
			content := fmt.Sprintf(`import { component%d } from '%s';
export const Component%d = () => {};
`, j, nextPkg, j)

			filePath := fmt.Sprintf("%s/src/Component%d.tsx", pkgPath, j)
			files[filePath] = []byte(content)
		}
	}

	return workspace, files
}

// generateCycleWorkspace creates a workspace with a cycle for overhead testing.
func generateCycleWorkspace() *types.WorkspaceData {
	return &types.WorkspaceData{
		RootPath:      "/workspace",
		WorkspaceType: types.WorkspaceTypePnpm,
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {
				Name:    "@mono/ui",
				Version: "1.0.0",
				Path:    "packages/ui",
				Dependencies: map[string]string{
					"@mono/api": "^1.0.0",
				},
			},
			"@mono/api": {
				Name:    "@mono/api",
				Version: "1.0.0",
				Path:    "packages/api",
				Dependencies: map[string]string{
					"@mono/core": "^1.0.0",
				},
			},
			"@mono/core": {
				Name:    "@mono/core",
				Version: "1.0.0",
				Path:    "packages/core",
				Dependencies: map[string]string{
					"@mono/ui": "^1.0.0", // Creates cycle
				},
			},
		},
	}
}

// generateSourceFiles creates source files for a workspace.
func generateSourceFiles(workspace *types.WorkspaceData) map[string][]byte {
	files := make(map[string][]byte)

	// Create source files based on dependencies
	for pkgName, pkg := range workspace.Packages {
		for depName := range pkg.Dependencies {
			content := fmt.Sprintf(`import { dep } from '%s';
export const use%s = () => dep;
`, depName, pkgName)
			filePath := fmt.Sprintf("%s/src/index.ts", pkg.Path)
			files[filePath] = []byte(content)
		}
	}

	return files
}
