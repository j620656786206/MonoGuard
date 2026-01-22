// Package parser provides workspace configuration parsing for monorepos.
// This file contains tests for import statement parsing for Story 3.2.
package parser

import (
	"testing"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Import Parser Tests (Story 3.2)
// ========================================

func TestImportParser_ParseFile_ESMNamed(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/api": true, "@mono/core": true}

	tests := []struct {
		name        string
		content     string
		wantCount   int
		wantType    types.ImportType
		wantPackage string
		wantSymbols []string
	}{
		{
			name:        "single named import",
			content:     `import { foo } from '@mono/api';`,
			wantCount:   1,
			wantType:    types.ImportTypeESMNamed,
			wantPackage: "@mono/api",
			wantSymbols: []string{"foo"},
		},
		{
			name:        "multiple named imports",
			content:     `import { foo, bar, baz } from '@mono/api';`,
			wantCount:   1,
			wantType:    types.ImportTypeESMNamed,
			wantPackage: "@mono/api",
			wantSymbols: []string{"foo", "bar", "baz"},
		},
		{
			name:        "named import with alias",
			content:     `import { foo as f, bar } from '@mono/api';`,
			wantCount:   1,
			wantType:    types.ImportTypeESMNamed,
			wantPackage: "@mono/api",
			wantSymbols: []string{"foo as f", "bar"},
		},
		{
			name:        "double quoted import",
			content:     `import { foo } from "@mono/api";`,
			wantCount:   1,
			wantType:    types.ImportTypeESMNamed,
			wantPackage: "@mono/api",
		},
		{
			name:        "non-target package ignored",
			content:     `import { foo } from 'lodash';`,
			wantCount:   0,
			wantType:    "",
			wantPackage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traces := parser.ParseFile([]byte(tt.content), "test.ts", targets)

			if len(traces) != tt.wantCount {
				t.Errorf("ParseFile() returned %d traces, want %d", len(traces), tt.wantCount)
				return
			}

			if tt.wantCount > 0 {
				trace := traces[0]
				if trace.ImportType != tt.wantType {
					t.Errorf("ImportType = %s, want %s", trace.ImportType, tt.wantType)
				}
				if trace.ToPackage != tt.wantPackage {
					t.Errorf("ToPackage = %s, want %s", trace.ToPackage, tt.wantPackage)
				}
				if tt.wantSymbols != nil && len(trace.Symbols) != len(tt.wantSymbols) {
					t.Errorf("Symbols count = %d, want %d", len(trace.Symbols), len(tt.wantSymbols))
				}
			}
		})
	}
}

func TestImportParser_ParseFile_ESMDefault(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/core": true}

	tests := []struct {
		name        string
		content     string
		wantCount   int
		wantPackage string
	}{
		{
			name:        "default import",
			content:     `import core from '@mono/core';`,
			wantCount:   1,
			wantPackage: "@mono/core",
		},
		{
			name:        "default import double quotes",
			content:     `import core from "@mono/core";`,
			wantCount:   1,
			wantPackage: "@mono/core",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traces := parser.ParseFile([]byte(tt.content), "test.ts", targets)

			if len(traces) != tt.wantCount {
				t.Errorf("ParseFile() returned %d traces, want %d", len(traces), tt.wantCount)
				return
			}

			if tt.wantCount > 0 {
				trace := traces[0]
				if trace.ImportType != types.ImportTypeESMDefault {
					t.Errorf("ImportType = %s, want %s", trace.ImportType, types.ImportTypeESMDefault)
				}
				if trace.ToPackage != tt.wantPackage {
					t.Errorf("ToPackage = %s, want %s", trace.ToPackage, tt.wantPackage)
				}
			}
		})
	}
}

func TestImportParser_ParseFile_ESMNamespace(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/utils": true}

	content := `import * as utils from '@mono/utils';`
	traces := parser.ParseFile([]byte(content), "test.ts", targets)

	if len(traces) != 1 {
		t.Fatalf("ParseFile() returned %d traces, want 1", len(traces))
	}

	trace := traces[0]
	if trace.ImportType != types.ImportTypeESMNamespace {
		t.Errorf("ImportType = %s, want %s", trace.ImportType, types.ImportTypeESMNamespace)
	}
	if trace.ToPackage != "@mono/utils" {
		t.Errorf("ToPackage = %s, want @mono/utils", trace.ToPackage)
	}
}

