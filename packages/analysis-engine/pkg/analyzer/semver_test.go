package analyzer

import (
	"testing"
)

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantMajor  int
		wantMinor  int
		wantPatch  int
		wantPrerel string
		wantNil    bool
	}{
		// Exact versions
		{"exact full", "4.17.21", 4, 17, 21, "", false},
		{"exact major.minor", "4.17", 4, 17, 0, "", false},
		{"exact major only", "4", 4, 0, 0, "", false},
		{"with prerelease", "1.0.0-alpha", 1, 0, 0, "alpha", false},
		{"with prerelease beta.1", "2.0.0-beta.1", 2, 0, 0, "beta.1", false},

		// Caret ranges
		{"caret full", "^4.17.0", 4, 17, 0, "", false},
		{"caret partial", "^4.17", 4, 17, 0, "", false},
		{"caret major", "^4", 4, 0, 0, "", false},

		// Tilde ranges
		{"tilde full", "~4.17.0", 4, 17, 0, "", false},
		{"tilde partial", "~4.17", 4, 17, 0, "", false},

		// Comparison ranges
		{"gte", ">=4.0.0", 4, 0, 0, "", false},
		{"gt", ">4.0.0", 4, 0, 0, "", false},
		{"lte", "<=4.0.0", 4, 0, 0, "", false},
		{"lt", "<4.0.0", 4, 0, 0, "", false},
		{"eq", "=4.0.0", 4, 0, 0, "", false},

		// Complex ranges (takes first version)
		{"range gte lt", ">=4.0.0 <5.0.0", 4, 0, 0, "", false},

		// Wildcards
		{"wildcard minor", "4.x", 4, 0, 0, "", false},
		{"wildcard patch", "4.17.x", 4, 17, 0, "", false},
		{"wildcard uppercase", "4.X", 4, 0, 0, "", false},

		// Invalid/special cases
		{"empty", "", 0, 0, 0, "", true},
		{"latest", "latest", 0, 0, 0, "", true},
		{"next", "next", 0, 0, 0, "", true},
		{"star", "*", 0, 0, 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseSemVer(tt.input)

			if tt.wantNil {
				if result != nil {
					t.Errorf("ParseSemVer(%q) = %+v, want nil", tt.input, result)
				}
				return
			}

			if result == nil {
				t.Fatalf("ParseSemVer(%q) = nil, want non-nil", tt.input)
			}

			if result.Major != tt.wantMajor {
				t.Errorf("Major = %d, want %d", result.Major, tt.wantMajor)
			}
			if result.Minor != tt.wantMinor {
				t.Errorf("Minor = %d, want %d", result.Minor, tt.wantMinor)
			}
			if result.Patch != tt.wantPatch {
				t.Errorf("Patch = %d, want %d", result.Patch, tt.wantPatch)
			}
			if result.Prerelease != tt.wantPrerel {
				t.Errorf("Prerelease = %q, want %q", result.Prerelease, tt.wantPrerel)
			}
		})
	}
}

