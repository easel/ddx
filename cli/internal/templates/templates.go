package templates

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// Apply applies a template to the target directory
func Apply(templateName, targetDir string, variables map[string]string) error {
	// Comprehensive input validation
	if err := validateTemplateName(templateName); err != nil {
		return fmt.Errorf("invalid template name: %w", err)
	}
	if err := validateTargetDirectory(targetDir); err != nil {
		return fmt.Errorf("invalid target directory: %w", err)
	}
	if err := validateVariables(variables); err != nil {
		return fmt.Errorf("invalid variables: %w", err)
	}
	if variables == nil {
		variables = make(map[string]string)
	}

	// Get and validate DDx home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	ddxHome := filepath.Join(home, ".ddx")

	// Sanitize template name to prevent path traversal
	sanitizedTemplateName := sanitizePathComponent(templateName)
	templatePath := filepath.Join(ddxHome, "templates", sanitizedTemplateName)

	// Verify the template path is within expected directory
	if err := validateTemplatePath(templatePath, ddxHome); err != nil {
		return fmt.Errorf("invalid template path: %w", err)
	}

	// Check if template exists and is a directory
	info, err := os.Stat(templatePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("template '%s' not found", sanitizedTemplateName)
	}
	if err != nil {
		return fmt.Errorf("failed to check template: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("template '%s' is not a directory", sanitizedTemplateName)
	}

	// Validate template size and complexity before processing
	if err := validateTemplateStructure(templatePath); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	// Apply the template with security context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	return applyTemplateWithContext(ctx, templatePath, targetDir, variables)
}

// applyTemplateWithContext recursively applies template files with context and security checks
func applyTemplateWithContext(ctx context.Context, templateDir, targetDir string, variables map[string]string) error {
	// Clean and validate target directory
	cleanTargetDir := filepath.Clean(targetDir)
	if err := validateTargetDirectory(cleanTargetDir); err != nil {
		return fmt.Errorf("invalid target directory: %w", err)
	}

	// Ensure target directory exists with secure permissions
	if err := os.MkdirAll(cleanTargetDir, 0750); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return fmt.Errorf("error walking template directory: %w", err)
		}

		// Security check: validate file info
		if err := validateFileInfo(info); err != nil {
			return fmt.Errorf("invalid file in template: %w", err)
		}

		// Get relative path from template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Enhanced security: validate relative path
		if err := validateRelativePath(relPath); err != nil {
			return fmt.Errorf("invalid path in template: %w", err)
		}

		// Skip hidden files and directories (any component starting with .)
		pathParts := strings.Split(relPath, string(filepath.Separator))
		for _, part := range pathParts {
			if strings.HasPrefix(part, ".") {
				return nil
			}
		}

		// Process variables in the relative path for security
		processedRelPath := secureReplaceVariables(relPath, variables)
		targetPath := filepath.Join(cleanTargetDir, processedRelPath)

		// Critical security check: ensure target path is within target directory
		if err := validateTargetPath(targetPath, cleanTargetDir); err != nil {
			return fmt.Errorf("path traversal detected: %w", err)
		}

		if info.IsDir() {
			// Create directory with secure permissions (never more permissive than 0755)
			mode := info.Mode() & 0755
			return os.MkdirAll(targetPath, mode)
		}

		// Process file with context and security checks
		return processTemplateFileSecure(ctx, path, targetPath, variables, info)
	})
}

