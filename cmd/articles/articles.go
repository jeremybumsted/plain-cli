package articles

import (
	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/output"
	"github.com/jeremybumsted/plain-cli/internal/plain"
)

// ArticlesCmd represents the articles command group
type ArticlesCmd struct {
	List ListCmd `cmd:"" help:"List help center articles"`
	Get  GetCmd  `cmd:"" help:"Get article by ID"`
}

// Context holds global flags passed from the root command
type Context struct {
	ConfigPath string
	JSON       bool
	Quiet      bool
}

// Helper functions for article commands

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
