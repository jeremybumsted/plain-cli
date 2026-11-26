package customers

import (
	"fmt"
)

// GetByExtIDCmd represents the "customers get-by-ext-id" command
type GetByExtIDCmd struct {
	ExternalID string `arg:"" name:"external-id" help:"Customer external ID" required:""`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the get-by-ext-id command
func (cmd *GetByExtIDCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain config'")
	}

	customer, err := client.GetCustomerByExternalID(cmd.ExternalID)
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
