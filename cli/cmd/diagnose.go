package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/easel/ddx/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	diagnoseReport bool
	diagnoseFix    bool
)

type DiagnosticResult struct {
	DDx         DDxStatus     `json:"ddx"`
	Git         GitStatus     `json:"git"`
	Project     ProjectStatus `json:"project"`
	AI          AIStatus      `json:"ai"`
	Score       int           `json:"score"`
	Issues      []string      `json:"issues"`
	Suggestions []string      `json:"suggestions"`
	Timestamp   time.Time     `json:"timestamp"`
	ProjectPath string        `json:"project_path"`
}

type DDxStatus struct {
	Installed   bool   `json:"installed"`
	Initialized bool   `json:"initialized"`
	ConfigValid bool   `json:"config_valid"`
	Version     string `json:"version,omitempty"`
}

type GitStatus struct {
	Repository bool     `json:"repository"`
	Gitignore  bool     `json:"gitignore"`
	Readme     bool     `json:"readme"`
	Hooks      []string `json:"hooks"`
	DDxSubtree bool     `json:"ddx_subtree"`
}

type ProjectStatus struct {
	Type       string   `json:"type"`
	ConfigFile string   `json:"config_file,omitempty"`
	SourceDirs []string `json:"source_dirs"`
	TestDirs   []string `json:"test_dirs"`
}

type AIStatus struct {
	ClaudeFile    bool     `json:"claude_file"`
	Files         []string `json:"files"`
	Documentation []string `json:"documentation"`
}

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Analyze current project setup and suggest improvements",
	Long: `Analyze the current project setup and provide recommendations.

This command checks:
• DDx installation and configuration
• Git repository setup
• Project structure and conventions  
• AI integration and documentation
• Overall development environment health`,
	RunE: runDiagnose,
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)

	diagnoseCmd.Flags().BoolVarP(&diagnoseReport, "report", "r", false, "Generate detailed report")
	diagnoseCmd.Flags().BoolVarP(&diagnoseFix, "fix", "f", false, "Automatically fix common issues")
}

func runDiagnose(cmd *cobra.Command, args []string) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	cyan.Println("🔍 Diagnosing project setup...")
	fmt.Println()

	s := spinner.New(spinner.CharSets[14], 100)
	s.Prefix = "Analyzing project... "
	s.Start()

	pwd, _ := os.Getwd()
	result := &DiagnosticResult{
		Issues:      []string{},
		Suggestions: []string{},
		Timestamp:   time.Now(),
		ProjectPath: pwd,
	}

	// Check DDx setup
	checkDDxSetup(result)

	// Check Git setup
	checkGitSetup(result)

	// Check project structure
	checkProjectStructure(result)

	// Check AI integration
	checkAIIntegration(result)

	// Calculate overall score
	calculateScore(result)

	s.Stop()
	green.Println("✅ Diagnosis complete!")
	fmt.Println()

	// Display results
	displayResults(result)

	// Generate report if requested
	if diagnoseReport {
		fmt.Println()
		if err := generateReport(result); err != nil {
			yellow := color.New(color.FgYellow)
			yellow.Printf("⚠️  Failed to generate report: %v\n", err)
		}
	}

	// Auto-fix if requested
	if diagnoseFix {
		fmt.Println()
		return autoFix(result)
	}

	return nil
}

func checkDDxSetup(result *DiagnosticResult) {
	ddxHome := getDDxHome()
	result.DDx.Installed = fileExists(ddxHome)
	result.DDx.Initialized = isInitialized()

	if !result.DDx.Installed {
		result.Issues = append(result.Issues, "DDx toolkit not installed")
	}

	if result.DDx.Initialized {
		if _, err := config.LoadLocal(); err == nil {
			result.DDx.ConfigValid = true
		} else {
			result.DDx.ConfigValid = false
			result.Issues = append(result.Issues, "DDx configuration is invalid")
		}
	} else {
		result.Issues = append(result.Issues, "DDx not initialized in this project")
		result.Suggestions = append(result.Suggestions, "Run 'ddx init' to initialize DDx")
	}
}

