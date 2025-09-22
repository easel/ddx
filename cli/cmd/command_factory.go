package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CommandFactory creates fresh command instances without global state
type CommandFactory struct {
	// Configuration options
	Version string
	Commit  string
	Date    string

	// Custom viper instance for isolation
	viperInstance *viper.Viper
}

// NewCommandFactory creates a new command factory with default settings
func NewCommandFactory() *CommandFactory {
	return &CommandFactory{
		Version:       Version,
		Commit:        Commit,
		Date:          Date,
		viperInstance: viper.New(),
	}
}

// NewCommandFactoryWithViper creates a factory with a custom viper instance
func NewCommandFactoryWithViper(v *viper.Viper) *CommandFactory {
	return &CommandFactory{
		Version:       Version,
		Commit:        Commit,
		Date:          Date,
		viperInstance: v,
	}
}

// NewRootCommand creates a fresh root command with all subcommands
func (f *CommandFactory) NewRootCommand() *cobra.Command {
	// Local flag variables scoped to this command instance
	var cfgFile string
	var verbose bool
	var libraryPath string

	// Create fresh root command
	rootCmd := &cobra.Command{
		Use:   "ddx",
		Short: "Document-Driven Development eXperience - AI development toolkit",
		Long: color.New(color.FgCyan).Sprint(banner) + `
DDx is a toolkit for AI-assisted development that helps you:

• Share templates, prompts, and patterns across projects
• Maintain consistent development practices
• Integrate AI tooling seamlessly
• Contribute improvements back to the community

Get started:
  ddx init          Initialize DDx in your project
  ddx list          See available resources
  ddx diagnose      Analyze your project setup

More information:
  Documentation: https://github.com/easel/ddx
  Issues & Support: https://github.com/easel/ddx/issues`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				fmt.Printf("DDx %s (commit: %s, built: %s)\n", f.Version, f.Commit, f.Date)
			}
		},
	}

	// Setup flags - these are now local to this command instance
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ddx.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&libraryPath, "library-base-path", "", "override path for DDx library location")

	// Store flag values in command context for access by subcommands
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Initialize config with the local viper instance
		f.initConfig(cfgFile, libraryPath)

		// Call the original PersistentPreRun if it exists
		if rootCmd.PersistentPreRun != nil {
			rootCmd.PersistentPreRun(cmd, args)
		}
		return nil
	}

	// Add all subcommands
	f.registerSubcommands(rootCmd)

	return rootCmd
}

// initConfig initializes configuration for this command instance
func (f *CommandFactory) initConfig(cfgFile, libPath string) {
	// Store library path override if provided
	if libPath != "" {
		os.Setenv("DDX_LIBRARY_BASE_PATH", libPath)
	}

	if cfgFile != "" {
		// Use config file from the flag
		f.viperInstance.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err == nil {
			// Search for config in home directory with name ".ddx" (without extension)
			f.viperInstance.AddConfigPath(home)
			f.viperInstance.AddConfigPath(".")
			f.viperInstance.SetConfigType("yaml")
			f.viperInstance.SetConfigName(".ddx")
		}
	}

	f.viperInstance.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := f.viperInstance.ReadInConfig(); err == nil {
		if verbose := f.viperInstance.GetBool("verbose"); verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", f.viperInstance.ConfigFileUsed())
		}
	}
}

// registerSubcommands adds all subcommands to the root command
func (f *CommandFactory) registerSubcommands(rootCmd *cobra.Command) {
	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			// Display version with proper formatting
			version := f.Version
			if version == "dev" {
				version = "v0.0.1-dev" // Make dev version semantic for tests
			} else if !strings.HasPrefix(version, "v") {
				version = "v" + version
			}

			fmt.Fprintf(cmd.OutOrStdout(), "DDx %s\n", version)
			fmt.Fprintf(cmd.OutOrStdout(), "Commit: %s\n", f.Commit)
			fmt.Fprintf(cmd.OutOrStdout(), "Built: %s\n", f.Date)

			// Check for --no-check flag
			noCheck, _ := cmd.Flags().GetBool("no-check")
			if !noCheck {
				// TODO: Implement update checking
				// For now, just document that this would check for updates
			}
		},
	}
	versionCmd.Flags().Bool("no-check", false, "Skip checking for updates")
	rootCmd.AddCommand(versionCmd)

	// Completion command
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To configure your shell to load completions:

Bash:
  echo 'source <(ddx completion bash)' >> ~/.bashrc

Zsh:
  echo 'source <(ddx completion zsh)' >> ~/.zshrc

Fish:
  ddx completion fish | source

PowerShell:
  ddx completion powershell | Out-String | Invoke-Expression
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}
	rootCmd.AddCommand(completionCmd)

	// Register all other commands
	rootCmd.AddCommand(f.newInitCommand())
	rootCmd.AddCommand(f.newListCommand())
	rootCmd.AddCommand(f.newDiagnoseCommand())
	rootCmd.AddCommand(f.newUpdateCommand())
	rootCmd.AddCommand(f.newContributeCommand())
	rootCmd.AddCommand(f.newApplyCommand())
	rootCmd.AddCommand(f.newConfigCommand())
	rootCmd.AddCommand(f.newWorkflowCommand())
	rootCmd.AddCommand(f.newPersonaCommand())
	rootCmd.AddCommand(f.newMCPCommand())
	rootCmd.AddCommand(f.newUninstallCommand())

	// Add prompts command group
	promptsCmd := &cobra.Command{
		Use:     "prompts",
		Short:   "Manage AI prompts",
		Aliases: []string{"prompt"},
	}
	promptsCmd.AddCommand(f.newPromptsListCommand())
	promptsCmd.AddCommand(f.newPromptsShowCommand())
	rootCmd.AddCommand(promptsCmd)
}

// Helper function to get library path from environment or flag
func getLibraryPathFromEnv() string {
	return os.Getenv("DDX_LIBRARY_BASE_PATH")
}
