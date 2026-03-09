package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCache is a simple test cache structure
type TestCache struct {
	CacheMetadata
	Data []string `json:"data"`
}

func TestLoadCache_Success(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "test-cache.json")

	// Write test cache file
	testData := `{
		"version": 1,
		"updated_at": "2026-03-09T10:00:00Z",
		"ttl_hours": 24,
		"data": ["item1", "item2"]
	}`
	if err := os.WriteFile(cachePath, []byte(testData), 0600); err != nil {
		t.Fatalf("Failed to write test cache: %v", err)
	}

	// Load cache
	var cache TestCache
	err := LoadCache(cachePath, &cache)
	if err != nil {
		t.Fatalf("LoadCache failed: %v", err)
	}

	// Verify contents
	if cache.Version != 1 {
		t.Errorf("Expected version 1, got %d", cache.Version)
	}
	if cache.TTLHours != 24 {
		t.Errorf("Expected TTL 24, got %d", cache.TTLHours)
	}
	if len(cache.Data) != 2 {
		t.Errorf("Expected 2 data items, got %d", len(cache.Data))
	}
	if cache.Data[0] != "item1" {
		t.Errorf("Expected data[0] 'item1', got '%s'", cache.Data[0])
	}
}

func TestLoadCache_FileNotExists(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "nonexistent.json")

	var cache TestCache
	err := LoadCache(cachePath, &cache)
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}
	// Check if the error message contains "does not exist" since error is wrapped
	if err.Error() == "" || len(err.Error()) == 0 {
		t.Errorf("Expected non-empty error message, got: %v", err)
	}
}

func TestLoadCache_CorruptedFile(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "corrupted.json")

	// Write invalid JSON
	if err := os.WriteFile(cachePath, []byte("not valid json {"), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	var cache TestCache
	err := LoadCache(cachePath, &cache)
	if err == nil {
		t.Fatal("Expected error for corrupted file, got nil")
	}
}

func TestSaveCache_Success(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "test-cache.json")

	// Create test cache
	cache := TestCache{
		CacheMetadata: CacheMetadata{
			Version:   1,
			UpdatedAt: time.Now(),
			TTLHours:  24,
		},
		Data: []string{"item1", "item2"},
	}

	// Save cache
	err := SaveCache(cachePath, &cache)
	if err != nil {
		t.Fatalf("SaveCache failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatal("Cache file was not created")
	}

	// Verify file permissions
	info, err := os.Stat(cachePath)
	if err != nil {
		t.Fatalf("Failed to stat cache file: %v", err)
	}
	mode := info.Mode()
	if mode.Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", mode.Perm())
	}

	// Load and verify contents
	var loaded TestCache
	if err := LoadCache(cachePath, &loaded); err != nil {
		t.Fatalf("Failed to load saved cache: %v", err)
	}
	if loaded.Version != cache.Version {
		t.Errorf("Version mismatch: expected %d, got %d", cache.Version, loaded.Version)
	}
	if len(loaded.Data) != len(cache.Data) {
		t.Errorf("Data length mismatch: expected %d, got %d", len(cache.Data), len(loaded.Data))
	}
}

func TestSaveCache_CreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "subdir", "nested", "test-cache.json")

	cache := TestCache{
		CacheMetadata: CacheMetadata{
			Version:   1,
			UpdatedAt: time.Now(),
			TTLHours:  24,
		},
		Data: []string{"test"},
	}

	// Save cache (should create nested directories)
	err := SaveCache(cachePath, &cache)
	if err != nil {
		t.Fatalf("SaveCache failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Fatal("Cache file was not created in nested directory")
	}
}

func TestIsFresh_Fresh(t *testing.T) {
	// Cache updated 1 hour ago with 24 hour TTL
	updatedAt := time.Now().Add(-1 * time.Hour)
	if !IsFresh(updatedAt, 24) {
		t.Error("Expected cache to be fresh (1h old, 24h TTL)")
	}
}

func TestIsFresh_Stale(t *testing.T) {
	// Cache updated 25 hours ago with 24 hour TTL
	updatedAt := time.Now().Add(-25 * time.Hour)
	if IsFresh(updatedAt, 24) {
		t.Error("Expected cache to be stale (25h old, 24h TTL)")
	}
}

func TestIsFresh_JustExpired(t *testing.T) {
	// Cache updated exactly 24 hours + 1 second ago
	updatedAt := time.Now().Add(-24*time.Hour - 1*time.Second)
	if IsFresh(updatedAt, 24) {
		t.Error("Expected cache to be stale (just expired)")
	}
}

func TestIsFresh_ZeroTTL(t *testing.T) {
	updatedAt := time.Now()
	if IsFresh(updatedAt, 0) {
		t.Error("Expected cache to be stale with 0 TTL")
	}
}

func TestIsFresh_NegativeTTL(t *testing.T) {
	updatedAt := time.Now()
	if IsFresh(updatedAt, -1) {
		t.Error("Expected cache to be stale with negative TTL")
	}
}

func TestGetCacheDir(t *testing.T) {
	cacheDir, err := getCacheDir()
	if err != nil {
		t.Fatalf("getCacheDir failed: %v", err)
	}

	// Verify it ends with .config/plain-cli
	if !filepath.IsAbs(cacheDir) {
		t.Error("Expected absolute path")
	}
	if filepath.Base(cacheDir) != "plain-cli" {
		t.Errorf("Expected cache dir to end with 'plain-cli', got: %s", cacheDir)
	}
}

func TestGetCachePath(t *testing.T) {
	cachePath, err := GetCachePath("test-cache.json")
	if err != nil {
		t.Fatalf("GetCachePath failed: %v", err)
	}

	// Verify it's absolute
	if !filepath.IsAbs(cachePath) {
		t.Error("Expected absolute path")
	}

	// Verify it ends with the filename
	if filepath.Base(cachePath) != "test-cache.json" {
		t.Errorf("Expected path to end with 'test-cache.json', got: %s", cachePath)
	}

	// Verify it contains .config/plain-cli
	if !filepath.IsAbs(cachePath) {
		t.Error("Expected absolute path")
	}
}

func TestSaveAndLoadRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "roundtrip.json")

	// Create original cache
	original := TestCache{
		CacheMetadata: CacheMetadata{
			Version:   2,
			UpdatedAt: time.Now().Truncate(time.Second), // Truncate for comparison
			TTLHours:  48,
		},
		Data: []string{"alpha", "beta", "gamma"},
	}

	// Save
	if err := SaveCache(cachePath, &original); err != nil {
		t.Fatalf("SaveCache failed: %v", err)
	}

	// Load
	var loaded TestCache
	if err := LoadCache(cachePath, &loaded); err != nil {
		t.Fatalf("LoadCache failed: %v", err)
	}

	// Compare
	if loaded.Version != original.Version {
		t.Errorf("Version mismatch: expected %d, got %d", original.Version, loaded.Version)
	}
	if loaded.TTLHours != original.TTLHours {
		t.Errorf("TTL mismatch: expected %d, got %d", original.TTLHours, loaded.TTLHours)
	}
	if !loaded.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: expected %v, got %v", original.UpdatedAt, loaded.UpdatedAt)
	}
	if len(loaded.Data) != len(original.Data) {
		t.Errorf("Data length mismatch: expected %d, got %d", len(original.Data), len(loaded.Data))
	}
	for i, v := range original.Data {
		if loaded.Data[i] != v {
			t.Errorf("Data[%d] mismatch: expected '%s', got '%s'", i, v, loaded.Data[i])
		}
	}
}
