package threads

import (
	"fmt"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/cache"
	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

// LabelRemoveCmd represents the threads label remove command
type LabelRemoveCmd struct {
	ThreadID   string   `arg:"" help:"Thread ID or URL"`
	Labels     []string `arg:"" help:"Label names or type IDs to remove"`
	Yes        bool     `help:"Skip confirmation" short:"y"`
	ConfigPath string   `help:"Path to config file" default:""`
	Format     string   `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads label remove command
func (cmd *LabelRemoveCmd) Run() error {
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
		return fmt.Errorf("failed to load label cache: %w (try 'plain threads label refresh')", err)
	}

	// Auto-refresh if cache is stale
	if !labelCache.IsFresh() {
		labelTypes, err := client.ListLabelTypes(false)
		if err != nil {
			return fmt.Errorf("failed to refresh label cache: %w", err)
		}
		if err := cache.SaveLabelCache(labelTypes); err != nil {
			return fmt.Errorf("failed to save label cache: %w", err)
		}
		labelCache.LabelTypes = labelTypes
	}

	// Resolve label identifiers to label type IDs
	labelTypeIDs, err := labelCache.ResolveLabelIdentifiers(cmd.Labels)
	if err != nil {
		return fmt.Errorf("failed to resolve label identifiers: %w", err)
	}

	// Fetch the thread to get current labels
	thread, err := client.GetThread(threadID, false)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*mcp.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to fetch thread: %w", err)
	}

	// Match label types to find label instance IDs
	labelInstanceIDs := make([]string, 0)
	labelNames := make([]string, 0)
	notFoundLabelTypes := make([]string, 0)

	for _, labelTypeID := range labelTypeIDs {
		found := false
		for _, label := range thread.Labels {
			if label.LabelType != nil && label.LabelType.ID == labelTypeID {
				labelInstanceIDs = append(labelInstanceIDs, label.ID)
				labelNames = append(labelNames, label.LabelType.Name)
				found = true
				break
			}
		}
		if !found {
			// Get the name for error reporting
			labelType := labelCache.GetLabelTypeByID(labelTypeID)
			if labelType != nil {
				notFoundLabelTypes = append(notFoundLabelTypes, labelType.Name)
			} else {
				notFoundLabelTypes = append(notFoundLabelTypes, labelTypeID)
			}
		}
	}

	// Error if any label types were not found on the thread
	if len(notFoundLabelTypes) > 0 {
		return formatter.Error(fmt.Sprintf("Label(s) not found on thread: %s", strings.Join(notFoundLabelTypes, ", ")))
	}

	// If no labels to remove (shouldn't happen, but check anyway)
	if len(labelInstanceIDs) == 0 {
		return formatter.Error("No labels to remove")
	}

	// Confirm action unless --yes flag is provided
	confirmed, err := confirmAction(fmt.Sprintf("Remove labels: %s?", strings.Join(labelNames, ", ")), cmd.Yes)
	if err != nil {
		return fmt.Errorf("failed to read confirmation: %w", err)
	}
	if !confirmed {
		return formatter.Print("Operation cancelled")
	}

	// Remove labels
	err = client.RemoveLabels(labelInstanceIDs)
	if err != nil {
		return fmt.Errorf("failed to remove labels: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(threadID)
	}

	if cmd.Format == "json" {
		// Return a simple confirmation message as JSON
		result := map[string]interface{}{
			"thread_id":      threadID,
			"removed_labels": labelNames,
		}
		return formatter.PrintJSON(result)
	}

	// Default: success message
	return formatter.Print(fmt.Sprintf("Removed labels: %s", strings.Join(labelNames, ", ")))
}
