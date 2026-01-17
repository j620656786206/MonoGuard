// Package parser tests for pnpm-workspace.yaml parsing functionality.
package parser

import (
	"reflect"
	"testing"
)

func TestParsePnpmWorkspace(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *PnpmWorkspace
		wantErr bool
	}{
		{
			name: "valid pnpm-workspace.yaml with packages",
			input: []byte(`packages:
  - 'packages/*'
  - 'apps/*'
  - 'tools/*'
`),
			want: &PnpmWorkspace{
				Packages: []string{"packages/*", "apps/*", "tools/*"},
			},
			wantErr: false,
		},
		{
			name: "pnpm-workspace.yaml with negation patterns",
			input: []byte(`packages:
  - 'packages/*'
  - '!packages/experimental-*'
  - '!packages/deprecated-*'
`),
			want: &PnpmWorkspace{
				Packages: []string{"packages/*", "!packages/experimental-*", "!packages/deprecated-*"},
			},
			wantErr: false,
		},
		{
			name: "pnpm-workspace.yaml with glob patterns",
			input: []byte(`packages:
  - 'packages/**'
  - 'libs/shared-*'
`),
			want: &PnpmWorkspace{
				Packages: []string{"packages/**", "libs/shared-*"},
			},
			wantErr: false,
		},
		{
			name: "empty packages array",
			input: []byte(`packages: []
`),
			want: &PnpmWorkspace{
				Packages: []string{},
			},
			wantErr: false,
		},
		{
			name: "no packages field",
			input: []byte(`# empty config
`),
			want: &PnpmWorkspace{
				Packages: nil,
			},
			wantErr: false,
		},
		{
			name:    "invalid YAML",
			input:   []byte(`packages: [invalid yaml`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   []byte(``),
			want:    &PnpmWorkspace{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePnpmWorkspace(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePnpmWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePnpmWorkspace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPnpmWorkspacePackagePatterns(t *testing.T) {
	// Test realistic pnpm-workspace.yaml configurations
	input := []byte(`packages:
  # All packages in the packages directory
  - 'packages/*'
  # All apps
  - 'apps/*'
  # Exclude experimental packages
  - '!packages/experimental-*'
`)

	ws, err := ParsePnpmWorkspace(input)
	if err != nil {
		t.Fatalf("ParsePnpmWorkspace() error = %v", err)
	}

	if len(ws.Packages) != 3 {
		t.Errorf("Expected 3 patterns, got %d", len(ws.Packages))
	}

	// Verify specific patterns
	expectedPatterns := []string{"packages/*", "apps/*", "!packages/experimental-*"}
	for i, expected := range expectedPatterns {
		if i >= len(ws.Packages) {
			t.Errorf("Missing pattern at index %d", i)
			continue
		}
		if ws.Packages[i] != expected {
			t.Errorf("Pattern[%d] = %q, want %q", i, ws.Packages[i], expected)
		}
	}
}

func TestPnpmWorkspaceWithCatalog(t *testing.T) {
	// Test that we can parse pnpm-workspace.yaml with catalog feature (pnpm 9+)
	input := []byte(`packages:
  - 'packages/*'
catalog:
  react: ^18.2.0
  typescript: ^5.0.0
`)

	ws, err := ParsePnpmWorkspace(input)
	if err != nil {
		t.Fatalf("ParsePnpmWorkspace() error = %v", err)
	}

	// We only care about packages field for this story
	if len(ws.Packages) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(ws.Packages))
	}
	if ws.Packages[0] != "packages/*" {
		t.Errorf("Pattern = %q, want packages/*", ws.Packages[0])
	}
}
