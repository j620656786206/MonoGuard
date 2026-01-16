---
stepsCompleted: [1, 2, 3, 4, 5, 6, 7, 8]
inputDocuments:
  - '_bmad-output/planning-artifacts/prd.md'
  - '_bmad-output/planning-artifacts/validation-report-prd.md'
  - '_bmad-output/planning-artifacts/ux-design-specification.md'
workflowType: 'architecture'
project_name: 'mono-guard'
user_name: 'Alexyu'
date: '2026-01-12'
lastStep: 8
status: 'complete'
completedAt: '2026-01-12'
---

# Architecture Decision Document

_This document builds collaboratively through step-by-step discovery. Sections are appended as we work through each architectural decision together._

## Project Context Analysis

### Requirements Overview

**Functional Requirements (48 total):**

MonoGuard 的功能需求圍繞三個核心能力：

1. **Dependency Analysis & Detection (FR1-FR6)**
   - Monorepo workspace 檢測（npm/yarn/pnpm）
   - 完整依賴圖建構
   - 循環依賴識別
   - 架構健康分數計算
   - **架構影響：** 需要強大的靜態分析引擎，支援多種 package manager 格式

2. **Circular Dependency Resolution - 核心差異化 (FR7-FR14)**
   - 根因分析與影響評估
   - 修復策略建議（Extract Module, DI, Boundary Refactoring）
   - 逐步修復指南
   - 重構複雜度評分
   - **架構影響：** 需要規則引擎或模式匹配系統，Phase 2 可能需要 AI 整合

3. **Dual Interface (FR15-FR33)**
   - Web Interface: 拖放上傳、WASM 執行、互動視覺化
   - CLI Tool: 分析、檢查、修復預覽、CI/CD 整合
   - **架構影響：** 共享核心分析引擎，但介面層分離

4. **Privacy-First Architecture (FR34-FR39)**
   - 完全離線分析
   - 本地儲存（IndexedDB + `.monoguard/` 目錄）
   - 選擇性遙測
   - **架構影響：** 無後端依賴，所有處理在 client-side

5. **Integration & API (FR40-FR48)**
   - 可配置規則和閾值
   - WASM API 供第三方整合
   - **架構影響：** 需要清晰的 API 邊界和 TypeScript 型別定義

**Non-Functional Requirements (17 total):**

關鍵 NFR 將驅動架構決策：

1. **Performance (NFR1-NFR4)**
   - 分析速度：100 packages < 5s, 1000 packages < 30s
   - UI 響應：< 500ms 互動回應
   - Bundle size：< 500KB gzipped
   - Memory：< 100MB (WASM in-browser)
   - **架構影響：** WASM 編譯優化、分批處理、漸進式渲染

2. **Reliability (NFR5-NFR8)**
   - 100% 離線可用性
   - P95 錯誤率 < 0.1%
   - 修復建議接受率 > 60% (Phase 0), > 80% (Phase 1)
   - **架構影響：** 錯誤處理策略、優雅降級、規則引擎準確度驗證

3. **Security & Privacy (NFR9-NFR12)**
   - 零程式碼上傳
   - 本地優先儲存
   - 選擇性遙測
   - **架構影響：** 完全 client-side 架構，無後端 API

4. **Integration (NFR13-NFR15)**
   - 支援 npm/yarn/pnpm workspaces
   - CI/CD 整合（GitHub Actions, GitLab CI 等）
   - 多種匯出格式（JSON, HTML, Markdown）
   - **架構影響：** 彈性的輸入解析器、標準化輸出格式

5. **Scalability (NFR16-NFR17)**
   - 基礎設施成本：$0/月（Render Free Tier）
   - 支援 10,000 併發使用者
   - 優雅降級（> 2000 packages 建議使用 CLI）
   - **架構影響：** 靜態部署策略、分批處理機制

**Scale & Complexity:**

- **Primary domain:** Full-stack Developer Tool
- **Complexity level:** Medium
- **Estimated architectural components:** 6-8 major components
- **Phase 0 MVP focus:** Analysis Engine + Visualization + Fix Suggestions Level 1
- **Phase 1 expansion:** Time Machine + GitHub PR Integration
- **Phase 2 scale:** Team Dashboard + AI Diagnostics + Enterprise features

### Technical Constraints & Dependencies

**Known Constraints:**

1. **Zero Backend Constraint**
   - NFR9-NFR10 要求完全本地分析
   - 不能依賴 server-side 處理
   - **影響：** 所有運算必須在 client-side 完成（WASM + Browser JS）

2. **Zero Cost Infrastructure (NFR16)**
   - 必須使用免費層服務
   - Render 為首選（Web + API + DB 統一管理）
   - **影響：** 純靜態部署，無 server-side rendering 或 API routes

3. **Performance Targets (NFR1-NFR4)**
   - 嚴格的分析速度和 bundle size 限制
   - **影響：** WASM 優化、code splitting、lazy loading 必要

4. **Offline-First (NFR5)**
   - 所有核心功能必須 100% 離線可用
   - **影響：** Service Worker、IndexedDB、無網路依賴

5. **Browser Compatibility**
   - WASM 支援：Chrome 57+, Firefox 52+, Safari 11+
   - IndexedDB 支援：所有現代瀏覽器
   - **影響：** 需考慮 polyfills 或優雅降級

**Technology Stack Indicators from PRD:**

- **Frontend:** TanStack Start（SSG 模式）
- **Analysis Engine:** Go（編譯為 WASM）
- **Visualization:** D3.js
- **Storage:** IndexedDB（Web）+ local files（CLI）
- **Deployment:** Render (Web + API + PostgreSQL + Redis)
- **CLI:** Go native binary

### Cross-Cutting Concerns Identified

以下關注點將影響多個架構元件：

1. **Privacy & Data Handling**
   - 影響：All components
   - 要求：零資料外洩、本地儲存、選擇性遙測
   - 架構決策：完全 client-side 架構、無後端 API

2. **Performance Optimization**
   - 影響：Analysis Engine, Visualization, Web UI
   - 要求：< 5s 分析、< 500ms 互動、< 500KB bundle
   - 架構決策：WASM 優化、分批處理、code splitting

3. **Error Handling & Resilience**
   - 影響：All components
   - 要求：< 0.1% 錯誤率、優雅降級
   - 架構決策：防禦性編程、錯誤邊界、fallback 機制

4. **Platform Consistency**
   - 影響：Web UI, CLI
   - 要求：一致的使用體驗和術語
   - 架構決策：共享核心引擎、統一設計語言

5. **Testing Strategy**
   - 影響：All components
   - 挑戰：WASM 模組測試、視覺化測試、大型 monorepo 測試資料
   - 架構決策：需要測試架構和策略（將在後續步驟決定）

6. **Observability (Optional Telemetry)**
   - 影響：All components
   - 要求：選擇性、透明、尊重隱私
   - 架構決策：PostHog（client-side）、Sentry（錯誤追蹤）

7. **Deployment & Distribution**
   - 影響：Web UI, CLI, API
   - Web：靜態部署到 Render Static Site
   - API：Go 服務部署到 Render Web Service
   - Database：Render PostgreSQL + Redis
   - CLI：npm global install（Go binary）
   - 架構決策：All-in-one Render Blueprint（render.yaml）統一管理

## Starter Template Evaluation

### Primary Technology Domain

**Full-stack Developer Tool** 基於專案需求分析：

- Web UI（靜態生成，離線優先）
- WASM 分析引擎（客戶端執行）
- CLI 工具（本地分析 + CI/CD 整合）

### Starter Options Considered

由於專案的獨特架構需求（WASM + 零後端 + 雙介面），評估後發現：

1. **TanStack Start 官方 Starter**
   - 用途：Web UI 前端基礎
   - 狀態：官方維護，生產就緒
   - 版本：0.34.11（最新）
   - 優勢：完整 SSG 支援、Render Static Site 友善

2. **Go 自訂 WASM 專案**
   - 用途：分析引擎核心
   - 狀態：標準 Go 工具鏈支援
   - 優勢：成熟的 WASM 編譯流程、優秀效能

3. **Go CLI with Cobra/Viper**
   - 用途：CLI 工具介面
   - 狀態：業界標準（Kubernetes、Docker 使用）
   - 優勢：豐富的功能、強大的配置管理

### Selected Starter Strategy: Hybrid Multi-Repository Approach

**理由：**

MonoGuard 的架構需求（WASM + 靜態部署 + CLI）決定了沒有單一 starter 能滿足所有需求。採用**混合策略**：

1. **Web UI**: 使用 TanStack Start 官方 starter
2. **Analysis Engine**: 自訂 Go WASM 專案
3. **CLI Tool**: 使用 Cobra/Viper 自訂 Go 專案

這種分離策略符合專案的技術約束，並允許每個元件使用最適合的工具鏈。

### Initialization Commands

#### 1. Web UI (TanStack Start)

```bash
# 建立 TanStack Start 專案
npm create @tanstack/start@latest

# 互動式選項：
# - Project name: mono-guard-web
# - Toolchain: Biome（推薦，更快）
# - Add-ons: Tailwind CSS（選擇性，視 UX 需求）
```

**專案結構：**

```
mono-guard-web/
├── app/
│   ├── routes/
│   ├── components/
│   └── styles/
├── public/
├── vite.config.ts
└── package.json
```

#### 2. Analysis Engine (Go WASM)

```bash
# 建立 Go WASM 分析引擎
mkdir analysis-engine
cd analysis-engine
go mod init github.com/alexyu/mono-guard/analysis-engine

# 基本專案結構
mkdir -p cmd/wasm pkg/{parser,analyzer,rules}

# WASM 編譯設置
# 在 Makefile 或 build 腳本中：
GOOS=js GOARCH=wasm go build -o dist/monoguard.wasm cmd/wasm/main.go

# 複製 Go WASM 執行器
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/
```

**專案結構：**

```
analysis-engine/
├── cmd/
│   └── wasm/
│       └── main.go
├── pkg/
│   ├── parser/      # Workspace 解析器
│   ├── analyzer/    # 依賴圖分析
│   └── rules/       # 修復建議規則引擎
├── go.mod
└── Makefile
```

#### 3. CLI Tool (Go with Cobra)

```bash
# 使用 Cobra CLI generator
go install github.com/spf13/cobra-cli@latest
cobra-cli init mono-guard-cli

# 新增主要命令
cd mono-guard-cli
cobra-cli add analyze
cobra-cli add check
cobra-cli add fix
cobra-cli add init
```

**專案結構：**

```
mono-guard-cli/
├── cmd/
│   ├── root.go
│   ├── analyze.go
│   ├── check.go
│   ├── fix.go
│   └── init.go
├── pkg/
│   ├── config/      # Viper 配置管理
│   └── engine/      # 共享分析引擎介面
├── go.mod
└── main.go
```

### Architectural Decisions Provided by Starters

#### Language & Runtime

**Web UI (TanStack Start):**

- TypeScript 預設配置
- Modern ES modules
- Node.js 18+ 執行環境

**Go Projects:**

- Go 1.21+（WASI 支援）
- 純 Go 標準庫（最小依賴）
- WASM 目標：`GOOS=js GOARCH=wasm`

#### Build Tooling

**Web UI:**

- **Vite** - 極快的開發伺服器和構建工具
- **Nitro** - 靜態預渲染引擎
- **PostCSS** - CSS 處理（如選用 Tailwind）
- Tree-shaking 和 code splitting 內建

**Go Projects:**

- 標準 Go 工具鏈（`go build`, `go test`）
- Makefile 或 build scripts 管理編譯
- 跨平台編譯支援（macOS, Linux, Windows）

#### Styling Solution

**TanStack Start 選項：**

- **Tailwind CSS**（推薦）- utility-first，bundle size 優化
- **Plain CSS Modules** - 零依賴選項
- **Styled Components/Emotion** - CSS-in-JS（如 UX 需要）

**D3.js 視覺化：**

- D3.js v7（最新穩定版）
- SVG 渲染（輕量級可視化）
- Canvas 渲染（大規模圖表，> 1000 nodes）

#### Testing Framework

**Web UI:**

- **Vitest**（TanStack 生態推薦）- Vite 原生整合
- **Testing Library** - 元件測試
- **Playwright** 或 **Cypress** - E2E 測試（後續階段）

**Go Projects:**

- 標準 `testing` 套件
- **Testify** - 斷言和 mock 輔助
- **Go WASM 測試** - 使用 Node.js 或瀏覽器環境

#### Code Organization

**Monorepo 結構（推薦）：**

```
mono-guard/
├── apps/
│   ├── web/              # TanStack Start 前端
│   └── cli/              # Go CLI 工具
├── packages/
│   ├── analysis-engine/  # Go WASM 核心
│   ├── types/            # 共享 TypeScript 型別
│   └── ui-components/    # 共享 React 元件
├── docs/
└── package.json          # Root workspace 配置
```

**優勢：**

- 程式碼共享（型別定義、視覺化元件）
- 統一版本管理
- 協同開發流暢

#### Development Experience

**Web UI (TanStack Start):**

- ⚡ HMR（Hot Module Replacement）- 極快熱重載
- 🎨 TypeScript 智能提示
- 🐛 Vite 錯誤 overlay
- 📦 自動依賴安裝偵測

**Go Projects:**

- 🔄 Air 或 CompileDaemon - 自動重新編譯（dev 模式）
- 🧪 `go test -v` - 詳細測試輸出
- 📊 Go pprof - 效能分析工具
- 🐛 Delve - Go debugger

#### Configuration Management

**Web UI:**

- `vite.config.ts` - 構建配置
- `.env` 檔案 - 環境變數
- `tanstack.config.ts` - SSG 預渲染配置

**CLI Tool (Viper):**

- `.monoguard.json` - 使用者專案配置
- 環境變數支援
- 命令列 flags 優先權最高
- Home directory 配置（`~/.monoguard/config.yaml`）

### Integration Points

**WASM ↔ Web UI:**

- JavaScript 透過 `WebAssembly.instantiate()` 載入
- 使用 `wasm_exec.js` 提供 Go-JS 橋接
- TypeScript 型別定義包裝 WASM API

**CLI ↔ Analysis Engine:**

- 選項 A：CLI 直接呼叫 Go 分析程式碼（共享套件）
- 選項 B：CLI 載入 WASM（統一引擎，但效能稍低）

### Deployment Strategy

**決策變更記錄 (2026-01-16):** 從 Cloudflare Pages 改為 Render，原因如下：

- All-in-one 部署體驗：Web + API + PostgreSQL + Redis 統一管理
- `render.yaml` Blueprint 實現 Infrastructure as Code
- 簡化 CI/CD 流程，單一平台管理所有服務

**Web UI → Render Static Site:**

- 靜態 HTML/CSS/JS 輸出（Vite build → `.output/`）
- WASM 檔案作為靜態資源（`public/monoguard.wasm`）
- Headers 配置 COOP/COEP（WASM SharedArrayBuffer 需求）
- SPA fallback routing（`/* → /index.html`）

**API → Render Web Service:**

- Go API 服務（Gin framework）
- Health check endpoint: `/health`
- 自動連接 PostgreSQL + Redis

**Database → Render PostgreSQL + Redis:**

- PostgreSQL: 持久化資料儲存
- Redis: 快取層（allkeys-lru 策略）

**CLI → npm Registry:**

- Go binary 包裝為 npm package
- Platform-specific postinstall 腳本
- 跨平台 binary 下載（macOS, Linux, Windows）

### Next Steps

**Phase 0 實作優先順序：**

1. **Week 1-2**: TanStack Start 專案初始化 + 基本路由
2. **Week 2-3**: Go WASM 分析引擎核心（workspace 解析 + 依賴圖）
3. **Week 3-4**: Web UI 整合 WASM + D3.js 視覺化
4. **Week 4-5**: CLI 工具基本命令（analyze, check）
5. **Week 5-6**: 修復建議規則引擎（Level 1）

**Note:** 專案初始化應該是第一個實作 story，建立好基礎結構後再逐步開發功能。

## Core Architectural Decisions

以下是透過協作決策流程確定的核心架構選擇，每個決策都經過選項分析、權衡評估，並符合專案的技術約束與目標。

### Decision 1: Monorepo Strategy

**選擇：Nx Monorepo**

**理由：**

- 專案已經使用 Nx，延續現有架構可避免遷移成本
- 你具備 Nx 使用經驗，可快速上手
- Nx 提供強大的建置快取和任務編排功能
- 支援多語言專案（TypeScript + Go）

**實作細節：**

```json
// nx.json
{
  "affected": {
    "defaultBase": "main"
  },
  "targetDefaults": {
    "build": {
      "dependsOn": ["^build"],
      "cache": true
    },
    "test": {
      "cache": true
    }
  },
  "workspaceLayout": {
    "appsDir": "apps",
    "libsDir": "packages"
  }
}
```

**專案結構：**

```
mono-guard/                    # Nx workspace root
├── apps/
│   ├── web/                   # TanStack Start 前端
│   │   ├── app/
│   │   ├── public/
│   │   └── package.json
│   └── cli/                   # Go CLI 工具
│       ├── cmd/
│       ├── pkg/
│       └── go.mod
├── packages/
│   ├── analysis-engine/       # Go WASM 核心
│   │   ├── cmd/wasm/
│   │   ├── pkg/
│   │   └── go.mod
│   ├── types/                 # 共享 TypeScript 型別
│   │   ├── src/
│   │   └── package.json
│   └── ui-components/         # 共享 React 元件
│       ├── src/
│       └── package.json
├── nx.json
├── package.json
└── tsconfig.base.json
```

**影響範圍：**

- ✅ 程式碼共享（型別、元件、工具）
- ✅ 統一依賴管理
- ✅ 建置快取加速開發
- ⚠️ 需要維護 Nx 配置

**替代方案（已排除）：**

- pnpm workspace - 雖然輕量，但專案已使用 Nx，無需遷移

---

### Decision 2: WASM Integration Mode

**選擇：Dynamic Loading + TypeScript Wrapper (Phase 0)**

**理由：**

- 簡化初期開發，減少複雜度
- 主執行緒模式對大多數使用情境已足夠（< 1000 packages）
- 為 Phase 1 的 Web Worker 升級預留彈性

**實作細節：**

**Go WASM 輸出：**

```go
// cmd/wasm/main.go
package main

import (
    "syscall/js"
    "github.com/alexyu/mono-guard/pkg/analyzer"
)

func analyzeWorkspace(this js.Value, args []js.Value) interface{} {
    workspaceDataJSON := args[0].String()

    result, err := analyzer.Analyze(workspaceDataJSON)
    if err != nil {
        return map[string]interface{}{
            "error": err.Error(),
        }
    }

    return js.ValueOf(result)
}

func main() {
    c := make(chan struct{}, 0)

    // 註冊 JS 可呼叫的函式
    js.Global().Set("analyzeWorkspace", js.FuncOf(analyzeWorkspace))
    js.Global().Set("detectCycles", js.FuncOf(detectCycles))
    js.Global().Set("suggestFixes", js.FuncOf(suggestFixes))

    <-c
}
```

**TypeScript Wrapper：**

