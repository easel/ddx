package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/easel/ddx/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

var (
	testLibraryPath string
	testLibraryOnce sync.Once
)

// getTempDir returns the appropriate temporary directory
func getTempDir() string {
	// Check DDX_TEST_LIBRARY_PATH first (explicit override)
	if path := os.Getenv("DDX_TEST_LIBRARY_PATH"); path != "" {
		return path
	}
	// Check TMPDIR (standard on macOS/BSD)
	if path := os.Getenv("TMPDIR"); path != "" {
		return path
	}
	// Check TMP (Windows)
	if path := os.Getenv("TMP"); path != "" {
		return path
	}
	// Fall back to /tmp
	return "/tmp"
}

// GetTestLibraryPath returns the absolute path to a temporary git repository
// created from the test library fixture. This is a real git repository that
// tests can use with file:// URLs.
func GetTestLibraryPath() string {
	testLibraryOnce.Do(func() {
		// Find the test fixtures directory relative to this file
		_, filename, _, _ := runtime.Caller(0)
		cmdDir := filepath.Dir(filename)
		fixtureDir := filepath.Join(cmdDir, "..", "test", "fixtures", "ddx-library")
		fixtureDir, _ = filepath.Abs(fixtureDir)

		// Create temp directory for git repository
		tempDir := getTempDir()
		testLibraryPath = filepath.Join(tempDir, ".test-ddx-library")

		// Check if temp repo already exists and is valid
		gitDir := filepath.Join(testLibraryPath, ".git")
		repoExists := false
		if stat, err := os.Stat(gitDir); err == nil && stat.IsDir() {
			repoExists = true
		}

		if !repoExists {
			// Create new temp directory and initialize git repo
			if err := os.MkdirAll(testLibraryPath, 0755); err != nil {
				panic(fmt.Sprintf("Failed to create temp library directory: %v", err))
			}
		}

		// Sync files from fixture to temp repo
		if err := syncDirectory(fixtureDir, testLibraryPath); err != nil {
			panic(fmt.Sprintf("Failed to sync fixture files: %v", err))
		}

		if !repoExists {
			// Initialize git repository with master branch (for compatibility with tests)
			gitInit := exec.Command("git", "init", "-b", "master")
			gitInit.Dir = testLibraryPath
			if output, err := gitInit.CombinedOutput(); err != nil {
				panic(fmt.Sprintf("Failed to initialize git repository: %v\nOutput: %s", err, output))
			}

			// Configure git user
			gitConfigEmail := exec.Command("git", "config", "user.email", "test@example.com")
			gitConfigEmail.Dir = testLibraryPath
			if output, err := gitConfigEmail.CombinedOutput(); err != nil {
				panic(fmt.Sprintf("Failed to configure git user.email: %v\nOutput: %s", err, output))
			}

			gitConfigName := exec.Command("git", "config", "user.name", "Test User")
			gitConfigName.Dir = testLibraryPath
			if output, err := gitConfigName.CombinedOutput(); err != nil {
				panic(fmt.Sprintf("Failed to configure git user.name: %v\nOutput: %s", err, output))
			}

			// Create initial commit (allow empty if no files)
			gitAdd := exec.Command("git", "add", ".")
			gitAdd.Dir = testLibraryPath
			if output, err := gitAdd.CombinedOutput(); err != nil {
				panic(fmt.Sprintf("Failed to add files to git: %v\nOutput: %s", err, output))
			}

			gitCommit := exec.Command("git", "commit", "--allow-empty", "-m", "Test fixture")
			gitCommit.Dir = testLibraryPath
			if output, err := gitCommit.CombinedOutput(); err != nil {
				panic(fmt.Sprintf("Failed to commit files: %v\nOutput: %s\nPath: %s", err, output, testLibraryPath))
			}
		} else {
			// Check if there are changes and commit them
			gitStatus := exec.Command("git", "status", "--porcelain")
			gitStatus.Dir = testLibraryPath
			output, _ := gitStatus.Output()
			if len(output) > 0 {
				// There are changes, commit them
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = testLibraryPath
				_ = gitAdd.Run() // Best effort

				gitCommit := exec.Command("git", "commit", "-m", "Update test fixture")
				gitCommit.Dir = testLibraryPath
				_ = gitCommit.Run() // Best effort
			}
		}
	})
	return testLibraryPath
}

// syncDirectory copies files from src to dst, preserving directory structure
func syncDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directories from source
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Calculate relative path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = srcFile.Close() }()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer func() { _ = dstFile.Close() }()

		if _, err := dstFile.ReadFrom(srcFile); err != nil {
			return err
		}

		return os.Chmod(dstPath, info.Mode())
	})
}

// TestEnvironment provides isolated testing environment for .ddx/config.yaml
type TestEnvironment struct {
	Dir            string
	ConfigPath     string
	LibraryPath    string
	TestLibraryURL string
	GitInitialized bool

	// Optional customizations
	Platform     string
	Architecture string
	HomeDir      string

	t *testing.T
}

// TestEnvOption is a functional option for configuring TestEnvironment
type TestEnvOption func(*TestEnvironment)

// WithGitInit controls whether to initialize a git repository
func WithGitInit(init bool) TestEnvOption {
	return func(te *TestEnvironment) {
		te.GitInitialized = init
	}
}

// WithPlatform sets the platform and architecture for installation tests
func WithPlatform(platform, arch string) TestEnvOption {
	return func(te *TestEnvironment) {
		te.Platform = platform
		te.Architecture = arch
	}
}

// WithHomeDir sets a custom home directory
func WithHomeDir(homeDir string) TestEnvOption {
	return func(te *TestEnvironment) {
		te.HomeDir = homeDir
	}
}

