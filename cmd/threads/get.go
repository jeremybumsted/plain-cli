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
	if err := f.Print("Thread Details"); err != nil {
		return err
	}
	if err := f.Print("--------------"); err != nil {
		return err
	}
	if err := f.Printf("ID:          %s\n", thread.ID); err != nil {
		return err
	}
	if err := f.Printf("Title:       %s\n", thread.Title); err != nil {
		return err
	}
	if err := f.Printf("Status:      %s\n", thread.Status); err != nil {
		return err
	}
	if err := f.Printf("Priority:    %s\n", plain.FormatPriority(thread.Priority)); err != nil {
		return err
	}

	// Display assignee
	if thread.AssignedTo != nil {
		if err := f.Printf("Assignee:    %s (%s)\n", thread.AssignedTo.Email, thread.AssignedTo.FullName); err != nil {
			return err
		}
	} else {
		if err := f.Printf("Assignee:    Unassigned\n"); err != nil {
			return err
		}
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
			if err := f.Printf("Labels:      %s\n", strings.Join(labelNames, ", ")); err != nil {
				return err
			}
		} else {
			if err := f.Printf("Labels:      None\n"); err != nil {
				return err
			}
		}
	} else {
		if err := f.Printf("Labels:      None\n"); err != nil {
			return err
		}
	}

	// Format timestamps
	createdAt, _ := thread.CreatedAt.Time()
	updatedAt, _ := thread.UpdatedAt.Time()
	if err := f.Printf("Created:     %s\n", createdAt.Format("2006-01-02 15:04:05")); err != nil {
		return err
	}
	if err := f.Printf("Updated:     %s\n", updatedAt.Format("2006-01-02 15:04:05")); err != nil {
		return err
	}

	// Display description if present
	if thread.Description != "" {
		if err := f.Print("Description:"); err != nil {
			return err
		}
		// Indent description
		lines := strings.Split(thread.Description, "\n")
		for _, line := range lines {
			if err := f.Printf("  %s\n", line); err != nil {
				return err
			}
		}
	}

	// Display timeline if requested and available
	if includeTimeline && thread.Timeline != nil && len(thread.Timeline.Entries) > 0 {
		if err := f.Print("\nTimeline"); err != nil {
			return err
		}
		if err := f.Print("--------"); err != nil {
			return err
		}
		for _, entry := range thread.Timeline.Entries {
			timestamp, _ := entry.Timestamp.Time()
			timestampStr := timestamp.Format("2006-01-02 15:04:05")
			actorInfo := "System"
			if entry.Actor != nil {
				actorInfo = entry.Actor.Email
			}

			entryType := entry.Entry.Type
			if err := f.Printf("%s - %s by %s\n", timestampStr, entryType, actorInfo); err != nil {
				return err
			}

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
						if err := f.Print("  ..."); err != nil {
							return err
						}
						break
					}
					if len(line) > 100 {
						line = line[:97] + "..."
					}
					if err := f.Printf("  %s\n", line); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
