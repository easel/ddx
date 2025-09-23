---
title: "Implementation Plan - US-029 Automatic PATH Configuration"
type: implementation-plan
user_story_id: US-029
feature_id: FEAT-004
workflow_phase: build
artifact_type: implementation-plan
tags:
  - helix/build
  - helix/artifact/implementation
  - installation
  - path-configuration
  - shell-integration
related:
  - "[[US-029-automatic-path-configuration]]"
  - "[[SD-004-cross-platform-installation]]"
  - "[[TS-004-installation-test-specification]]"
status: draft
priority: P0
created: 2025-01-22
updated: 2025-01-22
---

# Implementation Plan: US-029 Automatic PATH Configuration

## User Story Reference

**US-029**: As a developer, I want DDX to automatically configure my PATH so that I can use the `ddx` command immediately after installation without manual setup.

**Acceptance Criteria**:
- Shell-specific configuration (bash, zsh, fish, PowerShell)
- Automatic detection of shell profiles
- Safe modification with backup
- PATH persistence across sessions

## Component Mapping

**Primary Component**: EnvironmentConfigurer (from SD-004)
**Supporting Components**:
- InstallationManager (coordination)
- InstallationValidator (verification)

## Implementation Strategy

### Overview
Create a `ddx setup path` command and integrate automatic PATH configuration into installation scripts. Support major shells across platforms with safe profile modification and rollback capability.

### Key Changes
1. Create new CLI command `cli/cmd/setup.go`
2. Implement shell detection and profile management
3. Add backup and rollback functionality
4. Integrate with installation scripts
5. Support cross-platform PATH management

## Detailed Implementation Steps

### Step 1: Create Setup Command

**File**: `cli/cmd/setup.go`

**Implementation**:
```go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure DDX environment and shell integration",
	Long: `Configure DDX environment settings including PATH configuration,
shell completions, and other integration features.`,
}

// pathCmd configures PATH for DDX
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Configure PATH to include DDX binary",
	Long: `Automatically configure your shell's PATH to include the DDX binary
