package analyzer

import (
	"testing"
)

// TestNewExclusionMatcher verifies matcher creation.
func TestNewExclusionMatcher(t *testing.T) {
	patterns := []string{
		"packages/legacy",
		"packages/deprecated-*",
		"regex:^@mono/test-.*$",
	}

	matcher, err := NewExclusionMatcher(patterns)
	if err != nil {
		t.Fatalf("NewExclusionMatcher failed: %v", err)
	}

	if matcher == nil {
		t.Fatal("Matcher is nil")
	}

	if !matcher.HasPatterns() {
		t.Error("HasPatterns should return true")
	}

	if matcher.PatternCount() != 3 {
		t.Errorf("PatternCount = %d, want 3", matcher.PatternCount())
	}
}

// TestNewExclusionMatcherInvalidRegex verifies error on invalid regex.
func TestNewExclusionMatcherInvalidRegex(t *testing.T) {
	patterns := []string{
		"regex:[invalid",
	}

	_, err := NewExclusionMatcher(patterns)
	if err == nil {
		t.Error("Expected error for invalid regex, got nil")
	}
}

// TestNewExclusionMatcherNilSafe verifies nil matcher is safe.
func TestNewExclusionMatcherNilSafe(t *testing.T) {
	var matcher *ExclusionMatcher

	if matcher.IsExcluded("any-package") {
		t.Error("Nil matcher should return false")
	}

	if matcher.HasPatterns() {
		t.Error("Nil matcher HasPatterns should return false")
	}

	if matcher.PatternCount() != 0 {
		t.Error("Nil matcher PatternCount should return 0")
	}
}

// TestExactMatching verifies exact package name matching.
func TestExactMatching(t *testing.T) {
	patterns := []string{
		"packages/legacy",
		"@mono/deprecated",
	}

	matcher, _ := NewExclusionMatcher(patterns)

	tests := []struct {
		name     string
		pkg      string
		excluded bool
	}{
		{"exact match packages/legacy", "packages/legacy", true},
		{"exact match @mono/deprecated", "@mono/deprecated", true},
		{"no match packages/legacy-v2", "packages/legacy-v2", false},
		{"no match @mono/deprecated-utils", "@mono/deprecated-utils", false},
		{"no match different package", "packages/core", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcher.IsExcluded(tt.pkg); got != tt.excluded {
				t.Errorf("IsExcluded(%q) = %v, want %v", tt.pkg, got, tt.excluded)
			}
		})
	}
}

// TestGlobMatching verifies glob pattern matching.
func TestGlobMatching(t *testing.T) {
	patterns := []string{
		"packages/deprecated-*",
		"**/test-*",
		"packages/legacy-?",
	}

	matcher, _ := NewExclusionMatcher(patterns)

	tests := []struct {
		name     string
		pkg      string
		excluded bool
	}{
		// deprecated-* pattern
		{"glob packages/deprecated-utils", "packages/deprecated-utils", true},
		{"glob packages/deprecated-api", "packages/deprecated-api", true},
		{"no glob packages/deprecated", "packages/deprecated", false},

		// **/test-* pattern
		{"glob packages/test-utils", "packages/test-utils", true},
		{"glob apps/web/test-helpers", "apps/web/test-helpers", true},
		{"glob test-runner", "test-runner", true},

		// legacy-? pattern (single char)
		{"glob packages/legacy-1", "packages/legacy-1", true},
		{"glob packages/legacy-a", "packages/legacy-a", true},
		{"no glob packages/legacy-12", "packages/legacy-12", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcher.IsExcluded(tt.pkg); got != tt.excluded {
				t.Errorf("IsExcluded(%q) = %v, want %v", tt.pkg, got, tt.excluded)
			}
		})
	}
}

// TestRegexMatching verifies regex pattern matching.
func TestRegexMatching(t *testing.T) {
	patterns := []string{
		"regex:^@mono/legacy-.*$",
		"regex:.*-deprecated$",
		"regex:^test-[0-9]+$",
	}

	matcher, _ := NewExclusionMatcher(patterns)

	tests := []struct {
		name     string
		pkg      string
		excluded bool
	}{
		// ^@mono/legacy-.*$ pattern
		{"regex @mono/legacy-v1", "@mono/legacy-v1", true},
		{"regex @mono/legacy-utils", "@mono/legacy-utils", true},
		{"no regex @mono/legacy", "@mono/legacy", false},

		// .*-deprecated$ pattern
		{"regex @mono/utils-deprecated", "@mono/utils-deprecated", true},
		{"regex packages/api-deprecated", "packages/api-deprecated", true},
		{"no regex @mono/deprecated-utils", "@mono/deprecated-utils", false},

		// ^test-[0-9]+$ pattern
		{"regex test-123", "test-123", true},
		{"regex test-1", "test-1", true},
		{"no regex test-abc", "test-abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcher.IsExcluded(tt.pkg); got != tt.excluded {
				t.Errorf("IsExcluded(%q) = %v, want %v", tt.pkg, got, tt.excluded)
			}
		})
	}
}

