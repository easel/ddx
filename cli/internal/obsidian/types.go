package obsidian

import (
	"strings"
	"time"
)

// Frontmatter represents the YAML frontmatter for any markdown file
type Frontmatter struct {
	Title   string    `yaml:"title"`
	Type    string    `yaml:"type"`
	Tags    []string  `yaml:"tags"`
	Created time.Time `yaml:"created"`
	Updated time.Time `yaml:"updated"`
	Status  string    `yaml:"status,omitempty"`
	Version string    `yaml:"version,omitempty"`
	Aliases []string  `yaml:"aliases,omitempty"`
	Related []string  `yaml:"related,omitempty"`

	// Phase-specific fields
	PhaseID   string     `yaml:"phase_id,omitempty"`
	PhaseNum  int        `yaml:"phase_number,omitempty"`
	NextPhase string     `yaml:"next_phase,omitempty"`
	PrevPhase string     `yaml:"previous_phase,omitempty"`
	Gates     *Gates     `yaml:"gates,omitempty"`
	Artifacts *Artifacts `yaml:"artifacts,omitempty"`

	// Artifact-specific fields
	ArtifactCategory string   `yaml:"artifact_category,omitempty"`
	Phase            string   `yaml:"phase,omitempty"`
	Complexity       string   `yaml:"complexity,omitempty"`
	Prerequisites    []string `yaml:"prerequisites,omitempty"`
	Outputs          []string `yaml:"outputs,omitempty"`
	TimeEstimate     string   `yaml:"time_estimate,omitempty"`
	SkillsRequired   []string `yaml:"skills_required,omitempty"`

	// Feature-specific fields
	FeatureID string `yaml:"feature_id,omitempty"`
	Priority  string `yaml:"priority,omitempty"`
	Owner     string `yaml:"owner,omitempty"`

	// Workflow-specific fields
	WorkflowPhase string `yaml:"workflow_phase,omitempty"`
	ArtifactType  string `yaml:"artifact_type,omitempty"`
}

type Gates struct {
	Entry []string `yaml:"entry,omitempty"`
	Exit  []string `yaml:"exit,omitempty"`
}

type Artifacts struct {
	Required []string `yaml:"required,omitempty"`
	Optional []string `yaml:"optional,omitempty"`
}

// MarkdownFile represents a markdown file with optional frontmatter
type MarkdownFile struct {
	Path        string
	Content     string
	Frontmatter *Frontmatter
	FileType    FileType
}

// FileType represents the type of HELIX file
type FileType string

const (
	FileTypeUnknown     FileType = "unknown"
	FileTypePhase       FileType = "phase"
	FileTypeEnforcer    FileType = "enforcer"
	FileTypeArtifact    FileType = "artifact"
	FileTypeTemplate    FileType = "template"
	FileTypePrompt      FileType = "prompt"
	FileTypeExample     FileType = "example"
	FileTypeCoordinator FileType = "coordinator"
	FileTypePrinciple   FileType = "principle"
	FileTypeFeature     FileType = "feature-specification"
)

// String returns the string representation of a FileType
func (ft FileType) String() string {
	return string(ft)
}

// IsPhase returns true if the file type is phase-related
func (ft FileType) IsPhase() bool {
	return ft == FileTypePhase || ft == FileTypeEnforcer
}

// IsArtifact returns true if the file type is an artifact
func (ft FileType) IsArtifact() bool {
	return ft == FileTypeArtifact || ft == FileTypeTemplate ||
		ft == FileTypePrompt || ft == FileTypeExample ||
		ft == FileTypeFeature
}

// IsCore returns true if the file type is a core workflow file
func (ft FileType) IsCore() bool {
	return ft == FileTypeCoordinator || ft == FileTypePrinciple
}

// GetHierarchicalTags returns the hierarchical tags for this file type
func (ft FileType) GetHierarchicalTags() []string {
	baseTags := []string{"helix"}

	switch ft {
	case FileTypePhase:
		return append(baseTags, "helix/phase")
	case FileTypeEnforcer:
		return append(baseTags, "helix/enforcer")
	case FileTypeTemplate:
		return append(baseTags, "helix/artifact", "helix/artifact/template")
	case FileTypePrompt:
		return append(baseTags, "helix/artifact", "helix/artifact/prompt")
	case FileTypeExample:
		return append(baseTags, "helix/artifact", "helix/artifact/example")
	case FileTypeCoordinator:
		return append(baseTags, "helix/core", "helix/coordinator")
	case FileTypePrinciple:
		return append(baseTags, "helix/core", "helix/principle")
	case FileTypeFeature:
		return append(baseTags, "helix/artifact", "helix/artifact/specification")
	default:
		return baseTags
	}
}

