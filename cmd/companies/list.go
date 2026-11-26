package companies

import "fmt"

// ListCmd represents `plain companies list`.
type ListCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
	Limit      int    `help:"Number of results" default:"50"`
}

func (cmd *ListCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: run 'plain config'")
	}

	companies, err := client.ListCompanies(cmd.Limit)
	if err != nil {
		return fmt.Errorf("failed to list companies: %w", err)
	}

	formatter := getFormatter(cmd.Format)
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(companies)
	case "quiet":
		for _, company := range companies {
			if err := formatter.Print(company.ID); err != nil {
				return err
			}
		}
		return nil
	default:
		if len(companies) == 0 {
			return formatter.Info("No companies found")
		}
		headers := []string{"ID", "Name", "Domain", "Tier", "Updated"}
		rows := make([][]string, 0, len(companies))
		for _, company := range companies {
			tier := "-"
			if company.Tier != nil {
				tier = company.Tier.Name
			}
			updated := company.UpdatedAt.ISO8601
			if len(updated) > 10 {
				updated = updated[:10]
			}
			rows = append(rows, []string{company.ID, company.Name, company.DomainName, tier, updated})
		}
		return formatter.PrintTable(headers, rows)
	}
}
