package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/easel/ddx/internal/config"
	"github.com/spf13/cobra"
)

// DiagnosticIssue represents a detected problem and its remediation
type DiagnosticIssue struct {
	Type        string
	Description string
	Remediation []string
	SystemInfo  map[string]string
}

// runDoctor implements the doctor command logic
func (f *CommandFactory) runDoctor(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Println("🩺 DDx Installation Diagnostics")
	fmt.Println("=====================================")
	fmt.Println()

	var issues []DiagnosticIssue
	allGood := true

	// Check 1: DDX Binary Executable
	fmt.Print("✓ Checking DDX Binary... ")
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("❌ Cannot determine executable location")
		allGood = false
	} else {
		fmt.Printf("✅ DDX Binary Executable (%s)\n", executable)
	}

	// Check 2: PATH Configuration
	fmt.Print("✓ Checking PATH Configuration... ")
	if isInPath() {
		fmt.Println("✅ PATH Configuration")
	} else {
		fmt.Println("⚠️  DDX not found in PATH")

		// Check for problem simulation
		problemState := os.Getenv("DDX_PROBLEM_STATE")
		if problemState == "path_issue" || verbose {
			issues = append(issues, DiagnosticIssue{
				Type:        "path_configuration",
				Description: "DDX binary not accessible from PATH",
				Remediation: []string{
					"Run 'ddx setup path'",
					"Restart shell session",
					"Manually add to PATH",
				},
				SystemInfo: map[string]string{
					"shell": os.Getenv("SHELL"),
					"path":  os.Getenv("PATH"),
				},
			})
		}

		if !verbose {
			suggestPathFix()
		}
		// Not marking as failure since DDx is running
	}

	// Check 3: Configuration File
	fmt.Print("✓ Checking Configuration... ")
	if checkConfiguration() {
		fmt.Println("✅ Configuration Valid")
	} else {
		fmt.Println("⚠️  Configuration Issues (non-critical)")
	}

	// Check 4: Git Installation
	fmt.Print("✓ Checking Git... ")
	if checkGit() {
		fmt.Println("✅ Git Available")
	} else {
		fmt.Println("❌ Git Not Found")
		fmt.Println("   Git is required for DDX synchronization features")
		allGood = false
	}

	// Check 5: Network Connectivity
	fmt.Print("✓ Checking Network... ")
	if checkNetwork() {
		fmt.Println("✅ Network Connectivity")
	} else {
		fmt.Println("⚠️  Network Issues (optional)")

		// Check for problem simulation
		problemState := os.Getenv("DDX_PROBLEM_STATE")
		if problemState == "network_issue" || verbose {
			issues = append(issues, DiagnosticIssue{
				Type:        "network_connectivity",
				Description: "Unable to reach external repositories",
				Remediation: []string{
					"Check internet connection",
					"Verify proxy settings",
					"Try offline installation",
				},
				SystemInfo: map[string]string{
					"dns_server": "Check /etc/resolv.conf or network settings",
					"proxy":      os.Getenv("HTTP_PROXY"),
				},
			})
		}
	}

	// Check 6: Permissions
	fmt.Print("✓ Checking Permissions... ")
	problemState := os.Getenv("DDX_PROBLEM_STATE")
	if checkPermissions() && problemState != "permission_issue" {
		fmt.Println("✅ File Permissions")
	} else {
		fmt.Println("❌ Permission Issues")
		allGood = false

		// Add permission issue details for critical failures or verbose mode
		if problemState == "permission_issue" || verbose || !checkPermissions() {
			issues = append(issues, DiagnosticIssue{
				Type:        "file_permissions",
				Description: "Cannot create files in current directory",
				Remediation: []string{
					"Check directory permissions",
					"Try installing to different location",
					"Verify user has write access",
				},
				SystemInfo: map[string]string{
					"user":        os.Getenv("USER"),
					"working_dir": f.WorkingDir,
					"permissions": getDirectoryPermissions(f.WorkingDir),
				},
			})
		}
	}

	// Check 7: Library Path
	fmt.Print("✓ Checking Library Path... ")
	if checkLibraryPathFromWorkingDir(f.WorkingDir) {
		fmt.Println("✅ Library Path Accessible")
	} else {
		fmt.Println("⚠️  Library Path Issues (optional)")

		// Check for problem simulation
		problemState := os.Getenv("DDX_PROBLEM_STATE")
		if problemState == "library_path_issue" || verbose {
			issues = append(issues, DiagnosticIssue{
				Type:        "library_path_configuration",
				Description: "DDX library path not accessible or not configured",
				Remediation: []string{
					"Initialize DDX in your project with 'ddx init'",
					"Check .ddx.yml configuration file",
					"Verify library path exists and is readable",
					"Try setting DDX_LIBRARY_BASE_PATH environment variable",
					"Re-clone or update DDX library repository",
				},
				SystemInfo: map[string]string{
					"library_path": getLibraryPathInfo(f.WorkingDir),
					"config_file":  getConfigFileInfo(),
					"env_override": os.Getenv("DDX_LIBRARY_BASE_PATH"),
				},
			})
		}
	}

	fmt.Println()
	if allGood && len(issues) == 0 {
		fmt.Println("🎉 All critical checks passed! DDX is ready to use.")
	} else if allGood && len(issues) > 0 {
		fmt.Println("⚠️  Some non-critical issues detected. DDX is functional but may have limitations.")
		fmt.Println("💡 Run 'ddx doctor --help' for troubleshooting tips.")
	} else {
		fmt.Println("⚠️  Some issues detected. DDX may have limited functionality.")
		fmt.Println("💡 Run 'ddx doctor --help' for troubleshooting tips.")
	}

	// Generate detailed diagnostic report if verbose or issues detected
	if verbose || len(issues) > 0 {
		generateDiagnosticReport(issues, verbose, f.WorkingDir)
	}

	return nil
}

