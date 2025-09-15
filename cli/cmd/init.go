package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/briandowns/spinner"
	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/templates"
	"github.com/spf13/cobra"
)

var (
	initTemplate string
	initForce    bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize DDx in current project",
	Long: `Initialize DDx in the current project directory.

This will:
• Create a local .ddx directory
• Copy selected resources from the master toolkit
• Create a local configuration file
• Optionally apply a project template`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initTemplate, "template", "t", "", "Use specific template")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Force initialization even if DDx already exists")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Fprint(cmd.OutOrStdout(), "🚀 Initializing DDx in current project...\n")
	fmt.Fprintln(cmd.OutOrStdout())

	// Check if config already exists
	configPath := ".ddx.yml"
	if _, err := os.Stat(configPath); err == nil && !initForce {
		// Config exists and --force not used - exit code 2 per contract
		fmt.Fprint(cmd.OutOrStdout(), "❌ DDx configuration already exists. Use --force to overwrite.\n")
		// Return error with exit code 2
		cmd.SilenceUsage = true
		return fmt.Errorf("exit code 2: configuration already exists")
	}

	// Check if DDx home exists
	ddxHome := getDDxHome()
	ddxHomeExists := true
	if _, err := os.Stat(ddxHome); os.IsNotExist(err) {
		ddxHomeExists = false
	}

	// Create local configuration even if DDx home doesn't exist
	pwd, _ := os.Getwd()
	projectName := filepath.Base(pwd)

	localConfig := &config.Config{
		Version: "1.0",
		Repository: config.Repository{
			URL:    "https://github.com/easel/ddx",
			Branch: "main",
			Path:   ".ddx/",
		},
		Includes: []string{
			"prompts/claude",
			"scripts/hooks",
			"templates/common",
		},
		Variables: map[string]string{
			"project_name": projectName,
			"ai_model":     "claude-3-opus",
		},
	}

	// Try to load existing config for more accurate defaults
	if ddxHomeExists {
		if cfg, err := config.Load(); err == nil {
			localConfig.Version = cfg.Version
			localConfig.Repository = cfg.Repository
			localConfig.Includes = cfg.Includes
			for k, v := range cfg.Variables {
				if k != "project_name" { // Keep project-specific name
					localConfig.Variables[k] = v
				}
			}
		}
	}

	// Save local configuration
	if err := config.SaveLocal(localConfig); err != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "❌ Failed to save configuration: %v\n", err)
		cmd.SilenceUsage = true
		return fmt.Errorf("exit code 1: failed to save configuration: %w", err)
	}

	// Copy resources if DDx home exists
	if ddxHomeExists {
		s := spinner.New(spinner.CharSets[14], 100)
		s.Prefix = "Setting up DDx... "
		s.Start()

		// Create local .ddx directory
		localDDxPath := ".ddx"
		if err := os.MkdirAll(localDDxPath, 0755); err != nil {
			s.Stop()
			fmt.Fprintf(cmd.OutOrStdout(), "❌ Failed to create .ddx directory: %v\n", err)
			cmd.SilenceUsage = true
			return fmt.Errorf("exit code 1: failed to create .ddx directory: %w", err)
		}

		// Copy selected resources
		for _, include := range localConfig.Includes {
			sourcePath := filepath.Join(ddxHome, include)
			targetPath := filepath.Join(localDDxPath, include)

			if _, err := os.Stat(sourcePath); err == nil {
				s.Suffix = fmt.Sprintf(" Copying %s...", include)
				if err := copyDir(sourcePath, targetPath); err != nil {
					s.Stop()
					fmt.Fprintf(cmd.OutOrStdout(), "❌ Failed to copy %s: %v\n", include, err)
					cmd.SilenceUsage = true
					return fmt.Errorf("exit code 1: failed to copy %s: %w", include, err)
				}
			}
		}

		// Apply template if specified
		if initTemplate != "" {
			s.Suffix = fmt.Sprintf(" Applying template: %s...", initTemplate)
			if err := templates.Apply(initTemplate, ".", localConfig.Variables); err != nil {
				s.Stop()
				fmt.Fprintf(cmd.OutOrStdout(), "❌ Failed to apply template: %v\n", err)
				cmd.SilenceUsage = true
				return fmt.Errorf("exit code 4: template not found: %w", err)
			}
		}

		s.Stop()
	}

	fmt.Fprint(cmd.OutOrStdout(), "✅ DDx initialized successfully!\n")
	fmt.Fprint(cmd.OutOrStdout(), "Initialized DDx in current project.\n")
	fmt.Fprintln(cmd.OutOrStdout())

	// Show next steps only if DDx home exists
	if ddxHomeExists {
		fmt.Fprint(cmd.OutOrStdout(), "Next steps:\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx list          - See available resources\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx apply <name>  - Apply templates or patterns\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx diagnose      - Analyze your project\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx update        - Update toolkit\n")
		fmt.Fprintln(cmd.OutOrStdout())
	}

	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create destination path
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	_, err = dstFile.ReadFrom(srcFile)
	return err
}
