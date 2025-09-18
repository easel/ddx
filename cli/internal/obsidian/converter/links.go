package converter

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/easel/ddx/internal/obsidian"
)

// LinkConverter converts markdown links to Obsidian wikilinks
type LinkConverter struct {
	fileIndex map[string]string                 // path -> title mapping
	aliases   map[string]string                 // alias -> canonical name
	pathIndex map[string]*obsidian.MarkdownFile // for reverse lookup
}

// NewLinkConverter creates a new link converter
func NewLinkConverter() *LinkConverter {
	return &LinkConverter{
		fileIndex: make(map[string]string),
		aliases:   make(map[string]string),
		pathIndex: make(map[string]*obsidian.MarkdownFile),
	}
}

// BuildIndex builds the file index for link resolution
func (c *LinkConverter) BuildIndex(files []*obsidian.MarkdownFile) {
	for _, file := range files {
		c.pathIndex[file.Path] = file

		// Map file path to title
		if file.Frontmatter != nil {
			c.fileIndex[file.Path] = file.Frontmatter.Title

			// Register aliases
			for _, alias := range file.Frontmatter.Aliases {
				c.aliases[alias] = file.Frontmatter.Title
			}
		} else {
			// Fallback to path-based title
			c.fileIndex[file.Path] = obsidian.ExtractTitleFromPath(file.Path)
		}
	}
}

// ConvertContent converts all links in markdown content to wikilinks
func (c *LinkConverter) ConvertContent(content string) string {
	// Pattern for markdown links: [text](path)
	linkPattern := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	content = linkPattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := linkPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		linkText := parts[1]
		linkPath := parts[2]

		// Skip external links
		if strings.HasPrefix(linkPath, "http://") || strings.HasPrefix(linkPath, "https://") {
			return match
		}

		// Skip anchor-only links
		if strings.HasPrefix(linkPath, "#") {
			return match
		}

		// Skip email links
		if strings.HasPrefix(linkPath, "mailto:") {
			return match
		}

		// Convert relative path to wikilink
		return c.convertToWikilink(linkText, linkPath)
	})

	// Convert common phase references
	content = c.convertPhaseReferences(content)

	// Convert common artifact references
	content = c.convertArtifactReferences(content)

	// Convert HELIX workflow references
	content = c.convertWorkflowReferences(content)

	return content
}

// convertToWikilink converts a single link to wikilink format
func (c *LinkConverter) convertToWikilink(text, path string) string {
	// Clean up the path
	originalPath := path
	path = strings.TrimSuffix(path, ".md")

	// Try to resolve relative path
	resolvedPath := c.resolvePath(originalPath)
	if resolvedPath != "" {
		if title, ok := c.fileIndex[resolvedPath]; ok {
			if text != title {
				// Use alias syntax if link text differs from title
				return fmt.Sprintf("[[%s|%s]]", title, text)
			}
			return fmt.Sprintf("[[%s]]", title)
		}
	}

	// Try direct path lookup
	if title, ok := c.fileIndex[originalPath]; ok {
		if text != title {
			return fmt.Sprintf("[[%s|%s]]", title, text)
		}
		return fmt.Sprintf("[[%s]]", title)
	}

	// Extract just the filename for common patterns
	filename := filepath.Base(path)

	// Handle common HELIX patterns
	if wikilink := c.handleCommonPatterns(text, path, filename); wikilink != "" {
		return wikilink
	}

	// Default: create wikilink with the text
	return fmt.Sprintf("[[%s]]", text)
}

// resolvePath attempts to resolve a relative path to an absolute path
func (c *LinkConverter) resolvePath(relativePath string) string {
	// Simple implementation - in practice you might want more sophisticated path resolution
	for fullPath := range c.fileIndex {
		if strings.HasSuffix(fullPath, relativePath) {
			return fullPath
		}
		if strings.Contains(fullPath, relativePath) {
			return fullPath
		}
	}
	return ""
}

