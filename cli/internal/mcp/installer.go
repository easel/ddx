package mcp

import (
	"encoding/json"
	"fmt"
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
}

// NewInstaller creates a new installer
func NewInstaller() *Installer {
	return &Installer{
		validator: NewValidator(),
	}
}

// Install installs an MCP server
func (i *Installer) Install(serverName string, opts InstallOptions) error {
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

	// Validate environment variables
	if err := i.validateEnvironment(server, opts.Environment); err != nil {
		return fmt.Errorf("validating environment: %w", err)
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
		claudePath = detected[0].ConfigPath
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
	}

	// Create server configuration
	serverConfig := ServerConfig{
		Command: server.Command.Executable,
		Args:    server.Command.Args,
		Env:     opts.Environment,
	}

	// Add server to config
	if err := i.config.AddServer(serverName, serverConfig); err != nil {
		return fmt.Errorf("adding server: %w", err)
	}

	// Save config (unless dry run)
	if !opts.DryRun {
		if err := i.config.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}
	}

	return nil
}

// validateEnvironment validates required environment variables
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
func (cm *ConfigManager) RemoveServer(name string) error {
	delete(cm.config.MCPServers, name)
	return nil
}

// GetServer gets a server configuration
func (cm *ConfigManager) GetServer(name string) (ServerConfig, bool) {
	config, exists := cm.config.MCPServers[name]
	return config, exists
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
