// Package auth provides authentication management for DDx operations
package auth

import (
	"context"
	"time"
)

// AuthMethod represents different authentication methods
type AuthMethod string

const (
	AuthMethodHTTPS AuthMethod = "https"
	AuthMethodOAuth AuthMethod = "oauth"
	AuthMethodToken AuthMethod = "token"
)

// Platform represents different Git hosting platforms
type Platform string

const (
	PlatformGitHub    Platform = "github"
	PlatformGitLab    Platform = "gitlab"
	PlatformBitbucket Platform = "bitbucket"
	PlatformGeneric   Platform = "generic"
)

// Credential represents stored authentication credentials
type Credential struct {
	ID        string            `json:"id"`
	Platform  Platform          `json:"platform"`
	Method    AuthMethod        `json:"method"`
	Username  string            `json:"username,omitempty"`
	Token     string            `json:"token,omitempty"`
	ExpiresAt *time.Time        `json:"expires_at,omitempty"`
	Scopes    []string          `json:"scopes,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// AuthRequest represents an authentication request
type AuthRequest struct {
	Platform    Platform   `json:"platform"`
	Repository  string     `json:"repository"`
	Method      AuthMethod `json:"method,omitempty"`
	Scopes      []string   `json:"scopes,omitempty"`
	Interactive bool       `json:"interactive"`
}

// AuthResult represents the result of an authentication attempt
type AuthResult struct {
	Success    bool        `json:"success"`
	Credential *Credential `json:"credential,omitempty"`
	Method     AuthMethod  `json:"method"`
	Message    string      `json:"message,omitempty"`
	Error      error       `json:"error,omitempty"`
}

// TwoFactorChallenge represents a 2FA challenge
type TwoFactorChallenge struct {
	Type        string `json:"type"` // "totp", "sms", "app"
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// TwoFactorResponse represents a 2FA response
type TwoFactorResponse struct {
	Code   string `json:"code"`
	Method string `json:"method"`
}

// Manager is the main authentication manager interface
type Manager interface {
	// Authenticate performs authentication for a given request
	Authenticate(ctx context.Context, req *AuthRequest) (*AuthResult, error)

	// ValidateCredentials validates existing credentials
	ValidateCredentials(ctx context.Context, platform Platform, repository string) error

	// GetCredential retrieves stored credentials for a platform/repository
	GetCredential(ctx context.Context, platform Platform, repository string) (*Credential, error)

	// StoreCredential securely stores authentication credentials
	StoreCredential(ctx context.Context, cred *Credential) error

	// DeleteCredential removes stored credentials
	DeleteCredential(ctx context.Context, platform Platform, repository string) error

	// ListCredentials lists all stored credentials
	ListCredentials(ctx context.Context) ([]*Credential, error)

	// RefreshCredential refreshes expired credentials
	RefreshCredential(ctx context.Context, platform Platform, repository string) (*Credential, error)
}

// Authenticator is the interface for platform-specific authenticators
type Authenticator interface {
	// Platform returns the platform this authenticator supports
	Platform() Platform

	// SupportedMethods returns the authentication methods supported
	SupportedMethods() []AuthMethod

	// Authenticate performs platform-specific authentication
	Authenticate(ctx context.Context, req *AuthRequest) (*AuthResult, error)

	// ValidateToken validates a token's format and scopes
	ValidateToken(ctx context.Context, token string, requiredScopes []string) error

	// RefreshToken refreshes an expired token
	RefreshToken(ctx context.Context, refreshToken string) (*Credential, error)

	// HandleTwoFactor handles 2FA challenges
	HandleTwoFactor(ctx context.Context, challenge *TwoFactorChallenge) (*TwoFactorResponse, error)
}

// Store is the interface for credential storage backends
type Store interface {
	// Get retrieves a credential by platform and repository
	Get(ctx context.Context, platform Platform, repository string) (*Credential, error)

	// Set stores a credential securely
	Set(ctx context.Context, cred *Credential) error

	// Delete removes a stored credential
	Delete(ctx context.Context, platform Platform, repository string) error

	// List returns all stored credentials
	List(ctx context.Context) ([]*Credential, error)

	// Clear removes all stored credentials
	Clear(ctx context.Context) error

	// IsAvailable checks if the store is available
	IsAvailable() bool
}

// CredentialHelper is the interface for system credential helpers
type CredentialHelper interface {
	// Name returns the name of the credential helper
	Name() string

	// IsAvailable checks if the helper is available on the system
	IsAvailable() bool

	// Get retrieves credentials from the helper
	Get(ctx context.Context, host string) (*Credential, error)

	// Store stores credentials in the helper
	Store(ctx context.Context, host string, cred *Credential) error

	// Erase removes credentials from the helper
	Erase(ctx context.Context, host string) error
}

// ValidationError represents credential validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// AuthError represents authentication errors
type AuthError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Hint    string `json:"hint,omitempty"`
}

func (e *AuthError) Error() string {
	return e.Message
}

// Common error types
const (
	ErrorTypeInvalidCredentials = "invalid_credentials"
	ErrorTypeExpiredToken       = "expired_token"
	ErrorTypeInsufficientScope  = "insufficient_scope"
	ErrorTypeNetworkError       = "network_error"
	ErrorTypeTwoFactorRequired  = "two_factor_required"
	ErrorTypeRateLimited        = "rate_limited"
	ErrorTypeNotFound           = "not_found"
	ErrorTypeStorageError       = "storage_error"
)
