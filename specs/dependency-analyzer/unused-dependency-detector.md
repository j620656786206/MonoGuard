# Unused Dependency Detector Technical Specification

## Overview

The Unused Dependency Detector identifies dependencies declared in package.json files but not actually used in the source code through sophisticated static analysis. This component helps reduce bundle sizes, improve build times, and maintain cleaner dependency lists by providing high-confidence recommendations for safe dependency removal.

## Technical Requirements

### REQ-UDD-001: Static Analysis Engine

**Description:** Implement comprehensive static analysis to detect actual dependency usage across TypeScript/JavaScript codebases with high accuracy.

**Core Data Structures:**
```go
type UnusedAnalysis struct {
    // Analysis metadata
    AnalysisID          string              `json:"analysis_id"`
    WorkspaceRoot       string              `json:"workspace_root"`
    AnalyzedPackages    int                 `json:"analyzed_packages"`
    TotalSourceFiles    int                 `json:"total_source_files"`
    
    // Detected unused dependencies
    UnusedDependencies  []*UnusedDependency `json:"unused_dependencies"`
    PotentiallyUnused   []*UnusedDependency `json:"potentially_unused"`
    SafeToRemove        []*UnusedDependency `json:"safe_to_remove"`
    
    // Usage analysis
    UsagePatterns       []*UsagePattern     `json:"usage_patterns"`
    ImportAnalysis      *ImportAnalysis     `json:"import_analysis"`
    
    // Impact metrics
    CleanupImpact       *CleanupImpact      `json:"cleanup_impact"`
    
    // Analysis metadata
    AnalysisTime        time.Time           `json:"analysis_time"`
    ProcessingTime      time.Duration       `json:"processing_time"`
    ConfidenceLevel     float64             `json:"confidence_level"`
    AnalysisDepth       AnalysisDepth       `json:"analysis_depth"`
}

type UnusedDependency struct {
    // Package identification
    PackageName         string              `json:"package_name"`
    Version             string              `json:"version"`
    DeclaredIn          string              `json:"declared_in"`        // package.json path
    DependencyType      DependencyType      `json:"dependency_type"`
    
    // Usage analysis
    UsageEvidence       *UsageEvidence      `json:"usage_evidence"`
    ConfidenceScore     float64             `json:"confidence_score"`   // 0.0 to 1.0
    RemovalSafety       RemovalSafety       `json:"removal_safety"`
    
    // Impact analysis
    SizeImpact          *SizeImpact         `json:"size_impact"`
    RuntimeImpact       *RuntimeImpact      `json:"runtime_impact"`
    
    // Context and recommendations
    AnalysisContext     *AnalysisContext    `json:"analysis_context"`
    RemovalRisk         *RemovalRisk        `json:"removal_risk"`
    Recommendations     []*RemovalRecommendation `json:"recommendations"`
    
    // Metadata
    DetectedAt          time.Time           `json:"detected_at"`
    LastVerified        time.Time           `json:"last_verified"`
    FalsePositiveRisk   float64             `json:"false_positive_risk"`
}

type UsageEvidence struct {
    // Direct import/require evidence
    DirectImports       []*ImportStatement  `json:"direct_imports"`
    DynamicImports      []*DynamicImport    `json:"dynamic_imports"`
    RequireStatements   []*RequireStatement `json:"require_statements"`
    
    // Indirect usage evidence
    TransitiveUsage     []*TransitiveUsage  `json:"transitive_usage"`
    RuntimeReferences   []*RuntimeReference `json:"runtime_references"`
    ConfigReferences    []*ConfigReference  `json:"config_references"`
    
    // Build tool usage
    BuildToolUsage      []*BuildToolUsage   `json:"build_tool_usage"`
    ScriptUsage         []*ScriptUsage      `json:"script_usage"`
    
    // Evidence strength
    EvidenceStrength    EvidenceStrength    `json:"evidence_strength"`
    EvidenceCount       int                 `json:"evidence_count"`
    LastSeenEvidence    time.Time           `json:"last_seen_evidence"`
}

type EvidenceStrength string

const (
    EvidenceNone        EvidenceStrength = "none"       // No evidence found
    EvidenceWeak        EvidenceStrength = "weak"       // Circumstantial evidence
    EvidenceMedium      EvidenceStrength = "medium"     // Some direct usage
    EvidenceStrong      EvidenceStrength = "strong"     // Clear direct usage
    EvidenceDefinitive  EvidenceStrength = "definitive" // Multiple strong evidences
)

type RemovalSafety string

const (
    SafetyHigh          RemovalSafety = "high"          // Very safe to remove
    SafetyMedium        RemovalSafety = "medium"        // Probably safe with testing
    SafetyLow           RemovalSafety = "low"           // Risky, needs careful review
    SafetyUnsafe        RemovalSafety = "unsafe"        // Do not remove
    SafetyUnknown       RemovalSafety = "unknown"       // Need more analysis
)

type ImportStatement struct {
    // Import details
    ImportPath          string              `json:"import_path"`
    ImportType          ImportType          `json:"import_type"`
    ImportedNames       []string            `json:"imported_names"`   // Named imports
    ImportAlias         string              `json:"import_alias"`     // Default import name
    
    // Location information
    SourceFile          string              `json:"source_file"`
    LineNumber          int                 `json:"line_number"`
    ColumnNumber        int                 `json:"column_number"`
    
    // Context
    IsTypeOnly          bool                `json:"is_type_only"`
    IsConditional       bool                `json:"is_conditional"`   // Inside if/try block
    Usage               *ImportUsage        `json:"usage"`
}

type ImportType string

const (
    ImportTypeESM       ImportType = "esm"              // import statement
    ImportTypeCJS       ImportType = "cjs"              // require()
    ImportTypeDynamic   ImportType = "dynamic"          // import()
    ImportTypeTypeOnly  ImportType = "type_only"        // import type
    ImportTypeJSX       ImportType = "jsx"              // JSX component
)

type ImportUsage struct {
    UsedNames           []string            `json:"used_names"`       // Actually used imports
    UnusedNames         []string            `json:"unused_names"`     // Declared but not used
    UsageLocations      []*UsageLocation    `json:"usage_locations"`
    UsageFrequency      int                 `json:"usage_frequency"`
    IsTreeShakable      bool                `json:"is_tree_shakable"`
}

type UsageLocation struct {
    SourceFile          string              `json:"source_file"`
    LineNumber          int                 `json:"line_number"`
    Context             string              `json:"context"`          // Function/class/block
    UsageType           UsageType           `json:"usage_type"`
}

type UsageType string

const (
    UsageTypeFunctionCall   UsageType = "function_call"
    UsageTypeObjectAccess   UsageType = "object_access"
    UsageTypeTypeAnnotation UsageType = "type_annotation"
    UsageTypeJSXElement     UsageType = "jsx_element"
    UsageTypeDecorator      UsageType = "decorator"
    UsageTypeStringLiteral  UsageType = "string_literal"    // Dynamic references
)
```

### REQ-UDD-002: Advanced Import Detection

**Description:** Detect all forms of dependency usage including dynamic imports, string literals, and build tool configurations.

