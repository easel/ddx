---
title: "Test Specification - Obsidian Integration for HELIX"
type: test-specification
feature_id: FEAT-014
workflow_phase: test
artifact_type: test-specification
tags:
  - helix/test
  - helix/artifact/test
  - helix/phase/test
  - obsidian
  - frontmatter
  - wikilinks
related:
  - "[[FEAT-014-obsidian-integration]]"
  - "[[FEAT-014-technical-design]]"
  - "[[FEAT-014-test-procedures]]"
status: draft
priority: P1
created: 2025-01-18
updated: 2025-01-18
---

# Test Specification: [[FEAT-014]] Obsidian Integration

## Test Strategy Overview

This test specification defines comprehensive test scenarios for the Obsidian integration feature. Following Test-Driven Development principles, these tests must be written and failing BEFORE any implementation begins.

### Test Categories

1. **Unit Tests** - Individual component functionality
2. **Integration Tests** - Component interaction testing
3. **End-to-End Tests** - Complete workflow validation
4. **CLI Tests** - Command-line interface testing
5. **Validation Tests** - Schema and format compliance

## Unit Test Specifications

### 1. File Type Detection Tests

#### Test Suite: `TestFileTypeDetector`

```go
// Path: cli/internal/obsidian/detector_test.go

func TestFileTypeDetector_DetectHelixFiles(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        expected FileType
    }{
        // Phase detection
        {
            name:     "Frame phase README",
            path:     "workflows/helix/phases/01-frame/README.md",
            expected: FileTypePhase,
        },
        {
            name:     "Design phase enforcer",
            path:     "workflows/helix/phases/02-design/enforcer.md",
            expected: FileTypeEnforcer,
        },

        // Artifact detection
        {
            name:     "Feature specification template",
            path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md",
            expected: FileTypeTemplate,
        },
        {
            name:     "User stories prompt",
            path:     "workflows/helix/phases/01-frame/artifacts/user-stories/prompt.md",
            expected: FileTypePrompt,
        },

        // Workflow coordination
        {
            name:     "HELIX coordinator",
            path:     "workflows/helix/coordinator.md",
            expected: FileTypeCoordinator,
        },
        {
            name:     "HELIX principles",
            path:     "workflows/helix/principles.md",
            expected: FileTypePrinciple,
        },

        // Feature specifications
        {
            name:     "Feature specification in docs",
            path:     "docs/01-frame/features/FEAT-014-obsidian-integration.md",
            expected: FileTypeFeature,
        },

        // Non-HELIX files should not be detected
        {
            name:     "Regular README",
            path:     "README.md",
            expected: FileTypeUnknown,
        },
        {
            name:     "Random markdown file",
            path:     "docs/other/random.md",
            expected: FileTypeUnknown,
        },
    }

    detector := NewFileTypeDetector()
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := detector.Detect(tt.path)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestFileTypeDetector_PhaseExtraction(t *testing.T) {
    tests := []struct {
        path     string
        expected string
    }{
        {"workflows/helix/phases/01-frame/README.md", "frame"},
        {"workflows/helix/phases/02-design/enforcer.md", "design"},
        {"docs/01-frame/features/FEAT-014.md", "frame"},
        {"docs/02-design/FEAT-014-technical-design.md", "design"},
        {"docs/03-test/test-spec.md", "test"},
        {"docs/04-build/implementation.md", "build"},
        {"docs/05-deploy/deployment.md", "deploy"},
        {"docs/06-iterate/retrospective.md", "iterate"},
        {"random/path/file.md", ""},
    }

    for _, tt := range tests {
        result := GetPhaseFromPath(tt.path)
        assert.Equal(t, tt.expected, result, "Path: %s", tt.path)
    }
}
```

#### Expected Behavior:
- **MUST** correctly identify HELIX file types based on path patterns
- **MUST** extract phase information from directory structure
- **MUST** return `FileTypeUnknown` for non-HELIX files
- **MUST** handle edge cases like missing directories gracefully

