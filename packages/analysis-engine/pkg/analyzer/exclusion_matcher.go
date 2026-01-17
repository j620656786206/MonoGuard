// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file implements package exclusion pattern matching for Story 2.6.
package analyzer

import (
	"regexp"
	"strings"
)

// Pattern type constants
const (
	PatternTypeExact = "exact"
	PatternTypeGlob  = "glob"
	PatternTypeRegex = "regex"
)

// ExclusionMatcher handles package exclusion pattern matching.
// It supports three types of patterns:
//   - Exact matches: "packages/legacy" matches only "packages/legacy"
//   - Glob patterns: "packages/deprecated-*" uses wildcards
//   - Regex patterns: "regex:^packages/test-.*$" uses regular expressions
type ExclusionMatcher struct {
	exactMatches  map[string]bool
	globPatterns  []string
	regexPatterns []*regexp.Regexp
}

// NewExclusionMatcher creates a matcher from exclusion patterns.
// Patterns are parsed and categorized by type for efficient matching.
// Returns an error if a regex pattern is invalid.
func NewExclusionMatcher(patterns []string) (*ExclusionMatcher, error) {
	em := &ExclusionMatcher{
		exactMatches:  make(map[string]bool),
		globPatterns:  []string{},
		regexPatterns: []*regexp.Regexp{},
	}

	for _, pattern := range patterns {
		patternType, cleanPattern := parsePattern(pattern)

		switch patternType {
		case PatternTypeExact:
			em.exactMatches[cleanPattern] = true
		case PatternTypeGlob:
			em.globPatterns = append(em.globPatterns, cleanPattern)
		case PatternTypeRegex:
			re, err := regexp.Compile(cleanPattern)
			if err != nil {
				return nil, err
			}
			em.regexPatterns = append(em.regexPatterns, re)
		}
	}

	return em, nil
}

// IsExcluded checks if a package name matches any exclusion pattern.
// Matching order: exact > glob > regex (for performance).
func (em *ExclusionMatcher) IsExcluded(packageName string) bool {
	if em == nil {
		return false
	}

	// Check exact matches first (fastest)
	if em.exactMatches[packageName] {
		return true
	}

	// Check glob patterns
	for _, pattern := range em.globPatterns {
		if matchGlob(pattern, packageName) {
			return true
		}
	}

	// Check regex patterns (slowest)
	for _, re := range em.regexPatterns {
		if re.MatchString(packageName) {
			return true
		}
	}

	return false
}

// parsePattern categorizes a pattern as exact, glob, or regex.
// Regex patterns are prefixed with "regex:".
// Glob patterns contain *, **, or ?.
// Everything else is an exact match.
func parsePattern(pattern string) (patternType, cleanPattern string) {
	// Check for regex prefix
	if strings.HasPrefix(pattern, "regex:") {
		return PatternTypeRegex, strings.TrimPrefix(pattern, "regex:")
	}

	// Check for glob characters
	if containsGlobChars(pattern) {
		return PatternTypeGlob, pattern
	}

	// Default to exact match
	return PatternTypeExact, pattern
}

// containsGlobChars checks if a pattern contains glob wildcard characters.
func containsGlobChars(pattern string) bool {
	return strings.ContainsAny(pattern, "*?")
}

// matchGlob matches a package name against a glob pattern.
// Supports:
//   - * matches any sequence of non-separator characters
//   - ** matches any sequence of characters including separators
//   - ? matches any single character
func matchGlob(pattern, name string) bool {
	return matchGlobRecursive(pattern, name)
}

// matchGlobRecursive implements glob matching recursively.
func matchGlobRecursive(pattern, name string) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			// Check for ** (matches everything including separators)
			if len(pattern) > 1 && pattern[1] == '*' {
				// Skip **
				pattern = pattern[2:]
				// Skip optional separator after **
				if len(pattern) > 0 && pattern[0] == '/' {
					pattern = pattern[1:]
				}
				// ** at end matches everything
				if len(pattern) == 0 {
					return true
				}
				// Try matching rest of pattern at every position
				for i := 0; i <= len(name); i++ {
					if matchGlobRecursive(pattern, name[i:]) {
						return true
					}
				}
				return false
			}

			// Single * - matches any sequence except separator
			pattern = pattern[1:]
			// * at end matches rest of segment
			if len(pattern) == 0 {
				// Check no more separators in name
				return !strings.Contains(name, "/")
			}
			// Try matching rest of pattern at every position (within segment)
			for i := 0; i <= len(name); i++ {
				// Don't cross separator for single *
				if i > 0 && name[i-1] == '/' {
					break
				}
				if matchGlobRecursive(pattern, name[i:]) {
					return true
				}
			}
			return false

		case '?':
			// ? matches any single character except separator
			if len(name) == 0 || name[0] == '/' {
				return false
			}
			pattern = pattern[1:]
			name = name[1:]

		default:
			// Regular character - must match exactly
			if len(name) == 0 || pattern[0] != name[0] {
				return false
			}
			pattern = pattern[1:]
			name = name[1:]
		}
	}

	// Pattern exhausted - name must also be exhausted
	return len(name) == 0
}

// HasPatterns returns true if the matcher has any patterns configured.
func (em *ExclusionMatcher) HasPatterns() bool {
	if em == nil {
		return false
	}
	return len(em.exactMatches) > 0 || len(em.globPatterns) > 0 || len(em.regexPatterns) > 0
}

// PatternCount returns the total number of patterns.
func (em *ExclusionMatcher) PatternCount() int {
	if em == nil {
		return 0
	}
	return len(em.exactMatches) + len(em.globPatterns) + len(em.regexPatterns)
}
