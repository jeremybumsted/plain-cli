package config

import (
	"fmt"
)

// ShowCmd shows current configuration
type ShowCmd struct{}

// Run executes the show command
func (cmd *ShowCmd) Run() error {
	cfg, err := getConfig("")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Plain CLI Configuration")
	fmt.Println()
	showCurrentSettings(cfg)
	fmt.Println()
	fmt.Printf("Config file: %s\n", cfg.GetConfigPath())

	return nil
}
