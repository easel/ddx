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
		// Phase files
		{"workflows/helix/phases/01-frame/README.md", FileTypePhase},
		{"workflows/helix/phases/02-design/README.md", FileTypePhase},
		{"workflows/helix/phases/test/README.md", FileTypePhase},

		// Enforcer files
		{"workflows/helix/phases/01-frame/enforcer.md", FileTypeEnforcer},
		{"workflows/helix/phases/02-design/enforcer.md", FileTypeEnforcer},

		// Artifact files
		{"workflows/helix/phases/01-frame/artifacts/feature-specification/template.md", FileTypeTemplate},
		{"workflows/helix/phases/01-frame/artifacts/feature-specification/prompt.md", FileTypePrompt},
		{"workflows/helix/phases/01-frame/artifacts/feature-specification/example.md", FileTypeExample},

		// Core workflow files
		{"workflows/helix/coordinator.md", FileTypeCoordinator},
		{"workflows/helix/principles.md", FileTypePrinciple},

		// Documentation files (FEAT files in phase directories are still feature specifications)
		{"docs/01-frame/features/FEAT-001-new-feature.md", FileTypeFeature},
		{"docs/02-design/FEAT-001-technical-design.md", FileTypeFeature},
		{"docs/03-test/FEAT-001-test-specification.md", FileTypeFeature},
		{"docs/04-build/FEAT-001-implementation.md", FileTypeFeature},

		// Unknown files
		{"workflows/helix/random.md", FileTypeUnknown},
		{"README.md", FileTypeUnknown},
		{"docs/other/file.md", FileTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := detector.Detect(tt.path)
			if result != tt.expected {
				t.Errorf("Detect(%s) = %s, expected %s", tt.path, result, tt.expected)
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
		{"workflows/helix/phases/04-build/template.md", "build"},
		{"docs/01-frame/features/FEAT-001.md", "frame"},
		{"docs/02-design/technical-design.md", "design"},
		{"docs/build/implementation.md", "build"},
		{"workflows/helix/README.md", ""},
		{"README.md", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetPhaseFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("GetPhaseFromPath(%s) = %s, expected %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetArtifactCategory(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"workflows/helix/phases/01-frame/artifacts/feature-specification/template.md", "feature-specification"},
		{"workflows/helix/phases/01-frame/artifacts/user-stories/prompt.md", "user-stories"},
		{"docs/01-frame/features/FEAT-001-new-feature.md", "feature-specification"},
		{"docs/02-design/FEAT-001-technical-design.md", "feature-specification"},
		{"docs/03-test/test-specification.md", ""},
		{"docs/04-build/implementation-guide.md", ""},
		{"workflows/helix/README.md", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetArtifactCategory(tt.path)
			if result != tt.expected {
				t.Errorf("GetArtifactCategory(%s) = %s, expected %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetPhaseNumber(t *testing.T) {
	tests := []struct {
		phase    string
		expected int
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
		t.Run(tt.phase, func(t *testing.T) {
			result := GetPhaseNumber(tt.phase)
			if result != tt.expected {
				t.Errorf("GetPhaseNumber(%s) = %d, expected %d", tt.phase, result, tt.expected)
			}
		})
	}
}

func TestGetNextPhase(t *testing.T) {
	tests := []struct {
		phase    string
		expected string
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
		t.Run(tt.phase, func(t *testing.T) {
			result := GetNextPhase(tt.phase)
			if result != tt.expected {
				t.Errorf("GetNextPhase(%s) = %s, expected %s", tt.phase, result, tt.expected)
			}
		})
	}
}

func TestGetPreviousPhase(t *testing.T) {
	tests := []struct {
		phase    string
		expected string
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
		t.Run(tt.phase, func(t *testing.T) {
			result := GetPreviousPhase(tt.phase)
			if result != tt.expected {
				t.Errorf("GetPreviousPhase(%s) = %s, expected %s", tt.phase, result, tt.expected)
			}
		})
	}
}

func TestExtractTitleFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"workflows/helix/phases/01-frame/README.md", "Frame Phase"},
		{"workflows/helix/phases/01-frame/enforcer.md", "Frame Phase Enforcer"},
		{"workflows/helix/phases/01-frame/artifacts/feature-specification/template.md", "Feature Specification Template"},
		{"workflows/helix/phases/01-frame/artifacts/user-stories/prompt.md", "User Stories Prompt"},
		{"workflows/helix/phases/01-frame/artifacts/risk-register/example.md", "Risk Register Example"},
		{"docs/01-frame/features/FEAT-001-user-authentication.md", "User Authentication"},
		{"docs/02-design/api-architecture.md", "Api Architecture"},
		{"workflows/helix/coordinator.md", "Coordinator"},
		{"workflows/helix/principles.md", "Principles"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := ExtractTitleFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("ExtractTitleFromPath(%s) = %s, expected %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsHelixFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"workflows/helix/phases/01-frame/README.md", true},
		{"workflows/helix/coordinator.md", true},
		{"docs/01-frame/features/FEAT-001.md", true},
		{"docs/02-design/technical-design.md", true},
		{"docs/03-test/test-spec.md", true},
		{"docs/04-build/implementation.md", true},
		{"docs/05-deploy/deployment.md", true},
		{"docs/06-iterate/retrospective.md", true},
		{"README.md", false},
		{"src/main.go", false},
		{"docs/general/notes.md", false},
		{"other/workflows/helix/file.md", true}, // Contains helix path
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := IsHelixFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsHelixFile(%s) = %t, expected %t", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetComplexityFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"workflows/helix/phases/01-frame/artifacts/simple-feature/template.md", "moderate"},
		{"workflows/helix/phases/01-frame/artifacts/example/template.md", "simple"},
		{"workflows/helix/phases/01-frame/artifacts/advanced-integration/template.md", "complex"},
		{"workflows/helix/phases/01-frame/artifacts/complex-system/template.md", "complex"},
		{"workflows/helix/phases/01-frame/artifacts/user-stories/template.md", "moderate"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetComplexityFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("GetComplexityFromPath(%s) = %s, expected %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestFileTypeGetHierarchicalTags(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected []string
	}{
		{FileTypePhase, []string{"helix", "helix/phase"}},
		{FileTypeEnforcer, []string{"helix", "helix/enforcer"}},
		{FileTypeTemplate, []string{"helix", "helix/artifact", "helix/artifact/template"}},
		{FileTypePrompt, []string{"helix", "helix/artifact", "helix/artifact/prompt"}},
		{FileTypeExample, []string{"helix", "helix/artifact", "helix/artifact/example"}},
		{FileTypeCoordinator, []string{"helix", "helix/core", "helix/coordinator"}},
		{FileTypePrinciple, []string{"helix", "helix/core", "helix/principle"}},
		{FileTypeFeature, []string{"helix", "helix/artifact", "helix/artifact/specification"}},
		{FileTypeUnknown, []string{"helix"}},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			result := tt.fileType.GetHierarchicalTags()
			if len(result) != len(tt.expected) {
				t.Errorf("GetHierarchicalTags(%s) returned %d tags, expected %d", tt.fileType, len(result), len(tt.expected))
				return
			}

			for i, tag := range result {
				if tag != tt.expected[i] {
					t.Errorf("GetHierarchicalTags(%s)[%d] = %s, expected %s", tt.fileType, i, tag, tt.expected[i])
				}
			}
		})
	}
}

func TestFileTypeCheckers(t *testing.T) {
	tests := []struct {
		fileType   FileType
		isPhase    bool
		isArtifact bool
		isCore     bool
	}{
		{FileTypePhase, true, false, false},
		{FileTypeEnforcer, true, false, false},
		{FileTypeTemplate, false, true, false},
		{FileTypePrompt, false, true, false},
		{FileTypeExample, false, true, false},
		{FileTypeFeature, false, true, false},
		{FileTypeCoordinator, false, false, true},
		{FileTypePrinciple, false, false, true},
		{FileTypeUnknown, false, false, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			if result := tt.fileType.IsPhase(); result != tt.isPhase {
				t.Errorf("IsPhase(%s) = %t, expected %t", tt.fileType, result, tt.isPhase)
			}
			if result := tt.fileType.IsArtifact(); result != tt.isArtifact {
				t.Errorf("IsArtifact(%s) = %t, expected %t", tt.fileType, result, tt.isArtifact)
			}
			if result := tt.fileType.IsCore(); result != tt.isCore {
				t.Errorf("IsCore(%s) = %t, expected %t", tt.fileType, result, tt.isCore)
			}
		})
	}
}
