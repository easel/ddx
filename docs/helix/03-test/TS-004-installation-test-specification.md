---
title: "Test Specification - Cross-Platform Installation"
type: test-specification
feature_id: FEAT-004
workflow_phase: test
artifact_type: test-specification
tags:
  - helix/test
  - helix/artifact/test
  - helix/phase/test
  - installation
  - cross-platform
  - package-manager
related:
  - "[[FEAT-004-cross-platform-installation]]"
  - "[[SD-004-cross-platform-installation]]"
  - "[[US-028-one-command-installation]]"
  - "[[US-029-automatic-path-configuration]]"
  - "[[US-030-installation-verification]]"
  - "[[US-031-package-manager-installation]]"
  - "[[US-032-upgrade-existing-installation]]"
  - "[[US-033-uninstall-ddx]]"
  - "[[US-034-offline-installation]]"
  - "[[US-035-installation-diagnostics]]"
status: draft
priority: P0
created: 2025-01-22
updated: 2025-01-22
---

# Test Specification: FEAT-004 Cross-Platform Installation

## Test Strategy Overview

This test specification defines comprehensive test scenarios for the cross-platform installation system. Following Test-Driven Development principles, these tests must be written and failing BEFORE any implementation begins.

### Test Categories

1. **Acceptance Tests** - User story validation (Red phase requirement)
2. **Integration Tests** - Installation workflow testing
3. **Contract Tests** - CLI command interface testing
4. **Platform Tests** - Cross-platform compatibility
5. **Performance Tests** - Installation speed and reliability

## Acceptance Test Specifications

### 1. US-028: One-Command Installation

#### Test Suite: `TestAcceptance_US028_OneCommandInstallation`

```go
// Path: cli/cmd/installation_acceptance_test.go

func TestAcceptance_US028_OneCommandInstallation(t *testing.T) {
    tests := []struct {
        name         string
        platform     string
        architecture string
        command      string
        expected     InstallationResult
    }{
        {
            name:         "unix_one_command_install",
            platform:     "linux",
            architecture: "amd64",
            command:      "curl -sSL https://ddx.dev/install | sh",
            expected:     InstallationResult{Success: true, BinaryPath: "~/.local/bin/ddx"},
        },
        {
            name:         "macos_one_command_install",
            platform:     "darwin",
            architecture: "arm64",
            command:      "curl -sSL https://ddx.dev/install | sh",
            expected:     InstallationResult{Success: true, BinaryPath: "~/.local/bin/ddx"},
        },
        {
            name:         "windows_one_command_install",
            platform:     "windows",
            architecture: "amd64",
            command:      "iwr -useb https://ddx.dev/install.ps1 | iex",
            expected:     InstallationResult{Success: true, BinaryPath: "%USERPROFILE%\\bin\\ddx.exe"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: I am on a target platform
            env := setupTestEnvironment(tt.platform, tt.architecture)
            defer env.Cleanup()

            // When: I execute the one-command installation
            result := env.ExecuteInstallCommand(tt.command)

            // Then: DDX is installed and ready to use
            assert.Equal(t, tt.expected.Success, result.Success)
            assert.FileExists(t, result.BinaryPath)
            assert.True(t, result.CompletedInUnder60Seconds)

            // And: DDX version command works
            version := env.RunCommand("ddx version")
            assert.Contains(t, version.Output, "ddx version")
            assert.Equal(t, 0, version.ExitCode)
        })
    }
}
```

**Test Data Requirements:**
- Mock GitHub releases with platform-specific binaries
- Test network connectivity scenarios
- Various platform/architecture combinations

**Validation Criteria:**
- Installation completes in <60 seconds
- No admin privileges required
- Binary is executable and version command works
- EXIT_SUCCESS (0) for successful installation

### 2. US-029: Automatic PATH Configuration

#### Test Suite: `TestAcceptance_US029_AutomaticPathConfiguration`

```go
func TestAcceptance_US029_AutomaticPathConfiguration(t *testing.T) {
    tests := []struct {
        name          string
        shell         string
        profileFile   string
        expectedPath  string
    }{
        {
            name:         "bash_path_configuration",
            shell:        "bash",
            profileFile:  "~/.bashrc",
            expectedPath: "~/.local/bin",
        },
        {
            name:         "zsh_path_configuration",
            shell:        "zsh",
            profileFile:  "~/.zshrc",
            expectedPath: "~/.local/bin",
        },
        {
            name:         "fish_path_configuration",
            shell:        "fish",
            profileFile:  "~/.config/fish/config.fish",
            expectedPath: "~/.local/bin",
        },
        {
            name:         "powershell_path_configuration",
            shell:        "powershell",
            profileFile:  "$PROFILE",
            expectedPath: "%USERPROFILE%\\bin",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: I have a shell environment
            env := setupShellEnvironment(tt.shell)
            defer env.Cleanup()

            // When: Installation completes
            result := env.RunInstallation()
            assert.True(t, result.Success)

            // Then: PATH is automatically configured
            pathContent := env.ReadFile(tt.profileFile)
            assert.Contains(t, pathContent, tt.expectedPath)

            // And: DDX is accessible in new shell sessions
            newShell := env.NewShellSession()
            ddxVersion := newShell.RunCommand("ddx version")
            assert.Equal(t, 0, ddxVersion.ExitCode)
        })
    }
}
```