func checkGitSetup(result *DiagnosticResult) {
	// Check if it's a git repository
	result.Git.Repository = isGitRepository()

	if result.Git.Repository {
		result.Git.Gitignore = fileExists(".gitignore")
		result.Git.Readme = fileExists("README.md")

		// Check git hooks
		hooksDir := ".git/hooks"
		if fileExists(hooksDir) {
			if entries, err := os.ReadDir(hooksDir); err == nil {
				for _, entry := range entries {
					if !strings.HasSuffix(entry.Name(), ".sample") {
						result.Git.Hooks = append(result.Git.Hooks, entry.Name())
					}
				}
			}
		}

		// Check for DDx subtree (simplified check)
		if out, err := exec.Command("git", "log", "--grep=git-subtree-dir: .ddx", "--oneline").Output(); err == nil {
			result.Git.DDxSubtree = len(strings.TrimSpace(string(out))) > 0
		}

	} else {
		result.Issues = append(result.Issues, "Not a Git repository")
		result.Suggestions = append(result.Suggestions, "Initialize Git repository with 'git init'")
	}
}

func checkProjectStructure(result *DiagnosticResult) {
	// Common project files that indicate project type
	projectFiles := map[string]string{
		"package.json":     "Node.js/JavaScript",
		"requirements.txt": "Python",
		"Cargo.toml":       "Rust",
		"go.mod":           "Go",
		"pom.xml":          "Java (Maven)",
		"build.gradle":     "Java (Gradle)",
	}

	result.Project.Type = "unknown"
	for file, projectType := range projectFiles {
		if fileExists(file) {
			result.Project.ConfigFile = file
			result.Project.Type = projectType
			break
		}
	}

	// Check for source directories
	sourceDirs := []string{"src", "lib", "app", "components"}
	for _, dir := range sourceDirs {
		if fileExists(dir) {
			result.Project.SourceDirs = append(result.Project.SourceDirs, dir)
		}
	}

	// Check for test directories
	testDirs := []string{"test", "tests", "__tests__", "spec"}
	for _, dir := range testDirs {
		if fileExists(dir) {
			result.Project.TestDirs = append(result.Project.TestDirs, dir)
		}
	}

	if len(result.Project.TestDirs) == 0 {
		result.Suggestions = append(result.Suggestions, "Consider adding a tests directory")
	}
}

func checkAIIntegration(result *DiagnosticResult) {
	result.AI.ClaudeFile = fileExists("CLAUDE.md")

	if !result.AI.ClaudeFile {
		result.Suggestions = append(result.Suggestions, "Create CLAUDE.md for AI context")
	}

	// Check for AI-related files
	aiFiles := []string{".cursor-rules", ".ai-instructions", "prompts.md"}
	for _, file := range aiFiles {
		if fileExists(file) {
			result.AI.Files = append(result.AI.Files, file)
		}
	}

	// Check for documentation
	docFiles := []string{"README.md", "docs", "CONTRIBUTING.md"}
	for _, file := range docFiles {
		if fileExists(file) {
			result.AI.Documentation = append(result.AI.Documentation, file)
		}
	}
}

func calculateScore(result *DiagnosticResult) {
	score := 0
	maxScore := 100

	// DDx setup (30 points)
	if result.DDx.Installed {
		score += 10
	}
	if result.DDx.Initialized {
		score += 10
	}
	if result.DDx.ConfigValid {
		score += 10
	}

	// Git setup (25 points)
	if result.Git.Repository {
		score += 10
	}
	if result.Git.Gitignore {
		score += 5
	}
	if result.Git.Readme {
		score += 5
	}
	if len(result.Git.Hooks) > 0 {
		score += 5
	}

	// Project structure (25 points)
	if result.Project.Type != "unknown" {
		score += 10
	}
	if len(result.Project.SourceDirs) > 0 {
		score += 8
	}
	if len(result.Project.TestDirs) > 0 {
		score += 7
	}

	// AI integration (20 points)
	if result.AI.ClaudeFile {
		score += 10
	}
	if len(result.AI.Files) > 0 {
		score += 5
	}
	if len(result.AI.Documentation) > 0 {
		score += 5
	}

	result.Score = (score * 100) / maxScore
}