### 2. Frontmatter Generation Tests

#### Test Suite: `TestFrontmatterGenerator`

```go
// Path: cli/internal/obsidian/generator_test.go

func TestFrontmatterGenerator_PhaseFiles(t *testing.T) {
    generator := NewFrontmatterGenerator()

    // Test phase README frontmatter generation
    file := &MarkdownFile{
        Path:     "workflows/helix/phases/01-frame/README.md",
        Content:  "# Frame Phase - Problem Definition\n\nThis phase...",
        FileType: FileTypePhase,
    }

    fm, err := generator.Generate(file)
    require.NoError(t, err)

    // Required fields
    assert.Equal(t, "Frame Phase - Problem Definition", fm.Title)
    assert.Equal(t, "phase", fm.Type)
    assert.Equal(t, "frame", fm.PhaseID)
    assert.Equal(t, 1, fm.PhaseNum)
    assert.Contains(t, fm.Tags, "helix")
    assert.Contains(t, fm.Tags, "helix/phase")
    assert.Contains(t, fm.Tags, "helix/phase/frame")

    // Phase-specific fields
    assert.Equal(t, "[[Design Phase]]", fm.NextPhase)
    assert.Nil(t, fm.PrevPhase) // First phase
    assert.NotNil(t, fm.Gates)
    assert.NotNil(t, fm.Artifacts)

    // Timestamps
    assert.False(t, fm.Created.IsZero())
    assert.False(t, fm.Updated.IsZero())
}

func TestFrontmatterGenerator_TemplateFiles(t *testing.T) {
    generator := NewFrontmatterGenerator()

    file := &MarkdownFile{
        Path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md",
        Content:  "# Feature Specification Template\n\nUse this template...",
        FileType: FileTypeTemplate,
    }

    fm, err := generator.Generate(file)
    require.NoError(t, err)

    assert.Equal(t, "template", fm.Type)
    assert.Equal(t, "frame", fm.Phase)
    assert.Equal(t, "feature-specification", fm.ArtifactCategory)
    assert.Contains(t, fm.Tags, "helix/artifact")
    assert.Contains(t, fm.Tags, "helix/artifact/feature-specification")
    assert.Equal(t, "30-60 minutes", fm.TimeEstimate)
}

func TestFrontmatterGenerator_FeatureFiles(t *testing.T) {
    generator := NewFrontmatterGenerator()

    file := &MarkdownFile{
        Path:     "docs/01-frame/features/FEAT-014-obsidian-integration.md",
        Content:  "# Feature Specification: FEAT-014 - Obsidian Integration\n\n**Priority**: P1\n**Owner**: Platform Team",
        FileType: FileTypeFeature,
    }

    fm, err := generator.Generate(file)
    require.NoError(t, err)

    assert.Equal(t, "FEAT-014", fm.FeatureID)
    assert.Equal(t, "P1", fm.Priority)
    assert.Equal(t, "Platform Team", fm.Owner)
    assert.Equal(t, "frame", fm.WorkflowPhase)
    assert.Equal(t, "feature-specification", fm.ArtifactType)
}
```

#### Expected Behavior:
- **MUST** generate valid YAML frontmatter for all file types
- **MUST** extract titles from content when present
- **MUST** fall back to path-based titles when content title missing
- **MUST** generate hierarchical tags based on file type and location
- **MUST** include required fields for each file type schema
- **MUST** extract metadata fields from content (Priority, Owner, etc.)

### 3. Wikilink Conversion Tests

#### Test Suite: `TestLinkConverter`

