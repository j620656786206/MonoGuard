// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file contains fix guide generator for Story 3.4.
package analyzer

import (
	"fmt"
	"strings"

	"github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"
)

// ========================================
// Fix Guide Generator (Story 3.4)
// ========================================

// FixGuideGenerator creates step-by-step guides for fix strategies.
type FixGuideGenerator struct {
	workspace      *types.WorkspaceData
	packageManager string // "npm", "yarn", "pnpm"
}

// NewFixGuideGenerator creates a new generator.
func NewFixGuideGenerator(workspace *types.WorkspaceData) *FixGuideGenerator {
	return &FixGuideGenerator{
		workspace:      workspace,
		packageManager: detectPackageManager(workspace),
	}
}

// Generate creates a fix guide for a strategy.
func (fgg *FixGuideGenerator) Generate(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.FixGuide {
	if cycle == nil || strategy == nil {
		return nil
	}

	switch strategy.Type {
	case types.FixStrategyExtractModule:
		return fgg.generateExtractModuleGuide(cycle, strategy)
	case types.FixStrategyDependencyInject:
		return fgg.generateDIGuide(cycle, strategy)
	case types.FixStrategyBoundaryRefactor:
		return fgg.generateBoundaryRefactorGuide(cycle, strategy)
	default:
		return nil
	}
}

// ========================================
// Extract Module Guide (Task 3)
// ========================================

// generateExtractModuleGuide creates guide for Extract Shared Module strategy.
func (fgg *FixGuideGenerator) generateExtractModuleGuide(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.FixGuide {
	steps := []types.FixStep{}
	stepNum := 1

	// Get new package name
	newPkgName := strategy.NewPackageName
	if newPkgName == "" {
		newPkgName = "@mono/shared"
	}
	pkgDirName := extractPkgDirName(newPkgName)

	// Step 1: Create new shared package directory
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Create new shared package",
		Description: fmt.Sprintf("Create a new package directory for shared code: %s", newPkgName),
		Command: &types.CommandStep{
			Command:          fmt.Sprintf("mkdir -p packages/%s/src", pkgDirName),
			WorkingDirectory: ".",
			Description:      "Create the package directory structure",
		},
		ExpectedOutcome: "New package directory exists",
	})
	stepNum++

	// Step 2: Create package.json for new package
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Create package.json",
		Description: "Initialize the new package with a package.json file",
		FilePath:    fmt.Sprintf("packages/%s/package.json", pkgDirName),
		CodeAfter: &types.CodeSnippet{
			Language: "json",
			Code:     fgg.generatePackageJson(newPkgName),
		},
		ExpectedOutcome: "New package is recognized by the workspace",
	})
	stepNum++

	// Step 3: Create index.ts for new package
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Create entry point",
		Description: "Create the main entry point for the shared package",
		FilePath:    fmt.Sprintf("packages/%s/src/index.ts", pkgDirName),
		CodeAfter: &types.CodeSnippet{
			Language: "typescript",
			Code:     "// Shared code extracted to break circular dependency\n\nexport {};\n",
		},
		ExpectedOutcome: "Package has a valid entry point",
	})
	stepNum++

	// Step 4: Move shared code (for each target package)
	for _, pkg := range strategy.TargetPackages {
		pkgPath := fgg.getPackagePath(pkg)
		steps = append(steps, types.FixStep{
			Number:      stepNum,
			Title:       fmt.Sprintf("Update imports in %s", pkg),
			Description: "Update import statements to use the new shared package",
			FilePath:    fmt.Sprintf("%s/src/index.ts", pkgPath),
			CodeBefore: &types.CodeSnippet{
				Language: "typescript",
				Code:     fgg.generateBeforeImport(pkg, cycle),
			},
			CodeAfter: &types.CodeSnippet{
				Language: "typescript",
				Code:     fgg.generateAfterImport(pkg, newPkgName),
			},
			ExpectedOutcome: fmt.Sprintf("%s now imports from shared package", pkg),
		})
		stepNum++
	}

	// Step 5: Update dependencies in package.json files
	for _, pkg := range strategy.TargetPackages {
		pkgPath := fgg.getPackagePath(pkg)
		steps = append(steps, types.FixStep{
			Number:      stepNum,
			Title:       fmt.Sprintf("Add dependency in %s", pkg),
			Description: fmt.Sprintf("Add %s as a dependency", newPkgName),
			FilePath:    fmt.Sprintf("%s/package.json", pkgPath),
			CodeAfter: &types.CodeSnippet{
				Language: "json",
				Code:     fgg.generateDependencyAddition(newPkgName),
			},
			ExpectedOutcome: "Dependency is declared",
		})
		stepNum++
	}

	// Step 6: Install dependencies
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Install dependencies",
		Description: "Update the workspace to recognize the new package",
		Command: &types.CommandStep{
			Command:          fgg.getInstallCommand(),
			WorkingDirectory: ".",
			Description:      "Install and link workspace packages",
		},
		ExpectedOutcome: "All packages are linked correctly",
	})

	return &types.FixGuide{
		StrategyType:  types.FixStrategyExtractModule,
		Title:         fmt.Sprintf("Extract Shared Module: %s", newPkgName),
		Summary:       "Create a new shared package to hold common dependencies and break the circular dependency.",
		Steps:         steps,
		Verification:  fgg.generateVerificationSteps(cycle),
		Rollback:      fgg.generateRollbackInstructions(),
		EstimatedTime: fgg.estimateTime(strategy.Effort),
	}
}

