# CLI: Analyze Command 功能規格

## 概述

`monoguard analyze` 命令提供本地 monorepo 相依性分析功能，可快速檢測技術債務並生成詳細報告。

## 功能細節

### 命令格式

```bash
monoguard analyze [options] [path]
```

### 參數說明

#### 位置參數

- `path` (可選) - 要分析的專案路徑，預設為當前目錄

#### 選項參數

**分析選項：**

- `--focus <type>` - 聚焦特定分析類型
  - 可選值: `dependencies`, `circular`, `architecture`, `all` (預設)
- `--depth <level>` - 分析深度 (1-3，預設 2)
- `--include-dev` - 包含 devDependencies 分析
- `--exclude <pattern>` - 排除特定 package (支援 glob pattern)

**輸出選項：**

- `-o, --output <file>` - 輸出檔案路徑
- `-f, --format <format>` - 報告格式: `json`, `html`, `markdown`, `text` (預設)
- `--verbose` - 詳細輸出模式
- `--quiet` - 靜默模式，僅輸出結果
- `--no-color` - 禁用顏色輸出

**快取選項：**

- `--no-cache` - 禁用快取，強制重新分析
- `--cache-dir <path>` - 自訂快取目錄

### 核心功能

#### 1. Workspace 自動偵測

```typescript
interface WorkspaceDetection {
  // 偵測邏輯
  detectPackageManager(): 'npm' | 'yarn' | 'pnpm' | 'unknown';
  parseWorkspaces(): Package[];
  validateStructure(): ValidationResult;
}
```

**偵測順序：**

1. 檢查 `pnpm-workspace.yaml` → pnpm
2. 檢查 `package.json` 中的 `workspaces` 欄位 → npm/yarn
3. 檢查 `lerna.json` → Lerna
4. 檢查根目錄 `package.json` → 單一專案

#### 2. 相依性分析

**重複相依檢測：**

```typescript
interface DuplicateAnalysis {
  packageName: string;
  versions: string[];
  locations: string[];
  totalSize: string;
  potentialSavings: string;
  recommendation: string;
}
```

**版本衝突檢測：**

```typescript
interface ConflictAnalysis {
  packageName: string;
  conflictType: 'peer' | 'version_range' | 'breaking';
  severity: 'critical' | 'high' | 'medium' | 'low';
  conflictingPackages: string[];
  suggestedResolution: string;
}
```

**未使用相依檢測：**

```typescript
interface UnusedAnalysis {
  packageName: string;
  declaredIn: string;
  confidence: number; // 0-100
  canAutoRemove: boolean;
  reasoning: string;
}
```

#### 3. 循環相依檢測

```typescript
interface CircularDependency {
  cycle: string[]; // ['pkg-a', 'pkg-b', 'pkg-c', 'pkg-a']
  severity: 'critical' | 'high' | 'medium';
  breakPoints: BreakPoint[];
}

interface BreakPoint {
  location: string;
  effort: 'low' | 'medium' | 'high';
  strategy: string;
  codeExample?: string;
}
```

#### 4. 架構違規檢測

```typescript
interface ArchitectureViolation {
  violationType: 'layer_breach' | 'forbidden_import' | 'circular';
  sourcePackage: string;
  targetPackage: string;
  rule: string;
  severity: 'error' | 'warning' | 'info';
  fixSuggestion: string;
}
```

#### 5. 進度顯示

```bash
🔍 Analyzing monorepo...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 100% | ETA: 0s

✓ Workspace detection      (2 packages found)
✓ Dependency resolution    (45 dependencies)
✓ Duplicate detection      (3 duplicates found)
✓ Circular analysis        (1 cycle detected)
✓ Architecture validation  (2 violations)

Analysis completed in 3.2s
```

### 輸出格式

#### Text 格式 (預設)

```
┌─────────────────────────────────────────────┐
│ MonoGuard Analysis Report                   │
├─────────────────────────────────────────────┤
│ Project: my-monorepo                        │
│ Packages: 45                                │
│ Health Score: 72/100                        │
└─────────────────────────────────────────────┘

🔴 Critical Issues (2)
────────────────────────────────────────────
1. Circular Dependency Detected
   libs/auth → libs/user → libs/auth

   Recommendation: Extract shared types to libs/shared
   Effort: 2-4 hours

2. Version Conflict: react
   apps/web: ^18.0.0
   apps/mobile: ^17.0.0

   Recommendation: Upgrade apps/mobile to React 18
   Effort: 4-8 hours

🟡 Warnings (3)
────────────────────────────────────────────
...

💡 Suggestions (5)
────────────────────────────────────────────
...
```