func displayResults(result *DiagnosticResult) {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	// Overall score
	bold.Println("📊 Diagnosis Results")
	fmt.Println()

	scoreColor := green
	if result.Score < 60 {
		scoreColor = red
	} else if result.Score < 80 {
		scoreColor = yellow
	}
	scoreColor.Printf("Overall Score: %d/100\n\n", result.Score)

	// DDx Status
	bold.Println("🔧 DDx Status:")
	fmt.Printf("  Installation: %s %s\n", getStatusIcon(result.DDx.Installed), boolToText(result.DDx.Installed, "Installed", "Not installed"))
	fmt.Printf("  Initialization: %s %s\n", getStatusIcon(result.DDx.Initialized), boolToText(result.DDx.Initialized, "Initialized", "Not initialized"))
	if result.DDx.Initialized {
		fmt.Printf("  Configuration: %s %s\n", getStatusIcon(result.DDx.ConfigValid), boolToText(result.DDx.ConfigValid, "Valid", "Invalid"))
	}
	fmt.Println()

	// Git Status
	bold.Println("📦 Git Status:")
	fmt.Printf("  Repository: %s %s\n", getStatusIcon(result.Git.Repository), boolToText(result.Git.Repository, "Initialized", "Not a Git repo"))
	if result.Git.Repository {
		fmt.Printf("  .gitignore: %s %s\n", getStatusIcon(result.Git.Gitignore), boolToText(result.Git.Gitignore, "Present", "Missing"))
		fmt.Printf("  README.md: %s %s\n", getStatusIcon(result.Git.Readme), boolToText(result.Git.Readme, "Present", "Missing"))
		fmt.Printf("  Git hooks: %s %d configured\n", getStatusIcon(len(result.Git.Hooks) > 0), len(result.Git.Hooks))
	}
	fmt.Println()

	// Project Structure
	bold.Println("🏗️  Project Structure:")
	fmt.Printf("  Type: %s\n", result.Project.Type)
	fmt.Printf("  Source directories: %s\n", stringSliceOrNone(result.Project.SourceDirs))
	fmt.Printf("  Test directories: %s\n", stringSliceOrNone(result.Project.TestDirs))
	fmt.Println()

	// AI Integration
	bold.Println("🤖 AI Integration:")
	fmt.Printf("  CLAUDE.md: %s %s\n", getStatusIcon(result.AI.ClaudeFile), boolToText(result.AI.ClaudeFile, "Present", "Missing"))
	fmt.Printf("  AI files: %s\n", stringSliceOrNone(result.AI.Files))
	fmt.Printf("  Documentation: %s\n", stringSliceOrNone(result.AI.Documentation))
	fmt.Println()

	// Issues
	if len(result.Issues) > 0 {
		red.Println("⚠️  Issues Found:")
		for _, issue := range result.Issues {
			red.Printf("  • %s\n", issue)
		}
		fmt.Println()
	}

	// Suggestions
	if len(result.Suggestions) > 0 {
		yellow.Println("💡 Suggestions:")
		for _, suggestion := range result.Suggestions {
			yellow.Printf("  • %s\n", suggestion)
		}
		fmt.Println()
	}
}

func autoFix(result *DiagnosticResult) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	cyan.Println("🔧 Auto-fixing issues...")
	fmt.Println()

	fixed := 0

	// Auto-fix missing .gitignore
	if result.Git.Repository && !result.Git.Gitignore {
		if err := createBasicGitignore(); err == nil {
			green.Println("✅ Created .gitignore file")
			fixed++
		}
	}

	// Auto-fix missing CLAUDE.md
	if !result.AI.ClaudeFile {
		if err := createBasicClaudeFile(); err == nil {
			green.Println("✅ Created CLAUDE.md template")
			fixed++
		}
	}

	if fixed == 0 {
		green.Println("No auto-fixable issues found!")
	} else {
		green.Printf("Fixed %d issues!\n", fixed)
	}

	return nil
}

// Helper functions
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isGitRepository() bool {
	_, err := os.Stat(".git")
	return err == nil
}

func getStatusIcon(status bool) string {
	if status {
		return color.GreenString("✓")
	}
	return color.RedString("✗")
}

func boolToText(value bool, trueText, falseText string) string {
	if value {
		return trueText
	}
	return falseText
}

func stringSliceOrNone(slice []string) string {
	if len(slice) == 0 {
		return "None found"
	}
	return strings.Join(slice, ", ")
}

func createBasicGitignore() error {
	content := `# Dependencies
node_modules/
vendor/
__pycache__/

# Build outputs
dist/
build/
target/
*.exe
*.dll

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Environment variables
.env
.env.local

# Logs
*.log
logs/
`
	return os.WriteFile(".gitignore", []byte(content), 0644)
}

func createBasicClaudeFile() error {
	pwd, _ := os.Getwd()
	projectName := filepath.Base(pwd)

	content := fmt.Sprintf(`# %s

## Overview
Brief description of what this project does.

## Tech Stack
- Language/Framework: 
- Key dependencies: 

## Architecture
Describe the high-level architecture and key components.

## Development Guidelines
- Code style: 
- Testing approach: 
- Key patterns to follow: 

## AI Development Notes
- Preferred AI model: Claude 3 Opus
- Key context for AI assistance: 
- Important files to review: README.md, package.json (or equivalent)
`, projectName)

	return os.WriteFile("CLAUDE.md", []byte(content), 0644)
}