// ========================================
// Dependency Injection Guide (Task 4)
// ========================================

// generateDIGuide creates guide for Dependency Injection strategy.
func (fgg *FixGuideGenerator) generateDIGuide(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.FixGuide {
	steps := []types.FixStep{}
	stepNum := 1

	// Get critical edge from root cause analysis
	var criticalEdge *types.RootCauseEdge
	if cycle.RootCause != nil && cycle.RootCause.CriticalEdge != nil {
		criticalEdge = cycle.RootCause.CriticalEdge
	}

	// Determine packages involved
	fromPkg := ""
	toPkg := ""
	if criticalEdge != nil {
		fromPkg = criticalEdge.From
		toPkg = criticalEdge.To
	} else if len(strategy.TargetPackages) >= 2 {
		fromPkg = strategy.TargetPackages[0]
		toPkg = strategy.TargetPackages[1]
	} else if len(cycle.Cycle) >= 2 {
		fromPkg = cycle.Cycle[0]
		toPkg = cycle.Cycle[1]
	}

	// Step 1: Create interface/type definition
	interfaceFilePath := fgg.getInterfaceFilePath(toPkg)
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Create interface for dependency",
		Description: "Define an interface that abstracts the dependency that needs to be inverted",
		FilePath:    interfaceFilePath,
		CodeAfter: &types.CodeSnippet{
			Language: "typescript",
			Code:     fgg.generateInterfaceCode(toPkg),
		},
		ExpectedOutcome: "Interface is defined and exported",
	})
	stepNum++

	// Step 2: Update the dependent package to use interface
	if fromPkg != "" {
		fromPkgPath := fgg.getPackagePath(fromPkg)
		steps = append(steps, types.FixStep{
			Number:      stepNum,
			Title:       fmt.Sprintf("Update %s to accept dependency via injection", fromPkg),
			Description: "Modify the code to receive the dependency through a parameter or constructor",
			FilePath:    fmt.Sprintf("%s/src/index.ts", fromPkgPath),
			CodeBefore: &types.CodeSnippet{
				Language: "typescript",
				Code:     fgg.generateDIBeforeCode(fromPkg, toPkg),
			},
			CodeAfter: &types.CodeSnippet{
				Language: "typescript",
				Code:     fgg.generateDIAfterCode(fromPkg, toPkg),
			},
			ExpectedOutcome: "Package no longer directly imports the problematic dependency",
		})
		stepNum++
	}

	// Step 3: Wire up the dependency at the composition root
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Wire up the dependency at the composition root",
		Description: "Inject the concrete implementation where the components are assembled",
		FilePath:    "apps/main/src/index.ts",
		CodeAfter: &types.CodeSnippet{
			Language: "typescript",
			Code:     fgg.generateCompositionCode(fromPkg, toPkg),
		},
		ExpectedOutcome: "Dependency is properly injected",
	})
	stepNum++

	// Step 4: Remove direct import
	if fromPkg != "" && toPkg != "" {
		fromPkgPath := fgg.getPackagePath(fromPkg)
		steps = append(steps, types.FixStep{
			Number:      stepNum,
			Title:       "Remove the circular import",
			Description: "Delete the import statement that was causing the cycle",
			FilePath:    fmt.Sprintf("%s/src/index.ts", fromPkgPath),
			CodeBefore: &types.CodeSnippet{
				Language: "typescript",
				Code:     fmt.Sprintf("import { something } from '%s';", toPkg),
			},
			CodeAfter: &types.CodeSnippet{
				Language: "typescript",
				Code:     "// Import removed - dependency is now injected",
			},
			ExpectedOutcome: "Circular import is eliminated",
		})
	}

	return &types.FixGuide{
		StrategyType:  types.FixStrategyDependencyInject,
		Title:         "Dependency Injection: Invert the Dependency",
		Summary:       "Break the cycle by inverting the problematic dependency using an interface and injection pattern.",
		Steps:         steps,
		Verification:  fgg.generateVerificationSteps(cycle),
		Rollback:      fgg.generateRollbackInstructions(),
		EstimatedTime: fgg.estimateTime(strategy.Effort),
	}
}

