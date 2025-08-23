# MonoGuard Performance Guide

## Performance Requirements

### Target Performance Metrics
- **Analysis Completion**: < 5 minutes for 100+ packages
- **API Response Time**: < 200ms for 95th percentile
- **Database Query Performance**: < 100ms for complex queries
- **Memory Usage**: < 512MB per analysis worker
- **Concurrent Users**: Support 1000+ concurrent users
- **Uptime**: 99.9% availability (8.77 hours downtime/year)

### Performance Testing Benchmarks
```bash
# Analysis engine benchmarks
Package Count | Target Time | Memory Usage
-------------|-------------|-------------
10 packages  | < 10 seconds| < 64MB
50 packages  | < 45 seconds| < 128MB
100 packages | < 2 minutes | < 256MB
500 packages | < 5 minutes | < 512MB
1000+ packages| < 10 minutes| < 1GB
```

## Analysis Engine Performance

### Concurrent Processing Strategy
```go
// Optimized analysis engine with worker pools
type AnalysisEngine struct {
    workerPool   *WorkerPool
    resultCache  *Cache
    maxWorkers   int
    timeout      time.Duration
}

type WorkerPool struct {
    jobs     chan AnalysisJob
    results  chan AnalysisResult
    workers  []*Worker
    ctx      context.Context
    cancel   context.CancelFunc
}

func NewAnalysisEngine(config EngineConfig) *AnalysisEngine {
    return &AnalysisEngine{
        workerPool:  NewWorkerPool(config.MaxWorkers),
        resultCache: NewLRUCache(config.CacheSize),
        maxWorkers:  config.MaxWorkers,
        timeout:     config.Timeout,
    }
}

func (e *AnalysisEngine) AnalyzeProject(ctx context.Context, project *Project) (*Analysis, error) {
    // Check cache first
    if cached := e.resultCache.Get(project.CacheKey()); cached != nil {
        return cached.(*Analysis), nil
    }
    
    // Create analysis jobs for concurrent processing
    jobs := e.createAnalysisJobs(project)
    
    // Process jobs concurrently
    results := make(chan AnalysisResult, len(jobs))
    for _, job := range jobs {
        select {
        case e.workerPool.jobs <- job:
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
    
    // Collect results with timeout
    analysis := &Analysis{ProjectID: project.ID}
    for i := 0; i < len(jobs); i++ {
        select {
        case result := <-results:
            if result.Error != nil {
                return nil, result.Error
            }
            analysis.MergeResult(result)
        case <-time.After(e.timeout):
            return nil, errors.New("analysis timeout")
        case <-ctx.Done():
            return nil, ctx.Err()
        }
    }
    
    // Cache successful analysis
    e.resultCache.Set(project.CacheKey(), analysis, time.Hour)
    return analysis, nil
}
```

### Incremental Analysis Implementation
```go
// Incremental analysis for large repositories
type IncrementalAnalyzer struct {
    baseline     *Analysis
    changeTracker *ChangeTracker
    analyzer     *AnalysisEngine
}

type ChangeTracker struct {
    lastAnalysis time.Time
    changedFiles map[string]FileChecksum
    addedFiles   []string
    deletedFiles []string
}

func (ia *IncrementalAnalyzer) AnalyzeChanges(ctx context.Context, project *Project) (*Analysis, error) {
    changes, err := ia.changeTracker.DetectChanges(project.Path)
    if err != nil {
        return nil, err
    }
    
    // If no changes, return cached baseline
    if len(changes.ChangedFiles) == 0 && len(changes.AddedFiles) == 0 && len(changes.DeletedFiles) == 0 {
        return ia.baseline, nil
    }
    
    // Only analyze changed/new files
    partialProject := &Project{
        ID:   project.ID,
        Path: project.Path,
        Files: append(changes.ChangedFiles, changes.AddedFiles...),
    }
    
    partialAnalysis, err := ia.analyzer.AnalyzeProject(ctx, partialProject)
    if err != nil {
        return nil, err
    }
    
    // Merge with baseline
    result := ia.baseline.Clone()
    result.Merge(partialAnalysis)
    result.RemoveDeletedFiles(changes.DeletedFiles)
    
    ia.baseline = result
    ia.changeTracker.UpdateLastAnalysis()
    
    return result, nil
}
```

