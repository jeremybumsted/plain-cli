package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func TestConfigCmd_Run_NotAuthenticated(t *testing.T) {
	// Clear environment variables
	if err := os.Unsetenv("PLAIN_API_TOKEN"); err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}

	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create unauthenticated config
	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// We can't fully test the interactive flow, but we can verify the logic paths
	if cfg.IsAuthenticated() {
		t.Error("Config should not be authenticated")
	}
}

func TestConfigCmd_Run_PartiallyConfigured(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create partially configured config (authenticated but missing workspace/help center)
	cfg := &config.Config{
		AccessToken: "test-token",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload to verify
	loadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if !loadedCfg.IsAuthenticated() {
		t.Error("Config should be authenticated")
	}

	if loadedCfg.IsFullyConfigured() {
		t.Error("Config should not be fully configured")
	}
}

func TestConfigCmd_Run_FullyConfigured(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create fully configured config
	cfg := &config.Config{
		AccessToken:  "test-token",
		WorkspaceID:  "ws_test123",
		HelpCenterID: "hc_test456",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload to verify
	loadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if !loadedCfg.IsFullyConfigured() {
		t.Error("Config should be fully configured")
	}
}

func TestRunConfigMenuWithReader_ExitOption(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &config.Config{
		AccessToken:  "test-token",
		WorkspaceID:  "ws_test123",
		HelpCenterID: "hc_test456",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test exit option
	input := strings.NewReader("0\n")
	err := runConfigMenuWithReader(cfg, input)
	if err != nil {
		t.Errorf("Expected no error for exit option, got: %v", err)
	}
}

func TestRunConfigMenuWithReader_InvalidSelection(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &config.Config{
		AccessToken:  "test-token",
		WorkspaceID:  "ws_test123",
		HelpCenterID: "hc_test456",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test invalid selection
	input := strings.NewReader("99\n")
	err := runConfigMenuWithReader(cfg, input)
	if err != nil {
		t.Errorf("Expected no error for invalid selection, got: %v", err)
	}
}

func TestRunConfigMenuWithReader_NoInput(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &config.Config{
		AccessToken:  "test-token",
		WorkspaceID:  "ws_test123",
		HelpCenterID: "hc_test456",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test no input (empty reader)
	input := strings.NewReader("")
	err := runConfigMenuWithReader(cfg, input)
	if err == nil {
		t.Error("Expected error for empty input")
	}
}

func TestGetConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Test loading non-existent config
	cfg, err := getConfig(configPath)
	if err != nil {
		t.Errorf("getConfig should not error on non-existent file: %v", err)
		return
	}
	if cfg == nil {
		t.Error("Config should not be nil")
		return
	}

	// Create and save config
	cfg.AccessToken = "test-token"
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test loading existing config
	loadedCfg, err := getConfig(configPath)
	if err != nil {
		t.Errorf("Failed to load existing config: %v", err)
		return
	}
	if loadedCfg.AccessToken != "test-token" {
		t.Errorf("Expected token 'test-token', got '%s'", loadedCfg.AccessToken)
	}
}

func TestShowCurrentSettings(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "Fully configured",
			config: &config.Config{
				AccessToken:  "test-token-1234567890",
				TokenType:    "Bearer",
				WorkspaceID:  "ws_test123",
				HelpCenterID: "hc_test456",
			},
		},
		{
			name: "Partially configured",
			config: &config.Config{
				AccessToken: "test-token-1234567890",
				TokenType:   "Bearer",
			},
		},
		{
			name:   "Not configured",
			config: &config.Config{},
		},
		{
			name: "With expiration",
			config: &config.Config{
				AccessToken: "test-token-1234567890",
				TokenType:   "Bearer",
				ExpiresAt:   time.Now().Add(24 * time.Hour),
			},
		},
		{
			name: "Short token",
			config: &config.Config{
				AccessToken: "short",
				TokenType:   "Bearer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test just verifies showCurrentSettings doesn't panic
			// We can't easily capture stdout without more complex mocking
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("showCurrentSettings panicked: %v", r)
				}
			}()

			showCurrentSettings(tt.config)
		})
	}
}

func TestConfigCmd_Run_WithExpiredToken(t *testing.T) {
	// Clear environment variables
	if err := os.Unsetenv("PLAIN_API_TOKEN"); err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create config with expired token
	cfg := &config.Config{
		AccessToken:  "expired-token",
		WorkspaceID:  "ws_test123",
		HelpCenterID: "hc_test456",
		ExpiresAt:    time.Now().Add(-1 * time.Hour),
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload and verify it's not authenticated
	loadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedCfg.IsAuthenticated() {
		t.Error("Config with expired token should not be authenticated")
	}
}

func TestGetConfig_EmptyPath(t *testing.T) {
	// Test with empty path - should use default config location
	cfg, err := getConfig("")
	if err != nil {
		t.Errorf("getConfig with empty path should not error: %v", err)
	}
	if cfg == nil {
		t.Error("Config should not be nil")
	}
}

func TestRunConfigMenuWithReader_WhitespaceInSelection(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &config.Config{
		AccessToken:  "test-token",
		WorkspaceID:  "ws_test123",
		HelpCenterID: "hc_test456",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test with whitespace - should still work (exit option)
	input := strings.NewReader("  0  \n")
	err := runConfigMenuWithReader(cfg, input)
	if err != nil {
		t.Errorf("Expected no error with whitespace, got: %v", err)
	}
}

func TestConfigCmd_ConfigPath_Persistence(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create config
	cfg := &config.Config{
		AccessToken: "test-token",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should exist after save")
	}

	// Load and verify
	loadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedCfg.GetConfigPath() != configPath {
		t.Errorf("Config path should be preserved, expected %s, got %s",
			configPath, loadedCfg.GetConfigPath())
	}
}

func TestShowCurrentSettingsWithUser(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		expectUserInfo bool
	}{
		{
			name: "With user configured",
			config: &config.Config{
				AccessToken:  "test-token-1234567890",
				WorkspaceID:  "ws_test123",
				HelpCenterID: "hc_test456",
				UserID:       "u_123456",
				UserEmail:    "test@example.com",
				UserFullName: "Test User",
			},
			expectUserInfo: true,
		},
		{
			name: "Without user configured",
			config: &config.Config{
				AccessToken:  "test-token-1234567890",
				WorkspaceID:  "ws_test123",
				HelpCenterID: "hc_test456",
			},
			expectUserInfo: false,
		},
		{
			name: "With partial user info (should show if UserID present)",
			config: &config.Config{
				AccessToken: "test-token-1234567890",
				UserID:      "u_789012",
				UserEmail:   "partial@example.com",
			},
			expectUserInfo: true,
		},
		{
			name: "With empty config",
			config: &config.Config{},
			expectUserInfo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies that showCurrentSettings doesn't panic
			// and properly handles user info display
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("showCurrentSettings panicked: %v", r)
				}
			}()

			// Verify the HasUserConfigured method works as expected
			hasUser := tt.config.HasUserConfigured()
			if hasUser != tt.expectUserInfo {
				t.Errorf("HasUserConfigured() = %v, want %v", hasUser, tt.expectUserInfo)
			}

			// Call the function to ensure it doesn't panic
			showCurrentSettings(tt.config)
		})
	}
}
