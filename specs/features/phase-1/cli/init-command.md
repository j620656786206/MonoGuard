# CLI: Init Command 功能規格

## 概述

`monoguard init` 命令提供互動式設定精靈，協助用戶快速建立 `.monoguard.yml` 配置檔案，並根據專案結構自動推薦架構規則。

## 功能細節

### 命令格式

```bash
monoguard init [options] [path]
```

### 參數說明

#### 位置參數

- `path` (可選) - 專案路徑，預設為當前目錄

#### 選項參數

**初始化選項：**

- `--template <name>` - 使用預設模板 (`react`, `angular`, `node`, `full-stack`)
- `--interactive` - 互動式模式 (預設)
- `--non-interactive` - 非互動模式，使用預設值
- `--force` - 覆蓋現有配置檔

**配置選項：**

- `--detect` - 自動偵測專案結構並建議規則
- `--minimal` - 產生最小化配置
- `--full` - 產生完整配置（包含所有可選項）

**輸出選項：**

- `-o, --output <file>` - 配置檔輸出路徑 (預設: `.monoguard.yml`)
- `--dry-run` - 預覽配置但不寫入檔案

### 核心功能

#### 1. 專案結構偵測

```typescript
interface ProjectDetection {
  detectPackageManager(): 'npm' | 'yarn' | 'pnpm' | 'unknown';
  detectFramework(): Framework[];
  detectMonorepoStructure(): MonorepoStructure;
  analyzeDirectoryLayout(): DirectoryAnalysis;
}

interface Framework {
  name: string;
  version?: string;
  confidence: number; // 0-100
  indicators: string[];
}

interface MonorepoStructure {
  type: 'nx' | 'lerna' | 'turborepo' | 'custom';
  workspaces: Workspace[];
  structure: 'apps-libs' | 'packages' | 'mixed';
}

interface DirectoryAnalysis {
  apps: string[];
  libs: string[];
  packages: string[];
  shared: string[];
  tools: string[];
}
```

**偵測邏輯：**

```typescript
// 1. 套件管理器偵測
if (exists('pnpm-workspace.yaml')) return 'pnpm';
if (exists('yarn.lock')) return 'yarn';
if (exists('package-lock.json')) return 'npm';

// 2. 框架偵測
const dependencies = readPackageJson().dependencies;
if (dependencies['react']) frameworks.push({ name: 'React', confidence: 100 });
if (dependencies['@angular/core'])
  frameworks.push({ name: 'Angular', confidence: 100 });
if (dependencies['express'])
  frameworks.push({ name: 'Express', confidence: 80 });

// 3. Monorepo 工具偵測
if (exists('nx.json')) return { type: 'nx' };
if (exists('lerna.json')) return { type: 'lerna' };
if (exists('turbo.json')) return { type: 'turborepo' };

// 4. 目錄結構分析
const dirs = glob('**/');
const apps = dirs.filter((d) => d.startsWith('apps/'));
const libs = dirs.filter((d) => d.startsWith('libs/'));
```

#### 2. 互動式配置精靈

```typescript
interface ConfigWizard {
  runInteractive(): Promise<Config>;
  askProjectInfo(): Promise<ProjectInfo>;
  askLayerStructure(): Promise<Layer[]>;
  askRules(): Promise<Rule[]>;
  confirmConfiguration(): Promise<boolean>;
}

interface Question {
  type: 'input' | 'select' | 'multiselect' | 'confirm';
  message: string;
  choices?: string[];
  default?: any;
  validate?: (input: any) => boolean | string;
}
```

**互動流程：**