// ========================================
// Boundary Refactoring Guide (Task 5)
// ========================================

// generateBoundaryRefactorGuide creates guide for Boundary Refactoring strategy.
func (fgg *FixGuideGenerator) generateBoundaryRefactorGuide(
	cycle *types.CircularDependencyInfo,
	strategy *types.FixStrategy,
) *types.FixGuide {
	steps := []types.FixStep{}
	stepNum := 1

	// Step 1: Analyze current package responsibilities
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Identify overlapping responsibilities",
		Description: "Review the packages involved to identify code that belongs to the wrong package. Look for functionality that is used by multiple packages and could be causing the circular dependency.",
		ExpectedOutcome: "Clear understanding of what code needs to move",
	})
	stepNum++

	// Step 2: For each package that needs restructuring
	for _, pkg := range strategy.TargetPackages {
		pkgPath := fgg.getPackagePath(pkg)
		steps = append(steps, types.FixStep{
			Number:      stepNum,
			Title:       fmt.Sprintf("Restructure %s boundaries", pkg),
			Description: fmt.Sprintf("Move code that doesn't belong in %s to the appropriate package. Consider what functionality is core to this package vs what was added for convenience.", pkg),
			FilePath:    fmt.Sprintf("%s/src/index.ts", pkgPath),
			ExpectedOutcome: fmt.Sprintf("%s only contains code appropriate to its domain", pkg),
		})
		stepNum++
	}

	// Step 3: Update exports
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Update package exports",
		Description: "Ensure each package exports only what belongs to its domain. Remove exports that have been moved to other packages.",
		ExpectedOutcome: "Clean public API for each package",
	})
	stepNum++

	// Step 4: Update consumers
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Update import statements in consumers",
		Description: "Update any code that imports moved functionality to use the new locations.",
		Command: &types.CommandStep{
			Command:          "grep -r \"from.*@mono\" --include='*.ts' --include='*.tsx'",
			WorkingDirectory: ".",
			Description:      "Find all imports that may need updating",
		},
		ExpectedOutcome: "All imports point to correct packages",
	})
	stepNum++

	// Step 5: Run build to verify
	steps = append(steps, types.FixStep{
		Number:      stepNum,
		Title:       "Verify build passes",
		Description: "Run the build to ensure all imports resolve correctly",
		Command: &types.CommandStep{
			Command:          fgg.getBuildCommand(),
			WorkingDirectory: ".",
			Description:      "Build all packages",
		},
		ExpectedOutcome: "Build completes without errors",
	})

	return &types.FixGuide{
		StrategyType:  types.FixStrategyBoundaryRefactor,
		Title:         "Module Boundary Refactoring",
		Summary:       "Restructure package boundaries to eliminate overlapping responsibilities and break the cycle.",
		Steps:         steps,
		Verification:  fgg.generateVerificationSteps(cycle),
		Rollback:      fgg.generateRollbackInstructions(),
		EstimatedTime: fgg.estimateTime(strategy.Effort),
	}
}