location, enabling you to use 'ddx' command from anywhere.`,
	RunE: setupPathCommand,
}

func init() {
	setupCmd.AddCommand(pathCmd)

	// Add flags
	pathCmd.Flags().String("shell", "", "Target shell (auto-detected if not specified)")
	pathCmd.Flags().Bool("force", false, "Force PATH update even if already configured")
	pathCmd.Flags().Bool("dry-run", false, "Show what would be done without making changes")
	pathCmd.Flags().Bool("backup", true, "Create backup of shell profile before modification")
}

func setupPathCommand(cmd *cobra.Command, args []string) error {
	shell, _ := cmd.Flags().GetString("shell")
	force, _ := cmd.Flags().GetBool("force")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	backup, _ := cmd.Flags().GetBool("backup")

	configurator := NewPathConfigurator(cmd.OutOrStdout())

	opts := PathConfigOptions{
		Shell:      shell,
		Force:      force,
		DryRun:     dryRun,
		CreateBackup: backup,
	}

	return configurator.ConfigurePath(opts)
}

// PathConfigOptions contains options for PATH configuration
type PathConfigOptions struct {
	Shell        string // Target shell (auto-detected if empty)
	Force        bool   // Force update even if already configured
	DryRun       bool   // Show what would be done
	CreateBackup bool   // Create backup of profile files
}

// PathConfigurator manages PATH configuration across shells
type PathConfigurator struct {
	out io.Writer
}

// NewPathConfigurator creates a new PATH configurator
func NewPathConfigurator(out io.Writer) *PathConfigurator {
	return &PathConfigurator{out: out}
}

// ConfigurePath configures the PATH for DDX
func (pc *PathConfigurator) ConfigurePath(opts PathConfigOptions) error {
	fmt.Fprintln(pc.out, "🔧 Configuring DDX PATH integration...")

	// Get DDX binary location
	ddxPath, err := pc.getDDXBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to locate DDX binary: %w", err)
	}

	ddxDir := filepath.Dir(ddxPath)
	fmt.Fprintf(pc.out, "📁 DDX binary found at: %s\n", ddxPath)

	// Detect shell if not specified
	shell := opts.Shell
	if shell == "" {
		shell = pc.detectShell()
		fmt.Fprintf(pc.out, "🐚 Detected shell: %s\n", shell)
	}

	// Get shell profile path
	profilePath, err := pc.getShellProfilePath(shell)
	if err != nil {
		return fmt.Errorf("failed to get shell profile: %w", err)
	}

	// Check if PATH already configured
	if !opts.Force {
		isConfigured, err := pc.isPathAlreadyConfigured(profilePath, ddxDir)
		if err != nil {
			return fmt.Errorf("failed to check existing PATH configuration: %w", err)
		}
		if isConfigured {
			fmt.Fprintln(pc.out, "✅ DDX is already configured in PATH")
			return nil
		}
	}

	if opts.DryRun {
		fmt.Fprintln(pc.out, "🔍 DRY RUN - would make the following changes:")
		fmt.Fprintf(pc.out, "  - Add to %s: export PATH=\"$PATH:%s\"\n", profilePath, ddxDir)
		return nil
	}

	// Create backup if requested
	if opts.CreateBackup {
		if err := pc.createBackup(profilePath); err != nil {
			fmt.Fprintf(pc.out, "⚠️  Warning: failed to create backup: %v\n", err)
		} else {
			fmt.Fprintf(pc.out, "💾 Created backup: %s.backup\n", profilePath)
		}
	}

	// Add PATH configuration
	if err := pc.addToPath(profilePath, ddxDir, shell); err != nil {
		return fmt.Errorf("failed to configure PATH: %w", err)
	}

	fmt.Fprintln(pc.out, "✅ PATH configured successfully!")
	fmt.Fprintf(pc.out, "🔄 Please restart your shell or run: source %s\n", profilePath)

	return nil
}

// Shell detection and profile management methods
func (pc *PathConfigurator) detectShell() string {
	// Try SHELL environment variable first
	if shell := os.Getenv("SHELL"); shell != "" {
		return filepath.Base(shell)
	}

	// Platform-specific detection
	switch runtime.GOOS {
	case "windows":
		return "powershell"
	default:
		// Default to bash on Unix-like systems
		return "bash"
	}
}

func (pc *PathConfigurator) getShellProfilePath(shell string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch shell {
	case "bash":
		// Try .bashrc first, then .bash_profile
		bashrc := filepath.Join(homeDir, ".bashrc")
		if _, err := os.Stat(bashrc); err == nil {
			return bashrc, nil
		}
		return filepath.Join(homeDir, ".bash_profile"), nil

	case "zsh":
		return filepath.Join(homeDir, ".zshrc"), nil

	case "fish":
		configDir := filepath.Join(homeDir, ".config", "fish")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return "", err
		}
		return filepath.Join(configDir, "config.fish"), nil

	case "powershell":
		if runtime.GOOS != "windows" {
			return "", fmt.Errorf("PowerShell profile not supported on %s", runtime.GOOS)
		}
		// PowerShell profile location varies, use Documents location
		docs := filepath.Join(homeDir, "Documents")
		return filepath.Join(docs, "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil

	default:
		// Fallback to .profile
		return filepath.Join(homeDir, ".profile"), nil
	}
}

func (pc *PathConfigurator) getDDXBinaryPath() (string, error) {
	// Try to find ddx in PATH first
	if path, err := exec.LookPath("ddx"); err == nil {
		return path, nil
	}

	// Check common installation locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	candidates := []string{
		filepath.Join(homeDir, ".local", "bin", "ddx"),
		filepath.Join(homeDir, "bin", "ddx"),
		"/usr/local/bin/ddx",
	}

	if runtime.GOOS == "windows" {
		candidates = append(candidates,
			filepath.Join(homeDir, "bin", "ddx.exe"),
			filepath.Join(homeDir, ".local", "bin", "ddx.exe"),
		)
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("ddx binary not found in common locations")
}

func (pc *PathConfigurator) isPathAlreadyConfigured(profilePath, ddxDir string) (bool, error) {
	content, err := os.ReadFile(profilePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// Check if the directory is already in PATH configuration
	return strings.Contains(string(content), ddxDir), nil
}

func (pc *PathConfigurator) createBackup(profilePath string) error {
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return nil // No file to backup
	}

	backupPath := profilePath + ".backup"
	return copyFile(profilePath, backupPath)
}

func (pc *PathConfigurator) addToPath(profilePath, ddxDir, shell string) error {
	// Ensure profile directory exists
	if err := os.MkdirAll(filepath.Dir(profilePath), 0755); err != nil {
		return err
	}

	// Generate shell-specific PATH export
	var pathExport string
	switch shell {
	case "fish":
		pathExport = fmt.Sprintf("\n# DDX CLI PATH\nset -gx PATH $PATH %s\n", ddxDir)
	case "powershell":
		pathExport = fmt.Sprintf("\n# DDX CLI PATH\n$env:PATH += \";%s\"\n", ddxDir)
	default:
		pathExport = fmt.Sprintf("\n# DDX CLI PATH\nexport PATH=\"$PATH:%s\"\n", ddxDir)
	}

	// Append to profile file
	file, err := os.OpenFile(profilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(pathExport)
	return err
}

// Utility function for file copying
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}
```

