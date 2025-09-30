# SD-015: Meta-Prompt Injection System

**Solution Design ID**: SD-015
**Feature**: FEAT-015 - Meta-Prompt Injection System
**Status**: Draft
**Created**: 2025-01-30
**Updated**: 2025-01-30
**Author**: Core Team

## Overview

This solution design details the technical approach for implementing automatic meta-prompt synchronization to CLAUDE.md, mirroring the existing persona injection system architecture. This eliminates manual prompt updates and ensures all projects use current behavioral guidance from the library.

## Design Goals

1. **Automatic Injection**: Meta-prompts inject on `ddx init`
2. **Automatic Sync**: Meta-prompts re-sync on `ddx update` (always)
3. **Health Monitoring**: `ddx doctor` detects sync drift
4. **Config Integration**: Config changes trigger re-sync
5. **Pattern Consistency**: Use proven persona injection patterns

## Current State Analysis

### Manual Meta-Prompt Management (Before)

```markdown
# CLAUDE.md (Current manual process)

[Project content]

<!-- DDX-META-PROMPT:START -->
<!-- Source: claude/system-prompts/focused.md -->
# System Instructions

[Manually copied from library, often out of date]
<!-- DDX-META-PROMPT:END -->

[More content]
```

### Problems
1. **Manual Updates Required**: Users must copy-paste prompt content
2. **Sync Drift**: Library updates don't propagate to projects
3. **Discovery Gap**: Users don't know when prompts change
4. **Inconsistent Adoption**: Projects use different prompt versions
5. **Asymmetric Systems**: Personas auto-inject, prompts don't

## Proposed Solution

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│ DDx CLI Commands                                             │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ddx init                                                    │
│    └─> injectInitialMetaPrompt()                            │
│          └─> MetaPromptInjector.InjectMetaPrompt()          │
│                                                              │
│  ddx update                                                  │
│    └─> syncLibrary()                                         │
│    └─> syncMetaPrompt()  ← ALWAYS, even if no changes       │
│          └─> MetaPromptInjector.InjectMetaPrompt()          │
│                                                              │
│  ddx doctor                                                  │
│    └─> checkMetaPromptSync()                                │
│          └─> MetaPromptInjector.IsInSync()                  │
│                                                              │
│  ddx config set system.meta_prompt <path>                   │
│    └─> resyncMetaPrompt()                                   │
│          └─> MetaPromptInjector.InjectMetaPrompt()          │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ MetaPromptInjector (cli/internal/metaprompt/injector.go)   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  type MetaPromptInjector interface {                        │
│      InjectMetaPrompt(promptPath string) error              │
│      RemoveMetaPrompt() error                               │
│      IsInSync() (bool, error)                               │
│      GetCurrentMetaPrompt() (string, error)                 │
│  }                                                           │
│                                                              │
│  type MetaPromptInjectorImpl struct {                       │
│      claudeFilePath string                                  │
│      libraryPath    string                                  │
│      workingDir     string                                  │
│  }                                                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ CLAUDE.md (Auto-managed)                                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  [Project content preserved]                                 │
│                                                              │
│  <!-- DDX-META-PROMPT:START -->                             │
│  <!-- Source: claude/system-prompts/focused.md -->          │
│  # System Instructions                                       │
│                                                              │
│  [Auto-injected from .ddx/library/prompts/{source}]         │
│  <!-- DDX-META-PROMPT:END -->                               │
│                                                              │
│  [More project content preserved]                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                            ↑
┌─────────────────────────────────────────────────────────────┐
│ Library Prompts (.ddx/library/prompts/)                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  claude/system-prompts/                                      │
│  ├── focused.md       (default - YAGNI/KISS/DOWITYTD)       │
│  ├── strict.md        (high quality, strict enforcement)     │
│  ├── creative.md      (exploratory, innovative)             │
│  ├── tdd.md           (test-first, red-green-refactor)      │
│  └── ...              (extensible)                           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                            ↑
┌─────────────────────────────────────────────────────────────┐
│ Config (.ddx/config.yaml)                                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  system:                                                     │
│    meta_prompt: "claude/system-prompts/focused.md"          │
│                                                              │
│  # Or disable:                                               │
│  # meta_prompt: null                                         │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Component Details

