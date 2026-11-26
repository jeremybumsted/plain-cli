package helpcenters

import "fmt"

// GroupGetCmd represents the help-centers group get command.
type GroupGetCmd struct {
	ID         string `arg:"" name:"id" help:"Article group ID"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the help-centers group get command.
func (cmd *GroupGetCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	group, err := client.GetHelpCenterArticleGroup(cmd.ID)
	formatter := getFormatter(cmd.Format)
	if err != nil {
		if handled := formatGroupNotFound(err, cmd.ID, formatter); handled != err {
			return handled
		}
		return fmt.Errorf("failed to fetch article group: %w", err)
	}

	if cmd.Format == "quiet" {
		return formatter.Print(group.ID)
	}
	if cmd.Format == "json" {
		return formatter.PrintJSON(group)
	}

	return displayGroup(formatter, group)
}