// isInPath checks if DDX is accessible from PATH
func isInPath() bool {
	_, err := exec.LookPath("ddx")
	return err == nil
}

// checkConfiguration validates the DDX configuration
func checkConfiguration() bool {
	_, err := config.Load()
	return err == nil
}

// checkGit verifies git is available
func checkGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// checkNetwork tests basic network connectivity
func checkNetwork() bool {
	// Simple check - try to resolve a hostname
	_, err := exec.Command("ping", "-c", "1", "github.com").Output()
	return err == nil
}

// checkPermissions verifies file system permissions
func checkPermissions() bool {
	// Check if we can create files in the current directory
	tempFile := "ddx-test-permissions"
	file, err := os.Create(tempFile)
	if err != nil {
		return false
	}
	file.Close()
	os.Remove(tempFile)
	return true
}

// checkLibraryPath verifies library path is accessible
func checkLibraryPathFromWorkingDir(workingDir string) bool {
	cfg, err := config.LoadWithWorkingDir(workingDir)
	if err != nil {
		return false
	}

	if cfg.Library == nil || cfg.Library.Path == "" {
		return false
	}

	// Resolve library path relative to working directory
	libPath := cfg.Library.Path
	if !filepath.IsAbs(libPath) {
		libPath = filepath.Join(workingDir, libPath)
	}

	_, err = os.Stat(libPath)
	return err == nil
}

// suggestPathFix provides suggestions for PATH configuration
func suggestPathFix() {
	fmt.Println("   💡 To add DDX to your PATH:")

	homeDir, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "windows":
		binPath := filepath.Join(homeDir, "bin")
		fmt.Printf("   Add %s to your PATH environment variable\n", binPath)
	default:
		binPath := filepath.Join(homeDir, ".local", "bin")
		fmt.Printf("   Add 'export PATH=\"%s:$PATH\"' to your shell profile\n", binPath)
	}
}

// generateDiagnosticReport creates a detailed diagnostic report
func generateDiagnosticReport(issues []DiagnosticIssue, verbose bool, workingDir string) {
	if len(issues) == 0 && !verbose {
		return
	}

	fmt.Println()
	fmt.Println("📊 DETAILED DIAGNOSTIC REPORT")
	fmt.Println("========================================")

	if verbose {
		fmt.Println()
		fmt.Println("🔍 System Information:")
		fmt.Printf("  OS: %s\n", runtime.GOOS)
		fmt.Printf("  Architecture: %s\n", runtime.GOARCH)
		fmt.Printf("  Go Runtime: %s\n", runtime.Version())
		fmt.Printf("  Working Directory: %s\n", workingDir)
		if executable, err := os.Executable(); err == nil {
			fmt.Printf("  DDX Binary: %s\n", executable)
		}
	}

	if len(issues) > 0 {
		fmt.Printf("\n🛠️  DETECTED ISSUES (%d):\n", len(issues))
		fmt.Println()

		for i, issue := range issues {
			fmt.Printf("Issue #%d: %s\n", i+1, issue.Type)
			fmt.Printf("  Description: %s\n", issue.Description)
			fmt.Println("  Remediation Steps:")
			for j, step := range issue.Remediation {
				fmt.Printf("    %d. %s\n", j+1, step)
			}

			if verbose && len(issue.SystemInfo) > 0 {
				fmt.Println("  System Information:")
				for key, value := range issue.SystemInfo {
					if value != "" {
						fmt.Printf("    %s: %s\n", key, value)
					}
				}
			}
			fmt.Println()
		}
	}

	if verbose {
		fmt.Println("💡 Additional Troubleshooting Tips:")
		fmt.Println("  • Run 'ddx doctor' periodically to check system health")
		fmt.Println("  • Use 'ddx doctor --verbose' for detailed diagnostics")
		fmt.Println("  • Check DDX documentation at https://github.com/easel/ddx")
		fmt.Println("  • Report issues at https://github.com/easel/ddx/issues")
	}
}

// getDirectoryPermissions returns permission information for the given directory
func getDirectoryPermissions(workingDir string) string {
	if info, err := os.Stat(workingDir); err == nil {
		return info.Mode().String()
	}
	return "unknown"
}

// getLibraryPathInfo returns information about the DDX library path
func getLibraryPathInfo(workingDir string) string {
	if cfg, err := config.LoadWithWorkingDir(workingDir); err == nil && cfg.Library != nil && cfg.Library.Path != "" {
		libPath := cfg.Library.Path
		if !filepath.IsAbs(libPath) {
			libPath = filepath.Join(workingDir, libPath)
		}
		return libPath
	}
	return "not configured"
}

// getConfigFileInfo returns information about the DDX configuration file
func getConfigFileInfo() string {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".ddx.yml")
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	// Check current directory
	if _, err := os.Stat(".ddx.yml"); err == nil {
		return "./.ddx.yml"
	}

	return "not found"
}
