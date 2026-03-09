package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

// SnoozeCmd represents the threads snooze command
type SnoozeCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	Until      string `help:"Snooze until (e.g., 2h, 1d, 3w, or ISO8601)" default:"1d"`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads snooze command
func (cmd *SnoozeCmd) Run() error {
	// Extract thread ID from URL if provided
	threadID := extractThreadID(cmd.ThreadID)

	// Parse the until time (relative or absolute)
	until, err := parseRelativeTime(cmd.Until)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
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

	// Format the until time for display
	formattedDate := until.Format("2006-01-02 15:04:05")

	// Confirm action unless --yes flag is set
	confirmed, err := confirmAction(fmt.Sprintf("Snooze thread until %s?", formattedDate), cmd.Yes)
	if err != nil {
		return fmt.Errorf("confirmation failed: %w", err)
	}

	if !confirmed {
		return formatter.Info("Operation cancelled")
	}

	// Snooze the thread
	thread, err := client.SnoozeThread(threadID, until)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*mcp.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to snooze thread: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(thread.ID)
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(thread)
	}

	// Default: success message
	return formatter.Success(fmt.Sprintf("Thread snoozed until %s", formattedDate))
}
