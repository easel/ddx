package obsidian

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrontmatterGenerator_PhaseFiles(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/phases/01-frame/README.md",
		Content:  "# Frame Phase - Problem Definition\n\nThis phase defines what we're building.",
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
	assert.Empty(t, fm.PrevPhase) // First phase has no previous
	assert.NotNil(t, fm.Gates)
	assert.NotNil(t, fm.Artifacts)

	// Timestamps should be set
	assert.False(t, fm.Created.IsZero())
	assert.False(t, fm.Updated.IsZero())
	assert.True(t, fm.Updated.Sub(fm.Created) >= 0)
}

func TestFrontmatterGenerator_TemplateFiles(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md",
		Content:  "# Feature Specification Template\n\nUse this template to create feature specifications.",
		FileType: FileTypeTemplate,
	}

	fm, err := generator.Generate(file)
	require.NoError(t, err)

	assert.Equal(t, "Feature Specification Template", fm.Title)
	assert.Equal(t, "template", fm.Type)
	assert.Equal(t, "frame", fm.Phase)
	assert.Equal(t, "feature-specification", fm.ArtifactCategory)
	assert.Contains(t, fm.Tags, "helix/artifact")
	assert.Contains(t, fm.Tags, "helix/artifact/feature/specification")
	assert.Equal(t, "30-60 minutes", fm.TimeEstimate)
	assert.Equal(t, "moderate", fm.Complexity)
}

func TestFrontmatterGenerator_PromptFiles(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/phases/01-frame/artifacts/user-stories/prompt.md",
		Content:  "# User Stories Prompt\n\nUse this prompt to generate user stories.",
		FileType: FileTypePrompt,
	}

	fm, err := generator.Generate(file)
	require.NoError(t, err)

	assert.Equal(t, "User Stories Prompt", fm.Title)
	assert.Equal(t, "prompt", fm.Type)
	assert.Equal(t, "frame", fm.Phase)
	assert.Equal(t, "user-stories", fm.ArtifactCategory)
	assert.Equal(t, "15-30 minutes", fm.TimeEstimate)
}

func TestFrontmatterGenerator_ExampleFiles(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/example.md",
		Content:  "# Feature Specification Example\n\nThis is an example of a well-written feature specification.",
		FileType: FileTypeExample,
	}

	fm, err := generator.Generate(file)
	require.NoError(t, err)

	assert.Equal(t, "Feature Specification Example", fm.Title)
	assert.Equal(t, "example", fm.Type)
	assert.Equal(t, "5-15 minutes", fm.TimeEstimate)
}

func TestFrontmatterGenerator_FeatureFiles(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "docs/01-frame/features/FEAT-014-obsidian-integration.md",
		Content:  "# Feature Specification: FEAT-014 - Obsidian Integration\n\n**Priority**: P1\n**Owner**: Platform Team\n**Status**: Draft",
		FileType: FileTypeFeature,
	}

	fm, err := generator.Generate(file)
	require.NoError(t, err)

	assert.Equal(t, "FEAT-014", fm.FeatureID)
	assert.Equal(t, "Obsidian Integration", fm.Title)
	assert.Equal(t, "P1", fm.Priority)
	assert.Equal(t, "Platform Team", fm.Owner)
	assert.Equal(t, "Draft", fm.Status)
	assert.Equal(t, "frame", fm.WorkflowPhase)
	assert.Equal(t, "feature-specification", fm.ArtifactType)
}

func TestFrontmatterGenerator_EnforcerFiles(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/phases/02-design/enforcer.md",
		Content:  "# Design Phase Enforcer\n\nThis enforcer ensures proper design practices.",
		FileType: FileTypeEnforcer,
	}

	fm, err := generator.Generate(file)
	require.NoError(t, err)

	assert.Equal(t, "Design Phase Enforcer", fm.Title)
	assert.Equal(t, "enforcer", fm.Type)
	assert.Equal(t, "design", fm.Phase)
	assert.Contains(t, fm.Tags, "helix/phase/design/enforcer")
	assert.Contains(t, fm.Aliases, "Design Phase Enforcer")
	assert.Contains(t, fm.Aliases, "Design Guardian")
}