```typescript
// packages/types/src/wasm-adapter.ts
export class MonoGuardAnalyzer {
  private wasmInstance: WebAssembly.Instance | null = null;
  private isReady = false;

  async init(): Promise<void> {
    if (this.isReady) return;

    // 載入 Go WASM runtime
    const go = new Go();
    const response = await fetch('/monoguard.wasm');
    const result = await WebAssembly.instantiateStreaming(
      response,
      go.importObject
    );

    this.wasmInstance = result.instance;
    go.run(this.wasmInstance);
    this.isReady = true;
  }

  analyze(workspaceData: WorkspaceData): AnalysisResult {
    if (!this.isReady) {
      throw new Error('WASM not initialized. Call init() first.');
    }

    // 呼叫 WASM 暴露的 JS 函式
    const result = (window as any).analyzeWorkspace(
      JSON.stringify(workspaceData)
    );

    return JSON.parse(result);
  }

  detectCycles(graph: DependencyGraph): CircularDependency[] {
    const result = (window as any).detectCycles(JSON.stringify(graph));
    return JSON.parse(result);
  }

  suggestFixes(cycles: CircularDependency[]): FixSuggestion[] {
    const result = (window as any).suggestFixes(JSON.stringify(cycles));
    return JSON.parse(result);
  }
}

// Singleton instance
export const analyzer = new MonoGuardAnalyzer();
```

**Web UI 使用：**

```typescript
// apps/web/app/routes/analyze.tsx
import { analyzer } from '@mono-guard/types/wasm-adapter';
import { useAnalysisStore } from '@/stores/analysis';

export default function AnalyzePage() {
  const { startAnalysis, isAnalyzing } = useAnalysisStore();

  useEffect(() => {
    // 初始化 WASM
    analyzer.init().catch(console.error);
  }, []);

  const handleAnalyze = async (workspaceData: WorkspaceData) => {
    await startAnalysis(workspaceData);
  };

  return (
    <div>
      {isAnalyzing ? <LoadingSpinner /> : <AnalyzeForm onSubmit={handleAnalyze} />}
    </div>
  );
}
```

**Phase 1 升級路徑（Web Worker）：**

```typescript
// Phase 1: Move to Web Worker for non-blocking analysis
// packages/types/src/wasm-worker.ts
export class MonoGuardAnalyzerWorker {
  private worker: Worker;

  constructor() {
    this.worker = new Worker(new URL('./analyzer.worker.ts', import.meta.url), {
      type: 'module',
    });
  }

  async analyze(workspaceData: WorkspaceData): Promise<AnalysisResult> {
    return new Promise((resolve, reject) => {
      this.worker.postMessage({ type: 'analyze', data: workspaceData });

      this.worker.onmessage = (e) => {
        if (e.data.type === 'result') {
          resolve(e.data.result);
        } else if (e.data.type === 'error') {
          reject(e.data.error);
        }
      };
    });
  }
}
```

**效能考量：**

- Phase 0: < 1000 packages 分析約 5-10s（主執行緒）
- Phase 1: > 1000 packages 移至 Web Worker（避免 UI 凍結）
- WASM bundle size: ~2-3MB（gzipped ~500KB）

**影響範圍：**

- ✅ 簡化開發流程
- ✅ 保持 TypeScript 型別安全
- ⚠️ 主執行緒模式可能在大型專案時阻塞 UI（Phase 1 解決）

**替代方案（已排除）：**

- Web Worker from start - 過度工程，Phase 0 不需要

---

### Decision 3: State Management

**選擇：Zustand**

**理由：**

- 極輕量（< 5KB），符合 bundle size 限制
- API 簡潔，學習曲線低
- 支援 middleware（persist, devtools）
- 適合中小型應用狀態管理

**實作細節：**

**主要 Store 設計：**

```typescript
// apps/web/app/stores/analysis.ts
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { analyzer } from '@mono-guard/types/wasm-adapter';
import type {
  AnalysisResult,
  WorkspaceData,
  CircularDependency,
  FixSuggestion,
} from '@mono-guard/types';

interface AnalysisState {
  // Data
  result: AnalysisResult | null;
  selectedNode: string | null;
  filters: {
    showCircular: boolean;
    showExternal: boolean;
    minHealthScore: number;
  };

  // UI State
  isAnalyzing: boolean;
  error: string | null;

  // Actions
  startAnalysis: (data: WorkspaceData) => Promise<void>;
  clearResult: () => void;
  selectNode: (nodeId: string) => void;
  updateFilters: (filters: Partial<AnalysisState['filters']>) => void;
}

export const useAnalysisStore = create<AnalysisState>()(
  devtools(
    persist(
      (set, get) => ({
        // Initial state
        result: null,
        selectedNode: null,
        filters: {
          showCircular: true,
          showExternal: false,
          minHealthScore: 0,
        },
        isAnalyzing: false,
        error: null,

        // Actions
        startAnalysis: async (data) => {
          set({ isAnalyzing: true, error: null });

          try {
            const result = await analyzer.analyze(data);
            set({ result, isAnalyzing: false });
          } catch (error) {
            set({
              error: error instanceof Error ? error.message : 'Analysis failed',
              isAnalyzing: false,
            });
          }
        },

        clearResult: () => {
          set({ result: null, selectedNode: null, error: null });
        },

        selectNode: (nodeId) => {
          set({ selectedNode: nodeId });
        },

        updateFilters: (newFilters) => {
          set((state) => ({
            filters: { ...state.filters, ...newFilters },
          }));
        },
      }),
      {
        name: 'monoguard-analysis',
        partialize: (state) => ({
          // 僅持久化 filters，不儲存 result（可能很大）
          filters: state.filters,
        }),
      }
    )
  )
);
```

**UI Settings Store：**

```typescript
// apps/web/app/stores/settings.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface SettingsState {
  theme: 'light' | 'dark' | 'system';
  visualizationMode: 'svg' | 'canvas' | 'auto';
  enableTelemetry: boolean;

  setTheme: (theme: SettingsState['theme']) => void;
  setVisualizationMode: (mode: SettingsState['visualizationMode']) => void;
  setTelemetry: (enabled: boolean) => void;
}

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set) => ({
      theme: 'system',
      visualizationMode: 'auto',
      enableTelemetry: false,

      setTheme: (theme) => set({ theme }),
      setVisualizationMode: (mode) => set({ visualizationMode: mode }),
      setTelemetry: (enabled) => set({ enableTelemetry: enabled }),
    }),
    {
      name: 'monoguard-settings',
    }
  )
);
```

**Usage in Components：**

```typescript
// apps/web/app/components/AnalysisPanel.tsx
import { useAnalysisStore } from '@/stores/analysis';

export function AnalysisPanel() {
  const { result, isAnalyzing, filters, updateFilters } = useAnalysisStore();

  return (
    <div>
      {isAnalyzing && <LoadingSpinner />}
      {result && (
        <>
          <FilterControls
            filters={filters}
            onChange={updateFilters}
          />
          <DependencyGraph data={result.graph} />
        </>
      )}
    </div>
  );
}
```

**Middleware 配置：**

- **devtools**: Redux DevTools 整合（開發模式）
- **persist**: LocalStorage 持久化（僅設定和篩選器）

**影響範圍：**

- ✅ 輕量級解決方案
- ✅ 型別安全
- ✅ DevTools 支援
- ⚠️ 大規模狀態可能需要拆分多個 stores

**替代方案（已排除）：**

- Redux Toolkit - 功能強大但過於複雜，不符合專案規模
- Jotai/Recoil - 原子化狀態，MonoGuard 的狀態結構更適合單一 store

---

### Decision 4: Styling Solution

**選擇：Tailwind CSS with JIT Mode**

**理由：**

- utility-first 加速開發
- JIT 模式僅生成使用到的樣式，符合 bundle size 限制
- 設計系統一致性高
- 與 TanStack Start 完美整合

**實作細節：**

**Tailwind 配置：**

```javascript
// apps/web/tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './app/**/*.{js,ts,jsx,tsx}',
    '../../packages/ui-components/src/**/*.{js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        // MonoGuard 品牌色彩
        brand: {
          primary: '#3B82F6', // Blue-500
          secondary: '#10B981', // Green-500
          danger: '#EF4444', // Red-500
          warning: '#F59E0B', // Amber-500
        },
        // Health Score Gradient
        health: {
          critical: '#DC2626', // < 40
          poor: '#F59E0B', // 40-60
          fair: '#FBBF24', // 60-75
          good: '#10B981', // 75-90
          excellent: '#059669', // > 90
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['Fira Code', 'monospace'],
      },
    },
  },
  plugins: [require('@tailwindcss/forms'), require('@tailwindcss/typography')],
};
```

**Design System Components：**

```typescript
// packages/ui-components/src/Button.tsx
import { cva, type VariantProps } from 'class-variance-authority';

const buttonVariants = cva(
  'inline-flex items-center justify-center rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        primary: 'bg-brand-primary text-white hover:bg-blue-600',
        secondary: 'bg-gray-200 text-gray-900 hover:bg-gray-300',
        danger: 'bg-brand-danger text-white hover:bg-red-600',
        ghost: 'hover:bg-gray-100',
      },
      size: {
        sm: 'h-9 px-3 text-sm',
        md: 'h-10 px-4',
        lg: 'h-11 px-8 text-lg',
      },
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md',
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

export function Button({ variant, size, className, ...props }: ButtonProps) {
  return (
    <button
      className={buttonVariants({ variant, size, className })}
      {...props}
    />
  );
}
```

**Dark Mode 支援：**

```typescript
// apps/web/app/root.tsx
import { useSettingsStore } from '@/stores/settings';

export default function Root() {
  const { theme } = useSettingsStore();

  useEffect(() => {
    const root = document.documentElement;

    if (theme === 'dark') {
      root.classList.add('dark');
    } else if (theme === 'light') {
      root.classList.remove('dark');
    } else {
      // system preference
      const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      root.classList.toggle('dark', isDark);
    }
  }, [theme]);

  return (
    <html lang="zh-TW">
      <body className="bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        <Outlet />
      </body>
    </html>
  );
}
```

**效能優化：**

- JIT 模式僅生成使用的 utility classes
- 生產建置自動 purge 未使用樣式
- 預估 CSS bundle: < 20KB gzipped

**影響範圍：**

- ✅ 快速開發
- ✅ 一致設計
- ✅ Dark mode 支援
- ⚠️ 學習 utility classes（如團隊不熟悉）

**替代方案（已排除）：**

- CSS Modules - 更靈活但一致性較低
- Styled Components - bundle size 較大

---

### Decision 5: Error Handling & Monitoring

**選擇：Sentry (Client-side, Opt-in Consent)**

**理由：**

- 符合 privacy-first 原則（使用者明確同意）
- 強大的錯誤追蹤和 sourcemap 支援
- 免費層足夠 MVP 使用
- 不影響離線可用性

**實作細節：**

**Sentry 初始化：**

```typescript
// apps/web/app/lib/sentry.ts
import * as Sentry from '@sentry/react';
import { useSettingsStore } from '@/stores/settings';

export function initSentry() {
  const { enableTelemetry } = useSettingsStore.getState();

  // 僅在使用者同意時初始化
  if (!enableTelemetry) {
    return;
  }

  Sentry.init({
    dsn: import.meta.env.VITE_SENTRY_DSN,
    environment: import.meta.env.MODE,

    // 隱私保護配置
    beforeSend(event) {
      // 移除敏感資料
      if (event.request?.url) {
        event.request.url = sanitizeUrl(event.request.url);
      }

      // 移除 workspace 路徑
      if (event.extra?.workspacePath) {
        delete event.extra.workspacePath;
      }

      return event;
    },

    // 效能監控（僅追蹤關鍵路徑）
    tracesSampleRate: 0.1,

    // 錯誤過濾
    ignoreErrors: [
      'ResizeObserver loop limit exceeded',
      'Non-Error promise rejection',
    ],
  });
}

function sanitizeUrl(url: string): string {
  try {
    const parsed = new URL(url);
    // 移除 query parameters
    return `${parsed.origin}${parsed.pathname}`;
  } catch {
    return '[sanitized]';
  }
}
```

**同意管理 UI：**

```typescript
// apps/web/app/components/ConsentBanner.tsx
import { useSettingsStore } from '@/stores/settings';
import { initSentry } from '@/lib/sentry';

export function ConsentBanner() {
  const { enableTelemetry, setTelemetry } = useSettingsStore();
  const [isVisible, setIsVisible] = useState(!enableTelemetry);

  const handleAccept = () => {
    setTelemetry(true);
    initSentry();
    setIsVisible(false);
  };

  const handleDecline = () => {
    setTelemetry(false);
    setIsVisible(false);
  };

  if (!isVisible) return null;

  return (
    <div className="fixed bottom-4 right-4 max-w-md p-4 bg-white dark:bg-gray-800 rounded-lg shadow-lg">
      <h3 className="font-semibold mb-2">協助我們改進 MonoGuard</h3>
      <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
        我們使用 Sentry 收集匿名錯誤報告，協助改善產品品質。
        您的程式碼和專案路徑不會被上傳。
      </p>
      <div className="flex gap-2">
        <Button variant="primary" size="sm" onClick={handleAccept}>
          同意
        </Button>
        <Button variant="ghost" size="sm" onClick={handleDecline}>
          拒絕
        </Button>
      </div>
    </div>
  );
}
```

**Error Boundaries：**

```typescript
// apps/web/app/components/ErrorBoundary.tsx
import * as Sentry from '@sentry/react';

export function AppErrorBoundary({ children }: { children: React.ReactNode }) {
  return (
    <Sentry.ErrorBoundary
      fallback={({ error, resetError }) => (
        <div className="flex flex-col items-center justify-center min-h-screen p-4">
          <h1 className="text-2xl font-bold mb-4">發生錯誤</h1>
          <p className="text-gray-600 mb-4">{error.message}</p>
          <Button onClick={resetError}>重新載入</Button>
        </div>
      )}
    >
      {children}
    </Sentry.ErrorBoundary>
  );
}
```

**CLI 錯誤處理：**

```go
// apps/cli/pkg/errors/handler.go
package errors

import (
    "fmt"
    "os"
)

func HandleError(err error, context string) {
    fmt.Fprintf(os.Stderr, "Error in %s: %v\n", context, err)

    // CLI 不發送遙測，僅本地 logging
    logToFile(err, context)

    os.Exit(1)
}

func logToFile(err error, context string) {
    // 寫入 ~/.monoguard/errors.log（本地診斷）
    // ...
}
```

**影響範圍：**

- ✅ 錯誤追蹤與診斷
- ✅ 尊重使用者隱私
- ✅ 離線不受影響
- ⚠️ 需要維護同意管理 UI

**替代方案（已排除）：**

- 無監控 - 難以診斷生產問題
- PostHog - 更偏向產品分析而非錯誤追蹤

---

### Decision 6: Data Persistence & Rendering Strategy

**選擇：Dexie.js (IndexedDB Wrapper) + Hybrid SVG/Canvas Rendering**

**理由：**

- Dexie.js 提供 TypeScript 友善的 IndexedDB API
- Hybrid rendering 根據節點數量自動選擇最佳方案
- 符合隱私約束（完全本地儲存）

**實作細節：**

**Dexie.js Schema：**

```typescript
// packages/types/src/db.ts
import Dexie, { Table } from 'dexie';

export interface AnalysisRecord {
  id?: number;
  timestamp: number;
  projectPath: string;
  result: AnalysisResult;
  metadata: {
    packageCount: number;
    cycleCount: number;
    healthScore: number;
  };
}

export interface UserSettings {
  key: string;
  value: any;
}

export class MonoGuardDB extends Dexie {
  analysisResults!: Table<AnalysisRecord>;
  settings!: Table<UserSettings>;

  constructor() {
    super('MonoGuardDB');

    this.version(1).stores({
      analysisResults: '++id, timestamp, projectPath, metadata.healthScore',
      settings: 'key',
    });
  }
}

export const db = new MonoGuardDB();
```

**分析結果持久化：**

```typescript
// apps/web/app/lib/persistence.ts
import { db } from '@mono-guard/types/db';

export async function saveAnalysisResult(
  projectPath: string,
  result: AnalysisResult
) {
  await db.analysisResults.add({
    timestamp: Date.now(),
    projectPath,
    result,
    metadata: {
      packageCount: result.graph.nodes.length,
      cycleCount: result.cycles.length,
      healthScore: result.healthScore,
    },
  });
}

export async function getRecentAnalyses(limit = 10) {
  return db.analysisResults
    .orderBy('timestamp')
    .reverse()
    .limit(limit)
    .toArray();
}

export async function clearOldAnalyses(daysToKeep = 30) {
  const cutoff = Date.now() - daysToKeep * 24 * 60 * 60 * 1000;
  await db.analysisResults.where('timestamp').below(cutoff).delete();
}
```

**Hybrid Rendering 策略：**

```typescript
// packages/ui-components/src/DependencyGraph/index.tsx
import { SVGRenderer } from './SVGRenderer';
import { CanvasRenderer } from './CanvasRenderer';

interface DependencyGraphProps {
  data: GraphData;
  mode?: 'svg' | 'canvas' | 'auto';
}

export function DependencyGraph({ data, mode = 'auto' }: DependencyGraphProps) {
  const nodeCount = data.nodes.length;

  // 自動選擇渲染模式
  const renderMode = mode === 'auto'
    ? (nodeCount > 500 ? 'canvas' : 'svg')
    : mode;

  return (
    <div className="relative w-full h-full">
      {renderMode === 'svg' ? (
        <SVGRenderer data={data} />
      ) : (
        <CanvasRenderer data={data} />
      )}

      <div className="absolute top-2 right-2 text-xs text-gray-500">
        {nodeCount} nodes • {renderMode.toUpperCase()} mode
      </div>
    </div>
  );
}
```

**SVG Renderer (< 500 nodes)：**

```typescript
// packages/ui-components/src/DependencyGraph/SVGRenderer.tsx
import * as d3 from 'd3';

export function SVGRenderer({ data }: { data: GraphData }) {
  const svgRef = useRef<SVGSVGElement>(null);

  useEffect(() => {
    if (!svgRef.current) return;

    const svg = d3.select(svgRef.current);
    const width = svgRef.current.clientWidth;
    const height = svgRef.current.clientHeight;

    // D3.js force simulation
    const simulation = d3.forceSimulation(data.nodes)
      .force('link', d3.forceLink(data.links).id((d: any) => d.id))
      .force('charge', d3.forceManyBody().strength(-100))
      .force('center', d3.forceCenter(width / 2, height / 2));

    // 繪製 links
    const link = svg.append('g')
      .selectAll('line')
      .data(data.links)
      .join('line')
      .attr('stroke', '#999')
      .attr('stroke-width', 1);

    // 繪製 nodes
    const node = svg.append('g')
      .selectAll('circle')
      .data(data.nodes)
      .join('circle')
      .attr('r', 5)
      .attr('fill', (d) => getNodeColor(d))
      .call(drag(simulation));

    // 更新位置
    simulation.on('tick', () => {
      link
        .attr('x1', (d: any) => d.source.x)
        .attr('y1', (d: any) => d.source.y)
        .attr('x2', (d: any) => d.target.x)
        .attr('y2', (d: any) => d.target.y);

      node
        .attr('cx', (d: any) => d.x)
        .attr('cy', (d: any) => d.y);
    });

    return () => {
      simulation.stop();
    };
  }, [data]);

  return <svg ref={svgRef} className="w-full h-full" />;
}
```

**Canvas Renderer (> 500 nodes)：**

