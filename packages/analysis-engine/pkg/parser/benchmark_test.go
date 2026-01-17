// Package parser benchmark tests for performance verification.
package parser

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// generateWorkspaceFiles generates a simulated workspace with n packages
func generateWorkspaceFiles(n int) map[string][]byte {
	files := make(map[string][]byte)

	// Root package.json with workspaces
	rootPkg := map[string]interface{}{
		"name":       "monorepo-root",
		"workspaces": []string{"packages/*"},
	}
	rootJSON, _ := json.Marshal(rootPkg)
	files["package.json"] = rootJSON

	// package-lock.json to indicate npm workspace
	files["package-lock.json"] = []byte(`{}`)

	// Generate n packages
	for i := 0; i < n; i++ {
		pkgName := fmt.Sprintf("@mono/pkg-%d", i)
		pkgPath := fmt.Sprintf("packages/pkg-%d/package.json", i)

		// Create dependencies to other packages (circular dependency simulation)
		deps := make(map[string]string)
		if i > 0 {
			deps[fmt.Sprintf("@mono/pkg-%d", i-1)] = "^1.0.0"
		}
		if i < n-1 {
			deps[fmt.Sprintf("@mono/pkg-%d", i+1)] = "^1.0.0"
		}

		pkg := map[string]interface{}{
			"name":            pkgName,
			"version":         "1.0.0",
			"dependencies":    deps,
			"devDependencies": map[string]string{"typescript": "^5.0.0"},
		}
		pkgJSON, _ := json.Marshal(pkg)
		files[pkgPath] = pkgJSON
	}

	return files
}

func BenchmarkParse10Packages(b *testing.B) {
	files := generateWorkspaceFiles(10)
	p := NewParser("/workspace")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Parse(files)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParse50Packages(b *testing.B) {
	files := generateWorkspaceFiles(50)
	p := NewParser("/workspace")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Parse(files)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func BenchmarkParse100Packages(b *testing.B) {
	files := generateWorkspaceFiles(100)
	p := NewParser("/workspace")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.Parse(files)
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}
	}
}

func TestParse100PackagesPerformance(t *testing.T) {
	// AC6: Given a workspace with 100 packages
	// When parsing completes
	// Then it finishes in < 1 second
	files := generateWorkspaceFiles(100)
	p := NewParser("/workspace")

	start := time.Now()
	result, err := p.Parse(files)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify timing constraint
	if elapsed >= time.Second {
		t.Errorf("Parsing 100 packages took %v, should be < 1 second", elapsed)
	}

	// Verify all packages were parsed
	if len(result.Packages) != 100 {
		t.Errorf("Expected 100 packages, got %d", len(result.Packages))
	}

	t.Logf("Parsed 100 packages in %v", elapsed)
}

func TestParse500PackagesPerformance(t *testing.T) {
	// Stress test with larger workspace
	files := generateWorkspaceFiles(500)
	p := NewParser("/workspace")

	start := time.Now()
	result, err := p.Parse(files)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// For 500 packages, allow up to 5 seconds
	if elapsed >= 5*time.Second {
		t.Errorf("Parsing 500 packages took %v, should be < 5 seconds", elapsed)
	}

	if len(result.Packages) != 500 {
		t.Errorf("Expected 500 packages, got %d", len(result.Packages))
	}

	t.Logf("Parsed 500 packages in %v", elapsed)
}
