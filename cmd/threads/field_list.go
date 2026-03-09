package threads

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/cache"
	"github.com/jeremybumsted/plain-cli/internal/mcp"
	"github.com/jeremybumsted/plain-cli/internal/output"
)

// FieldListCmd handles listing thread field schemas from cache
type FieldListCmd struct {
	Refresh    bool   `help:"Force refresh from API" short:"r"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the field list command
func (cmd *FieldListCmd) Run() error {
	// Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create API client
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Create formatter
	formatter := getFormatter(cmd.Format)

	var fieldSchemas []*mcp.ThreadFieldSchema
	var needsRefresh bool

	// Try to load cache if not forcing refresh
	if !cmd.Refresh {
		fieldCache, err := cache.LoadFieldCache()
		if err != nil {
			// Cache doesn't exist or is corrupted, need to refresh
			needsRefresh = true
			if !formatter.IsQuiet() {
				fmt.Fprintln(os.Stderr, "Cache not found, fetching from API...")
			}
		} else if !fieldCache.IsFresh() {
			// Cache is stale, need to refresh
			needsRefresh = true
			if !formatter.IsQuiet() {
				fmt.Fprintln(os.Stderr, "Cache is stale, refreshing from API...")
			}
		} else {
			// Cache is fresh, use it
			fieldSchemas = fieldCache.FieldSchemas
		}
	} else {
		needsRefresh = true
		if !formatter.IsQuiet() {
			fmt.Fprintln(os.Stderr, "Refreshing from API...")
		}
	}

	// Fetch from API if needed
	if needsRefresh {
		schemas, err := client.ListThreadFieldSchemas()
		if err != nil {
			return fmt.Errorf("failed to fetch field schemas: %w", err)
		}
		fieldSchemas = schemas

		// Save to cache
		if err := cache.SaveFieldCache(schemas); err != nil {
			// Non-fatal: warn but continue
			fmt.Fprintf(os.Stderr, "Warning: failed to save cache: %v\n", err)
		}
	}

	// Output based on format
	switch formatter.GetFormat() {
	case output.FormatJSON:
		// Output full field schema objects as JSON array
		return formatter.PrintJSON(fieldSchemas)

	case output.FormatQuiet:
		// Just IDs, one per line
		for _, schema := range fieldSchemas {
			fmt.Println(schema.ID)
		}
		return nil

	default:
		// Table format
		if len(fieldSchemas) == 0 {
			return formatter.Info("No field schemas found")
		}

		headers := []string{"ID", "Key", "Label", "Type", "Required", "Enum Values"}
		rows := make([][]string, 0, len(fieldSchemas))

		for _, schema := range fieldSchemas {
			// Format required indicator
			required := ""
			if schema.IsRequired {
				required = "*"
			}

			// Format enum values if present
			enumValues := ""
			if schema.Type == "ENUM" && len(schema.EnumValues) > 0 {
				enumValues = strings.Join(schema.EnumValues, ", ")
				// Truncate if too long
				if len(enumValues) > 50 {
					enumValues = enumValues[:47] + "..."
				}
			}

			row := []string{
				schema.ID,
				schema.Key,
				schema.Label,
				schema.Type,
				required,
				enumValues,
			}
			rows = append(rows, row)
		}

		if err := formatter.PrintTable(headers, rows); err != nil {
			return err
		}

		// Add helpful footer message
		fmt.Fprintf(os.Stderr, "\nShowing %d field schemas (* = required)\n", len(fieldSchemas))
		return nil
	}
}
