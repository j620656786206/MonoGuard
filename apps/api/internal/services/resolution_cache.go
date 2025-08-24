package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ResolutionCache provides multi-level caching for dependency resolution
type ResolutionCache struct {
	logger    *logrus.Logger
	l1Cache   *sync.Map                    // In-memory, most recent
	l2Cache   map[string]*CachedResolution // Persistent cache
	l3Cache   string                       // Disk-based cache path
	maxL1Size int
	maxL2Size int
	ttl       time.Duration
	stats     *CacheStats
	mutex     sync.RWMutex
}

// TreeNode represents a dependency tree node for caching
type TreeNode struct {
	Name               string                 `json:"name"`
	RequestedRange     string                 `json:"requestedRange"`
	ResolvedVersion    *SemanticVersion       `json:"resolvedVersion"`
	PackageInfo        *PackageInfo           `json:"packageInfo,omitempty"`
	Depth              int                    `json:"depth"`
	ResolutionPath     []string               `json:"resolutionPath"`
	IsWorkspacePackage bool                   `json:"isWorkspacePackage"`
	IsDevDependency    bool                   `json:"isDevDependency"`
	IsOptional         bool                   `json:"isOptional"`
	HasConflict        bool                   `json:"hasConflict"`
	ConflictInfo       *VersionConflictInfo   `json:"conflictInfo,omitempty"`
	Children           []*TreeNode            `json:"children"`
}

// CachedResolution represents a cached resolution result
type CachedResolution struct {
	Key         string        `json:"key"`
	Result      *TreeNode     `json:"result"`
	CreatedAt   time.Time     `json:"created_at"`
	AccessCount int           `json:"access_count"`
	InputHash   string        `json:"input_hash"`
	Invalidated bool          `json:"invalidated"`
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	L1Hits        int64 `json:"l1_hits"`
	L2Hits        int64 `json:"l2_hits"`
	L3Hits        int64 `json:"l3_hits"`
	Misses        int64 `json:"misses"`
	Evictions     int64 `json:"evictions"`
	Invalidations int64 `json:"invalidations"`
}

// NewResolutionCache creates a new resolution cache
func NewResolutionCache(logger *logrus.Logger) *ResolutionCache {
	cacheDir := filepath.Join(os.TempDir(), "monoguard-cache")
	os.MkdirAll(cacheDir, 0755)

	cache := &ResolutionCache{
		logger:    logger,
		l1Cache:   &sync.Map{},
		l2Cache:   make(map[string]*CachedResolution),
		l3Cache:   cacheDir,
		maxL1Size: 1000,  // Maximum entries in L1 cache
		maxL2Size: 10000, // Maximum entries in L2 cache
		ttl:       15 * time.Minute,
		stats:     &CacheStats{},
	}

	// Load L2 cache from disk on startup
	cache.loadL2Cache()

	// Start cache maintenance goroutine
	go cache.maintenanceLoop()

	return cache
}

// Get retrieves a cached resolution result
func (rc *ResolutionCache) Get(key string) (*TreeNode, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	// Try L1 cache first
	if value, exists := rc.l1Cache.Load(key); exists {
		if cached, ok := value.(*CachedResolution); ok {
			if !rc.isExpired(cached) && !cached.Invalidated {
				cached.AccessCount++
				rc.stats.L1Hits++
				rc.logger.WithField("key", key).Debug("L1 cache hit")
				return rc.cloneCachedNode(cached.Result), true
			} else {
				rc.l1Cache.Delete(key)
			}
		}
	}

	// Try L2 cache
	if cached, exists := rc.l2Cache[key]; exists {
		if !rc.isExpired(cached) && !cached.Invalidated {
			cached.AccessCount++
			rc.stats.L2Hits++
			// Promote to L1 cache
			rc.l1Cache.Store(key, cached)
			rc.logger.WithField("key", key).Debug("L2 cache hit, promoted to L1")
			return rc.cloneCachedNode(cached.Result), true
		} else {
			delete(rc.l2Cache, key)
		}
	}

	// Try L3 cache (disk)
	if node := rc.loadFromDisk(key); node != nil {
		cached := &CachedResolution{
			Key:         key,
			Result:      node,
			CreatedAt:   time.Now(),
			AccessCount: 1,
		}
		rc.stats.L3Hits++
		// Promote to L2 and L1
		rc.l2Cache[key] = cached
		rc.l1Cache.Store(key, cached)
		rc.logger.WithField("key", key).Debug("L3 cache hit, promoted to L1/L2")
		return rc.cloneCachedNode(node), true
	}

	rc.stats.Misses++
	return nil, false
}

