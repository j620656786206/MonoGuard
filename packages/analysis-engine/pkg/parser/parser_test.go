// Package parser tests for main workspace parsing functionality.
package parser

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

func TestNewParser(t *testing.T) {
	p := NewParser("/workspace")
	if p == nil {
		t.Fatal("NewParser() returned nil")
	}
	if p.rootPath != "/workspace" {
		t.Errorf("rootPath = %q, want /workspace", p.rootPath)
	}
}

func TestDetectWorkspaceType(t *testing.T) {
	tests := []struct {
		name  string
		files map[string][]byte
		want  types.WorkspaceType
	}{
		{
			name: "pnpm workspace detected by pnpm-workspace.yaml",
			files: map[string][]byte{
				"pnpm-workspace.yaml": []byte(`packages: ['packages/*']`),
				"package.json":        []byte(`{"name": "root"}`),
			},
			want: types.WorkspaceTypePnpm,
		},
		{
			name: "yarn workspace detected by yarn.lock",
			files: map[string][]byte{
				"yarn.lock":    []byte(``),
				"package.json": []byte(`{"name": "root", "workspaces": ["packages/*"]}`),
			},
			want: types.WorkspaceTypeYarn,
		},
		{
			name: "npm workspace detected by package-lock.json",
			files: map[string][]byte{
				"package-lock.json": []byte(`{}`),
				"package.json":      []byte(`{"name": "root", "workspaces": ["packages/*"]}`),
			},
			want: types.WorkspaceTypeNpm,
		},
		{
			name: "unknown workspace type",
			files: map[string][]byte{
				"package.json": []byte(`{"name": "root"}`),
			},
			want: types.WorkspaceTypeUnknown,
		},
		{
			name:  "empty files",
			files: map[string][]byte{},
			want:  types.WorkspaceTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser("/workspace")
			got := p.DetectWorkspaceType(tt.files)
			if got != tt.want {
				t.Errorf("DetectWorkspaceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseNpmWorkspace(t *testing.T) {
	files := map[string][]byte{
		"package.json": []byte(`{
			"name": "monorepo-root",
			"workspaces": ["packages/*"]
		}`),
		"package-lock.json":           []byte(`{}`),
		"packages/pkg-a/package.json": []byte(`{
			"name": "@mono/pkg-a",
			"version": "1.0.0",
			"dependencies": {"@mono/pkg-b": "^1.0.0"}
		}`),
		"packages/pkg-b/package.json": []byte(`{
			"name": "@mono/pkg-b",
			"version": "1.0.0",
			"devDependencies": {"typescript": "^5.0.0"}
		}`),
	}

	p := NewParser("/workspace")
	result, err := p.Parse(files)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify workspace type
	if result.WorkspaceType != types.WorkspaceTypeNpm {
		t.Errorf("WorkspaceType = %v, want npm", result.WorkspaceType)
	}

	// Verify root path
	if result.RootPath != "/workspace" {
		t.Errorf("RootPath = %q, want /workspace", result.RootPath)
	}

	// Verify packages count
	if len(result.Packages) != 2 {
		t.Errorf("Packages count = %d, want 2", len(result.Packages))
	}

	// Verify pkg-a
	pkgA, ok := result.Packages["@mono/pkg-a"]
	if !ok {
		t.Fatal("Missing @mono/pkg-a package")
	}
	if pkgA.Name != "@mono/pkg-a" {
		t.Errorf("pkg-a Name = %q, want @mono/pkg-a", pkgA.Name)
	}
	if pkgA.Version != "1.0.0" {
		t.Errorf("pkg-a Version = %q, want 1.0.0", pkgA.Version)
	}
	if pkgA.Dependencies["@mono/pkg-b"] != "^1.0.0" {
		t.Errorf("pkg-a dependency @mono/pkg-b = %q, want ^1.0.0", pkgA.Dependencies["@mono/pkg-b"])
	}
}

func TestParsePnpmWorkspaceIntegration(t *testing.T) {
	files := map[string][]byte{
		"pnpm-workspace.yaml": []byte(`packages:
  - 'packages/*'
  - 'apps/*'
`),
		"package.json":                []byte(`{"name": "monorepo-root"}`),
		"packages/core/package.json":  []byte(`{"name": "@mono/core", "version": "2.0.0"}`),
		"packages/utils/package.json": []byte(`{"name": "@mono/utils", "version": "1.0.0"}`),
		"apps/web/package.json": []byte(`{
			"name": "@mono/web",
			"version": "0.0.1",
			"dependencies": {"@mono/core": "^2.0.0", "@mono/utils": "^1.0.0"}
		}`),
	}

	p := NewParser("/workspace")
	result, err := p.Parse(files)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if result.WorkspaceType != types.WorkspaceTypePnpm {
		t.Errorf("WorkspaceType = %v, want pnpm", result.WorkspaceType)
	}

	if len(result.Packages) != 3 {
		t.Errorf("Packages count = %d, want 3", len(result.Packages))
	}
}

func TestParseWithNegation(t *testing.T) {
	files := map[string][]byte{
		"package.json": []byte(`{
			"name": "monorepo-root",
			"workspaces": ["packages/*", "!packages/deprecated-*"]
		}`),
		"package-lock.json":                  []byte(`{}`),
		"packages/active/package.json":       []byte(`{"name": "@mono/active", "version": "1.0.0"}`),
		"packages/deprecated-old/package.json": []byte(`{"name": "@mono/deprecated-old", "version": "0.0.1"}`),
	}

	p := NewParser("/workspace")
	result, err := p.Parse(files)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Only @mono/active should be included
	if len(result.Packages) != 1 {
		t.Errorf("Packages count = %d, want 1 (deprecated should be excluded)", len(result.Packages))
	}

	if _, ok := result.Packages["@mono/active"]; !ok {
		t.Error("Missing @mono/active package")
	}

	if _, ok := result.Packages["@mono/deprecated-old"]; ok {
		t.Error("@mono/deprecated-old should be excluded by negation pattern")
	}
}

func TestParseInvalidWorkspace(t *testing.T) {
	tests := []struct {
		name    string
		files   map[string][]byte
		wantErr bool
	}{
		{
			name:    "empty files",
			files:   map[string][]byte{},
			wantErr: true,
		},
		{
			name: "no root package.json",
			files: map[string][]byte{
				"packages/pkg-a/package.json": []byte(`{"name": "pkg-a"}`),
			},
			wantErr: true,
		},
		{
			name: "invalid JSON in root package.json",
			files: map[string][]byte{
				"package.json": []byte(`{invalid json}`),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser("/workspace")
			_, err := p.Parse(tt.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
