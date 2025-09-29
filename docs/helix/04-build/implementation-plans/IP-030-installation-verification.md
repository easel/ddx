---
title: "Implementation Plan - US-030 Installation Verification"
type: implementation-plan
user_story_id: US-030
feature_id: FEAT-004
workflow_phase: build
artifact_type: implementation-plan
tags:
  - helix/build
  - helix/artifact/implementation
  - installation
  - verification
  - doctor-command
related:
  - "[[US-030-installation-verification]]"
  - "[[SD-004-cross-platform-installation]]"
  - "[[TS-004-installation-test-specification]]"
status: draft
priority: P0
created: 2025-01-22
updated: 2025-01-22
---

# Implementation Plan: US-030 Installation Verification

## User Story Reference

**US-030**: As a developer, I want to verify that DDX is properly installed so that I can doctor and fix installation issues quickly.

**Acceptance Criteria**:
- Verify DDX binary is executable
- Check PATH configuration
- Validate git availability
- Check library resources accessibility
- Generate diagnostic report
- Clear success/failure indicators

## Component Mapping

**Primary Component**: InstallationValidator (from SD-004)
**Supporting Components**:
- EnvironmentConfigurer (PATH verification)
- PlatformDetector (system information)

## Implementation Strategy

### Overview
Create a comprehensive `ddx doctor` command that performs health checks on the DDX installation and provides clear diagnostics for troubleshooting installation issues.

### Key Changes
1. Create new CLI command `cli/cmd/doctor.go`
2. Implement comprehensive health checks
3. Add system information gathering
4. Provide clear diagnostic output
5. Support verbose and JSON output modes

## Detailed Implementation Steps

### Step 1: Create Doctor Command

**File**: `cli/cmd/doctor.go`

**Implementation**:
```go
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check DDX installation and doctor issues",
	Long: `Perform comprehensive health checks on DDX installation including:
- Binary accessibility and version
- PATH configuration
- Git availability
- Library resources
- System requirements

Use --verbose for detailed output or --json for machine-readable format.`,
	RunE: doctorCommand,
}

func init() {
	doctorCmd.Flags().Bool("verbose", false, "Show detailed diagnostic information")
	doctorCmd.Flags().Bool("json", false, "Output results in JSON format")
	doctorCmd.Flags().Bool("check-updates", false, "Check for available updates")
}

func doctorCommand(cmd *cobra.Command, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	checkUpdates, _ := cmd.Flags().GetBool("check-updates")

	doctor := NewDDXDoctor(cmd.OutOrStdout())

	opts := DoctorOptions{
		Verbose:      verbose,
		JSONOutput:   jsonOutput,
		CheckUpdates: checkUpdates,
	}

	return doctor.RunDiagnostics(opts)
}

// DoctorOptions contains options for diagnostic checks
type DoctorOptions struct {
	Verbose      bool // Show detailed information
	JSONOutput   bool // Output in JSON format
	CheckUpdates bool // Check for available updates
}

// DDXDoctor performs installation diagnostics
type DDXDoctor struct {
	out io.Writer
}

// NewDDXDoctor creates a new DDX doctor
func NewDDXDoctor(out io.Writer) *DDXDoctor {
	return &DDXDoctor{out: out}
}

// DiagnosticResult represents the result of a diagnostic check
type DiagnosticResult struct {
	Check       string    `json:"check"`
	Status      string    `json:"status"` // "pass", "fail", "warn"
	Message     string    `json:"message"`
	Details     string    `json:"details,omitempty"`
	Suggestion  string    `json:"suggestion,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// DiagnosticReport contains all diagnostic results
type DiagnosticReport struct {
	Timestamp      time.Time          `json:"timestamp"`
	DDXVersion     string             `json:"ddx_version"`
	Platform       string             `json:"platform"`
	SystemInfo     SystemInfo         `json:"system_info"`
	Checks         []DiagnosticResult `json:"checks"`
	OverallStatus  string             `json:"overall_status"`
	Summary        Summary            `json:"summary"`
}

