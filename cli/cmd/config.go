package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/easel/ddx/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configGlobal bool
	configLocal  bool
	configInit   bool
	configShow   bool
)

var configCmd = &cobra.Command{
	Use:   "config [get|set|validate] [key] [value]",
	Short: "Manage DDx configuration",
	Long: `Manage DDx configuration files.

Configuration is managed at two levels:
• Global: ~/.ddx.yml (affects all projects)
• Local: ./.ddx.yml (project-specific settings)

Local configuration takes precedence over global settings.

Subcommands:
• get <key>       - Get a specific configuration value
• set <key> <val> - Set a configuration value
• validate        - Validate configuration file`,
	RunE: runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&configGlobal, "global", "g", false, "Edit global configuration")
	configCmd.Flags().BoolVarP(&configLocal, "local", "l", false, "Edit local project configuration")
	configCmd.Flags().BoolVar(&configInit, "init", false, "Initialize configuration wizard")
	configCmd.Flags().BoolVar(&configShow, "show", false, "Show current configuration")
}

func runConfig(cmd *cobra.Command, args []string) error {
	// Handle flags first
	if configShow || len(args) == 0 && !configInit && !configGlobal && !configLocal {
		// Show config when no args or --show flag
		return showCurrentConfigWithWriter(cmd.OutOrStdout())
	}

	if configInit {
		return initConfigWizard()
	}

	if configGlobal && len(args) == 0 {
		return showGlobalConfigWithWriter(cmd.OutOrStdout())
	}

	if configLocal && len(args) == 0 {
		return editLocalConfig()
	}

	// Handle subcommands
	if len(args) > 0 {
		switch args[0] {
		case "get":
			if len(args) < 2 {
				return fmt.Errorf("get requires a key")
			}
			return getConfigValueWithWriter(args[1], cmd.OutOrStdout())
		case "set":
			if len(args) < 3 {
				return fmt.Errorf("set requires key and value")
			}
			return setConfigValueWithWriter(args[1], args[2], cmd.OutOrStdout())
		case "validate":
			return validateConfigWithWriter(cmd.OutOrStdout())
		case "export":
			return exportConfigWithWriter(cmd.OutOrStdout())
		case "import":
			if len(args) < 2 {
				return fmt.Errorf("import requires a file path")
			}
			return importConfigWithWriter(args[1], cmd.OutOrStdout())
		default:
			return fmt.Errorf("unknown subcommand: %s", args[0])
		}
	}

	// Should never reach here
	return showCurrentConfigWithWriter(cmd.OutOrStdout())
}

// showCurrentConfig displays the current effective configuration
func showCurrentConfig() error {
	return showCurrentConfigWithWriter(os.Stdout)
}

// showCurrentConfigWithWriter displays the current effective configuration to specified writer
func showCurrentConfigWithWriter(w io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		// If no config files exist, show defaults
		cfg = config.DefaultConfig
	}

	// Pretty print the configuration
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	fmt.Fprint(w, string(yamlData))
	return nil
}

// getConfigValue gets a specific configuration value
func getConfigValue(key string) error {
	return getConfigValueWithWriter(key, os.Stdout)
}

// getConfigValueWithWriter gets a specific configuration value and writes to specified writer
func getConfigValueWithWriter(key string, w io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Handle nested keys like "repository.url"
	value, err := getNestedValue(cfg, key)
	if err != nil {
		return fmt.Errorf("key not found: %s", key)
	}

	fmt.Fprintln(w, value)
	return nil
}

// setConfigValue sets a specific configuration value
func setConfigValue(key, value string) error {
	return setConfigValueWithWriter(key, value, os.Stdout)
}

// setConfigValueWithWriter sets a specific configuration value with specified writer
func setConfigValueWithWriter(key, value string, w io.Writer) error {
	// Load local config or create new one
	configPath := ".ddx.yml"
	var cfg *config.Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg = config.DefaultConfig
	} else {
		cfg, err = config.LoadLocal()
		if err != nil {
			return fmt.Errorf("failed to load local configuration: %w", err)
		}
	}

	// Handle nested keys like "variables.new_var"
	if err := setNestedValue(cfg, key, value); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	// Save the configuration
	if err := config.SaveLocal(cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Fprintf(w, "✅ Configuration updated: %s = %s\n", key, value)
	return nil
}

// validateConfig validates the configuration file
func validateConfig() error {
	return validateConfigWithWriter(os.Stdout)
}

// validateConfigWithWriter validates the configuration file and writes to specified writer
func validateConfigWithWriter(w io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(w, "❌ Configuration is invalid: %v\n", err)
		return err
	}

	// Basic validation
	if cfg.Version == "" {
		return fmt.Errorf("version is required")
	}

	if cfg.Repository.URL == "" {
		return fmt.Errorf("repository.url is required")
	}

	fmt.Fprintln(w, "✅ Configuration is valid")
	return nil
}

// showGlobalConfig shows only the global configuration
func showGlobalConfig() error {
	return showGlobalConfigWithWriter(os.Stdout)
}

