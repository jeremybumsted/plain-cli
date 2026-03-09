package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

const (
	// LabelCacheFilename is the filename for the label cache
	LabelCacheFilename = "label-cache.json"
	// LabelCacheTTLHours is the default TTL for label cache in hours
	LabelCacheTTLHours = 24
	// LabelCacheVersion is the current version of the label cache format
	LabelCacheVersion = 1
)

// LabelCache stores cached label types
type LabelCache struct {
	CacheMetadata
	LabelTypes []*mcp.LabelType `json:"label_types"`
}

// LoadLabelCache loads the label cache from disk
// Returns the cache if it exists, or an error if the file doesn't exist or is corrupted
func LoadLabelCache() (*LabelCache, error) {
	cachePath, err := GetCachePath(LabelCacheFilename)
	if err != nil {
		return nil, err
	}

	var cache LabelCache
	if err := LoadCache(cachePath, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// SaveLabelCache saves label types to the cache file
// This creates a new cache with current timestamp and default TTL
func SaveLabelCache(labelTypes []*mcp.LabelType) error {
	cache := &LabelCache{
		CacheMetadata: CacheMetadata{
			Version:   LabelCacheVersion,
			UpdatedAt: time.Now(),
			TTLHours:  LabelCacheTTLHours,
		},
		LabelTypes: labelTypes,
	}

	cachePath, err := GetCachePath(LabelCacheFilename)
	if err != nil {
		return err
	}

	return SaveCache(cachePath, cache)
}

// IsFresh checks if the label cache is within TTL
func (c *LabelCache) IsFresh() bool {
	return IsFresh(c.UpdatedAt, c.TTLHours)
}

// GetLabelTypeByID finds a label type by ID in the cache
// Returns nil if not found
func (c *LabelCache) GetLabelTypeByID(id string) *mcp.LabelType {
	for _, lt := range c.LabelTypes {
		if lt.ID == id {
			return lt
		}
	}
	return nil
}

// GetLabelTypeByName finds a label type by name (case-insensitive) in the cache
// Returns nil if not found
func (c *LabelCache) GetLabelTypeByName(name string) *mcp.LabelType {
	lowerName := strings.ToLower(name)
	for _, lt := range c.LabelTypes {
		if strings.ToLower(lt.Name) == lowerName {
			return lt
		}
	}
	return nil
}

// ResolveLabelIdentifier resolves a label identifier (name or ID) to a label type ID
// Supports both direct IDs (labelType_xxx) and case-insensitive names
// Returns the label type ID and an error if not found
func (c *LabelCache) ResolveLabelIdentifier(identifier string) (string, error) {
	// If it starts with "labelType_", assume it's already an ID
	if strings.HasPrefix(identifier, "labelType_") {
		// Verify it exists in cache
		if c.GetLabelTypeByID(identifier) != nil {
			return identifier, nil
		}
		return "", fmt.Errorf("label type ID not found in cache: %s (try 'plain threads label refresh')", identifier)
	}

	// Otherwise, try to resolve as a name
	labelType := c.GetLabelTypeByName(identifier)
	if labelType == nil {
		return "", fmt.Errorf("label type name not found in cache: %s (try 'plain threads label refresh')", identifier)
	}

	return labelType.ID, nil
}

// ResolveLabelIdentifiers resolves multiple label identifiers to label type IDs
// Returns a slice of IDs and an error if any identifier could not be resolved
func (c *LabelCache) ResolveLabelIdentifiers(identifiers []string) ([]string, error) {
	ids := make([]string, 0, len(identifiers))
	for _, identifier := range identifiers {
		id, err := c.ResolveLabelIdentifier(identifier)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetLabelNames returns a slice of label names for the given label type IDs
// This is useful for displaying user-friendly confirmation prompts
func (c *LabelCache) GetLabelNames(labelTypeIDs []string) []string {
	names := make([]string, 0, len(labelTypeIDs))
	for _, id := range labelTypeIDs {
		labelType := c.GetLabelTypeByID(id)
		if labelType != nil {
			names = append(names, labelType.Name)
		} else {
			names = append(names, id) // Fallback to ID if not found
		}
	}
	return names
}
