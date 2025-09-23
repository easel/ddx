package auth

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// DefaultManager implements the Manager interface
type DefaultManager struct {
	authenticators map[Platform]Authenticator
	stores         []Store
	helpers        []CredentialHelper
	sshAgent       SSHAgent
	mu             sync.RWMutex
}

// NewDefaultManager creates a new authentication manager
func NewDefaultManager() *DefaultManager {
	return &DefaultManager{
		authenticators: make(map[Platform]Authenticator),
		stores:         make([]Store, 0),
		helpers:        make([]CredentialHelper, 0),
	}
}

// RegisterAuthenticator registers a platform-specific authenticator
func (m *DefaultManager) RegisterAuthenticator(auth Authenticator) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.authenticators[auth.Platform()] = auth
}

// RegisterStore registers a credential storage backend
func (m *DefaultManager) RegisterStore(store Store) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if store.IsAvailable() {
		m.stores = append(m.stores, store)
		// Sort stores by priority (keychain first, then file storage)
		sort.Slice(m.stores, func(i, j int) bool {
			return m.getStorePriority(m.stores[i]) < m.getStorePriority(m.stores[j])
		})
	}
}

// RegisterCredentialHelper registers a system credential helper
func (m *DefaultManager) RegisterCredentialHelper(helper CredentialHelper) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if helper.IsAvailable() {
		m.helpers = append(m.helpers, helper)
	}
}

// SetSSHAgent sets the SSH agent interface
func (m *DefaultManager) SetSSHAgent(agent SSHAgent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sshAgent = agent
}

// Authenticate performs authentication for a given request
func (m *DefaultManager) Authenticate(ctx context.Context, req *AuthRequest) (*AuthResult, error) {
	m.mu.RLock()
	authenticator, exists := m.authenticators[req.Platform]
	m.mu.RUnlock()

	if !exists {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("No authenticator available for platform: %s", req.Platform),
			Code:    "AUTH_PLATFORM_NOT_SUPPORTED",
			Hint:    "Check if the platform is supported and properly configured",
		}
	}

	// Try existing credentials first
	if existingCred, err := m.GetCredential(ctx, req.Platform, req.Repository); err == nil {
		if err := m.ValidateCredentials(ctx, req.Platform, req.Repository); err == nil {
			return &AuthResult{
				Success:    true,
				Credential: existingCred,
				Method:     existingCred.Method,
				Message:    "Using existing valid credentials",
			}, nil
		}
	}

	// Try credential helpers
	for _, helper := range m.helpers {
		if cred, err := helper.Get(ctx, req.Repository); err == nil {
			if err := authenticator.ValidateToken(ctx, cred.Token, req.Scopes); err == nil {
				// Store in primary store for future use
				if err := m.StoreCredential(ctx, cred); err != nil {
					// Log warning but continue
				}
				return &AuthResult{
					Success:    true,
					Credential: cred,
					Method:     cred.Method,
					Message:    fmt.Sprintf("Authenticated using %s credential helper", helper.Name()),
				}, nil
			}
		}
	}

	// Perform new authentication
	result, err := authenticator.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	if result.Success && result.Credential != nil {
		// Store the new credential
		if err := m.StoreCredential(ctx, result.Credential); err != nil {
			// Log warning but don't fail authentication
		}
	}

	return result, nil
}

// ValidateCredentials validates existing credentials
func (m *DefaultManager) ValidateCredentials(ctx context.Context, platform Platform, repository string) error {
	cred, err := m.GetCredential(ctx, platform, repository)
	if err != nil {
		return err
	}

	// Check expiration
	if cred.ExpiresAt != nil && time.Now().After(*cred.ExpiresAt) {
		return &AuthError{
			Type:    ErrorTypeExpiredToken,
			Message: "Stored credentials have expired",
			Code:    "AUTH_CREDENTIALS_EXPIRED",
			Hint:    "Run authentication again to refresh credentials",
		}
	}

	m.mu.RLock()
	authenticator, exists := m.authenticators[platform]
	m.mu.RUnlock()

	if !exists {
		return &AuthError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("No authenticator available for platform: %s", platform),
			Code:    "AUTH_PLATFORM_NOT_SUPPORTED",
		}
	}

	// Platform-specific validation
	return authenticator.ValidateToken(ctx, cred.Token, cred.Scopes)
}

// GetCredential retrieves stored credentials for a platform/repository
func (m *DefaultManager) GetCredential(ctx context.Context, platform Platform, repository string) (*Credential, error) {
	m.mu.RLock()
	stores := make([]Store, len(m.stores))
	copy(stores, m.stores)
	m.mu.RUnlock()

	for _, store := range stores {
		if cred, err := store.Get(ctx, platform, repository); err == nil {
			return cred, nil
		}
	}

	return nil, &AuthError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("No credentials found for %s/%s", platform, repository),
		Code:    "AUTH_CREDENTIALS_NOT_FOUND",
		Hint:    "Run authentication to store credentials for this repository",
	}
}