### Step 2: Register Setup Command

**File**: `cli/cmd/command_factory.go` or `cli/cmd/root.go`

**Changes**:
```go
// Add to command registration
func (f *CommandFactory) NewRootCommand() *cobra.Command {
	// ... existing code ...

	// Add setup command
	rootCmd.AddCommand(setupCmd)

	// ... rest of commands ...
}
```

### Step 3: Integration with Installation Scripts

**File**: `install.sh` (update existing)

**Add PATH configuration call**:
```bash
# After binary installation, configure PATH
configure_path() {
    log "Configuring PATH for DDX..."

    # Use the installed binary to configure itself
    if [ -x "${LOCAL_BIN}/ddx" ]; then
        "${LOCAL_BIN}/ddx" setup path --backup
        if [ $? -eq 0 ]; then
            success "PATH configured successfully"
        else
            warn "PATH configuration failed. You may need to add ${LOCAL_BIN} to your PATH manually"
            show_manual_path_instructions
        fi
    else
        warn "DDX binary not found. Falling back to manual PATH configuration"
        configure_path_manually
    fi
}

# Manual PATH configuration fallback
configure_path_manually() {
    # Detect shell and add PATH
    SHELL_NAME=$(basename "$SHELL")
    case "$SHELL_NAME" in
        bash)
            RC_FILE="$HOME/.bashrc"
            [ ! -f "$RC_FILE" ] && RC_FILE="$HOME/.bash_profile"
            ;;
        zsh)
            RC_FILE="$HOME/.zshrc"
            ;;
        fish)
            RC_FILE="$HOME/.config/fish/config.fish"
            mkdir -p "$(dirname "$RC_FILE")"
            ;;
        *)
            RC_FILE="$HOME/.profile"
            ;;
    esac

    # Check if already configured
    if [ -f "$RC_FILE" ] && grep -q "${LOCAL_BIN}" "$RC_FILE"; then
        success "PATH already configured"
        return
    fi

    # Add PATH configuration
    echo "" >> "$RC_FILE"
    echo "# DDX CLI PATH" >> "$RC_FILE"
    if [ "$SHELL_NAME" = "fish" ]; then
        echo "set -gx PATH \$PATH ${LOCAL_BIN}" >> "$RC_FILE"
    else
        echo "export PATH=\"\$PATH:${LOCAL_BIN}\"" >> "$RC_FILE"
    fi

    success "Added DDX to PATH in $RC_FILE"
}

show_manual_path_instructions() {
    echo ""
    echo "Manual PATH Configuration:"
    echo "Add the following line to your shell profile:"
    echo "  export PATH=\"\$PATH:${LOCAL_BIN}\""
    echo ""
    echo "Shell profiles:"
    echo "  bash: ~/.bashrc or ~/.bash_profile"
    echo "  zsh:  ~/.zshrc"
    echo "  fish: ~/.config/fish/config.fish"
    echo ""
}
```

