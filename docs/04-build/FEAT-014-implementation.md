---
title: "Build Implementation - Obsidian Integration for HELIX"
type: implementation-guide
feature_id: FEAT-014
phase: build
tags:
  - helix/build
  - helix/artifact/implementation
  - obsidian
  - golang
related:
  - "[[FEAT-014-obsidian-integration]]"
  - "[[FEAT-014-technical-design]]"
created: 2025-01-18
updated: 2025-01-18
status: ready
---

# Build Implementation: [[FEAT-014]] Obsidian Integration

## Implementation Checklist

- [ ] Create Go module structure for Obsidian integration
- [ ] Implement file type detection logic
- [ ] Build frontmatter generation system
- [ ] Create wikilink conversion engine
- [ ] Develop navigation hub generator
- [ ] Add CLI commands for migration
- [ ] Create validation framework
- [ ] Write comprehensive tests
- [ ] Document usage and examples

## Module Implementation

### Step 1: Create Obsidian Package Structure

```bash
# Create package structure
mkdir -p cli/internal/obsidian/{schemas,converter,validator}
mkdir -p cli/cmd/obsidian
mkdir -p workflows/helix/obsidian/{schemas,templates}
```

### Step 2: Core Types and Interfaces

```go
// cli/internal/obsidian/types.go
package obsidian

import (
    "time"
    "gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML frontmatter for any markdown file
type Frontmatter struct {
    Title     string    `yaml:"title"`
    Type      string    `yaml:"type"`
    Tags      []string  `yaml:"tags"`
    Created   time.Time `yaml:"created"`
    Updated   time.Time `yaml:"updated"`
    Status    string    `yaml:"status,omitempty"`
    Version   string    `yaml:"version,omitempty"`
    Aliases   []string  `yaml:"aliases,omitempty"`
    Related   []string  `yaml:"related,omitempty"`

    // Phase-specific fields
    PhaseID   string    `yaml:"phase_id,omitempty"`
    PhaseNum  int       `yaml:"phase_number,omitempty"`
    NextPhase string    `yaml:"next_phase,omitempty"`
    PrevPhase string    `yaml:"previous_phase,omitempty"`
    Gates     *Gates    `yaml:"gates,omitempty"`
    Artifacts *Artifacts `yaml:"artifacts,omitempty"`

    // Artifact-specific fields
    ArtifactCategory string   `yaml:"artifact_category,omitempty"`
    Phase           string   `yaml:"phase,omitempty"`
    Complexity      string   `yaml:"complexity,omitempty"`
    Prerequisites   []string `yaml:"prerequisites,omitempty"`
    Outputs        []string `yaml:"outputs,omitempty"`
    TimeEstimate   string   `yaml:"time_estimate,omitempty"`
    SkillsRequired []string `yaml:"skills_required,omitempty"`
}

type Gates struct {
    Entry []string `yaml:"entry,omitempty"`
    Exit  []string `yaml:"exit,omitempty"`
}

type Artifacts struct {
    Required []string `yaml:"required,omitempty"`
    Optional []string `yaml:"optional,omitempty"`
}

// MarkdownFile represents a markdown file with optional frontmatter
type MarkdownFile struct {
    Path        string
    Content     string
    Frontmatter *Frontmatter
    FileType    FileType
}

// FileType represents the type of HELIX file
type FileType string

const (
    FileTypeUnknown     FileType = "unknown"
    FileTypePhase       FileType = "phase"
    FileTypeEnforcer    FileType = "enforcer"
    FileTypeArtifact    FileType = "artifact"
    FileTypeTemplate    FileType = "template"
    FileTypePrompt      FileType = "prompt"
    FileTypeExample     FileType = "example"
    FileTypeCoordinator FileType = "coordinator"
    FileTypePrinciple   FileType = "principle"
)
```

### Step 3: File Type Detection

