package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the DDx configuration
type Config struct {
	Version         string                `yaml:"version"`
	LibraryPath     string                `yaml:"library_path,omitempty"`
	Repository      Repository            `yaml:"repository"`
	Repositories    map[string]Repository `yaml:"repositories,omitempty"`
	Includes        []string              `yaml:"includes"`
	Resources       *ResourceSelection    `yaml:"resources,omitempty"`
	Overrides       map[string]string     `yaml:"overrides,omitempty"`
	Variables       map[string]string     `yaml:"variables"`
	PersonaBindings map[string]string     `yaml:"persona_bindings,omitempty"`
}

// Repository configuration
type Repository struct {
	URL      string       `yaml:"url"`
	Branch   string       `yaml:"branch"`
	Path     string       `yaml:"path"`
	Remote   string       `yaml:"remote,omitempty"`
	Protocol string       `yaml:"protocol,omitempty"`
	Priority int          `yaml:"priority,omitempty"`
	Auth     *AuthConfig  `yaml:"auth,omitempty"`
	Proxy    *ProxyConfig `yaml:"proxy,omitempty"`
	Sync     *SyncConfig  `yaml:"sync,omitempty"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Method    string `yaml:"method,omitempty"`
	KeyPath   string `yaml:"key_path,omitempty"`
	Token     string `yaml:"token,omitempty"`
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	TokenFile string `yaml:"token_file,omitempty"`
}

// ProxyConfig represents proxy configuration
type ProxyConfig struct {
	URL      string `yaml:"url,omitempty"`
	Auth     string `yaml:"auth,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	NoProxy  string `yaml:"no_proxy,omitempty"`
}

// SyncConfig represents synchronization configuration
type SyncConfig struct {
	Frequency   string `yaml:"frequency,omitempty"`
	AutoUpdate  bool   `yaml:"auto_update,omitempty"`
	Timeout     int    `yaml:"timeout,omitempty"`
	RetryCount  int    `yaml:"retry_count,omitempty"`
	CheckSum    bool   `yaml:"checksum,omitempty"`
	ForceUpdate bool   `yaml:"force_update,omitempty"`
}

// ResourceSelection defines resource filtering configuration
type ResourceSelection struct {
	Prompts   *ResourceFilter `yaml:"prompts,omitempty"`
	Templates *ResourceFilter `yaml:"templates,omitempty"`
	Patterns  *ResourceFilter `yaml:"patterns,omitempty"`
	Configs   *ResourceFilter `yaml:"configs,omitempty"`
	Scripts   *ResourceFilter `yaml:"scripts,omitempty"`
	Workflows *ResourceFilter `yaml:"workflows,omitempty"`
}

// ResourceFilter defines include/exclude patterns for a resource type
type ResourceFilter struct {
	Include []string `yaml:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty"`
}

// ConfigError represents configuration-related errors
type ConfigError struct {
	Field      string
	Value      string
	Message    string
	Suggestion string
}

func (e *ConfigError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("config validation error in '%s': %s. Suggestion: %s", e.Field, e.Message, e.Suggestion)
	}
	return fmt.Sprintf("config validation error in '%s': %s", e.Field, e.Message)
}

// ValidationError represents multiple validation errors
type ValidationError struct {
	Errors []*ConfigError
}

func (e *ValidationError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	messages := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		messages[i] = err.Error()
	}
	return fmt.Sprintf("multiple validation errors:\n- %s", strings.Join(messages, "\n- "))
}

// SupportedVersions defines the versions that can be migrated
var SupportedVersions = []string{"1.0", "1.1", "2.0"}
var CurrentVersion = "2.0"

// Default configuration
var DefaultConfig = &Config{
	Version: "1.0",
	Repository: Repository{
		URL:    "https://github.com/easel/ddx",
		Branch: "main",
		Path:   ".ddx/",
	},
	Includes: []string{
		"prompts/claude",
		"scripts/hooks",
		"templates/common",
	},
	Overrides: make(map[string]string),
	Variables: map[string]string{
		"ai_model": "claude-3-opus",
		"author":   "", // Will be populated from git config or environment
		"email":    "", // Will be populated from git config or environment
	},
}

// Load configuration from global, local, and environment-specific files
func Load() (*Config, error) {
	// Start with default config
	config := *DefaultConfig

	// Load global config
	globalConfig, globalErr := LoadGlobal()
	if globalErr != nil && !os.IsNotExist(globalErr) {
		return nil, fmt.Errorf("failed to load global config: %w", globalErr)
	}

	// Load local config
	localConfig, localErr := LoadLocal()
	if localErr != nil && !os.IsNotExist(localErr) {
		return nil, fmt.Errorf("failed to load local config: %w", localErr)
	}

	// Load environment-specific config if DDX_ENV is set
	envConfig, envErr := LoadEnvironmentConfig()
	if envErr != nil && !os.IsNotExist(envErr) {
		return nil, fmt.Errorf("failed to load environment config: %w", envErr)
	}

	// Merge configurations with proper inheritance
	// Order: defaults < global < local < environment
	if globalConfig != nil {
		config = *DefaultConfig.Merge(globalConfig)
	}
	if localConfig != nil {
		config = *config.Merge(localConfig)
	}
	if envConfig != nil {
		config = *config.Merge(envConfig)
	}

	// Set current directory as project name if not set
	if config.Variables == nil {
		config.Variables = make(map[string]string)
	}
	if config.Variables["project_name"] == "" {
		pwd, _ := os.Getwd()
		config.Variables["project_name"] = filepath.Base(pwd)
	}

	// Validate the final configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &config, nil
}

