package converter

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/easel/ddx/internal/obsidian"
)

// LinkConverter converts markdown links to Obsidian wikilinks
type LinkConverter struct {
	fileIndex map[string]*obsidian.MarkdownFile
	pathIndex map[string]*obsidian.MarkdownFile
}

// NewLinkConverter creates a new link converter
func NewLinkConverter() *LinkConverter {
	return &LinkConverter{
		fileIndex: make(map[string]*obsidian.MarkdownFile),
		pathIndex: make(map[string]*obsidian.MarkdownFile),
	}
}

// BuildIndex builds an index of files for link resolution
func (lc *LinkConverter) BuildIndex(files []*obsidian.MarkdownFile) {
	lc.fileIndex = make(map[string]*obsidian.MarkdownFile)
	lc.pathIndex = make(map[string]*obsidian.MarkdownFile)

	for _, file := range files {
		// Index by path
		lc.pathIndex[file.Path] = file

		// Index by title if available
		if file.HasFrontmatter() && file.Frontmatter.Title != "" {
			lc.fileIndex[file.Frontmatter.Title] = file
		}

		// Index by aliases if available
		if file.HasFrontmatter() && len(file.Frontmatter.Aliases) > 0 {
			for _, alias := range file.Frontmatter.Aliases {
				if alias != "" {
					lc.fileIndex[alias] = file
				}
			}
		}

		// Index by filename without extension
		filename := filepath.Base(file.Path)
		filename = strings.TrimSuffix(filename, ".md")
		if filename != "README" {
			lc.fileIndex[filename] = file
		}
	}
}

// ConvertContent converts markdown links to wikilinks in content
func (lc *LinkConverter) ConvertContent(content string) string {
	// First convert markdown links
	result := lc.ConvertLinks(content)

	// Then convert phase and artifact references
	result = lc.convertPhaseReferences(result)
	result = lc.convertArtifactReferences(result)
	result = lc.convertWorkflowReferences(result)

	return result
}

// ConvertLinks converts markdown links to wikilinks in content
func (lc *LinkConverter) ConvertLinks(content string) string {
	// Regex to match markdown links: [text](url)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	result := linkRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Skip if already inside a wikilink
		if lc.isInWikilink(content, match) {
			return match
		}

		// Extract text and URL
		parts := linkRegex.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		text := parts[1]
		url := parts[2]

		// Skip external links (http/https)
		if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
			return match
		}

		// Try to resolve the link
		if targetFile := lc.resolveLink(url); targetFile != nil {
			return lc.createWikilink(targetFile, text)
		}

		// Keep original if can't resolve
		return match
	})

	return result
}

// isInWikilink checks if a markdown link is already inside a wikilink
func (lc *LinkConverter) isInWikilink(content, match string) bool {
	index := strings.Index(content, match)
	if index == -1 {
		return false
	}

	// Look for [[ before the match
	before := content[:index]
	afterOpenBrackets := strings.LastIndex(before, "[[")
	afterCloseBrackets := strings.LastIndex(before, "]]")

	// If we found [[ and no ]] after it, we're inside a wikilink
	return afterOpenBrackets != -1 && (afterCloseBrackets == -1 || afterOpenBrackets > afterCloseBrackets)
}

// resolveLink attempts to resolve a relative path to a file
func (lc *LinkConverter) resolveLink(url string) *obsidian.MarkdownFile {
	originalUrl := url

	// Clean up the URL for index lookup
	cleanUrl := strings.TrimPrefix(url, "./")
	cleanUrl = strings.TrimPrefix(cleanUrl, "../")

	// Try exact path match first
	for path, file := range lc.pathIndex {
		if strings.HasSuffix(path, cleanUrl) {
			return file
		}
	}

	// Try filename match
	filename := filepath.Base(cleanUrl)
	filename = strings.TrimSuffix(filename, ".md")

	if file, exists := lc.fileIndex[filename]; exists {
		return file
	}

	// If not found in index, create a synthetic file based on path analysis
	// Use original URL to preserve relative path context
	return lc.createSyntheticFile(originalUrl)
}

// createSyntheticFile creates a synthetic MarkdownFile based on path analysis
func (lc *LinkConverter) createSyntheticFile(path string) *obsidian.MarkdownFile {
	// Generate title based on path patterns
	title := lc.generateTitleFromPath(path)
	if title == "" {
		return nil
	}

	return &obsidian.MarkdownFile{
		Path: path,
		Frontmatter: &obsidian.Frontmatter{
			Title: title,
		},
	}
}

