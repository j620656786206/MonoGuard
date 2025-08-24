package services

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// VersionRange represents a version range constraint
type VersionRange struct {
	Raw        string         `json:"raw"`
	Operator   string         `json:"operator"`
	Version    *SemanticVersion `json:"version"`
	Satisfies  func(*SemanticVersion) bool `json:"-"`
}

// SemanticVersion represents a semantic version
type SemanticVersion struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	Prerelease string `json:"prerelease,omitempty"`
	Build      string `json:"build,omitempty"`
	Raw        string `json:"raw"`
}

// VersionConflictInfo represents detailed information about version conflicts
type VersionConflictInfo struct {
	PackageName     string                    `json:"packageName"`
	ConflictType    string                    `json:"conflictType"` // "major", "minor", "patch", "prerelease"
	Versions        []*SemanticVersion        `json:"versions"`
	Ranges          []*VersionRange           `json:"ranges"`
	IsResolvable    bool                      `json:"isResolvable"`
	SuggestedFix    *SemanticVersion          `json:"suggestedFix,omitempty"`
	RiskAssessment  *ConflictRiskAssessment   `json:"riskAssessment"`
}

// ConflictRiskAssessment represents risk assessment for version conflicts
type ConflictRiskAssessment struct {
	Level       string   `json:"level"` // "low", "medium", "high", "critical"
	Reasons     []string `json:"reasons"`
	Impact      string   `json:"impact"`
	Difficulty  string   `json:"difficulty"` // "easy", "moderate", "hard"
}

// VersionParser handles parsing and analysis of version ranges
type VersionParser struct {
	// Compiled regex patterns for performance
	semverRegex    *regexp.Regexp
	rangeRegex     *regexp.Regexp
	prereleaseRegex *regexp.Regexp
}