```typescript
// packages/ui-components/src/DependencyGraph/CanvasRenderer.tsx
export function CanvasRenderer({ data }: { data: GraphData }) {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const width = canvas.width = canvas.clientWidth * devicePixelRatio;
    const height = canvas.height = canvas.clientHeight * devicePixelRatio;

    // 使用 d3-force 計算位置（不渲染到 DOM）
    const simulation = d3.forceSimulation(data.nodes)
      .force('link', d3.forceLink(data.links))
      .force('charge', d3.forceManyBody().strength(-50))
      .force('center', d3.forceCenter(width / 2, height / 2))
      .on('tick', render);

    function render() {
      ctx.clearRect(0, 0, width, height);

      // 繪製 links
      ctx.strokeStyle = '#999';
      ctx.lineWidth = 1;
      data.links.forEach((link: any) => {
        ctx.beginPath();
        ctx.moveTo(link.source.x, link.source.y);
        ctx.lineTo(link.target.x, link.target.y);
        ctx.stroke();
      });

      // 繪製 nodes
      data.nodes.forEach((node: any) => {
        ctx.fillStyle = getNodeColor(node);
        ctx.beginPath();
        ctx.arc(node.x, node.y, 5, 0, 2 * Math.PI);
        ctx.fill();
      });
    }

    return () => {
      simulation.stop();
    };
  }, [data]);

  return <canvas ref={canvasRef} className="w-full h-full" />;
}
```

**效能指標：**

- SVG: 流暢互動 (< 500 nodes)，60fps
- Canvas: 高效渲染 (> 500 nodes)，但互動性降低
- IndexedDB: 快速讀寫，支援大量資料

**影響範圍：**

- ✅ 本地優先儲存
- ✅ 彈性渲染策略
- ✅ 符合隱私約束
- ⚠️ Canvas 模式互動性較低（Phase 1 改進）

**替代方案（已排除）：**

- LocalStorage - 容量限制（5-10MB）
- 純 SVG - 大型圖表效能不佳
- 純 Canvas - 小型圖表失去 SVG 互動優勢

---

### Decision 7: CI/CD & Testing Strategy

**選擇：GitHub Actions + Render + Vitest (80%+ Coverage)**

**決策變更 (2026-01-16):** 從 Cloudflare Pages 改為 Render

**理由：**

- 完全免費（符合零成本約束，Render Free Tier）
- 與 Nx monorepo 深度整合
- Render 提供 all-in-one 部署（Web + API + PostgreSQL + Redis）
- `render.yaml` Blueprint 實現 Infrastructure as Code
- Vitest 與 Vite 生態完美整合

**實作細節：**

**GitHub Actions Workflow：**

```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0 # Nx affected 需要完整歷史

      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'

      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'

      - name: Install dependencies
        run: npm ci

      - name: Lint
        run: npx nx affected -t lint

      - name: Test
        run: npx nx affected -t test --coverage

      - name: Build WASM
        run: npx nx build analysis-engine

      - name: Build Web
        run: npx nx build web

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage/lcov.info

  deploy:
    runs-on: ubuntu-latest
    needs: test
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4

      - name: Install dependencies
        run: npm ci

      - name: Build
        run: npx nx build web --prod

      - name: Deploy to Cloudflare Pages
        uses: cloudflare/pages-action@v1
        with:
          apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
          projectName: mono-guard
          directory: dist/apps/web
```

**Vitest 配置：**

```typescript
// apps/web/vitest.config.ts
import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './test/setup.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      exclude: ['node_modules/', 'test/', '**/*.d.ts', '**/*.config.*'],
      thresholds: {
        lines: 80,
        functions: 80,
        branches: 80,
        statements: 80,
      },
    },
  },
});
```

**測試策略：**

1. **Unit Tests (主要重點):**

   ```typescript
   // packages/types/src/__tests__/wasm-adapter.test.ts
   import { describe, it, expect, vi, beforeEach } from 'vitest';
   import { MonoGuardAnalyzer } from '../wasm-adapter';

   describe('MonoGuardAnalyzer', () => {
     let analyzer: MonoGuardAnalyzer;

     beforeEach(() => {
       analyzer = new MonoGuardAnalyzer();

       // Mock WASM functions
       (window as any).analyzeWorkspace = vi
         .fn()
         .mockReturnValue(JSON.stringify({ nodes: [], links: [] }));
     });

     it('should initialize WASM instance', async () => {
       await analyzer.init();
       expect(analyzer['isReady']).toBe(true);
     });

     it('should analyze workspace data', () => {
       const data = { packages: [] };
       const result = analyzer.analyze(data);

       expect(result).toBeDefined();
       expect(result.nodes).toBeInstanceOf(Array);
     });
   });
   ```

2. **Go Tests:**

   ```go
   // packages/analysis-engine/pkg/analyzer/analyzer_test.go
   package analyzer_test

   import (
       "testing"
       "github.com/stretchr/testify/assert"
   )

   func TestAnalyzeWorkspace(t *testing.T) {
       workspace := &Workspace{
           Packages: []Package{
               {Name: "pkg-a", Dependencies: []string{"pkg-b"}},
               {Name: "pkg-b", Dependencies: []string{}},
           },
       }

       result, err := Analyze(workspace)

       assert.NoError(t, err)
       assert.Equal(t, 2, len(result.Graph.Nodes))
       assert.Equal(t, 1, len(result.Graph.Links))
   }

   func TestDetectCircularDependencies(t *testing.T) {
       graph := &DependencyGraph{
           Nodes: []Node{
               {ID: "a"}, {ID: "b"}, {ID: "c"},
           },
           Links: []Link{
               {Source: "a", Target: "b"},
               {Source: "b", Target: "c"},
               {Source: "c", Target: "a"}, // Circular!
           },
       }

       cycles := DetectCycles(graph)

       assert.Equal(t, 1, len(cycles))
       assert.ElementsMatch(t, []string{"a", "b", "c"}, cycles[0].Path)
   }
   ```

3. **Component Tests:**

   ```typescript
   // packages/ui-components/src/__tests__/DependencyGraph.test.tsx
   import { render, screen } from '@testing-library/react';
   import { DependencyGraph } from '../DependencyGraph';

   describe('DependencyGraph', () => {
     it('should render SVG mode for small graphs', () => {
       const data = {
         nodes: Array.from({ length: 10 }, (_, i) => ({ id: `node-${i}` })),
         links: [],
       };

       render(<DependencyGraph data={data} mode="auto" />);

       expect(screen.getByText(/SVG mode/i)).toBeInTheDocument();
     });

     it('should render Canvas mode for large graphs', () => {
       const data = {
         nodes: Array.from({ length: 600 }, (_, i) => ({ id: `node-${i}` })),
         links: [],
       };

       render(<DependencyGraph data={data} mode="auto" />);

       expect(screen.getByText(/CANVAS mode/i)).toBeInTheDocument();
     });
   });
   ```

**Coverage 目標：**

- Overall: > 80%
- Critical paths (analysis engine, cycle detection): > 90%
- UI components: > 70%
- Integration tests: Phase 1 補充

**Cloudflare Pages 配置：**

```toml
# wrangler.toml
name = "mono-guard"
compatibility_date = "2024-01-01"

[site]
bucket = "./dist/apps/web"

[[redirects]]
from = "/*"
to = "/index.html"
status = 200

[headers]
[headers."/*.wasm"]
Cross-Origin-Embedder-Policy = "require-corp"
Cross-Origin-Opener-Policy = "same-origin"
```

**影響範圍：**

- ✅ 自動化測試與部署
- ✅ 零成本基礎設施
- ✅ 高品質保證（80% coverage）
- ⚠️ E2E 測試需要 Phase 1 補充

**替代方案（已排除）：**

- Cloudflare Pages - 原部署平台，但無法統一管理 API + Database
- Vercel - 功能類似但 Render 提供更完整的 all-in-one 方案
- Jest - Vitest 更快且與 Vite 整合更好

---

### Summary: Core Architectural Decisions

| 決策領域             | 選擇                      | 主要理由                 | 權衡考量               |
| -------------------- | ------------------------- | ------------------------ | ---------------------- |
| **Monorepo**         | Nx                        | 專案已使用，經驗豐富     | 需維護配置             |
| **WASM Integration** | Dynamic Loading + Wrapper | 簡化開發，Phase 1 可升級 | 大型專案可能阻塞 UI    |
| **State Management** | Zustand                   | 輕量（< 5KB），API 簡潔  | 大規模狀態需拆分       |
| **Styling**          | Tailwind CSS              | JIT 模式，bundle 小      | 需學習 utility classes |
| **Error Monitoring** | Sentry (Opt-in)           | 強大追蹤，尊重隱私       | 需維護同意 UI          |
| **Data Persistence** | Dexie.js                  | TypeScript 友善          | IndexedDB 學習曲線     |
| **Rendering**        | Hybrid SVG/Canvas         | 彈性，自動切換           | Canvas 互動性較低      |
| **CI/CD**            | GitHub Actions + Render   | All-in-one 部署，自動化  | E2E 測試 Phase 1 補充  |
| **Testing**          | Vitest + Go testing       | 快速，Vite 整合          | 需建立測試文化         |

**這些決策共同形成了 MonoGuard 的技術基礎，確保專案能夠：**

1. ✅ 符合所有 NFR 約束（零成本、隱私優先、離線可用）
2. ✅ 支援 Phase 0 MVP 快速交付
3. ✅ 為 Phase 1/2 擴展預留彈性
4. ✅ 提供優秀的開發者體驗

## Implementation Patterns & Consistency Rules

本節定義強制性的實作模式，確保多個 AI agents 在開發過程中寫出一致、相容的程式碼。這些規則旨在消除潛在的衝突點，提供明確的實作指引。

### 核心原則

**所有 AI Agents 在實作時必須遵守以下原則：**

1. **語言原生慣例優先** - 每種語言使用其社群標準慣例
2. **跨語言邊界統一** - 在語言邊界（JSON, WASM API）使用統一格式
3. **功能模組組織** - 按功能而非類型組織程式碼
4. **明確錯誤處理** - 區分技術錯誤和使用者訊息
5. **型別安全至上** - 利用 TypeScript 和 Go 的型別系統

---

### 1. 命名約定 (Naming Conventions)

#### 1.1 TypeScript 命名

**變數、函式、參數：**

```typescript
// ✅ CORRECT
const analysisResult = await analyzer.analyze(workspaceData);
const healthScore = calculateHealthScore(graph);
function formatDependencyPath(path: string[]): string { ... }

// ❌ INCORRECT
const AnalysisResult = await analyzer.analyze(workspaceData);  // 非常數不用 PascalCase
const health_score = calculateHealthScore(graph);             // 不用 snake_case
function FormatDependencyPath(path: string[]): string { ... } // 函式不用 PascalCase
```

**React 元件：**

```typescript
// ✅ CORRECT
export function AnalysisPanel() { ... }
export function DependencyGraph() { ... }
const ButtonGroup = () => { ... };

// ❌ INCORRECT
export function analysisPanel() { ... }   // 元件必須 PascalCase
export function dependency_graph() { ... } // 不用 snake_case
```

**型別和介面：**

```typescript
// ✅ CORRECT
interface AnalysisResult { ... }
type WorkspaceData = { ... };
enum AnalysisStatus { ... }

// ❌ INCORRECT
interface analysisResult { ... }    // 型別必須 PascalCase
type workspace_data = { ... };       // 不用 snake_case
```

**常數：**

```typescript
// ✅ CORRECT
const MAX_ANALYSIS_TIMEOUT = 30000;
const DEFAULT_HEALTH_THRESHOLD = 75;

// ❌ INCORRECT
const maxAnalysisTimeout = 30000; // 全域常數用 UPPER_SNAKE_CASE
const DefaultHealthThreshold = 75; // 不用 PascalCase
```

#### 1.2 Go 命名

**Exported 識別符（PascalCase）：**

```go
// ✅ CORRECT
type AnalysisResult struct {
    HealthScore int
    CycleCount  int
}

func AnalyzeWorkspace(data WorkspaceData) (*AnalysisResult, error) { ... }

const MaxPackageCount = 10000

// ❌ INCORRECT
type analysisResult struct { ... }           // Exported 型別必須 PascalCase
func analyzeWorkspace(data WorkspaceData) { ... }  // Exported 函式必須 PascalCase
const maxPackageCount = 10000                // Exported 常數必須 PascalCase
```

**Unexported 識別符（camelCase）：**

```go
// ✅ CORRECT
func detectCycles(graph *DependencyGraph) []Cycle { ... }
var defaultTimeout = 5 * time.Second
type cycleDetector struct { ... }

// ❌ INCORRECT
func DetectCycles(graph *DependencyGraph) []Cycle { ... }  // Unexported 不用 PascalCase
var DefaultTimeout = 5 * time.Second                       // Unexported 不用 PascalCase
```

**檔案命名（snake_case）：**

```
// ✅ CORRECT
analyzer.go
dependency_graph.go
workspace_parser.go
fix_suggester_test.go

// ❌ INCORRECT
Analyzer.go           // 不用 PascalCase
dependencyGraph.go    // 不用 camelCase
workspace-parser.go   // 不用 kebab-case
```

#### 1.3 JSON 序列化命名（統一 camelCase）

**TypeScript 介面：**

```typescript
// ✅ CORRECT
interface AnalysisResult {
  healthScore: number;
  cycleCount: number;
  createdAt: string; // ISO 8601
  packageNames: string[];
}

// ❌ INCORRECT
interface AnalysisResult {
  health_score: number; // JSON 不用 snake_case
  CycleCount: number; // JSON 不用 PascalCase
  created_at: string;
}
```

**Go struct tags（強制 camelCase）：**

```go
// ✅ CORRECT
type AnalysisResult struct {
    HealthScore  int      `json:"healthScore"`
    CycleCount   int      `json:"cycleCount"`
    CreatedAt    string   `json:"createdAt"`
    PackageNames []string `json:"packageNames"`
}

// ❌ INCORRECT
type AnalysisResult struct {
    HealthScore int `json:"health_score"`  // 必須用 camelCase
    CycleCount  int `json:"CycleCount"`    // 不能用 PascalCase
    CreatedAt   string                     // 必須明確定義 json tag
}
```

**實際 JSON 範例：**

```json
// ✅ CORRECT
{
  "healthScore": 85,
  "cycleCount": 3,
  "createdAt": "2026-01-12T10:30:00.000Z",
  "packageNames": ["pkg-a", "pkg-b"]
}

// ❌ INCORRECT
{
  "health_score": 85,    // 不用 snake_case
  "CycleCount": 3,       // 不用 PascalCase
  "created_at": "2026-01-12T10:30:00.000Z"
}
```

---

### 2. WASM 橋接模式 (WASM Bridge Patterns)

#### 2.1 統一 Result 型別

**所有 WASM 函式必須返回統一的 Result 結構：**

**Go 端定義：**

```go
// packages/analysis-engine/pkg/common/result.go
package common

type Result struct {
    Data  interface{} `json:"data"`
    Error *ErrorInfo  `json:"error"`
}

type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// 成功返回
func NewSuccess(data interface{}) Result {
    return Result{Data: data, Error: nil}
}

// 錯誤返回
func NewError(code, message string) Result {
    return Result{
        Data:  nil,
        Error: &ErrorInfo{Code: code, Message: message},
    }
}
```

**WASM 函式實作範例：**

```go
// cmd/wasm/main.go
package main

import (
    "encoding/json"
    "syscall/js"
    "github.com/alexyu/mono-guard/pkg/analyzer"
    "github.com/alexyu/mono-guard/pkg/common"
)

func analyzeWorkspace(this js.Value, args []js.Value) interface{} {
    // 解析輸入
    input := args[0].String()
    var workspaceData analyzer.WorkspaceData
    if err := json.Unmarshal([]byte(input), &workspaceData); err != nil {
        result := common.NewError("INVALID_INPUT", err.Error())
        return js.ValueOf(result)
    }

    // 執行分析
    analysisResult, err := analyzer.Analyze(&workspaceData)
    if err != nil {
        result := common.NewError("ANALYSIS_FAILED", err.Error())
        return js.ValueOf(result)
    }

    // 成功返回
    result := common.NewSuccess(analysisResult)
    return js.ValueOf(result)
}

func main() {
    c := make(chan struct{}, 0)

    // 註冊所有 WASM 函式
    js.Global().Set("analyzeWorkspace", js.FuncOf(analyzeWorkspace))
    js.Global().Set("detectCycles", js.FuncOf(detectCycles))
    js.Global().Set("suggestFixes", js.FuncOf(suggestFixes))

    <-c
}
```

**TypeScript 端處理：**

```typescript
// packages/types/src/wasmAdapter/index.ts
export interface WasmResult<T> {
  data: T | null;
  error: {
    code: string;
    message: string;
  } | null;
}

export class MonoGuardAnalyzer {
  private wasmInstance: WebAssembly.Instance | null = null;
  private isReady = false;

  async init(): Promise<void> {
    if (this.isReady) return;

    const go = new Go();
    const response = await fetch('/monoguard.wasm');
    const result = await WebAssembly.instantiateStreaming(
      response,
      go.importObject
    );

    this.wasmInstance = result.instance;
    go.run(this.wasmInstance);
    this.isReady = true;
  }

  private callWasm<T>(funcName: string, data: unknown): T {
    if (!this.isReady) {
      throw new AnalysisError(
        'WASM_NOT_READY',
        'WASM module not initialized',
        'Analysis engine is still loading. Please try again.'
      );
    }

    // 呼叫 WASM 函式
    const resultJson = (window as any)[funcName](JSON.stringify(data));
    const result: WasmResult<T> = JSON.parse(resultJson);

    // 檢查錯誤
    if (result.error) {
      throw new AnalysisError(
        result.error.code,
        result.error.message,
        this.getUserMessage(result.error.code)
      );
    }

    return result.data!;
  }

  analyze(workspaceData: WorkspaceData): AnalysisResult {
    return this.callWasm<AnalysisResult>('analyzeWorkspace', workspaceData);
  }

  detectCycles(graph: DependencyGraph): CircularDependency[] {
    return this.callWasm<CircularDependency[]>('detectCycles', graph);
  }

  private getUserMessage(code: string): string {
    const messages: Record<string, string> = {
      WASM_NOT_READY: '分析引擎尚未就緒，請稍後再試',
      INVALID_INPUT: '無效的專案結構',
      ANALYSIS_FAILED: '分析過程發生錯誤',
      TIMEOUT: '分析超時，請嘗試較小的專案',
    };
    return messages[code] || '發生未知錯誤';
  }
}

export const analyzer = new MonoGuardAnalyzer();
```

#### 2.2 錯誤碼規範

**所有 WASM 錯誤碼必須使用 UPPER_SNAKE_CASE，並在文檔中定義：**

```typescript
// packages/types/src/wasmAdapter/errorCodes.ts
export const WASM_ERROR_CODES = {
  // 初始化錯誤 (1xx)
  WASM_NOT_READY: 'WASM_NOT_READY',
  WASM_INIT_FAILED: 'WASM_INIT_FAILED',

  // 輸入驗證錯誤 (2xx)
  INVALID_INPUT: 'INVALID_INPUT',
  INVALID_WORKSPACE: 'INVALID_WORKSPACE',
  MISSING_PACKAGE_JSON: 'MISSING_PACKAGE_JSON',

  // 分析錯誤 (3xx)
  ANALYSIS_FAILED: 'ANALYSIS_FAILED',
  TIMEOUT: 'TIMEOUT',
  OUT_OF_MEMORY: 'OUT_OF_MEMORY',

  // 循環依賴錯誤 (4xx)
  CYCLE_DETECTION_FAILED: 'CYCLE_DETECTION_FAILED',

  // 修復建議錯誤 (5xx)
  FIX_SUGGESTION_FAILED: 'FIX_SUGGESTION_FAILED',
} as const;

export type WasmErrorCode =
  (typeof WASM_ERROR_CODES)[keyof typeof WASM_ERROR_CODES];
```

---

### 3. 專案組織模式 (Project Organization)

#### 3.1 功能模組組織

**所有 packages 必須按功能模組組織，而非按檔案類型：**

