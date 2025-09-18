package converter

import (
	"testing"

	"github.com/easel/ddx/internal/obsidian"
)

func TestLinkConverter(t *testing.T) {
	converter := NewLinkConverter()

	// Build test index
	files := []*obsidian.MarkdownFile{
		{
			Path: "workflows/helix/phases/01-frame/README.md",
			Frontmatter: &obsidian.Frontmatter{
				Title:   "Frame Phase",
				Aliases: []string{"Frame", "Frame Phase"},
			},
		},
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
			Path: "workflows/helix/coordinator.md",
			Frontmatter: &obsidian.Frontmatter{
				Title: "HELIX Workflow Coordinator",
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
			name:     "Convert relative markdown link to wikilink",
			input:    "See the [Design Phase](../02-design/README.md) for details.",
			expected: "See the [[Design Phase]] for details.",
		},
		{
			name:     "Convert template link",
			input:    "Use the [feature specification template](./artifacts/feature-specification/template.md).",
			expected: "Use the [[Feature Specification Template|feature specification template]].",
		},
		{
			name:     "Keep external links unchanged",
			input:    "Visit [GitHub](https://github.com) for more info.",
			expected: "Visit [GitHub](https://github.com) for more info.",
		},
		{
			name:     "Keep anchor links unchanged",
			input:    "See [section below](#implementation) for details.",
			expected: "See [section below](#implementation) for details.",
		},
		{
			name:     "Convert phase references",
			input:    "The Frame phase comes before the Design phase.",
			expected: "The [[Frame Phase]] comes before the [[Design Phase]].",
		},
		{
			name:     "Convert artifact references",
			input:    "Create a feature specification using the template.",
			expected: "Create a [[Feature Specification]] using the template.",
		},
		{
			name:     "Convert HELIX workflow references",
			input:    "The HELIX workflow uses TDD principles.",
			expected: "The [[HELIX Workflow]] uses [[Test-Driven Development|TDD]] principles.",
		},
		{
			name:     "Don't double-convert existing wikilinks",
			input:    "See [[Design Phase]] for more info.",
			expected: "See [[Design Phase]] for more info.",
		},
		{
			name:     "Convert multiple links in one text",
			input:    "The [Frame Phase](../01-frame/README.md) leads to [Design Phase](../02-design/README.md).",
			expected: "The [[Frame Phase]] leads to [[Design Phase]].",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertContent(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertContent() = '%s', expected '%s'", result, tt.expected)
			}
		})
	}
}

func TestConvertToWikilink(t *testing.T) {
	converter := NewLinkConverter()

	// Build test index
	files := []*obsidian.MarkdownFile{
		{
			Path: "workflows/helix/phases/01-frame/README.md",
			Frontmatter: &obsidian.Frontmatter{
				Title: "Frame Phase",
			},
		},
		{
			Path: "workflows/helix/phases/01-frame/enforcer.md",
			Frontmatter: &obsidian.Frontmatter{
				Title: "Frame Phase Enforcer",
			},
		},
	}

	converter.BuildIndex(files)

	tests := []struct {
		text     string
		path     string
		expected string
	}{
		{
			text:     "Frame Phase",
			path:     "../01-frame/README.md",
			expected: "[[Frame Phase]]",
		},
		{
			text:     "frame phase",
			path:     "../01-frame/README.md",
			expected: "[[Frame Phase|frame phase]]",
		},
		{
			text:     "enforcer",
			path:     "./enforcer.md",
			expected: "[[Frame Phase Enforcer|enforcer]]",
		},
		{
			text:     "unknown link",
			path:     "./unknown.md",
			expected: "[[unknown link]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.text+"_"+tt.path, func(t *testing.T) {
			result := converter.convertToWikilink(tt.text, tt.path)
			if result != tt.expected {
				t.Errorf("convertToWikilink(%s, %s) = '%s', expected '%s'", tt.text, tt.path, result, tt.expected)
			}
		})
	}
}

