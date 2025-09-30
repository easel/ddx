package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader loads workflow definitions from the library
type Loader struct {
	libraryPath string
}

// NewLoader creates a new workflow loader
func NewLoader(libraryPath string) *Loader {
	return &Loader{
		libraryPath: libraryPath,
	}
}

// Load reads and parses a workflow.yml file
func (l *Loader) Load(workflowName string) (*Definition, error) {
	// Construct path: {libraryPath}/workflows/{workflowName}/workflow.yml
	workflowPath := filepath.Join(l.libraryPath, "workflows", workflowName, "workflow.yml")

	// Check if file exists
	if _, err := os.Stat(workflowPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("workflow '%s' not found at %s", workflowName, workflowPath)
		}
		return nil, fmt.Errorf("failed to access workflow file: %w", err)
	}

	// Read file
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	// Parse YAML
	var def Definition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("failed to parse workflow.yml: %w", err)
	}

	// Validate definition
	if err := def.Validate(); err != nil {
		return nil, fmt.Errorf("invalid workflow definition: %w", err)
	}

	return &def, nil
}

// MatchesTriggers checks if text matches the triggers for a given agent command
func (l *Loader) MatchesTriggers(def *Definition, subcommand string, text string) bool {
	// Get agent command
	cmd, exists := def.GetAgentCommand(subcommand)
	if !exists {
		return false
	}

	// No triggers = never matches
	if cmd.Triggers == nil {
		return false
	}

	// Normalize text for matching (lowercase, trim)
	normalizedText := strings.ToLower(strings.TrimSpace(text))

	// Check keyword matches
	for _, keyword := range cmd.Triggers.Keywords {
		normalizedKeyword := strings.ToLower(keyword)

		// Match as whole word (with word boundaries)
		if matchesKeyword(normalizedText, normalizedKeyword) {
			return true
		}
	}

	// Check pattern matches (simple substring match)
	for _, pattern := range cmd.Triggers.Patterns {
		normalizedPattern := strings.ToLower(pattern)
		if strings.Contains(normalizedText, normalizedPattern) {
			return true
		}
	}

	return false
}