**Detection Patterns:**
```go
type ImportDetector interface {
    // Detect static imports
    DetectStaticImports(sourceCode string, filePath string) ([]*ImportStatement, error)
    
    // Detect dynamic imports
    DetectDynamicImports(sourceCode string, filePath string) ([]*DynamicImport, error)
    
    // Detect string literal references
    DetectStringReferences(sourceCode string, filePath string, packageName string) ([]*StringReference, error)
    
    // Detect build tool usage
    DetectBuildToolUsage(configFiles []string, packageName string) ([]*BuildToolUsage, error)
    
    // Detect runtime usage patterns
    DetectRuntimeUsage(sourceCode string, packageName string) ([]*RuntimeReference, error)
    
    // Validate detection accuracy
    ValidateDetection(detectedUsage []*UsageEvidence, groundTruth []*UsageEvidence) (*ValidationMetrics, error)
}

type DynamicImport struct {
    // Import expression details
    ImportExpression    string              `json:"import_expression"`
    ResolvedModule      string              `json:"resolved_module"`    // If determinable
    IsDeterministic     bool                `json:"is_deterministic"`
    
    // Context
    SourceFile          string              `json:"source_file"`
    LineNumber          int                 `json:"line_number"`
    ConditionalContext  *ConditionalContext `json:"conditional_context"`
    
    // Analysis
    Confidence          float64             `json:"confidence"`         // 0.0 to 1.0
    PossibleModules     []string            `json:"possible_modules"`   // If non-deterministic
    RequiresRuntime     bool                `json:"requires_runtime"`   // Needs runtime analysis
}

type StringReference struct {
    // Reference details
    StringValue         string              `json:"string_value"`
    ReferenceType       StringReferenceType `json:"reference_type"`
    Context             string              `json:"context"`            // What contains the string
    
    // Location
    SourceFile          string              `json:"source_file"`
    LineNumber          int                 `json:"line_number"`
    
    // Analysis
    IsPackageReference  bool                `json:"is_package_reference"`
    Confidence          float64             `json:"confidence"`
    RequiredAtRuntime   bool                `json:"required_at_runtime"`
}

type StringReferenceType string

const (
    RefTypeModulePath       StringReferenceType = "module_path"        // require("package")
    RefTypeConfigValue      StringReferenceType = "config_value"       // Config files
    RefTypePluginName       StringReferenceType = "plugin_name"        // Build tool plugins
    RefTypeComment          StringReferenceType = "comment"            // Comments
    RefTypeTemplateString   StringReferenceType = "template_string"    // Template literals
    RefTypeJSONValue        StringReferenceType = "json_value"         // JSON files
)

type BuildToolUsage struct {
    // Build tool information
    ToolName            string              `json:"tool_name"`          // webpack, rollup, etc.
    ConfigFile          string              `json:"config_file"`
    
    // Usage details
    UsageType           BuildUsageType      `json:"usage_type"`
    Configuration       map[string]interface{} `json:"configuration"`
    
    // Dependencies
    DependsOnPackage    bool                `json:"depends_on_package"`
    IsEssential         bool                `json:"is_essential"`       // Core to build process
    Alternatives        []string            `json:"alternatives"`       // Alternative packages
    
    // Impact
    RemovalImpact       string              `json:"removal_impact"`     // build, runtime, none
    MigrationPath       []string            `json:"migration_path"`
}

type BuildUsageType string

const (
    BuildUsagePlugin        BuildUsageType = "plugin"
    BuildUsageLoader        BuildUsageType = "loader"
    BuildUsagePreset        BuildUsageType = "preset"
    BuildUsageTransformer   BuildUsageType = "transformer"
    BuildUsageResolver      BuildUsageType = "resolver"
    BuildUsageOptimizer     BuildUsageType = "optimizer"
)

type ConditionalContext struct {
    ConditionType       ConditionType       `json:"condition_type"`
    Condition           string              `json:"condition"`          // The actual condition
    ExecutionProbability float64            `json:"execution_probability"` // 0.0 to 1.0
    Environment         []string            `json:"environment"`        // When condition is met
    IsTestCode          bool                `json:"is_test_code"`
}

type ConditionType string

const (
    ConditionTypeIf         ConditionType = "if_statement"
    ConditionTypeTry        ConditionType = "try_catch"
    ConditionTypeEnvironment ConditionType = "environment"      // process.env checks
    ConditionTypeFeatureFlag ConditionType = "feature_flag"
    ConditionTypeAsync      ConditionType = "async_load"
)
```

### REQ-UDD-003: Usage Pattern Analysis

**Description:** Analyze usage patterns to distinguish between different types of dependencies and their removal safety.

**Pattern Analysis:**
```go
type UsagePattern struct {
    // Pattern identification
    PatternType         PatternType         `json:"pattern_type"`
    PackageName         string              `json:"package_name"`
    Description         string              `json:"description"`
    
    // Pattern characteristics
    UsageFrequency      int                 `json:"usage_frequency"`
    SourceFiles         []string            `json:"source_files"`
    ImportPaths         []string            `json:"import_paths"`
    
    // Analysis results
    IsEssential         bool                `json:"is_essential"`
    RemovalComplexity   ComplexityLevel     `json:"removal_complexity"`
    AlternativeApproaches []*Alternative    `json:"alternative_approaches"`
    
    // Risk assessment
    BusinessLogicImpact bool                `json:"business_logic_impact"`
    TestCoverage        float64             `json:"test_coverage"`      // 0.0 to 1.0
    ProductionUsage     bool                `json:"production_usage"`
    
    // Recommendations
    RecommendedAction   RecommendedAction   `json:"recommended_action"`
    MigrationSteps      []string            `json:"migration_steps"`
    ConfidenceLevel     float64             `json:"confidence_level"`
}

type PatternType string

const (
    PatternTypeUtility      PatternType = "utility"         // Utility functions
    PatternTypePolyfill     PatternType = "polyfill"        // Browser polyfills
    PatternTypeFramework    PatternType = "framework"       // Framework/library
    PatternTypeTypes        PatternType = "types"           // Type definitions
    PatternTypeTesting      PatternType = "testing"         // Test utilities
    PatternTypeBuild        PatternType = "build"           // Build tools
    PatternTypeDevelopment  PatternType = "development"     // Dev-only usage
    PatternTypeLegacy       PatternType = "legacy"          // Legacy code
    PatternTypeConditional  PatternType = "conditional"     // Conditional usage
)

type ComplexityLevel string

const (
    ComplexityTrivial   ComplexityLevel = "trivial"         // Just remove
    ComplexityLow       ComplexityLevel = "low"             // Simple refactoring
    ComplexityMedium    ComplexityLevel = "medium"          // Some code changes
    ComplexityHigh      ComplexityLevel = "high"            // Significant refactoring
    ComplexityBlocking  ComplexityLevel = "blocking"        // Cannot remove easily
)

type RecommendedAction string

const (
    ActionRemoveImmediately RecommendedAction = "remove_immediately"
    ActionRemoveAfterTesting RecommendedAction = "remove_after_testing"
    ActionReplaceWithAlternative RecommendedAction = "replace_with_alternative"
    ActionRefactorUsage     RecommendedAction = "refactor_usage"
    ActionKeepAsIs          RecommendedAction = "keep_as_is"
    ActionNeedsReview       RecommendedAction = "needs_review"
)

type Alternative struct {
    AlternativeType     AlternativeType     `json:"alternative_type"`
    Name                string              `json:"name"`
    Description         string              `json:"description"`
    ImplementationCost  string              `json:"implementation_cost"`
    Benefits            []string            `json:"benefits"`
    Drawbacks           []string            `json:"drawbacks"`
    MigrationGuide      []string            `json:"migration_guide"`
}

type AlternativeType string

const (
    AltTypeNativeBrowser    AlternativeType = "native_browser"     // Native browser API
    AltTypeNativeNode       AlternativeType = "native_node"        // Native Node.js API
    AltTypeBuiltinLibrary   AlternativeType = "builtin_library"    // Already available library
    AltTypeCustomImplementation AlternativeType = "custom_implementation" // Write custom code
    AltTypeSmallerLibrary   AlternativeType = "smaller_library"    // Smaller alternative
    AltTypeRemoveFeature    AlternativeType = "remove_feature"     // Remove the feature
)
```

