package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	// DefaultConfigDir is the default directory for storing config
	DefaultConfigDir = ".config/plain-cli"
	// DefaultConfigFile is the default config filename
	DefaultConfigFile = "config.json"
	// EnvVarToken is the environment variable for overriding the API token
	EnvVarToken = "PLAIN_API_TOKEN"
	// EnvVarHelpCenterID is the environment variable for overriding the help center ID
	EnvVarHelpCenterID = "PLAIN_HELP_CENTER_ID"
	// EnvVarWorkspaceID is the environment variable for overriding the workspace ID
	EnvVarWorkspaceID = "PLAIN_WORKSPACE_ID"
)

// Config holds the CLI configuration including OAuth credentials
type Config struct {
	// OAuth token fields
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`

	// User preferences
	DefaultFormat string `json:"default_format,omitempty"` // "table", "json"
	HelpCenterID  string `json:"help_center_id,omitempty"` // Default help center ID
	WorkspaceID   string `json:"workspace_id,omitempty"`   // Default workspace ID

	// Internal tracking
	configPath string // Not serialized, used for saving
}

// Load reads the configuration from the specified path.
// If path is empty, uses the default config location.
// Returns a Config with defaults if the file doesn't exist.
func Load(path string) (*Config, error) {
	if path == "" {
		var err error
		path, err = defaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	config := &Config{
		configPath:    path,
		DefaultFormat: "table", // Default to table output
	}

	// If file doesn't exist, return config with defaults
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	config.configPath = path
	return config, nil
}

// Save persists the configuration to disk
func (c *Config) Save() error {
	if c.configPath == "" {
		var err error
		c.configPath, err = defaultConfigPath()
		if err != nil {
			return err
		}
	}

	// Ensure config directory exists
	dir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(c.configPath, data, 0600); err != nil {
		return err
	}

	return nil
}

// IsAuthenticated checks if the user has valid authentication credentials
func (c *Config) IsAuthenticated() bool {
	// Check environment variable first
	if os.Getenv(EnvVarToken) != "" {
		return true
	}

	// Check stored access token
	if c.AccessToken == "" {
		return false
	}

	// Check if token is expired
	if !c.ExpiresAt.IsZero() && time.Now().After(c.ExpiresAt) {
		return false
	}

	return true
}

// GetToken returns the active API token, preferring environment variable
func (c *Config) GetToken() (string, error) {
	// Environment variable takes precedence
	if token := os.Getenv(EnvVarToken); token != "" {
		return token, nil
	}

	// Use stored token
	if c.AccessToken == "" {
		return "", errors.New("not authenticated: run 'plain auth login'")
	}

	// Check expiration
	if !c.ExpiresAt.IsZero() && time.Now().After(c.ExpiresAt) {
		return "", errors.New("token expired: run 'plain auth login'")
	}

	return c.AccessToken, nil
}

// SetTokens stores OAuth tokens in the config
func (c *Config) SetTokens(accessToken, refreshToken, tokenType string, expiresIn int) {
	c.AccessToken = accessToken
	c.RefreshToken = refreshToken
	c.TokenType = tokenType

	if expiresIn > 0 {
		c.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}
}

// Clear removes all authentication credentials
func (c *Config) Clear() {
	c.AccessToken = ""
	c.RefreshToken = ""
	c.TokenType = ""
	c.ExpiresAt = time.Time{}
}

// defaultConfigPath returns the default config file path
func defaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, DefaultConfigDir, DefaultConfigFile), nil
}

// GetConfigPath returns the current config file path
func (c *Config) GetConfigPath() string {
	if c.configPath == "" {
		path, _ := defaultConfigPath()
		return path
	}
	return c.configPath
}

// SetConfigPath sets the config file path (useful for testing)
func (c *Config) SetConfigPath(path string) {
	c.configPath = path
}

// GetHelpCenterID returns the configured help center ID, checking env var first
func (c *Config) GetHelpCenterID() (string, error) {
	// Environment variable takes precedence
	if id := os.Getenv(EnvVarHelpCenterID); id != "" {
		return id, nil
	}

	// Use stored help center ID
	if c.HelpCenterID == "" {
		return "", errors.New("no help center configured: run 'plain config'")
	}

	return c.HelpCenterID, nil
}

// SetHelpCenterID stores the help center ID in the config
func (c *Config) SetHelpCenterID(id string) error {
	c.HelpCenterID = id
	return c.Save()
}

// GetWorkspaceID returns the configured workspace ID, checking env var first
func (c *Config) GetWorkspaceID() (string, error) {
	// Environment variable takes precedence
	if id := os.Getenv(EnvVarWorkspaceID); id != "" {
		return id, nil
	}

	// Use stored workspace ID
	if c.WorkspaceID == "" {
		return "", errors.New("no workspace configured: run 'plain config'")
	}

	return c.WorkspaceID, nil
}

// SetWorkspaceID stores the workspace ID in the config
func (c *Config) SetWorkspaceID(id string) error {
	c.WorkspaceID = id
	return c.Save()
}

// IsFullyConfigured returns true if all required configuration is present
func (c *Config) IsFullyConfigured() bool {
	// Check if authenticated (token present and not expired)
	if !c.IsAuthenticated() {
		return false
	}

	// Check if workspace ID is configured
	_, err := c.GetWorkspaceID()
	if err != nil {
		return false
	}

	// Check if help center ID is configured
	_, err = c.GetHelpCenterID()
	return err == nil
}
