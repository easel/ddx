package converter

import (
	"testing"

	"github.com/easel/ddx/internal/obsidian"
	"github.com/stretchr/testify/assert"
)

func TestLinkConverter_MarkdownToWikilink(t *testing.T) {
	converter := NewLinkConverter()

	// Build test file index
	files := []*obsidian.MarkdownFile{
		{
			Path: "workflows/helix/phases/02-design/README.md",
			Frontmatter: &obsidian.Frontmatter{
				Title: "Design Phase",
			},
		},
		{
			Path: "workflows/helix/phases/01-frame/artifacts/feature-specification/template.md",
			Frontmatter: &obsidian.Frontmatter{
				Title: "Feature Specification Template",
			},
		},
		{
			Path: "workflows/helix/phases/01-frame/artifacts/user-stories/template.md",
			Frontmatter: &obsidian.Frontmatter{
				Title: "User Stories Template",
			},
		},
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
			name:     "Template link matching title",
			input:    "Use the [Feature Specification Template](./artifacts/feature-specification/template.md)",
			expected: "Use the [[Feature Specification Template]]",
		},
		{
			name:     "External link preserved",
			input:    "Visit [GitHub](https://github.com) for more info",
			expected: "Visit [GitHub](https://github.com) for more info",
		},
		{
			name:     "HTTPS link preserved",
			input:    "Check [documentation](https://docs.example.com/guide)",
			expected: "Check [documentation](https://docs.example.com/guide)",
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
		{
			name:     "Multiple links in same text",
			input:    "See [Design Phase](../02-design/README.md) and [User Stories Template](./artifacts/user-stories/template.md)",
			expected: "See [[Design Phase]] and [[User Stories Template]]",
		},
		{
			name:     "Mixed link types",
			input:    "Check [Design Phase](../02-design/README.md) and [GitHub](https://github.com)",
			expected: "Check [[Design Phase]] and [GitHub](https://github.com)",
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
		name     string
		input    string
		expected string
	}{
		{
			name:     "Frame phase reference",
			input:    "Complete the Frame phase first",
			expected: "Complete the [[Frame Phase|Frame phase]] first",
		},
		{
			name:     "Design phase planning",
			input:    "During Design phase planning",
			expected: "During [[Design Phase|Design phase]] planning",
		},
		{
			name:     "Build phase implementation",
			input:    "The Build phase implementation",
			expected: "The [[Build Phase|Build phase]] implementation",
		},
		{
			name:     "Multiple phase references",
			input:    "After Frame phase, move to Design phase",
			expected: "After [[Frame Phase|Frame phase]], move to [[Design Phase|Design phase]]",
		},
		{
			name:     "Already wikilinked phases",
			input:    "Already in [[Frame Phase]] work",
			expected: "Already in [[Frame Phase]] work", // Don't double-convert
		},
		{
			name:     "Mixed existing and new",
			input:    "From [[Frame Phase]] to Design phase",
			expected: "From [[Frame Phase]] to [[Design Phase|Design phase]]",
		},
		{
			name:     "Capitalized phase names",
			input:    "Complete the Frame Phase first",
			expected: "Complete the [[Frame Phase]] first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertContent(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLinkConverter_ArtifactReferences(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Feature specification reference",
			input:    "Create a feature specification for this work",
			expected: "Create a [[Feature Specification|feature specification]] for this work",
		},
		{
			name:     "Feature spec shorthand",
			input:    "Update the feature spec with new requirements",
			expected: "Update the [[Feature Specification|feature spec]] with new requirements",
		},
		{
			name:     "Technical design reference",
			input:    "The technical design should include...",
			expected: "The [[Technical Design|technical design]] should include...",
		},
		{
			name:     "User stories reference",
			input:    "Write user stories before implementation",
			expected: "Write [[User Stories|user stories]] before implementation",
		},
		{
			name:     "PRD reference",
			input:    "The PRD contains all requirements",
			expected: "The [[Product Requirements Document|PRD]] contains all requirements",
		},
		{
			name:     "Already wikilinked artifacts",
			input:    "Update the [[Feature Specification]] document",
			expected: "Update the [[Feature Specification]] document", // Don't double-convert
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertContent(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLinkConverter_WorkflowReferences(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HELIX workflow reference",
			input:    "Follow the HELIX workflow for best results",
			expected: "Follow the [[HELIX Workflow]] for best results",
		},
		{
			name:     "TDD reference",
			input:    "Use TDD for development",
			expected: "Use [[Test-Driven Development|TDD]] for development",
		},
		{
			name:     "Test-Driven Development reference",
			input:    "Test-Driven Development improves code quality",
			expected: "[[Test-Driven Development]] improves code quality",
		},
		{
			name:     "Mixed case TDD",
			input:    "Apply test-driven development principles",
			expected: "Apply [[Test-Driven Development|test-driven development]] principles",
		},
		{
			name:     "Already wikilinked workflow",
			input:    "The [[HELIX Workflow]] provides structure",
			expected: "The [[HELIX Workflow]] provides structure", // Don't double-convert
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertContent(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLinkConverter_PreventDoubleConversion(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Existing wikilinks preserved",
			input:    "See [[Design Phase]] and visit [[Frame Phase]] for details",
			expected: "See [[Design Phase]] and visit [[Frame Phase]] for details",
		},
		{
			name:     "Wikilinks with aliases preserved",
			input:    "Check [[Feature Specification|feature spec]] and [[Technical Design|tech design]]",
			expected: "Check [[Feature Specification|feature spec]] and [[Technical Design|tech design]]",
		},
		{
			name:     "Mixed existing and new conversions",
			input:    "From [[Frame Phase]] to Design phase via feature specification",
			expected: "From [[Frame Phase]] to [[Design Phase|Design phase]] via [[Feature Specification|feature specification]]",
		},
		{
			name:     "Embedded wikilinks",
			input:    "![[Workflow Diagram]] shows the process",
			expected: "![[Workflow Diagram]] shows the process", // Embeds preserved
		},
		{
			name:     "Wikilinks with headings",
			input:    "See [[Design Phase#Architecture]] section",
			expected: "See [[Design Phase#Architecture]] section",
		},
		{
			name:     "Wikilinks with block references",
			input:    "Reference [[Requirements^block-123]] for details",
			expected: "Reference [[Requirements^block-123]] for details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertContent(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLinkConverter_CommonPatterns(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "README to phase conversion",
			input:    "Check [README](../phases/01-frame/README.md)",
			expected: "Check [[Frame Phase|README]]",
		},
		{
			name:     "Template file conversion",
			input:    "Use [template](./artifacts/feature-specification/template.md)",
			expected: "Use [[Feature Specification Template|template]]",
		},
		{
			name:     "Prompt file conversion",
			input:    "Run [prompt](./artifacts/user-stories/prompt.md)",
			expected: "Run [[User Stories Prompt|prompt]]",
		},
		{
			name:     "Example file conversion",
			input:    "See [example](./artifacts/feature-specification/example.md)",
			expected: "See [[Feature Specification Example|example]]",
		},
		{
			name:     "Enforcer file conversion",
			input:    "Check [enforcer](./phases/02-design/enforcer.md)",
			expected: "Check [[Design Phase Enforcer|enforcer]]",
		},
		{
			name:     "Coordinator reference",
			input:    "Read [coordinator](./coordinator.md) for overview",
			expected: "Read [[HELIX Workflow Coordinator|coordinator]] for overview",
		},
		{
			name:     "Principles reference",
			input:    "Follow [principles](./principles.md) guidelines",
			expected: "Follow [[HELIX Principles|principles]] guidelines",
		},
		{
			name:     "Feature file conversion",
			input:    "Update [FEAT-001](../features/FEAT-001-auth.md)",
			expected: "Update [[FEAT-001]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertContent(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseWikilinks(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []*obsidian.ParsedLink
	}{
		{
			name:    "Simple wikilink",
			content: "See [[Design Phase]] for details",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Design Phase]]",
					Target:   "Design Phase",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Wikilink with alias",
			content: "Check [[Feature Specification|feature spec]] document",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Feature Specification|feature spec]]",
					Target:   "Feature Specification",
					Alias:    "feature spec",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Wikilink with heading",
			content: "Read [[Design Phase#Architecture]] section",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Design Phase#Architecture]]",
					Target:   "Design Phase",
					Heading:  "Architecture",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Wikilink with block reference",
			content: "See [[Requirements^block-123]] for context",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Requirements^block-123]]",
					Target:   "Requirements",
					BlockID:  "block-123",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Embedded wikilink",
			content: "![[Workflow Diagram]] shows the process",
			expected: []*obsidian.ParsedLink{
				{
					Original: "![[Workflow Diagram]]",
					Target:   "Workflow Diagram",
					IsEmbed:  true,
				},
			},
		},
		{
			name:    "Complex wikilink with alias and heading",
			content: "Check [[Technical Design#API Contracts|API docs]] section",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Technical Design#API Contracts|API docs]]",
					Target:   "Technical Design",
					Alias:    "API docs",
					Heading:  "API Contracts",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Multiple wikilinks",
			content: "From [[Frame Phase]] to [[Design Phase]] via [[Feature Specification]]",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Frame Phase]]",
					Target:   "Frame Phase",
					IsEmbed:  false,
				},
				{
					Original: "[[Design Phase]]",
					Target:   "Design Phase",
					IsEmbed:  false,
				},
				{
					Original: "[[Feature Specification]]",
					Target:   "Feature Specification",
					IsEmbed:  false,
				},
			},
		},
		{
			name:     "No wikilinks",
			content:  "Just regular text with [markdown](link.md) but no wikilinks",
			expected: []*obsidian.ParsedLink{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseWikilinks(tt.content)
			assert.Equal(t, len(tt.expected), len(result), "Wrong number of parsed links")

			for i, expected := range tt.expected {
				if i < len(result) {
					assert.Equal(t, expected.Original, result[i].Original)
					assert.Equal(t, expected.Target, result[i].Target)
					assert.Equal(t, expected.Alias, result[i].Alias)
					assert.Equal(t, expected.Heading, result[i].Heading)
					assert.Equal(t, expected.BlockID, result[i].BlockID)
					assert.Equal(t, expected.IsEmbed, result[i].IsEmbed)
				}
			}
		})
	}
}