**Test Data Requirements:**
- Clean shell environments for each shell type
- Original profile file backups
- PATH modification detection

**Validation Criteria:**
- Shell profile files are modified correctly
- PATH contains DDX installation directory
- New shell sessions can execute DDX commands
- Original configurations are preserved

### 3. US-030: Installation Verification

#### Test Suite: `TestAcceptance_US030_InstallationVerification`

```go
func TestAcceptance_US030_InstallationVerification(t *testing.T) {
    tests := []struct {
        name              string
        installationState InstallationState
        expectedChecks    []VerificationCheck
    }{
        {
            name:              "healthy_installation_verification",
            installationState: HealthyInstallation(),
            expectedChecks: []VerificationCheck{
                {Name: "DDX Binary Executable", Status: "✅ PASS"},
                {Name: "PATH Configuration", Status: "✅ PASS"},
                {Name: "Git Availability", Status: "✅ PASS"},
                {Name: "Library Resources", Status: "✅ PASS"},
            },
        },
        {
            name:              "broken_path_verification",
            installationState: BrokenPathInstallation(),
            expectedChecks: []VerificationCheck{
                {Name: "DDX Binary Executable", Status: "✅ PASS"},
                {Name: "PATH Configuration", Status: "❌ FAIL"},
                {Name: "Git Availability", Status: "✅ PASS"},
                {Name: "Library Resources", Status: "✅ PASS"},
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: DDX is installed with specific state
            env := setupInstallationState(tt.installationState)
            defer env.Cleanup()

            // When: I run DDX doctor command
            result := env.RunCommand("ddx doctor")

            // Then: Verification checks are performed
            for _, check := range tt.expectedChecks {
                assert.Contains(t, result.Output, check.Name)
                assert.Contains(t, result.Output, check.Status)
            }

            // And: Exit code reflects overall health
            if hasFailures(tt.expectedChecks) {
                assert.Equal(t, 1, result.ExitCode)
            } else {
                assert.Equal(t, 0, result.ExitCode)
            }
        })
    }
}
```

### 4. US-031: Package Manager Installation

#### Test Suite: `TestAcceptance_US031_PackageManagerInstallation`

```go
func TestAcceptance_US031_PackageManagerInstallation(t *testing.T) {
    tests := []struct {
        name           string
        packageManager string
        platform       string
        installCommand string
    }{
        {
            name:           "homebrew_installation",
            packageManager: "homebrew",
            platform:       "darwin",
            installCommand: "brew install ddx",
        },
        {
            name:           "apt_installation",
            packageManager: "apt",
            platform:       "ubuntu",
            installCommand: "sudo apt install ddx",
        },
        {
            name:           "chocolatey_installation",
            packageManager: "chocolatey",
            platform:       "windows",
            installCommand: "choco install ddx",
        },
        {
            name:           "scoop_installation",
            packageManager: "scoop",
            platform:       "windows",
            installCommand: "scoop install ddx",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: Platform has package manager available
            env := setupPackageManagerEnvironment(tt.packageManager, tt.platform)
            defer env.Cleanup()

            // When: I install via package manager
            result := env.RunCommand(tt.installCommand)

            // Then: Installation succeeds
            assert.Equal(t, 0, result.ExitCode)

            // And: DDX is available in PATH
            version := env.RunCommand("ddx version")
            assert.Equal(t, 0, version.ExitCode)
            assert.Contains(t, version.Output, "ddx version")
        })
    }
}
```

### 5. US-032: Upgrade Existing Installation

#### Test Suite: `TestAcceptance_US032_UpgradeExistingInstallation`

