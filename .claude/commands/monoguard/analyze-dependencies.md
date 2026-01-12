---
description: Analyze project dependency health and compliance
argument-hint: [package-path]
---

Analyze MonoGuard's dependency health, version compliance, and potential issues.

**Usage:**

- `/monoguard:analyze-dependencies` - Analyze all packages
- `/monoguard:analyze-dependencies packages/analysis-engine` - Analyze specific package
- `/monoguard:analyze-dependencies apps/web` - Analyze web app dependencies

**What This Analyzes:**

**1. Version Compliance:**

- ✅ React 19.0.0 (required for TypeScript 5.9+ compatibility)
- ✅ TypeScript 5.9.2+ (required for React 19 types)
- ✅ TanStack Start 0.34+ (required for stable SSG)
- ✅ Zustand 4.4+ (required for React 19 compatibility)
- ✅ Go 1.21+ (required for WASM)
- ✅ Node.js >= 18.0.0
- ✅ pnpm 10.14.0

**2. Dependency Health:**

- 🔴 Outdated major versions (breaking changes)
- 🟡 Outdated minor versions (new features available)
- 🟢 Up-to-date dependencies
- ⚠️ Security vulnerabilities
- 📦 Unused dependencies

**3. Architecture Compliance:**

- ❌ Server-side dependencies in client-only app (violates zero backend)
- ❌ Conflicting versions across packages
- ❌ Direct dependencies that should be workspace dependencies
- ⚠️ Large bundle sizes (>2MB for WASM)

**4. Monorepo Structure:**

- 📊 Dependency graph visualization
- 🔗 Inter-package dependencies
- 🚨 Circular dependencies between packages
- 📈 Bundle size analysis per package

**Analysis Process:**

1. **Read package.json files:**
   - Root workspace package.json
   - All app and package package.json files
   - Identify direct and dev dependencies

2. **Check version compliance:**
   - Compare against architecture.md requirements
   - Flag critical version mismatches
   - Check for peer dependency conflicts

3. **Security scan:**
   - Run `pnpm audit` for vulnerabilities
   - Check for known CVEs
   - Recommend updates for security patches

4. **Dependency analysis:**
   - Identify unused dependencies
   - Check for duplicate dependencies (different versions)
   - Analyze bundle impact

5. **Architecture validation:**
   - Ensure no server-side dependencies in client app
   - Validate WASM build dependencies (Go)
   - Check for forbidden dependencies

**Report Format:**

```
🔍 MonoGuard Dependency Analysis

📊 Overview:
- Total packages: 5
- Total dependencies: 45
- Critical issues: 2
- Warnings: 3

⚠️ Critical Issues:

❌ apps/web/package.json:12
   React version: 18.2.0 (Required: 19.0.0+)
   Fix: pnpm add react@19.0.0 --filter @monoguard/web

❌ packages/types/package.json:8
   TypeScript version: 5.3.0 (Required: 5.9.2+)
   Fix: pnpm add -D typescript@5.9.2 --filter @monoguard/types

🟡 Warnings:

⚠️ apps/web/package.json:15
   Zustand version: 4.3.8 (Recommended: 4.4.0+ for React 19)
   Fix: pnpm add zustand@latest --filter @monoguard/web

📦 Unused Dependencies:

⚠️ apps/web/package.json:20
   lodash: Not imported in any file
   Fix: pnpm remove lodash --filter @monoguard/web

🔐 Security:

✅ No known vulnerabilities found

📈 Bundle Analysis:

✅ WASM bundle: 1.8MB (Target: <2MB)
⚠️ Web app bundle: 450KB (Consider code splitting)

🎯 Recommendations:

1. Update React to 19.0.0 (required for architecture)
2. Update TypeScript to 5.9.2+ (required for React 19 types)
3. Remove unused dependencies to reduce bundle size
4. Run `pnpm dedupe` to remove duplicate dependencies
```

**Commands Generated:**

The analysis will provide ready-to-run commands to fix issues:

```bash
# Fix critical version issues
pnpm add react@19.0.0 --filter @monoguard/web
pnpm add -D typescript@5.9.2 --filter @monoguard/types

# Remove unused dependencies
pnpm remove lodash --filter @monoguard/web

# Update to recommended versions
pnpm add zustand@latest --filter @monoguard/web

# Dedupe dependencies
pnpm dedupe
```

Let me analyze the project dependencies: **$ARGUMENTS**

I'll read all package.json files and check for compliance with architecture requirements.
