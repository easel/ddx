package obsidian

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// NavigationGenerator generates navigation hubs and indices
type NavigationGenerator struct {
	tagTree *TagTree
}

// NewNavigationGenerator creates a new navigation generator
func NewNavigationGenerator() *NavigationGenerator {
	return &NavigationGenerator{
		tagTree: NewTagTree(),
	}
}

// GenerateNavigationHub creates a comprehensive navigation hub
func (g *NavigationGenerator) GenerateNavigationHub(files []*MarkdownFile) string {
	// Build tag tree and organize data
	hub := g.organizeFiles(files)

	var content strings.Builder

	// Write frontmatter
	content.WriteString(g.generateHubFrontmatter())
	content.WriteString("\n")

	// Write main navigation
	content.WriteString(g.generateMainNavigation(hub))
	content.WriteString("\n")

	// Write phase navigation
	content.WriteString(g.generatePhaseNavigation(hub))
	content.WriteString("\n")

	// Write artifact navigation
	content.WriteString(g.generateArtifactNavigation(hub))
	content.WriteString("\n")

	// Write feature navigation
	content.WriteString(g.generateFeatureNavigation(hub))
	content.WriteString("\n")

	// Write tag navigation
	content.WriteString(g.generateTagNavigation())
	content.WriteString("\n")

	// Write quick actions
	content.WriteString(g.generateQuickActions())
	content.WriteString("\n")

	// Write search helpers
	content.WriteString(g.generateSearchHelpers())

	return content.String()
}

// organizeFiles organizes files into the navigation hub structure
func (g *NavigationGenerator) organizeFiles(files []*MarkdownFile) *NavigationHub {
	hub := &NavigationHub{
		Phases:    []*PhaseInfo{},
		Artifacts: make(map[string][]*ArtifactInfo),
		Features:  []*FeatureInfo{},
		Tags:      g.tagTree,
	}

	// Process each file
	for _, file := range files {
		// Add to tag tree
		for _, tag := range file.GetTags() {
			g.tagTree.Add(tag, file)
		}

		// Process by file type
		switch file.FileType {
		case FileTypePhase:
			if phaseInfo := g.extractPhaseInfo(file); phaseInfo != nil {
				hub.Phases = append(hub.Phases, phaseInfo)
			}

		case FileTypeArtifact, FileTypeTemplate, FileTypePrompt, FileTypeExample:
			if artifactInfo := g.extractArtifactInfo(file); artifactInfo != nil {
				category := artifactInfo.Category
				if category == "" {
					category = "general"
				}
				hub.Artifacts[category] = append(hub.Artifacts[category], artifactInfo)
			}

		case FileTypeFeature:
			if featureInfo := g.extractFeatureInfo(file); featureInfo != nil {
				hub.Features = append(hub.Features, featureInfo)
			}
		}
	}

	// Sort phases by number
	sort.Slice(hub.Phases, func(i, j int) bool {
		return hub.Phases[i].Number < hub.Phases[j].Number
	})

	// Sort features by ID
	sort.Slice(hub.Features, func(i, j int) bool {
		return hub.Features[i].ID < hub.Features[j].ID
	})

	// Sort artifacts within each category
	for category := range hub.Artifacts {
		sort.Slice(hub.Artifacts[category], func(i, j int) bool {
			return hub.Artifacts[category][i].Title < hub.Artifacts[category][j].Title
		})
	}

	return hub
}

// extractPhaseInfo extracts phase information from a file
func (g *NavigationGenerator) extractPhaseInfo(file *MarkdownFile) *PhaseInfo {
	if file.Frontmatter == nil {
		return nil
	}

	fm := file.Frontmatter
	phase := &PhaseInfo{
		ID:     fm.PhaseID,
		Number: fm.PhaseNum,
		Title:  fm.Title,
		Status: fm.Status,
	}

	if fm.NextPhase != "" {
		phase.Next = fm.NextPhase
	}
	if fm.PrevPhase != "" {
		phase.Previous = fm.PrevPhase
	}
	if fm.Gates != nil {
		phase.Gates = fm.Gates
	}
	if fm.Artifacts != nil {
		phase.Artifacts = fm.Artifacts
	}

	return phase
}

// extractArtifactInfo extracts artifact information from a file
func (g *NavigationGenerator) extractArtifactInfo(file *MarkdownFile) *ArtifactInfo {
	if file.Frontmatter == nil {
		return nil
	}

	fm := file.Frontmatter
	return &ArtifactInfo{
		Title:      fm.Title,
		Category:   fm.ArtifactCategory,
		Phase:      fm.Phase,
		Type:       file.FileType,
		Complexity: fm.Complexity,
		Path:       file.Path,
	}
}

