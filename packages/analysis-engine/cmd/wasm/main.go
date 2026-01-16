//go:build js && wasm

// Package main is the WASM entry point for MonoGuard analysis engine.
// It exports functions to JavaScript via syscall/js.
package main

import (
	"syscall/js"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/internal/result"
)

// version is the current version of the analysis engine.
const version = "0.1.0"

// getVersion returns the analysis engine version.
// Returns: Result<{version: string}>
func getVersion(this js.Value, args []js.Value) interface{} {
	r := result.NewSuccess(map[string]string{"version": version})
	return r.ToJSON()
}

// analyze performs dependency analysis on the provided workspace data.
// Input: JSON string of workspace configuration
// Returns: Result<AnalysisResult> - placeholder for Epic 2 implementation
func analyze(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		r := result.NewError("INVALID_INPUT", "Missing JSON input")
		return r.ToJSON()
	}

	// Placeholder - will be implemented in Epic 2
	input := args[0].String()
	_ = input // Suppress unused warning

	r := result.NewSuccess(map[string]interface{}{
		"healthScore": 100,
		"packages":    0,
		"placeholder": true,
	})
	return r.ToJSON()
}

// check validates the workspace configuration against configured rules.
// Input: JSON string of workspace configuration
// Returns: Result<CheckResult> - placeholder for Epic 2 implementation
func check(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		r := result.NewError("INVALID_INPUT", "Missing JSON input")
		return r.ToJSON()
	}

	// Placeholder - will be implemented in Epic 2
	r := result.NewSuccess(map[string]interface{}{
		"passed":      true,
		"errors":      []string{},
		"placeholder": true,
	})
	return r.ToJSON()
}

func main() {
	// Create MonoGuard namespace object
	monoguard := make(map[string]interface{})
	monoguard["getVersion"] = js.FuncOf(getVersion)
	monoguard["analyze"] = js.FuncOf(analyze)
	monoguard["check"] = js.FuncOf(check)

	// Register MonoGuard namespace globally
	js.Global().Set("MonoGuard", monoguard)

	// Keep the Go program running to handle JavaScript calls
	<-make(chan bool)
}
