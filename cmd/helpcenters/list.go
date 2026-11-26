package helpcenters

import "fmt"

// ListCmd represents the help-centers list command
type ListCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the help-centers list command
func (cmd *ListCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	helpCenters, err := client.ListHelpCenters()
	if err != nil {
		return fmt.Errorf("failed to list help centers: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(helpCenters)
	case "quiet":
		for _, hc := range helpCenters {
			fmt.Println(hc.ID)
		}
		return nil
	default:
		if len(helpCenters) == 0 {
			fmt.Println("No help centers found")
			return nil
		}

		fmt.Printf("%-40s\n", "ID")
		fmt.Printf("%-40s\n", "----")
		for _, hc := range helpCenters {
			fmt.Println(hc.ID)
		}
		fmt.Printf("\nTotal: %d help centers\n", len(helpCenters))
		return nil
	}
}
