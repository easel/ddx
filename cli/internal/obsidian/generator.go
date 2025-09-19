package obsidian

import (
	"regexp"
	"strings"
	"time"
)

// FrontmatterGenerator generates YAML frontmatter for markdown files
type FrontmatterGenerator struct {
	titleExtractor *TitleExtractor
}

// NewFrontmatterGenerator creates a new frontmatter generator
func NewFrontmatterGenerator() *FrontmatterGenerator {
	return &FrontmatterGenerator{
		titleExtractor: NewTitleExtractor(),
	}
}

// Generate creates frontmatter for a markdown file
func (fg *FrontmatterGenerator) Generate(file *MarkdownFile) (*Frontmatter, error) {
	now := time.Now()

	fm := &Frontmatter{
		Title:        fg.extractTitle(file),
		Type:         file.FileType.String(),
		Tags:         fg.generateTags(file),
		Created:      now,
		Updated:      now,
		TimeEstimate: fg.getTimeEstimateByType(file.FileType),
	}

	// Add file type specific fields
	switch file.FileType {
	case FileTypePhase:
		fg.addPhaseMetadata(fm, file)
	case FileTypeEnforcer:
		fg.addEnforcerFields(fm, file)
	case FileTypeTemplate, FileTypePrompt, FileTypeExample:
		fg.addArtifactMetadata(fm, file)
	case FileTypeFeature:
		fg.addFeatureFields(fm, file)
	case FileTypeCoordinator:
		fg.addCoordinatorFields(fm, file)
	case FileTypePrinciple:
		fg.addPrincipleFields(fm, file)
	}

	return fm, nil
}

// extractTitle extracts or generates a title for the file
func (fg *FrontmatterGenerator) extractTitle(file *MarkdownFile) string {
	// Try to extract from content first
	if title := fg.titleExtractor.ExtractFromContent(file.Content); title != "" {
		return title
	}

	// Fallback to path-based title
	return ExtractTitleFromPath(file.Path)
}

// extractTitleFromContent extracts title from markdown content
func (fg *FrontmatterGenerator) extractTitleFromContent(content string) string {
	return fg.titleExtractor.ExtractFromContent(content)
}

// extractFeatureID extracts feature ID from file path and content
func (fg *FrontmatterGenerator) extractFeatureID(file *MarkdownFile) string {
	// Try filename first
	if featureID := fg.extractFeatureIDFromPath(file.Path); featureID != "" {
		return featureID
	}

	// Try content
	return fg.extractFeatureIDFromContent(file.Content)
}

// extractFeatureIDFromPath extracts feature ID from filename
func (fg *FrontmatterGenerator) extractFeatureIDFromPath(path string) string {
	re := regexp.MustCompile(`FEAT-(\d+)`)
	if matches := re.FindStringSubmatch(path); len(matches) > 1 {
		return "FEAT-" + matches[1]
	}
	return ""
}

// extractFeatureIDFromContent extracts feature ID from content
func (fg *FrontmatterGenerator) extractFeatureIDFromContent(content string) string {
	re := regexp.MustCompile(`FEAT-(\d+)`)
	if matches := re.FindStringSubmatch(content); len(matches) > 1 {
		return "FEAT-" + matches[1]
	}
	return ""
}

// extractFromContent extracts field values from content
func (fg *FrontmatterGenerator) extractFromContent(content, field string) string {
	// Pattern for field extraction: **Field**: value or Field: value
	pattern := `(?:^\*\*` + regexp.QuoteMeta(field) + `\*\*|^` + regexp.QuoteMeta(field) + `):\s*(?:\[([^\]]+)\]|([^\n]+))`
	re := regexp.MustCompile(`(?m)` + pattern)

	if matches := re.FindStringSubmatch(content); len(matches) > 1 {
		// Return first non-empty capture group
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				return strings.TrimSpace(matches[i])
			}
		}
	}
	return ""
}

// generateTags generates hierarchical tags for a file
func (fg *FrontmatterGenerator) generateTags(file *MarkdownFile) []string {
	tags := file.FileType.GetHierarchicalTags()

	// Add phase-specific tags
	if phase := GetPhaseFromPath(file.Path); phase != "" {
		tags = append(tags, "helix/phase/"+phase)
	}

	// Add artifact-specific tags
	if file.FileType.IsArtifact() {
		if category := GetArtifactCategory(file.Path); category != "" {
			// Convert hyphens to slashes for hierarchical tags
			categoryTag := "helix/artifact/" + strings.ReplaceAll(category, "-", "/")
			tags = append(tags, categoryTag)
		}
	}

	return tags
}

// getTimeEstimateByType returns time estimate based on file type
func (fg *FrontmatterGenerator) getTimeEstimateByType(fileType FileType) string {
	switch fileType {
	case FileTypeTemplate:
		return "30-60 minutes"
	case FileTypePrompt:
		return "15-30 minutes"
	case FileTypeExample:
		return "5-15 minutes"
	case FileTypeFeature, FileTypeArtifact, FileTypePhase:
		return "1-2 hours"
	default:
		return "1-2 hours"
	}
}