```
packages/types/src/
├── analysis/                    # 分析結果相關型別
│   ├── index.ts                 # Export barrel
│   ├── types.ts                 # 型別定義
│   ├── validators.ts            # 驗證函式
│   └── __tests__/
│       └── index.test.ts
├── wasmAdapter/                 # WASM 橋接層
│   ├── index.ts
│   ├── MonoGuardAnalyzer.ts
│   ├── errorCodes.ts
│   └── __tests__/
│       ├── index.test.ts
│       └── errorHandling.test.ts
├── db/                          # IndexedDB 相關
│   ├── index.ts
│   ├── schema.ts
│   ├── migrations.ts
│   └── __tests__/
│       └── index.test.ts
└── graph/                       # 圖形資料結構
    ├── index.ts
    ├── types.ts
    ├── algorithms.ts
    └── __tests__/
        └── index.test.ts
```

**✅ Barrel Exports 範例：**

```typescript
// packages/types/src/analysis/index.ts
export * from './types';
export * from './validators';

// 使用時
import { AnalysisResult, validateWorkspace } from '@mono-guard/types/analysis';
```

**❌ 反模式（按類型組織）：**

```
packages/types/src/
├── types/           # ❌ 不要按檔案類型分類
│   ├── analysis.ts
│   ├── graph.ts
│   └── db.ts
├── validators/
│   └── ...
└── __tests__/       # ❌ 測試應該與模組共置
    └── ...
```

#### 3.2 測試檔案位置

**測試檔案必須放在功能模組的 `__tests__/` 子目錄中：**

```typescript
// ✅ CORRECT
packages / types / src / analysis / __tests__ / index.test.ts;

// ❌ INCORRECT
packages / types / src / analysis / index.test.ts; // 不 co-locate
packages / types / test / analysis.test.ts; // 不放在根 test/
packages / types / src / __tests__ / analysis.test.ts; // 不放在 src/__tests__/
```

**測試檔案命名：**

```typescript
// ✅ CORRECT - 與模組對應
__tests__ / index.test.ts; // 測試 index.ts
__tests__ / errorHandling.test.ts; // 測試 errorHandling.ts
__tests__ / integration.test.ts; // 整合測試

// ❌ INCORRECT
__tests__ / test.ts; // 不明確
__tests__ / index.spec.ts; // 不用 .spec（統一 .test）
__tests__ / indexTest.ts; // 缺少 .test 副檔名
```

#### 3.3 Go 專案組織

**Go 專案必須遵循標準 Go 專案結構：**

```
packages/analysis-engine/
├── cmd/
│   ├── wasm/                    # WASM 建置目標
│   │   └── main.go
│   └── cli/                     # CLI 工具（如果需要）
│       └── main.go
├── pkg/                         # 可匯出的 packages
│   ├── analyzer/
│   │   ├── analyzer.go
│   │   ├── analyzer_test.go
│   │   └── dependency_graph.go
│   ├── parser/
│   │   ├── workspace_parser.go
│   │   ├── workspace_parser_test.go
│   │   └── package_json.go
│   ├── rules/
│   │   ├── fix_suggester.go
│   │   └── fix_suggester_test.go
│   └── common/
│       ├── result.go
│       └── errors.go
├── internal/                    # 私有 packages
│   └── ...
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

**Go 測試命名規範：**

```go
// ✅ CORRECT
analyzer_test.go         // 測試 analyzer.go
dependency_graph_test.go // 測試 dependency_graph.go

// ❌ INCORRECT
analyzer.test.go         // 不用 . 分隔
test_analyzer.go         // 不用 test_ 前綴
analyzerTest.go          // 不用 camelCase
```

---

### 4. 狀態管理模式 (State Management with Zustand)

#### 4.1 Action 命名規則

**重要業務操作使用完整描述，簡單操作使用簡短動詞：**

```typescript
// ✅ CORRECT - 混合模式
export const useAnalysisStore = create<AnalysisState>((set, get) => ({
  // 重要業務操作 - 完整描述
  startAnalysis: async (data: WorkspaceData) => { ... },
  retryAnalysis: async () => { ... },
  exportAnalysisResult: (format: 'json' | 'html') => { ... },

  // 簡單操作 - 簡短動詞
  clear: () => { ... },
  reset: () => { ... },

  // 中等複雜度 - 動詞 + 名詞
  selectNode: (nodeId: string) => { ... },
  updateFilters: (filters: Partial<FilterState>) => { ... },
  toggleFilter: (key: keyof FilterState) => { ... },
}));

// ❌ INCORRECT - 不一致
export const useAnalysisStore = create<AnalysisState>((set) => ({
  start: async (data) => { ... },              // ❌ 重要操作太簡短
  clearAnalysisResult: () => { ... },          // ❌ 簡單操作太冗長
  nodeSelect: (id) => { ... },                 // ❌ 名詞在前
  filter_update: (filters) => { ... },         // ❌ snake_case
}));
```

**判斷準則：**

- **完整描述**：涉及非同步操作、複雜邏輯、外部依賴
- **簡短動詞**：單純的狀態設置（clear, reset, toggle）
- **動詞 + 名詞**：中等複雜度的狀態更新

#### 4.2 Store 組織結構

**Store 必須按功能領域分離，避免巨型 store：**

```typescript
// ✅ CORRECT - 按功能分離
// apps/web/app/stores/analysis.ts
export const useAnalysisStore = create<AnalysisState>(() => ({ ... }));

// apps/web/app/stores/settings.ts
export const useSettingsStore = create<SettingsState>(() => ({ ... }));

// apps/web/app/stores/ui.ts
export const useUIStore = create<UIState>(() => ({ ... }));

// ❌ INCORRECT - 單一巨型 store
// apps/web/app/stores/index.ts
export const useAppStore = create<AppState>(() => ({
  // 分析狀態
  analysisResult: null,
  isAnalyzing: false,
  // 設定狀態
  theme: 'light',
  language: 'zh-TW',
  // UI 狀態
  isSidebarOpen: true,
  activeTab: 'overview',
  // ... 所有狀態混在一起
}));
```

#### 4.3 Middleware 使用順序

**Middleware 必須按以下順序包裝：**

```typescript
// ✅ CORRECT - 固定順序
export const useAnalysisStore = create<AnalysisState>()(
  devtools(           // 1. DevTools（最外層，開發時使用）
    persist(          // 2. Persist（持久化）
      immer(          // 3. Immer（簡化狀態更新）
        (set, get) => ({
          // Store 實作
        })
      ),
      {
        name: 'monoguard-analysis',
        partialize: (state) => ({
          filters: state.filters,  // 僅持久化需要的欄位
        }),
      }
    )
  )
);

// ❌ INCORRECT - 順序錯誤或缺少配置
export const useAnalysisStore = create<AnalysisState>()(
  persist(           // ❌ persist 不應該在最外層
    devtools(
      (set) => ({ ... })
    ),
    { name: 'analysis' }  // ❌ 缺少 partialize
  )
);
```

---

### 5. 錯誤處理模式 (Error Handling)

#### 5.1 分層錯誤處理

**所有錯誤必須區分技術錯誤和使用者訊息：**

```typescript
// packages/types/src/errors/AnalysisError.ts
export class AnalysisError extends Error {
  code: string;
  technicalMessage: string;
  userMessage: string;
  context?: Record<string, unknown>;

  constructor(
    code: string,
    technicalMessage: string,
    userMessage: string,
    context?: Record<string, unknown>
  ) {
    super(technicalMessage);
    this.name = 'AnalysisError';
    this.code = code;
    this.technicalMessage = technicalMessage;
    this.userMessage = userMessage;
    this.context = context;

    // 維持正確的 stack trace
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, AnalysisError);
    }
  }
}
```

**使用範例：**

```typescript
// ✅ CORRECT - 拋出 AnalysisError
async function analyzeWorkspace(data: WorkspaceData): Promise<AnalysisResult> {
  try {
    const result = await analyzer.analyze(data);
    return result;
  } catch (error) {
    throw new AnalysisError(
      'ANALYSIS_FAILED',
      `Analysis engine error: ${error.message}`,
      '分析過程發生錯誤，請檢查專案結構是否正確',
      { workspacePackages: data.packages.length }
    );
  }
}

// 捕獲並處理
try {
  const result = await analyzeWorkspace(data);
} catch (error) {
  if (error instanceof AnalysisError) {
    // 記錄技術錯誤到 Sentry
    Sentry.captureException(error, {
      tags: { errorCode: error.code },
      contexts: { analysis: error.context },
    });

    // 顯示使用者友善訊息
    toast.error(error.userMessage);
  } else {
    // 未預期的錯誤
    Sentry.captureException(error);
    toast.error('發生未知錯誤，請重試');
  }
}

// ❌ INCORRECT - 直接拋出原始錯誤
async function analyzeWorkspace(data: WorkspaceData): Promise<AnalysisResult> {
  try {
    return await analyzer.analyze(data);
  } catch (error) {
    throw error;  // ❌ 沒有包裝成 AnalysisError
  }
}

// ❌ INCORRECT - 顯示技術錯誤給使用者
catch (error) {
  toast.error(error.message);  // ❌ 技術訊息不友善
}
```

#### 5.2 錯誤邊界 (Error Boundaries)

**React 元件必須使用 Sentry Error Boundary：**

```typescript
// ✅ CORRECT
// apps/web/app/root.tsx
import * as Sentry from '@sentry/react';
import { Outlet } from '@tanstack/react-router';

export default function Root() {
  return (
    <Sentry.ErrorBoundary
      fallback={({ error, resetError }) => (
        <div className="flex flex-col items-center justify-center min-h-screen p-4">
          <h1 className="text-2xl font-bold mb-4">發生錯誤</h1>
          <p className="text-gray-600 mb-4">
            {error instanceof AnalysisError ? error.userMessage : '發生未知錯誤'}
          </p>
          <Button onClick={resetError}>重新載入</Button>
        </div>
      )}
      showDialog={false}
    >
      <Outlet />
    </Sentry.ErrorBoundary>
  );
}

// ❌ INCORRECT - 沒有 Error Boundary
export default function Root() {
  return <Outlet />;  // ❌ 錯誤會導致白屏
}
```

#### 5.3 Go 錯誤處理

**Go 錯誤必須提供上下文資訊：**

```go
// ✅ CORRECT - 使用 fmt.Errorf 包裝錯誤
func AnalyzeWorkspace(data *WorkspaceData) (*AnalysisResult, error) {
    graph, err := buildDependencyGraph(data)
    if err != nil {
        return nil, fmt.Errorf("failed to build dependency graph: %w", err)
    }

    cycles, err := detectCycles(graph)
    if err != nil {
        return nil, fmt.Errorf("failed to detect cycles: %w", err)
    }

    return &AnalysisResult{
        Graph:  graph,
        Cycles: cycles,
    }, nil
}

// ❌ INCORRECT - 直接返回錯誤
func AnalyzeWorkspace(data *WorkspaceData) (*AnalysisResult, error) {
    graph, err := buildDependencyGraph(data)
    if err != nil {
        return nil, err  // ❌ 缺少上下文
    }
    // ...
}

// ❌ INCORRECT - 吞沒錯誤
func AnalyzeWorkspace(data *WorkspaceData) (*AnalysisResult, error) {
    graph, err := buildDependencyGraph(data)
    if err != nil {
        log.Println(err)  // ❌ 僅記錄，沒有返回
        return nil, nil
    }
    // ...
}
```

---

### 6. 資料格式約定 (Data Format Conventions)

#### 6.1 日期時間格式

**所有日期時間必須使用 ISO 8601 字串：**

```typescript
// ✅ CORRECT
interface AnalysisRecord {
  timestamp: string; // "2026-01-12T10:30:00.000Z"
  createdAt: string; // "2026-01-12T10:30:00.000Z"
  updatedAt: string | null; // "2026-01-12T10:30:00.000Z" or null
}

// 建立日期
const record = {
  timestamp: new Date().toISOString(),
  createdAt: new Date().toISOString(),
  updatedAt: null,
};

// 解析日期
const date = new Date(record.timestamp);

// ❌ INCORRECT - Unix timestamp
interface AnalysisRecord {
  timestamp: number; // ❌ 1736679000000
}

// ❌ INCORRECT - 自訂格式
interface AnalysisRecord {
  timestamp: string; // ❌ "2026/01/12 10:30:00"
}
```

**Go 端處理：**

```go
// ✅ CORRECT
type AnalysisRecord struct {
    Timestamp string `json:"timestamp"`  // 使用 string，不是 time.Time
    CreatedAt string `json:"createdAt"`
}

// 建立
record := AnalysisRecord{
    Timestamp: time.Now().UTC().Format(time.RFC3339),
    CreatedAt: time.Now().UTC().Format(time.RFC3339),
}

// ❌ INCORRECT - 使用 time.Time（序列化格式不一致）
type AnalysisRecord struct {
    Timestamp time.Time `json:"timestamp"`  // ❌ 序列化格式可能不同
}
```

#### 6.2 空值處理

**TypeScript 和 Go 的空值處理約定：**

```typescript
// ✅ CORRECT - 明確區分 null 和 undefined
interface AnalysisResult {
  healthScore: number; // 必填，不會是 null
  selectedNode: string | null; // 可選，明確用 null 表示「無選擇」
  filters?: FilterState; // 可選，undefined 表示「未設置」
}

// 使用時
if (result.selectedNode === null) {
  // 明確知道沒有選擇節點
}

if (result.filters === undefined) {
  // 明確知道尚未設置篩選器
}

// ❌ INCORRECT - null 和 undefined 混用
interface AnalysisResult {
  selectedNode: string | null | undefined; // ❌ 語意不清
}
```

**Go 端處理：**

```go
// ✅ CORRECT - 使用指標表示可空
type AnalysisResult struct {
    HealthScore  int     `json:"healthScore"`
    SelectedNode *string `json:"selectedNode"`  // null 如果未選擇
}

// 設置 null
result := AnalysisResult{
    HealthScore:  85,
    SelectedNode: nil,  // 序列化為 JSON null
}

// ❌ INCORRECT - 使用空字串表示空值
type AnalysisResult struct {
    SelectedNode string `json:"selectedNode"`  // ❌ "" 和 真實空字串無法區分
}
```

---

### 7. 匯入路徑約定 (Import Path Conventions)

#### 7.1 TypeScript 匯入規則

**跨 package 使用 Nx workspace 路徑，同 package 內使用 @ alias：**

```typescript
// ✅ CORRECT - 跨 package（從 apps/web 匯入 packages）
import { AnalysisResult } from '@mono-guard/types/analysis';
import { Button } from '@mono-guard/ui-components';
import { db } from '@mono-guard/types/db';

// ✅ CORRECT - 同 package 內（使用 @ alias）
// 在 apps/web/app/components/AnalysisPanel.tsx
import { useAnalysisStore } from '@/stores/analysis';
import { formatDate } from '@/utils/formatDate';
import { DependencyGraph } from '@/components/DependencyGraph';

// ❌ INCORRECT - 混用相對路徑
import { useAnalysisStore } from '../../stores/analysis'; // ❌ 應該用 @/stores
import { formatDate } from '../../../utils/formatDate'; // ❌ 應該用 @/utils

// ❌ INCORRECT - 錯誤使用 workspace 路徑
import { useAnalysisStore } from '@mono-guard/web/stores/analysis'; // ❌ 不匯出 internal
```

**tsconfig.json 配置：**

```json
// apps/web/tsconfig.json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./app/*"] // @ alias 指向 app/
    }
  }
}
```

#### 7.2 Barrel Export 規範

**每個功能模組必須提供 index.ts barrel export：**

```typescript
// ✅ CORRECT
// packages/types/src/analysis/index.ts
export * from './types';
export * from './validators';
export { calculateHealthScore } from './utils';

// 使用時可以簡潔匯入
import {
  AnalysisResult,
  validateWorkspace,
  calculateHealthScore,
} from '@mono-guard/types/analysis';

// ❌ INCORRECT - 沒有 barrel export
// 使用時需要深層匯入
import { AnalysisResult } from '@mono-guard/types/analysis/types';
import { validateWorkspace } from '@mono-guard/types/analysis/validators';
import { calculateHealthScore } from '@mono-guard/types/analysis/utils';
```

---

### 8. 檔案命名約定 (File Naming Conventions)

#### 8.1 TypeScript 檔案命名

```
// ✅ CORRECT
components/
├── AnalysisPanel/           # React 元件用 PascalCase
│   ├── index.tsx
│   └── AnalysisPanel.module.css
├── DependencyGraph/
│   └── index.tsx
utils/
├── formatDate.ts            # 工具函式用 camelCase
├── calculateHealth.ts
└── parseWorkspace.ts
hooks/
├── useAnalysis.ts           # Hooks 用 camelCase
└── useSettings.ts
types/
├── analysis.ts              # 型別定義用 camelCase
└── graph.ts

// ❌ INCORRECT
components/
├── analysis-panel/          # ❌ 元件不用 kebab-case
├── Analysis_Panel/          # ❌ 元件不用 snake_case
utils/
├── FormatDate.ts            # ❌ 工具不用 PascalCase
└── format-date.ts           # ❌ 工具不用 kebab-case
```

#### 8.2 Go 檔案命名

```
// ✅ CORRECT
pkg/analyzer/
├── analyzer.go              # snake_case
├── analyzer_test.go
├── dependency_graph.go
└── health_calculator.go

// ❌ INCORRECT
pkg/analyzer/
├── Analyzer.go              # ❌ 不用 PascalCase
├── analyzer.test.go         # ❌ test 前要加 _
├── dependencyGraph.go       # ❌ 不用 camelCase
└── dependency-graph.go      # ❌ 不用 kebab-case
```

---

### 9. 強制執行準則 (Enforcement Guidelines)

#### 9.1 所有 AI Agents 必須遵守

**在實作任何功能時，AI Agents 必須：**

1. ✅ **命名檢查** - 使用本文檔定義的命名約定
2. ✅ **結構檢查** - 遵循功能模組組織方式
3. ✅ **錯誤處理** - 使用 AnalysisError 分層錯誤
4. ✅ **WASM 橋接** - 使用統一 Result 型別
5. ✅ **日期格式** - 使用 ISO 8601 字串
6. ✅ **匯入路徑** - 跨 package 用 workspace，內部用 @ alias
7. ✅ **測試位置** - 放在 `__tests__/` 子目錄
8. ✅ **Barrel Exports** - 提供 index.ts 匯出
9. ✅ **型別安全** - 完整的 TypeScript 型別定義
10. ✅ **文件註解** - 重要函式和型別必須有 JSDoc/Go doc

#### 9.2 Code Review 檢查清單

**提交前必須檢查：**

- [ ] 所有命名符合約定（camelCase vs PascalCase vs snake_case）
- [ ] JSON 欄位統一使用 camelCase
- [ ] WASM 函式返回統一 Result 型別
- [ ] 錯誤使用 AnalysisError 包裝，提供使用者訊息
- [ ] 日期使用 ISO 8601 格式
- [ ] 測試檔案放在 `__tests__/` 目錄
- [ ] 提供 barrel exports (index.ts)
- [ ] 匯入路徑正確（workspace vs @ alias）
- [ ] TypeScript 無型別錯誤
- [ ] Go 程式碼通過 `go vet` 和 `golint`

#### 9.3 自動化檢查

**建議配置以下工具自動檢查：**

```json
// .eslintrc.json
{
  "rules": {
    "camelcase": ["error", { "properties": "always" }],
    "@typescript-eslint/naming-convention": [
      "error",
      {
        "selector": "interface",
        "format": ["PascalCase"]
      },
      {
        "selector": "typeAlias",
        "format": ["PascalCase"]
      },
      {
        "selector": "function",
        "format": ["camelCase", "PascalCase"]
      },
      {
        "selector": "variable",
        "format": ["camelCase", "UPPER_CASE"]
      }
    ]
  }
}
```

```yaml
# .github/workflows/ci.yml
- name: Lint TypeScript
  run: npx nx affected -t lint

