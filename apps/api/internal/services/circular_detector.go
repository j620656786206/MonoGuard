package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/internal/repository"
	"github.com/sirupsen/logrus"
)

// CircularDetectorService handles circular dependency detection
type CircularDetectorService struct {
	projectRepo  repository.ProjectRepositoryInterface
	analysisRepo repository.AnalysisRepositoryInterface
	logger       *logrus.Logger
}

// NewCircularDetectorService creates a new circular detector service
func NewCircularDetectorService(
	projectRepo repository.ProjectRepositoryInterface,
	analysisRepo repository.AnalysisRepositoryInterface,
	logger *logrus.Logger,
) *CircularDetectorService {
	return &CircularDetectorService{
		projectRepo:  projectRepo,
		analysisRepo: analysisRepo,
		logger:       logger,
	}
}

// CircularDependency represents a circular dependency path
type CircularDependency struct {
	CyclePath       []string                `json:"cycle_path"`
	CycleLength     int                     `json:"cycle_length"`
	BreakPoints     []BreakPointSuggestion  `json:"break_points"`
	ImpactAnalysis  CycleImpactReport       `json:"impact_analysis"`
	ResolutionSteps []ResolutionStep        `json:"resolution_steps"`
}

// BreakPointSuggestion represents a suggested breakpoint for circular dependency
type BreakPointSuggestion struct {
	PackageName         string `json:"package_name"`
	ImportToRemove      string `json:"import_to_remove"`
	AlternativeApproach string `json:"alternative_approach"`
	EstimatedEffort     string `json:"estimated_effort"`
	RiskLevel           string `json:"risk_level"`
}

// CycleImpactReport represents the impact analysis of a circular dependency
type CycleImpactReport struct {
	AffectedPackages    int     `json:"affected_packages"`
	EstimatedRefactorTime string `json:"estimated_refactor_time"`
	BusinessRisk        string  `json:"business_risk"`
	TechnicalDebt       string  `json:"technical_debt"`
}

// ResolutionStep represents a step in resolving circular dependency
type ResolutionStep struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
	CodeExample string `json:"code_example,omitempty"`
	Effort      string `json:"effort"`
}

// PackageGraph represents the dependency graph
type PackageGraph struct {
	Nodes map[string]*CircularPackageNode
	Edges map[string][]string
}

// CircularPackageNode represents a package in the graph for circular detection
type CircularPackageNode struct {
	Name         string
	Path         string
	Dependencies []string
	Visited      bool
	InStack      bool
}

// DetectCircularDependencies detects circular dependencies in a project
func (s *CircularDetectorService) DetectCircularDependencies(ctx context.Context, projectID string, repoPath string) (*models.DependencyAnalysis, error) {
	s.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"repo_path":  repoPath,
	}).Info("Starting circular dependency detection")

	// Get project - temporarily not using project.Settings
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Create analysis record
	analysis := &models.DependencyAnalysis{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Status:    "in_progress",
		StartedAt: time.Now().UTC(),
		Results:   &models.DependencyAnalysisResults{},
		Metadata: &models.AnalysisMetadata{
			Version:          "1.0.0",
			FilesProcessed:   0,
			PackagesAnalyzed: 0,
			Configuration: map[string]interface{}{
				"analysis_type":     "circular_dependencies",
				"exclude_patterns": []string{}, // project.Settings.ExcludePatterns - temporarily disabled
				"include_patterns": []string{}, // project.Settings.IncludePatterns - temporarily disabled
			},
			Environment: models.AnalysisEnvironment{
				Platform: "linux",
			},
		},
	}

	// Save initial analysis
	if err := s.analysisRepo.CreateDependencyAnalysis(ctx, analysis); err != nil {
		return nil, fmt.Errorf("failed to create analysis: %w", err)
	}

	// Build dependency graph
	graph, err := s.buildDependencyGraph(repoPath, []string{}) // project.Settings.ExcludePatterns - temporarily disabled
	if err != nil {
		s.logger.WithError(err).Error("Failed to build dependency graph")
		analysis.Status = "failed"
		s.analysisRepo.UpdateDependencyAnalysis(ctx, analysis.ID, map[string]interface{}{
			"status": "failed",
		})
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Detect circular dependencies using DFS
	circularDeps := s.detectCircularDepsWithDFS(graph)

	// Generate resolution suggestions
	for i, cycle := range circularDeps {
		circularDeps[i].BreakPoints = s.generateBreakPointSuggestions(cycle)
		circularDeps[i].ImpactAnalysis = s.analyzeImpact(cycle)
		circularDeps[i].ResolutionSteps = s.generateResolutionSteps(cycle)
	}

	// Update analysis results
	analysis.Results.CircularDependencies = s.convertToModelCircularDeps(circularDeps)
	analysis.Results.Summary.CircularCount = len(circularDeps)
	analysis.Results.Summary.TotalPackages = len(graph.Nodes)
	
	// Calculate health score based on circular dependencies
	healthScore := s.calculateCircularHealthScore(len(circularDeps), len(graph.Nodes))
	analysis.Results.Summary.HealthScore = float64(healthScore)

	// Update metadata
	analysis.Metadata.PackagesAnalyzed = len(graph.Nodes)
	analysis.Metadata.Duration = int64(time.Since(analysis.StartedAt).Milliseconds())
	
	// Mark as completed
	now := time.Now().UTC()
	analysis.CompletedAt = &now
	analysis.Status = "completed"

	// Update analysis using Save instead of Updates to handle complex types
	if err := s.analysisRepo.SaveDependencyAnalysis(ctx, analysis); err != nil {
		return nil, fmt.Errorf("failed to update analysis: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"project_id":        projectID,
		"circular_count":    len(circularDeps),
		"total_packages":    len(graph.Nodes),
		"health_score":      healthScore,
		"duration_ms":       analysis.Metadata.Duration,
	}).Info("Circular dependency detection completed")

	return analysis, nil
}

