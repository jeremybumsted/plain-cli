package threads

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/mcp"
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
	priorityInt, err := parsePriorityString(cmd.Priority)
	if err != nil {
		return fmt.Errorf("invalid priority: %w", err)
	}

	// Get priority name for display
	priorityName := mcp.FormatPriority(priorityInt)

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
		if mcpErr, ok := err.(*mcp.Error); ok && mcpErr.StatusCode == 404 {
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

// parsePriorityString converts a priority string to an integer value
// Accepts: "urgent", "high", "normal", "low", "0", "1", "2", "3"
// Returns: 3 for urgent, 2 for high, 1 for normal, 0 for low
func parsePriorityString(priority string) (int, error) {
	// Normalize to lowercase
	priority = strings.ToLower(strings.TrimSpace(priority))

	// Try parsing as integer first
	if val, err := strconv.Atoi(priority); err == nil {
		if val >= 0 && val <= 3 {
			return val, nil
		}
		return 0, fmt.Errorf("priority must be between 0 and 3, got %d", val)
	}

	// Parse as string
	switch priority {
	case "urgent":
		return 3, nil
	case "high":
		return 2, nil
	case "normal":
		return 1, nil
	case "low":
		return 0, nil
	default:
		return 0, fmt.Errorf("invalid priority '%s', expected: urgent, high, normal, low, or 0-3", priority)
	}
}
