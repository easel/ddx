package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/easel/ddx/internal/auth"
	"github.com/spf13/cobra"
)

var authManager *auth.DefaultManager


// runAuthLogin implements the auth login command logic
func runAuthLogin(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	manager := getAuthManager()

	var repository string
	if len(args) > 0 {
		repository = args[0]
	} else {
		repository = "github.com" // Default
	}

	platform := detectPlatform(repository)
	method, _ := cmd.Flags().GetString("method")
	scopes, _ := cmd.Flags().GetStringSlice("scopes")

	req := &auth.AuthRequest{
		Platform:    platform,
		Repository:  repository,
		Method:      auth.AuthMethod(method),
		Scopes:      scopes,
		Interactive: true,
	}

	result, err := manager.Authenticate(ctx, req)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if result.Success {
		fmt.Printf("✅ Successfully authenticated with %s using %s\n", repository, result.Method)
		if result.Credential != nil && result.Credential.Username != "" {
			fmt.Printf("👤 Authenticated as: %s\n", result.Credential.Username)
		}
	} else {
		fmt.Printf("❌ Authentication failed: %s\n", result.Message)
		if result.Error != nil {
			return result.Error
		}
		return fmt.Errorf("authentication failed")
	}

	return nil
}

// runAuthStatus implements the auth status command logic
func runAuthStatus(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	manager := getAuthManager()

	credentials, err := manager.ListCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to list credentials: %w", err)
	}

	if len(credentials) == 0 {
		fmt.Println("No stored credentials found.")
		fmt.Println("\nTo authenticate, run:")
		fmt.Println("  ddx auth login <platform>")
		return nil
	}

	fmt.Printf("Authentication Status (%d credentials)\n\n", len(credentials))

	for _, cred := range credentials {
		fmt.Printf("🔑 %s (%s)\n", cred.ID, cred.Platform)
		fmt.Printf("   Method: %s\n", cred.Method)
		if cred.Username != "" {
			fmt.Printf("   User: %s\n", cred.Username)
		}
		fmt.Printf("   Created: %s\n", cred.CreatedAt.Format(time.RFC3339))
		fmt.Printf("   Updated: %s\n", cred.UpdatedAt.Format(time.RFC3339))

		if cred.ExpiresAt != nil {
			if time.Now().After(*cred.ExpiresAt) {
				fmt.Printf("   Status: ❌ EXPIRED\n")
			} else {
				fmt.Printf("   Status: ✅ Valid (expires %s)\n", cred.ExpiresAt.Format(time.RFC3339))
			}
		} else {
			// Validate current credentials
			if err := manager.ValidateCredentials(ctx, cred.Platform, cred.ID); err != nil {
				fmt.Printf("   Status: ❌ INVALID (%s)\n", err.Error())
			} else {
				fmt.Printf("   Status: ✅ Valid\n")
			}
		}

		if len(cred.Scopes) > 0 {
			fmt.Printf("   Scopes: %s\n", strings.Join(cred.Scopes, ", "))
		}

		fmt.Println()
	}

	return nil
}

// runAuthList implements the auth list command logic
func runAuthList(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	manager := getAuthManager()

	credentials, err := manager.ListCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to list credentials: %w", err)
	}

	if len(credentials) == 0 {
		fmt.Println("No stored credentials found.")
		return nil
	}

	format, _ := cmd.Flags().GetString("format")
	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(credentials)
	default:
		for _, cred := range credentials {
			fmt.Printf("%s\t%s\t%s\t%s\n", cred.Platform, cred.ID, cred.Method, cred.Username)
		}
	}

	return nil
}

// runAuthLogout implements the auth logout command logic
func runAuthLogout(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	manager := getAuthManager()

	repository := args[0]
	platform := detectPlatform(repository)

	if err := manager.DeleteCredential(ctx, platform, repository); err != nil {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}

	fmt.Printf("✅ Removed credentials for %s\n", repository)
	return nil
}

// runAuthToken implements the auth token command logic
func runAuthToken(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	manager := getAuthManager()

	repository := args[0]
	token := args[1]
	platform := detectPlatform(repository)

	// Validate token format
	if authenticator := getAuthenticator(platform); authenticator != nil {
		if err := authenticator.ValidateToken(ctx, token, nil); err != nil {
			return fmt.Errorf("invalid token: %w", err)
		}
	}

	// Create credential
	cred := &auth.Credential{
		ID:        repository,
		Platform:  platform,
		Method:    auth.AuthMethodToken,
		Token:     token,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := manager.StoreCredential(ctx, cred); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	fmt.Printf("✅ Stored personal access token for %s\n", repository)
	return nil
}


// getAuthManager returns the global authentication manager, initializing it if needed
func getAuthManager() *auth.DefaultManager {
	if authManager == nil {
		authManager = initializeAuthManager()
	}
	return authManager
}

// initializeAuthManager initializes the authentication manager with all backends
func initializeAuthManager() *auth.DefaultManager {
	manager := auth.NewDefaultManager()

	// Register authenticators
	manager.RegisterAuthenticator(auth.NewGitHubAuthenticator())
	// manager.RegisterAuthenticator(auth.NewGitLabAuthenticator())  // TODO: Implement
	// manager.RegisterAuthenticator(auth.NewBitbucketAuthenticator()) // TODO: Implement

	// Register storage backends
	homeDir, _ := os.UserHomeDir()

	// Try keychain first
	keychainStore := auth.NewKeychainStore("com.ddx.auth")
	manager.RegisterStore(keychainStore)

	// File storage as fallback
	credFile := filepath.Join(homeDir, ".ddx", "credentials.enc")
	passphrase := getOrCreatePassphrase()
	fileStore := auth.NewFileStore(credFile, passphrase)
	manager.RegisterStore(fileStore)

	// Register credential helpers
	if gitHelper := auth.NewGitCredentialHelper(); gitHelper.IsAvailable() {
		manager.RegisterCredentialHelper(gitHelper)
	}
	if ghHelper := auth.NewGitHubCLIHelper(); ghHelper.IsAvailable() {
		manager.RegisterCredentialHelper(ghHelper)
	}

	// Register SSH agent
	sshAgent := auth.NewDefaultSSHAgent()
	manager.SetSSHAgent(sshAgent)

	return manager
}

// detectPlatform detects the platform from a repository URL or hostname
func detectPlatform(repository string) auth.Platform {
	switch {
	case strings.Contains(repository, "github.com"):
		return auth.PlatformGitHub
	case strings.Contains(repository, "gitlab.com"):
		return auth.PlatformGitLab
	case strings.Contains(repository, "bitbucket.org"):
		return auth.PlatformBitbucket
	default:
		return auth.PlatformGeneric
	}
}

// getAuthenticator returns the authenticator for a platform
func getAuthenticator(platform auth.Platform) auth.Authenticator {
	// Create a new authenticator for the platform
	switch platform {
	case auth.PlatformGitHub:
		return auth.NewGitHubAuthenticator()
	default:
		return nil
	}
}

// getOrCreatePassphrase gets or creates a passphrase for file encryption
func getOrCreatePassphrase() string {
	// In a real implementation, this would:
	// 1. Try to get from environment variable
	// 2. Try to get from system keychain
	// 3. Generate and store a new one
	// For now, use a default (not secure!)
	return "ddx-auth-passphrase-change-me"
}