// buildDependencyGraph builds a dependency graph from the repository
func (s *CircularDetectorService) buildDependencyGraph(repoPath string, excludePatterns []string) (*PackageGraph, error) {
	graph := &PackageGraph{
		Nodes: make(map[string]*CircularPackageNode),
		Edges: make(map[string][]string),
	}

	// Parse all package.json files in the repository
	packages, err := s.discoverPackages(repoPath, excludePatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to discover packages: %w", err)
	}

	s.logger.WithField("package_count", len(packages)).Info("Discovered packages for dependency graph")

	// Create nodes for all packages
	for _, pkg := range packages {
		relativePath, err := filepath.Rel(repoPath, pkg.Path)
		if err != nil {
			relativePath = pkg.Path
		}
		
		packageName := pkg.Name
		if packageName == "" {
			// Use directory name if no name specified
			packageName = filepath.Base(filepath.Dir(pkg.Path))
		}

		graph.Nodes[packageName] = &CircularPackageNode{
			Name: packageName,
			Path: relativePath,
		}
	}

	// Build dependency edges
	for _, pkg := range packages {
		packageName := pkg.Name
		if packageName == "" {
			packageName = filepath.Base(filepath.Dir(pkg.Path))
		}

		var dependencies []string
		
		// Add production dependencies that reference other workspace packages
		for depName := range pkg.Dependencies {
			if _, exists := graph.Nodes[depName]; exists {
				dependencies = append(dependencies, depName)
			}
		}

		// Add dev dependencies for more complete analysis
		for depName := range pkg.DevDependencies {
			if _, exists := graph.Nodes[depName]; exists {
				dependencies = append(dependencies, depName)
			}
		}

		graph.Nodes[packageName].Dependencies = dependencies
		graph.Edges[packageName] = dependencies
	}

	s.logger.WithFields(logrus.Fields{
		"total_nodes": len(graph.Nodes),
		"total_edges": s.countEdges(graph),
	}).Info("Built dependency graph successfully")

	return graph, nil
}

