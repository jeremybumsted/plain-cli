package helpcenters

import "fmt"

// GetCmd represents the help-centers get command
type GetCmd struct {
	ID         string `arg:"" name:"id" help:"Help center ID"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the help-centers get command
func (cmd *GetCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	helpCenter, err := client.GetHelpCenter(cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch help center: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(helpCenter)
	case "quiet":
		fmt.Println(helpCenter.ID)
		return nil
	default:
		isDeleted := "no"
		if helpCenter.IsDeleted {
			isDeleted = "yes"
		}

		fmt.Printf("ID:            %s\n", helpCenter.ID)
		if helpCenter.Type != "" {
			fmt.Printf("Type:          %s\n", helpCenter.Type)
		}
		if helpCenter.PublicName != "" {
			fmt.Printf("Public Name:   %s\n", helpCenter.PublicName)
		}
		if helpCenter.InternalName != "" {
			fmt.Printf("Internal Name: %s\n", helpCenter.InternalName)
		}
		if helpCenter.Description != "" {
			fmt.Printf("Description:   %s\n", helpCenter.Description)
		}
		fmt.Printf("Deleted:       %s\n", isDeleted)

		createdAt, _ := helpCenter.CreatedAt.Time()
		updatedAt, _ := helpCenter.UpdatedAt.Time()
		fmt.Printf("Created:       %s\n", createdAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated:       %s\n", updatedAt.Format("2006-01-02 15:04:05"))

		return nil
	}
}
