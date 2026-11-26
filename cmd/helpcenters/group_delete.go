package helpcenters

import "fmt"

// GroupDeleteCmd represents the help-centers group delete command.
type GroupDeleteCmd struct {
	ID         string `arg:"" name:"id" help:"Article group ID"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the help-centers group delete command.
func (cmd *GroupDeleteCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	if err := client.DeleteHelpCenterArticleGroup(cmd.ID); err != nil {
		return fmt.Errorf("failed to delete article group: %w", err)
	}

	formatter := getFormatter(cmd.Format)
	if cmd.Format == "quiet" {
		return formatter.Print(cmd.ID)
	}
	if cmd.Format == "json" {
		return formatter.PrintJSON(map[string]string{"deletedId": cmd.ID})
	}

	return formatter.Success(fmt.Sprintf("Article group deleted: %s", cmd.ID))
}
