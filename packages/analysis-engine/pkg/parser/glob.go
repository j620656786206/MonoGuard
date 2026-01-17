// Package parser provides glob pattern matching for workspace configurations.
package parser

import (
	"path/filepath"
	"strings"
)

// MatchPattern checks if a path matches a glob pattern.
// Supports:
//   - * matches any sequence of non-separator characters
//   - ** matches any sequence of characters including separators
//   - ? matches any single non-separator character
func MatchPattern(pattern, path string) bool {
	// Handle empty cases
	if pattern == "" && path == "" {
		return true
	}
	if pattern == "" {
		return false
	}
	// A wildcard pattern should not match empty path
	if path == "" && (strings.Contains(pattern, "*") || strings.Contains(pattern, "?")) {
		return false
	}

	// Normalize separators
	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)

	// Handle ** (doublestar) by expanding to regex-like matching
	if strings.Contains(pattern, "**") {
		return matchDoublestar(pattern, path)
	}

	// Use filepath.Match for simple patterns
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false
	}
	return matched
}

// matchDoublestar handles patterns containing **
func matchDoublestar(pattern, path string) bool {
	// Split pattern by **
	parts := strings.Split(pattern, "**")

	if len(parts) == 1 {
		// No ** found, use simple match
		matched, _ := filepath.Match(pattern, path)
		return matched
	}

	// Handle pattern like "packages/**"
	if len(parts) == 2 && parts[1] == "" {
		prefix := strings.TrimSuffix(parts[0], "/")
		if prefix == "" {
			return true // ** matches everything
		}
		return strings.HasPrefix(path, prefix+"/") || path == prefix
	}

	// Handle pattern like "src/**/index.ts"
	if len(parts) == 2 {
		prefix := strings.TrimSuffix(parts[0], "/")
		suffix := strings.TrimPrefix(parts[1], "/")

		// Check if path starts with prefix
		if prefix != "" && !strings.HasPrefix(path, prefix+"/") && path != prefix {
			return false
		}

		// Check if path ends with suffix
		if suffix != "" {
			remainingPath := path
			if prefix != "" {
				remainingPath = strings.TrimPrefix(path, prefix+"/")
			}

			// Check if the remaining path ends with the suffix
			// For "index.ts" suffix, match "index.ts", "foo/index.ts", "foo/bar/index.ts"
			if remainingPath == suffix || strings.HasSuffix(remainingPath, "/"+suffix) {
				return true
			}
			return false
		}

		return true
	}

	// Complex patterns with multiple ** - simplified matching
	return strings.HasPrefix(path, strings.TrimSuffix(parts[0], "/"))
}

// IsNegationPattern checks if a pattern starts with ! (negation)
func IsNegationPattern(pattern string) bool {
	return len(pattern) > 0 && pattern[0] == '!'
}

// getNegationBase returns the pattern without the negation prefix
func getNegationBase(pattern string) string {
	if IsNegationPattern(pattern) {
		return pattern[1:]
	}
	return pattern
}

// FilterPaths filters a list of paths based on include/exclude patterns.
// Patterns starting with ! are exclusion patterns.
// Exclusion patterns are applied after inclusion patterns.
func FilterPaths(paths []string, patterns []string) []string {
	if len(patterns) == 0 {
		return []string{}
	}

	// Separate include and exclude patterns
	var includePatterns, excludePatterns []string
	for _, p := range patterns {
		if IsNegationPattern(p) {
			excludePatterns = append(excludePatterns, getNegationBase(p))
		} else {
			includePatterns = append(includePatterns, p)
		}
	}

	// First apply include patterns
	included := make(map[string]bool)
	for _, path := range paths {
		for _, pattern := range includePatterns {
			if MatchPattern(pattern, path) {
				included[path] = true
				break
			}
		}
	}

	// Then apply exclude patterns
	for _, path := range paths {
		for _, pattern := range excludePatterns {
			if MatchPattern(pattern, path) {
				delete(included, path)
				break
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(included))
	for path := range included {
		result = append(result, path)
	}

	return result
}

// ExpandGlobPatternsFromFiles expands workspace glob patterns against a map of files.
// Returns unique directory paths that contain package.json files.
// This is used when we have file contents but no filesystem access (WASM).
func ExpandGlobPatternsFromFiles(files map[string][]byte, patterns []string) []string {
	// Extract unique directory paths that have package.json
	packageDirs := make(map[string]bool)
	for filePath := range files {
		filePath = filepath.ToSlash(filePath)
		if filepath.Base(filePath) == "package.json" {
			dir := filepath.Dir(filePath)
			if dir != "." && dir != "" {
				packageDirs[dir] = true
			}
		}
	}

	// Convert to slice
	allDirs := make([]string, 0, len(packageDirs))
	for dir := range packageDirs {
		allDirs = append(allDirs, dir)
	}

	// Filter using patterns
	return FilterPaths(allDirs, patterns)
}
