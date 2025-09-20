---
title: "Test Procedures - Obsidian Integration for HELIX"
type: test-procedures
feature_id: FEAT-014
workflow_phase: test
artifact_type: test-procedures
tags:
  - helix/test
  - helix/artifact/test
  - helix/phase/test
  - obsidian
  - procedures
related:
  - "[[FEAT-014-test-specification]]"
  - "[[FEAT-014-obsidian-integration]]"
  - "[[FEAT-014-technical-design]]"
status: draft
priority: P1
created: 2025-01-18
updated: 2025-01-18
---

# Test Procedures: [[FEAT-014]] Obsidian Integration

## Test Execution Strategy

This document defines the specific procedures for executing the comprehensive test suite for the Obsidian integration feature. These procedures must be followed to ensure proper Test-Driven Development and HELIX compliance.

## Pre-Implementation Test Setup

### Step 1: Create Test Infrastructure

```bash
# Create test directory structure
mkdir -p cli/internal/obsidian
mkdir -p cli/cmd
mkdir -p cli/test/fixtures/helix-workflow

# Create basic test files structure
touch cli/internal/obsidian/detector_test.go
touch cli/internal/obsidian/generator_test.go
touch cli/internal/obsidian/converter/links_test.go
touch cli/internal/obsidian/navigation_test.go
touch cli/cmd/obsidian_test.go
touch cli/test/obsidian_e2e_test.go
```

### Step 2: Set Up Test Fixtures

Create comprehensive test data representing a realistic HELIX workflow:

```bash
# Setup test HELIX workflow structure
mkdir -p cli/test/fixtures/helix-workflow/workflows/helix/phases/{01-frame,02-design,03-test,04-build,05-deploy,06-iterate}
mkdir -p cli/test/fixtures/helix-workflow/docs/{01-frame/features,02-design,03-test,04-build,05-deploy,06-iterate}

# Create sample files (these will be used to test the migration)
cat > cli/test/fixtures/helix-workflow/workflows/helix/coordinator.md << 'EOF'
# HELIX Workflow Coordinator

The HELIX workflow provides...

See the [Frame Phase](./phases/01-frame/README.md) to get started.
EOF

cat > cli/test/fixtures/helix-workflow/workflows/helix/phases/01-frame/README.md << 'EOF'
# Frame Phase - Problem Definition

This is the foundation phase where we define **WHAT** we're building.

Move to [Design Phase](../02-design/README.md) after completing requirements.
EOF

cat > cli/test/fixtures/helix-workflow/docs/helix/01-frame/features/FEAT-001-authentication.md << 'EOF'
# Feature Specification: FEAT-001 - User Authentication

**Priority**: P1
**Owner**: Security Team

## Overview
Implement secure user authentication system...

Related: [technical design](../../02-design/FEAT-001-technical-design.md)
EOF
```

## Test Phase Execution (Red Phase)

### Step 3: Write Failing Unit Tests

#### 3.1 File Type Detection Tests

```bash
# File: cli/internal/obsidian/detector_test.go
cat > cli/internal/obsidian/detector_test.go << 'EOF'
package obsidian

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFileTypeDetector_DetectHelixFiles(t *testing.T) {
    tests := []struct {
        name     string
        path     string
        expected FileType
    }{
        {
            name:     "Frame phase README",
            path:     "workflows/helix/phases/01-frame/README.md",
            expected: FileTypePhase,
        },
        {
            name:     "Feature specification",
            path:     "docs/helix/01-frame/features/FEAT-014-obsidian-integration.md",
            expected: FileTypeFeature,
        },
        {
            name:     "HELIX coordinator",
            path:     "workflows/helix/coordinator.md",
            expected: FileTypeCoordinator,
        },
    }

    detector := NewFileTypeDetector() // This will fail - doesn't exist yet
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := detector.Detect(tt.path)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestGetPhaseFromPath(t *testing.T) {
    tests := []struct {
        path     string
        expected string
    }{
        {"workflows/helix/phases/01-frame/README.md", "frame"},
        {"docs/helix/02-design/FEAT-001-technical-design.md", "design"},
        {"random/path/file.md", ""},
    }

    for _, tt := range tests {
        result := GetPhaseFromPath(tt.path) // This will fail - doesn't exist yet
        assert.Equal(t, tt.expected, result)
    }
}
EOF
```

