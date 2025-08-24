# MonoGuard Branch Naming and Development Workflow Strategy

## Executive Summary

This document establishes a comprehensive branch naming convention and development workflow strategy for the MonoGuard monorepo project. MonoGuard is a sophisticated multi-language project (Go backend API, Next.js frontend, Node.js CLI) with complex architectural components for dependency analysis, architecture validation, and project health monitoring.

## Project Architecture Overview

### Core Applications
- **API Service** (`apps/api/`): Go backend with Gin framework, PostgreSQL, Redis
- **Frontend Application** (`apps/frontend/`): Next.js 15 with React 19, Tailwind CSS, Radix UI
- **CLI Tool** (`apps/cli/`): Node.js command-line interface
- **E2E Tests** (`apps/frontend-e2e/`): Playwright end-to-end testing

### Shared Libraries
- **shared-types** (`libs/shared-types/`): TypeScript type definitions and API contracts
- **ui** (`libs/ui/`): Shared UI components with Radix UI and Tailwind CSS

### Key Service Components (Go Backend)
- **Dependency Analyzer**: Package.json parsing, dependency tree resolution, duplicate detection
- **Architecture Validator**: Layer validation, circular dependency detection
- **Project Manager**: Project lifecycle management, analysis orchestration
- **Health Monitor**: System health checks and performance monitoring

## Branch Naming Convention

### Format Structure
```
<type>/<scope>/<description>
```

### Branch Types
- `feature/` - New functionality or enhancements
- `fix/` - Bug fixes and patches
- `hotfix/` - Critical production fixes
- `refactor/` - Code refactoring without functional changes
- `docs/` - Documentation-only changes
- `test/` - Test-related changes
- `chore/` - Maintenance tasks (dependencies, build configs)
- `perf/` - Performance improvements
- `security/` - Security-related changes

### Scope Definitions

#### Application Scopes
- `api` - Go backend API service
- `frontend` - Next.js frontend application
- `cli` - Node.js CLI tool
- `e2e` - End-to-end testing
- `shared-types` - Shared TypeScript types
- `ui` - Shared UI components library

#### Feature Domain Scopes (API Services)
- `dependency-analyzer` - Dependency analysis engine
- `architecture-validator` - Architecture validation service
- `project-manager` - Project management service
- `health-monitor` - Health monitoring service
- `circular-detector` - Circular dependency detection
- `layer-validator` - Layer architecture validation
- `package-parser` - Package.json parsing service
- `tree-resolver` - Dependency tree resolution
- `duplicate-detector` - Duplicate dependency detection
- `unused-detector` - Unused dependency detection

#### Infrastructure Scopes
- `database` - Database-related changes (migrations, models)
- `redis` - Redis caching implementation
- `middleware` - HTTP middleware components
- `auth` - Authentication and authorization
- `config` - Configuration management
- `deployment` - Deployment and infrastructure
- `docker` - Docker containerization
- `ci-cd` - Continuous integration/deployment

#### Frontend Feature Scopes
- `dashboard` - Main dashboard interface
- `project-view` - Project detail views
- `analysis-view` - Analysis result visualization
- `health-dashboard` - Health score dashboard
- `dependency-graph` - D3.js dependency visualization
- `architecture-graph` - Architecture diagram components
- `config-manager` - Configuration management UI
- `report-export` - Report generation and export
- `user-management` - User authentication UI

### Branch Naming Examples

#### Feature Branches
```bash
# API Backend Features
feature/api/dependency-analyzer-engine
feature/api/circular-dependency-detection
feature/api/package-json-parser-enhancement
feature/api/redis-caching-layer
feature/api/health-check-endpoints

# Frontend Features
feature/frontend/dashboard-health-scores
feature/frontend/dependency-graph-visualization
feature/frontend/project-management-ui
feature/frontend/analysis-report-export
feature/frontend/architecture-validator-ui

# CLI Features
feature/cli/local-analysis-command
feature/cli/configuration-validation
feature/cli/report-generation

# Shared Library Features
feature/shared-types/api-contract-updates
feature/ui/analysis-visualization-components

# Cross-Application Features
feature/full-stack/project-analysis-workflow
feature/full-stack/user-authentication-system
```

#### Bug Fix Branches
```bash
# Component-Specific Fixes
fix/api/dependency-resolver-memory-leak
fix/frontend/dashboard-loading-performance
fix/cli/config-file-parsing-error
fix/database/migration-rollback-issue

# Service-Specific Fixes
fix/dependency-analyzer/version-conflict-detection
fix/architecture-validator/layer-rule-validation
fix/circular-detector/false-positive-reduction
```

