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
			actorName := entry.Actor.GetActorDisplayName()
			entryType := entry.Entry.GetEntryDisplayText()

			// Display entry header
			if err := f.Printf("%s - %s by %s\n", timestampStr, entryType, actorName); err != nil {
				return err
			}

			// Display entry details
			details := formatEntryDetails(&entry.Entry)
			if details != "" {
				lines := strings.Split(details, "\n")
				for i, line := range lines {
					if i >= 5 { // Limit to 5 lines
						if err := f.Print("  ..."); err != nil {
							return err
						}
						break
					}
					if len(line) > 120 {
						line = line[:117] + "..."
					}
					if err := f.Printf("  %s\n", line); err != nil {
						return err
					}
				}
			}

			// Display attachments
			if len(entry.Entry.Attachments) > 0 {
				if err := f.Printf("  Attachments (%d):\n", len(entry.Entry.Attachments)); err != nil {
					return err
				}
				for _, att := range entry.Entry.Attachments {
					sizeStr := formatFileSize(att.FileSize.Bytes)
					if err := f.Printf("    - %s (%s)\n", att.FileName, sizeStr); err != nil {
						return err
					}
				}
			}

			// Add blank line between entries
			if err := f.Print(""); err != nil {
				return err
			}
		}
	}

	return nil
}

// formatFileSize formats bytes into a human-readable file size
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatEntryDetails returns formatted details for a timeline entry based on its type
func formatEntryDetails(entry *plain.EntryData) string {
	switch entry.Typename {
	case "NoteEntry":
		if entry.Markdown != "" {
			return entry.Markdown
		}
		if entry.NoteText != "" {
			return entry.NoteText
		}
		return entry.Text
	case "ChatEntry":
		text := entry.ChatText
		if text == "" {
			text = entry.Text
		}
		details := text
		if entry.CustomerReadAt != nil {
			readTime, _ := entry.CustomerReadAt.Time()
			details += fmt.Sprintf("\n(Read by customer at %s)", readTime.Format("2006-01-02 15:04:05"))
		}
		return details
	case "EmailEntry":
		details := fmt.Sprintf("Subject: %s", entry.Subject)
		if entry.From != nil {
			details += fmt.Sprintf("\nFrom: %s <%s>", entry.From.Name, entry.From.Email)
		}
		if entry.To != nil {
			details += fmt.Sprintf("\nTo: %s <%s>", entry.To.Name, entry.To.Email)
		}
		return details
	case "ThreadStatusTransitionedEntry":
		return fmt.Sprintf("%s → %s", entry.PreviousStatus, entry.NextStatus)
	case "ThreadAssignmentTransitionedEntry":
		prevName := "Unassigned"
		if entry.PreviousAssignee != nil && entry.PreviousAssignee.Email != "" {
			prevName = entry.PreviousAssignee.Email
		}
		nextName := "Unassigned"
		if entry.NextAssignee != nil && entry.NextAssignee.Email != "" {
			nextName = entry.NextAssignee.Email
		}
		return fmt.Sprintf("%s → %s", prevName, nextName)
	case "ThreadPriorityChangedEntry":
		prevPri := "-"
		if entry.PreviousPriority != nil {
			prevPri = plain.FormatPriority(*entry.PreviousPriority)
		}
		nextPri := "-"
		if entry.NextPriority != nil {
			nextPri = plain.FormatPriority(*entry.NextPriority)
		}
		return fmt.Sprintf("%s → %s", prevPri, nextPri)
	case "SlackMessageEntry":
		text := entry.SlackText
		if text == "" {
			text = entry.Text
		}
		details := text
		if entry.SlackWebMessageLink != "" {
			details += fmt.Sprintf("\nLink: %s", entry.SlackWebMessageLink)
		}
		return details
	case "SlackReplyEntry":
		text := entry.SlackReplyText
		if text == "" {
			text = entry.Text
		}
		details := text
		if entry.SlackWebMessageLink != "" {
			details += fmt.Sprintf("\nLink: %s", entry.SlackWebMessageLink)
		}
		return details
	case "ThreadDiscussionMessageEntry":
		if entry.DiscussionText != "" {
			return entry.DiscussionText
		}
		return entry.Text
	case "ThreadDiscussionResolvedEntry":
		return ""
	case "ThreadDiscussionEntry":
		return ""
	case "ThreadLabelsChangedEntry":
		var parts []string
		if len(entry.AddedLabelTypes) > 0 {
			added := make([]string, 0, len(entry.AddedLabelTypes))
			for _, labelType := range entry.AddedLabelTypes {
				added = append(added, labelType.Name)
			}
			if len(added) > 0 {
				parts = append(parts, fmt.Sprintf("Added: %s", strings.Join(added, ", ")))
			}
		}
		if len(entry.RemovedLabelTypes) > 0 {
			removed := make([]string, 0, len(entry.RemovedLabelTypes))
			for _, labelType := range entry.RemovedLabelTypes {
				removed = append(removed, labelType.Name)
			}
			if len(removed) > 0 {
				parts = append(parts, fmt.Sprintf("Removed: %s", strings.Join(removed, ", ")))
			}
		}
		return strings.Join(parts, "\n")
	case "ThreadAdditionalAssigneesTransitionedEntry":
		var parts []string
		// Calculate who was added by finding who's in next but not previous
		addedMap := make(map[string]string)
		for _, next := range entry.NextAssignees {
			found := false
			for _, prev := range entry.PreviousAssignees {
				if next.ID == prev.ID {
					found = true
					break
				}
			}
			if !found {
				addedMap[next.ID] = next.FullName
			}
		}
		// Calculate who was removed by finding who's in previous but not next
		removedMap := make(map[string]string)
		for _, prev := range entry.PreviousAssignees {
			found := false
			for _, next := range entry.NextAssignees {
				if prev.ID == next.ID {
					found = true
					break
				}
			}
			if !found {
				removedMap[prev.ID] = prev.FullName
			}
		}

		if len(addedMap) > 0 {
			added := make([]string, 0, len(addedMap))
			for _, name := range addedMap {
				added = append(added, name)
			}
			parts = append(parts, fmt.Sprintf("Added: %s", strings.Join(added, ", ")))
		}
		if len(removedMap) > 0 {
			removed := make([]string, 0, len(removedMap))
			for _, name := range removedMap {
				removed = append(removed, name)
			}
			parts = append(parts, fmt.Sprintf("Removed: %s", strings.Join(removed, ", ")))
		}
		return strings.Join(parts, "\n")
	case "ThreadLinkCreatedEntry":
		if entry.Thread != nil {
			return fmt.Sprintf("Linked to: %s (%s)", entry.Thread.Title, entry.Thread.ID)
		}
		return ""
	case "CustomEntry":
		return entry.Title
	default:
		return ""
	}
}