// LoadGlobal loads only the global configuration with caching
func LoadGlobal() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	globalConfigPath := filepath.Join(home, ".ddx.yml")
	if _, err := os.Stat(globalConfigPath); os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}

	// Check cache first
	if cached := getCachedConfig(globalConfigPath, &globalConfigCache); cached != nil {
		return cached, nil
	}

	config := &Config{}
	if err := loadConfigFile(globalConfigPath, config); err != nil {
		return nil, err
	}

	// Cache the result
	cacheConfig(globalConfigPath, config, &globalConfigCache)

	return config, nil
}

// LoadLocal loads only the local configuration with caching
func LoadLocal() (*Config, error) {
	localConfigPath := ".ddx.yml"
	if _, err := os.Stat(localConfigPath); os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}

	// Check cache first
	if cached := getCachedConfig(localConfigPath, &localConfigCache); cached != nil {
		return cached, nil
	}

	config := &Config{}
	if err := loadConfigFile(localConfigPath, config); err != nil {
		return nil, err
	}

	// Cache the result
	cacheConfig(localConfigPath, config, &localConfigCache)

	return config, nil
}

// LoadEnvironmentConfig loads environment-specific configuration based on DDX_ENV
func LoadEnvironmentConfig() (*Config, error) {
	envName := os.Getenv("DDX_ENV")
	if envName == "" {
		// No environment specified, return not found
		return nil, os.ErrNotExist
	}

	// Build environment config file path: .ddx.{env}.yml
	envConfigPath := fmt.Sprintf(".ddx.%s.yml", envName)
	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		// Environment config file doesn't exist, return not found
		return nil, os.ErrNotExist
	}

	// Check cache first - use environment name in cache key
	cacheKey := fmt.Sprintf("env:%s:%s", envName, envConfigPath)
	if cached := getCachedConfig(cacheKey, &localConfigCache); cached != nil {
		return cached, nil
	}

	config := &Config{}
	if err := loadConfigFile(envConfigPath, config); err != nil {
		return nil, err
	}

	// Cache the result
	cacheConfig(cacheKey, config, &localConfigCache)

	return config, nil
}

// LoadFromFile loads configuration from a specific file path
func LoadFromFile(configPath string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Start with default config to ensure all required fields have defaults
	config := *DefaultConfig

	// Load the specified file
	if err := loadConfigFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to load configuration from %s: %w", configPath, err)
	}

	// Security: validate config before returning
	if err := validateConfigSecurity(&config); err != nil {
		return nil, fmt.Errorf("configuration security validation failed: %w", err)
	}

	return &config, nil
}

// Save configuration to global file
func Save(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".ddx.yml")
	return saveConfigFile(configPath, config)
}

// SaveLocal saves configuration to local file
func SaveLocal(config *Config) error {
	configPath := ".ddx.yml"
	return saveConfigFile(configPath, config)
}

// loadConfigFile loads configuration from a YAML file with security checks
func loadConfigFile(path string, config *Config) error {
	// Validate path to prevent path traversal
	if err := validateConfigPath(path); err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	// Check file size and permissions
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Security: limit config file size
	if info.Size() > maxConfigFileSize {
		return fmt.Errorf("config file too large: %d bytes (max %d)", info.Size(), maxConfigFileSize)
	}

	// Security: check file permissions
	mode := info.Mode()
	if mode.Perm()&0077 != 0 {
		// Config file is readable by group or others - potential security risk
		// Log warning but don't fail
	}

	// Read file with size limit
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	limitedReader := io.LimitReader(file, maxConfigFileSize)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return err
	}

	// Security: validate YAML content before parsing
	if err := validateYAMLContent(data); err != nil {
		return fmt.Errorf("config content validation failed: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config YAML: %w", err)
	}

	// Post-parse security validation
	return validateConfigSecurity(config)
}

// saveConfigFile saves configuration to a YAML file with security checks
func saveConfigFile(path string, config *Config) error {
	// Validate path to prevent path traversal
	if err := validateConfigPath(path); err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	// Security: validate config before saving
	if err := validateConfigSecurity(config); err != nil {
		return fmt.Errorf("config security validation failed: %w", err)
	}

	// Create a copy with sanitized sensitive data for marshaling
	sanitizedConfig := sanitizeConfigForSave(config)

	data, err := yaml.Marshal(sanitizedConfig)
	if err != nil {
		return err
	}

	// Security: validate serialized size
	if len(data) > maxConfigFileSize {
		return fmt.Errorf("serialized config too large: %d bytes (max %d)", len(data), maxConfigFileSize)
	}

	// Create directory if it doesn't exist with secure permissions
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	// Write file with secure permissions (readable only by owner)
	return os.WriteFile(path, data, 0600)
}

