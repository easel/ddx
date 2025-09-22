package mcp

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ClaudeWrapper wraps the Claude CLI for MCP server management
type ClaudeWrapper struct {
	executable string
}

// NewClaudeWrapper creates a new Claude CLI wrapper
func NewClaudeWrapper() *ClaudeWrapper {
	return &ClaudeWrapper{
		executable: "claude",
	}
}

// IsAvailable checks if Claude CLI is available
func (c *ClaudeWrapper) IsAvailable() error {
	_, err := exec.LookPath(c.executable)
	if err != nil {
		return fmt.Errorf("claude CLI not found in PATH: %w", err)
	}
	return nil
}

// ListServers lists installed MCP servers via claude mcp list
func (c *ClaudeWrapper) ListServers() (map[string]bool, error) {
	if err := c.IsAvailable(); err != nil {
		return nil, err
	}

	cmd := exec.Command(c.executable, "mcp", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list MCP servers: %w", err)
	}

	servers := make(map[string]bool)

	// Parse output to extract server names
	// Format: "servername: /path/to/server args - ✓ Connected"
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Checking") {
			continue
		}

		// Extract server name before ":"
		if idx := strings.Index(line, ":"); idx > 0 {
			serverName := strings.TrimSpace(line[:idx])
			// Check if it's connected (has ✓)
			connected := strings.Contains(line, "✓")
			servers[serverName] = connected
		}
	}

	return servers, nil
}

// AddServer adds an MCP server via claude mcp add
func (c *ClaudeWrapper) AddServer(name string, command string, args []string, env map[string]string) error {
	if err := c.IsAvailable(); err != nil {
		return err
	}

	if len(env) > 0 {
		// Use claude mcp add-json for servers with environment variables
		return c.addServerWithJSON(name, command, args, env)
	}

	// Use claude mcp add for simple servers
	cmdArgs := []string{"mcp", "add", name, command, "--"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(c.executable, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add MCP server %s: %w\nOutput: %s", name, err, string(output))
	}

	return nil
}

// addServerWithJSON adds an MCP server using JSON configuration
func (c *ClaudeWrapper) addServerWithJSON(name string, command string, args []string, env map[string]string) error {
	// Build JSON configuration
	config := map[string]interface{}{
		"type":    "stdio",
		"command": command,
		"args":    args,
		"env":     env,
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal server config: %w", err)
	}

	cmd := exec.Command(c.executable, "mcp", "add-json", name, string(jsonData))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add MCP server %s with JSON: %w\nOutput: %s", name, err, string(output))
	}

	return nil
}

// RemoveServer removes an MCP server via claude mcp remove
func (c *ClaudeWrapper) RemoveServer(name string) error {
	if err := c.IsAvailable(); err != nil {
		return err
	}

	cmd := exec.Command(c.executable, "mcp", "remove", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove MCP server %s: %w\nOutput: %s", name, err, string(output))
	}

	return nil
}

// GetServerStatus gets the status of a specific server
func (c *ClaudeWrapper) GetServerStatus(name string) (*ServerStatus, error) {
	servers, err := c.ListServers()
	if err != nil {
		return nil, err
	}

	connected, installed := servers[name]
	return &ServerStatus{
		Name:      name,
		Installed: installed,
		Running:   connected,
		Version:   "unknown",
		Errors:    []string{},
	}, nil
}

// CheckClaude checks if Claude CLI is available and working
func CheckClaude() error {
	wrapper := NewClaudeWrapper()
	return wrapper.IsAvailable()
}

// DetectClaude detects available Claude installations
func DetectClaude() ([]ClaudeInstallation, error) {
	installations := []ClaudeInstallation{}

	// Try to detect Claude Code installation
	if err := CheckClaude(); err == nil {
		installations = append(installations, ClaudeInstallation{
			Type:       ClaudeCode,
			ConfigPath: "~/.claude/config.json", // placeholder
			Version:    "unknown",
			Detected:   time.Now(),
		})
	}

	return installations, nil
}
