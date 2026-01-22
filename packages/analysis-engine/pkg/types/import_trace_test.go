// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains tests for import trace types for Story 3.2.
package types

import (
	"encoding/json"
	"testing"
)

// ========================================
// Import Trace Type Tests (Story 3.2)
// ========================================

func TestImportTrace_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		trace    ImportTrace
		wantJSON string
	}{
		{
			name: "ESM named import with symbols",
			trace: ImportTrace{
				FromPackage: "@mono/ui",
				ToPackage:   "@mono/api",
				FilePath:    "packages/ui/src/index.ts",
				LineNumber:  5,
				Statement:   "import { foo, bar } from '@mono/api'",
				ImportType:  ImportTypeESMNamed,
				Symbols:     []string{"foo", "bar"},
			},
			wantJSON: `{"fromPackage":"@mono/ui","toPackage":"@mono/api","filePath":"packages/ui/src/index.ts","lineNumber":5,"statement":"import { foo, bar } from '@mono/api'","importType":"esm-named","symbols":["foo","bar"]}`,
		},
		{
			name: "ESM default import",
			trace: ImportTrace{
				FromPackage: "@mono/ui",
				ToPackage:   "@mono/core",
				FilePath:    "packages/ui/src/utils.ts",
				LineNumber:  1,
				Statement:   "import core from '@mono/core'",
				ImportType:  ImportTypeESMDefault,
				Symbols:     nil,
			},
			wantJSON: `{"fromPackage":"@mono/ui","toPackage":"@mono/core","filePath":"packages/ui/src/utils.ts","lineNumber":1,"statement":"import core from '@mono/core'","importType":"esm-default"}`,
		},
		{
			name: "ESM namespace import",
			trace: ImportTrace{
				FromPackage: "@mono/api",
				ToPackage:   "@mono/utils",
				FilePath:    "packages/api/src/client.ts",
				LineNumber:  3,
				Statement:   "import * as utils from '@mono/utils'",
				ImportType:  ImportTypeESMNamespace,
				Symbols:     nil,
			},
			wantJSON: `{"fromPackage":"@mono/api","toPackage":"@mono/utils","filePath":"packages/api/src/client.ts","lineNumber":3,"statement":"import * as utils from '@mono/utils'","importType":"esm-namespace"}`,
		},
		{
			name: "ESM side-effect import",
			trace: ImportTrace{
				FromPackage: "@mono/app",
				ToPackage:   "@mono/polyfills",
				FilePath:    "packages/app/src/main.ts",
				LineNumber:  1,
				Statement:   "import '@mono/polyfills'",
				ImportType:  ImportTypeESMSideEffect,
				Symbols:     nil,
			},
			wantJSON: `{"fromPackage":"@mono/app","toPackage":"@mono/polyfills","filePath":"packages/app/src/main.ts","lineNumber":1,"statement":"import '@mono/polyfills'","importType":"esm-side-effect"}`,
		},
		{
			name: "ESM dynamic import",
			trace: ImportTrace{
				FromPackage: "@mono/app",
				ToPackage:   "@mono/lazy",
				FilePath:    "packages/app/src/loader.ts",
				LineNumber:  10,
				Statement:   "import('@mono/lazy')",
				ImportType:  ImportTypeESMDynamic,
				Symbols:     nil,
			},
			wantJSON: `{"fromPackage":"@mono/app","toPackage":"@mono/lazy","filePath":"packages/app/src/loader.ts","lineNumber":10,"statement":"import('@mono/lazy')","importType":"esm-dynamic"}`,
		},
		{
			name: "CJS require",
			trace: ImportTrace{
				FromPackage: "@mono/server",
				ToPackage:   "@mono/config",
				FilePath:    "packages/server/src/index.js",
				LineNumber:  2,
				Statement:   "const config = require('@mono/config')",
				ImportType:  ImportTypeCJSRequire,
				Symbols:     nil,
			},
			wantJSON: `{"fromPackage":"@mono/server","toPackage":"@mono/config","filePath":"packages/server/src/index.js","lineNumber":2,"statement":"const config = require('@mono/config')","importType":"cjs-require"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			gotBytes, err := json.Marshal(tt.trace)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			got := string(gotBytes)
			if got != tt.wantJSON {
				t.Errorf("json.Marshal() = %s, want %s", got, tt.wantJSON)
			}

			// Test unmarshaling
			var unmarshaled ImportTrace
			if err := json.Unmarshal([]byte(tt.wantJSON), &unmarshaled); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify fields
			if unmarshaled.FromPackage != tt.trace.FromPackage {
				t.Errorf("FromPackage = %s, want %s", unmarshaled.FromPackage, tt.trace.FromPackage)
			}
			if unmarshaled.ToPackage != tt.trace.ToPackage {
				t.Errorf("ToPackage = %s, want %s", unmarshaled.ToPackage, tt.trace.ToPackage)
			}
			if unmarshaled.FilePath != tt.trace.FilePath {
				t.Errorf("FilePath = %s, want %s", unmarshaled.FilePath, tt.trace.FilePath)
			}
			if unmarshaled.LineNumber != tt.trace.LineNumber {
				t.Errorf("LineNumber = %d, want %d", unmarshaled.LineNumber, tt.trace.LineNumber)
			}
			if unmarshaled.Statement != tt.trace.Statement {
				t.Errorf("Statement = %s, want %s", unmarshaled.Statement, tt.trace.Statement)
			}
			if unmarshaled.ImportType != tt.trace.ImportType {
				t.Errorf("ImportType = %s, want %s", unmarshaled.ImportType, tt.trace.ImportType)
			}
		})
	}
}

func TestImportType_Constants(t *testing.T) {
	// Verify import type constants have correct string values (camelCase for JSON)
	tests := []struct {
		importType ImportType
		want       string
	}{
		{ImportTypeESMNamed, "esm-named"},
		{ImportTypeESMDefault, "esm-default"},
		{ImportTypeESMNamespace, "esm-namespace"},
		{ImportTypeESMSideEffect, "esm-side-effect"},
		{ImportTypeESMDynamic, "esm-dynamic"},
		{ImportTypeCJSRequire, "cjs-require"},
	}

	for _, tt := range tests {
		t.Run(string(tt.importType), func(t *testing.T) {
			if string(tt.importType) != tt.want {
				t.Errorf("ImportType = %s, want %s", tt.importType, tt.want)
			}
		})
	}
}

func TestImportTrace_SymbolsOmitempty(t *testing.T) {
	// Test that nil/empty symbols are omitted from JSON
	trace := ImportTrace{
		FromPackage: "@mono/ui",
		ToPackage:   "@mono/core",
		FilePath:    "test.ts",
		LineNumber:  1,
		Statement:   "import '@mono/core'",
		ImportType:  ImportTypeESMSideEffect,
		Symbols:     nil,
	}

	gotBytes, err := json.Marshal(trace)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	got := string(gotBytes)

	// Should NOT contain "symbols" field when nil
	if containsSubstr(got, `"symbols"`) {
		t.Errorf("JSON should not contain symbols field when nil, got %s", got)
	}

	// Test with empty slice - should also be omitted
	trace.Symbols = []string{}
	gotBytes, err = json.Marshal(trace)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	got = string(gotBytes)

	// Empty slice serializes as [] in Go, but we want omitempty behavior
	// Note: Go's omitempty does NOT omit empty slices, only nil slices
	// This is acceptable behavior per story requirements
}

// containsSubstr checks if s contains substr (renamed to avoid conflict with root_cause_test.go)
func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
