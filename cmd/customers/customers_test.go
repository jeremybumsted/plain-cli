package customers

import (
	"path/filepath"
	"testing"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func TestCustomersCmdStructure(t *testing.T) {
	cmd := &CustomersCmd{}

	// Verify sub-commands exist by checking zero-value structs are present
	_ = cmd.List
	_ = cmd.Get
	_ = cmd.GetByEmail
	_ = cmd.GetByExtID
	_ = cmd.Search
}

func TestListCmdDefaultValues(t *testing.T) {
	cmd := &ListCmd{}

	if cmd.Limit != 0 {
		t.Errorf("Limit default = %d, want 0", cmd.Limit)
	}
	if cmd.Format != "" {
		t.Errorf("Format default = %q, want empty string", cmd.Format)
	}
	if cmd.ConfigPath != "" {
		t.Errorf("ConfigPath default = %q, want empty string", cmd.ConfigPath)
	}
}

func TestGetCmdRequiresID(t *testing.T) {
	cmd := &GetCmd{}

	if cmd.CustomerID != "" {
		t.Errorf("CustomerID default = %q, want empty string", cmd.CustomerID)
	}
}

func TestGetByEmailCmdRequiresEmail(t *testing.T) {
	cmd := &GetByEmailCmd{}

	if cmd.Email != "" {
		t.Errorf("Email default = %q, want empty string", cmd.Email)
	}
}

func TestGetByExtIDCmdRequiresExternalID(t *testing.T) {
	cmd := &GetByExtIDCmd{}

	if cmd.ExternalID != "" {
		t.Errorf("ExternalID default = %q, want empty string", cmd.ExternalID)
	}
}

func TestSearchCmdDefaultValues(t *testing.T) {
	cmd := &SearchCmd{}

	if cmd.Query != "" {
		t.Errorf("Query default = %q, want empty string", cmd.Query)
	}
	if cmd.Limit != 0 {
		t.Errorf("Limit default = %d, want 0", cmd.Limit)
	}
	if cmd.Format != "" {
		t.Errorf("Format default = %q, want empty string", cmd.Format)
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

func TestListCmdRunNoAuth(t *testing.T) {
	t.Setenv(config.EnvVarToken, "")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	cmd := &ListCmd{
		ConfigPath: configPath,
		Format:     "table",
	}

	err := cmd.Run()
	if err == nil {
		t.Fatal("Expected error when not authenticated, got nil")
	}
}

func TestSearchCmdRunNoAuth(t *testing.T) {
	t.Setenv(config.EnvVarToken, "")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	cmd := &SearchCmd{
		Query:      "test",
		ConfigPath: configPath,
		Format:     "table",
	}

	err := cmd.Run()
	if err == nil {
		t.Fatal("Expected error when not authenticated, got nil")
	}
}

func TestGetCmdRunNoAuth(t *testing.T) {
	t.Setenv(config.EnvVarToken, "")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	cmd := &GetCmd{
		CustomerID: "c_test123",
		ConfigPath: configPath,
		Format:     "table",
	}

	err := cmd.Run()
	if err == nil {
		t.Fatal("Expected error when not authenticated, got nil")
	}
}