// ========================================
// Verification Steps (Task 6)
// ========================================

// generateVerificationSteps creates steps to verify the fix worked.
func (fgg *FixGuideGenerator) generateVerificationSteps(cycle *types.CircularDependencyInfo) []types.FixStep {
	cycleStr := formatCycle(cycle.Cycle)

	return []types.FixStep{
		{
			Number:      1,
			Title:       "Run MonoGuard analysis",
			Description: "Verify the circular dependency is resolved",
			Command: &types.CommandStep{
				Command:          "npx monoguard analyze",
				WorkingDirectory: ".",
				Description:      "Run dependency analysis",
			},
			ExpectedOutcome: fmt.Sprintf("Cycle %s should no longer appear in results", cycleStr),
		},
		{
			Number:      2,
			Title:       "Run build",
			Description: "Ensure the project still builds successfully",
			Command: &types.CommandStep{
				Command:          fgg.getBuildCommand(),
				WorkingDirectory: ".",
				Description:      "Build all packages",
			},
			ExpectedOutcome: "Build completes without errors",
		},
		{
			Number:      3,
			Title:       "Run tests",
			Description: "Verify no regressions were introduced",
			Command: &types.CommandStep{
				Command:          fgg.getTestCommand(),
				WorkingDirectory: ".",
				Description:      "Run test suite",
			},
			ExpectedOutcome: "All tests pass",
		},
	}
}

// ========================================
// Rollback Instructions (Task 7)
// ========================================

// generateRollbackInstructions creates instructions to undo changes.
func (fgg *FixGuideGenerator) generateRollbackInstructions() *types.RollbackInstructions {
	return &types.RollbackInstructions{
		GitCommands: []string{
			"git stash  # Save any uncommitted changes",
			"git checkout .  # Discard all changes",
			"# OR to revert a specific commit:",
			"git revert <commit-hash>",
		},
		ManualSteps: []string{
			"1. Restore original import statements",
			"2. Delete any newly created packages",
			"3. Remove added dependencies from package.json files",
			"4. Run install command to update lockfile",
		},
		Warning: "If you've already pushed changes, coordinate with your team before reverting.",
	}
}

// ========================================
// Package Manager Detection (Task 8)
// ========================================

// detectPackageManager determines npm/yarn/pnpm from workspace.
func detectPackageManager(workspace *types.WorkspaceData) string {
	if workspace == nil {
		return "npm"
	}

	switch workspace.WorkspaceType {
	case types.WorkspaceTypePnpm:
		return "pnpm"
	case types.WorkspaceTypeYarn:
		return "yarn"
	case types.WorkspaceTypeNpm:
		return "npm"
	default:
		return "npm" // Default fallback
	}
}

// getInstallCommand returns the install command for the package manager.
func (fgg *FixGuideGenerator) getInstallCommand() string {
	switch fgg.packageManager {
	case "pnpm":
		return "pnpm install"
	case "yarn":
		return "yarn install"
	default:
		return "npm install"
	}
}

// getBuildCommand returns the build command for the package manager.
func (fgg *FixGuideGenerator) getBuildCommand() string {
	switch fgg.packageManager {
	case "pnpm":
		return "pnpm run build"
	case "yarn":
		return "yarn build"
	default:
		return "npm run build"
	}
}

// getTestCommand returns the test command for the package manager.
func (fgg *FixGuideGenerator) getTestCommand() string {
	switch fgg.packageManager {
	case "pnpm":
		return "pnpm test"
	case "yarn":
		return "yarn test"
	default:
		return "npm test"
	}
}

// ========================================
// Time Estimation
// ========================================

// estimateTime returns estimated time based on effort level.
func (fgg *FixGuideGenerator) estimateTime(effort types.EffortLevel) string {
	switch effort {
	case types.EffortLow:
		return "15-30 minutes"
	case types.EffortMedium:
		return "30-60 minutes"
	case types.EffortHigh:
		return "60-120 minutes"
	default:
		return "30-60 minutes"
	}
}

// ========================================
// Helper Functions
// ========================================

