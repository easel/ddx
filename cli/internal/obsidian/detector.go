package obsidian

import (
	"path/filepath"
	"regexp"
	"strings"
)

// FileTypeDetector detects the type of markdown files
type FileTypeDetector struct {
	patterns map[string]FileType
}

// NewFileTypeDetector creates a new file type detector
func NewFileTypeDetector() *FileTypeDetector {
	return &FileTypeDetector{
		patterns: make(map[string]FileType),
	}
}

// Detect determines the file type based on path patterns
func (d *FileTypeDetector) Detect(path string) FileType {
	// Normalize path
	path = filepath.ToSlash(path)
	filename := filepath.Base(path)

	// Check for feature specification pattern (generic pattern)
	if matched, _ := regexp.MatchString(`FEAT-\d+`, filename); matched {
		return FileTypeFeature
	}

	// Check exact filename matches (generic patterns)
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
		// Determine type based on directory context
		if strings.Contains(path, "/phases/") {
			return FileTypePhase
		}
		if strings.Contains(path, "/artifacts/") {
			return FileTypeArtifact
		}
	}

	return FileTypeUnknown
}

// GetPhaseFromPath extracts the phase name from a file path (generic implementation)
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
		// Check for numbered directory patterns
		if matched, _ := regexp.MatchString(`^\d+-`, part); matched {
			if idx := strings.Index(part, "-"); idx > 0 {
				return part[idx+1:]
			}
		}
		// Check for common phase names as directory names
		if part == "frame" || part == "design" || part == "test" ||
			part == "build" || part == "deploy" || part == "iterate" {
			return part
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

	// For feature files, return feature-specification
	filename := filepath.Base(path)
	if strings.Contains(filename, "FEAT-") {
		return "feature-specification"
	}

	return ""
}

// GetPhaseNumber returns the phase number for a given phase name
// This could be made configurable, but for now uses common patterns
func GetPhaseNumber(phaseName string) int {
	// Common phase numbering - could be made configurable
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
	// Common phase order - could be made configurable
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

// GetComplexityFromPath estimates complexity based on path structure
func GetComplexityFromPath(path string) string {
	// Simple heuristic based on path patterns
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
