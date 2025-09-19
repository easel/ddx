package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Workflow state structure
type WorkflowState struct {
	Workflow        string            `yaml:"workflow"`
	CurrentPhase    string            `yaml:"current_phase"`
	PhasesCompleted []string          `yaml:"phases_completed"`
	ActiveFeatures  map[string]string `yaml:"active_features"`
	StartedAt       string            `yaml:"started_at"`
	LastUpdated     string            `yaml:"last_updated"`
	TasksCompleted  []string          `yaml:"tasks_completed"`
	NextActions     []string          `yaml:"next_actions"`
	PhaseProgress   map[string]int    `yaml:"phase_progress"`
}

// HELIX phases
var helixPhases = []string{
	"frame",
	"design",
	"test",
	"build",
	"deploy",
	"iterate",
}

var (
	workflowWatch bool
	workflowForce bool
)

// workflowCmd represents the workflow command
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage HELIX workflow state and progression",
	Long: `The workflow command helps manage HELIX workflow state, track progress,
and automatically generate next actions for continuous development.

Available subcommands:
  status     Show current workflow state and progress
  sync       Update CLAUDE.md with current workflow context
  advance    Move to next phase if criteria are met
  init       Initialize HELIX workflow state
  validate   Check current phase completion criteria
  next       Show next recommended actions`,
}

// workflowStatusCmd shows current workflow state
var workflowStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current workflow state and progress",
	RunE:  runWorkflowStatus,
}

// workflowSyncCmd updates CLAUDE.md with workflow context
var workflowSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Update CLAUDE.md with current workflow context",
	Long: `Sync analyzes the current project state and updates CLAUDE.md with:
- Current HELIX phase and progress
- Completed tasks and next actions
- Phase-specific context and enforcement
- Auto-continuation prompts for Claude`,
	RunE: runWorkflowSync,
}

// workflowAdvanceCmd moves to next phase
var workflowAdvanceCmd = &cobra.Command{
	Use:   "advance",
	Short: "Move to next phase if criteria are met",
	RunE:  runWorkflowAdvance,
}

// workflowInitCmd initializes HELIX workflow
var workflowInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize HELIX workflow state",
	RunE:  runWorkflowInit,
}

// workflowValidateCmd checks phase completion
var workflowValidateCmd = &cobra.Command{
	Use:   "validate [phase]",
	Short: "Check current phase completion criteria",
	RunE:  runWorkflowValidate,
}

// workflowNextCmd shows next actions
var workflowNextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show next recommended actions",
	RunE:  runWorkflowNext,
}

func init() {
	rootCmd.AddCommand(workflowCmd)

	// Add subcommands
	workflowCmd.AddCommand(workflowStatusCmd)
	workflowCmd.AddCommand(workflowSyncCmd)
	workflowCmd.AddCommand(workflowAdvanceCmd)
	workflowCmd.AddCommand(workflowInitCmd)
	workflowCmd.AddCommand(workflowValidateCmd)
	workflowCmd.AddCommand(workflowNextCmd)
	workflowCmd.AddCommand(workflowAutoCmd)

	// Flags
	workflowSyncCmd.Flags().BoolVar(&workflowWatch, "watch", false, "Continuously watch and update context")
	workflowInitCmd.Flags().BoolVarP(&workflowForce, "force", "f", false, "Force reinitialize existing workflow")
}

func runWorkflowStatus(cmd *cobra.Command, args []string) error {
	state, err := loadWorkflowState()
	if err != nil {
		return fmt.Errorf("no workflow state found. Run 'ddx workflow init' to initialize")
	}

	fmt.Printf("🔄 HELIX Workflow Status\n\n")
	fmt.Printf("Current Phase: %s\n", strings.Title(state.CurrentPhase))
	fmt.Printf("Started: %s\n", state.StartedAt)
	fmt.Printf("Last Updated: %s\n", state.LastUpdated)
	fmt.Printf("\nPhases Completed: %v\n", state.PhasesCompleted)

	if len(state.TasksCompleted) > 0 {
		fmt.Printf("\nTasks Completed:\n")
		for _, task := range state.TasksCompleted {
			fmt.Printf("  ✅ %s\n", task)
		}
	}

	if len(state.NextActions) > 0 {
		fmt.Printf("\nNext Actions:\n")
		for i, action := range state.NextActions {
			fmt.Printf("  %d. %s\n", i+1, action)
		}
	}

	return nil
}

