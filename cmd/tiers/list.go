package tiers

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// ListCmd represents the tiers list command
type ListCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the tiers list command
func (cmd *ListCmd) Run() error {
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	tiers, err := client.ListTiers()
	if err != nil {
		return fmt.Errorf("failed to list tiers: %w", err)
	}

	formatter := getFormatter(cmd.Format)

	switch cmd.Format {
	case "json":
		return formatter.PrintJSON(tiers)
	case "quiet":
		for _, tier := range tiers {
			fmt.Println(tier.ID)
		}
		return nil
	default:
		return outputTable(tiers)
	}
}

func outputTable(tiers []*plain.Tier) error {
	if len(tiers) == 0 {
		fmt.Println("No tiers found")
		return nil
	}

	fmt.Printf("%-20s %-25s %-10s %-10s %-10s %s\n", "ID", "NAME", "DEFAULT", "MACHINE", "PRIORITY", "UPDATED")
	fmt.Printf("%-20s %-25s %-10s %-10s %-10s %s\n", "----", "----", "-------", "-------", "--------", "-------")

	for _, tier := range tiers {
		isDefault := "no"
		if tier.IsDefault {
			isDefault = "yes"
		}
		isMachine := "no"
		if tier.IsMachineTier {
			isMachine = "yes"
		}
		priority := plain.FormatPriority(tier.DefaultThreadPriority)
		updated := tier.UpdatedAt.ISO8601
		if len(updated) > 10 {
			updated = updated[:10]
		}

		fmt.Printf("%-20s %-25s %-10s %-10s %-10s %s\n", tier.ID, tier.Name, isDefault, isMachine, priority, updated)
	}

	fmt.Printf("\nTotal: %d tiers\n", len(tiers))
	return nil
}