func TestStripRange(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"4.17.21", "4.17.21"},
		{"^4.17.0", "4.17.0"},
		{"~4.17.0", "4.17.0"},
		{">=4.0.0", "4.0.0"},
		{">4.0.0", "4.0.0"},
		{"<=4.0.0", "4.0.0"},
		{"<4.0.0", "4.0.0"},
		{"=4.0.0", "4.0.0"},
		{">=4.0.0 <5.0.0", "4.0.0"},
		{"  ^4.17.0  ", "4.17.0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := StripRange(tt.input)
			if result != tt.want {
				t.Errorf("StripRange(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want VersionDifference
	}{
		// Same versions
		{"identical", "4.17.21", "4.17.21", VersionDifferenceNone},
		{"identical with caret", "^4.17.21", "^4.17.21", VersionDifferenceNone},

		// Patch differences
		{"patch diff", "4.17.19", "4.17.21", VersionDifferencePatch},
		{"patch diff caret", "^4.17.19", "^4.17.21", VersionDifferencePatch},

		// Minor differences
		{"minor diff", "4.17.0", "4.18.0", VersionDifferenceMinor},
		{"minor diff with patch", "4.17.21", "4.18.5", VersionDifferenceMinor},

		// Major differences
		{"major diff", "3.0.0", "4.0.0", VersionDifferenceMajor},
		{"major diff large", "1.0.0", "5.0.0", VersionDifferenceMajor},

		// Prerelease differences
		{"prerelease diff", "1.0.0-alpha", "1.0.0-beta", VersionDifferencePatch},
		{"prerelease vs release", "1.0.0-alpha", "1.0.0", VersionDifferencePatch},

		// Nil handling
		{"v1 unparseable", "invalid", "4.0.0", VersionDifferenceMajor},
		{"v2 unparseable", "4.0.0", "invalid", VersionDifferenceMajor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := ParseSemVer(tt.v1)
			v2 := ParseSemVer(tt.v2)
			result := CompareVersions(v1, v2)
			if result != tt.want {
				t.Errorf("CompareVersions(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.want)
			}
		})
	}
}

func TestFindMaxDifference(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		want     VersionDifference
	}{
		{"empty", []string{}, VersionDifferenceNone},
		{"single", []string{"4.17.21"}, VersionDifferenceNone},
		{"identical", []string{"4.17.21", "4.17.21"}, VersionDifferenceNone},
		{"patch only", []string{"4.17.19", "4.17.21"}, VersionDifferencePatch},
		{"minor only", []string{"4.17.0", "4.18.0"}, VersionDifferenceMinor},
		{"major only", []string{"3.0.0", "4.0.0"}, VersionDifferenceMajor},
		{"mixed patch minor", []string{"4.17.19", "4.17.21", "4.18.0"}, VersionDifferenceMinor},
		{"mixed all", []string{"3.0.0", "4.17.19", "4.18.0"}, VersionDifferenceMajor},
		{"three major versions", []string{"1.0.0", "2.0.0", "3.0.0"}, VersionDifferenceMajor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindMaxDifference(tt.versions)
			if result != tt.want {
				t.Errorf("FindMaxDifference(%v) = %d, want %d", tt.versions, result, tt.want)
			}
		})
	}
}

func TestFindHighestVersion(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		want     string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"4.17.21"}, "4.17.21"},
		{"two versions", []string{"4.17.19", "4.17.21"}, "4.17.21"},
		{"three versions", []string{"4.17.19", "4.17.21", "4.17.20"}, "4.17.21"},
		{"major difference", []string{"3.0.0", "4.0.0"}, "4.0.0"},
		{"minor difference", []string{"4.17.0", "4.18.0"}, "4.18.0"},
		{"with caret", []string{"^4.17.19", "^4.17.21"}, "^4.17.21"},
		{"mixed formats", []string{"^4.17.19", "~4.17.21", "4.17.20"}, "~4.17.21"},
		{"descending order", []string{"5.0.0", "4.0.0", "3.0.0"}, "5.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindHighestVersion(tt.versions)
			if result != tt.want {
				t.Errorf("FindHighestVersion(%v) = %q, want %q", tt.versions, result, tt.want)
			}
		})
	}
}

// Test edge cases for version parsing
func TestParseSemVer_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *SemVer)
	}{
		{
			name:  "preserves raw version",
			input: "^4.17.0",
			check: func(t *testing.T, sv *SemVer) {
				if sv.Raw != "^4.17.0" {
					t.Errorf("Raw = %q, want %q", sv.Raw, "^4.17.0")
				}
			},
		},
		{
			name:  "zero version",
			input: "0.0.0",
			check: func(t *testing.T, sv *SemVer) {
				if sv.Major != 0 || sv.Minor != 0 || sv.Patch != 0 {
					t.Errorf("Got %d.%d.%d, want 0.0.0", sv.Major, sv.Minor, sv.Patch)
				}
			},
		},
		{
			name:  "large version numbers",
			input: "100.200.300",
			check: func(t *testing.T, sv *SemVer) {
				if sv.Major != 100 || sv.Minor != 200 || sv.Patch != 300 {
					t.Errorf("Got %d.%d.%d, want 100.200.300", sv.Major, sv.Minor, sv.Patch)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv := ParseSemVer(tt.input)
			if sv == nil {
				t.Fatalf("ParseSemVer(%q) = nil", tt.input)
			}
			tt.check(t, sv)
		})
	}
}
