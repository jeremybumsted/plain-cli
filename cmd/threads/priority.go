package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// PriorityCmd represents the threads priority command
type PriorityCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	Priority   string `arg:"" help:"Priority (urgent, high, normal, low, or 0-3)" required:""`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads priority command
func (cmd *PriorityCmd) Run() error {
	// Extract thread ID from URL if provided
	threadID := extractThreadID(cmd.ThreadID)

	// Parse priority string to integer
	priorityInt, err := plain.ParsePriority(cmd.Priority)
	if err != nil {
		return fmt.Errorf("invalid priority: %w", err)
	}

	// Get priority name for display
	priorityName := plain.FormatPriority(priorityInt)

	// Confirm action unless --yes flag is set
	confirmed, err := confirmAction(fmt.Sprintf("Change priority to %s?", priorityName), cmd.Yes)
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

	// Change thread priority
	thread, err := client.ChangeThreadPriority(threadID, priorityInt)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*plain.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to change thread priority: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(thread.ID)
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(thread)
	}

	// Default: success message
	return formatter.Success(fmt.Sprintf("Priority changed to %s", priorityName))
}
