package threads

import (
	"fmt"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/plain"
	"github.com/jeremybumsted/plain-cli/internal/util"
)

// ListCmd represents the list threads command
type ListCmd struct {
	Status        string `help:"Filter by status" default:""`
	Assignee      string `help:"Filter by assignee ID" default:""`
	Priority      string `help:"Filter by priority" default:""`
	Label         string `help:"Filter by label IDs (comma-separated)" default:""`
	CreatedAfter  string `help:"Filter threads created after this date (ISO8601, relative like '7d', or 'yesterday')" default:""`
	CreatedBefore string `help:"Filter threads created before this date (ISO8601 or relative)" default:""`
	UpdatedAfter  string `help:"Filter threads updated after this date (ISO8601 or relative)" default:""`
	UpdatedBefore string `help:"Filter threads updated before this date (ISO8601 or relative)" default:""`
	Mine          bool   `help:"Show only threads assigned to me"`
	Limit         int    `help:"Number of results" default:"50"`
	Offset        int    `help:"Pagination offset" default:"0"`
	ConfigPath    string `help:"Path to config file" default:""`
	Format        string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the list threads command
func (cmd *ListCmd) Run() error {
	// 1. Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Get authenticated Plain API client
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain auth login'")
	}

	// 3. Validate flags - ensure --mine and --assignee are mutually exclusive
	if cmd.Mine && cmd.Assignee != "" {
		return fmt.Errorf("cannot use both --mine and --assignee flags")
	}

	// 4. Parse and validate date filters
	var createdAfterISO, createdBeforeISO, updatedAfterISO, updatedBeforeISO string

	// Parse CreatedAfter
	if cmd.CreatedAfter != "" {
		parsed, err := util.ParseDateToISO8601(cmd.CreatedAfter)
		if err != nil {
			return fmt.Errorf("invalid created-after date: %w", err)
		}
		createdAfterISO = parsed
	}

	// Parse CreatedBefore
	if cmd.CreatedBefore != "" {
		parsed, err := util.ParseDateToISO8601(cmd.CreatedBefore)
		if err != nil {
			return fmt.Errorf("invalid created-before date: %w", err)
		}
		createdBeforeISO = parsed
	}

	// Validate created date range if both are provided
	if createdAfterISO != "" && createdBeforeISO != "" {
		if err := util.ValidateDateRange(createdAfterISO, createdBeforeISO); err != nil {
			return fmt.Errorf("invalid created date range: %w", err)
		}
	}

	// Parse UpdatedAfter
	if cmd.UpdatedAfter != "" {
		parsed, err := util.ParseDateToISO8601(cmd.UpdatedAfter)
		if err != nil {
			return fmt.Errorf("invalid updated-after date: %w", err)
		}
		updatedAfterISO = parsed
	}

	// Parse UpdatedBefore
	if cmd.UpdatedBefore != "" {
		parsed, err := util.ParseDateToISO8601(cmd.UpdatedBefore)
		if err != nil {
			return fmt.Errorf("invalid updated-before date: %w", err)
		}
		updatedBeforeISO = parsed
	}

	// Validate updated date range if both are provided
	if updatedAfterISO != "" && updatedBeforeISO != "" {
		if err := util.ValidateDateRange(updatedAfterISO, updatedBeforeISO); err != nil {
			return fmt.Errorf("invalid updated date range: %w", err)
		}
	}

	// 5. Parse filters from flags
	filters := &plain.ThreadFilters{
		Status:        cmd.Status,
		Priority:      cmd.Priority,
		Limit:         cmd.Limit,
		Offset:        cmd.Offset,
		CreatedAfter:  createdAfterISO,
		CreatedBefore: createdBeforeISO,
		UpdatedAfter:  updatedAfterISO,
		UpdatedBefore: updatedBeforeISO,
	}

	// Handle --mine flag
	if cmd.Mine {
		userID, err := cfg.GetUserID()
		if err != nil {
			return fmt.Errorf("cannot use --mine: %w", err)
		}
		filters.AssigneeID = userID
	} else if cmd.Assignee != "" {
		filters.AssigneeID = cmd.Assignee
	}

	// Parse label IDs (comma-separated)
	if cmd.Label != "" {
		labelIDs := strings.Split(cmd.Label, ",")
		for i, id := range labelIDs {
			labelIDs[i] = strings.TrimSpace(id)
		}
		filters.LabelIDs = labelIDs
	}

	// 6. Call client.ListThreads()
	response, err := client.ListThreads(filters)
	if err != nil {
		return fmt.Errorf("failed to list threads: %w", err)
	}

	// 7. Format and display results
	formatter := getFormatter(cmd.Format)

	// Handle different output formats
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(response)
	case "quiet":
		// Print only thread IDs in quiet mode
		if len(response.Threads) == 0 {
			return nil
		}
		for _, thread := range response.Threads {
			if err := formatter.Print(thread.ID); err != nil {
				return err
			}
		}
		return nil
	default:
		// Table format
		return cmd.printTable(formatter, response)
	}
}

// printTable formats and prints threads as a table
func (cmd *ListCmd) printTable(formatter interface{ PrintTable([]string, [][]string) error; Info(string) error }, response *plain.ThreadsResponse) error {
	if len(response.Threads) == 0 {
		return formatter.Info("No threads found")
	}

	headers := []string{"ID", "Title", "Status", "Priority", "Assignee", "Updated"}
	rows := make([][]string, 0, len(response.Threads))

	for _, thread := range response.Threads {
		assignee := "-"
		if thread.AssignedTo != nil {
			if thread.AssignedTo.FullName != "" {
				assignee = thread.AssignedTo.FullName
			} else {
				assignee = thread.AssignedTo.Email
			}
		}

		priority := plain.FormatPriority(thread.Priority)

		// Format the title - truncate if too long
		title := thread.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}
		if title == "" {
			title = "(no title)"
		}

		// Format updated time
		updatedTime, _ := thread.UpdatedAt.Time()
		updatedAt := formatTime(updatedTime)

		row := []string{
			thread.ID,
			title,
			thread.Status,
			priority,
			assignee,
			updatedAt,
		}
		rows = append(rows, row)
	}

	return formatter.PrintTable(headers, rows)
}
