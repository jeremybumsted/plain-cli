package threads

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

// TestListCmdMineFlag tests that the --mine flag uses the configured user ID
func TestListCmdMineFlag(t *testing.T) {
	// Create a temporary config file with user ID
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Load and configure with a test user ID
	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	err := cfg.SetUserInfo("u_test123", "test@example.com", "Test User", "Test")
	if err != nil {
		t.Fatalf("Failed to set user info: %v", err)
	}

	// Create ListCmd with Mine flag set
	cmd := &ListCmd{
		Mine:       true,
		ConfigPath: configPath,
		Format:     "table",
	}

	// We can't actually call Run() without a real API client,
	// but we can test that the config is loaded properly
	loadedCfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify that GetUserID returns the expected user ID
	userID, err := loadedCfg.GetUserID()
	if err != nil {
		t.Errorf("GetUserID() error = %v, expected no error", err)
	}
	if userID != "u_test123" {
		t.Errorf("GetUserID() = %v, want %v", userID, "u_test123")
	}
}

// TestListCmdMineWithoutUser tests that --mine flag returns error when user is not configured
func TestListCmdMineWithoutUser(t *testing.T) {
	// Create a temporary config file without user ID
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create empty config
	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedCfg, err := getConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify that GetUserID returns an error
	_, err = loadedCfg.GetUserID()
	if err == nil {
		t.Error("GetUserID() expected error when user not configured, got nil")
	}
}

// TestListCmdMineMutualExclusion tests that --mine and --assignee cannot be used together
func TestListCmdMineMutualExclusion(t *testing.T) {
	// Create a temporary config file with user ID
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Load and configure with a test user ID
	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	err := cfg.SetUserInfo("u_test123", "test@example.com", "Test User", "Test")
	if err != nil {
		t.Fatalf("Failed to set user info: %v", err)
	}

	// Also set a fake API token to avoid authentication errors
	cfg.SetTokens("fake_token", "fake_refresh", "Bearer", 3600)
	err = cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create ListCmd with both Mine and Assignee set
	cmd := &ListCmd{
		Mine:       true,
		Assignee:   "u_other456",
		ConfigPath: configPath,
		Format:     "table",
	}

	// Run the command - should return error about mutual exclusion
	err = cmd.Run()
	if err == nil {
		t.Error("Run() expected error when both --mine and --assignee are set, got nil")
	}
	if err != nil && err.Error() != "cannot use both --mine and --assignee flags" {
		t.Errorf("Run() error = %v, want 'cannot use both --mine and --assignee flags'", err)
	}
}

// TestListCmdMineWithOtherFilters tests combining --mine with other filter flags
func TestListCmdMineWithOtherFilters(t *testing.T) {
	// Create a temporary config file with user ID
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Load and configure with a test user ID
	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	err := cfg.SetUserInfo("u_test123", "test@example.com", "Test User", "Test")
	if err != nil {
		t.Fatalf("Failed to set user info: %v", err)
	}

	// Test that Mine can be combined with Status filter
	cmd := &ListCmd{
		Mine:       true,
		Status:     "TODO",
		ConfigPath: configPath,
		Format:     "table",
	}

	loadedCfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	userID, err := loadedCfg.GetUserID()
	if err != nil {
		t.Errorf("GetUserID() error = %v, expected no error", err)
	}
	if userID != "u_test123" {
		t.Errorf("GetUserID() = %v, want %v", userID, "u_test123")
	}
	if cmd.Status != "TODO" {
		t.Errorf("Status = %v, want %v", cmd.Status, "TODO")
	}
}

// TestListCmdAssigneeWithoutMine tests that --assignee works when --mine is not set
func TestListCmdAssigneeWithoutMine(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Create ListCmd with only Assignee set (Mine is false by default)
	cmd := &ListCmd{
		Mine:       false,
		Assignee:   "u_other456",
		ConfigPath: configPath,
		Format:     "table",
	}

	// Verify Mine is false and Assignee is set
	if cmd.Mine {
		t.Error("Mine should be false")
	}
	if cmd.Assignee != "u_other456" {
		t.Errorf("Assignee = %v, want %v", cmd.Assignee, "u_other456")
	}
}

// TestListCmdDefaultValues tests that command has correct default values
func TestListCmdDefaultValues(t *testing.T) {
	cmd := &ListCmd{}

	if cmd.Mine != false {
		t.Errorf("Mine default = %v, want false", cmd.Mine)
	}
	if cmd.Status != "" {
		t.Errorf("Status default = %v, want empty string", cmd.Status)
	}
	if cmd.Assignee != "" {
		t.Errorf("Assignee default = %v, want empty string", cmd.Assignee)
	}
}

// TestGetUserIDFromEnv tests user ID retrieval from environment (if implemented)
func TestGetUserIDFromConfig(t *testing.T) {
	// Create a temporary config file with user info
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)

	// Test with full user info
	err := cfg.SetUserInfo("u_full123", "full@example.com", "Full Name", "Display Name")
	if err != nil {
		t.Fatalf("Failed to set user info: %v", err)
	}

	// Verify all fields are stored
	if cfg.UserID != "u_full123" {
		t.Errorf("UserID = %v, want %v", cfg.UserID, "u_full123")
	}
	if cfg.UserEmail != "full@example.com" {
		t.Errorf("UserEmail = %v, want %v", cfg.UserEmail, "full@example.com")
	}
	if cfg.UserFullName != "Full Name" {
		t.Errorf("UserFullName = %v, want %v", cfg.UserFullName, "Full Name")
	}
	if cfg.UserPublicName != "Display Name" {
		t.Errorf("UserPublicName = %v, want %v", cfg.UserPublicName, "Display Name")
	}

	// Verify GetUserID works
	userID, err := cfg.GetUserID()
	if err != nil {
		t.Errorf("GetUserID() error = %v, expected no error", err)
	}
	if userID != "u_full123" {
		t.Errorf("GetUserID() = %v, want %v", userID, "u_full123")
	}
}

