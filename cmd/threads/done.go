package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

// DoneCmd represents the threads done command
type DoneCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads done command
func (cmd *DoneCmd) Run() error {
	// Extract thread ID from URL if provided
	threadID := extractThreadID(cmd.ThreadID)

	// Confirm action unless --yes flag is provided
	confirmed, err := confirmAction(fmt.Sprintf("Mark thread %s as done?", threadID), cmd.Yes)
	if err != nil {
		return fmt.Errorf("confirmation failed: %w", err)
	}
	if !confirmed {
		return fmt.Errorf("operation cancelled")
	}

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

	// Change thread status to DONE
	thread, err := client.ChangeThreadStatus(threadID, "DONE")
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*mcp.Error); ok && mcpErr.StatusCode == 404 {
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
	return formatter.Success("Thread marked as done")
}