```go
// cli/internal/obsidian/detector.go
package obsidian

import (
    "path/filepath"
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
            "phases/*/README.md":       FileTypePhase,
            "phases/*/enforcer.md":     FileTypeEnforcer,
            "*/template.md":            FileTypeTemplate,
            "*/prompt.md":              FileTypePrompt,
            "*/example.md":             FileTypeExample,
            "coordinator.md":           FileTypeCoordinator,
            "principles.md":            FileTypePrinciple,
            "artifacts/*/README.md":    FileTypeArtifact,
        },
    }
}

// Detect determines the file type based on path patterns
func (d *FileTypeDetector) Detect(path string) FileType {
    // Normalize path
    path = filepath.ToSlash(path)

    // Check exact matches first
    filename := filepath.Base(path)
    if fileType, ok := d.patterns[filename]; ok {
        return fileType
    }

    // Check path patterns
    for pattern, fileType := range d.patterns {
        if matched, _ := filepath.Match(pattern, path); matched {
            return fileType
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
            if strings.HasSuffix(path, "/template.md") {
                return FileTypeTemplate
            }
            if strings.HasSuffix(path, "/prompt.md") {
                return FileTypePrompt
            }
            if strings.HasSuffix(path, "/example.md") {
                return FileTypeExample
            }
        }
    }

    return FileTypeUnknown
}

// GetPhaseFromPath extracts the phase name from a file path
func GetPhaseFromPath(path string) string {
    parts := strings.Split(filepath.ToSlash(path), "/")
    for i, part := range parts {
        if part == "phases" && i+1 < len(parts) {
            phaseName := parts[i+1]
            // Remove number prefix if present (e.g., "01-frame" -> "frame")
            if idx := strings.Index(phaseName, "-"); idx > 0 {
                return phaseName[idx+1:]
            }
            return phaseName
        }
    }
    return ""
}

// GetArtifactCategory extracts the artifact category from a file path
func GetArtifactCategory(path string) string {
    parts := strings.Split(filepath.ToSlash(path), "/")
    for i, part := range parts {
        if part == "artifacts" && i+1 < len(parts) {
            return parts[i+1]
        }
    }
    return ""
}
```

### Step 4: Frontmatter Generator

```go
// cli/internal/obsidian/generator.go
package obsidian

import (
    "fmt"
    "path/filepath"
    "regexp"
    "strings"
    "time"
)

// FrontmatterGenerator generates appropriate frontmatter for files
type FrontmatterGenerator struct {
    detector *FileTypeDetector
}

// NewFrontmatterGenerator creates a new frontmatter generator
func NewFrontmatterGenerator() *FrontmatterGenerator {
    return &FrontmatterGenerator{
        detector: NewFileTypeDetector(),
    }
}

// Generate creates frontmatter for a given file
func (g *FrontmatterGenerator) Generate(file *MarkdownFile) (*Frontmatter, error) {
    fm := &Frontmatter{
        Created: time.Now(),
        Updated: time.Now(),
        Tags:    []string{},
    }

    // Extract title from content
    fm.Title = g.extractTitle(file.Content)

    // Set type based on file type
    fm.Type = string(file.FileType)

    // Generate tags based on file type and location
    fm.Tags = g.generateTags(file)

    // Add type-specific metadata
    switch file.FileType {
    case FileTypePhase:
        g.addPhaseMetadata(fm, file)
    case FileTypeArtifact, FileTypeTemplate, FileTypePrompt, FileTypeExample:
        g.addArtifactMetadata(fm, file)
    case FileTypeEnforcer:
        g.addEnforcerMetadata(fm, file)
    case FileTypeCoordinator:
        fm.Tags = append(fm.Tags, "helix/core", "helix/coordinator")
    case FileTypePrinciple:
        fm.Tags = append(fm.Tags, "helix/core", "helix/principle")
    }

    return fm, nil
}

// extractTitle extracts the title from markdown content
func (g *FrontmatterGenerator) extractTitle(content string) string {
    // Look for first H1 heading
    re := regexp.MustCompile(`(?m)^#\s+(.+)$`)
    matches := re.FindStringSubmatch(content)
    if len(matches) > 1 {
        return strings.TrimSpace(matches[1])
    }
    return "Untitled"
}

