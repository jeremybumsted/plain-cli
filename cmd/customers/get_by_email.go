package customers

import (
	"fmt"
)

// GetByEmailCmd represents the "customers get-by-email" command
type GetByEmailCmd struct {
	Email      string `arg:"" name:"email" help:"Customer email address" required:""`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the get-by-email command
func (cmd *GetByEmailCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain config'")
	}

	customer, err := client.GetCustomerByEmail(cmd.Email)
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
