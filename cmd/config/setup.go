package config

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func runFirstTimeSetup(cfg *config.Config) error {
	// Step 1: API Token (if not authenticated)
	if !cfg.IsAuthenticated() {
		if err := configureAPIToken(cfg); err != nil {
			return err
		}
	}

	// Step 2: Workspace
	if _, err := cfg.GetWorkspaceID(); err != nil {
		if err := configureWorkspace(cfg); err != nil {
			return err
		}
	}

	// Step 3: Help Center
	if _, err := cfg.GetHelpCenterID(); err != nil {
		if err := configureHelpCenter(cfg); err != nil {
			return err
		}
	}

	fmt.Println("\n✓ Configuration complete!")
	fmt.Printf("Config saved to: %s\n\n", cfg.GetConfigPath())
	fmt.Println("You can now use Plain CLI commands:")
	fmt.Println("  plain articles list")
	fmt.Println("  plain threads list")

	return nil
}
