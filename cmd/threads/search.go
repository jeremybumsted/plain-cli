package threads

import (
	"fmt"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/plain"
	"github.com/jeremybumsted/plain-cli/internal/util"
)

// SearchCmd represents the threads search command
type SearchCmd struct {
	Query         string `arg:"" help:"Search query" required:""`
	Status        string `help:"Filter by status" default:""`
	Priority      string `help:"Filter by priority" default:""`
	Assignee      string `help:"Filter by assignee ID" default:""`
	Label         string `help:"Filter by label IDs (comma-separated)" default:""`
	CreatedAfter  string `help:"Filter threads created after this date (ISO8601, relative like '7d', or 'yesterday')" default:""`
	CreatedBefore string `help:"Filter threads created before this date (ISO8601 or relative)" default:""`
	UpdatedAfter  string `help:"Filter threads updated after this date (ISO8601 or relative)" default:""`
	UpdatedBefore string `help:"Filter threads updated before this date (ISO8601 or relative)" default:""`
	Limit         int    `help:"Number of results" default:"50"`
	Offset        int    `help:"Pagination offset" default:"0"`
	ConfigPath    string `help:"Path to config file" default:""`
	Format        string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the search command
func (cmd *SearchCmd) Run() error {
	// Load config and get client
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Get formatter
	formatter := getFormatter(cmd.Format)

	// Build filters
	filters := &plain.ThreadFilters{
		Status:   cmd.Status,
		Priority: cmd.Priority,
		Limit:    cmd.Limit,
		Offset:   cmd.Offset,
	}

	// Parse assignee if provided
	if cmd.Assignee != "" {
		filters.AssigneeID = cmd.Assignee
	}

	// Parse label IDs if provided
	if cmd.Label != "" {
		labelIDs := strings.Split(cmd.Label, ",")
		for i := range labelIDs {
			labelIDs[i] = strings.TrimSpace(labelIDs[i])
		}
		filters.LabelIDs = labelIDs
	}

	// Parse and validate date filters
	if cmd.CreatedAfter != "" {
		parsedDate, err := util.ParseDateToISO8601(cmd.CreatedAfter)
		if err != nil {
			return fmt.Errorf("invalid --created-after date: %w", err)
		}
		filters.CreatedAfter = parsedDate
	}

	if cmd.CreatedBefore != "" {
		parsedDate, err := util.ParseDateToISO8601(cmd.CreatedBefore)
		if err != nil {
			return fmt.Errorf("invalid --created-before date: %w", err)
		}
		filters.CreatedBefore = parsedDate
	}

	if cmd.UpdatedAfter != "" {
		parsedDate, err := util.ParseDateToISO8601(cmd.UpdatedAfter)
		if err != nil {
			return fmt.Errorf("invalid --updated-after date: %w", err)
		}
		filters.UpdatedAfter = parsedDate
	}

	if cmd.UpdatedBefore != "" {
		parsedDate, err := util.ParseDateToISO8601(cmd.UpdatedBefore)
		if err != nil {
			return fmt.Errorf("invalid --updated-before date: %w", err)
		}
		filters.UpdatedBefore = parsedDate
	}

	// Validate date ranges
	if err := util.ValidateDateRange(cmd.CreatedAfter, cmd.CreatedBefore); err != nil {
		return fmt.Errorf("invalid created date range: %w", err)
	}

	if err := util.ValidateDateRange(cmd.UpdatedAfter, cmd.UpdatedBefore); err != nil {
		return fmt.Errorf("invalid updated date range: %w", err)
	}

	// Search threads
	response, err := client.SearchThreads(cmd.Query, filters)
	if err != nil {
		return fmt.Errorf("failed to search threads: %w", err)
	}

	// Display search header (unless in quiet mode)
	if !formatter.IsQuiet() {
		if err := formatter.Info(fmt.Sprintf("Search results for: \"%s\"\n", cmd.Query)); err != nil {
			return err
		}
	}

	// Handle JSON output
	if formatter.GetFormat() == "json" {
		return formatter.PrintJSON(response)
	}

	// Handle empty results
	if len(response.Threads) == 0 {
		return formatter.Info(fmt.Sprintf("No threads found matching '%s'", cmd.Query))
	}

	// Handle quiet mode - just print thread IDs
	if formatter.IsQuiet() {
		for _, thread := range response.Threads {
			if err := formatter.Print(thread.ID); err != nil {
				return err
			}
		}
		return nil
	}

	// Build table for display
	headers := []string{"ID", "Title", "Status", "Priority", "Assignee", "Labels", "Updated"}
	rows := make([][]string, 0, len(response.Threads))

	for _, thread := range response.Threads {
		// Format assignee
		assignee := "-"
		if thread.AssignedTo != nil {
			if thread.AssignedTo.FullName != "" {
				assignee = thread.AssignedTo.FullName
			} else {
				assignee = thread.AssignedTo.Email
			}
		}

		// Format labels
		labels := "-"
		if len(thread.Labels) > 0 {
			labelNames := make([]string, 0, len(thread.Labels))
			for _, label := range thread.Labels {
				if label.LabelType != nil {
					labelNames = append(labelNames, label.LabelType.Name)
				}
			}
			if len(labelNames) > 0 {
				labels = strings.Join(labelNames, ", ")
			}
		}

		// Format priority
		priority := plain.FormatPriority(thread.Priority)

		// Format updated time
		updatedTime, _ := thread.UpdatedAt.Time()
		updatedAt := formatTime(updatedTime)

		row := []string{
			thread.ID,
			truncateString(thread.Title, 40),
			thread.Status,
			priority,
			truncateString(assignee, 20),
			truncateString(labels, 30),
			updatedAt,
		}
		rows = append(rows, row)
	}

	// Print table
	if err := formatter.PrintTable(headers, rows); err != nil {
		return fmt.Errorf("failed to display results: %w", err)
	}

	// Print summary
	resultCount := len(response.Threads)
	if response.Total > resultCount {
		return formatter.Info(fmt.Sprintf("\nShowing %d of %d results", resultCount, response.Total))
	}
	return formatter.Info(fmt.Sprintf("\n%d result(s) found", resultCount))
}
