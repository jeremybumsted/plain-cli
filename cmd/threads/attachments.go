package threads

import (
	"fmt"
	"time"

	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// AttachmentsCmd represents the threads attachments command group
type AttachmentsCmd struct {
	List     ListAttachmentsCmd     `cmd:"" help:"List all attachments in a thread"`
	Download DownloadAttachmentCmd  `cmd:"" help:"Download a specific attachment"`
}

// ListAttachmentsCmd lists all attachments in a thread
type ListAttachmentsCmd struct {
	ThreadID   string `arg:"" help:"Thread ID or URL" required:""`
	ConfigPath string `help:"Path to config file" default:""`
	Format     string `help:"Output format (table, json, quiet)" default:"table"`
}

// Run executes the threads attachments list command
func (cmd *ListAttachmentsCmd) Run() error {
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

	// Fetch thread with timeline to get attachments
	thread, err := client.GetThread(threadID, true) // includeTimeline=true
	if err != nil {
		return fmt.Errorf("failed to fetch thread: %w", err)
	}

	// Collect all attachments from timeline entries
	var allAttachments []AttachmentInfo
	if thread.Timeline != nil {
		for _, entry := range thread.Timeline.Entries {
			for _, att := range entry.Entry.Attachments {
				timestamp, _ := entry.Timestamp.Time()
				allAttachments = append(allAttachments, AttachmentInfo{
					Attachment: att,
					EntryID:    entry.ID,
					EntryType:  entry.Entry.GetEntryDisplayText(),
					Timestamp:  timestamp,
				})
			}
		}
	}

	if len(allAttachments) == 0 {
		return formatter.Print("No attachments found in thread")
	}

	// Handle different output formats
	if cmd.Format == "quiet" {
		// Just output attachment IDs
		for _, att := range allAttachments {
			if err := formatter.Print(att.Attachment.ID); err != nil {
				return err
			}
		}
		return nil
	}

	if cmd.Format == "json" {
		return formatter.PrintJSON(allAttachments)
	}

	// Table format
	return displayAttachmentTable(formatter, allAttachments)
}

// AttachmentInfo combines attachment with context information
type AttachmentInfo struct {
	Attachment plain.Attachment `json:"attachment"`
	EntryID    string           `json:"entryId"`
	EntryType  string           `json:"entryType"`
	Timestamp  time.Time        `json:"timestamp"`
}

// DownloadAttachmentCmd downloads a specific attachment
type DownloadAttachmentCmd struct {
	AttachmentID string `arg:"" help:"Attachment ID" required:""`
	Output       string `help:"Output file path (default: use original filename)" optional:""`
	ConfigPath   string `help:"Path to config file" default:""`
}

// Run executes the threads attachments download command
func (cmd *DownloadAttachmentCmd) Run() error {
	// Load config and get client
	cfg, err := getConfig(cmd.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client, err := getClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Get download URL first to get filename
	downloadInfo, err := client.CreateAttachmentDownloadUrl(cmd.AttachmentID)
	if err != nil {
		return fmt.Errorf("failed to get attachment info: %w", err)
	}

	// Determine output path
	outputPath := cmd.Output
	if outputPath == "" {
		outputPath = downloadInfo.Attachment.FileName
	}

	fmt.Printf("Downloading %s...\n", downloadInfo.Attachment.FileName)

	// Download the file
	if err := client.DownloadAttachment(cmd.AttachmentID, outputPath); err != nil {
		return fmt.Errorf("failed to download attachment: %w", err)
	}

	fmt.Printf("Downloaded to %s (%s)\n",
		outputPath,
		formatFileSize(downloadInfo.Attachment.FileSize.Bytes))

	return nil
}

// displayAttachmentTable displays attachments in a table format
func displayAttachmentTable(formatter interface{}, attachments []AttachmentInfo) error {
	type Formatter interface {
		Print(message string) error
		Printf(format string, args ...interface{}) error
	}
	f := formatter.(Formatter)

	// Print header
	if err := f.Print("Attachments"); err != nil {
		return err
	}
	if err := f.Print("-----------"); err != nil {
		return err
	}
	if err := f.Print(""); err != nil {
		return err
	}

	for _, info := range attachments {
		att := info.Attachment
		if err := f.Printf("ID:       %s\n", att.ID); err != nil {
			return err
		}
		if err := f.Printf("Filename: %s\n", att.FileName); err != nil {
			return err
		}
		if err := f.Printf("Size:     %s\n", formatFileSize(att.FileSize.Bytes)); err != nil {
			return err
		}
		if att.FileExtension != "" {
			if err := f.Printf("Type:     %s\n", att.FileExtension); err != nil {
				return err
			}
		}
		if err := f.Printf("From:     %s (%s)\n", info.EntryType, info.Timestamp.Format("2006-01-02 15:04")); err != nil {
			return err
		}
		if err := f.Print(""); err != nil {
			return err
		}
	}

	return nil
}
