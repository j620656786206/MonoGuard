// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains import trace types for Story 3.2.
package types

// ========================================
// Import Trace Types (Story 3.2)
// ========================================

// ImportTrace represents a single import statement that contributes to a cycle.
// Matches @monoguard/types ImportTrace interface.
type ImportTrace struct {
	// FromPackage is the package containing the import
	FromPackage string `json:"fromPackage"`

	// ToPackage is the package being imported
	ToPackage string `json:"toPackage"`

	// FilePath is the relative path to the file containing the import
	FilePath string `json:"filePath"`

	// LineNumber is the 1-based line number of the import statement
	LineNumber int `json:"lineNumber"`

	// Statement is the actual import/require statement text
	Statement string `json:"statement"`

	// ImportType classifies the import style
	ImportType ImportType `json:"importType"`

	// Symbols are the specific imports (empty for namespace/side-effect imports)
	Symbols []string `json:"symbols,omitempty"`
}

// ImportType classifies the import style.
// Matches @monoguard/types ImportType union type.
type ImportType string

const (
	// ImportTypeESMNamed is for named imports: import { foo } from 'bar'
	ImportTypeESMNamed ImportType = "esm-named"

	// ImportTypeESMDefault is for default imports: import foo from 'bar'
	ImportTypeESMDefault ImportType = "esm-default"

	// ImportTypeESMNamespace is for namespace imports: import * as foo from 'bar'
	ImportTypeESMNamespace ImportType = "esm-namespace"

	// ImportTypeESMSideEffect is for side-effect imports: import 'bar'
	ImportTypeESMSideEffect ImportType = "esm-side-effect"

	// ImportTypeESMDynamic is for dynamic imports: import('bar')
	ImportTypeESMDynamic ImportType = "esm-dynamic"

	// ImportTypeCJSRequire is for CommonJS require: require('bar')
	ImportTypeCJSRequire ImportType = "cjs-require"
)
