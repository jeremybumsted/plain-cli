package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

func TestSaveLabelCache(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Create test label types
	labelTypes := []*mcp.LabelType{
		{ID: "labelType_1", Name: "Bug"},
		{ID: "labelType_2", Name: "Feature Request"},
	}

	// Save cache
	err := SaveLabelCache(labelTypes)
	if err != nil {
		t.Fatalf("SaveLabelCache failed: %v", err)
	}

	// Verify file exists
	cachePath := filepath.Join(tempDir, LabelCacheFilename)
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

func TestLoadLabelCache(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Create and save test cache
	labelTypes := []*mcp.LabelType{
		{ID: "labelType_1", Name: "Bug"},
		{ID: "labelType_2", Name: "Feature Request"},
	}
	if err := SaveLabelCache(labelTypes); err != nil {
		t.Fatalf("SaveLabelCache failed: %v", err)
	}

	// Load cache
	cache, err := LoadLabelCache()
	if err != nil {
		t.Fatalf("LoadLabelCache failed: %v", err)
	}

	// Verify contents
	if cache.Version != LabelCacheVersion {
		t.Errorf("Expected version %d, got %d", LabelCacheVersion, cache.Version)
	}
	if cache.TTLHours != LabelCacheTTLHours {
		t.Errorf("Expected TTL %d, got %d", LabelCacheTTLHours, cache.TTLHours)
	}
	if len(cache.LabelTypes) != 2 {
		t.Errorf("Expected 2 label types, got %d", len(cache.LabelTypes))
	}
	if cache.LabelTypes[0].Name != "Bug" {
		t.Errorf("Expected first label 'Bug', got '%s'", cache.LabelTypes[0].Name)
	}
}

func TestLoadLabelCache_NotExists(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Try to load non-existent cache
	_, err := LoadLabelCache()
	if err == nil {
		t.Fatal("Expected error for non-existent cache, got nil")
	}
}

func TestLabelCache_IsFresh(t *testing.T) {
	// Fresh cache (1 hour old)
	freshCache := &LabelCache{
		CacheMetadata: CacheMetadata{
			UpdatedAt: time.Now().Add(-1 * time.Hour),
			TTLHours:  24,
		},
	}
	if !freshCache.IsFresh() {
		t.Error("Expected cache to be fresh (1h old, 24h TTL)")
	}

	// Stale cache (25 hours old)
	staleCache := &LabelCache{
		CacheMetadata: CacheMetadata{
			UpdatedAt: time.Now().Add(-25 * time.Hour),
			TTLHours:  24,
		},
	}
	if staleCache.IsFresh() {
		t.Error("Expected cache to be stale (25h old, 24h TTL)")
	}
}

func TestLabelCache_GetLabelTypeByID(t *testing.T) {
	cache := &LabelCache{
		LabelTypes: []*mcp.LabelType{
			{ID: "labelType_1", Name: "Bug"},
			{ID: "labelType_2", Name: "Feature"},
		},
	}

	// Test found
	result := cache.GetLabelTypeByID("labelType_1")
	if result == nil {
		t.Fatal("Expected to find label type")
	}
	if result.Name != "Bug" {
		t.Errorf("Expected name 'Bug', got '%s'", result.Name)
	}

	// Test not found
	result = cache.GetLabelTypeByID("labelType_999")
	if result != nil {
		t.Error("Expected nil for non-existent label type")
	}
}

func TestLabelCache_GetLabelTypeByName(t *testing.T) {
	cache := &LabelCache{
		LabelTypes: []*mcp.LabelType{
			{ID: "labelType_1", Name: "Bug"},
			{ID: "labelType_2", Name: "Feature Request"},
		},
	}

	tests := []struct {
		name     string
		expected string
	}{
		{"Bug", "labelType_1"},
		{"bug", "labelType_1"},           // Case insensitive
		{"BUG", "labelType_1"},           // Case insensitive
		{"Feature Request", "labelType_2"},
		{"feature request", "labelType_2"}, // Case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cache.GetLabelTypeByName(tt.name)
			if result == nil {
				t.Fatalf("Expected to find label type for '%s'", tt.name)
			}
			if result.ID != tt.expected {
				t.Errorf("Expected ID '%s', got '%s'", tt.expected, result.ID)
			}
		})
	}

	// Test not found
	result := cache.GetLabelTypeByName("NonExistent")
	if result != nil {
		t.Error("Expected nil for non-existent label type")
	}
}