// Set stores a resolution result in the cache
func (rc *ResolutionCache) Set(key string, node *TreeNode, inputHash string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	cached := &CachedResolution{
		Key:         key,
		Result:      rc.cloneCachedNode(node),
		CreatedAt:   time.Now(),
		AccessCount: 0,
		InputHash:   inputHash,
		Invalidated: false,
	}

	// Store in L1 cache
	rc.l1Cache.Store(key, cached)

	// Store in L2 cache with size limit
	if len(rc.l2Cache) >= rc.maxL2Size {
		rc.evictLeastRecentlyUsed()
	}
	rc.l2Cache[key] = cached

	// Async store to disk
	go rc.saveToDisk(key, node)

	rc.logger.WithField("key", key).Debug("Cached resolution result")
}

// Invalidate marks cache entries as invalid based on input hash
func (rc *ResolutionCache) Invalidate(inputHash string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	invalidated := 0

	// Invalidate L1 cache
	rc.l1Cache.Range(func(key, value interface{}) bool {
		if cached, ok := value.(*CachedResolution); ok {
			if cached.InputHash == inputHash {
				cached.Invalidated = true
				invalidated++
			}
		}
		return true
	})

	// Invalidate L2 cache
	for _, cached := range rc.l2Cache {
		if cached.InputHash == inputHash {
			cached.Invalidated = true
			invalidated++
		}
	}

	rc.stats.Invalidations += int64(invalidated)
	rc.logger.WithFields(logrus.Fields{
		"input_hash": inputHash,
		"invalidated": invalidated,
	}).Debug("Invalidated cache entries")
}

// Clear clears all cache levels
func (rc *ResolutionCache) Clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.l1Cache = &sync.Map{}
	rc.l2Cache = make(map[string]*CachedResolution)
	
	// Clear disk cache
	os.RemoveAll(rc.l3Cache)
	os.MkdirAll(rc.l3Cache, 0755)

	rc.stats = &CacheStats{}
	rc.logger.Info("Cleared all cache levels")
}

// GetStats returns cache statistics
func (rc *ResolutionCache) GetStats() *CacheStats {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	return &CacheStats{
		L1Hits:        rc.stats.L1Hits,
		L2Hits:        rc.stats.L2Hits,
		L3Hits:        rc.stats.L3Hits,
		Misses:        rc.stats.Misses,
		Evictions:     rc.stats.Evictions,
		Invalidations: rc.stats.Invalidations,
	}
}

// GenerateKey generates a cache key from package name and version range
func (rc *ResolutionCache) GenerateKey(packageName string, versionRange *VersionRange) string {
	data := fmt.Sprintf("%s@%s", packageName, versionRange.Raw)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)[:16] // Use first 16 chars of hash
}