#### Hotfix Branches
```bash
hotfix/api/security-vulnerability-patch
hotfix/frontend/critical-dashboard-crash
hotfix/database/connection-pool-exhaustion
```

#### Maintenance Branches
```bash
# Dependency Updates
chore/api/go-module-updates
chore/frontend/npm-dependency-upgrade
chore/shared-types/type-definition-updates

# Configuration Changes
chore/deployment/zeabur-config-update
chore/ci-cd/github-actions-optimization
chore/docker/multi-stage-build-optimization
```

## Development Workflow Guidelines

### Branch Creation Rules - 前後端嚴格分離策略

#### 1. 前後端分離開發原則
**重要：前後端程式碼必須分開到不同分支，不能混合開發**

#### 2. Backend (API) Development - 後端專用分支
```bash
# 後端服務開發 - 只能修改 apps/api/ 相關檔案
git checkout -b feature/api/dependency-analyzer-engine
git checkout -b feature/api/circular-dependency-detection
git checkout -b feature/api/health-monitoring-service
git checkout -b feature/api/authentication-middleware
git checkout -b feature/database/analysis-result-models
git checkout -b feature/database/user-management-schema
```

#### 3. Frontend Development - 前端專用分支
```bash
# 前端元件開發 - 只能修改 apps/frontend/ 相關檔案
git checkout -b feature/frontend/dependency-graph-visualization
git checkout -b feature/frontend/project-dashboard-layout
git checkout -b feature/frontend/user-authentication-ui
git checkout -b feature/frontend/analysis-report-components
git checkout -b feature/frontend/health-dashboard-charts
```

#### 4. Shared Libraries Development - 共享函式庫分支
```bash
# UI 元件庫開發 - 只能修改 libs/ui/ 相關檔案
git checkout -b feature/ui/analysis-result-cards
git checkout -b feature/ui/dependency-graph-components
git checkout -b feature/ui/dashboard-layout-components

# 共享類型定義 - 只能修改 libs/shared-types/ 相關檔案
git checkout -b feature/shared-types/api-contract-updates
git checkout -b feature/shared-types/analysis-result-types
```

#### 5. Integration Branches - 整合分支（特殊情況）
```bash
# 只有在必須同時修改前後端 API 契約時才使用
git checkout -b integration/full-stack/new-analysis-endpoints
git checkout -b integration/full-stack/authentication-system
git checkout -b integration/full-stack/websocket-real-time-updates

# 整合分支使用規則：
# 1. 先完成後端 API 開發並測試
# 2. 再進行前端對接開發
# 3. 最後進行整合測試
```

#### 6. Bug Fixes - 按前後端分離
```bash
# 後端 Bug 修復 - 只修改 apps/api/ 相關檔案
git checkout -b fix/api/dependency-analyzer-memory-leak
git checkout -b fix/api/architecture-validator-false-positive
git checkout -b fix/database/connection-timeout-handling
git checkout -b fix/redis/cache-invalidation-logic

# 前端 Bug 修復 - 只修改 apps/frontend/ 相關檔案  
git checkout -b fix/frontend/chart-rendering-performance
git checkout -b fix/frontend/dashboard-loading-state
git checkout -b fix/frontend/dependency-graph-layout

# 共享函式庫 Bug 修復
git checkout -b fix/ui/button-component-accessibility
git checkout -b fix/shared-types/api-response-typing
```

#### 7. 禁止跨元件變更分支
```bash
# ❌ 禁止：不能同時修改前後端
# git checkout -b feature/full-stack/enhanced-health-scoring  # 禁止！
# git checkout -b refactor/full-stack/error-handling-standardization  # 禁止！

# ✅ 正確：分開處理
git checkout -b feature/api/enhanced-health-scoring-backend
git checkout -b feature/frontend/enhanced-health-scoring-ui

# 如果必須協調，使用 integration/ 分支並遵循嚴格流程
git checkout -b integration/api-frontend/health-scoring-coordination
```

### Commit Message Convention

#### Format
```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

#### Examples
```bash
# Feature commits
feat(api): implement dependency tree resolver service
feat(frontend): add interactive dependency graph visualization
feat(cli): add project configuration validation command

# Bug fix commits
fix(dependency-analyzer): resolve memory leak in package parser
fix(frontend): correct dashboard loading state management
fix(database): handle connection pool exhaustion gracefully

# Performance commits
perf(api): optimize dependency analysis algorithm performance
perf(frontend): implement virtual scrolling for large dependency lists

