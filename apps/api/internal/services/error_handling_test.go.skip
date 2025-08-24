package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandler_NewErrorHandler(t *testing.T) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	assert.NotNil(t, eh)
	assert.NotNil(t, eh.logger)
	assert.NotNil(t, eh.errorRegistry)
	assert.NotNil(t, eh.recoverStrategies)
	assert.NotNil(t, eh.errorCounts)
	assert.NotNil(t, eh.errorHistory)
	assert.Equal(t, 1000, eh.maxHistorySize)
	
	// Should have registered default error definitions
	assert.True(t, len(eh.errorRegistry) > 0)
	assert.Contains(t, eh.errorRegistry, ErrorTypeFileSystem)
	assert.Contains(t, eh.errorRegistry, ErrorTypeParsingJSON)
	assert.Contains(t, eh.errorRegistry, ErrorTypeWorkspaceConfig)
	
	// Should have registered default recovery strategies
	assert.True(t, len(eh.recoverStrategies) > 0)
	assert.Contains(t, eh.recoverStrategies, ErrorTypeFileSystem)
	assert.Contains(t, eh.recoverStrategies, ErrorTypeParsingJSON)
}

func TestErrorHandler_ClassifyError(t *testing.T) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	tests := []struct {
		name         string
		err          error
		expectedType ErrorType
	}{
		{
			name:         "parsing error",
			err:          &ParsingError{ErrorType: ErrorTypeParsingJSON},
			expectedType: ErrorTypeParsingJSON,
		},
		{
			name:         "workspace error",
			err:          &WorkspaceError{ErrorType: ErrorTypeWorkspaceConfig},
			expectedType: ErrorTypeWorkspaceConfig,
		},
		{
			name:         "dependency error",
			err:          &DependencyError{ErrorType: ErrorTypeDependencyTree},
			expectedType: ErrorTypeDependencyTree,
		},
		{
			name:         "file system error by message",
			err:          fmt.Errorf("permission denied: cannot access file"),
			expectedType: ErrorTypeFileSystem,
		},
		{
			name:         "json error by message",
			err:          fmt.Errorf("failed to unmarshal JSON"),
			expectedType: ErrorTypeParsingJSON,
		},
		{
			name:         "yaml error by message",
			err:          fmt.Errorf("invalid yaml syntax"),
			expectedType: ErrorTypeParsingYAML,
		},
		{
			name:         "workspace error by message",
			err:          fmt.Errorf("workspace configuration conflict"),
			expectedType: ErrorTypeWorkspaceConfig,
		},
		{
			name:         "version error by message",
			err:          fmt.Errorf("invalid version range: ^1.0.0"),
			expectedType: ErrorTypeVersionRange,
		},
		{
			name:         "timeout error by message",
			err:          fmt.Errorf("operation timeout exceeded"),
			expectedType: ErrorTypeTimeout,
		},
		{
			name:         "memory error by message",
			err:          fmt.Errorf("out of memory: cannot allocate"),
			expectedType: ErrorTypeMemoryLimit,
		},
		{
			name:         "unknown error",
			err:          fmt.Errorf("some unknown error"),
			expectedType: ErrorTypeUnknown,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorType := eh.classifyError(tt.err)
			assert.Equal(t, tt.expectedType, errorType)
		})
	}
}

func TestErrorHandler_GetErrorDefinition(t *testing.T) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	// Test known error type
	def := eh.getErrorDefinition(ErrorTypeFileSystem)
	assert.NotNil(t, def)
	assert.Equal(t, ErrorTypeFileSystem, def.Type)
	assert.Equal(t, SeverityError, def.Severity)
	assert.True(t, def.Recoverable)
	assert.NotEmpty(t, def.UserMessage)
	assert.NotEmpty(t, def.TechnicalMessage)
	assert.True(t, len(def.CommonCauses) > 0)
	assert.True(t, len(def.Solutions) > 0)
	
	// Test unknown error type
	unknownDef := eh.getErrorDefinition(ErrorType("non-existent"))
	assert.NotNil(t, unknownDef)
	assert.Equal(t, ErrorTypeUnknown, unknownDef.Type)
	assert.Equal(t, SeverityError, unknownDef.Severity)
	assert.False(t, unknownDef.Recoverable)
}

