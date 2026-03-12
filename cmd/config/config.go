package config

import (
	"fmt"

	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// ConfigCmd represents the config command group
type ConfigCmd struct {
	Show  ShowCmd  `cmd:"" help:"Show current configuration"`
	Reset ResetCmd `cmd:"" help:"Reset configuration"`
}

// Run executes the interactive configuration
func (cmd *ConfigCmd) Run() error {
	cfg, err := getConfig("")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if this is first-time setup (no auth or incomplete config)
	if !cfg.IsAuthenticated() {
		fmt.Println("Welcome to Plain CLI!")
		fmt.Println("Let's get you set up.")
		fmt.Println()
		return runFirstTimeSetup(cfg)
	}

	// Check if config is incomplete
	if !cfg.IsFullyConfigured() {
		fmt.Println("Configuration is incomplete.")
		fmt.Println("Let's complete your setup.")
		fmt.Println()
		return runFirstTimeSetup(cfg)
	}

	// Show current config (default behavior)
	return (&ShowCmd{}).Run()
}

// Helper functions
func getConfig(configPath string) (*config.Config, error) {
	return config.Load(configPath)
}

func getClient(cfg *config.Config) (*plain.Client, error) {
	token, err := cfg.GetToken()
	if err != nil {
		return nil, err
	}
	return plain.NewClient(token), nil
}

func runConfigMenu(cfg *config.Config) error {
	fmt.Println("Plain CLI Configuration")
	fmt.Println()

	// Show current settings
	showCurrentSettings(cfg)

	fmt.Println("\nWhat would you like to configure?")
	fmt.Println("  1. API Token")
	fmt.Println("  2. Workspace")
	fmt.Println("  3. Help Center")
	fmt.Println("  4. Configure All")
	fmt.Println("  0. Exit")
	fmt.Println()

	var selection string
	fmt.Print("Select option [0-4]: ")
	if _, err := fmt.Scanln(&selection); err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	switch selection {
	case "1":
		return configureAPIToken(cfg)
	case "2":
		return configureWorkspace(cfg)
	case "3":
		return configureHelpCenter(cfg)
	case "4":
		return runFirstTimeSetup(cfg)
	case "0":
		return nil
	default:
		fmt.Println("Invalid selection")
		return nil
	}
}

func showCurrentSettings(cfg *config.Config) {
	fmt.Println("Current configuration:")

	// API Token
	if cfg.IsAuthenticated() {
		token := cfg.AccessToken
		if len(token) > 8 {
			token = "..." + token[len(token)-4:]
		}
		fmt.Printf("  API Token: %s", token)
		if !cfg.ExpiresAt.IsZero() {
			fmt.Printf(" (expires: %s)", cfg.ExpiresAt.Format("2006-01-02"))
		}
		fmt.Println()
	} else {
		fmt.Println("  API Token: Not configured")
	}

	// Workspace
	wsID, err := cfg.GetWorkspaceID()
	if err == nil {
		fmt.Printf("  Workspace: %s\n", wsID)
	} else {
		fmt.Println("  Workspace: Not configured")
	}

	// Help Center
	hcID, err := cfg.GetHelpCenterID()
	if err == nil {
		fmt.Printf("  Help Center: %s\n", hcID)
	} else {
		fmt.Println("  Help Center: Not configured")
	}
}