```go
// Path: cli/internal/obsidian/converter/links_test.go

func TestLinkConverter_MarkdownToWikilink(t *testing.T) {
    converter := NewLinkConverter()

    // Build test file index
    files := []*MarkdownFile{
        {Path: "workflows/helix/phases/02-design/README.md", Frontmatter: &Frontmatter{Title: "Design Phase"}},
        {Path: "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md", Frontmatter: &Frontmatter{Title: "Feature Specification Template"}},
    }
    converter.BuildIndex(files)

    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Phase reference link",
            input:    "See the [Design Phase](../02-design/README.md) for details",
            expected: "See the [[Design Phase]] for details",
        },
        {
            name:     "Template link with alias",
            input:    "Use the [feature spec template](./artifacts/feature-specification/template.md)",
            expected: "Use the [[Feature Specification Template|feature spec template]]",
        },
        {
            name:     "External link preserved",
            input:    "Visit [GitHub](https://github.com) for more info",
            expected: "Visit [GitHub](https://github.com) for more info",
        },
        {
            name:     "Anchor link preserved",
            input:    "Jump to [section](#implementation)",
            expected: "Jump to [section](#implementation)",
        },
        {
            name:     "Email link preserved",
            input:    "Contact [support](mailto:help@example.com)",
            expected: "Contact [support](mailto:help@example.com)",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := converter.ConvertContent(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestLinkConverter_PhaseReferences(t *testing.T) {
    converter := NewLinkConverter()

    tests := []struct {
        input    string
        expected string
    }{
        {"Complete the Frame phase first", "Complete the [[Frame Phase|Frame phase]] first"},
        {"During Design phase planning", "During [[Design Phase|Design phase]] planning"},
        {"The Build phase implementation", "The [[Build Phase|Build phase]] implementation"},
        {"Already in [[Frame Phase]] work", "Already in [[Frame Phase]] work"}, // Don't double-convert
    }

    for _, tt := range tests {
        result := converter.ConvertContent(tt.input)
        assert.Equal(t, tt.expected, result)
    }
}

func TestLinkConverter_PreventDoubleConversion(t *testing.T) {
    converter := NewLinkConverter()

    input := "See [[Design Phase]] and visit [[Frame Phase]] for details"
    result := converter.ConvertContent(input)

    // Should remain unchanged - already wikilinks
    assert.Equal(t, input, result)
}
```

#### Expected Behavior:
- **MUST** convert relative markdown links to wikilinks
- **MUST** preserve external links unchanged
- **MUST** preserve anchor links unchanged
- **MUST** prevent double-conversion of existing wikilinks
- **MUST** use aliases when link text differs from target title
- **MUST** resolve paths correctly in file index

## Integration Test Specifications

### 4. CLI Command Tests

#### Test Suite: `TestObsidianCommand`

