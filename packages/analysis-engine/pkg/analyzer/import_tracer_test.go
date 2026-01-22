// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains tests for import tracing for Story 3.2.
package analyzer

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Import Tracer Tests (Story 3.2)
// ========================================

func TestNewImportTracer(t *testing.T) {
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {Name: "@mono/ui", Path: "packages/ui"},
		},
	}
	files := map[string][]byte{
		"packages/ui/src/index.ts": []byte(`import { api } from '@mono/api';`),
	}

	tracer := NewImportTracer(workspace, files)

	if tracer == nil {
		t.Fatal("NewImportTracer() returned nil")
	}
	if tracer.workspace != workspace {
		t.Error("workspace not set correctly")
	}
	if tracer.files == nil {
		t.Error("files not set correctly")
	}
}

func TestImportTracer_Trace_BasicCycle(t *testing.T) {
	// Create a simple cycle: @mono/ui → @mono/api → @mono/core → @mono/ui
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":   {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/api":  {Name: "@mono/api", Path: "packages/api"},
			"@mono/core": {Name: "@mono/core", Path: "packages/core"},
		},
	}

	files := map[string][]byte{
		"packages/ui/src/index.ts":   []byte(`import { api } from '@mono/api';`),
		"packages/api/src/client.ts": []byte(`import { core } from '@mono/core';`),
		"packages/core/src/utils.ts": []byte(`import { ui } from '@mono/ui';`),
	}

	tracer := NewImportTracer(workspace, files)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/core", "@mono/ui"},
	}

	traces := tracer.Trace(cycle)

	// Should find 3 traces (one for each edge in the cycle)
	if len(traces) != 3 {
		t.Fatalf("Trace() returned %d traces, want 3", len(traces))
	}

	// Verify traces are in cycle order
	expectedFrom := []string{"@mono/ui", "@mono/api", "@mono/core"}
	expectedTo := []string{"@mono/api", "@mono/core", "@mono/ui"}

	for i, trace := range traces {
		if trace.FromPackage != expectedFrom[i] {
			t.Errorf("traces[%d].FromPackage = %s, want %s", i, trace.FromPackage, expectedFrom[i])
		}
		if trace.ToPackage != expectedTo[i] {
			t.Errorf("traces[%d].ToPackage = %s, want %s", i, trace.ToPackage, expectedTo[i])
		}
	}
}

func TestImportTracer_Trace_DirectCycle(t *testing.T) {
	// Create a direct cycle: @mono/a → @mono/b → @mono/a
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/a": {Name: "@mono/a", Path: "packages/a"},
			"@mono/b": {Name: "@mono/b", Path: "packages/b"},
		},
	}

	files := map[string][]byte{
		"packages/a/src/index.ts": []byte(`import { b } from '@mono/b';`),
		"packages/b/src/index.ts": []byte(`import { a } from '@mono/a';`),
	}

	tracer := NewImportTracer(workspace, files)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/a", "@mono/b", "@mono/a"},
	}

	traces := tracer.Trace(cycle)

	if len(traces) != 2 {
		t.Fatalf("Trace() returned %d traces, want 2", len(traces))
	}

	// Verify traces
	if traces[0].FromPackage != "@mono/a" || traces[0].ToPackage != "@mono/b" {
		t.Errorf("First trace incorrect: %s → %s", traces[0].FromPackage, traces[0].ToPackage)
	}
	if traces[1].FromPackage != "@mono/b" || traces[1].ToPackage != "@mono/a" {
		t.Errorf("Second trace incorrect: %s → %s", traces[1].FromPackage, traces[1].ToPackage)
	}
}

func TestImportTracer_Trace_EmptyFiles(t *testing.T) {
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/api": {Name: "@mono/api", Path: "packages/api"},
		},
	}

	// Empty files map - graceful degradation
	files := map[string][]byte{}

	tracer := NewImportTracer(workspace, files)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}

	traces := tracer.Trace(cycle)

	// Should return empty slice, not nil
	if traces == nil {
		t.Error("Trace() returned nil, want empty slice")
	}
	if len(traces) != 0 {
		t.Errorf("Trace() returned %d traces, want 0 for empty files", len(traces))
	}
}

func TestImportTracer_Trace_NilCycle(t *testing.T) {
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/ui": {Name: "@mono/ui", Path: "packages/ui"},
		},
	}
	files := map[string][]byte{
		"packages/ui/src/index.ts": []byte(`import { api } from '@mono/api';`),
	}

	tracer := NewImportTracer(workspace, files)
	traces := tracer.Trace(nil)

	// Should return empty slice for nil cycle
	if traces == nil {
		t.Error("Trace() returned nil, want empty slice")
	}
	if len(traces) != 0 {
		t.Errorf("Trace() returned %d traces, want 0 for nil cycle", len(traces))
	}
}

