package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Apply applies a template to the target directory
func Apply(templateName, targetDir string, variables map[string]string) error {
	// Get DDx home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	ddxHome := filepath.Join(home, ".ddx")

	templatePath := filepath.Join(ddxHome, "templates", templateName)
	
	// Check if template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Apply the template
	return applyTemplate(templatePath, targetDir, variables)
}

// applyTemplate recursively applies template files
func applyTemplate(templateDir, targetDir string, variables map[string]string) error {
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		targetPath := filepath.Join(targetDir, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Process file
		return processTemplateFile(path, targetPath, variables)
	})
}

// processTemplateFile processes a single template file
func processTemplateFile(sourcePath, targetPath string, variables map[string]string) error {
	// Read source file
	content, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Replace variables in content
	processedContent := replaceVariables(string(content), variables)

	// Create target directory if needed
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Write processed content
	return os.WriteFile(targetPath, []byte(processedContent), 0644)
}

// replaceVariables replaces template variables in content
func replaceVariables(content string, variables map[string]string) string {
	result := content
	
	for key, value := range variables {
		// Replace {{key}} patterns
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
		result = strings.ReplaceAll(result, "{{ "+key+" }}", value)
		
		// Replace ${KEY} patterns (uppercase)
		upperKey := strings.ToUpper(key)
		result = strings.ReplaceAll(result, "${"+upperKey+"}", value)
	}
	
	return result
}

// List returns available templates
func List() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	templatesPath := filepath.Join(home, ".ddx", "templates")
	
	if _, err := os.Stat(templatesPath); os.IsNotExist(err) {
		return []string{}, nil
	}
	
	entries, err := os.ReadDir(templatesPath)
	if err != nil {
		return nil, err
	}
	
	var templates []string
	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}
	
	return templates, nil
}