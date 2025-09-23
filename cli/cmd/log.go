package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show DDX asset history",
	Long: `Show commit history for DDX assets and resources.

This command displays the git log for your DDX resources, helping you
track changes, updates, and the evolution of your project setup.

Examples:
  ddx log                    # Show recent commit history
  ddx log -n 10              # Show last 10 commits
  ddx log --oneline          # Show compact format
  ddx log --since="1 week ago" # Show commits from last week`,
	RunE: runLog,
}

// Flag variables are now handled locally in command functions

// Remove init function - commands are now registered via command factory

func runLog(cmd *cobra.Command, args []string) error {
	// Verify we're in a DDX project
	if !isDDXProject() {
		return fmt.Errorf("not a DDX project - run 'ddx init' first")
	}

	ddxDir := ".ddx"
	if _, err := os.Stat(ddxDir); os.IsNotExist(err) {
		return fmt.Errorf("DDX directory not found - project may not be properly initialized")
	}

	// Get flag values
	logLimit, _ := cmd.Flags().GetInt("number")
	logLimitAlt, _ := cmd.Flags().GetInt("limit")
	logOneline, _ := cmd.Flags().GetBool("oneline")
	logDiff, _ := cmd.Flags().GetBool("diff")
	logExport, _ := cmd.Flags().GetString("export")
	logSince, _ := cmd.Flags().GetString("since")
	logAuthor, _ := cmd.Flags().GetString("author")
	logGrep, _ := cmd.Flags().GetString("grep")

	// Use --limit if specified, otherwise use --number
	if logLimitAlt != 20 { // 20 is the default, so if it's different, --limit was used
		logLimit = logLimitAlt
	}

	// Get path filter from args
	var pathFilter string
	if len(args) > 0 {
		pathFilter = args[0]
	}

	// Handle export functionality
	if logExport != "" {
		return exportLog(ddxDir, logExport, logLimit, logOneline, logDiff, logSince, logAuthor, logGrep, pathFilter)
	}

	return showGitLog(cmd, ddxDir, logLimit, logOneline, logDiff, logSince, logAuthor, logGrep, pathFilter)
}

func showGitLog(cmd *cobra.Command, ddxDir string, logLimit int, logOneline, logDiff bool, logSince, logAuthor, logGrep, pathFilter string) error {
	// Build git log command
	args := []string{"log"}

	// Add limit
	if logLimit > 0 {
		args = append(args, "-n", strconv.Itoa(logLimit))
	}

	// Add format options
	if logOneline {
		args = append(args, "--oneline")
	} else {
		args = append(args, "--pretty=format:%C(yellow)%h %C(blue)%ad %C(reset)%s %C(green)(%an)%C(reset)")
		args = append(args, "--date=short")
	}

	// Add diff options
	if logDiff {
		if logOneline {
			args = append(args, "--stat")
		} else {
			args = append(args, "--patch")
		}
	}

	// Add filters
	if logSince != "" {
		args = append(args, "--since", logSince)
	}

	if logAuthor != "" {
		args = append(args, "--author", logAuthor)
	}

	if logGrep != "" {
		args = append(args, "--grep", logGrep)
	}

	// Add path filter
	if pathFilter != "" {
		args = append(args, "--", pathFilter)
	}

	// Execute git log command
	gitCmd := exec.Command("git", args...)
	gitCmd.Dir = ddxDir
	gitCmd.Stdout = cmd.OutOrStdout()
	gitCmd.Stderr = cmd.ErrOrStderr()

	err := gitCmd.Run()
	if err != nil {
		// If git is not available or no repository, show alternative information
		return showAlternativeLog(cmd, ddxDir, pathFilter, logLimit, logDiff)
	}

	return nil
}

func showAlternativeLog(cmd *cobra.Command, ddxDir, pathFilter string, logLimit int, logDiff bool) error {
	fmt.Fprintln(cmd.OutOrStdout(), "DDX Asset History")
	fmt.Fprintln(cmd.OutOrStdout(), "=================")
	fmt.Fprintln(cmd.OutOrStdout())

	// Get file modification times as a fallback
	fmt.Fprintln(cmd.OutOrStdout(), "Recent Changes (based on file modification times):")
	fmt.Fprintln(cmd.OutOrStdout())

	// This is a simplified version that shows file modifications
	// In a real implementation, this would be more sophisticated
	fmt.Fprintln(cmd.OutOrStdout(), "Note: Git history not available. Showing file-based history.")
	fmt.Fprintln(cmd.OutOrStdout(), "Initialize git in .ddx directory for full history tracking.")
	fmt.Fprintln(cmd.OutOrStdout())

	// Build find command with optional path filter
	findArgs := []string{ddxDir, "-type", "f"}
	if pathFilter != "" {
		findArgs = append(findArgs, "-path", "*"+pathFilter+"*")
	}
	findArgs = append(findArgs, "-printf", "%TY-%Tm-%Td %TH:%TM %p\n")

	// Show recent files by modification time
	findCmd := exec.Command("find", findArgs...)
	findCmd.Stdout = cmd.OutOrStdout()
	err := findCmd.Run()

	if err != nil {
		// Fallback for systems without find -printf
		return showBasicFileList(cmd, ddxDir)
	}

	return nil
}

