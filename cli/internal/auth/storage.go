package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
)

// KeychainStore implements credential storage using OS keychain
type KeychainStore struct {
	serviceName string
}

// NewKeychainStore creates a new keychain-based credential store
func NewKeychainStore(serviceName string) *KeychainStore {
	return &KeychainStore{
		serviceName: serviceName,
	}
}

// IsAvailable checks if keychain is available on the current platform
func (s *KeychainStore) IsAvailable() bool {
	// Temporarily disable keychain until platform-specific implementations are complete
	// This will make file storage the primary storage method
	return false
}

// Get retrieves a credential from the keychain
func (s *KeychainStore) Get(ctx context.Context, platform Platform, repository string) (*Credential, error) {
	key := s.makeKey(platform, repository)

	data, err := s.getFromKeychain(key)
	if err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeNotFound,
			Message: fmt.Sprintf("Credential not found in keychain for %s/%s", platform, repository),
			Code:    "KEYCHAIN_NOT_FOUND",
		}
	}

	var cred Credential
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to decode credential from keychain",
			Code:    "KEYCHAIN_DECODE_ERROR",
		}
	}

	return &cred, nil
}

// Set stores a credential in the keychain
func (s *KeychainStore) Set(ctx context.Context, cred *Credential) error {
	key := s.makeKey(cred.Platform, cred.ID)

	data, err := json.Marshal(cred)
	if err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to encode credential for keychain",
			Code:    "KEYCHAIN_ENCODE_ERROR",
		}
	}

	if err := s.setInKeychain(key, data); err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to store credential in keychain",
			Code:    "KEYCHAIN_STORE_ERROR",
		}
	}

	return nil
}

// Delete removes a credential from the keychain
func (s *KeychainStore) Delete(ctx context.Context, platform Platform, repository string) error {
	key := s.makeKey(platform, repository)

	if err := s.deleteFromKeychain(key); err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to delete credential from keychain",
			Code:    "KEYCHAIN_DELETE_ERROR",
		}
	}

	return nil
}

// List returns all stored credentials from the keychain
func (s *KeychainStore) List(ctx context.Context) ([]*Credential, error) {
	keys, err := s.listKeychainKeys()
	if err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to list credentials from keychain",
			Code:    "KEYCHAIN_LIST_ERROR",
		}
	}

	var credentials []*Credential
	for _, key := range keys {
		if data, err := s.getFromKeychain(key); err == nil {
			var cred Credential
			if err := json.Unmarshal(data, &cred); err == nil {
				credentials = append(credentials, &cred)
			}
		}
	}

	return credentials, nil
}

// Clear removes all stored credentials from the keychain
func (s *KeychainStore) Clear(ctx context.Context) error {
	keys, err := s.listKeychainKeys()
	if err != nil {
		return err
	}

	for _, key := range keys {
		_ = s.deleteFromKeychain(key) // Best effort deletion
	}

	return nil
}

// makeKey creates a keychain key for the credential
func (s *KeychainStore) makeKey(platform Platform, repository string) string {
	return fmt.Sprintf("%s.%s.%s", s.serviceName, platform, repository)
}

// Platform-specific keychain implementations would go here

func (s *KeychainStore) getFromKeychain(key string) ([]byte, error) {
	// Keychain access not fully implemented yet - fail gracefully
	return nil, &AuthError{
		Type:    ErrorTypeNotFound,
		Message: "Keychain integration not available",
		Code:    "KEYCHAIN_NOT_AVAILABLE",
	}
}

func (s *KeychainStore) setInKeychain(key string, data []byte) error {
	// Keychain access not fully implemented yet - fail gracefully
	return &AuthError{
		Type:    ErrorTypeStorageError,
		Message: "Keychain integration not available",
		Code:    "KEYCHAIN_NOT_AVAILABLE",
	}
}

func (s *KeychainStore) deleteFromKeychain(key string) error {
	// Keychain access not fully implemented yet - fail gracefully
	return &AuthError{
		Type:    ErrorTypeStorageError,
		Message: "Keychain integration not available",
		Code:    "KEYCHAIN_NOT_AVAILABLE",
	}
}

func (s *KeychainStore) listKeychainKeys() ([]string, error) {
	// Keychain access not fully implemented yet - fail gracefully
	return []string{}, nil
}

// FileStore implements encrypted file-based credential storage
type FileStore struct {
	filePath   string
	passphrase string
}

// NewFileStore creates a new file-based credential store
func NewFileStore(filePath, passphrase string) *FileStore {
	return &FileStore{
		filePath:   filePath,
		passphrase: passphrase,
	}
}