// PackageInfo represents a discovered package
type CircularPackageInfo struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Path            string            `json:"path"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// discoverPackages finds all package.json files in the repository
func (s *CircularDetectorService) discoverPackages(rootPath string, excludePatterns []string) ([]*CircularPackageInfo, error) {
	var packages []*CircularPackageInfo

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip common directories that shouldn't contain workspace packages
		if info.IsDir() {
			dirName := filepath.Base(path)
			if dirName == "node_modules" || dirName == ".git" || dirName == "dist" || 
			   dirName == "build" || dirName == ".next" || dirName == "coverage" || 
			   dirName == "tmp" || dirName == "__pycache__" {
				return filepath.SkipDir
			}
		}

		// Check exclude patterns
		for _, pattern := range excludePatterns {
			if matched, _ := filepath.Match(pattern, path); matched {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Process package.json files
		if info.Name() == "package.json" {
			pkg, parseErr := s.parsePackageJSON(path)
			if parseErr != nil {
				s.logger.WithError(parseErr).WithField("path", path).Warn("Failed to parse package.json, skipping")
				return nil // Continue processing other files
			}
			packages = append(packages, pkg)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking repository: %w", err)
	}

	return packages, nil
}

// parsePackageJSON parses a single package.json file
func (s *CircularDetectorService) parsePackageJSON(filePath string) (*CircularPackageInfo, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var rawPackage map[string]interface{}
	if err := json.Unmarshal(data, &rawPackage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	pkg := &CircularPackageInfo{
		Path:            filePath,
		Dependencies:    make(map[string]string),
		DevDependencies: make(map[string]string),
	}

	// Extract basic package information
	if name, ok := rawPackage["name"].(string); ok {
		pkg.Name = name
	}

	if version, ok := rawPackage["version"].(string); ok {
		pkg.Version = version
	}

	// Parse dependencies
	if deps, ok := rawPackage["dependencies"].(map[string]interface{}); ok {
		for name, version := range deps {
			if versionStr, ok := version.(string); ok {
				pkg.Dependencies[name] = versionStr
			}
		}
	}

	if devDeps, ok := rawPackage["devDependencies"].(map[string]interface{}); ok {
		for name, version := range devDeps {
			if versionStr, ok := version.(string); ok {
				pkg.DevDependencies[name] = versionStr
			}
		}
	}

	return pkg, nil
}

// countEdges counts the total number of edges in the graph
func (s *CircularDetectorService) countEdges(graph *PackageGraph) int {
	count := 0
	for _, edges := range graph.Edges {
		count += len(edges)
	}
	return count
}

// detectCircularDepsWithDFS detects circular dependencies using Depth-First Search
func (s *CircularDetectorService) detectCircularDepsWithDFS(graph *PackageGraph) []CircularDependency {
	var cycles []CircularDependency
	var stack []string
	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	// DFS function
	var dfs func(string) bool
	dfs = func(node string) bool {
		if inStack[node] {
			// Found a cycle, extract it from stack
			cycleStart := -1
			for i, n := range stack {
				if n == node {
					cycleStart = i
					break
				}
			}
			if cycleStart >= 0 {
				cyclePath := append(stack[cycleStart:], node)
				cycles = append(cycles, CircularDependency{
					CyclePath:   cyclePath,
					CycleLength: len(cyclePath) - 1,
				})
			}
			return true
		}

		if visited[node] {
			return false
		}

		visited[node] = true
		inStack[node] = true
		stack = append(stack, node)

		// Visit all dependencies
		for _, dep := range graph.Edges[node] {
			if dfs(dep) {
				// Cycle found in subtree
			}
		}

		// Remove from stack
		inStack[node] = false
		if len(stack) > 0 && stack[len(stack)-1] == node {
			stack = stack[:len(stack)-1]
		}

		return false
	}

	// Run DFS from each unvisited node
	for node := range graph.Nodes {
		if !visited[node] {
			dfs(node)
		}
	}

	return cycles
}

// generateBreakPointSuggestions generates breakpoint suggestions for a cycle
func (s *CircularDetectorService) generateBreakPointSuggestions(cycle CircularDependency) []BreakPointSuggestion {
	var suggestions []BreakPointSuggestion

	for i := 0; i < len(cycle.CyclePath)-1; i++ {
		from := cycle.CyclePath[i]
		to := cycle.CyclePath[i+1]

		suggestion := BreakPointSuggestion{
			PackageName:    from,
			ImportToRemove: to,
			RiskLevel:     "medium",
			EstimatedEffort: "2-4 hours",
		}

		// Generate specific suggestions based on package types
		if strings.Contains(from, "libs/") && strings.Contains(to, "libs/") {
			suggestion.AlternativeApproach = "Extract shared interface to a common package"
			suggestion.RiskLevel = "low"
		} else if strings.Contains(from, "apps/") {
			suggestion.AlternativeApproach = "Use dependency injection or event-driven architecture"
			suggestion.RiskLevel = "high"
			suggestion.EstimatedEffort = "4-8 hours"
		} else {
			suggestion.AlternativeApproach = "Refactor to use dependency inversion principle"
		}

		suggestions = append(suggestions, suggestion)
	}

	return suggestions
}

// analyzeImpact analyzes the impact of a circular dependency
func (s *CircularDetectorService) analyzeImpact(cycle CircularDependency) CycleImpactReport {
	return CycleImpactReport{
		AffectedPackages:      cycle.CycleLength,
		EstimatedRefactorTime: fmt.Sprintf("%d-%.0f hours", cycle.CycleLength*2, float64(cycle.CycleLength)*4.5),
		BusinessRisk:          s.assessBusinessRisk(cycle),
		TechnicalDebt:         s.assessTechnicalDebt(cycle),
	}
}

// assessBusinessRisk assesses the business risk of a circular dependency
func (s *CircularDetectorService) assessBusinessRisk(cycle CircularDependency) string {
	hasAppsDeps := false
	for _, pkg := range cycle.CyclePath {
		if strings.Contains(pkg, "apps/") {
			hasAppsDeps = true
			break
		}
	}

	if hasAppsDeps {
		return "High - affects application-level dependencies"
	} else if cycle.CycleLength > 3 {
		return "Medium - complex dependency chain"
	}
	return "Low - library-level circular dependency"
}

// assessTechnicalDebt assesses the technical debt of a circular dependency
func (s *CircularDetectorService) assessTechnicalDebt(cycle CircularDependency) string {
	if cycle.CycleLength <= 2 {
		return "Low technical debt - simple circular reference"
	} else if cycle.CycleLength <= 4 {
		return "Medium technical debt - moderate complexity"
	}
	return "High technical debt - complex circular dependency chain"
}

// generateResolutionSteps generates step-by-step resolution instructions
func (s *CircularDetectorService) generateResolutionSteps(cycle CircularDependency) []ResolutionStep {
	steps := []ResolutionStep{
		{
			Step:        1,
			Description: "Identify the circular dependency chain",
			CodeExample: fmt.Sprintf("// Cycle detected: %s", strings.Join(cycle.CyclePath, " → ")),
			Effort:      "15 minutes",
		},
		{
			Step:        2,
			Description: "Analyze import statements in each package",
			Effort:      "30 minutes",
		},
		{
			Step:        3,
			Description: "Choose the best breakpoint based on business logic separation",
			Effort:      "1 hour",
		},
		{
			Step:        4,
			Description: "Refactor code to break the circular dependency",
			CodeExample: `// Example: Extract interface
