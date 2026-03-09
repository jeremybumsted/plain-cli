package threads

// FieldCmd represents the threads field command group
type FieldCmd struct {
	List FieldListCmd `cmd:"" help:"List available thread field schemas"`
}