// IsAvailable checks if file storage is available
func (s *FileStore) IsAvailable() bool {
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return false
	}

	// Test write access
	testFile := filepath.Join(dir, ".ddx-auth-test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		return false
	}
	_ = os.Remove(testFile) // Best effort cleanup

	return true
}

// Get retrieves a credential from the encrypted file
func (s *FileStore) Get(ctx context.Context, platform Platform, repository string) (*Credential, error) {
	credentials, err := s.loadCredentials()
	if err != nil {
		return nil, err
	}

	key := s.makeKey(platform, repository)
	if cred, exists := credentials[key]; exists {
		return cred, nil
	}

	return nil, &AuthError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("Credential not found for %s/%s", platform, repository),
		Code:    "FILE_NOT_FOUND",
	}
}

// Set stores a credential in the encrypted file
func (s *FileStore) Set(ctx context.Context, cred *Credential) error {
	credentials, err := s.loadCredentials()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if credentials == nil {
		credentials = make(map[string]*Credential)
	}

	key := s.makeKey(cred.Platform, cred.ID)
	credentials[key] = cred

	return s.saveCredentials(credentials)
}

// Delete removes a credential from the encrypted file
func (s *FileStore) Delete(ctx context.Context, platform Platform, repository string) error {
	credentials, err := s.loadCredentials()
	if err != nil {
		return err
	}

	key := s.makeKey(platform, repository)
	delete(credentials, key)

	return s.saveCredentials(credentials)
}

// List returns all stored credentials from the encrypted file
func (s *FileStore) List(ctx context.Context) ([]*Credential, error) {
	credentials, err := s.loadCredentials()
	if err != nil {
		return nil, err
	}

	result := make([]*Credential, 0, len(credentials))
	for _, cred := range credentials {
		result = append(result, cred)
	}

	return result, nil
}

// Clear removes all stored credentials from the encrypted file
func (s *FileStore) Clear(ctx context.Context) error {
	if err := os.Remove(s.filePath); err != nil && !os.IsNotExist(err) {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to clear credential file",
			Code:    "FILE_CLEAR_ERROR",
		}
	}
	return nil
}

// makeKey creates a storage key for the credential
func (s *FileStore) makeKey(platform Platform, repository string) string {
	return fmt.Sprintf("%s/%s", platform, repository)
}

// loadCredentials loads and decrypts credentials from file
func (s *FileStore) loadCredentials() (map[string]*Credential, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return make(map[string]*Credential), nil
	}

	encryptedData, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to read credential file",
			Code:    "FILE_READ_ERROR",
		}
	}

	if len(encryptedData) == 0 {
		return make(map[string]*Credential), nil
	}

	decryptedData, err := s.decrypt(encryptedData)
	if err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to decrypt credential file",
			Code:    "FILE_DECRYPT_ERROR",
		}
	}

	var credentials map[string]*Credential
	if err := json.Unmarshal(decryptedData, &credentials); err != nil {
		return nil, &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to parse credential file",
			Code:    "FILE_PARSE_ERROR",
		}
	}

	return credentials, nil
}

// saveCredentials encrypts and saves credentials to file
func (s *FileStore) saveCredentials(credentials map[string]*Credential) error {
	data, err := json.Marshal(credentials)
	if err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to encode credentials",
			Code:    "FILE_ENCODE_ERROR",
		}
	}

	encryptedData, err := s.encrypt(data)
	if err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to encrypt credentials",
			Code:    "FILE_ENCRYPT_ERROR",
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0700); err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to create credential directory",
			Code:    "FILE_DIRECTORY_ERROR",
		}
	}

	if err := os.WriteFile(s.filePath, encryptedData, 0600); err != nil {
		return &AuthError{
			Type:    ErrorTypeStorageError,
			Message: "Failed to write credential file",
			Code:    "FILE_WRITE_ERROR",
		}
	}

	return nil
}

// encrypt encrypts data using AES-GCM with the passphrase
func (s *FileStore) encrypt(data []byte) ([]byte, error) {
	key := s.deriveKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM with the passphrase
func (s *FileStore) decrypt(data []byte) ([]byte, error) {
	key := s.deriveKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("invalid encrypted data")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// deriveKey derives an encryption key from the passphrase
func (s *FileStore) deriveKey() []byte {
	// Use PBKDF2 to derive a key from the passphrase
	salt := []byte("ddx-auth-salt") // In production, use a random salt per file
	return pbkdf2.Key([]byte(s.passphrase), salt, 100000, 32, sha256.New)
}
