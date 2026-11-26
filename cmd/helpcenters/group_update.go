package helpcenters

import "fmt"

// GroupUpdateCmd represents the help-centers group update command.
type GroupUpdateCmd struct {
	ID         string `arg:"" name:"id" help:"Article group ID"`
	Name       string `help:"New group name" required:""`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the help-centers group update command.
func (cmd *GroupUpdateCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	group, err := client.UpdateHelpCenterArticleGroup(cmd.ID, cmd.Name)
	formatter := getFormatter(cmd.Format)
	if err != nil {
		return fmt.Errorf("failed to update article group: %w", err)
	}

	if cmd.Format == "quiet" {
		return formatter.Print(group.ID)
	}
	if cmd.Format == "json" {
		return formatter.PrintJSON(group)
	}

	if err := formatter.Success("Article group updated"); err != nil {
		return err
	}
	return displayGroup(formatter, group)
}