// handleCommonPatterns handles common HELIX file patterns
func (c *LinkConverter) handleCommonPatterns(text, path, filename string) string {
	// Remove extension for cleaner matching
	baseFilename := strings.TrimSuffix(filename, ".md")

	switch baseFilename {
	case "README":
		// Try to determine phase from path
		if strings.Contains(path, "/phases/") {
			phase := obsidian.GetPhaseFromPath(path)
			if phase != "" {
				return fmt.Sprintf("[[%s Phase]]", strings.Title(phase))
			}
		}
		// Try to determine artifact from path
		if strings.Contains(path, "/artifacts/") {
			artifact := obsidian.GetArtifactCategory(path)
			if artifact != "" {
				artifactName := strings.Title(strings.ReplaceAll(artifact, "-", " "))
				return fmt.Sprintf("[[%s]]", artifactName)
			}
		}

	case "template":
		artifact := extractArtifactFromPath(path)
		if artifact != "" {
			title := fmt.Sprintf("%s Template", strings.Title(strings.ReplaceAll(artifact, "-", " ")))
			// Use alias only if the text is more descriptive than just the base filename
			if text != title && text != baseFilename {
				return fmt.Sprintf("[[%s|%s]]", title, text)
			}
			return fmt.Sprintf("[[%s]]", title)
		}

	case "prompt":
		artifact := extractArtifactFromPath(path)
		if artifact != "" {
			title := fmt.Sprintf("%s Prompt", strings.Title(strings.ReplaceAll(artifact, "-", " ")))
			if text != title && text != baseFilename {
				return fmt.Sprintf("[[%s|%s]]", title, text)
			}
			return fmt.Sprintf("[[%s]]", title)
		}

	case "example":
		artifact := extractArtifactFromPath(path)
		if artifact != "" {
			title := fmt.Sprintf("%s Example", strings.Title(strings.ReplaceAll(artifact, "-", " ")))
			if text != title && text != baseFilename {
				return fmt.Sprintf("[[%s|%s]]", title, text)
			}
			return fmt.Sprintf("[[%s]]", title)
		}

	case "enforcer":
		phase := obsidian.GetPhaseFromPath(path)
		if phase != "" {
			title := fmt.Sprintf("%s Phase Enforcer", strings.Title(phase))
			if text != title && text != baseFilename {
				return fmt.Sprintf("[[%s|%s]]", title, text)
			}
			return fmt.Sprintf("[[%s]]", title)
		}
	}

	// Check for special files
	if strings.Contains(filename, "coordinator") {
		return "[[HELIX Workflow Coordinator]]"
	}

	if strings.Contains(filename, "principle") {
		return "[[HELIX Principles]]"
	}

	// Check for feature files
	if strings.Contains(filename, "FEAT-") {
		// Extract feature number
		re := regexp.MustCompile(`FEAT-(\d+)`)
		matches := re.FindStringSubmatch(filename)
		if len(matches) > 1 {
			return fmt.Sprintf("[[FEAT-%s]]", matches[1])
		}
	}

	return ""
}

// convertPhaseReferences converts common phase references to wikilinks
func (c *LinkConverter) convertPhaseReferences(content string) string {
	phases := []string{"Frame", "Design", "Test", "Build", "Deploy", "Iterate"}

	for _, phase := range phases {
		// Convert "Frame phase" -> "[[Frame Phase]]" (case insensitive)
		patterns := []string{
			fmt.Sprintf(`(?i)\b%s phase\b`, phase),
			fmt.Sprintf(`(?i)\b%s Phase\b`, phase),
		}

		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			content = re.ReplaceAllStringFunc(content, func(match string) string {
				// Don't replace if already in wikilink
				if c.isInWikilink(content, match) {
					return match
				}
				return fmt.Sprintf("[[%s Phase]]", phase)
			})
		}

		// Convert standalone phase names when they clearly refer to phases
		// Note: Go doesn't support lookahead/lookbehind, so we use simpler patterns
		contextPatterns := []string{
			fmt.Sprintf(`\bthe %s\b`, phase),
			fmt.Sprintf(`\bto %s\b`, phase),
			fmt.Sprintf(`\b%s phase\b`, strings.ToLower(phase)),
		}

		for _, pattern := range contextPatterns {
			re := regexp.MustCompile(pattern)
			content = re.ReplaceAllStringFunc(content, func(match string) string {
				// Check if it's already in a wikilink
				if c.isInWikilink(content, match) {
					return match
				}

				// Extract just the phase name from the match
				phaseName := strings.TrimPrefix(match, "the ")
				phaseName = strings.TrimPrefix(phaseName, "to ")
				phaseName = strings.TrimSuffix(phaseName, " phase")
				phaseName = strings.Title(phaseName)

				return strings.Replace(match, strings.Title(phase), fmt.Sprintf("[[%s Phase|%s]]", phase, strings.Title(phase)), 1)
			})
		}
	}

	return content
}

