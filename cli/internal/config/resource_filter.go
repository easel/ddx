package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ResourceFilterEngine handles resource selection logic
type ResourceFilterEngine struct {
	config *Config
	libPath string
}

// NewResourceFilterEngine creates a new resource filter engine
func NewResourceFilterEngine(config *Config, libraryPath string) *ResourceFilterEngine {
	return &ResourceFilterEngine{
		config:  config,
		libPath: libraryPath,
	}
}

// FilterResources filters resources based on configuration
func (rfe *ResourceFilterEngine) FilterResources(resourceType string, resourcePaths []string) ([]string, error) {
	if rfe.config.Resources == nil {
		// No resource filtering configured - return all
		return resourcePaths, nil
	}

	filter := rfe.getFilterForType(resourceType)
	if filter == nil {
		// No filter for this resource type - return all
		return resourcePaths, nil
	}

	var filtered []string

	for _, resourcePath := range resourcePaths {
		// Remove library path prefix to get relative path
		relPath := strings.TrimPrefix(resourcePath, rfe.libPath)
		relPath = strings.TrimPrefix(relPath, "/")
		relPath = strings.TrimPrefix(relPath, resourceType+"/")

		included := rfe.matchesIncludePatterns(relPath, filter.Include)
		excluded := rfe.matchesExcludePatterns(relPath, filter.Exclude)


		// Include if matches include patterns and doesn't match exclude patterns
		if included && !excluded {
			filtered = append(filtered, resourcePath)
		}
	}

	return filtered, nil
}

// DiscoverResourcesWithFilter discovers and filters resources of a specific type
func (rfe *ResourceFilterEngine) DiscoverResourcesWithFilter(resourceType string) ([]string, error) {
	resourceDir := filepath.Join(rfe.libPath, resourceType)

	if _, err := os.Stat(resourceDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Discover all resources
	var allResources []string
	err := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory
		if path == resourceDir {
			return nil
		}

		// For now, only include files in filtering, not directories
		// Directories will be included by the list command if they contain included files
		if !info.IsDir() {
			allResources = append(allResources, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover resources in %s: %w", resourceDir, err)
	}

	// Apply filtering
	return rfe.FilterResources(resourceType, allResources)
}

// ValidateResourceConfiguration validates resource selection configuration
func (rfe *ResourceFilterEngine) ValidateResourceConfiguration() []string {
	var warnings []string

	if rfe.config.Resources == nil {
		return warnings
	}

	resourceTypes := map[string]*ResourceFilter{
		"prompts":   rfe.config.Resources.Prompts,
		"templates": rfe.config.Resources.Templates,
		"patterns":  rfe.config.Resources.Patterns,
		"configs":   rfe.config.Resources.Configs,
		"scripts":   rfe.config.Resources.Scripts,
		"workflows": rfe.config.Resources.Workflows,
	}

	for resourceType, filter := range resourceTypes {
		if filter == nil {
			continue
		}

		warnings = append(warnings, rfe.validateResourceFilter(resourceType, filter)...)
	}

	return warnings
}

// PreviewResourceSelection shows what resources would be selected
func (rfe *ResourceFilterEngine) PreviewResourceSelection() (map[string][]string, map[string][]string, error) {
	included := make(map[string][]string)
	excluded := make(map[string][]string)

	if rfe.config.Resources == nil {
		return included, excluded, nil
	}

	resourceTypes := []string{"prompts", "templates", "patterns", "configs", "scripts", "workflows"}

	for _, resourceType := range resourceTypes {
		filter := rfe.getFilterForType(resourceType)
		if filter == nil {
			continue
		}

		// Discover all available resources
		resourceDir := filepath.Join(rfe.libPath, resourceType)
		if _, err := os.Stat(resourceDir); os.IsNotExist(err) {
			continue
		}

		var allResources []string
		err := filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == resourceDir {
				return nil
			}

			// Get relative path
			relPath := strings.TrimPrefix(path, resourceDir+"/")
			allResources = append(allResources, relPath)
			return nil
		})

		if err != nil {
			continue
		}

		// Categorize resources
		for _, relPath := range allResources {
			isIncluded := rfe.matchesIncludePatterns(relPath, filter.Include)
			isExcluded := rfe.matchesExcludePatterns(relPath, filter.Exclude)

			if isIncluded && !isExcluded {
				included[resourceType] = append(included[resourceType], relPath)
			} else if isExcluded {
				excluded[resourceType] = append(excluded[resourceType], relPath)
			}
		}
	}

	return included, excluded, nil
}

// getFilterForType returns the ResourceFilter for a given resource type
func (rfe *ResourceFilterEngine) getFilterForType(resourceType string) *ResourceFilter {
	if rfe.config.Resources == nil {
		return nil
	}

	switch resourceType {
	case "prompts":
		return rfe.config.Resources.Prompts
	case "templates":
		return rfe.config.Resources.Templates
	case "patterns":
		return rfe.config.Resources.Patterns
	case "configs":
		return rfe.config.Resources.Configs
	case "scripts":
		return rfe.config.Resources.Scripts
	case "workflows":
		return rfe.config.Resources.Workflows
	default:
		return nil
	}
}