```bash
$ monoguard init

🎯 MonoGuard Setup Wizard
Let's configure your monorepo architecture rules!

📦 Project Information
? Project name: my-awesome-monorepo
? Package manager: pnpm
? Monorepo tool: Nx

🏗️  Architecture Structure
We detected the following directory structure:
  - apps/ (3 applications)
  - libs/ (12 libraries)

? How would you like to organize layers?
  ❯ Auto-detect from structure (recommended)
    Custom layer definition
    Use template (React/Angular/Node)

📊 Layer Detection Results:
  ✓ Apps Layer: apps/*
  ✓ UI Components: libs/ui/*
  ✓ Business Logic: libs/business/*
  ✓ Shared Utilities: libs/shared/*

? Configure import rules for "Apps Layer":
  ✓ Can import from: libs/*
  ✓ Cannot import from: apps/*

? Configure import rules for "UI Components":
  ✓ Can import from: libs/shared/*
  ✓ Cannot import from: libs/business/*, apps/*

[继续配置其他层...]

🔒 Architecture Rules
? Enable circular dependency detection? Yes
? Detect unused dependencies? Yes
? Enforce version consistency? Yes

? Rule severity for circular dependencies:
  ❯ error (blocks CI/CD)
    warning (allows merge)
    info (informational only)

📝 Configuration Preview:
┌─────────────────────────────────────────────┐
│ Project: my-awesome-monorepo                │
│ Layers: 4                                   │
│ Rules: 6                                    │
│ Package Manager: pnpm                       │
└─────────────────────────────────────────────┘

architecture:
  layers:
    - name: 'Applications'
      pattern: 'apps/*'
      can_import: ['libs/*']
      cannot_import: ['apps/*']
    ...

? Save configuration to .monoguard.yml? Yes

✅ Configuration saved successfully!

Next steps:
  1. Run 'monoguard validate' to check your architecture
  2. Run 'monoguard analyze' for dependency analysis
  3. Add 'monoguard validate' to your CI/CD pipeline
```

#### 3. 模板系統

```typescript
interface Template {
  name: string;
  description: string;
  framework: string[];
  config: Partial<Config>;
}

const templates: Record<string, Template> = {
  react: {
    name: 'React Application',
    description: 'Standard React monorepo with apps and shared libraries',
    framework: ['react'],
    config: {
      architecture: {
        layers: [
          {
            name: 'Applications',
            pattern: 'apps/*',
            can_import: ['libs/*'],
            cannot_import: ['apps/*'],
          },
          {
            name: 'UI Components',
            pattern: 'libs/ui/*',
            can_import: ['libs/shared/*'],
            cannot_import: ['libs/business/*', 'apps/*'],
          },
          {
            name: 'Business Logic',
            pattern: 'libs/business/*',
            can_import: ['libs/shared/*', 'libs/data/*'],
            cannot_import: ['libs/ui/*', 'apps/*'],
          },
          {
            name: 'Shared Utilities',
            pattern: 'libs/shared/*',
            can_import: [],
            cannot_import: ['apps/*', 'libs/ui/*', 'libs/business/*'],
          },
        ],
        rules: [
          {
            name: 'No circular dependencies',
            severity: 'error',
          },
          {
            name: 'UI purity',
            severity: 'error',
            description: 'UI components cannot contain business logic',
          },
        ],
      },
    },
  },

  'full-stack': {
    name: 'Full-Stack Application',
    description: 'Frontend + Backend monorepo structure',
    framework: ['react', 'express', 'nest'],
    config: {
      architecture: {
        layers: [
          {
            name: 'Frontend Apps',
            pattern: 'apps/frontend/*',
            can_import: ['libs/shared/*', 'libs/ui/*'],
            cannot_import: ['apps/backend/*', 'libs/backend/*'],
          },
          {
            name: 'Backend Apps',
            pattern: 'apps/backend/*',
            can_import: ['libs/shared/*', 'libs/backend/*'],
            cannot_import: ['apps/frontend/*', 'libs/ui/*'],
          },
          {
            name: 'Shared Libraries',
            pattern: 'libs/shared/*',
            can_import: [],
            cannot_import: ['apps/*', 'libs/ui/*', 'libs/backend/*'],
          },
        ],
      },
    },
  },
};
```

#### 4. 智慧推薦系統

```typescript
interface RecommendationEngine {
  analyzeStructure(project: ProjectDetection): Recommendation[];
  suggestLayers(structure: MonorepoStructure): LayerSuggestion[];
  suggestRules(frameworks: Framework[]): RuleSuggestion[];
  validateRecommendations(suggestions: Recommendation[]): ValidationResult;
}

interface Recommendation {
  type: 'layer' | 'rule' | 'structure';
  confidence: number;
  reasoning: string;
  suggestion: any;
  examples?: string[];
}
```

