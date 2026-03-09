package threads

import (
	"fmt"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/cache"
	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

// LabelAddCmd represents the threads label add command
type LabelAddCmd struct {
	ThreadID   string   `arg:"" help:"Thread ID or URL"`
	Labels     []string `arg:"" help:"Label names or IDs to add"`
	Yes        bool     `help:"Skip confirmation" short:"y"`
	ConfigPath string   `help:"Path to config file" default:""`
	Format     string   `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads label add command
func (cmd *LabelAddCmd) Run() error {
	// Extract thread ID from URL if provided
	threadID := extractThreadID(cmd.ThreadID)

	// Load config and get client
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Get formatter
	formatter := getFormatter(cmd.Format)

	// Load label cache
	labelCache, err := cache.LoadLabelCache()
	if err != nil {
		// Cache doesn't exist or is corrupted - fetch fresh data
		labelTypes, fetchErr := client.ListLabelTypes(false)
		if fetchErr != nil {
			return fmt.Errorf("failed to fetch label types: %w", fetchErr)
		}

		// Save to cache
		if saveErr := cache.SaveLabelCache(labelTypes); saveErr != nil {
			// Non-fatal: log but continue
			fmt.Printf("Warning: failed to save label cache: %v\n", saveErr)
		}

		labelCache = &cache.LabelCache{
			LabelTypes: labelTypes,
		}
	} else if !labelCache.IsFresh() {
		// Cache is stale - refresh it
		labelTypes, fetchErr := client.ListLabelTypes(false)
		if fetchErr != nil {
			// Non-fatal: use stale cache
			fmt.Printf("Warning: failed to refresh label cache, using stale data: %v\n", fetchErr)
		} else {
			// Save refreshed cache
			if saveErr := cache.SaveLabelCache(labelTypes); saveErr != nil {
				fmt.Printf("Warning: failed to save label cache: %v\n", saveErr)
			}
			labelCache.LabelTypes = labelTypes
		}
	}

	// Resolve label identifiers to IDs
	labelTypeIDs, err := labelCache.ResolveLabelIdentifiers(cmd.Labels)
	if err != nil {
		return formatter.Error(err.Error())
	}

	// Get label names for confirmation
	labelNames := labelCache.GetLabelNames(labelTypeIDs)

	// Confirm action unless --yes flag is provided
	message := fmt.Sprintf("Add labels: %s?", strings.Join(labelNames, ", "))
	confirmed, err := confirmAction(message, cmd.Yes)
	if err != nil {
		return fmt.Errorf("failed to read confirmation: %w", err)
	}
	if !confirmed {
		return formatter.Print("Operation cancelled")
	}

	// Add labels to thread
	labels, err := client.AddLabels(threadID, labelTypeIDs)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*mcp.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to add labels: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(threadID)
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(labels)
	}

	// Default: success message with label names
	addedNames := make([]string, 0, len(labels))
	for _, label := range labels {
		if label.LabelType != nil {
			addedNames = append(addedNames, label.LabelType.Name)
		}
	}

	if len(addedNames) > 0 {
		return formatter.Print(fmt.Sprintf("Added labels: %s", strings.Join(addedNames, ", ")))
	}

	return formatter.Print("Labels added successfully")
}
