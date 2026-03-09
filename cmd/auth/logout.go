package auth

import (
	"fmt"
)

// LogoutCmd handles logging out
type LogoutCmd struct {
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
}

// Run executes the logout command
func (cmd *LogoutCmd) Run() error {
	formatter := getFormatter(cmd.Format)

	// Load config
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if already logged out
	if !cfg.IsAuthenticated() {
		formatter.Info("Already logged out")
		return nil
	}

	// Confirm logout (unless -y flag is used)
	if !cmd.Yes {
		formatter.Warning("This will remove your stored credentials.")
		fmt.Print("Continue? [y/N]: ")

		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" && response != "yes" {
			formatter.Info("Logout cancelled")
			return nil
		}
	}

	// Clear credentials
	cfg.Clear()

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	formatter.Success("Successfully logged out")
	formatter.Info("\nTo log back in, run:")
	formatter.Info("  plain auth login")

	return nil
}
