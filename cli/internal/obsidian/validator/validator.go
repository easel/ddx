package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/easel/ddx/internal/obsidian"
)

// ValidationError represents a validation error
type ValidationError struct {
	File     string
	Field    string
	Message  string
	Severity Severity
}

type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInfo
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "ERROR"
	case SeverityWarning:
		return "WARNING"
	case SeverityInfo:
		return "INFO"
	default:
		return "UNKNOWN"
	}
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s - %s", e.Severity, e.File, e.Field, e.Message)
}

// Validator validates Obsidian format in markdown files
type Validator struct {
	requiredFields map[obsidian.FileType][]string
	validTags      map[string]bool
	rules          []ValidationRule
}

// ValidationRule represents a validation rule
type ValidationRule interface {
	Validate(file *obsidian.MarkdownFile) []ValidationError
	Name() string
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	v := &Validator{
		requiredFields: initializeRequiredFields(),
		validTags:      initializeValidTags(),
		rules:          []ValidationRule{},
	}

	// Add built-in validation rules
	v.addRule(&FrontmatterRule{})
	v.addRule(&TagRule{v.validTags})
	v.addRule(&WikilinkRule{})
	v.addRule(&PhaseRule{})
	v.addRule(&FeatureRule{})

	return v
}

// addRule adds a validation rule
func (v *Validator) addRule(rule ValidationRule) {
	v.rules = append(v.rules, rule)
}

// ValidateFile validates a single markdown file
func (v *Validator) ValidateFile(file *obsidian.MarkdownFile) []ValidationError {
	var errors []ValidationError

	// Run all validation rules
	for _, rule := range v.rules {
		ruleErrors := rule.Validate(file)
		errors = append(errors, ruleErrors...)
	}

	return errors
}

// ValidateFiles validates multiple files and returns a summary
func (v *Validator) ValidateFiles(files []*obsidian.MarkdownFile) ValidationSummary {
	summary := ValidationSummary{
		TotalFiles:     len(files),
		ValidatedFiles: 0,
		Errors:         []ValidationError{},
		Warnings:       []ValidationError{},
		Info:           []ValidationError{},
	}

	for _, file := range files {
		if file.FileType == obsidian.FileTypeUnknown {
			continue
		}

		summary.ValidatedFiles++
		errors := v.ValidateFile(file)

		for _, err := range errors {
			switch err.Severity {
			case SeverityError:
				summary.Errors = append(summary.Errors, err)
			case SeverityWarning:
				summary.Warnings = append(summary.Warnings, err)
			case SeverityInfo:
				summary.Info = append(summary.Info, err)
			}
		}
	}

	return summary
}

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalFiles     int
	ValidatedFiles int
	Errors         []ValidationError
	Warnings       []ValidationError
	Info           []ValidationError
}

// HasErrors returns true if there are any errors
func (s ValidationSummary) HasErrors() bool {
	return len(s.Errors) > 0
}

// HasWarnings returns true if there are any warnings
func (s ValidationSummary) HasWarnings() bool {
	return len(s.Warnings) > 0
}

// FrontmatterRule validates frontmatter presence and required fields
type FrontmatterRule struct{}

func (r *FrontmatterRule) Name() string {
	return "frontmatter"
}

func (r *FrontmatterRule) Validate(file *obsidian.MarkdownFile) []ValidationError {
	var errors []ValidationError

	// Check for frontmatter presence
	if file.Frontmatter == nil {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "frontmatter",
			Message:  "missing frontmatter",
			Severity: SeverityError,
		})
		return errors
	}

	fm := file.Frontmatter

	// Check required fields based on file type
	requiredFields := getRequiredFields(file.FileType)
	for _, field := range requiredFields {
		if !hasField(fm, field) {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    field,
				Message:  "required field missing",
				Severity: SeverityError,
			})
		}
	}

	// Check field formats
	if fm.Title == "" {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "title",
			Message:  "title cannot be empty",
			Severity: SeverityError,
		})
	}

	if fm.Type == "" {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "type",
			Message:  "type cannot be empty",
			Severity: SeverityError,
		})
	}

	// Check dates
	if fm.Created.IsZero() {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "created",
			Message:  "created date missing or invalid",
			Severity: SeverityWarning,
		})
	}

	if fm.Updated.IsZero() {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "updated",
			Message:  "updated date missing or invalid",
			Severity: SeverityWarning,
		})
	}

	return errors
}

// TagRule validates tag format and hierarchy
type TagRule struct {
	validTags map[string]bool
}

func (r *TagRule) Name() string {
	return "tags"
}

