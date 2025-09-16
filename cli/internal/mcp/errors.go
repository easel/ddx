package mcp

import "errors"

// Common errors
var (
	// Registry errors
	ErrInvalidPath     = errors.New("invalid path")
	ErrEmptyServerName = errors.New("empty server name")
	ErrServerNotFound  = errors.New("server not found")
	ErrMissingVersion  = errors.New("missing registry version")
	ErrMissingName     = errors.New("missing server name")
	ErrDuplicateServer = errors.New("duplicate server name")

	// Installation errors
	ErrAlreadyInstalled = errors.New("server already installed")
	ErrClaudeNotFound   = errors.New("Claude installation not found")
	ErrInvalidConfig    = errors.New("invalid configuration")
	ErrConfigCorrupted  = errors.New("configuration file corrupted")
	ErrPermissionDenied = errors.New("permission denied")
	ErrBackupFailed     = errors.New("backup creation failed")

	// Environment errors
	ErrMissingRequired  = errors.New("missing required environment variable")
	ErrInvalidValue     = errors.New("invalid environment variable value")
	ErrValidationFailed = errors.New("validation failed")

	// Security errors
	ErrPathTraversal    = errors.New("path traversal detected")
	ErrInjectionAttempt = errors.New("injection attempt detected")
	ErrInsecureValue    = errors.New("insecure value detected")

	// Network errors
	ErrRegistryUnreachable = errors.New("registry unreachable")
	ErrDownloadFailed      = errors.New("download failed")
	ErrChecksumMismatch    = errors.New("checksum mismatch")
)

// MCPError provides detailed error information
type MCPError struct {
	Code       string                 // Error code for programmatic handling
	Message    string                 // User-friendly error message
	Details    map[string]interface{} // Additional error details
	Resolution string                 // Suggested resolution
	Wrapped    error                  // Underlying error
}

// Error implements the error interface
func (e *MCPError) Error() string {
	if e.Resolution != "" {
		return e.Message + ". " + e.Resolution
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *MCPError) Unwrap() error {
	return e.Wrapped
}

// NewMCPError creates a new MCP error
func NewMCPError(code, message, resolution string, wrapped error) *MCPError {
	return &MCPError{
		Code:       code,
		Message:    message,
		Resolution: resolution,
		Wrapped:    wrapped,
		Details:    make(map[string]interface{}),
	}
}

// Common error codes
const (
	ErrCodeServerNotFound   = "MCP_SERVER_NOT_FOUND"
	ErrCodeAlreadyInstalled = "MCP_ALREADY_INSTALLED"
	ErrCodeClaudeNotFound   = "MCP_CLAUDE_NOT_FOUND"
	ErrCodeInvalidEnv       = "MCP_INVALID_ENV"
	ErrCodeConfigCorrupt    = "MCP_CONFIG_CORRUPT"
	ErrCodePermissionDenied = "MCP_PERMISSION_DENIED"
	ErrCodeNetworkError     = "MCP_NETWORK_ERROR"
	ErrCodeValidationError  = "MCP_VALIDATION_ERROR"
	ErrCodeSecurityError    = "MCP_SECURITY_ERROR"
)
