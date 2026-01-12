---
description: Validate code against architecture.md requirements
argument-hint: [file-path or directory]
---

Validate that code follows the architectural decisions documented in `_bmad-output/planning-artifacts/architecture.md`.

**Usage:**

- `/monoguard:check-architecture` - Check all modified files
- `/monoguard:check-architecture apps/web` - Check specific directory
- `/monoguard:check-architecture apps/web/app/stores/analysis.ts` - Check specific file

**This command will verify:**

1. **Technology Stack Compliance:**
   - Using TanStack Start 0.34+ (not Next.js SSR features)
   - React 19.0.0 with TypeScript 5.9+
   - Zustand 4.4+ for state management
   - D3.js 7.x for visualization
   - Dexie.js 5.x for storage (NOT localStorage)

2. **Zero Backend Architecture (NFR9-NFR10):**
   - No server-side API endpoints
   - No external API calls with user code
   - All analysis happens client-side via WASM

3. **WASM Bridge Pattern:**
   - Go functions return Result<T> type
   - JSON uses camelCase (NOT snake_case)
   - Proper error handling with AnalysisError class

4. **File Structure:**
   - Follows Nx workspace conventions
   - Feature-based organization
   - Tests in correct locations (**tests**/ for TS, \*\_test.go for Go)

**Analysis Process:**

1. Read the architecture document
2. Analyze target files for architectural compliance
3. Report violations with specific file:line references
4. Suggest fixes based on architecture decisions

**Report Format:**

- ✅ Compliant patterns found
- ⚠️ Warnings (minor deviations)
- ❌ Violations (must fix)
- 💡 Recommendations

Let me analyze the code for architecture compliance.

**Target files:** $ARGUMENTS

I'll read the architecture document and check the specified files.