```go
// Path: cli/cmd/obsidian_test.go

func TestObsidianMigrate_DryRun(t *testing.T) {
    // Setup test directory with sample HELIX files
    testDir := setupTestHelixWorkflow(t)

    // Run dry-run migration
    cmd := exec.Command("ddx", "obsidian", "migrate", "--dry-run", "--path", testDir)
    output, err := cmd.CombinedOutput()
    require.NoError(t, err)

    // Verify output shows what would be done
    assert.Contains(t, string(output), "Would add frontmatter to")
    assert.Contains(t, string(output), "Would convert links in")
    assert.Contains(t, string(output), "Would generate navigation hub")

    // Verify no files were actually modified
    assertNoFilesModified(t, testDir)
}

func TestObsidianMigrate_FullMigration(t *testing.T) {
    testDir := setupTestHelixWorkflow(t)

    // Run full migration
    cmd := exec.Command("ddx", "obsidian", "migrate", "--path", testDir)
    err := cmd.Run()
    require.NoError(t, err)

    // Verify all markdown files have frontmatter
    err = filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
        if strings.HasSuffix(path, ".md") && isHelixFile(path) {
            content, err := os.ReadFile(path)
            require.NoError(t, err)

            // Must start with frontmatter
            assert.True(t, strings.HasPrefix(string(content), "---\n"))

            // Must contain required fields
            assert.Contains(t, string(content), "title:")
            assert.Contains(t, string(content), "type:")
            assert.Contains(t, string(content), "tags:")
            assert.Contains(t, string(content), "created:")
            assert.Contains(t, string(content), "updated:")
        }
        return nil
    })
    require.NoError(t, err)

    // Verify navigation hub was created
    navPath := filepath.Join(testDir, "HELIX-NAVIGATOR.md")
    assert.FileExists(t, navPath)
}

func TestObsidianValidate(t *testing.T) {
    testDir := setupMigratedHelixWorkflow(t)

    // Run validation
    cmd := exec.Command("ddx", "obsidian", "validate", "--path", testDir)
    output, err := cmd.CombinedOutput()
    require.NoError(t, err)

    // Should report no errors for valid files
    assert.Contains(t, string(output), "✅ All files valid")
    assert.NotContains(t, string(output), "❌ Validation failed")
}

func TestObsidianValidate_InvalidFiles(t *testing.T) {
    testDir := setupInvalidHelixFiles(t)

    cmd := exec.Command("ddx", "obsidian", "validate", "--path", testDir)
    output, err := cmd.CombinedOutput()

    // Should exit with error code
    assert.Error(t, err)

    // Should report specific validation errors
    assert.Contains(t, string(output), "❌ Missing required field")
    assert.Contains(t, string(output), "❌ Invalid tag format")
    assert.Contains(t, string(output), "❌ Broken wikilink")
}
```

#### Expected Behavior:
- **MUST** provide `migrate` command with dry-run option
- **MUST** provide `validate` command for checking compliance
- **MUST** provide `revert` command for rollback
- **MUST** support path specification for non-standard locations
- **MUST** provide clear progress indicators
- **MUST** generate comprehensive validation reports

### 5. Navigation Hub Tests

#### Test Suite: `TestNavigationHub`

```go
// Path: cli/internal/obsidian/navigation_test.go

func TestNavigationHub_Generation(t *testing.T) {
    files := setupTestWorkflowFiles()

    hub, err := GenerateNavigationHub(files)
    require.NoError(t, err)

    // Verify phase organization
    assert.Len(t, hub.Phases, 6) // All 6 HELIX phases

    framePhase := findPhase(hub.Phases, "frame")
    require.NotNil(t, framePhase)
    assert.Equal(t, 1, framePhase.Number)
    assert.Equal(t, "Design", framePhase.Next)
    assert.Empty(t, framePhase.Previous)

    // Verify artifact categorization
    assert.Contains(t, hub.Artifacts, "feature-specification")
    assert.Contains(t, hub.Artifacts, "technical-design")
    assert.Contains(t, hub.Artifacts, "test-specification")

    // Verify tag tree structure
    assert.NotNil(t, hub.Tags.Get("helix/phase/frame"))
    assert.NotNil(t, hub.Tags.Get("helix/artifact/template"))
    assert.True(t, len(hub.Tags.GetFilesByTag("helix/phase/frame")) > 0)
}

func TestNavigationHub_MarkdownGeneration(t *testing.T) {
    hub := &NavigationHub{
        Phases: []*PhaseInfo{
            {ID: "frame", Number: 1, Title: "Frame Phase", Status: "completed"},
            {ID: "design", Number: 2, Title: "Design Phase", Status: "in_progress"},
        },
        Artifacts: map[string][]*ArtifactInfo{
            "feature-specification": {
                {Title: "Feature Spec Template", Path: "template.md", Complexity: "moderate"},
            },
        },
    }

    markdown := hub.GenerateMarkdown()

    // Verify structure
    assert.Contains(t, markdown, "# HELIX Workflow Navigator")
    assert.Contains(t, markdown, "## Phase Progress")
    assert.Contains(t, markdown, "1. [[Frame Phase]] ✅")
    assert.Contains(t, markdown, "2. [[Design Phase]] 🚧")
    assert.Contains(t, markdown, "### Feature Specifications")
    assert.Contains(t, markdown, "[[Feature Spec Template]]")
}
```

