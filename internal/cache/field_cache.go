package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

const (
	// FieldCacheFilename is the filename for the field cache
	FieldCacheFilename = "field-cache.json"
	// FieldCacheTTLHours is the default TTL for field cache in hours
	FieldCacheTTLHours = 24
	// FieldCacheVersion is the current version of the field cache format
	FieldCacheVersion = 1
)

// FieldCache stores cached thread field schemas
type FieldCache struct {
	CacheMetadata
	FieldSchemas []*plain.ThreadFieldSchema `json:"field_schemas"`
}

// LoadFieldCache loads the field cache from disk
// Returns the cache if it exists, or an error if the file doesn't exist or is corrupted
func LoadFieldCache() (*FieldCache, error) {
	cachePath, err := GetCachePath(FieldCacheFilename)
	if err != nil {
		return nil, err
	}

	var cache FieldCache
	if err := LoadCache(cachePath, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// SaveFieldCache saves field schemas to the cache file
// This creates a new cache with current timestamp and default TTL
func SaveFieldCache(schemas []*plain.ThreadFieldSchema) error {
	cache := &FieldCache{
		CacheMetadata: CacheMetadata{
			Version:   FieldCacheVersion,
			UpdatedAt: time.Now(),
			TTLHours:  FieldCacheTTLHours,
		},
		FieldSchemas: schemas,
	}

	cachePath, err := GetCachePath(FieldCacheFilename)
	if err != nil {
		return err
	}

	return SaveCache(cachePath, cache)
}

// IsFresh checks if the field cache is within TTL
func (c *FieldCache) IsFresh() bool {
	return IsFresh(c.UpdatedAt, c.TTLHours)
}

// GetFieldSchemaByID finds a field schema by ID in the cache
// Returns nil if not found
func (c *FieldCache) GetFieldSchemaByID(id string) *plain.ThreadFieldSchema {
	for _, fs := range c.FieldSchemas {
		if fs.ID == id {
			return fs
		}
	}
	return nil
}

// GetFieldSchemaByKey finds a field schema by key (case-insensitive) in the cache
// Returns nil if not found
func (c *FieldCache) GetFieldSchemaByKey(key string) *plain.ThreadFieldSchema {
	lowerKey := strings.ToLower(key)
	for _, fs := range c.FieldSchemas {
		if strings.ToLower(fs.Key) == lowerKey {
			return fs
		}
	}
	return nil
}

// ResolveFieldIdentifier resolves a field identifier (key or ID) to a field schema ID
// Supports both direct IDs (threadFieldSchema_xxx) and case-insensitive keys
// Returns the field schema ID and an error if not found
func (c *FieldCache) ResolveFieldIdentifier(identifier string) (string, error) {
	// If it starts with "threadFieldSchema_", assume it's already an ID
	if strings.HasPrefix(identifier, "threadFieldSchema_") {
		// Verify it exists in cache
		if c.GetFieldSchemaByID(identifier) != nil {
			return identifier, nil
		}
		return "", fmt.Errorf("field schema ID not found in cache: %s (try 'plain threads field refresh')", identifier)
	}

	// Otherwise, try to resolve as a key
	schema := c.GetFieldSchemaByKey(identifier)
	if schema == nil {
		return "", fmt.Errorf("field schema key not found in cache: %s (try 'plain threads field refresh')", identifier)
	}

	return schema.ID, nil
}