func TestLabelCache_ResolveLabelIdentifier(t *testing.T) {
	cache := &LabelCache{
		LabelTypes: []*mcp.LabelType{
			{ID: "labelType_1", Name: "Bug"},
			{ID: "labelType_2", Name: "Feature Request"},
		},
	}

	tests := []struct {
		input       string
		expectedID  string
		expectError bool
	}{
		// Direct IDs
		{"labelType_1", "labelType_1", false},
		{"labelType_2", "labelType_2", false},
		{"labelType_999", "", true}, // Non-existent ID

		// Names (case insensitive)
		{"Bug", "labelType_1", false},
		{"bug", "labelType_1", false},
		{"BUG", "labelType_1", false},
		{"Feature Request", "labelType_2", false},
		{"feature request", "labelType_2", false},

		// Non-existent names
		{"NonExistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			id, err := cache.ResolveLabelIdentifier(tt.input)
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

func TestLabelCache_ResolveLabelIdentifiers(t *testing.T) {
	cache := &LabelCache{
		LabelTypes: []*mcp.LabelType{
			{ID: "labelType_1", Name: "Bug"},
			{ID: "labelType_2", Name: "Feature Request"},
		},
	}

	// Test successful resolution
	ids, err := cache.ResolveLabelIdentifiers([]string{"Bug", "labelType_2", "feature request"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := []string{"labelType_1", "labelType_2", "labelType_2"}
	if len(ids) != len(expected) {
		t.Fatalf("Expected %d IDs, got %d", len(expected), len(ids))
	}
	for i, id := range ids {
		if id != expected[i] {
			t.Errorf("Expected ID[%d] '%s', got '%s'", i, expected[i], id)
		}
	}

	// Test with error (non-existent label)
	_, err = cache.ResolveLabelIdentifiers([]string{"Bug", "NonExistent"})
	if err == nil {
		t.Error("Expected error for non-existent label, got nil")
	}
}

func TestLabelCache_GetLabelNames(t *testing.T) {
	cache := &LabelCache{
		LabelTypes: []*mcp.LabelType{
			{ID: "labelType_1", Name: "Bug"},
			{ID: "labelType_2", Name: "Feature Request"},
		},
	}

	// Test with all valid IDs
	names := cache.GetLabelNames([]string{"labelType_1", "labelType_2"})
	expected := []string{"Bug", "Feature Request"}
	if len(names) != len(expected) {
		t.Fatalf("Expected %d names, got %d", len(expected), len(names))
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("Expected name[%d] '%s', got '%s'", i, expected[i], name)
		}
	}

	// Test with non-existent ID (should fallback to ID)
	names = cache.GetLabelNames([]string{"labelType_1", "labelType_999"})
	if len(names) != 2 {
		t.Fatalf("Expected 2 names, got %d", len(names))
	}
	if names[0] != "Bug" {
		t.Errorf("Expected name[0] 'Bug', got '%s'", names[0])
	}
	if names[1] != "labelType_999" {
		t.Errorf("Expected name[1] 'labelType_999', got '%s'", names[1])
	}
}

func TestLabelCache_RoundTrip(t *testing.T) {
	// Use temp directory for test
	tempDir := t.TempDir()
	oldGetCachePath := GetCachePath
	GetCachePath = func(filename string) (string, error) {
		return filepath.Join(tempDir, filename), nil
	}
	defer func() { GetCachePath = oldGetCachePath }()

	// Create test label types
	original := []*mcp.LabelType{
		{ID: "labelType_1", Name: "Bug", Icon: "🐛", Color: "red", IsArchived: false},
		{ID: "labelType_2", Name: "Feature", Icon: "✨", Color: "blue", IsArchived: false},
		{ID: "labelType_3", Name: "Archived", Icon: "📦", Color: "gray", IsArchived: true},
	}

	// Save
	if err := SaveLabelCache(original); err != nil {
		t.Fatalf("SaveLabelCache failed: %v", err)
	}

	// Load
	cache, err := LoadLabelCache()
	if err != nil {
		t.Fatalf("LoadLabelCache failed: %v", err)
	}

	// Verify
	if len(cache.LabelTypes) != len(original) {
		t.Fatalf("Expected %d label types, got %d", len(original), len(cache.LabelTypes))
	}
	for i, lt := range original {
		loaded := cache.LabelTypes[i]
		if loaded.ID != lt.ID {
			t.Errorf("ID mismatch at index %d: expected '%s', got '%s'", i, lt.ID, loaded.ID)
		}
		if loaded.Name != lt.Name {
			t.Errorf("Name mismatch at index %d: expected '%s', got '%s'", i, lt.Name, loaded.Name)
		}
		if loaded.Icon != lt.Icon {
			t.Errorf("Icon mismatch at index %d: expected '%s', got '%s'", i, lt.Icon, loaded.Icon)
		}
		if loaded.Color != lt.Color {
			t.Errorf("Color mismatch at index %d: expected '%s', got '%s'", i, lt.Color, loaded.Color)
		}
		if loaded.IsArchived != lt.IsArchived {
			t.Errorf("IsArchived mismatch at index %d: expected %v, got %v", i, lt.IsArchived, loaded.IsArchived)
		}
	}
}