// processTemplateFileSecure processes a single template file with security checks
func processTemplateFileSecure(ctx context.Context, sourcePath, targetPath string, variables map[string]string, info os.FileInfo) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Security validation: check file size limits
	if info.Size() > maxTemplateFileSize {
		return fmt.Errorf("template file too large: %d bytes (max %d)", info.Size(), maxTemplateFileSize)
	}

	// Read source file with size limit
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open template file: %w", err)
	}
	defer file.Close()

	// Use limited reader to prevent memory exhaustion
	limitedReader := io.LimitReader(file, maxTemplateFileSize)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Validate content is valid UTF-8
	if !utf8.Valid(content) {
		return fmt.Errorf("template file contains invalid UTF-8")
	}

	// Security scan: check for potentially dangerous content
	if err := scanForDangerousContent(string(content)); err != nil {
		return fmt.Errorf("dangerous content detected in template: %w", err)
	}

	// Replace variables in content with security checks
	processedContent, err := secureReplaceVariablesInContent(string(content), variables)
	if err != nil {
		return fmt.Errorf("variable replacement failed: %w", err)
	}

	// Create target directory if needed with secure permissions
	if err := os.MkdirAll(filepath.Dir(targetPath), 0750); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write processed content with secure permissions
	// Preserve original permissions but cap at 0644 for security
	mode := info.Mode()
	if mode > 0644 {
		// If original is more permissive than 0644, cap it
		mode = 0644
	}
	// Ensure at least 0644 permissions
	if mode&0600 == 0 {
		mode = 0644
	}

	return os.WriteFile(targetPath, []byte(processedContent), mode)
}

// secureReplaceVariablesInContent replaces template variables in content with security checks
func secureReplaceVariablesInContent(content string, variables map[string]string) (string, error) {
	result := content
	maxReplacements := 1000 // Prevent infinite loops or excessive processing
	replacementCount := 0

	for key, value := range variables {
		// Validate key and value
		if err := validateVariableName(key); err != nil {
			return "", fmt.Errorf("invalid variable name '%s': %w", key, err)
		}
		if err := validateVariableValue(value); err != nil {
			return "", fmt.Errorf("invalid variable value for '%s': %w", key, err)
		}

		// Sanitize value to prevent injection
		sanitizedValue := sanitizeVariableValue(value)

		// Replace {{key}} patterns (exact match)
		oldResult := result
		result = strings.ReplaceAll(result, "{{"+key+"}}", sanitizedValue)
		if result != oldResult {
			replacementCount++
		}

		// Replace {{ key }} patterns (with spaces)
		oldResult = result
		result = strings.ReplaceAll(result, "{{ "+key+" }}", sanitizedValue)
		if result != oldResult {
			replacementCount++
		}

		// Replace {{key }} patterns (space after)
		oldResult = result
		result = strings.ReplaceAll(result, "{{"+key+" }}", sanitizedValue)
		if result != oldResult {
			replacementCount++
		}

		// Replace {{ key}} patterns (space before)
		oldResult = result
		result = strings.ReplaceAll(result, "{{ "+key+"}}", sanitizedValue)
		if result != oldResult {
			replacementCount++
		}

		// Replace ${KEY} patterns (uppercase)
		upperKey := strings.ToUpper(key)
		oldResult = result
		result = strings.ReplaceAll(result, "${"+upperKey+"}", sanitizedValue)
		if result != oldResult {
			replacementCount++
		}

		// Replace ${key} patterns (original case)
		oldResult = result
		result = strings.ReplaceAll(result, "${"+key+"}", sanitizedValue)
		if result != oldResult {
			replacementCount++
		}

		// Check replacement limit
		if replacementCount > maxReplacements {
			return "", fmt.Errorf("too many variable replacements (max %d)", maxReplacements)
		}
	}

	return result, nil
}

// List returns available templates with caching for performance
func List() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	templatesPath := filepath.Join(home, ".ddx", "templates")

	// Security validation: ensure templates path is safe
	if err := validateTemplatesDirectory(templatesPath); err != nil {
		return nil, fmt.Errorf("invalid templates directory: %w", err)
	}

	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to check templates directory: %w", err)
	}

	// Check cache first
	if cached := getCachedTemplateList(templatesPath); cached != nil {
		return cached, nil
	}

	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []string
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden directories and validate name
		if entry.IsDir() && !strings.HasPrefix(name, ".") {
			// Validate template name for security
			if err := validateTemplateName(name); err == nil {
				templates = append(templates, name)
			}
		}
	}

	// Cache the result
	cacheTemplateList(templatesPath, templates)

	return templates, nil
}