```go
func TestAcceptance_US032_UpgradeExistingInstallation(t *testing.T) {
    tests := []struct {
        name            string
        currentVersion  string
        targetVersion   string
        upgradeMethod   string
    }{
        {
            name:           "self_update_to_latest",
            currentVersion: "v1.0.0",
            targetVersion:  "v1.1.0",
            upgradeMethod:  "ddx self-update",
        },
        {
            name:           "self_update_to_specific_version",
            currentVersion: "v1.0.0",
            targetVersion:  "v1.0.5",
            upgradeMethod:  "ddx self-update --version v1.0.5",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: DDX is installed with current version
            env := setupDDXVersion(tt.currentVersion)
            defer env.Cleanup()

            // When: I run self-update command
            result := env.RunCommand(tt.upgradeMethod)

            // Then: Upgrade succeeds
            assert.Equal(t, 0, result.ExitCode)
            assert.Contains(t, result.Output, "Upgrade completed")

            // And: Version is updated
            version := env.RunCommand("ddx version")
            assert.Contains(t, version.Output, tt.targetVersion)

            // And: User configurations are preserved
            config := env.ReadFile("~/.ddx.yml")
            assert.NotEmpty(t, config) // Existing config should remain
        })
    }
}
```

### 6. US-033: Uninstall DDX

#### Test Suite: `TestAcceptance_US033_UninstallDDX`

```go
func TestAcceptance_US033_UninstallDDX(t *testing.T) {
    tests := []struct {
        name              string
        uninstallOptions  []string
        preserveUserData  bool
    }{
        {
            name:             "uninstall_preserve_data",
            uninstallOptions: []string{"--preserve-data"},
            preserveUserData: true,
        },
        {
            name:             "uninstall_remove_all",
            uninstallOptions: []string{"--remove-all"},
            preserveUserData: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: DDX is fully installed
            env := setupFullDDXInstallation()
            defer env.Cleanup()

            // When: I run uninstall command
            command := append([]string{"ddx uninstall"}, tt.uninstallOptions...)
            result := env.RunCommand(strings.Join(command, " "))

            // Then: Uninstallation succeeds
            assert.Equal(t, 0, result.ExitCode)

            // And: Binary is removed
            assert.False(t, env.FileExists("~/.local/bin/ddx"))

            // And: PATH configuration is cleaned
            bashrc := env.ReadFile("~/.bashrc")
            assert.NotContains(t, bashrc, "DDX CLI PATH")

            // And: User data handling follows options
            if tt.preserveUserData {
                assert.True(t, env.FileExists("~/.ddx.yml"))
            } else {
                assert.False(t, env.FileExists("~/.ddx.yml"))
            }
        })
    }
}
```

### 7. US-034: Offline Installation

#### Test Suite: `TestAcceptance_US034_OfflineInstallation`

```go
func TestAcceptance_US034_OfflineInstallation(t *testing.T) {
    tests := []struct {
        name         string
        platform     string
        packageType  string
    }{
        {
            name:        "offline_linux_tarball",
            platform:    "linux",
            packageType: "tar.gz",
        },
        {
            name:        "offline_windows_zip",
            platform:    "windows",
            packageType: "zip",
        },
        {
            name:        "offline_macos_tarball",
            platform:    "darwin",
            packageType: "tar.gz",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: No network connectivity
            env := setupOfflineEnvironment(tt.platform)
            defer env.Cleanup()

            // And: Offline package is available
            packagePath := env.PrepareOfflinePackage(tt.packageType)

            // When: I install from offline package
            result := env.RunOfflineInstall(packagePath)

            // Then: Installation succeeds
            assert.Equal(t, 0, result.ExitCode)

            // And: DDX is functional
            version := env.RunCommand("ddx version")
            assert.Equal(t, 0, version.ExitCode)

            // And: All library resources are included
            assert.True(t, env.DirectoryExists("~/.ddx/library"))
        })
    }
}
```

### 8. US-035: Installation Diagnostics

#### Test Suite: `TestAcceptance_US035_InstallationDiagnostics`

```go
func TestAcceptance_US035_InstallationDiagnostics(t *testing.T) {
    tests := []struct {
        name           string
        problemState   ProblemState
        expectedFixes  []string
    }{
        {
            name:         "network_connectivity_issue",
            problemState: NetworkConnectivityIssue(),
            expectedFixes: []string{
                "Check internet connection",
                "Verify proxy settings",
                "Try offline installation",
            },
        },
        {
            name:         "permission_issue",
            problemState: PermissionIssue(),
            expectedFixes: []string{
                "Check directory permissions",
                "Try installing to different location",
                "Verify user has write access",
            },
        },
        {
            name:         "path_configuration_issue",
            problemState: PathConfigurationIssue(),
            expectedFixes: []string{
                "Run 'ddx setup path'",
                "Restart shell session",
                "Manually add to PATH",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Given: Installation has specific problem
            env := setupProblemState(tt.problemState)
            defer env.Cleanup()

            // When: I run installation diagnostics
            result := env.RunCommand("ddx doctor --verbose")

            // Then: Problem is detected
            assert.Contains(t, result.Output, "Issues detected")

            // And: Remediation steps are suggested
            for _, fix := range tt.expectedFixes {
                assert.Contains(t, result.Output, fix)
            }

            // And: Diagnostic report is generated
            assert.Contains(t, result.Output, "Diagnostic Report")
            assert.Contains(t, result.Output, "System Information")
        })
    }
}
```