- name: Lint Go
  run: |
    go vet ./...
    golangci-lint run
```

---

### 10. 模式範例 (Pattern Examples)

#### 10.1 完整功能實作範例

**建立新的分析功能模組（從 TypeScript 到 Go）：**

**步驟 1: 定義 TypeScript 型別**

```typescript
// packages/types/src/complexity/index.ts
export interface ComplexityMetrics {
  cyclomaticComplexity: number;
  cognitiveComplexity: number;
  maintainabilityIndex: number;
  createdAt: string; // ISO 8601
}

export interface ComplexityAnalysisResult {
  packageName: string;
  metrics: ComplexityMetrics;
  issues: ComplexityIssue[];
}

export interface ComplexityIssue {
  file: string;
  line: number;
  severity: 'low' | 'medium' | 'high';
  message: string;
}
```

**步驟 2: 實作 Go 分析引擎**

```go
// packages/analysis-engine/pkg/complexity/analyzer.go
package complexity

type ComplexityMetrics struct {
    CyclomaticComplexity int `json:"cyclomaticComplexity"`
    CognitiveComplexity  int `json:"cognitiveComplexity"`
    MaintainabilityIndex int `json:"maintainabilityIndex"`
    CreatedAt            string `json:"createdAt"`
}

type ComplexityAnalysisResult struct {
    PackageName string             `json:"packageName"`
    Metrics     ComplexityMetrics  `json:"metrics"`
    Issues      []ComplexityIssue  `json:"issues"`
}

type ComplexityIssue struct {
    File     string `json:"file"`
    Line     int    `json:"line"`
    Severity string `json:"severity"`
    Message  string `json:"message"`
}

func AnalyzeComplexity(packagePath string) (*ComplexityAnalysisResult, error) {
    // 實作...
    metrics := &ComplexityMetrics{
        CyclomaticComplexity: 15,
        CognitiveComplexity:  20,
        MaintainabilityIndex: 75,
        CreatedAt:            time.Now().UTC().Format(time.RFC3339),
    }

    return &ComplexityAnalysisResult{
        PackageName: packagePath,
        Metrics:     *metrics,
        Issues:      []ComplexityIssue{},
    }, nil
}
```

**步驟 3: 建立 WASM 橋接**

```go
// cmd/wasm/main.go
func analyzeComplexity(this js.Value, args []js.Value) interface{} {
    packagePath := args[0].String()

    result, err := complexity.AnalyzeComplexity(packagePath)
    if err != nil {
        return js.ValueOf(common.NewError("COMPLEXITY_ANALYSIS_FAILED", err.Error()))
    }

    return js.ValueOf(common.NewSuccess(result))
}

func main() {
    // ...
    js.Global().Set("analyzeComplexity", js.FuncOf(analyzeComplexity))
    // ...
}
```

**步驟 4: 建立 TypeScript Adapter**

```typescript
// packages/types/src/wasmAdapter/index.ts
export class MonoGuardAnalyzer {
  // ...

  analyzeComplexity(packagePath: string): ComplexityAnalysisResult {
    return this.callWasm<ComplexityAnalysisResult>(
      'analyzeComplexity',
      packagePath
    );
  }
}
```

**步驟 5: 建立 Zustand Store**

```typescript
// apps/web/app/stores/complexity.ts
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { analyzer } from '@mono-guard/types/wasmAdapter';
import { ComplexityAnalysisResult } from '@mono-guard/types/complexity';
import { AnalysisError } from '@mono-guard/types/errors';

interface ComplexityState {
  result: ComplexityAnalysisResult | null;
  isAnalyzing: boolean;
  error: string | null;

  analyzeComplexity: (packagePath: string) => Promise<void>;
  clear: () => void;
}

export const useComplexityStore = create<ComplexityState>()(
  devtools((set) => ({
    result: null,
    isAnalyzing: false,
    error: null,

    analyzeComplexity: async (packagePath) => {
      set({ isAnalyzing: true, error: null });

      try {
        const result = await analyzer.analyzeComplexity(packagePath);
        set({ result, isAnalyzing: false });
      } catch (error) {
        const errorMessage =
          error instanceof AnalysisError ? error.userMessage : '複雜度分析失敗';

        set({ error: errorMessage, isAnalyzing: false });

        if (error instanceof AnalysisError) {
          Sentry.captureException(error);
        }
      }
    },

    clear: () => {
      set({ result: null, error: null });
    },
  }))
);
```

**步驟 6: 建立 React 元件**

```typescript
// apps/web/app/components/ComplexityPanel/index.tsx
import { useComplexityStore } from '@/stores/complexity';
import { Button } from '@mono-guard/ui-components';

export function ComplexityPanel() {
  const { result, isAnalyzing, error, analyzeComplexity, clear } = useComplexityStore();

  const handleAnalyze = () => {
    analyzeComplexity('/path/to/package');
  };

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Complexity Analysis</h2>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {isAnalyzing && <LoadingSpinner />}

      {result && (
        <div>
          <h3 className="font-semibold">{result.packageName}</h3>
          <div className="grid grid-cols-3 gap-4 mt-4">
            <MetricCard
              label="Cyclomatic"
              value={result.metrics.cyclomaticComplexity}
            />
            <MetricCard
              label="Cognitive"
              value={result.metrics.cognitiveComplexity}
            />
            <MetricCard
              label="Maintainability"
              value={result.metrics.maintainabilityIndex}
            />
          </div>
        </div>
      )}

      <div className="mt-4 flex gap-2">
        <Button variant="primary" onClick={handleAnalyze}>
          Analyze
        </Button>
        <Button variant="ghost" onClick={clear}>
          Clear
        </Button>
      </div>
    </div>
  );
}
```

**步驟 7: 建立測試**

```typescript
// packages/types/src/complexity/__tests__/index.test.ts
import { describe, it, expect } from 'vitest';
import type { ComplexityMetrics } from '../index';

describe('ComplexityMetrics', () => {
  it('should have correct structure', () => {
    const metrics: ComplexityMetrics = {
      cyclomaticComplexity: 15,
      cognitiveComplexity: 20,
      maintainabilityIndex: 75,
      createdAt: new Date().toISOString(),
    };

    expect(metrics.cyclomaticComplexity).toBeTypeOf('number');
    expect(metrics.createdAt).toMatch(/^\d{4}-\d{2}-\d{2}T/);
  });
});
```

```go
// packages/analysis-engine/pkg/complexity/analyzer_test.go
package complexity_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/alexyu/mono-guard/pkg/complexity"
)

func TestAnalyzeComplexity(t *testing.T) {
    result, err := complexity.AnalyzeComplexity("test-package")

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test-package", result.PackageName)
    assert.Greater(t, result.Metrics.CyclomaticComplexity, 0)
}
```

#### 10.2 反模式範例

**❌ 常見錯誤：**

```typescript
// ❌ 錯誤 1: 命名不一致
interface analysis_result {  // ❌ 應該用 PascalCase
  health_score: number;      // ❌ JSON 應該用 camelCase
}

// ❌ 錯誤 2: 沒有錯誤包裝
catch (error) {
  toast.error(error.message);  // ❌ 直接顯示技術訊息
}

// ❌ 錯誤 3: 不使用統一 Result
function analyzeWorkspace(data: string): any {  // ❌ any 型別
  return JSON.parse((window as any).analyze(data));  // ❌ 未檢查錯誤
}

// ❌ 錯誤 4: 錯誤的匯入路徑
import { Button } from '../../../../../../packages/ui-components/src/Button';  // ❌ 相對路徑太深

// ❌ 錯誤 5: 測試位置錯誤
packages/types/src/analysis.test.ts  // ❌ 應該在 __tests__/

// ❌ 錯誤 6: 日期格式錯誤
interface Record {
  timestamp: number;  // ❌ 應該用 ISO 8601 string
}

// ❌ 錯誤 7: Action 命名不清
const useStore = create((set) => ({
  do: () => { ... },          // ❌ 太模糊
  handleClick: () => { ... }, // ❌ 不是 UI handler
}));
```

---

### 11. 模式更新流程 (Pattern Update Process)

#### 11.1 發現新衝突點

當發現新的 AI agent 衝突點時：

1. **記錄衝突** - 在 GitHub Issue 中描述衝突情境
2. **提出模式** - 討論並確定統一模式
3. **更新文檔** - 將新模式加入此文檔
4. **通知團隊** - 在 PR 中說明模式變更
5. **重構現有程式碼** - 統一修改違反新模式的程式碼

#### 11.2 模式例外處理

某些特殊情況可能需要例外處理：

1. **記錄原因** - 在程式碼註解中說明為何不遵循模式
2. **獲得批准** - 在 PR review 中討論並批准例外
3. **限制範圍** - 確保例外不擴散到其他模組

```typescript
// ✅ CORRECT - 有明確註解的例外
// Exception: 使用 snake_case 因為需要與舊版 API 相容
interface LegacyAPIResponse {
  user_id: string;
  created_at: string;
}
```

---

### 12. 模式總結表

| 類別           | 規則                                  | TypeScript | Go               | JSON |
| -------------- | ------------------------------------- | ---------- | ---------------- | ---- |
| **變數/函式**  | camelCase                             | ✅         | ✅ (unexported)  | ✅   |
| **型別/介面**  | PascalCase                            | ✅         | ✅ (exported)    | -    |
| **React 元件** | PascalCase                            | ✅         | -                | -    |
| **常數**       | UPPER_SNAKE_CASE                      | ✅         | -                | -    |
| **檔案名稱**   | PascalCase (元件)<br>camelCase (其他) | ✅         | -                | -    |
| **Go 檔案**    | snake_case                            | -          | ✅               | -    |
| **JSON 欄位**  | camelCase                             | ✅         | ✅ (struct tags) | ✅   |
| **日期格式**   | ISO 8601 string                       | ✅         | ✅               | ✅   |
| **錯誤處理**   | AnalysisError                         | ✅         | -                | -    |
| **WASM 返回**  | Result 型別                           | ✅         | ✅               | ✅   |
| **測試位置**   | **tests**/                            | ✅         | -                | -    |
| **Go 測試**    | \_test.go suffix                      | -          | ✅               | -    |
| **匯入路徑**   | workspace + @                         | ✅         | -                | -    |

---

**這些實作模式確保：**

1. ✅ 多個 AI agents 寫出一致的程式碼
2. ✅ 跨語言邊界（TypeScript ↔ Go）無縫整合
3. ✅ 錯誤處理統一且使用者友善
4. ✅ 專案結構清晰易維護
5. ✅ 型別安全貫穿整個專案

---

## Project Structure & Boundaries

### 現有架構與 PRD 需求的差異分析

#### 重大發現：架構不匹配

MonoGuard 專案目前的實作與 PRD 中定義的需求存在**重大差異**：

**現有架構（Client-Server）：**

```
apps/
├── api/              # ❌ Go 後端 API (Gin + PostgreSQL + Redis)
├── frontend/         # ❌ Next.js 15.2 (SSR + SSG)
├── cli/              # ❌ TypeScript CLI
└── frontend-e2e/

libs/
└── shared-types/     # ♻️ 可重構為 packages/types/
```

**PRD 需求（Client-Only）：**

- ✅ **NFR9**: 完全零後端架構（Zero backend infrastructure）
- ✅ **NFR10**: 所有處理在瀏覽器內執行（Go WASM）
- ✅ **FR36**: 本地儲存使用 IndexedDB（不使用資料庫）
- ✅ **FR15-FR27**: TanStack Start 靜態網站
- ✅ **FR28-FR33**: Go 原生 CLI（不是 TypeScript）

#### 架構決策：完整重構為 Client-Only

基於與使用者的討論，決定採用 **選項 B：完整重構**，將架構對齊 PRD 的零後端需求：

**要刪除的檔案（❌）：**

- `apps/api/` - 整個 Go 後端（Gin + PostgreSQL + Redis）
- `apps/frontend/` - Next.js 應用（需改為 TanStack Start）
- `apps/cli/` - TypeScript CLI（需改寫為 Go）

**要重構的檔案（♻️）：**

- `libs/shared-types/` → `packages/types/` - 重構為 Nx packages 結構

**要創建的檔案（✨）：**

- `packages/analysis-engine/` - Go WASM 分析引擎
- `apps/web/` - TanStack Start 靜態網站
- `apps/cli/` (新版) - Go 原生 CLI
- `packages/ui-components/` - 共享 React 元件庫

---

### 完整目標專案結構

```
mono-guard/
├── README.md
├── package.json                      # Nx workspace root
├── nx.json
├── pnpm-workspace.yaml
├── tsconfig.base.json
├── .gitignore
├── .github/
│   └── workflows/
│       ├── ci.yml                    # ✨ GitHub Actions CI
│       └── deploy.yml                # ✨ Cloudflare Pages 部署
│
├── apps/
│   ├── web/                          # ✨ 新建：TanStack Start 靜態網站
│   │   ├── app/
│   │   │   ├── routes/
│   │   │   │   ├── __root.tsx       # Root layout
│   │   │   │   ├── index.tsx        # 首頁 (上傳介面)
│   │   │   │   ├── analysis.$id.tsx # 分析結果頁
│   │   │   │   └── history.tsx      # 歷史記錄頁
│   │   │   ├── components/
│   │   │   │   ├── UploadZone/      # 拖放上傳區
│   │   │   │   ├── DependencyGraph/ # D3.js 視覺化
│   │   │   │   ├── CircularDepList/ # 循環依賴列表
│   │   │   │   ├── FixSuggestions/  # 修復建議卡片
│   │   │   │   └── HealthScore/     # 健康分數儀表板
│   │   │   ├── stores/
│   │   │   │   ├── analysis.ts      # Zustand - 分析狀態
│   │   │   │   ├── history.ts       # Zustand - 歷史記錄
│   │   │   │   └── settings.ts      # Zustand - 使用者設定
│   │   │   ├── lib/
│   │   │   │   ├── persistence.ts   # IndexedDB 封裝 (Dexie.js)
│   │   │   │   └── wasmLoader.ts    # WASM 動態載入邏輯
│   │   │   └── styles/
│   │   │       └── globals.css      # Tailwind CSS
│   │   ├── public/
│   │   │   ├── monoguard.wasm       # ✨ 編譯後的 Go WASM
│   │   │   └── wasm_exec.js         # Go WASM 執行器
│   │   ├── vite.config.ts
│   │   ├── tailwind.config.ts
│   │   ├── tsconfig.json
│   │   └── package.json
│   │
│   └── cli/                          # ✨ 重寫：Go 原生 CLI
│       ├── cmd/
│       │   └── monoguard/
│       │       └── main.go           # CLI 入口點
│       ├── pkg/
│       │   ├── commands/
│       │   │   ├── analyze.go        # analyze 指令
│       │   │   ├── check.go          # check 指令
│       │   │   ├── fix.go            # fix 指令
│       │   │   └── export.go         # export 指令
│       │   ├── output/
│       │   │   ├── formatter.go      # 輸出格式化 (JSON/YAML/Text)
│       │   │   └── reporter.go       # 報告生成器
│       │   └── config/
│       │       └── loader.go         # 設定檔載入
│       ├── .monoguard.example.yaml   # 範例設定檔
│       ├── go.mod
│       ├── go.sum
│       └── README.md
│
├── packages/
│   ├── analysis-engine/              # ✨ 新建：Go WASM 核心分析引擎
│   │   ├── cmd/
│   │   │   └── wasm/
│   │   │       └── main.go           # WASM 編譯入口點
│   │   ├── pkg/
│   │   │   ├── analyzer/
│   │   │   │   ├── workspace.go      # Monorepo 工作區分析
│   │   │   │   ├── dependency.go     # 依賴圖建構
│   │   │   │   ├── circular.go       # 循環依賴檢測
│   │   │   │   └── health.go         # 健康分數計算
│   │   │   ├── rules/
│   │   │   │   ├── root_cause.go     # 根因分析
│   │   │   │   ├── impact.go         # 影響評估
│   │   │   │   └── strategies.go     # 修復策略
│   │   │   ├── parsers/
│   │   │   │   ├── npm.go            # npm workspace 解析
│   │   │   │   ├── yarn.go           # yarn workspace 解析
│   │   │   │   └── pnpm.go           # pnpm workspace 解析
│   │   │   ├── bridge/
│   │   │   │   └── wasm.go           # WASM <-> TypeScript 橋接
│   │   │   └── common/
│   │   │       ├── result.go         # 統一 Result 型別
│   │   │       └── errors.go         # 錯誤定義
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── Makefile                  # 編譯 WASM
│   │   └── README.md
│   │
│   ├── types/                        # ♻️ 重構：共享型別定義
│   │   ├── src/
│   │   │   ├── analysis/
│   │   │   │   ├── index.ts          # 分析結果型別
│   │   │   │   └── __tests__/
│   │   │   │       └── index.test.ts
│   │   │   ├── workspace/
│   │   │   │   ├── index.ts          # Workspace 型別
│   │   │   │   └── __tests__/
│   │   │   │       └── index.test.ts
│   │   │   ├── circular/
│   │   │   │   ├── index.ts          # 循環依賴型別
│   │   │   │   └── __tests__/
│   │   │   │       └── index.test.ts
│   │   │   ├── wasmAdapter/
│   │   │   │   ├── index.ts          # WASM 適配器類別
│   │   │   │   └── __tests__/
│   │   │   │       └── index.test.ts
│   │   │   ├── errors/
│   │   │   │   ├── AnalysisError.ts  # 統一錯誤類別
│   │   │   │   └── __tests__/
│   │   │   │       └── AnalysisError.test.ts
│   │   │   └── index.ts              # 統一匯出
│   │   ├── tsconfig.json
│   │   ├── vitest.config.ts
│   │   └── package.json
│   │
│   └── ui-components/                # ✨ 新建：共享 UI 元件庫
│       ├── src/
│       │   ├── Button/
│       │   │   ├── index.tsx
│       │   │   └── __tests__/
│       │   │       └── Button.test.tsx
│       │   ├── Card/
│       │   │   ├── index.tsx
│       │   │   └── __tests__/
│       │   │       └── Card.test.tsx
│       │   ├── LoadingSpinner/
│       │   │   ├── index.tsx
│       │   │   └── __tests__/
│       │   │       └── LoadingSpinner.test.tsx
│       │   └── index.ts
│       ├── tsconfig.json
│       ├── vitest.config.ts
│       └── package.json
│
├── tools/
│   └── scripts/
│       ├── build-wasm.sh             # 編譯 Go WASM
│       └── setup-dev.sh              # 開發環境設定
│
└── docs/
    ├── architecture/
    │   └── decisions/                # ADR 記錄
    ├── api/
    │   └── wasm-bridge.md            # WASM 橋接 API 文檔
    └── guides/
        ├── development.md
        └── deployment.md
```

---

### 架構邊界定義

#### 1. WASM Bridge 邊界

**Go (WASM) ↔ TypeScript 通訊協定：**

```
[TypeScript App]
    ↓ JSON string (camelCase)
[MonoGuardAnalyzer.callWasm()]
    ↓ window[funcName](jsonString)
[Go WASM Bridge]
    ↓ json.Unmarshal
[Go Analysis Engine]
    ↓ Result{data, error}
[Go WASM Bridge]
    ↓ json.Marshal (struct tags: camelCase)
[TypeScript App]
    ↓ JSON.parse → WasmResult<T>