#### Expected Behavior:
- **MUST** organize phases with proper progression indicators
- **MUST** categorize artifacts by type and phase
- **MUST** build hierarchical tag structure
- **MUST** generate navigable markdown with wikilinks
- **MUST** include status indicators and progress tracking

## End-to-End Test Specifications

### 6. Complete Workflow Tests

#### Test Suite: `TestObsidianWorkflow`

```go
// Path: cli/test/obsidian_e2e_test.go

func TestCompleteObsidianWorkflow(t *testing.T) {
    // Skip if not in integration test mode
    if testing.Short() {
        t.Skip("Skipping end-to-end test in short mode")
    }

    // Create temporary project with HELIX workflow
    projectDir := createTestProject(t)
    defer os.RemoveAll(projectDir)

    // Initialize DDX in project
    runCommand(t, projectDir, "ddx", "init", "--workflow", "helix")

    // Add some feature specifications
    createFeatureSpec(t, projectDir, "FEAT-001", "User Authentication")
    createFeatureSpec(t, projectDir, "FEAT-002", "Data Export")

    // Migrate to Obsidian format
    runCommand(t, projectDir, "ddx", "obsidian", "migrate")

    // Validate migration
    runCommand(t, projectDir, "ddx", "obsidian", "validate")

    // Verify specific behaviors
    t.Run("frontmatter_present", func(t *testing.T) {
        assertAllMarkdownHasFrontmatter(t, projectDir)
    })

    t.Run("wikilinks_converted", func(t *testing.T) {
        assertWikilinksPresent(t, projectDir)
    })

    t.Run("navigation_hub_exists", func(t *testing.T) {
        navPath := filepath.Join(projectDir, "workflows/helix/HELIX-NAVIGATOR.md")
        assert.FileExists(t, navPath)

        content := readFile(t, navPath)
        assert.Contains(t, content, "[[FEAT-001]]")
        assert.Contains(t, content, "[[FEAT-002]]")
    })

    t.Run("obsidian_compatibility", func(t *testing.T) {
        // Verify format is compatible with Obsidian
        assertObsidianCompatible(t, projectDir)
    })
}

func TestObsidianRevert(t *testing.T) {
    projectDir := createMigratedProject(t)
    defer os.RemoveAll(projectDir)

    // Capture original state
    originalFiles := captureFileStates(t, projectDir)

    // Revert migration
    runCommand(t, projectDir, "ddx", "obsidian", "revert")

    // Verify files restored to original state
    currentFiles := captureFileStates(t, projectDir)
    assert.Equal(t, originalFiles, currentFiles)
}
```

#### Expected Behavior:
- **MUST** work end-to-end in realistic project scenarios
- **MUST** handle projects with existing content gracefully
- **MUST** maintain file integrity throughout migration
- **MUST** provide reliable revert functionality
- **MUST** be compatible with Obsidian software

## Performance Test Specifications

### 7. Scalability Tests

```go
func TestObsidianMigration_LargeProject(t *testing.T) {
    // Create project with 1000+ markdown files
    projectDir := createLargeTestProject(t, 1000)
    defer os.RemoveAll(projectDir)

    start := time.Now()
    runCommand(t, projectDir, "ddx", "obsidian", "migrate")
    duration := time.Since(start)

    // Should complete in reasonable time (< 30 seconds for 1000 files)
    assert.Less(t, duration, 30*time.Second, "Migration took too long: %v", duration)

    // Verify all files processed correctly
    count := countProcessedFiles(t, projectDir)
    assert.GreaterOrEqual(t, count, 1000)
}

func TestObsidianValidation_Performance(t *testing.T) {
    projectDir := createLargeValidProject(t, 500)
    defer os.RemoveAll(projectDir)

    start := time.Now()
    runCommand(t, projectDir, "ddx", "obsidian", "validate")
    duration := time.Since(start)

    // Validation should be fast (< 10 seconds for 500 files)
    assert.Less(t, duration, 10*time.Second)
}
```

