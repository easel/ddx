package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Installer manages MCP server installation
type Installer struct {
	registry  *Registry
	claude    *ClaudeWrapper
	validator *Validator
	out       io.Writer
}

// NewInstaller creates a new installer
func NewInstaller() *Installer {
	return NewInstallerWithWriter(os.Stdout)
}

// NewInstallerWithWriter creates a new installer with a custom output writer
func NewInstallerWithWriter(w io.Writer) *Installer {
	return &Installer{
		claude:    NewClaudeWrapper(),
		validator: NewValidator(),
		out:       w,
	}
}

// Install installs an MCP server
func (i *Installer) Install(serverName string, opts InstallOptions) error {
	return i.InstallWithLibraryPath(serverName, opts, "")
}

// InstallWithLibraryPath installs an MCP server with a specific library path
func (i *Installer) InstallWithLibraryPath(serverName string, opts InstallOptions, libraryPath string) error {
	// Show installation start message
	fmt.Fprintf(i.out, "🔧 Installing %s MCP Server...\n\n", serverName)

	// Load registry
	if i.registry == nil {
		wd, _ := os.Getwd()
		registry, err := LoadRegistryWithLibraryPath("", wd, libraryPath)
		if err != nil {
			return fmt.Errorf("loading registry: %w", err)
		}
		i.registry = registry
	}

	// Get server definition
	server, err := i.registry.GetServer(serverName)
	if err != nil {
		return fmt.Errorf("getting server %s: %w", serverName, err)
	}

	// Show server information
	fmt.Fprintf(i.out, "📦 %s - %s\n", server.Name, server.Description)
	if len(server.Environment) > 0 {
		fmt.Fprintln(i.out, "\nℹ️ This server requires configuration:")
		for _, env := range server.Environment {
			if env.Required {
				fmt.Fprintf(i.out, "  - %s: %s\n", env.Name, env.Description)
			}
		}
		fmt.Fprintln(i.out)
	}

	// Validate and prompt for environment variables
	if len(server.Environment) > 0 {
		fmt.Fprintf(i.out, "🔧 Configuring server environment...\n")
	}
	if err := i.validateAndPromptEnvironment(server, opts.Environment, opts.Yes); err != nil {
		return fmt.Errorf("configuring environment: %w", err)
	}

	// Detect package manager and display
	packageManager := i.detectPackageManager()
	fmt.Fprintf(i.out, "Using package manager: %s\n", packageManager)

	// Check Claude CLI availability
	if err := i.claude.IsAvailable(); err != nil {
		return fmt.Errorf("Claude CLI not available: %w", err)
	}
	fmt.Fprintf(i.out, "📍 Detected Claude CLI available\n")

	// Check if already installed
	isInstalled := false
	if opts.ConfigPath != "" {
		// Check config file if path specified
		isInstalled = i.isServerInConfigFile(serverName, opts.ConfigPath)
	} else {
		// Check via Claude CLI
		if status, err := i.claude.GetServerStatus(serverName); err == nil && status.Installed {
			isInstalled = true
		}
	}

	if isInstalled && !opts.DryRun {
		fmt.Fprintf(i.out, "⚠️  %s is already installed.\n", serverName)
		fmt.Fprintf(i.out, "💡 To upgrade or reinstall, use: ddx mcp upgrade %s\n", serverName)
		return fmt.Errorf("%w: %s", ErrAlreadyInstalled, serverName)
	}

	if opts.DryRun {
		fmt.Fprintf(i.out, "🔍 Dry run mode - showing what would be done:\n")
		fmt.Fprintf(i.out, "   - Server: %s (%s)\n", serverName, server.Description)
		fmt.Fprintf(i.out, "   - Command: %s %s\n", server.Command.Executable, strings.Join(server.Command.Args, " "))
		if len(opts.Environment) > 0 {
			fmt.Fprintf(i.out, "   - Environment variables: %d configured\n", len(opts.Environment))
		}
		fmt.Fprintf(i.out, "   - Claude CLI: claude mcp add %s\n", serverName)
		return nil
	}

	// Install server via Claude CLI
	fmt.Fprintf(i.out, "📦 Installing server via Claude CLI...\n")
	if err := i.addServerWithConfig(serverName, server, opts); err != nil {
		return fmt.Errorf("installing server: %w", err)
	}

	// Success message with next steps
	fmt.Fprintf(i.out, "✅ %s MCP server installed Successfully!\n\n", serverName)
	fmt.Fprintf(i.out, "🚀 Next steps:\n")
	fmt.Fprintf(i.out, "  1. Restart Claude Code to load the server\n")
	fmt.Fprintf(i.out, "  2. Look for %s in Claude Code's MCP section\n", serverName)
	if serverName == "github" {
		fmt.Fprintf(i.out, "  3. Test with: \"Show my recent commits\"\n")
	} else if serverName == "filesystem" {
		fmt.Fprintf(i.out, "  3. Test with: \"List files in current directory\"\n")
	} else {
		fmt.Fprintf(i.out, "  3. Test the server functionality\n")
	}

	return nil
}