## API Interfaces

### REQ-UDD-004: Unused Dependency Detector Interface

**Primary Interface:**
```go
type UnusedDependencyDetector interface {
    // Analyze workspace for unused dependencies
    AnalyzeWorkspace(workspace *WorkspaceConfig, options AnalysisOptions) (*UnusedAnalysis, error)
    
    // Analyze specific packages
    AnalyzePackages(packages []*PackageInfo, options AnalysisOptions) ([]*UnusedDependency, error)
    
    // Detect usage of specific dependency
    DetectUsage(packageName string, sourceFiles []string, options DetectionOptions) (*UsageEvidence, error)
    
    // Validate removal safety
    ValidateRemoval(dependency *UnusedDependency, context *ValidationContext) (*RemovalValidation, error)
    
    // Generate removal plan
    GenerateRemovalPlan(dependencies []*UnusedDependency, preferences *RemovalPreferences) (*RemovalPlan, error)
    
    // Apply removal plan
    ApplyRemovalPlan(plan *RemovalPlan, workspace *WorkspaceConfig) (*RemovalResult, error)
}

type AnalysisOptions struct {
    // Analysis scope
    IncludeDevDependencies  bool                `json:"include_dev_dependencies"`
    IncludePeerDependencies bool                `json:"include_peer_dependencies"`
    IncludeOptionalDeps     bool                `json:"include_optional_deps"`
    
    // Analysis depth
    AnalysisDepth           AnalysisDepth       `json:"analysis_depth"`
    IncludeTypeUsage        bool                `json:"include_type_usage"`
    AnalyzeTestFiles        bool                `json:"analyze_test_files"`
    AnalyzeBuildConfigs     bool                `json:"analyze_build_configs"`
    
    // File filtering
    IncludePatterns         []string            `json:"include_patterns"`    // Files to analyze
    ExcludePatterns         []string            `json:"exclude_patterns"`    // Files to skip
    MaxFileSize             int64               `json:"max_file_size"`       // Skip large files
    
    // Detection sensitivity
    MinConfidence           float64             `json:"min_confidence"`      // Minimum confidence to report
    AllowStringReferences   bool                `json:"allow_string_references"`
    AnalyzeDynamicImports   bool                `json:"analyze_dynamic_imports"`
    
    // Performance options
    ConcurrencyLevel        int                 `json:"concurrency_level"`
    TimeoutPerFile          time.Duration       `json:"timeout_per_file"`
    EnableCaching           bool                `json:"enable_caching"`
    CacheDir               string              `json:"cache_dir"`
    
    // Safety settings
    SafetyThreshold         float64             `json:"safety_threshold"`    // Min safety score
    RequireTestCoverage     bool                `json:"require_test_coverage"`
    FalsePositiveLimit      float64             `json:"false_positive_limit"`
}

type DetectionOptions struct {
    // Detection methods
    EnableStaticAnalysis    bool                `json:"enable_static_analysis"`
    EnableRuntimeAnalysis   bool                `json:"enable_runtime_analysis"`
    EnableStringMatching    bool                `json:"enable_string_matching"`
    
    // AST parsing options
    ParseJavaScript         bool                `json:"parse_javascript"`
    ParseTypeScript         bool                `json:"parse_typescript"`
    ParseJSX                bool                `json:"parse_jsx"`
    ParseTSX                bool                `json:"parse_tsx"`
    ParseJSON               bool                `json:"parse_json"`
    ParseYAML               bool                `json:"parse_yaml"`
    
    // Context awareness
    UnderstandRequireContext bool               `json:"understand_require_context"`
    ResolveModulePaths      bool                `json:"resolve_module_paths"`
    FollowSymlinks          bool                `json:"follow_symlinks"`
    
    // Error handling
    ContinueOnParseError    bool                `json:"continue_on_parse_error"`
    MaxErrors               int                 `json:"max_errors"`
    ReportParsingIssues     bool                `json:"report_parsing_issues"`
}

type ValidationContext struct {
    // Workspace context
    WorkspaceConfig         *WorkspaceConfig    `json:"workspace_config"`
    BuildConfiguration      *BuildConfig        `json:"build_configuration"`
    TestConfiguration       *TestConfig         `json:"test_configuration"`
    
    // Validation requirements
    RequiredTests           []TestType          `json:"required_tests"`
    RiskTolerance          RiskTolerance       `json:"risk_tolerance"`
    ValidationDepth        ValidationDepth     `json:"validation_depth"`
    
    // Environment constraints
    ProductionConstraints  []string            `json:"production_constraints"`
    DevelopmentConstraints []string            `json:"development_constraints"`
    CIConstraints          []string            `json:"ci_constraints"`
}

type ValidationDepth string

const (
    ValidationBasic         ValidationDepth = "basic"           // Basic safety checks
    ValidationStandard      ValidationDepth = "standard"        // Include build verification
    ValidationThorough      ValidationDepth = "thorough"        // Include test execution
    ValidationComprehensive ValidationDepth = "comprehensive"   // Full integration testing
)

type TestType string

const (
    TestTypeUnit            TestType = "unit"
    TestTypeIntegration     TestType = "integration"
    TestTypeE2E             TestType = "e2e"
    TestTypeBuild           TestType = "build"
    TestTypeLint            TestType = "lint"
    TestTypeTypeCheck       TestType = "type_check"
)
```

### REQ-UDD-005: Source Code Analysis Interface