#### 1. MetaPromptInjector Interface

```go
// cli/internal/metaprompt/injector.go

package metaprompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MetaPromptInjector manages meta-prompt injection into CLAUDE.md
type MetaPromptInjector interface {
	// InjectMetaPrompt injects a meta-prompt from library into CLAUDE.md
	// promptPath is relative to .ddx/library/prompts/ (e.g., "claude/system-prompts/focused.md")
	InjectMetaPrompt(promptPath string) error

	// RemoveMetaPrompt removes the meta-prompt section from CLAUDE.md
	RemoveMetaPrompt() error

	// IsInSync checks if CLAUDE.md prompt matches library version
	// Returns true if in sync, false if out of sync or not found
	IsInSync() (bool, error)

	// GetCurrentMetaPrompt returns the currently injected prompt info
	// Returns source path and content hash, or error if not found
	GetCurrentMetaPrompt() (string, error)
}

// MetaPromptInjectorImpl implements MetaPromptInjector
type MetaPromptInjectorImpl struct {
	claudeFilePath string // Path to CLAUDE.md (typically "CLAUDE.md")
	libraryPath    string // Path to library root (typically ".ddx/library")
	workingDir     string // Working directory for relative path resolution
}

// NewMetaPromptInjector creates a new injector with default paths
func NewMetaPromptInjector() MetaPromptInjector {
	return &MetaPromptInjectorImpl{
		claudeFilePath: "CLAUDE.md",
		libraryPath:    ".ddx/library",
		workingDir:     ".",
	}
}

// NewMetaPromptInjectorWithPaths creates an injector with custom paths
func NewMetaPromptInjectorWithPaths(claudeFile, libraryPath, workingDir string) MetaPromptInjector {
	return &MetaPromptInjectorImpl{
		claudeFilePath: claudeFile,
		libraryPath:    libraryPath,
		workingDir:     workingDir,
	}
}

// Constants for marker handling (same pattern as persona system)
const (
	MetaPromptStartMarker = "<!-- DDX-META-PROMPT:START -->"
	MetaPromptEndMarker   = "<!-- DDX-META-PROMPT:END -->"
	MaxMetaPromptSize     = 1024 * 512 // 512KB max
)
```

#### 2. Injection Algorithm

```go
func (m *MetaPromptInjectorImpl) InjectMetaPrompt(promptPath string) error {
	// 1. Validate prompt path
	if strings.TrimSpace(promptPath) == "" {
		return fmt.Errorf("prompt path cannot be empty")
	}

	// 2. Read prompt content from library
	promptFullPath := filepath.Join(m.workingDir, m.libraryPath, "prompts", promptPath)
	promptContent, err := os.ReadFile(promptFullPath)
	if err != nil {
		return fmt.Errorf("failed to read meta-prompt from %s: %w", promptFullPath, err)
	}

	// 3. Validate size
	if len(promptContent) > MaxMetaPromptSize {
		return fmt.Errorf("meta-prompt too large: %d bytes (max %d)", len(promptContent), MaxMetaPromptSize)
	}

	// 4. Read or create CLAUDE.md
	claudeFullPath := filepath.Join(m.workingDir, m.claudeFilePath)
	var claudeContent string
	if fileExists(claudeFullPath) {
		existing, err := os.ReadFile(claudeFullPath)
		if err != nil {
			return fmt.Errorf("failed to read CLAUDE.md: %w", err)
		}
		claudeContent = string(existing)
	} else {
		// Create default CLAUDE.md if doesn't exist
		claudeContent = "# CLAUDE.md\n\nThis file provides guidance to Claude when working with code in this repository.\n"
	}

	// 5. Remove existing meta-prompt section (if any)
	claudeContent = m.removeMetaPromptSection(claudeContent)

	// 6. Build new meta-prompt section
	metaPromptSection := m.buildMetaPromptSection(string(promptContent), promptPath)

	// 7. Append meta-prompt section to CLAUDE.md
	claudeContent = strings.TrimSpace(claudeContent) + "\n\n" + metaPromptSection

	// 8. Write updated CLAUDE.md
	if err := m.saveCLAUDEFile(claudeContent); err != nil {
		return fmt.Errorf("failed to save CLAUDE.md: %w", err)
	}

	return nil
}

func (m *MetaPromptInjectorImpl) removeMetaPromptSection(content string) string {
	// Same algorithm as persona system (proven reliable)
	startIdx := strings.Index(content, MetaPromptStartMarker)
	if startIdx == -1 {
		return content // No section found
	}

	endIdx := strings.Index(content[startIdx:], MetaPromptEndMarker)
	if endIdx == -1 {
		// Malformed section - remove from start marker to end
		return strings.TrimSpace(content[:startIdx])
	}

	// Calculate absolute end index
	endIdx = startIdx + endIdx + len(MetaPromptEndMarker)

	// Remove the section
	before := strings.TrimRight(content[:startIdx], " \t\n")
	after := strings.TrimLeft(content[endIdx:], " \t\n")

	if before != "" && after != "" {
		return before + "\n\n" + after
	} else if before != "" {
		return before
	} else if after != "" {
		return after
	}
	return ""
}

func (m *MetaPromptInjectorImpl) buildMetaPromptSection(promptContent, sourcePath string) string {
	var sections []string

	sections = append(sections, MetaPromptStartMarker)
	sections = append(sections, fmt.Sprintf("<!-- Source: %s -->", sourcePath))
	sections = append(sections, promptContent)
	sections = append(sections, MetaPromptEndMarker)

	return strings.Join(sections, "\n")
}
```

