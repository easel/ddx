package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/easel/ddx/internal/config"
	// "github.com/easel/ddx/internal/persona" // Will be used when persona package is implemented
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Command registration is now handled by command_factory.go
// This file only contains the run function implementation

// PersonaMetadata represents parsed persona frontmatter
type PersonaMetadata struct {
	Name        string   `yaml:"name"`
	Roles       []string `yaml:"roles"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
}

// runPersona implements the persona command logic
func runPersona(cmd *cobra.Command, args []string) error {
	// Get flag values locally
	listFlag, _ := cmd.Flags().GetBool("list")
	showFlag, _ := cmd.Flags().GetString("show")
	bindFlag, _ := cmd.Flags().GetString("bind")
	roleFlag, _ := cmd.Flags().GetString("role")

	// Handle subcommands
	if len(args) > 0 {
		switch args[0] {
		case "list":
			return listPersonas(cmd)
		case "show":
			if len(args) < 2 {
				return fmt.Errorf("persona name required")
			}
			return showPersona(cmd, args[1])
		case "bind":
			if len(args) < 3 {
				return fmt.Errorf("role and persona name required")
			}
			return bindPersona(cmd, args[1], args[2])
		case "load":
			return loadPersonas(cmd, args[1:]...)
		case "bindings":
			return showBindings(cmd)
		case "status":
			return showStatus(cmd)
		}
	}

	// Handle flags
	if listFlag {
		return listPersonas(cmd)
	}

	if showFlag != "" {
		return showPersona(cmd, showFlag)
	}

	if bindFlag != "" && roleFlag != "" {
		return bindPersona(cmd, roleFlag, bindFlag)
	}

	// Default to list
	return listPersonas(cmd)
}

func listPersonas(cmd *cobra.Command) error {
	// Get library path
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	personasDir := filepath.Join(libPath, "personas")

	// Check if personas directory exists
	if _, err := os.Stat(personasDir); os.IsNotExist(err) {
		fmt.Fprintln(cmd.OutOrStdout(), "No personas directory found")
		return nil
	}

	// Get filter flags
	roleFilter, _ := cmd.Flags().GetString("role")
	tagFilter, _ := cmd.Flags().GetString("tag")

	fmt.Fprintln(cmd.OutOrStdout(), "Available Personas:")
	fmt.Fprintln(cmd.OutOrStdout())

	// Create tabwriter for aligned output
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PERSONA\tROLE\tDESCRIPTION")
	fmt.Fprintln(w, "-------\t----\t-----------")

	// Read personas directory
	entries, err := os.ReadDir(personasDir)
	if err != nil {
		return fmt.Errorf("failed to read personas directory: %w", err)
	}

	// Track if any personas were found
	personasFound := false

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")

		// Read and parse persona file
		content, err := os.ReadFile(filepath.Join(personasDir, entry.Name()))
		if err != nil {
			continue
		}

		// Parse frontmatter
		metadata := parsePersonaMetadata(string(content))
		if metadata == nil {
			// Fallback to simple parsing
			metadata = &PersonaMetadata{
				Name:        name,
				Roles:       []string{"general"},
				Description: name,
			}
		}

		// Apply role filter
		if roleFilter != "" {
			hasRole := false
			for _, role := range metadata.Roles {
				if role == roleFilter {
					hasRole = true
					break
				}
			}
			if !hasRole {
				continue
			}
		}

		// Apply tag filter
		if tagFilter != "" {
			hasTag := false
			for _, tag := range metadata.Tags {
				if tag == tagFilter {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		// Display persona
		roleStr := "general"
		if len(metadata.Roles) > 0 {
			roleStr = metadata.Roles[0]
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", name, roleStr, metadata.Description)
		personasFound = true
	}

	w.Flush()

	// Show message if no personas found
	if !personasFound {
		if roleFilter != "" || tagFilter != "" {
			fmt.Fprintln(cmd.OutOrStdout(), "\nNo personas found matching criteria")
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "No personas found")
		}
	}

	return nil
}

// parsePersonaMetadata parses YAML frontmatter from persona content
func parsePersonaMetadata(content string) *PersonaMetadata {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return nil
	}

	// Find end of frontmatter
	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return nil
	}

	// Parse YAML frontmatter
	frontmatter := strings.Join(lines[1:endIdx], "\n")
	var metadata PersonaMetadata
	if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
		return nil
	}

	return &metadata
}

func showPersona(cmd *cobra.Command, personaName string) error {
	// Get library path
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	personaPath := filepath.Join(libPath, "personas", personaName+".md")

	// Check if persona exists
	if _, err := os.Stat(personaPath); os.IsNotExist(err) {
		return fmt.Errorf("persona '%s' not found", personaName)
	}

	// Read persona content
	content, err := os.ReadFile(personaPath)
	if err != nil {
		return fmt.Errorf("failed to read persona: %w", err)
	}

	// Parse metadata
	metadata := parsePersonaMetadata(string(content))
	if metadata != nil {
		// Display formatted metadata
		fmt.Fprintf(cmd.OutOrStdout(), "Name: %s\n", metadata.Name)
		fmt.Fprintf(cmd.OutOrStdout(), "Roles: %s\n", strings.Join(metadata.Roles, ", "))
		fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", metadata.Description)
		fmt.Fprintf(cmd.OutOrStdout(), "Tags: %s\n", strings.Join(metadata.Tags, ", "))

		// Display content after frontmatter
		lines := strings.Split(string(content), "\n")
		contentStart := 0
		foundEnd := false
		for i, line := range lines {
			if i > 0 && line == "---" {
				contentStart = i + 1
				foundEnd = true
				break
			}
		}
		if foundEnd && contentStart < len(lines) {
			fmt.Fprintln(cmd.OutOrStdout())
			fmt.Fprint(cmd.OutOrStdout(), strings.Join(lines[contentStart:], "\n"))
		}
	} else {
		// No frontmatter, display raw content
		fmt.Fprint(cmd.OutOrStdout(), string(content))
	}
	return nil
}

func bindPersona(cmd *cobra.Command, role, personaName string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if persona exists
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	personaPath := filepath.Join(libPath, "personas", personaName+".md")
	if _, err := os.Stat(personaPath); os.IsNotExist(err) {
		return fmt.Errorf("persona '%s' not found", personaName)
	}

	// Update persona bindings
	if cfg.PersonaBindings == nil {
		cfg.PersonaBindings = make(map[string]string)
	}
	cfg.PersonaBindings[role] = personaName

	// Save config
	// TODO: Add Save method to Config or use yaml.Marshal
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}
	if err := os.WriteFile(".ddx.yml", data, 0644); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✅ Bound role '%s' to persona '%s'\n", role, personaName)
	return nil
}

func loadPersonas(cmd *cobra.Command, personas ...string) error {
	// Check if config file exists when loading all personas
	if len(personas) == 0 {
		if _, err := os.Stat(".ddx.yml"); os.IsNotExist(err) {
			return fmt.Errorf("No .ddx.yml configuration found")
		}
	}

	// Load config to get persona bindings
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get library path
	libPath, err := config.GetLibraryPath(getLibraryPath())
	if err != nil {
		return fmt.Errorf("failed to get library path: %w", err)
	}

	// Read CLAUDE.md if it exists
	claudePath := "CLAUDE.md"
	var claudeContent string
	if data, err := os.ReadFile(claudePath); err == nil {
		claudeContent = string(data)
	} else {
		// Create new CLAUDE.md
		claudeContent = "# CLAUDE.md\n\nProject guidance for my application."
	}

	// Remove existing persona section if present
	startMarker := "<!-- PERSONAS:START -->"
	endMarker := "<!-- PERSONAS:END -->"
	startIdx := strings.Index(claudeContent, startMarker)
	if startIdx != -1 {
		endIdx := strings.Index(claudeContent, endMarker)
		if endIdx != -1 {
			claudeContent = claudeContent[:startIdx] + claudeContent[endIdx+len(endMarker):]
		}
	}

	// Build persona content
	var personaSection strings.Builder
	personaSection.WriteString("\n" + startMarker + "\n")
	personaSection.WriteString("## Active Personas\n\n")

	// Track loaded personas
	loadedPersonas := []string{}

	// If specific personas requested, load those; otherwise load all bound personas
	if len(personas) > 0 {
		// Load specific personas
		for _, personaName := range personas {
			personaPath := filepath.Join(libPath, "personas", personaName+".md")
			if content, err := os.ReadFile(personaPath); err == nil {
				// Validate persona content if it has frontmatter
				contentStr := string(content)
				if strings.HasPrefix(contentStr, "---\n") {
					// Try to parse the frontmatter to validate it
					lines := strings.Split(contentStr, "\n")
					endIdx := -1
					for i := 1; i < len(lines); i++ {
						if lines[i] == "---" {
							endIdx = i
							break
						}
					}
					if endIdx > 0 {
						frontmatter := strings.Join(lines[1:endIdx], "\n")
						var metadata PersonaMetadata
						if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
							return fmt.Errorf("failed to parse YAML frontmatter in persona '%s': %w", personaName, err)
						}
					}
				}
				// Just add the content - personas have their own titles
				personaSection.WriteString(string(content) + "\n")
				loadedPersonas = append(loadedPersonas, personaName)
			} else if os.IsNotExist(err) {
				return fmt.Errorf("persona '%s' not found", personaName)
			}
		}
	} else {
		// Load all bound personas from config
		if cfg.PersonaBindings != nil {
			for role, personaName := range cfg.PersonaBindings {
				personaPath := filepath.Join(libPath, "personas", personaName+".md")
				if content, err := os.ReadFile(personaPath); err == nil {
					// Validate persona content if it has frontmatter
					contentStr := string(content)
					if strings.HasPrefix(contentStr, "---\n") {
						// Try to parse the frontmatter to validate it
						lines := strings.Split(contentStr, "\n")
						endIdx := -1
						for i := 1; i < len(lines); i++ {
							if lines[i] == "---" {
								endIdx = i
								break
							}
						}
						if endIdx > 0 {
							frontmatter := strings.Join(lines[1:endIdx], "\n")
							var metadata PersonaMetadata
							if err := yaml.Unmarshal([]byte(frontmatter), &metadata); err != nil {
								return fmt.Errorf("failed to parse YAML frontmatter in persona '%s': %w", personaName, err)
							}
						}
					}
					// Add role header with proper capitalization
					capitalizedRole := strings.Title(strings.ReplaceAll(role, "-", " "))
					personaSection.WriteString(fmt.Sprintf("### %s: %s\n", capitalizedRole, personaName))
					personaSection.WriteString(string(content) + "\n")
					loadedPersonas = append(loadedPersonas, personaName)
				}
			}
		}
	}

	personaSection.WriteString(endMarker + "\n")

	// Append persona section to CLAUDE.md
	claudeContent += personaSection.String()

	// Write updated CLAUDE.md
	if err := os.WriteFile(claudePath, []byte(claudeContent), 0644); err != nil {
		return fmt.Errorf("failed to write CLAUDE.md: %w", err)
	}

	// Provide detailed output
	if len(personas) > 0 {
		// Specific personas loaded
		if len(loadedPersonas) == 1 {
			fmt.Fprintf(cmd.OutOrStdout(), "✅ Loaded persona '%s' into CLAUDE.md\n", loadedPersonas[0])
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "✅ Loaded %d personas into CLAUDE.md\n", len(loadedPersonas))
		}
	} else {
		// All bound personas loaded
		if len(loadedPersonas) > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "✅ Loaded %d personas (%s) into CLAUDE.md\n",
				len(loadedPersonas), strings.Join(loadedPersonas, ", "))
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "No bound personas to load")
		}
	}
	return nil
}

// showBindings displays current persona-role bindings
func showBindings(cmd *cobra.Command) error {
	// Check if config file exists
	if _, err := os.Stat(".ddx.yml"); os.IsNotExist(err) {
		return fmt.Errorf("no .ddx.yml configuration found")
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if cfg.PersonaBindings == nil || len(cfg.PersonaBindings) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No persona bindings configured")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Current Persona Bindings:")
	fmt.Fprintln(cmd.OutOrStdout())

	// Create tabwriter for aligned output
	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ROLE\tPERSONA")
	fmt.Fprintln(w, "----\t-------")

	for role, persona := range cfg.PersonaBindings {
		fmt.Fprintf(w, "%s\t%s\n", role, persona)
	}

	w.Flush()
	return nil
}

// showStatus displays the status of active personas
func showStatus(cmd *cobra.Command) error {
	// Check if CLAUDE.md exists
	claudePath := "CLAUDE.md"
	if _, err := os.Stat(claudePath); os.IsNotExist(err) {
		fmt.Fprintln(cmd.OutOrStdout(), "No CLAUDE.md file found - no personas loaded")
		return nil
	}

	// Read CLAUDE.md
	content, err := os.ReadFile(claudePath)
	if err != nil {
		return fmt.Errorf("failed to read CLAUDE.md: %w", err)
	}

	// Check for persona markers
	claudeStr := string(content)
	if strings.Contains(claudeStr, "<!-- PERSONAS:START -->") &&
		strings.Contains(claudeStr, "<!-- PERSONAS:END -->") {
		// Extract persona section
		startIdx := strings.Index(claudeStr, "<!-- PERSONAS:START -->")
		endIdx := strings.Index(claudeStr, "<!-- PERSONAS:END -->")
		if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
			personaSection := claudeStr[startIdx:endIdx]

			// Parse loaded personas
			loadedPersonas := []string{}
			loadedRoles := []string{}
			lines := strings.Split(personaSection, "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "### ") {
					// Parse role and persona from header
					header := strings.TrimPrefix(line, "### ")
					parts := strings.Split(header, ": ")
					if len(parts) == 2 {
						loadedRoles = append(loadedRoles, parts[0])
						loadedPersonas = append(loadedPersonas, parts[1])
					}
				}
			}

			if len(loadedPersonas) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Loaded Personas:")
				for i, persona := range loadedPersonas {
					fmt.Fprintf(cmd.OutOrStdout(), "  - %s (%s)\n", persona, loadedRoles[i])
				}
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "No personas currently loaded")
			}
		}
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), "No personas currently loaded")
	}

	// Show bindings from config
	cfg, err := config.Load()
	if err == nil && cfg.PersonaBindings != nil && len(cfg.PersonaBindings) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "\n%d persona binding(s) configured\n", len(cfg.PersonaBindings))
	}

	return nil
}
