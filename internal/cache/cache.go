package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CacheMetadata is the shared metadata structure for all caches
type CacheMetadata struct {
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
	TTLHours  int       `json:"ttl_hours"`
}

// LoadCache loads any cache from disk into the provided interface
// The interface v must be a pointer to the cache struct
func LoadCache(cachePath string, v interface{}) error {
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache file does not exist: %w", err)
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse cache file (may be corrupted): %w", err)
	}

	return nil
}

// SaveCache saves any cache to disk with current timestamp
// The cache is written atomically using a temp file and rename
func SaveCache(cachePath string, v interface{}) error {
	// Ensure cache directory exists
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Marshal cache to JSON with indentation for readability
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	// Write to temporary file first
	tempPath := cachePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Atomically rename temp file to actual cache file
	if err := os.Rename(tempPath, cachePath); err != nil {
		os.Remove(tempPath) // Clean up temp file on error
		return fmt.Errorf("failed to save cache file: %w", err)
	}

	return nil
}

// IsFresh checks if cache metadata is within TTL
func IsFresh(updatedAt time.Time, ttlHours int) bool {
	if ttlHours <= 0 {
		return false
	}
	expiresAt := updatedAt.Add(time.Duration(ttlHours) * time.Hour)
	return time.Now().Before(expiresAt)
}

// getCacheDir returns the cache directory path (~/.config/plain-cli/)
func getCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "plain-cli"), nil
}

// GetCachePath returns the full path to a cache file
// This is a variable so it can be overridden in tests
var GetCachePath = func(filename string) (string, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, filename), nil
}