// extractPkgDirName extracts directory name from package name.
// e.g., "@mono/shared" -> "shared"
func extractPkgDirName(pkgName string) string {
	parts := strings.Split(pkgName, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return pkgName
}

// getPackagePath returns the path to a package.
func (fgg *FixGuideGenerator) getPackagePath(pkgName string) string {
	if fgg.workspace != nil && fgg.workspace.Packages != nil {
		if pkg, ok := fgg.workspace.Packages[pkgName]; ok && pkg.Path != "" {
			return pkg.Path
		}
	}
	// Fallback: derive from package name
	dirName := extractPkgDirName(pkgName)
	return fmt.Sprintf("packages/%s", dirName)
}

// generatePackageJson generates a basic package.json for a new package.
func (fgg *FixGuideGenerator) generatePackageJson(pkgName string) string {
	return fmt.Sprintf(`{
  "name": "%s",
  "version": "0.1.0",
  "main": "src/index.ts",
  "types": "src/index.ts",
  "dependencies": {}
}`, pkgName)
}

// generateBeforeImport generates sample before import code.
func (fgg *FixGuideGenerator) generateBeforeImport(pkg string, cycle *types.CircularDependencyInfo) string {
	// Find what this package imports from the cycle
	for i := 0; i < len(cycle.Cycle)-1; i++ {
		if cycle.Cycle[i] == pkg && i+1 < len(cycle.Cycle) {
			target := cycle.Cycle[i+1]
			return fmt.Sprintf("import { sharedFunction } from '%s';", target)
		}
	}
	return "import { sharedFunction } from './other-package';"
}

// generateAfterImport generates sample after import code.
func (fgg *FixGuideGenerator) generateAfterImport(pkg string, newPkgName string) string {
	return fmt.Sprintf("import { sharedFunction } from '%s';", newPkgName)
}

// generateDependencyAddition generates the dependency addition snippet.
func (fgg *FixGuideGenerator) generateDependencyAddition(pkgName string) string {
	return fmt.Sprintf(`"dependencies": {
  "%s": "workspace:*"
}`, pkgName)
}

// getInterfaceFilePath returns the path for the interface file.
func (fgg *FixGuideGenerator) getInterfaceFilePath(pkgName string) string {
	pkgPath := fgg.getPackagePath(pkgName)
	return fmt.Sprintf("%s/src/types.ts", pkgPath)
}

// generateInterfaceCode generates the interface code.
func (fgg *FixGuideGenerator) generateInterfaceCode(pkgName string) string {
	interfaceName := generateInterfaceName(pkgName)
	return fmt.Sprintf(`export interface %s {
  // Define the contract for the dependency
  execute(): void;
}
`, interfaceName)
}

// generateInterfaceName creates an interface name from package name.
func generateInterfaceName(pkgName string) string {
	dirName := extractPkgDirName(pkgName)
	// Capitalize first letter
	if len(dirName) > 0 {
		return strings.ToUpper(dirName[:1]) + dirName[1:] + "Handler"
	}
	return "Handler"
}

// generateDIBeforeCode generates before code for DI.
func (fgg *FixGuideGenerator) generateDIBeforeCode(fromPkg, toPkg string) string {
	return fmt.Sprintf(`import { something } from '%s';

export function doWork() {
  something();
}`, toPkg)
}

// generateDIAfterCode generates after code for DI.
func (fgg *FixGuideGenerator) generateDIAfterCode(fromPkg, toPkg string) string {
	interfaceName := generateInterfaceName(toPkg)
	return fmt.Sprintf(`import type { %s } from '%s/types';

let handler: %s | null = null;

export function setHandler(h: %s) {
  handler = h;
}

export function doWork() {
  if (handler) {
    handler.execute();
  }
}`, interfaceName, toPkg, interfaceName, interfaceName)
}

// generateCompositionCode generates composition root code.
func (fgg *FixGuideGenerator) generateCompositionCode(fromPkg, toPkg string) string {
	return fmt.Sprintf(`import { setHandler } from '%s';
import { concreteImplementation } from '%s';

// Wire up dependencies at the composition root
setHandler(concreteImplementation);`, fromPkg, toPkg)
}

// formatCycle formats a cycle for display.
func formatCycle(cycle []string) string {
	if len(cycle) == 0 {
		return ""
	}
	return strings.Join(cycle, " → ")
}
