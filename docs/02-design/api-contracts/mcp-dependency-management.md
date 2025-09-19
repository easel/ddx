# MCP Dependency Management API Contracts

> **Last Updated**: 2025-09-18
> **Status**: Design Phase
> **Owner**: DDx Team

## Overview

This document defines the API contracts for MCP (Model Context Protocol) server dependency management in DDx. These contracts specify how DDx manages npm packages for MCP servers, ensuring project-local installation and proper Claude configuration.

## Core Interfaces

### 1. Package Manager Abstraction

#### PackageManagerDetector Service

```go
type PackageManagerDetector interface {
    // DetectPackageManager identifies the package manager from lock files
    DetectPackageManager(projectPath string) (PackageManagerType, error)

    // GetConfiguredPackageManager reads package manager from .ddx.yml
    GetConfiguredPackageManager(projectPath string) (PackageManagerType, error)

    // ResolvePackageManager determines which package manager to use
    ResolvePackageManager(projectPath string) (PackageManagerType, error)
}

type PackageManagerType string

const (
    PackageManagerNPM  PackageManagerType = "npm"
    PackageManagerPNPM PackageManagerType = "pnpm"
    PackageManagerYarn PackageManagerType = "yarn"
    PackageManagerBun  PackageManagerType = "bun"
)
```

#### PackageManager Interface

```go
type PackageManager interface {
    // GetType returns the package manager type
    GetType() PackageManagerType

    // EnsurePackageJson ensures package.json exists in project root
    EnsurePackageJson(projectPath string) error

    // InstallDependency installs a package as devDependency
    InstallDependency(projectPath, packageName, version string) error

    // RemoveDependency removes a package
    RemoveDependency(projectPath, packageName string) error

    // IsInstalled checks if package is installed locally
    IsInstalled(projectPath, packageName string) (bool, string, error)

    // GetInstalledVersion returns installed package version
    GetInstalledVersion(projectPath, packageName string) (string, error)

    // ValidateInstallation verifies package can be executed
    ValidateInstallation(projectPath, packageName string) error

    // GetExecutorCommand returns the command to execute packages (npx, pnpx, yarn dlx, bunx)
    GetExecutorCommand() (string, []string)

    // GetLockFileName returns the lock file name for this package manager
    GetLockFileName() string
}
```

#### PackageManagerFactory

```go
type PackageManagerFactory interface {
    // CreatePackageManager creates appropriate package manager instance
    CreatePackageManager(pmType PackageManagerType) (PackageManager, error)
}
```

#### Package Management Requests/Responses

```go
// InstallDependencyRequest represents a package installation request
type InstallDependencyRequest struct {
    ProjectPath string `json:"projectPath" validate:"required,dir"`
    PackageName string `json:"packageName" validate:"required,npm_package"`
    Version     string `json:"version,omitempty" validate:"omitempty,semver"`
    Force       bool   `json:"force,omitempty"`
}

// InstallDependencyResponse represents the installation result
type InstallDependencyResponse struct {
    Success         bool   `json:"success"`
    InstalledVersion string `json:"installedVersion"`
    ExecutablePath   string `json:"executablePath"`
    ErrorMessage     string `json:"errorMessage,omitempty"`
}
```

### 2. MCP Server Registry Interface

Enhanced MCP server definitions to include npm package information:

```go
// MCPServerDefinition extends existing Server struct
type MCPServerDefinition struct {
    Name        string            `yaml:"name" validate:"required"`
    Description string            `yaml:"description" validate:"required"`
    Category    string            `yaml:"category" validate:"required"`

    // NPM Package Information
    Package     PackageInfo       `yaml:"package" validate:"required"`

    // Command configuration (for Claude)
    Command     CommandInfo       `yaml:"command" validate:"required"`

    // Environment requirements
    Environment []EnvironmentVar  `yaml:"environment,omitempty"`

    // Metadata
    Author      string            `yaml:"author,omitempty"`
    Version     string            `yaml:"version,omitempty"`
    Tags        []string          `yaml:"tags,omitempty"`
}

// PackageInfo defines npm package details
type PackageInfo struct {
    Name         string `yaml:"name" validate:"required,npm_package"`
    Version      string `yaml:"version,omitempty" validate:"omitempty,semver"`
    Registry     string `yaml:"registry,omitempty" validate:"omitempty,url"`
    InstallType  string `yaml:"installType" validate:"required,oneof=devDependency dependency"`
}

// CommandInfo defines how to execute the installed package
type CommandInfo struct {
    Executable string   `yaml:"executable" validate:"required"`
    Args       []string `yaml:"args,omitempty"`
    WorkingDir string   `yaml:"workingDir,omitempty"`
}
```

### 3. Installation Service Interface

```go
type MCPInstallationService interface {
    // InstallServer installs MCP server with dependencies
    InstallServer(req InstallServerRequest) (*InstallServerResponse, error)

    // UninstallServer removes MCP server and cleans up dependencies
    UninstallServer(req UninstallServerRequest) (*UninstallServerResponse, error)

    // UpdateServer updates MCP server to newer version
    UpdateServer(req UpdateServerRequest) (*UpdateServerResponse, error)

    // ValidateServer checks if server is properly installed and functional
    ValidateServer(req ValidateServerRequest) (*ValidateServerResponse, error)
}
```

#### Installation Request/Response Contracts

```go
// InstallServerRequest defines server installation parameters
type InstallServerRequest struct {
    ServerName    string            `json:"serverName" validate:"required"`
    ProjectPath   string            `json:"projectPath" validate:"required,dir"`
    Environment   map[string]string `json:"environment,omitempty"`
    Force         bool              `json:"force,omitempty"`
    DryRun        bool              `json:"dryRun,omitempty"`
    SkipValidation bool             `json:"skipValidation,omitempty"`
}

// InstallServerResponse provides installation results
type InstallServerResponse struct {
    Success        bool              `json:"success"`
    ServerName     string            `json:"serverName"`
    PackageInfo    InstalledPackage  `json:"packageInfo"`
    ConfigPath     string            `json:"configPath"`
    ValidationResult ValidationResult `json:"validationResult"`
    ErrorMessage   string            `json:"errorMessage,omitempty"`
    Warnings       []string          `json:"warnings,omitempty"`
}

// InstalledPackage contains package installation details
type InstalledPackage struct {
    Name           string `json:"name"`
    Version        string `json:"version"`
    InstallPath    string `json:"installPath"`
    ExecutablePath string `json:"executablePath"`
}

// ValidationResult contains server validation status
type ValidationResult struct {
    Reachable      bool     `json:"reachable"`
    Version        string   `json:"version,omitempty"`
    Capabilities   []string `json:"capabilities,omitempty"`
    ErrorMessage   string   `json:"errorMessage,omitempty"`
}
```

## Configuration Generation Contracts

### 1. Claude Configuration Interface

```go
type ClaudeConfigGenerator interface {
    // GenerateServerConfig creates Claude MCP server configuration
    GenerateServerConfig(server MCPServerDefinition, projectPath string, pmType PackageManagerType) (ServerConfig, error)

    // ValidateConfiguration checks if generated config is valid
    ValidateConfiguration(config ServerConfig) error

    // SubstituteVariables replaces template variables with actual values
    SubstituteVariables(args []string, projectPath string) ([]string, error)

    // AdaptCommandForPackageManager adjusts command for specific package manager
    AdaptCommandForPackageManager(command string, args []string, pmType PackageManagerType) (string, []string)
}

// ServerConfig represents Claude MCP server configuration
type ServerConfig struct {
    Command string            `json:"command" validate:"required"`
    Args    []string          `json:"args,omitempty"`
    Env     map[string]string `json:"env,omitempty"`
}

// PackageManagerCommands defines executor commands per package manager
type PackageManagerCommands struct {
    NPM  CommandExecutor `json:"npm"`
    PNPM CommandExecutor `json:"pnpm"`
    Yarn CommandExecutor `json:"yarn"`
    Bun  CommandExecutor `json:"bun"`
}

type CommandExecutor struct {
    Command string   `json:"command"`
    Prefix  []string `json:"prefix,omitempty"`
}
```