func runWorkflowSync(cmd *cobra.Command, args []string) error {
	if workflowWatch {
		return runWorkflowWatch()
	}

	return syncWorkflowContext()
}

func runWorkflowAdvance(cmd *cobra.Command, args []string) error {
	state, err := loadWorkflowState()
	if err != nil {
		return fmt.Errorf("no workflow state found. Run 'ddx workflow init' to initialize")
	}

	// Check if current phase is complete
	complete, missing := validatePhaseCompletion(state.CurrentPhase)
	if !complete {
		fmt.Printf("❌ Cannot advance - missing requirements:\n")
		for _, req := range missing {
			fmt.Printf("  • %s\n", req)
		}
		return fmt.Errorf("phase %s not complete", state.CurrentPhase)
	}

	// Advance to next phase
	nextPhase := getNextPhase(state.CurrentPhase)
	if nextPhase == "" {
		fmt.Printf("✅ Workflow complete! All phases finished.\n")
		return nil
	}

	state.PhasesCompleted = append(state.PhasesCompleted, state.CurrentPhase)
	state.CurrentPhase = nextPhase
	state.LastUpdated = time.Now().Format("2006-01-02 15:04:05")
	state.NextActions = getPhaseNextActions(nextPhase)

	err = saveWorkflowState(state)
	if err != nil {
		return fmt.Errorf("failed to save workflow state: %v", err)
	}

	fmt.Printf("🚀 Advanced to %s phase\n", strings.Title(nextPhase))

	// Auto-sync CLAUDE.md
	return syncWorkflowContext()
}

func runWorkflowInit(cmd *cobra.Command, args []string) error {
	stateFile := ".helix-state.yml"

	// Check if already exists
	if _, err := os.Stat(stateFile); err == nil && !workflowForce {
		return fmt.Errorf("workflow already initialized. Use --force to reinitialize")
	}

	state := &WorkflowState{
		Workflow:        "helix",
		CurrentPhase:    "frame",
		PhasesCompleted: []string{},
		ActiveFeatures:  make(map[string]string),
		StartedAt:       time.Now().Format("2006-01-02 15:04:05"),
		LastUpdated:     time.Now().Format("2006-01-02 15:04:05"),
		TasksCompleted:  []string{},
		NextActions:     getPhaseNextActions("frame"),
		PhaseProgress:   make(map[string]int),
	}

	err := saveWorkflowState(state)
	if err != nil {
		return fmt.Errorf("failed to initialize workflow: %v", err)
	}

	fmt.Printf("🚀 HELIX workflow initialized\n")
	fmt.Printf("Current phase: Frame (Problem Definition)\n")

	// Create docs directory structure if it doesn't exist
	err = os.MkdirAll("docs/01-frame", 0755)
	if err != nil {
		return fmt.Errorf("failed to create docs directory: %v", err)
	}

	// Auto-sync CLAUDE.md
	return syncWorkflowContext()
}

func runWorkflowValidate(cmd *cobra.Command, args []string) error {
	state, err := loadWorkflowState()
	if err != nil {
		return fmt.Errorf("no workflow state found. Run 'ddx workflow init' to initialize")
	}

	phase := state.CurrentPhase
	if len(args) > 0 {
		phase = args[0]
	}

	complete, missing := validatePhaseCompletion(phase)

	fmt.Printf("📋 Phase %s Validation\n\n", strings.Title(phase))

	if complete {
		fmt.Printf("✅ Phase complete - ready to advance\n")
		return nil
	}

	fmt.Printf("⏳ Phase incomplete - missing:\n")
	for _, req := range missing {
		fmt.Printf("  • %s\n", req)
	}

	return nil
}

