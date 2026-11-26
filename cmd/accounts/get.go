package accounts

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// GetCmd represents `plain accounts get`.
type GetCmd struct {
	ID         string `arg:"" name:"id" help:"Account ID" required:""`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

func (cmd *GetCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain config'")
	}

	account, err := client.GetAccount(cmd.ID)
	if err != nil {
		if apiErr, ok := err.(*plain.Error); ok && apiErr.StatusCode == 404 {
			return getFormatter(cmd.Format).Error(apiErr.Message)
		}
		return fmt.Errorf("failed to fetch account: %w", err)
	}

	formatter := getFormatter(cmd.Format)
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(account)
	case "quiet":
		return formatter.Print(account.ID)
	default:
		tier := "-"
		if account.Tier != nil {
			tier = account.Tier.Name
		}
		return formatter.Printf("ID: %s\nName: %s\nExternal ID: %s\nURL: %s\nSource: %s\nTier: %s\nCreated: %s\nUpdated: %s\n", account.ID, account.Name, account.ExternalID, account.URL, account.Source, tier, account.CreatedAt.ISO8601, account.UpdatedAt.ISO8601)
	}
}
