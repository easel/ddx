package obsidian

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileType_String(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected string
	}{
		{FileTypePhase, "phase"},
		{FileTypeArtifact, "artifact"},
		{FileTypeTemplate, "template"},
		{FileTypePrompt, "prompt"},
		{FileTypeExample, "example"},
		{FileTypeCoordinator, "coordinator"},
		{FileTypePrinciple, "principle"},
		{FileTypeFeature, "feature-specification"},
		{FileTypeEnforcer, "enforcer"},
		{FileTypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			result := tt.fileType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileType_IsPhase(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected bool
	}{
		{FileTypePhase, true},
		{FileTypeEnforcer, true},
		{FileTypeArtifact, false},
		{FileTypeTemplate, false},
		{FileTypeCoordinator, false},
		{FileTypeUnknown, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			result := tt.fileType.IsPhase()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileType_IsArtifact(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected bool
	}{
		{FileTypeArtifact, true},
		{FileTypeTemplate, true},
		{FileTypePrompt, true},
		{FileTypeExample, true},
		{FileTypeFeature, true},
		{FileTypePhase, false},
		{FileTypeCoordinator, false},
		{FileTypeUnknown, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			result := tt.fileType.IsArtifact()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFileType_GetHierarchicalTags(t *testing.T) {
	tests := []struct {
		fileType FileType
		expected []string
	}{
		{
			FileTypePhase,
			[]string{"helix", "helix/phase"},
		},
		{
			FileTypeTemplate,
			[]string{"helix", "helix/artifact", "helix/artifact/template"},
		},
		{
			FileTypeFeature,
			[]string{"helix", "helix/artifact", "helix/artifact/specification"},
		},
		{
			FileTypeCoordinator,
			[]string{"helix", "helix/core", "helix/coordinator"},
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.fileType), func(t *testing.T) {
			result := tt.fileType.GetHierarchicalTags()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarkdownFile_HasFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		file     *MarkdownFile
		expected bool
	}{
		{
			name: "with frontmatter",
			file: &MarkdownFile{
				Frontmatter: &Frontmatter{Title: "Test"},
			},
			expected: true,
		},
		{
			name:     "without frontmatter",
			file:     &MarkdownFile{Frontmatter: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.file.HasFrontmatter()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarkdownFile_GetTitle(t *testing.T) {
	tests := []struct {
		name     string
		file     *MarkdownFile
		expected string
	}{
		{
			name: "from frontmatter",
			file: &MarkdownFile{
				Frontmatter: &Frontmatter{Title: "Test Title"},
			},
			expected: "Test Title",
		},
		{
			name: "no frontmatter",
			file: &MarkdownFile{
				Frontmatter: nil,
			},
			expected: "Untitled",
		},
		{
			name: "empty title in frontmatter",
			file: &MarkdownFile{
				Frontmatter: &Frontmatter{Title: ""},
			},
			expected: "Untitled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.file.GetTitle()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMarkdownFile_GetTags(t *testing.T) {
	tests := []struct {
		name     string
		file     *MarkdownFile
		expected []string
	}{
		{
			name: "with tags",
			file: &MarkdownFile{
				Frontmatter: &Frontmatter{
					Tags: []string{"helix", "phase", "frame"},
				},
			},
			expected: []string{"helix", "phase", "frame"},
		},
		{
			name: "no frontmatter",
			file: &MarkdownFile{
				Frontmatter: nil,
			},
			expected: []string{},
		},
		{
			name: "empty tags",
			file: &MarkdownFile{
				Frontmatter: &Frontmatter{Tags: []string{}},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.file.GetTags()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewTagTree(t *testing.T) {
	tree := NewTagTree()

	assert.NotNil(t, tree)
	assert.NotNil(t, tree.Root)
	assert.Equal(t, "helix", tree.Root.Name)
	assert.Equal(t, "helix", tree.Root.FullPath)
	assert.NotNil(t, tree.Root.Children)
	assert.NotNil(t, tree.AllTags)
	assert.NotNil(t, tree.TagCount)
}

func TestTagTree_Add(t *testing.T) {
	tree := NewTagTree()
	file := &MarkdownFile{Path: "test.md"}

	// Add hierarchical tag
	tree.Add("helix/phase/frame", file)

	// Verify tag count
	assert.Equal(t, 1, tree.TagCount["helix/phase/frame"])

	// Verify tree structure
	phaseNode := tree.Get("helix/phase")
	assert.NotNil(t, phaseNode)
	assert.Equal(t, "phase", phaseNode.Name)
	assert.Equal(t, "helix/phase", phaseNode.FullPath)

	frameNode := tree.Get("helix/phase/frame")
	assert.NotNil(t, frameNode)
	assert.Equal(t, "frame", frameNode.Name)
	assert.Equal(t, "helix/phase/frame", frameNode.FullPath)
	assert.Contains(t, frameNode.Files, file)

	// Verify parent-child relationship
	assert.Equal(t, phaseNode, frameNode.Parent)
	assert.Contains(t, phaseNode.Children, "frame")
}

func TestTagTree_GetFilesByTag(t *testing.T) {
	tree := NewTagTree()
	file1 := &MarkdownFile{Path: "test1.md"}
	file2 := &MarkdownFile{Path: "test2.md"}

	tree.Add("helix/phase/frame", file1)
	tree.Add("helix/phase/frame", file2)
	tree.Add("helix/phase/design", file2)

	// Get files by specific tag
	frameFiles := tree.GetFilesByTag("helix/phase/frame")
	assert.Len(t, frameFiles, 2)
	assert.Contains(t, frameFiles, file1)
	assert.Contains(t, frameFiles, file2)

	designFiles := tree.GetFilesByTag("helix/phase/design")
	assert.Len(t, designFiles, 1)
	assert.Contains(t, designFiles, file2)

	// Non-existent tag
	nonExistent := tree.GetFilesByTag("non/existent")
	assert.Nil(t, nonExistent)
}

func TestFrontmatter_Timestamps(t *testing.T) {
	now := time.Now()
	fm := &Frontmatter{
		Created: now,
		Updated: now.Add(time.Hour),
	}

	assert.Equal(t, now, fm.Created)
	assert.True(t, fm.Updated.After(fm.Created))
}