func TestImportTracer_Trace_MultipleImportsPerEdge(t *testing.T) {
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/api": {Name: "@mono/api", Path: "packages/api"},
		},
	}

	// Multiple files importing @mono/api
	files := map[string][]byte{
		"packages/ui/src/index.ts":     []byte(`import { api } from '@mono/api';`),
		"packages/ui/src/Button.tsx":   []byte(`import { useApi } from '@mono/api';`),
		"packages/api/src/index.ts":    []byte(`import { ui } from '@mono/ui';`),
	}

	tracer := NewImportTracer(workspace, files)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}

	traces := tracer.Trace(cycle)

	// Should find traces for both edges
	// Edge 1: @mono/ui → @mono/api (might have multiple traces from different files)
	// Edge 2: @mono/api → @mono/ui
	if len(traces) < 2 {
		t.Errorf("Trace() returned %d traces, want at least 2", len(traces))
	}

	// Verify we have traces for both edges
	hasUIToAPI := false
	hasAPIToUI := false
	for _, trace := range traces {
		if trace.FromPackage == "@mono/ui" && trace.ToPackage == "@mono/api" {
			hasUIToAPI = true
		}
		if trace.FromPackage == "@mono/api" && trace.ToPackage == "@mono/ui" {
			hasAPIToUI = true
		}
	}

	if !hasUIToAPI {
		t.Error("Missing trace for @mono/ui → @mono/api")
	}
	if !hasAPIToUI {
		t.Error("Missing trace for @mono/api → @mono/ui")
	}
}

func TestImportTracer_Trace_NestedPackagePaths(t *testing.T) {
	// Test with nested package paths
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Path: "packages/frontend/ui"},
			"@mono/api": {Name: "@mono/api", Path: "packages/backend/api"},
		},
	}

	files := map[string][]byte{
		"packages/frontend/ui/src/index.ts": []byte(`import { api } from '@mono/api';`),
		"packages/backend/api/src/index.ts": []byte(`import { ui } from '@mono/ui';`),
	}

	tracer := NewImportTracer(workspace, files)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}

	traces := tracer.Trace(cycle)

	if len(traces) != 2 {
		t.Fatalf("Trace() returned %d traces, want 2", len(traces))
	}
}

func TestImportTracer_Trace_LineNumbersAndStatements(t *testing.T) {
	workspace := &types.WorkspaceData{
		Packages: map[string]*types.PackageInfo{
			"@mono/ui":  {Name: "@mono/ui", Path: "packages/ui"},
			"@mono/api": {Name: "@mono/api", Path: "packages/api"},
		},
	}

	files := map[string][]byte{
		"packages/ui/src/index.ts": []byte(`// Header comment
import React from 'react';
import { api } from '@mono/api';

export const Component = () => {};
`),
		"packages/api/src/index.ts": []byte(`import { ui } from '@mono/ui';`),
	}

	tracer := NewImportTracer(workspace, files)
	cycle := &types.CircularDependencyInfo{
		Cycle: []string{"@mono/ui", "@mono/api", "@mono/ui"},
	}

	traces := tracer.Trace(cycle)

	if len(traces) != 2 {
		t.Fatalf("Trace() returned %d traces, want 2", len(traces))
	}

	// Find the @mono/ui → @mono/api trace
	var uiToApiTrace *types.ImportTrace
	for i := range traces {
		if traces[i].FromPackage == "@mono/ui" && traces[i].ToPackage == "@mono/api" {
			uiToApiTrace = &traces[i]
			break
		}
	}

	if uiToApiTrace == nil {
		t.Fatal("Could not find @mono/ui → @mono/api trace")
	}

	// Verify line number (should be line 3)
	if uiToApiTrace.LineNumber != 3 {
		t.Errorf("LineNumber = %d, want 3", uiToApiTrace.LineNumber)
	}

	// Verify statement contains the import
	if uiToApiTrace.Statement == "" {
		t.Error("Statement should not be empty")
	}
}

func TestIsSourceFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"index.ts", true},
		{"Component.tsx", true},
		{"script.js", true},
		{"App.jsx", true},
		{"module.mjs", true},
		{"module.cjs", true},
		{"styles.css", false},
		{"data.json", false},
		{"package.json", false},
		{"README.md", false},
		{"image.png", false},
		{"src/index.ts", true},
		{"packages/ui/src/Button.tsx", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := IsSourceFile(tt.path)
			if got != tt.want {
				t.Errorf("IsSourceFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
