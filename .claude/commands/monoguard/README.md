# MonoGuard Development Skills

MonoGuard 專案的 Claude Code skills 集合，提供開發輔助工具來確保代碼品質和架構一致性。

## 📋 Skills 列表

### 🔍 驗證類 (Validation)

#### `/monoguard:check-architecture`

檢查代碼是否符合 `architecture.md` 的架構決策。

**使用方式:**

```
/monoguard:check-architecture
/monoguard:check-architecture apps/web
/monoguard:check-architecture apps/web/app/stores/analysis.ts
```

**檢查項目:**

- ✅ 技術堆疊合規性 (React 19, TypeScript 5.9+, TanStack Start)
- ✅ Zero Backend 架構 (NFR9-NFR10)
- ✅ WASM Bridge 模式 (Result<T> type)
- ✅ 文件結構 (Nx workspace conventions)

---

#### `/monoguard:check-context`

驗證代碼是否遵循 `project-context.md` 的所有規則。

**使用方式:**

```
/monoguard:check-context
/monoguard:check-context packages/analysis-engine
/monoguard:check-context apps/web/app/components/DependencyGraph.tsx
```

**檢查項目:**

- ✅ 語言特定規則 (camelCase, PascalCase, snake_case)
- ✅ 框架特定規則 (React hooks, Zustand selectors, D3 cleanup)
- ✅ 測試規則 (Mock patterns, test organization)
- ✅ Critical Don't-Miss Rules (反模式檢查)

---

#### `/monoguard:validate-wasm`

專門驗證 Go WASM 代碼的規範。

**使用方式:**

```
/monoguard:validate-wasm
/monoguard:validate-wasm packages/analysis-engine/pkg/analyzer/workspace.go
```

**檢查項目:**

- ✅ Result<T> type 強制使用
- ✅ JSON 使用 camelCase (NOT snake_case)
- ✅ 日期使用 ISO 8601 格式
- ✅ 錯誤代碼使用 UPPER_SNAKE_CASE

---

### 🛠️ 生成類 (Scaffolding)

#### `/monoguard:generate-wasm-bridge`

生成完整的 WASM bridge 實作 (Go + TypeScript)。

**使用方式:**

```
/monoguard:generate-wasm-bridge AnalyzeWorkspace "Analyze Nx workspace structure"
/monoguard:generate-wasm-bridge DetectCycles "Detect circular dependencies"
```

**生成內容:**

- Go WASM 函數 (packages/analysis-engine/cmd/wasm/)
- TypeScript Bridge (apps/web/app/lib/wasmBridge.ts)
- TypeScript Types (packages/types/src/)
- Unit Test 模板 (Go + TypeScript)
- Integration Test 模板

---

#### `/monoguard:create-store`

生成 Zustand store 模板（帶 devtools + persist middleware）。

**使用方式:**

```
/monoguard:create-store analysis "Manage analysis state and results"
/monoguard:create-store settings "User preferences and settings"
```

**生成內容:**

- Store file with devtools + persist middleware
- Selector functions for performance
- Complete TypeScript types
- Unit tests with React Testing Library
- Usage examples

---

#### `/monoguard:scaffold-component`

生成 React 組件模板（支援 basic/d3/form 類型）。

**使用方式:**

```
/monoguard:scaffold-component DependencyGraph d3
/monoguard:scaffold-component AnalysisForm form
/monoguard:scaffold-component MetricCard basic
```

**組件類型:**

- **basic**: 標準 React functional component
- **d3**: D3.js integration with cleanup (React.memo)
- **form**: Form with validation and error handling

---

#### `/monoguard:create-test`

生成測試文件模板（帶 WASM/Zustand mocks）。

**使用方式:**

```
/monoguard:create-test apps/web/app/lib/wasmBridge.ts unit
/monoguard:create-test apps/web/app/components/AnalysisView.tsx integration
/monoguard:create-test apps/web-e2e/src/analysis-flow.spec.ts e2e
```

**測試類型:**

- **unit**: 單元測試 (All dependencies mocked, <1s)
- **integration**: 整合測試 (Multiple modules, 5-10s)
- **e2e**: E2E 測試 (Full user flow, Playwright)

