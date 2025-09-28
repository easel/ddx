package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Simulate the test setup
	testDir := "/tmp/debug_list"
	os.RemoveAll(testDir)
	os.MkdirAll(testDir, 0755)
	
	// Create the same structure as the test
	libraryDir := filepath.Join(testDir, "library")
	promptsDir := filepath.Join(libraryDir, "prompts")
	os.MkdirAll(filepath.Join(promptsDir, "claude"), 0755)
	os.WriteFile(filepath.Join(promptsDir, "claude", "prompt.md"), []byte("# Prompt"), 0644)
	
	// Check what paths exist
	fmt.Printf("testDir: %s\n", testDir)
	fmt.Printf("libraryDir: %s\n", libraryDir)
	fmt.Printf("promptsDir: %s\n", promptsDir)
	
	// Check if paths exist
	if _, err := os.Stat(libraryDir); err == nil {
		fmt.Printf("✓ libraryDir exists\n")
	} else {
		fmt.Printf("✗ libraryDir missing: %v\n", err)
	}
	
	if _, err := os.Stat(promptsDir); err == nil {
		fmt.Printf("✓ promptsDir exists\n")
	} else {
		fmt.Printf("✗ promptsDir missing: %v\n", err)
	}
	
	// List contents
	if entries, err := os.ReadDir(promptsDir); err == nil {
		fmt.Printf("promptsDir contents: %v\n", len(entries))
		for _, entry := range entries {
			fmt.Printf("  - %s\n", entry.Name())
		}
	}
	
	// Test relative vs absolute path resolution
	relPath := "./library"
	fmt.Printf("Relative path from %s: %s\n", testDir, relPath)
	
	// Change to testDir and check relative path
	oldDir, _ := os.Getwd()
	os.Chdir(testDir)
	if _, err := os.Stat(relPath); err == nil {
		fmt.Printf("✓ relative path works from testDir\n")
	} else {
		fmt.Printf("✗ relative path fails from testDir: %v\n", err)
	}
	os.Chdir(oldDir)
}