func TestLinkConverter_ValidateWikilinks(t *testing.T) {
	converter := NewLinkConverter()

	// Build test file index
	files := []*obsidian.MarkdownFile{
		{
			Path: "phase1.md",
			Frontmatter: &obsidian.Frontmatter{
				Title:   "Frame Phase",
				Aliases: []string{"Phase 1", "Foundation Phase"},
			},
		},
		{
			Path: "phase2.md",
			Frontmatter: &obsidian.Frontmatter{
				Title: "Design Phase",
			},
		},
	}
	converter.BuildIndex(files)

	tests := []struct {
		name           string
		content        string
		expectedBroken []string
	}{
		{
			name:           "All valid links",
			content:        "See [[Frame Phase]] and [[Design Phase]]",
			expectedBroken: []string{},
		},
		{
			name:           "Valid alias link",
			content:        "Check [[Phase 1]] for foundation",
			expectedBroken: []string{},
		},
		{
			name:           "Broken link",
			content:        "See [[Nonexistent Phase]] for details",
			expectedBroken: []string{"Nonexistent Phase"},
		},
		{
			name:           "Mixed valid and broken",
			content:        "From [[Frame Phase]] to [[Missing Phase]] via [[Design Phase]]",
			expectedBroken: []string{"Missing Phase"},
		},
		{
			name:           "Multiple broken links",
			content:        "Check [[Missing One]] and [[Missing Two]]",
			expectedBroken: []string{"Missing One", "Missing Two"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ValidateWikilinks(tt.content)
			assert.ElementsMatch(t, tt.expectedBroken, result)
		})
	}
}
