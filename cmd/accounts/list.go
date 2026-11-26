package accounts

import "fmt"

// ListCmd represents `plain accounts list`.
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

	accounts, err := client.ListAccounts(cmd.Limit)
	if err != nil {
		return fmt.Errorf("failed to list accounts: %w", err)
	}

	formatter := getFormatter(cmd.Format)
	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(accounts)
	case "quiet":
		for _, account := range accounts {
			if err := formatter.Print(account.ID); err != nil {
				return err
			}
		}
		return nil
	default:
		if len(accounts) == 0 {
			return formatter.Info("No accounts found")
		}
		headers := []string{"ID", "Name", "External ID", "Tier", "Updated"}
		rows := make([][]string, 0, len(accounts))
		for _, account := range accounts {
			tier := "-"
			if account.Tier != nil {
				tier = account.Tier.Name
			}
			updated := account.UpdatedAt.ISO8601
			if len(updated) > 10 {
				updated = updated[:10]
			}
			rows = append(rows, []string{account.ID, account.Name, account.ExternalID, tier, updated})
		}
		return formatter.PrintTable(headers, rows)
	}
}
