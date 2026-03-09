package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/cache"
	"github.com/jeremybumsted/plain-cli/internal/output"
)

// LabelRefreshCmd handles refreshing the label cache from the API
type LabelRefreshCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the label refresh command
func (cmd *LabelRefreshCmd) Run() error {
	// Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Create formatter
	formatter := getFormatter(cmd.Format)

	// Fetch label types from API (exclude archived)
	labelTypes, err := client.ListLabelTypes(false)
	if err != nil {
		return fmt.Errorf("failed to fetch label types: %w", err)
	}

	// Save to cache
	if err := cache.SaveLabelCache(labelTypes); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}

	// Output success message based on format
	count := len(labelTypes)

	switch formatter.GetFormat() {
	case output.FormatJSON:
		// Output array of label type objects
		return formatter.PrintJSON(labelTypes)
	case output.FormatQuiet:
		// Just the count number
		fmt.Println(count)
		return nil
	default:
		// Table format: success message with count
		fmt.Printf("Cache refreshed with %d label types\n", count)
		return nil
	}
}
