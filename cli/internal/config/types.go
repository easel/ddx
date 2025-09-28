package config

// NewConfig represents the simplified DDx configuration structure
// This aligns with the schema defined in ADR-005 and SD-003
type NewConfig struct {
	Version         string               `yaml:"version" json:"version"`
	LibraryBasePath string               `yaml:"library_base_path,omitempty" json:"library_base_path,omitempty"`
	Repository      *NewRepositoryConfig `yaml:"repository,omitempty" json:"repository,omitempty"`
	Variables       map[string]string    `yaml:"variables,omitempty" json:"variables,omitempty"`
	PersonaBindings map[string]string    `yaml:"persona_bindings,omitempty" json:"persona_bindings,omitempty"`
}

// NewRepositoryConfig represents repository settings for the new format
type NewRepositoryConfig struct {
	URL           string `yaml:"url,omitempty" json:"url,omitempty"`
	Branch        string `yaml:"branch,omitempty" json:"branch,omitempty"`
	SubtreePrefix string `yaml:"subtree_prefix,omitempty" json:"subtree_prefix,omitempty"`
}

// DefaultNewConfig returns a new config with default values applied
func DefaultNewConfig() *NewConfig {
	return &NewConfig{
		Version:         "1.0",
		LibraryBasePath: "./library",
		Repository: &NewRepositoryConfig{
			URL:           "https://github.com/easel/ddx",
			Branch:        "main",
			SubtreePrefix: "library",
		},
		Variables:       make(map[string]string),
		PersonaBindings: make(map[string]string),
	}
}

// DefaultConfig is an alias for DefaultNewConfig for compatibility
var DefaultConfig = DefaultNewConfig()

// ApplyDefaults ensures all required fields have default values
func (c *NewConfig) ApplyDefaults() {
	if c.Version == "" {
		c.Version = "1.0"
	}
	if c.LibraryBasePath == "" {
		c.LibraryBasePath = "./library"
	}
	if c.Repository == nil {
		c.Repository = &NewRepositoryConfig{
			URL:           "https://github.com/easel/ddx",
			Branch:        "main",
			SubtreePrefix: "library",
		}
	} else {
		if c.Repository.URL == "" {
			c.Repository.URL = "https://github.com/easel/ddx"
		}
		if c.Repository.Branch == "" {
			c.Repository.Branch = "main"
		}
		if c.Repository.SubtreePrefix == "" {
			c.Repository.SubtreePrefix = "library"
		}
	}
	if c.Variables == nil {
		c.Variables = make(map[string]string)
	}
}