#### 3. Sync Detection Algorithm

```go
func (m *MetaPromptInjectorImpl) IsInSync() (bool, error) {
	// 1. Read CLAUDE.md
	claudeFullPath := filepath.Join(m.workingDir, m.claudeFilePath)
	claudeContent, err := os.ReadFile(claudeFullPath)
	if err != nil {
		return false, fmt.Errorf("failed to read CLAUDE.md: %w", err)
	}

	// 2. Extract current meta-prompt section
	currentContent, sourcePath, err := m.extractCurrentMetaPrompt(string(claudeContent))
	if err != nil {
		return false, err
	}

	// 3. Read library prompt
	promptFullPath := filepath.Join(m.workingDir, m.libraryPath, "prompts", sourcePath)
	libraryContent, err := os.ReadFile(promptFullPath)
	if err != nil {
		// Library file missing or changed - definitely out of sync
		return false, nil
	}

	// 4. Normalize and compare
	currentNorm := normalizeWhitespace(currentContent)
	libraryNorm := normalizeWhitespace(string(libraryContent))

	return currentNorm == libraryNorm, nil
}

func (m *MetaPromptInjectorImpl) extractCurrentMetaPrompt(content string) (string, string, error) {
	// Find markers
	startIdx := strings.Index(content, MetaPromptStartMarker)
	if startIdx == -1 {
		return "", "", fmt.Errorf("meta-prompt section not found")
	}

	endIdx := strings.Index(content[startIdx:], MetaPromptEndMarker)
	if endIdx == -1 {
		return "", "", fmt.Errorf("malformed meta-prompt section (missing end marker)")
	}

	// Extract section
	endIdx = startIdx + endIdx
	section := content[startIdx:endIdx]

	// Extract source path from comment
	sourcePattern := `<!-- Source: (.+) -->`
	re := regexp.MustCompile(sourcePattern)
	matches := re.FindStringSubmatch(section)
	if len(matches) < 2 {
		return "", "", fmt.Errorf("source path not found in meta-prompt section")
	}
	sourcePath := strings.TrimSpace(matches[1])

	// Extract content (between source comment and end marker)
	sourceCommentIdx := strings.Index(section, "-->")
	if sourceCommentIdx == -1 {
		return "", "", fmt.Errorf("malformed source comment")
	}
	contentStart := sourceCommentIdx + len("-->")
	promptContent := strings.TrimSpace(section[contentStart:])

	return promptContent, sourcePath, nil
}

func normalizeWhitespace(s string) string {
	// Remove all whitespace for comparison (handles formatting differences)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
```