**Code Analysis API:**
```go
type SourceCodeAnalyzer interface {
    // Parse source file and extract imports
    ParseSourceFile(filePath string, options ParseOptions) (*SourceFileAnalysis, error)
    
    // Analyze AST for dependency usage
    AnalyzeAST(ast interface{}, packageName string) (*ASTAnalysis, error)
    
    // Find all references to a package
    FindPackageReferences(sourceCode string, packageName string) ([]*PackageReference, error)
    
    // Analyze import tree shaking potential
    AnalyzeTreeShaking(importStatement *ImportStatement, packageInfo *PackageInfo) (*TreeShakingAnalysis, error)
    
    // Validate analysis accuracy
    ValidateAnalysis(analysis *SourceFileAnalysis, groundTruth *SourceFileAnalysis) (*AnalysisValidation, error)
}

type SourceFileAnalysis struct {
    // File information
    FilePath            string              `json:"file_path"`
    FileSize            int64               `json:"file_size"`
    Language            SourceLanguage      `json:"language"`
    ParsedAt            time.Time           `json:"parsed_at"`
    
    // Import analysis
    StaticImports       []*ImportStatement  `json:"static_imports"`
    DynamicImports      []*DynamicImport    `json:"dynamic_imports"`
    StringReferences    []*StringReference  `json:"string_references"`
    
    // Usage analysis
    PackageUsage        map[string]*PackageUsageInfo `json:"package_usage"` // package -> usage info
    
    // AST information
    ASTNodes            int                 `json:"ast_nodes"`
    ParseErrors         []ParseError        `json:"parse_errors"`
    ParseWarnings       []ParseWarning      `json:"parse_warnings"`
    
    // Metadata
    HasTypeAnnotations  bool                `json:"has_type_annotations"`
    IsTestFile          bool                `json:"is_test_file"`
    IsBuildConfig       bool                `json:"is_build_config"`
    Complexity          int                 `json:"complexity"`         // Cyclomatic complexity
}

type SourceLanguage string

const (
    LangJavaScript      SourceLanguage = "javascript"
    LangTypeScript      SourceLanguage = "typescript"
    LangJSX             SourceLanguage = "jsx"
    LangTSX             SourceLanguage = "tsx"
    LangJSON            SourceLanguage = "json"
    LangYAML            SourceLanguage = "yaml"
    LangMarkdown        SourceLanguage = "markdown"
)

type PackageUsageInfo struct {
    // Usage statistics
    ImportCount         int                 `json:"import_count"`
    UsageCount          int                 `json:"usage_count"`
    UniqueUsages        []string            `json:"unique_usages"`      // Unique ways it's used
    
    // Import analysis
    ImportedSymbols     []string            `json:"imported_symbols"`
    UnusedSymbols       []string            `json:"unused_symbols"`
    TreeShakingPotential float64            `json:"tree_shaking_potential"` // 0.0 to 1.0
    
    // Usage patterns
    UsagePatterns       []UsagePattern      `json:"usage_patterns"`
    IsTypeOnly          bool                `json:"is_type_only"`
    IsConditionalUsage  bool                `json:"is_conditional_usage"`
    
    // Context
    UsageLocations      []*UsageLocation    `json:"usage_locations"`
    FirstUsage          *UsageLocation      `json:"first_usage"`
    LastUsage           *UsageLocation      `json:"last_usage"`
}

type ASTAnalysis struct {
    // AST traversal results
    TotalNodes          int                 `json:"total_nodes"`
    RelevantNodes       int                 `json:"relevant_nodes"`
    TraversalTime       time.Duration       `json:"traversal_time"`
    
    // Found references
    DirectReferences    []*ASTReference     `json:"direct_references"`
    IndirectReferences  []*ASTReference     `json:"indirect_references"`
    TypeReferences      []*ASTReference     `json:"type_references"`
    
    // Analysis metadata
    AnalysisDepth       int                 `json:"analysis_depth"`
    Confidence          float64             `json:"confidence"`
    Completeness        float64             `json:"completeness"`       // How much of AST was analyzed
}

type ASTReference struct {
    NodeType            string              `json:"node_type"`          // AST node type
    ReferenceType       ReferenceType       `json:"reference_type"`
    Location            *SourceLocation     `json:"location"`
    Context             string              `json:"context"`            // Surrounding code
    Metadata            map[string]string   `json:"metadata"`
}

type ReferenceType string

const (
    RefTypeImportDeclaration ReferenceType = "import_declaration"
    RefTypeCallExpression    ReferenceType = "call_expression"
    RefTypeMemberExpression  ReferenceType = "member_expression"
    RefTypeTypeAnnotation    ReferenceType = "type_annotation"
    RefTypeJSXElement        ReferenceType = "jsx_element"
    RefTypeStringLiteral     ReferenceType = "string_literal"
    RefTypeVariableDeclarator ReferenceType = "variable_declarator"
)

type SourceLocation struct {
    Line                int                 `json:"line"`
    Column              int                 `json:"column"`
    StartOffset         int                 `json:"start_offset"`
    EndOffset           int                 `json:"end_offset"`
}

type TreeShakingAnalysis struct {
    // Tree shaking assessment
    IsTreeShakable      bool                `json:"is_tree_shakable"`
    ShakingEffectiveness float64            `json:"shaking_effectiveness"` // 0.0 to 1.0
    
    // Import optimization
    UsedExports         []string            `json:"used_exports"`
    UnusedExports       []string            `json:"unused_exports"`
    OptimizedImport     string              `json:"optimized_import"`   // Suggested import
    
    // Size impact
    CurrentSize         int64               `json:"current_size"`       // Current import size
    OptimizedSize       int64               `json:"optimized_size"`     // After tree shaking
    PotentialSavings    int64               `json:"potential_savings"`
    
    // Implementation
    RequiredChanges     []string            `json:"required_changes"`
    CompatibilityIssues []string            `json:"compatibility_issues"`
    Confidence          float64             `json:"confidence"`
}
```

## Algorithm Specifications

### REQ-UDD-006: Static Analysis Algorithm

**Algorithm:** Multi-pass static analysis with AST traversal and pattern matching

```
ALGORITHM: AnalyzeUnusedDependencies(workspace, options)
INPUT: workspace (WorkspaceConfig) - Workspace configuration
       options (AnalysisOptions) - Analysis configuration
OUTPUT: UnusedAnalysis - Complete unused dependency analysis

1. INITIALIZE analysis context and result structures
2. COLLECT all declared dependencies from package.json files
3. COLLECT all source files matching inclusion patterns

// Phase 1: Parse source files and extract imports
4. FOR each source file in parallel:
   a. PARSE file using appropriate language parser
   b. EXTRACT static imports, dynamic imports, string references
   c. BUILD AST and identify relevant nodes
   d. STORE file analysis results

// Phase 2: Cross-reference declared vs used dependencies
5. CREATE usage map: declared_dependency -> usage_evidence[]
6. FOR each declared dependency:
   a. SEARCH for usage evidence across all source files
   b. ANALYZE usage patterns and contexts
   c. CALCULATE confidence scores
   d. CLASSIFY dependency as used/unused/potentially-unused

// Phase 3: Enhanced analysis for potentially unused
7. FOR each potentially unused dependency:
   a. PERFORM deep AST analysis
   b. CHECK for indirect usage through other dependencies
   c. ANALYZE build tool configurations
   d. CHECK for runtime-only usage patterns
   e. UPDATE classification and confidence

// Phase 4: Safety assessment
8. FOR each unused dependency:
   a. ANALYZE removal risk factors
   b. CHECK for test coverage
   c. IDENTIFY alternative implementations
   d. GENERATE removal recommendations
   e. CALCULATE safety scores

// Phase 5: Generate comprehensive results
9. AGGREGATE results and compute overall metrics
10. RANK unused dependencies by confidence and safety
11. GENERATE removal plan with prioritized recommendations
12. RETURN complete UnusedAnalysis
```

**Complexity:** O(n * m) where n = number of source files, m = number of declared dependencies

### REQ-UDD-007: Import Detection Algorithm

**Algorithm:** Comprehensive import detection with multiple parsing strategies