func (r *TagRule) Validate(file *obsidian.MarkdownFile) []ValidationError {
	var errors []ValidationError

	if file.Frontmatter == nil {
		return errors
	}

	// Check if at least one tag exists
	if len(file.Frontmatter.Tags) == 0 {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "tags",
			Message:  "no tags specified",
			Severity: SeverityWarning,
		})
		return errors
	}

	// Check tag format
	for _, tag := range file.Frontmatter.Tags {
		if !r.isValidTag(tag) {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "tags",
				Message:  fmt.Sprintf("invalid tag format: %s", tag),
				Severity: SeverityError,
			})
		}
	}

	// Check for required base tag
	hasHelixTag := false
	for _, tag := range file.Frontmatter.Tags {
		if tag == "helix" || strings.HasPrefix(tag, "helix/") {
			hasHelixTag = true
			break
		}
	}

	if !hasHelixTag {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "tags",
			Message:  "missing required 'helix' tag",
			Severity: SeverityError,
		})
	}

	return errors
}

func (r *TagRule) isValidTag(tag string) bool {
	// Must start with "helix"
	if !strings.HasPrefix(tag, "helix") {
		return false
	}

	// Check against known valid tags
	if r.validTags[tag] {
		return true
	}

	// Allow project-specific tags with prefix
	if strings.HasPrefix(tag, "helix/project/") {
		return true
	}

	// Allow custom tags with valid format
	if matched, _ := regexp.MatchString(`^helix(/[a-z0-9-]+)+$`, tag); matched {
		return true
	}

	return false
}

// WikilinkRule validates wikilink format and targets
type WikilinkRule struct{}

func (r *WikilinkRule) Name() string {
	return "wikilinks"
}

func (r *WikilinkRule) Validate(file *obsidian.MarkdownFile) []ValidationError {
	var errors []ValidationError

	// Extract wikilinks from content
	links := extractWikilinks(file.Content)

	for _, link := range links {
		// Check for empty links
		if link == "" {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "content",
				Message:  "empty wikilink found",
				Severity: SeverityError,
			})
			continue
		}

		// Check for malformed links
		if strings.Contains(link, "[[") || strings.Contains(link, "]]") {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "content",
				Message:  fmt.Sprintf("malformed wikilink: %s", link),
				Severity: SeverityError,
			})
		}

		// Check for suspicious patterns
		if strings.Contains(link, "http://") || strings.Contains(link, "https://") {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "content",
				Message:  fmt.Sprintf("wikilink contains URL: %s", link),
				Severity: SeverityWarning,
			})
		}
	}

	return errors
}

// PhaseRule validates phase-specific requirements
type PhaseRule struct{}

func (r *PhaseRule) Name() string {
	return "phase"
}

func (r *PhaseRule) Validate(file *obsidian.MarkdownFile) []ValidationError {
	var errors []ValidationError

	if file.FileType != obsidian.FileTypePhase {
		return errors
	}

	if file.Frontmatter == nil {
		return errors
	}

	fm := file.Frontmatter

	// Check phase-specific fields
	if fm.PhaseID == "" {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "phase_id",
			Message:  "phase ID missing",
			Severity: SeverityError,
		})
	}

	if fm.PhaseNum <= 0 || fm.PhaseNum > 6 {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "phase_number",
			Message:  fmt.Sprintf("invalid phase number: %d (must be 1-6)", fm.PhaseNum),
			Severity: SeverityError,
		})
	}

	// Check gates structure
	if fm.Gates == nil {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "gates",
			Message:  "gates definition missing",
			Severity: SeverityWarning,
		})
	} else {
		if len(fm.Gates.Entry) == 0 {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "gates.entry",
				Message:  "entry gates not defined",
				Severity: SeverityWarning,
			})
		}
		if len(fm.Gates.Exit) == 0 {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "gates.exit",
				Message:  "exit gates not defined",
				Severity: SeverityWarning,
			})
		}
	}

	return errors
}

// FeatureRule validates feature specification requirements
type FeatureRule struct{}

func (r *FeatureRule) Name() string {
	return "feature"
}

func (r *FeatureRule) Validate(file *obsidian.MarkdownFile) []ValidationError {
	var errors []ValidationError

	if file.FileType != obsidian.FileTypeFeature {
		return errors
	}

	if file.Frontmatter == nil {
		return errors
	}

	fm := file.Frontmatter

	// Check feature ID format
	if fm.FeatureID == "" {
		errors = append(errors, ValidationError{
			File:     file.Path,
			Field:    "feature_id",
			Message:  "feature ID missing",
			Severity: SeverityError,
		})
	} else {
		if matched, _ := regexp.MatchString(`^FEAT-\d+$`, fm.FeatureID); !matched {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "feature_id",
				Message:  fmt.Sprintf("invalid feature ID format: %s (expected FEAT-XXX)", fm.FeatureID),
				Severity: SeverityError,
			})
		}
	}

	// Check priority format
	if fm.Priority != "" {
		validPriorities := []string{"P0", "P1", "P2", "P3"}
		valid := false
		for _, p := range validPriorities {
			if fm.Priority == p {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "priority",
				Message:  fmt.Sprintf("invalid priority: %s (expected P0-P3)", fm.Priority),
				Severity: SeverityError,
			})
		}
	}

	// Check status values
	if fm.Status != "" {
		validStatuses := []string{"draft", "specified", "approved", "in_progress", "completed", "deprecated"}
		valid := false
		for _, s := range validStatuses {
			if fm.Status == s {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, ValidationError{
				File:     file.Path,
				Field:    "status",
				Message:  fmt.Sprintf("invalid status: %s", fm.Status),
				Severity: SeverityWarning,
			})
		}
	}

	return errors
}

