package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// GitHubAuthenticator implements authentication for GitHub
type GitHubAuthenticator struct {
	baseURL string
}

// NewGitHubAuthenticator creates a new GitHub authenticator
func NewGitHubAuthenticator() *GitHubAuthenticator {
	return &GitHubAuthenticator{
		baseURL: "https://api.github.com",
	}
}

// Platform returns the platform this authenticator supports
func (a *GitHubAuthenticator) Platform() Platform {
	return PlatformGitHub
}

// SupportedMethods returns the authentication methods supported
func (a *GitHubAuthenticator) SupportedMethods() []AuthMethod {
	return []AuthMethod{
		AuthMethodToken,
		AuthMethodOAuth,
	}
}

// Authenticate performs GitHub-specific authentication
func (a *GitHubAuthenticator) Authenticate(ctx context.Context, req *AuthRequest) (*AuthResult, error) {
	switch req.Method {
	case AuthMethodToken:
		return a.authenticateWithToken(ctx, req)
	case AuthMethodOAuth:
		return a.authenticateWithOAuth(ctx, req)
	default:
		// Try token authentication as default
		return a.authenticateWithToken(ctx, req)
	}
}

// ValidateToken validates a GitHub token's format and scopes
func (a *GitHubAuthenticator) ValidateToken(ctx context.Context, token string, requiredScopes []string) error {
	if !a.isValidTokenFormat(token) {
		return &ValidationError{
			Field:   "token",
			Message: "Invalid GitHub token format",
			Code:    "GITHUB_INVALID_TOKEN_FORMAT",
		}
	}

	// Check token validity by calling GitHub API
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", a.baseURL+"/user", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return &AuthError{
			Type:    ErrorTypeNetworkError,
			Message: "Failed to validate GitHub token: " + err.Error(),
			Code:    "GITHUB_NETWORK_ERROR",
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return &AuthError{
			Type:    ErrorTypeInvalidCredentials,
			Message: "GitHub token is invalid or expired",
			Code:    "GITHUB_INVALID_TOKEN",
			Hint:    "Generate a new personal access token at https://github.com/settings/tokens",
		}
	}

	if resp.StatusCode != 200 {
		return &AuthError{
			Type:    ErrorTypeInvalidCredentials,
			Message: fmt.Sprintf("GitHub API returned status %d", resp.StatusCode),
			Code:    "GITHUB_API_ERROR",
		}
	}

	// Check scopes if required
	if len(requiredScopes) > 0 {
		scopes := resp.Header.Get("X-OAuth-Scopes")
		if err := a.validateScopes(scopes, requiredScopes); err != nil {
			return err
		}
	}

	return nil
}

// RefreshToken refreshes an expired GitHub token
func (a *GitHubAuthenticator) RefreshToken(ctx context.Context, refreshToken string) (*Credential, error) {
	// GitHub personal access tokens don't have refresh capability
	// Return error to trigger re-authentication
	return nil, &AuthError{
		Type:    ErrorTypeExpiredToken,
		Message: "GitHub personal access tokens cannot be refreshed",
		Code:    "GITHUB_NO_REFRESH",
		Hint:    "Generate a new personal access token",
	}
}

// HandleTwoFactor handles GitHub 2FA challenges
func (a *GitHubAuthenticator) HandleTwoFactor(ctx context.Context, challenge *TwoFactorChallenge) (*TwoFactorResponse, error) {
	// GitHub 2FA is typically handled via personal access tokens
	// which bypass 2FA requirements
	return nil, &AuthError{
		Type:    ErrorTypeTwoFactorRequired,
		Message: "Two-factor authentication required",
		Code:    "GITHUB_2FA_REQUIRED",
		Hint:    "Use a personal access token which bypasses 2FA",
	}
}

// authenticateWithToken performs token-based authentication
func (a *GitHubAuthenticator) authenticateWithToken(ctx context.Context, req *AuthRequest) (*AuthResult, error) {
	if !req.Interactive {
		return &AuthResult{
			Success: false,
			Message: "Token authentication requires interactive mode",
		}, nil
	}

	// In a real implementation, this would prompt the user for a token
	// For now, return error to indicate manual token setup is needed
	return &AuthResult{
		Success: false,
		Message: "Please configure GitHub personal access token manually",
		Error: &AuthError{
			Type:    ErrorTypeInvalidCredentials,
			Message: "GitHub personal access token required",
			Code:    "GITHUB_TOKEN_REQUIRED",
			Hint:    "Generate token at https://github.com/settings/tokens with scopes: " + strings.Join(req.Scopes, ", "),
		},
	}, nil
}

// authenticateWithOAuth performs OAuth authentication
func (a *GitHubAuthenticator) authenticateWithOAuth(ctx context.Context, req *AuthRequest) (*AuthResult, error) {
	// OAuth flow implementation would go here
	return &AuthResult{
		Success: false,
		Message: "OAuth authentication not yet implemented",
		Error: &AuthError{
			Type:    ErrorTypeNotFound,
			Message: "OAuth authentication not implemented",
			Code:    "GITHUB_OAUTH_NOT_IMPLEMENTED",
		},
	}, nil
}

// isValidTokenFormat checks if the token has a valid GitHub format
func (a *GitHubAuthenticator) isValidTokenFormat(token string) bool {
	// GitHub personal access tokens start with "ghp_"
	// GitHub app tokens start with "ghs_"
	// GitHub OAuth tokens start with "gho_"
	patterns := []string{
		`^ghp_[A-Za-z0-9]{36}$`,         // Personal access token
		`^ghs_[A-Za-z0-9]{36}$`,         // App token
		`^gho_[A-Za-z0-9]{36}$`,         // OAuth token
		`^github_pat_[A-Za-z0-9_]{82}$`, // Fine-grained personal access token
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, token); matched {
			return true
		}
	}

	return false
}

// validateScopes checks if the token has the required scopes
func (a *GitHubAuthenticator) validateScopes(tokenScopes string, requiredScopes []string) error {
	if len(requiredScopes) == 0 {
		return nil
	}

	availableScopes := strings.Split(tokenScopes, ", ")
	scopeMap := make(map[string]bool)
	for _, scope := range availableScopes {
		scopeMap[strings.TrimSpace(scope)] = true
	}

	var missingScopes []string
	for _, required := range requiredScopes {
		if !scopeMap[required] {
			missingScopes = append(missingScopes, required)
		}
	}

	if len(missingScopes) > 0 {
		return &AuthError{
			Type:    ErrorTypeInsufficientScope,
			Message: fmt.Sprintf("Token missing required scopes: %s", strings.Join(missingScopes, ", ")),
			Code:    "GITHUB_INSUFFICIENT_SCOPE",
			Hint:    "Update your personal access token to include the required scopes",
		}
	}

	return nil
}

// GitHubUser represents a GitHub user from the API
type GitHubUser struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// getUserInfo retrieves user information from GitHub API
func (a *GitHubAuthenticator) getUserInfo(ctx context.Context, token string) (*GitHubUser, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", a.baseURL+"/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