### Integration Points

#### 1. Init Command Integration

```go
// cli/cmd/init.go

func initProject(workingDir string, opts InitOptions) (*InitResult, error) {
	// ... existing init logic ...

	// After library setup and config creation
	if !opts.NoGit {
		// Set up git subtree library
		if err := setupGitSubtreeLibraryPure(localConfig, workingDir); err != nil {
			return nil, NewExitError(1, fmt.Sprintf("Failed to setup library: %v", err))
		}

		// Inject initial meta-prompt
		if err := injectInitialMetaPrompt(localConfig, workingDir); err != nil {
			// Warn but don't fail - meta-prompt is optional
			fmt.Fprintf(os.Stderr, "Warning: Failed to inject meta-prompt: %v\n", err)
		}
	}

	return result, nil
}

func injectInitialMetaPrompt(cfg *config.Config, workingDir string) error {
	// Get meta-prompt path from config (with default)
	promptPath := cfg.GetMetaPrompt()
	if promptPath == "" {
		// Empty means disabled
		return nil
	}

	// Create injector
	injector := metaprompt.NewMetaPromptInjectorWithPaths(
		"CLAUDE.md",
		cfg.Library.Path,
		workingDir,
	)

	// Inject prompt
	if err := injector.InjectMetaPrompt(promptPath); err != nil {
		return fmt.Errorf("failed to inject meta-prompt: %w", err)
	}

	return nil
}
```

#### 2. Update Command Integration

```go
// cli/cmd/update.go

func performUpdate(workingDir string, opts *UpdateOptions) (*UpdateResult, error) {
	// ... existing update logic ...

	// After library sync (even if no changes)
	if err := syncMetaPrompt(cfg, workingDir); err != nil {
		// Warn but don't fail
		fmt.Fprintf(os.Stderr, "Warning: Failed to sync meta-prompt: %v\n", err)
	}

	return result, nil
}

func syncMetaPrompt(cfg *config.Config, workingDir string) error {
	// Get meta-prompt path from config
	promptPath := cfg.GetMetaPrompt()
	if promptPath == "" {
		// Disabled - remove meta-prompt section if exists
		injector := metaprompt.NewMetaPromptInjectorWithPaths(
			"CLAUDE.md",
			cfg.Library.Path,
			workingDir,
		)
		return injector.RemoveMetaPrompt()
	}

	// Create injector and sync
	injector := metaprompt.NewMetaPromptInjectorWithPaths(
		"CLAUDE.md",
		cfg.Library.Path,
		workingDir,
	)

	return injector.InjectMetaPrompt(promptPath)
}
```

#### 3. Doctor Command Integration

```go
// cli/cmd/doctor.go

func runDoctorChecks(workingDir string) ([]HealthCheckResult, error) {
	results := []HealthCheckResult{}

	// ... existing checks ...

	// Add meta-prompt sync check
	metaPromptCheck := checkMetaPromptSync(workingDir)
	results = append(results, metaPromptCheck)

	return results, nil
}

func checkMetaPromptSync(workingDir string) HealthCheckResult {
	cfg, err := config.LoadWithWorkingDir(workingDir)
	if err != nil {
		return HealthCheckResult{
			Name:   "Meta-prompt sync",
			Status: "error",
			Message: fmt.Sprintf("Failed to load config: %v", err),
		}
	}

	promptPath := cfg.GetMetaPrompt()
	if promptPath == "" {
		return HealthCheckResult{
			Name:   "Meta-prompt sync",
			Status: "skipped",
			Message: "Meta-prompt disabled in config",
		}
	}

	injector := metaprompt.NewMetaPromptInjectorWithPaths(
		"CLAUDE.md",
		cfg.Library.Path,
		workingDir,
	)

	inSync, err := injector.IsInSync()
	if err != nil {
		return HealthCheckResult{
			Name:   "Meta-prompt sync",
			Status: "warning",
			Message: fmt.Sprintf("Could not check sync: %v", err),
			Fix:     "Run 'ddx update' to sync meta-prompt",
		}
	}

	if inSync {
		return HealthCheckResult{
			Name:   "Meta-prompt sync",
			Status: "ok",
			Message: "Meta-prompt is in sync with library",
		}
	}

	return HealthCheckResult{
		Name:   "Meta-prompt sync",
		Status: "warning",
		Message: "Meta-prompt is out of sync with library",
		Fix:     "Run 'ddx update' to sync meta-prompt",
	}
}
```

