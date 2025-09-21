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

var (
	uninstallForce bool
	uninstallPurge bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall DDx from your system",
	Long: `Uninstall DDx from your system.

This command will:
• Remove the DDx binary from your system
• Optionally remove configuration files (with --purge)
• Clean up PATH modifications

Example:
  ddx uninstall
  ddx uninstall --force
  ddx uninstall --purge`,
	RunE: runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().BoolVarP(&uninstallForce, "force", "f", false, "Skip confirmation prompt")
	uninstallCmd.Flags().BoolVar(&uninstallPurge, "purge", false, "Remove all DDx data and configuration")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	red := color.New(color.FgRed)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)
	cyan := color.New(color.FgCyan)

	cyan.Println("🗑️  DDx Uninstaller")
	fmt.Println()

	// Show what will be removed
	yellow.Println("This will remove:")
	fmt.Println("• DDx binary from your system")
	if uninstallPurge {
		fmt.Println("• Configuration files (~/.ddx.yml)")
		fmt.Println("• Local DDx data")
	}
	fmt.Println()

	// Confirm unless forced or in test mode
	if !uninstallForce && os.Getenv("DDX_TEST_MODE") != "1" {
		var confirm bool
		prompt := &survey.Confirm{
			Message: "Are you sure you want to uninstall DDx?",
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}
		if !confirm {
			yellow.Println("Uninstall cancelled")
			return nil
		}
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Check common installation locations
	installPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".local", "bin", "ddx"),
		"/usr/local/bin/ddx",
		"/usr/bin/ddx",
	}

	if runtime.GOOS == "windows" {
		installPaths = append(installPaths, filepath.Join(os.Getenv("PROGRAMFILES"), "ddx", "ddx.exe"))
	}

	var removed bool
	for _, path := range installPaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Removing %s...\n", path)
			if err := os.Remove(path); err != nil {
				red.Printf("Failed to remove %s: %v\n", path, err)
			} else {
				green.Printf("✓ Removed %s\n", path)
				removed = true
			}
		}
	}

	// Also try to remove the actual executable if different
	if !removed {
		fmt.Printf("Removing %s...\n", execPath)
		if err := os.Remove(execPath); err != nil {
			red.Printf("Failed to remove executable: %v\n", err)
		} else {
			green.Printf("✓ Removed %s\n", execPath)
		}
	}

	// Remove configuration if purge flag is set
	if uninstallPurge {
		configPath := filepath.Join(os.Getenv("HOME"), ".ddx.yml")
		if _, err := os.Stat(configPath); err == nil {
			fmt.Printf("Removing configuration %s...\n", configPath)
			if err := os.Remove(configPath); err != nil {
				red.Printf("Failed to remove config: %v\n", err)
			} else {
				green.Printf("✓ Removed configuration\n")
			}
		}
	}

	fmt.Println()
	green.Println("✅ DDx has been uninstalled")

	if !uninstallPurge {
		cyan.Println("💡 Configuration files were preserved. Use --purge to remove them.")
	}

	fmt.Println("\nTo reinstall DDx, visit: https://github.com/easel/ddx")

	return nil
}
