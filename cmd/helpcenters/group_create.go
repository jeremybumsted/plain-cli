package helpcenters

import "fmt"

// GroupCreateCmd represents the help-centers group create command.
type GroupCreateCmd struct {
	HelpCenterID string `help:"Help center ID (overrides config)" optional:""`
	Name         string `help:"Group name" required:""`
	ParentID     string `name:"parent" help:"Parent article group ID" optional:""`
	ConfigPath   string `help:"Path to config file" default:""`
	Format       string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the help-centers group create command.
func (cmd *GroupCreateCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	helpCenterID := cmd.HelpCenterID
	if helpCenterID == "" {
		helpCenterID, err = cfg.GetHelpCenterID()
		if err != nil {
			return err
		}
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	group, err := client.CreateHelpCenterArticleGroup(helpCenterID, cmd.Name, cmd.ParentID)
	if err != nil {
		return fmt.Errorf("failed to create article group: %w", err)
	}

	formatter := getFormatter(cmd.Format)
	if cmd.Format == "quiet" {
		return formatter.Print(group.ID)
	}
	if cmd.Format == "json" {
		return formatter.PrintJSON(group)
	}

	if err := formatter.Success("Article group created"); err != nil {
		return err
	}
	return displayGroup(formatter, group)
}
