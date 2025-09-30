package config

import (
	"fmt"
	"strings"
)

// NewConfig represents the simplified DDx configuration structure
// This aligns with the schema defined in ADR-005 and SD-003
type NewConfig struct {
	Version         string             `yaml:"version" json:"version"`
	Library         *LibraryConfig     `yaml:"library" json:"library"`
	Workflows       WorkflowsConfig    `yaml:"workflows,omitempty" json:"workflows,omitempty"`
	System          *SystemConfig      `yaml:"system,omitempty" json:"system,omitempty"`
	PersonaBindings map[string]string  `yaml:"persona_bindings,omitempty" json:"persona_bindings,omitempty"`
	UpdateCheck     *UpdateCheckConfig `yaml:"update_check,omitempty" json:"update_check,omitempty"`
}

// SystemConfig represents system-level configuration settings
type SystemConfig struct {
	MetaPrompt *string `yaml:"meta_prompt,omitempty" json:"meta_prompt,omitempty"`
}

// LibraryConfig represents library configuration settings
type LibraryConfig struct {
	Path       string            `yaml:"path,omitempty" json:"path,omitempty"`
	Repository *RepositoryConfig `yaml:"repository" json:"repository"`
}

// RepositoryConfig represents repository settings for the new format
type RepositoryConfig struct {
	URL    string `yaml:"url" json:"url"`
	Branch string `yaml:"branch" json:"branch"`
}

// UpdateCheckConfig represents update checking settings
type UpdateCheckConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Frequency string `yaml:"frequency"` // Duration: "24h", "12h", etc.
}

// WorkflowsConfig represents workflow activation and settings
type WorkflowsConfig struct {
	// Active workflows in priority order (first match wins)
	Active []string `yaml:"active,omitempty" json:"active,omitempty"`

	// SafeWord prefix to bypass workflow engagement
	// Default: "NODDX"
	SafeWord string `yaml:"safe_word,omitempty" json:"safe_word,omitempty"`
}

// ApplyDefaults sets default values for workflow configuration
func (w *WorkflowsConfig) ApplyDefaults() {
	if w.SafeWord == "" {
		w.SafeWord = "NODDX"
	}
	if w.Active == nil {
		w.Active = []string{}
	}
}

// Validate ensures workflow configuration is valid
func (w *WorkflowsConfig) Validate() error {
	// Validate safe word is not empty and has no spaces
	if strings.TrimSpace(w.SafeWord) == "" {
		return fmt.Errorf("safe_word cannot be empty")
	}
	if strings.Contains(w.SafeWord, " ") {
		return fmt.Errorf("safe_word cannot contain spaces")
	}

	// Validate no duplicate workflows
	seen := make(map[string]bool)
	for _, name := range w.Active {
		if seen[name] {
			return fmt.Errorf("duplicate workflow: %s", name)
		}
		seen[name] = true
	}

	return nil
}

// DefaultNewConfig returns a new config with default values applied
func DefaultNewConfig() *NewConfig {
	return &NewConfig{
		Version: "1.0",
		Library: &LibraryConfig{
			Path: ".ddx/library",
			Repository: &RepositoryConfig{
				URL:    "https://github.com/easel/ddx-library",
				Branch: "main",
			},
		},
		PersonaBindings: make(map[string]string),
		UpdateCheck: &UpdateCheckConfig{
			Enabled:   true,
			Frequency: "24h",
		},
	}
}

// DefaultConfig is an alias for DefaultNewConfig for compatibility
var DefaultConfig = DefaultNewConfig()

// GetMetaPrompt returns the meta-prompt path, defaulting to focused.md if unset
// Returns empty string if explicitly set to null/empty (disabled)
func (c *NewConfig) GetMetaPrompt() string {
	if c.System == nil || c.System.MetaPrompt == nil {
		// Unset: return default
		return "claude/system-prompts/focused.md"
	}
	// Explicitly set (could be empty string to disable)
	return *c.System.MetaPrompt
}

// ApplyDefaults ensures all required fields have default values
func (c *NewConfig) ApplyDefaults() {
	if c.Version == "" {
		c.Version = "1.0"
	}
	if c.Library == nil {
		c.Library = &LibraryConfig{
			Path: ".ddx/library",
			Repository: &RepositoryConfig{
				URL:    "https://github.com/easel/ddx-library",
				Branch: "main",
			},
		}
	} else {
		if c.Library.Path == "" {
			c.Library.Path = ".ddx/library"
		}
		if c.Library.Repository == nil {
			c.Library.Repository = &RepositoryConfig{
				URL:    "https://github.com/easel/ddx-library",
				Branch: "main",
			}
		} else {
			if c.Library.Repository.URL == "" {
				c.Library.Repository.URL = "https://github.com/easel/ddx-library"
			}
			if c.Library.Repository.Branch == "" {
				c.Library.Repository.Branch = "main"
			}
		}
	}
	if c.PersonaBindings == nil {
		c.PersonaBindings = make(map[string]string)
	}
	if c.UpdateCheck == nil {
		c.UpdateCheck = &UpdateCheckConfig{
			Enabled:   true,
			Frequency: "24h",
		}
	} else {
		if c.UpdateCheck.Frequency == "" {
			c.UpdateCheck.Frequency = "24h"
		}
	}
	c.Workflows.ApplyDefaults()
}
