# CLI: Validate Command 功能規格

## 概述

`monoguard validate` 命令專門用於驗證架構規則，檢查 monorepo 是否符合 `.monoguard.yml` 中定義的分層架構和相依性規則。

## 功能細節

### 命令格式

```bash
monoguard validate [options] [path]
```

### 參數說明

#### 位置參數

- `path` (可選) - 專案路徑，預設為當前目錄

#### 選項參數

**驗證選項：**

- `-c, --config <file>` - 指定配置檔路徑 (預設: `.monoguard.yml`)
- `--rules <rules>` - 僅驗證特定規則 (逗號分隔)
- `--severity <level>` - 最低嚴重性級別 (`error`, `warning`, `info`)
- `--fix` - 自動修復可修復的違規項目

**輸出選項：**

- `-f, --format <format>` - 輸出格式: `text`, `json`, `junit` (預設: `text`)
- `-o, --output <file>` - 輸出檔案路徑
- `--fail-on <severity>` - 當遇到指定嚴重性時失敗 (預設: `error`)
- `--verbose` - 顯示詳細資訊
- `--quiet` - 靜默模式

**CI 模式：**

- `--ci` - CI 模式，優化輸出格式
- `--exit-code` - 根據違規數量設定 exit code

### 核心功能

#### 1. 配置檔驗證

**Schema 驗證：**

```typescript
interface ConfigValidation {
  validateSchema(): ValidationResult;
  validateLayerPatterns(): PatternValidation[];
  validateRuleDefinitions(): RuleValidation[];
  detectConflicts(): Conflict[];
}

interface ValidationResult {
  isValid: boolean;
  errors: ConfigError[];
  warnings: ConfigWarning[];
}
```

**驗證項目：**

- YAML 語法正確性
- Schema 結構完整性
- Glob pattern 有效性
- 規則邏輯一致性
- 循環引用檢測

#### 2. 分層架構驗證

```typescript
interface LayerValidation {
  layerName: string;
  pattern: string;
  matchedPackages: string[];
  violations: LayerViolation[];
}

interface LayerViolation {
  type: 'forbidden_import' | 'missing_import' | 'layer_breach';
  sourcePackage: string;
  targetPackage: string;
  importPath: string;
  filePath: string;
  lineNumber: number;
  rule: string;
  severity: 'error' | 'warning' | 'info';
  canAutoFix: boolean;
  fixSuggestion: string;
}
```

**驗證邏輯：**

```yaml
# 範例配置
architecture:
  layers:
    - name: 'UI Layer'
      pattern: 'libs/ui/*'
      can_import: ['libs/shared/*']
      cannot_import: ['libs/business/*', 'apps/*']
```

驗證時檢查：

1. `libs/ui/button` 是否只 import `libs/shared/*`
2. 是否錯誤 import `libs/business/*`
3. 是否符合禁止規則

#### 3. 相依性規則驗證

```typescript
interface DependencyRule {
  name: string;
  pattern: string;
  allowedDependencies?: string[];
  forbiddenDependencies?: string[];
  maxDependencies?: number;
  severity: 'error' | 'warning' | 'info';
}

interface DependencyViolation {
  packageName: string;
  rule: string;
  violationType: 'forbidden' | 'exceeds_limit' | 'missing_required';
  details: string;
  severity: 'error' | 'warning' | 'info';
}
```

**規則範例：**

```yaml
rules:
  - name: 'No circular dependencies'
    severity: 'error'
    description: 'Packages cannot form circular dependencies'

  - name: 'UI layer purity'
    severity: 'error'
    description: 'UI components cannot import business logic'

  - name: 'Dependency limit'
    severity: 'warning'
    description: 'Packages should not exceed 20 dependencies'
    max_dependencies: 20
```

#### 4. Import 語句分析

```typescript
interface ImportAnalyzer {
  parseImports(filePath: string): ImportStatement[];
  resolveImportPath(importPath: string): string;
  matchAgainstRules(imports: ImportStatement[]): Violation[];
}

interface ImportStatement {
  source: string;
  imported: string[];
  type: 'named' | 'default' | 'namespace' | 'dynamic';
  filePath: string;
  lineNumber: number;
}
```

**支援的 Import 語法：**

```typescript
// Named imports
import { Button, Input } from '@libs/ui';

// Default imports
import React from 'react';

// Namespace imports
import * as utils from '@libs/shared/utils';

// Dynamic imports
const module = await import('@libs/lazy');

// Re-exports
export { Button } from '@libs/ui';
export * from '@libs/shared';
```