**推薦邏輯：**

```typescript
// 範例：推薦層結構
function suggestLayers(analysis: DirectoryAnalysis): LayerSuggestion[] {
  const suggestions: LayerSuggestion[] = [];

  // 檢測到 apps/ 目錄
  if (analysis.apps.length > 0) {
    suggestions.push({
      layer: {
        name: 'Applications',
        pattern: 'apps/*',
        can_import: ['libs/*'],
        cannot_import: ['apps/*'],
      },
      confidence: 95,
      reasoning: 'Detected apps/ directory with multiple applications',
    });
  }

  // 檢測到 UI 組件庫
  if (analysis.libs.some((lib) => lib.includes('ui'))) {
    suggestions.push({
      layer: {
        name: 'UI Components',
        pattern: 'libs/ui/*',
        can_import: ['libs/shared/*'],
        cannot_import: ['libs/business/*'],
      },
      confidence: 90,
      reasoning: 'Detected UI component libraries',
    });
  }

  return suggestions;
}
```

#### 5. 配置驗證與優化

```typescript
interface ConfigOptimizer {
  validateConfig(config: Config): ValidationResult;
  optimizePatterns(patterns: string[]): string[];
  detectConflicts(layers: Layer[]): Conflict[];
  suggestImprovements(config: Config): Improvement[];
}

interface Improvement {
  type: 'performance' | 'clarity' | 'best_practice';
  message: string;
  before: any;
  after: any;
  impact: 'low' | 'medium' | 'high';
}
```

### 輸出格式

#### 生成的 .monoguard.yml

```yaml
# MonoGuard Configuration
# Generated on: 2025-01-09
# Template: React Application

project:
  name: my-awesome-monorepo
  version: 1.0.0

architecture:
  layers:
    - name: 'Applications'
      pattern: 'apps/*'
      description: 'Frontend applications'
      can_import:
        - 'libs/*'
      cannot_import:
        - 'apps/*'

    - name: 'UI Components'
      pattern: 'libs/ui/*'
      description: 'Reusable UI components'
      can_import:
        - 'libs/shared/*'
      cannot_import:
        - 'libs/business/*'
        - 'apps/*'

    - name: 'Business Logic'
      pattern: 'libs/business/*'
      description: 'Core business logic'
      can_import:
        - 'libs/shared/*'
        - 'libs/data/*'
      cannot_import:
        - 'libs/ui/*'
        - 'apps/*'

    - name: 'Shared Utilities'
      pattern: 'libs/shared/*'
      description: 'Common utilities and helpers'
      can_import: []
      cannot_import:
        - 'apps/*'
        - 'libs/ui/*'
        - 'libs/business/*'

  rules:
    - name: 'No circular dependencies'
      severity: 'error'
      description: 'Packages cannot form circular dependencies'

    - name: 'UI component purity'
      severity: 'error'
      description: 'UI components cannot import business logic'

    - name: 'Dependency version consistency'
      severity: 'warning'
      description: 'Same package should use consistent versions'

    - name: 'Maximum dependencies limit'
      severity: 'info'
      max_dependencies: 20
      description: 'Packages should not exceed 20 direct dependencies'

analysis:
  include_dev_dependencies: false
  detect_unused: true
  bundle_impact: true

ci:
  fail_on: 'error'
  output_format: 'json'
```

## User Stories

### User Story 1: 快速初始化配置

**As a** 新用戶
**I want to** 透過互動式精靈快速設定 MonoGuard
**So that** 我可以在 15 分鐘內開始使用

**Acceptance Criteria:**

- [ ] 互動式精靈引導所有必要設定
- [ ] 自動偵測專案結構並提供建議
- [ ] 提供配置預覽
- [ ] 完整流程 < 15 分鐘
- [ ] 生成的配置可直接使用

### User Story 2: 使用模板快速啟動

**As a** React 開發者
**I want to** 使用 React 模板初始化配置
**So that** 我可以遵循最佳實踐而不需要從零開始

