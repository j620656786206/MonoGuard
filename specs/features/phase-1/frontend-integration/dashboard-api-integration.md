# Frontend: Dashboard API Integration 功能規格

## 概述

將儀表板從使用模擬數據切換到真實 API 數據，實現完整的前後端數據流，提供即時的專案健康狀況和分析結果。

## 功能細節

### 1. API 整合架構

```typescript
// API Client 結構
class MonoGuardAPI {
  // 專案管理
  async getProjects(): Promise<Project[]>;
  async getProject(id: string): Promise<Project>;
  async createProject(data: CreateProjectRequest): Promise<Project>;

  // 分析結果
  async getLatestAnalysis(projectId: string): Promise<AnalysisResult>;
  async getHealthScore(projectId: string): Promise<HealthScore>;
  async getAnalysisHistory(projectId: string): Promise<AnalysisRun[]>;

  // 技術債務
  async getIssues(projectId: string, filters?: IssueFilters): Promise<Issue[]>;
  async getViolations(projectId: string): Promise<Violation[]>;

  // 統計數據
  async getDashboardStats(projectId: string): Promise<DashboardStats>;
  async getTrends(projectId: string, period: string): Promise<TrendData>;
}
```

### 2. 儀表板資料結構

#### DashboardStats

```typescript
interface DashboardStats {
  overview: {
    totalPackages: number;
    totalDependencies: number;
    healthScore: number;
    lastAnalyzedAt: string;
  };

  issues: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };

  trends: {
    healthScoreChange: number; // percentage
    issuesChange: number;
    period: '7d' | '30d' | '90d';
  };

  topIssues: Issue[];
  recentAnalyses: AnalysisRun[];
}
```

### 3. 資料載入策略

#### React Query 整合

```typescript
// hooks/useDashboard.ts
export const useDashboardStats = (projectId: string) => {
  return useQuery({
    queryKey: ['dashboard', projectId],
    queryFn: () => api.getDashboardStats(projectId),
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 30 * 1000, // Auto-refresh every 30s
    enabled: !!projectId,
  });
};

export const useHealthScore = (projectId: string) => {
  return useQuery({
    queryKey: ['healthScore', projectId],
    queryFn: () => api.getHealthScore(projectId),
    staleTime: 2 * 60 * 1000,
  });
};

export const useIssues = (projectId: string, filters: IssueFilters) => {
  return useQuery({
    queryKey: ['issues', projectId, filters],
    queryFn: () => api.getIssues(projectId, filters),
    keepPreviousData: true,
  });
};
```

#### Loading States

```typescript
interface LoadingState {
  isLoading: boolean;
  isError: boolean;
  error?: Error;
  isRefetching: boolean;
}

// 使用 Skeleton 組件
if (isLoading) {
  return <DashboardSkeleton />;
}

if (isError) {
  return <ErrorDisplay error={error} onRetry={refetch} />;
}
```

### 4. 即時更新機制

#### WebSocket 整合

