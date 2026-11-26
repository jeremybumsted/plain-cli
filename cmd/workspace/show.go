package workspace

import "fmt"

// ShowCmd represents the workspace show command
type ShowCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the workspace show command
func (cmd *ShowCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	ws, err := client.GetWorkspace()
	if err != nil {
		return fmt.Errorf("failed to fetch workspace: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(ws)
	case "quiet":
		fmt.Println(ws.ID)
		return nil
	default:
		isDemo := "no"
		if ws.IsDemoWorkspace {
			isDemo = "yes"
		}

		fmt.Printf("ID:             %s\n", ws.ID)
		fmt.Printf("Name:           %s\n", ws.Name)
		if ws.PublicName != "" {
			fmt.Printf("Public Name:    %s\n", ws.PublicName)
		}
		fmt.Printf("Demo Workspace: %s\n", isDemo)

		createdAt, _ := ws.CreatedAt.Time()
		updatedAt, _ := ws.UpdatedAt.Time()
		fmt.Printf("Created:        %s\n", createdAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:        %s\n", updatedAt.Format("2006-01-02 15:04:05"))

		return nil
	}
}