// WithRuntimeVariables creates a copy of the config with runtime variables merged
func (c *Config) WithRuntimeVariables(runtimeVars map[string]string) *Config {
	// Create a deep copy of the config
	copy := *c
	copy.Variables = make(map[string]string)

	// Copy existing variables
	for k, v := range c.Variables {
		copy.Variables[k] = v
	}

	// Override with runtime variables
	for k, v := range runtimeVars {
		copy.Variables[k] = v
	}

	return &copy
}

// ReplaceVariables replaces template variables in content with security checks
func (c *Config) ReplaceVariables(content string) string {
	result := content
	maxReplacements := 1000 // Prevent infinite loops
	replacementCount := 0

	// Replace {{variable}} patterns from config variables
	for key, value := range c.Variables {
		// Sanitize value but still replace it (we only mask in saved configs)
		sanitizedValue := sanitizeVariableValue(value)

		oldPattern := "{{" + key + "}}"
		newResult := replaceAll(result, oldPattern, sanitizedValue)
		if newResult != result {
			replacementCount++
		}
		result = newResult

		// Also handle with spaces
		oldPattern = "{{ " + key + " }}"
		newResult = replaceAll(result, oldPattern, sanitizedValue)
		if newResult != result {
			replacementCount++
		}
		result = newResult

		// Replace ${KEY} patterns (case-insensitive for config variables)
		upperKey := strings.ToUpper(key)
		oldPattern = "${" + upperKey + "}"
		newResult = replaceAll(result, oldPattern, sanitizedValue)
		if newResult != result {
			replacementCount++
		}
		result = newResult

		// Replace ${key} patterns (original case for config variables)
		oldPattern = "${" + key + "}"
		newResult = replaceAll(result, oldPattern, sanitizedValue)
		if newResult != result {
			replacementCount++
		}
		result = newResult

		// Prevent excessive replacements
		if replacementCount > maxReplacements {
			break
		}
	}

	// Process ${VAR} and ${VAR:-default} patterns for environment variables
	result = c.processEnvironmentVariables(result)

	return result
}

// Simple string replacement helper
func replaceAll(s, old, new string) string {
	// This is a simple implementation - in a real project you might use
	// strings.ReplaceAll or a more sophisticated template engine
	result := s
	for {
		newResult := ""
		found := false
		i := 0
		for i < len(result) {
			if i <= len(result)-len(old) && result[i:i+len(old)] == old {
				newResult += new
				i += len(old)
				found = true
			} else {
				newResult += string(result[i])
				i++
			}
		}
		result = newResult
		if !found {
			break
		}
	}
	return result
}

// processEnvironmentVariables processes ${VAR} and ${VAR:-default} patterns
func (c *Config) processEnvironmentVariables(content string) string {
	result := content
	i := 0

	for i < len(result) {
		// Look for ${
		if i < len(result)-1 && result[i] == '$' && result[i+1] == '{' {
			start := i
			i += 2 // Skip ${

			// Find the closing }
			braceEnd := -1
			for j := i; j < len(result); j++ {
				if result[j] == '}' {
					braceEnd = j
					break
				}
			}

			if braceEnd == -1 {
				// No closing brace found, skip this ${ and continue
				i++
				continue
			}

			// Extract the variable expression
			varExpr := result[i:braceEnd]

			// Check if it has a default value (contains :-)
			var varName, defaultValue string
			if strings.Contains(varExpr, ":-") {
				parts := strings.SplitN(varExpr, ":-", 2)
				varName = parts[0]
				defaultValue = parts[1]
			} else {
				varName = varExpr
				defaultValue = ""
			}

			// Get the value to substitute
			var replacement string

			// First check if it's a config variable (case-insensitive)
			if c.Variables != nil {
				// Check exact case first
				if value, exists := c.Variables[varName]; exists {
					replacement = sanitizeVariableValue(value)
				} else {
					// Check case-insensitive
					lowerVarName := strings.ToLower(varName)
					for key, value := range c.Variables {
						if strings.ToLower(key) == lowerVarName {
							replacement = sanitizeVariableValue(value)
							break
						}
					}
				}
			}

			// If not found in config variables, check environment
			if replacement == "" {
				if envValue := os.Getenv(varName); envValue != "" {
					replacement = sanitizeVariableValue(envValue)
				} else {
					// Use default value if provided
					replacement = defaultValue
				}
			}

			// Replace the entire ${...} expression
			wholeExpr := result[start : braceEnd+1]
			result = strings.ReplaceAll(result, wholeExpr, replacement)

			// Continue from where we left off, accounting for the replacement
			i = start + len(replacement)
		} else {
			i++
		}
	}

	return result
}