```
ALGORITHM: DetectAllImports(sourceCode, filePath, packageNames)
INPUT: sourceCode (string) - Source code to analyze
       filePath (string) - File path for context
       packageNames ([]string) - Packages to look for
OUTPUT: []*ImportStatement - All detected imports

1. INITIALIZE parsers for different import types
2. INITIALIZE import collection structures
3. DETERMINE file language from extension and content

// Static import detection
4. PARSE source code with appropriate language parser
5. TRAVERSE AST looking for import/require nodes:
   a. EXTRACT import declarations (ES modules)
   b. EXTRACT require calls (CommonJS)
   c. EXTRACT import() dynamic imports
   d. RECORD location and context information

// String literal analysis
6. SEARCH for string literals that might be package references:
   a. USE regex patterns for common patterns
   b. CHECK against known package names
   c. ANALYZE context to determine likelihood
   d. FILTER false positives

// Dynamic reference detection
7. LOOK for computed imports and dynamic requires:
   a. ANALYZE template literals
   b. CHECK for variable concatenation
   c. IDENTIFY conditional loading patterns
   d. ASSESS runtime determinability

// Context analysis
8. FOR each detected import:
   a. ANALYZE surrounding code context
   b. DETERMINE if import is actually used
   c. CHECK for conditional usage
   d. ASSESS tree-shaking potential
   e. CALCULATE confidence scores

// Validation and filtering
9. VALIDATE detected imports against false positive patterns
10. FILTER out imports that don't match target packages
11. MERGE duplicate detections with different contexts
12. RETURN comprehensive import list with metadata
```

### REQ-UDD-008: Usage Pattern Classification Algorithm

**Algorithm:** Machine learning-based pattern classification with heuristic validation

```
ALGORITHM: ClassifyUsagePatterns(usageEvidence, packageInfo)
INPUT: usageEvidence (UsageEvidence) - Collected usage evidence
       packageInfo (PackageInfo) - Package information
OUTPUT: PatternType - Classified usage pattern

1. INITIALIZE feature extraction and pattern classifiers
2. EXTRACT features from usage evidence:
   a. Import frequency and locations
   b. Usage contexts (test, production, build)
   c. Import types (default, named, dynamic)
   d. File types and patterns
   e. Temporal usage patterns

// Heuristic classification
3. APPLY rule-based classifiers:
   a. IF used only in test files: LIKELY testing pattern
   b. IF used in build configs only: LIKELY build pattern
   c. IF used with type annotations only: LIKELY types pattern
   d. IF used conditionally: LIKELY conditional pattern

// Pattern matching
4. COMPARE against known pattern signatures:
   a. CHECK against utility library patterns
   b. MATCH against polyfill usage patterns
   c. IDENTIFY framework/library patterns
   d. DETECT legacy code patterns

// Machine learning classification
5. IF ML model available:
   a. EXTRACT numerical features
   b. APPLY trained classification model
   c. GET confidence scores for each pattern type
   d. VALIDATE against heuristic results

// Risk assessment
6. FOR identified pattern:
   a. ASSESS removal complexity
   b. IDENTIFY potential alternatives
   c. CALCULATE business logic impact
   d. DETERMINE test coverage requirements

// Final classification
7. COMBINE heuristic and ML results
8. APPLY confidence threshold filtering
9. SELECT most likely pattern with metadata
10. RETURN classified pattern with confidence scores
```

### REQ-UDD-009: Safety Assessment Algorithm

**Algorithm:** Multi-factor safety assessment with risk modeling

```
ALGORITHM: AssessRemovalSafety(unusedDep, context)
INPUT: unusedDep (UnusedDependency) - Dependency to assess
       context (ValidationContext) - Assessment context
OUTPUT: RemovalSafety - Safety classification

1. INITIALIZE risk factor collection and scoring
2. INITIALIZE safety score = 1.0 (maximum safety)

// Evidence strength analysis
3. ANALYZE usage evidence strength:
   a. IF no evidence found: MAINTAIN high safety
   b. IF weak evidence: REDUCE safety moderately
   c. IF strong evidence: CLASSIFY as unsafe
   d. WEIGHT by evidence diversity and recency

// Context risk factors
4. ASSESS contextual risks:
   a. IF used in production code: MAJOR risk factor
   b. IF used in critical paths: HIGH risk factor
   c. IF used in test code only: MINOR risk factor
   d. IF used in deprecated code: REDUCED risk

// Build and runtime dependencies
5. CHECK build tool dependencies:
   a. ANALYZE if package is required for build
   b. CHECK if removing breaks compilation
   c. ASSESS impact on bundling process
   d. EVALUATE tree-shaking implications

// Transitive dependency analysis
6. ANALYZE transitive usage:
   a. CHECK if other dependencies depend on this package
   b. ASSESS indirect usage through other packages
   c. EVALUATE peer dependency relationships
   d. CONSIDER workspace linking implications

// Test coverage analysis
7. IF test analysis enabled:
   a. CHECK test coverage for areas using the package
   b. ASSESS test quality and reliability
   c. EVALUATE integration test coverage
   d. CONSIDER e2e test implications

// Historical analysis
8. IF historical data available:
   a. CHECK previous removal attempts and outcomes
   b. ANALYZE similar package removal success rates
   c. CONSIDER community feedback and experiences
   d. EVALUATE package maintenance status

// Risk scoring
9. CALCULATE overall risk score:
   risk_score = W1 * evidence_risk + 
                W2 * context_risk + 
                W3 * build_risk + 
                W4 * transitive_risk + 
                W5 * test_risk + 
                W6 * historical_risk

10. CONVERT risk score to safety classification:
    a. IF risk_score < 0.2: RemovalSafety = HIGH
    b. IF risk_score < 0.5: RemovalSafety = MEDIUM  
    c. IF risk_score < 0.8: RemovalSafety = LOW
    d. ELSE: RemovalSafety = UNSAFE

11. GENERATE risk mitigation recommendations
12. RETURN safety assessment with detailed rationale
```

## Data Structures

### REQ-UDD-010: Efficient Analysis Storage

**Memory-Optimized Structures:**
```go
// Compressed storage for large-scale analysis
type CompressedUnusedAnalysis struct {
    // String interning
    StringPool          *StringPool         `json:"string_pool"`
    
    // Compressed dependency data
    Dependencies        []CompressedUnusedDep `json:"dependencies"`
    DependencyIndex     map[uint32]int      `json:"dependency_index"`    // name_id -> index
    
    // Usage evidence compression
    EvidencePool        *EvidencePool       `json:"evidence_pool"`
    
    // Compressed source file data
    SourceFilePool      *SourceFilePool     `json:"source_file_pool"`
    
    // Analysis metadata
    Metadata            CompressedMetadata  `json:"metadata"`
    
    // Memory usage tracking
    MemoryUsage         int64               `json:"memory_usage"`
    CompressionRatio    float64             `json:"compression_ratio"`
}

type CompressedUnusedDep struct {
    NameID              uint32              `json:"name_id"`
    VersionID           uint32              `json:"version_id"`
    DeclaredInID        uint32              `json:"declared_in_id"`      // File path ID
    DepType             uint8               `json:"dep_type"`            // DependencyType as byte
    
    // Evidence references
    EvidenceIDs         []uint32            `json:"evidence_ids"`
    
    // Scores (scaled to uint16 for space efficiency)
    ConfidenceScore     uint16              `json:"confidence_score"`    // 0-65535 -> 0.0-1.0
    SafetyScore         uint16              `json:"safety_score"`
    
    // Flags
    Flags               uint32              `json:"flags"`               // Packed boolean flags
    
    // Impact data
    SizeImpactID        uint32              `json:"size_impact_id"`
    RiskAssessmentID    uint32              `json:"risk_assessment_id"`
}

type EvidencePool struct {
    // Evidence storage
    imports             []ImportStatement
    dynamicImports      []DynamicImport
    stringRefs          []StringReference
    buildUsage          []BuildToolUsage
    
    // Lookup tables
    lookupMaps          map[EvidenceType]map[string]uint32
    
    // Statistics
    totalEvidence       int
    uniqueEvidence      int
    duplicateEvidence   int
}

type SourceFilePool struct {
    // File data
    filePaths           []string
    fileContents        []string            // Optional: store parsed content
    fileHashes          []string            // For change detection
    
    // Analysis results
    analyses            []SourceFileAnalysis
    
    // Indexing
    pathToIndex         map[string]int
    hashToIndex         map[string]int
    
    // Statistics
    totalFiles          int
    totalLines          int
    totalSize           int64
}
```