### Memory Management and Pooling
```go
// Object pooling for frequently allocated objects
var (
    dependencyPool = sync.Pool{
        New: func() interface{} {
            return &Dependency{
                Children: make([]*Dependency, 0, 10),
                Metadata: make(map[string]interface{}),
            }
        },
    }
    
    astNodePool = sync.Pool{
        New: func() interface{} {
            return &ASTNode{
                Children: make([]*ASTNode, 0, 5),
                Attributes: make(map[string]string),
            }
        },
    }
)

func AcquireDependency() *Dependency {
    return dependencyPool.Get().(*Dependency)
}

func ReleaseDependency(d *Dependency) {
    // Reset object state
    d.Reset()
    dependencyPool.Put(d)
}

// Memory-conscious dependency tree builder
type DependencyTreeBuilder struct {
    nodePool    *sync.Pool
    maxDepth    int
    maxNodes    int
    currentNodes int
}

func (b *DependencyTreeBuilder) BuildTree(ctx context.Context, rootPath string) (*DependencyTree, error) {
    if b.currentNodes >= b.maxNodes {
        return nil, errors.New("maximum node limit reached")
    }
    
    root := b.nodePool.Get().(*DependencyNode)
    defer b.nodePool.Put(root)
    
    return b.buildTreeRecursive(ctx, rootPath, root, 0)
}

func (b *DependencyTreeBuilder) buildTreeRecursive(ctx context.Context, path string, node *DependencyNode, depth int) (*DependencyTree, error) {
    if depth > b.maxDepth {
        return nil, errors.New("maximum depth exceeded")
    }
    
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // Process node...
    b.currentNodes++
    
    // Recursive processing with memory bounds
    for _, child := range getChildren(path) {
        if b.currentNodes >= b.maxNodes {
            break
        }
        
        childNode := b.nodePool.Get().(*DependencyNode)
        _, err := b.buildTreeRecursive(ctx, child, childNode, depth+1)
        b.nodePool.Put(childNode)
        
        if err != nil {
            return nil, err
        }
    }
    
    return &DependencyTree{Root: node}, nil
}
```

## Database Performance Optimization

### Indexing Strategy
```sql
-- Performance-optimized database schema with indexes

-- Projects table with optimized indexes
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    repository_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    last_analyzed_at TIMESTAMPTZ,
    health_score DECIMAL(4,2) DEFAULT 0.00
);

-- Composite indexes for common query patterns
CREATE INDEX idx_projects_user_created ON projects(user_id, created_at DESC);
CREATE INDEX idx_projects_health_score ON projects(health_score DESC) WHERE health_score IS NOT NULL;
CREATE INDEX idx_projects_last_analyzed ON projects(last_analyzed_at DESC) WHERE last_analyzed_at IS NOT NULL;

-- Analyses table with partitioning for time-series data
CREATE TABLE analyses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    analysis_data JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Monthly partitions for analyses
CREATE TABLE analyses_2024_01 PARTITION OF analyses
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
    
-- Indexes on partitioned table
CREATE INDEX idx_analyses_project_status ON analyses(project_id, status);
CREATE INDEX idx_analyses_created_at ON analyses(created_at DESC);
CREATE INDEX idx_analyses_data_gin ON analyses USING GIN (analysis_data);

-- Dependencies table with efficient storage
CREATE TABLE dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    package_name VARCHAR(255) NOT NULL,
    version VARCHAR(100),
    dependency_type VARCHAR(50), -- 'direct', 'transitive', 'dev'
    size_bytes BIGINT,
    vulnerability_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_dependencies_analysis ON dependencies(analysis_id);
CREATE INDEX idx_dependencies_package ON dependencies(package_name, version);
CREATE INDEX idx_dependencies_vulns ON dependencies(vulnerability_count) WHERE vulnerability_count > 0;
```

