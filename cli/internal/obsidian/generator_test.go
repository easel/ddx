package obsidian

import (
	"strings"
	"testing"
)

func TestFrontmatterGenerator(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		name     string
		file     *MarkdownFile
		validate func(t *testing.T, fm *Frontmatter)
	}{
		{
			name: "Phase README generates correct frontmatter",
			file: &MarkdownFile{
				Path:     "workflows/helix/phases/01-frame/README.md",
				Content:  "# Frame Phase\n\nThis is the frame phase description.",
				FileType: FileTypePhase,
			},
			validate: func(t *testing.T, fm *Frontmatter) {
				if fm.Title != "Frame Phase" {
					t.Errorf("Expected title 'Frame Phase', got '%s'", fm.Title)
				}
				if fm.Type != "phase" {
					t.Errorf("Expected type 'phase', got '%s'", fm.Type)
				}
				if fm.PhaseID != "frame" {
					t.Errorf("Expected phase_id 'frame', got '%s'", fm.PhaseID)
				}
				if fm.PhaseNum != 1 {
					t.Errorf("Expected phase_number 1, got %d", fm.PhaseNum)
				}
				if !containsTest(fm.Tags, "helix") {
					t.Errorf("Expected 'helix' tag, got %v", fm.Tags)
				}
				if !containsTest(fm.Tags, "helix/phase") {
					t.Errorf("Expected 'helix/phase' tag, got %v", fm.Tags)
				}
				if !containsTest(fm.Tags, "helix/phase/frame") {
					t.Errorf("Expected 'helix/phase/frame' tag, got %v", fm.Tags)
				}
			},
		},
		{
			name: "Template generates artifact metadata",
			file: &MarkdownFile{
				Path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md",
				Content:  "# Feature Specification Template\n\nUse this template...",
				FileType: FileTypeTemplate,
			},
			validate: func(t *testing.T, fm *Frontmatter) {
				if fm.Title != "Feature Specification Template" {
					t.Errorf("Expected title 'Feature Specification Template', got '%s'", fm.Title)
				}
				if fm.Type != "template" {
					t.Errorf("Expected type 'template', got '%s'", fm.Type)
				}
				if fm.Phase != "frame" {
					t.Errorf("Expected phase 'frame', got '%s'", fm.Phase)
				}
				if fm.ArtifactCategory != "feature-specification" {
					t.Errorf("Expected artifact_category 'feature-specification', got '%s'", fm.ArtifactCategory)
				}
				if fm.Complexity == "" {
					t.Errorf("Expected complexity to be set")
				}
				if fm.TimeEstimate == "" {
					t.Errorf("Expected time_estimate to be set")
				}
			},
		},
		{
			name: "Feature specification generates feature metadata",
			file: &MarkdownFile{
				Path:     "docs/01-frame/features/FEAT-001-user-authentication.md",
				Content:  "# Feature Specification: FEAT-001 - User Authentication\n\n**Priority**: P1\n**Owner**: Engineering Team\n**Status**: draft",
				FileType: FileTypeFeature,
			},
			validate: func(t *testing.T, fm *Frontmatter) {
				if fm.Title != "User Authentication" {
					t.Errorf("Expected title 'User Authentication', got '%s'", fm.Title)
				}
				if fm.Type != "feature-specification" {
					t.Errorf("Expected type 'feature-specification', got '%s'", fm.Type)
				}
				if fm.FeatureID != "FEAT-001" {
					t.Errorf("Expected feature_id 'FEAT-001', got '%s'", fm.FeatureID)
				}
				if fm.Priority != "P1" {
					t.Errorf("Expected priority 'P1', got '%s'", fm.Priority)
				}
				if fm.Owner != "Engineering Team" {
					t.Errorf("Expected owner 'Engineering Team', got '%s'", fm.Owner)
				}
				if fm.Status != "draft" {
					t.Errorf("Expected status 'draft', got '%s'", fm.Status)
				}
				if fm.WorkflowPhase != "frame" {
					t.Errorf("Expected workflow_phase 'frame', got '%s'", fm.WorkflowPhase)
				}
				if fm.ArtifactType != "feature-specification" {
					t.Errorf("Expected artifact_type 'feature-specification', got '%s'", fm.ArtifactType)
				}
			},
		},
		{
			name: "Enforcer generates enforcer metadata",
			file: &MarkdownFile{
				Path:     "workflows/helix/phases/01-frame/enforcer.md",
				Content:  "# Frame Phase Enforcer\n\nYou are the Frame Phase Guardian...",
				FileType: FileTypeEnforcer,
			},
			validate: func(t *testing.T, fm *Frontmatter) {
				if fm.Title != "Frame Phase Enforcer" {
					t.Errorf("Expected title 'Frame Phase Enforcer', got '%s'", fm.Title)
				}
				if fm.Type != "enforcer" {
					t.Errorf("Expected type 'enforcer', got '%s'", fm.Type)
				}
				if fm.Phase != "frame" {
					t.Errorf("Expected phase 'frame', got '%s'", fm.Phase)
				}
				if !containsTest(fm.Tags, "helix/phase/frame/enforcer") {
					t.Errorf("Expected 'helix/phase/frame/enforcer' tag, got %v", fm.Tags)
				}
				if len(fm.Aliases) == 0 {
					t.Errorf("Expected aliases to be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, err := generator.Generate(tt.file)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			// Check basic fields
			if fm.Created.IsZero() {
				t.Errorf("Expected created timestamp to be set")
			}
			if fm.Updated.IsZero() {
				t.Errorf("Expected updated timestamp to be set")
			}
			if len(fm.Tags) == 0 {
				t.Errorf("Expected tags to be set")
			}

			// Run custom validation
			tt.validate(t, fm)
		})
	}
}

func TestExtractTitleFromContent(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		content  string
		expected string
	}{
		{"# Simple Title\n\nContent here.", "Simple Title"},
		{"# Feature Specification: FEAT-001 - User Auth\n\nContent.", "User Auth"},
		{"# Technical Design: FEAT-002 - API Gateway\n\nContent.", "API Gateway"},
		{"# Build Implementation: FEAT-003 - Payment System\n\nContent.", "Payment System"},
		{"# [[FEAT-004]] - Database Migration\n\nContent.", "Database Migration"},
		{"No heading content.", ""},
		{"## Secondary Heading\n\nNo primary heading.", ""},
		{"# \n\nEmpty heading.", ""},
	}

	for _, tt := range tests {
		t.Run(tt.content[:minTest(30, len(tt.content))], func(t *testing.T) {
			result := generator.extractTitleFromContent(tt.content)
			if result != tt.expected {
				t.Errorf("extractTitleFromContent() = '%s', expected '%s'", result, tt.expected)
			}
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
				Path:    "docs/01-frame/features/FEAT-123-new-feature.md",
				Content: "Some content",
			},
			expected: "FEAT-123",
		},
		{
			name: "Extract from content",
			file: &MarkdownFile{
				Path:    "docs/01-frame/features/new-feature.md",
				Content: "# Feature Specification: FEAT-456 - New Feature\n\nContent here.",
			},
			expected: "FEAT-456",
		},
		{
			name: "No feature ID",
			file: &MarkdownFile{
				Path:    "docs/01-frame/features/no-id.md",
				Content: "# Some Feature\n\nNo ID here.",
			},
			expected: "",
		},
		{
			name: "Prefer filename over content",
			file: &MarkdownFile{
				Path:    "docs/01-frame/features/FEAT-789-primary.md",
				Content: "# Feature Specification: FEAT-999 - Secondary\n\nContent.",
			},
			expected: "FEAT-789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.extractFeatureID(tt.file)
			if result != tt.expected {
				t.Errorf("extractFeatureID() = '%s', expected '%s'", result, tt.expected)
			}
		})
	}
}