func TestFrontmatterGenerator_CoordinatorFiles(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/coordinator.md",
		Content:  "# HELIX Workflow Coordinator\n\nThe central coordination point for the HELIX workflow.",
		FileType: FileTypeCoordinator,
	}

	fm, err := generator.Generate(file)
	require.NoError(t, err)

	assert.Equal(t, "HELIX Workflow Coordinator", fm.Title)
	assert.Equal(t, "coordinator", fm.Type)
	assert.Contains(t, fm.Aliases, "HELIX Coordinator")
	assert.Contains(t, fm.Aliases, "Workflow Coordinator")
}

func TestExtractTitleFromContent(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Simple H1 title",
			content:  "# Simple Title\n\nContent here.",
			expected: "Simple Title",
		},
		{
			name:     "Feature specification title",
			content:  "# Feature Specification: FEAT-001 - User Authentication\n\nContent.",
			expected: "User Authentication",
		},
		{
			name:     "Technical design title",
			content:  "# Technical Design: FEAT-002 - Data Export\n\nDesign details.",
			expected: "Data Export",
		},
		{
			name:     "Build implementation title",
			content:  "# Build Implementation: FEAT-003 - Payment System\n\nImplementation.",
			expected: "Payment System",
		},
		{
			name:     "Wikilink in title",
			content:  "# [[FEAT-004]] - Database Migration\n\nMigration details.",
			expected: "Database Migration",
		},
		{
			name:     "No heading content",
			content:  "Just some content without a heading.",
			expected: "",
		},
		{
			name:     "Secondary heading only",
			content:  "## Secondary Heading\n\nNo primary heading here.",
			expected: "",
		},
		{
			name:     "Empty heading",
			content:  "# \n\nEmpty heading.",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.extractTitleFromContent(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractFeatureID(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		name     string
		file     *MarkdownFile
		expected string
	}{
		{
			name: "Extract from filename",
			file: &MarkdownFile{
				Path:    "docs/01-frame/features/FEAT-123-user-auth.md",
				Content: "# User Authentication",
			},
			expected: "FEAT-123",
		},
		{
			name: "Extract from content",
			file: &MarkdownFile{
				Path:    "docs/features/user-auth.md",
				Content: "# Feature Specification: FEAT-456\n\nContent here.",
			},
			expected: "FEAT-456",
		},
		{
			name: "No feature ID",
			file: &MarkdownFile{
				Path:    "docs/general/overview.md",
				Content: "# General Overview\n\nNo feature ID here.",
			},
			expected: "",
		},
		{
			name: "Prefer filename over content",
			file: &MarkdownFile{
				Path:    "docs/features/FEAT-111-auth.md",
				Content: "# Feature Specification: FEAT-222\n\nContent.",
			},
			expected: "FEAT-111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.extractFeatureID(tt.file)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractFromContent(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		name     string
		content  string
		field    string
		expected string
	}{
		{
			name:     "Priority with bold markdown",
			content:  "**Priority**: P1\n\nOther content.",
			field:    "Priority",
			expected: "P1",
		},
		{
			name:     "Priority without bold",
			content:  "Priority: P2\n\nOther content.",
			field:    "Priority",
			expected: "P2",
		},
		{
			name:     "Owner with bold markdown",
			content:  "**Owner**: Engineering Team\n\nOther content.",
			field:    "Owner",
			expected: "Engineering Team",
		},
		{
			name:     "Owner without bold",
			content:  "Owner: John Doe\n\nOther content.",
			field:    "Owner",
			expected: "John Doe",
		},
		{
			name:     "Status with brackets",
			content:  "**Status**: [Draft]\n\nOther content.",
			field:    "Status",
			expected: "Draft",
		},
		{
			name:     "Status without brackets",
			content:  "Status: in progress\n\nOther content.",
			field:    "Status",
			expected: "in progress",
		},
		{
			name:     "No field found",
			content:  "Random text without the field we're looking for.",
			field:    "Priority",
			expected: "",
		},
		{
			name:     "Field with just bold marker",
			content:  "**Priority**\n\nNo value provided.",
			field:    "Priority",
			expected: "",
		},
		{
			name:     "Field in different case",
			content:  "**priority**: P3\n\nContent.",
			field:    "Priority",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.extractFromContent(tt.content, tt.field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateTags(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		name     string
		file     *MarkdownFile
		expected []string
	}{
		{
			name: "Phase file tags",
			file: &MarkdownFile{
				Path:     "workflows/helix/phases/01-frame/README.md",
				FileType: FileTypePhase,
			},
			expected: []string{"helix", "helix/phase", "helix/phase/frame"},
		},
		{
			name: "Template file tags",
			file: &MarkdownFile{
				Path:     "workflows/helix/phases/01-frame/artifacts/user-stories/template.md",
				FileType: FileTypeTemplate,
			},
			expected: []string{"helix", "helix/artifact", "helix/artifact/template", "helix/phase/frame", "helix/artifact/user/stories"},
		},
		{
			name: "Feature file tags",
			file: &MarkdownFile{
				Path:     "docs/01-frame/features/FEAT-001.md",
				FileType: FileTypeFeature,
			},
			expected: []string{"helix", "helix/artifact", "helix/artifact/specification", "helix/phase/frame", "helix/artifact/feature/specification"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.generateTags(tt.file)

			// Check that all expected tags are present
			for _, expectedTag := range tt.expected {
				assert.Contains(t, result, expectedTag, "Missing tag: %s", expectedTag)
			}

			// Check that helix is always the first tag
			assert.Equal(t, "helix", result[0])
		})
	}
}

func TestTimeEstimateByType(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected string
	}{
		{FileTypeTemplate, "30-60 minutes"},
		{FileTypePrompt, "15-30 minutes"},
		{FileTypeExample, "5-15 minutes"},
		{FileTypeFeature, "1-2 hours"},
		{FileTypeArtifact, "1-2 hours"},
		{FileTypePhase, "1-2 hours"},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			generator := NewFrontmatterGenerator()
			file := &MarkdownFile{FileType: tt.fileType}

			fm, err := generator.Generate(file)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, fm.TimeEstimate)
		})
	}
}

func TestAddPhaseMetadata(t *testing.T) {
	generator := NewFrontmatterGenerator()
	fm := &Frontmatter{Tags: []string{"helix"}}
	file := &MarkdownFile{
		Path:     "workflows/helix/phases/02-design/README.md",
		FileType: FileTypePhase,
	}

	generator.addPhaseMetadata(fm, file)

	assert.Equal(t, "design", fm.PhaseID)
	assert.Equal(t, 2, fm.PhaseNum)
	assert.Equal(t, "[[Test Phase]]", fm.NextPhase)
	assert.Equal(t, "[[Frame Phase]]", fm.PrevPhase)
	assert.NotNil(t, fm.Gates)
	assert.NotNil(t, fm.Artifacts)
	assert.Contains(t, fm.Aliases, "Design Phase")
}

func TestAddArtifactMetadata(t *testing.T) {
	generator := NewFrontmatterGenerator()
	fm := &Frontmatter{Tags: []string{"helix"}}
	file := &MarkdownFile{
		Path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md",
		FileType: FileTypeTemplate,
	}

	generator.addArtifactMetadata(fm, file)

	assert.Equal(t, "frame", fm.Phase)
	assert.Equal(t, "feature-specification", fm.ArtifactCategory)
	assert.Equal(t, "moderate", fm.Complexity)
	assert.NotNil(t, fm.Prerequisites)
	assert.NotNil(t, fm.Outputs)
}
