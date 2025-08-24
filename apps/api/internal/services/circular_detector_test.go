package services

import (
	"context"
	"testing"

	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepository is a mock implementation of ProjectRepository
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetAll(ctx context.Context, params *repository.QueryParams) ([]*models.Project, int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*models.Project), args.Get(1).(int64), args.Error(2)
}


func (m *MockProjectRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) UpdateStatus(ctx context.Context, id string, status models.Status, lastAnalysisAt *string) error {
	args := m.Called(ctx, id, status, lastAnalysisAt)
	return args.Error(0)
}

func (m *MockProjectRepository) GetByOwnerID(ctx context.Context, ownerID string, params *repository.QueryParams) ([]*models.Project, int64, error) {
	args := m.Called(ctx, ownerID, params)
	return args.Get(0).([]*models.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockProjectRepository) UpdateHealthScore(ctx context.Context, projectID string, score int) error {
	args := m.Called(ctx, projectID, score)
	return args.Error(0)
}

// MockAnalysisRepository is a mock implementation of AnalysisRepository
type MockAnalysisRepository struct {
	mock.Mock
}

func (m *MockAnalysisRepository) Create(ctx context.Context, analysis *models.DependencyAnalysis) error {
	args := m.Called(ctx, analysis)
	return args.Error(0)
}

func (m *MockAnalysisRepository) GetByID(ctx context.Context, id string) (*models.DependencyAnalysis, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DependencyAnalysis), args.Error(1)
}

func (m *MockAnalysisRepository) Update(ctx context.Context, analysis *models.DependencyAnalysis) error {
	args := m.Called(ctx, analysis)
	return args.Error(0)
}

func (m *MockAnalysisRepository) GetProjectAnalyses(ctx context.Context, projectID string, limit, offset int) ([]*models.DependencyAnalysis, int64, error) {
	args := m.Called(ctx, projectID, limit, offset)
	return args.Get(0).([]*models.DependencyAnalysis), args.Get(1).(int64), args.Error(2)
}

func (m *MockAnalysisRepository) UpdateDependencyAnalysis(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockAnalysisRepository) GetLatestDependencyAnalysis(ctx context.Context, projectID string) (*models.DependencyAnalysis, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DependencyAnalysis), args.Error(1)
}

func (m *MockAnalysisRepository) CreateArchitectureValidation(ctx context.Context, validation *models.ArchitectureValidation) error {
	args := m.Called(ctx, validation)
	return args.Error(0)
}

func (m *MockAnalysisRepository) GetArchitectureValidation(ctx context.Context, id string) (*models.ArchitectureValidation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ArchitectureValidation), args.Error(1)
}

func (m *MockAnalysisRepository) UpdateArchitectureValidation(ctx context.Context, validation *models.ArchitectureValidation) error {
	args := m.Called(ctx, validation)
	return args.Error(0)
}

func TestCircularDetectorService_DetectCircularDependencies(t *testing.T) {
	// Setup
	mockProjectRepo := new(MockProjectRepository)
	mockAnalysisRepo := new(MockAnalysisRepository)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	service := NewCircularDetectorService(mockProjectRepo, mockAnalysisRepo, logger)
	
	ctx := context.Background()
	projectID := "test-project-id"
	repoPath := "/test/repo/path"
	
	// Mock project
	project := &models.Project{
		ID:   projectID,
		Name: "Test Project",
		Settings: &models.ProjectSettings{
			ExcludePatterns: []string{"node_modules/**", "dist/**"},
			IncludePatterns: []string{"**/*.json"},
		},
	}
	
	// Setup mocks
	mockProjectRepo.On("GetByID", ctx, projectID).Return(project, nil)
	mockAnalysisRepo.On("Create", ctx, mock.AnythingOfType("*models.DependencyAnalysis")).Return(nil)
	mockAnalysisRepo.On("Update", ctx, mock.AnythingOfType("*models.DependencyAnalysis")).Return(nil)
	
	// Execute
	result, err := service.DetectCircularDependencies(ctx, projectID, repoPath)
	
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, projectID, result.ProjectID)
	assert.Equal(t, models.StatusCompleted, result.Status)
	assert.NotNil(t, result.Results)
	
	// Verify mocks were called
	mockProjectRepo.AssertExpectations(t)
	mockAnalysisRepo.AssertExpectations(t)
}

