package config

import (
	"path/filepath"
	"testing"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func TestUserCmd_Clear(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create config with user info
	cfg := &config.Config{
		AccessToken:    "test-token",
		WorkspaceID:    "ws_test123",
		HelpCenterID:   "hc_test456",
		UserID:         "u_test789",
		UserEmail:      "test@example.com",
		UserFullName:   "Test User",
		UserPublicName: "Test",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify user is configured
	if !cfg.HasUserConfigured() {
		t.Fatal("User should be configured before clear")
	}

	// Run clear command
	cmd := &UserCmd{Clear: true}
	// We can't run the full command without mocking, but we can test the logic
	loadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	loadedCfg.ClearUserInfo()
	if err := loadedCfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload and verify user info is cleared
	reloadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if reloadedCfg.HasUserConfigured() {
		t.Error("User should not be configured after clear")
	}

	if reloadedCfg.UserID != "" {
		t.Errorf("UserID should be empty, got: %s", reloadedCfg.UserID)
	}
	if reloadedCfg.UserEmail != "" {
		t.Errorf("UserEmail should be empty, got: %s", reloadedCfg.UserEmail)
	}
	if reloadedCfg.UserFullName != "" {
		t.Errorf("UserFullName should be empty, got: %s", reloadedCfg.UserFullName)
	}
	if reloadedCfg.UserPublicName != "" {
		t.Errorf("UserPublicName should be empty, got: %s", reloadedCfg.UserPublicName)
	}

	// Verify other config remains intact
	if reloadedCfg.AccessToken != "test-token" {
		t.Error("AccessToken should remain unchanged")
	}
	if reloadedCfg.WorkspaceID != "ws_test123" {
		t.Error("WorkspaceID should remain unchanged")
	}
	if reloadedCfg.HelpCenterID != "hc_test456" {
		t.Error("HelpCenterID should remain unchanged")
	}

	// Test that cmd is properly structured
	if cmd.Clear != true {
		t.Error("Clear flag should be true")
	}
}

func TestUserCmd_SetUserInfo(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create minimal config
	cfg := &config.Config{
		AccessToken:  "test-token",
		WorkspaceID:  "ws_test123",
		HelpCenterID: "hc_test456",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test setting user info
	err := cfg.SetUserInfo(
		"u_newuser123",
		"newuser@example.com",
		"New User Name",
		"NewUser",
	)
	if err != nil {
		t.Fatalf("Failed to set user info: %v", err)
	}

	// Reload and verify
	reloadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if !reloadedCfg.HasUserConfigured() {
		t.Error("User should be configured")
	}

	if reloadedCfg.UserID != "u_newuser123" {
		t.Errorf("Expected UserID 'u_newuser123', got: %s", reloadedCfg.UserID)
	}
	if reloadedCfg.UserEmail != "newuser@example.com" {
		t.Errorf("Expected UserEmail 'newuser@example.com', got: %s", reloadedCfg.UserEmail)
	}
	if reloadedCfg.UserFullName != "New User Name" {
		t.Errorf("Expected UserFullName 'New User Name', got: %s", reloadedCfg.UserFullName)
	}
	if reloadedCfg.UserPublicName != "NewUser" {
		t.Errorf("Expected UserPublicName 'NewUser', got: %s", reloadedCfg.UserPublicName)
	}
}

func TestUserCmd_Flags(t *testing.T) {
	tests := []struct {
		name     string
		cmd      UserCmd
		wantList bool
		wantClear bool
		wantEmail string
	}{
		{
			name:      "List flag set",
			cmd:       UserCmd{List: true},
			wantList:  true,
			wantClear: false,
			wantEmail: "",
		},
		{
			name:      "Clear flag set",
			cmd:       UserCmd{Clear: true},
			wantList:  false,
			wantClear: true,
			wantEmail: "",
		},
		{
			name:      "Email flag set",
			cmd:       UserCmd{Email: "test@example.com"},
			wantList:  false,
			wantClear: false,
			wantEmail: "test@example.com",
		},
		{
			name:      "No flags set",
			cmd:       UserCmd{},
			wantList:  false,
			wantClear: false,
			wantEmail: "",
		},
		{
			name:      "Multiple flags set",
			cmd:       UserCmd{List: true, Email: "test@example.com"},
			wantList:  true,
			wantClear: false,
			wantEmail: "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.List != tt.wantList {
				t.Errorf("List = %v, want %v", tt.cmd.List, tt.wantList)
			}
			if tt.cmd.Clear != tt.wantClear {
				t.Errorf("Clear = %v, want %v", tt.cmd.Clear, tt.wantClear)
			}
			if tt.cmd.Email != tt.wantEmail {
				t.Errorf("Email = %v, want %v", tt.cmd.Email, tt.wantEmail)
			}
		})
	}
}

func TestUserCmd_GetUserID(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		wantError bool
		wantID    string
	}{
		{
			name: "User configured",
			config: &config.Config{
				UserID:       "u_test123",
				UserEmail:    "test@example.com",
				UserFullName: "Test User",
			},
			wantError: false,
			wantID:    "u_test123",
		},
		{
			name:      "User not configured",
			config:    &config.Config{},
			wantError: true,
			wantID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := tt.config.GetUserID()

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if userID != tt.wantID {
				t.Errorf("Got UserID %s, want %s", userID, tt.wantID)
			}
		})
	}
}