// generateTags creates tags based on file type and path
func (g *FrontmatterGenerator) generateTags(file *MarkdownFile) []string {
    tags := []string{"helix"}

    // Add phase tag if applicable
    if phase := GetPhaseFromPath(file.Path); phase != "" {
        tags = append(tags, fmt.Sprintf("helix/phase/%s", phase))
    }

    // Add type-specific tags
    switch file.FileType {
    case FileTypePhase:
        tags = append(tags, "helix/phase")
    case FileTypeArtifact, FileTypeTemplate:
        tags = append(tags, "helix/artifact", "helix/artifact/template")
    case FileTypePrompt:
        tags = append(tags, "helix/artifact", "helix/artifact/prompt")
    case FileTypeExample:
        tags = append(tags, "helix/artifact", "helix/artifact/example")
    case FileTypeEnforcer:
        tags = append(tags, "helix/enforcer")
    }

    return tags
}

// addPhaseMetadata adds phase-specific metadata
func (g *FrontmatterGenerator) addPhaseMetadata(fm *Frontmatter, file *MarkdownFile) {
    phase := GetPhaseFromPath(file.Path)
    fm.PhaseID = phase

    // Map phase names to numbers
    phaseNumbers := map[string]int{
        "frame":   1,
        "design":  2,
        "test":    3,
        "build":   4,
        "deploy":  5,
        "iterate": 6,
    }

    if num, ok := phaseNumbers[phase]; ok {
        fm.PhaseNum = num

        // Set next/previous phases
        phaseNames := []string{"", "frame", "design", "test", "build", "deploy", "iterate"}
        if num > 1 {
            fm.PrevPhase = fmt.Sprintf("[[%s Phase]]", strings.Title(phaseNames[num-1]))
        }
        if num < 6 {
            fm.NextPhase = fmt.Sprintf("[[%s Phase]]", strings.Title(phaseNames[num+1]))
        }
    }

    // Add gates and artifacts placeholders
    fm.Gates = &Gates{
        Entry: []string{"[[TODO: Add entry gates]]"},
        Exit:  []string{"[[TODO: Add exit gates]]"},
    }
    fm.Artifacts = &Artifacts{
        Required: []string{"[[TODO: Add required artifacts]]"},
        Optional: []string{"[[TODO: Add optional artifacts]]"},
    }
}

// addArtifactMetadata adds artifact-specific metadata
func (g *FrontmatterGenerator) addArtifactMetadata(fm *Frontmatter, file *MarkdownFile) {
    fm.Phase = GetPhaseFromPath(file.Path)
    fm.ArtifactCategory = GetArtifactCategory(file.Path)

    // Set default complexity based on content length
    contentLength := len(file.Content)
    switch {
    case contentLength < 1000:
        fm.Complexity = "simple"
    case contentLength < 5000:
        fm.Complexity = "moderate"
    default:
        fm.Complexity = "complex"
    }

    // Add common artifact fields
    fm.Prerequisites = []string{}
    fm.Outputs = []string{}
}

// addEnforcerMetadata adds enforcer-specific metadata
func (g *FrontmatterGenerator) addEnforcerMetadata(fm *Frontmatter, file *MarkdownFile) {
    phase := GetPhaseFromPath(file.Path)
    fm.Phase = phase
    fm.Tags = append(fm.Tags, fmt.Sprintf("helix/phase/%s/enforcer", phase))
    fm.Aliases = []string{
        fmt.Sprintf("%s Phase Enforcer", strings.Title(phase)),
        fmt.Sprintf("%s Guardian", strings.Title(phase)),
    }
}
```

### Step 5: Wikilink Converter

```go
// cli/internal/obsidian/converter/links.go
package converter

import (
    "fmt"
    "path/filepath"
    "regexp"
    "strings"
)

// LinkConverter converts markdown links to Obsidian wikilinks
type LinkConverter struct {
    fileIndex map[string]string // path -> title mapping
    aliases   map[string]string // alias -> canonical name
}

// NewLinkConverter creates a new link converter
func NewLinkConverter() *LinkConverter {
    return &LinkConverter{
        fileIndex: make(map[string]string),
        aliases:   make(map[string]string),
    }
}