// HasFrontmatter returns true if the file has frontmatter
func (mf *MarkdownFile) HasFrontmatter() bool {
	return mf.Frontmatter != nil
}

// GetTitle returns the title from frontmatter or extracts from content
func (mf *MarkdownFile) GetTitle() string {
	if mf.HasFrontmatter() && mf.Frontmatter.Title != "" {
		return mf.Frontmatter.Title
	}
	// Fallback to extracting from content would go here
	return "Untitled"
}

// GetTags returns tags from frontmatter
func (mf *MarkdownFile) GetTags() []string {
	if mf.HasFrontmatter() {
		return mf.Frontmatter.Tags
	}
	return []string{}
}

// GetPhase returns the phase from frontmatter or path
func (mf *MarkdownFile) GetPhase() string {
	if mf.HasFrontmatter() && mf.Frontmatter.Phase != "" {
		return mf.Frontmatter.Phase
	}
	if mf.HasFrontmatter() && mf.Frontmatter.PhaseID != "" {
		return mf.Frontmatter.PhaseID
	}
	// Fallback to extracting from path would go here
	return ""
}

// GetArtifactCategory returns the artifact category from frontmatter
func (mf *MarkdownFile) GetArtifactCategory() string {
	if mf.HasFrontmatter() {
		return mf.Frontmatter.ArtifactCategory
	}
	return ""
}

// ParsedLink represents a parsed wikilink
type ParsedLink struct {
	Original string
	Target   string
	Alias    string
	Heading  string
	BlockID  string
	IsEmbed  bool
}

// String returns the wikilink syntax for this parsed link
func (pl *ParsedLink) String() string {
	var result string

	if pl.IsEmbed {
		result = "!"
	}

	result += "[["
	result += pl.Target

	if pl.Heading != "" {
		result += "#" + pl.Heading
	}

	if pl.BlockID != "" {
		result += "^" + pl.BlockID
	}

	if pl.Alias != "" && pl.Alias != pl.Target {
		result += "|" + pl.Alias
	}

	result += "]]"

	return result
}

// NavigationHub represents the structure for the navigation hub
type NavigationHub struct {
	Phases    []*PhaseInfo
	Artifacts map[string][]*ArtifactInfo
	Tags      *TagTree
	Features  []*FeatureInfo
}

// PhaseInfo contains information about a phase
type PhaseInfo struct {
	ID        string
	Number    int
	Title     string
	Status    string
	Next      string
	Previous  string
	Gates     *Gates
	Artifacts *Artifacts
}

// ArtifactInfo contains information about an artifact
type ArtifactInfo struct {
	Title      string
	Category   string
	Phase      string
	Type       FileType
	Complexity string
	Path       string
}

// FeatureInfo contains information about a feature
type FeatureInfo struct {
	ID       string
	Title    string
	Status   string
	Priority string
	Owner    string
	Phase    string
}

// TagTree represents a hierarchical tag structure
type TagTree struct {
	Root     *TagNode
	AllTags  map[string]*TagNode
	TagCount map[string]int
}

// TagNode represents a node in the tag hierarchy
type TagNode struct {
	Name     string
	FullPath string
	Parent   *TagNode
	Children map[string]*TagNode
	Files    []*MarkdownFile
}

// NewTagTree creates a new tag tree
func NewTagTree() *TagTree {
	return &TagTree{
		Root: &TagNode{
			Name:     "helix",
			FullPath: "helix",
			Children: make(map[string]*TagNode),
		},
		AllTags:  make(map[string]*TagNode),
		TagCount: make(map[string]int),
	}
}

// Add adds a tag to the tree
func (tt *TagTree) Add(tag string, file *MarkdownFile) {
	tt.TagCount[tag]++

	if node, exists := tt.AllTags[tag]; exists {
		node.Files = append(node.Files, file)
		return
	}

	// Create the tag hierarchy
	parts := strings.Split(tag, "/")
	current := tt.Root
	fullPath := ""

	for i, part := range parts {
		if i > 0 {
			fullPath += "/"
		}
		fullPath += part

		if child, exists := current.Children[part]; exists {
			current = child
		} else {
			newNode := &TagNode{
				Name:     part,
				FullPath: fullPath,
				Parent:   current,
				Children: make(map[string]*TagNode),
				Files:    []*MarkdownFile{},
			}
			current.Children[part] = newNode
			tt.AllTags[fullPath] = newNode
			current = newNode
		}
	}

	current.Files = append(current.Files, file)
}

// Get retrieves a tag node by its full path
func (tt *TagTree) Get(tag string) *TagNode {
	return tt.AllTags[tag]
}

// GetFilesByTag returns all files with a specific tag
func (tt *TagTree) GetFilesByTag(tag string) []*MarkdownFile {
	if node := tt.Get(tag); node != nil {
		return node.Files
	}
	return nil
}