// Validate validates the configuration structure and values
func (c *Config) Validate() error {
	var errors []*ConfigError

	// Validate version
	if c.Version == "" {
		errors = append(errors, &ConfigError{
			Field:      "version",
			Message:    "version is required",
			Suggestion: "add 'version: \"1.0\"' to your config",
		})
	} else if !isValidVersion(c.Version) {
		errors = append(errors, &ConfigError{
			Field:      "version",
			Value:      c.Version,
			Message:    "unsupported version",
			Suggestion: fmt.Sprintf("use one of: %s", strings.Join(SupportedVersions, ", ")),
		})
	}

	// Validate repository
	if c.Repository.URL == "" {
		errors = append(errors, &ConfigError{
			Field:      "repository.url",
			Message:    "repository URL is required",
			Suggestion: "add a valid Git repository URL",
		})
	} else if !isValidURL(c.Repository.URL) {
		errors = append(errors, &ConfigError{
			Field:      "repository.url",
			Value:      c.Repository.URL,
			Message:    "invalid URL format",
			Suggestion: "use a valid URL like 'https://github.com/user/repo'",
		})
	}

	if c.Repository.Branch == "" {
		errors = append(errors, &ConfigError{
			Field:      "repository.branch",
			Message:    "repository branch is required",
			Suggestion: "add 'branch: \"main\"' or another valid branch name",
		})
	}

	if c.Repository.Path == "" {
		errors = append(errors, &ConfigError{
			Field:      "repository.path",
			Message:    "repository path is required",
			Suggestion: "add 'path: \".ddx/\"' or another valid path",
		})
	}

	// Validate repository protocol
	if c.Repository.Protocol != "" && c.Repository.Protocol != "ssh" && c.Repository.Protocol != "https" {
		errors = append(errors, &ConfigError{
			Field:      "repository.protocol",
			Value:      c.Repository.Protocol,
			Message:    "invalid protocol",
			Suggestion: "use 'ssh' or 'https'",
		})
	}

	// Validate authentication configuration
	if c.Repository.Auth != nil {
		if c.Repository.Auth.Method != "" {
			validAuthMethods := []string{"ssh-key", "token", "password", "oauth"}
			if !contains(validAuthMethods, c.Repository.Auth.Method) {
				errors = append(errors, &ConfigError{
					Field:      "repository.auth.method",
					Value:      c.Repository.Auth.Method,
					Message:    "invalid authentication method",
					Suggestion: fmt.Sprintf("use one of: %s", strings.Join(validAuthMethods, ", ")),
				})
			}
		}

		// Validate SSH key path if method is ssh-key
		if c.Repository.Auth.Method == "ssh-key" && c.Repository.Auth.KeyPath == "" {
			errors = append(errors, &ConfigError{
				Field:      "repository.auth.key_path",
				Message:    "SSH key path is required when using ssh-key authentication",
				Suggestion: "add 'key_path: \"~/.ssh/id_rsa\"' or another valid key path",
			})
		}

		// Validate token if method is token
		if c.Repository.Auth.Method == "token" && c.Repository.Auth.Token == "" && c.Repository.Auth.TokenFile == "" {
			errors = append(errors, &ConfigError{
				Field:      "repository.auth.token",
				Message:    "token or token_file is required when using token authentication",
				Suggestion: "add 'token: \"your-token\"' or 'token_file: \"/path/to/token\"'",
			})
		}
	}

	// Validate proxy configuration
	if c.Repository.Proxy != nil && c.Repository.Proxy.URL != "" {
		if !isValidURL(c.Repository.Proxy.URL) {
			errors = append(errors, &ConfigError{
				Field:      "repository.proxy.url",
				Value:      c.Repository.Proxy.URL,
				Message:    "invalid proxy URL format",
				Suggestion: "use a valid URL like 'http://proxy.company.com:8080'",
			})
		}
	}

	// Validate sync configuration
	if c.Repository.Sync != nil {
		if c.Repository.Sync.Frequency != "" {
			validFrequencies := []string{"never", "manual", "hourly", "daily", "weekly"}
			if !contains(validFrequencies, c.Repository.Sync.Frequency) {
				errors = append(errors, &ConfigError{
					Field:      "repository.sync.frequency",
					Value:      c.Repository.Sync.Frequency,
					Message:    "invalid sync frequency",
					Suggestion: fmt.Sprintf("use one of: %s", strings.Join(validFrequencies, ", ")),
				})
			}
		}

		if c.Repository.Sync.Timeout < 0 || c.Repository.Sync.Timeout > 3600 {
			errors = append(errors, &ConfigError{
				Field:      "repository.sync.timeout",
				Value:      fmt.Sprintf("%d", c.Repository.Sync.Timeout),
				Message:    "invalid sync timeout",
				Suggestion: "use a value between 0 and 3600 seconds",
			})
		}

		if c.Repository.Sync.RetryCount < 0 || c.Repository.Sync.RetryCount > 10 {
			errors = append(errors, &ConfigError{
				Field:      "repository.sync.retry_count",
				Value:      fmt.Sprintf("%d", c.Repository.Sync.RetryCount),
				Message:    "invalid retry count",
				Suggestion: "use a value between 0 and 10",
			})
		}
	}

	// Validate multiple repositories if configured
	for name, repo := range c.Repositories {
		if repo.URL == "" {
			errors = append(errors, &ConfigError{
				Field:      fmt.Sprintf("repositories.%s.url", name),
				Message:    "repository URL is required",
				Suggestion: "add a valid Git repository URL",
			})
		} else if !isValidURL(repo.URL) {
			errors = append(errors, &ConfigError{
				Field:      fmt.Sprintf("repositories.%s.url", name),
				Value:      repo.URL,
				Message:    "invalid URL format",
				Suggestion: "use a valid URL like 'https://github.com/user/repo'",
			})
		}

		if repo.Priority < 0 || repo.Priority > 100 {
			errors = append(errors, &ConfigError{
				Field:      fmt.Sprintf("repositories.%s.priority", name),
				Value:      fmt.Sprintf("%d", repo.Priority),
				Message:    "invalid priority",
				Suggestion: "use a value between 0 and 100",
			})
		}
	}

	// Validate variable names and values
	for key, value := range c.Variables {
		if !isValidVariableName(key) {
			errors = append(errors, &ConfigError{
				Field:      fmt.Sprintf("variables.%s", key),
				Value:      key,
				Message:    "invalid variable name",
				Suggestion: "use only letters, numbers, and underscores",
			})
		}

		// Validate variable value length
		if len(value) > maxVariableLength {
			errors = append(errors, &ConfigError{
				Field:      fmt.Sprintf("variables.%s", key),
				Value:      fmt.Sprintf("%d characters", len(value)),
				Message:    "variable value too long",
				Suggestion: fmt.Sprintf("limit variable values to %d characters", maxVariableLength),
			})
		}
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}

	return nil
}