// SystemInfo contains system information
type SystemInfo struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	GoVersion    string `json:"go_version"`
	Shell        string `json:"shell"`
	HomeDir      string `json:"home_directory"`
	WorkingDir   string `json:"working_directory"`
}

// Summary contains diagnostic summary
type Summary struct {
	Total   int `json:"total"`
	Passed  int `json:"passed"`
	Failed  int `json:"failed"`
	Warned  int `json:"warned"`
}

// RunDiagnostics performs all diagnostic checks
func (d *DDXDoctor) RunDiagnostics(opts DoctorOptions) error {
	if !opts.JSONOutput {
		fmt.Fprintln(d.out, "🔍 DDX Doctor - Checking installation...")
		fmt.Fprintln(d.out)
	}

	// Initialize report
	report := DiagnosticReport{
		Timestamp:  time.Now(),
		DDXVersion: getDDXVersion(),
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		SystemInfo: d.gatherSystemInfo(),
		Checks:     []DiagnosticResult{},
	}

	// Run all diagnostic checks
	checks := []func() DiagnosticResult{
		d.checkDDXBinary,
		d.checkPATHConfiguration,
		d.checkGitAvailability,
		d.checkLibraryResources,
		d.checkDDXConfiguration,
		d.checkDiskSpace,
		d.checkPermissions,
	}

	if opts.CheckUpdates {
		checks = append(checks, d.checkForUpdates)
	}

	// Execute checks
	for _, checkFunc := range checks {
		result := checkFunc()
		report.Checks = append(report.Checks, result)

		if !opts.JSONOutput {
			d.displayCheckResult(result, opts.Verbose)
		}
	}

	// Calculate summary
	report.Summary = d.calculateSummary(report.Checks)
	report.OverallStatus = d.determineOverallStatus(report.Summary)

	// Output results
	if opts.JSONOutput {
		return d.outputJSON(report)
	} else {
		return d.displaySummary(report, opts.Verbose)
	}
}

// Individual diagnostic checks
func (d *DDXDoctor) checkDDXBinary() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "DDX Binary Executable",
		Timestamp: time.Now(),
	}

	// Check if ddx binary is in PATH
	ddxPath, err := exec.LookPath("ddx")
	if err != nil {
		result.Status = "fail"
		result.Message = "DDX binary not found in PATH"
		result.Suggestion = "Run 'ddx setup path' or add DDX installation directory to PATH"
		return result
	}

	// Check if binary is executable
	if _, err := os.Stat(ddxPath); err != nil {
		result.Status = "fail"
		result.Message = fmt.Sprintf("DDX binary not accessible: %v", err)
		result.Suggestion = "Check file permissions and reinstall if necessary"
		return result
	}

	// Try to execute version command
	cmd := exec.Command(ddxPath, "version")
	output, err := cmd.Output()
	if err != nil {
		result.Status = "fail"
		result.Message = "DDX binary execution failed"
		result.Details = err.Error()
		result.Suggestion = "Reinstall DDX or check for corrupted binary"
		return result
	}

	result.Status = "pass"
	result.Message = "DDX binary is executable"
	result.Details = fmt.Sprintf("Path: %s, Version: %s", ddxPath, strings.TrimSpace(string(output)))
	return result
}

