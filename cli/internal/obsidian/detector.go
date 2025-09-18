package obsidian

import (
	"path/filepath"
	"regexp"
	"strings"
)

// FileTypeDetector detects the type of a HELIX markdown file
type FileTypeDetector struct {
	patterns map[string]FileType
}

// NewFileTypeDetector creates a new file type detector
func NewFileTypeDetector() *FileTypeDetector {
	return &FileTypeDetector{
		patterns: map[string]FileType{
			"phases/*/README.md":    FileTypePhase,
			"phases/*/enforcer.md":  FileTypeEnforcer,
			"*/template.md":         FileTypeTemplate,
			"*/prompt.md":           FileTypePrompt,
			"*/example.md":          FileTypeExample,
			"coordinator.md":        FileTypeCoordinator,
			"principles.md":         FileTypePrinciple,
			"artifacts/*/README.md": FileTypeArtifact,
		},
	}
}

// Detect determines the file type based on path patterns and content
func (d *FileTypeDetector) Detect(path string) FileType {
	// Normalize path
	path = filepath.ToSlash(path)
	filename := filepath.Base(path)

	// Check for HELIX document structure patterns first
	if fileType := d.detectFromContent(path); fileType != FileTypeUnknown {
		return fileType
	}

	// Check exact filename matches
	switch filename {
	case "coordinator.md":
		return FileTypeCoordinator
	case "principles.md", "principle.md":
		return FileTypePrinciple
	case "enforcer.md":
		return FileTypeEnforcer
	case "template.md":
		return FileTypeTemplate
	case "prompt.md":
		return FileTypePrompt
	case "example.md":
		return FileTypeExample
	case "README.md":
		// README could be phase or artifact based on location
		if strings.Contains(path, "/phases/") {
			return FileTypePhase
		}
		if strings.Contains(path, "/artifacts/") {
			return FileTypeArtifact
		}
	}

	// Check directory-based patterns
	if strings.Contains(path, "/phases/") {
		if strings.HasSuffix(path, "/README.md") {
			return FileTypePhase
		}
		if strings.HasSuffix(path, "/enforcer.md") {
			return FileTypeEnforcer
		}
		if strings.Contains(path, "/artifacts/") {
			switch filename {
			case "template.md":
				return FileTypeTemplate
			case "prompt.md":
				return FileTypePrompt
			case "example.md":
				return FileTypeExample
			}
		}
	}

	// Check docs directory patterns (for HELIX documents)
	if strings.Contains(path, "/docs/") {
		if strings.Contains(path, "/01-frame/") || strings.Contains(path, "/frame/") {
			if strings.Contains(filename, "FEAT-") {
				return FileTypeFeature
			}
		}
		// Note: Phase-based docs don't get special file types - they use artifact types
	}

	// Use pattern matching as fallback
	for pattern, fileType := range d.patterns {
		if matched, _ := filepath.Match(pattern, path); matched {
			return fileType
		}
	}

	return FileTypeUnknown
}

// detectFromContent analyzes file content to determine type
func (d *FileTypeDetector) detectFromContent(path string) FileType {
	// This would ideally read the file content, but for now we'll use path-based detection
	// In a real implementation, you might read the first few lines to check for specific patterns

	filename := filepath.Base(path)

	// Check for feature specification pattern
	if matched, _ := regexp.MatchString(`FEAT-\d+`, filename); matched {
		return FileTypeFeature
	}

	return FileTypeUnknown
}

// GetPhaseFromPath extracts the phase name from a file path
func GetPhaseFromPath(path string) string {
	path = filepath.ToSlash(path)
	parts := strings.Split(path, "/")

	for i, part := range parts {
		if part == "phases" && i+1 < len(parts) {
			phaseName := parts[i+1]
			// Remove number prefix if present (e.g., "01-frame" -> "frame")
			if idx := strings.Index(phaseName, "-"); idx > 0 {
				return phaseName[idx+1:]
			}
			return phaseName
		}
		// Also check for docs structure
		if strings.HasPrefix(part, "01-") || part == "frame" {
			return "frame"
		}
		if strings.HasPrefix(part, "02-") || part == "design" {
			return "design"
		}
		if strings.HasPrefix(part, "03-") || part == "test" {
			return "test"
		}
		if strings.HasPrefix(part, "04-") || part == "build" {
			return "build"
		}
		if strings.HasPrefix(part, "05-") || part == "deploy" {
			return "deploy"
		}
		if strings.HasPrefix(part, "06-") || part == "iterate" {
			return "iterate"
		}
	}
	return ""
}