func TestExtractFromContent(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		content  string
		field    string
		expected string
	}{
		{"**Priority**: P1\n", "Priority", "P1"},
		{"Priority: P2\n", "Priority", "P2"},
		{"**Owner**: Engineering Team\n", "Owner", "Engineering Team"},
		{"Owner: John Doe\n", "Owner", "John Doe"},
		{"**Status**: [Draft]\n", "Status", "Draft"},
		{"Status: in_progress\n", "Status", "in_progress"},
		{"No priority field", "Priority", ""},
		{"**Priority**\n", "Priority", ""},
		{"Random text", "Owner", ""},
	}

	for _, tt := range tests {
		t.Run(tt.field+"_"+tt.content[:min(20, len(tt.content))], func(t *testing.T) {
			result := generator.extractFromContent(tt.content, tt.field)
			if result != tt.expected {
				t.Errorf("extractFromContent(%s, %s) = '%s', expected '%s'", tt.content, tt.field, result, tt.expected)
			}
		})
	}
}

func TestGenerateTags(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		name     string
		file     *MarkdownFile
		validate func(t *testing.T, tags []string)
	}{
		{
			name: "Phase file tags",
			file: &MarkdownFile{
				Path:     "workflows/helix/phases/01-frame/README.md",
				FileType: FileTypePhase,
			},
			validate: func(t *testing.T, tags []string) {
				expectedTags := []string{"helix", "helix/phase", "helix/phase/frame"}
				for _, expected := range expectedTags {
					if !containsTest(tags, expected) {
						t.Errorf("Expected tag '%s' not found in %v", expected, tags)
					}
				}
			},
		},
		{
			name: "Template file tags",
			file: &MarkdownFile{
				Path:     "workflows/helix/phases/01-frame/artifacts/user-stories/template.md",
				FileType: FileTypeTemplate,
			},
			validate: func(t *testing.T, tags []string) {
				// Note: some tags might not be present due to category parsing
				requiredTags := []string{"helix", "helix/artifact", "helix/artifact/template", "helix/phase/frame"}
				for _, expected := range requiredTags {
					if !containsTest(tags, expected) {
						t.Errorf("Expected tag '%s' not found in %v", expected, tags)
					}
				}
			},
		},
		{
			name: "Feature file tags",
			file: &MarkdownFile{
				Path:     "docs/01-frame/features/FEAT-001-auth.md",
				FileType: FileTypeFeature,
			},
			validate: func(t *testing.T, tags []string) {
				expectedTags := []string{"helix", "helix/artifact", "helix/artifact/specification", "helix/phase/frame"}
				for _, expected := range expectedTags {
					if !containsTest(tags, expected) {
						t.Errorf("Expected tag '%s' not found in %v", expected, tags)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tags := generator.generateTags(tt.file)
			if len(tags) == 0 {
				t.Errorf("Expected tags to be generated")
			}
			tt.validate(t, tags)
		})
	}
}

func TestAddPhaseMetadata(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/phases/02-design/README.md",
		FileType: FileTypePhase,
	}

	fm := &Frontmatter{}
	generator.addPhaseMetadata(fm, file)

	if fm.PhaseID != "design" {
		t.Errorf("Expected phase_id 'design', got '%s'", fm.PhaseID)
	}
	if fm.PhaseNum != 2 {
		t.Errorf("Expected phase_number 2, got %d", fm.PhaseNum)
	}
	if fm.NextPhase != "[[Test Phase]]" {
		t.Errorf("Expected next_phase '[[Test Phase]]', got '%s'", fm.NextPhase)
	}
	if fm.PrevPhase != "[[Frame Phase]]" {
		t.Errorf("Expected previous_phase '[[Frame Phase]]', got '%s'", fm.PrevPhase)
	}
	if fm.Gates == nil {
		t.Errorf("Expected gates to be set")
	}
	if fm.Artifacts == nil {
		t.Errorf("Expected artifacts to be set")
	}
	if len(fm.Aliases) == 0 {
		t.Errorf("Expected aliases to be set")
	}
}