```typescript
// hooks/useRealtimeUpdates.ts
export const useRealtimeUpdates = (projectId: string) => {
  const queryClient = useQueryClient();

  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}/projects/${projectId}/updates`);

    ws.onmessage = (event) => {
      const update: UpdateEvent = JSON.parse(event.data);

      if (update.type === 'analysis_complete') {
        // Invalidate and refetch dashboard data
        queryClient.invalidateQueries(['dashboard', projectId]);
        queryClient.invalidateQueries(['healthScore', projectId]);

        // Show notification
        toast.success('Analysis completed!');
      }
    };

    return () => ws.close();
  }, [projectId]);
};
```

### 5. 錯誤處理與重試

```typescript
// API Client with retry logic
axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Retry on network error
    if (!originalRequest._retry && error.message === 'Network Error') {
      originalRequest._retry = true;
      await new Promise((resolve) => setTimeout(resolve, 1000));
      return axios(originalRequest);
    }

    // Handle 401 - refresh token
    if (error.response?.status === 401) {
      try {
        await refreshAccessToken();
        return axios(originalRequest);
      } catch {
        // Redirect to login
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);
```

## User Stories

### User Story 1: 即時儀表板數據

**As a** 專案管理者
**I want to** 查看即時的專案健康狀況
**So that** 我可以快速了解當前的技術債務狀況

**Acceptance Criteria:**

- [ ] 儀表板載入時間 < 3 秒
- [ ] 顯示最新的分析結果
- [ ] 健康度評分即時更新
- [ ] 問題列表自動刷新
- [ ] 顯示上次更新時間

### User Story 2: 分析進度追蹤

**As a** 開發者
**I want to** 看到分析正在進行中
**So that** 我知道何時可以查看結果

**Acceptance Criteria:**

- [ ] 顯示分析進度條
- [ ] 即時更新分析狀態
- [ ] 分析完成後自動刷新數據
- [ ] 提供取消分析功能
- [ ] 顯示預估完成時間

### User Story 3: 錯誤處理與恢復

**As a** 用戶
**I want to** 在網路錯誤時看到友善的提示
**So that** 我知道發生了什麼並可以重試

**Acceptance Criteria:**

- [ ] 顯示清楚的錯誤訊息
- [ ] 提供重試按鈕
- [ ] 自動重試網路錯誤
- [ ] 離線時顯示快取數據
- [ ] 網路恢復後自動重新載入

## 測試項目

### 整合測試

```typescript
describe('Dashboard API Integration', () => {
  test('should load dashboard data on mount', async () => {
    render(<Dashboard projectId="test-project" />);

    // Should show loading state
    expect(screen.getByText(/loading/i)).toBeInTheDocument();

    // Should fetch and display data
    await waitFor(() => {
      expect(screen.getByText(/health score/i)).toBeInTheDocument();
    });

    // Should show actual data
    expect(screen.getByText('85')).toBeInTheDocument(); // Health score
  });

  test('should handle API errors gracefully', async () => {
    server.use(
      rest.get('/api/v1/dashboard/:id', (req, res, ctx) => {
        return res(ctx.status(500), ctx.json({ error: 'Server error' }));
      })
    );

    render(<Dashboard projectId="test-project" />);

    await waitFor(() => {
      expect(screen.getByText(/error/i)).toBeInTheDocument();
    });

    // Should show retry button
    expect(screen.getByText(/retry/i)).toBeInTheDocument();
  });

  test('should auto-refresh data', async () => {
    jest.useFakeTimers();
    render(<Dashboard projectId="test-project" />);

    await waitFor(() => {
      expect(screen.getByText(/health score/i)).toBeInTheDocument();
    });

    const initialScore = screen.getByTestId('health-score').textContent;

    // Advance time by 30 seconds (auto-refresh interval)
    act(() => {
      jest.advanceTimersByTime(30000);
    });

    // Should refetch and potentially show new data
    await waitFor(() => {
      const newScore = screen.getByTestId('health-score').textContent;
      // Data might have changed
    });

    jest.useRealTimers();
  });
});
```

### WebSocket 測試

```typescript
describe('Realtime Updates', () => {
  test('should update dashboard on analysis complete', async () => {
    const mockWs = new MockWebSocket();
    render(<Dashboard projectId="test-project" />);

    // Trigger analysis complete event
    act(() => {
      mockWs.send({
        type: 'analysis_complete',
        projectId: 'test-project',
        analysisId: 'new-analysis'
      });
    });

    // Should refetch dashboard data
    await waitFor(() => {
      expect(screen.getByText('Analysis completed!')).toBeInTheDocument();
    });
  });
});
```

## 完成標準 (Definition of Done)

- [ ] 所有 API 整合完成
- [ ] React Query 設定完善
- [ ] WebSocket 即時更新實作
- [ ] Loading 和 Error 狀態處理
- [ ] 自動重試機制
- [ ] Token 刷新邏輯
- [ ] 所有測試通過 (覆蓋率 ≥ 80%)
- [ ] 效能優化完成 (載入 < 3 秒)
- [ ] 錯誤處理完善
- [ ] 文檔更新
