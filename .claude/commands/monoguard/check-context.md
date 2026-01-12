---
description: Verify code follows project-context.md rules
argument-hint: [file-path or directory]
---

Verify that code follows ALL critical rules documented in `_bmad-output/project-context.md`.

**Usage:**

- `/monoguard:check-context` - Check all modified files
- `/monoguard:check-context packages/analysis-engine` - Check Go WASM code
- `/monoguard:check-context apps/web/app/components/DependencyGraph.tsx` - Check specific file

**This command will verify:**

**Language-Specific Rules:**

- ✅ TypeScript: camelCase variables, PascalCase types/components
- ✅ Go: PascalCase exports, camelCase unexported, snake_case files
- ✅ JSON: camelCase everywhere (NO snake_case)
- ✅ Dates: ISO 8601 strings (NOT Unix timestamps)
- ✅ Errors: AnalysisError with layered messages
- ✅ WASM: Result<T> type mandatory

**Framework-Specific Rules:**

- ✅ React: Hooks, React.memo() for D3 components
- ✅ Zustand: Selectors (NOT entire store subscription)
- ✅ TanStack Start: NO SSR features (getServerSideProps forbidden)
- ✅ D3.js: useEffect cleanup (remove event listeners)
- ✅ Dexie.js: IndexedDB for large data (NOT localStorage)

**Testing Rules:**

- ✅ Tests in **tests**/ (TypeScript) or \*\_test.go (Go)
- ✅ WASM mocks return Result<T> structure
- ✅ Zustand store mocks provided
- ✅ IndexedDB mocked with fake-indexeddb

**Critical Don't-Miss Rules:**

- ❌ NEVER use localStorage for analysis results
- ❌ NEVER use snake_case in JSON
- ❌ NEVER forget D3.js cleanup
- ❌ NEVER use SSR features in TanStack Start
- ❌ NEVER return raw Go errors to WASM

**Check Process:**

1. Read project-context.md for all rules
2. Analyze target files line-by-line
3. Report rule violations with context
4. Show correct patterns from project-context.md

**Report Format:**

- ✅ Rules followed correctly
- ❌ Rule violations (file:line with fix)
- 💡 Best practice suggestions
- 📚 Reference to relevant project-context.md sections

Let me check the code against project context rules.

**Target files:** $ARGUMENTS

I'll read project-context.md and analyze the specified files for compliance.
