# CLI: CI/CD Integration 功能規格

## 概述

整合 MonoGuard CLI 到 CI/CD 流程，提供自動化架構驗證、技術債務檢測和 PR 檢查功能。

## 功能細節

### 支援的 CI/CD 平台

#### 1. GitHub Actions

```yaml
name: MonoGuard Architecture Check

on: [pull_request, push]

jobs:
  architecture-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install MonoGuard
        run: npm install -g monoguard

      - name: Run Architecture Validation
        run: monoguard validate --ci --fail-on=error

      - name: Run Dependency Analysis
        run: monoguard analyze --format=json --output=analysis.json

      - name: Upload Analysis Results
        uses: actions/upload-artifact@v3
        with:
          name: monoguard-analysis
          path: analysis.json

      - name: Comment PR
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const analysis = JSON.parse(fs.readFileSync('analysis.json'));
            const body = `## MonoGuard Analysis Results

            - Health Score: ${analysis.healthScore}/100
            - Critical Issues: ${analysis.summary.criticalIssues}
            - Warnings: ${analysis.summary.warnings}
            `;
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: body
            });
```

#### 2. GitLab CI

```yaml
monoguard:
  stage: test
  image: node:18
  script:
    - npm install -g monoguard
    - monoguard validate --ci --format=junit --output=report.xml
    - monoguard analyze --format=json
  artifacts:
    reports:
      junit: report.xml
    paths:
      - analysis.json
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
```

#### 3. Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('Architecture Check') {
            steps {
                sh 'npm install -g monoguard'
                sh 'monoguard validate --ci --fail-on=error'
                sh 'monoguard analyze --output=analysis.json'
            }
        }
    }
    post {
        always {
            junit 'monoguard-report.xml'
            archiveArtifacts artifacts: 'analysis.json'
        }
    }
}
```

### CI 模式特性

#### 1. 優化的輸出格式

```typescript
interface CIOutput {
  // GitHub Actions Annotations
  github: {
    annotations: Annotation[];
    summary: string;
    conclusion: 'success' | 'failure' | 'neutral';
  };

  // GitLab CI Report
  gitlab: {
    junit: string;
    codequality: CodeQualityReport[];
  };

  // Generic JSON
  json: {
    violations: Violation[];
    summary: Summary;
    exitCode: number;
  };
}

interface Annotation {
  path: string;
  start_line: number;
  end_line: number;
  annotation_level: 'notice' | 'warning' | 'failure';
  message: string;
  title: string;
}
```

#### 2. 環境變數支援

```bash
# GitHub Actions
MONOGUARD_CI=true
MONOGUARD_GITHUB_TOKEN=${{ secrets.GITHUB_TOKEN }}
MONOGUARD_PR_NUMBER=${{ github.event.pull_request.number }}

# GitLab CI
MONOGUARD_CI=true
MONOGUARD_GITLAB_TOKEN=$CI_JOB_TOKEN
MONOGUARD_MR_IID=$CI_MERGE_REQUEST_IID

# Generic
MONOGUARD_FAIL_ON=error
MONOGUARD_OUTPUT_FORMAT=json
MONOGUARD_CACHE_DIR=/tmp/monoguard-cache
```

#### 3. 增量分析

```typescript
interface IncrementalAnalysis {
  getChangedFiles(): string[];
  analyzeOnlyChanged(): AnalysisResult;
  compareWithBaseBranch(base: string): Diff;
}

// 僅分析變更的檔案
const changedFiles = getChangedFiles();
const violations = analyzeFiles(changedFiles);

// 與 base branch 比較
const diff = compareWithBaseBranch('main');
console.log(`New violations: ${diff.newViolations.length}`);
console.log(`Resolved violations: ${diff.resolvedViolations.length}`);
```

## User Stories

### User Story 1: GitHub PR 自動檢查

**As a** 開發者
**I want to** 在 PR 中自動執行架構檢查
**So that** 違規程式碼無法被合併

**Acceptance Criteria:**

- [ ] PR 創建時自動觸發檢查
- [ ] 在 PR 中顯示檢查結果
- [ ] 違規時阻止合併
- [ ] 提供詳細的錯誤註解
- [ ] 顯示健康度趨勢

### User Story 2: 增量分析優化

**As a** DevOps 工程師
**I want to** 只分析變更的檔案
**So that** CI 執行時間更短

**Acceptance Criteria:**

- [ ] 自動偵測變更檔案
- [ ] 僅分析相關的 packages
- [ ] 執行時間減少 70%+
- [ ] 結果仍然準確
- [ ] 支援多個 base branch

### User Story 3: 失敗門檻設定

**As a** 團隊主管
**I want to** 自訂 CI 失敗門檻
**So that** 我可以控制何時阻止部署

**Acceptance Criteria:**

- [ ] `--fail-on=error` 僅在 error 時失敗
- [ ] `--threshold=80` 健康度低於 80 時失敗
- [ ] 支援自訂規則組合
- [ ] 提供豁免機制
- [ ] 記錄所有決策

## 測試項目

### 整合測試

#### 1. GitHub Actions 整合

```typescript
describe('GitHub Actions Integration', () => {
  test('should create annotations', async () => {
    process.env.GITHUB_ACTIONS = 'true';
    const result = await runCommand('monoguard validate --ci');

    expect(result.output).toContain('::error file=');
    expect(result.output).toContain('::warning file=');
  });

  test('should set exit code correctly', async () => {
    const result = await runCommand('monoguard validate --fail-on=error');
    expect(result.exitCode).toBe(1); // Has errors
  });
});
```

#### 2. 增量分析測試

```typescript
describe('Incremental Analysis', () => {
  test('should analyze only changed files', async () => {
    await git.checkout('feature-branch');
    const result = await runCommand('monoguard analyze --incremental');

    expect(result.analyzedFiles).toHaveLength(5); // Only changed
    expect(result.duration).toBeLessThan(30000); // < 30s
  });
});
```

## 完成標準 (Definition of Done)

- [ ] GitHub Actions 完整支援
- [ ] GitLab CI 完整支援
- [ ] Jenkins 基本支援
- [ ] 增量分析功能完整
- [ ] CI 文件完善
- [ ] 範例工作流程提供
- [ ] 效能優化完成