**Analysis Result Cache:**
```go
type AnalysisCache struct {
    // Multi-level caching
    FileAnalysisCache   *LRUCache           `json:"file_analysis_cache"`
    DependencyCache     *LRUCache           `json:"dependency_cache"`
    PatternCache        *LRUCache           `json:"pattern_cache"`
    
    // Persistent storage
    PersistentStore     string              `json:"persistent_store"`    // File path for disk cache
    
    // Cache configuration
    MaxMemoryUsage      int64               `json:"max_memory_usage"`
    MaxCacheAge         time.Duration       `json:"max_cache_age"`
    
    // Invalidation tracking
    FileWatcher         *FileWatcher        `json:"file_watcher"`
    InvalidationRules   []InvalidationRule  `json:"invalidation_rules"`
    
    // Statistics
    HitRate             float64             `json:"hit_rate"`
    MissRate            float64             `json:"miss_rate"`
    EvictionCount       int64               `json:"eviction_count"`
}

type CacheEntry struct {
    Key                 string              `json:"key"`
    Data                interface{}         `json:"data"`
    CreatedAt           time.Time           `json:"created_at"`
    LastAccessed        time.Time           `json:"last_accessed"`
    AccessCount         int                 `json:"access_count"`
    Size                int64               `json:"size"`
    
    // Validation
    Checksum            string              `json:"checksum"`
    Dependencies        []string            `json:"dependencies"`       // Files this entry depends on
    IsValid             bool                `json:"is_valid"`
}

type FileWatcher struct {
    WatchedFiles        map[string]*WatchInfo `json:"watched_files"`
    ChangeCallbacks     []func(string)      `json:"change_callbacks"`
    IsActive            bool                `json:"is_active"`
    
    // Statistics
    FilesWatched        int                 `json:"files_watched"`
    ChangesDetected     int64               `json:"changes_detected"`
    LastChangeTime      time.Time           `json:"last_change_time"`
}

type WatchInfo struct {
    FilePath            string              `json:"file_path"`
    LastModified        time.Time           `json:"last_modified"`
    FileSize            int64               `json:"file_size"`
    Checksum            string              `json:"checksum"`
    CacheEntries        []string            `json:"cache_entries"`     // Cache entries that depend on this file
}
```

### REQ-UDD-011: Pattern Recognition Database

**Pattern Storage and Matching:**
```go
type PatternDatabase struct {
    // Known patterns
    Patterns            []*KnownPattern     `json:"patterns"`
    PatternIndex        map[string]*KnownPattern `json:"pattern_index"`
    
    // Pattern matching
    Matchers            []PatternMatcher    `json:"matchers"`
    
    // Machine learning model
    MLModel             *PatternClassifier  `json:"ml_model"`
    
    // Pattern statistics
    MatchStatistics     map[string]*MatchStats `json:"match_statistics"`
    
    // Version and updates
    DatabaseVersion     int                 `json:"database_version"`
    LastUpdated         time.Time           `json:"last_updated"`
    UpdateSource        string              `json:"update_source"`
}

type KnownPattern struct {
    // Pattern identification
    ID                  string              `json:"id"`
    Name                string              `json:"name"`
    Type                PatternType         `json:"type"`
    Category            string              `json:"category"`
    
    // Pattern characteristics
    ImportPatterns      []string            `json:"import_patterns"`    // Regex patterns
    UsagePatterns       []string            `json:"usage_patterns"`     // Code patterns
    FilePatterns        []string            `json:"file_patterns"`      // File name patterns
    ContextPatterns     []string            `json:"context_patterns"`   // Context clues
    
    // Classification data
    Examples            []*PatternExample   `json:"examples"`
    Counterexamples     []*PatternExample   `json:"counterexamples"`
    Features            map[string]float64  `json:"features"`           // Feature weights
    
    // Metadata
    Confidence          float64             `json:"confidence"`
    Precision           float64             `json:"precision"`
    Recall              float64             `json:"recall"`
    LastValidated       time.Time           `json:"last_validated"`
    Source              string              `json:"source"`             // human, ml, community
}

type PatternExample struct {
    Code                string              `json:"code"`
    Context             string              `json:"context"`
    PackageName         string              `json:"package_name"`
    ExpectedClassification PatternType      `json:"expected_classification"`
    Metadata            map[string]string   `json:"metadata"`
}

type PatternMatcher interface {
    // Match pattern against usage evidence
    Match(evidence *UsageEvidence, context *MatchContext) (*PatternMatch, error)
    
    // Calculate match confidence
    CalculateConfidence(match *PatternMatch) float64
    
    // Update pattern based on feedback
    UpdatePattern(feedback *PatternFeedback) error
    
    // Get pattern statistics
    GetStatistics() *MatcherStats
}

type PatternMatch struct {
    PatternID           string              `json:"pattern_id"`
    Confidence          float64             `json:"confidence"`
    Evidence            []string            `json:"evidence"`           // What matched
    Context             string              `json:"context"`
    FeatureMatches      map[string]bool     `json:"feature_matches"`
    Explanation         string              `json:"explanation"`
}

type MatchContext struct {
    PackageInfo         *PackageInfo        `json:"package_info"`
    SourceFiles         []string            `json:"source_files"`
    BuildContext        *BuildConfig        `json:"build_context"`
    ProjectType         string              `json:"project_type"`
    FrameworkInfo       *FrameworkInfo      `json:"framework_info"`
}

type FrameworkInfo struct {
    Name                string              `json:"name"`               // React, Vue, Angular, etc.
    Version             string              `json:"version"`
    Type                FrameworkType       `json:"type"`
    Dependencies        []string            `json:"dependencies"`
    Patterns            []string            `json:"patterns"`           // Framework-specific patterns
}

type FrameworkType string

const (
    FrameworkReact      FrameworkType = "react"
    FrameworkVue        FrameworkType = "vue"
    FrameworkAngular    FrameworkType = "angular"
    FrameworkSvelte     FrameworkType = "svelte"
    FrameworkNext       FrameworkType = "next"
    FrameworkNuxt       FrameworkType = "nuxt"
    FrameworkExpress    FrameworkType = "express"
    FrameworkNest       FrameworkType = "nest"
)
```

