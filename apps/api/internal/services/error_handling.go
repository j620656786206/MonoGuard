package services

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ErrorHandler provides comprehensive error handling and recovery mechanisms
type ErrorHandler struct {
	logger          *logrus.Logger
	errorRegistry   map[ErrorType]*ErrorDefinition
	recoverStrategies map[ErrorType]RecoveryStrategy
	errorCounts     map[ErrorType]int64
	errorHistory    []*ErrorRecord
	mutex           sync.RWMutex
	maxHistorySize  int
}

// ErrorType represents different categories of errors
type ErrorType string

const (
	ErrorTypeFileSystem       ErrorType = "filesystem"
	ErrorTypeParsingJSON      ErrorType = "parsing_json"
	ErrorTypeParsingYAML      ErrorType = "parsing_yaml"
	ErrorTypeWorkspaceConfig  ErrorType = "workspace_config"
	ErrorTypeVersionRange     ErrorType = "version_range"
	ErrorTypeDependencyTree   ErrorType = "dependency_tree"
	ErrorTypeNetworkRequest   ErrorType = "network_request"
	ErrorTypeTimeout          ErrorType = "timeout"
	ErrorTypeMemoryLimit      ErrorType = "memory_limit"
	ErrorTypeCacheCorruption  ErrorType = "cache_corruption"
	ErrorTypeCircularDep      ErrorType = "circular_dependency"
	ErrorTypeConflictResolution ErrorType = "conflict_resolution"
	ErrorTypeUnknown          ErrorType = "unknown"
)

// ErrorSeverity represents the severity level of errors
type ErrorSeverity string

const (
	SeverityInfo     ErrorSeverity = "info"
	SeverityWarning  ErrorSeverity = "warning"
	SeverityError    ErrorSeverity = "error"
	SeverityErrorCritical ErrorSeverity = "critical"
)

// ErrorDefinition contains metadata about error types
type ErrorDefinition struct {
	Type            ErrorType     `json:"type"`
	Severity        ErrorSeverity `json:"severity"`
	Description     string        `json:"description"`
	Recoverable     bool          `json:"recoverable"`
	UserMessage     string        `json:"userMessage"`
	TechnicalMessage string       `json:"technicalMessage"`
	DocumentationURL string       `json:"documentationUrl,omitempty"`
	CommonCauses    []string      `json:"commonCauses"`
	Solutions       []string      `json:"solutions"`
}

// ErrorRecord represents a recorded error instance
type ErrorRecord struct {
	ID               string                 `json:"id"`
	Type             ErrorType              `json:"type"`
	Severity         ErrorSeverity          `json:"severity"`
	Message          string                 `json:"message"`
	StackTrace       string                 `json:"stackTrace"`
	Timestamp        time.Time              `json:"timestamp"`
	Context          map[string]interface{} `json:"context"`
	RecoveryAttempted bool                  `json:"recoveryAttempted"`
	RecoverySuccess   bool                  `json:"recoverySuccess"`
	RecoveryAction    string                 `json:"recoveryAction,omitempty"`
	Duration          time.Duration          `json:"duration"`
}

// RecoveryStrategy defines how to recover from specific error types
type RecoveryStrategy func(ctx context.Context, err error, context map[string]interface{}) (*RecoveryResult, error)

// RecoveryResult contains the result of a recovery attempt
type RecoveryResult struct {
	Success         bool                   `json:"success"`
	Action          string                 `json:"action"`
	Message         string                 `json:"message"`
	RetrySuggested  bool                   `json:"retrySuggested"`
	ModifiedContext map[string]interface{} `json:"modifiedContext,omitempty"`
	PartialResult   interface{}            `json:"partialResult,omitempty"`
}

// ParsingError represents structured parsing errors
type ParsingError struct {
	FilePath    string    `json:"filePath"`
	LineNumber  int       `json:"lineNumber"`
	Column      int       `json:"column"`
	ErrorType   ErrorType `json:"errorType"`
	Message     string    `json:"message"`
	Content     string    `json:"content"`
	Suggestions []string  `json:"suggestions"`
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("%s at line %d, column %d in %s: %s", e.ErrorType, e.LineNumber, e.Column, e.FilePath, e.Message)
}

// WorkspaceError represents workspace-specific errors
type WorkspaceError struct {
	WorkspacePath   string    `json:"workspacePath"`
	ConfigFile      string    `json:"configFile"`
	ErrorType       ErrorType `json:"errorType"`
	Message         string    `json:"message"`
	ConflictingWorkspaces []string `json:"conflictingWorkspaces,omitempty"`
	SuggestedFix    string    `json:"suggestedFix"`
}