### 2. Variable Substitution Rules

| Variable | Substitution | Example |
|----------|-------------|---------|
| `$PWD` | Project absolute path | `/home/user/project` |
| `$NODE_MODULES` | Local node_modules path | `/home/user/project/node_modules` |
| `$NPM_BIN` | Local npm bin directory | `/home/user/project/node_modules/.bin` |
| `$PROJECT_NAME` | Project directory name | `my-project` |

## Error Handling Contracts

### Standard Error Types

```go
// Error categories for MCP operations
const (
    ErrCategoryValidation    = "validation"
    ErrCategoryNetwork      = "network"
    ErrCategoryFileSystem   = "filesystem"
    ErrCategoryNPM          = "npm"
    ErrCategoryConfiguration = "configuration"
    ErrCategoryExecution    = "execution"
)

// MCPError provides structured error information
type MCPError struct {
    Category    string            `json:"category"`
    Code        string            `json:"code"`
    Message     string            `json:"message"`
    Details     map[string]string `json:"details,omitempty"`
    Recoverable bool              `json:"recoverable"`
    Suggestions []string          `json:"suggestions,omitempty"`
}
```

### Error Response Examples

```json
{
  "category": "npm",
  "code": "PACKAGE_NOT_FOUND",
  "message": "npm package '@modelcontextprotocol/server-invalid' not found",
  "details": {
    "packageName": "@modelcontextprotocol/server-invalid",
    "registry": "https://registry.npmjs.org/"
  },
  "recoverable": true,
  "suggestions": [
    "Check package name spelling",
    "Verify package exists in npm registry",
    "Try different package version"
  ]
}
```

## Package.json Management Contracts

### 1. Package.json Structure

```json
{
  "name": "project-name",
  "version": "1.0.0",
  "description": "Project with DDx MCP servers",
  "devDependencies": {
    "@modelcontextprotocol/server-filesystem": "^1.0.0",
    "@modelcontextprotocol/server-sequential-thinking": "^1.0.0",
    "@playwright/mcp": "^1.0.0"
  },
  "ddx": {
    "mcpServers": {
      "filesystem": {
        "package": "@modelcontextprotocol/server-filesystem",
        "version": "^1.0.0",
        "installedAt": "2025-09-18T10:30:00Z"
      }
    }
  }
}
```

### 2. Package.json Operations

```go
type PackageJsonManager interface {
    // Read package.json or create minimal structure
    ReadOrCreate(projectPath string) (*PackageJson, error)

    // Add MCP server dependency
    AddMCPDependency(pkg *PackageJson, serverName, packageName, version string) error

    // Remove MCP server dependency
    RemoveMCPDependency(pkg *PackageJson, serverName string) error

    // Write package.json back to filesystem
    Write(pkg *PackageJson, projectPath string) error

    // Validate package.json structure
    Validate(pkg *PackageJson) error
}
```

## Validation Contracts

### 1. Pre-Installation Validation

```go
type PreInstallValidator interface {
    // ValidateEnvironment checks system requirements
    ValidateEnvironment() error

    // ValidateProject checks project setup
    ValidateProject(projectPath string) error

    // ValidateServerDefinition checks server config
    ValidateServerDefinition(server MCPServerDefinition) error

    // ValidatePackageName checks npm package name format
    ValidatePackageName(packageName string) error
}
```

### 2. Post-Installation Validation

```go
type PostInstallValidator interface {
    // ValidateInstallation verifies package was installed correctly
    ValidateInstallation(projectPath, packageName string) error

    // ValidateExecution tests if server can be executed
    ValidateExecution(config ServerConfig, projectPath string) error

    // ValidateClaudeConfig checks Claude configuration validity
    ValidateClaudeConfig(configPath string) error
}
```