func TestHandleCommonPatterns(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		text     string
		path     string
		filename string
		expected string
	}{
		{
			text:     "Frame Phase",
			path:     "workflows/helix/phases/01-frame/README",
			filename: "README",
			expected: "[[Frame Phase]]",
		},
		{
			text:     "enforcer",
			path:     "workflows/helix/phases/01-frame/enforcer",
			filename: "enforcer",
			expected: "[[Frame Phase Enforcer]]",
		},
		{
			text:     "template",
			path:     "workflows/helix/phases/01-frame/artifacts/feature-specification/template",
			filename: "template",
			expected: "[[Feature Specification Template]]",
		},
		{
			text:     "prompt",
			path:     "workflows/helix/phases/01-frame/artifacts/user-stories/prompt",
			filename: "prompt",
			expected: "[[User Stories Prompt]]",
		},
		{
			text:     "example",
			path:     "workflows/helix/phases/01-frame/artifacts/risk-register/example",
			filename: "example",
			expected: "[[Risk Register Example]]",
		},
		{
			text:     "coordinator",
			path:     "workflows/helix/coordinator",
			filename: "coordinator",
			expected: "[[HELIX Workflow Coordinator]]",
		},
		{
			text:     "principles",
			path:     "workflows/helix/principles",
			filename: "principles",
			expected: "[[HELIX Principles]]",
		},
		{
			text:     "feature",
			path:     "docs/01-frame/features/FEAT-001-auth",
			filename: "FEAT-001-auth",
			expected: "[[FEAT-001]]",
		},
		{
			text:     "unknown",
			path:     "unknown/path",
			filename: "unknown",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := converter.handleCommonPatterns(tt.text, tt.path, tt.filename)
			if result != tt.expected {
				t.Errorf("handleCommonPatterns(%s, %s, %s) = '%s', expected '%s'", tt.text, tt.path, tt.filename, result, tt.expected)
			}
		})
	}
}

func TestConvertPhaseReferences(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "The Frame phase comes first.",
			expected: "The [[Frame Phase]] comes first.",
		},
		{
			input:    "Move to the Design Phase next.",
			expected: "Move to the [[Design Phase]] next.",
		},
		{
			input:    "The test phase should have failing tests.",
			expected: "The [[Test Phase]] should have failing tests.",
		},
		{
			input:    "During the build phase, implement the solution.",
			expected: "During the [[Build Phase]], implement the solution.",
		},
		{
			input:    "The Deploy phase handles releases.",
			expected: "The [[Deploy Phase]] handles releases.",
		},
		{
			input:    "In the iterate phase, we learn and improve.",
			expected: "In the [[Iterate Phase]], we learn and improve.",
		},
		{
			input:    "Frame has enforcer rules.",
			expected: "[[Frame Phase|Frame]] has enforcer rules.",
		},
		{
			input:    "Go to Design for architecture.",
			expected: "Go to [[Design Phase|Design]] for architecture.",
		},
		{
			input:    "The Build artifacts are ready.",
			expected: "The [[Build Phase|Build]] artifacts are ready.",
		},
		{
			input:    "Already in [[Frame Phase]] so no change.",
			expected: "Already in [[Frame Phase]] so no change.",
		},
		{
			input:    "Frame work is different from Frame phase.",
			expected: "Frame work is different from [[Frame Phase]].",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input[:minTest(30, len(tt.input))], func(t *testing.T) {
			result := converter.convertPhaseReferences(tt.input)
			if result != tt.expected {
				t.Errorf("convertPhaseReferences() = '%s', expected '%s'", result, tt.expected)
			}
		})
	}
}