#### 5. 自動修復功能

```typescript
interface AutoFixer {
  canFix(violation: Violation): boolean;
  fix(violation: Violation): FixResult;
  preview(violation: Violation): string;
}

interface FixResult {
  success: boolean;
  changes: FileChange[];
  message: string;
}

interface FileChange {
  filePath: string;
  oldContent: string;
  newContent: string;
  diff: string;
}
```

**可自動修復的違規：**

- 移除未使用的 import
- 更新 import 路徑
- 排序 import 語句
- 添加缺少的類型 import

### 輸出格式

#### Text 格式 (預設)

```
┌─────────────────────────────────────────────┐
│ Architecture Validation Report              │
├─────────────────────────────────────────────┤
│ Config: .monoguard.yml                      │
│ Layers: 4                                   │
│ Rules: 6                                    │
└─────────────────────────────────────────────┘

✓ Configuration is valid

🔴 Errors (2)
────────────────────────────────────────────
1. Layer Breach: UI → Business
   libs/ui/button/index.ts:5
   import { validateUser } from '@libs/business/auth';

   ❌ UI layer cannot import from business layer

   💡 Fix: Extract shared validation to libs/shared
   🔧 Auto-fix available: No

2. Forbidden Dependency
   apps/web/src/App.tsx:12
   import Database from 'better-sqlite3';

   ❌ Frontend apps cannot import database libraries

   💡 Fix: Use API calls instead of direct database access
   🔧 Auto-fix available: No

🟡 Warnings (3)
────────────────────────────────────────────
...

💡 Info (1)
────────────────────────────────────────────
...

Summary: 2 errors, 3 warnings, 1 info
Status: ❌ FAILED
```

#### JSON 格式

```json
{
  "configPath": ".monoguard.yml",
  "configValid": true,
  "summary": {
    "totalViolations": 6,
    "errors": 2,
    "warnings": 3,
    "info": 1
  },
  "violations": [
    {
      "id": "layer-breach-001",
      "type": "layer_breach",
      "severity": "error",
      "sourcePackage": "libs/ui/button",
      "targetPackage": "libs/business/auth",
      "filePath": "libs/ui/button/index.ts",
      "lineNumber": 5,
      "importPath": "@libs/business/auth",
      "rule": "UI layer cannot import from business layer",
      "fixSuggestion": "Extract shared validation to libs/shared",
      "canAutoFix": false
    }
  ],
  "exitCode": 1
}
```

#### JUnit XML 格式 (CI 整合)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="MonoGuard Architecture Validation" tests="6" failures="2">
  <testsuite name="Layer Architecture" tests="4" failures="1">
    <testcase name="UI Layer Integrity" classname="architecture.layers">
      <failure message="Layer breach detected">
        libs/ui/button/index.ts:5 imports from forbidden layer
      </failure>
    </testcase>
  </testsuite>
  <testsuite name="Dependency Rules" tests="2" failures="1">
    <testcase name="No Database in Frontend" classname="architecture.dependencies">
      <failure message="Forbidden dependency">
        apps/web imports better-sqlite3
      </failure>
    </testcase>
  </testsuite>
