package workflow

import (
	"fmt"
	"strings"
)

// Type aliases for backward compatibility with existing code
type WorkflowDefinition = Definition
type WorkflowPhase = Phase
type WorkflowVariable = Variable

// Definition represents a workflow definition from workflow.yml
type Definition struct {
	Name          string                  `yaml:"name"`
	Version       string                  `yaml:"version"`
	Description   string                  `yaml:"description"`
	Author        string                  `yaml:"author,omitempty"`
	Created       string                  `yaml:"created,omitempty"`
	Tags          []string                `yaml:"tags,omitempty"`
	Coordinator   string                  `yaml:"coordinator,omitempty"`
	AgentCommands map[string]AgentCommand `yaml:"agent_commands,omitempty"`
	Phases        []Phase                 `yaml:"phases"`
	Variables     []Variable              `yaml:"variables,omitempty"`
}

// AgentCommand defines a command that Claude can invoke
type AgentCommand struct {
	Enabled     bool      `yaml:"enabled"`
	Triggers    *Triggers `yaml:"triggers,omitempty"`
	Action      string    `yaml:"action"`
	Description string    `yaml:"description"`
}

// Triggers define patterns that activate a command
type Triggers struct {
	Keywords []string `yaml:"keywords,omitempty"`
	Patterns []string `yaml:"patterns,omitempty"`
}

// Phase represents a workflow phase
type Phase struct {
	ID                string   `yaml:"id"`
	Order             int      `yaml:"order"`
	Name              string   `yaml:"name"`
	Description       string   `yaml:"description"`
	RequiredRole      string   `yaml:"required_role,omitempty"`
	ExitCriteria      []string `yaml:"exit_criteria,omitempty"`
	EstimatedDuration string   `yaml:"estimated_duration,omitempty"`
}

// Variable represents a workflow variable
type Variable struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Prompt      string `yaml:"prompt,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
}

// Validate ensures the workflow definition is valid
func (d *Definition) Validate() error {
	if d.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	if d.Version == "" {
		return fmt.Errorf("workflow version is required")
	}

	// Validate agent commands
	for cmdName, cmd := range d.AgentCommands {
		if cmd.Enabled && cmd.Action == "" {
			return fmt.Errorf("agent command %s: action is required when enabled", cmdName)
		}
	}

	return nil
}

// SupportsAgentCommand checks if workflow supports a given agent subcommand
func (d *Definition) SupportsAgentCommand(subcommand string) bool {
	cmd, exists := d.AgentCommands[subcommand]
	return exists && cmd.Enabled
}

// GetAgentCommand returns the agent command definition if it exists and is enabled
func (d *Definition) GetAgentCommand(subcommand string) (*AgentCommand, bool) {
	cmd, exists := d.AgentCommands[subcommand]
	if !exists || !cmd.Enabled {
		return nil, false
	}
	return &cmd, true
}

// GetPhaseByID returns a phase by its ID
func (d *Definition) GetPhaseByID(id string) *Phase {
	for i := range d.Phases {
		if d.Phases[i].ID == id {
			return &d.Phases[i]
		}
	}
	return nil
}

// GetNextPhase returns the next phase after the given phase ID
func (d *Definition) GetNextPhase(currentPhaseID string) *Phase {
	currentPhase := d.GetPhaseByID(currentPhaseID)
	if currentPhase == nil {
		return nil
	}

	nextOrder := currentPhase.Order + 1
	for i := range d.Phases {
		if d.Phases[i].Order == nextOrder {
			return &d.Phases[i]
		}
	}
	return nil
}

// GetPhaseNames returns a list of phase IDs in order
func (d *Definition) GetPhaseNames() []string {
	names := make([]string, len(d.Phases))
	for i, phase := range d.Phases {
		names[i] = phase.ID
	}
	return names
}

// matchesKeyword checks if keyword appears as a whole word in text
func matchesKeyword(text, keyword string) bool {
	// Check if text is exactly the keyword
	if text == keyword {
		return true
	}

	// Check if keyword appears with word boundaries (space, start, end)
	// Start of text
	if strings.HasPrefix(text, keyword+" ") {
		return true
	}

	// End of text
	if strings.HasSuffix(text, " "+keyword) {
		return true
	}

	// Middle of text (surrounded by spaces)
	if strings.Contains(text, " "+keyword+" ") {
		return true
	}

	return false
}
