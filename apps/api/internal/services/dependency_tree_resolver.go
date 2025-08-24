package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/monoguard/api/internal/models"
	"github.com/sirupsen/logrus"
)

// DependencyTreeResolver provides enhanced dependency tree resolution with performance optimizations
type DependencyTreeResolver struct {
	logger    *logrus.Logger
	cache     *ResolverCache
	registry  *PackageRegistry
	metrics   *ResolverMetrics
}

// DependencyTree represents a resolved dependency tree
type DependencyTree struct {
	RootPackages []*PackageNode     `json:"rootPackages"`
	AllPackages  map[string]*PackageNode `json:"allPackages"`
	Conflicts    []*EnhancedConflict     `json:"conflicts"`
	Metadata     *TreeMetadata           `json:"metadata"`
}

// PackageNode represents a node in the dependency tree
type PackageNode struct {
	Name              string                 `json:"name"`
	Version           string                 `json:"version"`
	RequestedRange    string                 `json:"requestedRange"`
	ResolvedVersion   *SemanticVersion       `json:"resolvedVersion"`
	Dependencies      map[string]*PackageNode `json:"dependencies"`
	DevDependencies   map[string]*PackageNode `json:"devDependencies,omitempty"`
	PeerDependencies  map[string]*PackageNode `json:"peerDependencies,omitempty"`
	Path              string                 `json:"path"`
	IsWorkspace       bool                   `json:"isWorkspace"`
	Depth             int                    `json:"depth"`
	ConflictSeverity  models.Severity        `json:"conflictSeverity"`
	ResolutionSource  string                 `json:"resolutionSource"` // "registry", "workspace", "cache"
	Metadata          map[string]interface{} `json:"metadata"`
}

// EnhancedConflict represents an enhanced version conflict with resolution details
type EnhancedConflict struct {
	models.VersionConflict
	PackageNodes      []*PackageNode     `json:"packageNodes"`
	ResolutionOptions []*ResolutionOption `json:"resolutionOptions"`
	RecommendedFix    *ResolutionOption   `json:"recommendedFix"`
	AutoResolvable    bool               `json:"autoResolvable"`
}

// ResolutionOption represents a potential resolution for conflicts
type ResolutionOption struct {
	Strategy        ResolutionStrategy `json:"strategy"`
	TargetVersion   *SemanticVersion   `json:"targetVersion"`
	AffectedPackages []string          `json:"affectedPackages"`
	RiskAssessment  *ConflictRiskAssessment `json:"riskAssessment"`
	MigrationSteps  []*MigrationStep   `json:"migrationSteps"`
}

// TreeMetadata contains metadata about the dependency tree
type TreeMetadata struct {
	BuildTime         time.Duration `json:"buildTime"`
	TotalNodes        int           `json:"totalNodes"`
	ExternalPackages  int           `json:"externalPackages"`
	WorkspacePackages int           `json:"workspacePackages"`
	MaxDepth          int           `json:"maxDepth"`
	CacheHitRate      float64       `json:"cacheHitRate"`
	RegistryHitRate   float64       `json:"registryHitRate"`
	ConflictCount     int           `json:"conflictCount"`
	ResolvedConflicts int           `json:"resolvedConflicts"`
	Options           BuildOptions  `json:"options"`
}

// BuildOptions contains options for building the dependency tree
type BuildOptions struct {
	MaxDepth             int               `json:"maxDepth"`
	IncludeDevDeps       bool              `json:"includeDevDeps"`
	IncludePeerDeps      bool              `json:"includePeerDeps"`
	IncludeOptional      bool              `json:"includeOptional"`
	Strategy             ResolutionStrategy `json:"strategy"`
	PreferWorkspace      bool              `json:"preferWorkspace"`
	AllowPreRelease      bool              `json:"allowPreRelease"`
	EnableCaching        bool              `json:"enableCaching"`
	ConcurrencyLevel     int               `json:"concurrencyLevel"`
	TimeoutPerPackage    time.Duration     `json:"timeoutPerPackage"`
	UseNpmRegistry       bool              `json:"useNpmRegistry"`
	UseLocalCache        bool              `json:"useLocalCache"`
	ConflictThreshold    models.Severity   `json:"conflictThreshold"`
	AutoResolveConflicts bool              `json:"autoResolveConflicts"`
}

// ResolutionStrategy defines how conflicts should be resolved
type ResolutionStrategy string