### Query Optimization Patterns
```go
// Optimized repository layer with efficient queries
type OptimizedProjectRepository struct {
    db    *gorm.DB
    cache *Cache
}

// Use database connection pooling
func NewOptimizedProjectRepository(db *gorm.DB) *OptimizedProjectRepository {
    // Configure connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    
    return &OptimizedProjectRepository{
        db:    db,
        cache: NewRedisCache(),
    }
}

// Efficient project listing with pagination
func (r *OptimizedProjectRepository) GetProjectsPaginated(ctx context.Context, userID string, limit, offset int) (*ProjectPage, error) {
    cacheKey := fmt.Sprintf("projects:user:%s:page:%d:%d", userID, limit, offset)
    
    // Try cache first
    if cached := r.cache.Get(cacheKey); cached != nil {
        return cached.(*ProjectPage), nil
    }
    
    var projects []Project
    var total int64
    
    // Use efficient counting
    countQuery := r.db.Model(&Project{}).Where("user_id = ?", userID)
    if err := countQuery.Count(&total).Error; err != nil {
        return nil, err
    }
    
    // Optimized main query with selected fields
    query := r.db.Select("id, name, health_score, last_analyzed_at, created_at").
        Where("user_id = ?", userID).
        Order("last_analyzed_at DESC NULLS LAST").
        Limit(limit).
        Offset(offset)
    
    if err := query.Find(&projects).Error; err != nil {
        return nil, err
    }
    
    result := &ProjectPage{
        Projects:    projects,
        Total:       total,
        Limit:       limit,
        Offset:      offset,
        HasMore:     offset+limit < int(total),
    }
    
    // Cache for 5 minutes
    r.cache.Set(cacheKey, result, 5*time.Minute)
    return result, nil
}

// Batch loading to prevent N+1 queries
func (r *OptimizedProjectRepository) GetProjectsWithAnalyses(ctx context.Context, projectIDs []string) (map[string]*ProjectWithAnalysis, error) {
    if len(projectIDs) == 0 {
        return make(map[string]*ProjectWithAnalysis), nil
    }
    
    // Load projects in batch
    var projects []Project
    if err := r.db.Where("id IN ?", projectIDs).Find(&projects).Error; err != nil {
        return nil, err
    }
    
    // Load latest analyses in batch
    var analyses []Analysis
    subquery := r.db.Select("project_id, MAX(completed_at) as max_completed").
        Where("project_id IN ? AND status = ?", projectIDs, "completed").
        Group("project_id")
    
    if err := r.db.Joins("JOIN (?) as latest ON analyses.project_id = latest.project_id AND analyses.completed_at = latest.max_completed", subquery).
        Find(&analyses).Error; err != nil {
        return nil, err
    }
    
    // Combine results
    result := make(map[string]*ProjectWithAnalysis)
    analysisMap := make(map[string]*Analysis)
    
    for i := range analyses {
        analysisMap[analyses[i].ProjectID] = &analyses[i]
    }
    
    for i := range projects {
        result[projects[i].ID] = &ProjectWithAnalysis{
            Project:        &projects[i],
            LatestAnalysis: analysisMap[projects[i].ID],
        }
    }
    
    return result, nil
}
```

### Database Connection Optimization
```go
// Connection pooling and optimization
func ConfigureDatabase(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.New(
            log.New(os.Stdout, "\r\n", log.LstdFlags),
            logger.Config{
                SlowThreshold:             100 * time.Millisecond,
                LogLevel:                  logger.Warn,
                IgnoreRecordNotFoundError: true,
                Colorful:                  false,
            },
        ),
        PrepareStmt: true, // Cache prepared statements
    })
    
    if err != nil {
        return nil, err
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // Optimize connection pool
    sqlDB.SetMaxOpenConns(25)                // Maximum open connections
    sqlDB.SetMaxIdleConns(10)                // Maximum idle connections
    sqlDB.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
    sqlDB.SetConnMaxIdleTime(5 * time.Minute) // Idle connection timeout
    
    return db, nil
}
```

## Frontend Performance Optimization

### React Performance Patterns
```typescript
// Optimized React components with performance best practices

// Memoized component for expensive renders
const DependencyGraph = React.memo<DependencyGraphProps>(({ 
  dependencies, 
  onNodeClick,
  selectedNodeId 
}) => {
  // Memoize expensive calculations
  const processedData = useMemo(() => {
    return processDependencies(dependencies);
  }, [dependencies]);
  
  // Debounce search to prevent excessive renders
  const [searchTerm, setSearchTerm] = useState('');
  const debouncedSearchTerm = useDebounce(searchTerm, 300);
  
  // Filter nodes based on search
  const filteredNodes = useMemo(() => {
    if (!debouncedSearchTerm) return processedData.nodes;
    return processedData.nodes.filter(node => 
      node.name.toLowerCase().includes(debouncedSearchTerm.toLowerCase())
    );
  }, [processedData.nodes, debouncedSearchTerm]);
  
  // Use callback to prevent child re-renders
  const handleNodeClick = useCallback((nodeId: string) => {
    onNodeClick?.(nodeId);
  }, [onNodeClick]);
  
  return (
    <div className="dependency-graph">
      <SearchInput 
        value={searchTerm}
        onChange={setSearchTerm}
        placeholder="Search dependencies..."
      />
      <VirtualizedGraph
        nodes={filteredNodes}
        edges={processedData.edges}
        onNodeClick={handleNodeClick}
        selectedNodeId={selectedNodeId}
      />
    </div>
  );
});

// Custom hook for debouncing
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}
```

