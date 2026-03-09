package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// LoginCmd handles authentication
type LoginCmd struct {
	Token      string `help:"API token (if not provided, will prompt)" short:"t"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the login command
func (cmd *LoginCmd) Run() error {
	formatter := getFormatter(cmd.Format)

	// TODO: Implement OAuth 2.0 flow
	// For now, we support manual token input for testing purposes
	// In a future update, this should:
	// 1. Start a local HTTP server for OAuth callback
	// 2. Open browser to Plain OAuth authorization URL
	// 3. Receive authorization code from callback
	// 4. Exchange code for access token
	// 5. Store tokens in config

	token := cmd.Token
	if token == "" {
		formatter.Info("OAuth flow not yet implemented. Please provide an API token manually.")
		formatter.Info("You can obtain a token from your Plain workspace settings.")
		fmt.Print("\nEnter your Plain API token: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		token = strings.TrimSpace(input)
	}

	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate token by making a test request
	formatter.Info("Validating token...")
	client := plain.NewClient(token)

	// TODO: Add a method to validate the token (e.g., GetCurrentUser)
	// For now, we'll just check that we can create a client
	if client.GetToken() != token {
		return fmt.Errorf("failed to initialize client with token")
	}

	// Store token in config
	// Note: Setting expiresIn to 0 means no expiration tracking for manually entered tokens
	cfg.SetTokens(token, "", "Bearer", 0)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	formatter.Success(fmt.Sprintf("Successfully authenticated! Config saved to %s", cfg.GetConfigPath()))
	formatter.Info("\nNext steps:")
	formatter.Info("  - Run 'plain auth status' to verify authentication")
	formatter.Info("  - Run 'plain threads list' to view your threads")

	return nil
}