// BuildIndex builds the file index for link resolution
func (c *LinkConverter) BuildIndex(files []*MarkdownFile) {
    for _, file := range files {
        // Map file path to title
        if file.Frontmatter != nil {
            c.fileIndex[file.Path] = file.Frontmatter.Title

            // Register aliases
            for _, alias := range file.Frontmatter.Aliases {
                c.aliases[alias] = file.Frontmatter.Title
            }
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

        // Convert relative path to wikilink
        return c.convertToWikilink(linkText, linkPath)
    })

    // Also convert common phase references
    content = c.convertPhaseReferences(content)

    return content
}

// convertToWikilink converts a single link to wikilink format
func (c *LinkConverter) convertToWikilink(text, path string) string {
    // Remove .md extension
    path = strings.TrimSuffix(path, ".md")

    // Try to resolve to a known file title
    if title, ok := c.fileIndex[path]; ok {
        if text != title {
            // Use alias syntax if link text differs from title
            return fmt.Sprintf("[[%s|%s]]", title, text)
        }
        return fmt.Sprintf("[[%s]]", title)
    }

    // Extract just the filename for common patterns
    filename := filepath.Base(path)

    // Handle common patterns
    switch filename {
    case "README":
        // Try to determine phase from path
        if strings.Contains(path, "/phases/") {
            phase := extractPhaseFromPath(path)
            if phase != "" {
                return fmt.Sprintf("[[%s Phase]]", strings.Title(phase))
            }
        }
    case "template":
        artifact := extractArtifactFromPath(path)
        if artifact != "" {
            return fmt.Sprintf("[[%s Template]]", strings.Title(strings.ReplaceAll(artifact, "-", " ")))
        }
    case "prompt":
        artifact := extractArtifactFromPath(path)
        if artifact != "" {
            return fmt.Sprintf("[[%s Prompt]]", strings.Title(strings.ReplaceAll(artifact, "-", " ")))
        }
    case "example":
        artifact := extractArtifactFromPath(path)
        if artifact != "" {
            return fmt.Sprintf("[[%s Example]]", strings.Title(strings.ReplaceAll(artifact, "-", " ")))
        }
    case "enforcer":
        phase := extractPhaseFromPath(path)
        if phase != "" {
            return fmt.Sprintf("[[%s Phase Enforcer]]", strings.Title(phase))
        }
    }

    // Default: use the link text as wikilink
    return fmt.Sprintf("[[%s]]", text)
}

// convertPhaseReferences converts common phase references to wikilinks
func (c *LinkConverter) convertPhaseReferences(content string) string {
    phases := []string{"Frame", "Design", "Test", "Build", "Deploy", "Iterate"}

    for _, phase := range phases {
        // Convert "Frame phase" -> "[[Frame Phase]]"
        pattern := regexp.MustCompile(fmt.Sprintf(`\b%s phase\b`, phase))
        content = pattern.ReplaceAllString(content, fmt.Sprintf("[[%s Phase]]", phase))

        // Convert "Frame Phase" -> "[[Frame Phase]]" (if not already in wikilink)
        pattern = regexp.MustCompile(fmt.Sprintf(`(?<!\[\[)%s Phase(?!\]\])`, phase))
        content = pattern.ReplaceAllString(content, fmt.Sprintf("[[%s Phase]]", phase))
    }

    return content
}

// Helper functions
func extractPhaseFromPath(path string) string {
    parts := strings.Split(path, "/")
    for i, part := range parts {
        if part == "phases" && i+1 < len(parts) {
            phase := parts[i+1]
            if idx := strings.Index(phase, "-"); idx > 0 {
                return phase[idx+1:]
            }
            return phase
        }
    }
    return ""
}

func extractArtifactFromPath(path string) string {
    parts := strings.Split(path, "/")
    for i, part := range parts {
        if part == "artifacts" && i+1 < len(parts) {
            return parts[i+1]
        }
    }
    return ""
}
```

### Step 6: CLI Commands

```go
// cli/cmd/obsidian.go
package cmd