#### Expected Behavior:
- **MUST** handle large projects (1000+ files) efficiently
- **MUST** provide progress indicators for long operations
- **MUST** use parallel processing where possible
- **MUST** maintain reasonable memory usage
- **MUST** validate large projects quickly

## Error Handling Test Specifications

### 8. Error Condition Tests

```go
func TestObsidianMigration_ErrorHandling(t *testing.T) {
    tests := []struct {
        name          string
        setup         func(string) // Setup function to create error condition
        expectedError string
    }{
        {
            name: "invalid_yaml_frontmatter",
            setup: func(dir string) {
                createFileWithInvalidYAML(dir, "test.md")
            },
            expectedError: "invalid YAML frontmatter",
        },
        {
            name: "readonly_files",
            setup: func(dir string) {
                createReadOnlyFile(dir, "readonly.md")
            },
            expectedError: "permission denied",
        },
        {
            name: "missing_directory",
            setup: func(dir string) {
                // Don't create the directory
            },
            expectedError: "directory not found",
        },
        {
            name: "corrupted_markdown",
            setup: func(dir string) {
                createCorruptedMarkdownFile(dir, "corrupt.md")
            },
            expectedError: "failed to parse markdown",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testDir := t.TempDir()
            tt.setup(testDir)

            cmd := exec.Command("ddx", "obsidian", "migrate", "--path", testDir)
            output, err := cmd.CombinedOutput()

            assert.Error(t, err)
            assert.Contains(t, string(output), tt.expectedError)
        })
    }
}
```

#### Expected Behavior:
- **MUST** handle invalid YAML gracefully
- **MUST** handle file permission errors
- **MUST** handle missing directories/files
- **MUST** provide clear error messages
- **MUST** exit with appropriate error codes

## Test Data and Fixtures

### Required Test Fixtures

```bash
# Test data directory structure
test/fixtures/helix-workflow/
├── workflows/helix/
│   ├── coordinator.md
│   ├── principles.md
│   └── phases/
│       ├── 01-frame/
│       │   ├── README.md
│       │   ├── enforcer.md
│       │   └── artifacts/
│       │       ├── feature-specification/
│       │       │   ├── template.md
│       │       │   ├── prompt.md
│       │       │   └── example.md
│       │       └── user-stories/
│       └── 02-design/
├── docs/
│   ├── 01-frame/features/
│   │   ├── FEAT-001-authentication.md
│   │   └── FEAT-002-data-export.md
│   └── 02-design/
│       ├── FEAT-001-technical-design.md
│       └── FEAT-002-technical-design.md
└── .ddx.yml
```

## Success Criteria

All tests in this specification must:

1. **Pass consistently** - No flaky tests
2. **Run quickly** - Unit tests < 1s, integration tests < 30s
3. **Be maintainable** - Clear test names and documentation
4. **Cover edge cases** - Handle error conditions gracefully
5. **Follow TDD** - Tests written before implementation
6. **Use proper setup/teardown** - Clean test environment

## Test Execution Strategy

### Phase 1: Unit Tests (Red Phase)
- Write and run all unit tests
- Verify they all fail (no implementation exists)
- Commit failing tests to establish test baseline

### Phase 2: Integration Tests (Red Phase)
- Write CLI and navigation tests
- Verify they fail appropriately
- Document expected behavior clearly

### Phase 3: Implementation (Green Phase)
- Implement minimal code to make tests pass
- Follow test-driven development strictly
- One test at a time, simplest implementation first

### Phase 4: Refactoring (Refactor Phase)
- Improve code quality while keeping tests green
- Add performance optimizations
- Enhance error handling

This test specification ensures that the Obsidian integration will be thoroughly tested before any implementation begins, following proper HELIX Test phase governance.