func generateReport(result *DiagnosticResult) error {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)

	cyan.Println("📊 Generating diagnostic report...")

	// Generate JSON report
	reportData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report data: %w", err)
	}

	// Create reports directory if it doesn't exist
	reportsDir := ".ddx-reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := result.Timestamp.Format("2006-01-02_15-04-05")
	reportPath := filepath.Join(reportsDir, fmt.Sprintf("diagnostic-report_%s.json", timestamp))

	// Write JSON report
	if err := os.WriteFile(reportPath, reportData, 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	green.Printf("✅ Report saved: %s\n", reportPath)

	// Also generate a human-readable summary
	summaryPath := filepath.Join(reportsDir, fmt.Sprintf("diagnostic-summary_%s.md", timestamp))
	if err := generateMarkdownReport(result, summaryPath); err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	green.Printf("✅ Summary saved: %s\n", summaryPath)

	return nil
}

func generateMarkdownReport(result *DiagnosticResult, path string) error {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("# DDx Diagnostic Report\n\n"))
	content.WriteString(fmt.Sprintf("**Generated:** %s\n", result.Timestamp.Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("**Project Path:** %s\n", result.ProjectPath))
	content.WriteString(fmt.Sprintf("**Overall Score:** %d/100\n\n", result.Score))

	// DDx Status
	content.WriteString("## DDx Status\n\n")
	content.WriteString(fmt.Sprintf("- **Installed:** %s\n", boolToIcon(result.DDx.Installed)))
	content.WriteString(fmt.Sprintf("- **Initialized:** %s\n", boolToIcon(result.DDx.Initialized)))
	content.WriteString(fmt.Sprintf("- **Configuration Valid:** %s\n\n", boolToIcon(result.DDx.ConfigValid)))

	// Git Status
	content.WriteString("## Git Status\n\n")
	content.WriteString(fmt.Sprintf("- **Repository:** %s\n", boolToIcon(result.Git.Repository)))
	content.WriteString(fmt.Sprintf("- **Gitignore:** %s\n", boolToIcon(result.Git.Gitignore)))
	content.WriteString(fmt.Sprintf("- **README:** %s\n", boolToIcon(result.Git.Readme)))
	content.WriteString(fmt.Sprintf("- **Git Hooks:** %d configured\n", len(result.Git.Hooks)))
	content.WriteString(fmt.Sprintf("- **DDx Subtree:** %s\n\n", boolToIcon(result.Git.DDxSubtree)))

	// Project Status
	content.WriteString("## Project Structure\n\n")
	content.WriteString(fmt.Sprintf("- **Type:** %s\n", result.Project.Type))
	content.WriteString(fmt.Sprintf("- **Config File:** %s\n", result.Project.ConfigFile))
	content.WriteString(fmt.Sprintf("- **Source Directories:** %s\n", formatSlice(result.Project.SourceDirs)))
	content.WriteString(fmt.Sprintf("- **Test Directories:** %s\n\n", formatSlice(result.Project.TestDirs)))

	// AI Integration
	content.WriteString("## AI Integration\n\n")
	content.WriteString(fmt.Sprintf("- **CLAUDE.md:** %s\n", boolToIcon(result.AI.ClaudeFile)))
	content.WriteString(fmt.Sprintf("- **AI Files:** %s\n", formatSlice(result.AI.Files)))
	content.WriteString(fmt.Sprintf("- **Documentation:** %s\n\n", formatSlice(result.AI.Documentation)))

	// Issues
	if len(result.Issues) > 0 {
		content.WriteString("## Issues Found\n\n")
		for _, issue := range result.Issues {
			content.WriteString(fmt.Sprintf("- ❌ %s\n", issue))
		}
		content.WriteString("\n")
	}

	// Suggestions
	if len(result.Suggestions) > 0 {
		content.WriteString("## Suggestions\n\n")
		for _, suggestion := range result.Suggestions {
			content.WriteString(fmt.Sprintf("- 💡 %s\n", suggestion))
		}
		content.WriteString("\n")
	}

	return os.WriteFile(path, []byte(content.String()), 0644)
}

func boolToIcon(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}

func formatSlice(slice []string) string {
	if len(slice) == 0 {
		return "None"
	}
	return strings.Join(slice, ", ")
}
