package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Global variables have been removed - use local variables in runUninstall

// Command registration is now handled by command_factory.go
// This file only contains the run function implementation

// runUninstall implements the uninstall command logic
func runUninstall(cmd *cobra.Command, args []string) error {
	// Get flag values locally
	uninstallForce, _ := cmd.Flags().GetBool("force")
	uninstallPurge, _ := cmd.Flags().GetBool("purge")
	keepConfig, _ := cmd.Flags().GetBool("keep-config")
	keepProjects, _ := cmd.Flags().GetBool("keep-projects")

	red := color.New(color.FgRed, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	// Confirm uninstallation
	if !uninstallForce {
		confirm := false
		prompt := &survey.Confirm{
			Message: "Are you sure you want to uninstall DDx?",
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}
		if !confirm {
			_, _ = yellow.Fprintln(cmd.OutOrStdout(), "Uninstallation cancelled")
			return nil
		}
	}

	_, _ = red.Fprintln(cmd.OutOrStdout(), "🗑️  Uninstalling DDx...")
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	// Remove binary
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine binary path: %w", err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Removing binary: %s\n", binaryPath)
	if err := os.Remove(binaryPath); err != nil && !os.IsNotExist(err) {
		_, _ = yellow.Fprintf(cmd.OutOrStdout(), "⚠️  Failed to remove binary: %v\n", err)
	}

	// Remove configuration files
	if !keepConfig || uninstallPurge {
		home, err := os.UserHomeDir()
		if err == nil {
			configPath := filepath.Join(home, ".ddx.yml")
			if _, err := os.Stat(configPath); err == nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Removing config: %s\n", configPath)
				_ = os.Remove(configPath)
			}

			ddxDir := filepath.Join(home, ".ddx")
			if _, err := os.Stat(ddxDir); err == nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Removing directory: %s\n", ddxDir)
				_ = os.RemoveAll(ddxDir)
			}
		}
	}

	// Remove shell completions
	removeCompletions(cmd)

	// Clean up project directories if requested
	if uninstallPurge && !keepProjects {
		_, _ = yellow.Fprintln(cmd.OutOrStdout(), "⚠️  Purge mode: removing .ddx directories from projects")
		// Note: This would need to scan for projects, which is risky
		// For safety, we'll just inform the user
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Please manually remove .ddx directories from your projects if needed")
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout())
	_, _ = green.Fprintln(cmd.OutOrStdout(), "✅ DDx has been uninstalled")

	if keepConfig {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Configuration files were preserved")
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout())
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Thank you for using DDx!")
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "You can reinstall anytime from: https://github.com/ddx-tools/ddx")

	return nil
}

func removeCompletions(cmd *cobra.Command) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	completionFiles := map[string][]string{
		"bash": {
			filepath.Join(home, ".local/share/bash-completion/completions/ddx"),
			"/usr/local/share/bash-completion/completions/ddx",
		},
		"zsh": {
			filepath.Join(home, ".zsh/completions/_ddx"),
			"/usr/local/share/zsh/site-functions/_ddx",
		},
		"fish": {
			filepath.Join(home, ".config/fish/completions/ddx.fish"),
		},
	}

	// Detect shell
	shell := os.Getenv("SHELL")
	if shell == "" && runtime.GOOS == "windows" {
		shell = "powershell"
	}

	for shellName, paths := range completionFiles {
		if shell != "" && !contains(shell, shellName) {
			continue
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Removing %s completion: %s\n", shellName, path)
				_ = os.Remove(path)
			}
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