func TestImportParser_ParseFile_ESMSideEffect(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/polyfills": true}

	content := `import '@mono/polyfills';`
	traces := parser.ParseFile([]byte(content), "test.ts", targets)

	if len(traces) != 1 {
		t.Fatalf("ParseFile() returned %d traces, want 1", len(traces))
	}

	trace := traces[0]
	if trace.ImportType != types.ImportTypeESMSideEffect {
		t.Errorf("ImportType = %s, want %s", trace.ImportType, types.ImportTypeESMSideEffect)
	}
	if trace.ToPackage != "@mono/polyfills" {
		t.Errorf("ToPackage = %s, want @mono/polyfills", trace.ToPackage)
	}
}

func TestImportParser_ParseFile_ESMDynamic(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/lazy": true}

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "dynamic import",
			content: `const mod = import('@mono/lazy');`,
		},
		{
			name:    "await dynamic import",
			content: `const mod = await import('@mono/lazy');`,
		},
		{
			name:    "dynamic import in function",
			content: `async function load() { return import('@mono/lazy'); }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traces := parser.ParseFile([]byte(tt.content), "test.ts", targets)

			if len(traces) != 1 {
				t.Fatalf("ParseFile() returned %d traces, want 1", len(traces))
			}

			trace := traces[0]
			if trace.ImportType != types.ImportTypeESMDynamic {
				t.Errorf("ImportType = %s, want %s", trace.ImportType, types.ImportTypeESMDynamic)
			}
			if trace.ToPackage != "@mono/lazy" {
				t.Errorf("ToPackage = %s, want @mono/lazy", trace.ToPackage)
			}
		})
	}
}

func TestImportParser_ParseFile_CJSRequire(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/config": true}

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "standard require",
			content: `const config = require('@mono/config');`,
		},
		{
			name:    "destructured require",
			content: `const { foo } = require('@mono/config');`,
		},
		{
			name:    "direct require",
			content: `require('@mono/config');`,
		},
		{
			name:    "double quotes",
			content: `const config = require("@mono/config");`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traces := parser.ParseFile([]byte(tt.content), "test.js", targets)

			if len(traces) != 1 {
				t.Fatalf("ParseFile() returned %d traces, want 1", len(traces))
			}

			trace := traces[0]
			if trace.ImportType != types.ImportTypeCJSRequire {
				t.Errorf("ImportType = %s, want %s", trace.ImportType, types.ImportTypeCJSRequire)
			}
			if trace.ToPackage != "@mono/config" {
				t.Errorf("ToPackage = %s, want @mono/config", trace.ToPackage)
			}
		})
	}
}

func TestImportParser_ParseFile_MultipleImports(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{
		"@mono/api":  true,
		"@mono/core": true,
		"@mono/ui":   true,
	}

	content := `
import { api } from '@mono/api';
import core from '@mono/core';
import * as ui from '@mono/ui';
import { something } from 'lodash'; // Should be ignored
`
	traces := parser.ParseFile([]byte(content), "test.ts", targets)

	if len(traces) != 3 {
		t.Errorf("ParseFile() returned %d traces, want 3", len(traces))
	}

	// Verify all three imports were captured
	packages := make(map[string]bool)
	for _, trace := range traces {
		packages[trace.ToPackage] = true
	}

	expected := []string{"@mono/api", "@mono/core", "@mono/ui"}
	for _, pkg := range expected {
		if !packages[pkg] {
			t.Errorf("Missing import for package %s", pkg)
		}
	}
}

func TestImportParser_ParseFile_RelativeImportsIgnored(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"./local": true, "../parent": true}

	content := `
import { foo } from './local';
import bar from '../parent';
import baz from './nested/module';
`
	traces := parser.ParseFile([]byte(content), "test.ts", targets)

	// Relative imports should be filtered out
	if len(traces) != 0 {
		t.Errorf("ParseFile() returned %d traces for relative imports, want 0", len(traces))
	}
}

func TestImportParser_ParseFile_LineNumbers(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/api": true, "@mono/core": true}

	content := `// Comment line 1
// Comment line 2
import { api } from '@mono/api';
// Comment line 4
import core from '@mono/core';`

	traces := parser.ParseFile([]byte(content), "test.ts", targets)

	if len(traces) != 2 {
		t.Fatalf("ParseFile() returned %d traces, want 2", len(traces))
	}

	// First import should be on line 3
	if traces[0].LineNumber != 3 {
		t.Errorf("First import LineNumber = %d, want 3", traces[0].LineNumber)
	}

	// Second import should be on line 5
	if traces[1].LineNumber != 5 {
		t.Errorf("Second import LineNumber = %d, want 5", traces[1].LineNumber)
	}
}

func TestImportParser_ParseFile_Statement(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/api": true}

	content := `import { foo, bar } from '@mono/api';`
	traces := parser.ParseFile([]byte(content), "test.ts", targets)

	if len(traces) != 1 {
		t.Fatalf("ParseFile() returned %d traces, want 1", len(traces))
	}

	// Statement should contain the actual import text
	if traces[0].Statement == "" {
		t.Error("Statement should not be empty")
	}
}

func TestImportParser_ParseFile_FilePath(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/api": true}

	content := `import { foo } from '@mono/api';`
	filePath := "packages/ui/src/components/Button.tsx"
	traces := parser.ParseFile([]byte(content), filePath, targets)

	if len(traces) != 1 {
		t.Fatalf("ParseFile() returned %d traces, want 1", len(traces))
	}

	if traces[0].FilePath != filePath {
		t.Errorf("FilePath = %s, want %s", traces[0].FilePath, filePath)
	}
}

func TestImportParser_ParseFile_ReExports(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/core": true}

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "named re-export",
			content: `export { foo } from '@mono/core';`,
		},
		{
			name:    "star re-export",
			content: `export * from '@mono/core';`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traces := parser.ParseFile([]byte(tt.content), "test.ts", targets)

			// Re-exports should be treated as imports (they create dependencies)
			if len(traces) != 1 {
				t.Fatalf("ParseFile() returned %d traces, want 1", len(traces))
			}

			if traces[0].ToPackage != "@mono/core" {
				t.Errorf("ToPackage = %s, want @mono/core", traces[0].ToPackage)
			}
		})
	}
}

func TestExtractPackageName(t *testing.T) {
	tests := []struct {
		importPath string
		want       string
	}{
		{"@scope/pkg", "@scope/pkg"},
		{"@scope/pkg/sub", "@scope/pkg"},
		{"@scope/pkg/sub/deep", "@scope/pkg"},
		{"lodash", "lodash"},
		{"lodash/debounce", "lodash"},
		{"./local", ""},
		{"../parent", ""},
		{"../../up/up", ""},
	}

	for _, tt := range tests {
		t.Run(tt.importPath, func(t *testing.T) {
			got := ExtractPackageName(tt.importPath)
			if got != tt.want {
				t.Errorf("ExtractPackageName(%q) = %q, want %q", tt.importPath, got, tt.want)
			}
		})
	}
}

// TestImportParser_ParseFile_MultiLineImports verifies parsing of multi-line imports (Story 3.2 edge case).
func TestImportParser_ParseFile_MultiLineImports(t *testing.T) {
	parser := NewImportParser()
	targets := map[string]bool{"@mono/core": true}

	tests := []struct {
		name        string
		content     string
		wantCount   int
		wantSymbols []string
	}{
		{
			name: "multi-line named import",
			content: `import {
  foo,
  bar,
  baz
} from '@mono/core';`,
			wantCount:   1,
			wantSymbols: []string{"foo", "bar", "baz"},
		},
		{
			name: "multi-line with aliases",
			content: `import {
  foo as f,
  bar as b,
  default as Baz
} from '@mono/core';`,
			wantCount:   1,
			wantSymbols: []string{"foo as f", "bar as b", "default as Baz"},
		},
		{
			name: "multi-line re-export",
			content: `export {
  Component,
  utils
} from '@mono/core';`,
			wantCount: 1,
		},
		{
			name: "mixed single and multi-line",
			content: `import { single } from '@mono/core';
import {
  multi1,
  multi2
} from '@mono/core';`,
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traces := parser.ParseFile([]byte(tt.content), "test.ts", targets)

			if len(traces) != tt.wantCount {
				t.Errorf("ParseFile() returned %d traces, want %d", len(traces), tt.wantCount)
				return
			}

			// Verify symbols for first trace if expected
			if tt.wantCount > 0 && tt.wantSymbols != nil {
				if len(traces[0].Symbols) != len(tt.wantSymbols) {
					t.Errorf("Symbols count = %d, want %d. Got: %v", len(traces[0].Symbols), len(tt.wantSymbols), traces[0].Symbols)
				}
			}
		})
	}
}