# Refactor commits
refactor(services): standardize error handling across all services
refactor(frontend): consolidate API client configuration
```

### Merge Strategy Guidelines

#### 1. Feature Branch Merging
- **Use Squash and Merge** for feature branches to maintain clean history
- **Require PR review** from at least one team member
- **Ensure all tests pass** before merging
- **Update documentation** if API contracts change

#### 2. Branch Protection Rules
```yaml
# .github/branch-protection.yml
main:
  required_status_checks:
    - "build-api"
    - "build-frontend" 
    - "build-cli"
    - "test-unit"
    - "test-integration"
    - "test-e2e"
    - "lint-check"
    - "type-check"
  require_branches_to_be_up_to_date: true
  required_pull_request_reviews:
    required_approving_review_count: 1
  restrictions:
    users: []
    teams: ["core-developers"]
```

#### 3. Release Branch Strategy
```bash
# Create release branch
git checkout -b release/v0.2.0

# Prepare release
# - Update version numbers
# - Update changelog
# - Final testing

# Merge to main and tag
git checkout main
git merge --no-ff release/v0.2.0
git tag -a v0.2.0 -m "Release version 0.2.0"
```

## Surgical Task Executor Integration

### Surgical Task Executor - 前後端分離自動分支選擇

surgical-task-executor 必須嚴格遵循前後端分離原則，使用以下邏輯進行分支選擇：

#### 1. 前後端分離檔案檢測規則
```typescript
interface BranchMappingRule {
  pattern: RegExp;
  branchPrefix: string;
  scope: string;
  category: 'backend' | 'frontend' | 'shared' | 'infrastructure';
}