import (
    "fmt"
    "io/ioutil"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"

    "ddx/cli/internal/obsidian"
    "ddx/cli/internal/obsidian/converter"
    "ddx/cli/internal/obsidian/validator"
)

var obsidianCmd = &cobra.Command{
    Use:   "obsidian",
    Short: "Manage Obsidian integration for HELIX workflow",
    Long:  `Convert HELIX workflow documentation to Obsidian-compatible format with frontmatter and wikilinks.`,
}

var migrateCmd = &cobra.Command{
    Use:   "migrate [path]",
    Short: "Migrate HELIX files to Obsidian format",
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        path := "workflows/helix"
        if len(args) > 0 {
            path = args[0]
        }

        dryRun, _ := cmd.Flags().GetBool("dry-run")

        fmt.Printf("🔄 Starting Obsidian migration for %s...\n", path)

        // Step 1: Scan files
        files, err := scanMarkdownFiles(path)
        if err != nil {
            return fmt.Errorf("failed to scan files: %w", err)
        }
        fmt.Printf("📁 Found %d markdown files\n", len(files))

        // Step 2: Detect file types
        detector := obsidian.NewFileTypeDetector()
        for _, file := range files {
            file.FileType = detector.Detect(file.Path)
        }

        // Step 3: Generate frontmatter
        generator := obsidian.NewFrontmatterGenerator()
        for _, file := range files {
            fm, err := generator.Generate(file)
            if err != nil {
                fmt.Printf("⚠️  Failed to generate frontmatter for %s: %v\n", file.Path, err)
                continue
            }
            file.Frontmatter = fm
        }

        // Step 4: Convert links
        linkConverter := converter.NewLinkConverter()
        linkConverter.BuildIndex(files)

        for _, file := range files {
            file.Content = linkConverter.ConvertContent(file.Content)
        }

        // Step 5: Save files (unless dry-run)
        if !dryRun {
            for _, file := range files {
                if err := saveMarkdownFile(file); err != nil {
                    fmt.Printf("⚠️  Failed to save %s: %v\n", file.Path, err)
                    continue
                }
                fmt.Printf("✅ Updated %s\n", file.Path)
            }
        } else {
            fmt.Println("🔍 Dry run mode - no files were modified")
        }

        // Step 6: Generate navigation hub
        if !dryRun {
            if err := generateNavigationHub(path, files); err != nil {
                return fmt.Errorf("failed to generate navigation hub: %w", err)
            }
            fmt.Printf("🗺️  Generated navigation hub at %s/NAVIGATOR.md\n", path)
        }

        fmt.Println("✨ Migration complete!")
        return nil
    },
}

var validateCmd = &cobra.Command{
    Use:   "validate [path]",
    Short: "Validate Obsidian format in HELIX files",
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        path := "workflows/helix"
        if len(args) > 0 {
            path = args[0]
        }

        fmt.Printf("🔍 Validating Obsidian format in %s...\n", path)

        files, err := scanMarkdownFiles(path)
        if err != nil {
            return fmt.Errorf("failed to scan files: %w", err)
        }

        v := validator.NewValidator()
        errorCount := 0

        for _, file := range files {
            errors := v.ValidateFile(file)
            if len(errors) > 0 {
                fmt.Printf("\n❌ %s:\n", file.Path)
                for _, err := range errors {
                    fmt.Printf("  - %s\n", err)
                    errorCount++
                }
            }
        }

        if errorCount == 0 {
            fmt.Printf("\n✅ All %d files are valid!\n", len(files))
        } else {
            fmt.Printf("\n⚠️  Found %d validation errors\n", errorCount)
        }

        return nil
    },
}

func init() {
    migrateCmd.Flags().BoolP("dry-run", "d", false, "Preview changes without modifying files")
    migrateCmd.Flags().BoolP("validate", "v", false, "Validate after migration")

    obsidianCmd.AddCommand(migrateCmd)
    obsidianCmd.AddCommand(validateCmd)

    rootCmd.AddCommand(obsidianCmd)
}