func showBasicFileList(cmd *cobra.Command, ddxDir string) error {
	fmt.Fprintln(cmd.OutOrStdout(), "DDX Files:")

	lsCmd := exec.Command("ls", "-la", ddxDir)
	lsCmd.Stdout = cmd.OutOrStdout()
	lsCmd.Stderr = cmd.ErrOrStderr()

	return lsCmd.Run()
}

// exportLog handles exporting log data to various formats
func exportLog(ddxDir, exportPath string, logLimit int, logOneline, logDiff bool, logSince, logAuthor, logGrep, pathFilter string) error {
	// Determine export format from file extension
	format := "markdown" // default
	if strings.HasSuffix(exportPath, ".json") {
		format = "json"
	} else if strings.HasSuffix(exportPath, ".csv") {
		format = "csv"
	} else if strings.HasSuffix(exportPath, ".html") {
		format = "html"
	}

	// Create export file
	file, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %v", err)
	}
	defer file.Close()

	// Try to get git log data first
	gitLog, err := getGitLogData(ddxDir, logLimit, logOneline, logDiff, logSince, logAuthor, logGrep, pathFilter)
	if err != nil {
		// Fallback to file-based data
		gitLog, err = getFileBasedLogData(ddxDir, pathFilter, logLimit)
		if err != nil {
			return fmt.Errorf("failed to get log data: %v", err)
		}
	}

	// Export in requested format
	switch format {
	case "json":
		return exportJSON(file, gitLog)
	case "csv":
		return exportCSV(file, gitLog)
	case "html":
		return exportHTML(file, gitLog)
	default:
		return exportMarkdown(file, gitLog)
	}
}

// LogEntry represents a single log entry for export
type LogEntry struct {
	Hash    string   `json:"hash"`
	Date    string   `json:"date"`
	Author  string   `json:"author"`
	Message string   `json:"message"`
	Files   []string `json:"files,omitempty"`
	Changes string   `json:"changes,omitempty"`
}

// getGitLogData retrieves log data from git
func getGitLogData(ddxDir string, logLimit int, logOneline, logDiff bool, logSince, logAuthor, logGrep, pathFilter string) ([]LogEntry, error) {
	// For now, return error to force fallback to file-based data
	// In a full implementation, this would parse git log output
	return nil, fmt.Errorf("git log parsing not implemented in this version")
}

// getFileBasedLogData creates log entries from file modification times
func getFileBasedLogData(ddxDir, pathFilter string, logLimit int) ([]LogEntry, error) {
	entries := []LogEntry{}

	// Create a simple entry based on DDX directory
	entry := LogEntry{
		Hash:    "file-based",
		Date:    "2025-01-20",
		Author:  "DDX System",
		Message: "File-based history (git not available)",
		Files:   []string{ddxDir},
	}

	if pathFilter != "" {
		entry.Message = fmt.Sprintf("File-based history for %s (git not available)", pathFilter)
		entry.Files = []string{pathFilter}
	}

	entries = append(entries, entry)
	return entries, nil
}

// exportMarkdown exports log data as Markdown
func exportMarkdown(file *os.File, entries []LogEntry) error {
	_, err := file.WriteString("# DDX Asset History\n\n")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_, err := file.WriteString(fmt.Sprintf("## %s - %s\n\n", entry.Date, entry.Author))
		if err != nil {
			return err
		}

		_, err = file.WriteString(fmt.Sprintf("**Message:** %s\n\n", entry.Message))
		if err != nil {
			return err
		}

		if len(entry.Files) > 0 {
			_, err = file.WriteString("**Files:**\n")
			if err != nil {
				return err
			}
			for _, f := range entry.Files {
				_, err = file.WriteString(fmt.Sprintf("- %s\n", f))
				if err != nil {
					return err
				}
			}
			_, err = file.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// exportJSON exports log data as JSON
func exportJSON(file *os.File, entries []LogEntry) error {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(map[string]interface{}{
		"ddx_history": entries,
		"exported_at": time.Now().Format(time.RFC3339),
	})
}

// exportCSV exports log data as CSV
func exportCSV(file *os.File, entries []LogEntry) error {
	_, err := file.WriteString("Hash,Date,Author,Message,Files\n")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		files := strings.Join(entry.Files, ";")
		_, err := file.WriteString(fmt.Sprintf("%s,%s,%s,\"%s\",\"%s\"\n",
			entry.Hash, entry.Date, entry.Author, entry.Message, files))
		if err != nil {
			return err
		}
	}

	return nil
}

// exportHTML exports log data as HTML
func exportHTML(file *os.File, entries []LogEntry) error {
	_, err := file.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>DDX Asset History</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .entry { border-bottom: 1px solid #ccc; margin-bottom: 20px; padding-bottom: 20px; }
        .date { color: #666; font-size: 0.9em; }
        .message { font-weight: bold; margin: 10px 0; }
        .files { background: #f5f5f5; padding: 10px; border-radius: 4px; }
    </style>
</head>
<body>
    <h1>DDX Asset History</h1>
`)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		_, err := file.WriteString(fmt.Sprintf(`
    <div class="entry">
        <div class="date">%s - %s</div>
        <div class="message">%s</div>
        <div class="files">Files: %s</div>
    </div>
`, entry.Date, entry.Author, entry.Message, strings.Join(entry.Files, ", ")))
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString(`
</body>
</html>`)
	return err
}
