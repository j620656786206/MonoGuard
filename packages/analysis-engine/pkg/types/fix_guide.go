// Package types defines Go types that match TypeScript definitions in @monoguard/types.
// This file contains fix guide types for Story 3.4.
package types

// ========================================
// Fix Guide Types (Story 3.4)
// ========================================

// FixGuide provides step-by-step instructions for implementing a fix strategy.
// Matches @monoguard/types FixGuide interface.
type FixGuide struct {
	// StrategyType links this guide to a specific strategy
	StrategyType FixStrategyType `json:"strategyType"`

	// Title is the guide headline
	Title string `json:"title"`

	// Summary is a brief overview of what this guide accomplishes
	Summary string `json:"summary"`

	// Steps are the ordered implementation instructions
	Steps []FixStep `json:"steps"`

	// Verification contains steps to confirm the fix worked
	Verification []FixStep `json:"verification"`

	// Rollback contains instructions to undo the changes
	Rollback *RollbackInstructions `json:"rollback,omitempty"`

	// EstimatedTime is the approximate time to complete (e.g., "15-30 minutes")
	EstimatedTime string `json:"estimatedTime"`
}

// FixStep represents a single step in the fix guide.
type FixStep struct {
	// Number is the step number (1-based)
	Number int `json:"number"`

	// Title is a short description of this step
	Title string `json:"title"`

	// Description provides detailed instructions
	Description string `json:"description"`

	// FilePath is the file to modify (if applicable)
	FilePath string `json:"filePath,omitempty"`

	// CodeBefore shows the current code (if applicable)
	CodeBefore *CodeSnippet `json:"codeBefore,omitempty"`

	// CodeAfter shows the desired code (if applicable)
	CodeAfter *CodeSnippet `json:"codeAfter,omitempty"`

	// Command is a terminal command to run (if applicable)
	Command *CommandStep `json:"command,omitempty"`

	// ExpectedOutcome describes what should happen after this step
	ExpectedOutcome string `json:"expectedOutcome,omitempty"`
}

// CodeSnippet represents a code example.
type CodeSnippet struct {
	// Language is the syntax highlighting hint (e.g., "typescript", "json")
	Language string `json:"language"`

	// Code is the actual code content
	Code string `json:"code"`

	// StartLine is the approximate line number (for context)
	StartLine int `json:"startLine,omitempty"`
}

// CommandStep represents a terminal command.
type CommandStep struct {
	// Command is the exact command to run
	Command string `json:"command"`

	// WorkingDirectory is where to run the command (relative to workspace root)
	WorkingDirectory string `json:"workingDirectory,omitempty"`

	// Description explains what this command does
	Description string `json:"description,omitempty"`
}

// RollbackInstructions provides steps to undo changes.
type RollbackInstructions struct {
	// GitCommands are git commands to revert (if in a git repo)
	GitCommands []string `json:"gitCommands,omitempty"`

	// ManualSteps are non-git rollback instructions
	ManualSteps []string `json:"manualSteps,omitempty"`

	// Warning is a caution message about rollback
	Warning string `json:"warning,omitempty"`
}
