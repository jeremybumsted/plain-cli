package tiers

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// GetCmd represents the tiers get command
type GetCmd struct {
	ID         string `arg:"" name:"id" help:"Tier ID"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the tiers get command
func (cmd *GetCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	tier, err := client.GetTier(cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch tier: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(tier)
	case "quiet":
		fmt.Println(tier.ID)
		return nil
	default:
		return displayTier(tier)
	}
}

func displayTier(tier *plain.Tier) error {
	isDefault := "no"
	if tier.IsDefault {
		isDefault = "yes"
	}
	isMachine := "no"
	if tier.IsMachineTier {
		isMachine = "yes"
	}

	fmt.Printf("ID:               %s\n", tier.ID)
	fmt.Printf("Name:             %s\n", tier.Name)
	if tier.ExternalID != "" {
		fmt.Printf("External ID:      %s\n", tier.ExternalID)
	}
	if tier.Color != "" {
		fmt.Printf("Color:            %s\n", tier.Color)
	}
	fmt.Printf("Default:          %s\n", isDefault)
	fmt.Printf("Machine Tier:     %s\n", isMachine)
	fmt.Printf("Default Priority: %s\n", plain.FormatPriority(tier.DefaultThreadPriority))

	createdAt, _ := tier.CreatedAt.Time()
	updatedAt, _ := tier.UpdatedAt.Time()
	fmt.Printf("Created:          %s\n", createdAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:          %s\n", updatedAt.Format("2006-01-02 15:04:05"))

	return nil
}