func TestUserCmd_ClearPreservesOtherConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Create fully configured config
	cfg := &config.Config{
		AccessToken:    "test-token-abc123",
		RefreshToken:   "refresh-token-xyz",
		TokenType:      "Bearer",
		WorkspaceID:    "ws_test123",
		HelpCenterID:   "hc_test456",
		DefaultFormat:  "json",
		UserID:         "u_test789",
		UserEmail:      "test@example.com",
		UserFullName:   "Test User",
		UserPublicName: "Test",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Clear user info
	cfg.ClearUserInfo()
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save after clear: %v", err)
	}

	// Reload and verify all non-user config is preserved
	reloadedCfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	// Verify user info is cleared
	if reloadedCfg.UserID != "" {
		t.Error("UserID should be empty")
	}

	// Verify other config is preserved
	if reloadedCfg.AccessToken != "test-token-abc123" {
		t.Error("AccessToken should be preserved")
	}
	if reloadedCfg.RefreshToken != "refresh-token-xyz" {
		t.Error("RefreshToken should be preserved")
	}
	if reloadedCfg.TokenType != "Bearer" {
		t.Error("TokenType should be preserved")
	}
	if reloadedCfg.WorkspaceID != "ws_test123" {
		t.Error("WorkspaceID should be preserved")
	}
	if reloadedCfg.HelpCenterID != "hc_test456" {
		t.Error("HelpCenterID should be preserved")
	}
	if reloadedCfg.DefaultFormat != "json" {
		t.Error("DefaultFormat should be preserved")
	}
}

func TestUserCmd_MultipleSetOperations(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &config.Config{
		AccessToken: "test-token",
	}
	cfg.SetConfigPath(configPath)
	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// Set user info first time
	err := cfg.SetUserInfo("u_user1", "user1@example.com", "User One", "User1")
	if err != nil {
		t.Fatalf("Failed first SetUserInfo: %v", err)
	}

	// Reload and verify
	cfg, err = config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload after first set: %v", err)
	}
	if cfg.UserEmail != "user1@example.com" {
		t.Error("First user email not set correctly")
	}

	// Update user info (second time)
	err = cfg.SetUserInfo("u_user2", "user2@example.com", "User Two", "User2")
	if err != nil {
		t.Fatalf("Failed second SetUserInfo: %v", err)
	}

	// Reload and verify update
	cfg, err = config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload after second set: %v", err)
	}
	if cfg.UserEmail != "user2@example.com" {
		t.Error("Second user email not updated correctly")
	}
	if cfg.UserID != "u_user2" {
		t.Error("Second user ID not updated correctly")
	}
}

func TestUserCmd_EmptyStringHandling(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	cfg := &config.Config{
		AccessToken: "test-token",
	}
	cfg.SetConfigPath(configPath)

	// Set user info with some empty strings
	err := cfg.SetUserInfo("u_test", "test@example.com", "", "")
	if err != nil {
		t.Fatalf("Failed to set user info with empty strings: %v", err)
	}

	// Reload and verify
	cfg, err = config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to reload: %v", err)
	}

	if cfg.UserID != "u_test" {
		t.Error("UserID should be set")
	}
	if cfg.UserEmail != "test@example.com" {
		t.Error("UserEmail should be set")
	}
	if cfg.UserFullName != "" {
		t.Error("UserFullName should be empty string")
	}
	if cfg.UserPublicName != "" {
		t.Error("UserPublicName should be empty string")
	}

	// User should still be considered configured if UserID is set
	if !cfg.HasUserConfigured() {
		t.Error("User should be considered configured when UserID is set")
	}
}

func TestShowCurrentSettings_WithUser(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "With user configured",
			config: &config.Config{
				AccessToken:  "test-token",
				WorkspaceID:  "ws_test",
				HelpCenterID: "hc_test",
				UserID:       "u_test",
				UserEmail:    "test@example.com",
				UserFullName: "Test User",
			},
		},
		{
			name: "Without user configured",
			config: &config.Config{
				AccessToken:  "test-token",
				WorkspaceID:  "ws_test",
				HelpCenterID: "hc_test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify showCurrentSettings doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("showCurrentSettings panicked: %v", r)
				}
			}()

			showCurrentSettings(tt.config)
		})
	}
}

// Note: The following tests would require mocking the Plain API client
// and are documented here for manual/integration testing:
//
// TestUserCmd_List - Would test listing users via API
// - Requires live API or mock client
// - Should verify user list is fetched and displayed
// - Should test with different numbers of users
//
// TestUserCmd_WithEmail - Would test configuring user via email
// - Requires live API or mock client
// - Should verify user lookup by email
// - Should verify user info is saved to config
// - Should test error cases (user not found, API errors)
//
// TestUserCmd_InteractiveMode - Would test interactive email prompt
// - Requires input mocking
// - Should verify email is prompted and read correctly
// - Should test with various email formats
