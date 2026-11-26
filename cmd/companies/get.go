package companies

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// GetCmd represents `plain companies get`.
type GetCmd struct {
	ID         string `arg:"" name:"id" help:"Company ID" required:""`
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

	company, err := client.GetCompany(cmd.ID)
	if err != nil {
		if apiErr, ok := err.(*plain.Error); ok && apiErr.StatusCode == 404 {
			return getFormatter(cmd.Format).Error(apiErr.Message)
		}
		return fmt.Errorf("failed to fetch company: %w", err)
	}

	formatter := getFormatter(cmd.Format)
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(company)
	case "quiet":
		return formatter.Print(company.ID)
	default:
		tier := "-"
		if company.Tier != nil {
			tier = company.Tier.Name
		}
		return formatter.Printf("ID: %s\nName: %s\nDomain: %s\nTier: %s\nDeleted: %v\nCreated: %s\nUpdated: %s\n", company.ID, company.Name, company.DomainName, tier, company.IsDeleted, company.CreatedAt.ISO8601, company.UpdatedAt.ISO8601)
	}
}