// generateTitleFromPath generates a title based on file path patterns
func (lc *LinkConverter) generateTitleFromPath(path string) string {
	// Clean up path for analysis
	cleanPath := strings.TrimPrefix(path, "./")
	cleanPath = strings.TrimPrefix(cleanPath, "../")

	filename := filepath.Base(cleanPath)
	filename = strings.TrimSuffix(filename, ".md")

	// Handle common file patterns
	if strings.Contains(cleanPath, "phases/") && filename == "README" {
		// Extract phase from path
		if phase := lc.getPhaseFromPath(cleanPath); phase != "" {
			return strings.Title(phase) + " Phase"
		}
	}

	if strings.Contains(cleanPath, "phases/") && filename == "enforcer" {
		// Extract phase from path
		if phase := lc.getPhaseFromPath(cleanPath); phase != "" {
			return strings.Title(phase) + " Phase Enforcer"
		}
	}

	if strings.Contains(cleanPath, "artifacts/") {
		// Extract artifact category and type
		if category := lc.getArtifactCategoryFromPath(cleanPath); category != "" {
			categoryName := strings.Title(strings.ReplaceAll(category, "-", " "))
			switch filename {
			case "template":
				return categoryName + " Template"
			case "prompt":
				return categoryName + " Prompt"
			case "example":
				return categoryName + " Example"
			}
		}
	}

	if filename == "coordinator" {
		return "HELIX Workflow Coordinator"
	}

	if filename == "principles" {
		return "HELIX Principles"
	}

	// Handle feature files - extract just the FEAT-XXX part
	if strings.Contains(filename, "FEAT-") {
		// Extract FEAT-XXX from FEAT-001-auth
		parts := strings.Split(filename, "-")
		if len(parts) >= 2 {
			return parts[0] + "-" + parts[1]
		}
		return filename
	}

	return ""
}

// getPhaseFromPath extracts phase name from path
func (lc *LinkConverter) getPhaseFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "phases" && i+1 < len(parts) {
			phaseName := parts[i+1]
			// Remove number prefix if present
			if idx := strings.Index(phaseName, "-"); idx > 0 {
				return phaseName[idx+1:]
			}
			return phaseName
		}
	}
	return ""
}

// getArtifactCategoryFromPath extracts artifact category from path
func (lc *LinkConverter) getArtifactCategoryFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "artifacts" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// createWikilink creates a wikilink from a target file and display text
func (lc *LinkConverter) createWikilink(targetFile *obsidian.MarkdownFile, displayText string) string {
	// Get the title from frontmatter or generate from path
	title := "Untitled"
	if targetFile.HasFrontmatter() && targetFile.Frontmatter.Title != "" {
		title = targetFile.Frontmatter.Title
	} else {
		title = obsidian.ExtractTitleFromPath(targetFile.Path)
	}

	// If display text matches title, use simple wikilink
	if displayText == title {
		return "[[" + title + "]]"
	}

	// Otherwise use alias format
	return "[[" + title + "|" + displayText + "]]"
}

// ParseWikilinks extracts all wikilinks from content
func (lc *LinkConverter) ParseWikilinks(content string) []*obsidian.ParsedLink {
	var links []*obsidian.ParsedLink

	// Regex for wikilinks: [[target#heading^blockid|alias]] or ![[embed]]
	wikilinkRegex := regexp.MustCompile(`(!?)\[\[([^|\]]+)(?:\|([^\]]+))?\]\]`)

	matches := wikilinkRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		isEmbed := match[1] == "!"
		fullTarget := match[2]
		alias := ""
		if len(match) > 3 {
			alias = match[3]
		}

		// Parse target for heading and block references
		target, heading, blockID := lc.parseTarget(fullTarget)

		link := &obsidian.ParsedLink{
			Original: match[0],
			Target:   target,
			Alias:    alias,
			Heading:  heading,
			BlockID:  blockID,
			IsEmbed:  isEmbed,
		}

		links = append(links, link)
	}

	return links
}

// ParseWikilinks is a package-level convenience function
func ParseWikilinks(content string) []*obsidian.ParsedLink {
	converter := NewLinkConverter()
	return converter.ParseWikilinks(content)
}

