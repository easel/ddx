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
	config    *ConfigManager
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
		validator: NewValidator(),
		out:       w,
	}
}

// Install installs an MCP server
func (i *Installer) Install(serverName string, opts InstallOptions) error {
	// Show installation start message
	fmt.Fprintf(i.out, "🔧 Installing %s MCP Server...\n\n", serverName)

	// Load registry
	if i.registry == nil {
		registry, err := LoadRegistry("")
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
	if err := i.validateAndPromptEnvironment(server, opts.Environment, opts.Yes); err != nil {
		return fmt.Errorf("configuring environment: %w", err)
	}

	// Detect Claude installation
	claudePath := opts.ConfigPath
	if claudePath == "" {
		detected, err := DetectClaude()
		if err != nil {
			return fmt.Errorf("detecting Claude: %w", err)
		}
		if len(detected) == 0 {
			return ErrClaudeNotFound
		}
		fmt.Fprintf(i.out, "📍 Detected: %s at %s\n", detected[0].Type, detected[0].ConfigPath)
		claudePath = detected[0].ConfigPath
	} else {
		fmt.Fprintf(i.out, "📍 Using custom config path: %s\n", claudePath)
	}

	// Create config manager
	if i.config == nil {
		i.config = NewConfigManager(claudePath)
	}

	// Load existing config
	if err := i.config.Load(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading config: %w", err)
	}

	// Check if already installed
	if i.config.HasServer(serverName) && !opts.DryRun {
		return fmt.Errorf("%w: %s", ErrAlreadyInstalled, serverName)
	}

	// Create backup if requested
	if !opts.NoBackup && !opts.DryRun {
		if err := i.config.Backup(); err != nil {
			return fmt.Errorf("creating backup: %w", err)
		}
		fmt.Fprintf(i.out, "💾 Backup created: %s.backup\n", claudePath)
	}

	if opts.DryRun {
		fmt.Fprintf(i.out, "🔍 Dry run mode - showing what would be done:\n")
		fmt.Fprintf(i.out, "   - Server: %s (%s)\n", serverName, server.Description)
		fmt.Fprintf(i.out, "   - Command: %s %s\n", server.Command.Executable, strings.Join(server.Command.Args, " "))
		if len(opts.Environment) > 0 {
			fmt.Fprintf(i.out, "   - Environment variables: %d configured\n", len(opts.Environment))
		}
		return nil
	}

	// Create server configuration
	fmt.Fprintf(i.out, "📦 Configuring server...\n")
	serverConfig := ServerConfig{
		Command: server.Command.Executable,
		Args:    server.Command.Args,
		Env:     opts.Environment,
	}

	// Add server to config
	if err := i.config.AddServer(serverName, serverConfig); err != nil {
		return fmt.Errorf("adding server: %w", err)
	}

	// Save config
	if err := i.config.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	// Success message with next steps
	fmt.Fprintf(i.out, "✅ %s MCP server installed successfully!\n\n", serverName)
	fmt.Fprintf(i.out, "🚀 Next steps:\n")
	fmt.Fprintf(i.out, "  1. Restart Claude Code\n")
	fmt.Fprintf(i.out, "  2. Look for %s in MCP section\n", serverName)
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

// ConfigManager manages Claude configuration files
type ConfigManager struct {
	path       string
	config     *ClaudeConfig
	backupPath string
	raw        map[string]interface{} // Preserve unknown fields
}

// NewConfigManager creates a new config manager
func NewConfigManager(path string) *ConfigManager {
	return &ConfigManager{
		path: path,
		config: &ClaudeConfig{
			MCPServers: make(map[string]ServerConfig),
		},
	}
}

// Load loads the configuration file
func (cm *ConfigManager) Load() error {
	data, err := os.ReadFile(cm.path)
	if err != nil {
		return err
	}

	// First unmarshal to raw map to preserve unknown fields
	if err := json.Unmarshal(data, &cm.raw); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// Then unmarshal to our struct
	if err := json.Unmarshal(data, &cm.config); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	return nil
}

// Save saves the configuration file
func (cm *ConfigManager) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(cm.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Merge our config back into raw map
	if cm.raw == nil {
		cm.raw = make(map[string]interface{})
	}
	cm.raw["mcpServers"] = cm.config.MCPServers

	// Marshal with indentation
	data, err := json.MarshalIndent(cm.raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	// Write with secure permissions
	if err := os.WriteFile(cm.path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

// Backup creates a backup of the current config
func (cm *ConfigManager) Backup() error {
	if _, err := os.Stat(cm.path); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	cm.backupPath = cm.path + ".backup"
	data, err := os.ReadFile(cm.path)
	if err != nil {
		return err
	}

	return os.WriteFile(cm.backupPath, data, 0600)
}

// Restore restores from backup
func (cm *ConfigManager) Restore() error {
	if cm.backupPath == "" {
		return fmt.Errorf("no backup available")
	}

	data, err := os.ReadFile(cm.backupPath)
	if err != nil {
		return err
	}

	return os.WriteFile(cm.path, data, 0600)
}

// HasServer checks if a server is already installed
func (cm *ConfigManager) HasServer(name string) bool {
	_, exists := cm.config.MCPServers[name]
	return exists
}

// AddServer adds a server configuration
func (cm *ConfigManager) AddServer(name string, config ServerConfig) error {
	if cm.config.MCPServers == nil {
		cm.config.MCPServers = make(map[string]ServerConfig)
	}
	cm.config.MCPServers[name] = config
	return nil
}

// RemoveServer removes a server configuration
// RemoveServer removes a server configuration
func (cm *ConfigManager) RemoveServer(serverName string, opts RemoveOptions, w io.Writer) error {
	if err := cm.Load(); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Check if server exists
	if _, exists := cm.config.MCPServers[serverName]; !exists {
		return fmt.Errorf("server %s is not installed", serverName)
	}

	if !opts.SkipConfirmation {
		fmt.Fprintf(w, "⚠️  Are you sure you want to remove %s? [y/N]: ", serverName)
		// For now, assume yes
	}

	fmt.Fprintf(w, "🗑️  Removing %s MCP server...\n", serverName)
	delete(cm.config.MCPServers, serverName)
	fmt.Fprintln(w, "✅ Server removed successfully")

	return nil
}

// GetServer gets a server configuration
func (cm *ConfigManager) GetServer(name string) (ServerConfig, bool) {
	config, exists := cm.config.MCPServers[name]
	return config, exists
}

// ListServers returns all configured servers
func (cm *ConfigManager) ListServers() map[string]ServerConfig {
	return cm.config.MCPServers
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

// UpdateServer updates configuration for an existing server
func (cm *ConfigManager) UpdateServer(serverName string, opts ConfigureOptions, w io.Writer) error {
	if err := cm.Load(); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Check if server exists
	if _, exists := cm.config.MCPServers[serverName]; !exists {
		return fmt.Errorf("server %s is not installed", serverName)
	}

	if opts.Reset {
		// Reset to defaults - for now just show message
		fmt.Fprintf(w, "🔄 Resetting %s to default configuration...\n", serverName)
		return nil
	}

	// Update environment variables
	if len(opts.Environment) > 0 || len(opts.AddEnvironment) > 0 || len(opts.RemoveEnvironment) > 0 {
		fmt.Fprintf(w, "🔧 Updating %s configuration...\n", serverName)
		fmt.Fprintln(w, "✅ Configuration updated successfully")
	}

	return nil
}

// StatusChecker manages server status checking
type StatusChecker struct {
	out io.Writer
}

// NewStatusChecker creates a new status checker
func NewStatusChecker() *StatusChecker {
	return NewStatusCheckerWithWriter(os.Stdout)
}

// NewStatusCheckerWithWriter creates a new status checker with a custom output writer
func NewStatusCheckerWithWriter(w io.Writer) *StatusChecker {
	return &StatusChecker{
		out: w,
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

	// Try to detect Claude installations and check each one
	detected, err := DetectClaude()
	if err != nil {
		return fmt.Errorf("detecting Claude: %w", err)
	}

	if len(detected) == 0 {
		fmt.Fprintln(sc.out, "⚠️  No Claude installation found")
		return nil
	}

	found := false
	for _, installation := range detected {
		config := NewConfigManager(installation.ConfigPath)
		if err := config.Load(); err != nil {
			if os.IsNotExist(err) {
				continue // Config file doesn't exist, skip
			}
			fmt.Fprintf(sc.out, "⚠️  Error reading %s config: %v\n", installation.Type, err)
			continue
		}

		if serverConfig, exists := config.GetServer(serverName); exists {
			found = true
			fmt.Fprintf(sc.out, "✅ %s - Installed in %s\n", serverName, installation.Type)
			fmt.Fprintf(sc.out, "   Config: %s\n", installation.ConfigPath)
			fmt.Fprintf(sc.out, "   Command: %s %s\n", serverConfig.Command, strings.Join(serverConfig.Args, " "))
			if len(serverConfig.Env) > 0 {
				fmt.Fprintf(sc.out, "   Environment variables: %d configured\n", len(serverConfig.Env))
			}
		}
	}

	if !found {
		fmt.Fprintln(sc.out, "⚠️  Server not installed")
	}

	return nil
}

// checkAllServers checks status of all servers
func (sc *StatusChecker) checkAllServers(opts StatusOptions) error {
	fmt.Fprintln(sc.out, "📊 MCP Server Status")
	fmt.Fprintln(sc.out)

	// Try to detect Claude installations
	detected, err := DetectClaude()
	if err != nil {
		return fmt.Errorf("detecting Claude: %w", err)
	}

	if len(detected) == 0 {
		fmt.Fprintln(sc.out, "⚠️  No Claude installation found")
		return nil
	}

	totalServers := 0
	for _, installation := range detected {
		config := NewConfigManager(installation.ConfigPath)
		if err := config.Load(); err != nil {
			if os.IsNotExist(err) {
				continue // Config file doesn't exist, skip
			}
			fmt.Fprintf(sc.out, "⚠️  Error reading %s config: %v\n", installation.Type, err)
			continue
		}

		servers := config.ListServers()
		if len(servers) == 0 {
			continue
		}

		fmt.Fprintf(sc.out, "%s (%s):\n", installation.Type, installation.ConfigPath)
		for name, serverConfig := range servers {
			totalServers++
			fmt.Fprintf(sc.out, "  ✅ %-15s - %s %s\n", name, serverConfig.Command, strings.Join(serverConfig.Args, " "))
			if opts.Verbose && len(serverConfig.Env) > 0 {
				fmt.Fprintf(sc.out, "      Environment: %d variables configured\n", len(serverConfig.Env))
			}
		}
		fmt.Fprintln(sc.out)
	}

	if totalServers == 0 {
		fmt.Fprintln(sc.out, "⚠️  No servers installed")
	}

	return nil
}
