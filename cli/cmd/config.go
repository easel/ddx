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

// Command registration is now handled by command_factory.go
// This file only contains the run function implementation

// runConfig implements the config command logic
func runConfig(cmd *cobra.Command, args []string) error {
	// Get flag values locally
	showFlag, _ := cmd.Flags().GetBool("show")
	editFlag, _ := cmd.Flags().GetBool("edit")
	resetFlag, _ := cmd.Flags().GetBool("reset")
	wizardFlag, _ := cmd.Flags().GetBool("wizard")
	validateFlag, _ := cmd.Flags().GetBool("validate")
	globalFlag, _ := cmd.Flags().GetBool("global")

	// Handle flags
	if showFlag {
		return showConfig(cmd, globalFlag)
	}

	if editFlag {
		return editConfig(cmd, globalFlag)
	}

	if resetFlag {
		return resetConfig(cmd, globalFlag)
	}

	if wizardFlag {
		return initConfigWizard()
	}

	if validateFlag {
		return validateConfig(cmd)
	}

	// Handle subcommands
	if len(args) == 0 {
		// Default behavior: show config if no args or flags
		return showConfig(cmd, globalFlag)
	}

	subcommand := args[0]
	switch subcommand {
	case "get":
		if len(args) < 2 {
			return fmt.Errorf("key required for get command")
		}
		return getConfigValue(cmd, args[1], globalFlag)
	case "set":
		if len(args) < 3 {
			return fmt.Errorf("key and value required for set command")
		}
		return setConfigValue(cmd, args[1], args[2], globalFlag)
	case "validate":
		return validateConfig(cmd)
	case "export":
		return showConfig(cmd, globalFlag)
	case "import":
		// For now, just read from stdin
		return fmt.Errorf("import not yet implemented")
	default:
		return fmt.Errorf("unknown subcommand: %s", subcommand)
	}
}

func showConfig(cmd *cobra.Command, global bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Marshal config to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	fmt.Fprint(cmd.OutOrStdout(), string(data))
	return nil
}

func editConfig(cmd *cobra.Command, global bool) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	configPath := ".ddx.yml"
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(home, ".ddx.yml")
	}

	// Open editor
	fmt.Fprintf(cmd.OutOrStdout(), "Opening %s in %s...\n", configPath, editor)
	// In real implementation, would exec the editor
	return nil
}

func resetConfig(cmd *cobra.Command, global bool) error {
	configPath := ".ddx.yml"
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(home, ".ddx.yml")
	}

	// Create default config
	cfg := config.DefaultConfig
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal default configuration: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✅ Configuration reset to defaults: %s\n", configPath)
	return nil
}

func getConfigValue(cmd *cobra.Command, key string, global bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Handle nested keys
	if strings.HasPrefix(key, "variables.") {
		// Extract the variable name after "variables."
		varName := strings.TrimPrefix(key, "variables.")
		if val, ok := cfg.Variables[varName]; ok {
			fmt.Fprintln(cmd.OutOrStdout(), val)
		} else {
			return fmt.Errorf("key not found: %s", key)
		}
	} else {
		// Handle other known keys
		switch key {
		case "version":
			fmt.Fprintln(cmd.OutOrStdout(), cfg.Version)
		case "repository.url":
			fmt.Fprintln(cmd.OutOrStdout(), cfg.Repository.URL)
		case "repository.branch":
			fmt.Fprintln(cmd.OutOrStdout(), cfg.Repository.Branch)
		default:
			if val, ok := cfg.Variables[key]; ok {
				fmt.Fprintln(cmd.OutOrStdout(), val)
			} else {
				return fmt.Errorf("key not found: %s", key)
			}
		}
	}
	return nil
}

func setConfigValue(cmd *cobra.Command, key, value string, global bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Handle nested keys
	if strings.HasPrefix(key, "variables.") {
		// Extract the variable name after "variables."
		varName := strings.TrimPrefix(key, "variables.")
		if cfg.Variables == nil {
			cfg.Variables = make(map[string]string)
		}
		cfg.Variables[varName] = value
	} else {
		// Handle other known keys
		switch key {
		case "repository.url":
			cfg.Repository.URL = value
		case "repository.branch":
			cfg.Repository.Branch = value
		default:
			// If no prefix, assume it's a variable
			if cfg.Variables == nil {
				cfg.Variables = make(map[string]string)
			}
			cfg.Variables[key] = value
		}
	}

	// Save config
	configPath := ".ddx.yml"
	if global {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(home, ".ddx.yml")
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✅ Set %s = %s\n", key, value)
	return nil
}

func validateConfig(cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "✅ Configuration is valid")
	return nil
}

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
				Message: "Resources to include:",
				Options: []string{"prompts", "templates", "patterns", "configs", "scripts"},
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

	// Update config with answers
	cfg.Variables["ai_model"] = answers.AIModel
	cfg.Includes = answers.Includes

	// Marshal to YAML
	configYAML, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file
	if err := os.WriteFile(".ddx.yml", configYAML, 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	cyan.Println("✅ Configuration saved to .ddx.yml")
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}