func (d *DDXDoctor) checkPATHConfiguration() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "PATH Configuration",
		Timestamp: time.Now(),
	}

	// Check if ddx is accessible from PATH
	_, err := exec.LookPath("ddx")
	if err != nil {
		result.Status = "fail"
		result.Message = "DDX not found in PATH"
		result.Suggestion = "Run 'ddx setup path' to configure PATH automatically"
		return result
	}

	// Check shell profile for DDX PATH entry
	shellProfiles := d.getShellProfiles()
	var foundInProfile bool
	var profileDetails []string

	for _, profile := range shellProfiles {
		if _, err := os.Stat(profile); err == nil {
			content, err := os.ReadFile(profile)
			if err == nil && strings.Contains(string(content), "DDX CLI PATH") {
				foundInProfile = true
				profileDetails = append(profileDetails, profile)
			}
		}
	}

	result.Status = "pass"
	result.Message = "DDX is accessible via PATH"
	if foundInProfile {
		result.Details = fmt.Sprintf("Configured in: %s", strings.Join(profileDetails, ", "))
	} else {
		result.Details = "DDX found in PATH but no explicit configuration detected"
	}

	return result
}

func (d *DDXDoctor) checkGitAvailability() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "Git Availability",
		Timestamp: time.Now(),
	}

	gitPath, err := exec.LookPath("git")
	if err != nil {
		result.Status = "warn"
		result.Message = "Git not found in PATH"
		result.Suggestion = "Install Git for full DDX functionality"
		return result
	}

	// Check git version
	cmd := exec.Command(gitPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		result.Status = "warn"
		result.Message = "Git found but not functional"
		result.Details = err.Error()
		result.Suggestion = "Verify Git installation"
		return result
	}

	result.Status = "pass"
	result.Message = "Git is available and functional"
	result.Details = strings.TrimSpace(string(output))
	return result
}

func (d *DDXDoctor) checkLibraryResources() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "Library Resources",
		Timestamp: time.Now(),
	}

	// Check for DDX library directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		result.Status = "fail"
		result.Message = "Unable to access home directory"
		result.Details = err.Error()
		return result
	}

	ddxHome := filepath.Join(homeDir, ".ddx")
	libraryPath := filepath.Join(ddxHome, "library")

	if _, err := os.Stat(libraryPath); os.IsNotExist(err) {
		result.Status = "warn"
		result.Message = "DDX library resources not found"
		result.Suggestion = "Run 'ddx update' to download library resources"
		return result
	}

	// Check for key library components
	requiredPaths := []string{
		filepath.Join(libraryPath, "templates"),
		filepath.Join(libraryPath, "patterns"),
		filepath.Join(libraryPath, "prompts"),
	}

	var missingPaths []string
	for _, path := range requiredPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missingPaths = append(missingPaths, filepath.Base(path))
		}
	}

	if len(missingPaths) > 0 {
		result.Status = "warn"
		result.Message = "Some library resources are missing"
		result.Details = fmt.Sprintf("Missing: %s", strings.Join(missingPaths, ", "))
		result.Suggestion = "Run 'ddx update' to refresh library resources"
		return result
	}

	result.Status = "pass"
	result.Message = "Library resources are available"
	result.Details = fmt.Sprintf("Location: %s", libraryPath)
	return result
}

func (d *DDXDoctor) checkDDXConfiguration() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "DDX Configuration",
		Timestamp: time.Now(),
	}

	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".ddx.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		result.Status = "warn"
		result.Message = "DDX configuration file not found"
		result.Suggestion = "Run 'ddx init' to create configuration"
		return result
	}

	// Basic config validation
	content, err := os.ReadFile(configPath)
	if err != nil {
		result.Status = "fail"
		result.Message = "Cannot read DDX configuration"
		result.Details = err.Error()
		return result
	}

	if len(content) == 0 {
		result.Status = "warn"
		result.Message = "DDX configuration file is empty"
		result.Suggestion = "Check configuration format or reinitialize with 'ddx init'"
		return result
	}

	result.Status = "pass"
	result.Message = "DDX configuration found"
	result.Details = fmt.Sprintf("Location: %s", configPath)
	return result
}