// convertArtifactReferences converts common artifact references
func (c *LinkConverter) convertArtifactReferences(content string) string {
	artifacts := map[string]string{
		"feature specification":   "Feature Specification",
		"feature spec":            "Feature Specification",
		"technical design":        "Technical Design",
		"implementation guide":    "Implementation Guide",
		"test specification":      "Test Specification",
		"user stories":            "User Stories",
		"product requirements":    "Product Requirements",
		"PRD":                     "Product Requirements Document",
		"risk register":           "Risk Register",
		"feasibility study":       "Feasibility Study",
		"research plan":           "Research Plan",
		"compliance requirements": "Compliance Requirements",
	}

	for pattern, replacement := range artifacts {
		// Create case-insensitive pattern
		re := regexp.MustCompile(fmt.Sprintf(`(?i)\b%s\b`, regexp.QuoteMeta(pattern)))
		content = re.ReplaceAllStringFunc(content, func(match string) string {
			// Don't replace if already in wikilink
			if c.isInWikilink(content, match) {
				return match
			}

			// Use original case for the alias if different from replacement
			if strings.ToLower(match) != strings.ToLower(replacement) {
				return fmt.Sprintf("[[%s|%s]]", replacement, match)
			}
			return fmt.Sprintf("[[%s]]", replacement)
		})
	}

	return content
}

// convertWorkflowReferences converts HELIX workflow references
func (c *LinkConverter) convertWorkflowReferences(content string) string {
	// Convert "HELIX workflow" to "[[HELIX Workflow]]"
	re := regexp.MustCompile(`(?i)\bHELIX workflow\b`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Don't replace if already in wikilink
		if c.isInWikilink(content, match) {
			return match
		}

		return "[[HELIX Workflow]]"
	})

	// Convert "TDD" or "Test-Driven Development" references
	patterns := []string{
		`\bTDD\b`,
		`\bTest-Driven Development\b`,
		`\btest-driven development\b`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		content = re.ReplaceAllStringFunc(content, func(match string) string {
			if c.isInWikilink(content, match) {
				return match
			}

			if match == "TDD" {
				return "[[Test-Driven Development|TDD]]"
			}
			return "[[Test-Driven Development]]"
		})
	}

	return content
}

// extractArtifactFromPath extracts artifact name from path
func extractArtifactFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "artifacts" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

// isInWikilink checks if a match is already inside a wikilink
func (c *LinkConverter) isInWikilink(content, match string) bool {
	idx := strings.Index(content, match)
	if idx == -1 {
		return false
	}

	// Find the nearest [[ before the match
	beforeContent := content[:idx]
	lastOpenIdx := strings.LastIndex(beforeContent, "[[")

	// Find the nearest ]] after the match
	afterContent := content[idx+len(match):]
	nextCloseIdx := strings.Index(afterContent, "]]")

	// If we found both [[ before and ]] after, check if they form a valid wikilink
	if lastOpenIdx != -1 && nextCloseIdx != -1 {
		// Check if there's a closing ]] between the opening [[ and our match
		betweenContent := content[lastOpenIdx+2 : idx]
		if !strings.Contains(betweenContent, "]]") {
			// We're inside a wikilink
			return true
		}
	}

	return false
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ParseWikilinks extracts all wikilinks from content
func ParseWikilinks(content string) []*obsidian.ParsedLink {
	var links []*obsidian.ParsedLink

	// Pattern for wikilinks: [[target|alias]] or [[target#heading]] or [[target^blockid]] or ![[embed]]
	re := regexp.MustCompile(`(!?)\[\[([^\]]+)\]\]`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		isEmbed := match[1] == "!"
		linkContent := match[2]

		link := &obsidian.ParsedLink{
			Original: match[0],
			IsEmbed:  isEmbed,
		}

		// Parse target|alias
		if idx := strings.Index(linkContent, "|"); idx != -1 {
			link.Target = linkContent[:idx]
			link.Alias = linkContent[idx+1:]
		} else {
			link.Target = linkContent
		}

		// Parse target#heading
		if idx := strings.Index(link.Target, "#"); idx != -1 {
			link.Heading = link.Target[idx+1:]
			link.Target = link.Target[:idx]
		}

		// Parse target^blockid
		if idx := strings.Index(link.Target, "^"); idx != -1 {
			link.BlockID = link.Target[idx+1:]
			link.Target = link.Target[:idx]
		}

		links = append(links, link)
	}

	return links
}

// ValidateWikilinks checks if all wikilinks in content are valid
func (c *LinkConverter) ValidateWikilinks(content string) []string {
	var brokenLinks []string

	links := ParseWikilinks(content)
	for _, link := range links {
		// Check if the target exists in our file index
		found := false
		for _, title := range c.fileIndex {
			if title == link.Target {
				found = true
				break
			}
		}

		// Check aliases
		if !found {
			if _, exists := c.aliases[link.Target]; exists {
				found = true
			}
		}

		if !found {
			brokenLinks = append(brokenLinks, link.Target)
		}
	}

	return brokenLinks
}
