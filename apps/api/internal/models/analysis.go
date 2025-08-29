package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DependencyAnalysis represents a dependency analysis result
type DependencyAnalysis struct {
	ID          string                      `json:"id" gorm:"primaryKey"`
	ProjectID   string                      `json:"projectId" gorm:"column:project_id;not null"`
	Status      Status                      `json:"status" gorm:"not null;default:'pending'"`
	StartedAt   time.Time                   `json:"startedAt" gorm:"column:started_at"`
	CompletedAt *time.Time                  `json:"completedAt,omitempty" gorm:"column:completed_at"`
	Results     *DependencyAnalysisResults  `json:"results" gorm:"type:jsonb"`
	Metadata    *AnalysisMetadata           `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time                   `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time                   `json:"updatedAt" gorm:"column:updated_at"`

	// Associations
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// DependencyAnalysisResults contains the results of dependency analysis
type DependencyAnalysisResults struct {
	DuplicateDependencies []DuplicateDependency `json:"duplicateDependencies"`
	VersionConflicts      []VersionConflict     `json:"versionConflicts"`
	UnusedDependencies    []UnusedDependency    `json:"unusedDependencies"`
	CircularDependencies  []CircularDependency  `json:"circularDependencies"`
	BundleImpact          BundleImpactReport    `json:"bundleImpact"`
	Summary               AnalysisSummary       `json:"summary"`
}

// DuplicateDependency represents a duplicate dependency issue
type DuplicateDependency struct {
	PackageName       string    `json:"packageName"`
	Versions          []string  `json:"versions"`
	AffectedPackages  []string  `json:"affectedPackages"`
	EstimatedWaste    string    `json:"estimatedWaste"`
	RiskLevel         RiskLevel `json:"riskLevel"`
	Recommendation    string    `json:"recommendation"`
	MigrationSteps    []string  `json:"migrationSteps"`
}

// VersionConflict represents a version conflict issue
type VersionConflict struct {
	PackageName          string               `json:"packageName"`
	ConflictingVersions  []ConflictingVersion `json:"conflictingVersions"`
	RiskLevel            RiskLevel            `json:"riskLevel"`
	Resolution           string               `json:"resolution"`
	Impact               string               `json:"impact"`
}

// ConflictingVersion represents a conflicting version
type ConflictingVersion struct {
	Version    string   `json:"version"`
	Packages   []string `json:"packages"`
	IsBreaking bool     `json:"isBreaking"`
}

// UnusedDependency represents an unused dependency
type UnusedDependency struct {
	PackageName string     `json:"packageName"`
	Version     string     `json:"version"`
	PackagePath string     `json:"packagePath"`
	SizeImpact  string     `json:"sizeImpact"`
	LastUsed    *time.Time `json:"lastUsed,omitempty"`
	Confidence  float64    `json:"confidence"`
}

// CircularDependency represents a circular dependency issue
type CircularDependency struct {
	Cycle    []string `json:"cycle"`
	Type     string   `json:"type"` // "direct" or "indirect"
	Severity Severity `json:"severity"`
	Impact   string   `json:"impact"`
}

// BundleImpactReport represents the bundle impact analysis
type BundleImpactReport struct {
	TotalSize        string            `json:"totalSize"`
	DuplicateSize    string            `json:"duplicateSize"`
	UnusedSize       string            `json:"unusedSize"`
	PotentialSavings string            `json:"potentialSavings"`
	Breakdown        []BundleBreakdown `json:"breakdown"`
}

// BundleBreakdown represents a breakdown of bundle impact
type BundleBreakdown struct {
	PackageName string  `json:"packageName"`
	Size        string  `json:"size"`
	Percentage  float64 `json:"percentage"`
	Duplicates  int     `json:"duplicates"`
}

// AnalysisSummary represents a summary of the analysis
type AnalysisSummary struct {
	TotalPackages int     `json:"totalPackages"`
	DuplicateCount int    `json:"duplicateCount"`
	ConflictCount  int    `json:"conflictCount"`
	UnusedCount    int    `json:"unusedCount"`
	CircularCount  int    `json:"circularCount"`
	HealthScore    float64 `json:"healthScore"`
}

// ArchitectureValidation represents an architecture validation result
type ArchitectureValidation struct {
	ID          string                          `json:"id" gorm:"primaryKey"`
	ProjectID   string                          `json:"projectId" gorm:"column:project_id;not null"`
	Status      Status                          `json:"status" gorm:"not null;default:'pending'"`
	StartedAt   time.Time                       `json:"startedAt" gorm:"column:started_at"`
	CompletedAt *time.Time                      `json:"completedAt,omitempty" gorm:"column:completed_at"`
	Results     *ArchitectureValidationResults  `json:"results" gorm:"type:jsonb"`
	Metadata    *AnalysisMetadata               `json:"metadata" gorm:"type:jsonb"`
	CreatedAt   time.Time                       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time                       `json:"updatedAt" gorm:"column:updated_at"`

	// Associations
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// ArchitectureValidationResults contains the results of architecture validation
type ArchitectureValidationResults struct {
	Violations           []ArchitectureViolation `json:"violations"`
	LayerCompliance      []LayerCompliance       `json:"layerCompliance"`
	CircularDependencies []CircularDependency    `json:"circularDependencies"`
	Summary              ValidationSummary       `json:"summary"`
}

// ArchitectureViolation represents an architecture rule violation
type ArchitectureViolation struct {
	RuleName         string   `json:"ruleName"`
	Severity         Severity `json:"severity"`
	Description      string   `json:"description"`
	ViolatingFile    string   `json:"violatingFile"`
	ViolatingImport  string   `json:"violatingImport"`
	ExpectedLayer    string   `json:"expectedLayer"`
	ActualLayer      string   `json:"actualLayer"`
	Suggestion       string   `json:"suggestion"`
}

// LayerCompliance represents layer compliance metrics
type LayerCompliance struct {
	LayerName            string  `json:"layerName"`
	TotalFiles           int     `json:"totalFiles"`
	CompliantFiles       int     `json:"compliantFiles"`
	ViolationCount       int     `json:"violationCount"`
	CompliancePercentage float64 `json:"compliancePercentage"`
}

// ValidationSummary represents a summary of validation results
type ValidationSummary struct {
	TotalViolations     int     `json:"totalViolations"`
	CriticalViolations  int     `json:"criticalViolations"`
	WarningViolations   int     `json:"warningViolations"`
	LayersAnalyzed      int     `json:"layersAnalyzed"`
	OverallCompliance   float64 `json:"overallCompliance"`
}

// AnalysisMetadata contains metadata about the analysis
type AnalysisMetadata struct {
	Version          string                 `json:"version"`
	Duration         int64                  `json:"duration"` // in milliseconds
	FilesProcessed   int                    `json:"filesProcessed"`
	PackagesAnalyzed int                    `json:"packagesAnalyzed"`
	Configuration    map[string]interface{} `json:"configuration"`
	Environment      AnalysisEnvironment    `json:"environment"`
}

// AnalysisEnvironment contains environment information
type AnalysisEnvironment struct {
	NodeVersion  string `json:"nodeVersion"`
	Platform     string `json:"platform"`
	MemoryUsage  string `json:"memoryUsage"`
	CPUUsage     string `json:"cpuUsage"`
}

// HealthScore represents overall health metrics
type HealthScore struct {
	ID             string          `json:"id" gorm:"primaryKey"`
	ProjectID      string          `json:"projectId" gorm:"column:project_id;not null"`
	Overall        float64         `json:"overall"`
	Dependencies   float64         `json:"dependencies"`
	Architecture   float64         `json:"architecture"`
	Maintainability float64        `json:"maintainability"`
	Security       float64         `json:"security"`
	Performance    float64         `json:"performance"`
	LastUpdated    time.Time       `json:"lastUpdated" gorm:"column:last_updated"`
	Trend          string          `json:"trend"` // "improving", "stable", "declining"
	Factors        []HealthFactor  `json:"factors" gorm:"type:jsonb"`
	CreatedAt      time.Time       `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      time.Time       `json:"updatedAt" gorm:"column:updated_at"`

	// Associations
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// HealthFactor represents a factor contributing to health score
type HealthFactor struct {
	Name            string   `json:"name"`
	Score           float64  `json:"score"`
	Weight          float64  `json:"weight"`
	Description     string   `json:"description"`
	Recommendations []string `json:"recommendations"`
}

// BeforeCreate generates UUIDs for analysis models before creating
func (da *DependencyAnalysis) BeforeCreate(tx *gorm.DB) error {
	if da.ID == "" {
		da.ID = uuid.New().String()
	}
	if da.StartedAt.IsZero() {
		da.StartedAt = time.Now()
	}
	return nil
}

func (av *ArchitectureValidation) BeforeCreate(tx *gorm.DB) error {
	if av.ID == "" {
		av.ID = uuid.New().String()
	}
	if av.StartedAt.IsZero() {
		av.StartedAt = time.Now()
	}
	return nil
}

func (hs *HealthScore) BeforeCreate(tx *gorm.DB) error {
	if hs.ID == "" {
		hs.ID = uuid.New().String()
	}
	if hs.LastUpdated.IsZero() {
		hs.LastUpdated = time.Now()
	}
	return nil
}

// TableName methods
func (DependencyAnalysis) TableName() string {
	return "dependency_analyses"
}

func (ArchitectureValidation) TableName() string {
	return "architecture_validations"
}

func (HealthScore) TableName() string {
	return "health_scores"
}

// PackageJSONAnalysis represents enhanced package.json analysis results
type PackageJSONAnalysis struct {
	ID              string                    `json:"id" gorm:"primaryKey"`
	ProjectID       string                    `json:"projectId" gorm:"column:project_id;not null"`
	Status          Status                    `json:"status" gorm:"not null;default:'pending'"`
	StartedAt       time.Time                 `json:"startedAt" gorm:"column:started_at"`
	CompletedAt     *time.Time                `json:"completedAt,omitempty" gorm:"column:completed_at"`
	Results         *PackageJSONAnalysisResults `json:"results" gorm:"type:jsonb"`
	WorkspaceConfig *WorkspaceAnalysisResults `json:"workspaceConfig" gorm:"type:jsonb"`
	Metadata        *AnalysisMetadata         `json:"metadata" gorm:"type:jsonb"`
	CreatedAt       time.Time                 `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt       time.Time                 `json:"updatedAt" gorm:"column:updated_at"`

	// Associations
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// PackageJSONAnalysisResults contains comprehensive package.json analysis results
type PackageJSONAnalysisResults struct {
	WorkspaceDiscovery    WorkspaceDiscoveryResults `json:"workspaceDiscovery"`
	PackageAnalysis       PackageAnalysisResults    `json:"packageAnalysis"`
	DependencyTree        *DependencyTreeResults    `json:"dependencyTree,omitempty"`
	VersionConflicts      []EnhancedVersionConflict `json:"versionConflicts"`
	DuplicateDependencies []DuplicateDependencyInfo `json:"duplicateDependencies"`
	UnusedDependencies    []UnusedDependencyInfo    `json:"unusedDependencies"`
	SecurityVulnerabilities []SecurityVulnerability `json:"securityVulnerabilities"`
	PerformanceMetrics    PerformanceMetrics        `json:"performanceMetrics"`
	Summary               PackageJSONSummary        `json:"summary"`
}

// WorkspaceDiscoveryResults contains results of workspace discovery
type WorkspaceDiscoveryResults struct {
	TotalWorkspaces       int                       `json:"totalWorkspaces"`
	WorkspaceTypes        []WorkspaceTypeInfo       `json:"workspaceTypes"`
	ConflictingWorkspaces []WorkspaceConflictInfo   `json:"conflictingWorkspaces"`
	ResolvedWorkspaces    []ResolvedWorkspaceInfo   `json:"resolvedWorkspaces"`
	DiscoveryTime         int64                     `json:"discoveryTime"` // milliseconds
}

// WorkspaceTypeInfo contains information about workspace types found
type WorkspaceTypeInfo struct {
	Type          string   `json:"type"` // "npm", "pnpm", "lerna", "nx", "yarn"
	Count         int      `json:"count"`
	ConfigFiles   []string `json:"configFiles"`
	PackageCount  int      `json:"packageCount"`
	Priority      int      `json:"priority"`
}

// WorkspaceConflictInfo contains information about conflicting workspace configurations
type WorkspaceConflictInfo struct {
	RootPath          string   `json:"rootPath"`
	ConflictingTypes  []string `json:"conflictingTypes"`
	ConfigFiles       []string `json:"configFiles"`
	ResolutionApplied string   `json:"resolutionApplied"`
	ResolvedTo        string   `json:"resolvedTo"`
}

// ResolvedWorkspaceInfo contains information about resolved workspace
type ResolvedWorkspaceInfo struct {
	Type            string   `json:"type"`
	RootPath        string   `json:"rootPath"`
	ConfigPath      string   `json:"configPath"`
	PackagePatterns []string `json:"packagePatterns"`
	PackageCount    int      `json:"packageCount"`
	PackageManager  string   `json:"packageManager"`
}

// PackageAnalysisResults contains results of package analysis
type PackageAnalysisResults struct {
	TotalPackages         int                     `json:"totalPackages"`
	WorkspacePackages     int                     `json:"workspacePackages"`
	ExternalDependencies  int                     `json:"externalDependencies"`
	PackagesByType        map[string]int          `json:"packagesByType"`
	LargestPackages       []PackageSizeInfo       `json:"largestPackages"`
	DependencyDistribution DependencyDistribution `json:"dependencyDistribution"`
	VersionRangeAnalysis  VersionRangeAnalysis    `json:"versionRangeAnalysis"`
}

// PackageSizeInfo contains information about package sizes
type PackageSizeInfo struct {
	Name              string `json:"name"`
	Path              string `json:"path"`
	EstimatedSize     string `json:"estimatedSize"`
	DependencyCount   int    `json:"dependencyCount"`
	DevDependencyCount int   `json:"devDependencyCount"`
	IsWorkspace       bool   `json:"isWorkspace"`
}

// DependencyDistribution contains distribution analysis of dependencies
type DependencyDistribution struct {
	MostUsedDependencies []DependencyUsageInfo `json:"mostUsedDependencies"`
	UniqueDependencies   int                   `json:"uniqueDependencies"`
	AverageDepPerPackage float64               `json:"averageDepPerPackage"`
	MedianDepPerPackage  int                   `json:"medianDepPerPackage"`
	MaxDepPerPackage     int                   `json:"maxDepPerPackage"`
}

// DependencyUsageInfo contains usage information for dependencies
type DependencyUsageInfo struct {
	Name           string   `json:"name"`
	UsageCount     int      `json:"usageCount"`
	Versions       []string `json:"versions"`
	UsedByPackages []string `json:"usedByPackages"`
	IsDevDep       bool     `json:"isDevDep"`
}

// VersionRangeAnalysis contains analysis of version ranges used
type VersionRangeAnalysis struct {
	RangeTypes         map[string]int          `json:"rangeTypes"` // ^, ~, >=, etc.
	ExactVersions      int                     `json:"exactVersions"`
	FlexibleRanges     int                     `json:"flexibleRanges"`
	RestrictiveRanges  int                     `json:"restrictiveRanges"`
	ProblematicRanges  []ProblematicRangeInfo  `json:"problematicRanges"`
}

// ProblematicRangeInfo contains information about problematic version ranges
type ProblematicRangeInfo struct {
	Package     string   `json:"package"`
	Range       string   `json:"range"`
	Issue       string   `json:"issue"`
	Suggestions []string `json:"suggestions"`
	Severity    string   `json:"severity"`
}

// DependencyTreeResults contains dependency tree analysis results
type DependencyTreeResults struct {
	TotalNodes          int                      `json:"totalNodes"`
	MaxDepth            int                      `json:"maxDepth"`
	CircularDependencies []CircularDependencyInfo `json:"circularDependencies"`
	OrphanedPackages    []OrphanedPackageInfo    `json:"orphanedPackages"`
	TreeComplexity      TreeComplexityMetrics    `json:"treeComplexity"`
	ResolutionTime      int64                    `json:"resolutionTime"` // milliseconds
	CacheHitRate        float64                  `json:"cacheHitRate"`
}

// CircularDependencyInfo contains detailed information about circular dependencies
type CircularDependencyInfo struct {
	Cycle         []string `json:"cycle"`
	CycleType     string   `json:"cycleType"` // "direct", "indirect"
	Severity      string   `json:"severity"`
	Impact        string   `json:"impact"`
	BreakingSuggestion string `json:"breakingSuggestion"`
}

// OrphanedPackageInfo contains information about orphaned packages
type OrphanedPackageInfo struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Version     string   `json:"version"`
	LastUsed    *time.Time `json:"lastUsed,omitempty"`
	Suggestions []string `json:"suggestions"`
}

// TreeComplexityMetrics contains metrics about dependency tree complexity
type TreeComplexityMetrics struct {
	ComplexityScore    float64 `json:"complexityScore"`
	BranchingFactor    float64 `json:"branchingFactor"`
	AverageDepth       float64 `json:"averageDepth"`
	DeepestPath        []string `json:"deepestPath"`
	HighestFanOut      string  `json:"highestFanOut"`
	FanOutCount        int     `json:"fanOutCount"`
}

// EnhancedVersionConflict extends the basic version conflict with more details
type EnhancedVersionConflict struct {
	VersionConflict
	ConflictSeverity    string                  `json:"conflictSeverity"`
	AutoResolvable      bool                    `json:"autoResolvable"`
	ResolutionOptions   []ResolutionOptionInfo  `json:"resolutionOptions"`
	RecommendedFix      *ResolutionOptionInfo   `json:"recommendedFix,omitempty"`
	MigrationComplexity string                  `json:"migrationComplexity"`
	TestingRequirements []string                `json:"testingRequirements"`
}

// ResolutionOptionInfo contains information about conflict resolution options
type ResolutionOptionInfo struct {
	Strategy            string   `json:"strategy"`
	TargetVersion       string   `json:"targetVersion"`
	AffectedPackages    []string `json:"affectedPackages"`
	RiskLevel          string   `json:"riskLevel"`
	EstimatedEffort    string   `json:"estimatedEffort"`
	MigrationSteps     []string `json:"migrationSteps"`
	RollbackComplexity string   `json:"rollbackComplexity"`
}

// DuplicateDependencyInfo contains enhanced information about duplicate dependencies
type DuplicateDependencyInfo struct {
	DuplicateDependency
	WasteAnalysis       ResourceWasteAnalysis   `json:"wasteAnalysis"`
	ConsolidationPlan   ConsolidationPlanInfo   `json:"consolidationPlan"`
	ImpactAssessment    ImpactAssessmentInfo    `json:"impactAssessment"`
}

// ResourceWasteAnalysis contains detailed analysis of resource waste
type ResourceWasteAnalysis struct {
	DiskSpaceWaste      string  `json:"diskSpaceWaste"`
	BundleSizeWaste     string  `json:"bundleSizeWaste"`
	InstallTimeWaste    string  `json:"installTimeWaste"`
	MemoryWaste         string  `json:"memoryWaste"`
	NetworkWaste        string  `json:"networkWaste"`
	WastePercentage     float64 `json:"wastePercentage"`
	AnnualizedCost      string  `json:"annualizedCost,omitempty"`
}

// ConsolidationPlanInfo contains detailed consolidation planning information
type ConsolidationPlanInfo struct {
	RecommendedVersion  string                `json:"recommendedVersion"`
	MigrationSteps      []MigrationStepInfo   `json:"migrationSteps"`
	Prerequisites       []string              `json:"prerequisites"`
	RiskMitigation      []string              `json:"riskMitigation"`
	ValidationSteps     []string              `json:"validationSteps"`
	RollbackPlan        []string              `json:"rollbackPlan"`
	EstimatedTimeHours  float64               `json:"estimatedTimeHours"`
	ResourceRequirements []string             `json:"resourceRequirements"`
}

// MigrationStepInfo contains detailed migration step information
type MigrationStepInfo struct {
	StepNumber      int      `json:"stepNumber"`
	Description     string   `json:"description"`
	Command         string   `json:"command"`
	ExpectedOutput  string   `json:"expectedOutput"`
	VerificationCmd string   `json:"verificationCmd"`
	RollbackCmd     string   `json:"rollbackCmd,omitempty"`
	Risks           []string `json:"risks"`
	Prerequisites   []string `json:"prerequisites"`
}

// ImpactAssessmentInfo contains impact assessment information
type ImpactAssessmentInfo struct {
	BusinessImpact      string   `json:"businessImpact"`
	TechnicalImpact     string   `json:"technicalImpact"`
	UserImpact          string   `json:"userImpact"`
	SecurityImpact      string   `json:"securityImpact"`
	PerformanceImpact   string   `json:"performanceImpact"`
	AffectedComponents  []string `json:"affectedComponents"`
	DownstreamEffects   []string `json:"downstreamEffects"`
}

// UnusedDependencyInfo contains enhanced information about unused dependencies
type UnusedDependencyInfo struct {
	UnusedDependency
	UsageAnalysisDetails UsageAnalysisDetails    `json:"usageAnalysisDetails"`
	RemovalPlan         RemovalPlanInfo         `json:"removalPlan"`
	SafetyAssessment    SafetyAssessmentInfo    `json:"safetyAssessment"`
}

// UsageAnalysisDetails contains detailed usage analysis
type UsageAnalysisDetails struct {
	StaticAnalysis      StaticAnalysisResults   `json:"staticAnalysis"`
	DynamicAnalysis     *DynamicAnalysisResults `json:"dynamicAnalysis,omitempty"`
	ImportPatterns      []ImportPatternInfo     `json:"importPatterns"`
	UsageLocations      []UsageLocationInfo     `json:"usageLocations"`
	LastAccessTime      *time.Time              `json:"lastAccessTime,omitempty"`
	UsageFrequency      UsageFrequencyInfo      `json:"usageFrequency"`
}

// StaticAnalysisResults contains results of static code analysis
type StaticAnalysisResults struct {
	ImportStatements    []string `json:"importStatements"`
	RequireStatements   []string `json:"requireStatements"`
	TypeImports         []string `json:"typeImports"`
	ConfigReferences    []string `json:"configReferences"`
	TestReferences      []string `json:"testReferences"`
	DocumentationRefs   []string `json:"documentationRefs"`
}

// DynamicAnalysisResults contains results of dynamic analysis (if available)
type DynamicAnalysisResults struct {
	RuntimeUsage        []RuntimeUsageInfo      `json:"runtimeUsage"`
	LoadTime            int64                   `json:"loadTime"` // milliseconds
	MemoryFootprint     string                  `json:"memoryFootprint"`
	NetworkRequests     []NetworkRequestInfo    `json:"networkRequests"`
}

// ImportPatternInfo contains information about import patterns
type ImportPatternInfo struct {
	Pattern         string   `json:"pattern"`
	Count           int      `json:"count"`
	Files           []string `json:"files"`
	IsConditional   bool     `json:"isConditional"`
	IsTypeOnly      bool     `json:"isTypeOnly"`
	IsDevContext    bool     `json:"isDevContext"`
}

// UsageLocationInfo contains information about where dependencies are used
type UsageLocationInfo struct {
	File            string   `json:"file"`
	LineNumber      int      `json:"lineNumber"`
	Context         string   `json:"context"`
	UsageType       string   `json:"usageType"` // "import", "require", "type", "config"
	IsConditional   bool     `json:"isConditional"`
	Criticality     string   `json:"criticality"`
}

// UsageFrequencyInfo contains frequency analysis
type UsageFrequencyInfo struct {
	TotalOccurrences    int                         `json:"totalOccurrences"`
	FileCount           int                         `json:"fileCount"`
	AveragePerFile      float64                     `json:"averagePerFile"`
	FrequencyByType     map[string]int              `json:"frequencyByType"`
	TrendAnalysis       *UsageTrendInfo             `json:"trendAnalysis,omitempty"`
}

// RuntimeUsageInfo contains runtime usage information
type RuntimeUsageInfo struct {
	Timestamp       time.Time               `json:"timestamp"`
	ExecutionPath   string                  `json:"executionPath"`
	Context         map[string]interface{}  `json:"context"`
	Performance     PerformanceInfo         `json:"performance"`
}

// NetworkRequestInfo contains network request information
type NetworkRequestInfo struct {
	URL             string                  `json:"url"`
	Method          string                  `json:"method"`
	Timestamp       time.Time               `json:"timestamp"`
	ResponseTime    int64                   `json:"responseTime"`
	DataTransferred string                  `json:"dataTransferred"`
}

// UsageTrendInfo contains trend analysis information
type UsageTrendInfo struct {
	TrendDirection  string                  `json:"trendDirection"` // "increasing", "decreasing", "stable"
	ChangeRate      float64                 `json:"changeRate"`
	LastIncreaseAt  *time.Time              `json:"lastIncreaseAt,omitempty"`
	LastDecreaseAt  *time.Time              `json:"lastDecreaseAt,omitempty"`
	PredictedUsage  *PredictedUsageInfo     `json:"predictedUsage,omitempty"`
}

// PredictedUsageInfo contains predicted usage information
type PredictedUsageInfo struct {
	NextMonthUsage      float64     `json:"nextMonthUsage"`
	Confidence          float64     `json:"confidence"`
	FactorsConsidered   []string    `json:"factorsConsidered"`
}

// RemovalPlanInfo contains dependency removal planning information
type RemovalPlanInfo struct {
	CanRemove           bool                    `json:"canRemove"`
	RemovalSteps        []RemovalStepInfo       `json:"removalSteps"`
	Prerequisites       []string                `json:"prerequisites"`
	Alternatives        []AlternativeInfo       `json:"alternatives"`
	RollbackPlan        []string                `json:"rollbackPlan"`
	TestingRequirements []string                `json:"testingRequirements"`
}

// RemovalStepInfo contains removal step information
type RemovalStepInfo struct {
	StepNumber      int      `json:"stepNumber"`
	Description     string   `json:"description"`
	Command         string   `json:"command"`
	AffectedFiles   []string `json:"affectedFiles"`
	BackupRequired  bool     `json:"backupRequired"`
	RiskLevel       string   `json:"riskLevel"`
}

// AlternativeInfo contains information about alternative dependencies
type AlternativeInfo struct {
	Name            string   `json:"name"`
	Reason          string   `json:"reason"`
	Compatibility   string   `json:"compatibility"`
	MigrationNotes  []string `json:"migrationNotes"`
	Pros            []string `json:"pros"`
	Cons            []string `json:"cons"`
}

// SafetyAssessmentInfo contains safety assessment for dependency removal
type SafetyAssessmentInfo struct {
	SafetyScore         float64                 `json:"safetyScore"` // 0-100
	RiskFactors         []RiskFactorInfo        `json:"riskFactors"`
	SafetyChecks        []SafetyCheckInfo       `json:"safetyChecks"`
	RecommendedAction   string                  `json:"recommendedAction"`
	MonitoringPlan      []string                `json:"monitoringPlan"`
}

// RiskFactorInfo contains risk factor information
type RiskFactorInfo struct {
	Factor          string  `json:"factor"`
	Severity        string  `json:"severity"`
	Likelihood      string  `json:"likelihood"`
	Impact          string  `json:"impact"`
	Mitigation      string  `json:"mitigation"`
	Weight          float64 `json:"weight"`
}

// SafetyCheckInfo contains safety check information
type SafetyCheckInfo struct {
	Check           string  `json:"check"`
	Status          string  `json:"status"` // "passed", "failed", "warning"
	Details         string  `json:"details"`
	Recommendation  string  `json:"recommendation"`
}

// SecurityVulnerability contains security vulnerability information
type SecurityVulnerability struct {
	ID                  string                  `json:"id"`
	PackageName         string                  `json:"packageName"`
	Version             string                  `json:"version"`
	Severity            string                  `json:"severity"`
	CVSS                float64                 `json:"cvss,omitempty"`
	Description         string                  `json:"description"`
	References          []string                `json:"references"`
	Patches             []PatchInfo             `json:"patches"`
	Workarounds         []WorkaroundInfo        `json:"workarounds"`
	AffectedPaths       []string                `json:"affectedPaths"`
	ExploitabilityScore float64                 `json:"exploitabilityScore"`
	ImpactScore         float64                 `json:"impactScore"`
}

// PatchInfo contains patch information for vulnerabilities
type PatchInfo struct {
	Version         string  `json:"version"`
	ReleaseDate     string  `json:"releaseDate"`
	BreakingChanges bool    `json:"breakingChanges"`
	ChangelogURL    string  `json:"changelogUrl"`
}

// WorkaroundInfo contains workaround information
type WorkaroundInfo struct {
	Description     string   `json:"description"`
	Steps           []string `json:"steps"`
	Limitations     []string `json:"limitations"`
	Effectiveness   string   `json:"effectiveness"`
}

// PerformanceMetrics contains performance metrics for the analysis
type PerformanceMetrics struct {
	TotalAnalysisTime   int64                   `json:"totalAnalysisTime"` // milliseconds
	PhaseTimings        map[string]int64        `json:"phaseTimings"`
	MemoryUsage         MemoryUsageInfo         `json:"memoryUsage"`
	CachePerformance    CachePerformanceInfo    `json:"cachePerformance"`
	ConcurrencyMetrics  ConcurrencyMetricsInfo  `json:"concurrencyMetrics"`
	ResourceUtilization ResourceUtilizationInfo `json:"resourceUtilization"`
}

// MemoryUsageInfo contains memory usage information
type MemoryUsageInfo struct {
	PeakMemory      string  `json:"peakMemory"`
	AverageMemory   string  `json:"averageMemory"`
	StartMemory     string  `json:"startMemory"`
	EndMemory       string  `json:"endMemory"`
	GCPressure      float64 `json:"gcPressure"`
}

// CachePerformanceInfo contains cache performance information
type CachePerformanceInfo struct {
	HitRate         float64         `json:"hitRate"`
	MissRate        float64         `json:"missRate"`
	HitCount        int64           `json:"hitCount"`
	MissCount       int64           `json:"missCount"`
	CacheSize       string          `json:"cacheSize"`
	EvictionCount   int64           `json:"evictionCount"`
	AverageHitTime  int64           `json:"averageHitTime"`
	AverageMissTime int64           `json:"averageMissTime"`
}

// ConcurrencyMetricsInfo contains concurrency metrics
type ConcurrencyMetricsInfo struct {
	MaxGoroutines       int             `json:"maxGoroutines"`
	AverageGoroutines   float64         `json:"averageGoroutines"`
	TotalGoroutines     int64           `json:"totalGoroutines"`
	DeadlockCount       int             `json:"deadlockCount"`
	ContentionEvents    int             `json:"contentionEvents"`
	SchedulingLatency   int64           `json:"schedulingLatency"`
}

// ResourceUtilizationInfo contains resource utilization information
type ResourceUtilizationInfo struct {
	CPUUsage        float64         `json:"cpuUsage"`
	DiskIO          DiskIOInfo      `json:"diskIO"`
	NetworkIO       NetworkIOInfo   `json:"networkIO"`
	FileDescriptors int             `json:"fileDescriptors"`
	ThreadCount     int             `json:"threadCount"`
}

// DiskIOInfo contains disk I/O information
type DiskIOInfo struct {
	ReadBytes       int64   `json:"readBytes"`
	WriteBytes      int64   `json:"writeBytes"`
	ReadOps         int64   `json:"readOps"`
	WriteOps        int64   `json:"writeOps"`
	ReadLatency     int64   `json:"readLatency"`
	WriteLatency    int64   `json:"writeLatency"`
}

// NetworkIOInfo contains network I/O information
type NetworkIOInfo struct {
	BytesSent       int64   `json:"bytesSent"`
	BytesReceived   int64   `json:"bytesReceived"`
	PacketsSent     int64   `json:"packetsSent"`
	PacketsReceived int64   `json:"packetsReceived"`
	Connections     int     `json:"connections"`
	Errors          int     `json:"errors"`
}

// PerformanceInfo contains performance information for individual operations
type PerformanceInfo struct {
	Duration        int64           `json:"duration"` // microseconds
	CPUTime         int64           `json:"cpuTime"`
	MemoryAllocated string          `json:"memoryAllocated"`
	IOOperations    int             `json:"ioOperations"`
}

// PackageJSONSummary contains summary information for package.json analysis
type PackageJSONSummary struct {
	TotalPackages           int                     `json:"totalPackages"`
	WorkspacePackages       int                     `json:"workspacePackages"`
	ExternalDependencies    int                     `json:"externalDependencies"`
	TotalDependencies       int                     `json:"totalDependencies"`
	DuplicateCount          int                     `json:"duplicateCount"`
	ConflictCount           int                     `json:"conflictCount"`
	UnusedCount             int                     `json:"unusedCount"`
	VulnerabilityCount      int                     `json:"vulnerabilityCount"`
	CriticalVulnerabilities int                     `json:"criticalVulnerabilities"`
	HighVulnerabilities     int                     `json:"highVulnerabilities"`
	OverallHealthScore      float64                 `json:"overallHealthScore"`
	SecurityScore           float64                 `json:"securityScore"`
	MaintenanceScore        float64                 `json:"maintenanceScore"`
	PerformanceScore        float64                 `json:"performanceScore"`
	Recommendations         []RecommendationInfo    `json:"recommendations"`
	KeyMetrics              map[string]interface{}  `json:"keyMetrics"`
}

// RecommendationInfo contains recommendation information
type RecommendationInfo struct {
	Type            string   `json:"type"`
	Priority        string   `json:"priority"`
	Category        string   `json:"category"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	ActionItems     []string `json:"actionItems"`
	EstimatedEffort string   `json:"estimatedEffort"`
	ExpectedBenefit string   `json:"expectedBenefit"`
	RiskLevel       string   `json:"riskLevel"`
}

// WorkspaceAnalysisResults contains comprehensive workspace analysis
type WorkspaceAnalysisResults struct {
	DiscoveryResults    WorkspaceDiscoveryResults   `json:"discoveryResults"`
	ConfigurationHealth WorkspaceHealthInfo         `json:"configurationHealth"`
	BestPractices       BestPracticesAnalysis       `json:"bestPractices"`
	OptimizationSuggestions []OptimizationSuggestion `json:"optimizationSuggestions"`
}

// WorkspaceHealthInfo contains workspace health information
type WorkspaceHealthInfo struct {
	OverallScore            float64                     `json:"overallScore"`
	ConfigurationScore      float64                     `json:"configurationScore"`
	ConsistencyScore        float64                     `json:"consistencyScore"`
	PerformanceScore        float64                     `json:"performanceScore"`
	MaintenabilityScore     float64                     `json:"maintainabilityScore"`
	Issues                  []WorkspaceIssueInfo        `json:"issues"`
	Strengths               []string                    `json:"strengths"`
}

// WorkspaceIssueInfo contains workspace issue information
type WorkspaceIssueInfo struct {
	Type            string   `json:"type"`
	Severity        string   `json:"severity"`
	Description     string   `json:"description"`
	Location        string   `json:"location"`
	Impact          string   `json:"impact"`
	Resolution      string   `json:"resolution"`
	AutoFixable     bool     `json:"autoFixable"`
}

// BestPracticesAnalysis contains best practices analysis
type BestPracticesAnalysis struct {
	Compliance          float64                     `json:"compliance"` // 0-100
	PassedChecks        []BestPracticeCheckInfo     `json:"passedChecks"`
	FailedChecks        []BestPracticeCheckInfo     `json:"failedChecks"`
	Recommendations     []BestPracticeRecommendation `json:"recommendations"`
	IndustryComparison  IndustryComparisonInfo      `json:"industryComparison"`
}

// BestPracticeCheckInfo contains best practice check information
type BestPracticeCheckInfo struct {
	Check           string  `json:"check"`
	Category        string  `json:"category"`
	Status          string  `json:"status"`
	Score           float64 `json:"score"`
	Description     string  `json:"description"`
	Importance      string  `json:"importance"`
	References      []string `json:"references"`
}

// BestPracticeRecommendation contains best practice recommendations
type BestPracticeRecommendation struct {
	Practice        string   `json:"practice"`
	Current         string   `json:"current"`
	Recommended     string   `json:"recommended"`
	Benefits        []string `json:"benefits"`
	ImplementationSteps []string `json:"implementationSteps"`
	Difficulty      string   `json:"difficulty"`
	Priority        int      `json:"priority"`
}

// IndustryComparisonInfo contains industry comparison information
type IndustryComparisonInfo struct {
	IndustryAverage     float64                     `json:"industryAverage"`
	Percentile          float64                     `json:"percentile"`
	ComparisonAreas     []ComparisonAreaInfo        `json:"comparisonAreas"`
	BenchmarkSources    []string                    `json:"benchmarkSources"`
}

// ComparisonAreaInfo contains comparison area information
type ComparisonAreaInfo struct {
	Area            string  `json:"area"`
	YourScore       float64 `json:"yourScore"`
	IndustryAverage float64 `json:"industryAverage"`
	TopPercentile   float64 `json:"topPercentile"`
	Ranking         string  `json:"ranking"`
}

// OptimizationSuggestion contains workspace optimization suggestions
type OptimizationSuggestion struct {
	Type                string                  `json:"type"`
	Title               string                  `json:"title"`
	Description         string                  `json:"description"`
	Category            string                  `json:"category"`
	Priority            int                     `json:"priority"`
	EstimatedImpact     string                  `json:"estimatedImpact"`
	ImplementationPlan  ImplementationPlanInfo  `json:"implementationPlan"`
	ExpectedBenefits    []string                `json:"expectedBenefits"`
	Risks               []string                `json:"risks"`
	Prerequisites       []string                `json:"prerequisites"`
}

// ImplementationPlanInfo contains implementation plan information
type ImplementationPlanInfo struct {
	Phases              []ImplementationPhaseInfo   `json:"phases"`
	TotalEstimatedTime  string                      `json:"totalEstimatedTime"`
	ResourceRequirements []string                   `json:"resourceRequirements"`
	SuccessCriteria     []string                    `json:"successCriteria"`
	RollbackStrategy    string                      `json:"rollbackStrategy"`
}

// ImplementationPhaseInfo contains implementation phase information
type ImplementationPhaseInfo struct {
	Phase           int      `json:"phase"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Tasks           []string `json:"tasks"`
	EstimatedTime   string   `json:"estimatedTime"`
	Dependencies    []int    `json:"dependencies"`
	Deliverables    []string `json:"deliverables"`
	RiskMitigation  []string `json:"riskMitigation"`
}

func (pja *PackageJSONAnalysis) BeforeCreate(tx *gorm.DB) error {
	if pja.ID == "" {
		pja.ID = uuid.New().String()
	}
	if pja.StartedAt.IsZero() {
		pja.StartedAt = time.Now()
	}
	return nil
}

func (PackageJSONAnalysis) TableName() string {
	return "package_json_analyses"
}