// WithCustomLibraryURL sets a custom library repository URL
func WithCustomLibraryURL(url string) TestEnvOption {
	return func(te *TestEnvironment) {
		te.TestLibraryURL = url
	}
}

// NewTestEnvironment creates a clean test environment with temp directory
// By default: creates git repo, uses test-library fixture via file:// URL
func NewTestEnvironment(t *testing.T, opts ...TestEnvOption) *TestEnvironment {
	t.Helper()

	tempDir := t.TempDir()
	ddxDir := filepath.Join(tempDir, ".ddx")
	configPath := filepath.Join(ddxDir, "config.yaml")

	te := &TestEnvironment{
		Dir:            tempDir,
		ConfigPath:     configPath,
		LibraryPath:    filepath.Join(ddxDir, "library"),
		GitInitialized: true, // default: init git
		t:              t,
	}

	// Apply options
	for _, opt := range opts {
		opt(te)
	}

	// Set default test library URL if not customized
	if te.TestLibraryURL == "" {
		te.TestLibraryURL = "file://" + GetTestLibraryPath()
	}

	// Create .ddx directory
	require.NoError(t, os.MkdirAll(ddxDir, 0755))

	// Initialize git repository if requested
	if te.GitInitialized {
		te.initGit()
	}

	return te
}

// initGit initializes a git repository in the test directory
func (te *TestEnvironment) initGit() {
	te.t.Helper()

	// git init
	gitInit := exec.Command("git", "init")
	gitInit.Dir = te.Dir
	require.NoError(te.t, gitInit.Run(), "git init should succeed")

	// git config user.email
	gitEmail := exec.Command("git", "config", "user.email", "test@example.com")
	gitEmail.Dir = te.Dir
	require.NoError(te.t, gitEmail.Run(), "git config user.email should succeed")

	// git config user.name
	gitName := exec.Command("git", "config", "user.name", "Test User")
	gitName.Dir = te.Dir
	require.NoError(te.t, gitName.Run(), "git config user.name should succeed")
}

// CreateConfig creates a config file with the given content
func (te *TestEnvironment) CreateConfig(content string) {
	te.t.Helper()
	require.NoError(te.t, os.WriteFile(te.ConfigPath, []byte(content), 0644))
}

// CreateDefaultConfig creates a minimal valid config file using test library
func (te *TestEnvironment) CreateDefaultConfig() {
	te.t.Helper()
	content := fmt.Sprintf(`version: "1.0"
library:
  path: .ddx/library
  repository:
    url: %s
    branch: master
persona_bindings: {}
`, te.TestLibraryURL)
	te.CreateConfig(content)
}

// CreateConfigWithCustomURL creates a config with a custom repository URL
func (te *TestEnvironment) CreateConfigWithCustomURL(url string) {
	te.t.Helper()
	content := fmt.Sprintf(`version: "1.0"
library:
  path: .ddx/library
  repository:
    url: %s
    branch: master
persona_bindings: {}
`, url)
	te.CreateConfig(content)
}

// RunCommand runs a DDx command in the test environment and returns output
func (te *TestEnvironment) RunCommand(args ...string) (string, error) {
	te.t.Helper()
	factory := NewCommandFactory(te.Dir)
	cmd := factory.NewRootCommand()
	cmd.SetArgs(args)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	return buf.String(), err
}

// InitWithDDx properly initializes DDx in the test environment using ddx init
// If git is initialized, uses real git subtree with file:// URL to test library
// Otherwise uses --no-git flag
func (te *TestEnvironment) InitWithDDx(flags ...string) {
	te.t.Helper()

	// Default flags if none provided
	if len(flags) == 0 {
		if te.GitInitialized {
			// Create initial commit (required for git subtree)
			te.CreateFile("README.md", "# Test Project")
			gitAdd := exec.Command("git", "add", ".")
			gitAdd.Dir = te.Dir
			require.NoError(te.t, gitAdd.Run())
			gitCommit := exec.Command("git", "commit", "-m", "Initial commit")
			gitCommit.Dir = te.Dir
			require.NoError(te.t, gitCommit.Run())

			// Create config with test library URL first, then use --force to initialize
			te.CreateDefaultConfig()
			flags = []string{"--force", "--silent"}
		} else {
			// No git repo, use --no-git
			flags = []string{"--no-git", "--silent"}
		}
	}

	args := append([]string{"init"}, flags...)
	output, err := te.RunCommand(args...)
	require.NoError(te.t, err, "init should succeed: %s", output)
}

// LoadConfig loads the config using ConfigLoader
func (te *TestEnvironment) LoadConfig() (*config.Config, error) {
	loader, err := config.NewConfigLoaderWithWorkingDir(te.Dir)
	if err != nil {
		return nil, err
	}
	return loader.LoadConfig()
}

// CreateFile creates any file in the test environment
func (te *TestEnvironment) CreateFile(relativePath, content string) {
	te.t.Helper()
	fullPath := filepath.Join(te.Dir, relativePath)
	dir := filepath.Dir(fullPath)
	require.NoError(te.t, os.MkdirAll(dir, 0755))
	require.NoError(te.t, os.WriteFile(fullPath, []byte(content), 0644))
}

// NewTestRootCommand creates a fresh root command for tests using isolated temp directory
// This is the preferred way to create test commands - it ensures test isolation.
func NewTestRootCommand(t *testing.T) *CommandFactory {
	t.Helper()
	tempDir := t.TempDir()
	return NewCommandFactory(tempDir)
}

// NewTestRootCommandWithDir creates a test command with a specific working directory
// Use this when your test needs to operate in a specific directory.
func NewTestRootCommandWithDir(dir string) *CommandFactory {
	return NewCommandFactory(dir)
}

// executeCommand is a helper to execute commands with captured output (legacy compatibility)
// New tests should use TestEnvironment.RunCommand() instead
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}