// matchesIncludePatterns checks if a path matches any include pattern
func (rfe *ResourceFilterEngine) matchesIncludePatterns(path string, includePatterns []string) bool {
	if len(includePatterns) == 0 {
		// No include patterns means include all
		return true
	}

	for _, pattern := range includePatterns {
		if rfe.matchesPattern(path, pattern) {
			return true
		}
	}
	return false
}

// matchesExcludePatterns checks if a path matches any exclude pattern
func (rfe *ResourceFilterEngine) matchesExcludePatterns(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		if rfe.matchesPattern(path, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a path matches a wildcard pattern
func (rfe *ResourceFilterEngine) matchesPattern(path, pattern string) bool {
	// Handle exact matches first
	if path == pattern {
		return true
	}

	// Convert wildcard pattern to regex
	regexPattern := rfe.wildcardToRegex(pattern)
	matched, err := regexp.MatchString(regexPattern, path)
	if err != nil {
		// If regex fails, fall back to simple string matching
		return strings.Contains(path, pattern)
	}

	return matched
}

// wildcardToRegex converts shell-style wildcards to regex
func (rfe *ResourceFilterEngine) wildcardToRegex(pattern string) string {
	// Escape special regex characters except wildcards
	escaped := regexp.QuoteMeta(pattern)

	// Replace escaped wildcards with regex equivalents
	escaped = strings.ReplaceAll(escaped, `\*\*`, `.*`)     // ** matches everything including /
	escaped = strings.ReplaceAll(escaped, `\*`, `[^/]*`)   // * matches anything except /
	escaped = strings.ReplaceAll(escaped, `\?`, `[^/]`)    // ? matches single character except /

	// Handle character classes [abc]
	escaped = strings.ReplaceAll(escaped, `\[`, `[`)
	escaped = strings.ReplaceAll(escaped, `\]`, `]`)

	// Handle brace expansion {option1,option2}
	braceRegex := regexp.MustCompile(`\\\{([^}]+)\\\}`)
	escaped = braceRegex.ReplaceAllStringFunc(escaped, func(match string) string {
		// Extract options from {option1,option2}
		content := match[2 : len(match)-2] // Remove \{ and \}
		options := strings.Split(content, ",")
		for i, opt := range options {
			options[i] = strings.TrimSpace(opt)
		}
		return "(" + strings.Join(options, "|") + ")"
	})

	// Anchor the pattern to match the full path
	return "^" + escaped + "$"
}

// validateResourceFilter validates a resource filter and returns warnings
func (rfe *ResourceFilterEngine) validateResourceFilter(resourceType string, filter *ResourceFilter) []string {
	var warnings []string

	resourceDir := filepath.Join(rfe.libPath, resourceType)
	if _, err := os.Stat(resourceDir); os.IsNotExist(err) {
		warnings = append(warnings, fmt.Sprintf("Resource directory '%s' does not exist", resourceType))
		return warnings
	}

	// Check include patterns
	for _, pattern := range filter.Include {
		if !rfe.patternHasMatches(resourceType, pattern) {
			warnings = append(warnings, fmt.Sprintf("Include pattern '%s' in %s matches no resources", pattern, resourceType))
		}
	}

	// Check for conflicting patterns
	for _, includePattern := range filter.Include {
		for _, excludePattern := range filter.Exclude {
			if rfe.patternsConflict(includePattern, excludePattern) {
				warnings = append(warnings, fmt.Sprintf("Conflicting patterns in %s: include '%s' conflicts with exclude '%s'", resourceType, includePattern, excludePattern))
			}
		}
	}

	return warnings
}

// patternHasMatches checks if a pattern matches any existing resources
func (rfe *ResourceFilterEngine) patternHasMatches(resourceType, pattern string) bool {
	resourceDir := filepath.Join(rfe.libPath, resourceType)

	hasMatches := false
	filepath.Walk(resourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == resourceDir {
			return nil
		}

		relPath := strings.TrimPrefix(path, resourceDir+"/")
		if rfe.matchesPattern(relPath, pattern) {
			hasMatches = true
			return filepath.SkipAll // Stop walking
		}
		return nil
	})

	return hasMatches
}

// patternsConflict checks if an include and exclude pattern would conflict
func (rfe *ResourceFilterEngine) patternsConflict(includePattern, excludePattern string) bool {
	// Simple heuristic: if exclude pattern is more specific than include pattern,
	// they might conflict. This is a basic implementation.

	// If exclude is a substring of include, they likely conflict
	if strings.Contains(includePattern, excludePattern) {
		return true
	}

	// If include and exclude are identical, they definitely conflict
	if includePattern == excludePattern {
		return true
	}

	return false
}