</testsuites>
```

## User Stories

### User Story 1: 開發時驗證架構

**As a** 前端開發者
**I want to** 在開發過程中快速驗證我的修改是否違反架構規則
**So that** 我可以在提交前修正問題

**Acceptance Criteria:**

- [ ] 命令執行時間 < 10 秒
- [ ] 清楚指出違規的檔案和行號
- [ ] 提供具體的修復建議
- [ ] 支援監看模式 (watch mode)
- [ ] 整合到 IDE (VS Code extension)

### User Story 2: CI/CD 架構守門員

**As a** DevOps 工程師
**I want to** 在 CI 流程中強制執行架構規則
**So that** 違反規則的 PR 無法合併

**Acceptance Criteria:**

- [ ] 根據違規設定 exit code (0/1)
- [ ] 輸出 JUnit XML 格式供 CI 工具解析
- [ ] 支援 `--fail-on=error` 只在 error 時失敗
- [ ] 在 PR 中顯示違規詳情
- [ ] 提供違規趨勢追蹤

### User Story 3: 配置檔驗證

**As a** 架構師
**I want to** 驗證 `.monoguard.yml` 配置檔是否正確
**So that** 我可以確保規則定義沒有錯誤

**Acceptance Criteria:**

- [ ] 驗證 YAML 語法
- [ ] 驗證 schema 結構
- [ ] 檢測規則衝突
- [ ] 提供配置範例和模板
- [ ] 支援乾跑模式 (dry-run)

### User Story 4: 自動修復違規

**As a** 開發者
**I want to** 自動修復簡單的架構違規
**So that** 我可以節省手動修改的時間

**Acceptance Criteria:**

- [ ] `--fix` 旗標自動修復可修復的違規
- [ ] 修復前顯示預覽
- [ ] 支援批次修復
- [ ] 保留程式碼格式
- [ ] 產生修復報告

## 測試項目

### 單元測試

#### 1. 配置檔驗證測試

```typescript
describe('Config Validation', () => {
  test('should validate correct YAML syntax', () => {
    const config = loadConfig('valid-config.yml');
    const result = validateConfig(config);
    expect(result.isValid).toBe(true);
  });

  test('should detect invalid layer pattern', () => {
    const config = {
      architecture: {
        layers: [{ name: 'UI', pattern: '[invalid-glob' }],
      },
    };
    const result = validateConfig(config);
    expect(result.errors).toContainEqual(
      expect.objectContaining({ type: 'invalid_pattern' })
    );
  });

  test('should detect circular layer dependencies', () => {
    // Layer A can import B, B can import A
    const config = createCircularLayerConfig();
    const result = validateConfig(config);
    expect(result.errors).toContainEqual(
      expect.objectContaining({ type: 'circular_dependency' })
    );
  });
});
```

#### 2. 分層架構驗證測試

```typescript
describe('Layer Validation', () => {
  test('should detect layer breach', () => {
    // UI layer imports from business layer
    const violations = validateLayer({
      layerName: 'UI',
      pattern: 'libs/ui/*',
      cannot_import: ['libs/business/*'],
    });

    expect(violations).toHaveLength(1);
    expect(violations[0].type).toBe('layer_breach');
  });

  test('should allow valid imports', () => {
    // UI layer imports from shared layer (allowed)
    const violations = validateLayer({
      layerName: 'UI',
      pattern: 'libs/ui/*',
      can_import: ['libs/shared/*'],
    });

    expect(violations).toHaveLength(0);
  });

  test('should match packages correctly', () => {
    const matches = matchLayerPattern('libs/ui/*', [
      'libs/ui/button',
      'libs/ui/input',
      'libs/business/auth',
    ]);

    expect(matches).toEqual(['libs/ui/button', 'libs/ui/input']);
  });
});
```

#### 3. Import 分析測試

```typescript
describe('Import Analysis', () => {
  test('should parse named imports', () => {
    const code = `import { Button, Input } from '@libs/ui';`;
    const imports = parseImports(code);

    expect(imports).toEqual([
      {
        source: '@libs/ui',
        imported: ['Button', 'Input'],
        type: 'named',
      },
    ]);
  });

  test('should resolve import paths', () => {
    const path = resolveImportPath('@libs/ui/button', {
      baseUrl: './libs',
      paths: { '@libs/*': ['*'] },
    });

    expect(path).toBe('libs/ui/button');
  });

  test('should handle dynamic imports', () => {
    const code = `const module = await import('@libs/lazy');`;
    const imports = parseImports(code);

    expect(imports[0].type).toBe('dynamic');
  });
});
```

#### 4. 自動修復測試

```typescript
describe('Auto Fix', () => {
  test('should remove unused imports', () => {
    const code = `
      import { Button } from '@libs/ui';
      import { unused } from '@libs/shared';

      export const MyComponent = () => <Button />;
    `;

    const fixed = autoFix(code, { removeUnused: true });

    expect(fixed).not.toContain('unused');
    expect(fixed).toContain('Button');
  });

  test('should update import paths', () => {
    const code = `import Button from '../../../libs/ui/button';`;
    const fixed = autoFix(code, {
      updatePaths: true,
      useAliases: true,
    });

    expect(fixed).toBe(`import Button from '@libs/ui/button';`);
  });

  test('should sort imports', () => {
    const code = `
      import z from 'z';
      import a from 'a';
      import m from 'm';
    `;

    const fixed = autoFix(code, { sortImports: true });
    const lines = fixed.split('\n').filter((l) => l.trim());

    expect(lines[0]).toContain('a');
    expect(lines[1]).toContain('m');
    expect(lines[2]).toContain('z');
  });
});
```

### 整合測試

#### 1. 完整驗證流程

```typescript
describe('E2E Validation', () => {
  test('should validate entire monorepo', async () => {
    const result = await runCommand('monoguard validate ./fixtures/test-repo');

    expect(result.exitCode).toBe(1); // Has errors
    expect(result.violations).toHaveLength(5);
    expect(result.errors).toHaveLength(2);
  });

  test('should respect severity filter', async () => {
    const result = await runCommand('monoguard validate --fail-on=error');

    // Should only fail on errors, not warnings
    expect(result.exitCode).toBe(0);
  });

  test('should output JUnit XML', async () => {
    const result = await runCommand(
      'monoguard validate --format=junit --output=report.xml'
    );

    const xml = await fs.readFile('report.xml', 'utf-8');
    expect(xml).toContain('<?xml version="1.0"');
    expect(xml).toContain('<testsuites');
  });
});
```

#### 2. CI 整合測試

```typescript
describe('CI Integration', () => {
  test('should work in GitHub Actions', async () => {
    // Simulate GitHub Actions environment
    process.env.CI = 'true';
    process.env.GITHUB_ACTIONS = 'true';

    const result = await runCommand('monoguard validate --ci');

    expect(result.output).toContain('::error');
    expect(result.output).toContain('::warning');
  });

  test('should generate annotations', async () => {
    const result = await runCommand('monoguard validate --ci --format=github');

    expect(result.output).toMatch(/::error file=.*,line=\d+::/);
  });
});
```

### 效能測試

```typescript
describe('Performance Tests', () => {
  test('should validate 1000 files within 10 seconds', async () => {
    const startTime = Date.now();
    await runCommand('monoguard validate ./fixtures/large-repo');
    const duration = Date.now() - startTime;

    expect(duration).toBeLessThan(10000);
  });

  test('should use incremental validation', async () => {
    // First run
    await runCommand('monoguard validate');

    // Change one file
    await modifyFile('libs/ui/button/index.ts');

    // Second run should be faster
    const startTime = Date.now();
    await runCommand('monoguard validate');
    const duration = Date.now() - startTime;

    expect(duration).toBeLessThan(2000);
  });
});
```

### 錯誤處理測試

```typescript
describe('Error Handling', () => {
  test('should handle missing config file', async () => {
    const result = await runCommand('monoguard validate --config=missing.yml');

    expect(result.exitCode).toBe(1);
    expect(result.stderr).toContain('Config file not found');
  });

  test('should handle invalid YAML', async () => {
    const result = await runCommand('monoguard validate --config=invalid.yml');

    expect(result.exitCode).toBe(1);
    expect(result.stderr).toContain('Invalid YAML syntax');
  });

  test('should provide helpful suggestions', async () => {
    const result = await runCommand('monoguard validate');

    expect(result.output).toContain('💡 Fix:');
    expect(result.output).toContain('Suggestion:');
  });
});
```

## 技術實作細節

### 依賴套件

```json
{
  "dependencies": {
    "commander": "^11.0.0",
    "js-yaml": "^4.1.0",
    "@typescript-eslint/parser": "^6.0.0",
    "glob": "^10.0.0",
    "chalk": "^5.0.0",
    "ora": "^6.0.0",
    "fast-xml-parser": "^4.0.0"
  }
}
```

### 程式碼結構

```
apps/cli/src/commands/validate/
├── index.ts              # 主命令入口
├── config-validator.ts   # 配置驗證
├── layer-validator.ts    # 分層驗證
├── import-analyzer.ts    # Import 分析
├── auto-fixer.ts         # 自動修復
├── formatters/
│   ├── text.ts
│   ├── json.ts
│   └── junit.ts
└── rules/
    ├── layer-rules.ts
    ├── dependency-rules.ts
    └── circular-rules.ts
```

## 完成標準 (Definition of Done)

- [ ] 所有單元測試通過 (覆蓋率 ≥ 90%)
- [ ] 所有整合測試通過
- [ ] 支援所有主要 import 語法
- [ ] JUnit XML 輸出格式正確
- [ ] 自動修復功能完整
- [ ] 錯誤訊息清晰有幫助
- [ ] CI/CD 整合文件完整
- [ ] 與 GitHub Actions 整合測試通過
- [ ] 效能符合要求 (< 10 秒驗證 1000 檔案)