## Error Handling

### REQ-UDD-012: Comprehensive Error Management

**Error Classification:**
```go
type UnusedAnalysisError struct {
    Type                AnalysisErrorType   `json:"type"`
    Phase               AnalysisPhase       `json:"phase"`
    SourceFile          string              `json:"source_file"`
    PackageName         string              `json:"package_name"`
    Message             string              `json:"message"`
    Details             map[string]string   `json:"details"`
    
    // Error context
    Context             *ErrorContext       `json:"context"`
    StackTrace          []string            `json:"stack_trace"`
    
    // Recovery information
    IsRecoverable       bool                `json:"is_recoverable"`
    RecoveryActions     []RecoveryAction    `json:"recovery_actions"`
    PartialResults      interface{}         `json:"partial_results"`
    
    // Error metadata
    Severity            ErrorSeverity       `json:"severity"`
    Timestamp           time.Time           `json:"timestamp"`
    AnalysisId          string              `json:"analysis_id"`
}

type AnalysisErrorType string

const (
    ErrorTypeParseFailure       AnalysisErrorType = "parse_failure"
    ErrorTypeMemoryExhaustion   AnalysisErrorType = "memory_exhaustion"
    ErrorTypeTimeout            AnalysisErrorType = "timeout"
    ErrorTypeFileAccess         AnalysisErrorType = "file_access"
    ErrorTypeInvalidSyntax      AnalysisErrorType = "invalid_syntax"
    ErrorTypeUnsupportedLanguage AnalysisErrorType = "unsupported_language"
    ErrorTypeNetworkFailure     AnalysisErrorType = "network_failure"
    ErrorTypeCacheCorruption    AnalysisErrorType = "cache_corruption"
    ErrorTypeConfigurationError AnalysisErrorType = "configuration_error"
)

type AnalysisPhase string

const (
    PhaseFileDiscovery      AnalysisPhase = "file_discovery"
    PhaseSourceParsing      AnalysisPhase = "source_parsing"
    PhaseImportExtraction   AnalysisPhase = "import_extraction"
    PhaseUsageAnalysis      AnalysisPhase = "usage_analysis"
    PhasePatternMatching    AnalysisPhase = "pattern_matching"
    PhaseSafetyAssessment   AnalysisPhase = "safety_assessment"
    PhaseRecommendationGen  AnalysisPhase = "recommendation_generation"
)
```

**Recovery Strategies:**
```go
type ErrorRecoveryManager struct {
    // Recovery strategies
    Strategies          map[AnalysisErrorType]*RecoveryStrategy `json:"strategies"`
    
    // Error tracking
    ErrorHistory        []*ErrorRecord      `json:"error_history"`
    ErrorPatterns       map[string]int      `json:"error_patterns"`
    
    // Recovery statistics
    SuccessfulRecoveries int                `json:"successful_recoveries"`
    FailedRecoveries    int                 `json:"failed_recoveries"`
    
    // Configuration
    MaxRetryAttempts    int                 `json:"max_retry_attempts"`
    RetryDelay          time.Duration       `json:"retry_delay"`
    EnableAutoRecovery  bool                `json:"enable_auto_recovery"`
}

type RecoveryStrategy struct {
    // Strategy configuration
    StrategyType        RecoveryStrategyType `json:"strategy_type"`
    Priority            int                 `json:"priority"`
    MaxAttempts         int                 `json:"max_attempts"`
    
    // Recovery actions
    Actions             []RecoveryAction    `json:"actions"`
    Fallbacks           []FallbackOption    `json:"fallbacks"`
    
    // Success criteria
    SuccessConditions   []string            `json:"success_conditions"`
    ValidatorFunc       string              `json:"validator_func"`
    
    // Impact assessment
    DataLossRisk        DataLossRisk        `json:"data_loss_risk"`
    QualityImpact       QualityImpact       `json:"quality_impact"`
    PerformanceImpact   PerformanceImpact   `json:"performance_impact"`
}

type RecoveryStrategyType string

const (
    RecoveryRetryWithBackoff    RecoveryStrategyType = "retry_with_backoff"
    RecoverySkipProblematicFile RecoveryStrategyType = "skip_problematic_file"
    RecoveryUseBackupParser     RecoveryStrategyType = "use_backup_parser"
    RecoveryReduceComplexity    RecoveryStrategyType = "reduce_complexity"
    RecoveryUseCachedResults    RecoveryStrategyType = "use_cached_results"
    RecoveryPartialAnalysis     RecoveryStrategyType = "partial_analysis"
    RecoveryManualIntervention  RecoveryStrategyType = "manual_intervention"
)

type FallbackOption struct {
    Name                string              `json:"name"`
    Description         string              `json:"description"`
    Implementation      string              `json:"implementation"`
    QualityReduction    float64             `json:"quality_reduction"`   // 0.0 to 1.0
    ResourceSavings     float64             `json:"resource_savings"`    // 0.0 to 1.0
    AutoApplicable      bool                `json:"auto_applicable"`
}

type DataLossRisk string

const (
    DataLossNone        DataLossRisk = "none"
    DataLossMinimal     DataLossRisk = "minimal"      // <5% of results
    DataLossModerate    DataLossRisk = "moderate"     // 5-25% of results
    DataLossSignificant DataLossRisk = "significant"  // >25% of results
)

type QualityImpact string

const (
    QualityNoImpact     QualityImpact = "no_impact"
    QualityMinorImpact  QualityImpact = "minor_impact"   // Slight reduction in accuracy
    QualityMajorImpact  QualityImpact = "major_impact"   // Significant accuracy loss
    QualitySevereImpact QualityImpact = "severe_impact"  // Substantial quality degradation
)

type PerformanceImpact string

const (
    PerfNoImpact        PerformanceImpact = "no_impact"
    PerfImprovement     PerformanceImpact = "improvement"   // Faster due to reduced scope
    PerfMinorImpact     PerformanceImpact = "minor_impact"  // Slightly slower
    PerfMajorImpact     PerformanceImpact = "major_impact"  // Significantly slower
)
```

## Testing Requirements

### REQ-UDD-013: Comprehensive Test Coverage

