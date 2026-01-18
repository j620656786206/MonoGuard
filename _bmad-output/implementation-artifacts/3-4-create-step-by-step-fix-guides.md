# Story 3.4: Create Step-by-Step Fix Guides

Status: ready-for-dev

## Story

As a **user**,
I want **detailed step-by-step guides for each fix strategy**,
So that **I can follow clear instructions to resolve the circular dependency**.

## Acceptance Criteria

1. **AC1: Numbered Steps**
   - Given a fix strategy from Story 3.3
   - When I request the step-by-step guide
   - Then the guide includes:
     - Numbered steps (1, 2, 3...)
     - Clear action description for each step
     - Expected outcome after completing each step
   - And steps are ordered logically for implementation

2. **AC2: Specific File Paths**
   - Given the packages involved in the cycle
   - When generating the guide
   - Then each step that modifies files includes:
     - Specific file paths to modify (e.g., `packages/ui/src/index.ts`)
     - Package.json paths for dependency changes
   - And paths are relative to workspace root

3. **AC3: Code Snippets (Before/After)**
   - Given a step that requires code changes
   - When showing the change
   - Then include:
     - "Before" code snippet showing current state
     - "After" code snippet showing desired state
     - Syntax highlighting hint (language type)
   - And snippets are contextual to the actual packages

4. **AC4: Commands to Run**
   - Given a step that requires terminal commands
   - When generating the guide
   - Then include:
     - Exact command to run (e.g., `pnpm install`, `npm run build`)
     - Working directory for the command
     - Expected output or success criteria
   - And commands are specific to detected package manager

5. **AC5: Verification Steps**
   - Given a completed fix guide
   - When including verification
   - Then add steps to confirm the fix worked:
     - Run `monoguard analyze` to verify cycle is broken
     - Run build/test commands to ensure no regressions
     - Check specific imports are updated
   - And verification steps come at the end of the guide

6. **AC6: Rollback Instructions**
   - Given a fix guide
   - When generating instructions
   - Then include rollback section with:
     - Git commands to revert changes (if applicable)
     - Manual steps to undo changes
     - Warning about potential side effects
   - And rollback is clearly separated from main steps

7. **AC7: Integration with FixStrategy**
   - Given analysis results
   - When enriching FixStrategy
   - Then add `guide` field:
     ```go
     type FixStrategy struct {
         // ... existing fields ...
         Guide *FixGuide `json:"guide,omitempty"`
     }
     ```
   - And guide generation is optional (can be requested separately)

8. **AC8: Performance**
   - Given a workspace with 100 packages and 5 cycles
   - When generating all fix guides
   - Then generation completes in < 500ms additional overhead
   - And memory usage increase is < 10MB

## Tasks / Subtasks