#### 4. Config Command Integration

```go
// cli/cmd/config.go

func setConfigValue(key, value string, workingDir string) error {
	// ... existing config set logic ...

	// After setting value
	if key == "system.meta_prompt" {
		// Re-sync meta-prompt
		if err := resyncMetaPrompt(workingDir); err != nil {
			return fmt.Errorf("failed to re-sync meta-prompt: %w", err)
		}
	}

	return nil
}

func resyncMetaPrompt(workingDir string) error {
	cfg, err := config.LoadWithWorkingDir(workingDir)
	if err != nil {
		return err
	}

	return syncMetaPrompt(cfg, workingDir)
}
```

## Data Model Changes

### Config Schema (Already Exists)

```yaml
# .ddx/config.yaml

version: "1.0"
system:
  meta_prompt: "claude/system-prompts/focused.md"  # Path relative to .ddx/library/prompts/
  # Or disable:
  # meta_prompt: null
```

### CLAUDE.md Structure (After Injection)

```markdown
# CLAUDE.md

## Project Overview
[Project content preserved]

## Development Commands
[Project content preserved]

<!-- DDX-META-PROMPT:START -->
<!-- Source: claude/system-prompts/focused.md -->
# System Instructions

**Execute ONLY what is requested:**

- **YAGNI** (You Aren't Gonna Need It): Implement only specified features.
- **KISS** (Keep It Simple, Stupid): Choose the simplest solution.
- **DOWITYTD** (Do Only What I Told You To Do): Stop when complete.

[Rest of prompt content]
<!-- DDX-META-PROMPT:END -->

[More project content preserved]
```

## API/Interface Changes

### New Package: `cli/internal/metaprompt`

```go
package metaprompt

// Exported interface
type MetaPromptInjector interface { ... }

// Constructor functions
func NewMetaPromptInjector() MetaPromptInjector
func NewMetaPromptInjectorWithPaths(...) MetaPromptInjector

// Error types
type MetaPromptError struct {
	Type    string
	Message string
	Cause   error
}
```

### Modified Commands

**No breaking changes**. All existing commands work identically, with added behavior:

- `ddx init` - Now injects meta-prompt
- `ddx update` - Now syncs meta-prompt
- `ddx doctor` - Now checks meta-prompt sync
- `ddx config set` - Now syncs when `system.meta_prompt` changes

## Performance Considerations

### Injection Performance

**Typical Case**:
- Read CLAUDE.md: ~5ms (typical file ~50KB)
- Read prompt file: ~2ms (typical prompt ~10KB)
- String processing: ~1ms
- Write CLAUDE.md: ~5ms
- **Total: ~13ms** (negligible)

**Edge Cases**:
- Large CLAUDE.md (500KB): ~50ms
- Large prompt (512KB max): ~25ms
- **Worst case: ~75ms** (still acceptable)

### Sync Detection Performance

**Per Check**:
- Read CLAUDE.md: ~5ms
- Extract section: ~1ms
- Read library prompt: ~2ms
- Normalize + compare: ~1ms
- **Total: ~9ms** (negligible)

### Memory Usage

- CLAUDE.md in memory: ~50KB typical
- Prompt content: ~10KB typical
- Processing overhead: ~10KB
- **Total: ~70KB** (minimal)

## Security Considerations

### Input Validation

1. **Prompt Path**:
   - Validate no path traversal (../ etc.)
   - Must be relative to library/prompts/
   - File must exist in library

2. **File Size Limits**:
   - Max prompt size: 512KB
   - Prevents memory exhaustion
   - Prevents malicious large files

3. **Content Validation**:
   - No executable content validation (prompts are text)
   - Marker injection prevention (escape markers in content)

### File Permissions

- CLAUDE.md written with 0644 (readable by all, writable by owner)
- No sensitive data in prompts
- Prompts versioned in git (auditable)