## Integration Test Specifications

### Installation Workflow Tests

#### Test Suite: `TestInstallationWorkflow`

```go
func TestInstallationWorkflow_EndToEnd(t *testing.T) {
    // Test complete installation workflow:
    // Platform detection → Binary download → Installation → PATH config → Verification

    platforms := []string{"linux", "darwin", "windows"}

    for _, platform := range platforms {
        t.Run(fmt.Sprintf("workflow_%s", platform), func(t *testing.T) {
            env := setupCleanEnvironment(platform)
            defer env.Cleanup()

            // Execute complete workflow
            result := env.ExecuteFullInstallation()

            // Validate each workflow step
            assert.True(t, result.PlatformDetected)
            assert.True(t, result.BinaryDownloaded)
            assert.True(t, result.BinaryInstalled)
            assert.True(t, result.PathConfigured)
            assert.True(t, result.InstallationVerified)
        })
    }
}
```

## Contract Test Specifications

### CLI Command Contracts

#### Test Suite: `TestInstallationCommandContracts`

```go
func TestInstallationCommandContracts(t *testing.T) {
    commands := []struct {
        command    string
        flags      []string
        exitCodes  []int
        outputFormat string
    }{
        {
            command:      "ddx doctor",
            flags:        []string{"--verbose", "--json"},
            exitCodes:    []int{0, 1}, // Success or problems found
            outputFormat: "structured",
        },
        {
            command:      "ddx self-update",
            flags:        []string{"--check", "--force", "--version"},
            exitCodes:    []int{0, 1, 5}, // Success, error, network error
            outputFormat: "progress",
        },
    }

    for _, cmd := range commands {
        t.Run(cmd.command, func(t *testing.T) {
            // Test command contract compliance
            testCommandContract(t, cmd)
        })
    }
}
```

## Performance Test Specifications

### Installation Speed Tests

#### Test Suite: `TestInstallationPerformance`

```go
func TestInstallationPerformance(t *testing.T) {
    tests := []struct {
        name             string
        networkSpeed     string
        expectedMaxTime  time.Duration
    }{
        {
            name:            "installation_on_10mbps",
            networkSpeed:    "10mbps",
            expectedMaxTime: 60 * time.Second,
        },
        {
            name:            "installation_on_1mbps",
            networkSpeed:    "1mbps",
            expectedMaxTime: 180 * time.Second,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            env := setupNetworkSpeedEnvironment(tt.networkSpeed)
            defer env.Cleanup()

            start := time.Now()
            result := env.ExecuteInstallation()
            duration := time.Since(start)

            assert.True(t, result.Success)
            assert.Less(t, duration, tt.expectedMaxTime)
        })
    }
}
```

## Test Environment Setup

### Platform Simulation

```go
type TestEnvironment struct {
    Platform     string
    Architecture string
    TempDir      string
    MockNetwork  *MockNetworkService
    MockGitHub   *MockGitHubAPI
}

func setupTestEnvironment(platform, arch string) *TestEnvironment {
    return &TestEnvironment{
        Platform:     platform,
        Architecture: arch,
        TempDir:      createTempTestDir(),
        MockNetwork:  newMockNetworkService(),
        MockGitHub:   newMockGitHubAPI(),
    }
}
```

### Test Data Requirements

1. **Mock GitHub Releases**
   - Platform-specific binaries (linux-amd64, darwin-arm64, windows-amd64)
   - Version metadata and checksums
   - Release notes and changelog

2. **Test Binaries**
   - Minimal functional DDX binaries for each platform
   - Test signature/checksum files

3. **Environment Configurations**
   - Clean shell environments
   - Various PATH configurations
   - Package manager mock services

## Success Criteria

### Test Execution Requirements

1. **All Tests Initially Fail (Red Phase)**
   - Tests compile successfully
   - Tests execute and fail with clear error messages
   - Failure messages indicate what needs implementation

2. **Incremental Implementation (Green Phase)**
   - Implement features to make tests pass one by one
   - No feature implementation without corresponding test
   - Maintain test isolation and reliability

3. **Test Coverage Targets**
   - 100% of user story acceptance criteria covered
   - All platform combinations tested
   - All error scenarios validated
   - Performance benchmarks verified

### Quality Gates

- Installation success rate >99% across all platforms
- Installation time <60 seconds on 10Mbps connection
- No admin privileges required for any installation method
- Clean uninstallation leaves no artifacts (when requested)
- Offline installation works without network connectivity

---

This test specification serves as the definitive guide for implementing and validating the cross-platform installation system. All tests must be written and failing before any implementation begins, following strict TDD principles.