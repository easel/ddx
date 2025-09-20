package workflow

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// WorkflowDefinition represents a workflow loaded from YAML
type WorkflowDefinition struct {
	Name        string             `yaml:"name"`
	Version     string             `yaml:"version"`
	Description string             `yaml:"description"`
	Author      string             `yaml:"author"`
	Created     string             `yaml:"created"`
	Tags        []string           `yaml:"tags"`
	Variables   []WorkflowVariable `yaml:"variables"`
	Phases      []WorkflowPhase    `yaml:"phases"`
}

// WorkflowVariable represents a configurable variable
type WorkflowVariable struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Prompt      string `yaml:"prompt"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default,omitempty"`
}

// WorkflowPhase represents a phase in the workflow
type WorkflowPhase struct {
	ID                string   `yaml:"id"`
	Order             int      `yaml:"order"`
	Name              string   `yaml:"name"`
	Description       string   `yaml:"description"`
	RequiredRole      string   `yaml:"required_role,omitempty"`
	ExitCriteria      []string `yaml:"exit_criteria"`
	EstimatedDuration string   `yaml:"estimated_duration,omitempty"`
}

// LoadWorkflow loads a workflow definition from the library
func LoadWorkflow(workflowName string) (*WorkflowDefinition, error) {
	// Check in .ddx directory first (for local overrides)
	localPath := filepath.Join(".ddx", "workflows", workflowName, "workflow.yml")
	if _, err := os.Stat(localPath); err == nil {
		return loadWorkflowFromFile(localPath)
	}

	// Then check in library
	libraryPath := filepath.Join("library", "workflows", workflowName, "workflow.yml")
	if _, err := os.Stat(libraryPath); err == nil {
		return loadWorkflowFromFile(libraryPath)
	}

	// If not found locally, check if we're in the DDx project itself
	// This handles the case where we're developing DDx
	projectLibPath := filepath.Join("..", "library", "workflows", workflowName, "workflow.yml")
	if _, err := os.Stat(projectLibPath); err == nil {
		return loadWorkflowFromFile(projectLibPath)
	}

	return nil, fmt.Errorf("workflow '%s' not found", workflowName)
}

// loadWorkflowFromFile loads a workflow definition from a specific file
func loadWorkflowFromFile(path string) (*WorkflowDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow WorkflowDefinition
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("failed to parse workflow YAML: %w", err)
	}

	return &workflow, nil
}

// ListAvailableWorkflows returns a list of available workflow names
func ListAvailableWorkflows() ([]string, error) {
	var workflows []string

	// Check local .ddx directory
	localDir := filepath.Join(".ddx", "workflows")
	if entries, err := os.ReadDir(localDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				workflows = append(workflows, entry.Name())
			}
		}
	}

	// Check library directory
	libraryDir := filepath.Join("library", "workflows")
	if entries, err := os.ReadDir(libraryDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				// Don't duplicate if already in local
				found := false
				for _, w := range workflows {
					if w == entry.Name() {
						found = true
						break
					}
				}
				if !found {
					workflows = append(workflows, entry.Name())
				}
			}
		}
	}

	// Check project library if we're in DDx development
	projectLibDir := filepath.Join("..", "library", "workflows")
	if entries, err := os.ReadDir(projectLibDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				// Don't duplicate if already found
				found := false
				for _, w := range workflows {
					if w == entry.Name() {
						found = true
						break
					}
				}
				if !found {
					workflows = append(workflows, entry.Name())
				}
			}
		}
	}

	if len(workflows) == 0 {
		return nil, fmt.Errorf("no workflows found")
	}

	return workflows, nil
}

// GetPhaseByID returns a phase by its ID
func (w *WorkflowDefinition) GetPhaseByID(id string) *WorkflowPhase {
	for _, phase := range w.Phases {
		if phase.ID == id {
			return &phase
		}
	}
	return nil
}

// GetNextPhase returns the next phase after the given phase ID
func (w *WorkflowDefinition) GetNextPhase(currentPhaseID string) *WorkflowPhase {
	currentPhase := w.GetPhaseByID(currentPhaseID)
	if currentPhase == nil {
		return nil
	}

	nextOrder := currentPhase.Order + 1
	for _, phase := range w.Phases {
		if phase.Order == nextOrder {
			return &phase
		}
	}
	return nil
}

// GetPhaseNames returns a list of phase IDs in order
func (w *WorkflowDefinition) GetPhaseNames() []string {
	names := make([]string, len(w.Phases))
	for i, phase := range w.Phases {
		names[i] = phase.ID
	}
	return names
}