## Error Handling

### Injection Failures

| Scenario | Handling | User Impact |
|----------|----------|-------------|
| Prompt file not found | Error message, suggest `ddx update` | Init/update continues, warning shown |
| CLAUDE.md read-only | Error message, check permissions | Init/update fails with clear error |
| Prompt too large | Error message, size limit | Init/update continues, warning shown |
| Malformed markers | Remove + recreate section | Automatic recovery |

### Sync Detection Failures

| Scenario | Handling | Doctor Output |
|----------|----------|---------------|
| CLAUDE.md missing | Return "not found" | Warning: "CLAUDE.md not found" |
| Prompt missing | Return "out of sync" | Warning: "Library prompt not found" |
| No markers | Return "not found" | Info: "No meta-prompt injected" |
| Malformed section | Return "error" | Warning: "Malformed meta-prompt section" |

### Recovery Strategies

1. **Corrupted CLAUDE.md**: Run `ddx update` to re-inject
2. **Missing library**: Run `ddx update` to sync library
3. **Wrong config**: Run `ddx config set system.meta_prompt <path>`
4. **Complete reset**: Delete CLAUDE.md, run `ddx init --force`

## Testing Strategy

### Unit Tests (`cli/internal/metaprompt/injector_test.go`)

```go
func TestInjectMetaPrompt(t *testing.T) {
	tests := []struct {
		name          string
		existingFile  string
		promptPath    string
		promptContent string
		expectError   bool
		expectContent string
	}{
		{
			name:          "inject into new file",
			existingFile:  "",
			promptPath:    "claude/system-prompts/focused.md",
			promptContent: "# Test Prompt\nContent here",
			expectError:   false,
			expectContent: "<!-- DDX-META-PROMPT:START -->",
		},
		{
			name:          "inject into existing file",
			existingFile:  "# CLAUDE.md\n\nExisting content",
			promptPath:    "claude/system-prompts/focused.md",
			promptContent: "# Test Prompt",
			expectError:   false,
			expectContent: "Existing content",
		},
		{
			name:          "replace existing meta-prompt",
			existingFile:  "Content\n<!-- DDX-META-PROMPT:START -->\nOld\n<!-- DDX-META-PROMPT:END -->",
			promptPath:    "claude/system-prompts/strict.md",
			promptContent: "New prompt",
			expectError:   false,
			expectContent: "New prompt",
		},
		{
			name:          "prompt file not found",
			existingFile:  "",
			promptPath:    "nonexistent/prompt.md",
			expectError:   true,
		},
		{
			name:          "prompt too large",
			existingFile:  "",
			promptPath:    "claude/system-prompts/huge.md",
			promptContent: strings.Repeat("x", MaxMetaPromptSize+1),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
		})
	}
}

func TestIsInSync(t *testing.T) {
	tests := []struct {
		name          string
		claudeContent string
		libraryContent string
		expectInSync  bool
		expectError   bool
	}{
		{
			name:           "in sync",
			claudeContent:  buildCLAUDEWithPrompt("Test prompt"),
			libraryContent: "Test prompt",
			expectInSync:   true,
			expectError:    false,
		},
		{
			name:           "out of sync",
			claudeContent:  buildCLAUDEWithPrompt("Old prompt"),
			libraryContent: "New prompt",
			expectInSync:   false,
			expectError:    false,
		},
		{
			name:          "no meta-prompt section",
			claudeContent: "# CLAUDE.md\n\nNo meta-prompt",
			expectInSync:  false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
		})
	}
}

func TestRemoveMetaPrompt(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectContent string
	}{
		{
			name:          "remove existing section",
			input:         "Before\n<!-- DDX-META-PROMPT:START -->\nPrompt\n<!-- DDX-META-PROMPT:END -->\nAfter",
			expectContent: "Before\n\nAfter",
		},
		{
			name:          "no section to remove",
			input:         "# CLAUDE.md\n\nNo meta-prompt here",
			expectContent: "# CLAUDE.md\n\nNo meta-prompt here",
		},
		{
			name:          "malformed section",
			input:         "Content\n<!-- DDX-META-PROMPT:START -->\nNo end marker",
			expectContent: "Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
		})
	}
}
```

