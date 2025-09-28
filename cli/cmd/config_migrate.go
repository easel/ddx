package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// configMigrateCmd represents the config migrate command
var configMigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate configuration from legacy .ddx.yml to new .ddx/config.yaml format",
	Long: `Migrate configuration from the legacy .ddx.yml format to the new .ddx/config.yaml format.

This command will:
1. Read your existing .ddx.yml file
2. Convert it to the new simplified format
3. Save it as .ddx/config.yaml
4. Create a backup of the original file

The new format is simpler and more focused on essential configuration options.

Examples:
  ddx config migrate                    # Migrate configuration
  ddx config migrate --dry-run          # Show what would be migrated without changing files
  ddx config migrate --force            # Overwrite existing .ddx/config.yaml if present`,
	RunE: runConfigMigrate,
}

var (
	migrateDryRun bool
	migrateForce  bool
)

func init() {
	// Note: Command registration is handled by command_factory.go
	// This init function sets up flags only
	configMigrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Show migration preview without making changes")
	configMigrateCmd.Flags().BoolVar(&migrateForce, "force", false, "Overwrite existing .ddx/config.yaml")
}

func runConfigMigrate(cmd *cobra.Command, args []string) error {
	// Migration no longer needed - only .ddx/config.yaml format is supported
	return fmt.Errorf("migration command is not needed - DDx now only supports .ddx/config.yaml format. Please use 'ddx init' to create a new configuration.")
}