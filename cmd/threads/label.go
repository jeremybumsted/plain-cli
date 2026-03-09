package threads

// LabelCmd represents the threads label command group
type LabelCmd struct {
	Add     LabelAddCmd     `cmd:"" help:"Add labels to a thread"`
	Remove  LabelRemoveCmd  `cmd:"" help:"Remove labels from a thread"`
	List    LabelListCmd    `cmd:"" help:"List available label types"`
	Refresh LabelRefreshCmd `cmd:"" help:"Refresh label cache from API"`
}