// scanMarkdownFiles recursively scans for markdown files
func scanMarkdownFiles(root string) ([]*obsidian.MarkdownFile, error) {
    var files []*obsidian.MarkdownFile

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !strings.HasSuffix(path, ".md") {
            return nil
        }

        content, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }

        file := &obsidian.MarkdownFile{
            Path:    path,
            Content: string(content),
        }

        // Parse existing frontmatter if present
        if fm, content := extractFrontmatter(string(content)); fm != nil {
            file.Frontmatter = fm
            file.Content = content
        }

        files = append(files, file)
        return nil
    })

    return files, err
}

// extractFrontmatter extracts existing YAML frontmatter from content
func extractFrontmatter(content string) (*obsidian.Frontmatter, string) {
    if !strings.HasPrefix(content, "---\n") {
        return nil, content
    }

    endIdx := strings.Index(content[4:], "\n---\n")
    if endIdx == -1 {
        return nil, content
    }

    yamlContent := content[4 : endIdx+4]
    remainingContent := content[endIdx+9:]

    var fm obsidian.Frontmatter
    if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
        return nil, content
    }

    return &fm, remainingContent
}

// saveMarkdownFile saves a markdown file with frontmatter
func saveMarkdownFile(file *obsidian.MarkdownFile) error {
    var content strings.Builder

    // Write frontmatter
    if file.Frontmatter != nil {
        yamlBytes, err := yaml.Marshal(file.Frontmatter)
        if err != nil {
            return err
        }
        content.WriteString("---\n")
        content.WriteString(string(yamlBytes))
        content.WriteString("---\n\n")
    }

    // Write content
    content.WriteString(file.Content)

    return ioutil.WriteFile(file.Path, []byte(content.String()), 0644)
}

// generateNavigationHub creates a navigation hub file
func generateNavigationHub(basePath string, files []*obsidian.MarkdownFile) error {
    hub := generateHubContent(files)
    hubPath := filepath.Join(basePath, "NAVIGATOR.md")
    return ioutil.WriteFile(hubPath, []byte(hub), 0644)
}

func generateHubContent(files []*obsidian.MarkdownFile) string {
    var content strings.Builder

    // Write frontmatter
    content.WriteString(`---
title: "HELIX Workflow Navigator"
type: navigation-hub
tags:
  - helix/core
  - helix/navigation
  - dashboard
created: ` + time.Now().Format("2006-01-02") + `
updated: ` + time.Now().Format("2006-01-02") + `
---

# HELIX Workflow Navigator

## 🔄 Workflow Phases

1. [[Frame Phase]] - Define the problem and establish context
2. [[Design Phase]] - Architect the solution approach
3. [[Test Phase]] - Write failing tests (Red phase)
4. [[Build Phase]] - Implement code to pass tests (Green phase)
5. [[Deploy Phase]] - Release to production with monitoring
6. [[Iterate Phase]] - Learn and improve for next cycle

## 📋 Quick Actions

- [[Create Feature Specification]]
- [[Write User Stories]]
- [[Design Technical Architecture]]
- [[Write Test Suite]]
- [[Implement Solution]]
- [[Deploy to Production]]

## 📚 Artifacts by Phase

### Frame Phase
`)

    // Group files by phase and type
    phaseArtifacts := groupFilesByPhase(files)

    for _, phase := range []string{"frame", "design", "test", "build", "deploy", "iterate"} {
        if artifacts, ok := phaseArtifacts[phase]; ok && len(artifacts) > 0 {
            content.WriteString(fmt.Sprintf("\n### %s Phase\n", strings.Title(phase)))
            for _, file := range artifacts {
                if file.Frontmatter != nil {
                    content.WriteString(fmt.Sprintf("- [[%s]]\n", file.Frontmatter.Title))
                }
            }
        }
    }

    content.WriteString(`

## 🏷️ Browse by Tags

\`\`\`dataview
TABLE file.name as "Document", type as "Type", phase as "Phase"
FROM #helix
SORT phase, type
\`\`\`

## 🔍 Search Helpers

### Find by Complexity
- #helix/complexity/simple
- #helix/complexity/moderate
- #helix/complexity/complex

### Find by Type
- #helix/artifact/template
- #helix/artifact/prompt
- #helix/artifact/example

### Find by Phase
- #helix/phase/frame
- #helix/phase/design
- #helix/phase/test
- #helix/phase/build
- #helix/phase/deploy
- #helix/phase/iterate
`)

    return content.String()
}

