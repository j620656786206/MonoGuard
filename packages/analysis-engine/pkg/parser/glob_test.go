// Package parser tests for glob pattern matching functionality.
package parser

import (
	"reflect"
	"sort"
	"testing"
)

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		// Simple wildcard *
		{"star matches single segment", "packages/*", "packages/foo", true},
		{"star matches segment with dash", "packages/*", "packages/pkg-a", true},
		{"star matches scoped package", "packages/*", "packages/@mono", true},
		{"star does not match nested", "packages/*", "packages/foo/bar", false},
		{"star at end", "apps/*", "apps/web", true},

		// Double star **
		{"doublestar matches nested", "packages/**", "packages/foo/bar", true},
		{"doublestar matches deep nested", "packages/**", "packages/a/b/c/d", true},
		{"doublestar matches single", "packages/**", "packages/foo", true},

		// Pattern with doublestar in middle
		{"doublestar in middle", "src/**/index.ts", "src/components/Button/index.ts", true},
		{"doublestar in middle direct", "src/**/index.ts", "src/index.ts", true},

		// Question mark ?
		{"question mark single char", "pkg-?", "pkg-a", true},
		{"question mark not match multi", "pkg-?", "pkg-ab", false},

		// Exact match
		{"exact match", "packages/core", "packages/core", true},
		{"exact no match", "packages/core", "packages/web", false},

		// Multiple segments
		{"multiple segments match", "apps/*/src", "apps/web/src", true},
		{"multiple segments no match", "apps/*/src", "apps/web/lib", false},

		// Edge cases
		{"empty pattern", "", "", true},
		{"empty path", "*", "", false},
		{"root only", "/", "/", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchPattern(tt.pattern, tt.path)
			if got != tt.want {
				t.Errorf("MatchPattern(%q, %q) = %v, want %v", tt.pattern, tt.path, got, tt.want)
			}
		})
	}
}

func TestIsNegationPattern(t *testing.T) {
	tests := []struct {
		pattern string
		want    bool
	}{
		{"!packages/deprecated-*", true},
		{"!node_modules", true},
		{"packages/*", false},
		{"!**/test", true},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			got := IsNegationPattern(tt.pattern)
			if got != tt.want {
				t.Errorf("IsNegationPattern(%q) = %v, want %v", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestFilterPaths(t *testing.T) {
	allPaths := []string{
		"packages/pkg-a",
		"packages/pkg-b",
		"packages/deprecated-old",
		"packages/deprecated-legacy",
		"apps/web",
		"apps/mobile",
		"tools/cli",
	}

	tests := []struct {
		name     string
		patterns []string
		want     []string
	}{
		{
			name:     "single include pattern",
			patterns: []string{"packages/*"},
			want:     []string{"packages/deprecated-legacy", "packages/deprecated-old", "packages/pkg-a", "packages/pkg-b"},
		},
		{
			name:     "multiple include patterns",
			patterns: []string{"packages/*", "apps/*"},
			want:     []string{"apps/mobile", "apps/web", "packages/deprecated-legacy", "packages/deprecated-old", "packages/pkg-a", "packages/pkg-b"},
		},
		{
			name:     "include with negation",
			patterns: []string{"packages/*", "!packages/deprecated-*"},
			want:     []string{"packages/pkg-a", "packages/pkg-b"},
		},
		{
			name:     "all with negation",
			patterns: []string{"*/*", "!packages/deprecated-*"},
			want:     []string{"apps/mobile", "apps/web", "packages/pkg-a", "packages/pkg-b", "tools/cli"},
		},
		{
			name:     "empty patterns",
			patterns: []string{},
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterPaths(allPaths, tt.patterns)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandGlobPatternsFromFiles(t *testing.T) {
	// Test expanding glob patterns against a set of known file paths
	// This simulates what would happen when we have a map of files
	// Note: Only directories with package.json are considered packages

	files := map[string][]byte{
		"package.json":                         []byte(`{}`),
		"packages/pkg-a/package.json":          []byte(`{}`),
		"packages/pkg-b/package.json":          []byte(`{}`),
		"packages/nested/sub/package.json":     []byte(`{}`),
		"apps/web/package.json":                []byte(`{}`),
		"packages/deprecated-old/package.json": []byte(`{}`),
	}

	tests := []struct {
		name     string
		patterns []string
		want     []string
	}{
		{
			name:     "packages/* finds direct children with package.json",
			patterns: []string{"packages/*"},
			// Note: packages/nested is NOT included because packages/nested/package.json doesn't exist
			// Only packages/nested/sub has a package.json
			want: []string{"packages/deprecated-old", "packages/pkg-a", "packages/pkg-b"},
		},
		{
			name:     "packages/** finds all nested with package.json",
			patterns: []string{"packages/**"},
			// packages/nested/sub is included because it has package.json
			want: []string{"packages/deprecated-old", "packages/nested/sub", "packages/pkg-a", "packages/pkg-b"},
		},
		{
			name:     "multiple patterns",
			patterns: []string{"packages/*", "apps/*"},
			want:     []string{"apps/web", "packages/deprecated-old", "packages/pkg-a", "packages/pkg-b"},
		},
		{
			name:     "with negation",
			patterns: []string{"packages/*", "!packages/deprecated-*"},
			want:     []string{"packages/pkg-a", "packages/pkg-b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandGlobPatternsFromFiles(files, tt.patterns)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExpandGlobPatternsFromFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