func TestCircularDetectorService_detectCircularDepsWithDFS(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &CircularDetectorService{logger: logger}
	
	// Test case 1: Simple circular dependency A -> B -> A
	graph1 := &PackageGraph{
		Nodes: map[string]*PackageNode{
			"A": {Name: "A"},
			"B": {Name: "B"},
		},
		Edges: map[string][]string{
			"A": {"B"},
			"B": {"A"},
		},
	}
	
	cycles1 := service.detectCircularDepsWithDFS(graph1)
	assert.Len(t, cycles1, 1)
	assert.Equal(t, 2, cycles1[0].CycleLength)
	
	// Test case 2: No circular dependencies
	graph2 := &PackageGraph{
		Nodes: map[string]*PackageNode{
			"A": {Name: "A"},
			"B": {Name: "B"},
			"C": {Name: "C"},
		},
		Edges: map[string][]string{
			"A": {"B"},
			"B": {"C"},
			"C": {},
		},
	}
	
	cycles2 := service.detectCircularDepsWithDFS(graph2)
	assert.Len(t, cycles2, 0)
	
	// Test case 3: Complex circular dependency A -> B -> C -> A
	graph3 := &PackageGraph{
		Nodes: map[string]*PackageNode{
			"A": {Name: "A"},
			"B": {Name: "B"},
			"C": {Name: "C"},
		},
		Edges: map[string][]string{
			"A": {"B"},
			"B": {"C"},
			"C": {"A"},
		},
	}
	
	cycles3 := service.detectCircularDepsWithDFS(graph3)
	assert.Len(t, cycles3, 1)
	assert.Equal(t, 3, cycles3[0].CycleLength)
}

func TestCircularDetectorService_generateBreakPointSuggestions(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &CircularDetectorService{logger: logger}
	
	cycle := CircularDependency{
		CyclePath:   []string{"libs/business", "libs/ui", "libs/business"},
		CycleLength: 2,
	}
	
	suggestions := service.generateBreakPointSuggestions(cycle)
	
	assert.Len(t, suggestions, 2)
	assert.Equal(t, "libs/business", suggestions[0].PackageName)
	assert.Equal(t, "libs/ui", suggestions[0].ImportToRemove)
	assert.Contains(t, suggestions[0].AlternativeApproach, "Extract shared interface")
}

func TestCircularDetectorService_calculateCircularHealthScore(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &CircularDetectorService{logger: logger}
	
	// Test cases
	tests := []struct {
		name          string
		circularCount int
		totalPackages int
		expectedScore int
	}{
		{"No circular dependencies", 0, 10, 100},
		{"Low circular ratio", 1, 10, 92},  // 8% reduction
		{"High circular ratio", 5, 10, 60}, // 40% reduction
		{"All packages circular", 10, 10, 20}, // 80% reduction
		{"No packages", 0, 0, 100},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.calculateCircularHealthScore(tt.circularCount, tt.totalPackages)
			assert.Equal(t, tt.expectedScore, score)
		})
	}
}

func TestCircularDetectorService_assessBusinessRisk(t *testing.T) {
	// Setup
	logger := logrus.New()
	service := &CircularDetectorService{logger: logger}
	
	// Test cases
	tests := []struct {
		name         string
		cycle        CircularDependency
		expectedRisk string
	}{
		{
			"High risk - involves apps",
			CircularDependency{CyclePath: []string{"apps/frontend", "libs/ui", "apps/frontend"}},
			"High - affects application-level dependencies",
		},
		{
			"Medium risk - long chain",
			CircularDependency{CyclePath: []string{"libs/a", "libs/b", "libs/c", "libs/d", "libs/a"}},
			"Medium - complex dependency chain",
		},
		{
			"Low risk - simple libs",
			CircularDependency{CyclePath: []string{"libs/a", "libs/b", "libs/a"}},
			"Low - library-level circular dependency",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := service.assessBusinessRisk(tt.cycle)
			assert.Equal(t, tt.expectedRisk, risk)
		})
	}
}