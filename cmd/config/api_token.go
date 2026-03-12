package config

import (
	"fmt"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

func configureAPIToken(cfg *config.Config) error {
	fmt.Println("[1/3] API Token")
	fmt.Println("You can generate an API token at: https://app.plain.com/developer/api-keys")
	fmt.Println()

	var token string
	fmt.Print("Enter your Plain API token: ")
	if _, err := fmt.Scanln(&token); err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Validate token by making a test API call
	fmt.Println("Validating token...")
	testClient := plain.NewClient(token)

	// Try to list workspaces as a validation check
	_, err := testClient.ListWorkspaces()
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	// Token is valid, save it
	// Note: We're storing it as AccessToken for now
	// In a full OAuth implementation, this would be different
	cfg.AccessToken = token
	cfg.TokenType = "Bearer"
	// Don't set ExpiresAt for API tokens (they don't expire like OAuth tokens)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Token validated and saved")
	fmt.Println()

	return nil
}
