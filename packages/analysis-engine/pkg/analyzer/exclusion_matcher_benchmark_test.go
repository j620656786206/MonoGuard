package analyzer

import (
	"fmt"
	"testing"
	"time"
)

// generateTestPatterns creates a mix of exclusion patterns.
func generateTestPatterns(count int) []string {
	patterns := make([]string, count)
	for i := 0; i < count; i++ {
		switch i % 3 {
		case 0:
			patterns[i] = fmt.Sprintf("packages/legacy-%d", i)
		case 1:
			patterns[i] = fmt.Sprintf("packages/deprecated-*-%d", i)
		case 2:
			patterns[i] = fmt.Sprintf("regex:^@mono/test-%d-.*$", i)
		}
	}
	return patterns
}

// generateTestPackageNames creates package names for testing.
func generateTestPackageNames(count int) []string {
	names := make([]string, count)
	for i := 0; i < count; i++ {
		names[i] = fmt.Sprintf("@mono/pkg-%d", i)
	}
	return names
}

// BenchmarkExclusionMatcher100Packages20Patterns benchmarks exclusion matching.
// AC7: 100 packages × 20 patterns < 50ms
func BenchmarkExclusionMatcher100Packages20Patterns(b *testing.B) {
	patterns := generateTestPatterns(20)
	matcher, _ := NewExclusionMatcher(patterns)
	packages := generateTestPackageNames(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pkg := range packages {
			matcher.IsExcluded(pkg)
		}
	}
}

// BenchmarkExclusionMatcherExactOnly tests exact matching performance.
func BenchmarkExclusionMatcherExactOnly(b *testing.B) {
	patterns := make([]string, 20)
	for i := 0; i < 20; i++ {
		patterns[i] = fmt.Sprintf("packages/legacy-%d", i)
	}
	matcher, _ := NewExclusionMatcher(patterns)
	packages := generateTestPackageNames(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pkg := range packages {
			matcher.IsExcluded(pkg)
		}
	}
}

// BenchmarkExclusionMatcherGlobOnly tests glob matching performance.
func BenchmarkExclusionMatcherGlobOnly(b *testing.B) {
	patterns := make([]string, 20)
	for i := 0; i < 20; i++ {
		patterns[i] = fmt.Sprintf("packages/deprecated-*-%d", i)
	}
	matcher, _ := NewExclusionMatcher(patterns)
	packages := generateTestPackageNames(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pkg := range packages {
			matcher.IsExcluded(pkg)
		}
	}
}

// BenchmarkExclusionMatcherRegexOnly tests regex matching performance.
func BenchmarkExclusionMatcherRegexOnly(b *testing.B) {
	patterns := make([]string, 20)
	for i := 0; i < 20; i++ {
		patterns[i] = fmt.Sprintf("regex:^@mono/test-%d-.*$", i)
	}
	matcher, _ := NewExclusionMatcher(patterns)
	packages := generateTestPackageNames(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pkg := range packages {
			matcher.IsExcluded(pkg)
		}
	}
}

// TestPerformanceRequirement_ExclusionMatcher tests AC7:
// "Given 100 packages with 20 exclusion patterns, matching completes in < 50ms"
func TestPerformanceRequirement_ExclusionMatcher(t *testing.T) {
	patterns := generateTestPatterns(20)
	matcher, err := NewExclusionMatcher(patterns)
	if err != nil {
		t.Fatalf("NewExclusionMatcher failed: %v", err)
	}

	packages := generateTestPackageNames(100)

	// Warm up
	for _, pkg := range packages {
		matcher.IsExcluded(pkg)
	}

	// Time the actual run
	start := time.Now()
	for _, pkg := range packages {
		matcher.IsExcluded(pkg)
	}
	duration := time.Since(start)

	// AC7: Must complete in < 50ms
	if duration >= 50*time.Millisecond {
		t.Errorf("Exclusion matching took %v, want < 50ms", duration)
	}

	t.Logf("100 packages × 20 patterns: %v", duration)
}

// TestPerformanceRequirement_ExclusionMatcher500Packages tests scalability.
func TestPerformanceRequirement_ExclusionMatcher500Packages(t *testing.T) {
	patterns := generateTestPatterns(50)
	matcher, err := NewExclusionMatcher(patterns)
	if err != nil {
		t.Fatalf("NewExclusionMatcher failed: %v", err)
	}

	packages := generateTestPackageNames(500)

	start := time.Now()
	for _, pkg := range packages {
		matcher.IsExcluded(pkg)
	}
	duration := time.Since(start)

	// Even 500 packages should be fast
	if duration >= 100*time.Millisecond {
		t.Errorf("Exclusion matching took %v, want < 100ms for 500 packages", duration)
	}

	t.Logf("500 packages × 50 patterns: %v", duration)
}
