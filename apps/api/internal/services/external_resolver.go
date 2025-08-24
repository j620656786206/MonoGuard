package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ExternalResolver interface for resolving external package information
type ExternalResolver interface {
	ResolvePackageVersion(name string, vRange *VersionRange) (*SemanticVersion, error)
	GetPackageMetadata(name, version string) (*PackageMetadata, error)
	PackageExists(name, version string) (bool, error)
	GetAvailableVersions(name string) ([]*SemanticVersion, error)
	GetPackageDependencies(name, version string) (map[string]*VersionRange, error)
}

// NPMExternalResolver implements ExternalResolver for npm registry
type NPMExternalResolver struct {
	logger       *logrus.Logger
	client       *http.Client
	registryURL  string
	rateLimiter  chan struct{}
}

// PackageMetadata contains external package metadata
type PackageMetadata struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Description      string            `json:"description"`
	Homepage         string            `json:"homepage"`
	Repository       string            `json:"repository"`
	License          string            `json:"license"`
	Dependencies     map[string]string `json:"dependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Engines          map[string]string `json:"engines"`
	PublishedAt      time.Time         `json:"published_at"`
	UnpackedSize     int64             `json:"unpacked_size"`
	FileCount        int               `json:"file_count"`
	HasTypings       bool              `json:"has_typings"`
	Deprecated       bool              `json:"deprecated"`
	SecurityVulns    []string          `json:"security_vulns"`
}

// NPM Registry response structures
type npmPackageResponse struct {
	Name        string                        `json:"name"`
	Versions    map[string]npmVersionInfo     `json:"versions"`
	DistTags    map[string]string             `json:"dist-tags"`
	Description string                        `json:"description"`
	Homepage    string                        `json:"homepage"`
	Repository  npmRepository                 `json:"repository"`
	License     interface{}                   `json:"license"` // Can be string or object
	Time        map[string]string             `json:"time"`
}

type npmVersionInfo struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Description      string            `json:"description"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
	Engines          map[string]string `json:"engines"`
	License          interface{}       `json:"license"`
	Dist             npmDist           `json:"dist"`
	Deprecated       interface{}       `json:"deprecated"` // Can be string or boolean
	HasShrinkwrap    bool              `json:"hasShinkwrap"`
}

type npmRepository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type npmDist struct {
	Tarball      string `json:"tarball"`
	Shasum       string `json:"shasum"`
	Integrity    string `json:"integrity"`
	FileCount    int    `json:"fileCount"`
	UnpackedSize int64  `json:"unpackedSize"`
}

// NewNPMExternalResolver creates a new NPM external resolver
func NewNPMExternalResolver(logger *logrus.Logger) ExternalResolver {
	// Rate limiter to prevent overwhelming the npm registry (10 requests per second)
	rateLimiter := make(chan struct{}, 10)
	go func() {
		for {
			time.Sleep(100 * time.Millisecond) // 10 requests per second
			select {
			case rateLimiter <- struct{}{}:
			default:
			}
		}
	}()

	return &NPMExternalResolver{
		logger: logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		registryURL: "https://registry.npmjs.org",
		rateLimiter: rateLimiter,
	}
}

// ResolvePackageVersion resolves a package version from the npm registry
func (npr *NPMExternalResolver) ResolvePackageVersion(name string, vRange *VersionRange) (*SemanticVersion, error) {
	// Rate limiting
	<-npr.rateLimiter

	versions, err := npr.GetAvailableVersions(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get available versions for %s: %w", name, err)
	}

	// Find the best matching version
	bestMatch := npr.findBestMatchingVersion(vRange, versions)
	if bestMatch == nil {
		return nil, fmt.Errorf("no matching version found for %s with range %s", name, vRange.Raw)
	}

	return bestMatch, nil
}