interface IUserService {
  getUser(id: string): Promise<User>;
}

// Use interface instead of concrete implementation`,
			Effort: "2-4 hours",
		},
		{
			Step:        5,
			Description: "Run tests to ensure functionality is preserved",
			Effort:      "1 hour",
		},
	}

	return steps
}

// calculateCircularHealthScore calculates health score based on circular dependencies
func (s *CircularDetectorService) calculateCircularHealthScore(circularCount, totalPackages int) int {
	if totalPackages == 0 {
		return 100
	}

	// Base score of 100, reduce by circular dependency ratio
	circularRatio := float64(circularCount) / float64(totalPackages)
	penalty := int(circularRatio * 80) // Up to 80 points penalty

	score := 100 - penalty
	if score < 0 {
		score = 0
	}

	return score
}

// convertToModelCircularDeps converts internal circular deps to model format
func (s *CircularDetectorService) convertToModelCircularDeps(cycles []CircularDependency) []models.CircularDependency {
	var result []models.CircularDependency

	for _, cycle := range cycles {
		modelCycle := models.CircularDependency{
			Cycle:    cycle.CyclePath,
			Type:     "direct",
			Severity: s.determineSeverity(cycle),
			Impact:   fmt.Sprintf("Affects %d packages", cycle.CycleLength),
		}
		result = append(result, modelCycle)
	}

	return result
}

// determineSeverity determines the severity of a circular dependency
func (s *CircularDetectorService) determineSeverity(cycle CircularDependency) models.Severity {
	hasAppsDeps := false
	for _, pkg := range cycle.CyclePath {
		if strings.Contains(pkg, "apps/") {
			hasAppsDeps = true
			break
		}
	}

	if hasAppsDeps || cycle.CycleLength > 4 {
		return models.SeverityHigh
	} else if cycle.CycleLength > 2 {
		return models.SeverityMedium
	}
	return models.SeverityLow
}