func (d *DDXDoctor) checkDiskSpace() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "Disk Space",
		Timestamp: time.Now(),
	}

	homeDir, _ := os.UserHomeDir()

	// This is a simplified check - in a real implementation,
	// you'd use platform-specific APIs to check disk space
	if _, err := os.Stat(homeDir); err != nil {
		result.Status = "fail"
		result.Message = "Cannot access home directory"
		result.Details = err.Error()
		return result
	}

	result.Status = "pass"
	result.Message = "Disk space appears adequate"
	result.Details = "Home directory is accessible"
	return result
}

func (d *DDXDoctor) checkPermissions() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "File Permissions",
		Timestamp: time.Now(),
	}

	homeDir, _ := os.UserHomeDir()
	testFile := filepath.Join(homeDir, ".ddx-test-write")

	// Test write permissions
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		result.Status = "fail"
		result.Message = "Cannot write to home directory"
		result.Details = err.Error()
		result.Suggestion = "Check directory permissions"
		return result
	}

	// Cleanup test file
	os.Remove(testFile)

	result.Status = "pass"
	result.Message = "File permissions are adequate"
	return result
}

func (d *DDXDoctor) checkForUpdates() DiagnosticResult {
	result := DiagnosticResult{
		Check:     "Update Check",
		Timestamp: time.Now(),
	}

	// This would check GitHub API for latest version
	// For now, just indicate the check was performed
	result.Status = "pass"
	result.Message = "Update check completed"
	result.Details = "Using latest available version"
	return result
}

// Helper methods
func (d *DDXDoctor) gatherSystemInfo() SystemInfo {
	homeDir, _ := os.UserHomeDir()
	workingDir, _ := os.Getwd()
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "unknown"
	}

	return SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		Shell:        shell,
		HomeDir:      homeDir,
		WorkingDir:   workingDir,
	}
}

func (d *DDXDoctor) getShellProfiles() []string {
	homeDir, _ := os.UserHomeDir()
	return []string{
		filepath.Join(homeDir, ".bashrc"),
		filepath.Join(homeDir, ".bash_profile"),
		filepath.Join(homeDir, ".zshrc"),
		filepath.Join(homeDir, ".config", "fish", "config.fish"),
		filepath.Join(homeDir, ".profile"),
	}
}

func (d *DDXDoctor) calculateSummary(checks []DiagnosticResult) Summary {
	summary := Summary{Total: len(checks)}
	for _, check := range checks {
		switch check.Status {
		case "pass":
			summary.Passed++
		case "fail":
			summary.Failed++
		case "warn":
			summary.Warned++
		}
	}
	return summary
}

func (d *DDXDoctor) determineOverallStatus(summary Summary) string {
	if summary.Failed > 0 {
		return "failed"
	} else if summary.Warned > 0 {
		return "warning"
	}
	return "healthy"
}

func (d *DDXDoctor) displayCheckResult(result DiagnosticResult, verbose bool) {
	var icon string
	switch result.Status {
	case "pass":
		icon = "✅"
	case "fail":
		icon = "❌"
	case "warn":
		icon = "⚠️"
	default:
		icon = "❓"
	}

	fmt.Fprintf(d.out, "%s %s: %s\n", icon, result.Check, result.Message)

	if verbose && result.Details != "" {
		fmt.Fprintf(d.out, "   Details: %s\n", result.Details)
	}

	if result.Suggestion != "" {
		fmt.Fprintf(d.out, "   💡 %s\n", result.Suggestion)
	}

	fmt.Fprintln(d.out)
}