#### 3.2 Frontmatter Generation Tests

```bash
# File: cli/internal/obsidian/generator_test.go
cat > cli/internal/obsidian/generator_test.go << 'EOF'
package obsidian

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFrontmatterGenerator_PhaseFiles(t *testing.T) {
    generator := NewFrontmatterGenerator() // This will fail - doesn't exist yet

    file := &MarkdownFile{
        Path:     "workflows/helix/phases/01-frame/README.md",
        Content:  "# Frame Phase - Problem Definition\n\nThis phase...",
        FileType: FileTypePhase,
    }

    fm, err := generator.Generate(file) // This will fail - doesn't exist yet
    require.NoError(t, err)

    assert.Equal(t, "Frame Phase - Problem Definition", fm.Title)
    assert.Equal(t, "phase", fm.Type)
    assert.Equal(t, "frame", fm.PhaseID)
    assert.Equal(t, 1, fm.PhaseNum)
    assert.Contains(t, fm.Tags, "helix/phase/frame")
}

func TestFrontmatterGenerator_FeatureFiles(t *testing.T) {
    generator := NewFrontmatterGenerator()

    file := &MarkdownFile{
        Path:     "docs/helix/01-frame/features/FEAT-014-obsidian-integration.md",
        Content:  "# Feature Specification: FEAT-014\n\n**Priority**: P1\n**Owner**: Platform Team",
        FileType: FileTypeFeature,
    }

    fm, err := generator.Generate(file)
    require.NoError(t, err)

    assert.Equal(t, "FEAT-014", fm.FeatureID)
    assert.Equal(t, "P1", fm.Priority)
    assert.Equal(t, "Platform Team", fm.Owner)
}
EOF
```

#### 3.3 Link Conversion Tests

```bash
# File: cli/internal/obsidian/converter/links_test.go
mkdir -p cli/internal/obsidian/converter
cat > cli/internal/obsidian/converter/links_test.go << 'EOF'
package converter

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/easel/ddx/internal/obsidian"
)

func TestLinkConverter_MarkdownToWikilink(t *testing.T) {
    converter := NewLinkConverter() // This will fail - doesn't exist yet

    files := []*obsidian.MarkdownFile{
        {Path: "workflows/helix/phases/02-design/README.md", Frontmatter: &obsidian.Frontmatter{Title: "Design Phase"}},
    }
    converter.BuildIndex(files) // This will fail - doesn't exist yet

    tests := []struct {
        input    string
        expected string
    }{
        {
            input:    "See the [Design Phase](../02-design/README.md) for details",
            expected: "See the [[Design Phase]] for details",
        },
        {
            input:    "Visit [GitHub](https://github.com) for more info",
            expected: "Visit [GitHub](https://github.com) for more info", // External preserved
        },
    }

    for _, tt := range tests {
        result := converter.ConvertContent(tt.input) // This will fail - doesn't exist yet
        assert.Equal(t, tt.expected, result)
    }
}
EOF
```

### Step 4: Run Failing Tests

```bash
# Verify all tests fail as expected (Red phase)
cd cli
go test ./internal/obsidian/... -v
# Expected output: All tests should fail with compilation errors

# This confirms we're properly following TDD - tests are written first
```

### Step 5: CLI Command Tests