// GenerateInputHash generates a hash from input packages for invalidation
func (rc *ResolutionCache) GenerateInputHash(packages []*PackageInfo) string {
	data := ""
	for _, pkg := range packages {
		data += fmt.Sprintf("%s:%s;", pkg.PackageJSON.Name, pkg.PackageJSON.Version)
		for dep, version := range pkg.PackageJSON.Dependencies {
			data += fmt.Sprintf("%s:%s;", dep, version)
		}
	}
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// Private helper methods

func (rc *ResolutionCache) isExpired(cached *CachedResolution) bool {
	return time.Since(cached.CreatedAt) > rc.ttl
}

func (rc *ResolutionCache) cloneCachedNode(original *TreeNode) *TreeNode {
	if original == nil {
		return nil
	}

	clone := &TreeNode{
		Name:               original.Name,
		RequestedRange:     original.RequestedRange,
		ResolvedVersion:    original.ResolvedVersion,
		PackageInfo:        original.PackageInfo,
		Depth:              original.Depth,
		ResolutionPath:     make([]string, len(original.ResolutionPath)),
		IsWorkspacePackage: original.IsWorkspacePackage,
		IsDevDependency:    original.IsDevDependency,
		IsOptional:         original.IsOptional,
		HasConflict:        original.HasConflict,
		ConflictInfo:       original.ConflictInfo,
		Children:           make([]*TreeNode, 0), // Don't clone children to avoid deep recursion
	}

	copy(clone.ResolutionPath, original.ResolutionPath)
	return clone
}

func (rc *ResolutionCache) evictLeastRecentlyUsed() {
	// Find the least recently used entry
	var oldestKey string
	var oldestTime time.Time
	var oldestAccess int

	first := true
	for key, cached := range rc.l2Cache {
		if first || cached.CreatedAt.Before(oldestTime) || 
			(cached.CreatedAt.Equal(oldestTime) && cached.AccessCount < oldestAccess) {
			oldestKey = key
			oldestTime = cached.CreatedAt
			oldestAccess = cached.AccessCount
			first = false
		}
	}

	if oldestKey != "" {
		delete(rc.l2Cache, oldestKey)
		rc.l1Cache.Delete(oldestKey)
		rc.stats.Evictions++
		rc.logger.WithField("key", oldestKey).Debug("Evicted cache entry")
	}
}

func (rc *ResolutionCache) loadL2Cache() {
	cachePath := filepath.Join(rc.l3Cache, "l2cache.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			rc.logger.WithError(err).Warn("Failed to load L2 cache from disk")
		}
		return
	}

	var diskCache map[string]*CachedResolution
	if err := json.Unmarshal(data, &diskCache); err != nil {
		rc.logger.WithError(err).Warn("Failed to parse L2 cache data")
		return
	}

	// Filter out expired entries
	for key, cached := range diskCache {
		if !rc.isExpired(cached) && !cached.Invalidated {
			rc.l2Cache[key] = cached
		}
	}

	rc.logger.WithField("entries", len(rc.l2Cache)).Info("Loaded L2 cache from disk")
}

func (rc *ResolutionCache) saveL2Cache() {
	cachePath := filepath.Join(rc.l3Cache, "l2cache.json")
	
	// Only save non-expired entries
	validCache := make(map[string]*CachedResolution)
	for key, cached := range rc.l2Cache {
		if !rc.isExpired(cached) && !cached.Invalidated {
			validCache[key] = cached
		}
	}

	data, err := json.Marshal(validCache)
	if err != nil {
		rc.logger.WithError(err).Warn("Failed to marshal L2 cache data")
		return
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		rc.logger.WithError(err).Warn("Failed to save L2 cache to disk")
		return
	}

	rc.logger.WithField("entries", len(validCache)).Debug("Saved L2 cache to disk")
}

func (rc *ResolutionCache) loadFromDisk(key string) *TreeNode {
	filePath := filepath.Join(rc.l3Cache, key+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}

	var node TreeNode
	if err := json.Unmarshal(data, &node); err != nil {
		rc.logger.WithError(err).WithField("key", key).Debug("Failed to unmarshal cached node")
		return nil
	}

	return &node
}

func (rc *ResolutionCache) saveToDisk(key string, node *TreeNode) {
	filePath := filepath.Join(rc.l3Cache, key+".json")
	data, err := json.Marshal(node)
	if err != nil {
		rc.logger.WithError(err).WithField("key", key).Debug("Failed to marshal node for disk cache")
		return
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		rc.logger.WithError(err).WithField("key", key).Debug("Failed to save node to disk cache")
	}
}

func (rc *ResolutionCache) maintenanceLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rc.performMaintenance()
	}
}

func (rc *ResolutionCache) performMaintenance() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Clean expired entries from L1 cache
	rc.l1Cache.Range(func(key, value interface{}) bool {
		if cached, ok := value.(*CachedResolution); ok {
			if rc.isExpired(cached) || cached.Invalidated {
				rc.l1Cache.Delete(key)
			}
		}
		return true
	})

	// Clean expired entries from L2 cache
	for key, cached := range rc.l2Cache {
		if rc.isExpired(cached) || cached.Invalidated {
			delete(rc.l2Cache, key)
		}
	}

	// Save L2 cache to disk
	rc.saveL2Cache()

	// Clean expired files from L3 cache
	rc.cleanDiskCache()

	rc.logger.Debug("Performed cache maintenance")
}

func (rc *ResolutionCache) cleanDiskCache() {
	entries, err := os.ReadDir(rc.l3Cache)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(rc.l3Cache, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Remove files older than TTL
		if time.Since(info.ModTime()) > rc.ttl {
			os.Remove(filePath)
		}
	}
}