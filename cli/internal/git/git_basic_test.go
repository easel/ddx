package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test git repository
func setupTestGitRepo(t *testing.T) string {
	tempDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to init git repo: %s", string(output))

	// Configure git user for commits
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	require.NoError(t, cmd.Run())

	// Create initial commit
	testFile := filepath.Join(tempDir, "README.md")
	require.NoError(t, os.WriteFile(testFile, []byte("# Test Repo"), 0644))

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tempDir
	require.NoError(t, cmd.Run())

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "Failed to commit: %s", string(output))

	return tempDir
}

// TestIsRepository tests checking if a directory is a git repository
func TestIsRepository_Basic(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() string
		expected bool
	}{
		{
			name: "valid git repository",
			setup: func() string {
				return setupTestGitRepo(t)
			},
			expected: true,
		},
		{
			name: "non-git directory",
			setup: func() string {
				return t.TempDir()
			},
			expected: false,
		},
		{
			name: "non-existent directory",
			setup: func() string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := IsRepository(path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestHasSubtree tests checking for existing subtrees
func TestHasSubtree_Basic(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(repoDir))

	// Test when no subtree exists
	exists, err := HasSubtree(".ddx")
	assert.NoError(t, err)
	assert.False(t, exists)

	// Create a commit with git-subtree-dir marker to simulate subtree
	cmd := exec.Command("git", "commit", "--allow-empty", "-m",
		"Add subtree\n\ngit-subtree-dir: .ddx\ngit-subtree-split: abc123")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to create subtree commit: %s", string(output))

	// Test when subtree exists
	exists, err = HasSubtree(".ddx")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test different prefix
	exists, err = HasSubtree("other")
	assert.NoError(t, err)
	assert.False(t, exists)
}

// TestGetCurrentBranch tests getting the current branch name
func TestGetCurrentBranch_Basic(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(repoDir))

	// Check default branch (master or main depending on git version)
	branch, err := GetCurrentBranch()
	assert.NoError(t, err)
	assert.Contains(t, []string{"master", "main"}, branch)

	// Create and switch to a new branch
	cmd := exec.Command("git", "checkout", "-b", "feature-test")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to create branch: %s", string(output))

	branch, err = GetCurrentBranch()
	assert.NoError(t, err)
	assert.Equal(t, "feature-test", branch)
}

// TestHasUncommittedChanges tests checking for uncommitted changes
func TestHasUncommittedChanges_Basic(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(repoDir))

	// Clean repository
	hasChanges, err := HasUncommittedChanges(".")
	assert.NoError(t, err)
	assert.False(t, hasChanges)

	// Add a new file
	require.NoError(t, os.WriteFile("new.txt", []byte("new content"), 0644))

	hasChanges, err = HasUncommittedChanges(".")
	assert.NoError(t, err)
	assert.True(t, hasChanges)

	// Stage the file
	cmd := exec.Command("git", "add", "new.txt")
	require.NoError(t, cmd.Run())

	// Should still have uncommitted changes (staged but not committed)
	hasChanges, err = HasUncommittedChanges(".")
	assert.NoError(t, err)
	assert.True(t, hasChanges)

	// Commit the changes
	cmd = exec.Command("git", "commit", "-m", "Add new file")
	require.NoError(t, cmd.Run())

	// Now should be clean
	hasChanges, err = HasUncommittedChanges(".")
	assert.NoError(t, err)
	assert.False(t, hasChanges)
}

// TestCommitChanges tests committing changes
func TestCommitChanges_Basic(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	require.NoError(t, os.Chdir(repoDir))

	// Create a new file
	require.NoError(t, os.WriteFile("test.txt", []byte("test content"), 0644))

	// Commit the changes
	err := CommitChanges("Test commit")
	assert.NoError(t, err)

	// Verify the file was committed
	hasChanges, err := HasUncommittedChanges(".")
	assert.NoError(t, err)
	assert.False(t, hasChanges)

	// Verify commit message
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	output, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(output), "Test commit")
}