// validateAndPromptEnvironment validates and prompts for required environment variables
func (i *Installer) validateAndPromptEnvironment(server *Server, env map[string]string, skipPrompts bool) error {
	for _, envVar := range server.GetRequiredEnvironment() {
		value, exists := env[envVar.Name]

		// If value is missing and not in non-interactive mode, prompt for it
		if (!exists || value == "") && !skipPrompts {
			prompt := fmt.Sprintf("🔐 Enter %s", envVar.Name)
			if envVar.Description != "" {
				prompt = fmt.Sprintf("🔐 Enter %s (%s)", envVar.Name, envVar.Description)
			}

			if envVar.Sensitive {
				prompt += " (input will be masked)"
			}
			prompt += ": "

			fmt.Fprint(i.out, prompt)

			// For now, skip actual input reading in favor of clear error message
			return fmt.Errorf("interactive prompting not yet implemented - please provide %s via --env %s=value", envVar.Name, envVar.Name)
		}

		// If still missing after potential prompting, return error
		if !exists || value == "" {
			return fmt.Errorf("%w: %s", ErrMissingRequired, envVar.Name)
		}

		// Validate against regex if provided
		if envVar.Validation != "" {
			matched, err := regexp.MatchString(envVar.Validation, value)
			if err != nil {
				return fmt.Errorf("invalid validation regex for %s: %w", envVar.Name, err)
			}
			if !matched {
				return fmt.Errorf("%w: %s does not match pattern %s", ErrValidationFailed, envVar.Name, envVar.Validation)
			}
		}

		// Show validation success for sensitive variables
		if envVar.Sensitive && value != "" {
			fmt.Fprintf(i.out, "✅ %s validated\n", envVar.Name)
		}
	}
	return nil
}

// addServerWithConfig adds server and creates config file if path specified
func (i *Installer) addServerWithConfig(serverName string, server *Server, opts InstallOptions) error {
	// If config path is specified, create/update the Claude config file
	if opts.ConfigPath != "" {
		return i.createClaudeConfig(serverName, server, opts)
	}

	// Otherwise use regular Claude CLI
	return i.claude.AddServer(serverName, server.Command.Executable, server.Command.Args, opts.Environment)
}

