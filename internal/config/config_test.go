package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadNonExistentConfig(t *testing.T) {
	// Ensure environment variable is not set
	os.Unsetenv(EnvVarToken)

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
			os.Unsetenv(EnvVarToken)

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
		os.Setenv(EnvVarToken, "env-token")
		defer os.Unsetenv(EnvVarToken)

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
		os.Unsetenv(EnvVarToken)

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
		os.Unsetenv(EnvVarToken)

		config := &Config{}

		_, err := config.GetToken()
		if err == nil {
			t.Error("Expected error for missing token")
		}
	})

	// Test with expired token
	t.Run("Expired token", func(t *testing.T) {
		os.Unsetenv(EnvVarToken)

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