```bash
# File: cli/cmd/obsidian_test.go
cat > cli/cmd/obsidian_test.go << 'EOF'
package cmd

import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestObsidianMigrate_DryRun(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping CLI test in short mode")
    }

    // Setup test directory
    testDir := t.TempDir()
    setupTestHelixWorkflow(t, testDir)

    // Run dry-run migration - This will fail, no command exists yet
    cmd := exec.Command("ddx", "obsidian", "migrate", "--dry-run", "--path", testDir)
    output, err := cmd.CombinedOutput()

    require.NoError(t, err)
    assert.Contains(t, string(output), "Would add frontmatter to")
    assert.Contains(t, string(output), "Would convert links in")
}

func setupTestHelixWorkflow(t *testing.T, dir string) {
    // Create basic HELIX structure for testing
    workflowDir := filepath.Join(dir, "workflows/helix")
    err := os.MkdirAll(workflowDir, 0755)
    require.NoError(t, err)

    // Create sample coordinator file
    coordinatorPath := filepath.Join(workflowDir, "coordinator.md")
    content := "# HELIX Coordinator\n\nSee [Frame Phase](./phases/01-frame/README.md)"
    err = os.WriteFile(coordinatorPath, []byte(content), 0644)
    require.NoError(t, err)
}
EOF
```

## Implementation Phase (Green Phase)

### Step 6: Minimal Implementation Strategy

Once all tests are written and failing, begin implementation in this order:

#### 6.1 Core Types Definition

```bash
# File: cli/internal/obsidian/types.go
# Implement basic types needed by tests
# - FileType constants
# - MarkdownFile struct
# - Frontmatter struct
```

#### 6.2 File Type Detection

```bash
# File: cli/internal/obsidian/detector.go
# Implement:
# - NewFileTypeDetector()
# - FileTypeDetector.Detect()
# - GetPhaseFromPath()
# - GetArtifactCategory()
```

#### 6.3 Frontmatter Generation

```bash
# File: cli/internal/obsidian/generator.go
# Implement:
# - NewFrontmatterGenerator()
# - FrontmatterGenerator.Generate()
# - Title extraction logic
# - Tag generation logic
```

#### 6.4 Link Conversion

```bash
# File: cli/internal/obsidian/converter/links.go
# Implement:
# - NewLinkConverter()
# - LinkConverter.BuildIndex()
# - LinkConverter.ConvertContent()
```

#### 6.5 CLI Commands

```bash
# File: cli/cmd/obsidian.go
# Implement:
# - obsidian migrate command
# - obsidian validate command
# - obsidian revert command
```

### Step 7: Test-Driven Implementation Cycle

For each component:

1. **Run failing tests** to see what's needed
2. **Write minimal code** to make tests pass
3. **Run tests again** to verify they pass
4. **Refactor** if needed while keeping tests green
5. **Move to next component**

```bash
# Example cycle for detector:
go test ./internal/obsidian -run TestFileTypeDetector -v
# See what's missing, implement minimal detector
go test ./internal/obsidian -run TestFileTypeDetector -v
# Verify tests pass, refactor if needed
```

## Integration Testing Procedures

### Step 8: CLI Integration Tests

```bash
# Build the CLI first
go build -o ddx

# Test basic functionality
./ddx obsidian migrate --help
# Should show usage information

# Test with real test data
./ddx obsidian migrate --dry-run --path test/fixtures/helix-workflow
# Should show what would be migrated

# Test actual migration
cp -r test/fixtures/helix-workflow /tmp/test-migration
./ddx obsidian migrate --path /tmp/test-migration
# Should successfully migrate files

# Validate migration
./ddx obsidian validate --path /tmp/test-migration
# Should report validation success
```

### Step 9: End-to-End Testing

```bash
# Create a real HELIX project
mkdir /tmp/real-helix-test
cd /tmp/real-helix-test
../ddx init --workflow helix

# Add some content
echo "# Test Feature" > docs/helix/01-frame/features/FEAT-TEST.md
echo "See [Frame Phase](../phases/01-frame/README.md)" >> docs/helix/01-frame/features/FEAT-TEST.md

# Migrate to Obsidian
../ddx obsidian migrate

# Verify results
grep "^---" docs/helix/01-frame/features/FEAT-TEST.md
# Should show frontmatter was added

grep "\[\[.*\]\]" docs/helix/01-frame/features/FEAT-TEST.md
# Should show wikilinks were created
```

