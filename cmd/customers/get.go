package customers

import (
	"fmt"
)

// GetCmd represents the "customers get" command
type GetCmd struct {
	CustomerID string `arg:"" name:"id" help:"Customer ID" required:""`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the get command
func (cmd *GetCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain config'")
	}

	customer, err := client.GetCustomer(cmd.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(customer)
	case "quiet":
		return formatter.Print(customer.ID)
	default:
		return displayCustomer(formatter, customer)
	}
}