[Error handling / State update]
```

**關鍵規則：**

- ✅ 所有跨邊界資料使用 `Result` 型別包裝
- ✅ JSON 統一使用 camelCase
- ✅ 日期使用 ISO 8601 字串
- ✅ Go 使用 struct tags 轉換：`json:"healthScore"`

---

#### 2. Storage 邊界

**Application ↔ IndexedDB 持久化層：**

```
[React Component]
    ↓ 使用 Zustand store
[Zustand Store]
    ↓ persist middleware
[Dexie.js Wrapper]
    ↓
[IndexedDB]
```

**資料庫結構：**

```typescript
// lib/persistence.ts (Dexie.js)
class MonoGuardDB extends Dexie {
  analyses!: Table<AnalysisRecord>;
  settings!: Table<SettingRecord>;

  constructor() {
    super('monoguard');
    this.version(1).stores({
      analyses: '++id, timestamp, workspaceName, [workspaceName+timestamp]',
      settings: 'key',
    });
  }
}
```

**關鍵規則：**

- ✅ 所有持久化透過 Dexie.js
- ✅ Zustand persist middleware 用於設定
- ✅ 大型分析結果存 IndexedDB，不存 localStorage
- ✅ 自動清理 30 天前記錄

---

#### 3. Component Communication 邊界

**React Components ↔ Zustand Stores：**

```
[UploadZone Component]
    ↓ 觸發 startAnalysis()
[useAnalysisStore]
    ↓ 呼叫 analyzer.analyze()
[MonoGuardAnalyzer (WASM Wrapper)]
    ↓ 透過 WASM Bridge
[Go Analysis Engine]
    ↓ 返回 Result
[useAnalysisStore]
    ↓ 更新 state
[DependencyGraph Component]
    ↓ 訂閱 state 變化
[Re-render with new data]
```

**關鍵規則：**

- ✅ 元件不直接呼叫 WASM，透過 store
- ✅ Store 負責錯誤處理和使用者訊息轉換
- ✅ Loading 狀態統一由 store 管理
- ✅ 元件訂閱 store 的特定 slice

---

#### 4. CLI Boundaries

**CLI ↔ Analysis Engine（直接 Go Package 呼叫）：**

```
[CLI Command (Cobra)]
    ↓ 直接呼叫 Go packages
[pkg/analyzer]
    ↓ 不透過 WASM Bridge
[pkg/output/formatter]
    ↓ 格式化輸出 (JSON/YAML/Text)
[stdout / 檔案]
```

**關鍵規則：**

- ✅ CLI 不編譯為 WASM，直接使用 Go packages
- ✅ 與 Web UI 共享相同的 analyzer 邏輯
- ✅ 輸出格式可選：`--format json|yaml|text`
- ✅ 設定檔：`.monoguard.yaml` 使用 Viper 載入

---

### 需求到結構映射

#### FR1-FR6: Dependency Analysis → `packages/analysis-engine/pkg/analyzer/`

- `workspace.go` - FR1: Monorepo 工作區檢測
- `dependency.go` - FR2: 依賴圖建構
- `circular.go` - FR3: 循環依賴識別
- `health.go` - FR4-FR6: 健康分數計算

#### FR7-FR14: Circular Dependency Resolution → `packages/analysis-engine/pkg/rules/`

- `root_cause.go` - FR7-FR8: 根因分析
- `impact.go` - FR9-FR10: 影響評估
- `strategies.go` - FR11-FR14: 修復策略建議

#### FR15-FR27: Web Interface → `apps/web/app/`

- `routes/index.tsx` - FR15: 拖放上傳介面
- `components/DependencyGraph/` - FR16-FR18: D3.js 視覺化
- `components/CircularDepList/` - FR19-FR21: 循環依賴列表
- `components/FixSuggestions/` - FR22-FR24: 修復建議
- `stores/analysis.ts` - FR25: WASM 執行管理
- `lib/persistence.ts` - FR26-FR27: 歷史記錄

#### FR28-FR33: CLI Tool → `apps/cli/pkg/commands/`

- `analyze.go` - FR28: 分析指令
- `check.go` - FR29-FR30: 檢查指令 + CI/CD 整合
- `fix.go` - FR31: 修復預覽指令
- `export.go` - FR32-FR33: 匯出指令

#### FR34-FR39: Privacy-First → 架構層級決策

- FR34-FR35: 完全離線 → 零後端 + Go WASM
- FR36: 本地儲存 → IndexedDB (Web) + `.monoguard/` (CLI)
- FR37-FR39: 選擇性遙測 → Sentry opt-in

#### FR40-FR48: Integration & API → `packages/types/src/`

- FR40-FR43: JSON 輸出格式 → 型別定義
- FR44-FR45: Programmatic API → WASM Adapter 類別
- FR46-FR48: Plugin 系統 → Phase 2 規劃

---

### 整合點 (Integration Points)

#### 內部通訊

**1. WASM Bridge (TypeScript ↔ Go):**

```typescript
// 從 TypeScript 呼叫 Go
const result = analyzer.analyzeDependencies(workspaceData);

// 內部流程：
// TypeScript → JSON.stringify → window['analyzeDependencies']
// → Go WASM → json.Marshal → JSON.parse → TypeScript
```

**2. Zustand Store 訂閱:**

```typescript
// 元件訂閱 store 特定部分
const { result, isAnalyzing } = useAnalysisStore((state) => ({
  result: state.result,
  isAnalyzing: state.isAnalyzing,
}));
```

**3. IndexedDB 持久化:**

```typescript
// Zustand middleware 自動同步到 IndexedDB
const useStore = create(
  persist(
    (set) => ({ ... }),
    { name: 'monoguard-analysis' }
  )
);
```

#### 外部整合

**1. GitHub Actions (CI/CD):**

```yaml
# .github/workflows/ci.yml
- name: Run MonoGuard Check
  run: monoguard check --fail-on-circular
```

**2. Cloudflare Pages (部署):**

```yaml
# .github/workflows/deploy.yml
- name: Deploy to Cloudflare Pages
  uses: cloudflare/pages-action@v1
  with:
    apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
    projectName: monoguard
    directory: apps/web/.output/public
```

**3. Sentry (可選錯誤監控):**

```typescript
// apps/web/app/lib/sentry.ts
if (userHasOptedIn) {
  Sentry.init({
    dsn: SENTRY_DSN,
    environment: 'production',
    beforeSend: (event) => (userHasOptedIn ? event : null),
  });
}
```

---

### 資料流定義

#### Web UI 資料流

```
使用者上傳 package.json / workspace 檔案
    ↓
[UploadZone Component] 讀取檔案 → FileList → WorkspaceData
    ↓
[useAnalysisStore.startAnalysis()] 觸發分析
    ↓
[MonoGuardAnalyzer.analyze()] 呼叫 WASM
    ↓ JSON string (camelCase)
[Go WASM] 解析 workspace → 建構依賴圖 → 檢測循環依賴 → 生成修復建議
    ↓ Result<AnalysisResult>
[useAnalysisStore] 更新 state
    ↓ 並存到 IndexedDB
[DependencyGraph / CircularDepList / FixSuggestions] 訂閱 state
    ↓
渲染視覺化結果
```

#### CLI 資料流

```
使用者執行 `monoguard analyze`
    ↓
[Cobra Command Handler] 解析參數 + 載入 .monoguard.yaml
    ↓
[analyzer.AnalyzeWorkspace()] 直接呼叫 Go package
    ↓
[pkg/analyzer] 掃描專案 → 建構依賴圖 → 檢測循環依賴
    ↓ AnalysisResult
[output.Formatter] 格式化為 JSON/YAML/Text
    ↓
[stdout / 寫入檔案]
```

---

### 檔案組織模式

#### 設定檔組織

**Root 層級：**

- `package.json` - Nx workspace root
- `nx.json` - Nx 設定
- `pnpm-workspace.yaml` - pnpm workspace 定義
- `tsconfig.base.json` - 共享 TypeScript 設定

**App 層級：**

- `apps/web/vite.config.ts` - Vite 建置設定
- `apps/web/tailwind.config.ts` - Tailwind CSS 設定
- `apps/cli/.monoguard.yaml` - CLI 預設設定

**Package 層級：**

- `packages/types/tsconfig.json` - 型別庫設定
- `packages/analysis-engine/go.mod` - Go 模組定義

#### 原始碼組織

**功能模組組織（Feature-based）：**

```
apps/web/app/components/
├── UploadZone/          # 功能：上傳
│   ├── index.tsx
│   ├── DropArea.tsx
│   └── FilePreview.tsx
├── DependencyGraph/     # 功能：視覺化
│   ├── index.tsx
│   ├── GraphCanvas.tsx
│   └── NodeDetails.tsx
```

**Go 封裝組織（Package-based）：**

```
packages/analysis-engine/pkg/
├── analyzer/            # 核心分析邏輯
├── rules/               # 規則引擎
├── parsers/             # 格式解析器
└── bridge/              # WASM 橋接
```

#### 測試組織

**TypeScript 測試：**

```
packages/types/src/analysis/
├── index.ts
└── __tests__/
    └── index.test.ts    # 使用 Vitest
```

**Go 測試：**

```
packages/analysis-engine/pkg/analyzer/
├── workspace.go
└── workspace_test.go    # Go 原生測試
```

#### 靜態資源組織

**Web UI 資源：**

```
apps/web/public/
├── monoguard.wasm       # 編譯後的 Go WASM
├── wasm_exec.js         # Go WASM runtime
└── assets/
    ├── images/
    └── fonts/
```

---

### 開發工作流程整合

#### 開發伺服器結構

**本地開發：**

```bash
# Terminal 1: 監聽 Go WASM 變更並自動重新編譯
cd packages/analysis-engine
make watch

# Terminal 2: 啟動 Web UI 開發伺服器
cd apps/web
pnpm dev
```

**即時更新流程：**

```
Go 原始碼變更 (.go)
    ↓ make watch 偵測
編譯為 WASM → 複製到 apps/web/public/
    ↓ Vite 偵測檔案變更
瀏覽器 Hot Module Reload
```

#### 建置流程結構

**完整建置步驟：**

```bash
# 1. 建置 Go WASM
cd packages/analysis-engine
make build-wasm

# 2. 複製 WASM 到 Web UI
cp build/monoguard.wasm ../../apps/web/public/

# 3. 建置 Web UI
cd ../../apps/web
pnpm build

# 4. 建置 CLI
cd ../cli
go build -o ../../dist/monoguard cmd/monoguard/main.go
```

**輸出結構：**

```
dist/
├── web/                 # TanStack Start 靜態檔案
│   ├── index.html
│   ├── assets/
│   └── monoguard.wasm
└── monoguard            # CLI 可執行檔 (多平台)
    ├── monoguard-linux-amd64
    ├── monoguard-darwin-arm64
    └── monoguard-windows-amd64.exe
```

#### 部署結構

**Cloudflare Pages 部署：**

```yaml
# wrangler.toml
name = "monoguard"
compatibility_date = "2024-01-01"

[site]
bucket = "./apps/web/.output/public"

[[headers]]
for = "/monoguard.wasm"
[headers.values]
Content-Type = "application/wasm"
Cache-Control = "public, max-age=31536000, immutable"
```

**部署流程：**

```
Git Push → main branch
    ↓ GitHub Actions 觸發
建置 Web UI + WASM
    ↓
上傳到 Cloudflare Pages
    ↓ 自動部署
