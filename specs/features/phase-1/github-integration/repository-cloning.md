# GitHub Integration: Repository Cloning 功能規格

## 概述

實現 GitHub 倉庫的自動克隆、分析和結果展示功能，讓用戶可以直接輸入 GitHub URL 進行 monorepo 分析。

## 功能細節

### 1. 倉庫克隆流程

```typescript
interface GitHubAnalysisFlow {
  // 1. 驗證倉庫 URL
  validateRepoUrl(url: string): RepoInfo;

  // 2. 克隆倉庫
  cloneRepository(repoInfo: RepoInfo): CloneResult;

  // 3. 偵測 workspace
  detectWorkspace(repoPath: string): WorkspaceInfo;

  // 4. 執行分析
  analyzeRepository(repoPath: string): AnalysisResult;

  // 5. 清理臨時文件
  cleanup(repoPath: string): void;
}

interface RepoInfo {
  owner: string;
  repo: string;
  branch?: string;
  fullUrl: string;
  isPrivate: boolean;
}

interface CloneResult {
  success: boolean;
  localPath: string;
  commitHash: string;
  error?: string;
}
```

### 2. API 端點設計

#### POST /api/v1/analysis/github

```typescript
interface GitHubAnalysisRequest {
  repoUrl: string; // https://github.com/owner/repo
  branch?: string; // default: main
  accessToken?: string; // For private repos
  analyzeOptions?: {
    focus?: 'dependencies' | 'architecture' | 'all';
    includeDevDeps?: boolean;
  };
}

interface GitHubAnalysisResponse {
  jobId: string;
  status: 'queued' | 'cloning' | 'analyzing' | 'completed' | 'failed';
  repoInfo: RepoInfo;
  estimatedTime: number; // seconds
}
```

#### GET /api/v1/analysis/github/:jobId

```typescript
interface JobStatusResponse {
  jobId: string;
  status: 'queued' | 'cloning' | 'analyzing' | 'completed' | 'failed';
  progress: number; // 0-100
  currentStep: string;
  result?: AnalysisResult;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
}
```

### 3. 倉庫克隆實作 (Go)