// StoreCredential securely stores authentication credentials
func (m *DefaultManager) StoreCredential(ctx context.Context, cred *Credential) error {
	m.mu.RLock()
	stores := make([]Store, len(m.stores))
	copy(stores, m.stores)
	m.mu.RUnlock()

	if len(stores) == 0 {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "No credential storage backend available",
			Code:    "AUTH_NO_STORAGE",
			Hint:    "Check if OS keychain or file storage is available",
		}
	}

	// Update timestamps
	now := time.Now()
	if cred.CreatedAt.IsZero() {
		cred.CreatedAt = now
	}
	cred.UpdatedAt = now

	// Try to store in the highest priority store
	primaryStore := stores[0]
	if err := primaryStore.Set(ctx, cred); err != nil {
		// Try fallback stores
		for _, store := range stores[1:] {
			if err := store.Set(ctx, cred); err == nil {
				return nil
			}
		}
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to store credentials in any available backend",
			Code:    "AUTH_STORAGE_FAILED",
			Hint:    "Check storage backend permissions and availability",
		}
	}

	return nil
}

// DeleteCredential removes stored credentials
func (m *DefaultManager) DeleteCredential(ctx context.Context, platform Platform, repository string) error {
	m.mu.RLock()
	stores := make([]Store, len(m.stores))
	copy(stores, m.stores)
	m.mu.RUnlock()

	var lastErr error
	deleted := false

	for _, store := range stores {
		if err := store.Delete(ctx, platform, repository); err == nil {
			deleted = true
		} else {
			lastErr = err
		}
	}

	if !deleted && lastErr != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to delete credentials from storage",
			Code:    "AUTH_DELETE_FAILED",
		}
	}

	return nil
}

// ListCredentials lists all stored credentials
func (m *DefaultManager) ListCredentials(ctx context.Context) ([]*Credential, error) {
	m.mu.RLock()
	stores := make([]Store, len(m.stores))
	copy(stores, m.stores)
	m.mu.RUnlock()

	credMap := make(map[string]*Credential)

	for _, store := range stores {
		if creds, err := store.List(ctx); err == nil {
			for _, cred := range creds {
				key := fmt.Sprintf("%s/%s", cred.Platform, cred.ID)
				if existing, exists := credMap[key]; !exists || cred.UpdatedAt.After(existing.UpdatedAt) {
					credMap[key] = cred
				}
			}
		}
	}

	result := make([]*Credential, 0, len(credMap))
	for _, cred := range credMap {
		result = append(result, cred)
	}

	// Sort by platform and ID
	sort.Slice(result, func(i, j int) bool {
		if result[i].Platform != result[j].Platform {
			return result[i].Platform < result[j].Platform
		}
		return result[i].ID < result[j].ID
	})

	return result, nil
}

// RefreshCredential refreshes expired credentials
func (m *DefaultManager) RefreshCredential(ctx context.Context, platform Platform, repository string) (*Credential, error) {
	cred, err := m.GetCredential(ctx, platform, repository)
	if err != nil {
		return nil, err
	}

	m.mu.RLock()
	authenticator, exists := m.authenticators[platform]
	m.mu.RUnlock()

	if !exists {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("No authenticator available for platform: %s", platform),
			Code:    "AUTH_PLATFORM_NOT_SUPPORTED",
		}
	}

	// Try to refresh using the authenticator
	if refreshed, err := authenticator.RefreshToken(ctx, cred.Token); err == nil {
		// Store the refreshed credential
		if err := m.StoreCredential(ctx, refreshed); err != nil {
			return nil, err
		}
		return refreshed, nil
	}

	// If refresh fails, perform new authentication
	req := &AuthRequest{
		Platform:    platform,
		Repository:  repository,
		Scopes:      cred.Scopes,
		Interactive: true,
	}

	result, err := m.Authenticate(ctx, req)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, &AuthError{
			Type:    ErrorTypeInvalidCredentials,
			Message: "Failed to refresh credentials",
			Code:    "AUTH_REFRESH_FAILED",
		}
	}

	return result.Credential, nil
}

// getStorePriority returns the priority of a storage backend (lower = higher priority)
func (m *DefaultManager) getStorePriority(store Store) int {
	switch s := store.(type) {
	case *KeychainStore:
		return 1 // Highest priority
	case *FileStore:
		return 2
	default:
		_ = s    // suppress unused variable warning
		return 3 // Lowest priority
	}
}