func TestErrorHandler_HandleError(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	testErr := fmt.Errorf("permission denied: cannot read file")
	errorContext := map[string]interface{}{
		"file_path": "/test/path/package.json",
		"operation": "read_file",
	}
	
	result, err := eh.HandleError(ctx, testErr, errorContext)
	require.NoError(t, err)
	
	// Should have attempted recovery for recoverable error
	assert.NotNil(t, result)
	assert.Equal(t, "skip_file", result.Action)
	assert.False(t, result.Success)
	assert.False(t, result.RetrySuggested)
	
	// Should have recorded error
	stats := eh.GetErrorStatistics()
	assert.Equal(t, 1, stats["total_errors"])
	assert.Equal(t, int64(1), stats["error_counts"].(map[ErrorType]int64)[ErrorTypeFileSystem])
	
	// Check error history
	assert.Equal(t, 1, len(eh.errorHistory))
	record := eh.errorHistory[0]
	assert.Equal(t, ErrorTypeFileSystem, record.Type)
	assert.Equal(t, testErr.Error(), record.Message)
	assert.True(t, record.RecoveryAttempted)
	assert.False(t, record.RecoverySuccess)
	assert.Equal(t, "skip_file", record.RecoveryAction)
	assert.NotEmpty(t, record.StackTrace)
	assert.True(t, record.Duration > 0)
}

func TestErrorHandler_HandleNilError(t *testing.T) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	result, err := eh.HandleError(ctx, nil, nil)
	
	assert.NoError(t, err)
	assert.Nil(t, result)
	
	// Should not have recorded anything
	stats := eh.GetErrorStatistics()
	assert.Equal(t, 0, stats["total_errors"])
}

func TestErrorHandler_RecoveryStrategies(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	
	tests := []struct {
		name            string
		err             error
		context         map[string]interface{}
		expectedAction  string
		expectedSuccess bool
	}{
		{
			name: "file system error recovery",
			err:  fmt.Errorf("permission denied"),
			context: map[string]interface{}{
				"file_path": "/test/package.json",
			},
			expectedAction:  "skip_file",
			expectedSuccess: false,
		},
		{
			name: "json parsing error recovery",
			err:  fmt.Errorf("invalid json syntax"),
			context: map[string]interface{}{
				"file_path": "/test/malformed.json",
			},
			expectedAction:  "skip_malformed_json",
			expectedSuccess: false,
		},
		{
			name: "workspace config error recovery",
			err:  fmt.Errorf("workspace configuration conflict"),
			context: map[string]interface{}{
				"workspace_paths": []string{"/test/workspace1", "/test/workspace2"},
			},
			expectedAction:  "apply_workspace_priority",
			expectedSuccess: true,
		},
		{
			name: "version range error recovery",
			err:  fmt.Errorf("invalid version range"),
			context: map[string]interface{}{
				"version_range": "invalid-range",
			},
			expectedAction:  "fallback_version_parsing",
			expectedSuccess: true,
		},
		{
			name: "dependency tree error recovery",
			err:  fmt.Errorf("circular dependency detected"),
			context: map[string]interface{}{
				"dependency_path": []string{"pkg-a", "pkg-b", "pkg-a"},
			},
			expectedAction:  "build_partial_tree",
			expectedSuccess: true,
		},
		{
			name: "timeout error recovery",
			err:  fmt.Errorf("operation timeout"),
			context: map[string]interface{}{
				"timeout_duration": "30s",
			},
			expectedAction:  "timeout_suggestions",
			expectedSuccess: false,
		},
		{
			name: "memory limit error recovery",
			err:  fmt.Errorf("out of memory"),
			context: map[string]interface{}{
				"memory_usage": "2GB",
			},
			expectedAction:  "enable_memory_optimization",
			expectedSuccess: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := eh.HandleError(ctx, tt.err, tt.context)
			require.NoError(t, err)
			require.NotNil(t, result)
			
			assert.Equal(t, tt.expectedAction, result.Action)
			assert.Equal(t, tt.expectedSuccess, result.Success)
			assert.NotEmpty(t, result.Message)
		})
	}
}

func TestErrorHandler_GenerateErrorID(t *testing.T) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	id1 := eh.generateErrorID()
	id2 := eh.generateErrorID()
	
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	assert.Contains(t, id1, "err_")
	assert.Contains(t, id2, "err_")
}

func TestErrorHandler_CaptureStackTrace(t *testing.T) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	stackTrace := eh.captureStackTrace()
	
	assert.NotEmpty(t, stackTrace)
	assert.Contains(t, stackTrace, "TestErrorHandler_CaptureStackTrace")
}

