package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GitCredentialHelper implements git credential helper integration
type GitCredentialHelper struct {
	helperCmd string
}

// NewGitCredentialHelper creates a new git credential helper
func NewGitCredentialHelper() *GitCredentialHelper {
	return &GitCredentialHelper{
		helperCmd: "git",
	}
}

// Name returns the name of the credential helper
func (h *GitCredentialHelper) Name() string {
	return "git-credential-helper"
}

// IsAvailable checks if git credential helper is available
func (h *GitCredentialHelper) IsAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// Get retrieves credentials from git credential helper
func (h *GitCredentialHelper) Get(ctx context.Context, host string) (*Credential, error) {
	cmd := exec.CommandContext(ctx, "git", "credential", "fill")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("protocol=https\nhost=%s\n\n", host))

	output, err := cmd.Output()
	if err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: "No credentials found in git credential helper",
			Code:    "GIT_CREDENTIAL_NOT_FOUND",
		}
	}

	return h.parseCredentialOutput(string(output), host)
}

// Store stores credentials in git credential helper
func (h *GitCredentialHelper) Store(ctx context.Context, host string, cred *Credential) error {
	input := fmt.Sprintf("protocol=https\nhost=%s\nusername=%s\npassword=%s\n\n",
		host, cred.Username, cred.Token)

	cmd := exec.CommandContext(ctx, "git", "credential", "approve")
	cmd.Stdin = strings.NewReader(input)

	if err := cmd.Run(); err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to store credentials in git credential helper",
			Code:    "GIT_CREDENTIAL_STORE_ERROR",
		}
	}

	return nil
}

// Erase removes credentials from git credential helper
func (h *GitCredentialHelper) Erase(ctx context.Context, host string) error {
	cmd := exec.CommandContext(ctx, "git", "credential", "reject")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("protocol=https\nhost=%s\n\n", host))

	if err := cmd.Run(); err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to erase credentials from git credential helper",
			Code:    "GIT_CREDENTIAL_ERASE_ERROR",
		}
	}

	return nil
}