// Merge combines this config with another, with the other taking precedence
func (c *Config) Merge(other *Config) *Config {
	result := &Config{
		Version:         c.Version,
		LibraryPath:     c.LibraryPath,
		Includes:        make([]string, len(c.Includes)),
		Overrides:       make(map[string]string),
		Variables:       make(map[string]string),
		PersonaBindings: make(map[string]string),
	}

	// Copy includes
	copy(result.Includes, c.Includes)

	// Copy repository (base config)
	result.Repository = c.Repository

	// Override with other's values if they exist
	if other.Version != "" {
		result.Version = other.Version
	}
	if other.LibraryPath != "" {
		result.LibraryPath = other.LibraryPath
	}
	if other.Repository.URL != "" {
		result.Repository.URL = other.Repository.URL
	}
	if other.Repository.Branch != "" {
		result.Repository.Branch = other.Repository.Branch
	}
	if other.Repository.Path != "" {
		result.Repository.Path = other.Repository.Path
	}
	if other.Repository.Remote != "" {
		result.Repository.Remote = other.Repository.Remote
	}
	if other.Repository.Protocol != "" {
		result.Repository.Protocol = other.Repository.Protocol
	}
	if other.Repository.Priority != 0 {
		result.Repository.Priority = other.Repository.Priority
	}
	if other.Repository.Auth != nil {
		result.Repository.Auth = other.Repository.Auth
	}
	if other.Repository.Proxy != nil {
		result.Repository.Proxy = other.Repository.Proxy
	}
	if other.Repository.Sync != nil {
		result.Repository.Sync = other.Repository.Sync
	}

	// Copy repositories map
	result.Repositories = make(map[string]Repository)
	for k, v := range c.Repositories {
		result.Repositories[k] = v
	}
	// Override with other's repositories
	for k, v := range other.Repositories {
		result.Repositories[k] = v
	}

	// Merge includes (append without duplicates)
	for _, include := range other.Includes {
		if !contains(result.Includes, include) {
			result.Includes = append(result.Includes, include)
		}
	}

	// Copy and merge variables
	for k, v := range c.Variables {
		result.Variables[k] = v
	}
	for k, v := range other.Variables {
		result.Variables[k] = v // Override
	}

	// Copy and merge overrides
	for k, v := range c.Overrides {
		result.Overrides[k] = v
	}
	for k, v := range other.Overrides {
		result.Overrides[k] = v // Override
	}

	// Copy and merge persona bindings
	for k, v := range c.PersonaBindings {
		result.PersonaBindings[k] = v
	}
	for k, v := range other.PersonaBindings {
		result.PersonaBindings[k] = v // Override
	}

	// Merge resources configuration
	result.Resources = mergeResourceSelections(c.Resources, other.Resources)

	return result
}

// GetNestedValue retrieves a value using dot notation (e.g., "repository.url")
func (c *Config) GetNestedValue(key string) (interface{}, error) {
	parts := strings.Split(key, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty key")
	}

	// Use reflection to navigate the structure
	v := reflect.ValueOf(c).Elem()

	for i, part := range parts {
		if !v.IsValid() {
			return nil, fmt.Errorf("invalid path at '%s'", strings.Join(parts[:i], "."))
		}

		if v.Kind() == reflect.Map {
			// Handle map access (for variables, overrides)
			mapValue := v.Interface().(map[string]string)
			if value, exists := mapValue[part]; exists {
				return value, nil
			}
			return nil, fmt.Errorf("key '%s' not found in map", part)
		}

		if v.Kind() == reflect.Struct {
			// Handle struct field access
			fieldName := capitalizeFieldName(part)
			field := v.FieldByName(fieldName)
			if !field.IsValid() {
				return nil, fmt.Errorf("field '%s' not found", part)
			}
			v = field
		} else {
			return nil, fmt.Errorf("cannot navigate into %s at '%s'", v.Kind(), part)
		}
	}

	return v.Interface(), nil
}

