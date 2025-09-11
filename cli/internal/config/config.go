package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the DDx configuration
type Config struct {
	Version    string            `yaml:"version"`
	Repository Repository        `yaml:"repository"`
	Includes   []string          `yaml:"includes"`
	Overrides  map[string]string `yaml:"overrides,omitempty"`
	Variables  map[string]string `yaml:"variables"`
}

// Repository configuration
type Repository struct {
	URL    string `yaml:"url"`
	Branch string `yaml:"branch"`
	Path   string `yaml:"path"`
}

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
	},
}

// Load configuration from global and local files
func Load() (*Config, error) {
	config := *DefaultConfig // Copy default config

	// Load global config
	home, err := os.UserHomeDir()
	if err == nil {
		globalConfigPath := filepath.Join(home, ".ddx.yml")
		if err := loadConfigFile(globalConfigPath, &config); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Load local config (takes precedence)
	localConfigPath := ".ddx.yml"
	if err := loadConfigFile(localConfigPath, &config); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Set current directory as project name if not set
	if config.Variables["project_name"] == "" {
		pwd, _ := os.Getwd()
		config.Variables["project_name"] = filepath.Base(pwd)
	}

	return &config, nil
}

// LoadLocal loads only the local configuration
func LoadLocal() (*Config, error) {
	config := *DefaultConfig

	localConfigPath := ".ddx.yml"
	if err := loadConfigFile(localConfigPath, &config); err != nil {
		return nil, err
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

// loadConfigFile loads configuration from a YAML file
func loadConfigFile(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// saveConfigFile saves configuration to a YAML file
func saveConfigFile(path string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ReplaceVariables replaces template variables in content
func (c *Config) ReplaceVariables(content string) string {
	result := content

	// Replace {{variable}} patterns
	for key, value := range c.Variables {
		oldPattern := "{{" + key + "}}"
		result = replaceAll(result, oldPattern, value)

		// Also handle with spaces
		oldPattern = "{{ " + key + " }}"
		result = replaceAll(result, oldPattern, value)
	}

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