// Security constants and caching
const (
	maxTemplateFileSize   = 10 * 1024 * 1024 // 10MB per file
	maxTemplateFiles      = 1000             // Maximum number of files in a template
	maxTemplateDepth      = 20               // Maximum directory nesting depth
	maxVariableLength     = 1024             // Maximum variable value length
	maxVariableCount      = 100              // Maximum number of variables
	cacheExpiration       = 5 * time.Minute  // Cache expiration time
)

var (
	// Cache for template listings with timestamps
	templateListCache = sync.Map{}

	// Validation patterns
	validTemplateNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*[a-zA-Z0-9]$`)
	validVariableNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	dangerousPatterns        = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(<script|javascript:|vbscript:|onload=|onerror=)`),
		regexp.MustCompile(`(?i)(eval\s*\(|exec\s*\(|system\s*\(|shell_exec)`),
		regexp.MustCompile(`(?i)(rm\s+-rf|del\s+/|format\s+c:)`),
		regexp.MustCompile(`\$\{[^}]*\$\{`), // Nested variable substitution
	}
)

type templateListCacheEntry struct {
	templates []string
	timestamp time.Time
	pathHash  string
}

// Security validation functions

// validateTemplateName validates template names
func validateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if len(name) > 100 {
		return fmt.Errorf("template name too long (max 100 characters)")
	}

	// Allow single character names for root directories
	if len(name) == 1 {
		if !regexp.MustCompile(`^[a-zA-Z0-9]$`).MatchString(name) {
			return fmt.Errorf("single character template name must be alphanumeric")
		}
		return nil
	}

	if !validTemplateNamePattern.MatchString(name) {
		return fmt.Errorf("template name contains invalid characters (only alphanumeric, dots, underscores, and hyphens allowed)")
	}

	// Prevent certain dangerous names
	dangerousNames := []string{"con", "prn", "aux", "nul", "com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9"}
	lowerName := strings.ToLower(name)
	for _, dangerous := range dangerousNames {
		if lowerName == dangerous {
			return fmt.Errorf("template name '%s' is reserved", name)
		}
	}

	return nil
}

// validateTargetDirectory validates target directories
func validateTargetDirectory(dir string) error {
	if dir == "" {
		return fmt.Errorf("target directory cannot be empty")
	}

	if len(dir) > 1024 {
		return fmt.Errorf("target directory path too long (max 1024 characters)")
	}

	// Clean the path and check for traversal attempts
	cleanDir := filepath.Clean(dir)
	if strings.Contains(cleanDir, "..") {
		return fmt.Errorf("target directory cannot contain path traversal sequences")
	}

	// Prevent writing to system directories
	systemDirs := []string{"/etc", "/bin", "/sbin", "/usr/bin", "/usr/sbin", "/root", "/boot", "/sys", "/proc", "/dev"}
	for _, sysDir := range systemDirs {
		if strings.HasPrefix(cleanDir, sysDir) {
			return fmt.Errorf("cannot write to system directory: %s", sysDir)
		}
	}

	return nil
}

// validateVariables validates the variables map
func validateVariables(variables map[string]string) error {
	if variables == nil {
		return nil
	}

	if len(variables) > maxVariableCount {
		return fmt.Errorf("too many variables (max %d)", maxVariableCount)
	}

	for key, value := range variables {
		if err := validateVariableName(key); err != nil {
			return fmt.Errorf("invalid variable name '%s': %w", key, err)
		}
		if err := validateVariableValue(value); err != nil {
			return fmt.Errorf("invalid variable value for '%s': %w", key, err)
		}
	}

	return nil
}