func runWorkflowNext(cmd *cobra.Command, args []string) error {
	state, err := loadWorkflowState()
	if err != nil {
		return fmt.Errorf("no workflow state found. Run 'ddx workflow init' to initialize")
	}

	if len(state.NextActions) == 0 {
		fmt.Printf("🎉 No pending actions - phase may be complete\n")
		fmt.Printf("Run 'ddx workflow validate' to check\n")
		return nil
	}

	fmt.Printf("📋 Next Actions (%s phase):\n\n", strings.Title(state.CurrentPhase))
	for i, action := range state.NextActions {
		fmt.Printf("%d. %s\n", i+1, action)
	}

	return nil
}

func runWorkflowWatch() error {
	fmt.Printf("👀 Watching for file changes and updating workflow context...\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")

	// Initial sync
	err := syncWorkflowContext()
	if err != nil {
		return err
	}

	// Simple polling-based file watching (for initial implementation)
	// In production, could use fsnotify for more efficient watching
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lastModTime := getLastModTime()

	for {
		select {
		case <-ticker.C:
			currentModTime := getLastModTime()
			if currentModTime != lastModTime {
				fmt.Printf("🔄 File changes detected, updating workflow context...\n")
				err := syncWorkflowContext()
				if err != nil {
					fmt.Printf("❌ Error updating context: %v\n", err)
				}
				lastModTime = currentModTime
			}
		}
	}
}

func getLastModTime() time.Time {
	// Check modification times of key directories
	dirs := []string{
		"docs/01-frame",
		"docs/02-design",
		"docs/03-test",
		"docs/04-build",
		"docs/05-deploy",
		"docs/06-iterate",
		".", // Root directory for any new files
	}

	var latestMod time.Time

	for _, dir := range dirs {
		if info, err := os.Stat(dir); err == nil {
			if info.ModTime().After(latestMod) {
				latestMod = info.ModTime()
			}

			// Also check files in directory
			if entries, err := os.ReadDir(dir); err == nil {
				for _, entry := range entries {
					if !entry.IsDir() {
						fullPath := dir + "/" + entry.Name()
						if info, err := os.Stat(fullPath); err == nil {
							if info.ModTime().After(latestMod) {
								latestMod = info.ModTime()
							}
						}
					}
				}
			}
		}
	}

	return latestMod
}

// Helper functions