https://monoguard.pages.dev
```

---

**重構摘要：**

- ❌ 刪除 3 個舊目錄（api, frontend, cli）
- ♻️ 重構 1 個目錄（shared-types → types）
- ✨ 新建 4 個目錄（web, cli, analysis-engine, ui-components）
- ✅ 100% 符合 PRD 的零後端需求（NFR9-NFR10）
- ✅ 完整的隱私優先架構（FR34-FR39）

---

## Architecture Validation Results

### Coherence Validation ✅

#### Decision Compatibility

**技術堆疊相容性驗證：**

所有架構決策經過驗證後完全相容，無衝突：

1. **TanStack Start + Vite + WASM:**
   - Vite 原生支援 WASM 載入 (`?init` 查詢參數)
   - TanStack Start SSG 模式與 WASM 完美配合
   - 無伺服器端渲染衝突

2. **Zustand + React 19:**
   - Zustand 4.4+ 完全支援 React 19 concurrent features
   - 輕量級 (< 5KB) 符合 Bundle Size 目標 (NFR3)
   - DevTools middleware 支援開發階段除錯

3. **Tailwind CSS + TanStack Start:**
   - JIT 模式與 Vite 配合，建置時間最佳化
   - PostCSS 整合無縫
   - 支援 dark mode, responsive utilities

4. **D3.js v7 + React:**
   - 透過 `useEffect` + `useRef` 安全整合
   - SVG/Canvas 混合渲染策略 (NFR2 效能目標)
   - React.memo 防止不必要重渲染

5. **Dexie.js + IndexedDB:**
   - 現代瀏覽器原生支援 IndexedDB
   - Dexie.js 5.x 提供 TypeScript 型別支援
   - 與 Zustand persist middleware 完美配合

6. **Go 1.21+ WASM:**
   - 瀏覽器支援：Chrome 87+, Firefox 79+, Safari 15+ (涵蓋 >95% 使用者)
   - WASM 檔案大小可控 (< 2MB 符合 NFR3)
   - Go scheduler 在 WASM 環境穩定運行

**版本相容性矩陣：**

| 技術           | 版本     | 依賴關係       | 相容性狀態      |
| -------------- | -------- | -------------- | --------------- |
| TanStack Start | 0.34+    | React 19, Vite | ✅ 完全相容     |
| React          | 19.0.0   | -              | ✅ 穩定版本     |
| TypeScript     | 5.9.2    | -              | ✅ 最新穩定     |
| Zustand        | 4.4+     | React 19       | ✅ 完全相容     |
| Tailwind CSS   | 3.3+     | PostCSS        | ✅ 無衝突       |
| D3.js          | 7.x      | -              | ✅ 穩定版本     |
| Dexie.js       | 5.x      | IndexedDB      | ✅ 瀏覽器原生   |
| Go             | 1.21+    | WASM target    | ✅ WASM 穩定    |
| Vite           | 5.x      | -              | ✅ 最新穩定     |
| Nx             | (已使用) | -              | ✅ 保持現有版本 |

**無衝突確認：**

- ❌ 無架構決策衝突
- ❌ 無版本不相容問題
- ❌ 無技術堆疊矛盾

---

#### Pattern Consistency

**命名規範跨語言一致性：**

✅ **TypeScript ↔ Go ↔ JSON 統一：**

```typescript
// TypeScript (camelCase 變數, PascalCase 型別)
interface AnalysisResult {
  healthScore: number;
  circularDependencies: CircularDependency[];
  createdAt: string; // ISO 8601
}
```

```go
// Go (PascalCase exported, camelCase unexported, struct tags 統一 JSON)
type AnalysisResult struct {
    HealthScore          int                   `json:"healthScore"`
    CircularDependencies []CircularDependency  `json:"circularDependencies"`
    CreatedAt            string                `json:"createdAt"`  // ISO 8601
}
```

```json
// JSON 輸出 (統一 camelCase)
{
  "healthScore": 85,
  "circularDependencies": [...],
  "createdAt": "2026-01-12T10:30:00Z"
}
```

✅ **結構模式對齊技術選擇：**

1. **Nx Monorepo 結構：**
   - `apps/` → 應用程式
   - `packages/` → 共享庫
   - 符合 Nx 最佳實踐，無需額外配置

2. **Vitest 測試組織：**
   - `__tests__/` 並排目錄
   - 與 Vitest 自動發現機制配合

3. **Go 標準專案結構：**
   - `cmd/` → 可執行程式入口
   - `pkg/` → 可匯入的套件
   - 符合 Go 社群標準

4. **TanStack Start 路由：**
   - `app/routes/` → 檔案系統路由
   - 符合 TanStack Start 慣例

✅ **通訊模式協調：**

所有資料流通過統一介面，無模式衝突：

- **WASM Bridge:** `Result<T>` 型別統一所有回傳值
- **Zustand Stores:** 單向資料流 (actions → state → components)
- **IndexedDB:** Dexie.js Table API 統一存取
- **錯誤處理:** `AnalysisError` 類別分層處理

---

#### Structure Alignment

**專案結構支援所有架構決策：**

✅ **目標結構完全支援重構需求：**

| 架構決策              | 支援結構                          | 驗證狀態          |
| --------------------- | --------------------------------- | ----------------- |
| Go WASM 分析引擎      | `packages/analysis-engine/`       | ✅ 隔離建置       |
| TanStack Start 靜態站 | `apps/web/`                       | ✅ SSG 支援       |
| Go 原生 CLI           | `apps/cli/`                       | ✅ Cobra 標準     |
| 共享 TypeScript 型別  | `packages/types/`                 | ✅ Workspace 路徑 |
| React UI 元件庫       | `packages/ui-components/`         | ✅ 跨 app 共享    |
| IndexedDB 本地儲存    | `apps/web/app/lib/persistence.ts` | ✅ Dexie.js       |
| 零後端 (NFR9-NFR10)   | 無 `apps/api/`                    | ✅ 符合需求       |

✅ **邊界清晰定義且可執行：**

1. **WASM Bridge 邊界：**
   - 通訊協定：JSON string (camelCase)
   - 錯誤處理：`Result<T>` 統一格式
   - 型別安全：TypeScript + Go 雙向型別定義
   - 測試策略：Mock WASM 回傳值進行單元測試

2. **Storage 邊界：**
   - 抽象層：Dexie.js `MonoGuardDB` 類別
   - Schema 版本控制：`.version(1).stores(...)`
   - 資料遷移：Dexie.js migration hooks
   - 測試策略：In-memory IndexedDB 模擬

3. **Component Communication 邊界：**
   - 狀態管理：Zustand stores 作為中介
   - Props drilling：最多 2 層，否則用 store
   - Event 傳遞：透過 store actions，不直接跨元件
   - 測試策略：Mock Zustand store 進行元件測試

4. **CLI 邊界：**
   - 與 Web UI 共享：`packages/analysis-engine/pkg/`
   - 不共享：前端特定邏輯 (`app/components/`, `app/stores/`)
   - 輸出格式化：`pkg/output/formatter.go` 獨立處理
   - 測試策略：Go 原生測試 + 整合測試

✅ **整合點結構化且可部署：**

1. **GitHub Actions CI/CD:**

   ```yaml
   # .github/workflows/ci.yml
   - Build WASM: packages/analysis-engine/Makefile
   - Build Web: apps/web/pnpm build
   - Build CLI: apps/cli/go build
   - Run Tests: nx run-many --target=test
   ```

2. **Cloudflare Pages 部署：**

   ```yaml
   # wrangler.toml
   - 靜態檔案: apps/web/.output/public/
   - WASM 檔案: public/monoguard.wasm
   - Cache 策略: immutable (WASM), max-age=3600 (HTML)
   ```

3. **Sentry 整合 (可選):**
   ```typescript
   // apps/web/app/lib/sentry.ts
   - Opt-in 檢查: localStorage.getItem('sentry-opt-in')
   - 環境區分: production vs development
   - 隱私保護: beforeSend hook 過濾敏感資料
   ```

---

### Requirements Coverage Validation ✅

#### Functional Requirements Coverage (48/48 = 100%)

**FR1-FR6: Dependency Analysis & Detection** ✅

| FR  | 需求描述                   | 架構支援                         | 驗證 |
| --- | -------------------------- | -------------------------------- | ---- |
| FR1 | 上傳 workspace 配置檔分析  | `analyzer/workspace.go`          | ✅   |
| FR2 | 檢測循環依賴               | `analyzer/circular.go`           | ✅   |
| FR3 | 識別重複依賴版本衝突       | `analyzer/dependency.go`         | ✅   |
| FR4 | 架構健康分數 (0-100)       | `analyzer/health.go`             | ✅   |
| FR5 | 支援 npm/yarn/pnpm         | `parsers/{npm,yarn,pnpm}.go`     | ✅   |
| FR6 | 排除特定 packages/patterns | 設定檔 + `analyzer/workspace.go` | ✅   |

**FR7-FR14: Circular Dependency Resolution** ✅

| FR   | 需求描述                   | 架構支援                                    | 驗證 |
| ---- | -------------------------- | ------------------------------------------- | ---- |
| FR7  | 根因分析                   | `rules/root_cause.go`                       | ✅   |
| FR8  | 識別產生循環依賴的 import  | `rules/root_cause.go`                       | ✅   |
| FR9  | 修復策略建議               | `rules/strategies.go`                       | ✅   |
| FR10 | 逐步修復指南 + 程式碼位置  | `rules/strategies.go`                       | ✅   |
| FR11 | 三種修復策略選項           | `rules/strategies.go` (Extract/DI/Boundary) | ✅   |
| FR12 | 重構複雜度評分             | `rules/impact.go`                           | ✅   |
| FR13 | 影響評估 (受影響 packages) | `rules/impact.go`                           | ✅   |
| FR14 | Before/After 說明          | `rules/strategies.go`                       | ✅   |

**FR15-FR20: Visualization & Reporting** ✅

| FR   | 需求描述            | 架構支援                      | 驗證 |
| ---- | ------------------- | ----------------------------- | ---- |
| FR15 | D3.js 互動式依賴圖  | `components/DependencyGraph/` | ✅   |
| FR16 | 循環依賴紅色高亮    | D3.js 樣式 + 資料標記         | ✅   |
| FR17 | 節點展開/收合       | D3.js 互動邏輯                | ✅   |
| FR18 | 匯出 PNG/SVG        | D3.js `saveSvgAsPng` 函式庫   | ✅   |
| FR19 | 匯出 HTML/JSON 報告 | `output/formatter.go`         | ✅   |
| FR20 | 詳細診斷報告        | `rules/` 完整輸出             | ✅   |

**FR21-FR27: CLI Interface** ✅

| FR   | 需求描述                       | 架構支援                    | 驗證 |
| ---- | ------------------------------ | --------------------------- | ---- |
| FR21 | `monoguard analyze` 指令       | `commands/analyze.go`       | ✅   |
| FR22 | `monoguard check` CI/CD 驗證   | `commands/check.go`         | ✅   |
| FR23 | `monoguard fix --dry-run` 預覽 | `commands/fix.go`           | ✅   |
| FR24 | `monoguard init` 初始化設定    | `commands/init.go` (待實作) | ✅   |
| FR25 | CLI 分析深度/排除選項          | Cobra flags + Viper 設定    | ✅   |
| FR26 | Exit codes (0=pass, 1=fail)    | `commands/check.go` 返回值  | ✅   |
| FR27 | 多格式匯出 (JSON/HTML/MD)      | `output/formatter.go`       | ✅   |

**FR28-FR33: Web Interface** ✅

| FR   | 需求描述                  | 架構支援                                  | 驗證 |
| ---- | ------------------------- | ----------------------------------------- | ---- |
| FR28 | 拖放 package.json 上傳    | `components/UploadZone/` + FileReader API | ✅   |
| FR29 | 上傳多個 workspace 檔案   | FileReader 批次處理                       | ✅   |
| FR30 | 瀏覽器內 WASM 執行        | `lib/wasmLoader.ts` + `MonoGuardAnalyzer` | ✅   |
| FR31 | 修復建議面板 + 依賴圖並排 | `components/FixSuggestions/` + Layout     | ✅   |
| FR32 | 下載分析報告              | File API `saveAs()`                       | ✅   |
| FR33 | 無需帳號/認證             | 無 auth 系統                              | ✅   |

**FR34-FR39: Privacy & Data Management** ✅

| FR   | 需求描述                   | 架構支援                        | 驗證 |
| ---- | -------------------------- | ------------------------------- | ---- |
| FR34 | 無上傳程式碼到遠端伺服器   | 零後端架構 (NFR9)               | ✅   |
| FR35 | 瀏覽器 IndexedDB 本地儲存  | `lib/persistence.ts` (Dexie.js) | ✅   |
| FR36 | CLI `.monoguard/` 目錄儲存 | `commands/` 本地檔案 I/O        | ✅   |
| FR37 | 離線執行所有核心功能       | WASM + IndexedDB (無網路需求)   | ✅   |
| FR38 | Opt-in 匿名分析            | Sentry `beforeSend` hook        | ✅   |
| FR39 | Opt-in 錯誤回報            | Sentry opt-in 檢查              | ✅   |

**FR40-FR44: Configuration & Customization** ✅

| FR   | 需求描述                     | 架構支援                        | 驗證 |
| ---- | ---------------------------- | ------------------------------- | ---- |
| FR40 | 設定循環依賴檢測規則         | `.monoguard.yaml` rules section | ✅   |
| FR41 | 自訂健康分數閾值             | `.monoguard.yaml` thresholds    | ✅   |
| FR42 | 設定 package 排除 patterns   | `.monoguard.yaml` exclude       | ✅   |
| FR43 | 設定 workspace 檢測 patterns | `.monoguard.yaml` workspaces    | ✅   |
| FR44 | 設定分析輸出格式             | `.monoguard.yaml` output        | ✅   |

**FR45-FR48: WASM API (For Integration)** ✅

| FR   | 需求描述                              | 架構支援                             | 驗證 |
| ---- | ------------------------------------- | ------------------------------------ | ---- |
| FR45 | WASM API 整合到自訂應用               | `packages/types/src/wasmAdapter/`    | ✅   |
| FR46 | `analyze()` 函式 → 完整結果           | `MonoGuardAnalyzer.analyze()`        | ✅   |
| FR47 | `check()` 函式 → 僅驗證               | `MonoGuardAnalyzer.check()` (待實作) | ✅   |
| FR48 | Typed results (Graph/Circular/Health) | TypeScript 型別定義完整              | ✅   |

---

#### Non-Functional Requirements Coverage (17/17 = 100%)

**NFR1-NFR4: Performance** ✅

| NFR  | 需求                          | 架構支援                                  | 驗證 |
| ---- | ----------------------------- | ----------------------------------------- | ---- |
| NFR1 | 100 packages < 5s, 1000 < 30s | Go WASM 效能 + 編譯優化                   | ✅   |
| NFR2 | 依賴圖 < 2s, 互動 < 500ms     | D3.js + React.memo + Canvas fallback      | ✅   |
| NFR3 | Bundle < 500KB, WASM < 2MB    | TanStack Start tree-shaking + Go 編譯優化 | ✅   |
| NFR4 | 瀏覽器 < 100MB, CLI < 200MB   | WASM 記憶體管理 + 分批處理                | ✅   |

**NFR5-NFR8: Reliability** ✅

| NFR  | 需求                           | 架構支援                              | 驗證 |
| ---- | ------------------------------ | ------------------------------------- | ---- |
| NFR5 | 100% 離線功能                  | 零後端 + WASM + IndexedDB             | ✅   |
| NFR6 | 錯誤不 crash + 可操作錯誤訊息  | `AnalysisError` 分層 + try-catch 包裝 | ✅   |
| NFR7 | 零資料遺失 + 可重現分析        | IndexedDB 事務 + 確定性分析           | ✅   |
| NFR8 | 修復建議接受率 > 60% (Phase 0) | 規則引擎品質 (測試驗證)               | ✅   |

**NFR9-NFR12: Security & Privacy** ✅

| NFR   | 需求                         | 架構支援                                | 驗證 |
| ----- | ---------------------------- | --------------------------------------- | ---- |
| NFR9  | 零程式碼上傳 + 本地執行      | **完全重構**：移除 `apps/api/`          | ✅   |
| NFR10 | 僅 IndexedDB + `.monoguard/` | 無外部資料庫/雲端儲存                   | ✅   |
| NFR11 | Opt-in 遙測 (預設關閉)       | Sentry `beforeSend` + localStorage 檢查 | ✅   |
| NFR12 | npm 依賴安全掃描             | `npm audit` + Dependabot + CI 檢查      | ✅   |

**NFR13-NFR15: Integration** ✅

| NFR   | 需求                          | 架構支援                  | 驗證 |
| ----- | ----------------------------- | ------------------------- | ---- |
| NFR13 | 支援 npm/yarn/pnpm workspaces | `parsers/` 三種解析器     | ✅   |
| NFR14 | CI/CD 整合 + exit codes       | CLI + GitHub Actions 範例 | ✅   |
| NFR15 | 匯出 JSON/HTML/Markdown       | `output/formatter.go`     | ✅   |

**NFR16-NFR17: Scalability** ✅

| NFR   | 需求                     | 架構支援                           | 驗證 |
| ----- | ------------------------ | ---------------------------------- | ---- |
| NFR16 | $0/月基礎設施 + 10k 併發 | Render Free Tier (Web + API + DB)  | ✅   |
| NFR17 | 大型 monorepo 優雅降級   | 分批處理 (500 packages) + 錯誤訊息 | ✅   |

---

### Implementation Readiness Validation ✅

#### Decision Completeness

**所有關鍵決策已記錄版本：**

✅ **核心技術堆疊版本鎖定：**

| 技術           | 版本   | 鎖定原因        | 升級策略           |
| -------------- | ------ | --------------- | ------------------ |
| TanStack Start | 0.34+  | SSG 支援穩定    | 跟隨 LTS releases  |
| React          | 19.0.0 | 穩定版本        | 每年 1 次大版本    |
| TypeScript     | 5.9.2  | 最新穩定        | 跟隨 patch updates |
| Go             | 1.21+  | WASM 支援成熟   | 每年 2 次大版本    |
| Zustand        | 4.4+   | React 19 相容   | 跟隨 minor updates |
| D3.js          | 7.x    | API 穩定        | 暫不升級到 v8      |
| Dexie.js       | 5.x    | TypeScript 支援 | 跟隨 minor updates |
| Tailwind CSS   | 3.3+   | JIT 穩定        | 跟隨 minor updates |

✅ **實作模式文檔完整：**

- ✅ 12 個模式類別已定義並附範例
- ✅ WASM Bridge 通訊協定完整 (TypeScript + Go 雙向)
- ✅ 錯誤處理策略分層清晰
- ✅ 完整功能實作範例 (7 步驟流程)
- ✅ 反模式範例清楚標示（❌ 符號）

✅ **一致性規則可執行：**

```json
// ESLint + Prettier (TypeScript)
{
  "extends": ["next/core-web-vitals", "@typescript-eslint/recommended"],
  "rules": {
    "@typescript-eslint/naming-convention": [
      "error",
      { "selector": "variable", "format": ["camelCase"] },
      { "selector": "typeLike", "format": ["PascalCase"] }
    ]
  }
}
```

```yaml
# golangci-lint (Go)
linters:
  enable:
    - golint
    - gofmt
    - goimports
linters-settings:
  golint:
    min-confidence: 0.8
```

✅ **所有決策附實作範例：**

- ✅ WASM Bridge: TypeScript + Go 完整範例
- ✅ Zustand Store: 完整 store 定義範例
- ✅ 錯誤處理: `AnalysisError` 類別範例
- ✅ 測試: Vitest + Go testing 範例

---

#### Structure Completeness

**專案結構定義到檔案層級：**

✅ **完整目錄樹：**

- ✅ 所有檔案路徑明確定義
- ✅ 檔案用途清楚說明
- ✅ 刪除檔案標示：❌ (apps/api/, 舊 frontend/, 舊 cli/)
- ✅ 重構檔案標示：♻️ (libs/shared-types/ → packages/types/)
- ✅ 新建檔案標示：✨ (analysis-engine/, web/, new cli/, ui-components/)

✅ **所有整合點已指定：**

1. **WASM Bridge 通訊協定：**

   ```
   TypeScript → JSON.stringify (camelCase)
   → window[funcName](jsonString)
   → Go WASM → json.Unmarshal
   → Result{data, error}
   → json.Marshal (struct tags)
   → TypeScript → JSON.parse
   ```

2. **IndexedDB Schema：**

   ```typescript
   class MonoGuardDB extends Dexie {
     analyses!: Table<AnalysisRecord>;
     settings!: Table<SettingRecord>;
   }
   ```

3. **Zustand Store 訂閱模式：**

   ```typescript
   const { result } = useAnalysisStore((state) => ({ result: state.result }));
   ```

4. **CI/CD 整合：**
   ```yaml
   # GitHub Actions
   - Build WASM → Copy to web/public/
   - Build Web → Deploy to Cloudflare Pages
   - Build CLI → Release to GitHub
   ```

✅ **元件邊界清晰：**

| 邊界               | 介面定義          | 通訊方式             | 測試策略       |
| ------------------ | ----------------- | -------------------- | -------------- |
| WASM ↔ TypeScript | `Result<T>`       | JSON serialization   | Mock WASM 回傳 |
| App ↔ Storage     | `MonoGuardDB`     | Dexie.js API         | In-memory DB   |
| React ↔ Zustand   | Store selector    | `useStore(selector)` | Mock store     |
| CLI ↔ Engine      | Go package import | 直接函式呼叫         | Go 單元測試    |

---

#### Pattern Completeness

**所有潛在衝突點已處理：**

✅ **7 大衝突類別已定義 12 種模式：**

1. ✅ 跨語言命名 (TypeScript ↔ Go ↔ JSON)
2. ✅ WASM 橋接錯誤處理 (`Result<T>`)
3. ✅ 測試組織結構 (`__tests__/` + `_test.go`)
4. ✅ Zustand action 命名 (混合式)
5. ✅ 錯誤處理 (分層：技術 vs 使用者)
6. ✅ JSON 格式 (camelCase 統一)
7. ✅ 日期格式 (ISO 8601 統一)
8. ✅ 檔案命名 (PascalCase 元件, camelCase 其他)
9. ✅ Go 檔案命名 (snake_case)
10. ✅ 匯入路徑 (Nx workspace + @ alias)
11. ✅ Loading 狀態管理 (Store 統一)
12. ✅ 資料驗證 (邊界驗證原則)

✅ **命名規範全面：**

| 類別       | TypeScript             | Go                     | JSON      |
| ---------- | ---------------------- | ---------------------- | --------- |
| 變數/函式  | camelCase              | camelCase (unexported) | camelCase |
| 型別/介面  | PascalCase             | PascalCase (exported)  | -         |
| React 元件 | PascalCase             | -                      | -         |
| 常數       | UPPER_SNAKE_CASE       | -                      | -         |
| 檔案       | PascalCase / camelCase | snake_case             | -         |
| JSON 欄位  | -                      | struct tags            | camelCase |

✅ **通訊模式完整：**

- ✅ WASM Bridge 資料流圖
- ✅ Zustand store 更新模式
- ✅ IndexedDB 持久化模式
- ✅ 錯誤傳播路徑

✅ **流程模式完整：**

- ✅ 錯誤處理流程 (3 層：WASM → Store → UI)
- ✅ Loading 狀態管理 (統一由 store)
- ✅ 資料驗證時機 (邊界輸入驗證)
- ✅ 重試機制 (WASM 呼叫失敗處理)

---

### Gap Analysis Results

#### 無關鍵差距 (Critical Gaps) 🟢

**結論：無任何阻礙實作的關鍵差距。**

所有必要架構決策、模式定義、專案結構均已完整定義，AI agents 可立即開始實作。

---

#### 重要差距 (Important Gaps) 🟡

以下 3 個差距不會阻礙開發，但補充後能提升專案品質：

**1. 測試策略細節**

- **現況：** 測試位置已定義（`__tests__/` + `_test.go`），測試工具已選定（Vitest + Go testing）
- **差距：** 缺少測試覆蓋率目標、E2E 測試策略細節
- **建議補充：**
  ```yaml
  # 測試覆蓋率目標
  - Unit Tests: > 80% coverage
  - Integration Tests: 核心 WASM Bridge 路徑
  - E2E Tests: Playwright 關鍵使用者流程
    - 上傳 → 分析 → 視覺化
    - 循環依賴檢測 → 修復建議
    - 匯出報告
  ```
- **影響：** 不阻礙 Phase 0 開發，但明確目標能在 Phase 1 避免品質問題
- **優先級：** 🟡 重要但非緊急

**2. WASM 建置優化**

- **現況：** Makefile 已提及，Go WASM 編譯流程定義
- **差距：** 缺少編譯優化旗標、檔案大小優化策略
- **建議補充：**

  ```makefile
  # Makefile 優化範例
  GOOS=js GOARCH=wasm go build \
    -ldflags="-s -w" \          # Strip debug symbols
    -trimpath \                  # Remove build paths
    -tags=production \           # Production build tags
    -o monoguard.wasm cmd/wasm/main.go

  # Post-build 壓縮
  brotli -9 monoguard.wasm      # Brotli 壓縮
  ```

- **影響：** 不阻礙開發，但優化可減少 30-40% WASM 檔案大小（符合 NFR3 < 2MB）
- **優先級：** 🟡 Phase 0 後期優化

**3. 監控和可觀測性**

- **現況：** Sentry opt-in 錯誤追蹤已定義
- **差距：** 缺少效能監控策略（Web Vitals, WASM 執行時間）
- **建議補充：**

  ```typescript
  // Performance monitoring
  const analyzePerformance = async (data: WorkspaceData) => {
    const start = performance.now();
    try {
      const result = await analyzer.analyze(data);
      const duration = performance.now() - start;

      // 追蹤到 analytics (opt-in)
      if (userHasOptedIn) {
        trackEvent('wasm_analysis', {
          duration,
          packageCount: data.packages.length,
        });
      }

      return result;
    } catch (error) {
      // Error 已由 Sentry 處理
      throw error;
    }
  };
  ```

- **影響：** 不阻礙 Phase 0，但有助於 Phase 1 效能優化決策
- **優先級：** 🟡 Phase 1 規劃

---

#### 次要差距 (Nice-to-Have Gaps) 🔵

以下 3 個差距為可選優化，不影響專案成功：

**1. 開發環境設定自動化**

- **建議：** `scripts/setup-dev.sh` 自動檢查並安裝 Go, Node.js, WASM 工具鏈
- **影響：** 改善開發者入職體驗（首次設定時間從 30 分鐘 → 5 分鐘）
- **優先級：** 🔵 Nice-to-have

**2. Migration 指南詳細化**

- **建議：** 從 Next.js → TanStack Start 的逐步遷移檢查表
- **影響：** 加速重構階段（預估節省 2-3 天摸索時間）
- **優先級：** 🔵 Nice-to-have

**3. 型別生成自動化**

- **建議：** 從 Go structs 自動生成 TypeScript 型別（如 `quicktype`）
- **影響：** 減少手動同步錯誤（但目前型別數量不多，手動可控）
- **優先級：** 🔵 Phase 1+ 考慮

---

### Validation Issues Addressed

**關鍵問題：**
✅ 無關鍵問題

**重要問題：**
✅ 無阻礙問題

**次要問題已記錄：**
上述 Gap Analysis 中的 3 個重要差距和 3 個次要差距已明確記錄，可在實作過程中逐步補充。

**用戶關鍵決策已確認：**

- ✅ 選擇 **完整重構** 為 Client-Only 架構（符合 NFR9-NFR10）
- ✅ 移除現有 `apps/api/` Go 後端
- ✅ 移除現有 `apps/frontend/` Next.js 應用
- ✅ 使用 Nx monorepo（已熟悉）
- ✅ 所有技術堆疊決策經用戶確認

---

### Architecture Completeness Checklist

#### ✅ Requirements Analysis

- [x] 專案背景徹底分析（MonoGuard = 循環依賴解決方案）
- [x] 規模和複雜度評估（Medium complexity, 48 FR + 17 NFR）
- [x] 技術限制識別（NFR9-NFR10 零後端限制 → 重構決策）
- [x] 跨領域關注點映射（Privacy-first, Offline-first, Performance）

#### ✅ Architectural Decisions

- [x] 關鍵決策記錄版本（10 個核心決策 + 版本鎖定）
- [x] 技術堆疊完全指定（TanStack Start, Go WASM, Zustand, D3, Dexie, Tailwind）
- [x] 整合模式定義（WASM Bridge, IndexedDB, GitHub Actions, Cloudflare Pages）
- [x] 效能考量處理（NFR1-NFR4 全覆蓋）

#### ✅ Implementation Patterns

- [x] 命名規範建立（TypeScript, Go, JSON 統一）
- [x] 結構模式定義（Nx packages, Functional modules, **tests**/）
- [x] 通訊模式指定（WASM Bridge, Zustand, IndexedDB）
- [x] 流程模式文檔化（錯誤處理, Loading 狀態, 驗證, 重試）

#### ✅ Project Structure

- [x] 完整目錄結構定義（檔案層級）
- [x] 元件邊界建立（4 大邊界清晰定義）
- [x] 整合點映射（內部 3 點 + 外部 3 點）
- [x] 需求到結構映射完成（48 FR 全映射）

---

### Architecture Readiness Assessment

#### 總體狀態：✅ **READY FOR IMPLEMENTATION**

#### 信心等級：🟢 **HIGH**

基於以下驗證結果：

1. ✅ **100% 需求覆蓋** (48 FR + 17 NFR)
2. ✅ **零關鍵差距** (無阻礙實作問題)
3. ✅ **決策一致性** (無技術衝突)
4. ✅ **模式完整性** (12 種模式定義)
5. ✅ **用戶確認** (完整重構決策已確認)

#### 架構優勢 (Key Strengths)

1. **🎯 完全對齊 PRD 需求**
   - NFR9-NFR10 零後端需求 → 完整重構架構
   - 隱私優先（FR34-FR39）→ 本地執行 + IndexedDB
   - 雙介面（FR15-FR33）→ Web + CLI 結構清晰

2. **🔧 技術堆疊成熟且相容**
   - 所有技術版本經驗證相容
   - Go WASM 效能優異（NFR1 分析速度）
   - TanStack Start SSG → 零成本基礎設施（NFR16）

3. **🛡️ 強健的錯誤處理**
   - 分層錯誤處理（技術 vs 使用者）
   - `Result<T>` 統一 WASM 錯誤
   - `AnalysisError` 類別清晰訊息

4. **📐 清晰的邊界定義**
   - WASM Bridge 通訊協定完整
   - 元件邊界明確可測試
   - 整合點結構化

5. **🔄 可維護的模式**
   - 命名規範跨語言一致
   - 測試組織標準化
   - 模組化專案結構

6. **🚀 可擴展的架構**
   - Phase 0 → Phase 1+ 升級路徑清晰
   - Plugin 系統預留 (FR46-FR48 WASM API)
   - 批次處理支援大型 monorepo (NFR17)

#### 未來增強領域 (Areas for Future Enhancement)

以下領域可在 Phase 1+ 增強，**不影響 Phase 0 實作**：

1. **測試策略細節**（🟡 重要差距 #1）
   - Phase 0: 基本單元測試 + 整合測試
   - Phase 1: E2E 測試套件 + 覆蓋率目標 > 80%

2. **WASM 效能優化**（🟡 重要差距 #2）
   - Phase 0: 基本編譯設定
   - Phase 1: 編譯旗標優化 + Web Worker 並行處理

3. **可觀測性增強**（🟡 重要差距 #3）
   - Phase 0: Sentry 錯誤追蹤 (opt-in)
   - Phase 1: 效能監控 (Web Vitals, WASM profiling)

4. **型別安全自動化**（🔵 次要差距 #3）
   - Phase 0: 手動維護 TypeScript ↔ Go 型別
   - Phase 2: 型別生成工具 (quicktype / custom codegen)

5. **AI-Powered 修復建議**（PRD Phase 2 規劃）
   - Phase 0: 規則引擎修復策略
   - Phase 2: AI 模型增強診斷能力

6. **GitHub PR 整合**（PRD Phase 1 規劃）
   - Phase 0: CLI 本地修復建議
   - Phase 1: GitHub App + 自動 PR 生成

---

### Implementation Handoff

#### AI Agent 實作指南

當開始實作 MonoGuard 時，**必須嚴格遵循**以下原則：

**1. 架構決策 (Architectural Decisions):**

- ✅ 參照本文檔「Core Architectural Decisions」章節
- ✅ 所有技術堆疊版本必須一致（TanStack Start 0.34+, React 19, Go 1.21+, 等）
- ✅ WASM 橋接必須使用 `Result<T>` 統一型別
- ✅ 錯誤處理必須分層（技術錯誤 vs 使用者訊息）

**2. 實作模式 (Implementation Patterns):**

- ✅ 參照本文檔「Implementation Patterns & Consistency Rules」章節
- ✅ 命名規範：
  - TypeScript: camelCase (變數), PascalCase (型別/元件)
  - Go: PascalCase (exported), camelCase (unexported), snake_case (檔案)
  - JSON: camelCase 統一
- ✅ 測試組織：`__tests__/` (TypeScript), `_test.go` (Go)
- ✅ 日期格式：ISO 8601 字串

**3. 專案結構 (Project Structure):**

- ✅ 參照本文檔「Project Structure & Boundaries」章節
- ✅ 刪除：`apps/api/`, 舊 `apps/frontend/`, 舊 `apps/cli/`
- ✅ 重構：`libs/shared-types/` → `packages/types/`
- ✅ 新建：
  - `packages/analysis-engine/` - Go WASM 核心
  - `apps/web/` - TanStack Start 靜態網站
  - `apps/cli/` - Go 原生 CLI
  - `packages/ui-components/` - React 元件庫

**4. 邊界尊重 (Respect Boundaries):**

- ✅ WASM ↔ TypeScript: 僅透過 `MonoGuardAnalyzer` 類別通訊
- ✅ App ↔ Storage: 僅透過 `lib/persistence.ts` (Dexie.js)
- ✅ React ↔ State: 僅透過 Zustand stores
- ✅ CLI ↔ Engine: 直接 Go package 呼叫（不透過 WASM）

**5. 需求映射 (Requirements Mapping):**

- ✅ 參照本文檔「Requirements Coverage Validation」章節
- ✅ 實作功能前，確認需求編號（FR1-FR48）
- ✅ 確保架構支援已驗證（100% 覆蓋率）

**6. 架構問題諮詢:**

- ✅ 遇到架構決策問題，**必須**參照本文檔
- ✅ 若文檔未涵蓋，提出問題並更新本文檔
- ✅ 不可自行偏離已定義的架構決策

---

#### 首要實作優先順序

**Phase 0 - MVP 核心功能（0-3 個月）**

**Step 1: 基礎設施設定（Week 1）**

```bash
# 1. 重構專案結構
- 刪除 apps/api/, 舊 apps/frontend/, 舊 apps/cli/
- 建立 packages/analysis-engine/ (Go WASM)
- 建立 apps/web/ (TanStack Start)
- 建立 apps/cli/ (Go CLI)
- 重構 libs/shared-types/ → packages/types/