// SetNestedValue sets a value using dot notation
func (c *Config) SetNestedValue(key string, value interface{}) error {
	parts := strings.Split(key, ".")
	if len(parts) == 0 {
		return fmt.Errorf("empty key")
	}

	// Handle special cases for maps
	if len(parts) == 2 {
		switch parts[0] {
		case "variables":
			if c.Variables == nil {
				c.Variables = make(map[string]string)
			}
			c.Variables[parts[1]] = fmt.Sprintf("%v", value)
			return nil
		case "overrides":
			if c.Overrides == nil {
				c.Overrides = make(map[string]string)
			}
			c.Overrides[parts[1]] = fmt.Sprintf("%v", value)
			return nil
		}
	}

	// Use reflection for struct fields
	v := reflect.ValueOf(c).Elem()

	for i, part := range parts[:len(parts)-1] {
		if !v.IsValid() {
			return fmt.Errorf("invalid path at '%s'", strings.Join(parts[:i], "."))
		}

		if v.Kind() == reflect.Struct {
			fieldName := capitalizeFieldName(part)
			field := v.FieldByName(fieldName)
			if !field.IsValid() {
				return fmt.Errorf("field '%s' not found", part)
			}
			v = field
		} else {
			return fmt.Errorf("cannot navigate into %s at '%s'", v.Kind(), part)
		}
	}

	// Set the final value
	lastPart := parts[len(parts)-1]
	if v.Kind() == reflect.Struct {
		fieldName := capitalizeFieldName(lastPart)
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			return fmt.Errorf("field '%s' not found", lastPart)
		}
		if !field.CanSet() {
			return fmt.Errorf("field '%s' cannot be set", lastPart)
		}

		// Convert value to appropriate type
		val := reflect.ValueOf(value)
		if val.Type().ConvertibleTo(field.Type()) {
			field.Set(val.Convert(field.Type()))
		} else {
			// Try string conversion
			strValue := fmt.Sprintf("%v", value)
			field.SetString(strValue)
		}
		return nil
	}

	return fmt.Errorf("cannot set value at path '%s'", key)
}

// MigrateVersion migrates configuration to a newer version
func (c *Config) MigrateVersion(targetVersion string) (*Config, []string, error) {
	if !isValidVersion(targetVersion) {
		return nil, nil, fmt.Errorf("unsupported target version: %s", targetVersion)
	}

	if c.Version == targetVersion {
		// Even if no migration is needed, we should still apply common migrations
		migrated := &Config{
			Version:     c.Version,
			LibraryPath: c.LibraryPath,
			Repository:  c.Repository,
			Includes:    make([]string, len(c.Includes)),
			Overrides:   make(map[string]string),
			Variables:   make(map[string]string),
		}

		copy(migrated.Includes, c.Includes)

		// Copy maps
		for k, v := range c.Variables {
			migrated.Variables[k] = v
		}
		for k, v := range c.Overrides {
			migrated.Overrides[k] = v
		}

		// Apply common migrations even for same version
		var warnings []string
		if targetVersion == "2.0" {
			if _, exists := migrated.Variables["security_scan"]; !exists {
				migrated.Variables["security_scan"] = "true"
				warnings = append(warnings, "Added default security_scan setting")
			}
		}

		return migrated, warnings, nil
	}

	var warnings []string
	migrated := &Config{
		Version:     targetVersion,
		LibraryPath: c.LibraryPath,
		Repository:  c.Repository,
		Includes:    make([]string, len(c.Includes)),
		Overrides:   make(map[string]string),
		Variables:   make(map[string]string),
	}

	copy(migrated.Includes, c.Includes)

	// Copy maps
	for k, v := range c.Variables {
		migrated.Variables[k] = v
	}
	for k, v := range c.Overrides {
		migrated.Overrides[k] = v
	}

	// Version-specific migrations
	switch {
	case c.Version == "1.0" && (targetVersion == "1.1" || targetVersion == "2.0"):
		// Migrate from 1.0 to 1.1 or 2.0
		if migrated.Variables["ai_model"] == "claude-3-opus" {
			migrated.Variables["ai_model"] = "claude-3-5-sonnet"
			warnings = append(warnings, "Updated ai_model from claude-3-opus to claude-3-5-sonnet")
		}

	case c.Version == "1.1" && targetVersion == "2.0":
		// Migrate from 1.1 to 2.0
		// Add any 1.1 -> 2.0 specific migrations here
		warnings = append(warnings, "Migrated to version 2.0")
	}

	// Common migrations for all versions
	if targetVersion == "2.0" {
		// Add security configurations in v2.0
		if _, exists := migrated.Variables["security_scan"]; !exists {
			migrated.Variables["security_scan"] = "true"
			warnings = append(warnings, "Added default security_scan setting")
		}
	}

	return migrated, warnings, nil
}