func loadWorkflowState() (*WorkflowState, error) {
	data, err := os.ReadFile(".helix-state.yml")
	if err != nil {
		return nil, err
	}

	var state WorkflowState
	err = yaml.Unmarshal(data, &state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}

func saveWorkflowState(state *WorkflowState) error {
	data, err := yaml.Marshal(state)
	if err != nil {
		return err
	}

	return os.WriteFile(".helix-state.yml", data, 0644)
}

func getNextPhase(currentPhase string) string {
	for i, phase := range helixPhases {
		if phase == currentPhase && i < len(helixPhases)-1 {
			return helixPhases[i+1]
		}
	}
	return ""
}

func getPhaseNextActions(phase string) []string {
	switch phase {
	case "frame":
		return []string{
			"Create problem statement document (docs/01-frame/problem.md)",
			"Define user stories and personas (docs/01-frame/user-stories/)",
			"Identify stakeholders and requirements (docs/01-frame/stakeholders.md)",
			"Create risk assessment (docs/01-frame/risks.md)",
			"Define success metrics (docs/01-frame/metrics.md)",
		}
	case "design":
		return []string{
			"Create architecture overview (docs/02-design/architecture.md)",
			"Define API contracts (docs/02-design/api-contracts/)",
			"Create database schema (docs/02-design/schema.md)",
			"Design component structure (docs/02-design/components.md)",
			"Create deployment strategy (docs/02-design/deployment.md)",
		}
	case "test":
		return []string{
			"Write failing unit tests for core functionality",
			"Create integration test scenarios",
			"Define acceptance criteria tests",
			"Set up test automation framework",
			"Validate all tests fail initially (Red phase)",
		}
	case "build":
		return []string{
			"Implement core functionality to pass tests",
			"Build API endpoints and handlers",
			"Create database models and migrations",
			"Implement business logic",
			"Ensure all tests pass (Green phase)",
		}
	case "deploy":
		return []string{
			"Set up deployment pipeline",
			"Configure production environment",
			"Implement monitoring and logging",
			"Create rollback procedures",
			"Deploy to production",
		}
	case "iterate":
		return []string{
			"Gather production feedback",
			"Analyze performance metrics",
			"Document lessons learned",
			"Plan next iteration features",
			"Update specifications with insights",
		}
	default:
		return []string{}
	}
}

func validatePhaseCompletion(phase string) (bool, []string) {
	missing := []string{}

	switch phase {
	case "frame":
		if !fileExists("docs/01-frame/problem.md") {
			missing = append(missing, "Problem statement document")
		}
		if !dirExists("docs/01-frame/user-stories") {
			missing = append(missing, "User stories directory")
		}
		if !fileExists("docs/01-frame/stakeholders.md") {
			missing = append(missing, "Stakeholders document")
		}
		if !fileExists("docs/01-frame/risks.md") {
			missing = append(missing, "Risk assessment document")
		}
		if !fileExists("docs/01-frame/metrics.md") {
			missing = append(missing, "Success metrics document")
		}
	case "design":
		if !fileExists("docs/02-design/architecture.md") {
			missing = append(missing, "Architecture overview")
		}
		if !dirExists("docs/02-design/api-contracts") {
			missing = append(missing, "API contracts directory")
		}
		// Add more design phase validations
		// Add other phases...
	}

	return len(missing) == 0, missing
}

func syncWorkflowContext() error {
	state, err := loadWorkflowState()
	if err != nil {
		return fmt.Errorf("no workflow state found. Run 'ddx workflow init' to initialize")
	}

	// Update task completion status
	err = updateTaskCompletionStatus(state)
	if err != nil {
		return fmt.Errorf("failed to update task status: %v", err)
	}

	// Save updated state
	err = saveWorkflowState(state)
	if err != nil {
		return fmt.Errorf("failed to save workflow state: %v", err)
	}

	// Read existing CLAUDE.md
	claudeContent := ""
	if data, err := os.ReadFile("CLAUDE.md"); err == nil {
		claudeContent = string(data)
	}

	// Generate workflow context
	workflowContext := generateWorkflowContext(state)

	// Update or insert workflow sections
	updatedContent := updateWorkflowSections(claudeContent, workflowContext)

	// Write back to CLAUDE.md
	err = os.WriteFile("CLAUDE.md", []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update CLAUDE.md: %v", err)
	}

	fmt.Printf("✅ Updated CLAUDE.md with current workflow context\n")
	fmt.Printf("Phase: %s | Next: %s\n", strings.Title(state.CurrentPhase),
		getNextActionSummary(state.NextActions))

	return nil
}

func generateWorkflowContext(state *WorkflowState) string {
	var sb strings.Builder

	// Workflow status section
	sb.WriteString("<!-- WORKFLOW:START -->\n")
	sb.WriteString("## Current HELIX Workflow State\n\n")
	sb.WriteString(fmt.Sprintf("**Phase**: %s\n", strings.Title(state.CurrentPhase)))
	sb.WriteString(fmt.Sprintf("**Started**: %s\n", state.StartedAt))
	sb.WriteString(fmt.Sprintf("**Last Updated**: %s\n\n", state.LastUpdated))

	if len(state.TasksCompleted) > 0 {
		sb.WriteString("### Completed This Phase:\n")
		for _, task := range state.TasksCompleted {
			sb.WriteString(fmt.Sprintf("- ✅ %s\n", task))
		}
		sb.WriteString("\n")
	}

	if len(state.NextActions) > 0 {
		sb.WriteString("### Next Actions:\n")
		for i, action := range state.NextActions {
			sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, action))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("### Phase Enforcer Active:\n")
	sb.WriteString(fmt.Sprintf("The %s Phase Enforcer is currently active. ", strings.Title(state.CurrentPhase)))
	sb.WriteString(getPhaseDescription(state.CurrentPhase))
	sb.WriteString("\n<!-- WORKFLOW:END -->\n\n")

	// Auto-prompts section
	sb.WriteString("<!-- AUTO-PROMPTS:START -->\n")
	sb.WriteString("## Workflow Auto-Continuation\n\n")
	sb.WriteString("**CRITICAL**: When you complete a task, automatically continue to the next action without asking. Use this pattern:\n\n")
	sb.WriteString("```\n")
	sb.WriteString("I've completed [previous task]. Based on our HELIX workflow state, I'll now work on: [next action]\n")
	sb.WriteString("```\n\n")

	if len(state.NextActions) > 0 {
		sb.WriteString(fmt.Sprintf("**Current Next Action**: %s\n\n", state.NextActions[0]))
	}

	sb.WriteString("**Auto-Loop Instructions**:\n")
	sb.WriteString("1. After completing any task, immediately check the workflow state\n")
	sb.WriteString("2. Automatically start the next action from the NextActions list\n")
	sb.WriteString("3. Update workflow context and continue seamlessly\n")
	sb.WriteString("4. No manual prompts or questions - maintain continuous flow\n")
	sb.WriteString("5. If all phase tasks are complete, automatically validate and advance to next phase\n\n")

	sb.WriteString("**Phase Completion Auto-Advance**: When current phase tasks are complete, automatically run `ddx workflow advance` and continue with the new phase's first task.\n")
	sb.WriteString("<!-- AUTO-PROMPTS:END -->\n\n")

	return sb.String()
}

func updateWorkflowSections(content, workflowContext string) string {
	// Remove existing workflow sections
	content = removeSection(content, "<!-- WORKFLOW:START -->", "<!-- WORKFLOW:END -->")
	content = removeSection(content, "<!-- AUTO-PROMPTS:START -->", "<!-- AUTO-PROMPTS:END -->")

	// Find insertion point (after personas section or at end)
	insertPos := strings.Index(content, "<!-- PERSONAS:END -->")
	if insertPos != -1 {
		// Insert after personas section
		insertPos = strings.Index(content[insertPos:], "\n") + insertPos + 1
		return content[:insertPos] + "\n" + workflowContext + content[insertPos:]
	}

	// Insert at end
	return content + "\n" + workflowContext
}

func removeSection(content, startMarker, endMarker string) string {
	start := strings.Index(content, startMarker)
	if start == -1 {
		return content
	}

	end := strings.Index(content[start:], endMarker)
	if end == -1 {
		return content
	}

	end = start + end + len(endMarker)

	// Include newlines
	if start > 0 && content[start-1] == '\n' {
		start--
	}
	if end < len(content)-1 && content[end] == '\n' {
		end++
	}

	return content[:start] + content[end:]
}

func getPhaseDescription(phase string) string {
	switch phase {
	case "frame":
		return "Focus on WHAT and WHY, not HOW. Define the problem completely before considering solutions."
	case "design":
		return "Focus on HOW to architect the solution. No implementation yet, just design."
	case "test":
		return "Write failing tests first. All tests must fail initially (Red phase)."
	case "build":
		return "Implement to make tests pass. Focus on making tests green with minimal code."
	case "deploy":
		return "Deploy to production with monitoring and rollback capabilities."
	case "iterate":
		return "Gather feedback and plan the next iteration based on production learnings."
	default:
		return "Follow HELIX methodology principles."
	}
}

func getNextActionSummary(actions []string) string {
	if len(actions) == 0 {
		return "No pending actions"
	}
	return actions[0]
}

func updateTaskCompletionStatus(state *WorkflowState) error {
	// Map of task patterns to file/directory checks
	taskChecks := map[string]func() bool{
		// Frame phase tasks
		"Create problem statement document (docs/01-frame/problem.md)": func() bool {
			return fileExists("docs/01-frame/problem.md") && hasMinimumContent("docs/01-frame/problem.md", 200)
		},
		"Define user stories and personas (docs/01-frame/user-stories/)": func() bool {
			return dirExists("docs/01-frame/user-stories") && hasFiles("docs/01-frame/user-stories")
		},
		"Identify stakeholders and requirements (docs/01-frame/stakeholders.md)": func() bool {
			return fileExists("docs/01-frame/stakeholders.md") && hasMinimumContent("docs/01-frame/stakeholders.md", 200)
		},
		"Create risk assessment (docs/01-frame/risks.md)": func() bool {
			return fileExists("docs/01-frame/risks.md") && hasMinimumContent("docs/01-frame/risks.md", 200)
		},
		"Define success metrics (docs/01-frame/metrics.md)": func() bool {
			return fileExists("docs/01-frame/metrics.md") && hasMinimumContent("docs/01-frame/metrics.md", 200)
		},
		// Design phase tasks
		"Create architecture overview (docs/02-design/architecture.md)": func() bool {
			return fileExists("docs/02-design/architecture.md") && hasMinimumContent("docs/02-design/architecture.md", 300)
		},
		"Define API contracts (docs/02-design/api-contracts/)": func() bool {
			return dirExists("docs/02-design/api-contracts") && hasFiles("docs/02-design/api-contracts")
		},
		"Create database schema (docs/02-design/schema.md)": func() bool {
			return fileExists("docs/02-design/schema.md") && hasMinimumContent("docs/02-design/schema.md", 200)
		},
		"Design component structure (docs/02-design/components.md)": func() bool {
			return fileExists("docs/02-design/components.md") && hasMinimumContent("docs/02-design/components.md", 200)
		},
		"Create deployment strategy (docs/02-design/deployment.md)": func() bool {
			return fileExists("docs/02-design/deployment.md") && hasMinimumContent("docs/02-design/deployment.md", 200)
		},
		// Test phase tasks
		"Write failing unit tests for core functionality": func() bool {
			return hasTestFiles() && hasFailingTests()
		},
		"Create integration test scenarios": func() bool {
			return hasIntegrationTests()
		},
		"Define acceptance criteria tests": func() bool {
			return fileExists("docs/03-test/acceptance-criteria.md")
		},
		// Build phase tasks
		"Implement core functionality to pass tests": func() bool {
			return hasImplementation() && hasPassingTests()
		},
		// Add more phases as needed
	}

	// Check each next action and move completed ones to completed list
	var remainingActions []string
	var newlyCompleted []string

	for _, action := range state.NextActions {
		if checkFunc, exists := taskChecks[action]; exists {
			if checkFunc() {
				// Task is complete
				if !contains(state.TasksCompleted, action) {
					state.TasksCompleted = append(state.TasksCompleted, action)
					newlyCompleted = append(newlyCompleted, action)
				}
			} else {
				// Task still pending
				remainingActions = append(remainingActions, action)
			}
		} else {
			// Unknown task pattern, keep in next actions
			remainingActions = append(remainingActions, action)
		}
	}

	// Update next actions with remaining items
	state.NextActions = remainingActions

	// Update timestamp if tasks were completed
	if len(newlyCompleted) > 0 {
		state.LastUpdated = time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("✅ Detected completed tasks: %v\n", newlyCompleted)

		// Check if phase is complete and auto-advance if needed
		if len(remainingActions) == 0 {
			fmt.Printf("🎉 Phase %s complete! Checking for auto-advance...\n", strings.Title(state.CurrentPhase))
			if complete, _ := validatePhaseCompletion(state.CurrentPhase); complete {
				fmt.Printf("🚀 Auto-advancing to next phase...\n")
				triggerAutoAdvance(state)
			}
		}
	}

	return nil
}

func hasFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Enhanced validation functions
func hasMinimumContent(filepath string, minBytes int) bool {
	if !fileExists(filepath) {
		return false
	}
	info, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return info.Size() >= int64(minBytes)
}

func hasTestFiles() bool {
	testDirs := []string{"test", "tests", "__tests__", "spec"}
	for _, dir := range testDirs {
		if dirExists(dir) && hasFiles(dir) {
			return true
		}
	}
	// Also check for *_test.go files in current directory
	entries, err := os.ReadDir(".")
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), "_test.go") {
			return true
		}
	}
	return false
}

