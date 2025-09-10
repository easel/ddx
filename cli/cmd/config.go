package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/easel/ddx/internal/config"
	"gopkg.in/yaml.v3"
)

var (
	configGlobal bool
	configLocal  bool
	configInit   bool
	configShow   bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage DDx configuration",
	Long: `Manage DDx configuration files.

Configuration is managed at two levels:
• Global: ~/.ddx.yml (affects all projects)
• Local: ./.ddx.yml (project-specific settings)

Local configuration takes precedence over global settings.`,
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
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	bold := color.New(color.Bold)

	if configShow {
		return showCurrentConfig()
	}

	if configInit {
		return initConfigWizard()
	}

	if configGlobal {
		return editGlobalConfig()
	}

	if configLocal {
		return editLocalConfig()
	}

	// Default behavior - show current config
	cyan.Println("🔧 DDx Configuration")
	fmt.Println()
	
	bold.Println("Current Configuration:")
	return showCurrentConfig()
}

// showCurrentConfig displays the current effective configuration
func showCurrentConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Pretty print the configuration
	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	fmt.Printf("%s\n", string(yamlData))
	
	// Show configuration file sources
	gray := color.New(color.FgHiBlack)
	fmt.Println()
	gray.Println("Configuration sources:")
	
	home, _ := os.UserHomeDir()
	globalPath := home + "/.ddx.yml"
	localPath := ".ddx.yml"
	
	if _, err := os.Stat(globalPath); err == nil {
		gray.Printf("  • Global: %s\n", globalPath)
	}
	
	if _, err := os.Stat(localPath); err == nil {
		gray.Printf("  • Local:  %s\n", localPath)
	}

	return nil
}

// initConfigWizard runs an interactive configuration wizard
func initConfigWizard() error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	
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
		green.Println("✅ Global configuration saved!")
	} else {
		if err := config.SaveLocal(&cfg); err != nil {
			return fmt.Errorf("failed to save local configuration: %w", err)
		}
		green.Println("✅ Local configuration saved!")
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