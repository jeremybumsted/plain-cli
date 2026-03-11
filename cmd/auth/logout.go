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
		return formatter.Info("Already logged out")
	}

	// Confirm logout (unless -y flag is used)
	if !cmd.Yes {
		if err := formatter.Warning("This will remove your stored credentials."); err != nil {
			return err
		}
		fmt.Print("Continue? [y/N]: ")

		var response string
		_, _ = fmt.Scanln(&response)

		if response != "y" && response != "Y" && response != "yes" {
			return formatter.Info("Logout cancelled")
		}
	}

	// Clear credentials
	cfg.Clear()

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := formatter.Success("Successfully logged out"); err != nil {
		return err
	}
	if err := formatter.Info("\nTo log back in, run:"); err != nil {
		return err
	}
	if err := formatter.Info("  plain auth login"); err != nil {
		return err
	}

	return nil
}