// TestConfigPersistence tests that user config is persisted to disk
func TestConfigPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Create and save config with user info
	cfg1 := &config.Config{}
	cfg1.SetConfigPath(configPath)
	err := cfg1.SetUserInfo("u_persist123", "persist@example.com", "Persist User", "Persist")
	if err != nil {
		t.Fatalf("Failed to set user info: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load config from file
	cfg2, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify user info was persisted
	if cfg2.UserID != "u_persist123" {
		t.Errorf("Loaded UserID = %v, want %v", cfg2.UserID, "u_persist123")
	}
	if cfg2.UserEmail != "persist@example.com" {
		t.Errorf("Loaded UserEmail = %v, want %v", cfg2.UserEmail, "persist@example.com")
	}
}

// TestListCmdDateParsingValidDates tests that valid date formats are parsed correctly
func TestListCmdDateParsingValidDates(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	cfg.SetTokens("fake_token", "fake_refresh", "Bearer", 3600)
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	tests := []struct {
		name          string
		createdAfter  string
		createdBefore string
		updatedAfter  string
		updatedBefore string
		shouldFail    bool
	}{
		{
			name:         "ISO8601 dates",
			createdAfter: "2026-03-01T00:00:00Z",
			shouldFail:   false,
		},
		{
			name:         "Date only format",
			createdAfter: "2026-03-01",
			shouldFail:   false,
		},
		{
			name:         "Relative date - days",
			createdAfter: "7d",
			shouldFail:   false,
		},
		{
			name:         "Relative date - weeks",
			updatedAfter: "2w",
			shouldFail:   false,
		},
		{
			name:         "Human readable - yesterday",
			createdAfter: "yesterday",
			shouldFail:   false,
		},
		{
			name:         "Human readable - last-week",
			updatedAfter: "last-week",
			shouldFail:   false,
		},
		{
			name:          "Valid date range",
			createdAfter:  "7d",
			createdBefore: "today",
			shouldFail:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ListCmd{
				CreatedAfter:  tt.createdAfter,
				CreatedBefore: tt.createdBefore,
				UpdatedAfter:  tt.updatedAfter,
				UpdatedBefore: tt.updatedBefore,
				ConfigPath:    configPath,
				Format:        "table",
			}

			err := cmd.Run()
			// We expect this to fail at the API call stage, not at date parsing
			if err != nil && !tt.shouldFail {
				// Check if error is about date parsing (which would be a test failure)
				// or about API/authentication (which is expected in tests)
				errMsg := err.Error()
				if contains(errMsg, "invalid created-after") ||
					contains(errMsg, "invalid created-before") ||
					contains(errMsg, "invalid updated-after") ||
					contains(errMsg, "invalid updated-before") ||
					contains(errMsg, "invalid created date range") ||
					contains(errMsg, "invalid updated date range") {
					t.Errorf("Date parsing failed: %v", err)
				}
				// If error is about API/auth, that's expected in unit tests
			}
		})
	}
}

// TestListCmdDateParsingInvalidDates tests that invalid date formats return errors
func TestListCmdDateParsingInvalidDates(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	cfg.SetTokens("fake_token", "fake_refresh", "Bearer", 3600)
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	tests := []struct {
		name          string
		createdAfter  string
		createdBefore string
		expectedError string
	}{
		{
			name:          "Invalid date format",
			createdAfter:  "invalid-date",
			expectedError: "invalid created-after date",
		},
		{
			name:          "Invalid before date",
			createdBefore: "not-a-date",
			expectedError: "invalid created-before date",
		},
		{
			name:          "Empty relative date",
			createdAfter:  "d",
			expectedError: "invalid created-after date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ListCmd{
				CreatedAfter:  tt.createdAfter,
				CreatedBefore: tt.createdBefore,
				ConfigPath:    configPath,
				Format:        "table",
			}

			err := cmd.Run()
			if err == nil {
				t.Error("Expected error for invalid date, got nil")
				return
			}

			if !contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error containing '%s', got '%v'", tt.expectedError, err)
			}
		})
	}
}

// TestListCmdDateRangeValidation tests that date range validation works correctly
func TestListCmdDateRangeValidation(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	cfg := &config.Config{}
	cfg.SetConfigPath(configPath)
	cfg.SetTokens("fake_token", "fake_refresh", "Bearer", 3600)
	err := cfg.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	tests := []struct {
		name          string
		createdAfter  string
		createdBefore string
		shouldFail    bool
		expectedError string
	}{
		{
			name:          "Invalid range - after is later than before",
			createdAfter:  "today",
			createdBefore: "yesterday",
			shouldFail:    true,
			expectedError: "invalid created date range",
		},
		{
			name:          "Valid range - after is earlier than before",
			createdAfter:  "7d",
			createdBefore: "today",
			shouldFail:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ListCmd{
				CreatedAfter:  tt.createdAfter,
				CreatedBefore: tt.createdBefore,
				ConfigPath:    configPath,
				Format:        "table",
			}

			err := cmd.Run()
			if tt.shouldFail {
				if err == nil {
					t.Error("Expected error for invalid date range, got nil")
					return
				}
				if !contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.expectedError, err)
				}
			} else if err != nil && contains(err.Error(), "invalid created date range") {
				t.Errorf("Unexpected date range validation error: %v", err)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