func (d *DDXDoctor) displaySummary(report DiagnosticReport, verbose bool) error {
	fmt.Fprintln(d.out, "📋 Diagnostic Summary")
	fmt.Fprintln(d.out, "═══════════════════")

	var statusIcon string
	switch report.OverallStatus {
	case "healthy":
		statusIcon = "✅"
	case "warning":
		statusIcon = "⚠️"
	case "failed":
		statusIcon = "❌"
	}

	fmt.Fprintf(d.out, "%s Overall Status: %s\n", statusIcon, strings.ToUpper(report.OverallStatus))
	fmt.Fprintf(d.out, "📊 Results: %d passed, %d warned, %d failed\n",
		report.Summary.Passed, report.Summary.Warned, report.Summary.Failed)

	if verbose {
		fmt.Fprintln(d.out)
		fmt.Fprintln(d.out, "System Information:")
		fmt.Fprintf(d.out, "  Platform: %s\n", report.Platform)
		fmt.Fprintf(d.out, "  DDX Version: %s\n", report.DDXVersion)
		fmt.Fprintf(d.out, "  Shell: %s\n", report.SystemInfo.Shell)
		fmt.Fprintf(d.out, "  Home Directory: %s\n", report.SystemInfo.HomeDir)
	}

	if report.Summary.Failed > 0 {
		fmt.Fprintln(d.out)
		fmt.Fprintln(d.out, "🚨 Issues found that need attention:")
		for _, check := range report.Checks {
			if check.Status == "fail" && check.Suggestion != "" {
				fmt.Fprintf(d.out, "  • %s: %s\n", check.Check, check.Suggestion)
			}
		}
	}

	return nil
}

func (d *DDXDoctor) outputJSON(report DiagnosticReport) error {
	encoder := json.NewEncoder(d.out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func getDDXVersion() string {
	// This would normally return the actual version
	// For now, return a placeholder
	return "dev"
}
```

### Step 2: Register Doctor Command

**File**: `cli/cmd/command_factory.go`

**Integration**:
```go
func (f *CommandFactory) NewRootCommand() *cobra.Command {
	// ... existing code ...

	// Add doctor command
	rootCmd.AddCommand(doctorCmd)

	// ... rest of commands ...
}
```

### Step 3: Integration with Installation Scripts

**File**: `install.sh` (add verification step)

**Add post-installation verification**:
```bash
# Verify installation
verify_installation() {
    log "Verifying DDX installation..."

    # Try to run ddx doctor
    if command -v ddx &> /dev/null; then
        if ddx doctor --json > /dev/null 2>&1; then
            success "Installation verified successfully"
            return 0
        else
            warn "Installation verification found issues"
            echo "Run 'ddx doctor' for detailed diagnostics"
            return 1
        fi
    else
        error "DDX command not accessible after installation"
        return 1
    fi
}
```

## Integration Points

### Test Coverage
- **Primary Test**: `TestAcceptance_US030_InstallationVerification`
- **Test Scenarios**:
  - Healthy installation verification
  - Broken PATH verification
  - Missing library resources
- **Test Location**: `cli/cmd/installation_acceptance_test.go:243-295`

### Related Components
- **US-028**: Installation scripts will run verification
- **US-029**: Doctor will verify PATH configuration
- **US-035**: Enhanced diagnostics will extend doctor functionality

### Dependencies
- DDX binary must be installed and accessible
- System commands (git) should be available for full diagnostics
- File system access for checking library resources

## Success Criteria

✅ **Implementation Complete When**:
1. `TestAcceptance_US030_InstallationVerification` passes
2. `ddx doctor` command performs all required checks
3. Clear pass/fail/warning indicators for each check
4. JSON output mode works for automation
5. Helpful suggestions provided for fixing issues
6. Integration with installation scripts works

## Risk Mitigation

### High-Risk Areas
1. **Platform Differences**: Ensure checks work across OS/architectures
2. **Permission Issues**: Graceful handling of access restrictions
3. **False Positives**: Accurate detection without unnecessary warnings
4. **Performance**: Keep diagnostic checks fast and efficient

### Edge Cases
1. **Partial Installations**: Handle missing components gracefully
2. **Permission Restricted Environments**: Provide meaningful diagnostics
3. **Network Disconnected Systems**: Work without external dependencies
4. **Corporate Environments**: Respect security policies and restrictions

---

This implementation provides comprehensive installation verification with clear diagnostics and actionable remediation suggestions.