### Integration Tests (`cli/cmd/*_test.go`)

```go
func TestInitCommand_InjectsMetaPrompt(t *testing.T) {
	// Setup test environment
	testEnv := setupTestEnvironment(t)
	defer testEnv.Cleanup()

	// Run init
	err := runInit(testEnv.Dir, InitOptions{})
	require.NoError(t, err)

	// Verify CLAUDE.md exists with meta-prompt
	claudeContent, err := os.ReadFile(filepath.Join(testEnv.Dir, "CLAUDE.md"))
	require.NoError(t, err)
	require.Contains(t, string(claudeContent), "<!-- DDX-META-PROMPT:START -->")
	require.Contains(t, string(claudeContent), "<!-- Source: claude/system-prompts/focused.md -->")
}

func TestUpdateCommand_SyncsMetaPrompt(t *testing.T) {
	// Setup test environment with old prompt
	testEnv := setupTestEnvironmentWithOldPrompt(t)
	defer testEnv.Cleanup()

	// Update library prompt
	updateLibraryPrompt(testEnv, "new prompt content")

	// Run update
	err := runUpdate(testEnv.Dir, UpdateOptions{})
	require.NoError(t, err)

	// Verify CLAUDE.md has new prompt
	claudeContent, err := os.ReadFile(filepath.Join(testEnv.Dir, "CLAUDE.md"))
	require.NoError(t, err)
	require.Contains(t, string(claudeContent), "new prompt content")
}

func TestDoctorCommand_DetectsOutOfSync(t *testing.T) {
	// Setup test environment with out-of-sync prompt
	testEnv := setupTestEnvironmentOutOfSync(t)
	defer testEnv.Cleanup()

	// Run doctor
	results, err := runDoctor(testEnv.Dir)
	require.NoError(t, err)

	// Find meta-prompt check
	var metaPromptCheck *HealthCheckResult
	for _, result := range results {
		if result.Name == "Meta-prompt sync" {
			metaPromptCheck = &result
			break
		}
	}

	require.NotNil(t, metaPromptCheck)
	require.Equal(t, "warning", metaPromptCheck.Status)
	require.Contains(t, metaPromptCheck.Message, "out of sync")
	require.Contains(t, metaPromptCheck.Fix, "ddx update")
}
```

### Acceptance Tests

```go
func TestAcceptance_MetaPromptLifecycle(t *testing.T) {
	// Full user workflow test
	t.Run("init -> check -> update -> check cycle", func(t *testing.T) {
		// 1. Init project
		testEnv := setupTestGitRepo(t)
		defer testEnv.Cleanup()

		err := runInit(testEnv.Dir, InitOptions{})
		require.NoError(t, err)

		// 2. Verify meta-prompt injected
		injector := metaprompt.NewMetaPromptInjectorWithPaths("CLAUDE.md", ".ddx/library", testEnv.Dir)
		inSync, err := injector.IsInSync()
		require.NoError(t, err)
		require.True(t, inSync, "Meta-prompt should be in sync after init")

		// 3. Manually edit meta-prompt section to simulate drift
		simulateLibraryUpdate(testEnv)

		// 4. Doctor should detect out of sync
		results, err := runDoctor(testEnv.Dir)
		require.NoError(t, err)
		metaPromptCheck := findCheck(results, "Meta-prompt sync")
		require.Equal(t, "warning", metaPromptCheck.Status)

		// 5. Update should fix sync
		err = runUpdate(testEnv.Dir, UpdateOptions{})
		require.NoError(t, err)

		// 6. Verify back in sync
		inSync, err = injector.IsInSync()
		require.NoError(t, err)
		require.True(t, inSync, "Meta-prompt should be in sync after update")
	})
}
```

## Migration Plan

### Rollout Strategy

**Phase 1: Implementation (Week 1)**
- Implement `cli/internal/metaprompt/injector.go`
- Add unit tests
- Integrate with init command
- Test in development

**Phase 2: Integration (Week 1)**
- Integrate with update command
- Integrate with doctor command
- Integrate with config command
- Add integration tests

