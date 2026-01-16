// Package result provides the unified Result<T> type for WASM function returns.
// All WASM exported functions MUST return this type to ensure consistent
// TypeScript interoperability.
package result

import "encoding/json"

// Result represents the unified response type for all WASM functions.
// This matches the TypeScript Result<T> type used in the frontend.
//
// JSON structure:
//
//	{
//	  "data": { ... } | null,
//	  "error": { "code": "ERROR_CODE", "message": "..." } | null
//	}
type Result struct {
	Data  interface{} `json:"data"`
	Error *Error      `json:"error"`
}

// Error represents an error response with code and message.
// Error codes MUST use UPPER_SNAKE_CASE convention.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewSuccess creates a successful Result with the given data.
// The data will be serialized to JSON with camelCase field names
// (as defined by struct tags).
func NewSuccess(data interface{}) *Result {
	return &Result{Data: data, Error: nil}
}

// NewError creates an error Result with the given code and message.
// The code should follow UPPER_SNAKE_CASE convention (e.g., "PARSE_ERROR").
func NewError(code, message string) *Result {
	return &Result{Data: nil, Error: &Error{Code: code, Message: message}}
}

// ToJSON serializes the Result to a JSON string.
// This is the format returned to JavaScript from WASM functions.
func (r *Result) ToJSON() string {
	b, err := json.Marshal(r)
	if err != nil {
		// Fallback to error result if marshaling fails
		return `{"data":null,"error":{"code":"MARSHAL_ERROR","message":"Failed to serialize result"}}`
	}
	return string(b)
}
