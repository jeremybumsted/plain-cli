package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/config"
)

func configureWorkspace(cfg *config.Config) error {
	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("not authenticated: %w", err)
	}

	fmt.Println("[2/3] Workspace")
	fmt.Println("Fetching workspaces...")

	workspaces, err := client.ListWorkspaces()
	if err != nil {
		return fmt.Errorf("failed to fetch workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		return fmt.Errorf("no workspaces found for your account")
	}

	if len(workspaces) == 1 {
		// Auto-select if only one workspace
		if err := cfg.SetWorkspaceID(workspaces[0].ID); err != nil {
			return fmt.Errorf("failed to save workspace: %w", err)
		}
		fmt.Printf("✓ Workspace configured: %s (%s)\n", workspaces[0].Name, workspaces[0].ID)
		fmt.Println()
		return nil
	}

	// Multiple workspaces - let user choose
	fmt.Println("\nAvailable workspaces:")
	for i, ws := range workspaces {
		fmt.Printf("  %d. %s (%s)\n", i+1, ws.Name, ws.ID)
	}
	fmt.Println()

	// Prompt for selection
	var selection int
	for {
		fmt.Printf("Select workspace [1-%d]: ", len(workspaces))
		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		selection, err = strconv.Atoi(strings.TrimSpace(input))
		if err != nil || selection < 1 || selection > len(workspaces) {
			fmt.Printf("Invalid selection. Please enter a number between 1 and %d.\n", len(workspaces))
			continue
		}
		break
	}

	// Save selection
	selectedWS := workspaces[selection-1]
	if err := cfg.SetWorkspaceID(selectedWS.ID); err != nil {
		return fmt.Errorf("failed to save workspace: %w", err)
	}

	fmt.Printf("\n✓ Workspace configured: %s (%s)\n", selectedWS.Name, selectedWS.ID)
	fmt.Println()

	return nil
}
