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
	"github.com/stretchr/testify/require"
)

var (
	testLibraryPath string
	testLibraryOnce sync.Once
)

// GetTestLibraryPath returns the absolute path to the test library fixture
// This is a real git repository that tests can use with file:// URLs
func GetTestLibraryPath() string {
	testLibraryOnce.Do(func() {
		// Find the test fixtures directory relative to this file
		_, filename, _, _ := runtime.Caller(0)
		cmdDir := filepath.Dir(filename)
		testLibraryPath = filepath.Join(cmdDir, "..", "test", "fixtures", "ddx-library")

		// Ensure it's an absolute path
		testLibraryPath, _ = filepath.Abs(testLibraryPath)

		// Verify it's a git repository
		gitDir := filepath.Join(testLibraryPath, ".git")
		if _, err := os.Stat(gitDir); os.IsNotExist(err) {
			panic(fmt.Sprintf("Test library fixture is not a git repository: %s", testLibraryPath))
		}
	})
	return testLibraryPath
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
// Note: Always uses --no-git for now until DDX_TEST_MODE is removed from production code
func (te *TestEnvironment) InitWithDDx(flags ...string) {
	te.t.Helper()

	// Default flags if none provided
	if len(flags) == 0 {
		// Always use --no-git until DDX_TEST_MODE is removed from production code
		// After Phase 2, we can use real git with file:// URLs
		flags = []string{"--no-git", "--silent"}
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