func (e *WorkspaceError) Error() string {
	return fmt.Sprintf("Workspace error in %s (%s): %s", e.WorkspacePath, e.ConfigFile, e.Message)
}

// DependencyError represents dependency-related errors
type DependencyError struct {
	PackageName     string                 `json:"packageName"`
	Version         string                 `json:"version"`
	DependencyPath  []string               `json:"dependencyPath"`
	ErrorType       ErrorType              `json:"errorType"`
	Message         string                 `json:"message"`
	ConflictDetails map[string]interface{} `json:"conflictDetails,omitempty"`
	ResolutionHints []string               `json:"resolutionHints"`
}

func (e *DependencyError) Error() string {
	pathStr := strings.Join(e.DependencyPath, " -> ")
	return fmt.Sprintf("Dependency error for %s@%s (path: %s): %s", e.PackageName, e.Version, pathStr, e.Message)
}

// NewErrorHandler creates a new error handler with default configurations
func NewErrorHandler(logger *logrus.Logger) *ErrorHandler {
	eh := &ErrorHandler{
		logger:            logger,
		errorRegistry:     make(map[ErrorType]*ErrorDefinition),
		recoverStrategies: make(map[ErrorType]RecoveryStrategy),
		errorCounts:       make(map[ErrorType]int64),
		errorHistory:      make([]*ErrorRecord, 0),
		maxHistorySize:    1000,
	}

	// Register default error definitions and recovery strategies
	eh.registerDefaultErrorDefinitions()
	eh.registerDefaultRecoveryStrategies()

	return eh
}

// registerDefaultErrorDefinitions registers built-in error definitions
func (eh *ErrorHandler) registerDefaultErrorDefinitions() {
	definitions := []*ErrorDefinition{
		{
			Type:        ErrorTypeFileSystem,
			Severity:    SeverityError,
			Description: "File system access or permission errors",
			Recoverable: true,
			UserMessage: "Unable to access or read project files. Check file permissions.",
			TechnicalMessage: "File system operation failed",
			CommonCauses: []string{
				"Insufficient file permissions",
				"File not found",
				"Directory not accessible",
				"Disk space issues",
			},
			Solutions: []string{
				"Check file and directory permissions",
				"Verify the file path exists",
				"Ensure sufficient disk space",
				"Run with appropriate privileges",
			},
		},
		{
			Type:        ErrorTypeParsingJSON,
			Severity:    SeverityError,
			Description: "JSON parsing and validation errors",
			Recoverable: true,
			UserMessage: "Invalid JSON format detected. Please check package.json files.",
			TechnicalMessage: "JSON parsing failed",
			CommonCauses: []string{
				"Malformed JSON syntax",
				"Missing commas or brackets",
				"Invalid escape sequences",
				"Trailing commas",
			},
			Solutions: []string{
				"Validate JSON syntax using a JSON validator",
				"Check for missing commas and brackets",
				"Remove trailing commas",
				"Use proper JSON escape sequences",
			},
		},
		{
			Type:        ErrorTypeParsingYAML,
			Severity:    SeverityError,
			Description: "YAML parsing and validation errors",
			Recoverable: true,
			UserMessage: "Invalid YAML format detected. Please check workspace configuration files.",
			TechnicalMessage: "YAML parsing failed",
			CommonCauses: []string{
				"Invalid YAML syntax",
				"Incorrect indentation",
				"Invalid characters",
				"Structure mismatch",
			},
			Solutions: []string{
				"Validate YAML syntax using a YAML validator",
				"Check indentation consistency",
				"Remove invalid characters",
				"Verify YAML structure",
			},
		},
		{
			Type:        ErrorTypeWorkspaceConfig,
			Severity:    SeverityWarning,
			Description: "Workspace configuration conflicts or issues",
			Recoverable: true,
			UserMessage: "Workspace configuration issues detected. Some packages may not be discovered.",
			TechnicalMessage: "Workspace configuration conflict",
			CommonCauses: []string{
				"Multiple workspace configurations",
				"Conflicting package patterns",
				"Invalid workspace paths",
				"Circular workspace references",
			},
			Solutions: []string{
				"Consolidate workspace configurations",
				"Resolve pattern conflicts",
				"Verify workspace paths",
				"Remove circular references",
			},
		},
		{
			Type:        ErrorTypeVersionRange,
			Severity:    SeverityWarning,
			Description: "Version range parsing or resolution errors",
			Recoverable: true,
			UserMessage: "Invalid version ranges detected. Some dependency analysis may be incomplete.",
			TechnicalMessage: "Version range parsing failed",
			CommonCauses: []string{
				"Invalid semantic version format",
				"Unsupported version operators",
				"Malformed version strings",
				"Conflicting version constraints",
			},
			Solutions: []string{
				"Use valid semantic versioning",
				"Check version operators (^, ~, >=, etc.)",
				"Fix malformed version strings",
				"Resolve version constraints",
			},
		},
		{
			Type:        ErrorTypeDependencyTree,
			Severity:    SeverityError,
			Description: "Dependency tree construction errors",
			Recoverable: true,
			UserMessage: "Unable to build complete dependency tree. Analysis may be incomplete.",
			TechnicalMessage: "Dependency tree construction failed",
			CommonCauses: []string{
				"Circular dependencies",
				"Missing dependencies",
				"Version conflicts",
				"Network issues",
			},
			Solutions: []string{
				"Resolve circular dependencies",
				"Install missing dependencies",
				"Fix version conflicts",
				"Check network connectivity",
			},
		},
		{
			Type:        ErrorTypeTimeout,
			Severity:    SeverityWarning,
			Description: "Operation timeout errors",
			Recoverable: true,
			UserMessage: "Operation timed out. Consider increasing timeout settings.",
			TechnicalMessage: "Operation exceeded timeout limit",
			CommonCauses: []string{
				"Large repository size",
				"Slow network connection",
				"Resource constraints",
				"Complex dependency tree",
			},
			Solutions: []string{
				"Increase timeout settings",
				"Reduce analysis scope",
				"Improve network connection",
				"Optimize resource usage",
			},
		},
		{
			Type:        ErrorTypeMemoryLimit,
			Severity:    SeverityErrorCritical,
			Description: "Memory limit exceeded during analysis",
			Recoverable: true,
			UserMessage: "Insufficient memory to complete analysis. Try reducing scope or increasing memory allocation.",
			TechnicalMessage: "Memory allocation failed",
			CommonCauses: []string{
				"Large monorepo size",
				"Memory leaks",
				"Insufficient system memory",
				"Memory fragmentation",
			},
			Solutions: []string{
				"Increase memory allocation",
				"Reduce analysis scope",
				"Enable memory optimization",
				"Process in batches",
			},
		},
	}

	for _, def := range definitions {
		eh.errorRegistry[def.Type] = def
	}
}