// Helper functions

func isValidVersion(version string) bool {
	for _, v := range SupportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

func isValidURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}

	// Check for SSH Git URL format (e.g., git@github.com:user/repo.git)
	if strings.Contains(rawURL, "@") && strings.Contains(rawURL, ":") {
		parts := strings.Split(rawURL, "@")
		if len(parts) == 2 {
			hostAndPath := parts[1]
			if strings.Contains(hostAndPath, ":") {
				hostParts := strings.Split(hostAndPath, ":")
				if len(hostParts) >= 2 && hostParts[0] != "" && hostParts[1] != "" {
					return true // Valid SSH Git URL
				}
			}
		}
	}

	// Check for standard URL format (e.g., https://github.com/user/repo)
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func isValidVariableName(name string) bool {
	// Allow letters, numbers, underscores
	matched, _ := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]*$", name)
	return matched
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// capitalizeFieldName properly capitalizes field names for reflection
func capitalizeFieldName(name string) string {
	if name == "" {
		return ""
	}
	// Handle special cases
	switch strings.ToLower(name) {
	case "url":
		return "URL"
	default:
		// Standard title case
		return strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	}
}

// Security constants and variables
const (
	maxConfigFileSize = 1024 * 1024 // 1MB max config file size
	maxVariableCount  = 100         // Maximum number of variables
	maxVariableLength = 1024        // Maximum variable value length
	cacheExpiration   = 5 * time.Minute
)