### Virtualization for Large Lists
```typescript
// Virtualized list for handling large dependency lists
import { FixedSizeList as List } from 'react-window';

interface VirtualizedDependencyListProps {
  dependencies: Dependency[];
  onDependencyClick: (dependency: Dependency) => void;
}

const VirtualizedDependencyList: React.FC<VirtualizedDependencyListProps> = ({
  dependencies,
  onDependencyClick
}) => {
  const Row = useCallback(({ index, style }: { index: number; style: CSSProperties }) => (
    <div style={style}>
      <DependencyItem
        dependency={dependencies[index]}
        onClick={onDependencyClick}
      />
    </div>
  ), [dependencies, onDependencyClick]);

  return (
    <List
      height={400}
      itemCount={dependencies.length}
      itemSize={50}
      width="100%"
    >
      {Row}
    </List>
  );
};

// Optimized dependency item component
const DependencyItem = React.memo<DependencyItemProps>(({ dependency, onClick }) => {
  const handleClick = useCallback(() => {
    onClick(dependency);
  }, [dependency, onClick]);

  return (
    <div 
      className="dependency-item"
      onClick={handleClick}
    >
      <span className="name">{dependency.name}</span>
      <span className="version">{dependency.version}</span>
      <Badge variant={dependency.riskLevel}>
        {dependency.riskLevel}
      </Badge>
    </div>
  );
});
```

### State Management Optimization
```typescript
// Optimized Zustand store with selectors
interface AppState {
  projects: Project[];
  currentProject: Project | null;
  analyses: Record<string, Analysis>;
  loading: Record<string, boolean>;
  errors: Record<string, string>;
}

interface AppActions {
  setProjects: (projects: Project[]) => void;
  setCurrentProject: (project: Project) => void;
  addAnalysis: (projectId: string, analysis: Analysis) => void;
  setLoading: (key: string, loading: boolean) => void;
  setError: (key: string, error: string | null) => void;
}

// Create store with optimized selectors
export const useAppStore = create<AppState & AppActions>()(
  subscribeWithSelector((set, get) => ({
    projects: [],
    currentProject: null,
    analyses: {},
    loading: {},
    errors: {},

    setProjects: (projects) => set({ projects }),
    setCurrentProject: (project) => set({ currentProject: project }),
    
    addAnalysis: (projectId, analysis) => 
      set((state) => ({
        analyses: {
          ...state.analyses,
          [projectId]: analysis
        }
      })),
    
    setLoading: (key, loading) =>
      set((state) => ({
        loading: {
          ...state.loading,
          [key]: loading
        }
      })),
    
    setError: (key, error) =>
      set((state) => ({
        errors: {
          ...state.errors,
          [key]: error || ''
        }
      }))
  }))
);

// Selective subscriptions to prevent unnecessary re-renders
export const useProjects = () => useAppStore(state => state.projects);
export const useCurrentProject = () => useAppStore(state => state.currentProject);
export const useAnalysis = (projectId: string) => useAppStore(
  state => state.analyses[projectId]
);
export const useLoading = (key: string) => useAppStore(
  state => state.loading[key] || false
);
```

## Caching Strategy

### Multi-Level Caching Architecture
```go
// Hierarchical caching system
type CacheManager struct {
    l1Cache *LRUCache    // In-memory cache (fastest)
    l2Cache *RedisCache  // Distributed cache (fast)
    l3Cache *DatabaseCache // Persistent cache (slower)
}

func (cm *CacheManager) Get(key string) (interface{}, error) {
    // Try L1 cache first
    if value := cm.l1Cache.Get(key); value != nil {
        return value, nil
    }
    
    // Try L2 cache
    if value, err := cm.l2Cache.Get(key); err == nil && value != nil {
        // Populate L1 cache
        cm.l1Cache.Set(key, value, time.Hour)
        return value, nil
    }
    
    // Try L3 cache (database)
    if value, err := cm.l3Cache.Get(key); err == nil && value != nil {
        // Populate L2 and L1 caches
        cm.l2Cache.Set(key, value, 24*time.Hour)
        cm.l1Cache.Set(key, value, time.Hour)
        return value, nil
    }
    
    return nil, errors.New("cache miss")
}

func (cm *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
    // Set in all cache levels
    cm.l1Cache.Set(key, value, ttl)
    cm.l2Cache.Set(key, value, ttl)
    cm.l3Cache.Set(key, value, ttl)
    
    return nil
}

// Cache invalidation strategy
func (cm *CacheManager) InvalidatePattern(pattern string) error {
    keys, err := cm.l2Cache.GetKeysByPattern(pattern)
    if err != nil {
        return err
    }
    
    for _, key := range keys {
        cm.l1Cache.Delete(key)
        cm.l2Cache.Delete(key)
        cm.l3Cache.Delete(key)
    }
    
    return nil
}
```