func TestConvertArtifactReferences(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Create a feature specification for the new API.",
			expected: "Create a [[Feature Specification]] for the new API.",
		},
		{
			input:    "The feature spec should include all requirements.",
			expected: "The [[Feature Specification|feature spec]] should include all requirements.",
		},
		{
			input:    "Write the technical design document.",
			expected: "Write the [[Technical Design]] document.",
		},
		{
			input:    "The implementation guide provides step-by-step instructions.",
			expected: "The [[Implementation Guide]] provides step-by-step instructions.",
		},
		{
			input:    "User stories define the requirements.",
			expected: "[[User Stories]] define the requirements.",
		},
		{
			input:    "The PRD contains product requirements.",
			expected: "The [[Product Requirements Document|PRD]] contains product requirements.",
		},
		{
			input:    "Update the risk register with new risks.",
			expected: "Update the [[Risk Register]] with new risks.",
		},
		{
			input:    "Already have [[Feature Specification]] linked.",
			expected: "Already have [[Feature Specification]] linked.",
		},
		{
			input:    "The specification is feature-specific.",
			expected: "The specification is feature-specific.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input[:minTest(30, len(tt.input))], func(t *testing.T) {
			result := converter.convertArtifactReferences(tt.input)
			if result != tt.expected {
				t.Errorf("convertArtifactReferences() = '%s', expected '%s'", result, tt.expected)
			}
		})
	}
}

func TestConvertWorkflowReferences(t *testing.T) {
	converter := NewLinkConverter()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "The HELIX workflow is test-driven.",
			expected: "The [[HELIX Workflow]] is test-driven.",
		},
		{
			input:    "We use TDD in our development.",
			expected: "We use [[Test-Driven Development|TDD]] in our development.",
		},
		{
			input:    "Test-Driven Development improves quality.",
			expected: "[[Test-Driven Development]] improves quality.",
		},
		{
			input:    "The helix workflow helps teams.",
			expected: "The [[HELIX Workflow]] helps teams.",
		},
		{
			input:    "Already using [[HELIX Workflow]] here.",
			expected: "Already using [[HELIX Workflow]] here.",
		},
		{
			input:    "TDD and [[Test-Driven Development]] both.",
			expected: "[[Test-Driven Development|TDD]] and [[Test-Driven Development]] both.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input[:minTest(30, len(tt.input))], func(t *testing.T) {
			result := converter.convertWorkflowReferences(tt.input)
			if result != tt.expected {
				t.Errorf("convertWorkflowReferences() = '%s', expected '%s'", result, tt.expected)
			}
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
			content: "See [[Frame Phase]] for details.",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Frame Phase]]",
					Target:   "Frame Phase",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Wikilink with alias",
			content: "The [[Frame Phase|frame]] is first.",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Frame Phase|frame]]",
					Target:   "Frame Phase",
					Alias:    "frame",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Wikilink with heading",
			content: "See [[Frame Phase#Overview]] section.",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Frame Phase#Overview]]",
					Target:   "Frame Phase",
					Heading:  "Overview",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Wikilink with block ID",
			content: "Reference [[Frame Phase^block123]] here.",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Frame Phase^block123]]",
					Target:   "Frame Phase",
					BlockID:  "block123",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Embedded wikilink",
			content: "![[Diagram]] shows the flow.",
			expected: []*obsidian.ParsedLink{
				{
					Original: "![[Diagram]]",
					Target:   "Diagram",
					IsEmbed:  true,
				},
			},
		},
		{
			name:    "Complex wikilink",
			content: "See [[Frame Phase#Overview|frame overview]] for more.",
			expected: []*obsidian.ParsedLink{
				{
					Original: "[[Frame Phase#Overview|frame overview]]",
					Target:   "Frame Phase",
					Heading:  "Overview",
					Alias:    "frame overview",
					IsEmbed:  false,
				},
			},
		},
		{
			name:    "Multiple wikilinks",
			content: "From [[Frame Phase]] to [[Design Phase]] and then [[Test Phase]].",
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
					Original: "[[Test Phase]]",
					Target:   "Test Phase",
					IsEmbed:  false,
				},
			},
		},
		{
			name:     "No wikilinks",
			content:  "Just regular text with [markdown links](./file.md).",
			expected: []*obsidian.ParsedLink{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseWikilinks(tt.content)

			if len(result) != len(tt.expected) {
				t.Errorf("ParseWikilinks() returned %d links, expected %d", len(result), len(tt.expected))
				return
			}

			for i, link := range result {
				expected := tt.expected[i]
				if link.Original != expected.Original {
					t.Errorf("Link %d: Original = '%s', expected '%s'", i, link.Original, expected.Original)
				}
				if link.Target != expected.Target {
					t.Errorf("Link %d: Target = '%s', expected '%s'", i, link.Target, expected.Target)
				}
				if link.Alias != expected.Alias {
					t.Errorf("Link %d: Alias = '%s', expected '%s'", i, link.Alias, expected.Alias)
				}
				if link.Heading != expected.Heading {
					t.Errorf("Link %d: Heading = '%s', expected '%s'", i, link.Heading, expected.Heading)
				}
				if link.BlockID != expected.BlockID {
					t.Errorf("Link %d: BlockID = '%s', expected '%s'", i, link.BlockID, expected.BlockID)
				}
				if link.IsEmbed != expected.IsEmbed {
					t.Errorf("Link %d: IsEmbed = %t, expected %t", i, link.IsEmbed, expected.IsEmbed)
				}
			}
		})
	}
}

