package auth

import (
	"fmt"
	"os"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

// StatusCmd checks authentication status
type StatusCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the status command
func (cmd *StatusCmd) Run() error {
	formatter := getFormatter(cmd.Format)

	// Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check environment variable override
	envToken := os.Getenv(config.EnvVarToken)
	if envToken != "" {
		formatter.Info(fmt.Sprintf("✓ Authenticated via environment variable %s", config.EnvVarToken))
		formatter.Info(fmt.Sprintf("  Token: %s...%s", envToken[:4], envToken[len(envToken)-4:]))
		return nil
	}

	// Check if authenticated
	if !cfg.IsAuthenticated() {
		formatter.Error("Not authenticated")
		formatter.Info("\nTo authenticate, run:")
		formatter.Info("  plain auth login")
		return fmt.Errorf("not authenticated")
	}

	// Get token info
	token, err := cfg.GetToken()
	if err != nil {
		formatter.Error(err.Error())
		return err
	}

	// Display status
	formatter.Success("Authenticated")

	// Show token preview (first 4 and last 4 characters)
	tokenPreview := token
	if len(token) > 8 {
		tokenPreview = fmt.Sprintf("%s...%s", token[:4], token[len(token)-4:])
	}

	pairs := map[string]string{
		"Config file": cfg.GetConfigPath(),
		"Token":       tokenPreview,
	}

	// Show expiration if available
	if !cfg.ExpiresAt.IsZero() {
		pairs["Expires"] = cfg.ExpiresAt.Format("2006-01-02 15:04:05")
	}

	if cmd.Format == "json" {
		// For JSON output, include more structured data
		data := map[string]interface{}{
			"authenticated": true,
			"config_path":   cfg.GetConfigPath(),
			"token_preview": tokenPreview,
		}
		if !cfg.ExpiresAt.IsZero() {
			data["expires_at"] = cfg.ExpiresAt
		}
		return formatter.PrintJSON(data)
	}

	formatter.Info("")
	return formatter.PrintKeyValue(pairs)
}
