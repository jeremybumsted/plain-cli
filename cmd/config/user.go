package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// UserCmd configures the user identity for personalized queries
type UserCmd struct {
	Email string `flag:"" help:"Your email address in Plain"`
	List  bool   `flag:"" help:"List all workspace users"`
	Clear bool   `flag:"" help:"Clear configured user info"`
}

func (cmd *UserCmd) Run() error {
	cfg, err := getConfig("")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Handle clear
	if cmd.Clear {
		cfg.ClearUserInfo()
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("User information cleared")
		return nil
	}

	// Get client
	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	// Handle list
	if cmd.List {
		return listUsers(client)
	}

	// Configure user by email
	email := cmd.Email
	if email == "" {
		// Prompt for email
		fmt.Print("Enter your email address: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf("failed to read email")
		}
		email = strings.TrimSpace(scanner.Text())
	}

	// Lookup user
	fmt.Printf("Looking up user with email: %s\n", email)
	user, err := client.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Save user info
	if err := cfg.SetUserInfo(user.ID, user.Email, user.FullName, user.PublicName); err != nil {
		return fmt.Errorf("failed to save user info: %w", err)
	}

	fmt.Printf("✓ User configured: %s (%s)\n", user.FullName, user.Email)
	fmt.Printf("  User ID: %s\n", user.ID)

	return nil
}

func listUsers(client *plain.Client) error {
	fmt.Println("Fetching workspace users...")

	users, err := client.ListUsers(true) // Only assignable users
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	fmt.Printf("\nFound %d users:\n\n", len(users))
	for _, user := range users {
		fmt.Printf("  %s (%s)\n", user.FullName, user.Email)
		fmt.Printf("    ID: %s\n\n", user.ID)
	}

	return nil
}