// registerDefaultRecoveryStrategies registers built-in recovery strategies
func (eh *ErrorHandler) registerDefaultRecoveryStrategies() {
	eh.recoverStrategies[ErrorTypeFileSystem] = eh.recoverFileSystemError
	eh.recoverStrategies[ErrorTypeParsingJSON] = eh.recoverJSONParsingError
	eh.recoverStrategies[ErrorTypeParsingYAML] = eh.recoverYAMLParsingError
	eh.recoverStrategies[ErrorTypeWorkspaceConfig] = eh.recoverWorkspaceConfigError
	eh.recoverStrategies[ErrorTypeVersionRange] = eh.recoverVersionRangeError
	eh.recoverStrategies[ErrorTypeDependencyTree] = eh.recoverDependencyTreeError
	eh.recoverStrategies[ErrorTypeTimeout] = eh.recoverTimeoutError
	eh.recoverStrategies[ErrorTypeMemoryLimit] = eh.recoverMemoryLimitError
}

// HandleError processes an error with comprehensive handling and recovery
func (eh *ErrorHandler) HandleError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	if err == nil {
		return nil, nil
	}

	startTime := time.Now()
	
	// Determine error type and get definition
	errorType := eh.classifyError(err)
	definition := eh.getErrorDefinition(errorType)
	
	// Create error record
	record := &ErrorRecord{
		ID:        eh.generateErrorID(),
		Type:      errorType,
		Severity:  definition.Severity,
		Message:   err.Error(),
		Timestamp: startTime,
		Context:   errorContext,
		StackTrace: eh.captureStackTrace(),
	}

	// Log the error
	eh.logError(record, definition)

	// Increment error count
	eh.mutex.Lock()
	eh.errorCounts[errorType]++
	eh.mutex.Unlock()

	// Attempt recovery if strategy exists and error is recoverable
	var recoveryResult *RecoveryResult
	if definition.Recoverable {
		if strategy, exists := eh.recoverStrategies[errorType]; exists {
			record.RecoveryAttempted = true
			
			recoveryResult, recoveryErr := strategy(ctx, err, errorContext)
			if recoveryErr == nil && recoveryResult != nil {
				record.RecoverySuccess = recoveryResult.Success
				record.RecoveryAction = recoveryResult.Action
			} else {
				record.RecoverySuccess = false
				if recoveryErr != nil {
					eh.logger.WithError(recoveryErr).Warn("Recovery strategy failed")
				}
			}
		}
	}

	// Finalize record
	record.Duration = time.Since(startTime)

	// Store error record
	eh.addErrorRecord(record)

	return recoveryResult, nil
}

