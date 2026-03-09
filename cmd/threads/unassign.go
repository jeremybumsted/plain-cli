package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

// UnassignCmd represents the threads unassign command
type UnassignCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads unassign command
func (cmd *UnassignCmd) Run() error {
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

	// Confirm action unless --yes flag is set
	confirmed, err := confirmAction("Unassign thread?", cmd.Yes)
	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}
	if !confirmed {
		return formatter.Info("Operation cancelled")
	}

	// Unassign thread by passing nil as userID
	thread, err := client.AssignThread(threadID, nil)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*mcp.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to unassign thread: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(thread.ID)
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(thread)
	}

	// Default: success message
	return formatter.Success("Thread unassigned")
}