func TestErrorHandler_ErrorHistory(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	
	// Add multiple errors
	errors := []error{
		fmt.Errorf("error 1"),
		fmt.Errorf("error 2"),
		fmt.Errorf("error 3"),
	}
	
	for _, err := range errors {
		_, handlerErr := eh.HandleError(ctx, err, nil)
		require.NoError(t, handlerErr)
	}
	
	// Check history
	assert.Equal(t, 3, len(eh.errorHistory))
	
	// Verify chronological order
	for i := 1; i < len(eh.errorHistory); i++ {
		assert.True(t, eh.errorHistory[i].Timestamp.After(eh.errorHistory[i-1].Timestamp) ||
			eh.errorHistory[i].Timestamp.Equal(eh.errorHistory[i-1].Timestamp))
	}
	
	// Test history size limit
	eh.maxHistorySize = 2
	_, handlerErr := eh.HandleError(ctx, fmt.Errorf("error 4"), nil)
	require.NoError(t, handlerErr)
	
	// Should maintain max size
	assert.Equal(t, 2, len(eh.errorHistory))
	// Should have kept the latest errors
	assert.Contains(t, eh.errorHistory[0].Message, "error 3")
	assert.Contains(t, eh.errorHistory[1].Message, "error 4")
}

func TestErrorHandler_GetErrorStatistics(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	
	// Initial state
	stats := eh.GetErrorStatistics()
	assert.Equal(t, 0, stats["total_errors"])
	assert.Equal(t, 0, stats["history_size"])
	assert.Equal(t, 1000, stats["max_history_size"])
	
	// Add some errors
	_, err := eh.HandleError(ctx, fmt.Errorf("permission denied"), nil)
	require.NoError(t, err)
	
	_, err = eh.HandleError(ctx, fmt.Errorf("workspace conflict"), nil)
	require.NoError(t, err)
	
	_, err = eh.HandleError(ctx, fmt.Errorf("another permission denied"), nil)
	require.NoError(t, err)
	
	// Check updated stats
	stats = eh.GetErrorStatistics()
	assert.Equal(t, 3, stats["total_errors"])
	assert.Equal(t, 3, stats["history_size"])
	
	errorCounts := stats["error_counts"].(map[ErrorType]int64)
	assert.Equal(t, int64(2), errorCounts[ErrorTypeFileSystem])
	assert.Equal(t, int64(1), errorCounts[ErrorTypeWorkspaceConfig])
	
	// Check recovery rate
	recoveryRate := stats["recovery_rate"].(float64)
	assert.True(t, recoveryRate >= 0 && recoveryRate <= 1)
}

func TestErrorHandler_ClearErrorHistory(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	
	// Add some errors
	_, err := eh.HandleError(ctx, fmt.Errorf("test error 1"), nil)
	require.NoError(t, err)
	_, err = eh.HandleError(ctx, fmt.Errorf("test error 2"), nil)
	require.NoError(t, err)
	
	// Verify errors exist
	stats := eh.GetErrorStatistics()
	assert.Equal(t, 2, stats["total_errors"])
	
	// Clear history
	eh.ClearErrorHistory()
	
	// Verify cleared
	stats = eh.GetErrorStatistics()
	assert.Equal(t, 0, stats["total_errors"])
	assert.Equal(t, 0, stats["history_size"])
	
	errorCounts := stats["error_counts"].(map[ErrorType]int64)
	assert.Equal(t, 0, len(errorCounts))
}

func TestParsingError_Error(t *testing.T) {
	err := &ParsingError{
		FilePath:   "/test/package.json",
		LineNumber: 5,
		Column:     10,
		ErrorType:  ErrorTypeParsingJSON,
		Message:    "unexpected token",
	}
	
	errorStr := err.Error()
	assert.Contains(t, errorStr, "parsing_json")
	assert.Contains(t, errorStr, "line 5")
	assert.Contains(t, errorStr, "column 10")
	assert.Contains(t, errorStr, "/test/package.json")
	assert.Contains(t, errorStr, "unexpected token")
}

func TestWorkspaceError_Error(t *testing.T) {
	err := &WorkspaceError{
		WorkspacePath: "/test/workspace",
		ConfigFile:    "pnpm-workspace.yaml",
		ErrorType:     ErrorTypeWorkspaceConfig,
		Message:       "conflicting patterns",
	}
	
	errorStr := err.Error()
	assert.Contains(t, errorStr, "/test/workspace")
	assert.Contains(t, errorStr, "pnpm-workspace.yaml")
	assert.Contains(t, errorStr, "conflicting patterns")
}

func TestDependencyError_Error(t *testing.T) {
	err := &DependencyError{
		PackageName:    "react",
		Version:        "18.0.0",
		DependencyPath: []string{"app", "package-a", "react"},
		ErrorType:      ErrorTypeDependencyTree,
		Message:        "circular dependency",
	}
	
	errorStr := err.Error()
	assert.Contains(t, errorStr, "react@18.0.0")
	assert.Contains(t, errorStr, "app -> package-a -> react")
	assert.Contains(t, errorStr, "circular dependency")
}