## Validation Procedures

### Step 10: Comprehensive Validation

```bash
# Run all tests
go test ./... -v

# Test with various project structures
for project in small-project medium-project large-project; do
    echo "Testing with $project"
    ./ddx obsidian migrate --path test/fixtures/$project
    ./ddx obsidian validate --path test/fixtures/$project
done

# Performance testing
time ./ddx obsidian migrate --path test/fixtures/large-project
# Should complete in reasonable time

# Memory usage testing
/usr/bin/time -v ./ddx obsidian migrate --path test/fixtures/large-project
# Should use reasonable memory
```

### Step 11: Error Condition Testing

```bash
# Test invalid input
./ddx obsidian migrate --path /nonexistent/path
# Should fail gracefully with clear error

# Test permission issues
chmod 444 test/fixtures/readonly-project/test.md
./ddx obsidian migrate --path test/fixtures/readonly-project
# Should handle permission errors gracefully

# Test corrupted files
echo "invalid yaml frontmatter" > test/fixtures/corrupt-project/bad.md
./ddx obsidian validate --path test/fixtures/corrupt-project
# Should report validation errors clearly
```

## Quality Assurance Procedures

### Step 12: Code Quality Checks

```bash
# Linting
go vet ./...
golangci-lint run

# Test coverage
go test ./... -cover
# Should achieve >80% coverage

# Security scanning
gosec ./...
# Should pass security checks

# Dependency checking
go mod verify
go mod tidy
```

### Step 13: Documentation Validation

```bash
# Generate CLI documentation
./ddx obsidian --help
./ddx obsidian migrate --help
./ddx obsidian validate --help
./ddx obsidian revert --help

# Verify all commands documented
# Verify examples work as documented
# Verify error messages are helpful
```

## Acceptance Testing Procedures

### Step 14: User Acceptance Testing

Manual testing procedures to verify user experience:

1. **New User Experience**
   - Clone a fresh HELIX project
   - Run `ddx obsidian migrate`
   - Verify intuitive output and progress
   - Open in Obsidian, verify graph view works

2. **Existing User Experience**
   - Take project with existing modifications
   - Run migration
   - Verify no content lost
   - Verify customizations preserved

3. **Error Recovery**
   - Intentionally break something
   - Run validation
   - Verify clear error messages
   - Test revert functionality

### Step 15: Obsidian Compatibility Testing

1. **Import in Obsidian**
   - Open migrated project in Obsidian
   - Verify all wikilinks resolve
   - Verify graph view shows connections
   - Verify tags work for filtering

2. **Obsidian Features**
   - Test backlink functionality
   - Test search with tags
   - Test navigation between phases
   - Test plugin compatibility

## Success Criteria Verification

### Final Validation Checklist

- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] All CLI commands work as specified
- [ ] Migration preserves all content
- [ ] Validation catches all error conditions
- [ ] Revert functionality works completely
- [ ] Performance meets requirements
- [ ] Memory usage is reasonable
- [ ] Error messages are clear and helpful
- [ ] Documentation is complete and accurate
- [ ] Obsidian compatibility confirmed
- [ ] Code quality standards met
- [ ] Security requirements satisfied

## Test Phase Completion

The Test phase is complete when:

1. **All tests written and failing** (Red phase established)
2. **Test infrastructure ready** for implementation
3. **Test procedures documented** and validated
4. **Success criteria defined** and measurable
5. **Implementation plan clear** from test requirements

At this point, we can transition to the Build phase with confidence that we have comprehensive test coverage defining exactly what needs to be implemented.

## Next Steps

After Test phase completion:
1. Execute Build phase following test-driven development
2. Implement minimal code to make each test pass
3. Refactor while keeping tests green
4. Validate all acceptance criteria met
5. Transition to Deploy phase for release preparation

This Test phase establishes a solid foundation for reliable, well-tested implementation of the Obsidian integration feature.