// TestMixedPatterns verifies mixed pattern types work together.
func TestMixedPatterns(t *testing.T) {
	patterns := []string{
		"packages/legacy",       // exact
		"packages/deprecated-*", // glob
		"regex:.*-test$",        // regex
	}

	matcher, _ := NewExclusionMatcher(patterns)

	tests := []struct {
		name     string
		pkg      string
		excluded bool
	}{
		{"exact match", "packages/legacy", true},
		{"glob match", "packages/deprecated-utils", true},
		{"regex match", "@mono/core-test", true},
		{"no match", "packages/core", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcher.IsExcluded(tt.pkg); got != tt.excluded {
				t.Errorf("IsExcluded(%q) = %v, want %v", tt.pkg, got, tt.excluded)
			}
		})
	}
}

// TestEmptyPatterns verifies empty pattern list works.
func TestEmptyPatterns(t *testing.T) {
	matcher, _ := NewExclusionMatcher([]string{})

	if matcher.HasPatterns() {
		t.Error("Empty matcher should not have patterns")
	}

	if matcher.IsExcluded("any-package") {
		t.Error("Empty matcher should not exclude anything")
	}
}

// TestParsePattern verifies pattern type detection.
func TestParsePattern(t *testing.T) {
	tests := []struct {
		pattern      string
		expectedType string
		clean        string
	}{
		{"packages/legacy", PatternTypeExact, "packages/legacy"},
		{"packages/deprecated-*", PatternTypeGlob, "packages/deprecated-*"},
		{"**/test-*", PatternTypeGlob, "**/test-*"},
		{"packages/legacy-?", PatternTypeGlob, "packages/legacy-?"},
		{"regex:^test-.*$", PatternTypeRegex, "^test-.*$"},
		{"regex:.*-deprecated$", PatternTypeRegex, ".*-deprecated$"},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			patternType, clean := parsePattern(tt.pattern)
			if patternType != tt.expectedType {
				t.Errorf("parsePattern(%q) type = %q, want %q", tt.pattern, patternType, tt.expectedType)
			}
			if clean != tt.clean {
				t.Errorf("parsePattern(%q) clean = %q, want %q", tt.pattern, clean, tt.clean)
			}
		})
	}
}

// TestDoubleStarGlob verifies ** matches across path separators.
func TestDoubleStarGlob(t *testing.T) {
	patterns := []string{
		"**/utils",
		"packages/**/test",
	}

	matcher, _ := NewExclusionMatcher(patterns)

	tests := []struct {
		name     string
		pkg      string
		excluded bool
	}{
		{"**/utils at root", "utils", true},
		{"**/utils one level", "packages/utils", true},
		{"**/utils deep", "apps/web/src/utils", true},
		{"packages/**/test one level", "packages/test", true},
		{"packages/**/test deep", "packages/core/src/test", true},
		{"no match", "packages/core", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcher.IsExcluded(tt.pkg); got != tt.excluded {
				t.Errorf("IsExcluded(%q) = %v, want %v", tt.pkg, got, tt.excluded)
			}
		})
	}
}

// TestSingleStarGlob verifies * only matches within path segment.
func TestSingleStarGlob(t *testing.T) {
	patterns := []string{
		"packages/*-utils",
	}

	matcher, _ := NewExclusionMatcher(patterns)

	tests := []struct {
		name     string
		pkg      string
		excluded bool
	}{
		{"single level match", "packages/core-utils", true},
		{"single level match 2", "packages/api-utils", true},
		{"does not cross separator", "packages/core/nested-utils", false},
		{"no match", "packages/core", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matcher.IsExcluded(tt.pkg); got != tt.excluded {
				t.Errorf("IsExcluded(%q) = %v, want %v", tt.pkg, got, tt.excluded)
			}
		})
	}
}