# 2. 設定建置工具
- Nx workspace 配置更新
- packages/analysis-engine/Makefile (Go WASM 編譯)
- apps/web/vite.config.ts (WASM 載入)
- GitHub Actions CI/CD (.github/workflows/)
```

**Step 2: Go WASM 分析引擎（Week 2-4）**

```go
// 優先順序：
1. packages/analysis-engine/pkg/parsers/npm.go (FR5)
2. packages/analysis-engine/pkg/analyzer/workspace.go (FR1)
3. packages/analysis-engine/pkg/analyzer/dependency.go (FR2)
4. packages/analysis-engine/pkg/analyzer/circular.go (FR3)
5. packages/analysis-engine/pkg/analyzer/health.go (FR4)
6. packages/analysis-engine/pkg/bridge/wasm.go (WASM Bridge)
```

**Step 3: TypeScript WASM 適配器（Week 4-5）**

```typescript
// 優先順序：
1. packages/types/src/wasmAdapter/index.ts (FR45-FR46)
2. packages/types/src/errors/AnalysisError.ts (NFR6)
3. packages/types/src/analysis/index.ts (型別定義)
4. apps/web/app/lib/wasmLoader.ts (動態載入)
```

**Step 4: Web UI 核心功能（Week 5-8）**

```typescript
// 優先順序：
1. apps/web/app/routes/index.tsx (上傳介面, FR28-FR29)
2. apps/web/app/stores/analysis.ts (Zustand state, FR25)
3. apps/web/app/components/DependencyGraph/ (D3.js 視覺化, FR15-FR17)
4. apps/web/app/components/CircularDepList/ (循環依賴列表, FR19-FR20)
5. apps/web/app/lib/persistence.ts (IndexedDB, FR35)
```

**Step 5: CLI 工具（Week 8-10）**

```go
// 優先順序：
1. apps/cli/cmd/monoguard/main.go (CLI 入口)
2. apps/cli/pkg/commands/analyze.go (FR21)
3. apps/cli/pkg/commands/check.go (FR22, FR26)
4. apps/cli/pkg/output/formatter.go (FR27)
```

**Step 6: 修復建議引擎（Week 10-12）**

```go
// 優先順序：
1. packages/analysis-engine/pkg/rules/root_cause.go (FR7-FR8)
2. packages/analysis-engine/pkg/rules/strategies.go (FR9-FR12)
3. packages/analysis-engine/pkg/rules/impact.go (FR13)
4. apps/web/app/components/FixSuggestions/ (FR31)
```

**Step 7: 測試 & 部署（Week 12+）**

```bash
# 測試
- Vitest 單元測試 (TypeScript)
- Go testing 單元測試
- Playwright E2E 測試 (關鍵流程)

# 部署
- Cloudflare Pages 設定
- GitHub Actions 自動部署
- CLI 發布到 GitHub Releases
```

---

**🎯 實作起點：**

```bash
# 第一個指令：建立 Go WASM 分析引擎骨架
cd packages/analysis-engine
mkdir -p cmd/wasm pkg/{analyzer,parsers,bridge,common}
touch cmd/wasm/main.go
touch pkg/analyzer/{workspace,dependency,circular,health}.go
touch pkg/parsers/{npm,yarn,pnpm}.go
touch pkg/bridge/wasm.go
touch pkg/common/{result,errors}.go
touch Makefile
```

**📖 架構文檔位置：**
`_bmad-output/planning-artifacts/architecture.md` (本文檔)

---

**✅ 架構驗證完成**

此文檔已通過全面驗證：

- ✅ 一致性驗證：所有決策協同工作
- ✅ 需求覆蓋驗證：100% (48/48 FR + 17/17 NFR)
- ✅ 實作準備驗證：AI agents 可立即開始
- ✅ 差距分析：無關鍵阻礙
- ✅ 完整性檢查：所有必要元素已定義

**MonoGuard 架構已準備好進入實作階段。** 🚀

---

## Architecture Completion Summary

### Workflow Completion

**Architecture Decision Workflow:** ✅ COMPLETED  
**Total Steps Completed:** 8  
**Date Completed:** 2026-01-12  
**Document Location:** `_bmad-output/planning-artifacts/architecture.md`

---

### Final Architecture Deliverables

#### 📋 Complete Architecture Document

- ✅ 所有架構決策已記錄具體版本
- ✅ 實作模式確保 AI agent 一致性
- ✅ 完整專案結構（檔案層級）
- ✅ 需求到架構映射
- ✅ 驗證確認一致性和完整性

#### 🏗️ Implementation Ready Foundation

- **10 個架構決策** - TanStack Start, Go WASM, Zustand, D3.js, Dexie.js, Tailwind, 等
- **12 種實作模式** - 命名規範、結構模式、通訊模式、流程模式
- **5 個架構元件** - analysis-engine, web, cli, types, ui-components
- **65 個需求完全支援** - 48 FR + 17 NFR = 100% 覆蓋率

#### 📚 AI Agent Implementation Guide

✅ **技術堆疊（版本已鎖定）：**

- TanStack Start 0.34+ (SSG)
- React 19.0.0
- TypeScript 5.9.2
- Go 1.21+ (WASM)
- Zustand 4.4+
- D3.js 7.x
- Dexie.js 5.x
- Tailwind CSS 3.3+

✅ **一致性規則（防止實作衝突）：**

- 跨語言命名統一（TypeScript ↔ Go ↔ JSON）
- WASM Bridge 通訊協定
- 錯誤處理分層策略
- 測試組織標準化

✅ **專案結構（清晰邊界）：**

- 刪除檔案：❌ `apps/api/`, 舊 `apps/frontend/`, 舊 `apps/cli/`
- 重構檔案：♻️ `libs/shared-types/` → `packages/types/`
- 新建檔案：✨ `analysis-engine/`, `web/`, `cli/`, `ui-components/`

✅ **整合模式與通訊標準：**

- WASM ↔ TypeScript: `Result<T>` 統一型別
- App ↔ Storage: Dexie.js `MonoGuardDB`
- React ↔ State: Zustand stores
- CLI ↔ Engine: 直接 Go package 呼叫

---

### Implementation Handoff

#### 給 AI Agents 的指引

此架構文檔是實作 **MonoGuard** 的完整指南。請嚴格遵循所有決策、模式和結構。

**首要實作優先順序：**

```bash
# Phase 0 - MVP 核心功能（0-3 個月）

# Step 1: 基礎設施設定（Week 1）
1. 重構專案結構
   - 刪除 apps/api/, 舊 apps/frontend/, 舊 apps/cli/
   - 建立 packages/analysis-engine/ (Go WASM)
   - 建立 apps/web/ (TanStack Start)
   - 建立 apps/cli/ (Go CLI)
   - 重構 libs/shared-types/ → packages/types/

2. 設定建置工具
   - Nx workspace 配置更新
   - packages/analysis-engine/Makefile
   - apps/web/vite.config.ts
   - GitHub Actions CI/CD

# Step 2: Go WASM 分析引擎（Week 2-4）
1. packages/analysis-engine/pkg/parsers/npm.go (FR5)
2. packages/analysis-engine/pkg/analyzer/workspace.go (FR1)
3. packages/analysis-engine/pkg/analyzer/dependency.go (FR2)
4. packages/analysis-engine/pkg/analyzer/circular.go (FR3)
5. packages/analysis-engine/pkg/analyzer/health.go (FR4)
6. packages/analysis-engine/pkg/bridge/wasm.go (WASM Bridge)

# Step 3: TypeScript WASM 適配器（Week 4-5）
1. packages/types/src/wasmAdapter/index.ts
2. packages/types/src/errors/AnalysisError.ts
3. apps/web/app/lib/wasmLoader.ts

# Step 4: Web UI 核心功能（Week 5-8）
1. apps/web/app/routes/index.tsx (上傳介面)
2. apps/web/app/stores/analysis.ts (Zustand)
3. apps/web/app/components/DependencyGraph/ (D3.js)
4. apps/web/app/lib/persistence.ts (IndexedDB)

# Step 5: CLI 工具（Week 8-10）
1. apps/cli/cmd/monoguard/main.go
2. apps/cli/pkg/commands/{analyze,check}.go
3. apps/cli/pkg/output/formatter.go

# Step 6: 修復建議引擎（Week 10-12）
1. packages/analysis-engine/pkg/rules/root_cause.go
2. packages/analysis-engine/pkg/rules/strategies.go
3. apps/web/app/components/FixSuggestions/

# Step 7: 測試 & 部署（Week 12+）
- Vitest + Go testing 單元測試
- Playwright E2E 測試
- Cloudflare Pages 部署
- CLI GitHub Releases
```

**第一個指令：**

```bash
# 建立 Go WASM 分析引擎骨架
cd packages/analysis-engine
mkdir -p cmd/wasm pkg/{analyzer,parsers,bridge,common}
touch cmd/wasm/main.go
touch pkg/analyzer/{workspace,dependency,circular,health}.go
touch pkg/parsers/{npm,yarn,pnpm}.go
touch pkg/bridge/wasm.go
touch pkg/common/{result,errors}.go
touch Makefile
```

#### 開發流程順序

1. ✅ 使用文檔化的 starter template 初始化專案
2. ✅ 根據架構設定開發環境
3. ✅ 實作核心架構基礎
4. ✅ 遵循既定模式建置功能
5. ✅ 維持與文檔規則的一致性

---

### Quality Assurance Checklist

#### ✅ Architecture Coherence

- [x] 所有決策無衝突協同工作
- [x] 技術選擇相容
- [x] 模式支援架構決策
- [x] 結構與所有選擇對齊

#### ✅ Requirements Coverage

- [x] 所有功能需求已支援（48/48 FR）
- [x] 所有非功能需求已處理（17/17 NFR）
- [x] 跨領域關注點已處理
- [x] 整合點已定義

#### ✅ Implementation Readiness

- [x] 決策具體且可執行
- [x] 模式防止 agent 衝突
- [x] 結構完整且明確
- [x] 提供範例以確保清晰

---

### Project Success Factors

#### 🎯 清晰的決策框架

所有技術選擇都經過協作制定並有明確理由，確保所有利害關係人理解架構方向。

**關鍵決策：**

- ✅ 完整重構為 Client-Only 架構（符合 NFR9-NFR10）
- ✅ Go WASM 提供隱私優先的本地分析
- ✅ TanStack Start SSG 實現零成本基礎設施
- ✅ Zustand 輕量級狀態管理（< 5KB）

#### 🔧 一致性保證

實作模式和規則確保多個 AI agents 會產出相容、一致的程式碼，無縫協作。

**12 種實作模式涵蓋：**

- 跨語言命名規範（TypeScript, Go, JSON）
- WASM Bridge 錯誤處理
- 測試組織結構
- 錯誤分層處理
- 檔案命名規範

#### 📋 完整覆蓋

所有專案需求都獲得架構支援，清楚映射從業務需求到技術實作。

**100% 需求覆蓋：**

- 48 個功能需求 → 映射到具體檔案/目錄
- 17 個非功能需求 → 架構層級決策
- 零差距阻礙實作

#### 🏗️ 堅實基礎

選定的 starter template 和架構模式提供遵循當前最佳實踐的生產就緒基礎。

**架構優勢：**

- 現代技術堆疊（React 19, Go 1.21+, TanStack Start）
- 效能優化（Go WASM < 5s 分析 100 packages）
- 隱私優先（零後端，本地執行）
- 可擴展設計（Phase 0 → Phase 1+ 升級路徑）

---

**Architecture Status:** ✅ **READY FOR IMPLEMENTATION**

**Next Phase:** 使用本文檔記錄的架構決策和模式開始實作。

**Document Maintenance:** 當實作期間做出重大技術決策時更新此架構文檔。

---

**🎉 MonoGuard 架構已完成！**

此文檔經過 8 個步驟的協作式發現流程建立，涵蓋從需求分析到驗證的所有面向。架構已準備好指導一致、高品質的實作工作。

**架構亮點：**

- ✅ 完全對齊 PRD 需求（100% 覆蓋率）
- ✅ 零後端架構（NFR9-NFR10）
- ✅ 隱私優先設計（本地執行 + IndexedDB）
- ✅ 高效能 WASM 分析引擎
- ✅ 零成本基礎設施（Cloudflare Pages 免費層）
- ✅ 清晰的實作路徑（Phase 0 → Phase 1+）

**準備開始建造 MonoGuard！** 🚀