func groupFilesByPhase(files []*obsidian.MarkdownFile) map[string][]*obsidian.MarkdownFile {
    grouped := make(map[string][]*obsidian.MarkdownFile)

    for _, file := range files {
        if file.Frontmatter != nil && file.Frontmatter.Phase != "" {
            grouped[file.Frontmatter.Phase] = append(grouped[file.Frontmatter.Phase], file)
        }
    }

    return grouped
}
```

### Step 7: Validation Framework

```go
// cli/internal/obsidian/validator/validator.go
package validator

import (
    "fmt"
    "strings"

    "ddx/cli/internal/obsidian"
)

// ValidationError represents a validation error
type ValidationError struct {
    File    string
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s - %s", e.File, e.Field, e.Message)
}

// Validator validates Obsidian format in markdown files
type Validator struct {
    requiredFields map[obsidian.FileType][]string
    validTags      map[string]bool
}

// NewValidator creates a new validator
func NewValidator() *Validator {
    return &Validator{
        requiredFields: map[obsidian.FileType][]string{
            obsidian.FileTypePhase: {"title", "type", "phase_id", "phase_number", "tags"},
            obsidian.FileTypeArtifact: {"title", "type", "artifact_category", "phase", "tags"},
            obsidian.FileTypeTemplate: {"title", "type", "tags"},
            obsidian.FileTypeEnforcer: {"title", "type", "phase", "tags"},
        },
        validTags: initializeValidTags(),
    }
}

// ValidateFile validates a single markdown file
func (v *Validator) ValidateFile(file *obsidian.MarkdownFile) []ValidationError {
    var errors []ValidationError

    // Check for frontmatter
    if file.Frontmatter == nil {
        errors = append(errors, ValidationError{
            File:    file.Path,
            Field:   "frontmatter",
            Message: "missing frontmatter",
        })
        return errors
    }

    fm := file.Frontmatter

    // Check required fields
    if required, ok := v.requiredFields[file.FileType]; ok {
        for _, field := range required {
            if !v.hasField(fm, field) {
                errors = append(errors, ValidationError{
                    File:    file.Path,
                    Field:   field,
                    Message: "required field missing",
                })
            }
        }
    }

    // Validate tags
    for _, tag := range fm.Tags {
        if !v.isValidTag(tag) {
            errors = append(errors, ValidationError{
                File:    file.Path,
                Field:   "tags",
                Message: fmt.Sprintf("invalid tag: %s", tag),
            })
        }
    }

    // Validate wikilinks
    errors = append(errors, v.validateWikilinks(file)...)

    return errors
}

// hasField checks if a frontmatter field is populated
func (v *Validator) hasField(fm *obsidian.Frontmatter, field string) bool {
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
    default:
        return false
    }
}

// isValidTag checks if a tag follows the correct format
func (v *Validator) isValidTag(tag string) bool {
    // Must start with "helix"
    if !strings.HasPrefix(tag, "helix") {
        return false
    }

    // Check against known valid tags
    if v.validTags[tag] {
        return true
    }

    // Allow project-specific tags with prefix
    if strings.HasPrefix(tag, "helix/project/") {
        return true
    }

    return false
}

// validateWikilinks checks that wikilinks are valid
func (v *Validator) validateWikilinks(file *obsidian.MarkdownFile) []ValidationError {
    var errors []ValidationError

    // Extract wikilinks from content
    links := extractWikilinks(file.Content)

    for _, link := range links {
        // Basic validation - check for empty links
        if link == "" {
            errors = append(errors, ValidationError{
                File:    file.Path,
                Field:   "content",
                Message: "empty wikilink found",
            })
        }

        // Check for malformed links
        if strings.Contains(link, "[[") || strings.Contains(link, "]]") {
            errors = append(errors, ValidationError{
                File:    file.Path,
                Field:   "content",
                Message: fmt.Sprintf("malformed wikilink: %s", link),
            })
        }
    }

    return errors
}

