package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	githubAPIURL     = "https://api.github.com/repos/easel/ddx/releases/latest"
	installScriptURL = "https://raw.githubusercontent.com/easel/ddx/main/install.sh"
)

// GitHubRelease represents the GitHub API response for a release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
}

func (f *CommandFactory) runUpgrade(cmd *cobra.Command, args []string) error {
	checkOnly, _ := cmd.Flags().GetBool("check")
	force, _ := cmd.Flags().GetBool("force")

	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)

	out := cmd.OutOrStdout()

	cyan.Fprintln(out, "🔍 Checking for DDx updates...")
	fmt.Fprintln(out)

	// Get current version
	currentVersion := f.Version
	if currentVersion == "" || currentVersion == "dev" {
		currentVersion = "v0.0.1-dev"
	}

	// Fetch latest release from GitHub
	latestRelease, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := latestRelease.TagName

	// Display current and latest versions
	fmt.Fprintf(out, "Current version: %s\n", currentVersion)
	fmt.Fprintf(out, "Latest version:  %s\n", latestVersion)
	fmt.Fprintln(out)

	// Compare versions
	needsUpgrade, err := needsUpgrade(currentVersion, latestVersion)
	if err != nil && !force {
		return fmt.Errorf("failed to compare versions: %w", err)
	}

	if !needsUpgrade && !force {
		green.Fprintln(out, "✅ You are already running the latest version of DDx!")
		return nil
	}

	if checkOnly {
		if needsUpgrade {
			yellow.Fprintln(out, "⬆️  A new version of DDx is available!")
			fmt.Fprintln(out)
			fmt.Fprintln(out, "To upgrade, run:")
			green.Fprintln(out, "  ddx upgrade")
		}
		return nil
	}

	// Perform upgrade
	if force {
		yellow.Fprintf(out, "⚠️  Force upgrading to %s...\n", latestVersion)
	} else {
		cyan.Fprintf(out, "⬆️  Upgrading DDx from %s to %s...\n", currentVersion, latestVersion)
	}
	fmt.Fprintln(out)

	// Download and execute install script
	if err := executeUpgrade(out); err != nil {
		return fmt.Errorf("upgrade failed: %w", err)
	}

	fmt.Fprintln(out)
	green.Fprintln(out, "✅ DDx has been upgraded successfully!")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Run 'ddx version' to verify the new version.")

	return nil
}

// fetchLatestRelease fetches the latest release information from GitHub
func fetchLatestRelease() (*GitHubRelease, error) {
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

// needsUpgrade compares two version strings and returns true if an upgrade is needed
func needsUpgrade(current, latest string) (bool, error) {
	// Normalize versions (remove 'v' prefix)
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Handle dev versions
	if strings.Contains(current, "dev") {
		// Dev versions should always allow upgrade
		return true, nil
	}

	// Parse semantic versions
	currentParts, err := parseVersion(current)
	if err != nil {
		return false, err
	}

	latestParts, err := parseVersion(latest)
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

// parseVersion parses a semantic version string into [major, minor, patch]
func parseVersion(version string) ([3]int, error) {
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

// executeUpgrade downloads and executes the install script
func executeUpgrade(out io.Writer) error {
	// Download install script
	resp, err := http.Get(installScriptURL)
	if err != nil {
		return fmt.Errorf("failed to download install script: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download install script: status %d", resp.StatusCode)
	}

	scriptContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read install script: %w", err)
	}

	// Write to temporary file
	tmpFile, err := os.CreateTemp("", "ddx-install-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(scriptContent); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write install script: %w", err)
	}

	if err := tmpFile.Chmod(0755); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to make script executable: %w", err)
	}
	tmpFile.Close()

	// Execute install script
	cmd := exec.Command("bash", tmpFile.Name())
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install script failed: %w", err)
	}

	return nil
}
