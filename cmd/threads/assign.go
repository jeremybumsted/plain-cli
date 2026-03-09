package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

// AssignCmd represents the threads assign command
type AssignCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	UserID     string `arg:"" help:"User ID to assign to" required:""`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads assign command
func (cmd *AssignCmd) Run() error {
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

	// Confirm action unless --yes flag is provided
	confirmed, err := confirmAction(fmt.Sprintf("Assign thread to user %s?", cmd.UserID), cmd.Yes)
	if err != nil {
		return fmt.Errorf("failed to read confirmation: %w", err)
	}
	if !confirmed {
		return formatter.Print("Operation cancelled")
	}

	// Assign thread to user
	thread, err := client.AssignThread(threadID, &cmd.UserID)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*mcp.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to assign thread: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(thread.ID)
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(thread)
	}

	// Default: success message
	userName := cmd.UserID
	if thread.AssignedTo != nil && thread.AssignedTo.FullName != "" {
		userName = thread.AssignedTo.FullName
	}

	return formatter.Print(fmt.Sprintf("Thread assigned to %s", userName))
}