// addPhaseMetadata adds phase-specific frontmatter fields
func (fg *FrontmatterGenerator) addPhaseMetadata(fm *Frontmatter, file *MarkdownFile) {
	phase := GetPhaseFromPath(file.Path)
	if phase == "" {
		return
	}

	fm.PhaseID = phase
	fm.PhaseNum = GetPhaseNumber(phase)

	// Set navigation links
	if next := GetNextPhase(phase); next != "" {
		fm.NextPhase = "[[" + strings.Title(next) + " Phase]]"
	}
	if prev := GetPreviousPhase(phase); prev != "" {
		fm.PrevPhase = "[[" + strings.Title(prev) + " Phase]]"
	}

	// Initialize gates and artifacts
	fm.Gates = &Gates{
		Entry: []string{},
		Exit:  []string{},
	}
	fm.Artifacts = &Artifacts{
		Required: []string{},
		Optional: []string{},
	}

	// Add aliases
	fm.Aliases = []string{strings.Title(phase) + " Phase"}
}

// addEnforcerFields adds enforcer-specific frontmatter fields
func (fg *FrontmatterGenerator) addEnforcerFields(fm *Frontmatter, file *MarkdownFile) {
	phase := GetPhaseFromPath(file.Path)
	if phase != "" {
		fm.Phase = phase
		// Add enforcer-specific tag
		fm.Tags = append(fm.Tags, "helix/phase/"+phase+"/enforcer")
		// Add aliases
		phaseName := strings.Title(phase)
		fm.Aliases = []string{
			phaseName + " Phase Enforcer",
			phaseName + " Guardian",
		}
	}
}

// addArtifactMetadata adds artifact-specific frontmatter fields
func (fg *FrontmatterGenerator) addArtifactMetadata(fm *Frontmatter, file *MarkdownFile) {
	category := GetArtifactCategory(file.Path)
	if category != "" {
		fm.ArtifactCategory = category
	}

	phase := GetPhaseFromPath(file.Path)
	if phase != "" {
		fm.Phase = phase
	}

	complexity := GetComplexityFromPath(file.Path)
	fm.Complexity = complexity

	// Initialize prerequisites and outputs
	fm.Prerequisites = []string{}
	fm.Outputs = []string{}
}

// addFeatureFields adds feature-specific frontmatter fields
func (fg *FrontmatterGenerator) addFeatureFields(fm *Frontmatter, file *MarkdownFile) {
	// Extract feature ID from filename
	if featureID := fg.extractFeatureID(file); featureID != "" {
		fm.FeatureID = featureID
	}

	phase := GetPhaseFromPath(file.Path)
	if phase != "" {
		fm.WorkflowPhase = phase
	}

	// Extract priority from content or use default
	if priority := fg.extractFromContent(file.Content, "Priority"); priority != "" {
		fm.Priority = priority
	} else {
		fm.Priority = "P2"
	}

	// Extract owner from content
	if owner := fg.extractFromContent(file.Content, "Owner"); owner != "" {
		fm.Owner = owner
	}

	// Extract status from content or use default
	if status := fg.extractFromContent(file.Content, "Status"); status != "" {
		fm.Status = status
	} else {
		fm.Status = "draft"
	}

	// Set artifact type for feature files
	fm.ArtifactType = "feature-specification"
}

// addCoordinatorFields adds coordinator-specific frontmatter fields
func (fg *FrontmatterGenerator) addCoordinatorFields(fm *Frontmatter, file *MarkdownFile) {
	// Add aliases
	fm.Aliases = []string{
		"HELIX Coordinator",
		"Workflow Coordinator",
	}
}

// addPrincipleFields adds principle-specific frontmatter fields
func (fg *FrontmatterGenerator) addPrincipleFields(fm *Frontmatter, file *MarkdownFile) {
	fm.Tags = append(fm.Tags, "helix/principle")
}

// TitleExtractor extracts titles from markdown content
type TitleExtractor struct {
	headingRegex *regexp.Regexp
}

// NewTitleExtractor creates a new title extractor
func NewTitleExtractor() *TitleExtractor {
	return &TitleExtractor{
		headingRegex: regexp.MustCompile(`^#\s+(.+)$`),
	}
}

// ExtractFromContent extracts the first H1 heading from markdown content
func (te *TitleExtractor) ExtractFromContent(content string) string {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := te.headingRegex.FindStringSubmatch(line); len(matches) > 1 {
			title := strings.TrimSpace(matches[1])
			if title != "" {
				return te.cleanTitle(title)
			}
		}
	}

	return ""
}

// cleanTitle extracts the core title from various title formats
func (te *TitleExtractor) cleanTitle(title string) string {
	// Handle patterns like "Feature Specification: FEAT-001 - User Authentication"
	// Only clean if it contains ": FEAT-" followed by " - "
	if strings.Contains(title, ": FEAT-") && strings.Contains(title, " - ") {
		parts := strings.Split(title, " - ")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[len(parts)-1])
		}
	}

	// Handle patterns like "[[FEAT-004]] - Database Migration"
	if strings.Contains(title, "]] - ") {
		parts := strings.Split(title, "]] - ")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1])
		}
	}

	// Handle wikilinks in title like "[[FEAT-004]] Database Migration"
	wikilinkRegex := regexp.MustCompile(`\[\[[^\]]+\]\]\s*(.*)`)
	if matches := wikilinkRegex.FindStringSubmatch(title); len(matches) > 1 {
		remaining := strings.TrimSpace(matches[1])
		if remaining != "" {
			return remaining
		}
	}

	return title
}
