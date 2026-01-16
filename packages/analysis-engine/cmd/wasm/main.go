//go:build js && wasm

// Package main is the WASM entry point for MonoGuard analysis engine.
// It exports functions to JavaScript via syscall/js.
// Business logic is delegated to internal/handlers for testability.
package main

import (
	"syscall/js"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/internal/handlers"
)

// getVersion returns the analysis engine version.
// Returns: Result<VersionInfo>
func getVersion(this js.Value, args []js.Value) interface{} {
	return handlers.GetVersion()
}

// analyze performs dependency analysis on the provided workspace data.
// Input: JSON string of workspace configuration
// Returns: Result<AnalysisResult>
func analyze(this js.Value, args []js.Value) interface{} {
	input := ""
	if len(args) >= 1 {
		input = args[0].String()
	}
	return handlers.Analyze(input)
}

// check validates the workspace configuration against configured rules.
// Input: JSON string of workspace configuration
// Returns: Result<CheckResult>
func check(this js.Value, args []js.Value) interface{} {
	input := ""
	if len(args) >= 1 {
		input = args[0].String()
	}
	return handlers.Check(input)
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
