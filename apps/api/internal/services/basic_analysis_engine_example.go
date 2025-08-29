package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// ExampleBasicAnalysisEngine demonstrates how to use the BasicAnalysisEngine
func ExampleBasicAnalysisEngine() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	// Create the analysis engine
	engine := NewBasicAnalysisEngine(logger)
	
	// Example usage
	projectID := "example-project-123"
	rootPath := "/path/to/monorepo"
	
	ctx := context.Background()
	results, err := engine.AnalyzeRepository(ctx, rootPath, projectID)
	if err != nil {
		logger.WithError(err).Error("Analysis failed")
		return
	}
	
	// Print results summary
	fmt.Printf("=== Analysis Results for Project: %s ===\n", projectID)
	fmt.Printf("Total Packages: %d\n", results.Summary.TotalPackages)
	fmt.Printf("Duplicate Dependencies: %d\n", results.Summary.DuplicateCount)
	fmt.Printf("Version Conflicts: %d\n", results.Summary.ConflictCount)
	fmt.Printf("Health Score: %.1f/100\n\n", results.Summary.HealthScore)
	
	// Print duplicate dependencies
	if len(results.DuplicateDependencies) > 0 {
		fmt.Println("=== Duplicate Dependencies ===")
		for i, duplicate := range results.DuplicateDependencies {
			fmt.Printf("%d. %s\n", i+1, duplicate.PackageName)
			fmt.Printf("   Versions: %v\n", duplicate.Versions)
			fmt.Printf("   Risk Level: %s\n", duplicate.RiskLevel)
			fmt.Printf("   Estimated Waste: %s\n", duplicate.EstimatedWaste)
			fmt.Printf("   Recommendation: %s\n", duplicate.Recommendation)
			fmt.Printf("   Affected Packages: %v\n\n", duplicate.AffectedPackages)
		}
	}
	
	// Print version conflicts
	if len(results.VersionConflicts) > 0 {
		fmt.Println("=== Version Conflicts ===")
		for i, conflict := range results.VersionConflicts {
			fmt.Printf("%d. %s\n", i+1, conflict.PackageName)
			fmt.Printf("   Risk Level: %s\n", conflict.RiskLevel)
			fmt.Printf("   Resolution: %s\n", conflict.Resolution)
			fmt.Printf("   Impact: %s\n", conflict.Impact)
			fmt.Println("   Conflicting Versions:")
			for _, version := range conflict.ConflictingVersions {
				fmt.Printf("     - %s (used by: %v, breaking: %v)\n", 
					version.Version, version.Packages, version.IsBreaking)
			}
			fmt.Println()
		}
	}
	
	// Print bundle impact
	fmt.Println("=== Bundle Impact ===")
	fmt.Printf("Total Size: %s\n", results.BundleImpact.TotalSize)
	fmt.Printf("Duplicate Waste: %s\n", results.BundleImpact.DuplicateSize)
	fmt.Printf("Potential Savings: %s\n", results.BundleImpact.PotentialSavings)
	
	if len(results.BundleImpact.Breakdown) > 0 {
		fmt.Println("\nTop Dependencies:")
		for i, breakdown := range results.BundleImpact.Breakdown {
			if i >= 5 { // Show top 5
				break
			}
			fmt.Printf("  %d. %s - %s (%.1f%%, %d duplicates)\n",
				i+1, breakdown.PackageName, breakdown.Size, 
				breakdown.Percentage, breakdown.Duplicates)
		}
	}
	
	// Convert to JSON for API responses
	jsonResults, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		logger.WithError(err).Error("Failed to marshal results to JSON")
		return
	}
	
	fmt.Printf("\n=== JSON Output (first 500 chars) ===\n%s...\n", 
		string(jsonResults[:min(500, len(jsonResults))]))
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// QuickAnalysisExample shows a simplified usage pattern
func QuickAnalysisExample(repoPath string) (*AnalysisQuickSummary, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel) // Reduce noise for quick analysis
	
	engine := NewBasicAnalysisEngine(logger)
	results, err := engine.AnalyzeRepository(context.Background(), repoPath, "quick-analysis")
	if err != nil {
		return nil, fmt.Errorf("analysis failed: %w", err)
	}
	
	// Create a simplified summary
	summary := &AnalysisQuickSummary{
		HealthScore:        results.Summary.HealthScore,
		TotalPackages:      results.Summary.TotalPackages,
		IssueCount:         results.Summary.DuplicateCount + results.Summary.ConflictCount,
		TopIssues:          make([]string, 0),
		RecommendedActions: make([]string, 0),
	}
	
	// Add top issues
	for _, duplicate := range results.DuplicateDependencies {
		if len(summary.TopIssues) < 3 {
			summary.TopIssues = append(summary.TopIssues, 
				fmt.Sprintf("Duplicate dependency: %s (%d versions)", 
					duplicate.PackageName, len(duplicate.Versions)))
		}
	}
	
	for _, conflict := range results.VersionConflicts {
		if len(summary.TopIssues) < 3 {
			summary.TopIssues = append(summary.TopIssues, 
				fmt.Sprintf("Version conflict: %s", conflict.PackageName))
		}
	}
	
	// Add recommendations
	if results.Summary.DuplicateCount > 0 {
		summary.RecommendedActions = append(summary.RecommendedActions, 
			"Consolidate duplicate dependencies to reduce bundle size")
	}
	
	if results.Summary.ConflictCount > 0 {
		summary.RecommendedActions = append(summary.RecommendedActions, 
			"Resolve version conflicts to ensure compatibility")
	}
	
	if results.Summary.HealthScore > 90 {
		summary.RecommendedActions = append(summary.RecommendedActions, 
			"Great job! Your dependency management is excellent")
	} else if results.Summary.HealthScore < 50 {
		summary.RecommendedActions = append(summary.RecommendedActions, 
			"Consider a comprehensive dependency audit and cleanup")
	}
	
	return summary, nil
}

// AnalysisQuickSummary provides a simplified view of analysis results
type AnalysisQuickSummary struct {
	HealthScore        float64  `json:"healthScore"`
	TotalPackages      int      `json:"totalPackages"`
	IssueCount         int      `json:"issueCount"`
	TopIssues          []string `json:"topIssues"`
	RecommendedActions []string `json:"recommendedActions"`
}