### Redis Optimization
```go
// Optimized Redis configuration
func NewRedisClient(addr, password string) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:            addr,
        Password:        password,
        DB:              0,
        MaxRetries:      3,
        PoolSize:        10,
        MinIdleConns:    5,
        ConnMaxIdleTime: 5 * time.Minute,
        ReadTimeout:     100 * time.Millisecond,
        WriteTimeout:    100 * time.Millisecond,
    })
}

// Batch operations for efficiency
func (c *RedisCache) SetMultiple(pairs map[string]interface{}, ttl time.Duration) error {
    pipe := c.client.Pipeline()
    
    for key, value := range pairs {
        serialized, err := json.Marshal(value)
        if err != nil {
            return err
        }
        pipe.Set(context.Background(), key, serialized, ttl)
    }
    
    _, err := pipe.Exec(context.Background())
    return err
}

func (c *RedisCache) GetMultiple(keys []string) (map[string]interface{}, error) {
    pipe := c.client.Pipeline()
    
    cmds := make([]*redis.StringCmd, len(keys))
    for i, key := range keys {
        cmds[i] = pipe.Get(context.Background(), key)
    }
    
    _, err := pipe.Exec(context.Background())
    if err != nil {
        return nil, err
    }
    
    results := make(map[string]interface{})
    for i, cmd := range cmds {
        if cmd.Err() == nil {
            var value interface{}
            if err := json.Unmarshal([]byte(cmd.Val()), &value); err == nil {
                results[keys[i]] = value
            }
        }
    }
    
    return results, nil
}
```

## Performance Monitoring

### Application Metrics
```go
// Prometheus metrics for performance monitoring
var (
    analysisStarted = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "monoguard_analysis_started_total",
            Help: "Total number of analyses started",
        },
        []string{"project_type"},
    )
    
    analysisDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "monoguard_analysis_duration_seconds",
            Help:    "Time spent on analysis",
            Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1s to 512s
        },
        []string{"project_type", "status"},
    )
    
    dbQueryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "monoguard_db_query_duration_seconds",
            Help:    "Database query duration",
            Buckets: prometheus.DefBuckets,
        },
        []string{"query_type"},
    )
    
    cacheHitRate = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "monoguard_cache_hit_rate",
            Help: "Cache hit rate percentage",
        },
        []string{"cache_level"},
    )
)

// Middleware for automatic metrics collection
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        status := fmt.Sprintf("%d", c.Writer.Status())
        
        httpDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
        ).Observe(duration.Seconds())
        
        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
        ).Inc()
    }
}
```

### Performance Testing Framework
```go
// Load testing utilities
func BenchmarkAnalysisEngine(b *testing.B) {
    engine := NewAnalysisEngine(DefaultConfig())
    projects := generateTestProjects(100) // Generate 100 test projects
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            project := projects[rand.Intn(len(projects))]
            _, err := engine.AnalyzeProject(context.Background(), project)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}

func BenchmarkDatabaseQueries(b *testing.B) {
    db := setupTestDatabase()
    repo := NewProjectRepository(db)
    
    // Seed database with test data
    seedTestData(db, 1000)
    
    b.Run("GetProjects", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            userID := generateRandomUserID()
            _, err := repo.GetProjectsPaginated(context.Background(), userID, 20, 0)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
    
    b.Run("GetAnalyses", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            projectID := generateRandomProjectID()
            _, err := repo.GetAnalyses(context.Background(), projectID, 10)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}

// Memory benchmarking
func BenchmarkMemoryUsage(b *testing.B) {
    var m1, m2 runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m1)
    
    engine := NewAnalysisEngine(DefaultConfig())
    
    for i := 0; i < b.N; i++ {
        project := generateLargeProject() // Project with 500+ packages
        _, err := engine.AnalyzeProject(context.Background(), project)
        if err != nil {
            b.Fatal(err)
        }
    }
    
    runtime.GC()
    runtime.ReadMemStats(&m2)
    
    b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "bytes/op")
    b.ReportMetric(float64(m2.Mallocs-m1.Mallocs)/float64(b.N), "allocs/op")
}
```

