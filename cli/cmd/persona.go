package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/easel/ddx/internal/persona"
	"github.com/spf13/cobra"
)

// personaCmd represents the persona command
var personaCmd = &cobra.Command{
	Use:   "persona",
	Short: "Manage AI personas for consistent role-based interactions",
	Long: `The persona command manages AI personalities that can be bound to specific roles
and loaded into CLAUDE.md for consistent, high-quality AI interactions.

Personas are defined as markdown files with YAML frontmatter containing:
- name: Unique identifier for the persona
- roles: List of roles this persona can fulfill
- description: Brief description of the persona's approach
- tags: Keywords for discovery and categorization

Examples:
  ddx persona list                    # List all available personas
  ddx persona show strict-reviewer    # Show details of a specific persona
  ddx persona bind code-reviewer strict-reviewer  # Bind persona to role
  ddx persona load                    # Load all bound personas into CLAUDE.md
  ddx persona status                  # Show currently loaded personas`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(personaCmd)

	// Add subcommands
	personaCmd.AddCommand(personaListCmd)
	personaCmd.AddCommand(personaShowCmd)
	personaCmd.AddCommand(personaBindCmd)
	personaCmd.AddCommand(personaUnbindCmd)
	personaCmd.AddCommand(personaBindingsCmd)
	personaCmd.AddCommand(personaLoadCmd)
	personaCmd.AddCommand(personaUnloadCmd)
	personaCmd.AddCommand(personaStatusCmd)
}

// personaListCmd lists available personas
var personaListCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List available personas",
	Long: `List all available personas with optional filtering by role or tags.

Examples:
  ddx persona list                    # List all personas
  ddx persona list --role code-reviewer  # List personas for specific role
  ddx persona list --tag security    # List personas with specific tag`,
	RunE: func(cmd *cobra.Command, args []string) error {
		role, _ := cmd.Flags().GetString("role")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		loader := persona.NewPersonaLoader()

		var personas []*persona.Persona
		var err error

		if role != "" {
			personas, err = loader.FindByRole(role)
		} else if len(tags) > 0 {
			personas, err = loader.FindByTags(tags)
		} else {
			personas, err = loader.ListPersonas()
		}

		if err != nil {
			return fmt.Errorf("failed to list personas: %w", err)
		}

		if len(personas) == 0 {
			if role != "" {
				cmd.Printf("No personas found for role '%s'\n", role)
			} else if len(tags) > 0 {
				cmd.Printf("No personas found with tags: %s\n", strings.Join(tags, ", "))
			} else {
				cmd.Println("No personas found")
			}
			return nil
		}

		// Display personas in table format
		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "Available Personas:")
		fmt.Fprintln(out)
		w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tROLES\tDESCRIPTION\tTAGS")
		fmt.Fprintln(w, "----\t-----\t-----------\t----")

		for _, p := range personas {
			rolesStr := strings.Join(p.Roles, ", ")
			tagsStr := strings.Join(p.Tags, ", ")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Name, rolesStr, p.Description, tagsStr)
		}

		return w.Flush()
	},
}

// personaShowCmd shows details of a specific persona
var personaShowCmd = &cobra.Command{
	Use:   "show <persona-name>",
	Short: "Show detailed information about a persona",
	Long: `Show the complete details of a specific persona including its content.

Examples:
  ddx persona show strict-reviewer    # Show details of strict-reviewer persona`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		personaName := args[0]
		loader := persona.NewPersonaLoader()

		p, err := loader.LoadPersona(personaName)
		if err != nil {
			// Check if it's a persona not found error
			if strings.Contains(err.Error(), "not found") {
				return NewExitError(ExitCodePersonaNotFound, fmt.Sprintf("persona '%s' not found", personaName))
			}
			return fmt.Errorf("failed to load persona '%s': %w", personaName, err)
		}

		// Display persona details
		cmd.Printf("Name: %s\n", p.Name)
		cmd.Printf("Description: %s\n", p.Description)
		cmd.Printf("Roles: %s\n", strings.Join(p.Roles, ", "))
		if len(p.Tags) > 0 {
			cmd.Printf("Tags: %s\n", strings.Join(p.Tags, ", "))
		}
		cmd.Println("\nContent:")
		cmd.Println(strings.Repeat("-", 40))
		cmd.Println(p.Content)

		return nil
	},
}