// createClaudeConfig creates a Claude config file at the specified path
func (i *Installer) createClaudeConfig(serverName string, server *Server, opts InstallOptions) error {
	configDir := filepath.Dir(opts.ConfigPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Create basic Claude MCP configuration
	config := map[string]interface{}{
		"mcpServers": map[string]interface{}{
			serverName: map[string]interface{}{
				"command": server.Command.Executable,
				"args":    server.Command.Args,
			},
		},
	}

	// Add environment variables if present
	if len(opts.Environment) > 0 {
		serverConfig := config["mcpServers"].(map[string]interface{})[serverName].(map[string]interface{})
		serverConfig["env"] = opts.Environment
	}

	// Write config file
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(opts.ConfigPath, jsonData, 0644); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// isServerInConfigFile checks if a server is already configured in the config file
func (i *Installer) isServerInConfigFile(serverName, configPath string) bool {
	// Check if config file exists
	if _, err := os.Stat(configPath); err != nil {
		return false // File doesn't exist, so server isn't installed
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return false // Can't read file, assume not installed
	}

	// Parse JSON config
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return false // Invalid JSON, assume not installed
	}

	// Check if mcpServers section exists and contains our server
	if mcpServers, ok := config["mcpServers"].(map[string]interface{}); ok {
		_, exists := mcpServers[serverName]
		return exists
	}

	return false
}

// detectPackageManager detects which package manager is in use
func (i *Installer) detectPackageManager() string {
	// Check for lock files in order of preference
	if _, err := os.Stat("pnpm-lock.yaml"); err == nil {
		return "pnpm"
	}
	if _, err := os.Stat("yarn.lock"); err == nil {
		return "yarn"
	}
	if _, err := os.Stat("package-lock.json"); err == nil {
		return "npm"
	}
	// Default to npm if no lock file found
	return "npm"
}

// ensurePackageJSON creates a basic package.json if it doesn't exist
func (i *Installer) ensurePackageJSON() error {
	// Check if package.json already exists
	if _, err := os.Stat("package.json"); err == nil {
		return nil // Already exists
	}

	// Create a basic package.json
	packageJSON := map[string]interface{}{
		"name":         "mcp-servers",
		"version":      "1.0.0",
		"description":  "MCP server dependencies for DDx",
		"private":      true,
		"dependencies": make(map[string]string),
	}

	data, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling package.json: %w", err)
	}

	if err := os.WriteFile("package.json", data, 0644); err != nil {
		return fmt.Errorf("writing package.json: %w", err)
	}

	return nil
}

// validateEnvironment validates required environment variables (legacy method)
func (i *Installer) validateEnvironment(server *Server, env map[string]string) error {
	for _, envVar := range server.GetRequiredEnvironment() {
		value, exists := env[envVar.Name]
		if !exists || value == "" {
			return fmt.Errorf("%w: %s", ErrMissingRequired, envVar.Name)
		}

		// Validate against regex if provided
		if envVar.Validation != "" {
			matched, err := regexp.MatchString(envVar.Validation, value)
			if err != nil {
				return fmt.Errorf("invalid validation regex for %s: %w", envVar.Name, err)
			}
			if !matched {
				return fmt.Errorf("%w: %s does not match pattern %s", ErrValidationFailed, envVar.Name, envVar.Validation)
			}
		}
	}
	return nil
}

// Validator provides input validation
type Validator struct {
	patterns map[string]*regexp.Regexp
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		patterns: map[string]*regexp.Regexp{
			"serverName": regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`),
			"envName":    regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`),
		},
	}
}

// ValidateServerName validates a server name
func (v *Validator) ValidateServerName(name string) error {
	if name == "" {
		return ErrEmptyServerName
	}

	if !v.patterns["serverName"].MatchString(name) {
		return fmt.Errorf("invalid server name: %s (must be lowercase with hyphens)", name)
	}

	// Check for path traversal
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return ErrPathTraversal
	}

	return nil
}

// ValidateEnvironment validates environment variables
func (v *Validator) ValidateEnvironment(env map[string]string) error {
	for key, value := range env {
		if !v.patterns["envName"].MatchString(key) {
			return fmt.Errorf("invalid environment variable name: %s", key)
		}

		// Check for injection attempts
		if containsShellInjection(value) {
			return fmt.Errorf("%w in environment variable %s", ErrInjectionAttempt, key)
		}
	}
	return nil
}

// ValidatePath validates a file path
func (v *Validator) ValidatePath(path string) error {
	// Check absolute path first
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %s", path)
	}

	// Check for path traversal
	if strings.Contains(path, "..") {
		return ErrPathTraversal
	}

	// Clean and check again
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return ErrPathTraversal
	}

	return nil
}

