// Package parser tests for package.json parsing functionality.
package parser

import (
	"reflect"
	"testing"
)

func TestParsePackageJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *PackageJSON
		wantErr bool
	}{
		{
			name: "valid package.json with all fields",
			input: []byte(`{
				"name": "@mono/pkg-a",
				"version": "1.0.0",
				"dependencies": {"lodash": "^4.17.21"},
				"devDependencies": {"typescript": "^5.0.0"},
				"peerDependencies": {"react": "^18.0.0"},
				"workspaces": ["packages/*"]
			}`),
			want: &PackageJSON{
				Name:             "@mono/pkg-a",
				Version:          "1.0.0",
				Dependencies:     map[string]string{"lodash": "^4.17.21"},
				DevDependencies:  map[string]string{"typescript": "^5.0.0"},
				PeerDependencies: map[string]string{"react": "^18.0.0"},
			},
			wantErr: false,
		},
		{
			name: "minimal package.json",
			input: []byte(`{
				"name": "@mono/simple",
				"version": "0.0.1"
			}`),
			want: &PackageJSON{
				Name:    "@mono/simple",
				Version: "0.0.1",
			},
			wantErr: false,
		},
		{
			name: "package.json without name",
			input: []byte(`{
				"version": "1.0.0"
			}`),
			want: &PackageJSON{
				Version: "1.0.0",
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid json`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte(``),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePackageJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePackageJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name != tt.want.Name {
					t.Errorf("Name = %q, want %q", got.Name, tt.want.Name)
				}
				if got.Version != tt.want.Version {
					t.Errorf("Version = %q, want %q", got.Version, tt.want.Version)
				}
			}
		})
	}
}

func TestExtractWorkspacePatterns(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []string
		wantErr bool
	}{
		{
			name: "array format workspaces",
			input: []byte(`{
				"name": "root",
				"workspaces": ["packages/*", "apps/*"]
			}`),
			want:    []string{"packages/*", "apps/*"},
			wantErr: false,
		},
		{
			name: "object format workspaces with packages",
			input: []byte(`{
				"name": "root",
				"workspaces": {
					"packages": ["packages/*", "libs/*"],
					"nohoist": ["**/react-native"]
				}
			}`),
			want:    []string{"packages/*", "libs/*"},
			wantErr: false,
		},
		{
			name: "no workspaces field",
			input: []byte(`{
				"name": "simple-package"
			}`),
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty workspaces array",
			input: []byte(`{
				"name": "root",
				"workspaces": []
			}`),
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   []byte(`{invalid`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := ParsePackageJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePackageJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got, err := ExtractWorkspacePatterns(pkg)
			if err != nil {
				t.Errorf("ExtractWorkspacePatterns() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractWorkspacePatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractWorkspacePatternsObjectFormat(t *testing.T) {
	// Test specifically the object format with packages and nohoist
	input := []byte(`{
		"name": "monorepo-root",
		"workspaces": {
			"packages": ["packages/*", "apps/*", "tools/*"],
			"nohoist": ["**/react-native", "**/react-native/**"]
		}
	}`)

	pkg, err := ParsePackageJSON(input)
	if err != nil {
		t.Fatalf("ParsePackageJSON() error = %v", err)
	}

	patterns, err := ExtractWorkspacePatterns(pkg)
	if err != nil {
		t.Fatalf("ExtractWorkspacePatterns() error = %v", err)
	}

	expected := []string{"packages/*", "apps/*", "tools/*"}
	if !reflect.DeepEqual(patterns, expected) {
		t.Errorf("patterns = %v, want %v", patterns, expected)
	}
}

func TestParsePackageJSONDependencies(t *testing.T) {
	input := []byte(`{
		"name": "@mono/complex",
		"version": "2.0.0",
		"dependencies": {
			"@mono/core": "^1.0.0",
			"lodash": "^4.17.21",
			"axios": "~1.0.0"
		},
		"devDependencies": {
			"typescript": "^5.0.0",
			"vitest": "^1.0.0",
			"@types/node": "^20.0.0"
		},
		"peerDependencies": {
			"react": "^17.0.0 || ^18.0.0",
			"react-dom": "^17.0.0 || ^18.0.0"
		}
	}`)

	pkg, err := ParsePackageJSON(input)
	if err != nil {
		t.Fatalf("ParsePackageJSON() error = %v", err)
	}

	// Verify dependencies
	if len(pkg.Dependencies) != 3 {
		t.Errorf("Dependencies count = %d, want 3", len(pkg.Dependencies))
	}
	if pkg.Dependencies["@mono/core"] != "^1.0.0" {
		t.Errorf("@mono/core version = %q, want ^1.0.0", pkg.Dependencies["@mono/core"])
	}

	// Verify devDependencies
	if len(pkg.DevDependencies) != 3 {
		t.Errorf("DevDependencies count = %d, want 3", len(pkg.DevDependencies))
	}
	if pkg.DevDependencies["typescript"] != "^5.0.0" {
		t.Errorf("typescript version = %q, want ^5.0.0", pkg.DevDependencies["typescript"])
	}

	// Verify peerDependencies
	if len(pkg.PeerDependencies) != 2 {
		t.Errorf("PeerDependencies count = %d, want 2", len(pkg.PeerDependencies))
	}
	if pkg.PeerDependencies["react"] != "^17.0.0 || ^18.0.0" {
		t.Errorf("react version = %q, want ^17.0.0 || ^18.0.0", pkg.PeerDependencies["react"])
	}
}