**File**: `install.ps1` (update)

**Add PowerShell PATH configuration**:
```powershell
# PATH configuration function
function Set-DDXPath {
    param([string]$DDXDir)

    Write-Log "Configuring PATH for DDX..."

    # Use DDX to configure itself
    $ddxExe = Join-Path $DDXDir "ddx.exe"
    if (Test-Path $ddxExe) {
        try {
            & $ddxExe setup path --backup
            Write-Success "PATH configured successfully"
            return
        }
        catch {
            Write-Log "DDX self-configuration failed. Using manual method..."
        }
    }

    # Manual PATH configuration
    Set-DDXPathManual $DDXDir
}

function Set-DDXPathManual {
    param([string]$DDXDir)

    # Get current user PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")

    # Check if already configured
    if ($currentPath -like "*$DDXDir*") {
        Write-Success "PATH already configured"
        return
    }

    # Add to user PATH
    $newPath = "$currentPath;$DDXDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")

    Write-Success "Added DDX to user PATH"
    Write-Log "Please restart your terminal or run 'refreshenv' to use DDX"
}
```

### Step 4: Add Shell Completion Support

**Extend setup command**:
```go
// completionCmd configures shell completions
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for DDX commands.

To load completions:

Bash:
  $ source <(ddx setup completion bash)

Zsh:
  $ ddx setup completion zsh > "${fpath[1]}/_ddx"

Fish:
  $ ddx setup completion fish | source

PowerShell:
  $ ddx setup completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
		}
		return nil
	},
}

func init() {
	setupCmd.AddCommand(completionCmd)
}
```

## Integration Points

### Test Coverage
- **Primary Test**: `TestAcceptance_US029_AutomaticPathConfiguration`
- **Test Scenarios**:
  - Bash PATH configuration
  - Zsh PATH configuration
  - Fish PATH configuration
  - PowerShell PATH configuration
- **Test Location**: `cli/cmd/installation_acceptance_test.go:181-241`

### Related Components
- **US-028**: Installation scripts will call `ddx setup path`
- **US-030**: Doctor command will verify PATH configuration
- **US-033**: Uninstall will clean up PATH configuration

### Dependencies
- Shell profile files must be writable
- Home directory must be accessible
- DDX binary must be installed before PATH configuration

## Success Criteria

✅ **Implementation Complete When**:
1. `TestAcceptance_US029_AutomaticPathConfiguration` passes on all platforms
2. `ddx setup path` command works for all supported shells
3. Installation scripts automatically configure PATH
4. PATH persists across shell sessions
5. Backup and rollback functionality works
6. Manual fallback works when automatic configuration fails

## Risk Mitigation

### High-Risk Areas
1. **Shell Profile Corruption**: Always create backups, validate before writing
2. **Permission Issues**: Graceful fallback to manual instructions
3. **Shell Variations**: Comprehensive testing across shell versions
4. **Cross-platform Differences**: Platform-specific handling

### Edge Cases
1. **Missing Profile Files**: Create them safely with proper permissions
2. **Read-only File Systems**: Provide clear error messages and alternatives
3. **Non-standard Shells**: Fallback to .profile or manual instructions
4. **Corporate Environments**: Respect existing PATH configurations

---

This implementation provides robust, cross-platform PATH configuration with comprehensive error handling and fallback mechanisms.