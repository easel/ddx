package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// State represents the current state of a workflow
type State struct {
	Workflow        string            `yaml:"workflow"`
	CurrentPhase    string            `yaml:"current_phase"`
	PhasesCompleted []string          `yaml:"phases_completed"`
	ActiveFeatures  map[string]string `yaml:"active_features,omitempty"`
	StartedAt       string            `yaml:"started_at"`
	LastUpdated     string            `yaml:"last_updated"`
	TasksCompleted  []string          `yaml:"tasks_completed,omitempty"`
	NextActions     []string          `yaml:"next_actions,omitempty"`
	PhaseProgress   map[string]int    `yaml:"phase_progress,omitempty"`
}

// LoadState loads the workflow state for a given workflow
func LoadState(workflowName string) (*State, error) {
	stateFile := fmt.Sprintf(".%s-state.yml", workflowName)

	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("workflow not initialized. Run 'ddx workflow init %s' first", workflowName)
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := yaml.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

// SaveState saves the workflow state
func SaveState(state *State) error {
	stateFile := fmt.Sprintf(".%s-state.yml", state.Workflow)

	state.LastUpdated = time.Now().Format("2006-01-02 15:04:05")

	data, err := yaml.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// InitializeState creates a new state for a workflow
func InitializeState(workflowName string, definition *WorkflowDefinition) (*State, error) {
	// Get the first phase
	var firstPhase *WorkflowPhase
	for _, phase := range definition.Phases {
		if phase.Order == 1 {
			firstPhase = &phase
			break
		}
	}

	if firstPhase == nil && len(definition.Phases) > 0 {
		firstPhase = &definition.Phases[0]
	}

	if firstPhase == nil {
		return nil, fmt.Errorf("workflow has no phases defined")
	}

	state := &State{
		Workflow:        workflowName,
		CurrentPhase:    firstPhase.ID,
		PhasesCompleted: []string{},
		ActiveFeatures:  make(map[string]string),
		StartedAt:       time.Now().Format("2006-01-02 15:04:05"),
		LastUpdated:     time.Now().Format("2006-01-02 15:04:05"),
		TasksCompleted:  []string{},
		NextActions:     getPhaseActions(firstPhase),
		PhaseProgress:   make(map[string]int),
	}

	return state, nil
}

// AdvancePhase moves the workflow to the next phase
func (s *State) AdvancePhase(definition *WorkflowDefinition) error {
	currentPhase := definition.GetPhaseByID(s.CurrentPhase)
	if currentPhase == nil {
		return fmt.Errorf("current phase '%s' not found in workflow definition", s.CurrentPhase)
	}

	nextPhase := definition.GetNextPhase(s.CurrentPhase)
	if nextPhase == nil {
		return fmt.Errorf("no phase after '%s' - workflow may be complete", s.CurrentPhase)
	}

	// Mark current phase as completed
	s.PhasesCompleted = append(s.PhasesCompleted, s.CurrentPhase)
	s.CurrentPhase = nextPhase.ID
	s.NextActions = getPhaseActions(nextPhase)
	s.LastUpdated = time.Now().Format("2006-01-02 15:04:05")

	return nil
}

// getPhaseActions returns suggested actions for a phase
func getPhaseActions(phase *WorkflowPhase) []string {
	actions := []string{}

	// Create actions from exit criteria
	for _, criteria := range phase.ExitCriteria {
		actions = append(actions, criteria)
	}

	// If no specific actions, create generic ones based on phase
	if len(actions) == 0 {
		basePath := filepath.Join("docs", "helix", fmt.Sprintf("%02d-%s", phase.Order, phase.ID))

		switch phase.ID {
		case "frame":
			actions = []string{
				fmt.Sprintf("Create problem statement (%s/problem.md)", basePath),
				fmt.Sprintf("Define user stories (%s/user-stories/)", basePath),
				fmt.Sprintf("Document requirements (%s/requirements.md)", basePath),
			}
		case "design":
			actions = []string{
				fmt.Sprintf("Create architecture overview (%s/architecture.md)", basePath),
				fmt.Sprintf("Define API contracts (%s/contracts/)", basePath),
				fmt.Sprintf("Document design decisions (%s/adr/)", basePath),
			}
		case "test":
			actions = []string{
				"Write failing contract tests",
				"Create integration test scenarios",
				"Define acceptance criteria",
			}
		case "build":
			actions = []string{
				"Implement functionality to pass tests",
				"Refactor code for quality",
				"Update documentation",
			}
		case "deploy":
			actions = []string{
				"Configure deployment pipeline",
				"Set up monitoring and alerts",
				"Deploy to production",
			}
		case "iterate":
			actions = []string{
				"Gather user feedback",
				"Analyze metrics and logs",
				"Plan next iteration",
			}
		default:
			actions = []string{
				fmt.Sprintf("Complete %s phase tasks", phase.Name),
			}
		}
	}

	return actions
}

// IsPhaseComplete checks if the given phase is in the completed list
func (s *State) IsPhaseComplete(phaseID string) bool {
	for _, completed := range s.PhasesCompleted {
		if completed == phaseID {
			return true
		}
	}
	return false
}

// GetProgress returns the overall workflow progress percentage
func (s *State) GetProgress(definition *WorkflowDefinition) int {
	if len(definition.Phases) == 0 {
		return 0
	}

	completedCount := len(s.PhasesCompleted)

	// Add partial credit for current phase if it has progress
	if progress, ok := s.PhaseProgress[s.CurrentPhase]; ok && progress > 0 {
		completedCount += progress / 100
	}

	return (completedCount * 100) / len(definition.Phases)
}