- [ ] **Task 1: Define FixGuide Types** (AC: #1, #2, #3, #4, #5, #6)
  - [ ] 1.1 Create `pkg/types/fix_guide.go`:
    ```go
    package types

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
    ```
  - [ ] 1.2 Add JSON serialization tests in `pkg/types/fix_guide_test.go`
  - [ ] 1.3 Ensure all JSON tags use camelCase

- [ ] **Task 2: Create FixGuideGenerator** (AC: #1, #7)
  - [ ] 2.1 Create `pkg/analyzer/fix_guide_generator.go`:
    ```go
    package analyzer

    import "github.com/j620656786206/MonoGuard/packages/analysis-engine/pkg/types"

    // FixGuideGenerator creates step-by-step guides for fix strategies.
    type FixGuideGenerator struct {
        workspace      *types.WorkspaceData
        packageManager string // "npm", "yarn", "pnpm"
    }

    // NewFixGuideGenerator creates a new generator.
    func NewFixGuideGenerator(workspace *types.WorkspaceData) *FixGuideGenerator

    // Generate creates a fix guide for a strategy.
    func (fgg *FixGuideGenerator) Generate(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.FixGuide

    // generateExtractModuleGuide creates guide for Extract Shared Module strategy.
    func (fgg *FixGuideGenerator) generateExtractModuleGuide(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.FixGuide

    // generateDIGuide creates guide for Dependency Injection strategy.
    func (fgg *FixGuideGenerator) generateDIGuide(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.FixGuide

    // generateBoundaryRefactorGuide creates guide for Boundary Refactoring strategy.
    func (fgg *FixGuideGenerator) generateBoundaryRefactorGuide(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.FixGuide

    // detectPackageManager determines npm/yarn/pnpm from workspace.
    func detectPackageManager(workspace *types.WorkspaceData) string
    ```
  - [ ] 2.2 Implement guide generation dispatch logic
  - [ ] 2.3 Create comprehensive tests in `pkg/analyzer/fix_guide_generator_test.go`

- [ ] **Task 3: Implement Extract Module Guide** (AC: #1, #2, #3, #4, #5)
  - [ ] 3.1 Implement `generateExtractModuleGuide`:
    ```go
    func (fgg *FixGuideGenerator) generateExtractModuleGuide(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.FixGuide {
        steps := []types.FixStep{}
        stepNum := 1

        // Step 1: Create new shared package directory
        steps = append(steps, types.FixStep{
            Number:      stepNum,
            Title:       "Create new shared package",
            Description: fmt.Sprintf("Create a new package directory for shared code: %s", strategy.NewPackageName),
            Command: &types.CommandStep{
                Command:          fmt.Sprintf("mkdir -p packages/%s/src", extractPkgDirName(strategy.NewPackageName)),
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
            FilePath:    fmt.Sprintf("packages/%s/package.json", extractPkgDirName(strategy.NewPackageName)),
            CodeAfter: &types.CodeSnippet{
                Language: "json",
                Code:     fgg.generatePackageJson(strategy.NewPackageName),
            },
            ExpectedOutcome: "New package is recognized by the workspace",
        })
        stepNum++

        // Step 3: Move shared code to new package
        for _, pkg := range strategy.TargetPackages {
            steps = append(steps, types.FixStep{
                Number:      stepNum,
                Title:       fmt.Sprintf("Update imports in %s", pkg),
                Description: fmt.Sprintf("Update import statements to use the new shared package"),
                FilePath:    fgg.getMainEntryFile(pkg),
                CodeBefore: &types.CodeSnippet{
                    Language: "typescript",
                    Code:     fgg.generateBeforeImport(pkg, cycle),
                },
                CodeAfter: &types.CodeSnippet{
                    Language: "typescript",
                    Code:     fgg.generateAfterImport(pkg, strategy.NewPackageName),
                },
                ExpectedOutcome: fmt.Sprintf("%s now imports from shared package", pkg),
            })
            stepNum++
        }

        // Step 4: Update dependencies in package.json files
        for _, pkg := range strategy.TargetPackages {
            steps = append(steps, types.FixStep{
                Number:      stepNum,
                Title:       fmt.Sprintf("Add dependency in %s", pkg),
                Description: fmt.Sprintf("Add %s as a dependency", strategy.NewPackageName),
                FilePath:    fgg.getPackageJsonPath(pkg),
                CodeAfter: &types.CodeSnippet{
                    Language: "json",
                    Code:     fgg.generateDependencyAddition(strategy.NewPackageName),
                },
                ExpectedOutcome: "Dependency is declared",
            })
            stepNum++
        }

        // Step 5: Install dependencies
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
            Title:         fmt.Sprintf("Extract Shared Module: %s", strategy.NewPackageName),
            Summary:       "Create a new shared package to hold common dependencies and break the circular dependency.",
            Steps:         steps,
            Verification:  fgg.generateVerificationSteps(cycle),
            Rollback:      fgg.generateRollbackInstructions(),
            EstimatedTime: fgg.estimateTime(strategy.Effort),
        }
    }
    ```
  - [ ] 3.2 Implement helper functions for code generation
  - [ ] 3.3 Add tests for Extract Module guide

- [ ] **Task 4: Implement Dependency Injection Guide** (AC: #1, #2, #3, #4, #5)
  - [ ] 4.1 Implement `generateDIGuide`:
    ```go
    func (fgg *FixGuideGenerator) generateDIGuide(
        cycle *types.CircularDependencyInfo,
        strategy *types.FixStrategy,
    ) *types.FixGuide {
        steps := []types.FixStep{}
        stepNum := 1

        // Get critical edge from root cause analysis
        var criticalEdge *types.DependencyEdge
        if cycle.RootCause != nil && cycle.RootCause.CriticalEdge != nil {
            criticalEdge = cycle.RootCause.CriticalEdge
        }

        // Step 1: Create interface/type definition
        steps = append(steps, types.FixStep{
            Number:      stepNum,
            Title:       "Create interface for dependency",
            Description: "Define an interface that abstracts the dependency that needs to be inverted",
            FilePath:    fgg.getInterfaceFilePath(criticalEdge),
            CodeAfter: &types.CodeSnippet{
                Language: "typescript",
                Code:     fgg.generateInterfaceCode(criticalEdge),
            },
            ExpectedOutcome: "Interface is defined and exported",
        })
        stepNum++

        // Step 2: Update the dependent package to use interface
        if criticalEdge != nil {
            steps = append(steps, types.FixStep{
                Number:      stepNum,
                Title:       fmt.Sprintf("Update %s to accept dependency via injection", criticalEdge.From),
                Description: "Modify the code to receive the dependency through a parameter or constructor",
                FilePath:    fgg.getMainEntryFile(criticalEdge.From),
                CodeBefore: &types.CodeSnippet{
                    Language: "typescript",
                    Code:     fgg.generateDIBeforeCode(criticalEdge),
                },
                CodeAfter: &types.CodeSnippet{
                    Language: "typescript",
                    Code:     fgg.generateDIAfterCode(criticalEdge),
                },
                ExpectedOutcome: "Package no longer directly imports the problematic dependency",
            })
            stepNum++

            // Step 3: Update the calling code to provide dependency
            steps = append(steps, types.FixStep{
                Number:      stepNum,
                Title:       "Wire up the dependency at the composition root",
                Description: "Inject the concrete implementation where the components are assembled",
                FilePath:    fgg.getCompositionRootPath(cycle),
                CodeAfter: &types.CodeSnippet{
                    Language: "typescript",
                    Code:     fgg.generateCompositionCode(criticalEdge),
                },
                ExpectedOutcome: "Dependency is properly injected",
            })
            stepNum++
        }

        // Step 4: Remove direct import
        steps = append(steps, types.FixStep{
            Number:      stepNum,
            Title:       "Remove the circular import",
            Description: "Delete the import statement that was causing the cycle",
            FilePath:    fgg.getCriticalImportFile(criticalEdge),
            CodeBefore: &types.CodeSnippet{
                Language: "typescript",
                Code:     fgg.generateCircularImportBefore(criticalEdge),
            },
            CodeAfter: &types.CodeSnippet{
                Language: "typescript",
                Code:     "// Import removed - dependency is now injected",
            },
            ExpectedOutcome: "Circular import is eliminated",
        })

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
    ```
  - [ ] 4.2 Implement DI-specific helper functions
  - [ ] 4.3 Add tests for DI guide

- [ ] **Task 5: Implement Boundary Refactoring Guide** (AC: #1, #2, #3, #4, #5)
  - [ ] 5.1 Implement `generateBoundaryRefactorGuide`:
    ```go
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
            Description: "Review the packages involved to identify code that belongs to the wrong package",
            ExpectedOutcome: "Clear understanding of what code needs to move",
        })
        stepNum++

        // Step 2: For each package that needs restructuring
        for _, pkg := range strategy.TargetPackages {
            steps = append(steps, types.FixStep{
                Number:      stepNum,
                Title:       fmt.Sprintf("Restructure %s boundaries", pkg),
                Description: fmt.Sprintf("Move code that doesn't belong in %s to the appropriate package", pkg),
                FilePath:    fgg.getMainEntryFile(pkg),
                ExpectedOutcome: fmt.Sprintf("%s only contains code appropriate to its domain", pkg),
            })
            stepNum++
        }

        // Step 3: Update exports
        steps = append(steps, types.FixStep{
            Number:      stepNum,
            Title:       "Update package exports",
            Description: "Ensure each package exports only what belongs to its domain",
            ExpectedOutcome: "Clean public API for each package",
        })
        stepNum++

        // Step 4: Update consumers
        steps = append(steps, types.FixStep{
            Number:      stepNum,
            Title:       "Update import statements in consumers",
            Description: "Update any code that imports moved functionality",
            Command: &types.CommandStep{
                Command:     "grep -r 'from.*@mono' --include='*.ts' --include='*.tsx'",
                Description: "Find all imports that may need updating",
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
    ```
  - [ ] 5.2 Implement boundary refactor helper functions
  - [ ] 5.3 Add tests for Boundary Refactoring guide

- [ ] **Task 6: Implement Verification Steps** (AC: #5)
  - [ ] 6.1 Implement `generateVerificationSteps`:
    ```go
    func (fgg *FixGuideGenerator) generateVerificationSteps(cycle *types.CircularDependencyInfo) []types.FixStep {
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
                ExpectedOutcome: fmt.Sprintf("Cycle %s should no longer appear in results", formatCycle(cycle.Cycle)),
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
    ```
  - [ ] 6.2 Add tests for verification steps

- [ ] **Task 7: Implement Rollback Instructions** (AC: #6)
  - [ ] 7.1 Implement `generateRollbackInstructions`:
    ```go
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
    ```
  - [ ] 7.2 Add tests for rollback instructions

- [ ] **Task 8: Implement Package Manager Detection** (AC: #4)
  - [ ] 8.1 Implement `detectPackageManager`:
    ```go
    func detectPackageManager(workspace *types.WorkspaceData) string {
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
    ```
  - [ ] 8.2 Add tests for package manager detection

- [ ] **Task 9: Integrate with FixStrategy** (AC: #7)
  - [ ] 9.1 Update `pkg/types/fix_strategy.go`:
    ```go
    type FixStrategy struct {
        // ... existing fields ...
        Guide *FixGuide `json:"guide,omitempty"` // NEW: Step-by-step guide
    }
    ```
  - [ ] 9.2 Verify existing tests still pass

- [ ] **Task 10: Wire to Analyzer Pipeline** (AC: all)
  - [ ] 10.1 Update `pkg/analyzer/analyzer.go`:
    ```go
    func (a *Analyzer) AnalyzeWithSources(...) (*types.AnalysisResult, error) {
        // ... existing analysis ...

        // Generate fix strategies (Story 3.3)
        fixGenerator := NewFixStrategyGenerator(graph, workspace)
        guideGenerator := NewFixGuideGenerator(workspace) // NEW

        for _, cycle := range cycles {
            strategies := fixGenerator.Generate(cycle)

            // NEW: Generate guides for each strategy
            for i := range strategies {
                strategies[i].Guide = guideGenerator.Generate(cycle, &strategies[i])
            }

            cycle.FixStrategies = strategies
        }

        return result, nil
    }
    ```
  - [ ] 10.2 Update analyzer tests

- [ ] **Task 11: Update TypeScript Types** (AC: #7)
  - [ ] 11.1 Update `packages/types/src/analysis/results.ts`:
    ```typescript
    export interface FixGuide {
      strategyType: FixStrategyType;
      title: string;
      summary: string;
      steps: FixStep[];
      verification: FixStep[];
      rollback?: RollbackInstructions;
      estimatedTime: string;
    }

    export interface FixStep {
      number: number;
      title: string;
      description: string;
      filePath?: string;
      codeBefore?: CodeSnippet;
      codeAfter?: CodeSnippet;
      command?: CommandStep;
      expectedOutcome?: string;
    }

    export interface CodeSnippet {
      language: string;
      code: string;
      startLine?: number;
    }

    export interface CommandStep {
      command: string;
      workingDirectory?: string;
      description?: string;
    }

    export interface RollbackInstructions {
      gitCommands?: string[];
      manualSteps?: string[];
      warning?: string;
    }

    export interface FixStrategy {
      // ... existing fields ...
      guide?: FixGuide; // NEW
    }
    ```
  - [ ] 11.2 Run `pnpm nx build types` to verify
  - [ ] 11.3 Update type tests if needed

- [ ] **Task 12: Performance Testing** (AC: #8)
  - [ ] 12.1 Create `pkg/analyzer/fix_guide_generator_benchmark_test.go`:
    ```go
    func BenchmarkFixGuideGeneration(b *testing.B) {
        workspace := generateWorkspace(100)
        cycles := generateCyclesWithStrategies(5)
        generator := NewFixGuideGenerator(workspace)

        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            for _, cycle := range cycles {
                for _, strategy := range cycle.FixStrategies {
                    generator.Generate(cycle, &strategy)
                }
            }
        }
    }
    ```
  - [ ] 12.2 Verify < 500ms for 100 packages with 5 cycles
  - [ ] 12.3 Document actual performance in completion notes

- [ ] **Task 13: Integration Verification** (AC: all)
  - [ ] 13.1 Run all tests: `cd packages/analysis-engine && make test`
  - [ ] 13.2 Build WASM: `pnpm nx build @monoguard/analysis-engine`
  - [ ] 13.3 Run affected CI checks: `pnpm nx affected --target=lint,test,type-check --base=main`
  - [ ] 13.4 Verify JSON output includes guide field

## Dev Notes

### Architecture Patterns & Constraints

**From Architecture Document:**
- **Location:** Fix guide generator in `pkg/analyzer/`
- **Pattern:** Generator pattern with workspace + strategy input
- **Integration:** Enriches FixStrategy with optional Guide
- **Dependency:** Requires Story 3.3 (FixStrategy) to be implemented first

**Critical Constraints:**
- **camelCase JSON:** All struct tags MUST use camelCase
- **Optional Fields:** Guide is `omitempty` - backward compatible
- **Contextual Output:** Code snippets must use actual package names, not placeholders
- **Package Manager Aware:** Commands should match detected npm/yarn/pnpm

### Critical Don't-Miss Rules

**From project-context.md:**

1. **JSON Naming Convention:**
   ```go
   // ✅ CORRECT: camelCase JSON tags
   type FixStep struct {
       FilePath        string       `json:"filePath,omitempty"`
       CodeBefore      *CodeSnippet `json:"codeBefore,omitempty"`
       ExpectedOutcome string       `json:"expectedOutcome,omitempty"`
   }

   // ❌ WRONG: snake_case JSON tags
   type FixStep struct {
       FilePath string `json:"file_path"` // WRONG!
   }
   ```

2. **Pointer vs Value for Optional Structs:**
   ```go
   // ✅ CORRECT: Use pointer for optional nested structs
   type FixStep struct {
       CodeBefore *CodeSnippet `json:"codeBefore,omitempty"`
       Command    *CommandStep `json:"command,omitempty"`
   }

   // ❌ WRONG: Value type with omitempty doesn't work for structs
   type FixStep struct {
       CodeBefore CodeSnippet `json:"codeBefore,omitempty"` // Won't omit!
   }
   ```

3. **Slice Initialization:**
   ```go
   // ✅ CORRECT: Initialize as empty slice
   steps := []types.FixStep{}

   // ❌ WRONG: Nil slice (serializes as null)
   var steps []types.FixStep // nil
   ```

### Project Structure Notes

**Target Directory Structure:**
```
packages/analysis-engine/
├── pkg/
│   ├── analyzer/
│   │   ├── analyzer.go                        # UPDATE: Add guide generation
│   │   ├── fix_strategy_generator.go          # From Story 3.3
│   │   ├── fix_guide_generator.go             # NEW: Guide generator
│   │   ├── fix_guide_generator_test.go        # NEW: Generator tests
│   │   └── fix_guide_generator_benchmark_test.go # NEW: Performance
│   └── types/
│       ├── fix_strategy.go                    # UPDATE: Add Guide field
│       ├── fix_guide.go                       # NEW: FixGuide types
│       └── fix_guide_test.go                  # NEW: Type tests
└── ...

packages/types/src/analysis/
├── results.ts                                 # UPDATE: Add TS types
└── ...
```

### Previous Story Intelligence

**From Story 3.3 (ready-for-dev):**
- FixStrategy has: type, name, description, suitability, effort, pros, cons, targetPackages, newPackageName
- FixStrategyType: extract-module, dependency-injection, boundary-refactoring
- EffortLevel: low, medium, high
- **Key Insight:** Use targetPackages for file path generation
- **Key Insight:** Use newPackageName for Extract Module guide

**From Story 3.1 (ready-for-dev):**
- RootCauseAnalysis has criticalEdge which is key for DI guide
- **Key Insight:** Use criticalEdge to determine injection point

**From Story 2.1 (done):**
- WorkspaceData has workspaceType (npm/yarn/pnpm)
- PackageInfo has path field for package location
- **Key Insight:** Use workspace.Packages[name].Path for file paths

### Guide Generation Examples

**Extract Module Guide Example:**
```
Title: Extract Shared Module: @mono/shared

Step 1: Create new shared package
  Command: mkdir -p packages/shared/src

Step 2: Create package.json
  File: packages/shared/package.json
  After: { "name": "@mono/shared", "version": "0.1.0", ... }

Step 3: Update imports in @mono/ui
  File: packages/ui/src/index.ts
  Before: import { helper } from '@mono/core'
  After: import { helper } from '@mono/shared'

Step 4: Add dependency in @mono/ui
  File: packages/ui/package.json
  After: "dependencies": { "@mono/shared": "workspace:*" }

Step 5: Install dependencies
  Command: pnpm install

Verification:
  1. Run monoguard analyze
  2. Run pnpm build
  3. Run pnpm test
```

**Dependency Injection Guide Example:**
```
Title: Dependency Injection: Invert the Dependency

Step 1: Create interface for dependency
  File: packages/core/src/types.ts
  After: export interface UIHandler { render(): void }

Step 2: Update @mono/core to accept dependency via injection
  File: packages/core/src/index.ts
  Before: import { render } from '@mono/ui'; render();
  After: export function init(handler: UIHandler) { handler.render(); }

Step 3: Wire up the dependency at the composition root
  File: packages/app/src/main.ts
  After: import { init } from '@mono/core'; init(uiHandler);

Step 4: Remove the circular import
  File: packages/core/src/index.ts
  Before: import { render } from '@mono/ui'
  After: // Import removed - dependency is now injected
```

### Input/Output Format

**Input (FixStrategy from Story 3.3):**
```json
{
  "type": "extract-module",
  "name": "Extract Shared Module",
  "suitability": 8,
  "effort": "medium",
  "targetPackages": ["@mono/ui", "@mono/api", "@mono/core"],
  "newPackageName": "@mono/shared"
}
```

**Output (FixStrategy with Guide):**
```json
{
  "type": "extract-module",
  "name": "Extract Shared Module",
  "suitability": 8,
  "effort": "medium",
  "targetPackages": ["@mono/ui", "@mono/api", "@mono/core"],
  "newPackageName": "@mono/shared",
  "guide": {
    "strategyType": "extract-module",
    "title": "Extract Shared Module: @mono/shared",
    "summary": "Create a new shared package to hold common dependencies.",
    "steps": [
      {
        "number": 1,
        "title": "Create new shared package",
        "description": "Create a new package directory for shared code",
        "command": {
          "command": "mkdir -p packages/shared/src",
          "workingDirectory": ".",
          "description": "Create the package directory structure"
        },
        "expectedOutcome": "New package directory exists"
      }
    ],
    "verification": [...],
    "rollback": {
      "gitCommands": ["git checkout .", "git revert <commit>"],
      "warning": "Coordinate with team before reverting pushed changes"
    },
    "estimatedTime": "30-60 minutes"
  }
}
```

### Test Scenarios

| Strategy | Guide Should Include |
|----------|---------------------|
| Extract Module | mkdir, package.json creation, import updates, dependency additions |
| Dependency Injection | Interface creation, injection setup, import removal |
| Boundary Refactoring | Responsibility analysis, code movement, export updates |

### References

- [Source: _bmad-output/planning-artifacts/epics.md#Story 3.4]
- [Source: _bmad-output/planning-artifacts/prd.md#FR10]
- [Source: _bmad-output/project-context.md#Language-Specific Rules]
- [Source: _bmad-output/implementation-artifacts/3-3-generate-fix-strategy-recommendations.md]
- [Refactoring to Patterns - Joshua Kerievsky](https://industriallogic.com/xp/refactoring/)

## Dev Agent Record

### Agent Model Used

{{agent_model_name_version}}

### Debug Log References

### Completion Notes List

### File List

### Change Log