// validateVariableName validates variable names
func validateVariableName(name string) error {
	if name == "" {
		return fmt.Errorf("variable name cannot be empty")
	}

	if len(name) > 64 {
		return fmt.Errorf("variable name too long (max 64 characters)")
	}

	if !validVariableNamePattern.MatchString(name) {
		return fmt.Errorf("variable name must start with letter or underscore and contain only letters, numbers, and underscores")
	}

	return nil
}

// validateVariableValue validates variable values
func validateVariableValue(value string) error {
	if len(value) > maxVariableLength {
		return fmt.Errorf("variable value too long (max %d characters)", maxVariableLength)
	}

	// Check for null bytes and other dangerous characters
	if strings.ContainsAny(value, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x0e\x0f") {
		return fmt.Errorf("variable value contains invalid control characters")
	}

	return nil
}

// validateTemplatePath ensures template path is within expected directory
func validateTemplatePath(templatePath, ddxHome string) error {
	absTemplatePath, err := filepath.Abs(templatePath)
	if err != nil {
		return fmt.Errorf("failed to resolve template path: %w", err)
	}

	absDdxHome, err := filepath.Abs(ddxHome)
	if err != nil {
		return fmt.Errorf("failed to resolve DDx home: %w", err)
	}

	rel, err := filepath.Rel(absDdxHome, absTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}

	if strings.HasPrefix(rel, "..") {
		return fmt.Errorf("template path is outside DDx directory")
	}

	return nil
}

// validateTemplateStructure validates template directory structure
func validateTemplateStructure(templatePath string) error {
	fileCount := 0
	maxDepth := 0

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Count files
		if !info.IsDir() {
			fileCount++
			if fileCount > maxTemplateFiles {
				return fmt.Errorf("template has too many files (max %d)", maxTemplateFiles)
			}
		}

		// Check depth
		rel, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}
		depth := strings.Count(rel, string(filepath.Separator))
		if depth > maxDepth {
			maxDepth = depth
		}
		if maxDepth > maxTemplateDepth {
			return fmt.Errorf("template directory too deep (max depth %d)", maxTemplateDepth)
		}

		return validateFileInfo(info)
	})

	return err
}

// validateFileInfo validates file information
func validateFileInfo(info os.FileInfo) error {
	// Check for suspicious file modes
	mode := info.Mode()

	// Reject setuid/setgid files for security
	if mode&os.ModeSetuid != 0 || mode&os.ModeSetgid != 0 {
		return fmt.Errorf("setuid/setgid files not allowed in templates")
	}

	// Reject device files, sockets, etc.
	if mode&os.ModeDevice != 0 || mode&os.ModeSocket != 0 || mode&os.ModeNamedPipe != 0 {
		return fmt.Errorf("special files not allowed in templates")
	}

	return nil
}

// validateRelativePath validates relative paths within templates
func validateRelativePath(relPath string) error {
	if strings.Contains(relPath, "..") {
		return fmt.Errorf("path traversal not allowed: %s", relPath)
	}

	// Check each component
	parts := strings.Split(relPath, string(filepath.Separator))
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}

		// Validate component length
		if len(part) > 255 {
			return fmt.Errorf("path component too long: %s", part)
		}

		// Check for invalid characters in filenames
		if strings.ContainsAny(part, "\x00<>:\"|?*") {
			return fmt.Errorf("invalid characters in path component: %s", part)
		}
	}

	return nil
}

// validateTargetPath ensures target path is within target directory
func validateTargetPath(targetPath, targetDir string) error {
	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve target directory: %w", err)
	}

	rel, err := filepath.Rel(absTargetDir, absTargetPath)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}

	if strings.HasPrefix(rel, "..") {
		return fmt.Errorf("target path is outside target directory: %s", rel)
	}

	return nil
}

// validateTemplatesDirectory validates the templates directory
func validateTemplatesDirectory(templatesPath string) error {
	absPath, err := filepath.Abs(templatesPath)
	if err != nil {
		return fmt.Errorf("failed to resolve templates path: %w", err)
	}

	// Ensure it's a reasonable path
	if strings.Contains(absPath, "..") {
		return fmt.Errorf("templates path contains traversal sequences")
	}

	return nil
}