**Acceptance Criteria:**

- [ ] `--template=react` 直接應用 React 最佳實踐
- [ ] 模板包含常見的層結構和規則
- [ ] 可以在模板基礎上自訂
- [ ] 提供模板說明文件
- [ ] 支援多種框架模板

### User Story 3: 自動偵測專案結構

**As a** 現有專案維護者
**I want to** 自動偵測現有專案結構並生成配置
**So that** 我不需要手動分析整個專案

**Acceptance Criteria:**

- [ ] 準確偵測目錄結構 (> 90%)
- [ ] 識別常見的分層模式
- [ ] 推薦合適的規則
- [ ] 提供信心度評分
- [ ] 允許手動調整建議

### User Story 4: 非互動模式

**As a** CI/CD 工程師
**I want to** 在自動化腳本中使用 init 命令
**So that** 我可以批次初始化多個專案

**Acceptance Criteria:**

- [ ] `--non-interactive` 使用預設值
- [ ] `--template` 直接應用模板
- [ ] 支援環境變數配置
- [ ] 可腳本化執行
- [ ] 提供錯誤處理和日誌

## 測試項目

### 單元測試

#### 1. 專案偵測測試

```typescript
describe('Project Detection', () => {
  test('should detect pnpm workspace', () => {
    const detector = new ProjectDetector('./fixtures/pnpm-mono');
    const result = detector.detectPackageManager();
    expect(result).toBe('pnpm');
  });

  test('should detect React framework', () => {
    const detector = new ProjectDetector('./fixtures/react-app');
    const frameworks = detector.detectFramework();
    expect(frameworks).toContainEqual(
      expect.objectContaining({ name: 'React', confidence: 100 })
    );
  });

  test('should analyze directory structure', () => {
    const detector = new ProjectDetector('./fixtures/standard-mono');
    const analysis = detector.analyzeDirectoryLayout();

    expect(analysis.apps).toEqual(['apps/web', 'apps/mobile']);
    expect(analysis.libs).toHaveLength(5);
  });
});
```

#### 2. 模板系統測試

```typescript
describe('Template System', () => {
  test('should load React template', () => {
    const template = loadTemplate('react');
    expect(template.name).toBe('React Application');
    expect(template.config.architecture.layers).toHaveLength(4);
  });

  test('should merge template with custom config', () => {
    const template = loadTemplate('react');
    const custom = { project: { name: 'my-project' } };
    const merged = mergeConfigs(template.config, custom);

    expect(merged.project.name).toBe('my-project');
    expect(merged.architecture).toBeDefined();
  });

  test('should validate template structure', () => {
    const template = loadTemplate('full-stack');
    const validation = validateTemplate(template);
    expect(validation.isValid).toBe(true);
  });
});
```

#### 3. 推薦引擎測試

```typescript
describe('Recommendation Engine', () => {
  test('should suggest layers based on structure', () => {
    const analysis = {
      apps: ['apps/web', 'apps/mobile'],
      libs: ['libs/ui', 'libs/business', 'libs/shared'],
    };

    const suggestions = suggestLayers(analysis);
    expect(suggestions).toHaveLength(4); // Apps, UI, Business, Shared
  });

  test('should calculate confidence scores', () => {
    const suggestion = suggestLayer({ apps: ['apps/web'] });
    expect(suggestion.confidence).toBeGreaterThan(90);
  });

  test('should provide reasoning for suggestions', () => {
    const suggestion = suggestLayer({ libs: ['libs/ui/button'] });
    expect(suggestion.reasoning).toContain('UI component');
  });
});
```

#### 4. 配置驗證測試

```typescript
describe('Config Validation', () => {
  test('should validate generated config', () => {
    const config = generateConfig({ template: 'react' });
    const validation = validateConfig(config);
    expect(validation.isValid).toBe(true);
  });

  test('should detect pattern conflicts', () => {
    const config = {
      layers: [
        { pattern: 'libs/*', can_import: [] },
        { pattern: 'libs/ui/*', can_import: ['libs/*'] },
      ],
    };

    const conflicts = detectConflicts(config);
    expect(conflicts).toHaveLength(0); // No conflicts
  });

  test('should suggest improvements', () => {
    const config = generateConfig({ minimal: true });
    const improvements = suggestImprovements(config);
    expect(improvements.length).toBeGreaterThan(0);
  });
});
```