**Test Suite Design:**
```go
// Comprehensive test scenarios for unused dependency detection
var UnusedDependencyTestSuite = []TestScenario{
    {
        Name: "basic-unused-detection",
        Description: "Simple unused dependency detection",
        TestData: &TestProject{
            Dependencies: map[string]string{
                "lodash": "^4.17.21",
                "moment": "^2.29.4",
                "uuid": "^9.0.0",
            },
            SourceFiles: map[string]string{
                "src/main.js": "import _ from 'lodash'; console.log(_.isEmpty({}));",
                "src/utils.js": "const moment = require('moment'); export default moment;",
            },
        },
        ExpectedUnused: []string{"uuid"},
        ExpectedUsed: []string{"lodash", "moment"},
    },
    {
        Name: "dynamic-imports",
        Description: "Detection with dynamic imports",
        TestData: &TestProject{
            Dependencies: map[string]string{
                "chart.js": "^3.9.1",
                "d3": "^7.6.1",
            },
            SourceFiles: map[string]string{
                "src/charts.js": `
                    async function loadChart() {
                        const Chart = await import('chart.js');
                        return Chart;
                    }
                `,
            },
        },
        ExpectedUnused: []string{"d3"},
        ExpectedUsed: []string{"chart.js"},
        ExpectedConfidence: map[string]float64{
            "chart.js": 0.85, // Lower confidence due to dynamic import
        },
    },
    {
        Name: "type-only-usage",
        Description: "TypeScript type-only imports",
        TestData: &TestProject{
            Dependencies: map[string]string{
                "express": "^4.18.2",
                "@types/express": "^4.17.17",
                "axios": "^1.4.0",
            },
            SourceFiles: map[string]string{
                "src/server.ts": `
                    import type { Request, Response } from 'express';
                    import axios from 'axios';
                    
                    function handler(req: Request, res: Response) {
                        return axios.get('/api');
                    }
                `,
            },
        },
        ExpectedUnused: []string{"express"},
        ExpectedUsed: []string{"axios", "@types/express"},
        AnalysisOptions: AnalysisOptions{
            IncludeTypeUsage: true,
        },
    },
    {
        Name: "build-tool-usage",
        Description: "Dependencies used only in build configuration",
        TestData: &TestProject{
            Dependencies: map[string]string{
                "webpack": "^5.88.2",
                "babel-loader": "^9.1.2",
                "react": "^18.2.0",
            },
            DevDependencies: map[string]string{
                "webpack-cli": "^5.1.4",
            },
            SourceFiles: map[string]string{
                "webpack.config.js": `
                    module.exports = {
                        module: {
                            rules: [
                                {
                                    test: /\.js$/,
                                    use: 'babel-loader'
                                }
                            ]
                        }
                    };
                `,
                "src/App.js": "import React from 'react'; export default function App() { return <div>Hello</div>; }",
            },
        },
        ExpectedUnused: []string{},
        ExpectedUsed: []string{"webpack", "babel-loader", "react", "webpack-cli"},
        AnalysisOptions: AnalysisOptions{
            AnalyzeBuildConfigs: true,
        },
    },
    {
        Name: "conditional-usage",
        Description: "Dependencies used conditionally",
        TestData: &TestProject{
            Dependencies: map[string]string{
                "debug": "^4.3.4",
                "compression": "^1.7.4",
            },
            SourceFiles: map[string]string{
                "src/server.js": `
                    if (process.env.NODE_ENV === 'development') {
                        const debug = require('debug')('app');
                        debug('Development mode');
                    }
                    
                    if (process.env.ENABLE_COMPRESSION) {
                        const compression = require('compression');
                        app.use(compression());
                    }
                `,
            },
        },
        ExpectedPotentiallyUnused: []string{"debug", "compression"},
        ExpectedSafetyLevel: map[string]RemovalSafety{
            "debug": SafetyMedium,
            "compression": SafetyMedium,
        },
    },
}

type TestScenario struct {
    Name                    string
    Description             string
    TestData               *TestProject
    ExpectedUnused         []string
    ExpectedUsed           []string
    ExpectedPotentiallyUnused []string
    ExpectedConfidence     map[string]float64
    ExpectedSafetyLevel    map[string]RemovalSafety
    AnalysisOptions        AnalysisOptions
    ValidationRules        []ValidationRule
}
```

**Performance Benchmarks:**
```go
func BenchmarkUnusedDetection(b *testing.B) {
    benchmarks := []struct {
        name         string
        sourceFiles  int
        dependencies int
        fileSize     string // small, medium, large
    }{
        {"small-project", 10, 20, "small"},
        {"medium-project", 50, 100, "medium"},
        {"large-project", 200, 300, "large"},
        {"enterprise-project", 1000, 500, "large"},
    }
    
    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            project := generateTestProject(bm.sourceFiles, bm.dependencies, bm.fileSize)
            detector := NewUnusedDependencyDetector()
            
            b.ResetTimer()
            b.ReportAllocs()
            
            for i := 0; i < b.N; i++ {
                analysis, err := detector.AnalyzeWorkspace(project.WorkspaceConfig, DefaultAnalysisOptions)
                if err != nil {
                    b.Fatal(err)
                }
                
                // Validate minimum expected results
                if len(analysis.UnusedDependencies) < 0 {
                    b.Errorf("Expected some analysis results")
                }
            }
        })
    }
}

func BenchmarkParseSourceFile(b *testing.B) {
    testFiles := []struct {
        name     string
        language SourceLanguage
        size     int // lines of code
    }{
        {"javascript-small", LangJavaScript, 100},
        {"typescript-medium", LangTypeScript, 1000},
        {"jsx-large", LangJSX, 5000},
    }
    
    for _, tf := range testFiles {
        b.Run(tf.name, func(b *testing.B) {
            sourceCode := generateSourceCode(tf.language, tf.size)
            analyzer := NewSourceCodeAnalyzer()
            
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                analysis, err := analyzer.ParseSourceFile("test."+string(tf.language), DefaultParseOptions)
                if err != nil {
                    b.Fatal(err)
                }
                
                if analysis.StaticImports == nil {
                    b.Error("Expected parsed imports")
                }
            }
        })
    }
}
```

### REQ-UDD-014: Accuracy Validation

**Accuracy Metrics:**
```go
type AccuracyValidation struct {
    // Test dataset
    GroundTruthDataset  []*GroundTruthEntry `json:"ground_truth_dataset"`
    
    // Validation results
    TruePositives       int                 `json:"true_positives"`
    TrueNegatives       int                 `json:"true_negatives"`
    FalsePositives      int                 `json:"false_positives"`
    FalseNegatives      int                 `json:"false_negatives"`
    
    // Calculated metrics
    Precision           float64             `json:"precision"`
    Recall              float64             `json:"recall"`
    F1Score             float64             `json:"f1_score"`
    Accuracy            float64             `json:"accuracy"`
    
    // Confidence analysis
    ConfidenceDistribution map[string]int   `json:"confidence_distribution"`
    CalibrationError    float64             `json:"calibration_error"`
    
    // Error analysis
    ErrorBreakdown      map[string]int      `json:"error_breakdown"`
    CommonErrors        []*ValidationError  `json:"common_errors"`
}

type GroundTruthEntry struct {
    ProjectPath         string              `json:"project_path"`
    PackageName         string              `json:"package_name"`
    IsActuallyUnused    bool                `json:"is_actually_unused"`
    UsageType           UsageType           `json:"usage_type"`
    RemovalSafety       RemovalSafety       `json:"removal_safety"`
    VerificationMethod  string              `json:"verification_method"` // manual, automated, community
    LastVerified        time.Time           `json:"last_verified"`
    VerifiedBy          string              `json:"verified_by"`
}

type ValidationError struct {
    ErrorType           string              `json:"error_type"`
    PackageName         string              `json:"package_name"`
    Predicted           bool                `json:"predicted"`
    Actual              bool                `json:"actual"`
    Explanation         string              `json:"explanation"`
    Frequency           int                 `json:"frequency"`
}
```

This comprehensive specification provides detailed technical requirements for implementing the Unused Dependency Detector as the final component of MonoGuard's Phase 1 Core Engine, focusing on accurate detection through sophisticated static analysis and high-confidence safety assessment for dependency removal recommendations.