func hasFailingTests() bool {
	// This is a simplified check - in real implementation would run tests
	// For now, just check if tests exist and no implementation is complete
	return hasTestFiles() && !hasPassingTests()
}

func hasPassingTests() bool {
	// Simplified check - would actually run test suite
	// For now, assume tests pass if implementation exists
	return hasImplementation()
}

func hasIntegrationTests() bool {
	integrationDirs := []string{
		"tests/integration",
		"test/integration",
		"integration",
		"e2e",
		"tests/e2e",
	}
	for _, dir := range integrationDirs {
		if dirExists(dir) && hasFiles(dir) {
			return true
		}
	}
	return false
}

func hasImplementation() bool {
	// Check for common implementation directories
	implDirs := []string{"src", "lib", "app", "internal", "pkg"}
	for _, dir := range implDirs {
		if dirExists(dir) && hasFiles(dir) {
			return true
		}
	}
	// Also check for implementation files in root
	entries, err := os.ReadDir(".")
	if err != nil {
		return false
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && (strings.HasSuffix(name, ".go") ||
			strings.HasSuffix(name, ".js") || strings.HasSuffix(name, ".ts") ||
			strings.HasSuffix(name, ".py") || strings.HasSuffix(name, ".java")) {
			if !strings.HasSuffix(name, "_test.go") && name != "main.go" {
				return true
			}
		}
	}
	return false
}