const BRANCH_MAPPING_RULES: BranchMappingRule[] = [
  // ===== BACKEND ONLY (後端專用) =====
  { pattern: /apps\/api\//, branchPrefix: 'feature', scope: 'api', category: 'backend' },
  { pattern: /apps\/api\/internal\/services\/dependency_analyzer/, branchPrefix: 'feature', scope: 'api/dependency-analyzer', category: 'backend' },
  { pattern: /apps\/api\/internal\/services\/circular_detector/, branchPrefix: 'feature', scope: 'api/circular-detector', category: 'backend' },
  { pattern: /apps\/api\/internal\/services\/layer_validator/, branchPrefix: 'feature', scope: 'api/layer-validator', category: 'backend' },
  { pattern: /apps\/api\/internal\/handlers\//, branchPrefix: 'feature', scope: 'api/handlers', category: 'backend' },
  { pattern: /apps\/api\/internal\/models\//, branchPrefix: 'feature', scope: 'database/models', category: 'backend' },
  { pattern: /apps\/api\/internal\/repository\//, branchPrefix: 'feature', scope: 'database/repository', category: 'backend' },
  { pattern: /apps\/api\/pkg\/database\//, branchPrefix: 'feature', scope: 'database/pkg', category: 'backend' },
  { pattern: /apps\/api\/internal\/middleware\//, branchPrefix: 'feature', scope: 'api/middleware', category: 'backend' },
  { pattern: /apps\/api\/go\.mod|apps\/api\/go\.sum/, branchPrefix: 'chore', scope: 'api/dependencies', category: 'backend' },
  
  // ===== FRONTEND ONLY (前端專用) =====
  { pattern: /apps\/frontend\//, branchPrefix: 'feature', scope: 'frontend', category: 'frontend' },
  { pattern: /apps\/frontend\/src\/app\/dashboard/, branchPrefix: 'feature', scope: 'frontend/dashboard', category: 'frontend' },
  { pattern: /apps\/frontend\/src\/components\//, branchPrefix: 'feature', scope: 'frontend/components', category: 'frontend' },
  { pattern: /apps\/frontend\/src\/lib\/api/, branchPrefix: 'feature', scope: 'frontend/api-client', category: 'frontend' },
  { pattern: /apps\/frontend\/src\/hooks\//, branchPrefix: 'feature', scope: 'frontend/hooks', category: 'frontend' },
  { pattern: /apps\/frontend\/src\/styles\//, branchPrefix: 'feature', scope: 'frontend/styles', category: 'frontend' },
  { pattern: /apps\/frontend\/package\.json/, branchPrefix: 'chore', scope: 'frontend/dependencies', category: 'frontend' },
  
  // ===== SHARED LIBRARIES (共享函式庫) =====
  { pattern: /libs\/shared-types\//, branchPrefix: 'feature', scope: 'shared-types', category: 'shared' },
  { pattern: /libs\/ui\//, branchPrefix: 'feature', scope: 'ui', category: 'shared' },
  
  // ===== CLI (獨立元件) =====
  { pattern: /apps\/cli\//, branchPrefix: 'feature', scope: 'cli', category: 'shared' },
  
  // ===== INFRASTRUCTURE (基礎設施) =====
  { pattern: /docker/, branchPrefix: 'chore', scope: 'docker', category: 'infrastructure' },
  { pattern: /\.github\/workflows/, branchPrefix: 'chore', scope: 'ci-cd', category: 'infrastructure' },
  { pattern: /pnpm-workspace\.yaml|package\.json$/, branchPrefix: 'chore', scope: 'workspace', category: 'infrastructure' },
];
```

#### 2. Task Type Detection
```typescript
interface TaskTypeRule {
  keywords: string[];
  branchType: string;
}

const TASK_TYPE_RULES: TaskTypeRule[] = [
  { keywords: ['fix', 'bug', 'error', 'issue'], branchType: 'fix' },
  { keywords: ['feature', 'add', 'implement', 'create'], branchType: 'feature' },
  { keywords: ['refactor', 'restructure', 'reorganize'], branchType: 'refactor' },
  { keywords: ['performance', 'optimize', 'speed'], branchType: 'perf' },
  { keywords: ['test', 'testing', 'spec'], branchType: 'test' },
  { keywords: ['docs', 'documentation'], branchType: 'docs' },
  { keywords: ['update', 'upgrade', 'dependency'], branchType: 'chore' },
  { keywords: ['security', 'vulnerability'], branchType: 'security' },
];
```

#### 3. 前後端分離分支命名演算法
```typescript
function generateBranchName(task: Task): string {
  const taskType = detectTaskType(task.description);
  const categoryScope = detectCategoryAndScope(task.affectedFiles);
  const description = sanitizeDescription(task.description);
  
  return `${taskType}/${categoryScope.scope}/${description}`;
}

function detectCategoryAndScope(files: string[]): { category: string, scope: string } {
  const categoryWeights = new Map<string, number>();
  const scopeWeights = new Map<string, number>();
  
  for (const file of files) {
    for (const rule of BRANCH_MAPPING_RULES) {
      if (rule.pattern.test(file)) {
        // Track categories
        const currentCategoryWeight = categoryWeights.get(rule.category) || 0;
        categoryWeights.set(rule.category, currentCategoryWeight + 1);
        
        // Track scopes within categories
        const currentScopeWeight = scopeWeights.get(rule.scope) || 0;
        scopeWeights.set(rule.scope, currentScopeWeight + 1);
      }
    }
  }
  
  // 前後端分離檢查 - 禁止混合
  const categories = Array.from(categoryWeights.keys());
  const hasBackend = categories.includes('backend');
  const hasFrontend = categories.includes('frontend');
  
  if (hasBackend && hasFrontend) {
    throw new Error(`
      ❌ 違反前後端分離原則！
      檔案同時包含前端和後端修改：
      - 後端檔案: ${files.filter(f => /apps\/api\//.test(f))}
      - 前端檔案: ${files.filter(f => /apps\/frontend\//.test(f))}
      
      ✅ 請分別建立分支：
      - feature/api/[功能名稱] - 用於後端修改
      - feature/frontend/[功能名稱] - 用於前端修改
    `);
  }
  
  // 確定主要範疇
  if (categoryWeights.size === 0) {
    return { category: 'general', scope: 'general' };
  }
  
  const dominantScope = Array.from(scopeWeights.entries())
    .sort((a, b) => b[1] - a[1])[0][0];
    
  const dominantCategory = Array.from(categoryWeights.entries())
    .sort((a, b) => b[1] - a[1])[0][0];
    
  return { category: dominantCategory, scope: dominantScope };
}

function detectTaskType(description: string): string {
  const lowerDesc = description.toLowerCase();
  
  for (const rule of TASK_TYPE_RULES) {
    if (rule.keywords.some(keyword => lowerDesc.includes(keyword))) {
      return rule.branchType;
    }
  }
  
  return 'feature'; // default
}
```

### 前後端分離工作流程決策矩陣

#### 後端專用變更 (Backend Only)
| 檔案模式 | 分支類型 | 範疇 | 範例 |
|----------|----------|------|------|
| `apps/api/internal/services/dependency_analyzer.go` | `feature` | `api/dependency-analyzer` | `feature/api/dependency-analyzer-optimization` |
| `apps/api/internal/handlers/project.go` | `feature` | `api/handlers` | `feature/api/project-endpoints` |
| `apps/api/pkg/database/` | `feature` | `database/pkg` | `feature/database/connection-pooling` |
| `apps/api/internal/models/` | `feature` | `database/models` | `feature/database/analysis-result-models` |

#### 前端專用變更 (Frontend Only)
| 檔案模式 | 分支類型 | 範疇 | 範例 |
|----------|----------|------|------|
| `apps/frontend/src/components/dashboard/` | `feature` | `frontend/dashboard` | `feature/frontend/dashboard-health-metrics` |
| `apps/frontend/src/lib/api/` | `feature` | `frontend/api-client` | `feature/frontend/api-integration` |
| `apps/frontend/src/app/` | `feature` | `frontend` | `feature/frontend/routing-optimization` |
| `apps/frontend/src/components/` | `feature` | `frontend/components` | `feature/frontend/chart-components` |

#### 共享函式庫變更 (Shared Libraries)
| 檔案模式 | 分支類型 | 範疇 | 範例 |
|----------|----------|------|------|
| `libs/shared-types/src/api.ts` | `feature` | `shared-types` | `feature/shared-types/analysis-result-types` |
| `libs/ui/src/components/` | `feature` | `ui` | `feature/ui/analysis-visualization-components` |
| `apps/cli/src/commands/` | `feature` | `cli` | `feature/cli/enhanced-analysis-command` |

#### Bug 修復 (分類別處理)
| 問題類型 | 分支類型 | 範疇 | 範例 |
|----------|----------|------|------|
| 後端記憶體洩漏 | `fix` | `api/dependency-analyzer` | `fix/api/dependency-analyzer-memory-leak` |
| 前端渲染問題 | `fix` | `frontend/components` | `fix/frontend/dashboard-rendering-performance` |
| CLI 配置解析 | `fix` | `cli` | `fix/cli/yaml-configuration-parsing` |
| 資料庫遷移 | `fix` | `database/models` | `fix/database/migration-constraint-error` |

#### ❌ 禁止的跨元件變更
| 禁止情況 | 原因 | 建議做法 |
|----------|------|----------|
| `API + Frontend` 同時修改 | 違反前後端分離 | 分別建立 `feature/api/xxx` 和 `feature/frontend/xxx` |
| `API + Database + Types` 同時修改 | 過於複雜 | 先做 `feature/database/xxx`，再做 `feature/api/xxx` |
| `All Apps` 同時修改 | 風險過高 | 按元件分階段處理 |

#### ✅ 特殊整合情況 (僅在必要時使用)
| 整合類型 | 分支類型 | 範疇 | 使用時機 |
|----------|----------|------|----------|
| API 契約變更 | `integration` | `api-contract` | 當 API 介面變更需要前後端協調時 |
| 認證系統 | `integration` | `authentication` | 涉及前後端認證流程變更時 |
| WebSocket 即時更新 | `integration` | `real-time` | 需要前後端即時通訊功能時 |

## Best Practices for Development

### 1. Branch Lifecycle Management
- **Create branches** from latest `main` branch
- **Keep branches focused** on single features or fixes
- **Merge branches promptly** after completion
- **Delete merged branches** to keep repository clean

### 2. Testing Requirements
- **Unit tests** required for all service components
- **Integration tests** for API endpoints
- **E2E tests** for critical user workflows
- **Performance tests** for analysis algorithms

### 3. Code Review Guidelines
- **API changes** require architecture review
- **Database changes** require migration review
- **Frontend changes** require UI/UX review
- **Security changes** require security team review

### 4. Documentation Updates
- Update API documentation for endpoint changes
- Update component documentation for UI changes
- Update CLI help text for command changes
- Update architecture diagrams for structural changes

## Integration with MonoGuard Features

### Self-Analysis Integration
Since MonoGuard analyzes itself as a dogfooding example, the branch strategy supports:

1. **Continuous Self-Analysis**: Each feature branch triggers MonoGuard analysis of itself
2. **Architecture Validation**: New branches validate against MonoGuard's own architecture rules
3. **Performance Benchmarking**: Branch changes measure impact on MonoGuard's own performance
4. **Real-World Testing**: Complex monorepo structure provides authentic testing scenarios

### CI/CD Pipeline Integration
```yaml
# .github/workflows/branch-analysis.yml
name: Branch Analysis
on:
  pull_request:
    branches: [main]

jobs:
  self-analysis:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run MonoGuard Analysis
        run: |
          npm run build:cli
          ./dist/monoguard analyze --output=github-action
      - name: Validate Architecture Rules
        run: |
          ./dist/monoguard validate-architecture --config=.monoguard.yml
```

This comprehensive branch strategy ensures systematic development, clear component boundaries, and seamless integration with MonoGuard's self-validating architecture.