## Integration Test Contracts

### 1. Test Scenarios

```go
type MCPInstallationTestSuite interface {
    // Test complete installation flow
    TestFullInstallationFlow(serverName string) error

    // Test installation with missing package.json
    TestInstallationWithoutPackageJson(serverName string) error

    // Test installation with existing dependencies
    TestInstallationWithExistingDeps(serverName string) error

    // Test server functionality after installation
    TestServerFunctionality(serverName string) error

    // Test uninstallation and cleanup
    TestUninstallationCleanup(serverName string) error
}
```

### 2. Mock Interfaces

```go
type MockNPMClient interface {
    NPMClient

    // Test-specific methods
    SetPackageAvailable(packageName, version string, available bool)
    SetInstallSuccess(packageName string, success bool)
    GetInstallHistory() []InstallRequest
}

type MockClaudeConfig interface {
    ClaudeConfigManager

    // Test-specific methods
    SetConfigPath(path string)
    GetConfigUpdates() []ConfigUpdate
    SimulateClaudeRestart()
}
```

## Compatibility Contracts

### 1. Version Compatibility

```go
// CompatibilityMatrix defines version requirements
type CompatibilityMatrix struct {
    DDXVersion     string   `yaml:"ddxVersion"`
    NodeVersions   []string `yaml:"nodeVersions"`
    NPMVersions    []string `yaml:"npmVersions"`
    ClaudeVersions []string `yaml:"claudeVersions"`
    Platforms      []string `yaml:"platforms"`
}

// CompatibilityChecker validates environment compatibility
type CompatibilityChecker interface {
    CheckDDXVersion(required string) error
    CheckNodeVersion(required []string) error
    CheckNPMVersion(required []string) error
    CheckPlatform(supported []string) error
    CheckOverallCompatibility(matrix CompatibilityMatrix) error
}
```

### 2. Migration Contracts

```go
// MigrationService handles upgrades from old MCP configurations
type MigrationService interface {
    // DetectLegacyConfiguration identifies old-style MCP configs
    DetectLegacyConfiguration(projectPath string) (*LegacyConfig, error)

    // MigrateToLocalDependencies converts to new dependency-managed approach
    MigrateToLocalDependencies(legacy *LegacyConfig, projectPath string) error

    // BackupConfiguration creates backup before migration
    BackupConfiguration(projectPath string) (string, error)

    // ValidateMigration ensures migration was successful
    ValidateMigration(projectPath string) error
}
```

## CLI Contract Extensions

### 1. Enhanced MCP Commands

```bash
# Install with dependency management
ddx mcp install <server-name> [--env KEY=VALUE] [--force] [--dry-run]

# Show installation status including dependencies
ddx mcp status [server-name] [--check-deps] [--verbose]

# Update server and dependencies
ddx mcp update [server-name] [--check-only] [--force]

# Remove server and clean dependencies
ddx mcp remove <server-name> [--keep-deps] [--purge]

# List available servers with dependency info
ddx mcp list [--show-deps] [--category CATEGORY]

# Validate local installation
ddx mcp validate [server-name] [--fix-issues]
```

### 2. Command Output Contracts

```json
{
  "command": "ddx mcp install filesystem",
  "status": "success",
  "result": {
    "serverName": "filesystem",
    "package": {
      "name": "@modelcontextprotocol/server-filesystem",
      "version": "1.0.0",
      "installedAt": "2025-09-18T10:30:00Z"
    },
    "configuration": {
      "path": ".claude/settings.local.json",
      "updated": true
    },
    "validation": {
      "reachable": true,
      "capabilities": ["filesystem_read", "filesystem_write"]
    }
  },
  "warnings": [],
  "nextSteps": [
    "Restart Claude Code to load the new MCP server",
    "Test with: 'List files in current directory'"
  ]
}
```

These contracts ensure consistent, testable, and maintainable MCP server dependency management across the DDx toolkit.