package obsidian

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// FrontmatterGenerator generates appropriate frontmatter for files
type FrontmatterGenerator struct {
	detector *FileTypeDetector
}

// NewFrontmatterGenerator creates a new frontmatter generator
func NewFrontmatterGenerator() *FrontmatterGenerator {
	return &FrontmatterGenerator{
		detector: NewFileTypeDetector(),
	}
}

// Generate creates frontmatter for a given file
func (g *FrontmatterGenerator) Generate(file *MarkdownFile) (*Frontmatter, error) {
	fm := &Frontmatter{
		Created: time.Now(),
		Updated: time.Now(),
		Tags:    []string{},
	}

	// Extract title from content or path
	fm.Title = g.extractTitle(file)

	// Set type based on file type
	fm.Type = string(file.FileType)

	// Generate tags based on file type and location
	fm.Tags = g.generateTags(file)

	// Add type-specific metadata
	switch file.FileType {
	case FileTypePhase:
		g.addPhaseMetadata(fm, file)
	case FileTypeArtifact, FileTypeTemplate, FileTypePrompt, FileTypeExample:
		g.addArtifactMetadata(fm, file)
	case FileTypeEnforcer:
		g.addEnforcerMetadata(fm, file)
	case FileTypeCoordinator:
		g.addCoordinatorMetadata(fm, file)
	case FileTypePrinciple:
		g.addPrincipleMetadata(fm, file)
	case FileTypeFeature:
		g.addFeatureMetadata(fm, file)
	}

	return fm, nil
}

// extractTitle extracts the title from markdown content or generates from path
func (g *FrontmatterGenerator) extractTitle(file *MarkdownFile) string {
	// First try to extract from content
	if title := g.extractTitleFromContent(file.Content); title != "" {
		return title
	}

	// Fallback to generating from path
	return ExtractTitleFromPath(file.Path)
}

