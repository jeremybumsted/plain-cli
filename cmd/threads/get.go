package threads

import (
	"fmt"
	"strings"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// GetCmd represents the threads get command
type GetCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	Timeline   bool   `help:"Include timeline entries" default:"false"`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads get command
func (cmd *GetCmd) Run() error {
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

	// Fetch thread
	thread, err := client.GetThread(threadID, cmd.Timeline)
	if err != nil {
		// Handle 404 gracefully
		if mcpErr, ok := err.(*plain.Error); ok && mcpErr.StatusCode == 404 {
			return formatter.Error(fmt.Sprintf("Thread not found: %s", threadID))
		}
		return fmt.Errorf("failed to fetch thread: %w", err)
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		return formatter.Print(thread.ID)
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(thread)
	}

	// Default: key-value format
	return displayThread(formatter, thread, cmd.Timeline)
}

// displayThread displays thread details in a human-readable format
func displayThread(formatter interface{}, thread *plain.Thread, includeTimeline bool) error {
	type Formatter interface {
		Print(message string) error
		Printf(format string, args ...interface{}) error
	}
	f := formatter.(Formatter)

	// Display thread details
	f.Print("Thread Details")
	f.Print("--------------")
	f.Printf("ID:          %s\n", thread.ID)
	f.Printf("Title:       %s\n", thread.Title)
	f.Printf("Status:      %s\n", thread.Status)
	f.Printf("Priority:    %s\n", plain.FormatPriority(thread.Priority))

	// Display assignee
	if thread.AssignedTo != nil {
		f.Printf("Assignee:    %s (%s)\n", thread.AssignedTo.Email, thread.AssignedTo.FullName)
	} else {
		f.Printf("Assignee:    Unassigned\n")
	}

	// Display labels
	if len(thread.Labels) > 0 {
		labelNames := make([]string, 0, len(thread.Labels))
		for _, label := range thread.Labels {
			if label.LabelType != nil {
				labelNames = append(labelNames, label.LabelType.Name)
			}
		}
		if len(labelNames) > 0 {
			f.Printf("Labels:      %s\n", strings.Join(labelNames, ", "))
		} else {
			f.Printf("Labels:      None\n")
		}
	} else {
		f.Printf("Labels:      None\n")
	}

	// Format timestamps
	createdAt, _ := thread.CreatedAt.Time()
	updatedAt, _ := thread.UpdatedAt.Time()
	f.Printf("Created:     %s\n", createdAt.Format("2006-01-02 15:04:05"))
	f.Printf("Updated:     %s\n", updatedAt.Format("2006-01-02 15:04:05"))

	// Display description if present
	if thread.Description != "" {
		f.Print("Description:")
		// Indent description
		lines := strings.Split(thread.Description, "\n")
		for _, line := range lines {
			f.Printf("  %s\n", line)
		}
	}

	// Display timeline if requested and available
	if includeTimeline && thread.Timeline != nil && len(thread.Timeline.Entries) > 0 {
		f.Print("\nTimeline")
		f.Print("--------")
		for _, entry := range thread.Timeline.Entries {
			timestamp, _ := entry.Timestamp.Time()
			timestampStr := timestamp.Format("2006-01-02 15:04:05")
			actorInfo := "System"
			if entry.Actor != nil {
				actorInfo = entry.Actor.Email
			}

			entryType := entry.Entry.Type
			f.Printf("%s - %s by %s\n", timestampStr, entryType, actorInfo)

			// Display entry content if available
			content := ""
			if entry.Entry.Text != "" {
				content = entry.Entry.Text
			} else if entry.Entry.Content != "" {
				content = entry.Entry.Content
			}

			if content != "" {
				// Indent content and truncate if too long
				lines := strings.Split(content, "\n")
				for i, line := range lines {
					if i >= 3 {
						f.Print("  ...")
						break
					}
					if len(line) > 100 {
						line = line[:97] + "..."
					}
					f.Printf("  %s\n", line)
				}
			}
		}
	}

	return nil
}