// extractFeatureInfo extracts feature information from a file
func (g *NavigationGenerator) extractFeatureInfo(file *MarkdownFile) *FeatureInfo {
	if file.Frontmatter == nil {
		return nil
	}

	fm := file.Frontmatter
	return &FeatureInfo{
		ID:       fm.FeatureID,
		Title:    fm.Title,
		Status:   fm.Status,
		Priority: fm.Priority,
		Owner:    fm.Owner,
		Phase:    fm.WorkflowPhase,
	}
}

// generateHubFrontmatter generates frontmatter for the navigation hub
func (g *NavigationGenerator) generateHubFrontmatter() string {
	return fmt.Sprintf(`---
title: "HELIX Workflow Navigator"
type: navigation-hub
tags:
  - helix/core
  - helix/navigation
  - dashboard
created: %s
updated: %s
aliases:
  - "HELIX Navigator"
  - "Workflow Dashboard"
  - "HELIX Dashboard"
---`, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
}

// generateMainNavigation generates the main workflow navigation
func (g *NavigationGenerator) generateMainNavigation(hub *NavigationHub) string {
	var content strings.Builder

	content.WriteString("# HELIX Workflow Navigator\n\n")
	content.WriteString("Welcome to the HELIX (Human-Enhanced Learning and Implementation eXperience) workflow navigator. ")
	content.WriteString("This dashboard provides quick access to all phases, artifacts, and resources in the workflow.\n\n")

	content.WriteString("## 🔄 The HELIX Spiral\n\n")
	content.WriteString("```mermaid\n")
	content.WriteString("graph TB\n")
	content.WriteString(`    subgraph "The HELIX Spiral"
        F[📋 FRAME] --> D[🏗️ DESIGN]
        D --> T[🧪 TEST]
        T --> B[⚙️ BUILD]
        B --> DP[🚀 DEPLOY]
        DP --> I[🔄 ITERATE]
        I -.->|Next Cycle| F
    end

    style F fill:#e1f5fe
    style D fill:#f3e5f5
    style T fill:#ffebee
    style B fill:#e8f5e9
    style DP fill:#fff3e0
    style I fill:#fce4ec`)
	content.WriteString("\n```\n\n")

	return content.String()
}

// generatePhaseNavigation generates phase-by-phase navigation
func (g *NavigationGenerator) generatePhaseNavigation(hub *NavigationHub) string {
	var content strings.Builder

	content.WriteString("## 📋 Workflow Phases\n\n")

	if len(hub.Phases) > 0 {
		content.WriteString("| Phase | Status | Description | Next Phase |\n")
		content.WriteString("|-------|--------|-------------|------------|\n")

		for _, phase := range hub.Phases {
			status := "⏸️ Not Started"
			if phase.Status == "in_progress" {
				status = "🔄 In Progress"
			} else if phase.Status == "completed" {
				status = "✅ Completed"
			}

			description := g.getPhaseDescription(phase.ID)
			nextPhase := "—"
			if phase.Next != "" {
				nextPhase = phase.Next
			}

			content.WriteString(fmt.Sprintf("| [[%s]] | %s | %s | %s |\n",
				phase.Title, status, description, nextPhase))
		}
	} else {
		// Default phase structure if no phase files found
		phases := []struct {
			name        string
			emoji       string
			description string
		}{
			{"Frame Phase", "📋", "Define the problem and establish context"},
			{"Design Phase", "🏗️", "Architect the solution approach"},
			{"Test Phase", "🧪", "Write failing tests (Red phase)"},
			{"Build Phase", "⚙️", "Implement code to pass tests (Green phase)"},
			{"Deploy Phase", "🚀", "Release to production with monitoring"},
			{"Iterate Phase", "🔄", "Learn and improve for next cycle"},
		}

		for i, phase := range phases {
			nextPhase := "—"
			if i < len(phases)-1 {
				nextPhase = fmt.Sprintf("[[%s]]", phases[i+1].name)
			}

			content.WriteString(fmt.Sprintf("| [[%s]] | ⏸️ Not Started | %s %s | %s |\n",
				phase.name, phase.emoji, phase.description, nextPhase))
		}
	}

	content.WriteString("\n")
	return content.String()
}

// generateArtifactNavigation generates artifact navigation by category
func (g *NavigationGenerator) generateArtifactNavigation(hub *NavigationHub) string {
	var content strings.Builder

	content.WriteString("## 📚 Artifacts by Category\n\n")

	if len(hub.Artifacts) > 0 {
		// Sort categories
		var categories []string
		for category := range hub.Artifacts {
			categories = append(categories, category)
		}
		sort.Strings(categories)

		for _, category := range categories {
			artifacts := hub.Artifacts[category]
			if len(artifacts) == 0 {
				continue
			}

			categoryTitle := strings.Title(strings.ReplaceAll(category, "-", " "))
			content.WriteString(fmt.Sprintf("### %s\n\n", categoryTitle))

			for _, artifact := range artifacts {
				complexity := ""
				if artifact.Complexity != "" {
					switch artifact.Complexity {
					case "simple":
						complexity = " 🟢"
					case "moderate":
						complexity = " 🟡"
					case "complex":
						complexity = " 🔴"
					}
				}

				phase := ""
				if artifact.Phase != "" {
					phase = fmt.Sprintf(" *(Phase: %s)*", strings.Title(artifact.Phase))
				}

				content.WriteString(fmt.Sprintf("- [[%s]]%s%s\n",
					artifact.Title, complexity, phase))
			}
			content.WriteString("\n")
		}
	} else {
		content.WriteString("*No artifacts found. Use `ddx obsidian migrate` to scan for artifacts.*\n\n")
	}

	return content.String()
}

// generateFeatureNavigation generates feature navigation
func (g *NavigationGenerator) generateFeatureNavigation(hub *NavigationHub) string {
	var content strings.Builder

	content.WriteString("## 🎯 Active Features\n\n")

	if len(hub.Features) > 0 {
		content.WriteString("| Feature ID | Title | Status | Priority | Owner | Phase |\n")
		content.WriteString("|------------|-------|--------|----------|-------|-------|\n")

		for _, feature := range hub.Features {
			status := "📝 Draft"
			switch feature.Status {
			case "specified":
				status = "📋 Specified"
			case "approved":
				status = "✅ Approved"
			case "in_progress":
				status = "🔄 In Progress"
			case "completed":
				status = "✅ Completed"
			case "deprecated":
				status = "❌ Deprecated"
			}

			priority := feature.Priority
			if priority == "" {
				priority = "P2"
			}

			owner := feature.Owner
			if owner == "" {
				owner = "—"
			}

			phase := feature.Phase
			if phase == "" {
				phase = "—"
			}

			content.WriteString(fmt.Sprintf("| %s | [[%s]] | %s | %s | %s | %s |\n",
				feature.ID, feature.Title, status, priority, owner, phase))
		}
	} else {
		content.WriteString("*No features found. Create feature specifications to see them here.*\n")
	}

	content.WriteString("\n")
	return content.String()
}

// generateTagNavigation generates tag-based navigation
func (g *NavigationGenerator) generateTagNavigation() string {
	var content strings.Builder

	content.WriteString("## 🏷️ Browse by Tags\n\n")

	content.WriteString("### Phase Tags\n")
	phases := []string{"frame", "design", "test", "build", "deploy", "iterate"}
	for _, phase := range phases {
		tag := fmt.Sprintf("#helix/phase/%s", phase)
		content.WriteString(fmt.Sprintf("- %s\n", tag))
	}

	content.WriteString("\n### Artifact Tags\n")
	artifactTypes := []string{"specification", "design", "test", "implementation", "template", "prompt", "example"}
	for _, aType := range artifactTypes {
		tag := fmt.Sprintf("#helix/artifact/%s", aType)
		content.WriteString(fmt.Sprintf("- %s\n", tag))
	}

	content.WriteString("\n### Complexity Tags\n")
	complexities := []string{"simple", "moderate", "complex"}
	for _, complexity := range complexities {
		tag := fmt.Sprintf("#helix/complexity/%s", complexity)
		content.WriteString(fmt.Sprintf("- %s\n", tag))
	}

	content.WriteString("\n")
	return content.String()
}

// generateQuickActions generates quick action links
func (g *NavigationGenerator) generateQuickActions() string {
	var content strings.Builder

	content.WriteString("## ⚡ Quick Actions\n\n")

	actions := []struct {
		title       string
		description string
	}{
		{"Create Feature Specification", "Start defining a new feature"},
		{"Write User Stories", "Document user requirements and acceptance criteria"},
		{"Design Technical Architecture", "Plan the technical implementation"},
		{"Write Test Suite", "Create failing tests that define behavior"},
		{"Implement Solution", "Write code to make tests pass"},
		{"Deploy to Production", "Release and monitor the solution"},
		{"Review Phase Gates", "Check criteria for phase progression"},
		{"Update Documentation", "Keep artifacts current and accurate"},
	}

	for _, action := range actions {
		content.WriteString(fmt.Sprintf("- [[%s]] - %s\n", action.title, action.description))
	}

	content.WriteString("\n")
	return content.String()
}

// generateSearchHelpers generates search and filter helpers
func (g *NavigationGenerator) generateSearchHelpers() string {
	var content strings.Builder

	content.WriteString("## 🔍 Search and Filter\n\n")

	content.WriteString("### Dataview Queries\n\n")
	content.WriteString("#### All HELIX Documents\n")
	content.WriteString("```dataview\n")
	content.WriteString("TABLE file.name as \"Document\", type as \"Type\", phase as \"Phase\", status as \"Status\"\n")
	content.WriteString("FROM #helix\n")
	content.WriteString("SORT phase, type, file.name\n")
	content.WriteString("```\n\n")

	content.WriteString("#### Current Phase Artifacts\n")
	content.WriteString("```dataview\n")
	content.WriteString("LIST\n")
	content.WriteString("FROM #helix/phase/frame\n")
	content.WriteString("WHERE type != \"phase\"\n")
	content.WriteString("SORT file.name\n")
	content.WriteString("```\n\n")

	content.WriteString("#### Templates by Complexity\n")
	content.WriteString("```dataview\n")
	content.WriteString("TABLE complexity as \"Complexity\", time_estimate as \"Time\", phase as \"Phase\"\n")
	content.WriteString("FROM #helix/artifact/template\n")
	content.WriteString("SORT complexity, phase\n")
	content.WriteString("```\n\n")

	content.WriteString("#### Features by Status\n")
	content.WriteString("```dataview\n")
	content.WriteString("TABLE priority as \"Priority\", owner as \"Owner\", status as \"Status\"\n")
	content.WriteString("FROM #helix/artifact/specification\n")
	content.WriteString("WHERE feature_id\n")
	content.WriteString("SORT priority, status\n")
	content.WriteString("```\n\n")

	content.WriteString("### Graph Navigation\n\n")
	content.WriteString("Use Obsidian's graph view to visualize:\n")
	content.WriteString("- **Workflow progression**: See how phases connect\n")
	content.WriteString("- **Artifact relationships**: Understand dependencies\n")
	content.WriteString("- **Feature flow**: Track features through phases\n")
	content.WriteString("- **Knowledge clusters**: Identify related concepts\n\n")

	content.WriteString("### Tag Filters\n\n")
	content.WriteString("Click on any tag to filter the graph:\n")
	content.WriteString("- `#helix/phase/frame` - Frame phase artifacts\n")
	content.WriteString("- `#helix/artifact/template` - Template files\n")
	content.WriteString("- `#helix/complexity/simple` - Simple artifacts\n")
	content.WriteString("- `#helix/status/draft` - Draft documents\n\n")

	content.WriteString("---\n\n")
	content.WriteString("*This navigator is automatically generated. Run `ddx obsidian migrate` to update.*\n")

	return content.String()
}

// getPhaseDescription returns a description for a phase
func (g *NavigationGenerator) getPhaseDescription(phaseID string) string {
	descriptions := map[string]string{
		"frame":   "Define the problem and establish context",
		"design":  "Architect the solution approach",
		"test":    "Write failing tests (Red phase)",
		"build":   "Implement code to pass tests (Green phase)",
		"deploy":  "Release to production with monitoring",
		"iterate": "Learn and improve for next cycle",
	}

	if desc, ok := descriptions[phaseID]; ok {
		return desc
	}
	return "HELIX workflow phase"
}

// GeneratePhaseIndex generates an index for a specific phase
func (g *NavigationGenerator) GeneratePhaseIndex(phase string, files []*MarkdownFile) string {
	var content strings.Builder

	phaseTitle := strings.Title(phase)
	content.WriteString(fmt.Sprintf("# %s Phase Index\n\n", phaseTitle))

	// Filter files for this phase
	phaseFiles := []*MarkdownFile{}
	for _, file := range files {
		if file.GetPhase() == phase || GetPhaseFromPath(file.Path) == phase {
			phaseFiles = append(phaseFiles, file)
		}
	}

	if len(phaseFiles) == 0 {
		content.WriteString("*No artifacts found for this phase.*\n")
		return content.String()
	}

	// Group by type
	byType := make(map[FileType][]*MarkdownFile)
	for _, file := range phaseFiles {
		byType[file.FileType] = append(byType[file.FileType], file)
	}

	// Output by type
	typeOrder := []FileType{FileTypePhase, FileTypeEnforcer, FileTypeTemplate, FileTypePrompt, FileTypeExample, FileTypeArtifact}
	for _, fileType := range typeOrder {
		if files, ok := byType[fileType]; ok && len(files) > 0 {
			content.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(string(fileType))))
			for _, file := range files {
				title := file.GetTitle()
				content.WriteString(fmt.Sprintf("- [[%s]]\n", title))
			}
			content.WriteString("\n")
		}
	}

	return content.String()
}
