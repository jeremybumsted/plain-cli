package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// MineCmd lists threads assigned to the authenticated user
//
// NOTE: This command is currently DISABLED because it requires OAuth user authentication.
// The Plain API's `myUser` query does not work with machine user tokens.
// To enable this command:
// 1. Implement OAuth 2.0 flow in cmd/auth/login.go (currently TODO)
// 2. Uncomment the Mine field in cmd/threads/threads.go
// 3. Users will need to authenticate with `plain auth login` using OAuth
//
// Alternative: Users can use `plain threads list --assignee=<user-id>` as a workaround
type MineCmd struct {
	Status     string `help:"Filter by status (default: active)" default:""`
	Limit      int    `help:"Number of results" default:"50"`
	Offset     int    `help:"Pagination offset" default:"0"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the mine command
func (cmd *MineCmd) Run() error {
	// 1. Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Create client
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// 3. Get formatter
	formatter := getFormatter(cmd.Format)

	// 4. Set default status if not provided
	status := cmd.Status
	if status == "" {
		status = "active"
	}

	// 5. Fetch threads assigned to current user
	formatter.Info(fmt.Sprintf("Fetching your %s threads...", status))
	response, err := client.GetMyThreads(status, cmd.Limit, cmd.Offset)
	if err != nil {
		return fmt.Errorf("failed to fetch threads: %w", err)
	}

	// 6. Handle JSON output
	if cmd.Format == "json" {
		return formatter.PrintJSON(response)
	}

	// 7. Handle quiet output (just IDs)
	if cmd.Format == "quiet" {
		for _, thread := range response.Threads {
			formatter.Print(thread.ID)
		}
		return nil
	}

	// 8. Format table output
	if len(response.Threads) == 0 {
		return formatter.Info("No threads found")
	}

	// Prepare table headers
	headers := []string{"ID", "Title", "Status", "Priority", "Updated"}
	rows := make([][]string, 0, len(response.Threads))

	// Convert threads to table rows
	for _, thread := range response.Threads {
		// Format the updated time
		updatedTime, _ := thread.UpdatedAt.Time()
		updatedAt := formatTime(updatedTime)

		rows = append(rows, []string{
			thread.ID,
			truncateString(thread.Title, 50),
			thread.Status,
			plain.FormatPriority(thread.Priority),
			updatedAt,
		})
	}

	// Print table
	if err := formatter.PrintTable(headers, rows); err != nil {
		return fmt.Errorf("failed to print table: %w", err)
	}

	// Print summary
	formatter.Info(fmt.Sprintf("\nShowing %d of %d total threads", len(response.Threads), response.Total))

	return nil
}
