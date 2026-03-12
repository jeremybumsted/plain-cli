package threads

import (
	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/plain"
	"github.com/jeremybumsted/plain-cli/internal/output"
)

// ThreadsCmd represents the threads command group
type ThreadsCmd struct {
	List        ListCmd        `cmd:"" help:"List threads with filters"`
	Get         GetCmd         `cmd:"" help:"Get thread details"`
	Search      SearchCmd      `cmd:"" help:"Search threads by text"`
	Assign      AssignCmd      `cmd:"" help:"Assign thread to a user"`
	Done        DoneCmd        `cmd:"" help:"Mark a thread as done"`
	Todo        TodoCmd        `cmd:"" help:"Mark thread as todo"`
	Note        NoteCmd        `cmd:"" help:"Add a note to a thread"`
	Unassign    UnassignCmd    `cmd:"" help:"Unassign a thread"`
	Snooze      SnoozeCmd      `cmd:"" help:"Snooze a thread"`
	Priority    PriorityCmd    `cmd:"" help:"Change thread priority"`
	Label       LabelCmd       `cmd:"" help:"Manage thread labels"`
	Field       FieldCmd       `cmd:"" help:"View thread field schemas"`
	Attachments AttachmentsCmd `cmd:"" help:"List and download thread attachments"`
}

// Helper functions for thread commands

// getConfig loads the configuration from the specified path
func getConfig(configPath string) (*config.Config, error) {
	return config.Load(configPath)
}

// getClient creates an authenticated Plain API client
func getClient(cfg *config.Config) (*plain.Client, error) {
	token, err := cfg.GetToken()
	if err != nil {
		return nil, err
	}
	return plain.NewClient(token), nil
}

// getFormatter creates an output formatter based on the format string
func getFormatter(format string) *output.Formatter {
	var fmt output.Format
	switch format {
	case "json":
		fmt = output.FormatJSON
	case "quiet":
		fmt = output.FormatQuiet
	default:
		fmt = output.FormatTable
	}
	return output.New(fmt)
}