// Helper functions

func initializeRequiredFields() map[obsidian.FileType][]string {
	return map[obsidian.FileType][]string{
		obsidian.FileTypePhase:       {"title", "type", "phase_id", "phase_number", "tags"},
		obsidian.FileTypeEnforcer:    {"title", "type", "phase", "tags"},
		obsidian.FileTypeArtifact:    {"title", "type", "artifact_category", "phase", "tags"},
		obsidian.FileTypeTemplate:    {"title", "type", "tags"},
		obsidian.FileTypePrompt:      {"title", "type", "tags"},
		obsidian.FileTypeExample:     {"title", "type", "tags"},
		obsidian.FileTypeFeature:     {"title", "type", "feature_id", "tags"},
		obsidian.FileTypeCoordinator: {"title", "type", "tags"},
		obsidian.FileTypePrinciple:   {"title", "type", "tags"},
	}
}

func getRequiredFields(fileType obsidian.FileType) []string {
	fields := initializeRequiredFields()
	if required, ok := fields[fileType]; ok {
		return required
	}
	return []string{"title", "type", "tags"}
}

func hasField(fm *obsidian.Frontmatter, field string) bool {
	switch field {
	case "title":
		return fm.Title != ""
	case "type":
		return fm.Type != ""
	case "tags":
		return len(fm.Tags) > 0
	case "phase_id":
		return fm.PhaseID != ""
	case "phase_number":
		return fm.PhaseNum > 0
	case "artifact_category":
		return fm.ArtifactCategory != ""
	case "phase":
		return fm.Phase != ""
	case "feature_id":
		return fm.FeatureID != ""
	default:
		return false
	}
}

func extractWikilinks(content string) []string {
	var links []string

	start := 0
	for {
		idx := strings.Index(content[start:], "[[")
		if idx == -1 {
			break
		}

		startIdx := start + idx + 2
		endIdx := strings.Index(content[startIdx:], "]]")
		if endIdx == -1 {
			break
		}

		link := content[startIdx : startIdx+endIdx]
		// Extract just the target (before | or # or ^)
		if pipeIdx := strings.Index(link, "|"); pipeIdx != -1 {
			link = link[:pipeIdx]
		}
		if hashIdx := strings.Index(link, "#"); hashIdx != -1 {
			link = link[:hashIdx]
		}
		if caretIdx := strings.Index(link, "^"); caretIdx != -1 {
			link = link[:caretIdx]
		}

		links = append(links, strings.TrimSpace(link))
		start = startIdx + endIdx + 2
	}

	return links
}

func initializeValidTags() map[string]bool {
	tags := make(map[string]bool)

	// Core tags
	coreTags := []string{
		"helix",
		"helix/core",
		"helix/workflow",
		"helix/coordinator",
		"helix/principle",
		"helix/enforcer",
		"helix/navigation",
		"dashboard",
	}

	// Phase tags
	phases := []string{"frame", "design", "test", "build", "deploy", "iterate"}
	for _, phase := range phases {
		tags[fmt.Sprintf("helix/phase/%s", phase)] = true
		tags[fmt.Sprintf("helix/phase/%s/enforcer", phase)] = true
	}

	// Artifact tags
	artifactTypes := []string{
		"specification", "design", "test", "implementation",
		"deployment", "monitoring", "template", "prompt", "example",
	}
	for _, aType := range artifactTypes {
		tags[fmt.Sprintf("helix/artifact/%s", aType)] = true
	}

	// Complexity tags
	tags["helix/complexity/simple"] = true
	tags["helix/complexity/moderate"] = true
	tags["helix/complexity/complex"] = true

	// Status tags
	statusTypes := []string{"draft", "review", "approved", "deprecated"}
	for _, status := range statusTypes {
		tags[fmt.Sprintf("helix/status/%s", status)] = true
	}

	// Add core tags
	for _, tag := range coreTags {
		tags[tag] = true
	}

	return tags
}
