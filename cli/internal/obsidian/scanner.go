package obsidian

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileScanner scans directories for markdown files and loads them
type FileScanner struct {
	detector *FileTypeDetector
}

// NewFileScanner creates a new file scanner
func NewFileScanner() *FileScanner {
	return &FileScanner{
		detector: NewFileTypeDetector(),
	}
}

// ScanDirectory scans a directory for markdown files and returns loaded MarkdownFile objects
func (s *FileScanner) ScanDirectory(dir string) ([]*MarkdownFile, error) {
	var files []*MarkdownFile

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Only process markdown files
		if !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		// Load the file
		file, err := s.LoadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load file %s: %w", path, err)
		}

		files = append(files, file)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// LoadFile loads a single markdown file and parses its frontmatter
func (s *FileScanner) LoadFile(path string) (*MarkdownFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file := &MarkdownFile{
		Path:    path,
		Content: string(content),
	}

	// Detect file type
	file.FileType = s.detector.Detect(path)

	// Parse frontmatter if present
	err = s.parseFrontmatter(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter in %s: %w", path, err)
	}

	return file, nil
}

// parseFrontmatter extracts YAML frontmatter from markdown content
func (s *FileScanner) parseFrontmatter(file *MarkdownFile) error {
	content := file.Content

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		// No frontmatter present
		return nil
	}

	// Find the end of frontmatter
	lines := strings.Split(content, "\n")
	var frontmatterLines []string
	var contentLines []string
	inFrontmatter := false
	frontmatterEnd := -1

	for i, line := range lines {
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFrontmatter = true
			continue
		}

		if inFrontmatter && strings.TrimSpace(line) == "---" {
			frontmatterEnd = i
			inFrontmatter = false
			contentLines = lines[i+1:]
			break
		}

		if inFrontmatter {
			frontmatterLines = append(frontmatterLines, line)
		}
	}

	// If we didn't find the closing ---, treat as no frontmatter
	if frontmatterEnd == -1 {
		return nil
	}

	// Parse YAML frontmatter
	if len(frontmatterLines) > 0 {
		yamlContent := strings.Join(frontmatterLines, "\n")

		var fm Frontmatter
		err := yaml.Unmarshal([]byte(yamlContent), &fm)
		if err != nil {
			return fmt.Errorf("invalid YAML frontmatter: %w", err)
		}

		file.Frontmatter = &fm
	}

	// Update content to exclude frontmatter
	file.Content = strings.Join(contentLines, "\n")

	return nil
}

// WriteToFile writes the MarkdownFile back to disk with frontmatter
func (file *MarkdownFile) WriteToFile() error {
	var content strings.Builder

	// Write frontmatter if present
	if file.HasFrontmatter() {
		yamlBytes, err := yaml.Marshal(file.Frontmatter)
		if err != nil {
			return fmt.Errorf("failed to marshal frontmatter: %w", err)
		}

		content.WriteString("---\n")
		content.Write(yamlBytes)
		content.WriteString("---\n")
	}

	// Write main content
	content.WriteString(file.Content)

	// Ensure the directory exists
	dir := filepath.Dir(file.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write to file
	err = os.WriteFile(file.Path, []byte(content.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", file.Path, err)
	}

	return nil
}

// ScanFiles scans multiple files and returns loaded MarkdownFile objects
func (s *FileScanner) ScanFiles(paths []string) ([]*MarkdownFile, error) {
	var files []*MarkdownFile

	for _, path := range paths {
		file, err := s.LoadFile(path)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

// FindFilesByType finds all files of a specific type in a directory
func (s *FileScanner) FindFilesByType(dir string, fileType FileType) ([]*MarkdownFile, error) {
	allFiles, err := s.ScanDirectory(dir)
	if err != nil {
		return nil, err
	}

	var filtered []*MarkdownFile
	for _, file := range allFiles {
		if file.FileType == fileType {
			filtered = append(filtered, file)
		}
	}

	return filtered, nil
}

// FindFilesByTag finds all files with a specific tag
func (s *FileScanner) FindFilesByTag(dir string, tag string) ([]*MarkdownFile, error) {
	allFiles, err := s.ScanDirectory(dir)
	if err != nil {
		return nil, err
	}

	var filtered []*MarkdownFile
	for _, file := range allFiles {
		tags := file.GetTags()
		for _, t := range tags {
			if t == tag {
				filtered = append(filtered, file)
				break
			}
		}
	}

	return filtered, nil
}

// GetStats returns statistics about scanned files
func (s *FileScanner) GetStats(files []*MarkdownFile) *ScanStats {
	stats := &ScanStats{
		Total:              len(files),
		ByType:             make(map[FileType]int),
		WithFrontmatter:    0,
		WithoutFrontmatter: 0,
	}

	for _, file := range files {
		stats.ByType[file.FileType]++

		if file.HasFrontmatter() {
			stats.WithFrontmatter++
		} else {
			stats.WithoutFrontmatter++
		}
	}

	return stats
}

// ScanStats holds statistics about scanned files
type ScanStats struct {
	Total              int
	ByType             map[FileType]int
	WithFrontmatter    int
	WithoutFrontmatter int
}

// Print outputs the stats in a readable format
func (stats *ScanStats) Print() {
	fmt.Printf("Scan Statistics:\n")
	fmt.Printf("  Total files: %d\n", stats.Total)
	fmt.Printf("  With frontmatter: %d\n", stats.WithFrontmatter)
	fmt.Printf("  Without frontmatter: %d\n", stats.WithoutFrontmatter)
	fmt.Printf("  By type:\n")

	for fileType, count := range stats.ByType {
		if count > 0 {
			fmt.Printf("    %s: %d\n", fileType, count)
		}
	}
}

// ValidateFileContent checks if file content is valid markdown
func ValidateFileContent(content string) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for common issues
		if strings.Contains(line, "\x00") {
			return fmt.Errorf("line %d contains null bytes", lineNum)
		}

		// Check for extremely long lines (potential binary content)
		if len(line) > 10000 {
			return fmt.Errorf("line %d is suspiciously long (%d characters)", lineNum, len(line))
		}
	}

	return scanner.Err()
}