const (
	StrategyLatest      ResolutionStrategy = "latest"
	StrategyOldest      ResolutionStrategy = "oldest"
	StrategyMajor       ResolutionStrategy = "major"
	StrategyMinor       ResolutionStrategy = "minor"
	StrategyPatch       ResolutionStrategy = "patch"
	StrategyWorkspace   ResolutionStrategy = "workspace"
	StrategyUpgradeAll  ResolutionStrategy = "upgrade_all"
	StrategyManual      ResolutionStrategy = "manual"
)


// ResolverCache provides caching for resolved packages
type ResolverCache struct {
	packageCache map[string]*CacheEntry `json:"packageCache"`
	treesCache   map[string]*DependencyTree `json:"treesCache"`
	mutex        sync.RWMutex
	ttl          time.Duration
}

// CacheEntry represents a cached package resolution
type CacheEntry struct {
	Package   *PackageNode `json:"package"`
	CachedAt  time.Time    `json:"cachedAt"`
	ExpiresAt time.Time    `json:"expiresAt"`
	HitCount  int          `json:"hitCount"`
}

// PackageRegistry provides access to package registry information
type PackageRegistry struct {
	baseURL      string
	client       interface{} // HTTP client would go here
	rateLimiter  interface{} // Rate limiter would go here
	cache        map[string]*RegistryInfo
	cacheMutex   sync.RWMutex
}