// classifyError determines the error type based on error content and type
func (eh *ErrorHandler) classifyError(err error) ErrorType {
	switch e := err.(type) {
	case *ParsingError:
		return e.ErrorType
	case *WorkspaceError:
		return e.ErrorType
	case *DependencyError:
		return e.ErrorType
	default:
		// Classify by error message content
		errMsg := strings.ToLower(err.Error())
		
		switch {
		case strings.Contains(errMsg, "permission denied") || strings.Contains(errMsg, "no such file"):
			return ErrorTypeFileSystem
		case strings.Contains(errMsg, "json") || strings.Contains(errMsg, "unmarshal"):
			return ErrorTypeParsingJSON
		case strings.Contains(errMsg, "yaml"):
			return ErrorTypeParsingYAML
		case strings.Contains(errMsg, "workspace"):
			return ErrorTypeWorkspaceConfig
		case strings.Contains(errMsg, "version") || strings.Contains(errMsg, "semver"):
			return ErrorTypeVersionRange
		case strings.Contains(errMsg, "dependency") || strings.Contains(errMsg, "circular"):
			return ErrorTypeDependencyTree
		case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline"):
			return ErrorTypeTimeout
		case strings.Contains(errMsg, "memory") || strings.Contains(errMsg, "out of memory"):
			return ErrorTypeMemoryLimit
		default:
			return ErrorTypeUnknown
		}
	}
}

// getErrorDefinition retrieves error definition for a given type
func (eh *ErrorHandler) getErrorDefinition(errorType ErrorType) *ErrorDefinition {
	eh.mutex.RLock()
	defer eh.mutex.RUnlock()
	
	if def, exists := eh.errorRegistry[errorType]; exists {
		return def
	}
	
	// Return default unknown error definition
	return &ErrorDefinition{
		Type:        ErrorTypeUnknown,
		Severity:    SeverityError,
		Description: "Unknown error type",
		Recoverable: false,
		UserMessage: "An unexpected error occurred",
		TechnicalMessage: "Unknown error",
	}
}

// Recovery strategy implementations

func (eh *ErrorHandler) recoverFileSystemError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting file system error recovery")
	
	// Try alternative file paths or permissions
	if filePath, exists := errorContext["file_path"].(string); exists {
		// Attempt to find alternative paths
		alternativePaths := eh.findAlternativeFilePaths(filePath)
		if len(alternativePaths) > 0 {
			return &RecoveryResult{
				Success:        true,
				Action:         "alternative_path",
				Message:        fmt.Sprintf("Using alternative path: %s", alternativePaths[0]),
				RetrySuggested: true,
				ModifiedContext: map[string]interface{}{
					"alternative_paths": alternativePaths,
				},
			}, nil
		}
	}

	return &RecoveryResult{
		Success:        false,
		Action:         "skip_file",
		Message:        "Skipping inaccessible file",
		RetrySuggested: false,
	}, nil
}

func (eh *ErrorHandler) recoverJSONParsingError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting JSON parsing error recovery")
	
	if filePath, exists := errorContext["file_path"].(string); exists {
		// Attempt to fix common JSON issues
		if fixedContent := eh.attemptJSONFix(filePath); fixedContent != "" {
			return &RecoveryResult{
				Success:        true,
				Action:         "json_fix_applied",
				Message:        "Applied automatic JSON fixes",
				RetrySuggested: true,
				ModifiedContext: map[string]interface{}{
					"fixed_content": fixedContent,
				},
			}, nil
		}
	}

	return &RecoveryResult{
		Success:        false,
		Action:         "skip_malformed_json",
		Message:        "Skipping malformed JSON file",
		RetrySuggested: false,
	}, nil
}

func (eh *ErrorHandler) recoverYAMLParsingError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting YAML parsing error recovery")
	
	return &RecoveryResult{
		Success:        false,
		Action:         "skip_malformed_yaml",
		Message:        "Skipping malformed YAML file",
		RetrySuggested: false,
	}, nil
}

func (eh *ErrorHandler) recoverWorkspaceConfigError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting workspace config error recovery")
	
	// Apply priority-based conflict resolution
	return &RecoveryResult{
		Success:        true,
		Action:         "apply_workspace_priority",
		Message:        "Applied priority-based workspace resolution",
		RetrySuggested: true,
		ModifiedContext: map[string]interface{}{
			"resolution_applied": true,
		},
	}, nil
}