// parseCredentialOutput parses git credential helper output
func (h *GitCredentialHelper) parseCredentialOutput(output, host string) (*Credential, error) {
	lines := strings.Split(output, "\n")
	var username, password string

	for _, line := range lines {
		if strings.HasPrefix(line, "username=") {
			username = strings.TrimPrefix(line, "username=")
		} else if strings.HasPrefix(line, "password=") {
			password = strings.TrimPrefix(line, "password=")
		}
	}

	if username == "" || password == "" {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: "Incomplete credentials from git credential helper",
			Code:    "GIT_CREDENTIAL_INCOMPLETE",
		}
	}

	return &Credential{
		ID:        host,
		Platform:  h.detectPlatform(host),
		Method:    AuthMethodToken,
		Username:  username,
		Token:     password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// detectPlatform detects the platform from hostname
func (h *GitCredentialHelper) detectPlatform(host string) Platform {
	switch {
	case strings.Contains(host, "github.com"):
		return PlatformGitHub
	case strings.Contains(host, "gitlab.com"):
		return PlatformGitLab
	case strings.Contains(host, "bitbucket.org"):
		return PlatformBitbucket
	default:
		return PlatformGeneric
	}
}

// GitHubCLIHelper implements GitHub CLI credential helper
type GitHubCLIHelper struct{}

// NewGitHubCLIHelper creates a new GitHub CLI credential helper
func NewGitHubCLIHelper() *GitHubCLIHelper {
	return &GitHubCLIHelper{}
}

// Name returns the name of the credential helper
func (h *GitHubCLIHelper) Name() string {
	return "github-cli"
}

// IsAvailable checks if GitHub CLI is available
func (h *GitHubCLIHelper) IsAvailable() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// Get retrieves credentials from GitHub CLI
func (h *GitHubCLIHelper) Get(ctx context.Context, host string) (*Credential, error) {
	if !strings.Contains(host, "github.com") {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: "GitHub CLI only supports github.com",
			Code:    "GH_CLI_UNSUPPORTED_HOST",
		}
	}

	cmd := exec.CommandContext(ctx, "gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: "No GitHub token found in GitHub CLI",
			Code:    "GH_CLI_TOKEN_NOT_FOUND",
			Hint:    "Run 'gh auth login' to authenticate with GitHub CLI",
		}
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: "Empty token from GitHub CLI",
			Code:    "GH_CLI_EMPTY_TOKEN",
		}
	}

	// Get user info
	userInfo, err := h.getUserInfo(ctx)
	if err != nil {
		// Use token without user info
		return &Credential{
			ID:        host,
			Platform:  PlatformGitHub,
			Method:    AuthMethodToken,
			Token:     token,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}

	return &Credential{
		ID:       host,
		Platform: PlatformGitHub,
		Method:   AuthMethodToken,
		Username: userInfo.Login,
		Token:    token,
		Metadata: map[string]string{
			"user_id":   fmt.Sprintf("%d", userInfo.ID),
			"user_name": userInfo.Name,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Store is not supported by GitHub CLI helper (read-only)
func (h *GitHubCLIHelper) Store(ctx context.Context, host string, cred *Credential) error {
	return &AuthError{
		Type:    ErrorTypeStorageError,
		Message: "GitHub CLI credential helper is read-only",
		Code:    "GH_CLI_READ_ONLY",
		Hint:    "Use 'gh auth login' to update GitHub credentials",
	}
}

// Erase is not supported by GitHub CLI helper (read-only)
func (h *GitHubCLIHelper) Erase(ctx context.Context, host string) error {
	return &AuthError{
		Type:    ErrorTypeStorageError,
		Message: "GitHub CLI credential helper is read-only",
		Code:    "GH_CLI_READ_ONLY",
		Hint:    "Use 'gh auth logout' to remove GitHub credentials",
	}
}

// getUserInfo gets user info from GitHub CLI
func (h *GitHubCLIHelper) getUserInfo(ctx context.Context) (*GitHubUser, error) {
	cmd := exec.CommandContext(ctx, "gh", "api", "user")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var user GitHubUser
	if err := json.Unmarshal(output, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// DefaultSSHAgent implements SSH agent integration
type DefaultSSHAgent struct{}

// NewDefaultSSHAgent creates a new SSH agent interface
func NewDefaultSSHAgent() *DefaultSSHAgent {
	return &DefaultSSHAgent{}
}

// IsAvailable checks if SSH agent is running
func (a *DefaultSSHAgent) IsAvailable() bool {
	// Check for SSH_AUTH_SOCK environment variable
	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		// No SSH agent running
		return false
	}

	// Test connection to SSH agent
	cmd := exec.Command("ssh-add", "-l")
	output, err := cmd.CombinedOutput()

	// ssh-add -l returns:
	// - 0: keys found
	// - 1: no keys but agent is running
	// - 2: cannot connect to agent
	if err != nil {
		if cmd.ProcessState.ExitCode() == 1 {
			// Agent is running but no keys - still available
			return true
		}
		// Cannot connect to agent
		return false
	}

	// Check for "The agent has no identities" message
	if strings.Contains(string(output), "no identities") {
		return true // Agent is running, just no keys
	}

	return true // Agent is running and has keys
}

// ListKeys returns available SSH keys
func (a *DefaultSSHAgent) ListKeys(ctx context.Context) ([]SSHKey, error) {
	cmd := exec.CommandContext(ctx, "ssh-add", "-l")
	output, err := cmd.CombinedOutput()

	outputStr := string(output)

	// Handle the case where no keys are loaded
	if err != nil && cmd.ProcessState.ExitCode() == 1 {
		// Agent is running but no keys loaded
		return []SSHKey{}, nil
	}

	if err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: "Cannot connect to SSH agent",
			Code:    "SSH_AGENT_CONNECTION_ERROR",
			Hint:    "Ensure SSH agent is running with 'ssh-agent'",
		}
	}

	// Check for "no identities" message
	if strings.Contains(outputStr, "no identities") || strings.Contains(outputStr, "The agent has no identities") {
		return []SSHKey{}, nil
	}

	return a.parseSSHKeys(outputStr)
}

// Sign signs data with the specified key
func (a *DefaultSSHAgent) Sign(ctx context.Context, keyPath string, data []byte) ([]byte, error) {
	// SSH agent signing would typically be done through the SSH protocol
	// This is a simplified implementation
	return nil, &AuthError{
		Type:    ErrorTypeNotFound,
		Message: "SSH signing not implemented",
		Code:    "SSH_SIGN_NOT_IMPLEMENTED",
	}
}

// parseSSHKeys parses ssh-add -l output
func (a *DefaultSSHAgent) parseSSHKeys(output string) ([]SSHKey, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var keys []SSHKey

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		// Format: "size fingerprint comment (type)"
		key := SSHKey{
			Fingerprint: parts[1],
			Comment:     strings.Join(parts[2:], " "),
		}

		// Extract type from comment if present
		if strings.Contains(key.Comment, "(") && strings.Contains(key.Comment, ")") {
			start := strings.LastIndex(key.Comment, "(")
			end := strings.LastIndex(key.Comment, ")")
			if start < end {
				key.Type = key.Comment[start+1 : end]
				key.Comment = strings.TrimSpace(key.Comment[:start])
			}
		}

		keys = append(keys, key)
	}

	return keys, nil
}