## Performance Troubleshooting

### Common Performance Issues

#### Issue: Slow Analysis Performance
```go
// Diagnostic tools for analysis performance
func (e *AnalysisEngine) DiagnosePerformance(ctx context.Context, project *Project) (*PerformanceDiagnostic, error) {
    diagnostic := &PerformanceDiagnostic{
        ProjectID:     project.ID,
        StartTime:     time.Now(),
    }
    
    // Measure different phases
    phases := []struct{
        name string
        fn   func() error
    }{
        {"parsing", func() error { return e.parseProject(ctx, project) }},
        {"dependency_resolution", func() error { return e.resolveDependencies(ctx, project) }},
        {"analysis", func() error { return e.analyzeIssues(ctx, project) }},
        {"report_generation", func() error { return e.generateReport(ctx, project) }},
    }
    
    for _, phase := range phases {
        start := time.Now()
        err := phase.fn()
        duration := time.Since(start)
        
        diagnostic.Phases = append(diagnostic.Phases, PhaseMetric{
            Name:     phase.name,
            Duration: duration,
            Error:    err,
        })
        
        if err != nil {
            diagnostic.ErrorPhase = phase.name
            break
        }
    }
    
    diagnostic.TotalDuration = time.Since(diagnostic.StartTime)
    return diagnostic, nil
}
```

#### Issue: Database Query Performance
```sql
-- Performance analysis queries
-- Find slow queries
SELECT query, calls, total_time, mean_time, rows
FROM pg_stat_statements
WHERE mean_time > 100  -- queries taking more than 100ms on average
ORDER BY mean_time DESC
LIMIT 10;

-- Find missing indexes
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE schemaname = 'public'
  AND n_distinct > 100
  AND correlation < 0.1;

-- Analyze table statistics
ANALYZE VERBOSE projects;
ANALYZE VERBOSE analyses;
ANALYZE VERBOSE dependencies;
```

#### Issue: Memory Leaks
```go
// Memory leak detection
func DetectMemoryLeaks() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    var lastAlloc uint64
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        currentAlloc := m.TotalAlloc
        if lastAlloc > 0 {
            growth := currentAlloc - lastAlloc
            growthMB := float64(growth) / 1024 / 1024
            
            if growthMB > 10 { // More than 10MB growth
                log.Printf("Memory growth detected: %.2f MB", growthMB)
                
                // Trigger garbage collection
                runtime.GC()
                
                // Log heap dump for analysis
                writeHeapProfile()
            }
        }
        
        lastAlloc = currentAlloc
    }
}

func writeHeapProfile() {
    f, err := os.Create("heap.prof")
    if err != nil {
        log.Printf("Could not create heap profile: %v", err)
        return
    }
    defer f.Close()
    
    if err := pprof.WriteHeapProfile(f); err != nil {
        log.Printf("Could not write heap profile: %v", err)
    }
}
```

### Performance Optimization Checklist

#### Backend Optimization
- [ ] **Profiling**: Use `go tool pprof` for CPU and memory profiling
- [ ] **Database**: Optimize queries, add indexes, use connection pooling
- [ ] **Caching**: Implement multi-level caching with appropriate TTLs
- [ ] **Concurrency**: Use goroutines and channels for parallel processing
- [ ] **Memory**: Implement object pooling for frequently allocated objects
- [ ] **Timeouts**: Set appropriate timeouts for all operations
- [ ] **Resource Limits**: Set memory and CPU limits for containers

#### Frontend Optimization
- [ ] **Code Splitting**: Lazy load components and routes
- [ ] **Virtualization**: Use virtual scrolling for large lists
- [ ] **Memoization**: Cache expensive calculations with useMemo
- [ ] **Bundle Size**: Analyze and optimize JavaScript bundle size
- [ ] **Image Optimization**: Use WebP format and lazy loading
- [ ] **CDN**: Serve static assets from CDN
- [ ] **Service Worker**: Implement caching for offline functionality

#### Database Optimization
- [ ] **Query Analysis**: Regular query performance analysis
- [ ] **Index Optimization**: Create indexes for common query patterns
- [ ] **Partitioning**: Partition large tables by date or other criteria
- [ ] **Connection Pooling**: Optimize database connection pool settings
- [ ] **Read Replicas**: Use read replicas for reporting queries
- [ ] **Query Caching**: Cache expensive query results
- [ ] **Statistics**: Keep table statistics up to date