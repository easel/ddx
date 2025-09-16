package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// StandardLocations defines where to look for Claude configurations
var StandardLocations = []struct {
	Type       ClaudeType
	ConfigPath string
	Platform   string
	Priority   int
}{
	// Claude Code
	{ClaudeCode, "~/.claude/settings.local.json", "darwin", 1},
	{ClaudeCode, "~/.claude/settings.local.json", "linux", 1},
	{ClaudeCode, "%USERPROFILE%\\.claude\\settings.local.json", "windows", 1},

	// Claude Desktop
	{ClaudeDesktop, "~/Library/Application Support/Claude/claude_desktop_config.json", "darwin", 2},
	{ClaudeDesktop, "~/.config/Claude/claude_desktop_config.json", "linux", 2},
	{ClaudeDesktop, "%APPDATA%\\Claude\\claude_desktop_config.json", "windows", 2},
}

// DetectClaude detects all Claude installations
func DetectClaude() ([]ClaudeInstallation, error) {
	var installations []ClaudeInstallation

	// Check environment variables first
	if path := os.Getenv("CLAUDE_CODE_CONFIG"); path != "" {
		if installation := detectFromPath(path, ClaudeCode); installation != nil {
			installations = append(installations, *installation)
		}
	}

	if path := os.Getenv("CLAUDE_DESKTOP_CONFIG"); path != "" {
		if installation := detectFromPath(path, ClaudeDesktop); installation != nil {
			installations = append(installations, *installation)
		}
	}

	// Check standard locations
	for _, loc := range StandardLocations {
		if runtime.GOOS != loc.Platform {
			continue
		}

		path := expandPath(loc.ConfigPath)
		if installation := detectFromPath(path, loc.Type); installation != nil {
			// Check if we already have this installation
			duplicate := false
			for _, existing := range installations {
				if existing.ConfigPath == installation.ConfigPath {
					duplicate = true
					break
				}
			}
			if !duplicate {
				installations = append(installations, *installation)
			}
		}
	}

	return installations, nil
}

// DetectClaudeCode detects Claude Code installation
func DetectClaudeCode() (*ClaudeInstallation, error) {
	// Check environment variable
	if path := os.Getenv("CLAUDE_CODE_CONFIG"); path != "" {
		return detectFromPath(path, ClaudeCode), nil
	}

	// Check standard locations
	for _, loc := range StandardLocations {
		if loc.Type != ClaudeCode || runtime.GOOS != loc.Platform {
			continue
		}

		path := expandPath(loc.ConfigPath)
		if installation := detectFromPath(path, ClaudeCode); installation != nil {
			return installation, nil
		}
	}

	return nil, ErrClaudeNotFound
}

// DetectClaudeDesktop detects Claude Desktop installation
func DetectClaudeDesktop() (*ClaudeInstallation, error) {
	// Check environment variable
	if path := os.Getenv("CLAUDE_DESKTOP_CONFIG"); path != "" {
		return detectFromPath(path, ClaudeDesktop), nil
	}

	// Check standard locations
	for _, loc := range StandardLocations {
		if loc.Type != ClaudeDesktop || runtime.GOOS != loc.Platform {
			continue
		}

		path := expandPath(loc.ConfigPath)
		if installation := detectFromPath(path, ClaudeDesktop); installation != nil {
			return installation, nil
		}
	}

	return nil, ErrClaudeNotFound
}

// detectFromPath checks if a Claude installation exists at the given path
func detectFromPath(path string, claudeType ClaudeType) *ClaudeInstallation {
	// Check if file exists or parent directory exists
	if _, err := os.Stat(path); err != nil {
		// Check if parent directory exists (file might not exist yet)
		parentDir := filepath.Dir(path)
		if _, err := os.Stat(parentDir); err != nil {
			return nil
		}
	}

	installation := &ClaudeInstallation{
		Type:       claudeType,
		ConfigPath: path,
		Detected:   time.Now(),
	}

	// Try to detect version from config
	if version := detectVersion(path); version != "" {
		installation.Version = version
	}

	return installation
}

// detectVersion attempts to detect Claude version from config
func detectVersion(configPath string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "unknown"
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return "unknown"
	}

	// Check for version field
	if version, ok := config["version"].(string); ok {
		return version
	}

	// Check for mcpServers to determine if it's a newer version
	if _, ok := config["mcpServers"]; ok {
		if strings.Contains(configPath, "desktop") {
			return "1.0+"
		}
		return "0.2+"
	}

	return "unknown"
}

// expandPath expands ~ and environment variables in a path
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	return os.ExpandEnv(path)
}

// GetConfigPath returns the default config path for a Claude type
func GetConfigPath(claudeType ClaudeType) string {
	for _, loc := range StandardLocations {
		if loc.Type == claudeType && runtime.GOOS == loc.Platform {
			return expandPath(loc.ConfigPath)
		}
	}
	return ""
}

// SelectInstallation prompts the user to select from multiple installations
func SelectInstallation(installations []ClaudeInstallation) (*ClaudeInstallation, error) {
	if len(installations) == 0 {
		return nil, ErrClaudeNotFound
	}

	if len(installations) == 1 {
		return &installations[0], nil
	}

	// In non-interactive mode, prefer Claude Code over Desktop
	for _, installation := range installations {
		if installation.Type == ClaudeCode {
			return &installation, nil
		}
	}

	return &installations[0], nil
}

// FormatInstallation formats an installation for display
func FormatInstallation(installation ClaudeInstallation) string {
	typeStr := "Claude Code"
	if installation.Type == ClaudeDesktop {
		typeStr = "Claude Desktop"
	}

	return fmt.Sprintf("%s at %s (version: %s)", typeStr, installation.ConfigPath, installation.Version)
}