// Test concurrent access to error handler

func TestErrorHandler_ConcurrentAccess(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	const numGoroutines = 10
	const errorsPerGoroutine = 20
	
	done := make(chan bool)
	
	// Concurrent error handling
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			for j := 0; j < errorsPerGoroutine; j++ {
				testErr := fmt.Errorf("concurrent error %d-%d", id, j)
				_, err := eh.HandleError(ctx, testErr, map[string]interface{}{
					"goroutine_id": id,
					"error_num":    j,
				})
				require.NoError(t, err)
			}
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// Verify all errors were recorded
	stats := eh.GetErrorStatistics()
	expectedTotal := numGoroutines * errorsPerGoroutine
	assert.Equal(t, expectedTotal, stats["total_errors"])
	assert.Equal(t, expectedTotal, stats["history_size"])
}

// Test error severity logging

func TestErrorHandler_LogError(t *testing.T) {
	// Create a logger with a custom hook to capture log entries
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	
	// We can't easily test the actual logging without a custom hook,
	// but we can test that the method doesn't panic
	eh := NewErrorHandler(logger)
	
	record := &ErrorRecord{
		ID:       "test-id",
		Type:     ErrorTypeFileSystem,
		Severity: SeverityError,
		Message:  "test error",
		Context: map[string]interface{}{
			"test_key": "test_value",
		},
		Duration: time.Millisecond * 100,
	}
	
	definition := eh.getErrorDefinition(ErrorTypeFileSystem)
	
	// Should not panic
	assert.NotPanics(t, func() {
		eh.logError(record, definition)
	})
}

// Test error definitions completeness

func TestErrorHandler_ErrorDefinitionsCompleteness(t *testing.T) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	expectedErrorTypes := []ErrorType{
		ErrorTypeFileSystem,
		ErrorTypeParsingJSON,
		ErrorTypeParsingYAML,
		ErrorTypeWorkspaceConfig,
		ErrorTypeVersionRange,
		ErrorTypeDependencyTree,
		ErrorTypeTimeout,
		ErrorTypeMemoryLimit,
	}
	
	for _, errorType := range expectedErrorTypes {
		t.Run(string(errorType), func(t *testing.T) {
			def := eh.getErrorDefinition(errorType)
			assert.Equal(t, errorType, def.Type)
			assert.NotEmpty(t, def.Description)
			assert.NotEmpty(t, def.UserMessage)
			assert.NotEmpty(t, def.TechnicalMessage)
			assert.True(t, len(def.CommonCauses) > 0)
			assert.True(t, len(def.Solutions) > 0)
			assert.True(t, def.Severity != "")
			
			// Check if recovery strategy exists for recoverable errors
			if def.Recoverable {
				_, hasStrategy := eh.recoverStrategies[errorType]
				assert.True(t, hasStrategy, "Recoverable error type should have recovery strategy")
			}
		})
	}
}

// Benchmark tests

func BenchmarkErrorHandler_HandleError(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	testErr := fmt.Errorf("benchmark error")
	errorContext := map[string]interface{}{
		"benchmark": true,
		"iteration": 0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		errorContext["iteration"] = i
		_, err := eh.HandleError(ctx, testErr, errorContext)
		require.NoError(b, err)
	}
}

func BenchmarkErrorHandler_ClassifyError(b *testing.B) {
	logger := logrus.New()
	eh := NewErrorHandler(logger)
	
	errors := []error{
		fmt.Errorf("permission denied"),
		fmt.Errorf("invalid json syntax"),
		fmt.Errorf("workspace conflict"),
		fmt.Errorf("version error"),
		fmt.Errorf("timeout exceeded"),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := errors[i%len(errors)]
		eh.classifyError(err)
	}
}

// Test recovery strategy edge cases

func TestErrorHandler_RecoveryStrategyEdgeCases(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	eh := NewErrorHandler(logger)
	
	ctx := context.Background()
	
	// Test recovery with empty context
	result, err := eh.recoverFileSystemError(ctx, fmt.Errorf("test error"), nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	
	// Test recovery with malformed context
	badContext := map[string]interface{}{
		"file_path": 123, // Wrong type
	}
	result, err = eh.recoverFileSystemError(ctx, fmt.Errorf("test error"), badContext)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	
	// Test JSON recovery with empty context
	result, err = eh.recoverJSONParsingError(ctx, fmt.Errorf("json error"), nil)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
}