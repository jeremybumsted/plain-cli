package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func configureHelpCenter(cfg *config.Config) error {
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: %w", err)
	}

	fmt.Println("[3/3] Help Center")
	fmt.Println("Fetching help centers...")

	helpCenters, err := client.ListHelpCenters()
	if err != nil {
		return fmt.Errorf("failed to fetch help centers: %w", err)
	}

	if len(helpCenters) == 0 {
		return fmt.Errorf("no help centers found in your workspace")
	}

	// Display options
	fmt.Println("\nAvailable help centers:")
	for i, hc := range helpCenters {
		fmt.Printf("  %d. %s\n", i+1, hc.ID)
	}
	fmt.Println()

	// Prompt for selection
	var selection int
	for {
		fmt.Printf("Select help center [1-%d]: ", len(helpCenters))
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		selection, err = strconv.Atoi(strings.TrimSpace(input))
		if err != nil || selection < 1 || selection > len(helpCenters) {
			fmt.Printf("Invalid selection. Please enter a number between 1 and %d.\n", len(helpCenters))
			continue
		}
		break
	}

	// Save selection
	selectedHC := helpCenters[selection-1]
	if err := cfg.SetHelpCenterID(selectedHC.ID); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("\n✓ Help center configured: %s\n", selectedHC.ID)
	fmt.Println()

	return nil
}