// NewVersionParser creates a new version parser with compiled regex patterns
func NewVersionParser() *VersionParser {
	return &VersionParser{
		semverRegex: regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<build>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`),
		rangeRegex: regexp.MustCompile(`^([\^~>=<]*)(.+)$`),
		prereleaseRegex: regexp.MustCompile(`^[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*$`),
	}
}

// ParseVersion parses a version string into a SemanticVersion
func (vp *VersionParser) ParseVersion(versionStr string) (*SemanticVersion, error) {
	if versionStr == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Clean the version string (remove leading 'v')
	cleanVersion := strings.TrimPrefix(versionStr, "v")
	
	matches := vp.semverRegex.FindStringSubmatch(cleanVersion)
	if matches == nil {
		return nil, fmt.Errorf("invalid semantic version: %s", versionStr)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	version := &SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
		Raw:   versionStr,
	}

	// Extract prerelease if present
	if len(matches) > 4 && matches[4] != "" {
		version.Prerelease = matches[4]
	}

	// Extract build metadata if present
	if len(matches) > 5 && matches[5] != "" {
		version.Build = matches[5]
	}

	return version, nil
}

// ParseVersionRange parses a version range string into a VersionRange
func (vp *VersionParser) ParseVersionRange(rangeStr string) (*VersionRange, error) {
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	matches := vp.rangeRegex.FindStringSubmatch(rangeStr)
	if matches == nil {
		return nil, fmt.Errorf("invalid version range: %s", rangeStr)
	}

	operator := matches[1]
	versionStr := matches[2]

	// Default operator is exact match
	if operator == "" {
		operator = "="
	}

	version, err := vp.ParseVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version in range: %w", err)
	}

	versionRange := &VersionRange{
		Raw:      rangeStr,
		Operator: operator,
		Version:  version,
	}

	// Set the satisfies function based on the operator
	versionRange.Satisfies = vp.createSatisfiesFunc(operator, version)

	return versionRange, nil
}

// createSatisfiesFunc creates a function that checks if a version satisfies the range
func (vp *VersionParser) createSatisfiesFunc(operator string, rangeVersion *SemanticVersion) func(*SemanticVersion) bool {
	return func(version *SemanticVersion) bool {
		switch operator {
		case "^":
			return vp.satisfiesCaretRange(version, rangeVersion)
		case "~":
			return vp.satisfiesTildeRange(version, rangeVersion)
		case ">=":
			return vp.compareVersions(version, rangeVersion) >= 0
		case "<=":
			return vp.compareVersions(version, rangeVersion) <= 0
		case ">":
			return vp.compareVersions(version, rangeVersion) > 0
		case "<":
			return vp.compareVersions(version, rangeVersion) < 0
		case "=", "":
			return vp.compareVersions(version, rangeVersion) == 0
		default:
			return false
		}
	}
}

// satisfiesCaretRange checks if version satisfies caret range (^1.2.3)
func (vp *VersionParser) satisfiesCaretRange(version, rangeVersion *SemanticVersion) bool {
	// ^1.2.3 := >=1.2.3 <2.0.0 (Same Major)
	// ^0.2.3 := >=0.2.3 <0.3.0 (Same Major.Minor if Major is 0)
	// ^0.0.3 := >=0.0.3 <0.0.4 (Exact if Major and Minor are 0)
	
	if rangeVersion.Major > 0 {
		return version.Major == rangeVersion.Major && 
			   vp.compareVersions(version, rangeVersion) >= 0
	}
	
	if rangeVersion.Minor > 0 {
		return version.Major == rangeVersion.Major && 
			   version.Minor == rangeVersion.Minor && 
			   vp.compareVersions(version, rangeVersion) >= 0
	}
	
	// Both major and minor are 0
	return version.Major == rangeVersion.Major && 
		   version.Minor == rangeVersion.Minor && 
		   version.Patch == rangeVersion.Patch
}

// satisfiesTildeRange checks if version satisfies tilde range (~1.2.3)
func (vp *VersionParser) satisfiesTildeRange(version, rangeVersion *SemanticVersion) bool {
	// ~1.2.3 := >=1.2.3 <1.3.0 (Same Major.Minor)
	// ~1.2 := >=1.2.0 <1.3.0
	// ~1 := >=1.0.0 <2.0.0
	
	return version.Major == rangeVersion.Major && 
		   version.Minor == rangeVersion.Minor && 
		   vp.compareVersions(version, rangeVersion) >= 0
}

// compareVersions compares two semantic versions
func (vp *VersionParser) compareVersions(v1, v2 *SemanticVersion) int {
	// Compare major
	if v1.Major != v2.Major {
		return v1.Major - v2.Major
	}

	// Compare minor
	if v1.Minor != v2.Minor {
		return v1.Minor - v2.Minor
	}

	// Compare patch
	if v1.Patch != v2.Patch {
		return v1.Patch - v2.Patch
	}

	// Compare prerelease
	return vp.comparePrereleases(v1.Prerelease, v2.Prerelease)
}

// comparePrereleases compares prerelease versions
func (vp *VersionParser) comparePrereleases(pre1, pre2 string) int {
	// No prerelease is greater than any prerelease
	if pre1 == "" && pre2 == "" {
		return 0
	}
	if pre1 == "" {
		return 1
	}
	if pre2 == "" {
		return -1
	}

	// Split prerelease parts
	parts1 := strings.Split(pre1, ".")
	parts2 := strings.Split(pre2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 string
		if i < len(parts1) {
			p1 = parts1[i]
		}
		if i < len(parts2) {
			p2 = parts2[i]
		}

		// Empty part is less than any part
		if p1 == "" && p2 == "" {
			continue
		}
		if p1 == "" {
			return -1
		}
		if p2 == "" {
			return 1
		}

		// Try to parse as numbers
		n1, err1 := strconv.Atoi(p1)
		n2, err2 := strconv.Atoi(p2)

		if err1 == nil && err2 == nil {
			// Both are numbers
			if n1 != n2 {
				return n1 - n2
			}
		} else if err1 == nil {
			// p1 is number, p2 is string (number < string)
			return -1
		} else if err2 == nil {
			// p1 is string, p2 is number (string > number)
			return 1
		} else {
			// Both are strings
			if p1 != p2 {
				if p1 < p2 {
					return -1
				}
				return 1
			}
		}
	}

	return 0
}

// AnalyzeVersionConflicts analyzes version conflicts for a set of dependencies
func (vp *VersionParser) AnalyzeVersionConflicts(dependencies map[string][]string) ([]*VersionConflictInfo, error) {
	var conflicts []*VersionConflictInfo

	for packageName, versionStrs := range dependencies {
		if len(versionStrs) <= 1 {
			continue // No conflict with single version
		}

		conflictInfo, err := vp.analyzePackageVersionConflict(packageName, versionStrs)
		if err != nil {
			continue // Skip packages with parse errors
		}

		if conflictInfo != nil {
			conflicts = append(conflicts, conflictInfo)
		}
	}

	return conflicts, nil
}

// analyzePackageVersionConflict analyzes version conflicts for a single package
func (vp *VersionParser) analyzePackageVersionConflict(packageName string, versionStrs []string) (*VersionConflictInfo, error) {
	var versions []*SemanticVersion
	var ranges []*VersionRange

	// Parse all versions and ranges
	for _, versionStr := range versionStrs {
		if version, err := vp.ParseVersion(versionStr); err == nil {
			versions = append(versions, version)
		}

		if versionRange, err := vp.ParseVersionRange(versionStr); err == nil {
			ranges = append(ranges, versionRange)
		}
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no valid versions found for package %s", packageName)
	}

	// Determine conflict type and severity
	conflictType := vp.determineConflictType(versions)
	isResolvable := vp.isConflictResolvable(ranges)
	suggestedFix := vp.suggestVersionFix(versions, ranges)
	riskAssessment := vp.assessConflictRisk(packageName, versions, conflictType)

	return &VersionConflictInfo{
		PackageName:    packageName,
		ConflictType:   conflictType,
		Versions:       versions,
		Ranges:         ranges,
		IsResolvable:   isResolvable,
		SuggestedFix:   suggestedFix,
		RiskAssessment: riskAssessment,
	}, nil
}

// determineConflictType determines the type of version conflict
func (vp *VersionParser) determineConflictType(versions []*SemanticVersion) string {
	if len(versions) <= 1 {
		return "none"
	}

	// Check for major version differences
	majorVersions := make(map[int]bool)
	minorVersions := make(map[string]bool)
	patchVersions := make(map[string]bool)

	for _, v := range versions {
		majorVersions[v.Major] = true
		minorVersions[fmt.Sprintf("%d.%d", v.Major, v.Minor)] = true
		patchVersions[fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)] = true
	}

	// Check for prerelease differences
	hasPrerelease := false
	for _, v := range versions {
		if v.Prerelease != "" {
			hasPrerelease = true
			break
		}
	}

	if len(majorVersions) > 1 {
		return "major"
	}
	if len(minorVersions) > 1 {
		return "minor"
	}
	if len(patchVersions) > 1 {
		return "patch"
	}
	if hasPrerelease {
		return "prerelease"
	}

	return "none"
}

// isConflictResolvable determines if version ranges can be resolved
func (vp *VersionParser) isConflictResolvable(ranges []*VersionRange) bool {
	if len(ranges) <= 1 {
		return true
	}

	// Generate a test version space and check if there's any version
	// that satisfies all ranges
	testVersions := vp.generateTestVersions(ranges)
	
	for _, testVersion := range testVersions {
		satisfiesAll := true
		for _, versionRange := range ranges {
			if !versionRange.Satisfies(testVersion) {
				satisfiesAll = false
				break
			}
		}
		if satisfiesAll {
			return true
		}
	}

	return false
}

// generateTestVersions generates a set of test versions based on ranges
func (vp *VersionParser) generateTestVersions(ranges []*VersionRange) []*SemanticVersion {
	var testVersions []*SemanticVersion
	
	// Collect all base versions from ranges
	for _, versionRange := range ranges {
		testVersions = append(testVersions, versionRange.Version)
		
		// Generate additional test versions based on range type
		switch versionRange.Operator {
		case "^":
			// Test next major version boundary
			nextMajor := &SemanticVersion{
				Major: versionRange.Version.Major + 1,
				Minor: 0,
				Patch: 0,
			}
			testVersions = append(testVersions, nextMajor)
		case "~":
			// Test next minor version boundary
			nextMinor := &SemanticVersion{
				Major: versionRange.Version.Major,
				Minor: versionRange.Version.Minor + 1,
				Patch: 0,
			}
			testVersions = append(testVersions, nextMinor)
		}
	}

	return vp.deduplicateVersions(testVersions)
}

// suggestVersionFix suggests a version that resolves conflicts
func (vp *VersionParser) suggestVersionFix(versions []*SemanticVersion, ranges []*VersionRange) *SemanticVersion {
	if len(ranges) == 0 {
		return vp.findLatestVersion(versions)
	}

	// Find the highest version that satisfies all ranges
	testVersions := vp.generateTestVersions(ranges)
	
	// Add current versions to test set
	testVersions = append(testVersions, versions...)
	testVersions = vp.deduplicateVersions(testVersions)
	
	// Sort versions in descending order (latest first)
	sort.Slice(testVersions, func(i, j int) bool {
		return vp.compareVersions(testVersions[i], testVersions[j]) > 0
	})

	for _, testVersion := range testVersions {
		satisfiesAll := true
		for _, versionRange := range ranges {
			if !versionRange.Satisfies(testVersion) {
				satisfiesAll = false
				break
			}
		}
		if satisfiesAll {
			return testVersion
		}
	}

	// If no compatible version found, suggest the latest
	return vp.findLatestVersion(versions)
}

// assessConflictRisk assesses the risk of version conflicts
func (vp *VersionParser) assessConflictRisk(packageName string, versions []*SemanticVersion, conflictType string) *ConflictRiskAssessment {
	var level string
	var reasons []string
	var impact string
	var difficulty string

	switch conflictType {
	case "major":
		level = "critical"
		reasons = append(reasons, "Major version differences may introduce breaking changes")
		impact = "High likelihood of runtime errors or API incompatibilities"
		difficulty = "hard"
	case "minor":
		level = "high"
		reasons = append(reasons, "Minor version differences may introduce new features or deprecations")
		impact = "Possible behavioral changes or warnings"
		difficulty = "moderate"
	case "patch":
		level = "medium"
		reasons = append(reasons, "Patch version differences typically contain bug fixes")
		impact = "Generally safe, but may affect bug-dependent code"
		difficulty = "easy"
	case "prerelease":
		level = "high"
		reasons = append(reasons, "Prerelease versions may be unstable")
		impact = "Unpredictable behavior in production"
		difficulty = "moderate"
	default:
		level = "low"
		impact = "Minimal impact expected"
		difficulty = "easy"
	}

	// Additional risk factors
	if len(versions) > 3 {
		level = vp.escalateRiskLevel(level)
		reasons = append(reasons, fmt.Sprintf("Many conflicting versions (%d)", len(versions)))
	}

	// Check for zero major versions (experimental packages)
	hasZeroMajor := false
	for _, v := range versions {
		if v.Major == 0 {
			hasZeroMajor = true
			break
		}
	}
	if hasZeroMajor {
		level = vp.escalateRiskLevel(level)
		reasons = append(reasons, "Package is in experimental phase (0.x.x)")
	}

	return &ConflictRiskAssessment{
		Level:      level,
		Reasons:    reasons,
		Impact:     impact,
		Difficulty: difficulty,
	}
}

// escalateRiskLevel escalates risk level to the next higher level
func (vp *VersionParser) escalateRiskLevel(currentLevel string) string {
	switch currentLevel {
	case "low":
		return "medium"
	case "medium":
		return "high"
	case "high":
		return "critical"
	default:
		return "critical"
	}
}

// findLatestVersion finds the latest version in a slice
func (vp *VersionParser) findLatestVersion(versions []*SemanticVersion) *SemanticVersion {
	if len(versions) == 0 {
		return nil
	}

	latest := versions[0]
	for _, version := range versions[1:] {
		if vp.compareVersions(version, latest) > 0 {
			latest = version
		}
	}

	return latest
}

// deduplicateVersions removes duplicate versions from a slice
func (vp *VersionParser) deduplicateVersions(versions []*SemanticVersion) []*SemanticVersion {
	seen := make(map[string]bool)
	var unique []*SemanticVersion

	for _, version := range versions {
		key := fmt.Sprintf("%d.%d.%d-%s+%s", version.Major, version.Minor, version.Patch, version.Prerelease, version.Build)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, version)
		}
	}

	return unique
}

// String methods for pretty printing
func (sv *SemanticVersion) String() string {
	version := fmt.Sprintf("%d.%d.%d", sv.Major, sv.Minor, sv.Patch)
	if sv.Prerelease != "" {
		version += "-" + sv.Prerelease
	}
	if sv.Build != "" {
		version += "+" + sv.Build
	}
	return version
}

func (vr *VersionRange) String() string {
	return vr.Raw
}

// IsStable returns true if the version is stable (no prerelease)
func (sv *SemanticVersion) IsStable() bool {
	return sv.Prerelease == ""
}

// IsBackwardsCompatible checks if two versions are backwards compatible
func (vp *VersionParser) IsBackwardsCompatible(v1, v2 *SemanticVersion) bool {
	// Major version changes are never backwards compatible
	if v1.Major != v2.Major {
		return false
	}

	// For 0.x.x versions, minor changes are not backwards compatible
	if v1.Major == 0 && v1.Minor != v2.Minor {
		return false
	}

	// Minor and patch increases are generally backwards compatible
	return vp.compareVersions(v1, v2) <= 0
}