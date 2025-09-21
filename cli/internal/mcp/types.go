// Package mcp provides MCP server management functionality for DDx
package mcp

import "time"

// Server represents an MCP server definition
type Server struct {
	Name          string           `yaml:"name" json:"name"`
	Description   string           `yaml:"description" json:"description"`
	Category      string           `yaml:"category" json:"category"`
	Author        string           `yaml:"author" json:"author"`
	Version       string           `yaml:"version" json:"version"`
	Tags          []string         `yaml:"tags" json:"tags"`
	Command       CommandSpec      `yaml:"command" json:"command"`
	Environment   []EnvironmentVar `yaml:"environment" json:"environment"`
	Documentation Documentation    `yaml:"documentation" json:"documentation"`
	Compatibility Compatibility    `yaml:"compatibility" json:"compatibility"`
	Security      SecurityConfig   `yaml:"security" json:"security"`
	Verification  Verification     `yaml:"verification" json:"verification"`
	Links         Links            `yaml:"links" json:"links"`
}

// CommandSpec defines how to run the MCP server
type CommandSpec struct {
	Executable string   `yaml:"executable" json:"executable"`
	Args       []string `yaml:"args" json:"args"`
}

// EnvironmentVar defines a required or optional environment variable
type EnvironmentVar struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Required    bool   `yaml:"required" json:"required"`
	Sensitive   bool   `yaml:"sensitive" json:"sensitive"`
	Validation  string `yaml:"validation" json:"validation"`
	Default     string `yaml:"default" json:"default"`
	Example     string `yaml:"example" json:"example"`
}

// Documentation contains setup and usage information
type Documentation struct {
	Setup         string   `yaml:"setup" json:"setup"`
	Permissions   []string `yaml:"permissions" json:"permissions"`
	Examples      []string `yaml:"examples" json:"examples"`
	SecurityNotes string   `yaml:"security_notes" json:"security_notes"`
}

// Compatibility defines platform and version requirements
type Compatibility struct {
	Platforms      []string `yaml:"platforms" json:"platforms"`
	ClaudeVersions []string `yaml:"claude_versions" json:"claude_versions"`
	MinDDXVersion  string   `yaml:"min_ddx_version" json:"min_ddx_version"`
	NodeVersion    string   `yaml:"node_version" json:"node_version"`
}

// SecurityConfig defines security requirements and warnings
type SecurityConfig struct {
	Sandbox       string   `yaml:"sandbox" json:"sandbox"`
	NetworkAccess string   `yaml:"network_access" json:"network_access"`
	FileAccess    string   `yaml:"file_access" json:"file_access"`
	DataHandling  string   `yaml:"data_handling" json:"data_handling"`
	Warnings      []string `yaml:"warnings" json:"warnings"`
}

// Verification contains test commands for validating installation
type Verification struct {
	TestCommand      string `yaml:"test_command" json:"test_command"`
	ExpectedResponse string `yaml:"expected_response" json:"expected_response"`
}

// Links contains relevant URLs
type Links struct {
	Homepage      string `yaml:"homepage" json:"homepage"`
	Documentation string `yaml:"documentation" json:"documentation"`
	Issues        string `yaml:"issues" json:"issues"`
}

// Registry represents the MCP server registry
type Registry struct {
	Version    string              `yaml:"version" json:"version"`
	Updated    time.Time           `yaml:"updated" json:"updated"`
	Servers    []ServerReference   `yaml:"servers" json:"servers"`
	Categories map[string]Category `yaml:"categories" json:"categories"`

	// Private fields for caching
	cache    map[string]*Server
	cacheTTL time.Time

	// Config manager for checking installed servers
	config *ConfigManager
}

// ServerReference is a lightweight reference in the registry
type ServerReference struct {
	Name        string `yaml:"name" json:"name"`
	File        string `yaml:"file" json:"file"`
	Category    string `yaml:"category" json:"category"`
	Description string `yaml:"description" json:"description"`
}

// Category defines a server category
type Category struct {
	Description string `yaml:"description" json:"description"`
	Icon        string `yaml:"icon" json:"icon"`
}

// ClaudeType represents the type of Claude installation
type ClaudeType string

const (
	// ClaudeCode represents Claude Code installation
	ClaudeCode ClaudeType = "code"
	// ClaudeDesktop represents Claude Desktop installation
	ClaudeDesktop ClaudeType = "desktop"
)

// ClaudeInstallation represents a detected Claude installation
type ClaudeInstallation struct {
	Type       ClaudeType `json:"type"`
	ConfigPath string     `json:"config_path"`
	Version    string     `json:"version"`
	Detected   time.Time  `json:"detected"`
}

// ClaudeConfig represents the Claude configuration file
type ClaudeConfig struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
	// Preserve other fields when reading/writing
	Other map[string]interface{} `json:"-"`
}

// ServerConfig is the installed server configuration in Claude
type ServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

// InstallOptions contains options for server installation
type InstallOptions struct {
	ConfigPath  string            // Custom config file path
	Environment map[string]string // Environment variables
	NoBackup    bool              // Skip backup creation
	DryRun      bool              // Show what would be done
	Yes         bool              // Skip confirmations
	ClaudeType  ClaudeType        // Target Claude type
}

// ServerStatus represents the status of an installed server
type ServerStatus struct {
	Name        string            `json:"name"`
	Installed   bool              `json:"installed"`
	Configured  bool              `json:"configured"`
	Version     string            `json:"version"`
	ClaudeType  ClaudeType        `json:"claude_type"`
	ConfigPath  string            `json:"config_path"`
	Environment map[string]string `json:"environment"`
	Errors      []string          `json:"errors"`
}

// ListOptions contains options for listing servers
type ListOptions struct {
	Category  string // Filter by category
	Search    string // Search term
	Installed bool   // Show only installed
	Available bool   // Show only available
	Verbose   bool   // Detailed output
	Format    string // Output format (table/json/yaml)
}

// ConfigureOptions contains options for server configuration
type ConfigureOptions struct {
	Environment       map[string]string // Set environment variables
	AddEnvironment    map[string]string // Add environment variables
	RemoveEnvironment []string          // Remove environment variable keys
	Reset             bool              // Reset to defaults
}

// RemoveOptions contains options for server removal
type RemoveOptions struct {
	SkipConfirmation bool // Skip confirmation prompts
	NoBackup         bool // Skip backup creation
	Purge            bool // Remove all related data
}

// StatusOptions contains options for status checking
type StatusOptions struct {
	ServerName string // Specific server to check (empty for all)
	Check      bool   // Verify server connectivity
	Verbose    bool   // Show detailed information
	Format     string // Output format
}

// UpdateOptions contains options for registry updates
type UpdateOptions struct {
	Force  bool   // Force update even if current
	Server string // Update specific server
	Check  bool   // Check for updates only
}
