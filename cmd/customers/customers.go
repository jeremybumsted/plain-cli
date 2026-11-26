package customers

import (
	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/output"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// CustomersCmd represents the customers command group
type CustomersCmd struct {
	List       ListCmd       `cmd:"" help:"List customers"`
	Get        GetCmd        `cmd:"" help:"Get customer by ID"`
	GetByEmail GetByEmailCmd `cmd:"" help:"Get customer by email"`
	GetByExtID GetByExtIDCmd `cmd:"" help:"Get customer by external ID"`
	Search     SearchCmd     `cmd:"" help:"Search customers"`
}

// Helper functions for customer commands

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