// personaBindCmd binds a persona to a role
var personaBindCmd = &cobra.Command{
	Use:   "bind <role> <persona-name>",
	Short: "Bind a persona to a role",
	Long: `Bind a specific persona to a role in the project configuration.
This creates a mapping in .ddx.yml that will be used when loading personas.

Examples:
  ddx persona bind code-reviewer strict-reviewer
  ddx persona bind test-engineer tdd-specialist`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		role := args[0]
		personaName := args[1]

		// Verify persona exists
		loader := persona.NewPersonaLoader()
		p, err := loader.LoadPersona(personaName)
		if err != nil {
			return fmt.Errorf("persona '%s' not found: %w", personaName, err)
		}

		// Check if persona can fulfill the role
		canFulfill := false
		for _, personaRole := range p.Roles {
			if personaRole == role {
				canFulfill = true
				break
			}
		}

		if !canFulfill {
			return fmt.Errorf("persona '%s' cannot fulfill role '%s'. Available roles: %s",
				personaName, role, strings.Join(p.Roles, ", "))
		}

		// Create binding
		manager := persona.NewBindingManager()
		if err := manager.SetBinding(role, personaName); err != nil {
			// Check if it's a missing config error
			if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "No such file") {
				return NewExitError(ExitCodeNoConfig, "No .ddx.yml configuration found")
			}
			return fmt.Errorf("failed to create binding: %w", err)
		}

		cmd.Printf("✓ Bound role '%s' to persona '%s'\n", role, personaName)
		return nil
	},
}

// personaUnbindCmd removes a persona binding
var personaUnbindCmd = &cobra.Command{
	Use:   "unbind <role>",
	Short: "Remove the binding for a role",
	Long: `Remove the persona binding for a specific role from the project configuration.

Examples:
  ddx persona unbind code-reviewer
  ddx persona unbind test-engineer`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		role := args[0]

		manager := persona.NewBindingManager()

		// Check if binding exists
		currentPersona, err := manager.GetBinding(role)
		if err != nil {
			return fmt.Errorf("failed to check binding: %w", err)
		}

		if currentPersona == "" {
			cmd.Printf("No binding exists for role '%s'\n", role)
			return nil
		}

		// Remove binding
		if err := manager.RemoveBinding(role); err != nil {
			return fmt.Errorf("failed to remove binding: %w", err)
		}

		cmd.Printf("✓ Removed binding for role '%s' (was '%s')\n", role, currentPersona)
		return nil
	},
}

// personaBindingsCmd shows current bindings
var personaBindingsCmd = &cobra.Command{
	Use:   "bindings",
	Short: "Show current persona bindings",
	Long: `Display all current role-persona bindings in the project configuration.

Examples:
  ddx persona bindings               # Show all current bindings`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := persona.NewBindingManager()

		bindings, err := manager.GetAllBindings()
		if err != nil {
			return fmt.Errorf("failed to get bindings: %w", err)
		}

		if len(bindings) == 0 {
			cmd.Println("No persona bindings configured")
			return nil
		}

		// Display bindings in table format
		cmd.Println("Current Persona Bindings:")
		cmd.Println()
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ROLE\tPERSONA")
		fmt.Fprintln(w, "----\t-------")

		for role, personaName := range bindings {
			fmt.Fprintf(w, "%s\t%s\n", role, personaName)
		}

		return w.Flush()
	},
}