// extractTitleFromContent extracts title from markdown content
func (g *FrontmatterGenerator) extractTitleFromContent(content string) string {
	// Look for first H1 heading
	re := regexp.MustCompile(`(?m)^#\s+(.+)$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		// Clean up common patterns
		title = strings.TrimPrefix(title, "Feature Specification: ")
		title = strings.TrimPrefix(title, "Technical Design: ")
		title = strings.TrimPrefix(title, "Build Implementation: ")
		// Remove FEAT-XXX patterns
		if idx := strings.Index(title, " - "); idx > 0 {
			if strings.HasPrefix(title, "FEAT-") || strings.HasPrefix(title, "[[FEAT-") {
				title = title[idx+3:]
			}
		}
		return strings.TrimSpace(title)
	}
	return ""
}

// generateTags creates tags based on file type and path
func (g *FrontmatterGenerator) generateTags(file *MarkdownFile) []string {
	tags := []string{"helix"}

	// Add hierarchical tags based on file type
	hierarchicalTags := file.FileType.GetHierarchicalTags()
	for _, tag := range hierarchicalTags {
		if tag != "helix" { // Don't duplicate base tag
			tags = append(tags, tag)
		}
	}

	// Add phase tag if applicable
	if phase := GetPhaseFromPath(file.Path); phase != "" {
		phaseTag := fmt.Sprintf("helix/phase/%s", phase)
		if !contains(tags, phaseTag) {
			tags = append(tags, phaseTag)
		}
	}

	// Add artifact category tag
	if category := GetArtifactCategory(file.Path); category != "" {
		categoryTag := fmt.Sprintf("helix/artifact/%s", strings.ReplaceAll(category, "-", "/"))
		if !contains(tags, categoryTag) {
			tags = append(tags, categoryTag)
		}
	}

	return tags
}

// addPhaseMetadata adds phase-specific metadata
func (g *FrontmatterGenerator) addPhaseMetadata(fm *Frontmatter, file *MarkdownFile) {
	phase := GetPhaseFromPath(file.Path)
	fm.PhaseID = phase
	fm.PhaseNum = GetPhaseNumber(phase)

	// Set next/previous phases
	if nextPhase := GetNextPhase(phase); nextPhase != "" {
		fm.NextPhase = fmt.Sprintf("[[%s Phase]]", strings.Title(nextPhase))
	}
	if prevPhase := GetPreviousPhase(phase); prevPhase != "" {
		fm.PrevPhase = fmt.Sprintf("[[%s Phase]]", strings.Title(prevPhase))
	}

	// Add gates and artifacts placeholders
	fm.Gates = &Gates{
		Entry: []string{"[[TODO: Define entry criteria]]"},
		Exit:  []string{"[[TODO: Define exit criteria]]"},
	}
	fm.Artifacts = &Artifacts{
		Required: []string{"[[TODO: List required artifacts]]"},
		Optional: []string{"[[TODO: List optional artifacts]]"},
	}

	// Add phase-specific aliases
	fm.Aliases = []string{
		fmt.Sprintf("%s Phase", strings.Title(phase)),
	}
}

// addArtifactMetadata adds artifact-specific metadata
func (g *FrontmatterGenerator) addArtifactMetadata(fm *Frontmatter, file *MarkdownFile) {
	fm.Phase = GetPhaseFromPath(file.Path)
	fm.ArtifactCategory = GetArtifactCategory(file.Path)

	// Set complexity based on content or path
	fm.Complexity = GetComplexityFromPath(file.Path)

	// Add common artifact fields
	fm.Prerequisites = []string{}
	fm.Outputs = []string{}

	// Estimate time based on file type
	switch file.FileType {
	case FileTypeTemplate:
		fm.TimeEstimate = "30-60 minutes"
	case FileTypePrompt:
		fm.TimeEstimate = "15-30 minutes"
	case FileTypeExample:
		fm.TimeEstimate = "5-15 minutes"
	default:
		fm.TimeEstimate = "1-2 hours"
	}

	// Add relevant aliases
	artifactName := strings.Title(strings.ReplaceAll(fm.ArtifactCategory, "-", " "))
	switch file.FileType {
	case FileTypeTemplate:
		fm.Aliases = []string{
			fmt.Sprintf("%s Template", artifactName),
		}
	case FileTypePrompt:
		fm.Aliases = []string{
			fmt.Sprintf("%s Prompt", artifactName),
		}
	case FileTypeExample:
		fm.Aliases = []string{
			fmt.Sprintf("%s Example", artifactName),
		}
	}
}

// addEnforcerMetadata adds enforcer-specific metadata
func (g *FrontmatterGenerator) addEnforcerMetadata(fm *Frontmatter, file *MarkdownFile) {
	phase := GetPhaseFromPath(file.Path)
	fm.Phase = phase

	// Add enforcer-specific tags
	if phase != "" {
		enforcerTag := fmt.Sprintf("helix/phase/%s/enforcer", phase)
		if !contains(fm.Tags, enforcerTag) {
			fm.Tags = append(fm.Tags, enforcerTag)
		}
	}

	fm.Aliases = []string{
		fmt.Sprintf("%s Phase Enforcer", strings.Title(phase)),
		fmt.Sprintf("%s Guardian", strings.Title(phase)),
	}
}

// addCoordinatorMetadata adds coordinator-specific metadata
func (g *FrontmatterGenerator) addCoordinatorMetadata(fm *Frontmatter, file *MarkdownFile) {
	fm.Aliases = []string{
		"HELIX Coordinator",
		"Workflow Coordinator",
	}
}

// addPrincipleMetadata adds principle-specific metadata
func (g *FrontmatterGenerator) addPrincipleMetadata(fm *Frontmatter, file *MarkdownFile) {
	fm.Aliases = []string{
		"HELIX Principles",
		"Workflow Principles",
	}
}

// addFeatureMetadata adds feature specification metadata
func (g *FrontmatterGenerator) addFeatureMetadata(fm *Frontmatter, file *MarkdownFile) {
	fm.WorkflowPhase = "frame"
	fm.ArtifactType = "feature-specification"

	// Try to extract feature ID from filename or content
	if featureID := g.extractFeatureID(file); featureID != "" {
		fm.FeatureID = featureID
	}

	// Extract priority, owner, status from content if present
	fm.Priority = g.extractFromContent(file.Content, "Priority")
	fm.Owner = g.extractFromContent(file.Content, "Owner")
	fm.Status = g.extractFromContent(file.Content, "Status")

	// Default values if not found
	if fm.Priority == "" {
		fm.Priority = "P2"
	}
	if fm.Status == "" {
		fm.Status = "draft"
	}

	// Add related artifacts
	if fm.FeatureID != "" {
		fm.Related = []string{
			fmt.Sprintf("[[%s-technical-design]]", fm.FeatureID),
			fmt.Sprintf("[[%s-implementation]]", fm.FeatureID),
		}
	}
}

// addDesignMetadata adds technical design metadata
func (g *FrontmatterGenerator) addDesignMetadata(fm *Frontmatter, file *MarkdownFile) {
	fm.WorkflowPhase = "design"
	fm.ArtifactType = "technical-design"

	// Try to extract feature ID from filename
	if featureID := g.extractFeatureID(file); featureID != "" {
		fm.FeatureID = featureID
		fm.Related = []string{
			fmt.Sprintf("[[%s-obsidian-integration]]", featureID),
			fmt.Sprintf("[[%s-implementation]]", featureID),
		}
	}

	fm.Status = "draft"
}

// addTestMetadata adds test specification metadata
func (g *FrontmatterGenerator) addTestMetadata(fm *Frontmatter, file *MarkdownFile) {
	fm.WorkflowPhase = "test"
	fm.ArtifactType = "test-specification"

	// Try to extract feature ID from filename
	if featureID := g.extractFeatureID(file); featureID != "" {
		fm.FeatureID = featureID
	}

	fm.Status = "draft"
}

// addImplementationMetadata adds implementation guide metadata
func (g *FrontmatterGenerator) addImplementationMetadata(fm *Frontmatter, file *MarkdownFile) {
	fm.WorkflowPhase = "build"
	fm.ArtifactType = "implementation-guide"

	// Try to extract feature ID from filename
	if featureID := g.extractFeatureID(file); featureID != "" {
		fm.FeatureID = featureID
		fm.Related = []string{
			fmt.Sprintf("[[%s-obsidian-integration]]", featureID),
			fmt.Sprintf("[[%s-technical-design]]", featureID),
		}
	}

	fm.Status = "ready"
}

// extractFeatureID extracts feature ID from filename or content
func (g *FrontmatterGenerator) extractFeatureID(file *MarkdownFile) string {
	// Check filename first
	re := regexp.MustCompile(`FEAT-(\d+)`)
	matches := re.FindStringSubmatch(file.Path)
	if len(matches) > 1 {
		return fmt.Sprintf("FEAT-%s", matches[1])
	}

	// Check content
	matches = re.FindStringSubmatch(file.Content)
	if len(matches) > 1 {
		return fmt.Sprintf("FEAT-%s", matches[1])
	}

	return ""
}

// extractFromContent extracts a field value from markdown content
func (g *FrontmatterGenerator) extractFromContent(content, field string) string {
	// Look for patterns like "**Priority**: P1" or "Priority: P1"
	patterns := []string{
		fmt.Sprintf(`\*\*%s\*\*:\s*([^\n\r]+)`, field),
		fmt.Sprintf(`%s:\s*([^\n\r]+)`, field),
		fmt.Sprintf(`\*\*%s\*\*\s*([^\n\r]+)`, field),
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			value := strings.TrimSpace(matches[1])
			// Clean up common markdown formatting
			value = strings.Trim(value, "[]")
			return value
		}
	}

	return ""
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