// GetArtifactCategory extracts the artifact category from a file path
func GetArtifactCategory(path string) string {
	path = filepath.ToSlash(path)
	parts := strings.Split(path, "/")

	for i, part := range parts {
		if part == "artifacts" && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	// For docs structure, try to infer from filename
	filename := filepath.Base(path)

	// Check for FEAT pattern - these are always feature specifications
	if strings.Contains(filename, "FEAT-") {
		return "feature-specification"
	}

	return ""
}

// GetPhaseNumber returns the phase number for a given phase name
func GetPhaseNumber(phaseName string) int {
	phaseNumbers := map[string]int{
		"frame":   1,
		"design":  2,
		"test":    3,
		"build":   4,
		"deploy":  5,
		"iterate": 6,
	}

	if num, ok := phaseNumbers[phaseName]; ok {
		return num
	}
	return 0
}

// GetNextPhase returns the next phase name
func GetNextPhase(phaseName string) string {
	phaseOrder := []string{"frame", "design", "test", "build", "deploy", "iterate"}

	for i, phase := range phaseOrder {
		if phase == phaseName && i < len(phaseOrder)-1 {
			return phaseOrder[i+1]
		}
	}
	return ""
}

// GetPreviousPhase returns the previous phase name
func GetPreviousPhase(phaseName string) string {
	phaseOrder := []string{"frame", "design", "test", "build", "deploy", "iterate"}

	for i, phase := range phaseOrder {
		if phase == phaseName && i > 0 {
			return phaseOrder[i-1]
		}
	}
	return ""
}

// ExtractTitleFromPath generates a reasonable title from the file path
func ExtractTitleFromPath(path string) string {
	filename := filepath.Base(path)
	filename = strings.TrimSuffix(filename, ".md")

	// Handle special cases
	if filename == "README" {
		if phase := GetPhaseFromPath(path); phase != "" {
			return strings.Title(phase) + " Phase"
		}
		if artifact := GetArtifactCategory(path); artifact != "" {
			return strings.Title(strings.ReplaceAll(artifact, "-", " "))
		}
		return "README"
	}

	if filename == "enforcer" {
		if phase := GetPhaseFromPath(path); phase != "" {
			return strings.Title(phase) + " Phase Enforcer"
		}
		return "Phase Enforcer"
	}

	if filename == "template" {
		if artifact := GetArtifactCategory(path); artifact != "" {
			return strings.Title(strings.ReplaceAll(artifact, "-", " ")) + " Template"
		}
		return "Template"
	}

	if filename == "prompt" {
		if artifact := GetArtifactCategory(path); artifact != "" {
			return strings.Title(strings.ReplaceAll(artifact, "-", " ")) + " Prompt"
		}
		return "Prompt"
	}

	if filename == "example" {
		if artifact := GetArtifactCategory(path); artifact != "" {
			return strings.Title(strings.ReplaceAll(artifact, "-", " ")) + " Example"
		}
		return "Example"
	}

	// Handle feature files
	if strings.Contains(filename, "FEAT-") {
		// Extract feature name from filename
		parts := strings.Split(filename, "-")
		if len(parts) > 2 {
			name := strings.Join(parts[2:], " ")
			return strings.Title(name)
		}
		return "Feature Specification"
	}

	// Default: clean up the filename
	title := strings.ReplaceAll(filename, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	return strings.Title(title)
}

// IsHelixFile determines if a file is part of the HELIX workflow
func IsHelixFile(path string) bool {
	path = filepath.ToSlash(path)

	// Check if it's in a HELIX directory
	if strings.Contains(path, "/workflows/helix/") ||
		strings.Contains(path, "/helix/") {
		return true
	}

	// Check if it's in docs with HELIX structure
	if strings.Contains(path, "/docs/") &&
		(strings.Contains(path, "/01-frame/") ||
			strings.Contains(path, "/02-design/") ||
			strings.Contains(path, "/03-test/") ||
			strings.Contains(path, "/04-build/") ||
			strings.Contains(path, "/05-deploy/") ||
			strings.Contains(path, "/06-iterate/")) {
		return true
	}

	return false
}

// GetComplexityFromPath estimates complexity based on path structure
func GetComplexityFromPath(path string) string {
	// This is a simple heuristic - in practice you might want to analyze content
	if strings.Contains(path, "/example") ||
		strings.Contains(path, "/simple") {
		return "simple"
	}

	if strings.Contains(path, "/advanced") ||
		strings.Contains(path, "/complex") {
		return "complex"
	}

	// Default to moderate
	return "moderate"
}
