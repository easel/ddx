package cmd

import (
	"os"
	"testing"

	"github.com/easel/ddx/internal/config"
)

// Test that demonstrates the functional design principles
func TestInitProject_FunctionalDesign(t *testing.T) {
	// Setup: Create a temporary directory
	tempDir, err := os.MkdirTemp("", "ddx-init-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set test mode to avoid interactive prompts and network calls
	os.Setenv("DDX_TEST_MODE", "1")
	defer os.Unsetenv("DDX_TEST_MODE")

	// Test 1: Test with git validation (should fail in non-git directory)
	opts := InitOptions{
		Force:    false,
		NoGit:    false, // Enable git validation for test
		Template: "",
	}

	result, err := initProject(tempDir, opts)
	// Should fail because tempDir is not a git repository
	if err == nil {
		t.Error("Expected error for non-git directory, but got none")
	}
	if result != nil {
		t.Error("Expected nil result on error")
	}

	// Test 2: Test with NoGit option (should succeed)
	opts.NoGit = true
	result, err = initProject(tempDir, opts)

	// Should work now that we skip git validation
	if err != nil {
		t.Errorf("Expected success with NoGit=true, but got error: %v", err)
	}
	if result == nil {
		t.Error("Expected result struct, got nil")
	}
}

// Test the pure git validation function
func TestValidateGitRepo_Pure(t *testing.T) {
	// Test with non-git directory
	tempDir, err := os.MkdirTemp("", "ddx-git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Should return error for non-git directory
	err = validateGitRepo(tempDir)
	if err == nil {
		t.Error("Expected error for non-git directory, but got none")
	}
}

// Test InitOptions struct
func TestInitOptions_Structure(t *testing.T) {
	opts := InitOptions{
		Force:    true,
		NoGit:    false,
		Template: "nextjs",
	}

	// Verify all fields are accessible
	if !opts.Force {
		t.Error("Force field not set correctly")
	}
	if opts.NoGit {
		t.Error("NoGit field not set correctly")
	}
	if opts.Template != "nextjs" {
		t.Error("Template field not set correctly")
	}
}

// Test InitResult struct
func TestInitResult_Structure(t *testing.T) {
	result := &InitResult{
		ConfigCreated: true,
		BackupPath:    "/tmp/backup",
		LibraryExists: false,
		IsDDxRepo:     true,
	}

	// Verify all fields are accessible
	if !result.ConfigCreated {
		t.Error("ConfigCreated field not set correctly")
	}
	if result.BackupPath != "/tmp/backup" {
		t.Error("BackupPath field not set correctly")
	}
	if result.LibraryExists {
		t.Error("LibraryExists field not set correctly")
	}
	if !result.IsDDxRepo {
		t.Error("IsDDxRepo field not set correctly")
	}
}

// Test that business logic functions don't depend on cobra.Command
func TestBusinessLogicIndependence(t *testing.T) {
	// These functions should work without cobra.Command

	// Test validateGitRepo - pure function
	tempDir, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(tempDir)

	err := validateGitRepo(tempDir)
	if err == nil {
		t.Error("validateGitRepo should work without cobra.Command")
	}

	// Test initializeSynchronizationPure - pure function
	cfg := &config.Config{
		Repository: &config.NewRepositoryConfig{
			URL:    "",
			Branch: "",
		},
	}
	err = initializeSynchronizationPure(cfg)
	if err == nil {
		t.Error("initializeSynchronizationPure should return error for empty URL")
	}
}