#### JSON 格式

```json
{
  "metadata": {
    "analyzedAt": "2025-01-09T10:30:00Z",
    "version": "0.1.0",
    "projectPath": "/path/to/project",
    "duration": 3.2
  },
  "summary": {
    "packageCount": 45,
    "healthScore": 72,
    "criticalIssues": 2,
    "warnings": 3,
    "suggestions": 5
  },
  "duplicates": [...],
  "conflicts": [...],
  "circular": [...],
  "violations": [...]
}
```

#### HTML 格式

生成互動式 HTML 報告，包含：

- 視覺化圖表
- 可折疊的詳細資訊
- 搜尋和過濾功能
- 可點擊的相依性圖

## User Stories

### User Story 1: 快速本地分析

**As a** 前端開發者
**I want to** 在本地快速執行 monorepo 分析
**So that** 我可以在提交 PR 前發現潛在問題

**Acceptance Criteria:**

- [ ] 命令執行時間 < 5 分鐘 (100 packages)
- [ ] 顯示即時進度條
- [ ] 準確偵測 workspace 類型
- [ ] 輸出易讀的文字報告
- [ ] 支援快取機制，第二次分析 < 30 秒

### User Story 2: CI/CD 整合分析

**As a** DevOps 工程師
**I want to** 在 CI 流程中自動執行分析
**So that** 我可以在部署前攔截架構問題

**Acceptance Criteria:**

- [ ] 支援 `--fail-on-error` 旗標
- [ ] 根據嚴重性設定 exit code (0 或 1)
- [ ] JSON 輸出可被 CI 工具解析
- [ ] 支援靜默模式避免過多日誌
- [ ] 提供分析摘要統計

### User Story 3: 聚焦特定問題分析

**As a** 架構師
**I want to** 只分析特定類型的問題（如循環相依）
**So that** 我可以快速定位特定技術債務

**Acceptance Criteria:**

- [ ] `--focus=circular` 僅執行循環相依分析
- [ ] `--focus=dependencies` 僅執行相依性分析
- [ ] 執行時間縮短 50%+
- [ ] 報告僅包含相關資訊
- [ ] 支援多個 focus 組合

### User Story 4: 產生詳細報告

**As a** 技術主管
**I want to** 生成 HTML 報告分享給團隊
**So that** 團隊可以視覺化了解技術債務狀況

**Acceptance Criteria:**

- [ ] HTML 報告包含互動式圖表
- [ ] 支援匯出為獨立 HTML 檔案
- [ ] 報告包含所有分析細節
- [ ] 支援深色/淺色主題切換
- [ ] 報告大小 < 5MB

## 測試項目

### 單元測試

#### 1. Workspace 偵測測試

```typescript
describe('Workspace Detection', () => {
  test('should detect pnpm workspace', () => {
    // 測試 pnpm-workspace.yaml 偵測
  });

  test('should detect npm workspace', () => {
    // 測試 package.json workspaces 偵測
  });

  test('should detect yarn workspace', () => {
    // 測試 yarn workspaces 偵測
  });

  test('should handle invalid workspace structure', () => {
    // 測試錯誤處理
  });
});
```

#### 2. 相依性分析測試

```typescript
describe('Duplicate Detection', () => {
  test('should detect duplicate dependencies', () => {
    // Given: monorepo with lodash@4.17.21 and lodash@4.17.15
    // When: analyze runs
    // Then: should report 1 duplicate with version details
  });

  test('should calculate bundle impact', () => {
    // 測試 bundle size 計算
  });

  test('should provide migration steps', () => {
    // 測試遷移建議生成
  });
});

describe('Version Conflict Detection', () => {
  test('should detect peer dependency conflicts', () => {
    // 測試 peer dependency 衝突
  });

  test('should detect semver range conflicts', () => {
    // 測試語義化版本衝突
  });
});
```

#### 3. 循環相依測試

```typescript
describe('Circular Dependency Detection', () => {
  test('should detect simple circular dependency', () => {
    // A → B → A
  });

  test('should detect complex circular dependency', () => {
    // A → B → C → D → B
  });

  test('should suggest optimal break points', () => {
    // 測試中斷點建議
  });

  test('should handle no circular dependencies', () => {
    // 測試正常情況
  });
});
```

#### 4. 輸出格式測試

