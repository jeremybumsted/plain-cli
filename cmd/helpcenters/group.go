package helpcenters

// GroupCmd represents the help-centers group command group.
type GroupCmd struct {
	Get    GroupGetCmd    `cmd:"" help:"Get article group details"`
	Create GroupCreateCmd `cmd:"" help:"Create an article group"`
	Update GroupUpdateCmd `cmd:"" help:"Update an article group"`
	Delete GroupDeleteCmd `cmd:"" help:"Delete an article group"`
}
