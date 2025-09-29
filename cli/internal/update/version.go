package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	githubAPIURL = "https://api.github.com/repos/easel/ddx/releases/latest"
)

// FetchLatestRelease fetches the latest release information from GitHub
func FetchLatestRelease() (*GitHubRelease, error) {
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	return &release, nil
}

// NeedsUpgrade compares two version strings and returns true if an upgrade is needed
func NeedsUpgrade(current, latest string) (bool, error) {
	// Normalize versions (remove 'v' prefix)
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Handle dev versions
	if strings.Contains(current, "dev") {
		// Dev versions should always allow upgrade
		return true, nil
	}

	// Parse semantic versions
	currentParts, err := ParseVersion(current)
	if err != nil {
		return false, err
	}

	latestParts, err := ParseVersion(latest)
	if err != nil {
		return false, err
	}

	// Compare major.minor.patch
	for i := 0; i < 3; i++ {
		if latestParts[i] > currentParts[i] {
			return true, nil
		}
		if latestParts[i] < currentParts[i] {
			return false, nil
		}
	}

	// Versions are equal
	return false, nil
}

// ParseVersion parses a semantic version string into [major, minor, patch]
func ParseVersion(version string) ([3]int, error) {
	var parts [3]int

	// Remove any suffixes like -dev, -beta, etc.
	version = regexp.MustCompile(`[+-].*`).ReplaceAllString(version, "")

	// Split by dots
	components := strings.Split(version, ".")
	if len(components) < 1 || len(components) > 3 {
		return parts, fmt.Errorf("invalid version format: %s", version)
	}

	// Parse each component
	for i := 0; i < len(components) && i < 3; i++ {
		num, err := strconv.Atoi(components[i])
		if err != nil {
			return parts, fmt.Errorf("invalid version number: %s", components[i])
		}
		parts[i] = num
	}

	return parts, nil
}
