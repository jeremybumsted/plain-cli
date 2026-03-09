package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

func TestSaveFieldCache(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Create test field schemas
	schemas := []*mcp.ThreadFieldSchema{
		{ID: "field_1", Key: "priority", Label: "Priority", Type: "STRING"},
		{ID: "field_2", Key: "category", Label: "Category", Type: "ENUM"},
	}

	// Save cache
	err := SaveFieldCache(schemas)
	if err != nil {
		t.Fatalf("SaveFieldCache failed: %v", err)
	}

	// Verify file exists
	cachePath := filepath.Join(tempDir, FieldCacheFilename)
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatal("Cache file was not created")
	}

	// Verify file permissions
	info, err := os.Stat(cachePath)
	if err != nil {
		t.Fatalf("Failed to stat cache file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestLoadFieldCache(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Create and save test cache
	schemas := []*mcp.ThreadFieldSchema{
		{ID: "field_1", Key: "priority", Label: "Priority", Type: "STRING"},
		{ID: "field_2", Key: "category", Label: "Category", Type: "ENUM"},
	}
	if err := SaveFieldCache(schemas); err != nil {
		t.Fatalf("SaveFieldCache failed: %v", err)
	}

	// Load cache
	cache, err := LoadFieldCache()
	if err != nil {
		t.Fatalf("LoadFieldCache failed: %v", err)
	}

	// Verify contents
	if cache.Version != FieldCacheVersion {
		t.Errorf("Expected version %d, got %d", FieldCacheVersion, cache.Version)
	}
	if cache.TTLHours != FieldCacheTTLHours {
		t.Errorf("Expected TTL %d, got %d", FieldCacheTTLHours, cache.TTLHours)
	}
	if len(cache.FieldSchemas) != 2 {
		t.Errorf("Expected 2 field schemas, got %d", len(cache.FieldSchemas))
	}
	if cache.FieldSchemas[0].Key != "priority" {
		t.Errorf("Expected first field key 'priority', got '%s'", cache.FieldSchemas[0].Key)
	}
}

func TestLoadFieldCache_NotExists(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Try to load non-existent cache
	_, err := LoadFieldCache()
	if err == nil {
		t.Fatal("Expected error for non-existent cache, got nil")
	}
}

func TestFieldCache_IsFresh(t *testing.T) {
	// Fresh cache (1 hour old)
	freshCache := &FieldCache{
		CacheMetadata: CacheMetadata{
			UpdatedAt: time.Now().Add(-1 * time.Hour),
			TTLHours:  24,
		},
	}
	if !freshCache.IsFresh() {
		t.Error("Expected cache to be fresh (1h old, 24h TTL)")
	}

	// Stale cache (25 hours old)
	staleCache := &FieldCache{
		CacheMetadata: CacheMetadata{
			UpdatedAt: time.Now().Add(-25 * time.Hour),
			TTLHours:  24,
		},
	}
	if staleCache.IsFresh() {
		t.Error("Expected cache to be stale (25h old, 24h TTL)")
	}
}

func TestFieldCache_GetFieldSchemaByID(t *testing.T) {
	cache := &FieldCache{
		FieldSchemas: []*mcp.ThreadFieldSchema{
			{ID: "field_1", Key: "priority", Label: "Priority"},
			{ID: "field_2", Key: "category", Label: "Category"},
		},
	}

	// Test found
	result := cache.GetFieldSchemaByID("field_1")
	if result == nil {
		t.Fatal("Expected to find field schema")
	}
	if result.Key != "priority" {
		t.Errorf("Expected key 'priority', got '%s'", result.Key)
	}

	// Test not found
	result = cache.GetFieldSchemaByID("field_999")
	if result != nil {
		t.Error("Expected nil for non-existent field schema")
	}
}

func TestFieldCache_GetFieldSchemaByKey(t *testing.T) {
	cache := &FieldCache{
		FieldSchemas: []*mcp.ThreadFieldSchema{
			{ID: "field_1", Key: "priority", Label: "Priority"},
			{ID: "field_2", Key: "customerType", Label: "Customer Type"},
		},
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"priority", "field_1"},
		{"Priority", "field_1"},           // Case insensitive
		{"PRIORITY", "field_1"},           // Case insensitive
		{"customerType", "field_2"},
		{"customertype", "field_2"},       // Case insensitive
		{"CUSTOMERTYPE", "field_2"},       // Case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := cache.GetFieldSchemaByKey(tt.key)
			if result == nil {
				t.Fatalf("Expected to find field schema for '%s'", tt.key)
			}
			if result.ID != tt.expected {
				t.Errorf("Expected ID '%s', got '%s'", tt.expected, result.ID)
			}
		})
	}

	// Test not found
	result := cache.GetFieldSchemaByKey("nonExistent")
	if result != nil {
		t.Error("Expected nil for non-existent field schema")
	}
}