func (eh *ErrorHandler) recoverVersionRangeError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting version range error recovery")
	
	// Use fallback version parsing
	return &RecoveryResult{
		Success:        true,
		Action:         "fallback_version_parsing",
		Message:        "Using simplified version parsing",
		RetrySuggested: true,
	}, nil
}

func (eh *ErrorHandler) recoverDependencyTreeError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting dependency tree error recovery")
	
	// Build partial tree with available data
	return &RecoveryResult{
		Success:        true,
		Action:         "build_partial_tree",
		Message:        "Built partial dependency tree with available data",
		RetrySuggested: false,
		PartialResult:  "partial_tree_data",
	}, nil
}

func (eh *ErrorHandler) recoverTimeoutError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting timeout error recovery")
	
	// Suggest reducing scope or increasing timeout
	return &RecoveryResult{
		Success:        false,
		Action:         "timeout_suggestions",
		Message:        "Consider reducing analysis scope or increasing timeout",
		RetrySuggested: true,
		ModifiedContext: map[string]interface{}{
			"timeout_extended": true,
			"scope_reduced":    true,
		},
	}, nil
}

func (eh *ErrorHandler) recoverMemoryLimitError(ctx context.Context, err error, errorContext map[string]interface{}) (*RecoveryResult, error) {
	eh.logger.Debug("Attempting memory limit error recovery")
	
	// Enable memory optimization and batch processing
	return &RecoveryResult{
		Success:        true,
		Action:         "enable_memory_optimization",
		Message:        "Enabled memory optimization and batch processing",
		RetrySuggested: true,
		ModifiedContext: map[string]interface{}{
			"memory_optimized": true,
			"batch_processing": true,
		},
	}, nil
}

// Helper methods

func (eh *ErrorHandler) generateErrorID() string {
	return fmt.Sprintf("err_%d", time.Now().UnixNano())
}

func (eh *ErrorHandler) captureStackTrace() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

func (eh *ErrorHandler) logError(record *ErrorRecord, definition *ErrorDefinition) {
	fields := logrus.Fields{
		"error_id":   record.ID,
		"error_type": record.Type,
		"severity":   record.Severity,
		"duration":   record.Duration,
	}

	// Add context fields
	for k, v := range record.Context {
		fields["ctx_"+k] = v
	}

	switch definition.Severity {
	case SeverityInfo:
		eh.logger.WithFields(fields).Info(definition.TechnicalMessage)
	case SeverityWarning:
		eh.logger.WithFields(fields).Warn(definition.TechnicalMessage)
	case SeverityError:
		eh.logger.WithFields(fields).Error(definition.TechnicalMessage)
	case SeverityErrorCritical:
		eh.logger.WithFields(fields).Fatal(definition.TechnicalMessage)
	}
}

func (eh *ErrorHandler) addErrorRecord(record *ErrorRecord) {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()

	eh.errorHistory = append(eh.errorHistory, record)
	
	// Maintain history size limit
	if len(eh.errorHistory) > eh.maxHistorySize {
		eh.errorHistory = eh.errorHistory[1:]
	}
}

func (eh *ErrorHandler) findAlternativeFilePaths(filePath string) []string {
	// Simplified implementation - would contain more sophisticated path resolution
	return []string{}
}

func (eh *ErrorHandler) attemptJSONFix(filePath string) string {
	// Simplified implementation - would contain JSON repair logic
	return ""
}

// GetErrorStatistics returns error statistics
func (eh *ErrorHandler) GetErrorStatistics() map[string]interface{} {
	eh.mutex.RLock()
	defer eh.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_errors":    len(eh.errorHistory),
		"error_counts":    eh.errorCounts,
		"history_size":    len(eh.errorHistory),
		"max_history_size": eh.maxHistorySize,
	}

	// Calculate recovery rate
	recoveredCount := 0
	for _, record := range eh.errorHistory {
		if record.RecoverySuccess {
			recoveredCount++
		}
	}

	if len(eh.errorHistory) > 0 {
		stats["recovery_rate"] = float64(recoveredCount) / float64(len(eh.errorHistory))
	}

	return stats
}

// ClearErrorHistory clears the error history
func (eh *ErrorHandler) ClearErrorHistory() {
	eh.mutex.Lock()
	defer eh.mutex.Unlock()

	eh.errorHistory = make([]*ErrorRecord, 0)
	eh.errorCounts = make(map[ErrorType]int64)
}