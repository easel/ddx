package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/workflow"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// newAgentRequestCommand creates the agent request command
func newAgentRequestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request [message...]",
		Short: "Process agent request and route to appropriate workflow",
		Long: `Processes an incoming message from Claude and determines if it should
be handled by a workflow. Returns routing information or NO_HANDLER.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return handleAgentRequest(cmd, args)
		},
	}

	return cmd
}

func handleAgentRequest(cmd *cobra.Command, args []string) error {
	return handleAgentRequestWithDir(cmd, args, ".")
}

func handleAgentRequestWithDir(cmd *cobra.Command, args []string, workingDir string) error {
	// Join args into a single message
	message := strings.Join(args, " ")

	// Load config
	cfg, err := loadConfigFrom(workingDir)
	if err != nil || cfg == nil {
		// No config found - return NO_HANDLER
		fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
		return nil
	}

	// Apply defaults
	cfg.ApplyDefaults()

	// Check if there are active workflows
	if len(cfg.Workflows.Active) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
		return nil
	}

	// Check for safe word
	if len(args) > 0 && isSafeWord(args[0], cfg.Workflows.SafeWord) {
		// Remove safe word from message
		remainingMessage := strings.Join(args[1:], " ")
		fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
		fmt.Fprintf(cmd.OutOrStdout(), "SAFE_WORD: %s\n", cfg.Workflows.SafeWord)
		fmt.Fprintf(cmd.OutOrStdout(), "MESSAGE: %s\n", remainingMessage)
		return nil
	}

	// Get library path
	libraryPath := cfg.Library.Path
	if !filepath.IsAbs(libraryPath) {
		// Make relative to working directory
		libraryPath = filepath.Join(workingDir, libraryPath)
	}

	// Create workflow loader
	loader := workflow.NewLoader(libraryPath)

	// Filter out questions/discussions - messages starting with question words
	if isQuestion(message) {
		fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
		return nil
	}

	// First pass: check for pattern matches (more specific)
	for _, workflowName := range cfg.Workflows.Active {
		def, err := loader.Load(workflowName)
		if err != nil || !def.SupportsAgentCommand("request") {
			continue
		}

		if matchesPattern(def, "request", message) {
			return outputWorkflowMatch(cmd, workflowName, def, message)
		}
	}

	// Second pass: check for keyword matches (less specific)
	for _, workflowName := range cfg.Workflows.Active {
		def, err := loader.Load(workflowName)
		if err != nil || !def.SupportsAgentCommand("request") {
			continue
		}

		if matchesKeyword(def, "request", message) {
			return outputWorkflowMatch(cmd, workflowName, def, message)
		}
	}

	// No workflow matched
	fmt.Fprintln(cmd.OutOrStdout(), "NO_HANDLER")
	return nil
}

// isSafeWord checks if the first argument matches the safe word
func isSafeWord(arg, safeWord string) bool {
	// Strip trailing colon if present (e.g., "NODDX:" -> "NODDX")
	cleanArg := strings.TrimSuffix(arg, ":")
	return strings.EqualFold(cleanArg, safeWord)
}

// isQuestion checks if message starts with a question word
func isQuestion(message string) bool {
	message = strings.ToLower(strings.TrimSpace(message))
	questionWords := []string{
		"should", "would", "could", "can", "do", "does",
		"is", "are", "was", "were",
		"what", "when", "where", "why", "how", "which", "who",
	}
	for _, word := range questionWords {
		if strings.HasPrefix(message, word+" ") {
			return true
		}
	}
	return false
}

// matchesPattern checks if message matches any pattern in workflow triggers
func matchesPattern(def *workflow.Definition, subcommand, message string) bool {
	cmd, ok := def.GetAgentCommand(subcommand)
	if !ok || cmd.Triggers == nil {
		return false
	}

	normalized := strings.ToLower(strings.TrimSpace(message))
	for _, pattern := range cmd.Triggers.Patterns {
		normalizedPattern := strings.ToLower(pattern)
		if strings.Contains(normalized, normalizedPattern) {
			return true
		}
	}
	return false
}

// matchesKeyword checks if message matches any keyword in workflow triggers
func matchesKeyword(def *workflow.Definition, subcommand, message string) bool {
	cmd, ok := def.GetAgentCommand(subcommand)
	if !ok || cmd.Triggers == nil {
		return false
	}

	normalized := strings.ToLower(strings.TrimSpace(message))
	for _, keyword := range cmd.Triggers.Keywords {
		normalizedKeyword := strings.ToLower(keyword)
		// Match as whole word (with word boundaries)
		if matchWord(normalized, normalizedKeyword) {
			return true
		}
	}
	return false
}

// matchWord checks if keyword appears as a whole word in text
func matchWord(text, keyword string) bool {
	if text == keyword {
		return true
	}
	if strings.HasPrefix(text, keyword+" ") {
		return true
	}
	if strings.HasSuffix(text, " "+keyword) {
		return true
	}
	if strings.Contains(text, " "+keyword+" ") {
		return true
	}
	return false
}

// outputWorkflowMatch outputs the workflow routing information
func outputWorkflowMatch(cmd *cobra.Command, workflowName string, def *workflow.Definition, message string) error {
	agentCmd, ok := def.GetAgentCommand("request")
	if !ok {
		return fmt.Errorf("workflow %s missing request command", workflowName)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "WORKFLOW: %s\n", workflowName)
	fmt.Fprintf(cmd.OutOrStdout(), "SUBCOMMAND: request\n")
	fmt.Fprintf(cmd.OutOrStdout(), "ACTION: %s\n", agentCmd.Action)
	fmt.Fprintf(cmd.OutOrStdout(), "COMMAND: ddx workflow %s execute %s %s\n",
		workflowName, agentCmd.Action, quoteMessage(message))
	return nil
}

// quoteMessage quotes the message for shell command
func quoteMessage(message string) string {
	// Simple quoting - wrap in double quotes if contains special chars
	if strings.ContainsAny(message, " \t\n\"'\\$") {
		// Escape double quotes and backslashes
		escaped := strings.ReplaceAll(message, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
		return "\"" + escaped + "\""
	}
	return message
}

// loadConfig loads the DDx configuration file from current directory
func loadConfig() (*config.NewConfig, error) {
	return loadConfigFrom(".")
}

// loadConfigFrom loads the DDx configuration file from specified directory
func loadConfigFrom(workingDir string) (*config.NewConfig, error) {
	// Look for .ddx/config.yaml in specified directory
	configPath := filepath.Join(workingDir, ".ddx", "config.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	var cfg config.NewConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
