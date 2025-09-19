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
		// Phase detection
		{
			name:     "Frame phase README",
			path:     "workflows/helix/phases/01-frame/README.md",
			expected: FileTypePhase,
		},
		{
			name:     "Design phase README",
			path:     "workflows/helix/phases/02-design/README.md",
			expected: FileTypePhase,
		},
		{
			name:     "Test phase README",
			path:     "workflows/helix/phases/test/README.md",
			expected: FileTypePhase,
		},

		// Enforcer detection
		{
			name:     "Frame phase enforcer",
			path:     "workflows/helix/phases/01-frame/enforcer.md",
			expected: FileTypeEnforcer,
		},
		{
			name:     "Design phase enforcer",
			path:     "workflows/helix/phases/02-design/enforcer.md",
			expected: FileTypeEnforcer,
		},

		// Artifact templates
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
		{
			name:     "Technical design example",
			path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/example.md",
			expected: FileTypeExample,
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
			path:     "docs/01-frame/features/FEAT-001-new-feature.md",
			expected: FileTypeFeature,
		},
		{
			name:     "Design document",
			path:     "docs/02-design/FEAT-001-technical-design.md",
			expected: FileTypeFeature,
		},
		{
			name:     "Test specification",
			path:     "docs/03-test/FEAT-001-test-specification.md",
			expected: FileTypeFeature,
		},
		{
			name:     "Implementation guide",
			path:     "docs/04-build/FEAT-001-implementation.md",
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
			path:     "docs/other/file.md",
			expected: FileTypeUnknown,
		},
		{
			name:     "Random file in workflows",
			path:     "workflows/helix/random.md",
			expected: FileTypeUnknown,
		},
	}

	detector := NewFileTypeDetector()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detector.Detect(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s", tt.path)
		})
	}
}

func TestGetPhaseFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		// Workflow structure
		{"workflows/helix/phases/01-frame/README.md", "frame"},
		{"workflows/helix/phases/02-design/enforcer.md", "design"},
		{"workflows/helix/phases/test/artifacts/test.md", "test"},
		{"workflows/helix/phases/04-build/template.md", "build"},

		// Docs structure
		{"docs/01-frame/features/FEAT-001.md", "frame"},
		{"docs/02-design/technical-design.md", "design"},
		{"docs/build/implementation.md", "build"},

		// Non-phase paths
		{"workflows/helix/README.md", ""},
		{"README.md", ""},
		{"random/path/file.md", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetPhaseFromPath(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s", tt.path)
		})
	}
}

func TestGetArtifactCategory(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		// Workflow artifacts
		{"workflows/helix/phases/01-frame/artifacts/feature-specification/template.md", "feature-specification"},
		{"workflows/helix/phases/01-frame/artifacts/user-stories/prompt.md", "user-stories"},

		// Feature files
		{"docs/01-frame/features/FEAT-001-new-feature.md", "feature-specification"},
		{"docs/02-design/FEAT-001-technical-design.md", "feature-specification"},
		{"docs/03-test/test-specification.md", ""},
		{"docs/04-build/implementation-guide.md", ""},

		// Non-artifact paths
		{"workflows/helix/README.md", ""},
		{"docs/other/file.md", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetArtifactCategory(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s", tt.path)
		})
	}
}

func TestGetPhaseNumber(t *testing.T) {
	tests := []struct {
		phaseName string
		expected  int
	}{
		{"frame", 1},
		{"design", 2},
		{"test", 3},
		{"build", 4},
		{"deploy", 5},
		{"iterate", 6},
		{"unknown", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.phaseName, func(t *testing.T) {
			result := GetPhaseNumber(tt.phaseName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetNextPhase(t *testing.T) {
	tests := []struct {
		phaseName string
		expected  string
	}{
		{"frame", "design"},
		{"design", "test"},
		{"test", "build"},
		{"build", "deploy"},
		{"deploy", "iterate"},
		{"iterate", ""},
		{"unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.phaseName, func(t *testing.T) {
			result := GetNextPhase(tt.phaseName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetPreviousPhase(t *testing.T) {
	tests := []struct {
		phaseName string
		expected  string
	}{
		{"frame", ""},
		{"design", "frame"},
		{"test", "design"},
		{"build", "test"},
		{"deploy", "build"},
		{"iterate", "deploy"},
		{"unknown", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.phaseName, func(t *testing.T) {
			result := GetPreviousPhase(tt.phaseName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractTitleFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		// README files
		{"workflows/helix/phases/01-frame/README.md", "Frame Phase"},
		{"workflows/helix/phases/02-design/README.md", "Design Phase"},
		{"workflows/helix/artifacts/feature-specification/README.md", "Feature Specification"},

		// Special files
		{"workflows/helix/phases/01-frame/enforcer.md", "Frame Phase Enforcer"},
		{"workflows/helix/phases/02-design/enforcer.md", "Design Phase Enforcer"},
		{"workflows/helix/artifacts/feature-specification/template.md", "Feature Specification Template"},
		{"workflows/helix/artifacts/user-stories/prompt.md", "User Stories Prompt"},
		{"workflows/helix/artifacts/technical-design/example.md", "Technical Design Example"},

		// Feature files
		{"docs/01-frame/features/FEAT-001-user-authentication.md", "User Authentication"},
		{"docs/02-design/FEAT-002-data-export-system.md", "Data Export System"},

		// Regular files
		{"some-random-file.md", "Some Random File"},
		{"docs/user-guide.md", "User Guide"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ExtractTitleFromPath(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s", tt.path)
		})
	}
}

// Removed IsHelixFile test as we're making DDX workflow-agnostic
// Workflow-specific file detection should be handled through configuration

func TestGetComplexityFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		// Simple complexity
		{"workflows/helix/artifacts/example/template.md", "simple"},
		{"docs/simple/feature.md", "simple"},

		// Complex complexity
		{"workflows/helix/artifacts/advanced/template.md", "complex"},
		{"docs/complex/implementation.md", "complex"},

		// Default to moderate
		{"workflows/helix/artifacts/feature-specification/template.md", "moderate"},
		{"docs/01-frame/features/FEAT-001.md", "moderate"},
		{"regular/file.md", "moderate"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetComplexityFromPath(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s", tt.path)
		})
	}
}