// extractWikilinks finds all wikilinks in content
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
        links = append(links, link)

        start = startIdx + endIdx + 2
    }

    return links
}

// initializeValidTags returns the set of valid tags
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
```

## Testing Implementation

### Unit Tests

```go
// cli/internal/obsidian/detector_test.go
package obsidian

import (
    "testing"
)

func TestFileTypeDetector(t *testing.T) {
    detector := NewFileTypeDetector()

    tests := []struct {
        path     string
        expected FileType
    }{
        {"workflows/helix/phases/01-frame/README.md", FileTypePhase},
        {"workflows/helix/phases/01-frame/enforcer.md", FileTypeEnforcer},
        {"workflows/helix/phases/01-frame/artifacts/feature-specification/template.md", FileTypeTemplate},
        {"workflows/helix/phases/01-frame/artifacts/feature-specification/prompt.md", FileTypePrompt},
        {"workflows/helix/phases/01-frame/artifacts/feature-specification/example.md", FileTypeExample},
        {"workflows/helix/coordinator.md", FileTypeCoordinator},
        {"workflows/helix/principles.md", FileTypePrinciple},
        {"workflows/helix/random.md", FileTypeUnknown},
    }

    for _, tt := range tests {
        t.Run(tt.path, func(t *testing.T) {
            result := detector.Detect(tt.path)
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}

func TestGetPhaseFromPath(t *testing.T) {
    tests := []struct {
        path     string
        expected string
    }{
        {"workflows/helix/phases/01-frame/README.md", "frame"},
        {"workflows/helix/phases/02-design/enforcer.md", "design"},
        {"workflows/helix/phases/test/artifacts/test.md", "test"},
        {"workflows/helix/README.md", ""},
    }

    for _, tt := range tests {
        t.Run(tt.path, func(t *testing.T) {
            result := GetPhaseFromPath(tt.path)
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

## Makefile Updates

```makefile
# Add to cli/Makefile

# Obsidian integration targets
.PHONY: obsidian-migrate
obsidian-migrate: build
	./$(BINARY_NAME) obsidian migrate workflows/helix

.PHONY: obsidian-validate
obsidian-validate: build
	./$(BINARY_NAME) obsidian validate workflows/helix

.PHONY: obsidian-test
obsidian-test:
	go test ./internal/obsidian/... -v
	go test ./internal/obsidian/converter/... -v
	go test ./internal/obsidian/validator/... -v

.PHONY: obsidian-dry-run
obsidian-dry-run: build
	./$(BINARY_NAME) obsidian migrate --dry-run workflows/helix
```

## Usage Examples

```bash
# Migrate HELIX workflow to Obsidian format
ddx obsidian migrate

# Preview migration without making changes
ddx obsidian migrate --dry-run

# Migrate a specific directory
ddx obsidian migrate workflows/helix/phases/01-frame

# Validate Obsidian format
ddx obsidian validate

# Validate specific directory
ddx obsidian validate workflows/helix/phases
```

## Rollback Procedure

```bash
#!/bin/bash
# rollback-obsidian.sh

echo "🔄 Rolling back Obsidian migration..."

# Option 1: Restore from backup (if created)
if [ -d "workflows/helix.backup" ]; then
    rm -rf workflows/helix
    mv workflows/helix.backup workflows/helix
    echo "✅ Restored from backup"
    exit 0
fi

# Option 2: Strip frontmatter only
echo "Removing frontmatter from files..."
find workflows/helix -name "*.md" -exec sh -c '
    for file do
        # Remove frontmatter (content between --- markers at start)
        sed -i.bak "/^---$/,/^---$/d" "$file"
        rm "${file}.bak"
    done
' sh {} +

echo "✅ Frontmatter removed"

# Option 3: Revert wikilinks to markdown links
echo "Converting wikilinks back to markdown format..."
# This would require a reverse conversion tool
# ddx obsidian revert --format markdown

echo "✅ Rollback complete"
```

This completes the build implementation for the Obsidian integration feature!