// GetPackageMetadata retrieves package metadata from npm registry
func (npr *NPMExternalResolver) GetPackageMetadata(name, version string) (*PackageMetadata, error) {
	// Rate limiting
	<-npr.rateLimiter

	escapedName := url.QueryEscape(name)
	url := fmt.Sprintf("%s/%s", npr.registryURL, escapedName)

	resp, err := npr.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch package metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("package %s not found", name)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("npm registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var npmResp npmPackageResponse
	if err := json.Unmarshal(body, &npmResp); err != nil {
		return nil, fmt.Errorf("failed to parse npm response: %w", err)
	}

	// Find the specific version
	versionInfo, exists := npmResp.Versions[version]
	if !exists {
		return nil, fmt.Errorf("version %s not found for package %s", version, name)
	}

	// Parse license
	license := npr.parseLicense(versionInfo.License)

	// Parse repository URL
	repository := ""
	if npmResp.Repository.URL != "" {
		repository = npmResp.Repository.URL
	}

	// Parse published date
	publishedAt := time.Now()
	if timeStr, exists := npmResp.Time[version]; exists {
		if parsed, err := time.Parse(time.RFC3339, timeStr); err == nil {
			publishedAt = parsed
		}
	}

	// Check for deprecation
	deprecated := false
	if versionInfo.Deprecated != nil {
		deprecated = true
	}

	// Check for TypeScript definitions
	hasTypings := false
	if versionInfo.Dependencies != nil {
		if _, exists := versionInfo.Dependencies["@types/"+name]; exists {
			hasTypings = true
		}
	}
	// Also check if the package itself has types
	if strings.Contains(strings.ToLower(versionInfo.Description), "typescript") ||
		strings.Contains(strings.ToLower(versionInfo.Description), "@types") {
		hasTypings = true
	}

	return &PackageMetadata{
		Name:             name,
		Version:          version,
		Description:      versionInfo.Description,
		Homepage:         npmResp.Homepage,
		Repository:       repository,
		License:          license,
		Dependencies:     versionInfo.Dependencies,
		PeerDependencies: versionInfo.PeerDependencies,
		Engines:          versionInfo.Engines,
		PublishedAt:      publishedAt,
		UnpackedSize:     versionInfo.Dist.UnpackedSize,
		FileCount:        versionInfo.Dist.FileCount,
		HasTypings:       hasTypings,
		Deprecated:       deprecated,
		SecurityVulns:    []string{}, // Would need additional security API calls
	}, nil
}

// PackageExists checks if a package version exists
func (npr *NPMExternalResolver) PackageExists(name, version string) (bool, error) {
	// Rate limiting
	<-npr.rateLimiter

	escapedName := url.QueryEscape(name)
	url := fmt.Sprintf("%s/%s/%s", npr.registryURL, escapedName, version)

	resp, err := npr.client.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

// GetAvailableVersions gets all available versions for a package
func (npr *NPMExternalResolver) GetAvailableVersions(name string) ([]*SemanticVersion, error) {
	// Rate limiting
	<-npr.rateLimiter

	escapedName := url.QueryEscape(name)
	url := fmt.Sprintf("%s/%s", npr.registryURL, escapedName)

	resp, err := npr.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch package versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("package %s not found", name)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("npm registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var npmResp npmPackageResponse
	if err := json.Unmarshal(body, &npmResp); err != nil {
		return nil, fmt.Errorf("failed to parse npm response: %w", err)
	}

	var versions []*SemanticVersion
	for versionStr := range npmResp.Versions {
		version, err := npr.parseVersion(versionStr)
		if err != nil {
			npr.logger.WithError(err).WithField("version", versionStr).Debug("Failed to parse version")
			continue
		}
		versions = append(versions, version)
	}

	// Sort versions
	npr.sortVersions(versions)

	return versions, nil
}

// GetPackageDependencies gets dependencies for a specific package version
func (npr *NPMExternalResolver) GetPackageDependencies(name, version string) (map[string]*VersionRange, error) {
	metadata, err := npr.GetPackageMetadata(name, version)
	if err != nil {
		return nil, err
	}

	dependencies := make(map[string]*VersionRange)
	
	// Convert dependencies to version ranges
	for depName, versionRange := range metadata.Dependencies {
		vRange, err := npr.parseVersionRange(versionRange)
		if err != nil {
			npr.logger.WithError(err).WithFields(logrus.Fields{
				"package":   depName,
				"range":     versionRange,
			}).Debug("Failed to parse dependency version range")
			continue
		}
		dependencies[depName] = vRange
	}

	return dependencies, nil
}

// Helper methods

func (npr *NPMExternalResolver) parseVersion(versionStr string) (*SemanticVersion, error) {
	// Remove any prefixes
	clean := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(versionStr, "v"), "^"), "~")
	
	// Handle pre-release versions
	parts := strings.Split(clean, "-")
	versionPart := parts[0]
	preRelease := ""
	if len(parts) > 1 {
		preRelease = strings.Join(parts[1:], "-")
	}
	
	// Split version parts
	versionParts := strings.Split(versionPart, ".")
	
	var major, minor, patch int
	var err error
	
	if len(versionParts) > 0 {
		major, err = strconv.Atoi(versionParts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid major version: %s", versionParts[0])
		}
	}
	
	if len(versionParts) > 1 {
		minor, err = strconv.Atoi(versionParts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", versionParts[1])
		}
	}
	
	if len(versionParts) > 2 {
		patch, err = strconv.Atoi(versionParts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", versionParts[2])
		}
	}
	
	return &SemanticVersion{
		Raw:        versionStr,
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: preRelease,
	}, nil
}

func (npr *NPMExternalResolver) parseVersionRange(versionRange string) (*VersionRange, error) {
	// Clean the version range
	clean := strings.TrimSpace(versionRange)
	
	// Extract operator
	operator := "="
	version := clean
	
	// Use regex to properly parse complex ranges
	operatorRegex := regexp.MustCompile(`^(\^|~|>=|<=|>|<|=)?(.+)$`)
	matches := operatorRegex.FindStringSubmatch(clean)
	
	if len(matches) >= 3 {
		if matches[1] != "" {
			operator = matches[1]
		}
		version = matches[2]
	}
	
	// Parse version parts
	parts := strings.Split(version, ".")
	var major, minor, patch int
	var preRelease string
	
	if len(parts) > 0 {
		majorStr := strings.Split(parts[0], "-")[0]
		major, _ = strconv.Atoi(majorStr)
	}
	if len(parts) > 1 {
		minorStr := strings.Split(parts[1], "-")[0]
		minor, _ = strconv.Atoi(minorStr)
	}
	if len(parts) > 2 {
		// Handle pre-release versions
		patchParts := strings.Split(parts[2], "-")
		patch, _ = strconv.Atoi(patchParts[0])
		if len(patchParts) > 1 {
			preRelease = strings.Join(patchParts[1:], "-")
		}
	}
	
	parsedVersion := &SemanticVersion{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: preRelease,
		Raw:        version,
	}
	
	return &VersionRange{
		Raw:        clean,
		Operator:   operator,
		Version:    parsedVersion,
	}, nil
}

func (npr *NPMExternalResolver) findBestMatchingVersion(vRange *VersionRange, versions []*SemanticVersion) *SemanticVersion {
	var candidates []*SemanticVersion
	
	for _, version := range versions {
		if npr.versionSatisfiesRange(version, vRange) {
			candidates = append(candidates, version)
		}
	}
	
	if len(candidates) == 0 {
		return nil
	}
	
	// Return the latest matching version
	npr.sortVersions(candidates)
	return candidates[len(candidates)-1]
}

func (npr *NPMExternalResolver) versionSatisfiesRange(version *SemanticVersion, vRange *VersionRange) bool {
	switch vRange.Operator {
	case "^":
		// Compatible within same major version
		return version.Major == vRange.Version.Major && 
			(version.Minor > vRange.Version.Minor || 
				(version.Minor == vRange.Version.Minor && version.Patch >= vRange.Version.Patch))
	case "~":
		// Compatible within same major.minor
		return version.Major == vRange.Version.Major && 
			version.Minor == vRange.Version.Minor && 
			version.Patch >= vRange.Version.Patch
	case ">=":
		return npr.compareVersions(version, &SemanticVersion{
			Major: vRange.Version.Major,
			Minor: vRange.Version.Minor,
			Patch: vRange.Version.Patch,
		}) >= 0
	case "<=":
		return npr.compareVersions(version, &SemanticVersion{
			Major: vRange.Version.Major,
			Minor: vRange.Version.Minor,
			Patch: vRange.Version.Patch,
		}) <= 0
	case ">":
		return npr.compareVersions(version, &SemanticVersion{
			Major: vRange.Version.Major,
			Minor: vRange.Version.Minor,
			Patch: vRange.Version.Patch,
		}) > 0
	case "<":
		return npr.compareVersions(version, &SemanticVersion{
			Major: vRange.Version.Major,
			Minor: vRange.Version.Minor,
			Patch: vRange.Version.Patch,
		}) < 0
	case "=", "":
		return version.Major == vRange.Version.Major && 
			version.Minor == vRange.Version.Minor && 
			version.Patch == vRange.Version.Patch
	default:
		return false
	}
}

func (npr *NPMExternalResolver) compareVersions(v1, v2 *SemanticVersion) int {
	if v1.Major != v2.Major {
		return v1.Major - v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor - v2.Minor
	}
	if v1.Patch != v2.Patch {
		return v1.Patch - v2.Patch
	}
	
	// Handle pre-release versions
	if v1.Prerelease == "" && v2.Prerelease != "" {
		return 1 // Release > pre-release
	}
	if v1.Prerelease != "" && v2.Prerelease == "" {
		return -1 // Pre-release < release
	}
	if v1.Prerelease != "" && v2.Prerelease != "" {
		return strings.Compare(v1.Prerelease, v2.Prerelease)
	}
	
	return 0
}

func (npr *NPMExternalResolver) sortVersions(versions []*SemanticVersion) {
	// Sort in ascending order (oldest first)
	for i := 0; i < len(versions)-1; i++ {
		for j := i + 1; j < len(versions); j++ {
			if npr.compareVersions(versions[i], versions[j]) > 0 {
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}
}

func (npr *NPMExternalResolver) parseLicense(license interface{}) string {
	if license == nil {
		return ""
	}
	
	switch v := license.(type) {
	case string:
		return v
	case map[string]interface{}:
		if licenseType, exists := v["type"]; exists {
			if typeStr, ok := licenseType.(string); ok {
				return typeStr
			}
		}
	}
	
	return ""
}