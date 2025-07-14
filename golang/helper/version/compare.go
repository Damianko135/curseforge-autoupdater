package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a semantic version
type Version struct {
	Major int
	Minor int
	Patch int
	Pre   string
	Build string
}

// Parse parses a version string into a Version struct
func Parse(version string) (*Version, error) {
	// Remove leading 'v' if present
	version = strings.TrimPrefix(version, "v")

	// Regex to match semantic version pattern
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
	matches := re.FindStringSubmatch(version)

	if len(matches) < 4 {
		// Try simpler pattern (major.minor)
		re = regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
		matches = re.FindStringSubmatch(version)

		if len(matches) < 3 {
			return nil, fmt.Errorf("invalid version format: %s", version)
		}
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch := 0
	if len(matches) > 3 && matches[3] != "" {
		patch, err = strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", matches[3])
		}
	}

	pre := ""
	if len(matches) > 4 && matches[4] != "" {
		pre = matches[4]
	}

	build := ""
	if len(matches) > 5 && matches[5] != "" {
		build = matches[5]
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		Pre:   pre,
		Build: build,
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.Pre != "" {
		version += "-" + v.Pre
	}

	if v.Build != "" {
		version += "+" + v.Build
	}

	return version
}

// Compare compares two versions
// Returns:
//
//	-1 if v < other
//	 0 if v == other
//	 1 if v > other
func (v *Version) Compare(other *Version) int {
	// Compare major
	if v.Major < other.Major {
		return -1
	} else if v.Major > other.Major {
		return 1
	}

	// Compare minor
	if v.Minor < other.Minor {
		return -1
	} else if v.Minor > other.Minor {
		return 1
	}

	// Compare patch
	if v.Patch < other.Patch {
		return -1
	} else if v.Patch > other.Patch {
		return 1
	}

	// Compare pre-release
	if v.Pre == "" && other.Pre != "" {
		return 1 // No pre-release is greater than with pre-release
	} else if v.Pre != "" && other.Pre == "" {
		return -1 // With pre-release is less than without pre-release
	} else if v.Pre != "" && other.Pre != "" {
		return strings.Compare(v.Pre, other.Pre)
	}

	// Versions are equal
	return 0
}

// IsNewer checks if this version is newer than the other
func (v *Version) IsNewer(other *Version) bool {
	return v.Compare(other) > 0
}

// IsOlder checks if this version is older than the other
func (v *Version) IsOlder(other *Version) bool {
	return v.Compare(other) < 0
}

// IsEqual checks if this version is equal to the other
func (v *Version) IsEqual(other *Version) bool {
	return v.Compare(other) == 0
}

// CompareVersions compares two version strings
func CompareVersions(v1, v2 string) (int, error) {
	version1, err := Parse(v1)
	if err != nil {
		return 0, fmt.Errorf("failed to parse version %s: %w", v1, err)
	}

	version2, err := Parse(v2)
	if err != nil {
		return 0, fmt.Errorf("failed to parse version %s: %w", v2, err)
	}

	return version1.Compare(version2), nil
}

// IsNewer checks if version1 is newer than version2
func IsNewer(v1, v2 string) (bool, error) {
	cmp, err := CompareVersions(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp > 0, nil
}

// IsOlder checks if version1 is older than version2
func IsOlder(v1, v2 string) (bool, error) {
	cmp, err := CompareVersions(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp < 0, nil
}

// IsEqual checks if version1 is equal to version2
func IsEqual(v1, v2 string) (bool, error) {
	cmp, err := CompareVersions(v1, v2)
	if err != nil {
		return false, err
	}
	return cmp == 0, nil
}

// IsValidVersion checks if a version string is valid
func IsValidVersion(version string) bool {
	_, err := Parse(version)
	return err == nil
}

// GetLatestVersion returns the latest version from a slice of version strings
func GetLatestVersion(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions provided")
	}

	latest := versions[0]
	latestVersion, err := Parse(latest)
	if err != nil {
		return "", fmt.Errorf("failed to parse version %s: %w", latest, err)
	}

	for _, v := range versions[1:] {
		currentVersion, err := Parse(v)
		if err != nil {
			continue // Skip invalid versions
		}

		if currentVersion.IsNewer(latestVersion) {
			latest = v
			latestVersion = currentVersion
		}
	}

	return latest, nil
}

// SortVersions sorts version strings in ascending order
func SortVersions(versions []string) ([]string, error) {
	type versionPair struct {
		original string
		parsed   *Version
	}

	var pairs []versionPair
	for _, v := range versions {
		parsed, err := Parse(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse version %s: %w", v, err)
		}
		pairs = append(pairs, versionPair{original: v, parsed: parsed})
	}

	// Sort pairs by parsed version
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].parsed.IsNewer(pairs[j].parsed) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	// Extract original version strings
	result := make([]string, len(pairs))
	for i, pair := range pairs {
		result[i] = pair.original
	}

	return result, nil
}

// ExtractVersionFromString attempts to extract a version from a string
func ExtractVersionFromString(s string) string {
	// Common version patterns
	patterns := []string{
		`v?(\d+\.\d+\.\d+(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?)`,
		`v?(\d+\.\d+(?:\.\d+)?(?:-[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?)`,
		`(\d+\.\d+\.\d+)`,
		`(\d+\.\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(s)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// IsPreRelease checks if a version is a pre-release
func IsPreRelease(version string) bool {
	v, err := Parse(version)
	if err != nil {
		return false
	}
	return v.Pre != ""
}

// GetMajorVersion returns the major version number
func GetMajorVersion(version string) (int, error) {
	v, err := Parse(version)
	if err != nil {
		return 0, err
	}
	return v.Major, nil
}

// GetMinorVersion returns the minor version number
func GetMinorVersion(version string) (int, error) {
	v, err := Parse(version)
	if err != nil {
		return 0, err
	}
	return v.Minor, nil
}

// GetPatchVersion returns the patch version number
func GetPatchVersion(version string) (int, error) {
	v, err := Parse(version)
	if err != nil {
		return 0, err
	}
	return v.Patch, nil
}
