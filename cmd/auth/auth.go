package auth

import (
	"github.com/jeremybumsted/plain-cli/internal/config"
	"github.com/jeremybumsted/plain-cli/internal/output"
)

// AuthCmd represents the auth command group
type AuthCmd struct {
	Login  LoginCmd  `cmd:"" help:"Authenticate with Plain"`
	Status StatusCmd `cmd:"" help:"Check authentication status"`
	Logout LogoutCmd `cmd:"" help:"Log out and clear credentials"`
}

// getConfig loads the configuration from the specified path
func getConfig(configPath string) (*config.Config, error) {
	return config.Load(configPath)
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