// personaLoadCmd loads personas into CLAUDE.md
var personaLoadCmd = &cobra.Command{
	Use:   "load [persona-name] [flags]",
	Short: "Load personas into CLAUDE.md",
	Long: `Load personas into CLAUDE.md for the AI assistant to use.

Without arguments, loads all bound personas based on project configuration.
With a persona name, loads that specific persona.
With --role flag, loads the persona bound to that role.

Examples:
  ddx persona load                    # Load all bound personas
  ddx persona load strict-reviewer    # Load specific persona
  ddx persona load --role code-reviewer  # Load persona for specific role`,
	RunE: func(cmd *cobra.Command, args []string) error {
		role, _ := cmd.Flags().GetString("role")

		loader := persona.NewPersonaLoader()
		manager := persona.NewBindingManager()
		injector := persona.NewClaudeInjector()

		if len(args) == 1 {
			// Load specific persona
			personaName := args[0]
			p, err := loader.LoadPersona(personaName)
			if err != nil {
				return fmt.Errorf("failed to load persona '%s': %w", personaName, err)
			}

			// Use the first role from the persona
			roleToUse := p.Roles[0]
			if err := injector.InjectPersona(p, roleToUse); err != nil {
				return fmt.Errorf("failed to inject persona: %w", err)
			}

			cmd.Printf("✓ Loaded persona '%s' for role '%s' into CLAUDE.md\n", personaName, roleToUse)
			return nil
		}

		if role != "" {
			// Load persona for specific role
			personaName, err := manager.GetBinding(role)
			if err != nil {
				// Check if it's a missing config error
				if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "No such file") {
					return NewExitError(ExitCodeNoConfig, "No .ddx.yml configuration found")
				}
				return fmt.Errorf("failed to get binding for role '%s': %w", role, err)
			}

			if personaName == "" {
				return fmt.Errorf("no persona bound to role '%s'", role)
			}

			p, err := loader.LoadPersona(personaName)
			if err != nil {
				return fmt.Errorf("failed to load persona '%s': %w", personaName, err)
			}

			if err := injector.InjectPersona(p, role); err != nil {
				return fmt.Errorf("failed to inject persona: %w", err)
			}

			cmd.Printf("✓ Loaded persona '%s' for role '%s' into CLAUDE.md\n", personaName, role)
			return nil
		}

		// Load all bound personas
		bindings, err := manager.GetAllBindings()
		if err != nil {
			// Check if it's a missing config error
			if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "No such file") ||
				strings.Contains(err.Error(), "file does not exist") {
				return NewExitError(ExitCodeNoConfig, "No .ddx.yml configuration found")
			}
			return fmt.Errorf("failed to get bindings: %w", err)
		}

		if len(bindings) == 0 {
			cmd.Println("No persona bindings configured. Use 'ddx persona bind' to create bindings.")
			return nil
		}

		// Load all personas
		personas := make(map[string]*persona.Persona)
		var lastError error
		for role, personaName := range bindings {
			p, err := loader.LoadPersona(personaName)
			if err != nil {
				cmd.Printf("Warning: Failed to load persona '%s' for role '%s': %v\n", personaName, role, err)
				lastError = err
				continue
			}
			personas[role] = p
		}

		if len(personas) == 0 {
			if lastError != nil {
				return fmt.Errorf("failed to load personas: %w", lastError)
			}
			return fmt.Errorf("no personas could be loaded")
		}

		if err := injector.InjectMultiple(personas); err != nil {
			return fmt.Errorf("failed to inject personas: %w", err)
		}

		cmd.Printf("✓ Loaded %d personas into CLAUDE.md\n", len(personas))
		for role, p := range personas {
			cmd.Printf("  - %s: %s\n", role, p.Name)
		}

		return nil
	},
}

// personaUnloadCmd removes personas from CLAUDE.md
var personaUnloadCmd = &cobra.Command{
	Use:   "unload",
	Short: "Remove all personas from CLAUDE.md",
	Long: `Remove all personas from CLAUDE.md, leaving other content intact.

Examples:
  ddx persona unload                 # Remove all personas from CLAUDE.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		injector := persona.NewClaudeInjector()

		if err := injector.RemovePersonas(); err != nil {
			return fmt.Errorf("failed to remove personas: %w", err)
		}

		cmd.Println("✓ Removed all personas from CLAUDE.md")
		return nil
	},
}

// personaStatusCmd shows currently loaded personas
var personaStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show currently loaded personas",
	Long: `Display which personas are currently loaded in CLAUDE.md.

Examples:
  ddx persona status                 # Show loaded personas`,
	RunE: func(cmd *cobra.Command, args []string) error {
		injector := persona.NewClaudeInjector()

		// Check if CLAUDE.md exists
		if _, err := os.Stat("CLAUDE.md"); os.IsNotExist(err) {
			cmd.Println("No CLAUDE.md file found")
			return nil
		}

		loadedPersonas, err := injector.GetLoadedPersonas()
		if err != nil {
			return fmt.Errorf("failed to get loaded personas: %w", err)
		}

		if len(loadedPersonas) == 0 {
			cmd.Println("No personas currently loaded in CLAUDE.md")
			return nil
		}

		cmd.Println("Loaded Personas:")
		cmd.Println()
		for role, personaName := range loadedPersonas {
			cmd.Printf("  %s: %s\n", role, personaName)
		}

		return nil
	},
}

func init() {
	// Add flags to list command
	personaListCmd.Flags().String("role", "", "Filter personas by role")
	personaListCmd.Flags().StringSlice("tag", []string{}, "Filter personas by tags")

	// Add flags to load command
	personaLoadCmd.Flags().String("role", "", "Load persona for specific role")
}