func TestParsedLinkString(t *testing.T) {
	tests := []struct {
		link     *obsidian.ParsedLink
		expected string
	}{
		{
			link: &obsidian.ParsedLink{
				Target: "Frame Phase",
			},
			expected: "[[Frame Phase]]",
		},
		{
			link: &obsidian.ParsedLink{
				Target: "Frame Phase",
				Alias:  "frame",
			},
			expected: "[[Frame Phase|frame]]",
		},
		{
			link: &obsidian.ParsedLink{
				Target:  "Frame Phase",
				Heading: "Overview",
			},
			expected: "[[Frame Phase#Overview]]",
		},
		{
			link: &obsidian.ParsedLink{
				Target:  "Frame Phase",
				BlockID: "block123",
			},
			expected: "[[Frame Phase^block123]]",
		},
		{
			link: &obsidian.ParsedLink{
				Target:  "Diagram",
				IsEmbed: true,
			},
			expected: "![[Diagram]]",
		},
		{
			link: &obsidian.ParsedLink{
				Target:  "Frame Phase",
				Heading: "Overview",
				Alias:   "frame overview",
			},
			expected: "[[Frame Phase#Overview|frame overview]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.link.String()
			if result != tt.expected {
				t.Errorf("ParsedLink.String() = '%s', expected '%s'", result, tt.expected)
			}
		})
	}
}

func TestValidateWikilinks(t *testing.T) {
	converter := NewLinkConverter()

	// Build test index
	files := []*obsidian.MarkdownFile{
		{
			Path: "workflows/helix/phases/01-frame/README.md",
			Frontmatter: &obsidian.Frontmatter{
				Title:   "Frame Phase",
				Aliases: []string{"Frame"},
			},
		},
		{
			Path: "workflows/helix/phases/02-design/README.md",
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
			name:           "Valid wikilinks",
			content:        "See [[Frame Phase]] and [[Design Phase]] for details.",
			expectedBroken: []string{},
		},
		{
			name:           "Valid alias reference",
			content:        "The [[Frame]] is the first phase.",
			expectedBroken: []string{},
		},
		{
			name:           "Broken wikilink",
			content:        "See [[Unknown Phase]] for details.",
			expectedBroken: []string{"Unknown Phase"},
		},
		{
			name:           "Mixed valid and broken",
			content:        "From [[Frame Phase]] to [[Unknown Phase]] to [[Design Phase]].",
			expectedBroken: []string{"Unknown Phase"},
		},
		{
			name:           "Multiple broken links",
			content:        "See [[Unknown A]] and [[Unknown B]] for details.",
			expectedBroken: []string{"Unknown A", "Unknown B"},
		},
		{
			name:           "No wikilinks",
			content:        "Just regular text here.",
			expectedBroken: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ValidateWikilinks(tt.content)

			if len(result) != len(tt.expectedBroken) {
				t.Errorf("ValidateWikilinks() returned %d broken links, expected %d", len(result), len(tt.expectedBroken))
				t.Errorf("Got: %v, Expected: %v", result, tt.expectedBroken)
				return
			}

			for i, broken := range result {
				if broken != tt.expectedBroken[i] {
					t.Errorf("Broken link %d = '%s', expected '%s'", i, broken, tt.expectedBroken[i])
				}
			}
		})
	}
}

// Helper function for tests
func minTest(a, b int) int {
	if a < b {
		return a
	}
	return b
}