```typescript
describe('Output Formatting', () => {
  test('should generate valid JSON output', () => {
    // 測試 JSON 格式驗證
  });

  test('should generate HTML report', () => {
    // 測試 HTML 生成
  });

  test('should generate markdown report', () => {
    // 測試 Markdown 生成
  });

  test('should support color output', () => {
    // 測試顏色輸出
  });

  test('should support no-color mode', () => {
    // 測試無顏色模式
  });
});
```

### 整合測試

#### 1. 端對端分析流程

```typescript
describe('E2E Analysis Flow', () => {
  test('should analyze real monorepo project', async () => {
    // Given: A real monorepo with known issues
    const result = await runCommand('monoguard analyze ./fixtures/test-repo');

    // Then: Should detect expected issues
    expect(result.duplicates).toHaveLength(3);
    expect(result.circular).toHaveLength(1);
    expect(result.healthScore).toBe(72);
  });

  test('should use cache on second run', async () => {
    // First run
    const firstRun = await runCommand('monoguard analyze');

    // Second run
    const secondRun = await runCommand('monoguard analyze');

    expect(secondRun.duration).toBeLessThan(firstRun.duration * 0.3);
  });
});
```

#### 2. CLI 參數組合測試

```typescript
describe('CLI Options Combinations', () => {
  test('should work with --focus and --format together', () => {
    // monoguard analyze --focus=circular --format=json
  });

  test('should work with --output and --verbose', () => {
    // monoguard analyze --output=report.html --verbose
  });

  test('should respect --exclude pattern', () => {
    // monoguard analyze --exclude="**/test/**"
  });
});
```

### 效能測試

#### 1. 大型 Monorepo 測試

```typescript
describe('Performance Tests', () => {
  test('should analyze 100 packages within 5 minutes', async () => {
    const startTime = Date.now();
    await runCommand('monoguard analyze ./fixtures/large-repo');
    const duration = Date.now() - startTime;

    expect(duration).toBeLessThan(5 * 60 * 1000);
  });

  test('should use < 2GB memory for 500 packages', async () => {
    // 測試記憶體使用
  });
});
```

#### 2. 快取效能測試

```typescript
describe('Cache Performance', () => {
  test('cached analysis should be 10x faster', async () => {
    // First run without cache
    const firstRun = await runCommand('monoguard analyze --no-cache');

    // Second run with cache
    const cachedRun = await runCommand('monoguard analyze');

    expect(cachedRun.duration).toBeLessThan(firstRun.duration / 10);
  });
});
```

### 錯誤處理測試

```typescript
describe('Error Handling', () => {
  test('should handle missing package.json', () => {
    // 測試找不到 package.json 的情況
  });

  test('should handle invalid workspace configuration', () => {
    // 測試無效配置
  });

  test('should handle network errors gracefully', () => {
    // 測試網路錯誤
  });

  test('should provide helpful error messages', () => {
    // 測試錯誤訊息品質
  });
});
```

## 技術實作細節

### 依賴套件

```json
{
  "dependencies": {
    "commander": "^11.0.0",
    "ora": "^6.0.0",
    "chalk": "^5.0.0",
    "glob": "^10.0.0",
    "semver": "^7.5.0",
    "js-yaml": "^4.1.0"
  }
}
```

### 程式碼結構

```
apps/cli/src/
├── commands/
│   ├── analyze.ts          # 主要命令實作
│   ├── options.ts          # 命令選項定義
│   └── validators.ts       # 參數驗證
├── analyzers/
│   ├── workspace.ts        # Workspace 偵測
│   ├── dependencies.ts     # 相依性分析
│   ├── circular.ts         # 循環相依
│   └── architecture.ts     # 架構驗證
├── formatters/
│   ├── text.ts            # Text 輸出
│   ├── json.ts            # JSON 輸出
│   ├── html.ts            # HTML 輸出
│   └── markdown.ts        # Markdown 輸出
└── utils/
    ├── cache.ts           # 快取管理
    ├── progress.ts        # 進度顯示
    └── api-client.ts      # API 客戶端
```

## 完成標準 (Definition of Done)

- [ ] 所有單元測試通過 (覆蓋率 ≥ 90%)
- [ ] 所有整合測試通過
- [ ] 效能測試符合要求
- [ ] 錯誤處理完善
- [ ] 說明文件完整
- [ ] Code review 完成
- [ ] 與後端 API 整合測試通過
- [ ] 支援 macOS, Linux, Windows 三大平台
- [ ] CI/CD 流程驗證通過