**Phase 3: Testing (Week 1)**
- Run full test suite
- Manual testing across platforms
- Edge case testing
- Performance validation

**Phase 4: Release (Week 2)**
- Merge to main
- Release new version
- Monitor for issues
- Gather user feedback

### Backward Compatibility

- **New projects**: Auto-inject meta-prompt on `ddx init`
- **Existing projects**: First `ddx update` injects meta-prompt
- **Manual markers**: If users already have markers, content is replaced
- **No markers**: Meta-prompt section appended to CLAUDE.md

**No breaking changes**. All existing functionality preserved.

### Rollback Plan

If issues occur:
1. Revert changes to CLI commands
2. Remove `cli/internal/metaprompt` package
3. Restore previous command behavior
4. Git commit revert

Time to rollback: <10 minutes

## Monitoring & Observability

### Success Metrics

- **Injection Success Rate**: Track `ddx init` meta-prompt injection success
- **Sync Detection Accuracy**: Track false positives/negatives in doctor
- **Update Success Rate**: Track `ddx update` meta-prompt sync success
- **User Adoption**: Count projects with auto-injected prompts

### Logging

```go
// Log all meta-prompt operations
log.Info("Injecting meta-prompt",
	"path", promptPath,
	"source", sourcePath,
	"size", len(promptContent))

log.Warn("Meta-prompt injection failed",
	"path", promptPath,
	"error", err)

log.Info("Meta-prompt sync check",
	"in_sync", inSync,
	"source", sourcePath)
```

### Error Tracking

- Track injection failures by error type
- Track sync detection failures
- Track file permission issues
- Report via telemetry (if enabled)

## Open Questions

1. **Q**: Should we validate prompt content format (e.g., must be markdown)?
   **A**: No, keep validation minimal. Just size check.

2. **Q**: Should we support multiple meta-prompts in CLAUDE.md?
   **A**: No, single meta-prompt only. Keep it simple.

3. **Q**: What if users manually edit injected content?
   **A**: Update will overwrite. Document this behavior clearly.

4. **Q**: Should sync check be strict or lenient on whitespace?
   **A**: Lenient. Normalize whitespace for comparison.

5. **Q**: Should we add a --skip-meta-prompt flag to commands?
   **A**: No, use config: `system.meta_prompt: null` to disable.

## Dependencies

### Internal
- Config system with `system.meta_prompt` support (✓ already exists)
- Library management system (✓ already exists)
- Persona injection system (✓ reference implementation)

### External
- None

## Related Documents

- **Feature Spec**: `docs/helix/01-frame/features/FEAT-015-meta-prompt-injection.md`
- **User Story**: `docs/helix/01-frame/user-stories/US-045-meta-prompt-auto-sync.md`
- **Reference**: `cli/internal/persona/claude.go` (persona injection)
- **Config**: `cli/internal/config/types.go` (SystemConfig)

## Decision Log

### Decision 1: Injection Timing
**Options**:
1. Only on explicit command
2. On init + update
3. On every command

**Choice**: On init + update + config change
**Rationale**: Balance between automatic and non-intrusive

### Decision 2: Sync Detection Method
**Options**:
1. Hash-based comparison
2. Content comparison
3. Timestamp-based

**Choice**: Content comparison with whitespace normalization
**Rationale**: Most reliable, handles formatting changes

### Decision 3: Error Handling Strategy
**Options**:
1. Fail command on injection error
2. Warn and continue
3. Silent ignore

**Choice**: Warn and continue
**Rationale**: Meta-prompt is enhancement, shouldn't block operations

### Decision 4: Marker Pattern
**Options**:
1. Create new marker format
2. Use persona marker pattern
3. Use generic DDX markers

**Choice**: Use same pattern as persona system
**Rationale**: Consistency, proven reliability

## Next Steps

1. **Test Phase**: Write comprehensive test suite (TS-015)
2. **Build Phase**: Implement injector and integrations
3. **Validation**: Test on multiple projects
4. **Documentation**: Update all user docs
5. **Release**: Deploy with monitoring

---
*Status: Ready for test phase*