```go
// services/github_service.go
type GitHubService struct {
    tmpDir string
    git    *GitClient
    auth   *AuthService
}

func (s *GitHubService) CloneRepository(req GitHubAnalysisRequest) (*CloneResult, error) {
    // 1. Parse GitHub URL
    repoInfo, err := parseGitHubURL(req.RepoUrl)
    if err != nil {
        return nil, fmt.Errorf("invalid GitHub URL: %w", err)
    }

    // 2. Create temporary directory
    tmpPath := filepath.Join(s.tmpDir, fmt.Sprintf("repo-%s", uuid.New().String()))
    if err := os.MkdirAll(tmpPath, 0755); err != nil {
        return nil, err
    }

    // 3. Clone repository
    cloneOptions := &git.CloneOptions{
        URL:      repoInfo.FullUrl,
        Progress: os.Stdout,
        Depth:    1, // Shallow clone for faster speed
    }

    // Add authentication if provided
    if req.AccessToken != "" {
        cloneOptions.Auth = &http.BasicAuth{
            Username: "oauth2",
            Password: req.AccessToken,
        }
    }

    // Clone with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    repo, err := git.PlainCloneContext(ctx, tmpPath, false, cloneOptions)
    if err != nil {
        os.RemoveAll(tmpPath)
        return nil, fmt.Errorf("failed to clone repository: %w", err)
    }

    // 4. Get commit hash
    ref, err := repo.Head()
    if err != nil {
        return nil, err
    }

    return &CloneResult{
        Success:    true,
        LocalPath:  tmpPath,
        CommitHash: ref.Hash().String(),
    }, nil
}

func parseGitHubURL(url string) (*RepoInfo, error) {
    // Support multiple formats:
    // - https://github.com/owner/repo
    // - https://github.com/owner/repo.git
    // - git@github.com:owner/repo.git
    // - github.com/owner/repo

    patterns := []string{
        `github\.com[:/]([^/]+)/([^/\.]+)`,
        `https://github\.com/([^/]+)/([^/\.]+)`,
    }

    for _, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindStringSubmatch(url)
        if len(matches) == 3 {
            return &RepoInfo{
                Owner:   matches[1],
                Repo:    matches[2],
                FullUrl: fmt.Sprintf("https://github.com/%s/%s", matches[1], matches[2]),
            }, nil
        }
    }

    return nil, errors.New("invalid GitHub URL format")
}
```

### 4. 異步處理機制

```go
// Background job processing
type AnalysisJob struct {
    ID          string
    RepoURL     string
    Status      string
    Progress    int
    CurrentStep string
    Result      *AnalysisResult
    Error       *JobError
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func (s *GitHubService) ProcessAnalysisJob(jobID string, req GitHubAnalysisRequest) {
    job := &AnalysisJob{
        ID:      jobID,
        RepoURL: req.RepoUrl,
        Status:  "queued",
    }

    // Update job status
    updateJob := func(status string, progress int, step string) {
        job.Status = status
        job.Progress = progress
        job.CurrentStep = step
        job.UpdatedAt = time.Now()
        s.saveJob(job)
    }

    // Step 1: Clone repository
    updateJob("cloning", 10, "Cloning repository...")
    cloneResult, err := s.CloneRepository(req)
    if err != nil {
        job.Error = &JobError{Code: "CLONE_FAILED", Message: err.Error()}
        updateJob("failed", 0, "")
        return
    }
    defer os.RemoveAll(cloneResult.LocalPath)

    // Step 2: Detect workspace
    updateJob("analyzing", 30, "Detecting workspace structure...")
    workspace, err := s.DetectWorkspace(cloneResult.LocalPath)
    if err != nil {
        job.Error = &JobError{Code: "WORKSPACE_DETECTION_FAILED", Message: err.Error()}
        updateJob("failed", 0, "")
        return
    }

    // Step 3: Run analysis
    updateJob("analyzing", 50, "Analyzing dependencies...")
    analysisResult, err := s.AnalyzeRepository(cloneResult.LocalPath, workspace)
    if err != nil {
        job.Error = &JobError{Code: "ANALYSIS_FAILED", Message: err.Error()}
        updateJob("failed", 0, "")
        return
    }

    // Complete
    job.Result = analysisResult
    updateJob("completed", 100, "Analysis completed")
}
```

### 5. 前端整合

```typescript
// components/GitHubAnalyzer.tsx
export const GitHubAnalyzer: React.FC = () => {
  const [repoUrl, setRepoUrl] = useState('');
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [jobId, setJobId] = useState<string | null>(null);

  const startAnalysis = async () => {
    setIsAnalyzing(true);
    try {
      const response = await api.post('/analysis/github', { repoUrl });
      setJobId(response.data.jobId);

      // Start polling for status
      pollJobStatus(response.data.jobId);
    } catch (error) {
      toast.error('Failed to start analysis');
      setIsAnalyzing(false);
    }
  };

  const pollJobStatus = async (jobId: string) => {
    const interval = setInterval(async () => {
      const status = await api.get(`/analysis/github/${jobId}`);

      if (status.data.status === 'completed') {
        clearInterval(interval);
        setIsAnalyzing(false);
        // Show results
        router.push(`/analysis/${status.data.result.id}`);
      } else if (status.data.status === 'failed') {
        clearInterval(interval);
        setIsAnalyzing(false);
        toast.error(status.data.error.message);
      }
    }, 2000); // Poll every 2 seconds
  };

  return (
    <div>
      <input
        type="text"
        placeholder="https://github.com/owner/repo"
        value={repoUrl}
        onChange={(e) => setRepoUrl(e.target.value)}
      />
      <button onClick={startAnalysis} disabled={isAnalyzing}>
        {isAnalyzing ? 'Analyzing...' : 'Analyze Repository'}
      </button>

      {isAnalyzing && <ProgressIndicator jobId={jobId} />}
    </div>
  );
};
```

## User Stories

### User Story 1: 分析公開倉庫

**As a** 開發者
**I want to** 輸入 GitHub URL 直接分析公開倉庫
**So that** 我不需要手動下載和上傳代碼

**Acceptance Criteria:**

- [ ] 支援標準 GitHub URL 格式
- [ ] 自動克隆倉庫
- [ ] 顯示分析進度
- [ ] 分析完成後顯示結果
- [ ] 克隆時間 < 2 分鐘

### User Story 2: 分析私有倉庫

**As a** 企業用戶
**I want to** 使用 Personal Access Token 分析私有倉庫
**So that** 我可以分析公司的內部專案

**Acceptance Criteria:**

- [ ] 支援 GitHub Token 認證
- [ ] Token 安全存儲
- [ ] 驗證 Token 權限
- [ ] 顯示清楚的權限錯誤
- [ ] Token 過期提醒

### User Story 3: 追蹤分析進度

**As a** 用戶
**I want to** 看到分析的詳細進度
**So that** 我知道當前正在做什麼

**Acceptance Criteria:**

- [ ] 顯示當前步驟 (克隆/偵測/分析)
- [ ] 顯示進度百分比
- [ ] 預估剩餘時間
- [ ] 支援取消分析
- [ ] 分析失敗時顯示詳細錯誤

## 測試項目

### 單元測試

```go
func TestParseGitHubURL(t *testing.T) {
    tests := []struct {
        input    string
        expected *RepoInfo
        hasError bool
    }{
        {
            input: "https://github.com/owner/repo",
            expected: &RepoInfo{
                Owner: "owner",
                Repo:  "repo",
            },
            hasError: false,
        },
        {
            input:    "invalid-url",
            expected: nil,
            hasError: true,
        },
    }

    for _, tt := range tests {
        result, err := parseGitHubURL(tt.input)
        if tt.hasError {
            assert.Error(t, err)
        } else {
            assert.NoError(t, err)
            assert.Equal(t, tt.expected.Owner, result.Owner)
            assert.Equal(t, tt.expected.Repo, result.Repo)
        }
    }
}
```

### 整合測試

```go
func TestCloneRepository(t *testing.T) {
    service := NewGitHubService()

    result, err := service.CloneRepository(GitHubAnalysisRequest{
        RepoUrl: "https://github.com/facebook/react",
    })

    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.NotEmpty(t, result.LocalPath)
    assert.NotEmpty(t, result.CommitHash)

    // Cleanup
    defer os.RemoveAll(result.LocalPath)

    // Verify files exist
    assert.DirExists(t, result.LocalPath)
    assert.FileExists(t, filepath.Join(result.LocalPath, "package.json"))
}
```

## 完成標準 (Definition of Done)

- [ ] 倉庫克隆功能完成
- [ ] 支援公開和私有倉庫
- [ ] 異步處理機制實作
- [ ] 進度追蹤功能
- [ ] 前端整合完成
- [ ] 錯誤處理完善
- [ ] 所有測試通過
- [ ] 臨時文件清理機制
- [ ] 安全性驗證通過
- [ ] 文檔更新
