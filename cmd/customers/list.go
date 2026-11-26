package customers

import (
	"fmt"
)

// ListCmd represents the "customers list" command
type ListCmd struct {
	Limit      int    `help:"Number of results" default:"50"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the list command
func (cmd *ListCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain config'")
	}

	response, err := client.ListCustomers(cmd.Limit, "")
	if err != nil {
		return fmt.Errorf("failed to list customers: %w", err)
	}

	formatter := getFormatter(cmd.Format)

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
		return printCustomerTable(formatter, response.Customers)
	}
}
