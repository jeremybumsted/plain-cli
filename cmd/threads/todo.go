package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// TodoCmd represents the threads todo command
type TodoCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads todo command
func (cmd *TodoCmd) Run() error {
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

	// Ask for confirmation unless --yes flag is set
	confirmed, err := confirmAction(fmt.Sprintf("Mark thread %s as TODO?", threadID), cmd.Yes)
	if err != nil {
		return fmt.Errorf("confirmation failed: %w", err)
	}

	if !confirmed {
		return formatter.Info("Operation cancelled")
	}

	// Change thread status to TODO
	thread, err := client.ChangeThreadStatus(threadID, "TODO")
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*plain.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to change thread status: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(thread.ID)
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(thread)
	}

	// Default: success message
	return formatter.Success("Thread marked as todo")
}
