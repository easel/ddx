package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// IsRepository checks if the current directory is a git repository
func IsRepository(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// HasSubtree checks if a subtree exists for the given prefix
func HasSubtree(prefix string) (bool, error) {
	cmd := exec.Command("git", "log", "--grep=git-subtree-dir: "+prefix, "--oneline")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// SubtreeAdd adds a subtree
func SubtreeAdd(prefix, repoURL, branch string) error {
	cmd := exec.Command("git", "subtree", "add", 
		"--prefix="+prefix, 
		repoURL, 
		branch, 
		"--squash")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git subtree add failed: %s", string(output))
	}
	
	return nil
}

// SubtreePull pulls updates for a subtree
func SubtreePull(prefix, repoURL, branch string) error {
	cmd := exec.Command("git", "subtree", "pull", 
		"--prefix="+prefix, 
		repoURL, 
		branch, 
		"--squash")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git subtree pull failed: %s", string(output))
	}
	
	return nil
}

// SubtreePush pushes subtree changes to a branch
func SubtreePush(prefix, repoURL, branch string) error {
	cmd := exec.Command("git", "subtree", "push", 
		"--prefix="+prefix, 
		repoURL, 
		branch)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git subtree push failed: %s", string(output))
	}
	
	return nil
}

// SubtreeReset resets subtree to remote state
func SubtreeReset(prefix, repoURL, branch string) error {
	// Remove the subtree
	cmd := exec.Command("git", "rm", "-r", prefix)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove subtree: %w", err)
	}
	
	// Commit the removal
	cmd = exec.Command("git", "commit", "-m", "Remove DDx subtree for reset")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit subtree removal: %w", err)
	}
	
	// Add it back fresh
	return SubtreeAdd(prefix, repoURL, branch)
}

// CheckBehind checks how many commits behind the subtree is
func CheckBehind(prefix, repoURL, branch string) (int, error) {
	// This is a simplified implementation
	// In a real implementation, you'd want to fetch and compare commit hashes
	
	// Fetch the latest from remote
	cmd := exec.Command("git", "fetch", repoURL, branch)
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to fetch: %w", err)
	}
	
	// Get the commit count difference
	cmd = exec.Command("git", "rev-list", "--count", "HEAD..FETCH_HEAD", "--", prefix)
	output, err := cmd.Output()
	if err != nil {
		// If this fails, assume we're up to date
		return 0, nil
	}
	
	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, nil
	}
	
	return count, nil
}

// HasUncommittedChanges checks if there are uncommitted changes in a directory
func HasUncommittedChanges(path string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain", path)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// GetCurrentBranch returns the current git branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(output)), nil
}

// CommitChanges commits changes with a message
func CommitChanges(message string) error {
	// Add all changes
	cmd := exec.Command("git", "add", "-A")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}
	
	// Commit
	cmd = exec.Command("git", "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to commit: %s", string(output))
	}
	
	return nil
}