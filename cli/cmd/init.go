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
		cmd.SilenceUsage = true
		return NewExitError(2, ".ddx.yml already exists. Use --force to overwrite.")
	}

	// Check if library path exists
	libPath, err := config.GetLibraryPath(libraryPath)
	libraryExists := true
	if err != nil || libPath == "" {
		libraryExists = false
	} else if _, err := os.Stat(libPath); os.IsNotExist(err) {
		libraryExists = false
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

	// Check if we're in the DDx repository itself
	if isDDxRepository() {
		// For DDx repo, point directly to the library directory
		localConfig.LibraryPath = "../library"
		fmt.Fprint(cmd.OutOrStdout(), "📚 Detected DDx repository - configuring library_path to use ../library\n")
	}

	// Try to load existing config for more accurate defaults
	if libraryExists {
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
		cmd.SilenceUsage = true
		return NewExitError(1, fmt.Sprintf("Failed to save configuration: %v", err))
	}

	// Copy resources if library exists
	if libraryExists {
		s := spinner.New(spinner.CharSets[14], 100)
		s.Prefix = "Setting up DDx... "
		s.Start()

		// Create local .ddx directory
		localDDxPath := ".ddx"
		if err := os.MkdirAll(localDDxPath, 0755); err != nil {
			s.Stop()
			cmd.SilenceUsage = true
			return NewExitError(1, fmt.Sprintf("Failed to create .ddx directory: %v", err))
		}

		// Copy selected resources
		for _, include := range localConfig.Includes {
			sourcePath := filepath.Join(libPath, include)
			targetPath := filepath.Join(localDDxPath, include)

			if _, err := os.Stat(sourcePath); err == nil {
				s.Suffix = fmt.Sprintf(" Copying %s...", include)
				if err := copyDir(sourcePath, targetPath); err != nil {
					s.Stop()
					cmd.SilenceUsage = true
					return NewExitError(1, fmt.Sprintf("Failed to copy %s: %v", include, err))
				}
			}
		}

		// Apply template if specified
		if initTemplate != "" {
			s.Suffix = fmt.Sprintf(" Applying template: %s...", initTemplate)
			if err := templates.Apply(initTemplate, ".", localConfig.Variables); err != nil {
				s.Stop()
				cmd.SilenceUsage = true
				return NewExitError(4, fmt.Sprintf("Template '%s' not found. Run 'ddx list templates' to see available templates.", initTemplate))
			}
		}

		s.Stop()
	}

	fmt.Fprint(cmd.OutOrStdout(), "✅ DDx initialized successfully!\n")
	fmt.Fprint(cmd.OutOrStdout(), "Initialized DDx in current project.\n")
	fmt.Fprintln(cmd.OutOrStdout())

	// Show next steps only if library exists
	if libraryExists {
		fmt.Fprint(cmd.OutOrStdout(), "Next steps:\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx list          - See available resources\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx apply <name>  - Apply templates or patterns\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx diagnose      - Analyze your project\n")
		fmt.Fprint(cmd.OutOrStdout(), "  ddx update        - Update toolkit\n")
		fmt.Fprintln(cmd.OutOrStdout())
	}

	return nil
}

// isDDxRepository checks if we're in the DDx repository
func isDDxRepository() bool {
	// Check for identifying files that indicate this is the DDx repo
	// Look for cli/main.go and library/ directory
	pwd, err := os.Getwd()
	if err != nil {
		return false
	}

	// Check if we're in the cli directory of DDx repo
	if filepath.Base(pwd) == "cli" {
		// Check for main.go
		if _, err := os.Stat("main.go"); err == nil {
			// Check for ../library directory
			if _, err := os.Stat("../library"); err == nil {
				return true
			}
		}
	}

	// Check if we're at the root of DDx repo
	if _, err := os.Stat("cli/main.go"); err == nil {
		if _, err := os.Stat("library"); err == nil {
			return true
		}
	}

	return false
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
