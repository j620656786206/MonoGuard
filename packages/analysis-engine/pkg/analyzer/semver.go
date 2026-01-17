// Package analyzer provides dependency graph analysis for monorepo workspaces.
// This file implements semantic version parsing and comparison for Story 2.4.
package analyzer

import (
	"regexp"
	"strconv"
	"strings"
)

// SemVer represents a parsed semantic version.
type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Raw        string // Original unparsed string
}

// VersionDifference represents the type of difference between two versions.
type VersionDifference int

const (
	// VersionDifferenceNone means the versions are the same.
	VersionDifferenceNone VersionDifference = iota
	// VersionDifferencePatch means only the patch version differs.
	VersionDifferencePatch
	// VersionDifferenceMinor means the minor (or patch) version differs.
	VersionDifferenceMinor
	// VersionDifferenceMajor means the major version differs.
	VersionDifferenceMajor
)

// semverRegex matches semantic version patterns.
// Handles: 1.2.3, 1.2.3-alpha, 1.2.3-beta.1, etc.
var semverRegex = regexp.MustCompile(`^(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:-([\w.-]+))?`)

// rangeRegex matches version range prefixes like ^, ~, >=, <=, >, <, =.
var rangeRegex = regexp.MustCompile(`^[\^~>=<]+\s*`)

// wildcardRegex matches wildcard patterns like x, X, *.
var wildcardRegex = regexp.MustCompile(`^(\d+)\.(?:x|X|\*)$|^(\d+)\.(\d+)\.(?:x|X|\*)$`)

// ParseSemVer parses a version string like "4.17.21" or "^4.17.0".
// Returns nil if the version cannot be parsed.
func ParseSemVer(version string) *SemVer {
	if version == "" {
		return nil
	}

	// Store original
	raw := version

	// Strip range prefixes
	version = StripRange(version)

	// Handle special cases
	version = strings.TrimSpace(version)
	if version == "" || version == "latest" || version == "next" || version == "*" {
		return nil
	}

	// Handle wildcard patterns (e.g., "4.x", "4.17.x")
	if matches := wildcardRegex.FindStringSubmatch(version); len(matches) > 0 {
		// Pattern: X.x or X.Y.x
		if matches[1] != "" {
			// X.x pattern
			major, _ := strconv.Atoi(matches[1])
			return &SemVer{Major: major, Minor: 0, Patch: 0, Raw: raw}
		}
		if matches[2] != "" && matches[3] != "" {
			// X.Y.x pattern
			major, _ := strconv.Atoi(matches[2])
			minor, _ := strconv.Atoi(matches[3])
			return &SemVer{Major: major, Minor: minor, Patch: 0, Raw: raw}
		}
	}

	// Parse standard semver
	matches := semverRegex.FindStringSubmatch(version)
	if matches == nil {
		return nil
	}

	sv := &SemVer{Raw: raw}

	// Major (required)
	if matches[1] != "" {
		sv.Major, _ = strconv.Atoi(matches[1])
	}

	// Minor (optional)
	if matches[2] != "" {
		sv.Minor, _ = strconv.Atoi(matches[2])
	}

	// Patch (optional)
	if matches[3] != "" {
		sv.Patch, _ = strconv.Atoi(matches[3])
	}

	// Prerelease (optional)
	if matches[4] != "" {
		sv.Prerelease = matches[4]
	}

	return sv
}

// StripRange removes semver range prefixes (^, ~, >=, <=, >, <, =, etc.).
func StripRange(version string) string {
	// Handle complex range like ">=4.0.0 <5.0.0" - take first version
	if strings.Contains(version, " ") {
		parts := strings.Fields(version)
		if len(parts) > 0 {
			version = parts[0]
		}
	}

	// Remove range prefixes
	version = rangeRegex.ReplaceAllString(version, "")

	return strings.TrimSpace(version)
}

// CompareVersions returns the type of difference between two versions.
// Returns VersionDifferenceMajor, VersionDifferenceMinor, VersionDifferencePatch, or VersionDifferenceNone.
func CompareVersions(v1, v2 *SemVer) VersionDifference {
	if v1 == nil || v2 == nil {
		// Cannot compare, treat as major difference
		return VersionDifferenceMajor
	}

	// Check major version
	if v1.Major != v2.Major {
		return VersionDifferenceMajor
	}

	// Check minor version
	if v1.Minor != v2.Minor {
		return VersionDifferenceMinor
	}

	// Check patch version
	if v1.Patch != v2.Patch {
		return VersionDifferencePatch
	}

	// Check prerelease
	if v1.Prerelease != v2.Prerelease {
		// Different prereleases of same version are minor differences
		return VersionDifferencePatch
	}

	return VersionDifferenceNone
}

// FindMaxDifference finds the maximum version difference among a list of versions.
func FindMaxDifference(versions []string) VersionDifference {
	if len(versions) < 2 {
		return VersionDifferenceNone
	}

	// Parse all versions
	parsed := make([]*SemVer, 0, len(versions))
	for _, v := range versions {
		if sv := ParseSemVer(v); sv != nil {
			parsed = append(parsed, sv)
		}
	}

	if len(parsed) < 2 {
		return VersionDifferenceNone
	}

	maxDiff := VersionDifferenceNone

	// Compare all pairs
	for i := 0; i < len(parsed); i++ {
		for j := i + 1; j < len(parsed); j++ {
			diff := CompareVersions(parsed[i], parsed[j])
			if diff > maxDiff {
				maxDiff = diff
			}
		}
	}

	return maxDiff
}

// FindHighestVersion finds the highest version from a list of version strings.
// Returns the original version string (with any range prefix).
func FindHighestVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	highest := versions[0]
	highestParsed := ParseSemVer(highest)

	for _, v := range versions[1:] {
		parsed := ParseSemVer(v)
		if parsed == nil {
			continue
		}

		if highestParsed == nil {
			highest = v
			highestParsed = parsed
			continue
		}

		// Compare: higher major > higher minor > higher patch
		if parsed.Major > highestParsed.Major ||
			(parsed.Major == highestParsed.Major && parsed.Minor > highestParsed.Minor) ||
			(parsed.Major == highestParsed.Major && parsed.Minor == highestParsed.Minor && parsed.Patch > highestParsed.Patch) {
			highest = v
			highestParsed = parsed
		}
	}

	return highest
}