var (
	// Cache for config validations
	configValidationCache = sync.Map{}

	// Sensitive variable patterns
	sensitiveVariablePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(password|secret|private_key|token|credential|auth)$`),
		regexp.MustCompile(`(?i)(api_key|apikey|access_token)$`),
	}

	// Dangerous content patterns in config
	dangerousConfigPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
		regexp.MustCompile(`(?i)(eval|exec|system|shell_exec)\s*\(`),
		regexp.MustCompile(`\$\{[^}]*\$\{`), // Nested variable substitution
	}
)

// Security validation functions

// validateConfigPath validates configuration file paths
func validateConfigPath(path string) error {
	if path == "" {
		return fmt.Errorf("config path cannot be empty")
	}

	if len(path) > 1024 {
		return fmt.Errorf("config path too long (max 1024 characters)")
	}

	// Clean path and check for traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal not allowed in config path")
	}

	// Validate file extension
	ext := filepath.Ext(cleanPath)
	allowedExtensions := []string{".yml", ".yaml", ".json"}
	validExt := false
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			validExt = true
			break
		}
	}
	if !validExt {
		return fmt.Errorf("invalid config file extension: %s (allowed: %s)", ext, strings.Join(allowedExtensions, ", "))
	}

	return nil
}

// validateYAMLContent validates YAML content for security issues
func validateYAMLContent(data []byte) error {
	content := string(data)

	// Check for dangerous patterns
	for _, pattern := range dangerousConfigPatterns {
		if pattern.MatchString(content) {
			return fmt.Errorf("potentially dangerous pattern detected in config")
		}
	}

	// Check for excessive nesting (YAML bomb protection)
	bracesCount := strings.Count(content, "{") + strings.Count(content, "[")
	if bracesCount > 1000 {
		return fmt.Errorf("config structure too complex (potential YAML bomb)")
	}

	// Check for excessively long lines (potential DoS)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if len(line) > 10000 {
			return fmt.Errorf("line %d too long (potential DoS): %d characters", i+1, len(line))
		}
	}

	return nil
}

// validateConfigSecurity performs security validation on parsed config
func validateConfigSecurity(config *Config) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Validate variable count
	if len(config.Variables) > maxVariableCount {
		return fmt.Errorf("too many variables (max %d)", maxVariableCount)
	}

	// Validate variable names and values
	for key, value := range config.Variables {
		if len(key) > 64 {
			return fmt.Errorf("variable name too long: %s (max 64 characters)", key)
		}
		if len(value) > maxVariableLength {
			return fmt.Errorf("variable value too long for %s (max %d characters)", key, maxVariableLength)
		}

		// Check for control characters
		if strings.ContainsAny(value, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x0e\x0f") {
			return fmt.Errorf("variable %s contains invalid control characters", key)
		}
	}

	// Validate includes paths
	for _, include := range config.Includes {
		if len(include) > 512 {
			return fmt.Errorf("include path too long: %s (max 512 characters)", include)
		}
		if strings.Contains(include, "..") {
			return fmt.Errorf("path traversal not allowed in includes: %s", include)
		}
	}

	return nil
}

// sanitizeConfigForSave creates a sanitized copy of config for saving
func sanitizeConfigForSave(config *Config) *Config {
	sanitized := *config
	sanitized.Variables = make(map[string]string)

	// Copy variables, masking sensitive ones
	for key, value := range config.Variables {
		if isSensitiveVariable(key) {
			// Mask sensitive variables in saved config
			sanitized.Variables[key] = "[REDACTED]"
		} else {
			sanitized.Variables[key] = sanitizeVariableValue(value)
		}
	}

	return &sanitized
}

// isSensitiveVariable checks if a variable name indicates sensitive data
func isSensitiveVariable(name string) bool {
	lowerName := strings.ToLower(name)
	for _, pattern := range sensitiveVariablePatterns {
		if pattern.MatchString(lowerName) {
			return true
		}
	}
	return false
}

// sanitizeVariableValue sanitizes variable values
func sanitizeVariableValue(value string) string {
	// Remove null bytes and control characters except newlines and tabs
	var result strings.Builder
	for _, r := range value {
		if r >= 32 || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Performance: Add caching for expensive operations

type configCacheEntry struct {
	config    *Config
	timestamp time.Time
	hash      string
}

var (
	globalConfigCache = sync.Map{}
	localConfigCache  = sync.Map{}
)

// getCachedConfig retrieves cached config if valid
func getCachedConfig(path string, cacheMap *sync.Map) *Config {
	if cached, exists := cacheMap.Load(path); exists {
		entry := cached.(configCacheEntry)
		if time.Since(entry.timestamp) < cacheExpiration {
			// Verify file hasn't changed
			if info, err := os.Stat(path); err == nil {
				currentHash := hashFileInfo(info)
				if currentHash == entry.hash {
					return entry.config
				}
			}
		}
		// Cache expired or file changed, remove it
		cacheMap.Delete(path)
	}
	return nil
}

// cacheConfig stores config in cache
func cacheConfig(path string, config *Config, cacheMap *sync.Map) {
	if info, err := os.Stat(path); err == nil {
		entry := configCacheEntry{
			config:    config,
			timestamp: time.Now(),
			hash:      hashFileInfo(info),
		}
		cacheMap.Store(path, entry)
	}
}

// hashFileInfo creates a hash of file info for cache validation
func hashFileInfo(info os.FileInfo) string {
	data := fmt.Sprintf("%s-%d-%s", info.Name(), info.Size(), info.ModTime().String())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// mergeResourceSelections merges two ResourceSelection configurations
func mergeResourceSelections(base, override *ResourceSelection) *ResourceSelection {
	if base == nil && override == nil {
		return nil
	}
	if base == nil {
		return copyResourceSelection(override)
	}
	if override == nil {
		return copyResourceSelection(base)
	}

	result := &ResourceSelection{
		Prompts:   mergeResourceFilter(base.Prompts, override.Prompts),
		Templates: mergeResourceFilter(base.Templates, override.Templates),
		Patterns:  mergeResourceFilter(base.Patterns, override.Patterns),
		Configs:   mergeResourceFilter(base.Configs, override.Configs),
		Scripts:   mergeResourceFilter(base.Scripts, override.Scripts),
		Workflows: mergeResourceFilter(base.Workflows, override.Workflows),
	}

	return result
}

// mergeResourceFilter merges two ResourceFilter configurations
func mergeResourceFilter(base, override *ResourceFilter) *ResourceFilter {
	if base == nil && override == nil {
		return nil
	}
	if base == nil {
		return copyResourceFilter(override)
	}
	if override == nil {
		return copyResourceFilter(base)
	}

	result := &ResourceFilter{}

	// Merge includes (override completely replaces base if present)
	if len(override.Include) > 0 {
		result.Include = make([]string, len(override.Include))
		copy(result.Include, override.Include)
	} else {
		result.Include = make([]string, len(base.Include))
		copy(result.Include, base.Include)
	}

	// Merge excludes (override completely replaces base if present)
	if len(override.Exclude) > 0 {
		result.Exclude = make([]string, len(override.Exclude))
		copy(result.Exclude, override.Exclude)
	} else {
		result.Exclude = make([]string, len(base.Exclude))
		copy(result.Exclude, base.Exclude)
	}

	return result
}

// copyResourceSelection creates a deep copy of ResourceSelection
func copyResourceSelection(rs *ResourceSelection) *ResourceSelection {
	if rs == nil {
		return nil
	}

	return &ResourceSelection{
		Prompts:   copyResourceFilter(rs.Prompts),
		Templates: copyResourceFilter(rs.Templates),
		Patterns:  copyResourceFilter(rs.Patterns),
		Configs:   copyResourceFilter(rs.Configs),
		Scripts:   copyResourceFilter(rs.Scripts),
		Workflows: copyResourceFilter(rs.Workflows),
	}
}

// copyResourceFilter creates a deep copy of ResourceFilter
func copyResourceFilter(rf *ResourceFilter) *ResourceFilter {
	if rf == nil {
		return nil
	}

	result := &ResourceFilter{}

	if rf.Include != nil {
		result.Include = make([]string, len(rf.Include))
		copy(result.Include, rf.Include)
	}

	if rf.Exclude != nil {
		result.Exclude = make([]string, len(rf.Exclude))
		copy(result.Exclude, rf.Exclude)
	}

	return result
}