func TestFieldCache_ResolveFieldIdentifier(t *testing.T) {
	cache := &FieldCache{
		FieldSchemas: []*mcp.ThreadFieldSchema{
			{ID: "threadFieldSchema_1", Key: "priority", Label: "Priority"},
			{ID: "threadFieldSchema_2", Key: "customerType", Label: "Customer Type"},
		},
	}

	tests := []struct {
		input       string
		expectedID  string
		expectError bool
	}{
		// Direct IDs
		{"threadFieldSchema_1", "threadFieldSchema_1", false},
		{"threadFieldSchema_2", "threadFieldSchema_2", false},
		{"threadFieldSchema_999", "", true}, // Non-existent ID

		// Keys (case insensitive)
		{"priority", "threadFieldSchema_1", false},
		{"Priority", "threadFieldSchema_1", false},
		{"PRIORITY", "threadFieldSchema_1", false},
		{"customerType", "threadFieldSchema_2", false},
		{"customertype", "threadFieldSchema_2", false},

		// Non-existent keys
		{"nonExistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			id, err := cache.ResolveFieldIdentifier(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input '%s', got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input '%s': %v", tt.input, err)
				}
				if id != tt.expectedID {
					t.Errorf("Expected ID '%s', got '%s'", tt.expectedID, id)
				}
			}
		})
	}
}

func TestFieldCache_RoundTrip(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Create test field schemas
	original := []*mcp.ThreadFieldSchema{
		{
			ID:                  "field_1",
			Key:                 "priority",
			Label:               "Priority",
			Type:                "STRING",
			Description:         "Task priority",
			IsRequired:          true,
			IsAiAutoFillEnabled: false,
		},
		{
			ID:                  "field_2",
			Key:                 "category",
			Label:               "Category",
			Type:                "ENUM",
			Description:         "Issue category",
			EnumValues:          []string{"Bug", "Feature", "Question"},
			IsRequired:          false,
			IsAiAutoFillEnabled: true,
		},
	}

	// Save
	if err := SaveFieldCache(original); err != nil {
		t.Fatalf("SaveFieldCache failed: %v", err)
	}

	// Load
	cache, err := LoadFieldCache()
	if err != nil {
		t.Fatalf("LoadFieldCache failed: %v", err)
	}

	// Verify
	if len(cache.FieldSchemas) != len(original) {
		t.Fatalf("Expected %d field schemas, got %d", len(original), len(cache.FieldSchemas))
	}
	for i, fs := range original {
		loaded := cache.FieldSchemas[i]
		if loaded.ID != fs.ID {
			t.Errorf("ID mismatch at index %d: expected '%s', got '%s'", i, fs.ID, loaded.ID)
		}
		if loaded.Key != fs.Key {
			t.Errorf("Key mismatch at index %d: expected '%s', got '%s'", i, fs.Key, loaded.Key)
		}
		if loaded.Label != fs.Label {
			t.Errorf("Label mismatch at index %d: expected '%s', got '%s'", i, fs.Label, loaded.Label)
		}
		if loaded.Type != fs.Type {
			t.Errorf("Type mismatch at index %d: expected '%s', got '%s'", i, fs.Type, loaded.Type)
		}
		if loaded.Description != fs.Description {
			t.Errorf("Description mismatch at index %d: expected '%s', got '%s'", i, fs.Description, loaded.Description)
		}
		if loaded.IsRequired != fs.IsRequired {
			t.Errorf("IsRequired mismatch at index %d: expected %v, got %v", i, fs.IsRequired, loaded.IsRequired)
		}
		if loaded.IsAiAutoFillEnabled != fs.IsAiAutoFillEnabled {
			t.Errorf("IsAiAutoFillEnabled mismatch at index %d: expected %v, got %v", i, fs.IsAiAutoFillEnabled, loaded.IsAiAutoFillEnabled)
		}
		if len(loaded.EnumValues) != len(fs.EnumValues) {
			t.Errorf("EnumValues length mismatch at index %d: expected %d, got %d", i, len(fs.EnumValues), len(loaded.EnumValues))
		}
	}
}
