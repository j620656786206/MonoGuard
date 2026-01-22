// Package parser provides workspace configuration parsing for monorepos.
// This file contains import statement parsing for Story 3.2.
package parser

import (
	"regexp"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Import Parser (Story 3.2)
// ========================================

// ImportParser extracts import statements from source files.
type ImportParser struct {
	// Regex patterns for different import types
	esmNamedPattern      *regexp.Regexp
	esmDefaultPattern    *regexp.Regexp
	esmNamespacePattern  *regexp.Regexp
	esmSideEffectPattern *regexp.Regexp
	esmDynamicPattern    *regexp.Regexp
	cjsRequirePattern    *regexp.Regexp
	// Re-export patterns (create dependencies too)
	reExportNamedPattern *regexp.Regexp
	reExportStarPattern  *regexp.Regexp
}

// NewImportParser creates a new parser with compiled regex patterns.
func NewImportParser() *ImportParser {
	return &ImportParser{
		// ESM Named: import { foo, bar } from 'package'
		// Captures: group(1) = imports, group(2) = package
		esmNamedPattern: regexp.MustCompile(`import\s*\{([^}]+)\}\s*from\s*['"]([^'"]+)['"]`),

		// ESM Default: import foo from 'package'
		// Must NOT match "import * as" or "import {"
		// Captures: group(1) = name, group(2) = package
		esmDefaultPattern: regexp.MustCompile(`import\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s+from\s*['"]([^'"]+)['"]`),

		// ESM Namespace: import * as foo from 'package'
		// Captures: group(1) = alias, group(2) = package
		esmNamespacePattern: regexp.MustCompile(`import\s*\*\s*as\s+(\w+)\s+from\s*['"]([^'"]+)['"]`),

		// ESM Side-effect: import 'package' (standalone line, no identifiers)
		// Captures: group(1) = package
		esmSideEffectPattern: regexp.MustCompile(`(?m)^\s*import\s*['"]([^'"]+)['"]\s*;?\s*$`),

		// ESM Dynamic: import('package') or await import('package')
		// Captures: group(1) = package
		esmDynamicPattern: regexp.MustCompile(`import\s*\(\s*['"]([^'"]+)['"]\s*\)`),

		// CJS Require: require('package')
		// Captures: group(1) = package
		cjsRequirePattern: regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`),

		// Re-export Named: export { foo } from 'package'
		// Captures: group(1) = exports, group(2) = package
		reExportNamedPattern: regexp.MustCompile(`export\s*\{([^}]+)\}\s*from\s*['"]([^'"]+)['"]`),

		// Re-export Star: export * from 'package'
		// Captures: group(1) = package
		reExportStarPattern: regexp.MustCompile(`export\s*\*\s*from\s*['"]([^'"]+)['"]`),
	}
}

// ParseFile extracts all imports from a source file.
// Returns imports that reference the specified target packages.
func (ip *ImportParser) ParseFile(content []byte, filePath string, targetPackages map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace
	contentStr := string(content)

	// Parse all import types
	traces = append(traces, ip.parseESMNamedImports(contentStr, filePath, targetPackages)...)
	traces = append(traces, ip.parseESMNamespaceImports(contentStr, filePath, targetPackages)...)
	traces = append(traces, ip.parseESMDefaultImports(contentStr, filePath, targetPackages)...)
	traces = append(traces, ip.parseESMSideEffectImports(contentStr, filePath, targetPackages)...)
	traces = append(traces, ip.parseESMDynamicImports(contentStr, filePath, targetPackages)...)
	traces = append(traces, ip.parseCJSRequires(contentStr, filePath, targetPackages)...)
	traces = append(traces, ip.parseReExports(contentStr, filePath, targetPackages)...)

	return traces
}

// parseESMNamedImports extracts ESM named import statements.
func (ip *ImportParser) parseESMNamedImports(content string, filePath string, targets map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace

	matches := ip.esmNamedPattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		imports := content[match[2]:match[3]]
		packagePath := content[match[4]:match[5]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := content[match[0]:match[1]]
		lineNumber := getLineNumber(content, match[0])
		symbols := parseSymbols(imports)

		traces = append(traces, types.ImportTrace{
			FromPackage: "", // Will be filled by ImportTracer
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeESMNamed,
			Symbols:     symbols,
		})
	}

	return traces
}

// parseESMDefaultImports extracts ESM default import statements.
func (ip *ImportParser) parseESMDefaultImports(content string, filePath string, targets map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace

	matches := ip.esmDefaultPattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		// Check if this is actually a namespace import (starts with "* as")
		// or a named import (starts with "{")
		beforeMatch := ""
		if match[0] > 7 {
			beforeMatch = content[match[0]-7 : match[0]]
		}
		if strings.Contains(beforeMatch, "* as") || strings.Contains(beforeMatch, "{") {
			continue
		}

		packagePath := content[match[4]:match[5]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := content[match[0]:match[1]]
		lineNumber := getLineNumber(content, match[0])

		traces = append(traces, types.ImportTrace{
			FromPackage: "",
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeESMDefault,
			Symbols:     nil,
		})
	}

	return traces
}

// parseESMNamespaceImports extracts ESM namespace import statements.
func (ip *ImportParser) parseESMNamespaceImports(content string, filePath string, targets map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace

	matches := ip.esmNamespacePattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		packagePath := content[match[4]:match[5]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := content[match[0]:match[1]]
		lineNumber := getLineNumber(content, match[0])

		traces = append(traces, types.ImportTrace{
			FromPackage: "",
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeESMNamespace,
			Symbols:     nil,
		})
	}

	return traces
}

// parseESMSideEffectImports extracts ESM side-effect import statements.
func (ip *ImportParser) parseESMSideEffectImports(content string, filePath string, targets map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace

	matches := ip.esmSideEffectPattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		packagePath := content[match[2]:match[3]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := strings.TrimSpace(content[match[0]:match[1]])
		lineNumber := getLineNumber(content, match[0])

		traces = append(traces, types.ImportTrace{
			FromPackage: "",
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeESMSideEffect,
			Symbols:     nil,
		})
	}

	return traces
}

// parseESMDynamicImports extracts ESM dynamic import expressions.
func (ip *ImportParser) parseESMDynamicImports(content string, filePath string, targets map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace

	matches := ip.esmDynamicPattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		packagePath := content[match[2]:match[3]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := content[match[0]:match[1]]
		lineNumber := getLineNumber(content, match[0])

		traces = append(traces, types.ImportTrace{
			FromPackage: "",
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeESMDynamic,
			Symbols:     nil,
		})
	}

	return traces
}

// parseCJSRequires extracts CommonJS require statements.
func (ip *ImportParser) parseCJSRequires(content string, filePath string, targets map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace

	matches := ip.cjsRequirePattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		packagePath := content[match[2]:match[3]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := content[match[0]:match[1]]
		lineNumber := getLineNumber(content, match[0])

		traces = append(traces, types.ImportTrace{
			FromPackage: "",
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeCJSRequire,
			Symbols:     nil,
		})
	}

	return traces
}

// parseReExports extracts re-export statements (they create dependencies too).
func (ip *ImportParser) parseReExports(content string, filePath string, targets map[string]bool) []types.ImportTrace {
	var traces []types.ImportTrace

	// Named re-exports: export { foo } from 'package'
	matches := ip.reExportNamedPattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		packagePath := content[match[4]:match[5]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := content[match[0]:match[1]]
		lineNumber := getLineNumber(content, match[0])
		exports := content[match[2]:match[3]]
		symbols := parseSymbols(exports)

		traces = append(traces, types.ImportTrace{
			FromPackage: "",
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeESMNamed,
			Symbols:     symbols,
		})
	}

	// Star re-exports: export * from 'package'
	matches = ip.reExportStarPattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		packagePath := content[match[2]:match[3]]
		packageName := ExtractPackageName(packagePath)

		if packageName == "" || !targets[packageName] {
			continue
		}

		statement := content[match[0]:match[1]]
		lineNumber := getLineNumber(content, match[0])

		traces = append(traces, types.ImportTrace{
			FromPackage: "",
			ToPackage:   packageName,
			FilePath:    filePath,
			LineNumber:  lineNumber,
			Statement:   statement,
			ImportType:  types.ImportTypeESMNamespace,
			Symbols:     nil,
		})
	}

	return traces
}

// getLineNumber returns the 1-based line number for a byte offset in content.
func getLineNumber(content string, offset int) int {
	lines := 1
	for i := 0; i < offset && i < len(content); i++ {
		if content[i] == '\n' {
			lines++
		}
	}
	return lines
}

// parseSymbols extracts individual symbols from an import clause.
// e.g., "foo, bar as b, baz" -> ["foo", "bar as b", "baz"]
func parseSymbols(imports string) []string {
	var symbols []string
	parts := strings.Split(imports, ",")
	for _, part := range parts {
		symbol := strings.TrimSpace(part)
		if symbol != "" {
			symbols = append(symbols, symbol)
		}
	}
	return symbols
}

// ExtractPackageName extracts the package name from an import path.
// Handles scoped packages (@scope/pkg) and subpath imports (pkg/submodule).
// Returns empty string for relative imports.
//
// Examples:
//
//	'@scope/pkg'       → '@scope/pkg'
//	'@scope/pkg/sub'   → '@scope/pkg'
//	'lodash'           → 'lodash'
//	'lodash/debounce'  → 'lodash'
//	'./local'          → '' (relative import, skip)
//	'../parent'        → '' (relative import, skip)
func ExtractPackageName(importPath string) string {
	// Skip relative imports
	if strings.HasPrefix(importPath, ".") {
		return ""
	}

	// Handle scoped packages (@scope/pkg)
	if strings.HasPrefix(importPath, "@") {
		parts := strings.SplitN(importPath, "/", 3)
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
		return importPath
	}

	// Handle regular packages
	parts := strings.SplitN(importPath, "/", 2)
	return parts[0]
}