---

### 📊 分析類 (Analysis)

#### `/monoguard:analyze-dependencies`

分析專案依賴健康度和合規性。

**使用方式:**

```
/monoguard:analyze-dependencies
/monoguard:analyze-dependencies packages/analysis-engine
/monoguard:analyze-dependencies apps/web
```

**分析項目:**

- ✅ 版本合規性 (React 19, TypeScript 5.9+, etc.)
- ✅ 依賴健康度 (Outdated, security vulnerabilities)
- ✅ 架構合規性 (Zero backend, no server-side deps)
- ✅ Monorepo 結構 (Circular dependencies, bundle size)

---

#### `/monoguard:check-coverage`

檢查測試覆蓋率是否達到 >80% 目標。

**使用方式:**

```
/monoguard:check-coverage
/monoguard:check-coverage apps/web
/monoguard:check-coverage packages/analysis-engine
```

**檢查項目:**

- ✅ Unit Tests: >80% coverage (Line, Branch, Function)
- ✅ Integration Tests: Core WASM bridge paths
- ✅ E2E Tests: 3-5 critical user flows
- ✅ Critical Path Coverage: WASM bridge, Zustand stores, IndexedDB

---

## 🚀 快速開始

### 典型開發流程

**1. 開始新功能前 - 檢查架構:**

```
/monoguard:check-architecture
/monoguard:check-context
```

**2. 生成代碼模板:**

```
/monoguard:generate-wasm-bridge AnalyzeCircularDeps "Detect circular dependencies"
/monoguard:create-store circularDeps "Manage circular dependency detection"
/monoguard:scaffold-component CircularDepsView d3
```

**3. 實作完成後 - 生成測試:**

```
/monoguard:create-test packages/analysis-engine/pkg/analyzer/circular.go unit
/monoguard:create-test apps/web/app/components/CircularDepsView.tsx integration
```

**4. 提交前檢查:**

```
/monoguard:validate-wasm packages/analysis-engine
/monoguard:check-coverage
/monoguard:analyze-dependencies
```

---

## 💡 最佳實踐

### 驗證類 Skills

- 在 **PR 前** 運行 check-architecture 和 check-context
- 每次修改 Go WASM 代碼後運行 validate-wasm
- 定期運行 analyze-dependencies 保持依賴健康

### 生成類 Skills

- 使用 generate-wasm-bridge 確保 WASM bridge 一致性
- 使用 create-store 生成符合規範的 Zustand stores
- 使用 scaffold-component 生成帶 cleanup 的 D3 組件

### 分析類 Skills

- 每週運行 analyze-dependencies 檢查依賴更新
- 每次提交前運行 check-coverage 確保覆蓋率
- CI/CD 集成這些 skills 進行自動檢查

---

## 🔧 Skills 集成到 CI/CD

可以在 GitHub Actions 中使用這些 skills：

```yaml
name: MonoGuard Quality Checks

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Check Architecture
        run: /monoguard:check-architecture

      - name: Check Context Rules
        run: /monoguard:check-context

      - name: Validate WASM
        run: /monoguard:validate-wasm

      - name: Check Coverage
        run: /monoguard:check-coverage

      - name: Analyze Dependencies
        run: /monoguard:analyze-dependencies
```

---

## 📚 相關文檔

- **Architecture Document**: `_bmad-output/planning-artifacts/architecture.md`
- **Project Context**: `_bmad-output/project-context.md`
- **PRD**: `_bmad-output/planning-artifacts/prd.md`

---

## 🆘 疑難排解

### Skill 無法使用？

確認 Claude Code 已載入 skills：

1. 重啟 Claude Code
2. 確認 `.claude/commands/monoguard/commands.json` 存在
3. 輸入 `/monoguard:` 查看自動完成建議

### Skills 檢查失敗？

1. 確認已讀取最新的 architecture.md 和 project-context.md
2. 檢查文件路徑是否正確（使用絕對路徑）
3. 運行 `/monoguard:check-architecture` 查看具體錯誤

---

**Created:** 2026-01-12
**Version:** 1.0.0
**Maintainer:** MonoGuard Team