// containsShellInjection checks for potential shell injection
func containsShellInjection(value string) bool {
	dangerousPatterns := []string{
		"$(", "`", ";", "&&", "||", ">", "<", "|",
		"\n", "\r", "\x00",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(value, pattern) {
			return true
		}
	}

	return false
}

// MaskSensitive masks sensitive values for display
func MaskSensitive(value string, sensitive bool) string {
	if !sensitive || value == "" {
		return value
	}

	// Show first few characters for debugging
	if len(value) > 8 {
		return value[:4] + "***"
	}
	return "***"
}

// StatusChecker manages server status checking
type StatusChecker struct {
	out    io.Writer
	claude *ClaudeWrapper
}

// NewStatusChecker creates a new status checker
func NewStatusChecker() *StatusChecker {
	return NewStatusCheckerWithWriter(os.Stdout)
}

// NewStatusCheckerWithWriter creates a new status checker with a custom output writer
func NewStatusCheckerWithWriter(w io.Writer) *StatusChecker {
	return &StatusChecker{
		out:    w,
		claude: NewClaudeWrapper(),
	}
}

// Check checks the status of servers
func (sc *StatusChecker) Check(opts StatusOptions) error {
	if opts.ServerName != "" {
		return sc.checkServer(opts.ServerName, opts)
	}
	return sc.checkAllServers(opts)
}

// checkServer checks status of a specific server
func (sc *StatusChecker) checkServer(serverName string, opts StatusOptions) error {
	fmt.Fprintf(sc.out, "📊 Status for %s MCP server:\n", serverName)

	// Check if Claude CLI is available
	if err := sc.claude.IsAvailable(); err != nil {
		fmt.Fprintf(sc.out, "⚠️  Claude CLI not available: %v\n", err)
		return nil
	}

	// Get server status from Claude CLI
	status, err := sc.claude.GetServerStatus(serverName)
	if err != nil {
		fmt.Fprintf(sc.out, "⚠️  Error checking server status: %v\n", err)
		return nil
	}

	if status.Installed {
		fmt.Fprintf(sc.out, "✅ %s - Installed\n", serverName)
		if status.Running {
			fmt.Fprintf(sc.out, "   Status: Running ✓\n")
		} else {
			fmt.Fprintf(sc.out, "   Status: Installed but not connected\n")
		}
		if status.Version != "unknown" {
			fmt.Fprintf(sc.out, "   Version: %s\n", status.Version)
		}
		if len(status.Errors) > 0 {
			fmt.Fprintf(sc.out, "   Errors: %s\n", strings.Join(status.Errors, ", "))
		}
	} else {
		fmt.Fprintln(sc.out, "⚠️  Server not installed")
	}

	return nil
}

// checkAllServers checks status of all servers
func (sc *StatusChecker) checkAllServers(opts StatusOptions) error {
	fmt.Fprintln(sc.out, "📊 MCP Server Status")
	fmt.Fprintln(sc.out)

	// Check if Claude CLI is available
	if err := sc.claude.IsAvailable(); err != nil {
		fmt.Fprintf(sc.out, "⚠️  Claude CLI not available: %v\n", err)
		return nil
	}

	// List installed servers from Claude CLI
	servers, err := sc.claude.ListServers()
	if err != nil {
		fmt.Fprintf(sc.out, "⚠️  Error listing servers: %v\n", err)
		return nil
	}

	if len(servers) == 0 {
		fmt.Fprintln(sc.out, "⚠️  No servers installed")
		return nil
	}

	fmt.Fprintln(sc.out, "Installed servers:")
	for name, connected := range servers {
		status := "Installed"
		if connected {
			status = "Running ✓"
		}
		fmt.Fprintf(sc.out, "  ✅ %-15s - %s\n", name, status)
	}
	fmt.Fprintln(sc.out)

	return nil
}