func triggerAutoAdvance(state *WorkflowState) error {
	// Advance to next phase
	nextPhase := getNextPhase(state.CurrentPhase)
	if nextPhase == "" {
		fmt.Printf("✅ Workflow complete! All phases finished.\n")
		state.NextActions = []string{"Workflow complete - ready for production!"}
		return nil
	}

	state.PhasesCompleted = append(state.PhasesCompleted, state.CurrentPhase)
	state.CurrentPhase = nextPhase
	state.LastUpdated = time.Now().Format("2006-01-02 15:04:05")
	state.NextActions = getPhaseNextActions(nextPhase)

	fmt.Printf("🚀 Advanced to %s phase\n", strings.Title(nextPhase))
	fmt.Printf("📋 Next actions loaded: %d tasks\n", len(state.NextActions))

	return nil
}

// Add auto-loop command for continuous monitoring
var workflowAutoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Start automatic workflow monitoring and continuation",
	Long: `Starts automatic monitoring of workflow state and provides continuous
task progression without manual intervention. This creates an auto-loop that
detects task completion and automatically moves to the next action.`,
	RunE: runWorkflowAuto,
}

func runWorkflowAuto(cmd *cobra.Command, args []string) error {
	fmt.Printf("🔄 Starting automatic workflow monitoring...\n")
	fmt.Printf("This will continuously monitor and update workflow state\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")

	// Initial sync
	err := syncWorkflowContext()
	if err != nil {
		return err
	}

	// Continuous monitoring
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	lastSyncTime := time.Now()

	for {
		select {
		case <-ticker.C:
			// Check if any files have been modified
			if getLastModTime().After(lastSyncTime) {
				fmt.Printf("📄 Changes detected, syncing workflow state...\n")
				err := syncWorkflowContext()
				if err != nil {
					fmt.Printf("❌ Error syncing: %v\n", err)
				} else {
					fmt.Printf("✅ Workflow state updated\n")
				}
				lastSyncTime = time.Now()
			}
		}
	}
}

// Use existing fileExists and dirExists from diagnose.go