// RegistryInfo contains package information from the registry
type RegistryInfo struct {
	Name        string                 `json:"name"`
	Versions    []string               `json:"versions"`
	LatestVersion string               `json:"latest"`
	Tags        map[string]string      `json:"dist-tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	FetchedAt   time.Time              `json:"fetchedAt"`
}

// ResolverMetrics tracks performance metrics
type ResolverMetrics struct {
	CacheHits      int64         `json:"cacheHits"`
	CacheMisses    int64         `json:"cacheMisses"`
	RegistryHits   int64         `json:"registryHits"`
	TotalRequests  int64         `json:"totalRequests"`
	AverageLatency time.Duration `json:"averageLatency"`
	ErrorCount     int64         `json:"errorCount"`
	mutex          sync.RWMutex
}

// NewDependencyTreeResolver creates a new dependency tree resolver
func NewDependencyTreeResolver(logger *logrus.Logger) *DependencyTreeResolver {
	return &DependencyTreeResolver{
		logger: logger,
		cache: &ResolverCache{
			packageCache: make(map[string]*CacheEntry),
			treesCache:   make(map[string]*DependencyTree),
			ttl:          time.Hour * 2, // 2 hour cache TTL
		},
		registry: &PackageRegistry{
			baseURL: "https://registry.npmjs.org",
			cache:   make(map[string]*RegistryInfo),
		},
		metrics: &ResolverMetrics{},
	}
}

// BuildTree builds a complete dependency tree for the given packages
func (dtr *DependencyTreeResolver) BuildTree(ctx context.Context, packages []*PackageInfo, options BuildOptions) (*DependencyTree, error) {
	startTime := time.Now()
	
	dtr.logger.WithFields(logrus.Fields{
		"package_count":    len(packages),
		"max_depth":        options.MaxDepth,
		"include_dev_deps": options.IncludeDevDeps,
		"strategy":         options.Strategy,
	}).Info("Starting dependency tree build")

	// Check cache first
	if options.EnableCaching {
		cacheKey := dtr.generateCacheKey(packages, options)
		if cachedTree := dtr.getCachedTree(cacheKey); cachedTree != nil {
			dtr.logger.Debug("Returning cached dependency tree")
			return cachedTree, nil
		}
	}

	// Initialize the tree
	tree := &DependencyTree{
		RootPackages: make([]*PackageNode, 0),
		AllPackages:  make(map[string]*PackageNode),
		Conflicts:    make([]*EnhancedConflict, 0),
		Metadata: &TreeMetadata{
			Options: options,
		},
	}

	// Build workspace package map for quick lookup
	workspacePackages := make(map[string]*PackageInfo)
	for _, pkg := range packages {
		workspacePackages[pkg.PackageJSON.Name] = pkg
	}

	// Resolve each root package concurrently
	var wg sync.WaitGroup
	rootNodesChan := make(chan *PackageNode, len(packages))
	errorsChan := make(chan error, len(packages))
	
	// Use semaphore for concurrency control
	sem := make(chan struct{}, options.ConcurrencyLevel)

	for _, pkg := range packages {
		wg.Add(1)
		go func(pkg *PackageInfo) {
			defer wg.Done()
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			rootNode, err := dtr.resolvePackageNode(ctx, pkg.PackageJSON.Name, pkg.PackageJSON.Version, "", 0, workspacePackages, options)
			if err != nil {
				errorsChan <- fmt.Errorf("failed to resolve root package %s: %w", pkg.PackageJSON.Name, err)
				return
			}

			rootNode.IsWorkspace = true
			rootNode.Path = pkg.Path
			rootNodesChan <- rootNode
		}(pkg)
	}

	// Wait for all resolutions to complete
	wg.Wait()
	close(rootNodesChan)
	close(errorsChan)

	// Check for errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 && len(errors) == len(packages) {
		return nil, fmt.Errorf("failed to resolve any packages: %v", errors)
	}

	// Collect root nodes
	for rootNode := range rootNodesChan {
		tree.RootPackages = append(tree.RootPackages, rootNode)
		dtr.collectAllPackages(rootNode, tree.AllPackages)
	}

	// Count external vs workspace packages
	workspaceCount := 0
	for _, node := range tree.AllPackages {
		if node.IsWorkspace {
			workspaceCount++
		}
	}
	tree.Metadata.WorkspacePackages = workspaceCount
	tree.Metadata.ExternalPackages = len(tree.AllPackages) - workspaceCount

	// Detect and analyze conflicts
	conflicts := dtr.detectConflicts(tree.AllPackages)
	tree.Conflicts = dtr.enhanceConflicts(conflicts, tree.AllPackages, options)

	// Auto-resolve conflicts if enabled
	if options.AutoResolveConflicts {
		resolvedCount := dtr.autoResolveConflicts(tree, options)
		tree.Metadata.ResolvedConflicts = resolvedCount
	}

	// Calculate metadata
	tree.Metadata.BuildTime = time.Since(startTime)
	tree.Metadata.TotalNodes = len(tree.AllPackages)
	tree.Metadata.MaxDepth = dtr.calculateMaxDepth(tree.RootPackages)
	tree.Metadata.ConflictCount = len(tree.Conflicts)
	tree.Metadata.CacheHitRate = dtr.calculateCacheHitRate()
	tree.Metadata.RegistryHitRate = dtr.calculateRegistryHitRate()

	// Cache the result
	if options.EnableCaching {
		cacheKey := dtr.generateCacheKey(packages, options)
		dtr.cacheTree(cacheKey, tree)
	}

	dtr.logger.WithFields(logrus.Fields{
		"build_time":     tree.Metadata.BuildTime,
		"total_nodes":    tree.Metadata.TotalNodes,
		"max_depth":      tree.Metadata.MaxDepth,
		"conflicts":      tree.Metadata.ConflictCount,
		"cache_hit_rate": tree.Metadata.CacheHitRate,
	}).Info("Dependency tree build completed")

	return tree, nil
}

// resolvePackageNode resolves a single package node and its dependencies
func (dtr *DependencyTreeResolver) resolvePackageNode(ctx context.Context, name, version, requestedRange string, depth int, workspacePackages map[string]*PackageInfo, options BuildOptions) (*PackageNode, error) {
	// Check depth limit
	if depth > options.MaxDepth {
		return nil, fmt.Errorf("maximum depth exceeded for package %s", name)
	}

	// Check cache first
	cacheKey := fmt.Sprintf("%s@%s:%s", name, version, requestedRange)
	if options.EnableCaching {
		if cached := dtr.getCachedPackage(cacheKey); cached != nil {
			dtr.metrics.incrementCacheHits()
			return cached, nil
		}
	}

	dtr.metrics.incrementCacheMisses()

	// Create the package node
	node := &PackageNode{
		Name:           name,
		Version:        version,
		RequestedRange: requestedRange,
		Dependencies:   make(map[string]*PackageNode),
		Depth:          depth,
		Metadata:       make(map[string]interface{}),
	}

	// Parse semantic version
	versionParser := NewVersionParser()
	if resolvedVersion, err := versionParser.ParseVersion(version); err == nil {
		node.ResolvedVersion = resolvedVersion
	}

	// Prefer workspace packages
	if options.PreferWorkspace {
		if workspacePkg, exists := workspacePackages[name]; exists {
			node.IsWorkspace = true
			node.Path = workspacePkg.Path
			node.ResolutionSource = "workspace"
			node.Version = workspacePkg.PackageJSON.Version
			
			// Resolve workspace dependencies
			return dtr.resolveWorkspaceDependencies(ctx, node, workspacePkg, workspacePackages, options)
		}
	}

	// Resolve from registry (simplified implementation)
	node.ResolutionSource = "registry"
	
	// Cache the result
	if options.EnableCaching {
		dtr.cachePackage(cacheKey, node)
	}

	return node, nil
}

// resolveWorkspaceDependencies resolves dependencies for a workspace package
func (dtr *DependencyTreeResolver) resolveWorkspaceDependencies(ctx context.Context, node *PackageNode, pkg *PackageInfo, workspacePackages map[string]*PackageInfo, options BuildOptions) (*PackageNode, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	
	// Resolve production dependencies
	for depName, depVersion := range pkg.PackageJSON.Dependencies {
		wg.Add(1)
		go func(name, version string) {
			defer wg.Done()
			
			depNode, err := dtr.resolvePackageNode(ctx, name, version, version, node.Depth+1, workspacePackages, options)
			if err != nil {
				dtr.logger.WithError(err).WithFields(logrus.Fields{
					"package": node.Name,
					"dependency": name,
				}).Warn("Failed to resolve dependency")
				return
			}

			mutex.Lock()
			node.Dependencies[name] = depNode
			mutex.Unlock()
		}(depName, depVersion)
	}

	// Resolve dev dependencies if enabled
	if options.IncludeDevDeps {
		node.DevDependencies = make(map[string]*PackageNode)
		for depName, depVersion := range pkg.PackageJSON.DevDependencies {
			wg.Add(1)
			go func(name, version string) {
				defer wg.Done()
				
				depNode, err := dtr.resolvePackageNode(ctx, name, version, version, node.Depth+1, workspacePackages, options)
				if err != nil {
					dtr.logger.WithError(err).WithFields(logrus.Fields{
						"package": node.Name,
						"dev_dependency": name,
					}).Warn("Failed to resolve dev dependency")
					return
				}

				mutex.Lock()
				node.DevDependencies[name] = depNode
				mutex.Unlock()
			}(depName, depVersion)
		}
	}

	// Resolve peer dependencies if enabled
	if options.IncludePeerDeps && len(pkg.PackageJSON.PeerDependencies) > 0 {
		node.PeerDependencies = make(map[string]*PackageNode)
		for depName, depVersion := range pkg.PackageJSON.PeerDependencies {
			wg.Add(1)
			go func(name, version string) {
				defer wg.Done()
				
				depNode, err := dtr.resolvePackageNode(ctx, name, version, version, node.Depth+1, workspacePackages, options)
				if err != nil {
					dtr.logger.WithError(err).WithFields(logrus.Fields{
						"package": node.Name,
						"peer_dependency": name,
					}).Warn("Failed to resolve peer dependency")
					return
				}

				mutex.Lock()
				node.PeerDependencies[name] = depNode
				mutex.Unlock()
			}(depName, depVersion)
		}
	}

	wg.Wait()
	return node, nil
}

// detectConflicts detects version conflicts in the resolved tree
func (dtr *DependencyTreeResolver) detectConflicts(allPackages map[string]*PackageNode) []*models.VersionConflict {
	conflictMap := make(map[string]map[string][]*PackageNode)
	
	// Group packages by name and version
	for _, node := range allPackages {
		if conflictMap[node.Name] == nil {
			conflictMap[node.Name] = make(map[string][]*PackageNode)
		}
		conflictMap[node.Name][node.Version] = append(conflictMap[node.Name][node.Version], node)
	}

	var conflicts []*models.VersionConflict
	versionParser := NewVersionParser()

	// Identify conflicts
	for packageName, versions := range conflictMap {
		if len(versions) > 1 {
			var conflictingVersions []models.ConflictingVersion
			var allPackageNames []string
			
			for version, nodes := range versions {
				packageNames := make([]string, len(nodes))
				for i, node := range nodes {
					packageNames[i] = node.Name
				}
				
				// Check if this is a breaking change
				isBreaking := false
				if len(versions) > 1 {
					// Simplified breaking change detection
					for otherVersion := range versions {
						if otherVersion != version {
							if v1, err1 := versionParser.ParseVersion(version); err1 == nil {
								if v2, err2 := versionParser.ParseVersion(otherVersion); err2 == nil {
									if v1.Major != v2.Major {
										isBreaking = true
										break
									}
								}
							}
						}
					}
				}

				conflictingVersions = append(conflictingVersions, models.ConflictingVersion{
					Version:    version,
					Packages:   packageNames,
					IsBreaking: isBreaking,
				})
				
				allPackageNames = append(allPackageNames, packageNames...)
			}

			// Calculate risk level
			riskLevel := models.RiskLevelMedium
			for _, cv := range conflictingVersions {
				if cv.IsBreaking {
					riskLevel = models.RiskLevelHigh
					break
				}
			}

			conflicts = append(conflicts, &models.VersionConflict{
				PackageName:         packageName,
				ConflictingVersions: conflictingVersions,
				RiskLevel:           riskLevel,
				Resolution:          fmt.Sprintf("Align all packages to use the same version of %s", packageName),
				Impact:              fmt.Sprintf("Affects %d packages in the monorepo", len(allPackageNames)),
			})
		}
	}

	return conflicts
}

// enhanceConflicts enhances basic conflicts with additional resolution information
func (dtr *DependencyTreeResolver) enhanceConflicts(conflicts []*models.VersionConflict, allPackages map[string]*PackageNode, options BuildOptions) []*EnhancedConflict {
	var enhanced []*EnhancedConflict

	for _, conflict := range conflicts {
		enhancedConflict := &EnhancedConflict{
			VersionConflict: *conflict,
			PackageNodes:    make([]*PackageNode, 0),
			ResolutionOptions: make([]*ResolutionOption, 0),
			AutoResolvable:  false,
		}

		// Collect relevant package nodes
		for _, node := range allPackages {
			if node.Name == conflict.PackageName {
				enhancedConflict.PackageNodes = append(enhancedConflict.PackageNodes, node)
			}
		}

		// Generate resolution options
		enhancedConflict.ResolutionOptions = dtr.generateResolutionOptions(conflict, enhancedConflict.PackageNodes, options)
		
		// Select recommended fix
		if len(enhancedConflict.ResolutionOptions) > 0 {
			enhancedConflict.RecommendedFix = enhancedConflict.ResolutionOptions[0]
			enhancedConflict.AutoResolvable = enhancedConflict.RecommendedFix.RiskAssessment.Level == "low"
		}

		enhanced = append(enhanced, enhancedConflict)
	}

	return enhanced
}

// generateResolutionOptions generates possible resolution options for conflicts
func (dtr *DependencyTreeResolver) generateResolutionOptions(conflict *models.VersionConflict, nodes []*PackageNode, options BuildOptions) []*ResolutionOption {
	var resolutionOptions []*ResolutionOption
	versionParser := NewVersionParser()

	// Collect all versions
	var versions []*SemanticVersion
	for _, cv := range conflict.ConflictingVersions {
		if version, err := versionParser.ParseVersion(cv.Version); err == nil {
			versions = append(versions, version)
		}
	}

	if len(versions) == 0 {
		return resolutionOptions
	}

	// Strategy 1: Use latest version
	latestVersion := versionParser.findLatestVersion(versions)
	if latestVersion != nil {
		resolutionOptions = append(resolutionOptions, &ResolutionOption{
			Strategy:      StrategyLatest,
			TargetVersion: latestVersion,
			RiskAssessment: &ConflictRiskAssessment{
				Level:      "medium",
				Impact:     "May introduce new features or breaking changes",
				Difficulty: "moderate",
			},
		})
	}

	// Strategy 2: Use oldest stable version
	oldestVersion := versions[0]
	for _, v := range versions[1:] {
		if versionParser.compareVersions(v, oldestVersion) < 0 {
			oldestVersion = v
		}
	}
	
	resolutionOptions = append(resolutionOptions, &ResolutionOption{
		Strategy:      StrategyOldest,
		TargetVersion: oldestVersion,
		RiskAssessment: &ConflictRiskAssessment{
			Level:      "low",
			Impact:     "Conservative approach, minimal risk",
			Difficulty: "easy",
		},
	})

	return resolutionOptions
}

// Helper methods for cache management and metrics

func (dtr *DependencyTreeResolver) generateCacheKey(packages []*PackageInfo, options BuildOptions) string {
	// Simplified cache key generation
	data, _ := json.Marshal(struct {
		PackageCount int          `json:"package_count"`
		Options      BuildOptions `json:"options"`
	}{
		PackageCount: len(packages),
		Options:      options,
	})
	return string(data)
}

func (dtr *DependencyTreeResolver) getCachedTree(cacheKey string) *DependencyTree {
	dtr.cache.mutex.RLock()
	defer dtr.cache.mutex.RUnlock()
	
	if tree, exists := dtr.cache.treesCache[cacheKey]; exists {
		return tree
	}
	return nil
}

func (dtr *DependencyTreeResolver) cacheTree(cacheKey string, tree *DependencyTree) {
	dtr.cache.mutex.Lock()
	defer dtr.cache.mutex.Unlock()
	
	dtr.cache.treesCache[cacheKey] = tree
}

func (dtr *DependencyTreeResolver) getCachedPackage(cacheKey string) *PackageNode {
	dtr.cache.mutex.RLock()
	defer dtr.cache.mutex.RUnlock()
	
	if entry, exists := dtr.cache.packageCache[cacheKey]; exists {
		if time.Now().Before(entry.ExpiresAt) {
			entry.HitCount++
			return entry.Package
		}
		// Cache expired, remove entry
		delete(dtr.cache.packageCache, cacheKey)
	}
	return nil
}

func (dtr *DependencyTreeResolver) cachePackage(cacheKey string, node *PackageNode) {
	dtr.cache.mutex.Lock()
	defer dtr.cache.mutex.Unlock()
	
	dtr.cache.packageCache[cacheKey] = &CacheEntry{
		Package:   node,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(dtr.cache.ttl),
		HitCount:  0,
	}
}

func (dtr *DependencyTreeResolver) collectAllPackages(node *PackageNode, allPackages map[string]*PackageNode) {
	key := fmt.Sprintf("%s@%s", node.Name, node.Version)
	allPackages[key] = node

	// Recursively collect dependencies
	for _, dep := range node.Dependencies {
		dtr.collectAllPackages(dep, allPackages)
	}

	if node.DevDependencies != nil {
		for _, dep := range node.DevDependencies {
			dtr.collectAllPackages(dep, allPackages)
		}
	}

	if node.PeerDependencies != nil {
		for _, dep := range node.PeerDependencies {
			dtr.collectAllPackages(dep, allPackages)
		}
	}
}

func (dtr *DependencyTreeResolver) calculateMaxDepth(rootNodes []*PackageNode) int {
	maxDepth := 0
	for _, root := range rootNodes {
		depth := dtr.getNodeMaxDepth(root)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

func (dtr *DependencyTreeResolver) getNodeMaxDepth(node *PackageNode) int {
	maxDepth := node.Depth

	for _, dep := range node.Dependencies {
		depth := dtr.getNodeMaxDepth(dep)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

func (dtr *DependencyTreeResolver) autoResolveConflicts(tree *DependencyTree, options BuildOptions) int {
	// Simplified auto-resolution implementation
	resolved := 0
	for _, conflict := range tree.Conflicts {
		if conflict.AutoResolvable && conflict.RecommendedFix != nil {
			// Apply the recommended fix (implementation would modify the tree)
			resolved++
		}
	}
	return resolved
}

// Metrics methods
func (rm *ResolverMetrics) incrementCacheHits() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.CacheHits++
	rm.TotalRequests++
}

func (rm *ResolverMetrics) incrementCacheMisses() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.CacheMisses++
	rm.TotalRequests++
}

func (dtr *DependencyTreeResolver) calculateCacheHitRate() float64 {
	dtr.metrics.mutex.RLock()
	defer dtr.metrics.mutex.RUnlock()
	
	if dtr.metrics.TotalRequests == 0 {
		return 0.0
	}
	return float64(dtr.metrics.CacheHits) / float64(dtr.metrics.TotalRequests)
}

func (dtr *DependencyTreeResolver) calculateRegistryHitRate() float64 {
	dtr.metrics.mutex.RLock()
	defer dtr.metrics.mutex.RUnlock()
	
	if dtr.metrics.TotalRequests == 0 {
		return 0.0
	}
	return float64(dtr.metrics.RegistryHits) / float64(dtr.metrics.TotalRequests)
}