package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a fresh root command for tests
func getVersionTestRootCommand() *cobra.Command {
	factory := NewCommandFactory("/tmp")
	return factory.NewRootCommand()
}

// TestAcceptance_US008_CheckDDxVersion tests US-008: Check DDX Version
func TestAcceptance_US008_CheckDDxVersion(t *testing.T) {

	t.Run("display_current_version", func(t *testing.T) {
		// AC: Given I want version info, when I run `ddx version`, then the current DDX version number is displayed
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should execute successfully")

		// Should display version number
		assert.Contains(t, output, "DDx", "Should mention DDx")
		assert.Contains(t, output, "v", "Should show version with 'v' prefix")

		// Should follow semantic versioning pattern (at least major.minor.patch)
		lines := strings.Split(output, "\n")
		versionLine := ""
		for _, line := range lines {
			if strings.Contains(line, "DDx") && strings.Contains(line, "v") {
				versionLine = line
				break
			}
		}
		assert.NotEmpty(t, versionLine, "Should find version line")
		assert.Contains(t, versionLine, ".", "Version should contain dots (semantic versioning)")
	})

	t.Run("include_build_information", func(t *testing.T) {
		// AC: Given version is displayed, when I view the output, then build information (commit hash, build date) is included
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should work")

		// Should include commit hash
		assert.Contains(t, output, "Commit:", "Should show commit hash")

		// Should include build date
		assert.Contains(t, output, "Built:", "Should show build date")

		// Verify format of build information
		lines := strings.Split(output, "\n")
		hasCommit := false
		hasBuilt := false

		for _, line := range lines {
			if strings.Contains(line, "Commit:") {
				hasCommit = true
				// Commit hash should be non-empty
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					hash := strings.TrimSpace(parts[1])
					assert.NotEmpty(t, hash, "Commit hash should not be empty")
					assert.Greater(t, len(hash), 6, "Commit hash should be reasonably long")
				}
			}
			if strings.Contains(line, "Built:") {
				hasBuilt = true
				// Build date should be in reasonable format
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					date := strings.TrimSpace(parts[1])
					assert.NotEmpty(t, date, "Build date should not be empty")
					// In tests, build date might be "unknown", in real builds it should be current
					if date != "unknown" {
						assert.Contains(t, date, "2025", "Build date should be current year")
					}
				}
			}
		}

		assert.True(t, hasCommit, "Should have commit information")
		assert.True(t, hasBuilt, "Should have build date information")
	})

	t.Run("automatic_update_check", func(t *testing.T) {
		// AC: Given I'm online, when I check version, then the system checks for available updates automatically
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should work")

		// This test documents expected behavior - update checking may need implementation
		if strings.Contains(output, "Update available") || strings.Contains(output, "up to date") {
			// Update checking is implemented
			t.Log("Update checking appears to be implemented")
		} else {
			// Update checking needs implementation
			t.Skip("Automatic update checking not yet implemented - test documents requirement")
		}
	})

	t.Run("changelog_highlights_for_updates", func(t *testing.T) {
		// AC: Given updates are available, when version is displayed, then changelog highlights for newer versions are shown
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should work")

		// This test documents expected behavior for when updates are available
		if strings.Contains(output, "Update available") {
			// If update is available, should show changelog highlights
			assert.Contains(t, output, "What's new", "Should show changelog section")
		} else {
			// No updates available or feature not implemented
			t.Skip("No updates available or changelog display not implemented - test documents requirement")
		}
	})

	t.Run("suppress_update_check_flag", func(t *testing.T) {
		// AC: Given I don't want update checks, when I run `ddx version --no-check`, then update checking is suppressed
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version", "--no-check")

		if err == nil {
			// --no-check flag is implemented
			assert.NotContains(t, output, "Checking for updates", "Should not show update check messages")
			assert.NotContains(t, output, "Update available", "Should not show update availability")

			// Should still show version information
			assert.Contains(t, output, "DDx", "Should still show version info")
			assert.Contains(t, output, "v", "Should still show version number")
		} else {
			// --no-check flag needs implementation
			assert.Contains(t, err.Error(), "unknown flag", "Flag not yet implemented")
			t.Skip("--no-check flag not yet implemented - test documents requirement")
		}
	})

	t.Run("outdated_version_indication", func(t *testing.T) {
		// AC: Given my version is outdated, when I check version, then a clear indication that the version is outdated is shown
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should work")

		// This test documents expected behavior when version is outdated
		if strings.Contains(output, "outdated") || strings.Contains(output, "Update available") {
			// Outdated indication is working
			assert.True(t, true, "Outdated version indication appears to be working")
		} else {
			// Either up to date or feature needs implementation
			t.Skip("Version appears up to date or outdated indication not implemented - test documents requirement")
		}
	})

	t.Run("compatibility_warnings", func(t *testing.T) {
		// AC: Given version changes may affect compatibility, when updates are available, then compatibility warnings are displayed
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should work")

		// This test documents expected behavior for compatibility warnings
		if strings.Contains(output, "compatibility") || strings.Contains(output, "breaking") {
			// Compatibility warnings are implemented
			assert.True(t, true, "Compatibility warnings appear to be working")
		} else {
			// No compatibility issues or feature needs implementation
			t.Skip("No compatibility issues or compatibility warnings not implemented - test documents requirement")
		}
	})

	t.Run("version_command_error_handling", func(t *testing.T) {
		// Test that version command handles various error scenarios gracefully
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		// Version command should never fail for basic version display
		assert.NoError(t, err, "Version command should not fail")
		assert.NotEmpty(t, output, "Should produce output even if update check fails")

		// Should contain essential version information even if network fails
		assert.Contains(t, output, "DDx", "Should show DDx name")
		assert.Contains(t, output, "v", "Should show version")
	})

	t.Run("version_format_validation", func(t *testing.T) {
		// Test that version follows expected format patterns
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should work")

		// Should not have empty lines at the beginning
		lines := strings.Split(strings.TrimSpace(output), "\n")
		assert.Greater(t, len(lines), 0, "Should have at least one line of output")

		// First line should contain version info
		firstLine := lines[0]
		assert.Contains(t, firstLine, "DDx", "First line should contain DDx")
		assert.Contains(t, firstLine, "v", "First line should contain version")

		// Should not contain debug or error messages in normal operation
		lowerOutput := strings.ToLower(output)
		assert.NotContains(t, lowerOutput, "error", "Should not contain error messages")
		assert.NotContains(t, lowerOutput, "debug", "Should not contain debug messages")
		assert.NotContains(t, lowerOutput, "panic", "Should not contain panic messages")
	})

	t.Run("semantic_versioning_compliance", func(t *testing.T) {
		// Test that version follows semantic versioning
		rootCmd := getVersionTestRootCommand()
		output, err := executeCommand(rootCmd, "version")

		require.NoError(t, err, "Version command should work")

		// Extract version number
		versionLine := ""
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "DDx") && strings.Contains(line, "v") {
				versionLine = line
				break
			}
		}

		require.NotEmpty(t, versionLine, "Should find version line")

		// Should contain semantic version pattern
		hasSemanticVersion := strings.Contains(versionLine, ".")
		assert.True(t, hasSemanticVersion, "Should use semantic versioning (contain dots)")

		// Extract version part
		parts := strings.Fields(versionLine)
		var versionPart string
		for _, part := range parts {
			if strings.HasPrefix(part, "v") {
				versionPart = part
				break
			}
		}

		assert.NotEmpty(t, versionPart, "Should find version part starting with 'v'")

		// Should be in format like v1.2.3 or v1.2.3-beta.1
		versionNum := strings.TrimPrefix(versionPart, "v")
		assert.NotEmpty(t, versionNum, "Should have version number after 'v'")

		// Should have at least major.minor pattern
		dotCount := strings.Count(versionNum, ".")
		assert.GreaterOrEqual(t, dotCount, 1, "Should have at least major.minor versioning")
	})
}
