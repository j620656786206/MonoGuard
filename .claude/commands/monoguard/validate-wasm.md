---
description: Validate Go WASM code compliance (Result<T>, camelCase JSON)
argument-hint: [go-file-path]
---

Validate Go WASM code against MonoGuard's strict WASM bridge requirements.

**Usage:**

- `/monoguard:validate-wasm` - Check all Go files in packages/analysis-engine
- `/monoguard:validate-wasm packages/analysis-engine/pkg/analyzer/workspace.go` - Check specific file

**Critical WASM Rules:**

**1. Result<T> Type is MANDATORY:**

```go
// ✅ CORRECT: All WASM functions return Result type
func AnalyzeWorkspace(input string) string {
    if err != nil {
        return toJSON(NewError("PARSE_ERROR", err.Error()))
    }
    return toJSON(NewSuccess(result))
}

// ❌ WRONG: Raw error string
func AnalyzeWorkspace(input string) string {
    if err != nil {
        return err.Error() // VIOLATION
    }
}
```

**2. JSON Must Use camelCase:**

```go
// ✅ CORRECT: camelCase struct tags
type AnalysisResult struct {
    HealthScore int    `json:"healthScore"`
    CreatedAt   string `json:"createdAt"`
}

// ❌ WRONG: snake_case struct tags
type AnalysisResult struct {
    HealthScore int `json:"health_score"` // VIOLATION
}
```

**3. Dates Must Be ISO 8601:**

```go
// ✅ CORRECT: ISO 8601 string
CreatedAt: time.Now().Format(time.RFC3339)

// ❌ WRONG: Unix timestamp
CreatedAt: time.Now().Unix() // VIOLATION
```

**4. Error Codes Must Be UPPER_SNAKE_CASE:**

```go
// ✅ CORRECT
return NewError("PARSE_ERROR", msg)
return NewError("CIRCULAR_DETECTED", msg)

// ❌ WRONG
return NewError("parseError", msg) // VIOLATION
```

**Validation Checks:**

1. **Function Signatures:**
   - All exported functions return `string` (JSON)
   - Input parameters are validated
   - No panics (use error returns)

2. **JSON Structure:**
   - All struct tags use camelCase
   - Date fields are string (ISO 8601)
   - Error responses match Result<T> format

3. **Error Handling:**
   - All errors wrapped in Result type
   - Error codes are UPPER_SNAKE_CASE
   - User-friendly error messages provided

4. **Memory Safety:**
   - No unbounded allocations
   - Large data processed in chunks
   - Memory limit checks for WASM constraints

**Report Format:**

- ✅ WASM-compliant patterns
- ❌ Violations with file:line:column
- 🔧 Auto-fixable issues
- 📖 Reference to correct patterns

Let me validate the Go WASM code.

**Target files:** $ARGUMENTS

I'll analyze the Go code for WASM bridge compliance.