// showGlobalConfigWithWriter shows only the global configuration to specified writer
func showGlobalConfigWithWriter(w io.Writer) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	globalConfigPath := filepath.Join(home, ".ddx.yml")
	if _, err := os.Stat(globalConfigPath); os.IsNotExist(err) {
		fmt.Fprintln(w, "No global configuration found")
		return nil
	}

	data, err := os.ReadFile(globalConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read global configuration: %w", err)
	}

	fmt.Fprint(w, string(data))
	return nil
}

// Helper functions for nested key handling
func getNestedValue(cfg *config.Config, key string) (string, error) {
	switch key {
	case "version":
		return cfg.Version, nil
	case "repository.url":
		return cfg.Repository.URL, nil
	case "repository.branch":
		return cfg.Repository.Branch, nil
	case "repository.path":
		return cfg.Repository.Path, nil
	default:
		// Check variables
		if strings.HasPrefix(key, "variables.") {
			varKey := strings.TrimPrefix(key, "variables.")
			if value, exists := cfg.Variables[varKey]; exists {
				return value, nil
			}
		}
		return "", fmt.Errorf("key not found")
	}
}

func setNestedValue(cfg *config.Config, key, value string) error {
	switch key {
	case "version":
		cfg.Version = value
	case "repository.url":
		cfg.Repository.URL = value
	case "repository.branch":
		cfg.Repository.Branch = value
	case "repository.path":
		cfg.Repository.Path = value
	default:
		// Check variables
		if strings.HasPrefix(key, "variables.") {
			varKey := strings.TrimPrefix(key, "variables.")
			if cfg.Variables == nil {
				cfg.Variables = make(map[string]string)
			}
			cfg.Variables[varKey] = value
		} else {
			return fmt.Errorf("unknown key: %s", key)
		}
	}
	return nil
}

// initConfigWizard runs an interactive configuration wizard
func initConfigWizard() error {
	cyan := color.New(color.FgCyan)

	cyan.Println("🧙 DDx Configuration Wizard")
	fmt.Println()

	// Start with default config
	cfg := *config.DefaultConfig

	// Ask configuration questions
	questions := []*survey.Question{
		{
			Name: "ai_model",
			Prompt: &survey.Select{
				Message: "Preferred AI model:",
				Options: []string{"claude-3-opus", "claude-3-sonnet", "gpt-4", "gpt-3.5-turbo"},
				Default: cfg.Variables["ai_model"],
			},
		},
		{
			Name: "includes",
			Prompt: &survey.MultiSelect{
				Message: "Select resources to include by default:",
				Options: []string{
					"prompts/claude",
					"prompts/general",
					"scripts/hooks",
					"scripts/setup",
					"templates/common",
					"patterns/error-handling",
					"patterns/testing",
					"configs/eslint",
					"configs/prettier",
				},
				Default: cfg.Includes,
			},
		},
	}

	answers := struct {
		AIModel  string   `survey:"ai_model"`
		Includes []string `survey:"includes"`
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}

	// Update configuration
	cfg.Variables["ai_model"] = answers.AIModel
	cfg.Includes = answers.Includes

	// Ask where to save
	var saveGlobal bool
	prompt := &survey.Confirm{
		Message: "Save as global configuration (affects all projects)?",
		Default: true,
	}
	if err := survey.AskOne(prompt, &saveGlobal); err != nil {
		return err
	}

	// Save configuration
	if saveGlobal {
		if err := config.Save(&cfg); err != nil {
			return fmt.Errorf("failed to save global configuration: %w", err)
		}
		fmt.Println("✅ Global configuration saved!")
	} else {
		if err := config.SaveLocal(&cfg); err != nil {
			return fmt.Errorf("failed to save local configuration: %w", err)
		}
		fmt.Println("✅ Local configuration saved!")
	}

	return nil
}

// editGlobalConfig opens the global config file for editing
func editGlobalConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := home + "/.ddx.yml"

	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := config.Save(config.DefaultConfig); err != nil {
			return fmt.Errorf("failed to create default configuration: %w", err)
		}
	}

	return openEditor(configPath)
}

// editLocalConfig opens the local config file for editing
func editLocalConfig() error {
	configPath := ".ddx.yml"

	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := *config.DefaultConfig
		// Set project-specific defaults
		cfg.Variables["project_name"] = "{{PROJECT_NAME}}"

		if err := config.SaveLocal(&cfg); err != nil {
			return fmt.Errorf("failed to create local configuration: %w", err)
		}
	}

	return openEditor(configPath)
}

// openEditor opens a file in the user's preferred editor
func openEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // fallback
	}

	fmt.Printf("Opening %s in %s...\n", filePath, editor)

	// Note: In a real implementation, you'd want to use exec.Command
	// to open the editor properly. This is simplified for the example.
	color.Yellow("Please edit the file manually: %s", filePath)

	return nil
}

// exportConfigWithWriter exports the current configuration to the specified writer
func exportConfigWithWriter(w io.Writer) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Export as YAML
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	fmt.Fprint(w, string(yamlData))
	return nil
}

// importConfigWithWriter imports configuration from a file
func importConfigWithWriter(filePath string, w io.Writer) error {
	// Read the import file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	// Parse the YAML
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Validate the imported config
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Save to local config
	if err := config.SaveLocal(&cfg); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Fprintf(w, "Configuration imported successfully from %s\n", filePath)
	return nil
}
