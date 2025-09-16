package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetLibraryPath returns the path to the DDx library with the following priority:
// 1. Override path (from command flag)
// 2. DDX_LIBRARY_BASE_PATH environment variable
// 3. Git repository root /library (for development)
// 4. Nearest .ddx/library/ directory (project-specific)
// 5. ~/.ddx/library/ (global fallback)
func GetLibraryPath(overridePath string) (string, error) {
	// 1. Check override path (from command flag)
	if overridePath != "" {
		absPath, err := filepath.Abs(overridePath)
		if err != nil {
			return "", fmt.Errorf("invalid override path: %w", err)
		}
		if !dirExists(absPath) {
			return "", fmt.Errorf("library path does not exist: %s", absPath)
		}
		return absPath, nil
	}

	// 2. Check environment variable
	if envPath := os.Getenv("DDX_LIBRARY_BASE_PATH"); envPath != "" {
		absPath, err := filepath.Abs(envPath)
		if err != nil {
			return "", fmt.Errorf("invalid DDX_LIBRARY_BASE_PATH: %w", err)
		}
		if !dirExists(absPath) {
			return "", fmt.Errorf("DDX_LIBRARY_BASE_PATH does not exist: %s", absPath)
		}
		return absPath, nil
	}

	// 3. Check if we're in DDx development (git repo with library/)
	if gitRoot := findGitRoot(); gitRoot != "" {
		libPath := filepath.Join(gitRoot, "library")
		if dirExists(libPath) {
			// Verify this is the DDx repository by checking for cli/main.go
			if fileExists(filepath.Join(gitRoot, "cli", "main.go")) {
				return libPath, nil
			}
		}
	}

	// 4. Find nearest .ddx/library/ (project-specific)
	if projectLib := findNearestDDxLibrary(); projectLib != "" {
		return projectLib, nil
	}

	// 5. Global fallback
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(homeDir, ".ddx", "library"), nil
}

// findGitRoot finds the root of the git repository
func findGitRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if dirExists(gitDir) {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// findNearestDDxLibrary finds the nearest .ddx/library directory
func findNearestDDxLibrary() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		ddxLib := filepath.Join(dir, ".ddx", "library")
		if dirExists(ddxLib) {
			return ddxLib
		}

		// Also check for .ddx.yml to identify project root
		ddxConfig := filepath.Join(dir, ".ddx.yml")
		if fileExists(ddxConfig) {
			// Project has .ddx.yml but no library yet
			// Return the expected path (it will be created if needed)
			return filepath.Join(dir, ".ddx", "library")
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// GetPersonasPath returns the path to the personas directory
func GetPersonasPath(libraryOverride string) (string, error) {
	libPath, err := GetLibraryPath(libraryOverride)
	if err != nil {
		return "", err
	}
	return filepath.Join(libPath, "personas"), nil
}

// GetMCPServersPath returns the path to the MCP servers directory
func GetMCPServersPath(libraryOverride string) (string, error) {
	libPath, err := GetLibraryPath(libraryOverride)
	if err != nil {
		return "", err
	}
	return filepath.Join(libPath, "mcp-servers"), nil
}

// GetTemplatesPath returns the path to the templates directory
func GetTemplatesPath(libraryOverride string) (string, error) {
	libPath, err := GetLibraryPath(libraryOverride)
	if err != nil {
		return "", err
	}
	return filepath.Join(libPath, "templates"), nil
}

// GetPatternsPath returns the path to the patterns directory
func GetPatternsPath(libraryOverride string) (string, error) {
	libPath, err := GetLibraryPath(libraryOverride)
	if err != nil {
		return "", err
	}
	return filepath.Join(libPath, "patterns"), nil
}

// GetPromptsPath returns the path to the prompts directory
func GetPromptsPath(libraryOverride string) (string, error) {
	libPath, err := GetLibraryPath(libraryOverride)
	if err != nil {
		return "", err
	}
	return filepath.Join(libPath, "prompts"), nil
}

// GetConfigsPath returns the path to the configs directory
func GetConfigsPath(libraryOverride string) (string, error) {
	libPath, err := GetLibraryPath(libraryOverride)
	if err != nil {
		return "", err
	}
	return filepath.Join(libPath, "configs"), nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil && info.IsDir()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// LibraryPathFlag is the flag name for overriding library path
const LibraryPathFlag = "library-base-path"

// ResolveLibraryResource resolves a resource path relative to the library
func ResolveLibraryResource(resourcePath string, libraryOverride string) (string, error) {
	libPath, err := GetLibraryPath(libraryOverride)
	if err != nil {
		return "", err
	}

	// Clean the resource path to prevent directory traversal
	cleanPath := filepath.Clean(resourcePath)
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("invalid resource path: %s", resourcePath)
	}

	fullPath := filepath.Join(libPath, cleanPath)

	// Verify the path exists
	if !fileExists(fullPath) && !dirExists(fullPath) {
		return "", fmt.Errorf("resource not found: %s", resourcePath)
	}

	return fullPath, nil
}

