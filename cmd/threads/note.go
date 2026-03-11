package threads

import (
	"fmt"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// NoteCmd represents the threads note command
type NoteCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	Text       string `help:"Note text (if not provided, opens editor)" default:""`
	Yes        bool   `help:"Skip confirmation prompt" short:"y"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads note command
func (cmd *NoteCmd) Run() error {
	// Extract thread ID from URL if provided
	threadID := extractThreadID(cmd.ThreadID)

	// Load config and get client
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Get formatter
	formatter := getFormatter(cmd.Format)

	// Get note text - either from flag or editor
	noteText := strings.TrimSpace(cmd.Text)
	if noteText == "" {
		// Open editor for text input
		editorText, err := openEditor("")
		if err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}
		noteText = strings.TrimSpace(editorText)
	}

	// Validate that we have text
	if noteText == "" {
		return fmt.Errorf("note text cannot be empty")
	}

	// Confirm action unless --yes flag is set
	confirmed, err := confirmAction("Add note to thread?", cmd.Yes)
	if err != nil {
		return fmt.Errorf("failed to read confirmation: %w", err)
	}
	if !confirmed {
		return formatter.Error("Cancelled")
	}

	// Create the note
	note, err := client.CreateNote(threadID, noteText)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*plain.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to create note: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		// In quiet mode, show note ID
		return formatter.Print(note.ID)
	}

	if cmd.Format == "json" {
		// In JSON mode, show full note object with ID and timestamp
		return formatter.PrintJSON(note)
	}

	// Default: success message with details
	timestamp, _ := note.CreatedAt.Time()
	timestampStr := timestamp.Format("2006-01-02 15:04:05")

	if err := formatter.Print("Note added to thread"); err != nil {
		return err
	}
	if err := formatter.Printf("ID:        %s\n", note.ID); err != nil {
		return err
	}
	if err := formatter.Printf("Created:   %s\n", timestampStr); err != nil {
		return err
	}
	if note.CreatedBy.Email != "" {
		if err := formatter.Printf("Author:    %s\n", note.CreatedBy.Email); err != nil {
			return err
		}
	}

	return nil
}