### 整合測試

#### 1. 完整初始化流程

```typescript
describe('E2E Initialization', () => {
  test('should initialize with template', async () => {
    await runCommand('monoguard init --template=react --force');

    const configExists = await fs.exists('.monoguard.yml');
    expect(configExists).toBe(true);

    const config = await loadYaml('.monoguard.yml');
    expect(config.architecture.layers).toHaveLength(4);
  });

  test('should initialize with auto-detection', async () => {
    await runCommand('monoguard init --detect --non-interactive');

    const config = await loadYaml('.monoguard.yml');
    expect(config.architecture.layers.length).toBeGreaterThan(0);
  });

  test('should not overwrite existing config', async () => {
    await fs.writeFile('.monoguard.yml', 'existing: config');

    const result = await runCommand('monoguard init');

    expect(result.exitCode).toBe(1);
    expect(result.stderr).toContain('already exists');
  });

  test('should overwrite with --force flag', async () => {
    await fs.writeFile('.monoguard.yml', 'old: config');

    await runCommand('monoguard init --force --template=react');

    const config = await loadYaml('.monoguard.yml');
    expect(config.old).toBeUndefined();
    expect(config.architecture).toBeDefined();
  });
});
```

#### 2. 互動式模式測試

```typescript
describe('Interactive Mode', () => {
  test('should handle user input', async () => {
    const inputs = [
      'my-project', // project name
      'pnpm', // package manager
      'y', // enable circular detection
      'error', // severity
      'y', // confirm save
    ];

    const result = await runCommandWithInput(
      'monoguard init --interactive',
      inputs
    );

    expect(result.exitCode).toBe(0);
    expect(result.output).toContain('Configuration saved');
  });

  test('should allow cancellation', async () => {
    const inputs = ['my-project', 'n']; // Cancel at confirm

    const result = await runCommandWithInput('monoguard init', inputs);

    expect(result.exitCode).toBe(0);
    expect(result.output).toContain('Cancelled');
  });
});
```

### 效能測試

```typescript
describe('Performance Tests', () => {
  test('should complete initialization within 30 seconds', async () => {
    const startTime = Date.now();
    await runCommand('monoguard init --template=react --non-interactive');
    const duration = Date.now() - startTime;

    expect(duration).toBeLessThan(30000);
  });

  test('should detect large monorepo structure quickly', async () => {
    const startTime = Date.now();
    const detector = new ProjectDetector('./fixtures/large-monorepo');
    await detector.analyzeDirectoryLayout();
    const duration = Date.now() - startTime;

    expect(duration).toBeLessThan(5000);
  });
});
```

## 技術實作細節

### 依賴套件

```json
{
  "dependencies": {
    "commander": "^11.0.0",
    "inquirer": "^9.0.0",
    "js-yaml": "^4.1.0",
    "glob": "^10.0.0",
    "chalk": "^5.0.0",
    "ora": "^6.0.0",
    "validate": "^5.0.0"
  }
}
```

### 程式碼結構

```
apps/cli/src/commands/init/
├── index.ts                 # 主命令入口
├── detector.ts              # 專案偵測
├── wizard.ts                # 互動式精靈
├── templates/
│   ├── react.ts
│   ├── angular.ts
│   ├── node.ts
│   └── full-stack.ts
├── recommender.ts           # 推薦引擎
├── config-generator.ts      # 配置生成
└── validator.ts             # 配置驗證
```

## 完成標準 (Definition of Done)

- [ ] 所有單元測試通過 (覆蓋率 ≥ 90%)
- [ ] 互動式精靈流程完整
- [ ] 至少 4 個模板可用
- [ ] 自動偵測準確度 ≥ 90%
- [ ] 配置驗證完善
- [ ] 使用文件完整
- [ ] 新用戶可在 15 分鐘內完成設定
- [ ] 支援非互動模式
- [ ] 錯誤處理完善