// scanForDangerousContent scans template content for dangerous patterns
func scanForDangerousContent(content string) error {
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(content) {
			return fmt.Errorf("potentially dangerous pattern detected")
		}
	}

	// Check for excessive variable nesting (potential ReDoS)
	varCount := strings.Count(content, "{{") + strings.Count(content, "${")
	if varCount > 500 {
		return fmt.Errorf("too many variable placeholders (max 500)")
	}

	return nil
}

// Helper functions for variable replacement and caching

// secureReplaceVariables replaces variables in paths with security checks
func secureReplaceVariables(input string, variables map[string]string) string {
	result := input
	for key, value := range variables {
		// Only replace in paths if the value is safe for filenames
		if isValidFilenameComponent(value) {
			result = strings.ReplaceAll(result, "{{"+key+"}}", value)
			result = strings.ReplaceAll(result, "${"+key+"}", value)
		}
	}
	return result
}

// isValidFilenameComponent checks if a string is safe for use in filenames
func isValidFilenameComponent(s string) bool {
	if s == "" || s == "." || s == ".." {
		return false
	}

	// Check for invalid filename characters
	if strings.ContainsAny(s, "/\\<>:\"|?*\x00") {
		return false
	}

	// Check length
	if len(s) > 255 {
		return false
	}

	return true
}

// sanitizeVariableValue sanitizes variable values
func sanitizeVariableValue(value string) string {
	// Remove null bytes and control characters except newlines and tabs
	var result strings.Builder
	for _, r := range value {
		if r >= 32 || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// sanitizePathComponent sanitizes path components
func sanitizePathComponent(component string) string {
	// Remove dangerous characters
	result := strings.ReplaceAll(component, "..", "")
	result = strings.ReplaceAll(result, "/", "_")
	result = strings.ReplaceAll(result, "\\", "_")
	result = strings.ReplaceAll(result, "\x00", "")

	// Limit length
	if len(result) > 100 {
		result = result[:100]
	}

	return result
}

// Caching functions

// getCachedTemplateList retrieves cached template list
func getCachedTemplateList(templatesPath string) []string {
	pathHash := hashPath(templatesPath)
	if cached, exists := templateListCache.Load(pathHash); exists {
		entry := cached.(templateListCacheEntry)
		if time.Since(entry.timestamp) < cacheExpiration {
			return entry.templates
		}
		// Cache expired, remove it
		templateListCache.Delete(pathHash)
	}
	return nil
}

// cacheTemplateList caches template list
func cacheTemplateList(templatesPath string, templates []string) {
	pathHash := hashPath(templatesPath)
	entry := templateListCacheEntry{
		templates: templates,
		timestamp: time.Now(),
		pathHash:  pathHash,
	}
	templateListCache.Store(pathHash, entry)
}

// hashPath creates a hash of a path for caching
func hashPath(path string) string {
	hash := sha256.Sum256([]byte(path))
	return hex.EncodeToString(hash[:])
}

// Backwards compatibility functions for tests

// replaceVariables provides backwards compatibility for tests
func replaceVariables(content string, variables map[string]string) string {
	result, _ := secureReplaceVariablesInContent(content, variables)
	return result
}

// processTemplateFile provides backwards compatibility for tests
func processTemplateFile(sourcePath, targetPath string, variables map[string]string) error {
	// For backward compatibility with tests, replace variables in target path
	processedTargetPath := replaceVariables(targetPath, variables)

	ctx := context.Background()
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}
	return processTemplateFileSecure(ctx, sourcePath, processedTargetPath, variables, info)
}

// applyTemplate provides backwards compatibility for tests
func applyTemplate(templateDir, targetDir string, variables map[string]string) error {
	ctx := context.Background()
	return applyTemplateWithContext(ctx, templateDir, targetDir, variables)
}
