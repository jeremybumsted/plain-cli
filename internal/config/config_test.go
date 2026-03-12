package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadNonExistentConfig(t *testing.T) {
	// Ensure environment variable is not set
	if err := os.Unsetenv(EnvVarToken); err != nil {
		t.Fatalf("Failed to unset environment variable: %v", err)
	}

	// Create temp path that doesn't exist
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "nonexistent", "config.json")

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load should not error on nonexistent file: %v", err)
	}

	if config == nil {
		t.Fatal("Config should not be nil")
	}

	if config.DefaultFormat != "table" {
		t.Errorf("Expected default format 'table', got '%s'", config.DefaultFormat)
	}

	if config.IsAuthenticated() {
		t.Error("New config should not be authenticated")
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create and save config
	config := &Config{
		configPath:    configPath,
		DefaultFormat: "json",
	}
	config.SetTokens("test-access-token", "test-refresh-token", "Bearer", 3600)

	if err := config.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load config
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values
	if loadedConfig.AccessToken != "test-access-token" {
		t.Errorf("Expected access token 'test-access-token', got '%s'", loadedConfig.AccessToken)
	}

	if loadedConfig.RefreshToken != "test-refresh-token" {
		t.Errorf("Expected refresh token 'test-refresh-token', got '%s'", loadedConfig.RefreshToken)
	}

	if loadedConfig.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got '%s'", loadedConfig.TokenType)
	}

	if loadedConfig.DefaultFormat != "json" {
		t.Errorf("Expected default format 'json', got '%s'", loadedConfig.DefaultFormat)
	}
}

func TestIsAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected bool
	}{
		{
			name:     "No token",
			config:   &Config{},
			expected: false,
		},
		{
			name: "Valid token, no expiration",
			config: &Config{
				AccessToken: "test-token",
			},
			expected: true,
		},
		{
			name: "Valid token, future expiration",
			config: &Config{
				AccessToken: "test-token",
				ExpiresAt:   time.Now().Add(1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "Expired token",
			config: &Config{
				AccessToken: "test-token",
				ExpiresAt:   time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variable
			if err := os.Unsetenv(EnvVarToken); err != nil {
				t.Fatalf("Failed to unset environment variable: %v", err)
			}

			result := tt.config.IsAuthenticated()
			if result != tt.expected {
				t.Errorf("Expected IsAuthenticated() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	// Test with environment variable
	t.Run("Environment variable override", func(t *testing.T) {
		if err := os.Setenv(EnvVarToken, "env-token"); err != nil {
			t.Fatalf("Failed to set environment variable: %v", err)
		}
		defer func() { _ = os.Unsetenv(EnvVarToken) }()

		config := &Config{
			AccessToken: "config-token",
		}

		token, err := config.GetToken()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if token != "env-token" {
			t.Errorf("Expected env token 'env-token', got '%s'", token)
		}
	})

	// Test with config token
	t.Run("Config token", func(t *testing.T) {
		if err := os.Unsetenv(EnvVarToken); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		config := &Config{
			AccessToken: "config-token",
		}

		token, err := config.GetToken()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if token != "config-token" {
			t.Errorf("Expected config token 'config-token', got '%s'", token)
		}
	})

	// Test with no token
	t.Run("No token", func(t *testing.T) {
		if err := os.Unsetenv(EnvVarToken); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		config := &Config{}

		_, err := config.GetToken()
		if err == nil {
			t.Error("Expected error for missing token")
		}
	})

	// Test with expired token
	t.Run("Expired token", func(t *testing.T) {
		if err := os.Unsetenv(EnvVarToken); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		config := &Config{
			AccessToken: "expired-token",
			ExpiresAt:   time.Now().Add(-1 * time.Hour),
		}

		_, err := config.GetToken()
		if err == nil {
			t.Error("Expected error for expired token")
		}
	})
}

func TestClear(t *testing.T) {
	config := &Config{
		AccessToken:  "test-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	config.Clear()

	if config.AccessToken != "" {
		t.Error("AccessToken should be empty after Clear()")
	}

	if config.RefreshToken != "" {
		t.Error("RefreshToken should be empty after Clear()")
	}

	if config.TokenType != "" {
		t.Error("TokenType should be empty after Clear()")
	}

	if !config.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should be zero after Clear()")
	}
}

func TestGetHelpCenterID(t *testing.T) {
	// Test with environment variable
	t.Run("Environment variable override", func(t *testing.T) {
		if err := os.Setenv(EnvVarHelpCenterID, "env-hc-id"); err != nil {
			t.Fatalf("Failed to set environment variable: %v", err)
		}
		defer func() { _ = os.Unsetenv(EnvVarHelpCenterID) }()

		config := &Config{
			HelpCenterID: "config-hc-id",
		}

		id, err := config.GetHelpCenterID()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if id != "env-hc-id" {
			t.Errorf("Expected env help center ID 'env-hc-id', got '%s'", id)
		}
	})

	// Test with config help center ID
	t.Run("Config help center ID", func(t *testing.T) {
		if err := os.Unsetenv(EnvVarHelpCenterID); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		config := &Config{
			HelpCenterID: "config-hc-id",
		}

		id, err := config.GetHelpCenterID()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if id != "config-hc-id" {
			t.Errorf("Expected config help center ID 'config-hc-id', got '%s'", id)
		}
	})

	// Test with no help center ID
	t.Run("No help center ID", func(t *testing.T) {
		if err := os.Unsetenv(EnvVarHelpCenterID); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		config := &Config{}

		_, err := config.GetHelpCenterID()
		if err == nil {
			t.Error("Expected error for missing help center ID")
		}
	})
}

func TestSetHelpCenterID(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	config := &Config{
		configPath: configPath,
	}

	// Set help center ID
	err := config.SetHelpCenterID("hc_test123")
	if err != nil {
		t.Fatalf("Failed to set help center ID: %v", err)
	}

	// Verify it was saved
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedConfig.HelpCenterID != "hc_test123" {
		t.Errorf("Expected help center ID 'hc_test123', got '%s'", loadedConfig.HelpCenterID)
	}
}

func TestGetWorkspaceID(t *testing.T) {
	// Test with environment variable
	t.Run("Environment variable override", func(t *testing.T) {
		if err := os.Setenv(EnvVarWorkspaceID, "env-ws-id"); err != nil {
			t.Fatalf("Failed to set environment variable: %v", err)
		}
		defer func() { _ = os.Unsetenv(EnvVarWorkspaceID) }()

		config := &Config{
			WorkspaceID: "config-ws-id",
		}

		id, err := config.GetWorkspaceID()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if id != "env-ws-id" {
			t.Errorf("Expected env workspace ID 'env-ws-id', got '%s'", id)
		}
	})

	// Test with config workspace ID
	t.Run("Config workspace ID", func(t *testing.T) {
		if err := os.Unsetenv(EnvVarWorkspaceID); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		config := &Config{
			WorkspaceID: "config-ws-id",
		}

		id, err := config.GetWorkspaceID()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if id != "config-ws-id" {
			t.Errorf("Expected config workspace ID 'config-ws-id', got '%s'", id)
		}
	})

	// Test with no workspace ID
	t.Run("No workspace ID", func(t *testing.T) {
		if err := os.Unsetenv(EnvVarWorkspaceID); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		config := &Config{}

		_, err := config.GetWorkspaceID()
		if err == nil {
			t.Error("Expected error for missing workspace ID")
		}
	})
}

func TestSetWorkspaceID(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	config := &Config{
		configPath: configPath,
	}

	// Set workspace ID
	err := config.SetWorkspaceID("ws_test123")
	if err != nil {
		t.Fatalf("Failed to set workspace ID: %v", err)
	}

	// Verify it was saved
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loadedConfig.WorkspaceID != "ws_test123" {
		t.Errorf("Expected workspace ID 'ws_test123', got '%s'", loadedConfig.WorkspaceID)
	}
}

func TestIsFullyConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		setupEnv func()
		expected bool
	}{
		{
			name:     "No configuration",
			config:   &Config{},
			setupEnv: func() {},
			expected: false,
		},
		{
			name: "Only access token",
			config: &Config{
				AccessToken: "test-token",
			},
			setupEnv: func() {},
			expected: false,
		},
		{
			name: "Access token and workspace ID",
			config: &Config{
				AccessToken: "test-token",
				WorkspaceID: "ws_test123",
			},
			setupEnv: func() {},
			expected: false,
		},
		{
			name: "Access token and help center ID",
			config: &Config{
				AccessToken:  "test-token",
				HelpCenterID: "hc_test123",
			},
			setupEnv: func() {},
			expected: false,
		},
		{
			name: "All required fields configured",
			config: &Config{
				AccessToken:  "test-token",
				WorkspaceID:  "ws_test123",
				HelpCenterID: "hc_test123",
			},
			setupEnv: func() {},
			expected: true,
		},
		{
			name: "All required fields configured with future expiration",
			config: &Config{
				AccessToken:  "test-token",
				WorkspaceID:  "ws_test123",
				HelpCenterID: "hc_test123",
				ExpiresAt:    time.Now().Add(1 * time.Hour),
			},
			setupEnv: func() {},
			expected: true,
		},
		{
			name: "Expired token",
			config: &Config{
				AccessToken:  "test-token",
				WorkspaceID:  "ws_test123",
				HelpCenterID: "hc_test123",
				ExpiresAt:    time.Now().Add(-1 * time.Hour),
			},
			setupEnv: func() {},
			expected: false,
		},
		{
			name: "Config values with env var overrides",
			config: &Config{
				AccessToken:  "test-token",
				WorkspaceID:  "",
				HelpCenterID: "",
			},
			setupEnv: func() {
				_ = os.Setenv(EnvVarWorkspaceID, "env-ws-id")
				_ = os.Setenv(EnvVarHelpCenterID, "env-hc-id")
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			_ = os.Unsetenv(EnvVarToken)
			_ = os.Unsetenv(EnvVarWorkspaceID)
			_ = os.Unsetenv(EnvVarHelpCenterID)

			// Setup environment if needed
			tt.setupEnv()

			// Cleanup after test
			defer func() {
				_ = os.Unsetenv(EnvVarWorkspaceID)
				_ = os.Unsetenv(EnvVarHelpCenterID)
			}()

			result := tt.config.IsFullyConfigured()
			if result != tt.expected {
				t.Errorf("Expected IsFullyConfigured() = %v, got %v", tt.expected, result)
			}
		})
	}
}
