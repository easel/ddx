package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/easel/ddx/internal/config"
	"github.com/easel/ddx/internal/templates"
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
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	cyan.Println("🚀 Initializing DDx in current project...")
	fmt.Println()

	// Check if already initialized
	if isInitialized() && !initForce {
		var proceed bool
		prompt := &survey.Confirm{
			Message: "DDx is already initialized. Do you want to proceed anyway?",
			Default: false,
		}
		if err := survey.AskOne(prompt, &proceed); err != nil {
			return err
		}
		
		if !proceed {
			yellow.Println("Initialization cancelled.")
			return nil
		}
	}

	s := spinner.New(spinner.CharSets[14], 100)
	s.Prefix = "Setting up DDx... "
	s.Start()

	// Check if DDx home exists
	ddxHome := getDDxHome()
	if _, err := os.Stat(ddxHome); os.IsNotExist(err) {
		s.Stop()
		red.Println("❌ DDx toolkit not found. Please run the installation script first.")
		return fmt.Errorf("DDx not installed at %s", ddxHome)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		s.Stop()
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create local .ddx directory
	localDDxPath := ".ddx"
	if err := os.MkdirAll(localDDxPath, 0755); err != nil {
		s.Stop()
		return fmt.Errorf("failed to create .ddx directory: %w", err)
	}

	// Copy selected resources
	for _, include := range cfg.Includes {
		sourcePath := filepath.Join(ddxHome, include)
		targetPath := filepath.Join(localDDxPath, include)
		
		if _, err := os.Stat(sourcePath); err == nil {
			s.Suffix = fmt.Sprintf(" Copying %s...", include)
			if err := copyDir(sourcePath, targetPath); err != nil {
				s.Stop()
				return fmt.Errorf("failed to copy %s: %w", include, err)
			}
		}
	}

	// Create local configuration
	pwd, _ := os.Getwd()
	projectName := filepath.Base(pwd)
	
	localConfig := &config.Config{
		Version:  cfg.Version,
		Includes: cfg.Includes,
		Variables: map[string]string{
			"project_name": projectName,
			"ai_model":     cfg.Variables["ai_model"],
		},
	}

	if err := config.SaveLocal(localConfig); err != nil {
		s.Stop()
		return fmt.Errorf("failed to save local configuration: %w", err)
	}

	// Apply template if specified
	if initTemplate != "" {
		s.Suffix = fmt.Sprintf(" Applying template: %s...", initTemplate)
		if err := templates.Apply(initTemplate, ".", cfg.Variables); err != nil {
			s.Stop()
			return fmt.Errorf("failed to apply template: %w", err)
		}
	}

	s.Stop()
	green.Println("✅ DDx initialized successfully!")
	fmt.Println()
	
	// Show next steps
	fmt.Println(color.New(color.Bold).Sprint("Next steps:"))
	cyan.Printf("  ddx list          - See available resources\n")
	cyan.Printf("  ddx apply <name>  - Apply templates or patterns\n")
	cyan.Printf("  ddx diagnose      - Analyze your project\n")
	cyan.Printf("  ddx update        - Update toolkit\n")
	fmt.Println()

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