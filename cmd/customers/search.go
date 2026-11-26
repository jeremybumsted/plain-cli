package customers

import (
	"fmt"
)

// SearchCmd represents the "customers search" command
type SearchCmd struct {
	Query      string `arg:"" help:"Search query" required:""`
	Limit      int    `help:"Number of results" default:"50"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the search command
func (cmd *SearchCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain config'")
	}

	response, err := client.SearchCustomers(cmd.Query, cmd.Limit)
	if err != nil {
		return fmt.Errorf("failed to search customers: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	if !formatter.IsQuiet() {
		if err := formatter.Info(fmt.Sprintf("Search results for: \"%s\"\n", cmd.Query)); err != nil {
			return err
		}
	}

	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(response)
	case "quiet":
		for _, customer := range response.Customers {
			if err := formatter.Print(customer.ID); err != nil {
				return err
			}
		}
		return nil
	default:
		if len(response.Customers) == 0 {
			return formatter.Info(fmt.Sprintf("No customers found matching '%s'", cmd.Query))
		}
		if err := printCustomerTable(formatter, response.Customers); err != nil {
			return err
		}
		return formatter.Info(fmt.Sprintf("\n%d result(s) found", len(response.Customers)))
	}
}
