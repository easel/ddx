package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Type aliases for smooth transition
type Config = NewConfig
type Repository = NewRepositoryConfig

// ConfigError represents a single configuration error
type ConfigError struct {
	Field      string
	Value      string
	Message    string
	Suggestion string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationError represents multiple configuration errors
type ValidationError struct {
	Errors []*ConfigError
}

func (e *ValidationError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("%d configuration errors found", len(e.Errors))
}

// Load loads configuration using the new simplified approach
func Load() (*Config, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	return LoadWithWorkingDir(workingDir)
}

// LoadWithWorkingDir loads configuration from a specific working directory
func LoadWithWorkingDir(workingDir string) (*Config, error) {
	if workingDir == "" {
		var err error
		workingDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	// Use the new ConfigLoader to load from .ddx/config.yaml only
	loader, err := NewConfigLoaderWithWorkingDir(workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create config loader: %w", err)
	}

	config, err := loader.LoadConfig()
	if err != nil {
		// If no config file exists, return default config
		if os.IsNotExist(err) || strings.Contains(err.Error(), "no configuration file found") {
			config = DefaultNewConfig()
		} else {
			return nil, err
		}
	}

	// Ensure variables map exists and set project name
	if config.Variables == nil {
		config.Variables = make(map[string]string)
	}
	if config.Variables["project_name"] == "" {
		config.Variables["project_name"] = filepath.Base(workingDir)
	}

	// Apply defaults to ensure complete configuration
	config.ApplyDefaults()

	// Override library path with environment variable if set
	if envLibraryPath := os.Getenv("DDX_LIBRARY_BASE_PATH"); envLibraryPath != "" {
		config.LibraryBasePath = envLibraryPath
	}

	// Validate the final configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate validates the configuration structure and values (simplified)
func (c *Config) Validate() error {
	var errors []*ConfigError

	// Validate version
	if c.Version == "" {
		errors = append(errors, &ConfigError{
			Field:      "version",
			Message:    "version is required",
			Suggestion: "add 'version: \"1.0\"' to your config",
		})
	}

	// Validate repository
	if c.Repository != nil {
		if c.Repository.URL == "" {
			errors = append(errors, &ConfigError{
				Field:      "repository.url",
				Message:    "repository URL is required",
				Suggestion: "add a valid Git repository URL",
			})
		}

		if c.Repository.Branch == "" {
			errors = append(errors, &ConfigError{
				Field:      "repository.branch",
				Message:    "repository branch is required",
				Suggestion: "add 'branch: \"main\"' or another valid branch name",
			})
		}
	}

	// Validate variables
	if c.Variables != nil {
		for key, value := range c.Variables {
			if err := validateVariableKey(key); err != nil {
				errors = append(errors, &ConfigError{
					Field:      fmt.Sprintf("variables.%s", key),
					Value:      key,
					Message:    err.Error(),
					Suggestion: "use only alphanumeric characters, hyphens, and underscores",
				})
			}

			if err := validateVariableValue(value); err != nil {
				errors = append(errors, &ConfigError{
					Field:      fmt.Sprintf("variables.%s", key),
					Value:      value,
					Message:    err.Error(),
					Suggestion: "avoid control characters and extremely long values",
				})
			}
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
		LibraryBasePath: c.LibraryBasePath,
		Repository:      c.Repository,
		Variables:       make(map[string]string),
	}

	// Copy variables from base
	for k, v := range c.Variables {
		result.Variables[k] = v
	}

	// Override with other's values
	if other.Version != "" {
		result.Version = other.Version
	}
	if other.LibraryBasePath != "" {
		result.LibraryBasePath = other.LibraryBasePath
	}
	if other.Repository != nil {
		result.Repository = other.Repository
	}
	if other.Variables != nil {
		for k, v := range other.Variables {
			result.Variables[k] = v
		}
	}

	return result
}

// validateVariableKey validates a variable key
func validateVariableKey(key string) error {
	if key == "" {
		return fmt.Errorf("variable key cannot be empty")
	}
	if len(key) > 100 {
		return fmt.Errorf("variable key too long (max 100 characters)")
	}
	return nil
}

// validateVariableValue validates a variable value
func validateVariableValue(value string) error {
	if len(value) > 10000 {
		return fmt.Errorf("variable value too long (max 10000 characters)")
	}
	return nil
}

// ResolveLibraryResource resolves a library resource path
// NOTE: This function is now legacy. New code should load config and use cfg.LibraryBasePath directly.
func ResolveLibraryResource(resourcePath, configPath, workingDir string) (string, error) {
	// Load config to get the authoritative library path (which includes env var override)
	cfg, err := LoadWithWorkingDir(workingDir)
	if err == nil && cfg.LibraryBasePath != "" {
		// Check if it's an absolute path
		if filepath.IsAbs(resourcePath) {
			return resourcePath, nil
		}

		// Try relative to library path from config
		configPath := filepath.Join(cfg.LibraryBasePath, resourcePath)
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		// Return the config path even if it doesn't exist (for consistency)
		return configPath, nil
	}

	// Fall back to original logic if config loading fails
	if workingDir == "" {
		workingDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	// Check if it's an absolute path
	if filepath.IsAbs(resourcePath) {
		return resourcePath, nil
	}

	// Try relative to working directory first
	fullPath := filepath.Join(workingDir, resourcePath)
	if _, err := os.Stat(fullPath); err == nil {
		return fullPath, nil
	}

	// Try relative to library directory
	libraryPath := filepath.Join(workingDir, "library", resourcePath)
	if _, err := os.Stat(libraryPath); err == nil {
		return libraryPath, nil
	}

	// Return the original path even if it doesn't exist
	return fullPath, nil
}

// LoadFromFile loads configuration from a specific file path
func LoadFromFile(configPath string) (*Config, error) {
	// Use ConfigLoader to load the file
	loader, err := NewConfigLoaderWithWorkingDir(filepath.Dir(configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create config loader: %w", err)
	}
	return loader.LoadConfigFromPath(configPath)
}