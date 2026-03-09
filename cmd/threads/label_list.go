package threads

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/cache"
	"github.com/jeremybumsted/plain-cli/internal/mcp"
)

// LabelListCmd represents the threads label list command
type LabelListCmd struct {
	Refresh    bool   `help:"Force refresh from API" short:"r"`
	Archived   bool   `help:"Include archived labels" short:"a"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads label list command
// This lists all available label types (templates) that can be used with the add command
func (cmd *LabelListCmd) Run() error {
	// 1. Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Get authenticated MCP client
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// 3. Create formatter
	formatter := getFormatter(cmd.Format)

	// 4. Determine if we need to fetch from API
	var labelTypes []*mcp.LabelType
	needsRefresh := cmd.Refresh

	// Try to load from cache if not forcing refresh
	if !needsRefresh {
		labelCache, err := cache.LoadLabelCache()
		if err != nil {
			// Cache doesn't exist or is corrupted - need to fetch from API
			needsRefresh = true
		} else if !labelCache.IsFresh() {
			// Cache is stale - need to refresh
			needsRefresh = true
		} else {
			// Cache is valid - use it
			labelTypes = labelCache.LabelTypes
		}
	}

	// 5. Fetch from API if needed and update cache
	if needsRefresh {
		fetchedLabels, err := client.ListLabelTypes(cmd.Archived)
		if err != nil {
			return fmt.Errorf("failed to fetch label types: %w", err)
		}
		labelTypes = fetchedLabels

		// Save to cache (non-fatal if it fails)
		if err := cache.SaveLabelCache(labelTypes); err != nil {
			// Log warning but continue
			if !formatter.IsQuiet() {
				fmt.Printf("Warning: failed to save cache: %v\n", err)
			}
		}
	}

	// 6. Filter archived labels if --archived not set
	if !cmd.Archived {
		filtered := make([]*mcp.LabelType, 0, len(labelTypes))
		for _, lt := range labelTypes {
			if !lt.IsArchived {
				filtered = append(filtered, lt)
			}
		}
		labelTypes = filtered
	}

	// 7. Output based on format
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(labelTypes)
	case "quiet":
		// Just IDs, one per line
		for _, lt := range labelTypes {
			formatter.Print(lt.ID)
		}
		return nil
	default:
		// Table format
		return cmd.printTable(formatter, labelTypes)
	}
}

// printTable formats and prints label types as a table
func (cmd *LabelListCmd) printTable(formatter interface{ PrintTable([]string, [][]string) error; Info(string) error }, labelTypes []*mcp.LabelType) error {
	if len(labelTypes) == 0 {
		return formatter.Info("No label types found")
	}

	headers := []string{"ID", "Name", "Icon", "Color", "Archived"}
	rows := make([][]string, 0, len(labelTypes))

	for _, lt := range labelTypes {
		// Handle empty values
		icon := lt.Icon
		if icon == "" {
			icon = "-"
		}

		color := lt.Color
		if color == "" {
			color = "-"
		}

		archived := "No"
		if lt.IsArchived {
			archived = "Yes"
		}

		row := []string{
			lt.ID,
			lt.Name,
			icon,
			color,
			archived,
		}
		rows = append(rows, row)
	}

	return formatter.PrintTable(headers, rows)
}