func TestAddArtifactMetadata(t *testing.T) {
	generator := NewFrontmatterGenerator()

	file := &MarkdownFile{
		Path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md",
		FileType: FileTypeTemplate,
		Content:  strings.Repeat("a", 2000), // Moderate length content
	}

	fm := &Frontmatter{}
	generator.addArtifactMetadata(fm, file)

	if fm.Phase != "frame" {
		t.Errorf("Expected phase 'frame', got '%s'", fm.Phase)
	}
	if fm.ArtifactCategory != "feature-specification" {
		t.Errorf("Expected artifact_category 'feature-specification', got '%s'", fm.ArtifactCategory)
	}
	if fm.Complexity != "moderate" {
		t.Errorf("Expected complexity 'moderate', got '%s'", fm.Complexity)
	}
	if fm.TimeEstimate == "" {
		t.Errorf("Expected time_estimate to be set")
	}
	if fm.Prerequisites == nil {
		t.Errorf("Expected prerequisites to be initialized")
	}
	if fm.Outputs == nil {
		t.Errorf("Expected outputs to be initialized")
	}
}

func TestTimeEstimateByType(t *testing.T) {
	generator := NewFrontmatterGenerator()

	tests := []struct {
		fileType FileType
		expected string
	}{
		{FileTypeTemplate, "30-60 minutes"},
		{FileTypePrompt, "15-30 minutes"},
		{FileTypeExample, "5-15 minutes"},
		{FileTypeFeature, "1-2 hours"},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			file := &MarkdownFile{
				Path:     "test/path.md",
				FileType: tt.fileType,
			}

			fm := &Frontmatter{}
			generator.addArtifactMetadata(fm, file)

			if fm.TimeEstimate != tt.expected {
				t.Errorf("Expected time_estimate '%s' for %s, got '%s'", tt.expected, tt.fileType, fm.TimeEstimate)
			}
		})
	}
}

// Helper function for tests
func containsTest(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func minTest(a, b int) int {
	if a < b {
		return a
	}
	return b
}
