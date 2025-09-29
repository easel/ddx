package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/easel/ddx/internal/update"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	installScriptURL = "https://raw.githubusercontent.com/easel/ddx/main/install.sh"
)

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
	latestRelease, err := update.FetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	latestVersion := latestRelease.TagName

	// Display current and latest versions
	fmt.Fprintf(out, "Current version: %s\n", currentVersion)
	fmt.Fprintf(out, "Latest version:  %s\n", latestVersion)
	fmt.Fprintln(out)

	// Compare versions
	needsUpgrade, err := update.NeedsUpgrade(currentVersion, latestVersion)
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