// ValidateWikilinks validates wikilinks and returns broken links
func (lc *LinkConverter) ValidateWikilinks(content string) []string {
	var brokenLinks []string
	links := lc.ParseWikilinks(content)

	for _, link := range links {
		// Check if target exists in file index
		if _, exists := lc.fileIndex[link.Target]; !exists {
			// Check path index
			found := false
			for path := range lc.pathIndex {
				if strings.Contains(path, link.Target) {
					found = true
					break
				}
			}
			if !found {
				brokenLinks = append(brokenLinks, link.Target)
			}
		}
	}

	return brokenLinks
}

// parseTarget parses a wikilink target to extract file, heading, and block references
func (lc *LinkConverter) parseTarget(target string) (file, heading, blockID string) {
	// Split on # for headings
	if idx := strings.Index(target, "#"); idx != -1 {
		file = target[:idx]
		rest := target[idx+1:]

		// Split on ^ for block IDs
		if blockIdx := strings.Index(rest, "^"); blockIdx != -1 {
			heading = rest[:blockIdx]
			blockID = rest[blockIdx+1:]
		} else {
			heading = rest
		}
	} else if idx := strings.Index(target, "^"); idx != -1 {
		// Only block ID, no heading
		file = target[:idx]
		blockID = target[idx+1:]
	} else {
		// Just a file reference
		file = target
	}

	return file, heading, blockID
}

// convertPhaseReferences converts plain text phase references to wikilinks
func (lc *LinkConverter) convertPhaseReferences(content string) string {
	// Common phase names - order matters, specific before general
	phases := []string{"Frame Phase", "Design Phase", "Test Phase", "Build Phase", "Deploy Phase", "Iterate Phase"}
	lowercase_phases := []string{"Frame phase", "Design phase", "Test phase", "Build phase", "Deploy phase", "Iterate phase"}

	// Skip if already in wikilink
	for i, phase := range phases {
		lowercase := lowercase_phases[i]

		// Convert "Frame Phase" → "[[Frame Phase]]"
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(phase) + `\b`)
		content = pattern.ReplaceAllStringFunc(content, func(match string) string {
			if lc.isInWikilink(content, match) {
				return match
			}
			return "[[" + phase + "]]"
		})

		// Convert "Frame phase" → "[[Frame Phase|Frame phase]]"
		pattern = regexp.MustCompile(`\b` + regexp.QuoteMeta(lowercase) + `\b`)
		content = pattern.ReplaceAllStringFunc(content, func(match string) string {
			if lc.isInWikilink(content, match) {
				return match
			}
			return "[[" + phase + "|" + lowercase + "]]"
		})
	}

	return content
}

// convertArtifactReferences converts plain text artifact references to wikilinks
func (lc *LinkConverter) convertArtifactReferences(content string) string {
	// Common artifact references
	artifacts := map[string]string{
		"Feature specification": "Feature Specification",
		"feature specification": "Feature Specification|feature specification",
		"feature spec":          "Feature Specification|feature spec",
		"Technical design":      "Technical Design",
		"technical design":      "Technical Design|technical design",
		"User stories":          "User Stories",
		"user stories":          "User Stories|user stories",
		"PRD":                   "Product Requirements Document|PRD",
		"Test specification":    "Test Specification",
		"test specification":    "Test Specification|test specification",
	}

	for text, wikilink := range artifacts {
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(text) + `\b`)
		content = pattern.ReplaceAllStringFunc(content, func(match string) string {
			if lc.isInWikilink(content, match) {
				return match
			}
			if strings.Contains(wikilink, "|") {
				return "[[" + wikilink + "]]"
			}
			return "[[" + wikilink + "]]"
		})
	}

	return content
}

// convertWorkflowReferences converts plain text workflow references to wikilinks
func (lc *LinkConverter) convertWorkflowReferences(content string) string {
	// Common workflow references
	workflows := map[string]string{
		"HELIX workflow":          "HELIX Workflow",
		"HELIX Workflow":          "HELIX Workflow",
		"TDD":                     "Test-Driven Development|TDD",
		"Test-Driven Development": "Test-Driven Development",
		"test-driven development": "Test-Driven Development|test-driven development",
	}

	for text, wikilink := range workflows {
		pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(text) + `\b`)
		content = pattern.ReplaceAllStringFunc(content, func(match string) string {
			if lc.isInWikilink(content, match) {
				return match
			}
			if strings.Contains(wikilink, "|") {
				return "[[" + wikilink + "]]"
			}
			return "[[" + wikilink + "]]"
		})
	}

	return content
}
