package config

import (
	"fmt"
	"os"
)

// ResetCmd resets configuration
type ResetCmd struct {
	Yes bool `help:"Skip confirmation" short:"y"`
}

// Run executes the reset command
func (cmd *ResetCmd) Run() error {
	cfg, err := getConfig("")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Confirm unless --yes flag is used
	if !cmd.Yes {
		fmt.Print("Are you sure you want to reset configuration? [y/N]: ")
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	// Clear all configuration
	cfg.Clear()
	cfg.WorkspaceID = ""
	cfg.HelpCenterID = ""

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Also try to delete the config file
	configPath := cfg.GetConfigPath()
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: Could not delete config file: %v\n", err)
	}

	fmt.Println("✓ Configuration reset")
	fmt.Println("Run 'plain config' to set up again")

	return nil
}
